package blockchain

import (
	"context"
	"contracts/internal/core/port/driven"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"contracts/internal/adapter/blockchain/binding"
)

// EthereumWriter is an adapter for writing to the blockchain
type EthereumWriter struct {
	client          *ethclient.Client
	contract        *binding.SLAEnforcer
	contractAddress common.Address
	privateKey      *ecdsa.PrivateKey
	publicAddress   common.Address
	chainID         *big.Int
	//logger          *zap.Logger
}

// NewEthereumWriter creates a new EthereumWriter instance
func NewEthereumWriter(
	rpcURL string,
	contractAddress string,
	privateKeyHex string,
	// logger *zap.Logger,
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

	//logger.Info("Ethereum writer initialized",
	//	zap.String("rpcURL", rpcURL),
	//	zap.String("contractAddress", addr.Hex()),
	//	zap.String("publicAddress", publicAddress.Hex()),
	//	zap.String("chainID", chainID.String()),
	//)

	return &EthereumWriter{
		client:          client,
		contract:        contract,
		contractAddress: addr,
		privateKey:      privateKey,
		publicAddress:   publicAddress,
		chainID:         chainID,
		//logger:          logger,
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

	//ew.logger.Debug("Balance retrieved",
	//	zap.String("address", ew.publicAddress.Hex()),
	//	zap.String("balance", balance.String()),
	//)

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
	//ew.logger.Info("Adding contract",
	//	zap.String("contractID", contractID),
	//	zap.String("customerID", customerID),
	//)

	auth, err := ew.createTransactor(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := ew.contract.AddContract(auth, contractID, path, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	//ew.logger.Info("Transaction sent",
	//	zap.String("txHash", tx.Hash().Hex()),
	//)

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

	//if receipt.Status == 1 {
	//	ew.logger.Info("Contract added successfully",
	//		zap.String("contractID", contractID),
	//		zap.String("txHash", txReceipt.TxHash),
	//		zap.Uint64("blockNumber", txReceipt.BlockNumber),
	//	)
	//} else {
	//	ew.logger.Error("Contract addition failed",
	//		zap.String("contractID", contractID),
	//		zap.String("txHash", txReceipt.TxHash),
	//	)
	//}

	return txReceipt, nil
}

// AddSLA adds a new SLA to an existing contract
func (ew *EthereumWriter) AddSLA(ctx context.Context, contractID, slaID, name, description string, target *big.Int, comparator uint8) (*driven.TransactionReceipt, error) {
	//ew.logger.Info("Adding SLA",
	//	zap.String("contractID", contractID),
	//	zap.String("slaID", slaID),
	//	zap.String("name", name),
	//)

	auth, err := ew.createTransactor(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := ew.contract.AddSLA(auth, contractID, slaID, name, description, target, comparator)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	//ew.logger.Info("Transaction sent",
	//	zap.String("txHash", tx.Hash().Hex()),
	//)

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

	//if receipt.Status == 1 {
	//	ew.logger.Info("SLA added successfully",
	//		zap.String("slaID", slaID),
	//		zap.String("txHash", txReceipt.TxHash),
	//		zap.Uint64("blockNumber", txReceipt.BlockNumber),
	//	)
	//} else {
	//	ew.logger.Error("SLA addition failed",
	//		zap.String("slaID", slaID),
	//		zap.String("txHash", txReceipt.TxHash),
	//	)
	//}

	return txReceipt, nil
}

// SetSLAStatus updates the status of an SLA
func (ew *EthereumWriter) SetSLAStatus(ctx context.Context, contractID string, slaIndex uint64, status uint8) (*driven.TransactionReceipt, error) {
	//ew.logger.Info("Setting SLA status",
	//	zap.String("contractID", contractID),
	//	zap.Uint64("slaIndex", slaIndex),
	//	zap.Uint8("status", status),
	//)

	auth, err := ew.createTransactor(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := ew.contract.SetSLAStatus(auth, contractID, big.NewInt(int64(slaIndex)), status)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	//ew.logger.Info("Transaction sent",
	//	zap.String("txHash", tx.Hash().Hex()),
	//)

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

	//if receipt.Status == 1 {
	//	ew.logger.Info("SLA status updated successfully",
	//		zap.String("contractID", contractID),
	//		zap.Uint64("slaIndex", slaIndex),
	//		zap.String("txHash", txReceipt.TxHash),
	//	)
	//} else {
	//	ew.logger.Error("SLA status update failed",
	//		zap.String("contractID", contractID),
	//		zap.Uint64("slaIndex", slaIndex),
	//		zap.String("txHash", txReceipt.TxHash),
	//	)
	//}

	return txReceipt, nil
}

// CheckSLA checks an SLA against an actual value
func (ew *EthereumWriter) CheckSLA(ctx context.Context, contractID string, slaIndex uint64, actualValue *big.Int) (*driven.TransactionReceipt, error) {
	//ew.logger.Info("Checking SLA",
	//	zap.String("contractID", contractID),
	//	zap.Uint64("slaIndex", slaIndex),
	//	zap.String("actualValue", actualValue.String()),
	//)

	auth, err := ew.createTransactor(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := ew.contract.CheckSLA(auth, contractID, big.NewInt(int64(slaIndex)), actualValue)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	//ew.logger.Info("Transaction sent",
	//	zap.String("txHash", tx.Hash().Hex()),
	//)

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

	//if receipt.Status == 1 {
	//	ew.logger.Info("SLA checked successfully",
	//		zap.String("contractID", contractID),
	//		zap.Uint64("slaIndex", slaIndex),
	//		zap.String("txHash", txReceipt.TxHash),
	//	)
	//} else {
	//	ew.logger.Error("SLA check failed",
	//		zap.String("contractID", contractID),
	//		zap.Uint64("slaIndex", slaIndex),
	//		zap.String("txHash", txReceipt.TxHash),
	//	)
	//}

	return txReceipt, nil
}
