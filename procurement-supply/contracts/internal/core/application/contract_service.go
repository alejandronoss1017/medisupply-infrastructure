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
)

// ContractService handles business logic for contract management
type ContractService struct {
	repository       driven.Repository[string, domain.Contract]
	blockchainWriter driven.BlockchainWriter
	logger           *logger.Logger
	idSeq            int64
}

// Ensure ContractService implements the driver.ContractService interface
var _ driver.ContractService = (*ContractService)(nil)

// NewContractService creates a new contract service
func NewContractService(repository driven.Repository[string, domain.Contract], blockchainWriter driven.BlockchainWriter) *ContractService {
	return &ContractService{
		repository:       repository,
		blockchainWriter: blockchainWriter,
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

	// Check if contract already exists
	if s.repository.Exists(contract.ID) {
		return nil, fmt.Errorf("contract with id %s already exists", contract.ID)
	}

	// Save to repository
	if err := s.repository.Create(contract); err != nil {
		return nil, fmt.Errorf("failed to create contract in repository: %w", err)
	}

	// Register on blockchain
	s.logger.Info("Registering contract %s on blockchain", contract.ID)
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

	s.logger.Info("Contract %s successfully registered on blockchain (tx: %s, block: %d)",
		contract.ID, receipt.TxHash, receipt.BlockNumber)

	return &contract, nil
}

// RetrieveContract retrieves a contract by ID
func (s *ContractService) RetrieveContract(id string) (*domain.Contract, error) {
	contract, err := s.repository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("contract not found: %w", err)
	}
	return contract, nil
}

// RetrieveContracts retrieves all contracts
func (s *ContractService) RetrieveContracts() ([]domain.Contract, error) {
	contracts, err := s.repository.FindAll()
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

	s.logger.Info("Contract %s updated successfully", contract.ID)
	return &contract, nil
}

// DeleteContract removes a contract
func (s *ContractService) DeleteContract(id string) error {
	// Verify contract exists
	if !s.repository.Exists(id) {
		return fmt.Errorf("contract with id %s not found", id)
	}

	// Delete from repository
	if err := s.repository.Delete(id); err != nil {
		return fmt.Errorf("failed to delete contract: %w", err)
	}

	s.logger.Info("Contract %s deleted successfully", id)
	return nil
}
