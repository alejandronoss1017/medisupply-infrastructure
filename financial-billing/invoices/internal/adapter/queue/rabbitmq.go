package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/invoices/internal/core/domain"
	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/invoices/internal/core/port/driver"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ implements the MessagePublisher interface
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQ creates a new RabbitMQ adapter
func NewRabbitMQ(user, password, host string) (*RabbitMQ, error) {
	if strings.TrimSpace(user) == "" {
		return nil, fmt.Errorf("user cannot be empty or whitespace")
	}

	if strings.TrimSpace(password) == "" {
		return nil, fmt.Errorf("password cannot be empty or whitespace")
	}

	if strings.TrimSpace(host) == "" {
		return nil, fmt.Errorf("host cannot be empty or whitespace")
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:5672/", user, password, host)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

// Publish publishes a message to an exchange with a routing key
func (r *RabbitMQ) Publish(exchange, routingKey string, body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.channel.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // persistent messages
			Timestamp:    time.Now(),
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// Consume starts consuming messages from a queue with connection monitoring
func (r *RabbitMQ) Consume(queueName string, handler func([]byte) error, done chan error) error {
	messages, err := r.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack (false for manual acknowledgment)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	// Monitor connection and channel closures
	connClosed := r.conn.NotifyClose(make(chan *amqp.Error))
	chanClosed := r.channel.NotifyClose(make(chan *amqp.Error))

	go func() {
		for {
			select {
			case msg, ok := <-messages:
				if !ok {
					log.Println("Messages channel closed")
					done <- fmt.Errorf("messages channel closed")
					return
				}
				if err = handler(msg.Body); err != nil {
					log.Printf("Error handling message: %v", err)
					// Negative acknowledgment - requeue the message
					msg.Nack(false, true)
				} else {
					// Positive acknowledgment
					msg.Ack(false)
				}
			case err := <-connClosed:
				log.Printf("RabbitMQ connection closed: %v", err)
				done <- fmt.Errorf("connection closed: %w", err)
				return
			case err := <-chanClosed:
				log.Printf("RabbitMQ channel closed: %v", err)
				done <- fmt.Errorf("channel closed: %w", err)
				return
			}
		}
	}()

	log.Printf("Started consuming from queue: %s", queueName)
	return nil
}

// Close closes the channel and connection
func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			return fmt.Errorf("failed to close channel: %w", err)
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}
	return nil
}

// IsConnected checks if the connection is still alive
func (r *RabbitMQ) IsConnected() bool {
	return r.conn != nil && !r.conn.IsClosed()
}

// HandleInvoiceEvent processes incoming invoice events
func (r *RabbitMQ) HandleInvoiceEvent(body []byte) error {
	var event domain.InvoiceEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("ERROR: Failed to unmarshal invoice event: %v", err)
		return err
	}

	// Process the event based on type
	switch event.EventType {
	case domain.InvoiceCreatedEvent:
		log.Printf("✓ CREATED - Invoice ID: %s | Buyer: %s | Subtotal: %.2f | Discount: %.2f | Taxes: %.2f | Total: %.2f",
			event.Invoice.ID,
			event.Invoice.Buyer,
			event.Invoice.Subtotal,
			event.Invoice.Discount,
			event.Invoice.Taxes,
			event.Invoice.Total,
		)

	case domain.InvoiceUpdatedEvent:
		log.Printf("↻ UPDATED - Invoice ID: %s | Buyer: %s | Subtotal: %.2f | Discount: %.2f | Taxes: %.2f | Total: %.2f",
			event.Invoice.ID,
			event.Invoice.Buyer,
			event.Invoice.Subtotal,
			event.Invoice.Discount,
			event.Invoice.Taxes,
			event.Invoice.Total,
		)

	case domain.InvoiceDeletedEvent:
		log.Printf("✗ DELETED - Invoice ID: %s",
			event.Invoice.ID,
		)

	default:
		log.Printf("? UNKNOWN - Event Type: %s", event.EventType)
	}

	return nil
}

// HandlePurchaseEvent processes incoming purchase events from the purchases microservice
func (r *RabbitMQ) HandlePurchaseEvent(body []byte, invoiceService driver.InvoiceService) error {
	var event domain.PurchaseEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("ERROR: Failed to unmarshal purchase event: %v", err)
		return err
	}

	// Process the event based on type
	switch event.EventType {
	case domain.PurchaseCreatedEvent:

	case domain.PurchaseUpdatedEvent:
		log.Printf("📦 PURCHASE UPDATED - Purchase ID: %s | Price: %.2f | Quantity: %d | Total: %.2f",
			event.Purchase.ID,
			event.Purchase.Price,
			event.Purchase.Quantity,
			event.Purchase.Total,
		)
		// Call the service to handle the update
		if err := invoiceService.ProcessPurchaseUpdated(event.Purchase); err != nil {
			log.Printf("ERROR: Failed to handle purchase update: %v", err)
			return err
		}

	case domain.PurchaseDeletedEvent:
		log.Printf("📦 PURCHASE DELETED - Purchase ID: %s",
			event.Purchase.ID,
		)
		// Call the service to handle the deletion if needed
		if err := invoiceService.ProcessPurchaseDeleted(event.Purchase.ID); err != nil {
			log.Printf("ERROR: Failed to handle purchase deletion: %v", err)
			return err
		}

	default:
		log.Printf("EVENT RECEIVED - %s", string(body))
	}

	return nil
}
