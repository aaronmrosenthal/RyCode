# VHS Tape for RyCode Splash Screen Demo
# Install VHS: https://github.com/charmbracelet/vhs
# Run with: vhs splash_demo.tape

# Output settings
Output splash_demo.gif
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

# Run RyCode with splash enabled
Type "./rycode --splash"
Enter

# Let splash animation play
Sleep 6s

# Demo complete - the splash auto-closes or user can press any key
