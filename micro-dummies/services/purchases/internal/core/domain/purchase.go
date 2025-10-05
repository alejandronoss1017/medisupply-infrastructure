package domain

import "time"

type Purchase struct {
	ID        string    `json:"id"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Total     float64   `json:"total"`
	CreatedAt time.Time `json:"created_at"`
}
