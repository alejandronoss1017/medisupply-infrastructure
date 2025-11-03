# Multi-stage Dockerfile for Blockchain Notification Service
# Stage 1: Builder - Compile the Go application
# Stage 2: Runtime - Minimal image with non-root user

# ==============================================================================
# Stage 1: Builder
# ==============================================================================
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy dependency files first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
# - CGO_ENABLED=0: Build a statically-linked binary (no external dependencies)
# - GOOS=linux: Target Linux
# - GOARCH=amd64: Target 64-bit architecture
# - -ldflags: Linker flags to strip debug information and reduce binary size
#   -s: Omit symbol table and debug information
#   -w: Omit DWARF symbol table
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-s -w -extldflags "-static"' \
    -a \
    -installsuffix cgo \
    -o listener \
    ./cmd/listener/main.go

# Verify the binary was created and is executable
RUN chmod +x /build/listener && \
    /build/listener --version || echo "Binary built successfully"

# ==============================================================================
# Stage 2: Runtime
# ==============================================================================
FROM alpine:3.19

# Install runtime dependencies
# - ca-certificates: Required for HTTPS connections (blockchain RPC, AWS SNS)
# - tzdata: Timezone data for proper timestamp handling
RUN apk --no-cache add ca-certificates tzdata && \
    update-ca-certificates

# Create a non-root user and group
# - UID/GID 1001 is commonly used for application users
# - No home directory needed for service
# - No shell access for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup -D -H -s /sbin/nologin

# Set working directory
WORKDIR /app

# Copy the compiled binary from builder stage
COPY --from=builder --chown=appuser:appgroup /build/listener /app/listener

# Copy any additional required files (if needed)
# COPY --from=builder --chown=appuser:appgroup /build/abi.json /app/abi.json

# Set proper permissions
RUN chmod 755 /app/listener

# Switch to non-root user
USER appuser

# Set environment variables for Go runtime
ENV GODEBUG=netdns=go

# Health check (optional - adjust based on your service capabilities)
# HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
#     CMD ["/app/listener", "--health"] || exit 1

# Document that this service doesn't expose ports (it's an event listener)
# If you add a health check endpoint, expose that port:
# EXPOSE 8080

# Set the binary as the entrypoint
ENTRYPOINT ["/app/listener"]

# Metadata
LABEL maintainer="your-team@example.com" \
      description="Blockchain Event Listener for SLA Monitoring" \
      version="1.0.0" \
      org.opencontainers.image.source="https://github.com/your-org/notifications"
