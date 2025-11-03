package driver

import (
	"context"
	"contracts/internal/core/domain"
)

// ContractService defines the interface for contract business logic
type ContractService interface {
	// CreateContract creates a new contract and registers it on the blockchain
	CreateContract(ctx context.Context, contract domain.Contract) (*domain.Contract, error)

	// RetrieveContract retrieves a contract by ID
	RetrieveContract(ctx context.Context, id string) (*domain.Contract, error)

	// RetrieveContracts retrieves all contracts
	RetrieveContracts(ctx context.Context) ([]*domain.Contract, error)

	// UpdateContract updates an existing contract
	UpdateContract(contract domain.Contract) (*domain.Contract, error)

	// DeleteContract removes a contract
	DeleteContract(id string) error

	// CreateSLA append an SLA to a specific contract
	CreateSLA(ctx context.Context, id string, sla domain.SLA) (*domain.SLA, error)

	// RetrieveSLAs retrieves all SLAs related to a specific contract
	RetrieveSLAs(ctx context.Context, id string) ([]domain.SLA, error)
}

// SLAService defines the interface for SLA business logic
type SLAService interface {
	// CreateSLA creates a new SLA
	CreateSLA(sla domain.SLA) (*domain.SLA, error)

	// RetrieveSLA retrieves an SLA by ID
	RetrieveSLA(id string) (*domain.SLA, error)

	// RetrieveSLAs retrieves all SLAs
	RetrieveSLAs() ([]domain.SLA, error)

	// UpdateSLA updates an existing SLA
	UpdateSLA(sla domain.SLA) (*domain.SLA, error)

	// DeleteSLA removes an SLA
	DeleteSLA(id string) error
}

// CustomerService defines the interface for customer business logic
type CustomerService interface {
	// CreateCustomer creates a new customer
	CreateCustomer(customer domain.Customer) (*domain.Customer, error)

	// RetrieveCustomer retrieves a customer by ID
	RetrieveCustomer(id string) (*domain.Customer, error)

	// RetrieveCustomers retrieves all customers
	RetrieveCustomers() ([]domain.Customer, error)

	// UpdateCustomer updates an existing customer
	UpdateCustomer(customer domain.Customer) (*domain.Customer, error)

	// DeleteCustomer removes a customer
	DeleteCustomer(id string) error
}
