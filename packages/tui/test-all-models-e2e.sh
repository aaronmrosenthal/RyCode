#!/bin/bash

# E2E Test Script for All CLI Provider Models
# Tests each model by attempting to send a simple message

set -e

echo "=== Testing All CLI Provider Models E2E ==="
echo ""

# Get all providers and models
PROVIDERS_JSON=$(cd ../rycode && bun run src/auth/cli.ts cli-providers 2>/dev/null)

# Test results
WORKING_MODELS=()
FAILED_MODELS=()

# Function to test a model
test_model() {
    local provider=$1
    local model=$2

    echo -n "Testing $provider / $model ... "

    # Create a test session with the model
    RESPONSE=$(curl -s -X POST http://127.0.0.1:4096/session \
        -H "Content-Type: application/json" \
        -d "{\"providerID\":\"$provider\",\"modelID\":\"$model\"}" 2>&1)

    # Check if response contains an error
    if echo "$RESPONSE" | grep -q "ProviderModelNotFoundError\|error\|Error"; then
        echo "❌ FAILED"
        FAILED_MODELS+=("$provider/$model")
        echo "  Error: $(echo "$RESPONSE" | jq -r '.error // .message // .' 2>/dev/null || echo "$RESPONSE")"
        return 1
    else
        # Extract session ID
        SESSION_ID=$(echo "$RESPONSE" | jq -r '.id' 2>/dev/null)

        if [ "$SESSION_ID" != "null" ] && [ -n "$SESSION_ID" ]; then
            echo "✅ SUCCESS (session: $SESSION_ID)"
            WORKING_MODELS+=("$provider/$model")

            # Clean up - delete the test session
            curl -s -X DELETE "http://127.0.0.1:4096/session/$SESSION_ID" >/dev/null 2>&1
            return 0
        else
            echo "❌ FAILED (no session ID)"
            FAILED_MODELS+=("$provider/$model")
            return 1
        fi
    fi
}

# Parse JSON and test each model
echo "$PROVIDERS_JSON" | jq -r '.providers[] | "\(.provider) \(.models | join(" "))"' | while read -r provider models; do
    echo ""
    echo "=== Testing Provider: $provider ==="
    for model in $models; do
        test_model "$provider" "$model" || true
    done
done

echo ""
echo "=== Test Summary ==="
echo "Working models: ${#WORKING_MODELS[@]}"
echo "Failed models: ${#FAILED_MODELS[@]}"
echo ""

if [ ${#FAILED_MODELS[@]} -gt 0 ]; then
    echo "Failed models to remove from configuration:"
    printf '%s\n' "${FAILED_MODELS[@]}"
fi
