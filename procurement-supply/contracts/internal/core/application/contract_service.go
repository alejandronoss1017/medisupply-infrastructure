package application

import (
	"context"
	"contracts/internal/core/domain"
	"contracts/internal/core/port/driven"
	"contracts/internal/core/port/driver"
	"contracts/pkg/logger"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// ContractService handles business logic for contract management
type ContractService struct {
	repository       driven.Repository[string, domain.Contract]
	blockchainWriter driven.BlockchainWriter
	blockchainReader driven.BlockchainReader
	logger           *zap.SugaredLogger
	idSeq            int64
}

// Ensure ContractService implements the driver.ContractService interface
var _ driver.ContractService = (*ContractService)(nil)

// NewContractService creates a new contract service
func NewContractService(repository driven.Repository[string, domain.Contract], writer driven.BlockchainWriter, reader driven.BlockchainReader) *ContractService {
	return &ContractService{
		repository:       repository,
		blockchainWriter: writer,
		blockchainReader: reader,
		logger:           logger.New("CONTRACT-SERVICE"),
		idSeq:            time.Now().UnixNano(),
	}
}

// CreateContract creates a new contract and registers it on the blockchain
func (s *ContractService) CreateContract(ctx context.Context, contract domain.Contract) (*domain.Contract, error) {
	// Generate ID if empty
	if contract.ID == "" {
		s.idSeq++
		contract.ID = strconv.FormatInt(s.idSeq, 10)
	}

	// Check if the contract already exists
	if s.repository.Exists(contract.ID) {
		return nil, fmt.Errorf("contract with id %s already exists", contract.ID)
	}

	// Save to repository
	if err := s.repository.Create(contract); err != nil {
		return nil, fmt.Errorf("failed to create contract in repository: %w", err)
	}

	// Register on blockchain
	s.logger.Infow("Registering contract on blockchain",
		"contract_id", contract.ID,
		"customer_id", contract.CustomerID,
		"path", contract.Path,
	)
	receipt, err := s.blockchainWriter.AddContract(ctx, contract.ID, contract.Path, contract.CustomerID)
	if err != nil {
		// Rollback repository changes
		_ = s.repository.Delete(contract.ID)
		return nil, fmt.Errorf("failed to write contract to blockchain: %w", err)
	}

	// Check transaction status
	if receipt.Status != 1 {
		// Rollback repository changes
		_ = s.repository.Delete(contract.ID)
		return nil, fmt.Errorf("blockchain transaction failed (status: %d, tx: %s)", receipt.Status, receipt.TxHash)
	}

	s.logger.Infow("Contract successfully registered on blockchain",
		"contract_id", contract.ID,
		"tx_hash", receipt.TxHash,
		"block_number", receipt.BlockNumber,
		"gas_used", receipt.GasUsed,
	)

	return &contract, nil
}

// RetrieveContract retrieves a contract by ID
func (s *ContractService) RetrieveContract(ctx context.Context, id string) (*domain.Contract, error) {

	contract, err := s.blockchainReader.GetContract(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("contract not found: %w", err)
	}

	return contract, nil
}

// RetrieveContracts retrieves all contracts
func (s *ContractService) RetrieveContracts(ctx context.Context) ([]*domain.Contract, error) {

	contracts, err := s.blockchainReader.GetContracts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve contracts: %w", err)
	}

	return contracts, nil
}

// UpdateContract updates an existing contract
func (s *ContractService) UpdateContract(contract domain.Contract) (*domain.Contract, error) {
	// Verify contract exists
	if !s.repository.Exists(contract.ID) {
		return nil, fmt.Errorf("contract with id %s not found", contract.ID)
	}

	// Update in repository
	if err := s.repository.Update(contract); err != nil {
		return nil, fmt.Errorf("failed to update contract: %w", err)
	}

	s.logger.Infow("Contract updated successfully",
		"contract_id", contract.ID,
	)
	return &contract, nil
}

// DeleteContract removes a contract
func (s *ContractService) DeleteContract(id string) error {
	// Verify contract exists
	if !s.repository.Exists(id) {
		return fmt.Errorf("contract with id %s not found", id)
	}

	// Delete it from the repository
	if err := s.repository.Delete(id); err != nil {
		return fmt.Errorf("failed to delete contract: %w", err)
	}

	s.logger.Infow("Contract deleted successfully",
		"contract_id", id,
	)
	return nil
}

func (s *ContractService) CreateSLA(ctx context.Context, id string, sla domain.SLA) (*domain.SLA, error) {
	// Validate inputs
	if id == "" {
		return nil, fmt.Errorf("contract id cannot be empty")
	}
	if sla.Name == "" {
		return nil, fmt.Errorf("sla name cannot be empty")
	}
	if sla.Target == nil {
		return nil, fmt.Errorf("sla target cannot be nil")
	}

	// Generate SLA ID if empty
	if sla.ID == "" {
		s.idSeq++
		sla.ID = strconv.FormatInt(s.idSeq, 10)
	}

	s.logger.Infow("Appending SLA to contract on blockchain",
		"contract_id", id,
		"sla_id", sla.ID,
		"name", sla.Name,
	)

	receipt, err := s.blockchainWriter.AddSLA(ctx, id, sla.ID, sla.Name, sla.Description, sla.Target, uint8(sla.Comparator))
	if err != nil {
		return nil, fmt.Errorf("failed to add SLA to blockchain: %w", err)
	}
	if receipt.Status != 1 {
		return nil, fmt.Errorf("blockchain transaction failed (status: %d, tx: %s)", receipt.Status, receipt.TxHash)
	}

	s.logger.Infow("SLA successfully appended to contract on blockchain",
		"contract_id", id,
		"sla_id", sla.ID,
		"tx_hash", receipt.TxHash,
		"block_number", receipt.BlockNumber,
		"gas_used", receipt.GasUsed,
	)

	return &sla, nil
}

func (s *ContractService) RetrieveSLAs(ctx context.Context, id string) ([]domain.SLA, error) {
	// Validate input
	if id == "" {
		return nil, fmt.Errorf("contract id cannot be empty")
	}

	// Read SLAs from the blockchain
	slaPtrs, err := s.blockchainReader.GetSLAs(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve slas: %w", err)
	}

	// Convert []*domain.SLA to []domain.SLA
	slas := make([]domain.SLA, len(slaPtrs))
	for i, p := range slaPtrs {
		if p == nil {
			continue
		}
		slas[i] = *p
	}

	return slas, nil
}
