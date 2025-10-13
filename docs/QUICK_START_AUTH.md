# Quick Start: Adding Your API Keys to RyCode

## Recommended Method: Interactive CLI

The easiest way to add your API keys:

```bash
# Start the interactive setup
bun run packages/rycode/src/index.ts auth login
```

Then select your provider:
```
┌  Add credential
│
◆  Select provider
│  ● OpenAI          # For GPT, Codex
│  ● Anthropic       # For Claude (recommended)
│  ● Google          # For Gemini
│  ● GitHub Copilot  # If you have Copilot
│  ...
```

Enter your API key when prompted. Done!

## Verify Setup

```bash
# Check configured providers
bun run packages/rycode/src/index.ts auth list

# List available models
bun run packages/rycode/src/index.ts models
```

## Backup Method: Environment Variables

If you prefer environment variables or need them for CI/CD:

```bash
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."
export GOOGLE_API_KEY="..."
```

**Priority:** CLI auth is checked first, then environment variables.

## Get API Keys

- **OpenAI:** https://platform.openai.com/api-keys
- **Anthropic:** https://console.anthropic.com/settings/keys (recommended!)
- **Google:** https://makersuite.google.com/app/apikey

## Test It Works

```bash
# Test provider detection
bun run packages/rycode/test/provider-test.ts

# Test actual API calls (optional)
bun run packages/rycode/test/provider-e2e-test.ts
```

## Start Using RyCode

```bash
# Start server
bun run packages/rycode/src/index.ts serve --port 4096

# In another terminal
rycode
```

That's it! Your own API keys, your own costs, full control.
