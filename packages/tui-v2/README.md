# RyCode Matrix TUI v2

**The most stunning, intuitive, and powerful terminal user interface for AI-powered development.**

## Features

- 🎨 **Matrix Theme**: Cyberpunk aesthetics with neon green glow effects
- 📱 **Mobile-First**: Productive coding on phones, tablets, and desktops
- 👆 **Gesture-Driven**: Swipe, tap, pinch navigation
- 🎤 **Voice Input**: 95%+ accuracy voice commands
- 🤖 **Multi-Agent AI**: Claude, GPT-4, Gemini working together
- ⚡ **60 FPS**: Buttery smooth animations
- ♿ **Accessible**: WCAG 2.1 AAA compliant

## Quick Start

```bash
# Build
make build

# Run
../../packages/rycode/dist/rycode

# Demo mode
../../packages/rycode/dist/rycode --demo
```

## Development

```bash
# Install dependencies
go mod download

# Run tests
make test

# Run linter
make lint

# Generate coverage
make coverage
```

## Architecture

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [Chroma](https://github.com/alecthomas/chroma) - Syntax highlighting

## Project Structure

```
packages/tui-v2/
├── cmd/rycode/          # Entry point
├── internal/
│   ├── ui/              # UI models, components, views
│   ├── input/           # Gesture, voice, keyboard handling
│   ├── ai/              # AI providers and routing
│   ├── theme/           # Matrix theme and effects
│   ├── layout/          # Responsive layout system
│   ├── animation/       # Animation engine
│   └── util/            # Utilities
├── pkg/api/             # Public API
└── test/                # Tests
```

## License

MIT
