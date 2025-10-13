#!/bin/bash
# Simple splash screen demo script
# Creates a basic recording using macOS built-in tools
# No external dependencies required

set -e

PROJECT_ROOT="/Users/aaron/Code/RyCode/RyCode/packages/tui"
cd "$PROJECT_ROOT"

echo "🎬 RyCode Splash Screen Simple Recording"
echo "=========================================="
echo ""

# Check if rycode binary exists
if [ ! -f "./rycode" ]; then
  echo "📦 Building RyCode..."
  go build -o rycode ./cmd/rycode
  echo "✅ Build complete"
  echo ""
fi

echo "📸 Recording Instructions:"
echo ""
echo "This script will run the splash screen for you."
echo "Use one of these methods to capture it:"
echo ""
echo "Method 1: macOS Screenshot (Cmd+Shift+5)"
echo "  1. Press Cmd+Shift+5"
echo "  2. Select 'Record Selected Portion'"
echo "  3. Select the terminal window"
echo "  4. Click 'Record'"
echo "  5. Press Enter below to start splash"
echo "  6. After 6 seconds, stop recording (Stop button in menu bar)"
echo ""
echo "Method 2: QuickTime Screen Recording"
echo "  1. Open QuickTime Player"
echo "  2. File → New Screen Recording"
echo "  3. Select terminal area"
echo "  4. Press Enter below to start splash"
echo "  5. After 6 seconds, stop recording (Cmd+Control+Esc)"
echo ""
echo "Method 3: Simple GIF (requires ImageMagick)"
echo "  This will take 10 screenshots and combine them"
echo ""

read -p "Press Enter to start splash screen (or Ctrl+C to cancel)..."

# Clear screen for clean recording
clear

echo "Starting in 3..."
sleep 1
echo "Starting in 2..."
sleep 1
echo "Starting in 1..."
sleep 1
clear

# Run splash
./rycode --splash

echo ""
echo "✅ Splash complete!"
echo ""
echo "Your recording is saved according to the method you used:"
echo "  - Screenshot tool: ~/Desktop/Screen Recording YYYY-MM-DD.mov"
echo "  - QuickTime: Saved to location you chose"
echo ""
