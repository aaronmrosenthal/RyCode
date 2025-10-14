# Dynamic Provider Theming Specification

**Vision**: When users Tab between models, the entire TUI theme switches to match that provider's native CLI aesthetic. Users familiar with Claude Code, Gemini CLI, Codex, or Qwen CLI will instantly feel at home.

---

## User Story

> "As a developer who uses Claude Code daily, when I Tab to Claude models in RyCode, I want the interface to look and feel like Claude Code - warm orange borders, copper accents, and familiar typography. When I Tab to Gemini, I want the vibrant blue-pink gradient and colorful aesthetic I know from Gemini CLI."

---

## Provider Theme Definitions

### 1. Claude Theme (Based on Screenshots)

**Brand Identity**: Warm, approachable, developer-friendly

**Color Palette**:
```go
ClaudeTheme := theme.ProviderTheme{
    Name: "Claude",

    // Primary accents - warm orange/copper
    Primary:   "#D4754C", // Copper orange (borders, highlights)
    Secondary: "#B85C3C", // Darker copper
    Accent:    "#F08C5C", // Lighter warm orange

    // Background - dark warm
    Background:        "#1A1816", // Warm dark brown
    BackgroundPanel:   "#2C2622", // Slightly lighter panel
    BackgroundElement: "#3A3330", // Element backgrounds

    // Borders - distinctive orange
    BorderSubtle: "#4A3F38",
    Border:       "#D4754C", // Signature copper
    BorderActive: "#F08C5C", // Bright on focus

    // Text - warm tones
    Text:      "#E8D5C4", // Warm cream
    TextMuted: "#9C8373", // Muted warm gray

    // Status colors
    Success: "#6FA86F", // Muted green
    Error:   "#D47C7C", // Warm red
    Info:    "#D4754C", // Use primary
    Warning: "#E8A968", // Warm amber

    // Special - avatar background
    AvatarBg: "#D4754C",
}
```

**Typography**:
- Friendly, readable monospace
- Slightly rounded edges on panels
- Warm, inviting spacing

**UI Elements**:
- Orange border on active input
- Copper-colored model badges
- Warm glow on hover states
- Friendly avatar (pixelated character)

---

### 2. Gemini Theme (Based on Screenshot)

**Brand Identity**: Modern, vibrant, AI-forward

**Color Palette**:
```go
GeminiTheme := theme.ProviderTheme{
    Name: "Gemini",

    // Primary accents - Google blue to pink gradient
    Primary:   "#4285F4", // Google blue
    Secondary: "#9B72F2", // Purple midpoint
    Accent:    "#EA4335", // Google red/pink

    // Background - cool dark
    Background:        "#0D0D0D", // Pure black
    BackgroundPanel:   "#1A1A1A", // Dark gray
    BackgroundElement: "#2A2A2A", // Element backgrounds

    // Borders - gradient inspired
    BorderSubtle: "#2A2A45",
    Border:       "#4285F4", // Blue primary
    BorderActive: "#9B72F2", // Purple on focus

    // Text - cool tones
    Text:      "#E8EAED", // Light gray
    TextMuted: "#9AA0A6", // Medium gray

    // Status colors - vibrant
    Success: "#34A853", // Google green
    Error:   "#EA4335", // Google red
    Info:    "#4285F4", // Blue
    Warning: "#FBBC04", // Google yellow

    // Special - gradient for ASCII art
    GradientStart: "#4285F4", // Blue
    GradientMid:   "#9B72F2", // Purple
    GradientEnd:   "#EA4335", // Pink/red
}
```

**Typography**:
- Modern, sharp monospace
- Clean, minimal spacing
- Tech-forward aesthetic

**UI Elements**:
- Gradient border on panels
- Colorful ASCII art logo
- Blue-purple-pink gradient on active elements
- "Thinking" indicator with gradient animation
- Progress bars with gradient fill

---

### 3. Codex Theme (OpenAI)

**Brand Identity**: Professional, technical, precise

**Color Palette**:
```go
CodexTheme := theme.ProviderTheme{
    Name: "Codex",

    // Primary accents - OpenAI teal
    Primary:   "#10A37F", // OpenAI teal
    Secondary: "#0D8569", // Darker teal
    Accent:    "#1FC2AA", // Lighter teal

    // Background - neutral dark
    Background:        "#0E0E0E", // Almost black
    BackgroundPanel:   "#1C1C1C", // Dark gray
    BackgroundElement: "#2D2D2D", // Element backgrounds

    // Borders - teal accent
    BorderSubtle: "#2D3D38",
    Border:       "#10A37F", // Teal
    BorderActive: "#1FC2AA", // Bright teal

    // Text - clean neutrals
    Text:      "#ECECEC", // Off-white
    TextMuted: "#8E8E8E", // Medium gray

    // Status colors
    Success: "#10A37F", // Use primary
    Error:   "#EF4444", // Clean red
    Info:    "#3B82F6", // Blue
    Warning: "#F59E0B", // Amber

    // Special - technical feel
    CodeBlock: "#1C2D27", // Dark teal tint
}
```

**Typography**:
- Technical, precise monospace
- Tight, efficient spacing
- Professional aesthetic

**UI Elements**:
- Clean teal borders
- Minimalist badges
- Technical progress indicators
- Code-first interface design

---

### 4. Qwen Theme (Alibaba)

**Brand Identity**: Modern, innovative, Chinese tech aesthetic

**Color Palette**:
```go
QwenTheme := theme.ProviderTheme{
    Name: "Qwen",

    // Primary accents - Alibaba orange
    Primary:   "#FF6A00", // Alibaba orange
    Secondary: "#E55D00", // Darker orange
    Accent:    "#FF8533", // Lighter orange

    // Background - warm dark
    Background:        "#161410", // Warm black
    BackgroundPanel:   "#221E18", // Dark warm gray
    BackgroundElement: "#2F2A22", // Element backgrounds

    // Borders - orange accent
    BorderSubtle: "#3A352C",
    Border:       "#FF6A00", // Orange
    BorderActive: "#FF8533", // Bright orange

    // Text - neutral with warm tint
    Text:      "#F0E8DC", // Warm off-white
    TextMuted: "#A0947C", // Warm gray

    // Status colors
    Success: "#52C41A", // Chinese green
    Error:   "#FF4D4F", // Chinese red
    Info:    "#1890FF", // Chinese blue
    Warning: "#FAAD14", // Chinese gold

    // Special - Chinese design elements
    Accent2: "#FAAD14", // Gold accent
}
```

**Typography**:
- Modern, international monospace
- Balanced spacing
- Contemporary aesthetic

**UI Elements**:
- Orange/gold color scheme
- Modern Chinese design language
- Clean, international interface
- Elegant progress animations

---

## Implementation Architecture

### 1. Theme Structure

```go
// packages/tui/internal/theme/provider_themes.go

package theme

type ProviderTheme struct {
    Name string

    // Implement Theme interface
    BaseTheme

    // Provider-specific extensions
    LogoASCII       string
    LoadingSpinner  string
    WelcomeMessage  string
    TypingIndicator TypingStyle
}

type TypingStyle struct {
    Text      string // "Thinking..." or "Processing..."
    Animation string // "dots", "gradient", "pulse"
    Color     string // Override for typing indicator
}
```

### 2. Theme Manager

```go
// packages/tui/internal/theme/manager.go

type ThemeManager struct {
    themes         map[string]ProviderTheme
    currentTheme   ProviderTheme
    previousTheme  ProviderTheme

    // Animation
    transitionDuration time.Duration
    transitionEasing   string
}

func (tm *ThemeManager) SwitchToProvider(providerID string) {
    newTheme := tm.themes[providerID]

    // Animate transition
    tm.animateTransition(tm.currentTheme, newTheme)

    // Update current
    tm.previousTheme = tm.currentTheme
    tm.currentTheme = newTheme

    // Emit event for UI updates
    tm.emitThemeChanged(newTheme)
}

func (tm *ThemeManager) animateTransition(from, to ProviderTheme) {
    // Crossfade borders
    // Morph colors
    // Smooth gradient transition
    // 200-300ms transition
}
```

### 3. Model Selector Integration

```go
// packages/tui/internal/components/modelselector/modelselector.go

func (m *Model) switchProvider(providerID string) {
    // Switch model
    m.selectedProvider = providerID

    // TRIGGER THEME CHANGE
    m.themeManager.SwitchToProvider(providerID)

    // Update UI
    m.updateView()
}
```

### 4. View Updates

Every component that renders visual elements needs to react to theme changes:

```go
// packages/tui/internal/components/chat/messages.go

func (m *Model) View() string {
    // Get current theme
    theme := m.themeManager.CurrentTheme()

    // Apply theme colors
    borderStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(theme.Border())

    // Render with current theme
    return borderStyle.Render(m.content)
}
```

---

## Visual Transitions

### Theme Switch Animation

```
1. User presses Tab
2. Model selector advances to next provider
3. ThemeManager detects provider change
4. Crossfade animation begins (200ms):
   - Borders morph from old color to new
   - Background subtly shifts tone
   - Text colors smoothly transition
   - Logo/ASCII art fades in
5. New theme fully applied
6. User sees familiar provider aesthetic
```

### Example: Claude → Gemini Transition

```
Frame 1 (0ms):   Claude theme (orange borders)
Frame 2 (50ms):  Borders fade to purple
Frame 3 (100ms): Background shifts cooler
Frame 4 (150ms): Gemini colors emerge
Frame 5 (200ms): Gemini theme (blue-pink gradient)
```

---

## UI Elements Affected

### 1. Borders
- Main chat border
- Input box border
- Model selector border
- Panel borders

### 2. Badges
- Active model badge color
- Provider indicators
- Status badges

### 3. Text Colors
- Primary text (subtle tint)
- Secondary text
- Muted text
- Highlights

### 4. Special Elements
- Loading spinner (provider-specific animation)
- Typing indicator (matches provider style)
- Welcome screen (provider-specific ASCII art)
- Error messages (provider color scheme)

### 5. Interactive Elements
- Button hover states
- Focus indicators
- Selection highlights
- Progress bars

---

## Playwright Visual Testing

### Test Suite

```typescript
// packages/tui/tests/visual/theme-switching.spec.ts

import { test, expect } from '@playwright/test';

test.describe('Provider Theme Switching', () => {

  test('Claude theme matches native CLI', async ({ page }) => {
    // Launch RyCode with Claude
    await page.goto('rycode://claude');

    // Capture screenshot
    await page.screenshot({
      path: 'tests/visual/snapshots/claude-theme.png'
    });

    // Compare with reference
    await expect(page).toHaveScreenshot('claude-reference.png', {
      threshold: 0.1, // 10% tolerance
    });

    // Verify orange borders
    const border = await page.locator('.chat-border');
    const borderColor = await border.evaluate(
      el => getComputedStyle(el).borderColor
    );
    expect(borderColor).toContain('212, 117, 76'); // #D4754C RGB
  });

  test('Theme switches on Tab', async ({ page }) => {
    await page.goto('rycode://claude');

    // Capture initial state
    const claudeScreenshot = await page.screenshot();

    // Press Tab to switch to Gemini
    await page.keyboard.press('Tab');
    await page.waitForTimeout(250); // Wait for transition

    // Capture new state
    const geminiScreenshot = await page.screenshot();

    // Screenshots should be different
    expect(claudeScreenshot).not.toEqual(geminiScreenshot);

    // Verify Gemini colors
    const border = await page.locator('.chat-border');
    const borderColor = await border.evaluate(
      el => getComputedStyle(el).borderColor
    );
    expect(borderColor).toContain('66, 133, 244'); // #4285F4 RGB
  });

  test('Smooth transition animation', async ({ page }) => {
    await page.goto('rycode://claude');

    // Record transition frames
    const frames = [];
    page.on('framenavigated', () => {
      frames.push(page.screenshot());
    });

    await page.keyboard.press('Tab');
    await page.waitForTimeout(300);

    // Should have intermediate frames (smooth animation)
    expect(frames.length).toBeGreaterThan(3);
  });

  test('All providers have distinct themes', async ({ page }) => {
    const screenshots = {};

    for (const provider of ['claude', 'gemini', 'codex', 'qwen']) {
      await page.goto(`rycode://${provider}`);
      screenshots[provider] = await page.screenshot();
    }

    // All screenshots should be unique
    const hashes = Object.values(screenshots).map(img => hash(img));
    const uniqueHashes = new Set(hashes);
    expect(uniqueHashes.size).toBe(4);
  });
});
```

---

## Implementation Phases

### Phase 1: Theme Infrastructure (Week 1)
- [ ] Create `ProviderTheme` struct
- [ ] Implement `ThemeManager` with hot-swapping
- [ ] Define all 4 provider themes
- [ ] Add theme switching logic to model selector
- [ ] Basic color transitions (no animations)

### Phase 2: Visual Polish (Week 2)
- [ ] Add transition animations (200ms crossfade)
- [ ] Provider-specific ASCII art
- [ ] Custom loading spinners per provider
- [ ] Provider-specific typing indicators
- [ ] Welcome screen variants

### Phase 3: Testing & Refinement (Week 3)
- [ ] Playwright visual regression tests
- [ ] Screenshot comparison with native CLIs
- [ ] Performance optimization (GPU acceleration)
- [ ] Accessibility audit (color contrast)
- [ ] User testing with developers familiar with each CLI

### Phase 4: Documentation (Week 4)
- [ ] Theme customization guide
- [ ] Custom provider theme API
- [ ] Visual design system docs
- [ ] Developer onboarding (theme-aware)

---

## Success Criteria

### User Experience
- ✅ Users familiar with Claude Code immediately recognize Claude theme
- ✅ Gemini users see familiar blue-pink gradients and colorful aesthetic
- ✅ Theme switch on Tab feels instant and smooth (<300ms)
- ✅ Each provider theme has distinct personality
- ✅ Transitions are delightful, not jarring

### Technical
- ✅ Theme switching has no performance impact
- ✅ All colors meet WCAG AA contrast standards
- ✅ Visual tests catch theme regressions
- ✅ Themes load instantly (no flicker)

### Design
- ✅ Matches native CLI aesthetics (Playwright verified)
- ✅ Consistent design language within each theme
- ✅ Cohesive color palette per provider
- ✅ Professional, polished appearance

---

## Future Enhancements

### Custom Themes
Allow users to create custom provider themes:

```json
{
  "name": "My Custom Theme",
  "provider": "claude",
  "extends": "claude-base",
  "colors": {
    "primary": "#FF00FF",
    "background": "#000000"
  }
}
```

### Theme Marketplace
- Community-contributed themes
- Seasonal themes (dark mode variants)
- High-contrast accessibility themes

### Advanced Animations
- Particle effects on theme switch
- Ripple transitions from cursor position
- Provider logo morphing animations

---

## Notes

This specification brings **emotional design** to RyCode. When a developer who loves Claude Code sees the familiar warm orange borders and copper accents, they'll feel at home instantly. The theme becomes part of the tool's personality, making it memorable and delightful to use.

The key insight: **familiarity breeds confidence**. By matching the native CLI aesthetics, we reduce cognitive load and make Tab-switching feel like moving between trusted tools, not switching contexts entirely.
