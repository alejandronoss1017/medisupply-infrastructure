# Contracts Service

A microservice for managing contracts, SLAs (Service Level Agreements), and customers with Ethereum blockchain integration. Built following Hexagonal Architecture principles with Go generics.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Environment Variables](#environment-variables)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [Development](#development)
- [Project Structure](#project-structure)

## Overview

The Contracts Service provides:
- **Contract Management**: CRUD operations for contracts with blockchain registration
- **SLA Management**: Service level agreement tracking and management
- **Customer Management**: Customer information and relationship management
- **Blockchain Integration**: Ethereum smart contract interaction for contract immutability
- **Event Processing**: Kafka consumer for medicine-related events

## Architecture

This service follows **Hexagonal Architecture** (Ports & Adapters) with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                        Adapters (Infrastructure)             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   HTTP/REST  │  │    Kafka     │  │   Ethereum   │      │
│  │   Handlers   │  │   Consumer   │  │    Client    │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                  │                  │              │
│         ▼                  ▼                  ▼              │
├─────────────────────────────────────────────────────────────┤
│                    Ports (Interfaces)                        │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Driver Ports           │  Driven Ports              │   │
│  │  - Service Interfaces   │  - Repository[ID, T]       │   │
│  │  - Consumer Interface   │  - Event Handler           │   │
│  └─────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│                 Application (Business Logic)                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Contract    │  │     SLA      │  │   Customer   │      │
│  │   Service    │  │   Service    │  │    Service   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
├─────────────────────────────────────────────────────────────┤
│                    Domain (Core)                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Contract   │  │     SLA      │  │   Customer   │      │
│  │   Medicine   │  │    Events    │  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

### Key Design Patterns

- **Generic Repository Pattern**: Single `Repository[ID, T]` interface for all entity types
- **Service Interfaces**: All services implement driver port interfaces
- **Dependency Injection**: All dependencies injected via constructors
- **Interface Segregation**: Small, focused interfaces

## Features

### Contract Management
- Create contracts with automatic blockchain registration
- Retrieve contracts by ID or list all
- Update existing contracts
- Delete contracts
- Automatic rollback on blockchain failures

### SLA Management
- Create and manage service level agreements
- Track SLA status and compliance
- Support for multiple comparator types (greater than, less than, equal, etc.)
- Auto-generated IDs

### Customer Management
- Full CRUD operations for customer data
- Customer information tracking
- Relationship with contracts

### Blockchain Integration
- **Read Operations**: Query contract data from blockchain (no gas cost)
- **Write Operations**: Register contracts on blockchain with transaction support
- **EIP-1559 Support**: Modern Ethereum transaction format
- **Legacy Support**: Fallback for older networks

### Event-Driven Architecture
- Kafka consumer for medicine events
- Event routing and processing
- Support for `medicine.updated` and `medicine.deleted` events

## Prerequisites

- **Go**: 1.24.0 or higher
- **Ethereum Node**: Access to an Ethereum-compatible blockchain (Ganache, Hardhat, etc.)
- **Kafka**: For event processing (optional)
- **ABI File**: Smart contract ABI JSON file

## Environment Variables

### Required (for Web Server)

```bash
# Ethereum Configuration
RCP_URL=http://localhost:8545                    # Ethereum RPC endpoint
SMART_CONTRACT_ADDRESS=0x...                     # Contract address
PRIVATE_KEY=...                                  # Private key (without 0x prefix)
ABI_PATH=./abi.json                             # Path to ABI file
```

### Optional (for Kafka Consumer)

```bash
# Kafka Configuration
KAFKA_HOST=localhost:9092                        # Kafka bootstrap servers
KAFKA_GROUP_ID=contracts-service                 # Consumer group ID
KAFKA_TOPICS=medicine-events                     # Topics to subscribe
```

## Getting Started

### 1. Clone and Navigate

```bash
cd medisupply-infrastructure/procurement-supply/contracts
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Set Environment Variables

Create a `.env` file or export variables:

```bash
export RCP_URL="http://localhost:8545"
export SMART_CONTRACT_ADDRESS="0x5FbDB2315678afecb367f032d93F642f64180aa3"
export PRIVATE_KEY="ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
export ABI_PATH="./abi.json"
```

### 4. Run the Web Server

```bash
# Development
go run ./cmd/web/main.go

# Production build
go build -o web.exe ./cmd/web
./web.exe
```

Server starts on port **8080** by default.

### 5. Run the Event Consumer (Optional)

```bash
go run ./cmd/consumer/main.go
```

## API Endpoints

### Health Check

```
GET /ping
```

Returns server status.

### Contracts

```
GET    /contracts           # List all contracts
GET    /contracts/:id       # Get contract by ID
POST   /contracts           # Create new contract
PUT    /contracts/:id       # Update contract
DELETE /contracts/:id       # Delete contract
```

**Example Request:**
```bash
curl -X POST http://localhost:8080/contracts \
  -H "Content-Type: application/json" \
  -d '{
    "id": "contract-001",
    "path": "/contracts/medical-supplies",
    "customerId": "customer-123",
    "slas": [
      {
        "id": "sla-001",
        "name": "Delivery Time",
        "description": "Maximum delivery time",
        "target": "48",
        "comparator": 1,
        "status": 0
      }
    ]
  }'
```

### SLAs

```
GET    /slas                # List all SLAs
GET    /slas/:id            # Get SLA by ID
POST   /slas                # Create new SLA
PUT    /slas/:id            # Update SLA
DELETE /slas/:id            # Delete SLA
```

### Customers

```
GET    /customers           # List all customers
GET    /customers/:id       # Get customer by ID
POST   /customers           # Create new customer
PUT    /customers/:id       # Update customer
DELETE /customers/:id       # Delete customer
```

**Example Request:**
```bash
curl -X POST http://localhost:8080/customers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hospital San Rafael",
    "email": "procurement@sanrafael.com",
    "phone": "+1-555-0123"
  }'
```

## Development

### Building

```bash
# Web server
go build -o web.exe ./cmd/web

# Event consumer
go build -o consumer.exe ./cmd/consumer
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/core/application/...
```

### Code Structure

```bash
go fmt ./...        # Format code
go vet ./...        # Vet code
go mod tidy         # Clean dependencies
```

### Adding a New Repository Implementation

To add a database repository (e.g., PostgreSQL):

```go
// internal/adapter/storage/postgres/contract_repository.go
package postgres

import (
    "contracts/internal/core/domain"
    "contracts/internal/core/port/driven"
    "database/sql"
)

type ContractRepository struct {
    db *sql.DB
}

var _ driven.Repository[string, domain.Contract] = (*ContractRepository)(nil)

func NewContractRepository(db *sql.DB) *ContractRepository {
    return &ContractRepository{db: db}
}

func (r *ContractRepository) Create(contract domain.Contract) error {
    // Implement SQL INSERT
}

// Implement other methods...
```

Then update `main.go`:

```go
// Old
contractRepo := memory.NewContractRepository()

// New
db, _ := sql.Open("postgres", connectionString)
contractRepo := postgres.NewContractRepository(db)
```

## Project Structure

```
contracts/
├── abi.json                           # Smart contract ABI
├── cmd/
│   ├── consumer/
│   │   └── main.go                   # Kafka consumer entry point
│   └── web/
│       └── main.go                   # Web server entry point
├── internal/
│   ├── adapter/                      # Infrastructure adapters
│   │   ├── ethereum/                 # Blockchain client
│   │   ├── http/                     # HTTP handlers
│   │   ├── queue/                    # Kafka consumer
│   │   └── storage/
│   │       └── memory/               # In-memory repositories
│   └── core/
│       ├── application/              # Business logic services
│       ├── domain/                   # Domain entities
│       └── port/
│           ├── driven/               # Output port interfaces
│           └── driver/               # Input port interfaces
├── pkg/
│   └── logger/                       # Custom logger
├── go.mod
├── go.sum
├── CLAUDE.md                         # Architecture documentation
└── README.md
```

## Technologies

- **Language**: Go 1.24.0
- **Web Framework**: Gin
- **Blockchain**: go-ethereum
- **Messaging**: Confluent Kafka Go
- **Architecture**: Hexagonal Architecture
- **Patterns**: Repository Pattern (with Generics), Dependency Injection

## License

This project is for educational purposes.


