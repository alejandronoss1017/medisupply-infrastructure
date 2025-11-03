package main

import (
	"contracts/internal/adapter/blockchain"
	"contracts/internal/adapter/http"
	"contracts/internal/adapter/storage/memory"
	"contracts/internal/core/application"
	"contracts/pkg/logger"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	log := logger.New("WEB-API")
	defer func() {
		// Flush logger on exit
		_ = log.Sync()
	}()

	log.Info("Starting Contracts Web API Service...")

	// Load environment variables
	rcpURL := os.Getenv("RCP_URL")
	smartContractAddress := os.Getenv("SMART_CONTRACT_ADDRESS")
	privateKey := os.Getenv("PRIVATE_KEY")

	log.Infow("Configuration loaded",
		"rcp_url", rcpURL,
		"contract_address", smartContractAddress,
	)

	// Initialize blockchain writer (for state-changing operations)
	blockchainWriter, err := blockchain.NewEthereumWriter(rcpURL, smartContractAddress, privateKey, log)
	if err != nil {
		log.Fatalw("Failed to create blockchain writer",
			"error", err,
		)
	}
	defer blockchainWriter.Close()

	// Initialize blockchain reader (for read-only operations)
	blockchainReader, err := blockchain.NewEthereumReader(rcpURL, smartContractAddress, log)
	if err != nil {
		log.Fatalw("Failed to create blockchain reader",
			"error", err,
		)
	}
	defer blockchainReader.Close()

	// Initialize repositories (driven adapters)
	contractRepo := memory.NewContractRepository()
	customerRepo := memory.NewCustomerRepository()
	slaRepo := memory.NewSLARepository()

	// Initialize application services (business logic)
	contractService := application.NewContractService(contractRepo, blockchainWriter, blockchainReader)
	customerService := application.NewCustomerService(customerRepo)
	slaService := application.NewSLAService(slaRepo, blockchainWriter, blockchainReader)

	// Initialize HTTP handlers (driver adapters)
	contractHandler := http.NewContractHandler(contractService)
	customerHandler := http.NewCustomerHandler(customerService)
	slaHandler := http.NewSLAHandler(slaService)

	// Setup router with custom middleware
	router := gin.New()

	// Add Zap logger middleware for structured request logging
	router.Use(logger.GinLogger(log))
	// Add recovery middleware with Zap logging
	router.Use(logger.GinRecovery(log))

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
		contractsRoutes.GET("/:id/slas", contractHandler.GetSLAs)
		contractsRoutes.POST("/:id/slas", contractHandler.PostSLA)
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
	log.Infow("Starting HTTP server",
		"port", "8080",
		"address", ":8080",
	)
	if err := router.Run(); err != nil {
		log.Fatalw("Failed to start server",
			"error", err,
		)
	}
}
