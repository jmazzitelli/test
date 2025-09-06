# Build stage
FROM docker.io/golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wms-proxy ./cmd/proxy

# Final stage
FROM docker.io/alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -s /bin/sh wms-proxy

# Set working directory
WORKDIR /app

# Create certificates directory
RUN mkdir -p /app/certs && chown wms-proxy:wms-proxy /app/certs

# Copy binary from builder stage
COPY --from=builder /app/wms-proxy .

# Copy certificate generation script
COPY scripts/generate-certs.sh /app/

# Change ownership to non-root user
RUN chown wms-proxy:wms-proxy /app/wms-proxy /app/generate-certs.sh && \
    chmod +x /app/generate-certs.sh

# Switch to non-root user
USER wms-proxy

# Expose ports (HTTP and HTTPS)
EXPOSE 8080 8443

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./wms-proxy"]
