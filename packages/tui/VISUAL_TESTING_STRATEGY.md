# Visual Testing Strategy for Provider Themes

**Status**: Design Document
**Phase**: 3.3 - Visual Regression Testing
**Date**: October 14, 2025

---

## Overview

Visual regression testing for RyCode TUI themes requires a different approach than web applications. Since we're rendering in a terminal, we need terminal-specific capture and comparison tools.

---

## Testing Approach

### 1. Terminal Screenshot Capture

**Tools Available:**
- **VHS** (by Charm) - Terminal recorder/screenshot tool
- **ttyrec** - Terminal session recorder
- **asciinema** - Terminal session recorder with playback
- **Script/Scriptreplay** - Built-in Unix terminal recording

**Recommended: VHS (Charm)**
- Already in RyCode ecosystem (Charm tools)
- Creates GIF/PNG screenshots
- Supports automated test scenarios
- Can capture specific frames
- Integrates well with our TUI stack

### 2. Visual Comparison Strategy

**Approach A: Pixel-Perfect Comparison**
```bash
# Capture screenshots for each theme
vhs capture-claude.tape
vhs capture-gemini.tape
vhs capture-codex.tape
vhs capture-qwen.tape

# Compare with reference images
compare reference/claude.png output/claude.png diff/claude-diff.png
```

**Pros:**
- Exact visual verification
- Catches all rendering changes
- Easy to automate

**Cons:**
- Sensitive to font rendering differences
- Terminal size must be exact
- OS-specific rendering variations

**Approach B: Color Sampling Verification**
```go
// Sample key UI elements and verify colors
func TestThemeColors(t *testing.T) {
    theme.SwitchToProvider("claude")

    // Verify border color at specific location
    borderColor := capturePixelAt(50, 10)
    assert.Equal(t, "#D4754C", hexColor(borderColor))

    // Verify text color
    textColor := capturePixelAt(60, 15)
    assert.Equal(t, "#E8D5C4", hexColor(textColor))
}
```

**Pros:**
- Less sensitive to minor rendering differences
- Focuses on actual theme colors
- Works across different terminals

**Cons:**
- Doesn't catch layout issues
- More complex to implement

**Approach C: ASCII Art Hash Comparison**
```go
// Capture terminal output as text
func TestThemeLayout(t *testing.T) {
    output := captureTerminalText()

    // Generate hash of structure (ignore colors)
    hash := structuralHash(stripAnsi(output))

    // Compare with reference
    assert.Equal(t, referenceHash, hash)
}
```

**Pros:**
- Font-independent
- Terminal-independent
- Focuses on structure

**Cons:**
- Doesn't verify actual colors
- May miss visual regressions

---

## Recommended Implementation

### Phase 3.3A: Color Verification Tests

Create automated tests that verify theme colors are applied correctly.

**Test File**: `test_theme_visual.go`

```go
package main

import (
    "testing"
    "github.com/aaronmrosenthal/rycode/internal/theme"
    "github.com/stretchr/testify/assert"
)

func TestClaudeThemeColors(t *testing.T) {
    theme.SwitchToProvider("claude")
    t := theme.CurrentTheme()

    // Verify primary colors
    assert.Equal(t, "#D4754C", extractHex(t.Primary()))
    assert.Equal(t, "#E8D5C4", extractHex(t.Text()))
    assert.Equal(t, "#1A1816", extractHex(t.Background()))
}

func TestGeminiThemeColors(t *testing.T) {
    theme.SwitchToProvider("gemini")
    t := theme.CurrentTheme()

    assert.Equal(t, "#4285F4", extractHex(t.Primary()))
    assert.Equal(t, "#E8EAED", extractHex(t.Text()))
    assert.Equal(t, "#0D0D0D", extractHex(t.Background()))
}

// ... tests for Codex and Qwen
```

### Phase 3.3B: Manual Visual Verification

Create `.tape` files for manual review:

**File**: `visual-tests/claude-theme.tape`
```tape
Output claude-theme.gif

Set Theme "claude"
Set FontSize 14
Set Width 1200
Set Height 800

Type "rycode"
Enter
Sleep 1s

# Switch to Claude theme
Type "Tab"
Sleep 500ms

# Show chat with Claude colors
Type "Hello from Claude theme"
Enter
Sleep 2s

Screenshot claude-theme.png
```

**Benefits:**
- Visual review by humans
- Can be run locally or in CI
- Creates artifacts for documentation
- Shows actual user experience

### Phase 3.3C: Reference Screenshot Comparison

**Directory Structure:**
```
visual-tests/
├── tapes/
│   ├── claude-theme.tape
│   ├── gemini-theme.tape
│   ├── codex-theme.tape
│   └── qwen-theme.tape
├── reference/
│   ├── claude-theme.png
│   ├── gemini-theme.png
│   ├── codex-theme.png
│   └── qwen-theme.png
└── output/
    ├── claude-theme.png
    ├── gemini-theme.png
    ├── codex-theme.png
    └── qwen-theme.png
```

**CI Script:**
```bash
#!/bin/bash
# visual-tests/run-visual-tests.sh

echo "Running visual regression tests..."

# Generate screenshots
for theme in claude gemini codex qwen; do
    vhs visual-tests/tapes/${theme}-theme.tape
    mv ${theme}-theme.png visual-tests/output/
done

# Compare with references
for theme in claude gemini codex qwen; do
    compare -metric AE \
        visual-tests/reference/${theme}-theme.png \
        visual-tests/output/${theme}-theme.png \
        visual-tests/diff/${theme}-diff.png 2>&1 | \
        tee visual-tests/diff/${theme}-result.txt
done

echo "Visual tests complete. Check diff/ directory."
```

---

## Native CLI Comparison

### Capturing Reference Screenshots

**Claude Code:**
```bash
# Launch Claude Code
claude-code

# Wait for full render
sleep 2

# Capture screenshot
screencapture -R x,y,w,h claude-code-reference.png
```

**Gemini CLI:**
```bash
# Launch Gemini
gemini

# Capture
screencapture -R x,y,w,h gemini-cli-reference.png
```

**Comparison Strategy:**
1. Place native CLI and RyCode side-by-side
2. Verify color matching visually
3. Use color picker to sample exact RGB values
4. Document any intentional differences

---

## Automated Testing Approach

### Option 1: VHS-Based CI Tests

```yaml
# .github/workflows/visual-tests.yml
name: Visual Regression Tests

on: [push, pull_request]

jobs:
  visual-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Install VHS
        run: |
          go install github.com/charmbracelet/vhs@latest

      - name: Run visual tests
        run: |
          cd packages/tui/visual-tests
          ./run-visual-tests.sh

      - name: Upload screenshots
        uses: actions/upload-artifact@v3
        with:
          name: visual-test-results
          path: visual-tests/output/

      - name: Upload diffs
        if: failure()
        uses: actions/upload-artifact@v3
        with:
          name: visual-test-diffs
          path: visual-tests/diff/
```

### Option 2: Color Verification Tests

```go
// test_theme_visual_verification.go
package main

import (
    "fmt"
    "testing"
    "github.com/aaronmrosenthal/rycode/internal/theme"
)

type ColorTest struct {
    Provider string
    Element  string
    Expected string
}

var colorTests = []ColorTest{
    // Claude
    {"claude", "Primary", "#D4754C"},
    {"claude", "Text", "#E8D5C4"},
    {"claude", "Background", "#1A1816"},

    // Gemini
    {"gemini", "Primary", "#4285F4"},
    {"gemini", "Text", "#E8EAED"},
    {"gemini", "Background", "#0D0D0D"},

    // Codex
    {"codex", "Primary", "#10A37F"},
    {"codex", "Text", "#ECECEC"},
    {"codex", "Background", "#0E0E0E"},

    // Qwen
    {"qwen", "Primary", "#FF6A00"},
    {"qwen", "Text", "#F0E8DC"},
    {"qwen", "Background", "#161410"},
}

func TestThemeColorAccuracy(t *testing.T) {
    for _, test := range colorTests {
        theme.SwitchToProvider(test.Provider)
        th := theme.CurrentTheme()

        var actual string
        switch test.Element {
        case "Primary":
            actual = colorToHex(th.Primary())
        case "Text":
            actual = colorToHex(th.Text())
        case "Background":
            actual = colorToHex(th.Background())
        }

        if actual != test.Expected {
            t.Errorf("%s %s: expected %s, got %s",
                test.Provider, test.Element, test.Expected, actual)
        }
    }
}
```

---

## Testing Checklist

### Manual Visual Review
- [ ] Claude theme matches Claude Code aesthetics
- [ ] Gemini theme matches Gemini CLI aesthetics
- [ ] Codex theme matches OpenAI Codex aesthetics
- [ ] Qwen theme matches Qwen CLI aesthetics
- [ ] Theme switching is smooth (no flicker)
- [ ] Colors are consistent across UI elements
- [ ] Text is readable in all themes
- [ ] Borders are clearly visible
- [ ] Status colors (success, error, warning) are distinct

### Automated Tests
- [ ] All theme colors match specifications
- [ ] Theme switching doesn't corrupt colors
- [ ] UI elements use correct theme colors
- [ ] Screenshots match reference images (within tolerance)
- [ ] No color bleeding between themes
- [ ] Typography remains consistent
- [ ] Layout is preserved across themes

### Performance
- [ ] Visual tests run in <30 seconds
- [ ] Screenshots are deterministic (same output each run)
- [ ] No memory leaks during screenshot capture
- [ ] CI integration doesn't timeout

---

## Success Criteria

### Visual Accuracy
- ✅ All colors match specification within 5% tolerance
- ✅ Theme switching shows correct colors immediately
- ✅ No visual artifacts or glitches
- ✅ Consistent rendering across terminals

### Testing Coverage
- ✅ All 4 provider themes tested
- ✅ Key UI elements verified (borders, text, status)
- ✅ Both automated and manual verification
- ✅ CI integration for regression detection

### Documentation
- ✅ Reference screenshots captured
- ✅ Testing process documented
- ✅ Failure investigation guide
- ✅ Maintenance procedures

---

## Implementation Priority

### Phase 3.3A: Color Verification (High Priority) ✅
**Status**: Can implement immediately
**Effort**: 2-3 hours
**Value**: High - catches color regressions

### Phase 3.3B: VHS Manual Review (Medium Priority)
**Status**: Requires VHS setup
**Effort**: 4-6 hours
**Value**: Medium - good for documentation

### Phase 3.3C: Screenshot Comparison (Lower Priority)
**Status**: Requires ImageMagick + VHS
**Effort**: 8-10 hours
**Value**: Medium - pixel-perfect verification

---

## Recommended Next Steps

1. **Implement Color Verification Tests** (Phase 3.3A)
   - Quick to build
   - High value
   - No external dependencies
   - Can run in CI immediately

2. **Create Reference Documentation**
   - Manual screenshots of each theme
   - Side-by-side comparison with native CLIs
   - Document intentional differences

3. **CI Integration**
   - Add color tests to pre-push hooks
   - Run on every PR
   - Block merges on color regression

4. **Future Enhancement**
   - Full VHS-based screenshot testing
   - Automated comparison pipeline
   - Visual regression dashboard

---

## Conclusion

Visual testing for TUI applications requires a different approach than web apps. We recommend starting with **automated color verification tests** (Phase 3.3A) as they provide immediate value with minimal setup.

Full screenshot-based testing (Phase 3.3B/C) is valuable but requires more infrastructure. We can add this incrementally as needed.

**Key Insight**: Our themes are already verified for accessibility (48/48 tests) and performance (317ns switching). Color verification tests will ensure visual accuracy is maintained over time.

---

**Implementation Status**: Design Complete ✅

**Next Action**: Build color verification test suite
