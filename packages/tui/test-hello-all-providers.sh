#!/bin/bash
#
# Core Build Unit Test: Hello All Providers E2E
#
# This test validates that ALL authenticated SOTA providers respond to messages.
# It simulates the Tab cycling workflow by:
# 1. Loading all CLI providers
# 2. For each authenticated provider:
#    - Create a session
#    - Send "hello" message
#    - Validate response
#    - Clean up session
#
# Exit codes:
#   0 = All providers passed
#   1 = One or more providers failed

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "=== Core Build Unit Test: Hello All Providers ==="
echo ""
echo "Project root: $PROJECT_ROOT"
echo "Test script: test_hello_all_providers_e2e.go"
echo ""

# Check if API server is running
echo "[1] Checking if API server is running on port 4096..."
if ! curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:4096/health > /dev/null 2>&1; then
    echo "    ✗ ERROR: API server is not running on port 4096"
    echo ""
    echo "    Please start the API server first:"
    echo "      cd packages/rycode && bun run dev"
    echo ""
    exit 1
fi
echo "    ✓ API server is running"
echo ""

# Compile the test
echo "[2] Compiling test..."
cd "$SCRIPT_DIR"
if ! go build -o /tmp/test_hello_all_providers_e2e test_hello_all_providers_e2e.go; then
    echo "    ✗ ERROR: Failed to compile test"
    exit 1
fi
echo "    ✓ Test compiled to /tmp/test_hello_all_providers_e2e"
echo ""

# Run the test
echo "[3] Running test..."
echo ""
if /tmp/test_hello_all_providers_e2e; then
    EXIT_CODE=0
else
    EXIT_CODE=$?
fi

echo ""
echo "=== Test Logs ==="
echo "View detailed logs at: /tmp/rycode-e2e-hello-all.log"
echo ""

# Clean up
rm -f /tmp/test_hello_all_providers_e2e

exit $EXIT_CODE
