# ============================================================
# Multi-stage Dockerfile for production Go backend
# ============================================================

# Stage 1: Build
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with optimizations:
# -ldflags="-s -w"  strips debug info and symbol table (smaller binary)
# CGO_ENABLED=0     static binary (no libc dependency)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/server \
    .

# Stage 2: Runtime (minimal image)
FROM alpine:3.20

# Security: run as non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Install CA certificates for HTTPS
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Own the binary
RUN chown appuser:appgroup /app/server

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 4000

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:4000/health/live || exit 1

# Run the binary
ENTRYPOINT ["/app/server"]
