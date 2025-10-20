# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies including musl-dev for static compilation
RUN apk add --no-cache \
    gcc \
    musl-dev \
    librdkafka-dev \
    pkgconfig

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

# Build the application with dynamic linking (no vendored librdkafka)
RUN CGO_ENABLED=1 go build -tags dynamic -o main ./cmd/web

# Runtime stage
FROM golang:1.24-alpine AS runtime

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    librdkafka

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
