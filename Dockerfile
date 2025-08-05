# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git (needed for go mod)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o logaid .

# Final stage
FROM alpine:latest

# Install common tools that LogAid might help with
RUN apk --no-cache add \
    git \
    docker \
    curl \
    bash \
    ca-certificates \
    && rm -rf /var/cache/apk/*

WORKDIR /root/

# Copy the binary
COPY --from=builder /app/logaid .

# Create a non-root user
RUN adduser -D -s /bin/bash logaid
USER logaid
WORKDIR /home/logaid

# Copy binary to user directory
COPY --from=builder --chown=logaid:logaid /app/logaid ./logaid

# Make it executable
RUN chmod +x ./logaid

# Add to PATH
ENV PATH="/home/logaid:${PATH}"

ENTRYPOINT ["./logaid"]
