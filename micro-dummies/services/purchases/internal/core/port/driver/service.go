package driver

import "github.com/medisupply/medisupply-infrastructure/micro-dummies/services/purchases/internal/core/domain"

type PurchaseService interface {
	RetrievePurchases() ([]domain.Purchase, error)
	RetrievePurchase(id string) (*domain.Purchase, error)
	CreatePurchase(price float64, quantity int) (*domain.Purchase, error)
	UpdatePurchase(id string, price float64, quantity int) (*domain.Purchase, error)
	DeletePurchase(id string) error
}
