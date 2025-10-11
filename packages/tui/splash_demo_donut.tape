# VHS Tape for RyCode Infinite Donut Mode Demo
# Install VHS: https://github.com/charmbracelet/vhs
# Run with: vhs splash_demo_donut.tape

# Output settings
Output splash_demo_donut.gif
Set Shell bash
Set FontSize 14
Set Width 1200
Set Height 800
Set Padding 20
Set Theme "Dracula"

# Build RyCode
Type "go build -o rycode ./cmd/rycode"
Enter
Sleep 3s

# Clear screen
Type "clear"
Enter
Sleep 0.5s

# Run RyCode in infinite donut mode
Type "./rycode donut"
Enter

# Let donut spin for 10 seconds
Sleep 10s

# Press ? to show math equations
Type "?"
Sleep 5s

# Press ? again to hide
Type "?"
Sleep 3s

# Try Konami code for rainbow mode
Type@100ms "↑"
Type@100ms "↑"
Type@100ms "↓"
Type@100ms "↓"
Type@100ms "←"
Type@100ms "→"
Type@100ms "←"
Type@100ms "→"
Type "b"
Type "a"

# Rainbow mode active - show for 5 seconds
Sleep 5s

# Quit with Q
Type "q"
