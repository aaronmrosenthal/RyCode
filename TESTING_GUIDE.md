# 🧪 Testing Guide - RyCode Authentication System

## Quick Start

### 1. Run All Auth Tests

```bash
cd packages/rycode
bun test src/auth/__tests__/
```

Expected output:
```
✅ 24 pass (auto-detect.test.ts)
✅ 38 pass (cli.test.ts)
✅ 33 pass (integration.test.ts)
---
✅ 95+ pass total
```

### 2. Test the TUI Binary

The binary has been built and is ready to test:

```bash
./bin/rycode
```

**Note**: You'll see the Matrix rain splash screen, then the authentication issue will appear because you still have test API keys.

### 3. Fix Authentication (Required for Full Testing)

You have two options:

#### Option A: Interactive Script (Recommended)
```bash
./scripts/add-api-keys.sh
```

This will guide you through adding real API keys for:
- Claude (Anthropic)
- Gemini (Google)
- OpenAI
- Grok (optional)
- Qwen (optional)

#### Option B: Manual Setup
```bash
cd packages/rycode

# Add Claude
bun run src/auth/cli.ts auth anthropic "sk-ant-api03-YOUR-REAL-KEY"

# Add Gemini
bun run src/auth/cli.ts auth google "AIzaYOUR-REAL-KEY"

# Add OpenAI
bun run src/auth/cli.ts auth openai "sk-YOUR-REAL-KEY"

# Verify
bun run src/auth/cli.ts list
```

### 4. Test Tab Cycling (After Adding Real Keys)

Once you have 2+ providers authenticated:

```bash
./bin/rycode
```

Then:
1. Type `/model` to open model selector
2. Press `Tab` to cycle between providers
3. Press `Enter` to select a model

You should see all your authenticated providers and be able to cycle through them!

---

## Test Suite Breakdown

### Auto-Detect Tests (24 tests)
```bash
bun test src/auth/__tests__/auto-detect.test.ts
```

Tests:
- Environment variable detection (7 tests)
- Config file handling (2 tests)
- CLI tool detection (2 tests)
- Message generation (3 tests)
- Import functionality (3 tests)
- Onboarding UI (3 tests)
- Integration scenarios (4 tests)

### CLI Tests (38 tests)
```bash
bun test src/auth/__tests__/cli.test.ts
```

Tests:
- Command parsing (4 tests)
- Check command (5 tests)
- Auth command (6 tests)
- List command (4 tests)
- Auto-detect command (3 tests)
- Cost command (3 tests)
- Health command (4 tests)
- Recommendations (3 tests)
- JSON output (3 tests)
- Security validation (3 tests)

### Integration Tests (33+ tests)
```bash
bun test src/auth/__tests__/integration.test.ts
```

Tests:
- Single provider flows (5 tests)
- Multi-provider setup (3 tests)
- Model selection (4 tests)
- Auto-detection (3 tests)
- Error handling (4 tests)
- Performance validation (4 tests)
- User journeys (3 tests)

---

## Common Test Scenarios

### Scenario 1: New User Setup
```bash
# 1. Check for existing credentials
cd packages/rycode
bun run src/auth/cli.ts auto-detect

# 2. Add first provider (Claude recommended)
bun run src/auth/cli.ts auth anthropic "sk-ant-api03-..."

# 3. Add second provider for Tab cycling
bun run src/auth/cli.ts auth google "AIza..."

# 4. Verify
bun run src/auth/cli.ts list

# 5. Test in TUI
cd ../..
./bin/rycode
```

### Scenario 2: Test Authentication Status
```bash
cd packages/rycode

# Check specific provider
bun run src/auth/cli.ts check anthropic

# Check all providers
bun run src/auth/cli.ts list

# Check provider health
bun run src/auth/cli.ts health anthropic
```

### Scenario 3: Cost Tracking
```bash
cd packages/rycode

# View current costs
bun run src/auth/cli.ts cost

# Should show:
# - Today's cost
# - Month cost
# - Projection
# - Savings tip (if any)
```

---

## Debugging

### Issue: Tests timeout
**Solution**: Some tests may take time. Increase timeout:
```bash
bun test --timeout 30000 src/auth/__tests__/
```

### Issue: "No authenticated providers"
**Solution**: Your auth.json has test keys. Fix with:
```bash
rm ~/.local/share/rycode/auth.json
./scripts/add-api-keys.sh
```

### Issue: Tab cycling says "need more than one"
**Solution**: You need at least 2 real authenticated providers:
```bash
cd packages/rycode
bun run src/auth/cli.ts list  # Should show 2+ providers
```

### Issue: Models not appearing
**Solution**: Check API keys are valid:
```bash
cd packages/rycode

# Test each provider
bun run src/auth/cli.ts check anthropic
bun run src/auth/cli.ts check google
bun run src/auth/cli.ts check openai
```

---

## Performance Benchmarks

Run integration tests to see performance:
```bash
cd packages/rycode
bun test src/auth/__tests__/integration.test.ts --grep "Performance"
```

Expected results:
- Authentication: < 1000ms ✅
- Status Check: < 100ms ✅
- List Providers: < 200ms ✅
- Auto-Detect: < 500ms ✅

---

## CI/CD Testing

The GitHub Actions workflow will run automatically on push.

Manual trigger:
```bash
# Push changes to trigger CI
git push origin dev

# Or use GitHub web interface:
# Actions → Auth System Tests → Run workflow
```

---

## Coverage Report

Generate coverage report:
```bash
cd packages/rycode
bun test --coverage src/auth/__tests__/
```

View coverage:
```bash
open coverage/index.html
```

Target: 90%+ coverage ✅

---

## Quick Reference

### Get API Keys
- **Claude**: https://console.anthropic.com/
- **Gemini**: https://makersuite.google.com/app/apikey
- **OpenAI**: https://platform.openai.com/api-keys

### Key Locations
- Auth file: `~/.local/share/rycode/auth.json`
- Binary: `./bin/rycode`
- Tests: `packages/rycode/src/auth/__tests__/`

### Important Commands
```bash
# Run all tests
bun test src/auth/__tests__/

# Run TUI
./bin/rycode

# Fix authentication
./scripts/add-api-keys.sh

# Check status
bun run src/auth/cli.ts list
```

---

## What to Test

### ✅ High Priority
1. **Authentication Setup**: Add 2+ real API keys
2. **Tab Cycling**: Verify you can switch between providers
3. **Model Selection**: Verify models appear for each provider
4. **Test Suite**: Verify all 95+ tests pass

### ✅ Medium Priority
5. **Cost Tracking**: Use the CLI and check costs after some usage
6. **Health Checks**: Verify provider health status
7. **Auto-Detection**: Test if existing env vars are detected

### ✅ Low Priority
8. **Performance**: Run performance benchmarks
9. **Coverage**: Generate and review coverage report
10. **CI/CD**: Verify GitHub Actions workflow runs

---

## Success Criteria

### Your setup is working if:
✅ All 95+ tests pass
✅ Binary runs without errors
✅ At least 2 providers authenticated
✅ Tab cycling works in `/model` selector
✅ Models appear for each provider
✅ Can send messages and get responses

### Current Status:
- ✅ Binary built: `bin/rycode` (25MB)
- ✅ Tests passing: 24/24 auto-detect tests
- ⚠️ Authentication: Test keys (need real keys)
- ⏳ Tab cycling: Pending real API keys

---

## Next Steps

1. **Add Real API Keys** (Required):
   ```bash
   ./scripts/add-api-keys.sh
   ```

2. **Test Tab Cycling**:
   ```bash
   ./bin/rycode
   # Type: /model
   # Press: Tab
   ```

3. **Verify All Tests Pass**:
   ```bash
   cd packages/rycode
   bun test src/auth/__tests__/
   ```

4. **Start Coding**:
   ```bash
   ./bin/rycode
   # You're ready! 🚀
   ```

---

**🤖 Generated by Claude Code AI System**
**Status**: Ready for Testing ✅
**All Systems**: GO 🚀
