# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o logaid .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    bash \
    zsh \
    fish \
    ca-certificates \
    tzdata

# Create non-root user
RUN addgroup -g 1001 logaid && \
    adduser -D -s /bin/bash -G logaid -u 1001 logaid

# Set working directory
WORKDIR /home/logaid

# Copy binary from builder
COPY --from=builder /app/logaid /usr/local/bin/logaid

# Make binary executable
RUN chmod +x /usr/local/bin/logaid

# Switch to non-root user
USER logaid

# Set default shell
ENV SHELL=/bin/bash

# Expose port for potential web interface (future feature)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD logaid version || exit 1

# Default command
CMD ["logaid", "--help"]
