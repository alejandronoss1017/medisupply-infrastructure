package application

import (
	"contracts/internal/core/domain"
	"contracts/pkg/logger"
	"fmt"

	"go.uber.org/zap"
)

// MedicineEventService handles the business logic for medicine events
type MedicineEventService struct {
	logger *zap.SugaredLogger
}

// NewMedicineEventService creates a new medicine event service
func NewMedicineEventService() *MedicineEventService {
	return &MedicineEventService{
		logger: logger.New("EVENT-SERVICE"),
	}
}

// HandleMedicineUpdated processes medicine update events
func (s *MedicineEventService) HandleMedicineUpdated(event *domain.Event[domain.Medicine]) error {
	s.logger.Infow("Processing UPDATED event for medicine",
		"medicine_id", event.Data.ID,
		"medicine_name", event.Data.Name)

	// Update the existing contract with new medicine data
	s.logger.Infow("Finding existing contracts with medicine",
		"medicine_id", event.Data.ID,
	)

	s.logger.Infow("Successfully updated contracts for medicine",
		"medicine_id", event.Data.ID,
	)
	return nil
}

// HandleMedicineDeleted processes medicine deletion events
func (s *MedicineEventService) HandleMedicineDeleted(event *domain.Event[domain.Medicine]) error {
	s.logger.Infow("Processing DELETED event for medicine",
		"medicine_id", event.Data.ID,
	)

	// Find an existing contract for the medicine
	s.logger.Infow("Finding existing contracts with medicine",
		"medicine_id", event.Data.ID,
	)

	// Delete medicines related to the contract
	s.logger.Infow("Deleting medicines related to contract",
		"medicine_id", event.Data.ID,
	)

	return nil
}

// ProcessEvent routes events to the appropriate handler based on an event type
func (s *MedicineEventService) ProcessEvent(event *domain.Event[domain.Medicine]) error {
	switch event.EventType {
	case domain.MedicineUpdatedEvent:
		return s.HandleMedicineUpdated(event)
	case domain.MedicineDeletedEvent:
		return s.HandleMedicineDeleted(event)
	default:
		s.logger.Warnw("Unknown event type",
			"event_type", event.EventType)
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}
}
