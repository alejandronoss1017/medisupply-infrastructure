package application

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/purchases/internal/core/domain"
	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/purchases/internal/core/port/driven"
)

var (
	ErrPurchaseNotFound = errors.New("purchase not found")
	ErrInvalidInput     = errors.New("invalid input")
)

var data = []domain.Purchase{
	{
		ID:        "1",
		Price:     10.0,
		Quantity:  1,
		Total:     10.0,
		CreatedAt: time.Now(),
	},
	{
		ID:        "2",
		Price:     20.0,
		Quantity:  2,
		Total:     20.0,
		CreatedAt: time.Now(),
	},
	{
		ID:        "3",
		Price:     30.0,
		Quantity:  3,
		Total:     30.0,
		CreatedAt: time.Now(),
	},
	{
		ID:        "4",
		Price:     40.0,
		Quantity:  4,
		Total:     40.0,
		CreatedAt: time.Now(),
	},
	{
		ID:        "5",
		Price:     50.0,
		Quantity:  5,
		Total:     50.0,
		CreatedAt: time.Now(),
	},
}

type PurchaseService struct {
	messagePublisher driven.Publisher
}

func NewPurchaseService(messagePublisher driven.Publisher) *PurchaseService {
	return &PurchaseService{
		messagePublisher: messagePublisher,
	}
}

func (s *PurchaseService) RetrievePurchases() ([]domain.Purchase, error) {
	return data, nil
}

func (s *PurchaseService) RetrievePurchase(id string) (*domain.Purchase, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	for _, purchase := range data {
		if purchase.ID == id {
			return &purchase, nil
		}
	}

	return nil, ErrPurchaseNotFound
}

func (s *PurchaseService) CreatePurchase(price float64, quantity int) (*domain.Purchase, error) {
	if price <= 0 || quantity <= 0 {
		return nil, ErrInvalidInput
	}

	purchase := &domain.Purchase{
		ID:        strconv.Itoa(len(data) + 1),
		Price:     price,
		Quantity:  quantity,
		Total:     price * float64(quantity),
		CreatedAt: time.Now(),
	}

	data = append(data, *purchase)

	// Publish event
	s.publishEvent(domain.NewPurchaseCreatedEvent(*purchase), "purchases")

	return purchase, nil
}

func (s *PurchaseService) UpdatePurchase(id string, price float64, quantity int) (*domain.Purchase, error) {
	if id == "" || price <= 0 || quantity <= 0 {
		return nil, ErrInvalidInput
	}

	var existing *domain.Purchase

	for _, purchase := range data {
		if purchase.ID == id {
			existing = &purchase
			break
		}
	}

	if existing == nil {
		return nil, ErrPurchaseNotFound
	}

	purchase := &domain.Purchase{
		ID:        id,
		Price:     price,
		Quantity:  quantity,
		Total:     price * float64(quantity),
		CreatedAt: existing.CreatedAt,
	}

	for i, p := range data {
		if p.ID == id {
			data[i] = *purchase
			break
		}
	}

	// Publish event
	s.publishEvent(domain.NewPurchaseUpdatedEvent(*purchase), "purchases")

	return purchase, nil
}

func (s *PurchaseService) DeletePurchase(id string) error {
	if id == "" {
		return ErrInvalidInput
	}

	var existing *domain.Purchase

	for _, purchase := range data {
		if purchase.ID == id {
			existing = &purchase
			break
		}
	}

	if existing == nil {
		return ErrPurchaseNotFound
	}

	for i, purchase := range data {
		if purchase.ID == id {
			data = append(data[:i], data[i+1:]...)
			// Publish event
			s.publishEvent(domain.NewPurchaseDeletedEvent(purchase), "purchases")
			break
		}
	}

	return nil
}

// publishEvent publishes a domain event to RabbitMQ
func (s *PurchaseService) publishEvent(event domain.PurchaseEvent, routingKey string) {
	if s.messagePublisher == nil {
		log.Println("Message publisher not configured, skipping event publishing")
		return
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	err = s.messagePublisher.Publish("purchases_exchange", routingKey, body)
	if err != nil {
		log.Printf("Failed to publish event: %v", err)
	}
}
