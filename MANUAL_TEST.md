# Manual Testing Guide

## Quick Test (30 seconds)

Run this to verify everything works:

```bash
# 1. Verify CLI providers are detected
./test-list-providers.sh

# Expected output:
# ✅ 7 providers (3 API + 4 CLI)
# ✅ 47 models (19 API + 28 CLI)

# 2. Run the TUI
./bin/rycode

# 3. Open model selector
#    Press: Ctrl+X then m
#    (or type /models and press Tab)

# 4. You should see:
#    - 7 provider sections
#    - All CLI providers (claude, qwen, codex, gemini) with ✓ unlocked
#    - 47 total models listed

# 5. Close modal (press Esc)

# 6. Test Tab cycling
#    Press: Tab repeatedly
#    Should see toasts like:
#    - "Switched to Claude (6 models)"
#    - "Switched to Qwen (7 models)"
#    - "Switched to Codex (8 models)"
#    - etc.

# 7. Quit (press Ctrl+C twice)
```

## What You'll See

### In Model Selector (`/models`)

```
┌─ Select Model ────────────────────────────────────┐
│                                                    │
│  Recent                                            │
│    └─ (your recently used models)                 │
│                                                    │
│  Anthropic ✓                                       │
│    └─ Claude 3.5 Sonnet                           │
│    └─ Claude 3.5 Haiku                            │
│    └─ ... (6 models total)                        │
│                                                    │
│  Claude ✓  ← CLI Provider!                        │
│    └─ claude-sonnet-4-5                           │
│    └─ claude-opus-4-1                             │
│    └─ ... (6 models total)                        │
│                                                    │
│  Codex ✓  ← CLI Provider!                         │
│    └─ gpt-5                                        │
│    └─ gpt-5-mini                                   │
│    └─ o3                                           │
│    └─ ... (8 models total)                        │
│                                                    │
│  Gemini ✓  ← CLI Provider!                        │
│    └─ gemini-2.5-pro                              │
│    └─ gemini-2.5-flash                            │
│    └─ ... (7 models total)                        │
│                                                    │
│  Google ✓                                          │
│    └─ Gemini 2.0 Flash                            │
│    └─ ... (7 models total)                        │
│                                                    │
│  OpenAI ✓                                          │
│    └─ GPT-4o                                       │
│    └─ o1 Preview                                   │
│    └─ ... (6 models total)                        │
│                                                    │
│  Qwen ✓  ← CLI Provider!                          │
│    └─ qwen3-max                                    │
│    └─ qwen3-next                                   │
│    └─ ... (7 models total)                        │
│                                                    │
│  [d] auto-detect  [a] authenticate  [i] insights  │
└────────────────────────────────────────────────────┘
```

### Tab Cycling

When you press Tab:
```
[Toast] Switched to Claude (6 models) ✓
```

Press Tab again:
```
[Toast] Switched to Qwen (7 models) ✓
```

Press Tab again:
```
[Toast] Switched to Codex (8 models) ✓
```

And so on, cycling through all authenticated providers!

## Troubleshooting

### If CLI providers don't show:

1. Make sure the CLIs are actually running:
   ```bash
   # In separate terminals, run:
   claude
   qwen
   codex
   gemini
   ```

2. Check CLI detection:
   ```bash
   bun run packages/rycode/src/auth/cli.ts cli-providers
   ```

3. Check debug log:
   ```bash
   cat /tmp/rycode-debug.log
   ```

### If Tab cycling doesn't work:

1. Make sure you have authenticated providers:
   ```bash
   bun run packages/rycode/src/auth/cli.ts list
   ```

2. Try pressing `d` in the model selector to auto-detect

3. Check if you're in a modal - Tab only cycles when NOT in a modal

## Success Criteria

✅ You see 7 providers in `/models`
✅ All 4 CLI providers (claude, qwen, codex, gemini) show ✓ unlocked
✅ Total of 47 models across all providers
✅ Tab key cycles through authenticated providers
✅ Toasts show "Switched to [Provider] (X models)"

## Files Created for Testing

- `E2E_PROOF.md` - Complete technical proof with code snippets
- `test-list-providers.sh` - Verify provider data
- `test-tui-expect.exp` - Automated test (shows debug logs)
- `MANUAL_TEST.md` - This file

All debug logging is in place - check `/tmp/rycode-debug.log` after running the TUI.
