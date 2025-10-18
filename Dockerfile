# Multi-stage Dockerfile for tmdb-mcp
# Stage 1: Build the Go binary
FROM golang:1.25.2-alpine AS builder

WORKDIR /build

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary (CGO_ENABLED=0 for fully static binary)
RUN CGO_ENABLED=0 GOOS=linux go build -o tmdb-mcp ./cmd/tmdb-mcp

# Stage 2: Runtime image
FROM alpine:latest

# Install CA certificates (required for HTTPS requests to TMDB API)
RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/tmdb-mcp .

# Create config directory
RUN mkdir -p /root/.tmdb-mcp

# Expose SSE server port
EXPOSE 8910

# Health check using wget
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8910/health || exit 1

# Run the application
CMD ["./tmdb-mcp"]
