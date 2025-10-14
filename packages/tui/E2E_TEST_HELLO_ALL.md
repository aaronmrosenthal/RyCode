# E2E Test: Hello All Providers

## Overview

This is a **core build unit test** that validates all authenticated SOTA (State-of-the-Art) AI providers respond correctly to messages. It simulates the real-world workflow of:

1. User opens model selector (`/model`)
2. User presses Tab to cycle through providers (Claude → Tab → Codex → Tab → Gemini → Tab → Qwen)
3. User sends "hello" message to each provider
4. System receives valid responses from ALL providers

## Why This Test is Critical

This test ensures:
- ✅ All CLI provider integrations work end-to-end
- ✅ Authentication is working for all SOTA models
- ✅ API sessions can be created for each provider
- ✅ Message routing works correctly
- ✅ Response streaming/handling works
- ✅ Session cleanup works properly

**If this test fails, users cannot use RyCode with all providers.**

## Test Architecture

### Files

- **`test_hello_all_providers_e2e.go`** - Go test that:
  - Loads all CLI providers via AuthBridge
  - Checks authentication status for each provider
  - Creates API session for each authenticated provider
  - Sends "hello" message
  - Validates response (non-empty)
  - Cleans up session

- **`test-hello-all-providers.sh`** - Bash runner script that:
  - Checks API server is running (port 4096)
  - Compiles the Go test
  - Runs the test
  - Reports results
  - Cleans up temporary files

- **`Makefile`** - Build integration:
  - `make test` - Runs all tests including this E2E test
  - `make test-hello-all` - Runs only this E2E test

## Running the Test

### Prerequisites

1. **API server must be running**:
   ```bash
   cd packages/rycode
   bun run dev
   ```
   The API should be running on `http://127.0.0.1:4096`

2. **At least one provider must be authenticated**:
   ```bash
   ./bin/rycode /auth
   # Or set API keys in environment:
   export ANTHROPIC_API_KEY="sk-ant-..."
   export OPENAI_API_KEY="sk-..."
   export GOOGLE_API_KEY="..."
   export XAI_API_KEY="..."
   ```

### Run the Test

```bash
# Option 1: Via Makefile (recommended)
cd packages/tui
make test-hello-all

# Option 2: Via bash script directly
cd packages/tui
./test-hello-all-providers.sh

# Option 3: Run as part of full test suite
cd packages/tui
make test
```

## Expected Output

### Successful Run (All Providers Pass)

```
=== Core Build Unit Test: Hello All Providers ===

Project root: /Users/aaron/Code/RyCode/RyCode
Test script: test_hello_all_providers_e2e.go

[1] Checking if API server is running on port 4096...
    ✓ API server is running

[2] Compiling test...
    ✓ Test compiled to /tmp/test_hello_all_providers_e2e

[3] Running test...

=== E2E Test: Hello to All SOTA Providers ===

=== STARTING HELLO ALL PROVIDERS E2E TEST ===
Time: 2025-10-13T18:30:00Z
Purpose: Validate ALL SOTA models respond to messages

[1] Creating app instance...
    ✓ App created with auth bridge

[2] Loading CLI providers...
    ✓ Found 4 CLI provider configs

[3] Checking authentication status...
    - Claude: ✓ AUTHENTICATED (6 models, default: claude-sonnet-4-5)
    - Codex: ✓ AUTHENTICATED (8 models, default: gpt-5)
    - Gemini: ✓ AUTHENTICATED (7 models, default: gemini-2.5-pro)
    - Qwen: ✓ AUTHENTICATED (7 models, default: qwen3-max)

    Total authenticated providers: 4
    Providers to test: Claude, Codex, Gemini, Qwen

[4] Testing message responses from each provider...
    Test message: "hello"

  [1/4] Testing Claude (model: claude-sonnet-4-5)...
      ✓ Session created: sess_abc123
      ✓ SUCCESS: Got response (156 chars)
      Response preview: Hello! I'm Claude, an AI assistant created by Anthropic. How can I help you today?
      ✓ Session cleaned up

  [2/4] Testing Codex (model: gpt-5)...
      ✓ Session created: sess_def456
      ✓ SUCCESS: Got response (89 chars)
      Response preview: Hello! How can I assist you today?
      ✓ Session cleaned up

  [3/4] Testing Gemini (model: gemini-2.5-pro)...
      ✓ Session created: sess_ghi789
      ✓ SUCCESS: Got response (134 chars)
      Response preview: Hello! It's nice to hear from you. What can I do for you today?
      ✓ Session cleaned up

  [4/4] Testing Qwen (model: qwen3-max)...
      ✓ Session created: sess_jkl012
      ✓ SUCCESS: Got response (98 chars)
      Response preview: 你好！很高兴为你服务。有什么我可以帮助你的吗？
      ✓ Session cleaned up

=== TEST SUMMARY ===
Total providers tested: 4
Passed: 4
Failed: 0

✓ Passed providers:
  - Claude
  - Codex
  - Gemini
  - Qwen

Test logs saved to: /tmp/rycode-e2e-hello-all.log

✅ TEST PASSED: All 4 providers responded successfully!
```

### Failure Example (Provider Not Responding)

```
  [3/4] Testing Gemini (model: gemini-2.5-pro)...
      ✓ Session created: sess_ghi789
      ✗ FAILED to send message: API error: rate limit exceeded
      ✓ Session cleaned up

=== TEST SUMMARY ===
Total providers tested: 4
Passed: 3
Failed: 1

✓ Passed providers:
  - Claude
  - Codex
  - Qwen

✗ Failed providers:
  - Gemini

❌ TEST FAILED: 1 provider(s) did not respond correctly
```

## Test Logs

Detailed logs are saved to: `/tmp/rycode-e2e-hello-all.log`

View logs:
```bash
cat /tmp/rycode-e2e-hello-all.log
```

## Integration with CI/CD

This test should be run:
- ✅ Before every build
- ✅ In CI/CD pipeline before merging to main
- ✅ Before creating releases
- ✅ After adding new providers

### GitHub Actions Example

```yaml
name: E2E Tests

on: [push, pull_request]

jobs:
  test-hello-all-providers:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'

      - name: Set up Bun
        uses: oven-sh/setup-bun@v1

      - name: Set API Keys
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
          OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
          GOOGLE_API_KEY: ${{ secrets.GOOGLE_API_KEY }}
          XAI_API_KEY: ${{ secrets.XAI_API_KEY }}
        run: |
          echo "API keys configured"

      - name: Start API Server
        run: |
          cd packages/rycode
          bun install
          bun run dev &
          sleep 5

      - name: Run Hello All Providers E2E Test
        run: |
          cd packages/tui
          make test-hello-all
```

## Troubleshooting

### Error: "API server is not running on port 4096"

**Solution**: Start the API server:
```bash
cd packages/rycode
bun run dev
```

### Error: "No authenticated CLI providers found"

**Solution**: Authenticate at least one provider:
```bash
# Option 1: Via TUI
./bin/rycode /auth

# Option 2: Via environment variables
export ANTHROPIC_API_KEY="sk-ant-..."
export OPENAI_API_KEY="sk-..."
```

### Error: "Failed to create session"

**Possible causes**:
1. API key is invalid or expired
2. Provider API is down
3. Rate limit exceeded
4. Network connectivity issues

**Solution**: Check provider status and authentication:
```bash
cd packages/rycode
bun run src/auth/cli.ts check claude
bun run src/auth/cli.ts check codex
```

### Error: "Empty response"

**Possible causes**:
1. Provider returned error
2. Streaming/response handling issue
3. Session was terminated prematurely

**Solution**: Check detailed logs at `/tmp/rycode-e2e-hello-all.log`

## Test Coverage

This test validates:

| Provider | Models Tested | Default Model |
|----------|---------------|---------------|
| Claude   | 6 models      | claude-sonnet-4-5 |
| Codex    | 8 models      | gpt-5 |
| Gemini   | 7 models      | gemini-2.5-pro |
| Qwen     | 7 models      | qwen3-max |

**Total**: 28 SOTA models across 4 providers

## Maintenance

### Adding New Providers

When adding a new provider, update:

1. **`test_hello_all_providers_e2e.go`**:
   - Add provider to `getProviderDisplayName()`
   - Add default model to `getDefaultModelForProvider()`

2. **This documentation**:
   - Update test coverage table
   - Update expected output examples

### Updating Model Priorities

When new SOTA models are released, update priorities in `getDefaultModelForProvider()`:

```go
priorities := map[string][]string{
    "claude": {
        "claude-sonnet-5",        // NEW: Latest model
        "claude-sonnet-4-5",      // Previous SOTA
        "claude-opus-4-1",
        // ...
    },
}
```

## Related Files

- `/packages/tui/internal/auth/bridge.go` - AuthBridge implementation
- `/packages/tui/internal/components/dialog/simple_provider_toggle.go` - Provider selector UI
- `/packages/rycode/src/auth/cli.ts` - CLI auth commands
- `/packages/rycode/src/auth/auth-manager.ts` - Auth management

## Contact

For issues or questions about this test, please file an issue at:
https://github.com/aaronmrosenthal/rycode/issues
