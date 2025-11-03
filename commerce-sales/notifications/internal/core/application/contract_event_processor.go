package application

import (
	"context"
	"notifications/internal/core/domain"
	"notifications/internal/core/port/driven"

	"go.uber.org/zap"
)

// ContractEventProcessor is responsible for processing contract-related events
type ContractEventProcessor struct {
	logger   *zap.Logger
	notifier driven.Notifier
}

// NewContractEventProcessor creates a new ContractEventProcessor instance
func NewContractEventProcessor(logger *zap.Logger, notifier driven.Notifier) *ContractEventProcessor {
	return &ContractEventProcessor{
		logger:   logger,
		notifier: notifier,
	}
}

// Process handles contract-related events
func (p *ContractEventProcessor) Process(ctx context.Context, event interface{}) error {
	// Type asserts to ContractAddedEvent
	contractEvent, ok := event.(*domain.ContractAddedEvent)
	if !ok {
		// This processor only handles ContractAdded events, ignore others
		return nil
	}

	p.logger.Info("Processing ContractAdded event",
		zap.String("contractId", contractEvent.ContractID),
		zap.String("customerId", contractEvent.CustomerID),
		zap.Uint64("blockNumber", contractEvent.BlockNumber),
		zap.String("txHash", contractEvent.TxHash),
	)

	// Add any contract-specific business logic here,
	// For example, validate contract, enrich data, etc.

	// Notify all registered handlers
	//var handlerErrors []error
	//for _, handler := range p.handlers {
	//	if err := handler.HandleContractAdded(ctx, contractEvent); err != nil {
	//		p.logger.Error("Handler failed to process ContractAdded event",
	//			zap.Error(err),
	//			zap.String("contractId", contractEvent.ContractID),
	//		)
	//		handlerErrors = append(handlerErrors, err)
	//	}
	//}

	// If any handler failed, return a combined error
	//if len(handlerErrors) > 0 {
	//	return fmt.Errorf("failed to handle ContractAdded event: %d handler(s) failed", len(handlerErrors))
	//}

	// Send notification
	message := map[string]interface{}{
		"contractId": contractEvent.ContractID,
		"customerId": contractEvent.CustomerID,
		"eventType":  string(domain.EventTypeContractAdded),
	}

	if err := p.notifier.SendNotification(ctx, message, string(domain.EventTypeContractAdded)); err != nil {
		p.logger.Error("Failed to send notification for ContractAdded event",
			zap.Error(err),
			zap.String("contractId", contractEvent.ContractID),
		)
		return err
	}

	p.logger.Info("Successfully processed ContractAdded event",
		zap.String("contractId", contractEvent.ContractID),
	)

	return nil
}
