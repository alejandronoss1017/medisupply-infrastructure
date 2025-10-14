package domain

import "time"

const (
	MedicineUpdated = "medicine.updated"
)

type MedicineEvent struct {
	EventType string    `json:"event_type"`
	Medicine  Medicine  `json:"medicine"`
	Timestamp time.Time `json:"timestamp"`
}
