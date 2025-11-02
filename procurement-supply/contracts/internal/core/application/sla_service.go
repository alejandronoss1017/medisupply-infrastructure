package application

import (
	"contracts/internal/core/domain"
	"contracts/internal/core/port/driven"
	"contracts/internal/core/port/driver"
	"contracts/pkg/logger"
	"fmt"
	"strconv"
	"time"
)

// SLAService handles business logic for SLA management
type SLAService struct {
	repo   driven.Repository[string, domain.SLA]
	logger *logger.Logger
	idSeq  int64
}

// Ensure SLAService implements the driver.SLAService interface
var _ driver.SLAService = (*SLAService)(nil)

// NewSLAService creates a new SLA service
func NewSLAService(repo driven.Repository[string, domain.SLA]) *SLAService {
	return &SLAService{
		repo:   repo,
		logger: logger.New("SLA-SERVICE"),
		idSeq:  time.Now().UnixNano(),
	}
}

// CreateSLA creates a new SLA
func (s *SLAService) CreateSLA(sla domain.SLA) (*domain.SLA, error) {
	// Generate ID if empty
	if sla.ID == "" {
		s.idSeq++
		sla.ID = strconv.FormatInt(s.idSeq, 10)
	}

	// Check if SLA already exists
	if s.repo.Exists(sla.ID) {
		return nil, fmt.Errorf("sla with id %s already exists", sla.ID)
	}

	// Save to repository
	if err := s.repo.Create(sla); err != nil {
		return nil, fmt.Errorf("failed to create sla: %w", err)
	}

	s.logger.Info("SLA %s created successfully", sla.ID)
	return &sla, nil
}

// RetrieveSLA retrieves an SLA by ID
func (s *SLAService) RetrieveSLA(id string) (*domain.SLA, error) {
	sla, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("sla not found: %w", err)
	}
	return sla, nil
}

// RetrieveSLAs retrieves all SLAs
func (s *SLAService) RetrieveSLAs() ([]domain.SLA, error) {
	slas, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve slas: %w", err)
	}
	return slas, nil
}

// UpdateSLA updates an existing SLA
func (s *SLAService) UpdateSLA(sla domain.SLA) (*domain.SLA, error) {
	// Verify SLA exists
	if !s.repo.Exists(sla.ID) {
		return nil, fmt.Errorf("sla with id %s not found", sla.ID)
	}

	// Update in repository
	if err := s.repo.Update(sla); err != nil {
		return nil, fmt.Errorf("failed to update sla: %w", err)
	}

	s.logger.Info("SLA %s updated successfully", sla.ID)
	return &sla, nil
}

// DeleteSLA removes an SLA
func (s *SLAService) DeleteSLA(id string) error {
	// Verify SLA exists
	if !s.repo.Exists(id) {
		return fmt.Errorf("sla with id %s not found", id)
	}

	// Delete from repository
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete sla: %w", err)
	}

	s.logger.Info("SLA %s deleted successfully", id)
	return nil
}
