# Model Selector - Testing & UX Documentation

## Quick Start

### View the Interactive Demo
```bash
open test-model-selector-web.html
```

### Run Playwright Tests
```bash
# Install browsers (if needed)
bunx playwright install chromium

# Run all tests
bunx playwright test test-model-selector.spec.ts

# Run with UI
bunx playwright test test-model-selector.spec.ts --ui
```

### Test the Real TUI
```bash
# Run direct Go test (proves data layer works)
go run test_models_direct.go

# Run the actual TUI
../../bin/rycode
# Then press: Ctrl+X, m (to open model selector)
```

---

## What's Here

### Test Files
- **`test-model-selector-web.html`** - Interactive web visualization for Playwright testing
- **`test-model-selector.spec.ts`** - 26 comprehensive E2E tests
- **`test_models_direct.go`** - Direct Go test proving provider merging works

### Documentation
- **`../../docs/MODEL_SELECTOR_UX_ANALYSIS.md`** - Multi-agent UX analysis (Codex + Claude)
- **`../../PLAYWRIGHT_TEST_SUMMARY.md`** - Complete testing overview
- **`MODEL_SELECTOR_README.md`** - This file

---

## Test Results

### Direct Go Test Output
```
✅ Found 4 CLI providers (28 models)
✅ Found 1 API provider (2 models)
✅ ListProviders merged: 5 providers, 30 models total
```

### Playwright Test Coverage
- **26 tests** across 3 suites
- **Core functionality**: Provider detection, model loading, search, keyboard nav
- **Edge cases**: Empty results, locked providers, rapid interactions
- **Performance**: Load time, render speed, input responsiveness

---

## Key Features Tested

### ✅ Provider Detection
- All 5 providers displayed (Anthropic, OpenAI, Claude CLI, Qwen, Gemini)
- Authentication status (✓ for authenticated, 🔒 for locked)
- CLI providers distinguished with "CLI" badge

### ✅ Model Loading
- 30 models total across all providers
- Recent models section (3 most recent)
- Model metadata badges (⚡💰🔥🆕)

### ✅ Search & Filtering
- Fuzzy search through model names
- Search by provider name
- Empty results handling

### ✅ Keyboard Navigation
- `/` - Focus search
- `1-9` - Jump to provider
- `d` - Auto-detect credentials
- `?` - Show help
- `Tab` - Quick switch (in TUI)

### ✅ Authentication Flow
- Auto-detect CLI providers
- Inline auth prompts (web demo)
- Success/failure feedback

---

## UX Improvements Demonstrated

Based on multi-agent AI analysis (Codex + Claude), the web visualization shows:

### 1. Visual Hierarchy
- ✨ Persistent shortcut bar at top
- 📌 Recent models section
- 🗂️ Collapsible provider groups
- 🎨 Icon-based badges
- 💡 AI insights panel

### 2. Model Metadata
- ⚡ Speed indicator (fast vs reasoning)
- 💰 Cost tiers ($ to $$$$)
- 🔥 Popularity (top 10%)
- 🆕 Recency (< 30 days old)
- 📏 Context sizes, output limits

### 3. Accessibility
- ⌨️ All keyboard shortcuts documented
- 🔢 Number keys for provider jump
- ❓ Help overlay (`?` key)
- 🎯 Clear focus indicators

---

## Implementation Roadmap

### Phase 1: Critical Fixes (1-2 days)
1. Add persistent shortcut footer to TUI
2. Implement model metadata badges
3. Create collapsible provider groups

### Phase 2: Accessibility (1 day)
4. Number key navigation (1-9)
5. ARIA-equivalent labels
6. Help overlay (`?` key)

### Phase 3: Polish (2 days)
7. Inline authentication flow
8. Optimistic UI with progress
9. Search filters (provider:, cost:, speed:)

---

## Success Metrics

### Current State
- Time to select model: ~8 seconds
- Keyboard usage: ~30%
- Auth success rate: ~60%

### Target State (after improvements)
- Time to select model: < 3 seconds (**60% reduction**)
- Keyboard usage: 70% (**3x increase**)
- Auth success rate: 90% (**50% improvement**)

---

## Related Files

### Go Implementation
- `internal/app/app.go:1308-1373` - ListProviders() merging logic
- `internal/components/dialog/models.go` - Model dialog (957 lines)
- `internal/auth/bridge.go` - CLI provider bridge

### TypeScript Auth System
- `../../packages/rycode/src/auth/cli.ts` - CLI provider detection
- `../../packages/rycode/src/auth/cli-bridge.ts` - Bridge implementation

---

## Debugging

### Check Debug Logs
```bash
cat /tmp/rycode-debug.log
```

### Verify CLI Providers
```bash
bun run ../../packages/rycode/src/auth/cli.ts cli-providers
```

### Verify API Providers
```bash
bun run ../../packages/rycode/src/auth/cli.ts list
```

---

## Questions?

- See `../../docs/MODEL_SELECTOR_UX_ANALYSIS.md` for detailed UX recommendations
- See `../../PLAYWRIGHT_TEST_SUMMARY.md` for complete test overview
- See `../../E2E_PROOF.md` for code flow analysis
- See `../../MANUAL_TEST.md` for manual testing guide

---

**Ready to test?** Open `test-model-selector-web.html` and try the keyboard shortcuts!
