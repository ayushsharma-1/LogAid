# LogAid Makefile

# Variables
BINARY_NAME=logaid
BUILD_DIR=build
MAIN_PACKAGE=./cmd
VERSION=1.0.0
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Default target
.PHONY: all
all: clean test build

# Build the application
.PHONY: build
build:
	@echo "Building LogAid..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	GEMINI_API_KEY=test-key go test -v ./cmd
	cd tests && GEMINI_API_KEY=test-key go test -v .
	@echo "Tests complete"

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(BUILD_DIR)
	GEMINI_API_KEY=test-key go test -v -coverprofile=$(BUILD_DIR)/coverage.out ./cmd ./tests
	go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo "Coverage report generated: $(BUILD_DIR)/coverage.html"

# Run linting
.PHONY: lint
lint:
	@echo "Running linting..."
	go vet ./...
	go fmt ./...
	@echo "Linting complete"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	rm -f logaid_test
	@echo "Clean complete"

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download
	@echo "Dependencies installed"

# Run the application
.PHONY: run
run: build
	@echo "Running LogAid..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Development build (faster, no optimizations)
.PHONY: dev
dev:
	@echo "Building development version..."
	go build -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Check if environment is set up correctly
.PHONY: check-env
check-env:
	@echo "Checking environment..."
	@if [ -z "$(GEMINI_API_KEY)" ] && [ -z "$(OPENAI_API_KEY)" ]; then \
		echo "Warning: No API keys found in environment"; \
		echo "Set GEMINI_API_KEY or OPENAI_API_KEY"; \
	else \
		echo "Environment OK"; \
	fi

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	@echo "Multi-platform build complete"

# Package for distribution
.PHONY: package
package: build-all
	@echo "Creating distribution packages..."
	@mkdir -p $(BUILD_DIR)/dist
	cd $(BUILD_DIR) && tar -czf dist/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	cd $(BUILD_DIR) && tar -czf dist/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	cd $(BUILD_DIR) && zip -q dist/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "Distribution packages created in $(BUILD_DIR)/dist/"

# Help
.PHONY: help
help:
	@echo "LogAid Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  all          - Clean, test, and build (default)"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  lint         - Run linting and formatting"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  run          - Build and run the application"
	@echo "  dev          - Quick development build"
	@echo "  check-env    - Check environment setup"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  package      - Create distribution packages"
	@echo "  help         - Show this help message"
