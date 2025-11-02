package main

import (
	"context"
	"log"
	"math/big"
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

	log.Println(value)

	observed := big.NewInt(12345)
	slaId := big.NewInt(0)
	note := "Test note"

	tx, err := client.SendContractTransaction(ctx, "reportMetric", nil, slaId, observed, note)
	if err != nil {
		log.Fatalf("failed to report metric: %v", err)
	}

	_, err = client.WaitTransaction(ctx, tx)
	if err != nil {
		log.Fatalf("failed to wait for tx: %v", err)
	}

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/ping", http.PongHandler)

	if err := router.Run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
