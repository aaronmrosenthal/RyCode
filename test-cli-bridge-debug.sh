#!/bin/bash

echo "=== Testing CLI Bridge from Go ==="
echo ""

cd /Users/aaron/Code/RyCode/RyCode/packages/tui

echo "Running test_list_providers.go to see debug output..."
echo ""

go run test_list_providers.go 2>&1

echo ""
echo "=== End of test ==="
