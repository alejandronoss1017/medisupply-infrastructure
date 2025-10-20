package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/invoices/internal/adapter/queue"
	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/invoices/internal/core/application"
)

func main() {
	// Get RabbitMQ credentials from environment variables
	rabbitUser := os.Getenv("RABBITMQ_USER")
	rabbitPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitExchange := os.Getenv("RABBITMQ_EXCHANGE")

	// Queue for consuming purchase events from the purchases microservice
	rabbitQueue := os.Getenv("PURCHASES_QUEUE")
	if rabbitQueue == "" {
		rabbitQueue = "purchases.events" // Default queue name
	}

	// Initialize RabbitMQ adapter
	rabbitMQ, err := queue.NewRabbitMQ(rabbitUser, rabbitPassword, rabbitHost)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Initialize invoice service (needed for handling purchase updates)
	invoiceService := application.NewInvoiceService(rabbitMQ, rabbitExchange)

	log.Println("Starting Invoice Service Events Consumer...")
	log.Printf("  - Consuming from purchases queue: %s", rabbitQueue)

	// Channel to receive errors from consumer (connection issues)
	done := make(chan error, 1)

	// Start consuming purchase events
	if err = rabbitMQ.Consume(rabbitQueue, func(body []byte) error {
		return rabbitMQ.HandlePurchaseEvent(body, invoiceService)
	}, done); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	log.Println("Consumer started. Waiting for purchase update events...")

	// Wait for interrupt signal or connection error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Shutting down consumer gracefully...")
	case err := <-done:
		log.Printf("Consumer stopped due to error: %v", err)
		log.Println("Exiting... (process manager should restart the service)")
	}
}
