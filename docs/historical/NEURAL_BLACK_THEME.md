# Neural Black Theme

The definitive **"AI in control"** theme for Toolkit-CLI and RyCode, designed for terminal realism fused with next-gen clarity.

## 🎨 Color Palette

| Role | Color | HEX | Description |
|------|-------|-----|-------------|
| **Background** | Deep Black | `#0B0C10` | Matte black for focus and contrast |
| **Primary** | Electric Mint | `#00FF88` | Success, alignment, intelligence |
| **Secondary** | Cool Cyan | `#00AEEF` | System, context, thought |
| **Warning** | Gold | `#FFD166` | Drift, caution |
| **Error** | Neon Red | `#FF5555` | Failure or interruption |
| **Text** | Soft White | `#E0E0E0` | Primary narrative text |
| **Muted** | Gray | `#7B7F8B` | Secondary info, comments |
| **Purple** | Educational | `#B26FFF` | Persona commentary |
| **Selection** | Graphite | `#1F1F1F` | Subtle highlight panel |

## 🧠 Typography & Structure

- **Font**: JetBrains Mono / Fira Code Retina (ligatures ON)
- **Line Height**: 1.3 for spacious readouts
- **Indent Guides**: Thin cyan lines
- **Prompt Symbol**: `λ` (lambda — symbol of creation)

### Status Bullets

- ✅ Success / Complete
- ⚠️ Warning / Drift risk
- 💡 Info / Learning insight
- 🧩 Context update
- 🚀 Ship ready
- 🧠 Thinking / Analysis
- ❌ Error / Failure

## 📦 Installation

### Python Rich (toolkit-cli)

```python
from toolkit_theme import toolkit_theme, get_console

# Use the pre-configured console
console = get_console()
console.print("✅ Phase 3 Complete", style="success")

# Or create your own with the theme
from rich.console import Console
console = Console(theme=toolkit_theme)
```

### VS Code

1. Copy `.vscode/neural-black.json` to your VS Code user themes folder
2. Open Command Palette (Cmd/Ctrl+Shift+P)
3. Select "Preferences: Color Theme"
4. Choose "Neural Black"

### iTerm2 (macOS)

1. Open iTerm2 Preferences
2. Go to Profiles → Colors
3. Click "Color Presets..." → "Import..."
4. Select `.terminal-themes/neural-black-iterm2.itermcolors`
5. Select "Neural Black" from the presets

### Windows Terminal

1. Open Windows Terminal settings (Ctrl+,)
2. Click "Open JSON file"
3. Add the content from `.terminal-themes/neural-black-windows-terminal.json` to the `schemes` array
4. In your profile, set `"colorScheme": "Neural Black"`

### Alacritty

Add to your `alacritty.yml`:

```yaml
colors:
  primary:
    background: '#0B0C10'
    foreground: '#E0E0E0'

  cursor:
    text: '#0B0C10'
    cursor: '#00FF88'

  selection:
    text: '#E0E0E0'
    background: '#1F1F1F'

  normal:
    black:   '#0B0C10'
    red:     '#FF5555'
    green:   '#00FF88'
    yellow:  '#FFD166'
    blue:    '#00AEEF'
    magenta: '#B26FFF'
    cyan:    '#00AEEF'
    white:   '#E0E0E0'

  bright:
    black:   '#7B7F8B'
    red:     '#FF4C6A'
    green:   '#00FF88'
    yellow:  '#FFD166'
    blue:    '#00AEEF'
    magenta: '#B26FFF'
    cyan:    '#00AEEF'
    white:   '#FFFFFF'
```

### Kitty

Add to your `kitty.conf`:

```conf
# Neural Black Theme
foreground            #E0E0E0
background            #0B0C10
selection_foreground  #E0E0E0
selection_background  #1F1F1F
cursor                #00FF88
cursor_text_color     #0B0C10

# Black
color0   #0B0C10
color8   #7B7F8B

# Red
color1   #FF5555
color9   #FF4C6A

# Green
color2   #00FF88
color10  #00FF88

# Yellow
color3   #FFD166
color11  #FFD166

# Blue
color4   #00AEEF
color12  #00AEEF

# Magenta
color5   #B26FFF
color13  #B26FFF

# Cyan
color6   #00AEEF
color14  #00AEEF

# White
color7   #E0E0E0
color15  #FFFFFF
```

## 💬 UX Tone

The Neural Black theme narrative voice is:

**Analytical, confident, and slightly cinematic** — like an AI narrating its own reasoning.

### Example Output

```
🔍 Context Refresh Detected
💡 Spec alignment improved from 0.82 → 0.93
✅ Phase 3 verified — production ready
```

```
Write(AI-COMMAND-SUITE-PHASE3-COMPLETE.md)
  └─ Wrote 544 lines to AI-COMMAND-SUITE-PHASE3-COMPLETE.md
     # ✅ AI Command Suite – Phase 3 Complete
     **Status**: PRODUCTION READY
     **Duration**: 1 day (accelerated from 2-3 day estimate)
     **Test Pass Rate**: 100% (16/16 tests passing)
```

## 🪄 CLI Behavior Enhancements

### Typing Delay Simulation

Sub-60ms "AI typing" effect when generating diffs for cinematic feel.

### Progress Bar

Mint gradient (#00FF88) → cyan trail (#00AEEF)

```python
from rich.progress import Progress, SpinnerColumn, BarColumn, TextColumn

with Progress(
    SpinnerColumn(),
    TextColumn("[progress.description]{task.description}"),
    BarColumn(complete_style="#00FF88", finished_style="#00FF88"),
    TextColumn("[progress.percentage]{task.percentage:>3.0f}%"),
) as progress:
    task = progress.add_task("Processing...", total=100)
    # ... your work here
```

### Code Blocks

Double-border boxes with color-coded headers:

```python
from rich.panel import Panel
from rich.syntax import Syntax

code = '''
def process_data(items: List[str]) -> Dict:
    return {"status": "success"}
'''

syntax = Syntax(code, "python", theme="monokai", line_numbers=True)
panel = Panel(syntax, title="[bold cyan]Implementation[/]", border_style="cyan")
console.print(panel)
```

### Diff Sections

```diff
- **Current Phase**: Phase 2 Complete ✅
+ **Current Phase**: Phase 3 Complete ✅
- **Next Phase**: Phase 3 – /suggest-refactor
+ **Next Phase**: Phase 4 – /generate-tests
```

## 🎯 Rich Style Reference

### Available Styles

```python
# Status
"success"      # Bold electric mint
"warning"      # Bold gold
"error"        # Bold neon red
"info"         # Bold cool cyan

# Context
"context"      # Cool cyan
"comment"      # Dim gray
"meta"         # Italic gray
"educational"  # Purple glow

# Diffs
"diff.add"     # Black on mint background
"diff.remove"  # Black on red background
"diff.context" # Soft white

# UI
"prompt"       # Bold cyan
"title"        # Bold cyan
"subtitle"     # Muted gray
"file"         # Cool cyan
"path"         # Electric mint

# Code
"code.keyword"   # Cool cyan
"code.function"  # Electric mint
"code.string"    # Gold
"code.number"    # Purple
"code.comment"   # Muted gray
"code.type"      # Cool cyan
"code.operator"  # Electric mint
```

### Usage Examples

```python
from toolkit_theme import get_console, print_status, print_section

console = get_console()

# Status messages
print_status("Operation completed successfully", "success")
print_status("Spec drift detected", "warning")
print_status("Connection failed", "error")
print_status("Processing 15 files", "info")

# Section headers
print_section("Phase 3 Complete", "All tests passing with 100% coverage")

# Custom styling
console.print("File: [file]src/main.py[/]")
console.print("Path: [path]/Users/dev/project[/]")
console.print("[code.keyword]def[/] [code.function]process[/]([code.type]str[/]):")
```

## 📐 Design Principles

1. **High Contrast**: Deep black background with vibrant accents for terminal clarity
2. **Semantic Color**: Each color has meaning (mint=success, cyan=system, gold=caution)
3. **Consistent Icons**: Status bullets create visual hierarchy
4. **Readability First**: 1.3 line height, clear typography, spacious layout
5. **Cinematic Feel**: Subtle "AI typing" effects, gradient progress bars
6. **Accessibility**: WCAG AA compliant contrast ratios for all text

## 🔧 Customization

### Custom Console Configuration

```python
from rich.console import Console
from rich.theme import Theme

custom_theme = Theme({
    "success": "bold #00FF88",
    "my_custom": "#YOUR_COLOR",
})

console = Console(theme=custom_theme)
```

### Environment Variables

```bash
# Force color output even when piping
export FORCE_COLOR=1

# Set Rich theme
export RICH_THEME=neural-black
```

## 🚀 Testing the Theme

Run the theme demo to see all colors and styles:

```bash
python toolkit_theme.py
```

This will display:
- All status indicators
- Diff examples
- Code syntax highlighting
- Progress bars
- Icons and symbols

## 📝 License

Part of the RyCode project. Same license as main repository.

## 🤝 Contributing

To propose theme improvements:

1. Test changes across multiple terminals
2. Ensure WCAG AA contrast compliance
3. Maintain semantic color meaning
4. Submit with screenshots showing before/after

---

**Neural Black** — Where terminal realism meets AI aesthetics.
