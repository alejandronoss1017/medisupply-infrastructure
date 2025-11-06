package main

import (
	"context"
	"math/big"
	"novelties/internal/adapter/blockchain"
	"novelties/internal/adapter/http"
	"novelties/pkg/logger"
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

	rcpURL := os.Getenv("BLOCKCHAIN_RPC_URL")
	contractAddress := os.Getenv("CONTRACT_ADDRESS")
	privateKey := os.Getenv("PRIVATE_KEY")

	blockchainWriter, err := blockchain.NewEthereumWriter(rcpURL, contractAddress, privateKey, log)
	if err != nil {
		log.Fatalw("Failed to create blockchain writer",
			"error", err,
		)
	}
	defer blockchainWriter.Close()

	handler := http.NewNoveltyHandler(blockchainWriter, log)

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/ping", http.PongHandler)

	router.POST("/novelties", handler.PostNovelty)

	if err := router.Run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
