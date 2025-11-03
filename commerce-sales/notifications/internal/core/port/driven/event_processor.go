package driven

import (
	"context"
)

// EventProcessor defines the interface for processing blockchain events
// Each implementation handles only the events it cares about
type EventProcessor interface {
	// Process handles a blockchain event
	// Implementations should type-assert to the specific event type they handle
	Process(ctx context.Context, event interface{}) error
}
