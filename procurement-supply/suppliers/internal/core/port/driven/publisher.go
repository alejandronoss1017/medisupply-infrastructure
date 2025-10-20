package driven

import "suppliers/internal/core/domain"

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	PublishMedicineEvent(event *domain.Event[domain.Medicine]) error
	Close()
}
