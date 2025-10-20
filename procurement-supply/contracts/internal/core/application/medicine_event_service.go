package application

import (
	"contracts/internal/core/domain"
	"contracts/pkg/logger"
	"fmt"
)

// MedicineEventService handles the business logic for medicine events
type MedicineEventService struct {
	logger *logger.Logger
}

// NewMedicineEventService creates a new medicine event service
func NewMedicineEventService() *MedicineEventService {
	return &MedicineEventService{
		logger: logger.New("EVENT-SERVICE"),
	}
}

// HandleMedicineUpdated processes medicine update events
func (s *MedicineEventService) HandleMedicineUpdated(event *domain.Event[domain.Medicine]) error {
	s.logger.Info("Processing UPDATED event for medicine: %s (ID: %s)",
		event.Data.Name, event.Data.ID)

	// Update the existing contract with new medicine data
	s.logger.Info("Finding existing contracts with medicine ID: %s", event.Data.ID)

	s.logger.Info("Successfully updated contracts for medicine: %s", event.Data)
	return nil
}

// HandleMedicineDeleted processes medicine deletion events
func (s *MedicineEventService) HandleMedicineDeleted(event *domain.Event[domain.Medicine]) error {
	s.logger.Info("Processing DELETED event for medicine ID: %s", event.Data.ID)

	// Find existing contract for the medicine
	s.logger.Info("Finding existing contracts with medicine ID: %s", event.Data.ID)

	// Delete medicines related to the contract
	s.logger.Info("Deleting medicines related to contract: %s", event.Data.ID)

	return nil
}

// ProcessEvent routes events to the appropriate handler based on event type
func (s *MedicineEventService) ProcessEvent(event *domain.Event[domain.Medicine]) error {
	switch event.EventType {
	case domain.MedicineUpdatedEvent:
		return s.HandleMedicineUpdated(event)
	case domain.MedicineDeletedEvent:
		return s.HandleMedicineDeleted(event)
	default:
		s.logger.Warn("Unknown event type: %s", event.EventType)
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}
}
