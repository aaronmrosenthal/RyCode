# CLI Provider Testing Guide

## Overview

This guide explains how to test RyCode's native CLI-authenticated providers (OpenAI Codex, Google Gemini, and Anthropic Claude) without making actual API calls.

## Test Scripts

### 1. `provider-cli-test.ts`
Tests CLI authentication status via `~/.local/share/rycode/auth.json`

**What it tests**:
- Authentication configuration
- Provider availability
- Model counts
- Auth types (API key vs OAuth)

**Usage**:
```bash
bun run packages/rycode/test/provider-cli-test.ts
```

**Output**:
- Lists authenticated providers
- Shows auth type (api/oauth)
- Displays masked API keys
- Reports model availability

### 2. `provider-test.ts`
Tests auto-detection from environment variables

**What it tests**:
- Environment variable detection
- CLI tool configuration
- Config file detection
- Model availability

**Usage**:
```bash
# With environment variables
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."
bun run packages/rycode/test/provider-test.ts
```

### 3. `provider-server-test.ts`
Tests provider API through running RyCode server

**What it tests**:
- Server connectivity
- Provider API endpoints
- Model selection
- Authentication status via API

**Prerequisites**:
```bash
# Start RyCode server first
bun run packages/rycode/src/index.ts serve --port 4096
```

**Usage**:
```bash
bun run packages/rycode/test/provider-server-test.ts
```

### 4. `run-provider-tests.sh`
Runs all provider tests in sequence

**What it does**:
1. Tests CLI authentication
2. Tests environment variables
3. Tests server API (if running)
4. Generates comprehensive report

**Usage**:
```bash
./packages/rycode/test/run-provider-tests.sh
```

## Authentication Methods

### Method 1: Interactive CLI Auth (Recommended)

```bash
# Run interactive authentication
rycode auth login

# Or if rycode CLI is not installed:
bun run packages/rycode/src/index.ts auth login

# Follow prompts:
# 1. Select provider (OpenAI, Anthropic, Google)
# 2. Choose auth method (API Key or OAuth)
# 3. Enter credentials
```

**Stored in**: `~/.local/share/rycode/auth.json`

**Format**:
```json
{
  "openai": {
    "type": "api",
    "key": "sk-proj-..."
  },
  "anthropic": {
    "type": "oauth",
    "access": "...",
    "refresh": "...",
    "expires": "2025-10-12T12:00:00Z"
  }
}
```

### Method 2: Environment Variables

```bash
# Set environment variables
export OPENAI_API_KEY="sk-proj-..."
export ANTHROPIC_API_KEY="sk-ant-..."
export GOOGLE_API_KEY="..."

# Or use .env file in project root
cat > .env << 'EOF'
OPENAI_API_KEY=sk-proj-...
ANTHROPIC_API_KEY=sk-ant-...
GOOGLE_API_KEY=...
EOF
```

### Method 3: Config File

Add to `opencode.json` in project root:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "provider": {
    "openai": {
      "apiKey": "sk-proj-..."
    },
    "anthropic": {
      "apiKey": "sk-ant-..."
    }
  }
}
```

## Testing Workflow

### Step 1: Authenticate Providers

Choose one method above to authenticate. For testing, we recommend interactive CLI auth:

```bash
# Authenticate OpenAI
rycode auth login
# Select: OpenAI
# Enter API key: sk-proj-...

# Authenticate Anthropic
rycode auth login
# Select: Anthropic
# Enter API key: sk-ant-...

# Authenticate Google
rycode auth login
# Select: Google
# Enter API key: ...
```

### Step 2: Run Tests

```bash
# Run all tests
./packages/rycode/test/run-provider-tests.sh

# Or run individual tests
bun run packages/rycode/test/provider-cli-test.ts
```

### Step 3: Verify Results

Expected output for authenticated providers:
```
✓ OpenAI Codex (openai)
  Auth Type: api
  API Key: sk-proj-...xxxx

  OpenAI Codex:
  ✅ Ready to use
  Auth: api
  Models: 25 available
```

## Current Test Results

Based on the most recent test run:

**CLI Authentication**: ❌ 0/3 authenticated
- No `auth.json` file found
- No providers configured via CLI

**Environment Variables**: ❌ 0/4 detected
- No environment variables set
- OPENAI_API_KEY not found
- ANTHROPIC_API_KEY not found
- GOOGLE_API_KEY not found

**Action Required**:
You mentioned having Codex, Gemini, and Claude authenticated locally. Please:

1. **Check if you used a different authentication method**:
   ```bash
   # Check for auth.json
   ls -la ~/.local/share/rycode/auth.json

   # Check environment variables
   env | grep -i "api_key"

   # Check config file
   cat opencode.json
   ```

2. **Re-authenticate if needed**:
   ```bash
   rycode auth login
   ```

3. **Verify with the TUI**:
   ```bash
   # Start server
   bun run packages/rycode/src/index.ts serve --port 4096

   # Start TUI
   rycode

   # Press /model to see available models
   ```

## Troubleshooting

### No Providers Authenticated

**Symptom**: All tests show "Not authenticated"

**Solutions**:
1. Run `rycode auth login` to authenticate
2. Set environment variables (see Method 2 above)
3. Add to `opencode.json` (see Method 3 above)

### Authentication File Not Found

**Symptom**: "Auth file not found"

**Solutions**:
```bash
# Create auth directory
mkdir -p ~/.local/share/rycode

# Run authentication
rycode auth login
```

### Server Not Running

**Symptom**: "Cannot connect to server"

**Solutions**:
```bash
# Start server
bun run packages/rycode/src/index.ts serve --port 4096

# Run tests again
./packages/rycode/test/run-provider-tests.sh
```

### Models Not Available

**Symptom**: "0 models" or "Provider not found in models"

**Solutions**:
1. Restart the server after authentication
2. Check provider ID matches exactly (openai, anthropic, google)
3. Verify API key is valid

## No API Calls Made

**Important**: These tests do NOT make actual API calls to providers.

They only verify:
- ✅ Authentication configuration
- ✅ Provider availability
- ✅ Model lists
- ✅ Configuration validity

To test actual API calls, use:
```bash
bun run packages/rycode/test/provider-e2e-test.ts
```

⚠️ **Warning**: E2E tests make real API calls and will consume tokens/credits.

## Next Steps

1. **Authenticate your providers** using one of the methods above
2. **Run the test suite** to verify authentication
3. **Start the RyCode server** to test via API
4. **Use the TUI** to verify models are available

## Provider IDs

| Provider | ID | Auth Methods |
|----------|----|--------------|
| OpenAI Codex | `openai` | API Key, OAuth |
| Google Gemini | `google` | API Key, OAuth, gcloud CLI |
| Anthropic Claude | `anthropic` | API Key, OAuth |
| GitHub Copilot | `github-copilot` | OAuth (device flow) |
| Qwen (local) | `ollama` | No auth needed |

## Example: Full Authentication Flow

```bash
# Step 1: Authenticate OpenAI
rycode auth login
# > Select provider: OpenAI
# > Auth method: API Key
# > Enter key: sk-proj-...
# ✓ Authenticated successfully

# Step 2: Authenticate Anthropic
rycode auth login
# > Select provider: Anthropic
# > Auth method: OAuth
# > Browser opens, sign in
# ✓ Authenticated successfully

# Step 3: Verify
bun run packages/rycode/test/provider-cli-test.ts
# ✓ OpenAI Codex (openai) - 25 models
# ✓ Anthropic Claude (anthropic) - 8 models

# Step 4: Start server
bun run packages/rycode/src/index.ts serve --port 4096

# Step 5: Test server API
bun run packages/rycode/test/provider-server-test.ts
# ✅ All target providers authenticated and ready!
```

## CI/CD Integration

For continuous integration:

```yaml
# .github/workflows/test.yml
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: oven-sh/setup-bun@v1

      - name: Test provider detection
        env:
          OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
        run: |
          bun install
          ./packages/rycode/test/run-provider-tests.sh
```

## Summary

✅ **Tests Created**:
- `provider-cli-test.ts` - CLI auth testing
- `provider-server-test.ts` - Server API testing
- `run-provider-tests.sh` - Complete test suite

✅ **No API Calls**: Tests verify configuration only

✅ **Multiple Auth Methods**: CLI, environment, config file

❌ **Action Needed**: Authenticate your providers to enable testing

Run this to get started:
```bash
rycode auth login
```
