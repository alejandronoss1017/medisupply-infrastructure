package domain

import "time"

// Purchase represents a purchase entity (from purchases microservice)
type Purchase struct {
	ID        string    `json:"id"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Total     float64   `json:"total"`
	CreatedAt time.Time `json:"created_at"`
}

// Event types for purchases (consumed from purchases microservice)
const (
	PurchaseCreatedEvent = "purchase.created"
	PurchaseUpdatedEvent = "purchase.updated"
	PurchaseDeletedEvent = "purchase.deleted"
)

// PurchaseEvent represents a domain event from the purchases microservice
type PurchaseEvent struct {
	EventType string    `json:"event_type"`
	Purchase  Purchase  `json:"purchase"`
	Timestamp time.Time `json:"timestamp"`
}
