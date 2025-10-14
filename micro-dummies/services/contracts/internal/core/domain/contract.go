package domain

import "time"

type Contract struct {
	ID         string         `json:"id"`
	SupplierId string         `json:"supplierId"`
	Medicines  []Medicine     `json:"medicines"`
	Terms      string         `json:"terms"`
	Value      float64        `json:"value"`
	StartDate  time.Time      `json:"startDate"`
	EndDate    time.Time      `json:"endDate"`
	Status     string         `json:"status"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
}
