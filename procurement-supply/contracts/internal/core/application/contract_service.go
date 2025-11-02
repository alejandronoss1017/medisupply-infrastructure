package application

import (
	"context"
	"contracts/internal/adapter/ethereum"
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
	repo                driven.Repository[string, domain.Contract]
	smartContractClient *ethereum.SmartContractClient
	logger              *logger.Logger
	idSeq               int64
}

// Ensure ContractService implements the driver.ContractService interface
var _ driver.ContractService = (*ContractService)(nil)

// NewContractService creates a new contract service
func NewContractService(repo driven.Repository[string, domain.Contract], smartContractClient *ethereum.SmartContractClient) *ContractService {
	return &ContractService{
		repo:                repo,
		smartContractClient: smartContractClient,
		logger:              logger.New("CONTRACT-SERVICE"),
		idSeq:               time.Now().UnixNano(),
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
	if s.repo.Exists(contract.ID) {
		return nil, fmt.Errorf("contract with id %s already exists", contract.ID)
	}

	// Save to repository
	if err := s.repo.Create(contract); err != nil {
		return nil, fmt.Errorf("failed to create contract in repository: %w", err)
	}

	// Register on blockchain
	s.logger.Info("Registering contract %s on blockchain", contract.ID)
	tx, err := s.smartContractClient.SendContractTransaction(ctx, "createContract", nil, contract)
	if err != nil {
		// Rollback repository changes
		_ = s.repo.Delete(contract.ID)
		return nil, fmt.Errorf("failed to write contract to blockchain: %w", err)
	}

	// Wait for transaction to be mined
	receipt, err := s.smartContractClient.WaitTransaction(ctx, tx)
	if err != nil {
		// Rollback repository changes
		_ = s.repo.Delete(contract.ID)
		return nil, fmt.Errorf("transaction failed on blockchain: %w", err)
	}

	s.logger.Info("Contract %s successfully registered on blockchain (tx: %s)", contract.ID, receipt.TxHash.Hex())

	return &contract, nil
}

// RetrieveContract retrieves a contract by ID
func (s *ContractService) RetrieveContract(id string) (*domain.Contract, error) {
	contract, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("contract not found: %w", err)
	}
	return contract, nil
}

// RetrieveContracts retrieves all contracts
func (s *ContractService) RetrieveContracts() ([]domain.Contract, error) {
	contracts, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve contracts: %w", err)
	}
	return contracts, nil
}

// UpdateContract updates an existing contract
func (s *ContractService) UpdateContract(contract domain.Contract) (*domain.Contract, error) {
	// Verify contract exists
	if !s.repo.Exists(contract.ID) {
		return nil, fmt.Errorf("contract with id %s not found", contract.ID)
	}

	// Update in repository
	if err := s.repo.Update(contract); err != nil {
		return nil, fmt.Errorf("failed to update contract: %w", err)
	}

	s.logger.Info("Contract %s updated successfully", contract.ID)
	return &contract, nil
}

// DeleteContract removes a contract
func (s *ContractService) DeleteContract(id string) error {
	// Verify contract exists
	if !s.repo.Exists(id) {
		return fmt.Errorf("contract with id %s not found", id)
	}

	// Delete from repository
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete contract: %w", err)
	}

	s.logger.Info("Contract %s deleted successfully", id)
	return nil
}
