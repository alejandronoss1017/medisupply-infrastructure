package application

import (
	"context"
	"notifications/internal/core/domain"
	"notifications/internal/core/port/driven"

	"go.uber.org/zap"
)

// SLAEventProcessor is responsible for processing SLA-related events
type SLAEventProcessor struct {
	logger   *zap.Logger
	notifier driven.Notifier
}

// NewSLAEventProcessor creates a new SLAEventProcessor instance
func NewSLAEventProcessor(logger *zap.Logger, notifier driven.Notifier) *SLAEventProcessor {
	return &SLAEventProcessor{
		logger:   logger,
		notifier: notifier,
	}
}

// Process handles SLA-related events
func (p *SLAEventProcessor) Process(ctx context.Context, event interface{}) error {
	// Handle SLAAdded events
	if slaAddedEvent, ok := event.(*domain.SLAAddedEvent); ok {
		return p.processSLAAdded(ctx, slaAddedEvent)
	}

	// Handle SLAStatusUpdated events
	if slaStatusEvent, ok := event.(*domain.SLAStatusUpdatedEvent); ok {
		return p.processSLAStatusUpdated(ctx, slaStatusEvent)
	}

	// This processor only handles SLA events, ignore others
	return nil
}

// processSLAAdded processes an SLAAdded event
func (p *SLAEventProcessor) processSLAAdded(ctx context.Context, event *domain.SLAAddedEvent) error {
	p.logger.Info("Processing SLAAdded event",
		zap.String("contractId", event.ContractID),
		zap.String("slaId", event.SLAID),
		zap.Uint64("blockNumber", event.BlockNumber),
		zap.String("txHash", event.TxHash),
	)

	// Add any SLA-specific business logic here,
	// For example, validate SLA, track SLA metrics, etc.

	// Notify all registered handlers
	//var handlerErrors []error
	//for _, handler := range p.handlers {
	//	if err := handler.HandleSLAAdded(ctx, event); err != nil {
	//		p.logger.Error("Handler failed to process SLAAdded event",
	//			zap.Error(err),
	//			zap.String("contractId", event.ContractID),
	//			zap.String("slaId", event.SLAID),
	//		)
	//		handlerErrors = append(handlerErrors, err)
	//	}
	//}

	// If any handler failed, return a combined error
	//if len(handlerErrors) > 0 {
	//	return fmt.Errorf("failed to handle SLAAdded event: %d handler(s) failed", len(handlerErrors))
	//}

	//TODO: Send notification

	p.logger.Info("Successfully processed SLAAdded event",
		zap.String("contractId", event.ContractID),
		zap.String("slaId", event.SLAID),
	)

	return nil
}

// processSLAStatusUpdated processes an SLAStatusUpdated event
func (p *SLAEventProcessor) processSLAStatusUpdated(ctx context.Context, event *domain.SLAStatusUpdatedEvent) error {
	status := domain.SLAStatus(event.NewStatus)

	p.logger.Info("Processing SLAStatusUpdated event",
		zap.String("contractId", event.ContractID),
		zap.String("slaId", event.SLAID),
		zap.Uint8("newStatus", event.NewStatus),
		zap.String("statusName", status.String()),
		zap.Uint64("blockNumber", event.BlockNumber),
		zap.String("txHash", event.TxHash),
	)

	// Add any SLA status-specific business logic here,
	// For example, trigger alerts for violations, update metrics, etc.
	if status == domain.SLAStatusViolated {
		p.logger.Warn("⚠️ SLA VIOLATION DETECTED",
			zap.String("contractId", event.ContractID),
			zap.String("slaId", event.SLAID),
		)
	}

	// Notify all registered handlers
	//var handlerErrors []error
	//for _, handler := range p.handlers {
	//	if err := handler.HandleSLAStatusUpdated(ctx, event); err != nil {
	//		p.logger.Error("Handler failed to process SLAStatusUpdated event",
	//			zap.Error(err),
	//			zap.String("contractId", event.ContractID),
	//			zap.String("slaId", event.SLAID),
	//		)
	//		handlerErrors = append(handlerErrors, err)
	//	}
	//}

	// If any handler failed, return a combined error
	//if len(handlerErrors) > 0 {
	//	return fmt.Errorf("failed to handle SLAStatusUpdated event: %d handler(s) failed", len(handlerErrors))
	//}

	//TODO: Send notification

	p.logger.Info("Successfully processed SLAStatusUpdated event",
		zap.String("contractId", event.ContractID),
		zap.String("slaId", event.SLAID),
		zap.String("status", status.String()),
	)

	return nil
}
