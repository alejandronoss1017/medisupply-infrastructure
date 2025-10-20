# Build stage
FROM golang:1.24 AS builder

# Install librdkafka development libraries
RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Verify module and tidy dependencies
RUN go mod tidy

# Build the application with dynamic linking to librdkafka
RUN go build -tags dynamic -o main ./cmd/consumer

# Runtime stage
FROM golang:1.24-alpine AS runtime

# Install librdkafka runtime libraries
RUN apk add --no-cache librdkafka

# Create a non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["./main"]
