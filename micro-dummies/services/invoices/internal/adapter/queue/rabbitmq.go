package queue

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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

// Consume starts consuming messages from a queue
func (r *RabbitMQ) Consume(queueName string, handler func([]byte) error) error {
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

	go func() {
		for msg := range messages {
			if err = handler(msg.Body); err != nil {
				log.Printf("Error handling message: %v", err)
				// Negative acknowledgment - requeue the message
				msg.Nack(false, true)
			} else {
				// Positive acknowledgment
				msg.Ack(false)
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
	log.Printf("📄 INVOICE EVENT RECEIVED: %s", string(body))
	return nil
}

// HandlePurchaseEvent processes incoming purchase events from the purchases microservice
func (r *RabbitMQ) HandlePurchaseEvent(body []byte, invoiceService driver.InvoiceService) error {
	log.Printf("📦 PURCHASE EVENT RECEIVED: %s", string(body))
	return nil
}
