package main

import (
	"os"
	"suppliers/internal/adapter/http"
	"suppliers/internal/adapter/queue"
	"suppliers/internal/core/application"
	"suppliers/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize application logger
	log := logger.New("APP")
	log.Info("Starting Medicine Supplier Service...")

	// Get Env variables
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	// Setup router
	router := gin.Default()

	// Setup Kafka event publisher (driven adapter)
	log.Info("Initializing Kafka event publisher...")
	kafkaPublisher, err := queue.NewKafkaEventPublisher(kafkaHost, kafkaTopic)
	if err != nil {
		log.Fatal("Failed to create Kafka publisher: %v", err)
	}
	defer kafkaPublisher.Close()
	log.Info("Kafka event publisher initialized successfully")

	// Setup medicine service (core application)
	log.Info("Initializing medicine service...")
	medicineService := application.NewMedicineService(kafkaPublisher)

	// Setup medicine handler (driver adapter)
	log.Info("Initializing HTTP handlers...")
	medicineHandler := http.NewMedicineHandler(medicineService)

	// Health check
	router.GET("/ping", http.PongHandler)

	// Medicine routes
	medicineRoutes := router.Group("/medicines")
	{
		medicineRoutes.GET("", medicineHandler.GetMedicines)
		medicineRoutes.GET("/:id", medicineHandler.GetMedicine)
		medicineRoutes.POST("", medicineHandler.PostMedicine)
		medicineRoutes.PUT("/:id", medicineHandler.PutMedicine)
		medicineRoutes.DELETE("/:id", medicineHandler.DeleteMedicine)
	}

	log.Info("All components initialized successfully")
	log.Info("Server starting on :8080")

	if err = router.Run(); err != nil {
		log.Fatal("Failed to start server: %v", err)
	}
}
