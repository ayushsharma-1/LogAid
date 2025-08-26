# LogAid Makefile
.PHONY: all build test clean install uninstall deps fmt vet lint release help

# Variables
BINARY_NAME=logaid
PACKAGE=github.com/ayushsharma-1/LogAid
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Build directory
BUILD_DIR=build

# Default target
all: test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

# Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	
	# macOS
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-macos-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-macos-arm64 .
	
	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

# Install binary to system
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)

# Uninstall binary from system
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from /usr/local/bin..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# Create a release package
release: build-all
	@echo "Creating release packages..."
	@mkdir -p $(BUILD_DIR)/release
	
	# Create tarballs for each platform
	cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-$(VERSION)-macos-amd64.tar.gz $(BINARY_NAME)-macos-amd64
	cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-$(VERSION)-macos-arm64.tar.gz $(BINARY_NAME)-macos-arm64
	cd $(BUILD_DIR) && zip release/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	
	# Create checksums
	cd $(BUILD_DIR)/release && sha256sum * > checksums.txt
	
	@echo "Release packages created in $(BUILD_DIR)/release/"

# Development targets
dev-test:
	@echo "Running development tests..."
	./$(BUILD_DIR)/$(BINARY_NAME) test

dev-run:
	@echo "Running development version..."
	./$(BUILD_DIR)/$(BINARY_NAME) run

# Docker targets
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .

docker-run:
	@echo "Running Docker container..."
	docker run --rm -it $(BINARY_NAME):$(VERSION)

# Help target
help:
	@echo "LogAid Build System"
	@echo "=================="
	@echo ""
	@echo "Available targets:"
	@echo "  build       - Build binary for current platform"
	@echo "  build-all   - Build binaries for all supported platforms"
	@echo "  test        - Run tests"
	@echo "  deps        - Install/update dependencies"
	@echo "  fmt         - Format source code"
	@echo "  vet         - Run go vet"
	@echo "  lint        - Run linter (requires golangci-lint)"
	@echo "  clean       - Clean build artifacts"
	@echo "  install     - Install binary to /usr/local/bin"
	@echo "  uninstall   - Remove binary from /usr/local/bin"
	@echo "  release     - Create release packages"
	@echo "  dev-test    - Test the built binary"
	@echo "  dev-run     - Run the built binary"
	@echo "  docker-build- Build Docker image"
	@echo "  docker-run  - Run Docker container"
	@echo "  help        - Show this help message"
