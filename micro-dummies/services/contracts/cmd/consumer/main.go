package main

import (
	"context"
	"contracts/internal/adapter/queue"
	"contracts/internal/core/application"
	"contracts/pkg/logger"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	// Initialize logger
	log := logger.New("APP")

	log.Info("Starting Contracts Service...")

	kafkaHost := getEnv("KAFKA_HOST", "localhost:9092")
	kafkaGroupId := getEnv("KAFKA_GROUP_ID", "contracts-service")
	kafkaTopics := getEnv("KAFKA_TOPICS", "medicine-events")

	// Get configuration from environment variables with defaults
	config := kafka.ConfigMap{
		"bootstrap.servers":  kafkaHost,
		"group.id":           kafkaGroupId,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false, // Manual commit for better control
	}

	log.Info("Configuration loaded:")
	log.Info("  Kafka host: %s", kafkaHost)
	log.Info("  Group ID: %s", kafkaGroupId)
	log.Info("  Topics: %v", kafkaTopics)

	// Create application service (business logic layer)
	log.Info("Initializing application service...")
	eventService := application.NewMedicineEventService()

	// Create Kafka consumer adapter (infrastructure layer)
	log.Info("Initializing Kafka consumer...")
	consumer, err := queue.NewKafkaEventConsumer(config, []string{kafkaTopics}, eventService)
	if err != nil {
		log.Fatal("Failed to create Kafka consumer: %v", err)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start consumer in a goroutine
	errChan := make(chan error, 1)
	go func() {
		log.Info("Starting consumer...")
		if err = consumer.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			errChan <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case <-sigChan:
		log.Info("Shutdown signal received, stopping service...")
		cancel()
	case err = <-errChan:
		log.Error("Consumer error: %v", err)
		cancel()
	}

	// Stop consumer gracefully
	if err = consumer.Stop(); err != nil {
		log.Error("Error stopping consumer: %v", err)
	}

	log.Info("Contracts Service stopped successfully")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
