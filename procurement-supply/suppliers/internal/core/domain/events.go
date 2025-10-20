package domain

import (
	"time"
)

type EventType string

type Event[T any] struct {
	EventType `json:"event_type"`
	Data      T         `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

const (
	MedicineCreatedEvent EventType = "medicine.created"
	MedicineUpdatedEvent EventType = "medicine.updated"
	MedicineDeletedEvent EventType = "medicine.deleted"
)

// NewMedicineEvent creates a new medicine event with validation
func NewMedicineEvent(data Medicine, event EventType) *Event[Medicine] {
	return &Event[Medicine]{
		EventType: event,
		Data:      data,
		Timestamp: time.Now(),
	}
}
