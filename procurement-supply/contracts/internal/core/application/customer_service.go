package application

import (
	"contracts/internal/core/domain"
	"contracts/internal/core/port/driven"
	"contracts/internal/core/port/driver"
	"contracts/pkg/logger"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// CustomerService handles business logic for customer management
type CustomerService struct {
	repo   driven.Repository[string, domain.Customer]
	logger *zap.SugaredLogger
	idSeq  int64
}

// Ensure CustomerService implements the driver.CustomerService interface
var _ driver.CustomerService = (*CustomerService)(nil)

// NewCustomerService creates a new customer service
func NewCustomerService(repo driven.Repository[string, domain.Customer]) *CustomerService {
	return &CustomerService{
		repo:   repo,
		logger: logger.New("CUSTOMER-SERVICE"),
		idSeq:  time.Now().UnixNano(),
	}
}

// CreateCustomer creates a new customer
func (s *CustomerService) CreateCustomer(customer domain.Customer) (*domain.Customer, error) {
	// Generate ID if empty
	if customer.ID == "" {
		s.idSeq++
		customer.ID = strconv.FormatInt(s.idSeq, 10)
	}

	// Check if customer already exists
	if s.repo.Exists(customer.ID) {
		return nil, fmt.Errorf("customer with id %s already exists", customer.ID)
	}

	// Save to repository
	if err := s.repo.Create(customer); err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	s.logger.Infow("Customer created successfully",
		"customer_id", customer.ID,
		"name", customer.Name,
	)
	return &customer, nil
}

// RetrieveCustomer retrieves a customer by ID
func (s *CustomerService) RetrieveCustomer(id string) (*domain.Customer, error) {
	customer, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}
	return customer, nil
}

// RetrieveCustomers retrieves all customers
func (s *CustomerService) RetrieveCustomers() ([]domain.Customer, error) {
	customers, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve customers: %w", err)
	}
	return customers, nil
}

// UpdateCustomer updates an existing customer
func (s *CustomerService) UpdateCustomer(customer domain.Customer) (*domain.Customer, error) {
	// Verify customer exists
	if !s.repo.Exists(customer.ID) {
		return nil, fmt.Errorf("customer with id %s not found", customer.ID)
	}

	// Update in repository
	if err := s.repo.Update(customer); err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	s.logger.Infow("Customer updated successfully",
		"customer_id", customer.ID,
	)
	return &customer, nil
}

// DeleteCustomer removes a customer
func (s *CustomerService) DeleteCustomer(id string) error {
	// Verify customer exists
	if !s.repo.Exists(id) {
		return fmt.Errorf("customer with id %s not found", id)
	}

	// Delete it from the repository
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	s.logger.Infow("Customer deleted successfully",
		"customer_id", id,
	)
	return nil
}
