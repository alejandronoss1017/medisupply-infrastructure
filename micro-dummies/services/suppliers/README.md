# Medicine Suppliers Service

A microservice for managing medicine suppliers built with Go, following Hexagonal Architecture (Ports and Adapters pattern). This service provides a RESTful API for medicine operations and publishes domain events to Kafka.

## Architecture

This service implements **Hexagonal Architecture** with clear separation of concerns:

```
├── cmd/web/              # Application entry point
├── internal/
│   ├── core/            # Business logic (hexagon core)
│   │   ├── domain/      # Domain entities and events
│   │   ├── application/ # Use cases/services
│   │   └── port/        # Interfaces (ports)
│   │       ├── driver/  # Inbound ports (driving side)
│   │       └── driven/  # Outbound ports (driven side)
│   └── adapter/         # External adapters
│       ├── http/        # REST API handlers (driver adapter)
│       └── queue/       # Kafka publisher (driven adapter)
└── pkg/
    └── logger/          # Shared logging package
```

### Core Concepts

- **Domain**: Medicine entity with business rules
- **Ports**: Interfaces defining contracts (MedicineService, EventPublisher)
- **Adapters**:
  - HTTP handlers (Gin) for REST API
  - Kafka publisher for event streaming
- **Events**: Domain events (created, updated, deleted) published to Kafka

## Features

- RESTful API for medicine management (CRUD operations)
- Event-driven architecture with Kafka integration
- Hexagonal architecture for maintainability
- Structured logging with context
- Dockerized deployment
- Non-root container for security

## Tech Stack

- **Language**: Go 1.24
- **Web Framework**: Gin
- **Message Broker**: Apache Kafka (confluent-kafka-go)
- **Containerization**: Docker (multi-stage build)

## API Endpoints

### Health Check
- `GET /ping` - Health check endpoint

### Medicine Operations
- `GET /medicines` - Get all medicines
- `GET /medicines/:id` - Get medicine by ID
- `POST /medicines` - Create new medicine
- `PUT /medicines/:id` - Update medicine
- `DELETE /medicines/:id` - Delete medicine

### Medicine Schema

```json
{
  "id": "string (UUID)",
  "name": "string (required)",
  "description": "string (required)",
  "price": "number (required)",
  "strength": "string (required)",
  "category": "string (required)",
  "supplier_id": "string (required)",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `KAFKA_HOST` | Kafka bootstrap servers | - |
| `KAFKA_TOPIC` | Kafka topic for events | - |

## Getting Started

### Prerequisites

- Go 1.24+
- Docker (for containerized deployment)
- Kafka instance (for event publishing)

### Local Development

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Set environment variables**
   ```bash
   export KAFKA_HOST=localhost:9092
   export KAFKA_TOPIC=medicine-events
   ```

3. **Run the service**
   ```bash
   go run ./cmd/web
   ```

The service will start on port 8080.

### Docker Deployment

#### Build the image

```bash
docker build -t medisupply/procurement-supply-optimization/suppliers/web .
```

#### Run the container

```bash
docker run -d \
  -p 8080:8080 \
  -e KAFKA_HOST=kafka:9092 \
  -e KAFKA_TOPIC=medicine-events \
  --name suppliers-service \
  medisupply/procurement-supply-optimization/suppliers/web
```

### Docker Build Notes

The Dockerfile uses a multi-stage build:
- **Builder stage**: Compiles the Go application with dynamic linking to librdkafka
- **Runtime stage**: Minimal Alpine image with only runtime dependencies

**Important**: The service uses dynamic linking for Kafka (`-tags dynamic`) to ensure cross-platform compatibility, especially for ARM64 architectures (Apple Silicon).

## Event Publishing

The service publishes domain events to Kafka for the following operations:

### Event Types

- `medicine.created` - Published when a new medicine is created
- `medicine.updated` - Published when a medicine is updated
- `medicine.deleted` - Published when a medicine is deleted

### Event Structure

```json
{
  "event_type": "medicine.created",
  "data": {
    "id": "uuid",
    "name": "Aspirin",
    "description": "Pain reliever",
    "price": 9.99,
    "strength": "500mg",
    "category": "Analgesic",
    "supplier_id": "supplier-uuid",
    "created_at": "2025-10-19T10:00:00Z",
    "updated_at": "2025-10-19T10:00:00Z"
  },
  "timestamp": "2025-10-19T10:00:00Z"
}
```

## Logging

The service uses a custom structured logger with contextual prefixes:

- `[APP]` - Application-level logs
- `[HTTP]` - HTTP request/response logs
- `[KAFKA]` - Kafka event publishing logs

Log levels: `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`

## Project Structure Details

### Core Layer (`internal/core`)

- **domain/medicine.go**: Medicine entity with validation tags
- **domain/events.go**: Event types and factory methods
- **application/medicine.go**: Business logic and use cases
- **port/driver/service.go**: Inbound port interface
- **port/driven/publisher.go**: Outbound port interface

### Adapter Layer (`internal/adapter`)

- **http/medicine.go**: REST API handlers using Gin
- **http/pong.go**: Health check handler
- **queue/kafka.go**: Kafka event publisher implementation

## Development Guidelines

### Adding New Features

1. **Define domain entities** in `internal/core/domain`
2. **Create ports** (interfaces) in `internal/core/port`
3. **Implement business logic** in `internal/core/application`
4. **Create adapters** in `internal/adapter`
5. **Wire dependencies** in `cmd/web/main.go`

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Troubleshooting

### Docker Build Issues

**Issue**: `package suppliers/pkg/logger is not in std`
- **Solution**: Ensure `pkg/` is not in `.dockerignore`

**Issue**: Architecture mismatch with librdkafka
- **Solution**: Build with `-tags dynamic` flag and install librdkafka libraries

### Kafka Connection Issues

- Verify `KAFKA_HOST` environment variable
- Check network connectivity to Kafka broker
- Review Kafka logs for broker availability

## License

This project is part of the MediSupply infrastructure.

## Contributing

1. Follow hexagonal architecture principles
2. Keep domain logic pure (no external dependencies)
3. Use dependency injection via ports
4. Write tests for business logic
5. Document public APIs
