package application

import (
	"fmt"
	"suppliers/internal/core/domain"
	"suppliers/internal/core/port/driven"
	"time"

	"github.com/google/uuid"
)

type MedicineService struct {
	publisher driven.EventPublisher
}

func NewMedicineService(publisher driven.EventPublisher) *MedicineService {
	return &MedicineService{
		publisher: publisher,
	}
}

func (s *MedicineService) RetrieveMedicines() ([]domain.Medicine, error) {
	// TODO: Implement repository integration to fetch medicines from database
	return []domain.Medicine{}, nil
}

func (s *MedicineService) RetrieveMedicine(id string) (*domain.Medicine, error) {
	// TODO: Implement repository integration to fetch medicine by ID from database
	return nil, fmt.Errorf("medicine not found")
}

func (s *MedicineService) CreateMedicine(medicine *domain.Medicine) (*domain.Medicine, error) {
	// Generate ID and timestamps
	medicine.ID = uuid.New().String()
	medicine.CreatedAt = time.Now()
	medicine.UpdatedAt = time.Now()

	// TODO: Implement repository integration to save medicine to database

	// Create and publish the event
	event := domain.NewMedicineEvent(*medicine, domain.MedicineCreatedEvent)

	if err := s.publisher.PublishMedicineEvent(event); err != nil {
		return nil, fmt.Errorf("failed to publish event: %w", err)
	}

	return medicine, nil
}

func (s *MedicineService) UpdateMedicine(id string, medicine *domain.Medicine) (*domain.Medicine, error) {
	// TODO: Implement repository integration to check if medicine exists and update it

	// Update timestamp
	medicine.ID = id
	medicine.UpdatedAt = time.Now()

	// Create and publish the event
	event := domain.NewMedicineEvent(*medicine, domain.MedicineUpdatedEvent)

	if err := s.publisher.PublishMedicineEvent(event); err != nil {
		return nil, fmt.Errorf("failed to publish event: %w", err)
	}

	return medicine, nil
}

func (s *MedicineService) DeleteMedicine(id string) error {
	// TODO: Implement repository integration to fetch and delete medicine from database

	// For now, create a dummy medicine object for the event
	medicine := domain.Medicine{ID: id}

	// Create and publish the event
	event := domain.NewMedicineEvent(medicine, domain.MedicineDeletedEvent)

	if err := s.publisher.PublishMedicineEvent(event); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}
