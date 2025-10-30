package ethereum

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SmartContractClient wraps an Ethereum RPC client together with a target
// contract address and its ABI. It exposes helpers to call read-only methods
// and to send signed state-changing transactions using the configured private key.
// The zero value is not usable; construct one via NewSmartContractClient.
type SmartContractClient struct {
	client        *ethclient.Client
	address       common.Address
	abi           abi.ABI
	privateKey    *ecdsa.PrivateKey
	publicAddress common.Address
}

// NewSmartContractClient creates and returns a SmartContractClient for interacting with a
// specific on-chain contract using the provided RPC endpoint, contract address,
// private key, and ABI file path.
//
// Parameters:
//
//	rcpURL - RPC endpoint URL (HTTP(S) or WS) of the Ethereum node
//	Address - Hex-encoded contract address (0x-prefixed accepted)
//	Key - Hex-encoded ECDSA private key
//	Path - Filesystem path to the ABI JSON file
//
// Returns:
//
//	A configured SmartContractClient on success, or a non-nil error if connecting
//	to the node fails, the ABI file cannot be read or parsed, or the private key
//	is invalid.
func NewSmartContractClient(rcpURL, address, key, path string) (*SmartContractClient, error) {
	client, err := ethclient.Dial(rcpURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Ethereum network: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading ABI file from %q: %w", path, err)
	}

	abiContract, err := abi.JSON(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("error loading ABI from file %q: %w", path, err)
	}

	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, fmt.Errorf("error converting private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	publicAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &SmartContractClient{
		client:        client,
		address:       common.HexToAddress(address),
		privateKey:    privateKey,
		abi:           abiContract,
		publicAddress: publicAddress,
	}, nil
}

// InvoqueContract calls a constant (read-only) method on the configured contract
// and decodes the return values into out.
// Supported bindings:
// - Single return value: pass a pointer to the Go type (e.g. *uint64, *bool, *string, *common.Address)
// - Multiple return values: pass a pointer to a struct with fields in order or tagged with `abi:"<index or name>"`
// - Dynamic catch-all: pass a pointer to []interface{} to receive all values
//
// Parameters:
//
//	ctx - Context to use for the call
//	method - The name of the contract method to call
//	out - A pointer to the destination type that matches the ABI outputs
//	args - The arguments to pass to the contract method
//
// Return:
//
//	An error if the method does not exist in the ABI or if the call fails.
//
// Examples:
//
//	var bal *big.Int
//	if err := client.InvoqueContract(ctx, "balanceOf", &bal, addr); err != nil { /* handle */ }
//
//	type Info struct {
//	    Name  string          `abi:"0"` // or use output names if defined in the ABI
//	    Owner common.Address  `abi:"1"`
//	}
//
//	var info Info
//	if err := client.InvoqueContract(ctx, "getInfo", &info); err != nil { /* handle */ }
func (c *SmartContractClient) InvoqueContract(ctx context.Context, method string, out *any, args ...any) error {
	// Validate method exists (gives a better error than Pack for unknown method)
	input, err := c.packMethodInput(method, args...)
	if err != nil {
		return err
	}

	msg := ethereum.CallMsg{
		From: c.publicAddress,
		To:   &c.address,
		Data: input,
	}

	data, err := c.client.CallContract(ctx, msg, nil)
	if err != nil {
		return fmt.Errorf("call contract %s: %w", method, err)
	}

	// Decode into provided template
	if err := c.abi.UnpackIntoInterface(out, method, data); err != nil {
		return fmt.Errorf("unpack outputs for %s: %w", method, err)
	}
	return nil
}

// InvoqueContractWrite sends a state-changing transaction to the configured contract.
// It packs the given method and args using the loaded ABI, estimates gas, builds
// an EIP-1559 transaction (falling back to legacy if base fee is unavailable),
// signs it with the client's private key, and broadcasts it.
//
// Parameters:
//
//	ctx - Context for RPC calls and cancellation
//	method - Contract method name to invoke
//	value - Ether value to send with the transaction (can be nil for 0)
//	args - Method arguments per the ABI, this must be pointers to the expected
//		types (e.g. *big.Int, *common.Address, *string, etc.)
//
// Returns:
//
//	The broadcasted transaction and/or an error if any step fails (pack, gas
//	estimation, fee calculation, signing, or RPC broadcast).
func (c *SmartContractClient) InvoqueContractWrite(ctx context.Context, method string, value *big.Int, args ...any) (*types.Transaction, error) {
	input, err := c.packMethodInput(method, args...)
	if err != nil {
		return nil, err
	}

	// Default value to 0 if nil
	if value == nil {
		value = big.NewInt(0)
	}

	// Chain ID and nonce
	chainID, err := c.client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}
	nonce, err := c.client.PendingNonceAt(ctx, c.publicAddress)
	if err != nil {
		return nil, fmt.Errorf("get pending nonce: %w", err)
	}

	// Estimate gas
	call := ethereum.CallMsg{
		From:  c.publicAddress,
		To:    &c.address,
		Value: value,
		Data:  input,
	}
	gasLimit, err := c.client.EstimateGas(ctx, call)
	if err != nil {
		return nil, fmt.Errorf("estimate gas for %s: %w", method, err)
	}
	// Optional buffer (10%) to reduce underestimation failures
	gasLimit = uint64(float64(gasLimit) * 1.1)

	// Try EIP-1559 fees
	head, herr := c.client.HeaderByNumber(ctx, nil)
	if herr == nil && head != nil && head.BaseFee != nil {
		// EIP-1559 path
		tipCap, err := c.client.SuggestGasTipCap(ctx)
		if err != nil {
			return nil, fmt.Errorf("suggest gas tip cap: %w", err)
		}
		// feeCap = baseFee*2 + tipCap (simple rule of thumb)
		feeCap := new(big.Int).Add(new(big.Int).Mul(head.BaseFee, big.NewInt(2)), tipCap)

		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			To:        &c.address,
			Value:     value,
			Data:      input,
			Gas:       gasLimit,
			GasTipCap: tipCap,
			GasFeeCap: feeCap,
		})

		signer := types.LatestSignerForChainID(chainID)
		signedTx, err := types.SignTx(tx, signer, c.privateKey)
		if err != nil {
			return nil, fmt.Errorf("sign eip-1559 tx: %w", err)
		}
		if err := c.client.SendTransaction(ctx, signedTx); err != nil {
			return nil, fmt.Errorf("send eip-1559 tx: %w", err)
		}
		return signedTx, nil
	}

	// Legacy fallback (pre EIP-1559 or no base fee available)
	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("suggest gas price: %w", err)
	}
	legacy := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &c.address,
		Value:    value,
		Data:     input,
		Gas:      gasLimit,
		GasPrice: gasPrice,
	})
	signer := types.LatestSignerForChainID(chainID)
	signedLegacy, err := types.SignTx(legacy, signer, c.privateKey)
	if err != nil {
		return nil, fmt.Errorf("sign legacy tx: %w", err)
	}
	if err := c.client.SendTransaction(ctx, signedLegacy); err != nil {
		return nil, fmt.Errorf("send legacy tx: %w", err)
	}
	return signedLegacy, nil
}

// GetBalance returns the Ether balance (in wei) of the provided address at the
// latest known block. It wraps eth_getBalance via the underlying RPC client.
// It returns a non-nil error if the RPC call fails.
func (c *SmartContractClient) GetBalance(ctx context.Context, addr common.Address) (*big.Int, error) {
	balance, err := c.client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("get balance: %w", err)
	}
	return balance, nil
}

// WaitTransaction blocks until the given transaction is mined and then returns
// its receipt. It returns an error if waiting fails or if the mined receipt
// indicates failure (Status == 0).
func (c *SmartContractClient) WaitTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {

	receipt, err := bind.WaitMined(ctx, c.client, tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("wait for tx: %w", err)
	}

	if receipt.Status == 0 {
		return nil, fmt.Errorf("tx failed")
	}

	return receipt, nil
}

// packMethodInput validates that the method exists in the client's ABI and
// packs the provided arguments into call data for that method.
// It mirrors the previous inline error messages for consistency.
func (c *SmartContractClient) packMethodInput(method string, args ...any) ([]byte, error) {
	// Ensure the method exists in ABI
	if _, ok := c.abi.Methods[method]; !ok {
		return nil, fmt.Errorf("method %q not found in ABI", method)
	}
	// Pack input data
	input, err := c.abi.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("pack args for %s: %w", method, err)
	}
	return input, nil
}
