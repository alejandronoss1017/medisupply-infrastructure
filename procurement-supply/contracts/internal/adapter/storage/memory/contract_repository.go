package memory

import (
	"contracts/internal/core/domain"
	"contracts/internal/core/port/driven"
	"fmt"
	"sync"
)

// ContractRepository is an in-memory implementation of the Repository[string, domain.Contract] port
type ContractRepository struct {
	mu        sync.RWMutex
	contracts map[string]domain.Contract
}

// Ensure ContractRepository implements the Repository interface
var _ driven.Repository[string, domain.Contract] = (*ContractRepository)(nil)

// NewContractRepository creates a new in-memory contract repository
func NewContractRepository() *ContractRepository {
	return &ContractRepository{
		contracts: make(map[string]domain.Contract),
	}
}

// Create adds a new contract to the repository
func (r *ContractRepository) Create(contract domain.Contract) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.contracts[contract.ID]; exists {
		return fmt.Errorf("contract with id %s already exists", contract.ID)
	}

	r.contracts[contract.ID] = contract
	return nil
}

// FindByID retrieves a contract by its ID
func (r *ContractRepository) FindByID(id string) (*domain.Contract, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	contract, exists := r.contracts[id]
	if !exists {
		return nil, fmt.Errorf("contract with id %s not found", id)
	}

	return &contract, nil
}

// FindAll retrieves all contracts
func (r *ContractRepository) FindAll() ([]domain.Contract, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	contracts := make([]domain.Contract, 0, len(r.contracts))
	for _, contract := range r.contracts {
		contracts = append(contracts, contract)
	}

	return contracts, nil
}

// Update modifies an existing contract
func (r *ContractRepository) Update(contract domain.Contract) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.contracts[contract.ID]; !exists {
		return fmt.Errorf("contract with id %s not found", contract.ID)
	}

	r.contracts[contract.ID] = contract
	return nil
}

// Delete removes a contract by its ID
func (r *ContractRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.contracts[id]; !exists {
		return fmt.Errorf("contract with id %s not found", id)
	}

	delete(r.contracts, id)
	return nil
}

// Exists checks if a contract with the given ID exists
func (r *ContractRepository) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.contracts[id]
	return exists
}
