#!/usr/bin/env bash
#
# Quick API Key Setup Script
#
# Usage: ./scripts/add-api-keys.sh
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CLI_PATH="$PROJECT_ROOT/packages/rycode/src/auth/cli.ts"

echo ""
echo "🔐 RyCode API Key Setup"
echo "======================="
echo ""
echo "This script will help you add your API keys to RyCode."
echo "You need at least 2 providers for Tab cycling to work."
echo ""

# Function to add a key
add_key() {
    local provider=$1
    local provider_name=$2
    local key_prefix=$3
    local url=$4

    echo ""
    echo "📍 $provider_name"
    echo "   Get your API key: $url"
    echo ""
    read -p "   Enter your $provider_name API key (or press Enter to skip): " key

    if [ -z "$key" ]; then
        echo "   ⏭️  Skipped"
        return
    fi

    # Validate key format
    if [[ ! "$key" =~ ^$key_prefix ]]; then
        echo "   ⚠️  Warning: Key doesn't start with expected prefix '$key_prefix'"
        read -p "   Continue anyway? (y/N): " confirm
        if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
            echo "   ⏭️  Skipped"
            return
        fi
    fi

    # Add the key
    echo "   Adding key to RyCode..."
    if cd "$PROJECT_ROOT/packages/rycode" && bun run "$CLI_PATH" auth "$provider" "$key" > /dev/null 2>&1; then
        echo "   ✅ $provider_name authenticated successfully!"
    else
        echo "   ❌ Failed to add $provider_name key"
        echo "   You can try manually: cd packages/rycode && bun run src/auth/cli.ts auth $provider YOUR-KEY"
    fi
}

# Add keys for each provider
add_key "anthropic" "Claude (Anthropic)" "sk-ant-" "https://console.anthropic.com/"
add_key "google" "Gemini (Google)" "AIza" "https://makersuite.google.com/app/apikey"
add_key "openai" "OpenAI (GPT/Codex)" "sk-" "https://platform.openai.com/api-keys"
add_key "grok" "Grok (xAI)" "" "https://console.x.ai/"
add_key "qwen" "Qwen (Alibaba)" "" "https://dashscope.aliyun.com/"

# Verify
echo ""
echo "🔍 Verifying authentication status..."
echo ""
cd "$PROJECT_ROOT/packages/rycode"
if bun run "$CLI_PATH" list; then
    echo ""
    echo "✅ Setup complete!"
    echo ""
    echo "🎯 Next steps:"
    echo "   1. Run the TUI: cd packages/tui && go run cmd/rycode/main.go"
    echo "   2. Type /model to open the model selector"
    echo "   3. Press Tab to cycle between your authenticated providers"
    echo ""
else
    echo ""
    echo "❌ Failed to verify authentication"
    echo "Please check the error messages above"
    echo ""
fi
