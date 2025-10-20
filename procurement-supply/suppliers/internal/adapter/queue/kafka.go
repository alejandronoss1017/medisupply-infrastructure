package queue

import (
	"encoding/json"
	"fmt"
	"suppliers/internal/core/domain"
	"suppliers/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaEventPublisher implements the EventPublisher driven port
type KafkaEventPublisher struct {
	producer *kafka.Producer
	topic    string
	logger   *logger.Logger
}

// NewKafkaEventPublisher creates a new Kafka event publisher
func NewKafkaEventPublisher(bootstrapServers, topic string) (*KafkaEventPublisher, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &KafkaEventPublisher{
		producer: producer,
		topic:    topic,
		logger:   logger.New("KAFKA"),
	}, nil
}

// PublishMedicineEvent publishes a medicine event to Kafka
func (k *KafkaEventPublisher) PublishMedicineEvent(event *domain.Event[domain.Medicine]) error {
	k.logger.Info("Publishing event: type=%s, medicine_id=%s", event.EventType, event.Data.ID)

	// Serialize the event to JSON
	eventBytes, err := json.Marshal(event)
	if err != nil {
		k.logger.Error("Failed to marshal event: %v", err)
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	k.logger.Debug("Event serialized successfully, size=%d bytes", len(eventBytes))

	// Create the Kafka message
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &k.topic,
			Partition: kafka.PartitionAny,
		},
		Value: eventBytes,
		Key:   []byte(event.Data.ID), // Use medicine ID as the message key for partitioning
	}

	// Produce the message
	deliveryChan := make(chan kafka.Event)
	err = k.producer.Produce(message, deliveryChan)
	if err != nil {
		k.logger.Error("Failed to produce message: %v", err)
		return fmt.Errorf("failed to produce message: %w", err)
	}

	// Wait for delivery confirmation
	k.logger.Debug("Waiting for delivery confirmation...")
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		k.logger.Error("Delivery failed for medicine_id=%s: %v", event.Data.ID, m.TopicPartition.Error)
		return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
	}

	k.logger.Info("Event delivered successfully: topic=%s, partition=%d, offset=%d",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)

	close(deliveryChan)
	return nil
}

// Close closes the Kafka producer
func (k *KafkaEventPublisher) Close() {
	k.producer.Close()
}
