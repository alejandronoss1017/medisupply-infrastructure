package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"notifications/internal/core/port/driven"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	"notifications/internal/adapter/blockchain/binding"
	"notifications/internal/core/domain"
)

// EthereumListener is an adapter that listens to Ethereum blockchain events
type EthereumListener struct {
	client            *ethclient.Client
	contractAddress   common.Address
	contractABI       abi.ABI
	rpcURL            string
	reconnectInterval time.Duration
	startBlock        uint64
	logger            *zap.Logger
	processors        []driven.EventProcessor

	mu           sync.RWMutex
	processorsMu sync.RWMutex
	connected    bool
	stopCh       chan struct{}
}

// NewEthereumListener creates a new EthereumListener instance
// The parameter "processors" is optional and can be nil. Additional processors can be added via Subscribe().
func NewEthereumListener(
	rpcURL string,
	contractAddress string,
	startBlock uint64,
	reconnectInterval int,
	logger *zap.Logger,
	processors ...driven.EventProcessor,
) (*EthereumListener, error) {
	// Parse contract address
	if !common.IsHexAddress(contractAddress) {
		return nil, fmt.Errorf("invalid contract address: %s", contractAddress)
	}

	// Parse contract ABI
	contractABI, err := abi.JSON(strings.NewReader(binding.SLAEnforcerMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	return &EthereumListener{
		contractAddress:   common.HexToAddress(contractAddress),
		contractABI:       contractABI,
		rpcURL:            rpcURL,
		reconnectInterval: time.Duration(reconnectInterval) * time.Second,
		startBlock:        startBlock,
		logger:            logger,
		processors:        processors,
		stopCh:            make(chan struct{}),
	}, nil
}

// Start begins listening for blockchain events
func (el *EthereumListener) Start(ctx context.Context) error {
	el.logger.Info("Starting Ethereum listener",
		zap.String("rpcURL", el.rpcURL),
		zap.String("contractAddress", el.contractAddress.Hex()),
		zap.Uint64("startBlock", el.startBlock),
	)

	// Initial connection
	if err := el.connect(); err != nil {
		return fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	// Start listening in a goroutine
	go el.listen(ctx)

	return nil
}

// Stop stops the blockchain listener
func (el *EthereumListener) Stop() error {
	el.logger.Info("Stopping Ethereum listener")
	close(el.stopCh)

	el.mu.Lock()
	defer el.mu.Unlock()

	if el.client != nil {
		el.client.Close()
		el.connected = false
	}

	return nil
}

// Subscribe registers an event processor to receive blockchain events
// This method is thread-safe and can be called at any time, even after the listener has started.
func (el *EthereumListener) Subscribe(processor driven.EventProcessor) error {
	if processor == nil {
		return fmt.Errorf("processor cannot be nil")
	}

	el.processorsMu.Lock()
	defer el.processorsMu.Unlock()

	el.processors = append(el.processors, processor)
	el.logger.Info("Event processor subscribed",
		zap.Int("totalProcessors", len(el.processors)),
	)

	return nil
}

// IsConnected returns true if the listener is connected to the blockchain
func (el *EthereumListener) IsConnected() bool {
	el.mu.RLock()
	defer el.mu.RUnlock()
	return el.connected
}

// notifyProcessors sends an event to all registered processors
func (el *EthereumListener) notifyProcessors(ctx context.Context, event interface{}) error {
	el.processorsMu.RLock()
	defer el.processorsMu.RUnlock()

	if len(el.processors) == 0 {
		el.logger.Warn("No processors registered to handle event")
		return nil
	}

	var errors []error
	for i, processor := range el.processors {
		if err := processor.Process(ctx, event); err != nil {
			el.logger.Error("Processor failed to handle event",
				zap.Error(err),
				zap.Int("processorIndex", i),
			)
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to process event: %d processor(s) failed", len(errors))
	}

	return nil
}

// connect establishes a connection to the Ethereum node
func (el *EthereumListener) connect() error {
	el.logger.Info("Connecting to Ethereum node", zap.String("rpcURL", el.rpcURL))

	client, err := ethclient.Dial(el.rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	// Test connection by getting chain ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	el.mu.Lock()
	el.client = client
	el.connected = true
	el.mu.Unlock()

	el.logger.Info("Connected to Ethereum node",
		zap.String("chainID", chainID.String()),
	)

	return nil
}

// listen to for blockchain events with auto-reconnect
func (el *EthereumListener) listen(ctx context.Context) {
	for {
		select {
		case <-el.stopCh:
			el.logger.Info("Listener stopped")
			return
		case <-ctx.Done():
			el.logger.Info("Context cancelled, stopping listener")
			return
		default:
			if !el.IsConnected() {
				el.logger.Warn("Not connected, attempting to reconnect",
					zap.Duration("interval", el.reconnectInterval),
				)
				if err := el.connect(); err != nil {
					el.logger.Error("Failed to reconnect", zap.Error(err))
					time.Sleep(el.reconnectInterval)
					continue
				}
			}

			if err := el.subscribeToEvents(ctx); err != nil {
				el.logger.Error("Error subscribing to events", zap.Error(err))
				el.mu.Lock()
				el.connected = false
				if el.client != nil {
					el.client.Close()
				}
				el.mu.Unlock()
				time.Sleep(el.reconnectInterval)
			}
		}
	}
}

// subscribeToEvents subscribes to contract events
func (el *EthereumListener) subscribeToEvents(ctx context.Context) error {
	// Check if WebSocket is available by trying to subscribe
	if strings.HasPrefix(el.rpcURL, "ws://") || strings.HasPrefix(el.rpcURL, "wss://") {
		el.logger.Info("Using WebSocket subscription mode")
		return el.subscribeWebSocket(ctx)
	}

	// Fall back to HTTP polling
	el.logger.Info("Using HTTP polling mode")
	return el.pollHTTP(ctx)
}

// subscribeWebSocket subscribes to events via WebSocket
func (el *EthereumListener) subscribeWebSocket(ctx context.Context) error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{el.contractAddress},
	}

	if el.startBlock > 0 {
		query.FromBlock = big.NewInt(int64(el.startBlock))
	}

	logs := make(chan types.Log)
	sub, err := el.client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %w", err)
	}
	defer sub.Unsubscribe()

	el.logger.Info("Successfully subscribed to events, waiting for logs...")

	for {
		select {
		case <-el.stopCh:
			return nil
		case <-ctx.Done():
			return nil
		case err := <-sub.Err():
			return fmt.Errorf("subscription error: %w", err)
		case vLog := <-logs:
			if err := el.processLog(ctx, vLog); err != nil {
				el.logger.Error("Failed to process log", zap.Error(err))
			}
		}
	}
}

// pollHTTP polls for events via HTTP
func (el *EthereumListener) pollHTTP(ctx context.Context) error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{el.contractAddress},
	}

	// Get the current block number
	currentBlock, err := el.client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}

	// Start from a configured block or current block
	fromBlock := el.startBlock
	if fromBlock == 0 {
		fromBlock = currentBlock
	}

	el.logger.Info("Starting HTTP polling for events",
		zap.Uint64("fromBlock", fromBlock),
		zap.Uint64("currentBlock", currentBlock),
	)

	ticker := time.NewTicker(5 * time.Second) // Poll every 5 seconds
	defer ticker.Stop()

	lastProcessedBlock := fromBlock

	for {
		select {
		case <-el.stopCh:
			return nil
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			// Get the latest block
			latestBlock, err := el.client.BlockNumber(ctx)
			if err != nil {
				el.logger.Error("Failed to get latest block number", zap.Error(err))
				continue
			}

			// If no new blocks, continue
			if latestBlock <= lastProcessedBlock {
				continue
			}

			// Query logs from last processed block to latest
			query.FromBlock = big.NewInt(int64(lastProcessedBlock + 1))
			query.ToBlock = big.NewInt(int64(latestBlock))

			logs, err := el.client.FilterLogs(ctx, query)
			if err != nil {
				el.logger.Error("Failed to filter logs", zap.Error(err))
				continue
			}

			// Process each log
			for _, vLog := range logs {
				if err := el.processLog(ctx, vLog); err != nil {
					el.logger.Error("Failed to process log", zap.Error(err))
				}
			}

			if len(logs) > 0 {
				el.logger.Info("Processed logs",
					zap.Int("count", len(logs)),
					zap.Uint64("fromBlock", lastProcessedBlock+1),
					zap.Uint64("toBlock", latestBlock),
				)
			}

			lastProcessedBlock = latestBlock
		}
	}
}

// processLog processes a single log entry
func (el *EthereumListener) processLog(ctx context.Context, vLog types.Log) error {
	eventSignature := vLog.Topics[0].Hex()

	el.logger.Debug("Processing log",
		zap.String("eventSignature", eventSignature),
		zap.Uint64("blockNumber", vLog.BlockNumber),
	)

	switch eventSignature {
	case el.contractABI.Events["ContractAdded"].ID.Hex():
		return el.handleContractAdded(ctx, vLog)
	case el.contractABI.Events["SLAAdded"].ID.Hex():
		return el.handleSLAAdded(ctx, vLog)
	case el.contractABI.Events["SLAStatusUpdated"].ID.Hex():
		return el.handleSLAStatusUpdated(ctx, vLog)
	default:
		el.logger.Warn("Unknown event signature", zap.String("signature", eventSignature))
	}

	return nil
}

// handleContractAdded handles ContractAdded events
func (el *EthereumListener) handleContractAdded(ctx context.Context, vLog types.Log) error {
	var event binding.SLAEnforcerContractAdded

	err := el.contractABI.UnpackIntoInterface(&event, "ContractAdded", vLog.Data)
	if err != nil {
		return fmt.Errorf("failed to unpack ContractAdded event: %w", err)
	}

	// Extract indexed contractId from topics
	contractID := string(vLog.Topics[1].Bytes())

	domainEvent := &domain.ContractAddedEvent{
		BlockchainEvent: domain.BlockchainEvent{
			EventType:   domain.EventTypeContractAdded,
			BlockNumber: vLog.BlockNumber,
			TxHash:      vLog.TxHash.Hex(),
			Timestamp:   time.Now(),
		},
		ContractID: contractID,
		CustomerID: event.CustomerId,
	}

	return el.notifyProcessors(ctx, domainEvent)
}

// handleSLAAdded handles SLAAdded events
func (el *EthereumListener) handleSLAAdded(ctx context.Context, vLog types.Log) error {
	var event binding.SLAEnforcerSLAAdded

	err := el.contractABI.UnpackIntoInterface(&event, "SLAAdded", vLog.Data)
	if err != nil {
		return fmt.Errorf("failed to unpack SLAAdded event: %w", err)
	}

	// Extract indexed contractId from topics
	contractID := string(vLog.Topics[1].Bytes())

	domainEvent := &domain.SLAAddedEvent{
		BlockchainEvent: domain.BlockchainEvent{
			EventType:   domain.EventTypeSLAAdded,
			BlockNumber: vLog.BlockNumber,
			TxHash:      vLog.TxHash.Hex(),
			Timestamp:   time.Now(),
		},
		ContractID: contractID,
		SLAID:      event.SlaId,
	}

	return el.notifyProcessors(ctx, domainEvent)
}

// handleSLAStatusUpdated handles SLAStatusUpdated events
func (el *EthereumListener) handleSLAStatusUpdated(ctx context.Context, vLog types.Log) error {
	var event binding.SLAEnforcerSLAStatusUpdated

	err := el.contractABI.UnpackIntoInterface(&event, "SLAStatusUpdated", vLog.Data)
	if err != nil {
		return fmt.Errorf("failed to unpack SLAStatusUpdated event: %w", err)
	}

	// Extract indexed contractId and slaId from topics
	contractID := string(vLog.Topics[1].Bytes())
	slaID := string(vLog.Topics[2].Bytes())

	domainEvent := &domain.SLAStatusUpdatedEvent{
		BlockchainEvent: domain.BlockchainEvent{
			EventType:   domain.EventTypeSLAStatusUpdated,
			BlockNumber: vLog.BlockNumber,
			TxHash:      vLog.TxHash.Hex(),
			Timestamp:   time.Now(),
		},
		ContractID: contractID,
		SLAID:      slaID,
		NewStatus:  event.NewStatus,
	}

	return el.notifyProcessors(ctx, domainEvent)
}
