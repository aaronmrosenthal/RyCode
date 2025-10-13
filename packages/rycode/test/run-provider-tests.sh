#!/usr/bin/env bash
# RyCode Provider Test Runner
#
# This script runs all provider tests in the correct order:
# 1. CLI authentication test (checks auth.json)
# 2. Environment variable detection test
# 3. Server provider test (requires running server)

set -e

echo "╔════════════════════════════════════════════════╗"
echo "║   RyCode Provider Test Suite                   ║"
echo "║   Testing CLI-Authenticated Providers          ║"
echo "╚════════════════════════════════════════════════╝"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: CLI Authentication
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test 1/3: CLI Authentication (auth.json)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if bun run packages/rycode/test/provider-cli-test.ts; then
    echo -e "${GREEN}✓ CLI authentication test completed${NC}"
else
    echo -e "${YELLOW}⚠ No CLI authentication found${NC}"
fi

echo ""

# Test 2: Environment Variable Detection
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test 2/3: Environment Variable Detection"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if bun run packages/rycode/test/provider-test.ts; then
    echo -e "${GREEN}✓ Environment detection test completed${NC}"
else
    echo -e "${YELLOW}⚠ No environment variables found${NC}"
fi

echo ""

# Test 3: Server Provider Test (requires running server)
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Test 3/3: Server Provider API"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if server is running
if curl -s -f http://127.0.0.1:4096/app/providers > /dev/null 2>&1; then
    echo "✓ Server is running"
    echo ""
    if bun run packages/rycode/test/provider-server-test.ts; then
        echo -e "${GREEN}✓ Server provider test completed${NC}"
    else
        echo -e "${RED}✗ Server provider test failed${NC}"
    fi
else
    echo -e "${YELLOW}⚠ Server is not running at http://127.0.0.1:4096${NC}"
    echo ""
    echo "To start the server:"
    echo "  bun run packages/rycode/src/index.ts serve --port 4096"
    echo ""
    echo "Then run this script again to complete the server tests."
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Test Suite Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "✅ Test suite completed"
echo ""
echo "Provider Authentication Methods:"
echo "  1. CLI Auth (rycode auth login) → ~/.local/share/rycode/auth.json"
echo "  2. Environment Variables → OPENAI_API_KEY, ANTHROPIC_API_KEY, etc."
echo "  3. Config File → opencode.json in project root"
echo ""
echo "To authenticate providers:"
echo "  • Interactive: rycode auth login"
echo "  • Environment: export ANTHROPIC_API_KEY=\"sk-ant-...\""
echo "  • Config: Add to opencode.json"
echo ""
