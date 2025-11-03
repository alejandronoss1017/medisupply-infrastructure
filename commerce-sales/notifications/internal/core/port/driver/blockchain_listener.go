package driver

import (
	"context"
	"notifications/internal/core/port/driven"
)

// BlockchainListener defines the interface for listening to blockchain events
type BlockchainListener interface {
	// Start begins listening for blockchain events
	Start(ctx context.Context) error

	// Stop stops the blockchain listener
	Stop() error

	// Subscribe registers an event processor to receive blockchain events
	Subscribe(processor driven.EventProcessor) error

	// IsConnected returns true if the listener is connected to the blockchain
	IsConnected() bool
}
