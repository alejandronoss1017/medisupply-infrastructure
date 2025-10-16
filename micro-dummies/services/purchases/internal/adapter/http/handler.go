package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/gin-gonic/gin"

	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/purchases/internal/core/application"
	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/purchases/internal/core/port/driven"
	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/purchases/internal/core/port/driver"
)

type PurchaseHandler struct {
	service   driver.PurchaseService
	publisher driven.Publisher
	exchange  string
}

func NewPurchaseHandler(service driver.PurchaseService, publisher driven.Publisher, exchange string) *PurchaseHandler {
	return &PurchaseHandler{
		service:   service,
		publisher: publisher,
		exchange:  exchange,
	}
}

func PongHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *PurchaseHandler) GetPurchases(c *gin.Context) {
	purchases, err := h.service.RetrievePurchases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, purchases)
}

func (h *PurchaseHandler) GetPurchase(c *gin.Context) {
	id := c.Param("id")

	purchase, err := h.service.RetrievePurchase(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if purchase == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "purchase not found"})
		return
	}

	c.JSON(http.StatusOK, purchase)
}

func (h *PurchaseHandler) PostPurchase(c *gin.Context) {
    // Read raw request body since the schema is currently unknown
    bodyBytes, err := io.ReadAll(c.Request.Body)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
        return
    }

    // Log the raw body for observability
    log.Printf("PostPurchase received raw body: %s", string(bodyBytes))

    // Publish raw message to RabbitMQ with a generic routing key
    if h.publisher != nil {
        routingKey := "purchase.raw"
        if err := h.publisher.Publish(h.exchange, routingKey, bodyBytes); err != nil {
            log.Printf("Failed to publish raw purchase message: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish message"})
            return
        }
    } else {
        log.Printf("RabbitMQ publisher not configured, skipping publish")
    }

    c.JSON(http.StatusAccepted, gin.H{"message": "received and processed"})
}

func (h *PurchaseHandler) PutPurchase(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Price    float64 `json:"price" binding:"required"`
		Quantity int     `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	purchase, err := h.service.UpdatePurchase(id, req.Price, req.Quantity)
	if err != nil {
		if errors.Is(err, application.ErrPurchaseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, application.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, purchase)
}

func (h *PurchaseHandler) DeletePurchase(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeletePurchase(id)
	if err != nil {
		if errors.Is(err, application.ErrPurchaseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// HandleCloudEvent processes CloudEvents received from Knative triggers
func (h *PurchaseHandler) HandleCloudEvent(c *gin.Context) {
	//TODO: TEST THIS
	// Parse the CloudEvent from the HTTP request
	var cloudEvent event.Event

	bodyBytes, _ := io.ReadAll(c.Request.Body)

	if err := cloudEvent.UnmarshalJSON(bodyBytes); err != nil {
		log.Printf("Failed to unmarshal CloudEvent: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CloudEvent format"})
		return
	}

	// Log the received CloudEvent details
	log.Printf("=== CloudEvent Received ===")
	log.Printf("Event ID: %s", cloudEvent.ID())
	log.Printf("Event Type: %s", cloudEvent.Type())
	log.Printf("Event Source: %s", cloudEvent.Source())
	log.Printf("Event Subject: %s", cloudEvent.Subject())
	log.Printf("Event Time: %s", cloudEvent.Time())
	log.Printf("Event Data Content Type: %s", cloudEvent.DataContentType())

	// Log the event data
	if cloudEvent.Data() != nil {
		dataBytes, err := json.MarshalIndent(cloudEvent.Data(), "", "  ")
		if err != nil {
			log.Printf("Event Data (raw): %v", cloudEvent.Data())
		} else {
			log.Printf("Event Data:\n%s", string(dataBytes))
		}
	}

	// Log all extensions/attributes
	log.Printf("Event Extensions:")
	for key, value := range cloudEvent.Extensions() {
		log.Printf("  %s: %v", key, value)
	}

	// Log CloudEvent headers for debugging
	log.Printf("CloudEvent Headers:")
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			log.Printf("  %s: %s", key, values[0])
		}
	}

	log.Printf("=== End CloudEvent ===")

	// Publish the CloudEvent to RabbitMQ
	if h.publisher != nil {
		// Convert CloudEvent to JSON for RabbitMQ
		eventBytes, err := json.Marshal(cloudEvent)
		if err != nil {
			log.Printf("Failed to marshal CloudEvent for RabbitMQ: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process event"})
			return
		}

		// Determine routing key based on event type
		routingKey := "cloud.event"
		if cloudEvent.Type() != "" {
			routingKey = fmt.Sprintf("cloud.%s", cloudEvent.Type())
		}

		// Publish to RabbitMQ
		err = h.publisher.Publish(h.exchange, routingKey, eventBytes)
		if err != nil {
			log.Printf("Failed to publish CloudEvent to RabbitMQ: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish event"})
			return
		}

		log.Printf("✓ CloudEvent published to RabbitMQ - Exchange: %s, Routing Key: %s", h.exchange, routingKey)
	} else {
		log.Printf("⚠️ RabbitMQ publisher not configured, skipping event publishing")
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"message":   "CloudEvent received, logged, and published to RabbitMQ successfully",
		"eventId":   cloudEvent.ID(),
		"eventType": cloudEvent.Type(),
		"published": h.publisher != nil,
	})
}
