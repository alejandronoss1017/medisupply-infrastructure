# =============================================================================
# Multi-Stage Dockerfile for Contracts Kafka Consumer Service
# =============================================================================
# This Dockerfile builds the Kafka Consumer service that processes medicine
# events from Kafka topics.
#
# Services: Kafka Event Consumer
# Port: None (not an HTTP service)
# Dependencies: Kafka cluster, librdkafka
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Builder
# -----------------------------------------------------------------------------
FROM golang:1.24 AS builder

# Install librdkafka development libraries
# Required for building confluent-kafka-go client
RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy dependency manifests first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
# This layer will be cached unless go.mod or go.sum changes
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the consumer application
# - tags dynamic: Enable dynamic linking to librdkafka (required)
# - -o consumer: Output binary name
# - ./cmd/consumer: Consumer entry point
RUN go build -tags dynamic -o consumer ./cmd/consumer

# -----------------------------------------------------------------------------
# Stage 2: Runtime
# -----------------------------------------------------------------------------
FROM golang:1.24-alpine AS runtime

# Install librdkafka runtime libraries
# Required for running the Kafka consumer at runtime
RUN apk add --no-cache librdkafka

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/consumer .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Environment variables documentation
# Optional (with defaults):
#   - KAFKA_HOST: Kafka bootstrap servers (default: localhost:9092)
#   - KAFKA_GROUP_ID: Consumer group ID (default: contracts-service)
#   - KAFKA_TOPICS: Topics to subscribe (default: medicine-events)
#
# Optional (Logging):
#   - LOG_LEVEL: debug|info|warn|error|fatal (default: info)
#   - LOG_ENV: development|production (default: production)
#   - LOG_ENCODING: json|console (default: json in production, console in dev)

# Run the consumer application
ENTRYPOINT ["./consumer"]
