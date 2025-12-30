# API Service (agentd) Dockerfile
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Install swag for Swagger documentation generation
# Get GOPATH and add to PATH for swag command
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    GOPATH=$(go env GOPATH) && \
    echo "export PATH=\$PATH:\$GOPATH/bin" >> /root/.profile

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger documentation
# Use full path to swag or source profile to get PATH
RUN GOPATH=$(go env GOPATH) && \
    $GOPATH/bin/swag init -g cmd/agentd/main.go -o docs --parseDependency --parseInternal

# Build the application
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /build/bin/agentd \
    ./cmd/agentd

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget

# Create app user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/bin/agentd /app/agentd

# Copy config files
COPY configs/ /app/configs/

# Create data directory for artifacts
RUN mkdir -p /app/data/workdir && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3333

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3333/v1/healthz || exit 1

# Run the application
CMD ["/app/agentd"]

