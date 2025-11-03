# Blockchain Notification Service

A blockchain event listener service that monitors Ethereum-compatible smart contracts for events related to contracts and Service Level Agreements (SLAs). The service processes blockchain events in real-time and publishes notifications to AWS SNS.

## Overview

This service monitors the **SLAEnforcer** smart contract for three types of events:
- **ContractAdded**: Triggered when a new contract is added to the blockchain
- **SLAAdded**: Triggered when an SLA is added to a contract
- **SLAStatusUpdated**: Triggered when an SLA status changes (Pending → Met/Violated)

## Architecture

The project follows **Hexagonal Architecture** (Ports and Adapters pattern) for clean separation of concerns:

```
┌─────────────────────────────────────────────────────────┐
│                    Inbound Adapter                       │
│            (EthereumListener - Blockchain)               │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│                  Application Layer                       │
│         (Event Processors - Business Logic)              │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│                   Outbound Adapter                       │
│              (SNSNotifier - Notifications)               │
└─────────────────────────────────────────────────────────┘
```

### Key Components

- **Domain Layer** (`internal/core/domain/`): Core business entities and events
- **Application Layer** (`internal/core/application/`): Event processors with business logic
- **Ports** (`internal/core/port/`): Interfaces defining contracts between layers
- **Adapters** (`internal/adapter/`):
  - Blockchain adapter for Ethereum event listening
  - SNS adapter for AWS notifications

## Features

- Real-time blockchain event monitoring via WebSocket
- HTTP polling mode fallback for non-WebSocket endpoints
- Automatic reconnection with configurable intervals
- AWS SNS integration for event notifications
- Structured logging with configurable levels (Zap logger)
- SLA violation detection and alerting
- Graceful shutdown handling

## Prerequisites

- Go 1.21 or higher
- Access to an Ethereum-compatible RPC endpoint (WebSocket or HTTP)
- AWS account with SNS configured (optional, for notifications)
- Smart contract ABI and deployed contract address

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd notifications
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables (see Configuration section below)

## Configuration

Configuration is managed via environment variables. Set the following variables in your shell environment, Docker container, or deployment platform:

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `BLOCKCHAIN_RPC_URL` | Ethereum RPC endpoint (WebSocket or HTTP) | `wss://mainnet.infura.io/ws/v3/YOUR_KEY` |
| `CONTRACT_ADDRESS` | Smart contract address to monitor | `0x1234...abcd` |

### Optional Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `START_BLOCK` | `0` (current) | Block number to start listening from |
| `RECONNECT_INTERVAL` | `5` | Seconds between reconnection attempts |
| `SNS_ENABLED` | `false` | Enable AWS SNS notifications |
| `SNS_TOPIC_ARN` | - | SNS topic ARN (required if SNS enabled) |
| `AWS_REGION` | `us-east-1` | AWS region for SNS |
| `LOG_LEVEL` | `info` | Logging level: debug/info/warn/error |
| `LOG_FORMAT` | `console` | Log format: json/console |

### Setting Environment Variables

**Linux/macOS (Shell) - Minimal Configuration:**
```bash
# Required variables
export BLOCKCHAIN_RPC_URL="wss://mainnet.infura.io/ws/v3/YOUR_KEY"
export CONTRACT_ADDRESS="0x1234567890123456789012345678901234567890"
```

**Linux/macOS (Shell) - Full Configuration:**
```bash
# Required variables
export BLOCKCHAIN_RPC_URL="wss://mainnet.infura.io/ws/v3/YOUR_KEY"
export CONTRACT_ADDRESS="0x1234567890123456789012345678901234567890"

# Optional blockchain configuration
export START_BLOCK="0"
export RECONNECT_INTERVAL="5"

# Optional AWS SNS configuration
export SNS_ENABLED="false"
export SNS_TOPIC_ARN="arn:aws:sns:us-east-1:123456789012:blockchain-events"
export AWS_REGION="us-east-1"

# Optional logging configuration
export LOG_LEVEL="info"
export LOG_FORMAT="console"
```

**Windows (PowerShell):**
```powershell
# Required variables
$env:BLOCKCHAIN_RPC_URL="wss://mainnet.infura.io/ws/v3/YOUR_KEY"
$env:CONTRACT_ADDRESS="0x1234567890123456789012345678901234567890"

# Optional variables
$env:START_BLOCK="0"
$env:RECONNECT_INTERVAL="5"
$env:SNS_ENABLED="false"
$env:SNS_TOPIC_ARN="arn:aws:sns:us-east-1:123456789012:blockchain-events"
$env:AWS_REGION="us-east-1"
$env:LOG_LEVEL="info"
$env:LOG_FORMAT="console"
```

**Docker:**
```bash
docker run \
  -e BLOCKCHAIN_RPC_URL="wss://mainnet.infura.io/ws/v3/YOUR_KEY" \
  -e CONTRACT_ADDRESS="0x1234567890123456789012345678901234567890" \
  -e START_BLOCK="0" \
  -e RECONNECT_INTERVAL="5" \
  -e SNS_ENABLED="false" \
  -e SNS_TOPIC_ARN="arn:aws:sns:us-east-1:123456789012:blockchain-events" \
  -e AWS_REGION="us-east-1" \
  -e LOG_LEVEL="info" \
  -e LOG_FORMAT="console" \
  your-image
```

**Docker Compose:**
```yaml
services:
  listener:
    image: your-image
    environment:
      # Required
      - BLOCKCHAIN_RPC_URL=wss://mainnet.infura.io/ws/v3/YOUR_KEY
      - CONTRACT_ADDRESS=0x1234567890123456789012345678901234567890
      # Optional blockchain
      - START_BLOCK=0
      - RECONNECT_INTERVAL=5
      # Optional AWS SNS
      - SNS_ENABLED=false
      - SNS_TOPIC_ARN=arn:aws:sns:us-east-1:123456789012:blockchain-events
      - AWS_REGION=us-east-1
      # Optional logging
      - LOG_LEVEL=info
      - LOG_FORMAT=console
```

**Kubernetes:**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: blockchain-listener
spec:
  containers:
  - name: listener
    image: your-image
    env:
      # Required
      - name: BLOCKCHAIN_RPC_URL
        value: "wss://mainnet.infura.io/ws/v3/YOUR_KEY"
      - name: CONTRACT_ADDRESS
        value: "0x1234567890123456789012345678901234567890"
      # Optional blockchain
      - name: START_BLOCK
        value: "0"
      - name: RECONNECT_INTERVAL
        value: "5"
      # Optional AWS SNS
      - name: SNS_ENABLED
        value: "false"
      - name: SNS_TOPIC_ARN
        value: "arn:aws:sns:us-east-1:123456789012:blockchain-events"
      - name: AWS_REGION
        value: "us-east-1"
      # Optional logging
      - name: LOG_LEVEL
        value: "info"
      - name: LOG_FORMAT
        value: "console"
```

> **Note**: AWS Credentials (automatically loaded from environment or ~/.aws/credentials), you can also set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables or use IAM roles if running on EC2/ECS

## Usage

### Running the Service

```bash
# Build the service
go build -o listener ./cmd/listener

# Run the service
./listener
```

Or run directly:
```bash
go run ./cmd/listener/main.go
```

### Graceful Shutdown

Press `Ctrl+C` to gracefully shut down the service. The listener will:
1. Stop accepting new events
2. Complete processing of in-flight events
3. Close connections cleanly

## Development

### Project Structure

```
notifications/
├── cmd/
│   └── listener/           # Application entry point
│       └── main.go
├── internal/
│   ├── core/
│   │   ├── domain/         # Domain entities and events
│   │   ├── application/    # Business logic (event processors)
│   │   └── port/           # Interface definitions
│   │       ├── driven/     # Outbound ports (Notifier)
│   │       └── driver/     # Inbound ports (BlockchainListener)
│   └── adapter/
│       ├── blockchain/     # Ethereum listener adapter
│       │   └── binding/    # Auto-generated contract bindings
│       └── notification/   # SNS notifier adapter
├── config/                 # Configuration loading
├── pkg/
│   └── logger/            # Logging setup
├── abi.json               # Smart contract ABI
├── .env.example           # Example environment configuration
└── README.md
```

### Building

```bash
# Build the service
go build -o listener ./cmd/listener

# Build all packages
go build ./...
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Test specific package
go test ./internal/core/application/...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Lint (requires golangci-lint)
golangci-lint run

# Verify dependencies
go mod verify
go mod tidy
```

## Event Processing

### Supported Events

#### 1. ContractAdded
Emitted when a new contract is registered on the blockchain.

**Notification Payload:**
```json
{
  "contractId": "contract-123",
  "customerId": "customer-456",
  "eventType": "ContractAdded"
}
```

#### 2. SLAAdded
Emitted when a new SLA is added to a contract.

**Notification Payload:**
```json
{
  "contractId": "contract-123",
  "slaId": "sla-789",
  "eventType": "SLAAdded"
}
```

#### 3. SLAStatusUpdated
Emitted when an SLA status changes (Pending → Met/Violated).

**Notification Payload:**
```json
{
  "contractId": "contract-123",
  "slaId": "sla-789",
  "status": 2,
  "statusName": "Violated",
  "eventType": "SLAStatusUpdated"
}
```

**SLA Status Values:**
- `0` - Pending
- `1` - Met
- `2` - Violated

### Event Processing Flow

1. **EthereumListener** subscribes to blockchain events (WebSocket) or polls (HTTP)
2. Raw blockchain logs are parsed using auto-generated contract bindings
3. Events are converted to domain-specific event objects
4. Registered **EventProcessors** receive and process events:
   - Execute business logic
   - Send notifications via **Notifier**
5. **SNSNotifier** publishes notifications to AWS SNS

### Error Handling

- If notification sending fails, the entire event processing fails
- Events that fail processing will not be marked as processed
- The service will retry on reconnection if connection is lost
- All errors are logged with structured context

## Adding New Event Types

To add support for new blockchain events:

1. **Update the ABI**: Add event definition to `abi.json`
2. **Regenerate bindings**:
   ```bash
   abigen --abi=abi.json --pkg=binding --out=internal/adapter/blockchain/binding/sla_enforcer.go
   ```
3. **Define domain event**: Add to `internal/core/domain/event.go`
4. **Update EthereumListener**: Add event handler in `ethereum_listener.go`
5. **Create processor**: Implement processor in `internal/core/application/`
6. **Wire up**: Register processor in `cmd/listener/main.go`

## AWS SNS Integration

### Message Attributes

All SNS messages include the following attributes:
- `eventType`: String attribute identifying the event type
- `Subject`: "Blockchain Event: {EventType}"

### Message Format

Messages are published as JSON with core business fields (no blockchain metadata like txHash or blockNumber).

### IAM Permissions

The service requires the following AWS IAM permissions:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "sns:Publish"
      ],
      "Resource": "arn:aws:sns:REGION:ACCOUNT_ID:TOPIC_NAME"
    }
  ]
}
```

## Troubleshooting

### Connection Issues

**Problem**: Service can't connect to Ethereum node
- Verify `BLOCKCHAIN_RPC_URL` is correct
- Check network connectivity
- Ensure WebSocket endpoint supports subscriptions (or use HTTP for polling)

**Problem**: Frequent disconnections
- Increase `RECONNECT_INTERVAL`
- Check node provider rate limits
- Consider using a dedicated node

### Notification Failures

**Problem**: SNS notifications not being sent
- Verify `SNS_ENABLED=true` in configuration
- Check `SNS_TOPIC_ARN` is correct
- Verify AWS credentials are configured (via environment or IAM role)
- Check IAM permissions for SNS publish

### Performance

**Problem**: Event processing is slow
- Check network latency to Ethereum node
- Verify SNS endpoint response times
- Review log level (set to `info` or `warn` in production)

## Contributing

1. Follow Go best practices and conventions
2. Maintain hexagonal architecture principles
3. Write tests for new functionality
4. Update documentation for significant changes
5. Use structured logging with appropriate context

## License

[Add your license information here]
