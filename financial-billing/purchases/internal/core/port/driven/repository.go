package driven

import "github.com/medisupply/medisupply-infrastructure/micro-dummies/services/purchases/internal/core/domain"

type PurchaseRepository interface {
	GetAll() ([]domain.Purchase, error)
	GetByID(id string) (*domain.Purchase, error)
	Create(purchase *domain.Purchase) error
	Update(purchase *domain.Purchase) error
	Delete(id string) error
}
