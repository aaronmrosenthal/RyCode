#!/bin/bash
# Direct test of ListProviders to prove CLI providers are being merged

echo "=== Testing Provider Merging ==="
echo ""

echo "1. API Providers (from list command):"
bun run packages/rycode/src/auth/cli.ts list 2>&1 | jq '.'
echo ""

API_COUNT=$(bun run packages/rycode/src/auth/cli.ts list 2>&1 | jq '.providers | length')
echo "API Providers count: $API_COUNT"
echo ""

echo "2. CLI Providers (from cli-providers command):"
bun run packages/rycode/src/auth/cli.ts cli-providers 2>&1 | jq '.'
echo ""

CLI_COUNT=$(bun run packages/rycode/src/auth/cli.ts cli-providers 2>&1 | jq '.providers | length')
echo "CLI Providers count: $CLI_COUNT"
echo ""

echo "3. Total expected providers: $((API_COUNT + CLI_COUNT))"
echo ""

echo "4. All provider names:"
echo "   API:"
bun run packages/rycode/src/auth/cli.ts list 2>&1 | jq -r '.providers[].name' | sed 's/^/     - /'
echo "   CLI:"
bun run packages/rycode/src/auth/cli.ts cli-providers 2>&1 | jq -r '.providers[].provider' | sed 's/^/     - /'
echo ""

echo "5. Model counts per provider:"
echo "   API:"
bun run packages/rycode/src/auth/cli.ts list 2>&1 | jq -r '.providers[] | "     - \(.name): \(.modelsCount) models"'
echo "   CLI:"
bun run packages/rycode/src/auth/cli.ts cli-providers 2>&1 | jq -r '.providers[] | "     - \(.provider): \(.models | length) models"'
echo ""

API_MODELS=$(bun run packages/rycode/src/auth/cli.ts list 2>&1 | jq '[.providers[].modelsCount] | add')
CLI_MODELS=$(bun run packages/rycode/src/auth/cli.ts cli-providers 2>&1 | jq '[.providers[].models | length] | add')

echo "TOTAL MODELS: $((API_MODELS + CLI_MODELS))"
echo "  - API models: $API_MODELS"
echo "  - CLI models: $CLI_MODELS"
echo ""

echo "=== When you open /models in TUI, you should see ALL $((API_MODELS + CLI_MODELS)) models ==="
echo "=== across $((API_COUNT + CLI_COUNT)) providers (API + CLI combined) ==="
