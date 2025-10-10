package driver

import "github.com/medisupply/medisupply-infrastructure/micro-dummies/services/invoices/internal/core/domain"

type InvoiceService interface {
	RetrieveInvoices() ([]domain.Invoice, error)
	RetrieveInvoice(id string) (*domain.Invoice, error)
	CreateInvoice(purchasesID []string, buyer string, subtotal, discount, taxes float64) (*domain.Invoice, error)
	UpdateInvoice(id string, purchasesID []string, buyer string, subtotal, discount, taxes float64) (*domain.Invoice, error)
	DeleteInvoice(id string) error
	ProcessPurchaseUpdated(purchase domain.Purchase) error
	ProcessPurchaseDeleted(id string) error
}
