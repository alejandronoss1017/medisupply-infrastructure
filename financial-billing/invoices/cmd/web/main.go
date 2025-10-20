package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/invoices/internal/adapter/http"
	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/invoices/internal/adapter/queue"
	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/invoices/internal/core/application"
)

func main() {
	// Get RabbitMQ credentials from environment variables
	rabbitUser := os.Getenv("RABBITMQ_USER")
	rabbitPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitExchange := os.Getenv("RABBITMQ_EXCHANGE")

	// Initialize RabbitMQ adapter
	rabbitMQ, err := queue.NewRabbitMQ(rabbitUser, rabbitPassword, rabbitHost)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Initialize dependencies
	service := application.NewInvoiceService(rabbitMQ, rabbitExchange)
	handler := http.NewInvoiceHandler(service)

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/ping", http.PongHandler)

	// Invoice routes
	router.GET("/", handler.GetInvoices)
	router.GET("/:id", handler.GetInvoice)
	router.POST("/", handler.PostInvoice)
	router.PUT("/:id", handler.PutInvoice)
	router.DELETE("/:id", handler.DeleteInvoice)

	err = router.Run()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
