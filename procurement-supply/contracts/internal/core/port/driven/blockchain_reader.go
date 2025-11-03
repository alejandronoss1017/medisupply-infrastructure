package driven

import (
	"context"
	"contracts/internal/core/domain"
	"math/big"
)

// BlockchainReader defines the interface for reading blockchain state
type BlockchainReader interface {
	// GetContract retrieves a contract by ID with all its SLAs
	GetContract(ctx context.Context, contractID string) (*domain.Contract, error)

	// GetContracts retrieves all contracts stored in the blockchain
	GetContracts(ctx context.Context) ([]*domain.Contract, error)

	// GetContractCount returns the total number of contracts
	GetContractCount(ctx context.Context) (*big.Int, error)

	// GetContractByIndex retrieves a contract by its index
	GetContractByIndex(ctx context.Context, index uint64) (*domain.Contract, error)

	// GetSLA retrieves a specific SLA from a contract
	GetSLA(ctx context.Context, contractID string, slaIndex uint64) (*domain.SLA, error)

	// GetSLAs retrieves all SLAs for a contract
	GetSLAs(ctx context.Context, contractID string) ([]*domain.SLA, error)

	// GetBlockNumber returns the current block number
	GetBlockNumber(ctx context.Context) (uint64, error)

	// GetChainID returns the chain ID
	GetChainID(ctx context.Context) (*big.Int, error)
}
