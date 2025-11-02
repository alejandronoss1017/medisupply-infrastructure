package blockchain

import (
	"context"
	"contracts/internal/adapter/blockchain/binding"
	"contracts/internal/core/domain"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EthereumReader is an adapter for reading blockchain state
type EthereumReader struct {
	client          *ethclient.Client
	contract        *binding.SLAEnforcer
	contractAddress common.Address
}

// NewEthereumReader creates a new EthereumReader instance
func NewEthereumReader(
	rpcURL string,
	contractAddress string,
) (*EthereumReader, error) {
	// Validate contract address
	if !common.IsHexAddress(contractAddress) {
		return nil, fmt.Errorf("invalid contract address: %s", contractAddress)
	}

	// Connect to Ethereum node
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	addr := common.HexToAddress(contractAddress)

	// Create a contract instance
	contract, err := binding.NewSLAEnforcer(addr, client)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to instantiate contract: %w", err)
	}

	//logger.Info("Ethereum reader initialized",
	//	zap.String("rpcURL", rpcURL),
	//	zap.String("contractAddress", addr.Hex()),
	//)

	return &EthereumReader{
		client:          client,
		contract:        contract,
		contractAddress: addr,
	}, nil
}

// Close closes the Ethereum client connection
func (er *EthereumReader) Close() {
	if er.client != nil {
		er.client.Close()
	}
}

// GetContract retrieves a contract by ID with all its SLAs
func (er *EthereumReader) GetContract(ctx context.Context, contractID string) (*domain.Contract, error) {
	//er.logger.Debug("Getting contract", zap.String("contractID", contractID))

	result, err := er.contract.GetContract(&bind.CallOpts{Context: ctx}, contractID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	// Convert SLAs
	slas := make([]*domain.SLA, len(result.Slas))
	for i, sla := range result.Slas {
		slas[i] = &domain.SLA{
			ID:          sla.Id,
			Name:        sla.Name,
			Description: sla.Description,
			Target:      sla.Target,
			Comparator:  domain.Comparator(sla.Comparator),
			Status:      domain.SLAStatus(sla.Status),
		}
	}

	contract := &domain.Contract{
		ID:         result.Id,
		Path:       result.Path,
		CustomerID: result.CustomerId,
		SLAs:       slas,
	}

	//er.logger.Debug("Contract retrieved successfully",
	//	zap.String("contractID", contract.ID),
	//	zap.Int("slaCount", len(slas)),
	//)

	return contract, nil
}

// GetContractCount returns the total number of contracts
func (er *EthereumReader) GetContractCount(ctx context.Context) (*big.Int, error) {
	//er.logger.Debug("Getting contract count")

	count, err := er.contract.GetContractCount(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("failed to get contract count: %w", err)
	}

	//er.logger.Debug("Contract count retrieved", zap.String("count", count.String()))

	return count, nil
}

// GetContractByIndex retrieves a contract by its index
func (er *EthereumReader) GetContractByIndex(ctx context.Context, index uint64) (*domain.Contract, error) {
	//er.logger.Debug("Getting contract by index", zap.Uint64("index", index))

	result, err := er.contract.Contracts(&bind.CallOpts{Context: ctx}, big.NewInt(int64(index)))
	if err != nil {
		return nil, fmt.Errorf("failed to get contract by index: %w", err)
	}

	// Get full contract details with SLAs
	return er.GetContract(ctx, result.Id)
}

// GetSLA retrieves a specific SLA from a contract
func (er *EthereumReader) GetSLA(ctx context.Context, contractID string, slaIndex uint64) (*domain.SLA, error) {
	//er.logger.Debug("Getting SLA",
	//	zap.String("contractID", contractID),
	//	zap.Uint64("slaIndex", slaIndex),
	//)

	result, err := er.contract.GetSLA(&bind.CallOpts{Context: ctx}, contractID, big.NewInt(int64(slaIndex)))
	if err != nil {
		return nil, fmt.Errorf("failed to get SLA: %w", err)
	}

	sla := &domain.SLA{
		ID:          result.Id,
		Name:        result.Name,
		Description: result.Description,
		Target:      result.Target,
		Comparator:  domain.Comparator(result.Comparator),
		Status:      domain.SLAStatus(result.Status),
	}

	//er.logger.Debug("SLA retrieved successfully",
	//	zap.String("slaID", sla.ID),
	//	zap.String("status", sla.Status.String()),
	//)

	return sla, nil
}

// GetSLAs retrieves all SLAs for a contract
func (er *EthereumReader) GetSLAs(ctx context.Context, contractID string) ([]*domain.SLA, error) {
	//er.logger.Debug("Getting all SLAs", zap.String("contractID", contractID))

	results, err := er.contract.GetSLAs(&bind.CallOpts{Context: ctx}, contractID)
	if err != nil {
		return nil, fmt.Errorf("failed to get SLAs: %w", err)
	}

	slas := make([]*domain.SLA, len(results))
	for i, result := range results {
		slas[i] = &domain.SLA{
			ID:          result.Id,
			Name:        result.Name,
			Description: result.Description,
			Target:      result.Target,
			Comparator:  domain.Comparator(result.Comparator),
			Status:      domain.SLAStatus(result.Status),
		}
	}

	//er.logger.Debug("SLAs retrieved successfully",
	//	zap.String("contractID", contractID),
	//	zap.Int("count", len(slas)),
	//)

	return slas, nil
}

// GetBlockNumber returns the current block number
func (er *EthereumReader) GetBlockNumber(ctx context.Context) (uint64, error) {
	//er.logger.Debug("Getting current block number")

	blockNumber, err := er.client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get block number: %w", err)
	}

	//er.logger.Debug("Block number retrieved", zap.Uint64("blockNumber", blockNumber))

	return blockNumber, nil
}

// GetChainID returns the chain ID
func (er *EthereumReader) GetChainID(ctx context.Context) (*big.Int, error) {
	//er.logger.Debug("Getting chain ID")

	chainID, err := er.client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	//er.logger.Debug("Chain ID retrieved", zap.String("chainID", chainID.String()))

	return chainID, nil
}

// cleanString removes null bytes and trims whitespace from strings
func cleanString(s string) string {
	return strings.TrimRight(strings.TrimSpace(s), "\x00")
}
