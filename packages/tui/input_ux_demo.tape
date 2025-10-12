# RyCode Input UX Demo
# Showcases the polished input field with modern design

Output input_ux_demo.gif

Set FontSize 16
Set Width 1200
Set Height 800
Set Padding 20
Set Theme "Catppuccin Mocha"

Type "/tmp/rycode-polished"
Enter
Sleep 2s

# Show the beautiful placeholder text
Screenshot input_placeholder.png
Sleep 1s

# Type a prompt to show the modern input styling
Type "Help me build a REST API"
Sleep 500ms
Screenshot input_typing.png
Sleep 500ms

# Submit the prompt
Enter
Sleep 2s

# Show bash mode with $ prompt
Type "!"
Sleep 500ms
Screenshot input_bash_mode.png
Type "ls -la"
Sleep 500ms

# Cancel and go back to normal
Escape
Sleep 500ms

# Show the model selector with latest models
Type "/model"
Sleep 1s
Screenshot model_selector_latest.png
Sleep 2s

# Navigate through models
Down
Down
Sleep 500ms
Screenshot model_selector_navigation.png
