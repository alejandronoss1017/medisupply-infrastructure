package driven

import "contracts/internal/core/domain"

// EventProcessor defines the interface for processing events
// This allows adapters to depend on an abstraction rather than concrete implementation
type EventProcessor interface {
	ProcessEvent(event *domain.Event[domain.Medicine]) error
}
