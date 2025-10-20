package domain

import "time"

type Invoice struct {
	ID          string    `json:"id"`
	PurchasesID []string  `json:"purchases"`
	Buyer       string    `json:"buyer"`
	Subtotal    float64   `json:"subtotal"`
	Discount    float64   `json:"discount"`
	Taxes       float64   `json:"taxes"`
	Total       float64   `json:"total"`
	CreatedAt   time.Time `json:"created_at"`
}
