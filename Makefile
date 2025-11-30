.PHONY: help build build-all install test test-verbose test-coverage test-integration lint lint-fix install-lint clean fmt

# Variables
BINARY_NAME=go-pdf-extractor
EXAMPLE_BINARY=example
CMD_DIR=./cmd/example
PKG_DIR=./pkg/...
TEST_DIR=./tests/...

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
help:
	@echo "$(GREEN)go-pdf-extractor - Makefile targets$(NC)"
	@echo ""
	@echo "$(YELLOW)Build targets:$(NC)"
	@echo "  make build            - Build the example binary"
	@echo "  make build-all        - Build for multiple platforms (Linux, macOS, Windows)"
	@echo "  make install          - Download and install dependencies"
	@echo ""
	@echo "$(YELLOW)Test targets:$(NC)"
	@echo "  make test             - Run all tests"
	@echo "  make test-verbose     - Run tests with verbose output"
	@echo "  make test-coverage    - Run tests with coverage report"
	@echo "  make test-integration - Run integration tests (requires OPENAI_API_KEY)"
	@echo ""
	@echo "$(YELLOW)Lint targets:$(NC)"
	@echo "  make lint             - Run golangci-lint"
	@echo "  make lint-fix         - Run golangci-lint with auto-fix"
	@echo "  make install-lint     - Install golangci-lint"
	@echo ""
	@echo "$(YELLOW)Utility targets:$(NC)"
	@echo "  make clean            - Remove build artifacts and binaries"
	@echo "  make fmt              - Format code with gofmt"
	@echo "  make help             - Display this help message"
	@echo ""

## install: Download and install dependencies
install:
	@echo "$(GREEN)Installing dependencies...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)Dependencies installed successfully$(NC)"

## build: Build the example binary
build:
	@echo "$(GREEN)Building $(EXAMPLE_BINARY)...$(NC)"
	go build -o $(EXAMPLE_BINARY) $(CMD_DIR)
	@echo "$(GREEN)Build complete: $(EXAMPLE_BINARY)$(NC)"

## build-all: Build for multiple platforms
build-all:
	@echo "$(GREEN)Building for multiple platforms...$(NC)"
	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 go build -o $(EXAMPLE_BINARY)-linux-amd64 $(CMD_DIR)
	@echo "Building for Linux (arm64)..."
	GOOS=linux GOARCH=arm64 go build -o $(EXAMPLE_BINARY)-linux-arm64 $(CMD_DIR)
	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 go build -o $(EXAMPLE_BINARY)-darwin-amd64 $(CMD_DIR)
	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 go build -o $(EXAMPLE_BINARY)-darwin-arm64 $(CMD_DIR)
	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 go build -o $(EXAMPLE_BINARY)-windows-amd64.exe $(CMD_DIR)
	@echo "$(GREEN)Cross-platform builds complete$(NC)"

## test: Run all tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	go test $(TEST_DIR)
	@echo "$(GREEN)Tests completed$(NC)"

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "$(GREEN)Running tests with verbose output...$(NC)"
	go test -v $(TEST_DIR)

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	go test -cover $(TEST_DIR)
	@echo ""
	@echo "$(YELLOW)Generating detailed coverage report...$(NC)"
	go test -coverprofile=coverage.out $(TEST_DIR)
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## test-integration: Run integration tests (requires OPENAI_API_KEY)
test-integration:
	@echo "$(GREEN)Running integration tests...$(NC)"
	@if [ -z "$$OPENAI_API_KEY" ]; then \
		echo "$(RED)Error: OPENAI_API_KEY environment variable is not set$(NC)"; \
		echo "Please set it with: export OPENAI_API_KEY=your-api-key"; \
		exit 1; \
	fi
	INTEGRATION_TEST=true go test -v $(TEST_DIR)
	@echo "$(GREEN)Integration tests completed$(NC)"

## lint: Run golangci-lint
lint:
	@echo "$(GREEN)Running golangci-lint...$(NC)"
	@if ! command -v golangci-lint > /dev/null; then \
		echo "$(RED)golangci-lint is not installed$(NC)"; \
		echo "Run 'make install-lint' to install it"; \
		exit 1; \
	fi
	golangci-lint run --config .golangci.yml ./...
	@echo "$(GREEN)Linting completed$(NC)"

## lint-fix: Run golangci-lint with auto-fix
lint-fix:
	@echo "$(GREEN)Running golangci-lint with auto-fix...$(NC)"
	@if ! command -v golangci-lint > /dev/null; then \
		echo "$(RED)golangci-lint is not installed$(NC)"; \
		echo "Run 'make install-lint' to install it"; \
		exit 1; \
	fi
	golangci-lint run --config .golangci.yml --fix ./...
	@echo "$(GREEN)Linting with auto-fix completed$(NC)"

## install-lint: Install golangci-lint
install-lint:
	@echo "$(GREEN)Installing golangci-lint...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		echo "$(YELLOW)golangci-lint is already installed$(NC)"; \
		golangci-lint version; \
	else \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
		echo "$(GREEN)golangci-lint installed successfully$(NC)"; \
	fi

## fmt: Format code with gofmt
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	gofmt -s -w .
	@echo "$(GREEN)Code formatted$(NC)"

## clean: Remove build artifacts and binaries
clean:
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	rm -f $(EXAMPLE_BINARY)
	rm -f $(EXAMPLE_BINARY)-*
	rm -f coverage.out coverage.html
	go clean
	@echo "$(GREEN)Clean completed$(NC)"
