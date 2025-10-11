# VHS Tape for RyCode - Splash Screen to Model Selector
# Shows the 5 SOTA models in the selector interface
# Run with: vhs splash_model_selector.tape

Output splash_model_selector.gif
Set Shell bash
Set FontSize 14
Set Width 1200
Set Height 800
Set Padding 20
Set Theme "Dracula"

# Change to rycode directory
Type "cd /Users/aaron/Code/RyCode/RyCode/packages/rycode"
Enter
Sleep 1s

# Clear screen
Type "clear"
Enter
Sleep 0.5s

# Run rycode (this should show splash then TUI)
Type "./dist/rycode-darwin-arm64/bin/rycode"
Enter

# Wait for splash animation to complete
Sleep 6s

# Now at the TUI prompt - type /model to show selector
Type "/model"
Sleep 2s

# Model selector should be visible now - hold for viewing
Sleep 5s

# Exit gracefully
Escape
Sleep 0.5s
Ctrl+C
Sleep 0.5s
