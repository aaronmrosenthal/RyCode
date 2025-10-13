# Authentication Issue Fix Guide

## 🔍 Problem Diagnosed

Your `~/.local/share/rycode/auth.json` contains **TEST API keys**, not real credentials:

```json
{
  "openai": {
    "type": "api",
    "apiKey": "sk-test-key-for-testing"
  },
  "anthropic": {
    "type": "api",
    "apiKey": "sk-ant-test-key-for-testing"
  },
  "google": {
    "type": "api",
    "apiKey": "test-google-key"
  }
}
```

These test keys are being recognized as "authenticated" by the auth system, but they **cannot make real API calls**. This is why:
- Models aren't appearing properly
- Tab cycling shows "need more than one to tab"
- Providers appear authenticated but don't work

## ✅ Solution: Add Real API Keys

You need to obtain **real API keys** from each provider and add them to RyCode.

---

## 📋 Step-by-Step Instructions

### 1. **Claude (Anthropic)**

#### Get Your API Key:
1. Go to https://console.anthropic.com/
2. Sign in or create an account
3. Navigate to **API Keys** section
4. Click **Create Key**
5. Copy your key (starts with `sk-ant-api03-`)

#### Add to RyCode:
```bash
# From the RyCode project root:
cd packages/rycode
bun run src/auth/cli.ts auth anthropic "sk-ant-api03-YOUR-KEY-HERE"
```

---

### 2. **Gemini (Google)**

#### Get Your API Key:
1. Go to https://makersuite.google.com/app/apikey
2. Sign in with your Google account
3. Click **Create API Key**
4. Copy your key (starts with `AIza`)

#### Add to RyCode:
```bash
cd packages/rycode
bun run src/auth/cli.ts auth google "AIzaYOUR-KEY-HERE"
```

---

### 3. **OpenAI (GPT/Codex)**

#### Get Your API Key:
1. Go to https://platform.openai.com/api-keys
2. Sign in or create an account
3. Click **Create new secret key**
4. Give it a name and copy the key (starts with `sk-proj-` or `sk-`)

#### Add to RyCode:
```bash
cd packages/rycode
bun run src/auth/cli.ts auth openai "sk-YOUR-KEY-HERE"
```

---

### 4. **Grok (xAI)** *(Optional)*

#### Get Your API Key:
1. Go to https://console.x.ai/
2. Sign in or create an account
3. Navigate to **API Keys**
4. Create a new key
5. Copy your key

#### Add to RyCode:
```bash
cd packages/rycode
bun run src/auth/cli.ts auth grok "YOUR-GROK-KEY-HERE"
```

---

### 5. **Qwen (Alibaba Cloud)** *(Optional)*

#### Get Your API Key:
1. Go to https://dashscope.aliyun.com/
2. Sign in or create an account
3. Navigate to **API Keys**
4. Create a new key
5. Copy your key

#### Add to RyCode:
```bash
cd packages/rycode
bun run src/auth/cli.ts auth qwen "YOUR-QWEN-KEY-HERE"
```

---

## 🔄 After Adding Real Keys

### Verify Authentication:
```bash
cd packages/rycode
bun run src/auth/cli.ts list
```

You should see output like:
```json
{
  "providers": [
    {
      "id": "anthropic",
      "name": "anthropic",
      "modelsCount": 3
    },
    {
      "id": "google",
      "name": "google",
      "modelsCount": 4
    },
    {
      "id": "openai",
      "name": "openai",
      "modelsCount": 5
    }
  ]
}
```

### Test Tab Cycling:
1. Run RyCode TUI:
   ```bash
   cd packages/tui
   go run cmd/rycode/main.go
   ```

2. Type `/model` to open model selector

3. Press `Tab` to cycle between authenticated providers

4. You should see:
   - Multiple providers available
   - Ability to cycle through them with Tab
   - Models listed for each provider

---

## 🧪 Auto-Detection (Alternative Method)

If you already have API keys in your environment variables or standard locations, RyCode can auto-detect them:

```bash
cd packages/rycode
bun run src/auth/cli.ts auto-detect
```

This will search for:
- `ANTHROPIC_API_KEY` environment variable
- `OPENAI_API_KEY` environment variable
- `GOOGLE_API_KEY` environment variable
- `GROK_API_KEY` environment variable
- API keys in `~/.anthropic/` directory
- API keys in `~/.config/openai/` directory
- And more standard locations

---

## 📍 Where API Keys Are Stored

RyCode stores encrypted credentials in:
```
~/.local/share/rycode/auth.json
```

**Security Features:**
- ✅ Encrypted with `RYCODE_ENCRYPTION_KEY` (if set)
- ✅ Integrity verification
- ✅ File permissions set to `0600` (owner read/write only)
- ✅ Audit logging of all authentication events

---

## 🔐 Enable Encryption (Recommended)

For extra security, set an encryption key:

```bash
# Add to your ~/.bashrc or ~/.zshrc:
export RYCODE_ENCRYPTION_KEY="your-secure-random-key-here"
```

Generate a secure key:
```bash
# On macOS/Linux:
openssl rand -base64 32
```

Then migrate existing credentials to encrypted storage:
```bash
cd packages/rycode
bun run -e "import { authManager } from './src/auth/auth-manager'; await authManager.export(); console.log('Migrated to encrypted storage')"
```

---

## 🐛 Troubleshooting

### Issue: "No authenticated providers"
**Solution**: Run `bun run src/auth/cli.ts list` to verify credentials are stored

### Issue: "Only one provider authenticated"
**Solution**: Add at least 2 real API keys (see steps above)

### Issue: Tab cycling still not working
**Solution**:
1. Delete the test keys file:
   ```bash
   rm ~/.local/share/rycode/auth.json
   ```
2. Re-add your real API keys using the commands above

### Issue: Models still not showing
**Solution**: Check provider health:
```bash
cd packages/rycode
bun run src/auth/cli.ts health anthropic
bun run src/auth/cli.ts health google
bun run src/auth/cli.ts health openai
```

---

## 💰 Cost Tracking

RyCode tracks your API usage and costs:

```bash
cd packages/rycode
bun run src/auth/cli.ts cost
```

Output:
```json
{
  "todayCost": 0.15,
  "monthCost": 3.42,
  "projection": 15.20,
  "savingsTip": "Consider using Claude Haiku for simple tasks to reduce costs"
}
```

---

## 🎯 Next Steps

1. **Get API keys** from providers you want to use (at least 2)
2. **Add them to RyCode** using the `auth` CLI commands
3. **Verify** with `list` command
4. **Test Tab cycling** in the TUI
5. **Optional**: Enable encryption for extra security

---

## 📚 Related Documentation

- [OAUTH_AUTHENTICATION.md](./OAUTH_AUTHENTICATION.md) - OAuth setup guide
- [DEVELOPER_API_KEYS.md](./DEVELOPER_API_KEYS.md) - Developer API key management
- [AUTH_PRIORITY.md](./AUTH_PRIORITY.md) - Authentication priority system
- [QUICK_START_AUTH.md](./QUICK_START_AUTH.md) - Quick authentication setup

---

## ✨ Why You Need Multiple Providers

Having multiple authenticated providers allows you to:
- **Tab cycle** between providers quickly
- **Compare responses** from different models
- **Fallback** if one provider is down
- **Optimize costs** by choosing cheaper models for simple tasks
- **Leverage strengths** of different models (Claude for code, GPT-4 for reasoning, etc.)

---

**Status**: Ready to fix - follow the steps above to add your real API keys!
