package domain

import "time"

// Event types
const (
	PurchaseCreatedEvent = "purchase.created"
	PurchaseUpdatedEvent = "purchase.updated"
	PurchaseDeletedEvent = "purchase.deleted"
)

// PurchaseEvent represents a domain event for purchases
type PurchaseEvent struct {
	EventType string    `json:"event_type"`
	Purchase  Purchase  `json:"purchase"`
	Timestamp time.Time `json:"timestamp"`
}

// NewPurchaseCreatedEvent creates a new purchase created event
func NewPurchaseCreatedEvent(purchase Purchase) PurchaseEvent {
	return PurchaseEvent{
		EventType: PurchaseCreatedEvent,
		Purchase:  purchase,
		Timestamp: time.Now(),
	}
}

// NewPurchaseUpdatedEvent creates a new purchase updated event
func NewPurchaseUpdatedEvent(purchase Purchase) PurchaseEvent {
	return PurchaseEvent{
		EventType: PurchaseUpdatedEvent,
		Purchase:  purchase,
		Timestamp: time.Now(),
	}
}

// NewPurchaseDeletedEvent creates a new purchase deleted event
func NewPurchaseDeletedEvent(purchase Purchase) PurchaseEvent {
	return PurchaseEvent{
		EventType: PurchaseDeletedEvent,
		Purchase:  purchase,
		Timestamp: time.Now(),
	}
}
