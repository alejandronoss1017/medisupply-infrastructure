package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/invoices/internal/core/domain"
)

func main() {

	config := kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"), // Kafka broker address
		"group.id":          os.Getenv("KAFKA_GROUP_ID"),          // Consumer group ID
		"auto.offset.reset": "earliest",                           // Start reading at the earliest message
	}

	// Create a new consumer
	consumer, err := kafka.NewConsumer(&config)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}
	defer consumer.Close()

	topic := os.Getenv("KAFKA_TOPIC")

	if strings.TrimSpace(topic) == "" {
		fmt.Println("KAFKA_TOPIC environment variable is not set or is empty")
		os.Exit(1)
	}

	eventMeshEndpoint := os.Getenv("EVENT_MESH_ENDPOINT")
	if strings.TrimSpace(eventMeshEndpoint) == "" {
		fmt.Println("EVENT_MESH_ENDPOINT environment variable is not set or is empty")
		os.Exit(1)
	}

	// Subscribe to the topic
	err = consumer.SubscribeTopics([]string{topic}, nil)
	// Set up a channel for handling Ctrl-C, etc. (graceful shutdown)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))

			// Here you would typically process the message, e.g., unmarshal JSON and handle it
			// For example:
			var event domain.MedicineEvent
			if err := json.Unmarshal(ev.Value, &event); err != nil {
				fmt.Printf("Failed to unmarshal message: %s\n", err)
				continue
			}

			// Create a new request
			req, err := http.NewRequest("POST", eventMeshEndpoint+"/order-system/order-broker", bytes.NewBuffer(ev.Value))
			if err != nil {
				fmt.Printf("Failed to create request: %s\n", err)
				return
			}

			// Add CloudEvent headers
			req.Header.Set("Ce-Id", "12345") // You might want to generate a unique ID
			req.Header.Set("Ce-Specversion", "1.0")
			req.Header.Set("Ce-Type", "order.created")
			req.Header.Set("Ce-Source", "order-system/order-api")
			req.Header.Set("Content-Type", "application/json")

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Failed to send request: %s\n", err)
				return
			}

			fmt.Printf("Event sent to Event Mesh, response status: %s\n", resp.Status)
			fmt.Printf("Event sent to Event Mesh, response body: %s\n", resp.Body)
			// It's important to close the response body
			defer resp.Body.Close()
		}
	}
	fmt.Println("Consumer shut down.")
}
