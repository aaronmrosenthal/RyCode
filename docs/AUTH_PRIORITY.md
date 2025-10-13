# RyCode Authentication Priority

## Loading Order

RyCode loads credentials in this order (later sources override earlier ones):

1. **Environment Variables** (backup)
2. **CLI Auth** (`~/.local/share/rycode/auth.json`) ✅ **DEFAULT/PRIMARY**
3. Custom Provider Loaders
4. Plugin Auth
5. Config File (`opencode.json`)

## Priority Example

```bash
# Scenario: Both environment and CLI auth configured

# 1. Environment variable set
export OPENAI_API_KEY="sk-env-key"

# 2. CLI auth also configured
bun run packages/rycode/src/index.ts auth login
# (adds key to ~/.local/share/rycode/auth.json)

# Result: CLI auth key is used (takes priority)
```

## Recommended Setup Method

### Primary: CLI Auth (Recommended) ✅

```bash
# Interactive setup - stores in ~/.local/share/rycode/auth.json
bun run packages/rycode/src/index.ts auth login

# Select provider, enter API key
# Stored securely, persists across sessions
```

**Benefits:**
- ✅ Simple interactive setup
- ✅ Persists across terminal sessions
- ✅ Easy to manage (`auth list`, `auth logout`)
- ✅ Secure storage in user's home directory

### Backup: Environment Variables

```bash
# For CI/CD or temporary use
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."

# Or in .env file
echo "OPENAI_API_KEY=sk-..." > .env
```

**Benefits:**
- ✅ Good for CI/CD pipelines
- ✅ Temporary sessions
- ✅ Per-project configuration
- ✅ Works without running auth setup

## Use Cases

### For Regular Development
**Use CLI Auth:**
```bash
rycode auth login
# One-time setup, works everywhere
```

### For CI/CD
**Use Environment Variables:**
```yaml
# GitHub Actions
env:
  ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
```

### For Multiple Projects
**Use .env per project:**
```bash
# project-1/.env
OPENAI_API_KEY=sk-project1-key

# project-2/.env
OPENAI_API_KEY=sk-project2-key
```

## Verification

Check what RyCode detects:

```bash
# Shows CLI auth credentials
bun run packages/rycode/src/index.ts auth list

# Shows environment variables detected
bun run packages/rycode/test/provider-test.ts
```

## Source Code

The loading logic is in `packages/rycode/src/provider/provider.ts`:

```typescript
// Line 292-303: Load from environment
for (const [providerID, provider] of Object.entries(database)) {
  const apiKey = provider.env.map((item) => process.env[item]).at(0)
  if (!apiKey) continue
  mergeProvider(providerID, { apiKey }, "env")
}

// Line 305-311: Load from CLI auth (OVERRIDES env)
for (const [providerID, provider] of Object.entries(await Auth.all())) {
  if (provider.type === "api") {
    mergeProvider(providerID, { apiKey: provider.key }, "api")
  }
}
```

The `mergeProvider` function on line 233 overwrites the source with each merge, so **CLI auth takes priority**.

## Summary

✅ **CLI auth is the primary/default method** (recommended for users)
✅ **Environment variables work as backup** (good for CI/CD)
✅ **CLI auth overrides environment variables** when both are present
✅ **Easy to manage with `auth` commands**

Users should be guided to use `rycode auth login` first, with environment variables documented as an alternative for CI/CD or temporary use.
