package driver

import "context"

// Consumer represents the driver port for consuming events
// This interface defines what our application can do (consume events)
type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
}
