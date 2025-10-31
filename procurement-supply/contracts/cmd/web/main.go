package main

import (
	"contracts/internal/adapter/http"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Setup router
	router := gin.Default()

	contractHandler := http.NewContractHandler()
	customerHandler := http.NewCustomerHandler()
	slaHandler := http.NewSLAHandler()

	// Health check
	router.GET("/ping", http.PongHandler)

	contractsRoutes := router.Group("/contracts")
	{
		contractsRoutes.GET("", contractHandler.GetContracts)
		contractsRoutes.GET("/:id", contractHandler.GetContract)
		contractsRoutes.POST("", contractHandler.PostContract)
		contractsRoutes.PUT("/:id", contractHandler.PutContract)
		contractsRoutes.DELETE("/:id", contractHandler.DeleteContract)
	}

	customersRoutes := router.Group("/customers")
	{
		customersRoutes.GET("", customerHandler.GetCustomers)
		customersRoutes.GET("/:id", customerHandler.GetCustomer)
		customersRoutes.POST("", customerHandler.PostCustomer)
		customersRoutes.PUT("/:id", customerHandler.PutCustomer)
		customersRoutes.DELETE("/:id", customerHandler.DeleteCustomer)
	}

	slaRoutes := router.Group("/slas")
	{
		slaRoutes.GET("", slaHandler.GetSLAs)
		slaRoutes.GET("/:id", slaHandler.GetSLA)
		slaRoutes.POST("", slaHandler.PostSLA)
		slaRoutes.PUT("/:id", slaHandler.PutSLA)
		slaRoutes.DELETE("/:id", slaHandler.DeleteSLA)
	}

	if err := router.Run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
