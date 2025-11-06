package driven

import (
	"context"
	"math/big"
)

// TransactionReceipt represents the result of a blockchain transaction
type TransactionReceipt struct {
	TxHash      string
	BlockNumber uint64
	Status      uint64 // 1 = success, 0 = failure
	GasUsed     uint64
}

// BlockchainWriter defines the interface for writing to the blockchain
type BlockchainWriter interface {
	// AddContract adds a new contract to the blockchain
	AddContract(ctx context.Context, contractID, path, customerID string) (*TransactionReceipt, error)

	// AddSLA adds a new SLA to an existing contract
	AddSLA(ctx context.Context, contractID, slaID, name, description string, target *big.Int, comparator uint8) (*TransactionReceipt, error)

	// SetSLAStatus updates the status of an SLA
	SetSLAStatus(ctx context.Context, contractID string, slaIndex uint64, status uint8) (*TransactionReceipt, error)

	// CheckSLA checks an SLA against an actual value
	CheckSLA(ctx context.Context, contractID string, slaID string, actualValue *big.Int) (*TransactionReceipt, error)

	// GetAddress returns the address of the account used for transactions
	GetAddress() string

	// GetBalance returns the balance of the account
	GetBalance(ctx context.Context) (*big.Int, error)
}
