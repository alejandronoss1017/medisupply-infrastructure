package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/purchases/internal/adapter/http"
	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/purchases/internal/adapter/queue"
	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/purchases/internal/core/application"
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
	service := application.NewPurchaseService(rabbitMQ, rabbitExchange)
	handler := http.NewPurchaseHandler(service)

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/ping", http.PongHandler)

	// Purchase routes
	router.GET("/", handler.GetPurchases)
	router.GET("/:id", handler.GetPurchase)
	router.POST("/", handler.PostPurchase)
	router.PUT("/:id", handler.PutPurchase)
	router.DELETE("/:id", handler.DeletePurchase)

	err = router.Run()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
