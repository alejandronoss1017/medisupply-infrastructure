package application

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/invoices/internal/core/domain"
	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/invoices/internal/core/port/driven"
)

var (
	ErrInvoiceNotFound = errors.New("invoice not found")
	ErrInvalidInput    = errors.New("invalid input")
)

var data = []domain.Invoice{
	{
		ID:          "1",
		PurchasesID: []string{"1", "2"},
		Buyer:       "John Doe",
		Subtotal:    100.0,
		Discount:    10.0,
		Taxes:       13.5,
		Total:       103.5,
		CreatedAt:   time.Now(),
	},
	{
		ID:          "2",
		PurchasesID: []string{"3"},
		Buyer:       "Jane Smith",
		Subtotal:    200.0,
		Discount:    20.0,
		Taxes:       27.0,
		Total:       207.0,
		CreatedAt:   time.Now(),
	},
	{
		ID:          "3",
		PurchasesID: []string{"4", "5"},
		Buyer:       "Bob Johnson",
		Subtotal:    150.0,
		Discount:    15.0,
		Taxes:       20.25,
		Total:       155.25,
		CreatedAt:   time.Now(),
	},
}

type InvoiceService struct {
	messagePublisher driven.Publisher
	exchange         string
}

func NewInvoiceService(messagePublisher driven.Publisher, exchange string) *InvoiceService {
	return &InvoiceService{
		messagePublisher: messagePublisher,
		exchange:         exchange,
	}
}

func (s *InvoiceService) RetrieveInvoices() ([]domain.Invoice, error) {
	return data, nil
}

func (s *InvoiceService) RetrieveInvoice(id string) (*domain.Invoice, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	for _, invoice := range data {
		if invoice.ID == id {
			return &invoice, nil
		}
	}

	return nil, ErrInvoiceNotFound
}

func (s *InvoiceService) CreateInvoice(purchasesID []string, buyer string, subtotal, discount, taxes float64) (*domain.Invoice, error) {
	if len(purchasesID) == 0 || buyer == "" || subtotal < 0 || discount < 0 || taxes < 0 {
		return nil, ErrInvalidInput
	}

	total := subtotal - discount + taxes

	invoice := &domain.Invoice{
		ID:          strconv.Itoa(len(data) + 1),
		PurchasesID: purchasesID,
		Buyer:       buyer,
		Subtotal:    subtotal,
		Discount:    discount,
		Taxes:       taxes,
		Total:       total,
		CreatedAt:   time.Now(),
	}

	data = append(data, *invoice)

	// Publish event
	s.publishEvent(domain.NewInvoiceCreatedEvent(*invoice), s.exchange)

	return invoice, nil
}

func (s *InvoiceService) UpdateInvoice(id string, purchasesID []string, buyer string, subtotal, discount, taxes float64) (*domain.Invoice, error) {
	if id == "" || len(purchasesID) == 0 || buyer == "" || subtotal < 0 || discount < 0 || taxes < 0 {
		return nil, ErrInvalidInput
	}

	var existing *domain.Invoice

	for _, invoice := range data {
		if invoice.ID == id {
			existing = &invoice
			break
		}
	}

	if existing == nil {
		return nil, ErrInvoiceNotFound
	}

	total := subtotal - discount + taxes

	invoice := &domain.Invoice{
		ID:          id,
		PurchasesID: purchasesID,
		Buyer:       buyer,
		Subtotal:    subtotal,
		Discount:    discount,
		Taxes:       taxes,
		Total:       total,
		CreatedAt:   existing.CreatedAt,
	}

	for i, inv := range data {
		if inv.ID == id {
			data[i] = *invoice
			break
		}
	}

	// Publish event
	s.publishEvent(domain.NewInvoiceUpdatedEvent(*invoice), s.exchange)

	return invoice, nil
}

func (s *InvoiceService) DeleteInvoice(id string) error {
	if id == "" {
		return ErrInvalidInput
	}

	var existing *domain.Invoice

	for _, invoice := range data {
		if invoice.ID == id {
			existing = &invoice
			break
		}
	}

	if existing == nil {
		return ErrInvoiceNotFound
	}

	for i, invoice := range data {
		if invoice.ID == id {
			data = append(data[:i], data[i+1:]...)
			// Publish event
			s.publishEvent(domain.NewInvoiceDeletedEvent(invoice), s.exchange)
			break
		}
	}

	return nil
}

// ProcessPurchaseUpdated process purchase.updated events from the purchases microservice
// It updates all invoices that contain the modified purchase and recalculates totals
func (s *InvoiceService) ProcessPurchaseUpdated(purchase domain.Purchase) error {
	log.Printf("📦 Processing purchase.updated event for Purchase ID: %s", purchase.ID)

	updatedCount := 0
	for i, invoice := range data {
		// Check if this invoice contains the updated purchase
		containsPurchase := false
		for _, purchaseID := range invoice.PurchasesID {
			if purchaseID == purchase.ID {
				containsPurchase = true
				break
			}
		}

		if containsPurchase {
			// In a real system, you would recalculate the invoice based on updated purchase data
			// For now, we'll just log the update
			log.Printf("↻ Invoice ID %s contains Purchase ID %s - would recalculate totals", invoice.ID, purchase.ID)
			log.Printf("  Purchase updated: Price=%.2f, Quantity=%d, Total=%.2f", purchase.Price, purchase.Quantity, purchase.Total)

			// Optionally publish an event that this invoice was affected
			updatedInvoice := data[i]
			s.publishEvent(domain.NewInvoiceUpdatedEvent(updatedInvoice), s.exchange)
			updatedCount++
		}
	}

	if updatedCount > 0 {
		log.Printf("✓ Updated %d invoice(s) affected by Purchase ID: %s", updatedCount, purchase.ID)
	} else {
		log.Printf("ℹ No invoices found containing Purchase ID: %s", purchase.ID)
	}

	return nil
}

// ProcessPurchaseDeleted process purchase.deleted events from the purchases microservice
// It removes the deleted purchase from all invoices that contain it and recalculates totals
func (s *InvoiceService) ProcessPurchaseDeleted(id string) error {
	log.Printf("📦 Processing purchase.deleted event for Purchase ID: %s", id)

	updatedCount := 0
	for i, invoice := range data {
		// Check if this invoice contains the deleted purchase
		containsPurchase := false
		for _, pid := range invoice.PurchasesID {
			if pid == id {
				containsPurchase = true
				break
			}
		}

		if containsPurchase {
			// Remove the purchase ID from the invoice's PurchasesID slice
			var newPurchasesID []string
			for _, pid := range invoice.PurchasesID {
				if pid != id {
					newPurchasesID = append(newPurchasesID, pid)
				}
			}
			data[i].PurchasesID = newPurchasesID
		}
		// In a real system, you would recalculate the invoice based on remaining purchases
		// For now, we'll just log the update
		log.Printf("↻ Invoice ID %s contained deleted Purchase ID %s - would recalculate totals", invoice.ID, id)

		// Optionally publish an event that this invoice was affected
		updatedInvoice := data[i]
		s.publishEvent(domain.NewInvoiceUpdatedEvent(updatedInvoice), s.exchange)
		updatedCount++
	}
	if updatedCount > 0 {
		log.Printf("✓ Updated %d invoice(s) affected by deleted Purchase ID: %s", updatedCount, id)
	}
	return nil
}

// publishEvent publishes a domain event to RabbitMQ
func (s *InvoiceService) publishEvent(event domain.InvoiceEvent, routingKey string) {
	if s.messagePublisher == nil {
		log.Println("Message publisher not configured, skipping event publishing")
		return
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	err = s.messagePublisher.Publish(s.exchange, routingKey, body)
	if err != nil {
		log.Printf("Failed to publish event: %v", err)
	}
}
