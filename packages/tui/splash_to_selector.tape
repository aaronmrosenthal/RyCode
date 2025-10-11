# VHS Tape for RyCode - Splash Screen Reveals Model Selector
# Shows splash screen then reveals the model selector interface
# Run with: vhs splash_to_selector.tape

Output splash_to_selector.gif
Set Shell bash
Set FontSize 14
Set Width 1200
Set Height 800
Set Padding 20
Set Theme "Dracula"

# Navigate to rycode directory
Type "cd /Users/aaron/Code/RyCode/RyCode/packages/rycode"
Enter
Sleep 1s

# Build RyCode
Type "go build -o rycode ./cmd/rycode"
Enter
Sleep 3s

# Clear screen
Type "clear"
Enter
Sleep 0.5s

# Launch RyCode (WITH splash screen)
Type "./rycode"
Enter

# Wait for splash to play (adjust timing based on your splash duration)
Sleep 5s

# Now the model selector/TUI should be visible
# Show it for a few seconds
Sleep 3s

# Type /model to show the model selector
Type "/model"
Sleep 2s

# Show the model list
Sleep 3s

# Press Escape to go back
Type@100ms "\x1b"  # Escape key
Sleep 1s

# Quit
Ctrl+C
Sleep 0.5s
