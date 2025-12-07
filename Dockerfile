# Build stage
FROM golang:1.21-bullseye AS builder

# Install only essential build dependencies for CGO
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libc6-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download
RUN go mod tidy

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o sim-board .

# Runtime stage - use minimal Debian slim
FROM debian:bookworm-slim

# Install only runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libsqlite3-0 \
    curl \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# Create non-root user for security
RUN groupadd -r appuser && useradd -r -g appuser -u 1000 appuser

# Create app directory and data directory
WORKDIR /app
RUN mkdir -p /app/data && chown -R appuser:appuser /app

# Copy the binary from builder
COPY --from=builder /app/sim-board .

# Set ownership
RUN chown appuser:appuser /app/sim-board

# Expose port (default, can be overridden via PORT env var)
EXPOSE 8869

# Set data directory environment variable
ENV DATA_DIR=/app/data

# Switch to non-root user
USER appuser

# Health check (uses default port 8869, or PORT env var if set)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD sh -c 'curl -f http://localhost:${PORT:-8869}/health || exit 1'

# Run the application
CMD ["./sim-board"]

