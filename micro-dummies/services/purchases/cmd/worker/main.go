package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/purchases/internal/adapter/queue"
)

func main() {
	// Get RabbitMQ credentials from environment variables
	rabbitUser := os.Getenv("RABBITMQ_USER")
	rabbitPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitHost := os.Getenv("RABBITMQ_HOST")

	// Initialize RabbitMQ adapter
	rabbitMQ, err := queue.NewRabbitMQ(rabbitUser, rabbitPassword, rabbitHost)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	log.Println("Starting Purchase Events Consumer...")

	// Start consuming messages
	if err := rabbitMQ.Consume("purchases_queue", rabbitMQ.HandleEvent); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	log.Println("Consumer started. Waiting for messages...")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down consumer...")
}
