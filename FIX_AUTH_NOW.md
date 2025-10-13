# 🚨 Fix Authentication Issue NOW

## TL;DR - Your Problem

Your `~/.local/share/rycode/auth.json` has **TEST API keys**, not real ones:

```json
{
  "openai": { "apiKey": "sk-test-key-for-testing" },      // ❌ FAKE
  "anthropic": { "apiKey": "sk-ant-test-key-for-testing" },  // ❌ FAKE
  "google": { "apiKey": "test-google-key" }                 // ❌ FAKE
}
```

This is why models aren't showing up and Tab cycling says "need more than one to tab".

---

## ⚡ Quick Fix (2 Options)

### Option 1: Interactive Script (Easiest)

```bash
./scripts/add-api-keys.sh
```

The script will walk you through adding each provider's API key.

### Option 2: Manual Commands

```bash
cd packages/rycode

# Add Claude (REQUIRED - best coding model)
bun run src/auth/cli.ts auth anthropic "sk-ant-api03-YOUR-REAL-KEY"

# Add Gemini (REQUIRED - for Tab cycling to work)
bun run src/auth/cli.ts auth google "AIzaYOUR-REAL-KEY"

# Add OpenAI (Optional - GPT-4 / Codex)
bun run src/auth/cli.ts auth openai "sk-YOUR-REAL-KEY"

# Add Grok (Optional - xAI)
bun run src/auth/cli.ts auth grok "YOUR-GROK-KEY"

# Add Qwen (Optional - Alibaba)
bun run src/auth/cli.ts auth qwen "YOUR-QWEN-KEY"
```

---

## 🔑 Where to Get API Keys

| Provider | URL | Key Format |
|----------|-----|------------|
| **Claude** | https://console.anthropic.com/ | `sk-ant-api03-...` |
| **Gemini** | https://makersuite.google.com/app/apikey | `AIza...` |
| **OpenAI** | https://platform.openai.com/api-keys | `sk-...` or `sk-proj-...` |
| **Grok** | https://console.x.ai/ | varies |
| **Qwen** | https://dashscope.aliyun.com/ | varies |

---

## ✅ Verify It Works

```bash
# Check authentication status
cd packages/rycode
bun run src/auth/cli.ts list

# Should show something like:
# {
#   "providers": [
#     { "id": "anthropic", "modelsCount": 3 },
#     { "id": "google", "modelsCount": 4 }
#   ]
# }
```

---

## 🧪 Test Tab Cycling

```bash
# Run the TUI
cd packages/tui
go run cmd/rycode/main.go

# Then:
# 1. Type /model
# 2. Press Tab
# 3. You should cycle between your providers!
```

---

## 📚 Full Documentation

For more details, troubleshooting, and security features:
- **Full Guide**: [docs/AUTHENTICATION_FIX.md](docs/AUTHENTICATION_FIX.md)
- **OAuth Setup**: [docs/OAUTH_AUTHENTICATION.md](docs/OAUTH_AUTHENTICATION.md)
- **Quick Start**: [docs/QUICK_START_AUTH.md](docs/QUICK_START_AUTH.md)

---

## 🎯 Why You Need This

You mentioned you're "authed into Gemini, Qwen, Codex, and Claude" but the auth.json shows only test keys. This means:

1. You might have set up OAuth or environment variables elsewhere
2. But RyCode stores credentials in `~/.local/share/rycode/auth.json`
3. That file currently has fake test keys
4. You need to replace them with your real API keys

---

**Bottom line**: Run `./scripts/add-api-keys.sh` or add your real API keys manually, then test `/model` + Tab in the TUI!
