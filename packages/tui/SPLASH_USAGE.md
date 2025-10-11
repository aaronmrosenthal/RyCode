# RyCode Splash Screen - Usage Guide

> Complete guide to splash screen configuration, command-line options, and customization

---

## 🚀 Quick Start

**First launch:**
```bash
./rycode
# Splash shows automatically on first run
```

**Skip the splash:**
```bash
./rycode --no-splash
```

**Force show splash:**
```bash
./rycode --splash
```

**Infinite donut mode:**
```bash
./rycode donut
```

---

## ⌨️ Keyboard Controls

### During Splash

| Key | Action |
|-----|--------|
| `S` | Skip splash (continue to TUI) |
| `ESC` | Skip and disable forever |
| `?` | Toggle math equations |
| `Enter` / `Space` | Continue from closer screen |
| `↑↑↓↓←→←→BA` | Konami code (rainbow mode) |

### In Donut Mode

| Key | Action |
|-----|--------|
| `Q` | Quit donut mode |
| `?` | Toggle math equations |
| `↑↑↓↓←→←→BA` | Konami code (rainbow mode) |

---

## ⚙️ Configuration

### Config File Location

```
~/.rycode/config.json
```

### Default Configuration

```json
{
  "splash_enabled": true,
  "splash_frequency": "first",
  "reduced_motion": false,
  "color_mode": "auto"
}
```

### Configuration Options

#### `splash_enabled`
**Type:** `boolean`
**Default:** `true`
**Description:** Master switch for splash screen

**Values:**
- `true` - Splash can show (based on frequency)
- `false` - Splash never shows

**Example:**
```json
{
  "splash_enabled": false
}
```

---

#### `splash_frequency`
**Type:** `string`
**Default:** `"first"`
**Description:** How often the splash should appear

**Values:**
- `"first"` - Only on first run (default)
- `"always"` - Every launch
- `"random"` - 10% chance on each launch
- `"never"` - Never show (same as `splash_enabled: false`)

**Examples:**
```json
// Show every time
{
  "splash_frequency": "always"
}

// Show randomly (10% chance)
{
  "splash_frequency": "random"
}

// Never show
{
  "splash_frequency": "never"
}
```

---

#### `reduced_motion`
**Type:** `boolean`
**Default:** `false`
**Description:** Accessibility setting to respect reduced motion preference

**Values:**
- `true` - Disable splash (accessibility)
- `false` - Normal behavior

**Auto-detection:**
The splash respects the `PREFERS_REDUCED_MOTION` environment variable:
```bash
export PREFERS_REDUCED_MOTION=1
./rycode  # Splash will not show
```

**Example:**
```json
{
  "reduced_motion": true
}
```

---

#### `color_mode`
**Type:** `string`
**Default:** `"auto"`
**Description:** Terminal color support

**Values:**
- `"auto"` - Auto-detect terminal capabilities
- `"truecolor"` - 16 million colors (best)
- `"256"` - 256 colors
- `"16"` - 16 colors (basic)

**Auto-detection:**
- Checks `$COLORTERM` for truecolor support
- Falls back to `$TERM` for 256-color detection
- Respects `$NO_COLOR` environment variable

**Example:**
```json
{
  "color_mode": "256"
}
```

---

## 🖥️ Command-Line Flags

### `--splash`
**Force show the splash screen**

Overrides all configuration settings and shows the splash even if disabled.

**Usage:**
```bash
./rycode --splash
```

**Notes:**
- Does not update the marker file
- Can be used repeatedly
- Useful for demonstrations

---

### `--no-splash`
**Skip the splash screen**

Prevents splash from showing this launch only. Does not change configuration.

**Usage:**
```bash
./rycode --no-splash
```

**Notes:**
- Does not update config
- Temporary for this launch only
- Useful for automation/scripts

---

## 🎨 Fallback Modes

The splash automatically adapts to your terminal capabilities.

### Full Mode (Default)
**Requirements:**
- Terminal ≥ 80×24
- Truecolor or 256-color support
- Unicode support

**Features:**
- Full 3D cortex animation
- Cyan-magenta gradient
- All easter eggs enabled
- 30 FPS smooth animation

---

### Simplified Mode
**Triggers:**
- Terminal < 80×24
- Limited color support (16-color)
- No unicode support

**Features:**
- Text-only splash
- Simple model list
- No animation
- Centered layout

**Example output:**
```
═══════════════════════════════════
      RYCODE NEURAL CORTEX
═══════════════════════════════════

🧩 Claude  • Logical Reasoning
⚙️  Gemini  • System Architecture
💻 Codex   • Code Generation
🔎 Qwen   • Research Pipeline
🤖 Grok   • Humor & Chaos
✅ GPT    • Language Core

⚡ SIX MINDS. ONE COMMAND LINE.

Press any key to continue...
```

---

### Skipped Mode
**Triggers:**
- Terminal < 60×20 (extremely small)
- `splash_enabled: false` in config
- `--no-splash` flag used
- Pressed `ESC` on previous run

**Behavior:**
- No splash shown at all
- Direct launch to TUI

---

## 🔧 Advanced Usage

### Disable Splash Permanently

**Method 1: Press ESC during splash**
```
# Launch RyCode
./rycode

# When splash appears, press ESC
# Config is automatically updated to disable splash
```

**Method 2: Edit config manually**
```bash
# Edit config file
nano ~/.rycode/config.json

# Set splash_enabled to false
{
  "splash_enabled": false
}
```

**Method 3: Delete config**
```bash
# Remove config to reset to defaults
rm ~/.rycode/config.json
```

---

### Reset First-Run Status

To see the splash again as if it's the first run:

```bash
# Remove marker file
rm ~/.rycode/.splash_shown

# Next launch will show splash
./rycode
```

---

### Force Specific Frequency

```bash
# Edit config
nano ~/.rycode/config.json

# Set frequency
{
  "splash_frequency": "always"  // Every launch
  "splash_frequency": "random"  // 10% chance
  "splash_frequency": "first"   // Only first run
  "splash_frequency": "never"   // Never
}
```

---

### Environment Variables

The splash respects several environment variables:

#### `PREFERS_REDUCED_MOTION`
Accessibility setting to disable animations:
```bash
export PREFERS_REDUCED_MOTION=1
./rycode  # Splash disabled
```

#### `NO_COLOR`
Disable all colors:
```bash
export NO_COLOR=1
./rycode  # Monochrome splash
```

#### `COLORTERM`
Indicates truecolor support:
```bash
export COLORTERM=truecolor
./rycode  # Uses 16M colors
```

#### `TERM`
Terminal type indicator:
```bash
export TERM=xterm-256color
./rycode  # Uses 256 colors
```

---

## 📊 Performance Tuning

### Adaptive Frame Rate

The splash automatically adjusts FPS based on system performance:

**Fast systems (M1/M2/M3, modern Intel):**
- 30 FPS (33ms per frame)
- Smooth, fluid animation

**Slower systems (Raspberry Pi, old laptops):**
- Automatically drops to 15 FPS (66ms per frame)
- Still smooth, uses less CPU

**Detection:**
- Monitors last 30 frames
- If average frame time > 50ms, drops to 15 FPS
- No user configuration needed

---

### Manual Performance Control

If splash is laggy on your system:

**Option 1: Disable entirely**
```json
{
  "splash_enabled": false
}
```

**Option 2: Use --no-splash flag**
```bash
./rycode --no-splash
```

**Option 3: Reduce color depth**
```json
{
  "color_mode": "16"
}
```

---

## 🐛 Troubleshooting

### Splash doesn't show on first run

**Check:**
1. Config file exists? `cat ~/.rycode/config.json`
2. splash_enabled is true?
3. Marker file exists? `ls ~/.rycode/.splash_shown`

**Solution:**
```bash
# Reset first-run
rm ~/.rycode/.splash_shown
./rycode
```

---

### Splash shows every time (unwanted)

**Check config:**
```bash
cat ~/.rycode/config.json
```

**Expected:**
```json
{
  "splash_frequency": "first"  // Should be "first", not "always"
}
```

**Fix:**
```bash
# Edit config
nano ~/.rycode/config.json

# Change to:
{
  "splash_frequency": "first"
}
```

---

### Splash is laggy/slow

**Solutions:**

1. **Check terminal:** Some terminals are slower
   - Try: iTerm2, Windows Terminal, Alacritty
   - Avoid: xterm, Windows CMD

2. **SSH connection?** Remote sessions are slower
   - Splash will auto-adapt to 15 FPS
   - Or use: `./rycode --no-splash`

3. **Old hardware?** Adaptive FPS should help
   - Wait ~1 second for detection
   - Or disable: `"splash_enabled": false`

---

### Colors look wrong

**Check color mode:**
```bash
echo $COLORTERM
echo $TERM
```

**Force specific mode:**
```json
{
  "color_mode": "256"  // or "16" or "truecolor"
}
```

---

### Splash shows but immediately disappears

**This is normal if:**
- You previously pressed ESC (disabled forever)
- Config has `splash_enabled: false`
- Terminal is too small (<60×20)

**Check:**
```bash
cat ~/.rycode/config.json
# Look for splash_enabled: false
```

**Fix:**
```bash
# Edit config
nano ~/.rycode/config.json

# Set:
{
  "splash_enabled": true
}
```

---

## 📚 Examples

### Example 1: Development Workflow
```bash
# First time setup
./rycode              # Splash shows (first run)

# Daily use
./rycode              # No splash (already shown)

# Occasional check
./rycode --splash     # Force show splash
```

---

### Example 2: Presentation/Demo
```bash
# Always show splash for demos
nano ~/.rycode/config.json

{
  "splash_frequency": "always"
}

# Or use flag:
./rycode --splash
```

---

### Example 3: CI/CD Pipeline
```bash
# Never show splash in automation
./rycode --no-splash

# Or set config:
{
  "splash_enabled": false
}
```

---

### Example 4: Accessibility
```bash
# Disable animations
export PREFERS_REDUCED_MOTION=1
./rycode

# Or in config:
{
  "reduced_motion": true
}
```

---

## 🎓 Best Practices

1. **First run:** Let the splash show once (it's cool!)
2. **Daily use:** Default "first" frequency is perfect
3. **Presentations:** Use `--splash` flag or "always" frequency
4. **Automation:** Use `--no-splash` flag
5. **Accessibility:** Set `reduced_motion: true` if needed
6. **Slow systems:** Adaptive FPS handles it automatically

---

## 🤝 Contributing

Found a bug? Want a new feature?

- **Issues:** https://github.com/aaronmrosenthal/RyCode/issues
- **Discussions:** https://github.com/aaronmrosenthal/RyCode/discussions
- **PRs welcome!**

---

**🤖 Built with ❤️ by Claude AI**

*For more details, see:*
- **Easter Eggs:** [EASTER_EGGS.md](EASTER_EGGS.md)
- **Implementation Plan:** [SPLASH_IMPLEMENTATION_PLAN.md](SPLASH_IMPLEMENTATION_PLAN.md)
- **Task Breakdown:** [SPLASH_TASKS.md](SPLASH_TASKS.md)
