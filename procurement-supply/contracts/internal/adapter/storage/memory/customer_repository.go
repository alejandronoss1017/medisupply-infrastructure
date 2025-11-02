package memory

import (
	"contracts/internal/core/domain"
	"contracts/internal/core/port/driven"
	"fmt"
	"sync"
)

// CustomerRepository is an in-memory implementation of the Repository[string, domain.Customer] port
type CustomerRepository struct {
	mu        sync.RWMutex
	customers map[string]domain.Customer
}

// Ensure CustomerRepository implements the Repository interface
var _ driven.Repository[string, domain.Customer] = (*CustomerRepository)(nil)

// NewCustomerRepository creates a new in-memory customer repository
func NewCustomerRepository() *CustomerRepository {
	return &CustomerRepository{
		customers: make(map[string]domain.Customer),
	}
}

// Create adds a new customer to the repository
func (r *CustomerRepository) Create(customer domain.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.customers[customer.ID]; exists {
		return fmt.Errorf("customer with id %s already exists", customer.ID)
	}

	r.customers[customer.ID] = customer
	return nil
}

// FindByID retrieves a customer by its ID
func (r *CustomerRepository) FindByID(id string) (*domain.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	customer, exists := r.customers[id]
	if !exists {
		return nil, fmt.Errorf("customer with id %s not found", id)
	}

	return &customer, nil
}

// FindAll retrieves all customers
func (r *CustomerRepository) FindAll() ([]domain.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	customers := make([]domain.Customer, 0, len(r.customers))
	for _, customer := range r.customers {
		customers = append(customers, customer)
	}

	return customers, nil
}

// Update modifies an existing customer
func (r *CustomerRepository) Update(customer domain.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.customers[customer.ID]; !exists {
		return fmt.Errorf("customer with id %s not found", customer.ID)
	}

	r.customers[customer.ID] = customer
	return nil
}

// Delete removes a customer by its ID
func (r *CustomerRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.customers[id]; !exists {
		return fmt.Errorf("customer with id %s not found", id)
	}

	delete(r.customers, id)
	return nil
}

// Exists checks if a customer with the given ID exists
func (r *CustomerRepository) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.customers[id]
	return exists
}
