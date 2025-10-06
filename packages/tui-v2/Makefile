.PHONY: build test test-unit test-integration coverage lint clean install dev

BINARY_NAME=rycode
BUILD_DIR=../../packages/rycode/dist
GO=go
GOFLAGS=-v

# Build binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/rycode

# Run all tests
test:
	@echo "Running all tests..."
	$(GO) test -v ./...

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	$(GO) test -v ./internal/... ./pkg/...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	$(GO) test -v -tags=integration ./test/integration/...

# Generate coverage report
coverage:
	@echo "Generating coverage report..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found, installing..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)/$(BINARY_NAME)
	rm -f coverage.out coverage.html

# Install binary to user bin
install: build
	@echo "Installing to $(HOME)/bin..."
	@mkdir -p $(HOME)/bin
	cp $(BUILD_DIR)/$(BINARY_NAME) $(HOME)/bin/

# Development mode (watch for changes)
dev: build
	@echo "Running in dev mode..."
	$(BUILD_DIR)/$(BINARY_NAME)

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download

# Verify dependencies
verify:
	@echo "Verifying dependencies..."
	$(GO) mod verify

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GO) mod tidy

# Run demo
demo: build
	@echo "Running demo..."
	$(BUILD_DIR)/$(BINARY_NAME) --demo

# Run chat interface
chat: build
	@echo "Running chat interface..."
	$(BUILD_DIR)/$(BINARY_NAME) --chat

# Run workspace (FileTree + Chat)
workspace: build
	@echo "Running workspace..."
	$(BUILD_DIR)/$(BINARY_NAME) --workspace

# Show help
help:
	@echo "Available targets:"
	@echo "  build              - Build the binary"
	@echo "  test               - Run all tests"
	@echo "  test-unit          - Run unit tests only"
	@echo "  test-integration   - Run integration tests"
	@echo "  coverage           - Generate coverage report"
	@echo "  lint               - Run linter"
	@echo "  clean              - Clean build artifacts"
	@echo "  install            - Install binary to ~/bin"
	@echo "  dev                - Run in development mode"
	@echo "  fmt                - Format code"
	@echo "  deps               - Download dependencies"
	@echo "  verify             - Verify dependencies"
	@echo "  tidy               - Tidy dependencies"
	@echo "  demo               - Run theme demo mode"
	@echo "  chat               - Run interactive chat interface"
	@echo "  workspace          - Run workspace (FileTree + Chat)"
	@echo "  help               - Show this help"
