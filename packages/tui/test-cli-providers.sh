#!/bin/bash
#
# Core Build Unit Test: CLI Providers Authentication
#
# This test validates that ALL SOTA CLI providers are authenticated and accessible.
# This is the foundational test for the Tab cycling workflow.
#
# Exit codes:
#   0 = All providers authenticated
#   1 = One or more providers failed

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "=== Core Build Unit Test: CLI Providers Authentication ==="
echo ""
echo "Project root: $PROJECT_ROOT"
echo "Test script: test_cli_providers_e2e.go"
echo ""

# Compile the test
echo "[1] Compiling test..."
cd "$SCRIPT_DIR"
if ! go build -o /tmp/test_cli_providers_e2e test_cli_providers_e2e.go; then
    echo "    ✗ ERROR: Failed to compile test"
    exit 1
fi
echo "    ✓ Test compiled to /tmp/test_cli_providers_e2e"
echo ""

# Run the test
echo "[2] Running test..."
echo ""
if /tmp/test_cli_providers_e2e; then
    EXIT_CODE=0
else
    EXIT_CODE=$?
fi

echo ""
echo "=== Test Logs ==="
echo "View detailed logs at: /tmp/rycode-e2e-cli-providers.log"
echo ""

# Clean up
rm -f /tmp/test_cli_providers_e2e

exit $EXIT_CODE
