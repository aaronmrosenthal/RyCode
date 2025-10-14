#!/bin/bash
# Generate visual examples of all provider themes
# Requires: VHS (https://github.com/charmbracelet/vhs)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
OUTPUT_DIR="$PROJECT_ROOT/packages/tui/docs/visuals"

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo "=== Theme Visual Generator ==="
echo "Project root: $PROJECT_ROOT"
echo "Output dir: $OUTPUT_DIR"
echo ""

# Check if VHS is installed
if ! command -v vhs &> /dev/null; then
    echo "❌ VHS is not installed!"
    echo ""
    echo "Install VHS:"
    echo "  brew install vhs"
    echo "  # or"
    echo "  go install github.com/charmbracelet/vhs@latest"
    echo ""
    exit 1
fi

echo "✓ VHS is installed"
echo ""

# Theme colors for reference
declare -A THEME_COLORS
THEME_COLORS[claude]="#D4754C"
THEME_COLORS[gemini]="#4285F4"
THEME_COLORS[codex]="#10A37F"
THEME_COLORS[qwen]="#FF6A00"

# Generate VHS tape for each theme
for theme in claude gemini codex qwen; do
    echo "Generating tape for $theme theme..."

    TAPE_FILE="$OUTPUT_DIR/${theme}_theme.tape"

    cat > "$TAPE_FILE" <<EOF
# Theme: $theme
# Primary color: ${THEME_COLORS[$theme]}

Output $OUTPUT_DIR/${theme}_theme.gif
Output $OUTPUT_DIR/${theme}_theme.png

Set FontSize 14
Set Width 1200
Set Height 800
Set Theme "dark"

# Launch RyCode (mock for visual generation)
Type "echo '╭────────────────────────────────────────────────╮'"
Enter
Type "echo '│                                                │'"
Enter
Type "echo '│  RyCode TUI - $theme Theme                     │'"
Enter
Type "echo '│                                                │'"
Enter
Type "echo '│  Primary Color: ${THEME_COLORS[$theme]}                  │'"
Enter
Type "echo '│                                                │'"
Enter
Type "echo '│  ✓ Accessibility: WCAG AA Compliant           │'"
Enter
Type "echo '│  ⚡ Performance: 317ns theme switching         │'"
Enter
Type "echo '│  🎨 Provider-specific branding                 │'"
Enter
Type "echo '│                                                │'"
Enter
Type "echo '╰────────────────────────────────────────────────╯'"
Enter
Sleep 1s

Screenshot $OUTPUT_DIR/${theme}_theme.png
EOF

    echo "  Created: $TAPE_FILE"
done

echo ""
echo "=== Generating Screenshots ==="
echo ""

# Generate all visuals
for theme in claude gemini codex qwen; do
    TAPE_FILE="$OUTPUT_DIR/${theme}_theme.tape"
    echo "Running VHS for $theme..."
    vhs "$TAPE_FILE"
    echo "  ✓ Generated: ${theme}_theme.gif"
    echo "  ✓ Generated: ${theme}_theme.png"
done

echo ""
echo "=== Theme Comparison ==="
echo ""

# Create comparison tape showing all themes
COMPARISON_TAPE="$OUTPUT_DIR/theme_comparison.tape"

cat > "$COMPARISON_TAPE" <<'EOF'
Output theme_comparison.gif

Set FontSize 14
Set Width 1200
Set Height 900
Set Theme "dark"

Type "# RyCode Theme System - All Providers"
Enter
Enter
Type "echo 'Press Tab to cycle through themes:'"
Enter
Sleep 500ms

# Claude
Type "echo ''"
Enter
Type "echo '╭─ CLAUDE ──────────────────────────────────────╮'"
Enter
Type "echo '│ 🟠 Warm copper orange (#D4754C)               │'"
Enter
Type "echo '│ Developer-friendly, approachable aesthetic    │'"
Enter
Type "echo '╰───────────────────────────────────────────────╯'"
Enter
Sleep 1s

# Gemini
Type "echo ''"
Enter
Type "echo '╭─ GEMINI ──────────────────────────────────────╮'"
Enter
Type "echo '│ 🔵 Blue-pink gradient (#4285F4 → #EA4335)    │'"
Enter
Type "echo '│ Modern, vibrant, AI-forward aesthetic         │'"
Enter
Type "echo '╰───────────────────────────────────────────────╯'"
Enter
Sleep 1s

# Codex
Type "echo ''"
Enter
Type "echo '╭─ CODEX ───────────────────────────────────────╮'"
Enter
Type "echo '│ 🟢 OpenAI teal (#10A37F)                      │'"
Enter
Type "echo '│ Professional, technical, precise aesthetic    │'"
Enter
Type "echo '╰───────────────────────────────────────────────╯'"
Enter
Sleep 1s

# Qwen
Type "echo ''"
Enter
Type "echo '╭─ QWEN ────────────────────────────────────────╮'"
Enter
Type "echo '│ 🟠 Alibaba orange (#FF6A00)                   │'"
Enter
Type "echo '│ Modern, innovative, international aesthetic   │'"
Enter
Type "echo '╰───────────────────────────────────────────────╯'"
Enter
Sleep 2s
EOF

echo "Generating theme comparison..."
cd "$OUTPUT_DIR"
vhs "$COMPARISON_TAPE"
echo "  ✓ Generated: theme_comparison.gif"

echo ""
echo "=== Visual Assets Created ==="
echo ""
ls -lh "$OUTPUT_DIR"/*.{gif,png} 2>/dev/null || echo "No files generated"
echo ""
echo "✅ Visual generation complete!"
echo ""
echo "Generated files:"
echo "  - 4 theme GIFs (animated)"
echo "  - 4 theme PNGs (static)"
echo "  - 1 comparison GIF"
echo ""
echo "Location: $OUTPUT_DIR"
