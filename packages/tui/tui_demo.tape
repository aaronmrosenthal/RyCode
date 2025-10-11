# VHS Tape for RyCode TUI Demo - Multi-Agent Interface
# Shows the core feature: switching between AI agents
# Run with: vhs tui_demo.tape

Output tui_demo.gif
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

# Launch RyCode (without splash to show TUI immediately)
Type "./rycode --no-splash"
Enter

# Show TUI with agent tabs for 10 seconds
Sleep 10s

# Simulate Tab key to switch agents
Type@200ms "	"  # Tab key
Sleep 2s

Type@200ms "	"  # Tab again
Sleep 2s

Type@200ms "	"  # Tab again
Sleep 2s

# Show typing in prompt
Type "help"
Sleep 1s

# Quit
Type "q"
