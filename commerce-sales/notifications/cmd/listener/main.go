package main

import (
	"context"
	"fmt"
	"notifications/config"
	"notifications/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"notifications/internal/adapter/blockchain"
	"notifications/internal/adapter/notification"
	"notifications/internal/core/application"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(cfg.LogLevel, cfg.LogFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting Blockchain Event Listener",
		zap.String("rpcURL", cfg.BlockchainRPCURL),
		zap.String("contractAddress", cfg.ContractAddress),
		zap.Uint64("startBlock", cfg.StartBlock),
		zap.String("logLevel", cfg.LogLevel),
	)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize SNS notifier if enabled
	log.Info("Initializing SNS notifier",
		zap.String("topicARN", cfg.SNSTopicARN),
		zap.String("region", cfg.AWSRegion),
	)

	snsNotifier, err := notification.NewSNSNotifier(ctx, cfg.AWSRegion, cfg.SNSTopicARN, log)
	if err != nil {
		log.Fatal("Failed to create SNS notifier", zap.Error(err))
	}

	log.Info("SNS notifier initialized successfully")

	// Initialize specific event processors
	contractEventProcessor := application.NewContractEventProcessor(log, snsNotifier)
	slaEventProcessor := application.NewSLAEventProcessor(log, snsNotifier)

	// Initialize Ethereum listener (inbound adapter)
	// Method 1: Pass the first processor in the constructor
	ethListener, err := blockchain.NewEthereumListener(
		cfg.BlockchainRPCURL,
		cfg.ContractAddress,
		cfg.StartBlock,
		cfg.ReconnectInterval,
		log,
		slaEventProcessor, // First processor passed in constructor
		contractEventProcessor,
	)
	if err != nil {
		log.Fatal("Failed to create Ethereum listener", zap.Error(err))
	}

	log.Info("All event processors registered successfully")

	// Start listening for events
	if err := ethListener.Start(ctx); err != nil {
		log.Fatal("Failed to start Ethereum listener", zap.Error(err))
	}

	log.Info("Blockchain Event Listener is running. Press Ctrl+C to stop.")

	// Wait for the interrupt signal for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh

	log.Info("Received shutdown signal, gracefully stopping...")

	// Cancel context to stop all listeners
	cancel()

	// Stop the Ethereum listener
	if err := ethListener.Stop(); err != nil {
		log.Error("Error stopping Ethereum listener", zap.Error(err))
	}

	log.Info("Blockchain Event Listener stopped successfully")
}
