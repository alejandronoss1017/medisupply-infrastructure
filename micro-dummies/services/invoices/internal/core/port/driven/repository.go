package driven

import "github.com/medisupply/medisupply-infrastructure/micro-dummies/services/invoices/internal/core/domain"

type InvoiceRepository interface {
	GetAll() ([]domain.Invoice, error)
	GetByID(id string) (*domain.Invoice, error)
	Create(invoice *domain.Invoice) error
	Update(invoice *domain.Invoice) error
	Delete(id string) error
}
