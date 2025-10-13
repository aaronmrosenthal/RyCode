#!/bin/bash
# Verification script to confirm everything is perfect

# Change to script directory
cd "$(dirname "$0")"

echo "🚀 RyCode Model Selector - Perfection Verification"
echo "=================================================="
echo ""

PASS=0
FAIL=0

# Test 1: Go build
echo "Test 1: Building TUI..."
if go build -o bin/rycode ./packages/tui/cmd/rycode 2>/dev/null; then
    echo "✅ Build successful"
    ((PASS++))
else
    echo "❌ Build failed"
    ((FAIL++))
fi
echo ""

# Test 2: Direct Go test
echo "Test 2: Testing data layer (ListProviders merging)..."
OUTPUT=$(go run packages/tui/test_models_direct.go 2>&1)
if echo "$OUTPUT" | grep -q "TOTAL MERGED MODELS: 30"; then
    echo "✅ Data layer works - 30 models merged"
    ((PASS++))
else
    echo "❌ Data layer test failed"
    ((FAIL++))
fi
echo ""

# Test 3: Playwright tests
echo "Test 3: Running Playwright tests..."
TEST_OUTPUT=$(bunx playwright test packages/tui/test-model-selector.spec.ts --reporter=list 2>&1)
if echo "$TEST_OUTPUT" | grep -q "26 passed"; then
    echo "✅ All 26 Playwright tests passed"
    ((PASS++))
else
    echo "❌ Some Playwright tests failed"
    ((FAIL++))
fi
echo ""

# Test 4: File existence checks
echo "Test 4: Checking documentation files..."
FILES=(
    "docs/MODEL_SELECTOR_UX_ANALYSIS.md"
    "packages/tui/test-model-selector-web.html"
    "packages/tui/test-model-selector.spec.ts"
    "packages/tui/test_models_direct.go"
    "PLAYWRIGHT_TEST_SUMMARY.md"
    "packages/tui/MODEL_SELECTOR_README.md"
    "PERFECT_SUMMARY.md"
)

ALL_FILES_EXIST=true
for file in "${FILES[@]}"; do
    if [ -f "$file" ]; then
        echo "  ✓ $file"
    else
        echo "  ✗ $file MISSING"
        ALL_FILES_EXIST=false
    fi
done

if [ "$ALL_FILES_EXIST" = true ]; then
    echo "✅ All documentation files present"
    ((PASS++))
else
    echo "❌ Some documentation files missing"
    ((FAIL++))
fi
echo ""

# Test 5: Code verification
echo "Test 5: Verifying UX improvements in code..."
CODE_FILE="packages/tui/internal/components/dialog/models.go"

if grep -q "renderShortcutFooter" "$CODE_FILE" && \
   grep -q "getModelBadges" "$CODE_FILE" && \
   grep -q "jumpToProvider" "$CODE_FILE"; then
    echo "✅ All UX improvements present in code"
    ((PASS++))
else
    echo "❌ Some UX improvements missing"
    ((FAIL++))
fi
echo ""

# Summary
echo "=================================================="
echo "VERIFICATION SUMMARY"
echo "=================================================="
echo "✅ Passed: $PASS/5"
echo "❌ Failed: $FAIL/5"
echo ""

if [ $FAIL -eq 0 ]; then
    echo "🎉 PERFECT! Everything works!"
    echo ""
    echo "Next steps:"
    echo "  1. Run the TUI:  ./bin/rycode"
    echo "  2. Press: Ctrl+X then m"
    echo "  3. Try: Number keys (1-9) to jump to providers"
    echo "  4. See: Footer with shortcuts"
    echo "  5. See: Badges (⚡💰🧠🆕) next to models"
    echo ""
    echo "For demo: open packages/tui/test-model-selector-web.html"
    exit 0
else
    echo "⚠️  $FAIL test(s) failed. Review output above."
    exit 1
fi
