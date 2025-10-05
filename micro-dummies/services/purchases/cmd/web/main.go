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

	// Initialize RabbitMQ adapter
	rabbitMQ, err := queue.NewRabbitMQ(rabbitUser, rabbitPassword, rabbitHost)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Initialize dependencies
	service := application.NewPurchaseService(rabbitMQ)
	handler := http.NewPurchaseHandler(service)

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/ping", http.PongHandler)

	// Purchase routes
	router.GET("/purchases", handler.GetPurchases)
	router.GET("/purchases/:id", handler.GetPurchase)
	router.POST("/purchases", handler.PostPurchase)
	router.PUT("/purchases/:id", handler.PutPurchase)
	router.DELETE("/purchases/:id", handler.DeletePurchase)

	err = router.Run()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
