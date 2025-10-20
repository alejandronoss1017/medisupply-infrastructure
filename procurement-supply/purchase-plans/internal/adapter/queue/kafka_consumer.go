package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"purchase-plans/internal/core/domain"
	"purchase-plans/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaEventConsumer is the adapter that implements the EventConsumer driver port
type KafkaEventConsumer struct {
	consumer     *kafka.Consumer
	eventService EventProcessor
	logger       *logger.Logger
	topics       []string
}

// EventProcessor defines the interface for processing events
// This allows the adapter to depend on an abstraction rather than concrete implementation
type EventProcessor interface {
	ProcessEvent(event *domain.Event[domain.Medicine]) error
}

// Config holds the Kafka consumer configuration
type Config struct {
	BootstrapServers string
	GroupID          string
	Topics           []string
	AutoOffsetReset  string
}

// NewKafkaEventConsumer creates a new Kafka event consumer adapter
func NewKafkaEventConsumer(config kafka.ConfigMap, topics []string, eventService EventProcessor) (*KafkaEventConsumer, error) {
	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return &KafkaEventConsumer{
		consumer:     consumer,
		eventService: eventService,
		logger:       logger.New("KAFKA-CONSUMER"),
		topics:       topics,
	}, nil
}

// Start begins consuming messages from Kafka
func (k *KafkaEventConsumer) Start(ctx context.Context) error {
	// Subscribe to topics
	if err := k.consumer.SubscribeTopics(k.topics, nil); err != nil {
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	k.logger.Info("Subscribed to topics: %v", k.topics)
	k.logger.Info("Started consuming messages. Waiting for events...")

	// Message consumption loop
	for {
		select {
		case <-ctx.Done():
			k.logger.Info("Context cancelled, stopping consumer...")
			return ctx.Err()
		default:
			if err := k.pollAndProcess(ctx); err != nil {
				k.logger.Error("Error in poll and process: %v", err)
				// Continue processing other messages
			}
		}
	}
}

// pollAndProcess polls for a message and processes it
func (k *KafkaEventConsumer) pollAndProcess(ctx context.Context) error {
	ev := k.consumer.Poll(100)
	if ev == nil {
		return nil
	}

	switch e := ev.(type) {
	case *kafka.Message:
		return k.handleMessage(e)

	case kafka.Error:
		k.logger.Error("Kafka error: %v (code: %v)", e, e.Code())

		// Check if it's a fatal error
		if e.Code() == kafka.ErrAllBrokersDown {
			return fmt.Errorf("all brokers are down")
		}

	case *kafka.Stats:
		// Optionally log statistics
		k.logger.Debug("Kafka stats received")

	default:
		// Ignore other event types (partition assignment, etc.)
	}

	return nil
}

// handleMessage processes a single Kafka message
func (k *KafkaEventConsumer) handleMessage(msg *kafka.Message) error {
	k.logger.Debug("Received message on %s [partition %d] at offset %d",
		*msg.TopicPartition.Topic,
		msg.TopicPartition.Partition,
		msg.TopicPartition.Offset)

	// Parse the message
	var event domain.Event[domain.Medicine]
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		k.logger.Error("Failed to unmarshal message: %v", err)
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	k.logger.Info("Processing event: type=%s, medicine_id=%s, medicine_name=%s",
		event.EventType, event.Data.ID, event.Data.Name)

	// Process the event through the application service
	if err := k.eventService.ProcessEvent(&event); err != nil {
		k.logger.Error("Error processing event: %v", err)
		// Depending on your error handling strategy, you might want to:
		// - Return the error to stop processing
		// - Continue processing other messages
		// - Send to a dead letter queue
		return err
	}

	// Commit offset after successful processing
	if _, err := k.consumer.CommitMessage(msg); err != nil {
		k.logger.Error("Error committing offset: %v", err)
		return fmt.Errorf("failed to commit offset: %w", err)
	}

	k.logger.Debug("Successfully processed and committed message at offset %d",
		msg.TopicPartition.Offset)

	return nil
}

// Stop gracefully stops the consumer
func (k *KafkaEventConsumer) Stop() error {
	k.logger.Info("Closing Kafka consumer...")

	if err := k.consumer.Close(); err != nil {
		k.logger.Error("Error closing consumer: %v", err)
		return fmt.Errorf("failed to close consumer: %w", err)
	}

	k.logger.Info("Kafka consumer closed successfully")
	return nil
}
