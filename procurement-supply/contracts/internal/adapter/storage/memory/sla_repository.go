package memory

import (
	"contracts/internal/core/domain"
	"contracts/internal/core/port/driven"
	"fmt"
	"sync"
)

// SLARepository is an in-memory implementation of the Repository[string, domain.SLA] port
type SLARepository struct {
	mu   sync.RWMutex
	slas map[string]domain.SLA
}

// Ensure SLARepository implements the Repository interface
var _ driven.Repository[string, domain.SLA] = (*SLARepository)(nil)

// NewSLARepository creates a new in-memory SLA repository
func NewSLARepository() *SLARepository {
	return &SLARepository{
		slas: make(map[string]domain.SLA),
	}
}

// Create adds a new SLA to the repository
func (r *SLARepository) Create(sla domain.SLA) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.slas[sla.ID]; exists {
		return fmt.Errorf("sla with id %s already exists", sla.ID)
	}

	r.slas[sla.ID] = sla
	return nil
}

// FindByID retrieves an SLA by its ID
func (r *SLARepository) FindByID(id string) (*domain.SLA, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sla, exists := r.slas[id]
	if !exists {
		return nil, fmt.Errorf("sla with id %s not found", id)
	}

	return &sla, nil
}

// FindAll retrieves all SLAs
func (r *SLARepository) FindAll() ([]domain.SLA, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	slas := make([]domain.SLA, 0, len(r.slas))
	for _, sla := range r.slas {
		slas = append(slas, sla)
	}

	return slas, nil
}

// Update modifies an existing SLA
func (r *SLARepository) Update(sla domain.SLA) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.slas[sla.ID]; !exists {
		return fmt.Errorf("sla with id %s not found", sla.ID)
	}

	r.slas[sla.ID] = sla
	return nil
}

// Delete removes an SLA by its ID
func (r *SLARepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.slas[id]; !exists {
		return fmt.Errorf("sla with id %s not found", id)
	}

	delete(r.slas, id)
	return nil
}

// Exists checks if an SLA with the given ID exists
func (r *SLARepository) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.slas[id]
	return exists
}
