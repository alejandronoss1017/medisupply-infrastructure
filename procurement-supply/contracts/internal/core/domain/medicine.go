package domain

import "time"

type Medicine struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Strength    string    `json:"strength"`
	Category    string    `json:"category"`
	SupplierID  string    `json:"supplier_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
