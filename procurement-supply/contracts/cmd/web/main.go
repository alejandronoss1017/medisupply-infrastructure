package main

import (
	"contracts/internal/adapter/blockchain"
	"contracts/internal/adapter/http"
	"contracts/internal/adapter/storage/memory"
	"contracts/internal/core/application"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	// Load environment variables
	rcpURL := os.Getenv("RCP_URL")
	smartContractAddress := os.Getenv("SMART_CONTRACT_ADDRESS")
	privateKey := os.Getenv("PRIVATE_KEY")

	// Initialize blockchain writer (for state-changing operations)
	blockchainWriter, err := blockchain.NewEthereumWriter(rcpURL, smartContractAddress, privateKey)
	if err != nil {
		log.Fatalf("failed to create blockchain writer: %v", err)
	}
	defer blockchainWriter.Close()

	// Initialize blockchain reader (for read-only operations)
	blockchainReader, err := blockchain.NewEthereumReader(rcpURL, smartContractAddress)
	if err != nil {
		log.Fatalf("failed to create blockchain reader: %v", err)
	}
	defer blockchainReader.Close()

	// Initialize repositories (driven adapters)
	contractRepo := memory.NewContractRepository()
	customerRepo := memory.NewCustomerRepository()
	slaRepo := memory.NewSLARepository()

	// Initialize application services (business logic)
	contractService := application.NewContractService(contractRepo, blockchainWriter)
	customerService := application.NewCustomerService(customerRepo)
	slaService := application.NewSLAService(slaRepo)

	// Initialize HTTP handlers (driver adapters)
	contractHandler := http.NewContractHandler(contractService)
	customerHandler := http.NewCustomerHandler(customerService)
	slaHandler := http.NewSLAHandler(slaService)

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/ping", http.PongHandler)

	// Contract routes
	contractsRoutes := router.Group("/contracts")
	{
		contractsRoutes.GET("", contractHandler.GetContracts)
		contractsRoutes.GET("/:id", contractHandler.GetContract)
		contractsRoutes.POST("", contractHandler.PostContract)
		contractsRoutes.PUT("/:id", contractHandler.PutContract)
		contractsRoutes.DELETE("/:id", contractHandler.DeleteContract)
	}

	// Customer routes
	customersRoutes := router.Group("/customers")
	{
		customersRoutes.GET("", customerHandler.GetCustomers)
		customersRoutes.GET("/:id", customerHandler.GetCustomer)
		customersRoutes.POST("", customerHandler.PostCustomer)
		customersRoutes.PUT("/:id", customerHandler.PutCustomer)
		customersRoutes.DELETE("/:id", customerHandler.DeleteCustomer)
	}

	// SLA routes
	slaRoutes := router.Group("/slas")
	{
		slaRoutes.GET("", slaHandler.GetSLAs)
		slaRoutes.GET("/:id", slaHandler.GetSLA)
		slaRoutes.POST("", slaHandler.PostSLA)
		slaRoutes.PUT("/:id", slaHandler.PutSLA)
		slaRoutes.DELETE("/:id", slaHandler.DeleteSLA)
	}

	// Start server
	log.Println("Starting server on :8080...")
	if err := router.Run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
