package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"novelties/internal/core/port/driven"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	"novelties/internal/adapter/blockchain/binding"
)

// EthereumWriter is an adapter for writing to the blockchain
type EthereumWriter struct {
	client          *ethclient.Client
	contract        *binding.SLAEnforcer
	contractAddress common.Address
	privateKey      *ecdsa.PrivateKey
	publicAddress   common.Address
	chainID         *big.Int
	logger          *zap.SugaredLogger
}

// NewEthereumWriter creates a new EthereumWriter instance
func NewEthereumWriter(
	rpcURL string,
	contractAddress string,
	privateKeyHex string,
	logger *zap.SugaredLogger,
) (*EthereumWriter, error) {
	// Validate contract address
	if !common.IsHexAddress(contractAddress) {
		return nil, fmt.Errorf("invalid contract address: %s", contractAddress)
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	// Derive public address from private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	publicAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Connect to Ethereum node
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	// Get chain ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	addr := common.HexToAddress(contractAddress)

	// Create contract instance
	contract, err := binding.NewSLAEnforcer(addr, client)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to instantiate contract: %w", err)
	}

	logger.Infow("Ethereum writer initialized",
		"rpc_url", rpcURL,
		"contract_address", addr.Hex(),
		"public_address", publicAddress.Hex(),
		"chain_id", chainID.String(),
	)

	return &EthereumWriter{
		client:          client,
		contract:        contract,
		contractAddress: addr,
		privateKey:      privateKey,
		publicAddress:   publicAddress,
		chainID:         chainID,
		logger:          logger,
	}, nil
}

// Close closes the Ethereum client connection
func (ew *EthereumWriter) Close() {
	if ew.client != nil {
		ew.client.Close()
	}
}

// GetAddress returns the address of the account used for transactions
func (ew *EthereumWriter) GetAddress() string {
	return ew.publicAddress.Hex()
}

// GetBalance returns the balance of the account
func (ew *EthereumWriter) GetBalance(ctx context.Context) (*big.Int, error) {
	balance, err := ew.client.BalanceAt(ctx, ew.publicAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	ew.logger.Debugw("Balance retrieved",
		"address", ew.publicAddress.Hex(),
		"balance", balance.String(),
	)

	return balance, nil
}

// createTransactor creates a transactor with the current nonce and gas settings
func (ew *EthereumWriter) createTransactor(ctx context.Context) (*bind.TransactOpts, error) {
	nonce, err := ew.client.PendingNonceAt(ctx, ew.publicAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := ew.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(ew.privateKey, ew.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice
	auth.Context = ctx

	return auth, nil
}

// AddContract adds a new contract to the blockchain
func (ew *EthereumWriter) AddContract(ctx context.Context, contractID, path, customerID string) (*driven.TransactionReceipt, error) {
	ew.logger.Infow("Adding contract to blockchain",
		"contract_id", contractID,
		"customer_id", customerID,
		"path", path,
	)

	auth, err := ew.createTransactor(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := ew.contract.AddContract(auth, contractID, path, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	ew.logger.Infow("Transaction sent to blockchain",
		"tx_hash", tx.Hash().Hex(),
	)

	// Wait for transaction receipt
	receipt, err := bind.WaitMined(ctx, ew.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}

	txReceipt := &driven.TransactionReceipt{
		TxHash:      receipt.TxHash.Hex(),
		BlockNumber: receipt.BlockNumber.Uint64(),
		Status:      receipt.Status,
		GasUsed:     receipt.GasUsed,
	}

	if receipt.Status == 1 {
		ew.logger.Infow("Contract added successfully to blockchain",
			"contract_id", contractID,
			"tx_hash", txReceipt.TxHash,
			"block_number", txReceipt.BlockNumber,
			"gas_used", txReceipt.GasUsed,
		)
	} else {
		ew.logger.Errorw("Contract addition failed on blockchain",
			"contract_id", contractID,
			"tx_hash", txReceipt.TxHash,
			"status", receipt.Status,
		)
	}

	return txReceipt, nil
}

// AddSLA adds a new SLA to an existing contract
func (ew *EthereumWriter) AddSLA(ctx context.Context, contractID, slaID, name, description string, target *big.Int, comparator uint8) (*driven.TransactionReceipt, error) {
	ew.logger.Infow("Adding SLA to contract",
		"contract_id", contractID,
		"sla_id", slaID,
		"name", name,
		"target", target.String(),
		"comparator", comparator,
	)

	auth, err := ew.createTransactor(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := ew.contract.AddSLA(auth, contractID, slaID, name, description, target, comparator)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	ew.logger.Infow("Transaction sent to blockchain",
		"tx_hash", tx.Hash().Hex(),
	)

	// Wait for transaction receipt
	receipt, err := bind.WaitMined(ctx, ew.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}

	txReceipt := &driven.TransactionReceipt{
		TxHash:      receipt.TxHash.Hex(),
		BlockNumber: receipt.BlockNumber.Uint64(),
		Status:      receipt.Status,
		GasUsed:     receipt.GasUsed,
	}

	if receipt.Status == 1 {
		ew.logger.Infow("SLA added successfully to blockchain",
			"sla_id", slaID,
			"tx_hash", txReceipt.TxHash,
			"block_number", txReceipt.BlockNumber,
			"gas_used", txReceipt.GasUsed,
		)
	} else {
		ew.logger.Errorw("SLA addition failed on blockchain",
			"sla_id", slaID,
			"tx_hash", txReceipt.TxHash,
			"status", receipt.Status,
		)
	}

	return txReceipt, nil
}

// SetSLAStatus updates the status of an SLA
func (ew *EthereumWriter) SetSLAStatus(ctx context.Context, contractID string, slaIndex uint64, status uint8) (*driven.TransactionReceipt, error) {
	ew.logger.Infow("Setting SLA status",
		"contract_id", contractID,
		"sla_index", slaIndex,
		"status", status,
	)

	auth, err := ew.createTransactor(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := ew.contract.SetSLAStatus(auth, contractID, big.NewInt(int64(slaIndex)), status)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	ew.logger.Infow("Transaction sent to blockchain",
		"tx_hash", tx.Hash().Hex(),
	)

	// Wait for transaction receipt
	receipt, err := bind.WaitMined(ctx, ew.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}

	txReceipt := &driven.TransactionReceipt{
		TxHash:      receipt.TxHash.Hex(),
		BlockNumber: receipt.BlockNumber.Uint64(),
		Status:      receipt.Status,
		GasUsed:     receipt.GasUsed,
	}

	if receipt.Status == 1 {
		ew.logger.Infow("SLA status updated successfully on blockchain",
			"contract_id", contractID,
			"sla_index", slaIndex,
			"tx_hash", txReceipt.TxHash,
			"block_number", txReceipt.BlockNumber,
			"gas_used", txReceipt.GasUsed,
		)
	} else {
		ew.logger.Errorw("SLA status update failed on blockchain",
			"contract_id", contractID,
			"sla_index", slaIndex,
			"tx_hash", txReceipt.TxHash,
			"status", receipt.Status,
		)
	}

	return txReceipt, nil
}

// CheckSLA checks an SLA against an actual value
func (ew *EthereumWriter) CheckSLA(ctx context.Context, contractID string, slaID string, actualValue *big.Int) (*driven.TransactionReceipt, error) {
	ew.logger.Infow("Checking SLA",
		"contract_id", contractID,
		"sla_id", slaID,
		"actual_value", actualValue.String(),
	)

	auth, err := ew.createTransactor(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := ew.contract.CheckSLA(auth, contractID, slaID, actualValue)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	ew.logger.Infow("Transaction sent to blockchain",
		"tx_hash", tx.Hash().Hex(),
	)

	// Wait for transaction receipt
	receipt, err := bind.WaitMined(ctx, ew.client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction: %w", err)
	}

	txReceipt := &driven.TransactionReceipt{
		TxHash:      receipt.TxHash.Hex(),
		BlockNumber: receipt.BlockNumber.Uint64(),
		Status:      receipt.Status,
		GasUsed:     receipt.GasUsed,
	}

	if receipt.Status == 1 {
		ew.logger.Infow("SLA checked successfully on blockchain",
			"contract_id", contractID,
			"sla_id", slaID,
			"tx_hash", txReceipt.TxHash,
			"block_number", txReceipt.BlockNumber,
			"gas_used", txReceipt.GasUsed,
		)
	} else {
		ew.logger.Errorw("SLA check failed on blockchain",
			"contract_id", contractID,
			"sla_id", slaID,
			"tx_hash", txReceipt.TxHash,
			"status", receipt.Status,
		)
	}

	return txReceipt, nil
}
