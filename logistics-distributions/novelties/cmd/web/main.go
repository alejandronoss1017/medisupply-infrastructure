package main

import (
	"context"
	"log"
	"novelties/internal/adapter/ethereum"
	"novelties/internal/adapter/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	rcpURL := os.Getenv("RCP_URL")
	smartContractAddress := os.Getenv("SMART_CONTRACT_ADDRESS")
	privateKey := os.Getenv("PRIVATE_KEY")
	abiPath := os.Getenv("ABI_PATH")

	client, err := ethereum.NewSmartContractClient(rcpURL, smartContractAddress, privateKey, abiPath)
	if err != nil {
		log.Fatalf("failed to create Ethereum client: %v", err)
	}

	ctx := context.Background()

	var value any

	err = client.InvoqueContract(ctx, "retrieveNumber", &value)
	if err != nil {
		log.Fatalf("failed to retrieve number: %v", err)
	}

	log.Println(value)

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/ping", http.PongHandler)

	if err := router.Run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
