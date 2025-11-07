# =============================================================================
# Multi-Stage Dockerfile for Contracts Web API Service
# =============================================================================
# This Dockerfile builds the Web API service that provides REST endpoints for
# contract management with Ethereum smart contract integration.
#
# Services: REST API, Blockchain Writer/Reader
# Port: 8080
# Dependencies: Ethereum RPC endpoint, Smart Contract Address, Private Key
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Builder
# -----------------------------------------------------------------------------
FROM golang:1.24 AS builder

# Set working directory
WORKDIR /app

# Copy dependency manifests first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
# This layer will be cached unless go.mod or go.sum changes
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the web application
# - CGO_ENABLED=0: Static binary (no C dependencies needed)
# - GOOS=linux: Target Linux OS
# - -ldflags="-s -w": Strip debug information for smaller binary
# - -o web: Output binary name
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o web \
    ./cmd/web

# -----------------------------------------------------------------------------
# Stage 2: Runtime
# -----------------------------------------------------------------------------
FROM golang:1.24-alpine AS runtime

# Install runtime dependencies
# - ca-certificates: Required for HTTPS connections to Ethereum RPC endpoints
# - tzdata: Timezone data for proper logging timestamps
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/web .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose HTTP port
EXPOSE 8080

# Environment variables documentation
# Required:
#   - RCP_URL: Ethereum RPC endpoint URL (e.g., http://localhost:8545)
#   - SMART_CONTRACT_ADDRESS: Hex-encoded smart contract address (e.g., 0x...)
#   - PRIVATE_KEY: Hex-encoded ECDSA private key for transaction signing
#
# Optional (Logging):
#   - LOG_LEVEL: debug|info|warn|error|fatal (default: info)
#   - LOG_ENV: development|production (default: production)
#   - LOG_ENCODING: json|console (default: json in production, console in dev)

# Health check
# Checks if the web server is responding on the /ping endpoint
#HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ping || exit 1

# Run the web application
ENTRYPOINT ["./web"]
