# Using Your Own API Keys with RyCode

RyCode allows developers to use their own API keys from OpenAI, Anthropic, Google, and other providers. This means **you don't need to pay for tokens through RyCode** - you can use your existing accounts and billing.

## Quick Start

### Option 1: Environment Variables (Recommended for Development)

Set environment variables before running RyCode:

```bash
# OpenAI (for GPT models and Codex)
export OPENAI_API_KEY="sk-..."

# Anthropic (for Claude models)
export ANTHROPIC_API_KEY="sk-ant-..."

# Google (for Gemini models)
export GOOGLE_API_KEY="..."

# Qwen/Alibaba DashScope
export QWEN_API_KEY="..."
# or
export DASHSCOPE_API_KEY="..."

# Then start RyCode
bun run dev
```

### Option 2: .env File (Recommended for Projects)

Create a `.env` file in your project root:

```bash
# .env
OPENAI_API_KEY=sk-your-key-here
ANTHROPIC_API_KEY=sk-ant-your-key-here
GOOGLE_API_KEY=your-google-key-here
QWEN_API_KEY=your-qwen-key-here
```

RyCode will automatically detect and use these credentials.

### Option 3: Interactive Login

Use the RyCode CLI to add credentials interactively:

```bash
# Login to a provider
bun run packages/rycode/src/index.ts auth login

# Select provider (OpenAI, Anthropic, Google, etc.)
# Enter your API key when prompted

# List configured providers
bun run packages/rycode/src/index.ts auth list
```

Credentials are stored in `~/.local/share/rycode/auth.json`.

## Supported Providers

### OpenAI (GPT, Codex)
- **Environment Variable:** `OPENAI_API_KEY`
- **Get API Key:** https://platform.openai.com/api-keys
- **Models:** GPT-4, GPT-4-Turbo, GPT-3.5-Turbo, Codex

```bash
export OPENAI_API_KEY="sk-..."
```

### Anthropic (Claude)
- **Environment Variables:** `ANTHROPIC_API_KEY` or `CLAUDE_API_KEY`
- **Get API Key:** https://console.anthropic.com/settings/keys
- **Models:** Claude 3 Opus, Claude 3 Sonnet, Claude 3 Haiku

```bash
export ANTHROPIC_API_KEY="sk-ant-..."
```

### Google (Gemini)
- **Environment Variable:** `GOOGLE_API_KEY`
- **Get API Key:** https://makersuite.google.com/app/apikey
- **Models:** Gemini Pro, Gemini Pro Vision

```bash
export GOOGLE_API_KEY="..."
```

### Qwen (Alibaba DashScope)
- **Environment Variables:** `QWEN_API_KEY` or `DASHSCOPE_API_KEY`
- **Get API Key:** https://dashscope.console.aliyun.com/
- **Models:** Qwen-Turbo, Qwen-Plus, Qwen-Max

```bash
export QWEN_API_KEY="..."
```

### xAI (Grok)
- **Environment Variables:** `XAI_API_KEY` or `GROK_API_KEY`
- **Get API Key:** https://x.ai/api
- **Models:** Grok Code

```bash
export XAI_API_KEY="..."
```

## Local Models (No API Key Needed)

### Ollama
Configure in `opencode.json`:

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
          "name": "Qwen 2.5 Coder"
        },
        "codellama": {
          "name": "Code Llama"
        }
      }
    }
  }
}
```

### LM Studio
Configure in `opencode.json`:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "provider": {
    "lmstudio": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "LM Studio (local)",
      "options": {
        "baseURL": "http://127.0.0.1:1234/v1"
      },
      "models": {
        "local-model": {
          "name": "Local Model"
        }
      }
    }
  }
}
```

## Testing Your Configuration

Run the provider test script to verify your API keys are detected:

```bash
bun run packages/rycode/test/provider-test.ts
```

Expected output:
```
🔍 Testing Provider Auto-Detection...

Found 4 credential(s) in environment:

  ✓ openai (from env: OPENAI_API_KEY)
  ✓ anthropic (from env: ANTHROPIC_API_KEY)
  ✓ google (from env: GOOGLE_API_KEY)
  ✓ qwen (from env: QWEN_API_KEY)
```

## Verifying Model Availability

List all available models:

```bash
bun run packages/rycode/src/index.ts models
```

This will show all models from your configured providers.

## Troubleshooting

### API Keys Not Detected

1. **Check environment variables:**
   ```bash
   env | grep -E "OPENAI|ANTHROPIC|GOOGLE|QWEN"
   ```

2. **Verify .env file is in project root:**
   ```bash
   cat .env
   ```

3. **Check configured credentials:**
   ```bash
   bun run packages/rycode/src/index.ts auth list
   ```

### Models Not Showing Up

1. **Refresh models database:**
   ```bash
   bun run packages/rycode/src/index.ts models
   ```

2. **Check provider configuration in `opencode.json`**

3. **Verify API key has correct permissions**

## Security Best Practices

1. **Never commit API keys to git:**
   - Add `.env` to `.gitignore`
   - Use `.env.example` for documentation

2. **Use environment-specific keys:**
   - Development: Use lower-tier API keys
   - Production: Use production API keys with rate limits

3. **Rotate keys regularly:**
   - Especially after sharing code or screenshots

4. **Monitor usage:**
   - Check your provider dashboards for unexpected usage
   - Set up billing alerts

## Cost Optimization

Since developers use their own API keys:

1. **Choose cost-effective models:**
   - GPT-3.5-Turbo is cheaper than GPT-4
   - Claude Haiku is cheaper than Claude Opus
   - Gemini Pro is competitive on price

2. **Use local models for development:**
   - Ollama with Qwen 2.5 Coder
   - LM Studio with Code Llama
   - No API costs!

3. **Set up usage limits in provider dashboards:**
   - OpenAI: https://platform.openai.com/account/billing/limits
   - Anthropic: https://console.anthropic.com/settings/limits
   - Google: https://console.cloud.google.com/apis/api/generativelanguage.googleapis.com/quotas

## Example: Complete Setup

```bash
# 1. Create .env file
cat > .env << 'EOF'
OPENAI_API_KEY=sk-proj-...
ANTHROPIC_API_KEY=sk-ant-...
GOOGLE_API_KEY=AI...
QWEN_API_KEY=sk-...
EOF

# 2. Test provider detection
bun run packages/rycode/test/provider-test.ts

# 3. Start RyCode server
bun run packages/rycode/src/index.ts serve --port 4096

# 4. In another terminal, start TUI
rycode

# 5. Select your preferred model with /models
```

## FAQs

**Q: Do I need to use RyCode's OpenCode Zen service?**
A: No! OpenCode Zen is completely optional. You can use your own API keys from any provider.

**Q: Can I mix RyCode providers with my own API keys?**
A: Yes! You can use OpenCode Zen for some models and your own API keys for others.

**Q: Are there any rate limits?**
A: Rate limits depend on your provider account tier. RyCode doesn't impose additional limits.

**Q: Can I use multiple API keys for the same provider?**
A: Currently, RyCode uses one API key per provider. For multiple keys, use different projects or environment configurations.

**Q: How do I switch between providers during a session?**
A: Use the `/model` command in the TUI to switch between any configured provider's models.

## Contributing

Found an issue with provider authentication? Please report it:
- GitHub Issues: https://github.com/aaronmrosenthal/RyCode/issues
- Include: Provider name, error message, and steps to reproduce
