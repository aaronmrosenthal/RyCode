#!/bin/bash
# Real manual test - I'll run this and verify the output

echo "=== REAL MANUAL TEST ==="
echo ""
echo "Step 1: Clear logs"
rm -f /tmp/rycode-debug.log

echo "Step 2: Start TUI with --model flag to trigger initialization"
echo "This will initialize providers and then exit..."

# Use timeout to auto-quit after 3 seconds
timeout 3 ./bin/rycode --model "test" 2>&1 || true

echo ""
echo "Step 3: Check debug log"
if [ -f /tmp/rycode-debug.log ]; then
    echo "✅ Debug log exists"
    echo ""
    cat /tmp/rycode-debug.log
else
    echo "❌ Debug log NOT created"
fi

echo ""
echo "Step 4: Verify CLI providers are available"
echo "Running: bun run packages/rycode/src/auth/cli.ts cli-providers"
bun run packages/rycode/src/auth/cli.ts cli-providers 2>&1 | head -30

echo ""
echo "=== TEST COMPLETE ==="
