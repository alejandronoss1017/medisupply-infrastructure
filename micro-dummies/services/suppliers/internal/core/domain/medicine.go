package domain

import "time"

type Medicine struct {
	ID          string    `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Price       float64   `json:"price" binding:"required"`
	Strength    string    `json:"strength" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	SupplierID  string    `json:"supplier_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
