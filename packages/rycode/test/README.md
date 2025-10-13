# RyCode Provider Testing

This directory contains tests to verify that RyCode correctly supports developers using their own API keys from various providers (OpenAI, Anthropic, Google, Qwen, etc.).

## Test Files

### `provider-test.ts`
**Detection & Configuration Test**

Tests that RyCode can detect API keys from:
- Environment variables (`OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, etc.)
- Configuration files (`~/.local/share/rycode/auth.json`)
- Auto-detection from common credential locations

```bash
# Run detection test
bun run packages/rycode/test/provider-test.ts
```

Expected output:
```
🔍 Testing Provider Auto-Detection...
Found 3 credential(s) in environment:
  ✓ openai (from env: OPENAI_API_KEY)
  ✓ anthropic (from env: ANTHROPIC_API_KEY)
  ✓ google (from env: GOOGLE_API_KEY)
```

### `provider-e2e-test.ts`
**End-to-End API Test**

Makes actual API calls to verify providers work correctly with user keys.

```bash
# Test all configured providers
OPENAI_API_KEY=sk-... ANTHROPIC_API_KEY=sk-ant-... \
  bun run packages/rycode/test/provider-e2e-test.ts

# Test specific provider only
OPENAI_API_KEY=sk-... \
  bun run packages/rycode/test/provider-e2e-test.ts openai
```

Expected output:
```
🧪 Running provider tests...
  Testing OpenAI (GPT-3.5)... ✅ PASS - Received 52 chars (1234ms)
  Testing Anthropic (Claude)... ✅ PASS - Received 48 chars (892ms)
  Testing Google (Gemini)... ⏭️  SKIPPED (GOOGLE_API_KEY not set)
```

## 🆕 New: CLI Provider Tests

We've added new tests specifically for CLI-authenticated providers:

- **`provider-cli-test.ts`** - Tests authentication via `rycode auth login`
- **`provider-server-test.ts`** - Tests providers via RyCode server API
- **`run-provider-tests.sh`** - Runs all tests in sequence
- **`CLI_PROVIDER_TESTING.md`** - Complete CLI testing guide

**Quick test:**
```bash
./packages/rycode/test/run-provider-tests.sh
```

See [CLI_PROVIDER_TESTING.md](./CLI_PROVIDER_TESTING.md) for detailed guide.

---

## Quick Start

### 1. Set Up Environment Variables

```bash
# Option A: Export in terminal
export OPENAI_API_KEY="sk-proj-..."
export ANTHROPIC_API_KEY="sk-ant-..."
export GOOGLE_API_KEY="..."

# Option B: Create .env file
cat > .env << 'EOF'
OPENAI_API_KEY=sk-proj-...
ANTHROPIC_API_KEY=sk-ant-...
GOOGLE_API_KEY=...
EOF
```

### 2. Run Tests

```bash
# Step 1: Detection test
bun run packages/rycode/test/provider-test.ts

# Step 2: E2E test (makes real API calls)
bun run packages/rycode/test/provider-e2e-test.ts
```

### 3. Start RyCode

```bash
# Start server
bun run packages/rycode/src/index.ts serve --port 4096

# In another terminal, list available models
bun run packages/rycode/src/index.ts models

# Or start the TUI
rycode
```

## Supported Providers

| Provider | Environment Variable | Get API Key |
|----------|---------------------|-------------|
| OpenAI | `OPENAI_API_KEY` | https://platform.openai.com/api-keys |
| Anthropic | `ANTHROPIC_API_KEY` | https://console.anthropic.com/settings/keys |
| Google | `GOOGLE_API_KEY` | https://makersuite.google.com/app/apikey |
| Qwen | `QWEN_API_KEY` or `DASHSCOPE_API_KEY` | https://dashscope.console.aliyun.com/ |
| xAI | `XAI_API_KEY` or `GROK_API_KEY` | https://x.ai/api |

### Local Models (No API Key)

For local models (Ollama, LM Studio), configure in `opencode.json`:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "provider": {
    "ollama": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "Ollama (local)",
      "options": {
        "baseURL": "http://localhost:11434/v1"
      },
      "models": {
        "qwen2.5-coder": {
          "name": "Qwen 2.5 Coder (Local)"
        }
      }
    }
  }
}
```

## Testing Checklist

Use this checklist to verify provider support:

- [ ] Environment variables are detected (`provider-test.ts`)
- [ ] API keys work with actual API calls (`provider-e2e-test.ts`)
- [ ] Models appear in `rycode models` command
- [ ] Models work in the TUI
- [ ] Cost tracking shows user's provider accounts
- [ ] Documentation is clear for developers

## Troubleshooting

### API Keys Not Detected

```bash
# Check environment
env | grep -E "OPENAI|ANTHROPIC|GOOGLE|QWEN"

# Run detection test
bun run packages/rycode/test/provider-test.ts
```

### API Calls Failing

```bash
# Run E2E test to see error messages
OPENAI_API_KEY=sk-... bun run packages/rycode/test/provider-e2e-test.ts

# Common issues:
# - Invalid API key format
# - Expired API key
# - Insufficient credits/quota
# - Network connectivity
```

### Models Not Showing

```bash
# List configured providers
bun run packages/rycode/src/index.ts auth list

# List available models
bun run packages/rycode/src/index.ts models

# Check configuration
cat opencode.json
```

## Cost Considerations

Since developers use their own API keys:

1. **Monitor usage in provider dashboards:**
   - OpenAI: https://platform.openai.com/usage
   - Anthropic: https://console.anthropic.com/settings/billing
   - Google: https://console.cloud.google.com/billing

2. **Set up billing alerts to avoid surprises**

3. **Use cost-effective models for development:**
   - GPT-3.5-Turbo instead of GPT-4
   - Claude Haiku instead of Claude Opus
   - Local models (free!)

4. **Consider local models for testing:**
   - Ollama with Qwen 2.5 Coder
   - LM Studio with Code Llama
   - No API costs, no rate limits!

## CI/CD Integration

For testing in CI/CD pipelines:

```yaml
# GitHub Actions example
env:
  OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
  ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}

steps:
  - name: Test Provider Detection
    run: bun run packages/rycode/test/provider-test.ts

  - name: Test Provider E2E
    run: bun run packages/rycode/test/provider-e2e-test.ts
```

## Contributing

To add a new provider test:

1. Add environment variable to `auto-detect.ts`
2. Add test case to `provider-e2e-test.ts`
3. Update documentation in `DEVELOPER_API_KEYS.md`
4. Test with actual API key
5. Submit PR with test results

## See Also

- [Developer API Keys Guide](../../docs/DEVELOPER_API_KEYS.md) - Comprehensive guide for using own API keys
- [Provider Documentation](../../packages/web/src/content/docs/providers.mdx) - Official provider setup docs
- `.env.test.example` - Example environment variables file
