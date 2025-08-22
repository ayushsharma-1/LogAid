# LogAid Development Guide

## Quick Start

```bash
# Clone the repository
git clone https://github.com/ayushsharma-1/LogAid.git
cd LogAid

# Install dependencies
make deps

# Build the application
make build

# Run tests
make test

# Create .env file with your API keys
cp .env.example .env
# Edit .env with your actual API keys

# Run LogAid
./run_logaid.sh version
```

## Testing Strategy

### Unit Tests
Located in `cmd/main_test.go` - tests individual functions and components.

### Integration Tests
Located in `tests/integration_test.go` - tests the complete application flow.

### Running Tests
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run only integration tests
cd tests && go test -v .

# Run only unit tests
go test -v ./cmd
```

### Test Environment
Tests use mock API keys (`test-key`) to avoid hitting real APIs during testing.

## Build Targets

```bash
make build        # Single platform build
make build-all    # Multi-platform builds
make package      # Create distribution packages
make dev          # Fast development build
```

## Code Quality

```bash
make lint         # Run linting and formatting
go vet ./...      # Static analysis
go fmt ./...      # Format code
```

## Project Structure

```
LogAid/
├── cmd/                 # Main application entry point
├── pkg/                 # Core packages
│   ├── ai/             # AI provider implementations
│   ├── cli/            # CLI interface
│   ├── config/         # Configuration management
│   ├── logger/         # Logging system
│   ├── plugin/         # Plugin system
│   └── pty/            # PTY wrapper
├── tests/              # Integration tests
├── .github/workflows/  # CI/CD pipelines
├── build/              # Build artifacts
└── dist/               # Distribution packages
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run `make test` to ensure all tests pass
6. Run `make lint` to check code quality
7. Submit a pull request

## Environment Variables

See `.env.example` for all available configuration options.

Required:
- `GEMINI_API_KEY` or `OPENAI_API_KEY` (at least one)

Optional:
- `AI_PROVIDER` - Preferred AI provider
- `LOG_LEVEL` - Logging verbosity
- `ENABLE_COLORS` - Colored output
- And many more...

## Continuous Integration

GitHub Actions automatically:
- Runs tests on multiple Go versions
- Builds for multiple platforms
- Performs security scanning
- Creates releases for tagged versions

## Releasing

1. Update version in `Makefile`
2. Create and push a git tag: `git tag v1.0.1 && git push origin v1.0.1`
3. GitHub Actions will automatically create a release with binaries
