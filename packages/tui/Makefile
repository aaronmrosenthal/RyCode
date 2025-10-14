.PHONY: help build test test-hello-all test-unit test-integration clean install

# Default target
help:
	@echo "RyCode TUI - Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  make build              - Build the rycode binary"
	@echo "  make test               - Run all tests (unit + integration + E2E)"
	@echo "  make test-hello-all     - Run hello-all-providers E2E test (CORE BUILD TEST)"
	@echo "  make test-unit          - Run Go unit tests"
	@echo "  make test-integration   - Run integration tests"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make install            - Build and install to bin/"
	@echo ""

# Build the rycode binary
build:
	@echo "Building rycode..."
	@go build -ldflags="-s -w" -o bin/rycode ./cmd/rycode
	@echo "✓ Binary built: bin/rycode"

# Run all tests
test: test-unit test-integration test-cli-providers
	@echo ""
	@echo "✅ All tests passed!"

# Core build unit test: CLI providers authentication
test-cli-providers:
	@echo "Running core build test: CLI Providers Authentication..."
	@./test-cli-providers.sh

# Run Go unit tests
test-unit:
	@echo "Running unit tests..."
	@go test -v -race ./internal/...
	@echo "✓ Unit tests passed"

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@go test -v -tags=integration ./...
	@echo "✓ Integration tests passed"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f bin/rycode
	@rm -f /tmp/test_hello_all_providers_e2e
	@rm -f /tmp/rycode-e2e-*.log
	@rm -rf .next
	@echo "✓ Clean complete"

# Build and install
install: build
	@echo "Installing to project bin/..."
	@mkdir -p ../../bin
	@cp bin/rycode ../../bin/rycode
	@echo "✓ Installed to ../../bin/rycode"
