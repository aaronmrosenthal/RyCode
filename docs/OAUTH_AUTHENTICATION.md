# OAuth Authentication in RyCode

## Overview

RyCode supports **OAuth authentication** so you don't need to manually create or manage API keys. Just **"Login with Claude"**, **"Login with GitHub"**, etc.

## Supported OAuth Providers

### ✅ Anthropic (Claude Pro/Max)

**Best for:** Claude users with Pro or Max subscription

```bash
bun run packages/rycode/src/index.ts auth login
```

1. Select **Anthropic**
2. Choose **Claude Pro/Max** login method
3. Browser opens automatically
4. Sign in with your Anthropic account
5. Done! All Claude models available

**Benefits:**
- No API key management
- Uses your existing Claude Pro/Max subscription
- Automatic access to all models you have access to

### ✅ GitHub Copilot

**Best for:** Developers with GitHub Copilot subscription

```bash
bun run packages/rycode/src/index.ts auth login
```

1. Select **GitHub Copilot**
2. You'll see: "Please visit: https://github.com/login/device"
3. Enter the code shown (e.g., "8F43-6FCF")
4. Authorize in your browser
5. Done! Access to GPT-4, Claude, and other models via Copilot

**Benefits:**
- No API key needed
- Use your existing Copilot subscription
- Access to multiple model providers through one login

### ✅ Google (OAuth or API Key)

**OAuth method:**

```bash
bun run packages/rycode/src/index.ts auth login
```

1. Select **Google**
2. Choose OAuth login method
3. Sign in with your Google account
4. Grant permissions
5. Done! Access to Gemini models

**Or use CLI:**
```bash
gcloud auth application-default login
export GOOGLE_CLOUD_PROJECT="your-project"
```

## How OAuth Works

### Authentication Flow

```
1. User runs: rycode auth login
2. RyCode starts OAuth flow
3. Browser opens (or device code shown)
4. User signs in to provider
5. Provider sends tokens to RyCode
6. RyCode stores tokens securely in ~/.local/share/rycode/auth.json
7. Tokens auto-refresh when needed
```

### Token Storage

Tokens are stored in `~/.local/share/rycode/auth.json`:

```json
{
  "anthropic": {
    "type": "oauth",
    "access": "...",
    "refresh": "...",
    "expires": "2025-10-12T12:00:00Z"
  },
  "github-copilot": {
    "type": "oauth",
    "access": "...",
    "refresh": "...",
    "expires": "2025-10-12T13:00:00Z"
  }
}
```

### Token Refresh

RyCode automatically refreshes expired tokens. You'll never need to re-authenticate unless:
- Your refresh token expires (usually 30-90 days)
- You explicitly log out
- You revoke access from the provider's dashboard

## Quick Start Guide

### For Claude Users

```bash
# 1. Run auth
bun run packages/rycode/src/index.ts auth login

# 2. Select Anthropic > Claude Pro/Max
# 3. Sign in via browser

# 4. Start RyCode
bun run packages/rycode/src/index.ts serve --port 4096

# 5. Launch TUI
rycode

# 6. Use /models to select Claude models
```

### For GitHub Copilot Users

```bash
# 1. Run auth
bun run packages/rycode/src/index.ts auth login

# 2. Select GitHub Copilot
# 3. Go to github.com/login/device
# 4. Enter the code shown
# 5. Authorize

# 6. Start RyCode
bun run packages/rycode/src/index.ts serve --port 4096

# 7. Launch TUI
rycode

# 8. Access GPT-4, Claude, and more via Copilot
```

### For Google/Gemini Users

**Option 1: OAuth**
```bash
bun run packages/rycode/src/index.ts auth login
# Select Google > OAuth
```

**Option 2: CLI (easier)**
```bash
gcloud auth application-default login
export GOOGLE_CLOUD_PROJECT="your-project-id"
```

## Testing OAuth Setup

```bash
# Check configured providers
bun run packages/rycode/src/index.ts auth list

# Should show:
# ┌  Credentials ~/.local/share/rycode/auth.json
# │
# ├ Anthropic oauth
# ├ GitHub Copilot oauth
# │
# └  2 credentials
```

## OAuth vs API Keys

| Feature | OAuth | API Keys |
|---------|-------|----------|
| Setup | Click & sign in | Copy/paste key |
| Security | More secure (tokens expire) | Less secure (keys don't expire) |
| Rotation | Automatic | Manual |
| Revocation | From provider dashboard | Delete key |
| Cost Control | Via subscription | Via API billing |
| Best For | Personal use, subscriptions | CI/CD, automation |

## When to Use What

### Use OAuth When:
- ✅ You have Claude Pro/Max subscription
- ✅ You have GitHub Copilot subscription
- ✅ You want easy "Login with..." experience
- ✅ You don't want to manage API keys
- ✅ You want automatic token rotation

### Use API Keys When:
- ✅ Running in CI/CD
- ✅ Automation/scripts
- ✅ Pay-per-use billing model
- ✅ Need fine-grained access control
- ✅ Provider doesn't support OAuth

## FAQ

**Q: Do I need Claude Pro to use OAuth?**
A: Yes, OAuth currently works with Claude Pro/Max subscriptions. For API key access, you can create keys at console.anthropic.com.

**Q: Can I use both OAuth and API keys?**
A: Yes! You can OAuth into Claude and use API keys for OpenAI, for example.

**Q: How do I revoke access?**
A: Either run `rycode auth logout` or revoke from the provider's dashboard:
- Anthropic: https://console.anthropic.com/settings/tokens
- GitHub: https://github.com/settings/apps/authorizations
- Google: https://myaccount.google.com/permissions

**Q: What happens when my token expires?**
A: RyCode automatically refreshes it using the refresh token. You won't notice anything.

**Q: Can I share my OAuth tokens?**
A: No, OAuth tokens are personal and tied to your account. For teams, each developer should authenticate separately.

**Q: Does OAuth work for OpenAI?**
A: Not currently. OpenAI doesn't provide OAuth for desktop apps. Use API keys instead.

**Q: Can I use OAuth with local Qwen models?**
A: Local models don't need authentication - just configure the endpoint in opencode.json.

## Troubleshooting

### Browser doesn't open

```bash
# Manually copy the URL shown and open it
# Or use device code flow if available
```

### "Authorization failed"

1. Check your subscription is active
2. Try logging out and back in:
   ```bash
   bun run packages/rycode/src/index.ts auth logout
   bun run packages/rycode/src/index.ts auth login
   ```

### "Token expired" errors

```bash
# RyCode should auto-refresh, but you can manually re-auth:
bun run packages/rycode/src/index.ts auth logout
bun run packages/rycode/src/index.ts auth login
```

### Can't see models after OAuth

```bash
# Verify authentication
bun run packages/rycode/src/index.ts auth list

# Check available models
bun run packages/rycode/src/index.ts models

# Restart server
pkill -f "rycode.*serve"
bun run packages/rycode/src/index.ts serve --port 4096
```

## Implementation Details

OAuth is implemented via plugins in RyCode:

- `opencode-anthropic-auth` - Anthropic OAuth plugin
- `opencode-copilot-auth` - GitHub Copilot OAuth plugin

These plugins are loaded automatically when you run `rycode auth login`.

The OAuth flow code is in `packages/rycode/src/cli/cmd/auth.ts` (lines 142-226).

## Security

✅ Tokens stored securely in `~/.local/share/rycode/auth.json`
✅ File permissions set to 600 (user read/write only)
✅ Tokens auto-refresh before expiry
✅ CSRF protection for OAuth flows
✅ No plaintext passwords stored
✅ Tokens never logged

## Summary

**For the best developer experience:**
1. Use OAuth for Claude (if you have Pro/Max)
2. Use OAuth for GitHub Copilot (if you have subscription)
3. Use API keys or environment variables for other providers
4. Use local models (Ollama/Qwen) when possible (free!)

**You don't need to manage API keys - just click "Login" and you're done!** ✨
