package domain

import "time"

// Event types
const (
	InvoiceCreatedEvent = "invoice.created"
	InvoiceUpdatedEvent = "invoice.updated"
	InvoiceDeletedEvent = "invoice.deleted"
)

// InvoiceEvent represents a domain event for invoices
type InvoiceEvent struct {
	EventType string    `json:"event_type"`
	Invoice   Invoice   `json:"invoice"`
	Timestamp time.Time `json:"timestamp"`
}

// NewInvoiceCreatedEvent creates a new invoice created event
func NewInvoiceCreatedEvent(invoice Invoice) InvoiceEvent {
	return InvoiceEvent{
		EventType: InvoiceCreatedEvent,
		Invoice:   invoice,
		Timestamp: time.Now(),
	}
}

// NewInvoiceUpdatedEvent creates a new invoice updated event
func NewInvoiceUpdatedEvent(invoice Invoice) InvoiceEvent {
	return InvoiceEvent{
		EventType: InvoiceUpdatedEvent,
		Invoice:   invoice,
		Timestamp: time.Now(),
	}
}

// NewInvoiceDeletedEvent creates a new invoice deleted event
func NewInvoiceDeletedEvent(invoice Invoice) InvoiceEvent {
	return InvoiceEvent{
		EventType: InvoiceDeletedEvent,
		Invoice:   invoice,
		Timestamp: time.Now(),
	}
}
