# RyCode Matrix TUI v2

> The AI-Native, Mobile-First Terminal IDE with Matrix Cyberpunk Aesthetics

[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue)](https://golang.org/doc/install)
[![Tests](https://img.shields.io/badge/tests-134%20passing-brightgreen)](https://github.com/aaronmrosenthal/RyCode)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

## 🎯 Vision

RyCode Matrix TUI is the **most stunning, intuitive, and powerful terminal user interface** ever built for developers. Combining **Matrix-themed cyberpunk aesthetics** with **revolutionary mobile-first UX**, **gesture-based interactions**, and **AI-native workflows** to create an IDE that works seamlessly from phone to desktop.

This isn't just a TUI - it's **the future of coding**.

---

## ✨ Features

### 🎨 Matrix Cyberpunk Theme
- **20+ Neon Colors**: Matrix Green, Neon Cyan, Pink, Purple, Yellow, Orange
- **10+ Visual Effects**: Gradients, glow, pulse, Matrix rain, scanlines
- **Semantic Color System**: Error (pink), warning (yellow), success (green), info (cyan)
- **Syntax Highlighting**: 200+ languages via Chroma

### 📱 Mobile-First Responsive Design
- **6 Device Classes**: PhonePortrait → PhoneLandscape → Tablet → Desktop
- **Automatic Adaptation**: UI adapts to terminal size (40-160+ columns)
- **Smart Layouts**: Stack on mobile, split on tablet, multi-pane on desktop
- **Touch-Friendly**: Large targets, clear focus indicators

### 🗂️ File Tree Navigation
- **Vim-Style Shortcuts**: j/k navigate, h/l expand/collapse, g/G first/last
- **12+ File Type Icons**: Go 🔷, JS 📜, Python 🐍, Rust 🦀, JSON 📋, etc.
- **Git Status Indicators**: ?, M, A, D, R, ✓, • (color-coded)
- **Show/Hide Hidden Files**: Toggle with `.`
- **Smart Scrolling**: Auto-scroll to keep selection visible

### 💬 Interactive Chat Interface
- **Real AI Providers**: Claude Opus 4 & GPT-4o with streaming responses
- **Auto-Provider Selection**: Automatically uses available API keys
- **Mock AI Fallback**: Works without API keys for demos
- **Markdown Rendering**: Beautiful code blocks, lists, quotes
- **Ghost Text Suggestions**: Tab to accept predictions
- **Quick Actions**: Fix, Test, Explain, Refactor, Run buttons
- **15+ Keyboard Shortcuts**: Vim-style navigation

### 🔄 Workspace Management
- **Split-Pane Layout**: FileTree + Chat side-by-side
- **Focus Switching**: Ctrl+B to toggle between panes
- **Toggle Visibility**: Ctrl+T to show/hide FileTree
- **Adaptive Layout**: Auto-hide tree on mobile devices

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** ([install](https://golang.org/doc/install))
- **Git** (for version control features)
- **Terminal**: Any modern terminal (iTerm2, Alacritty, Windows Terminal, etc.)
- **AI API Key** (optional): [Claude](https://console.anthropic.com/) or [OpenAI](https://platform.openai.com/api-keys)

### Installation

```bash
# Clone the repository
git clone https://github.com/aaronmrosenthal/RyCode.git
cd RyCode/packages/tui-v2

# Install dependencies
make deps

# Build the binary
make build

# Run the workspace (FileTree + Chat)
make workspace

# Or run directly
../../packages/rycode/dist/rycode
```

### Alternative: Install to ~/bin

```bash
# Install to user bin
make install

# Run from anywhere
rycode
```

### Enable Real AI (Optional)

```bash
# For Claude (Anthropic)
export ANTHROPIC_API_KEY="sk-ant-..."

# For GPT-4 (OpenAI)
export OPENAI_API_KEY="sk-..."

# The TUI will auto-detect and use the first available key
# See AI_INTEGRATION.md for full details
```

---

## 📖 Usage

### Running Modes

```bash
# Workspace mode (default) - FileTree + Chat
rycode
rycode --workspace

# Chat only (no file tree)
rycode --chat

# Theme demo (showcase colors and effects)
rycode --demo

# Show help
rycode --help
make help
```

### Keyboard Shortcuts

#### Global
| Key | Action |
|-----|--------|
| `Ctrl+C` / `Esc` | Quit application |
| `Ctrl+B` | Switch focus (FileTree ↔ Chat) |
| `Ctrl+T` | Toggle FileTree visibility |

#### FileTree (when focused)
| Key | Action |
|-----|--------|
| `j` / `↓` | Select next |
| `k` / `↑` | Select previous |
| `g` | Go to first item |
| `G` | Go to last item |
| `h` / `←` / `Backspace` | Go to parent / Collapse folder |
| `l` / `→` / `Enter` | Expand folder / Open file |
| `.` | Toggle hidden files |
| `r` | Refresh file tree |
| `o` | Open selected file |

#### Chat (when focused)
| Key | Action |
|-----|--------|
| `Enter` | Send message |
| `Tab` | Accept ghost text suggestion |
| `Backspace` | Delete character before cursor |
| `Delete` | Delete character after cursor |
| `←` / `→` | Move cursor left/right |
| `Home` / `Ctrl+A` | Move cursor to start |
| `End` / `Ctrl+E` | Move cursor to end |
| `↑` / `↓` | Scroll messages up/down |
| `Ctrl+D` | Scroll to bottom |
| `Ctrl+L` | Clear all messages |

---

## 🧪 Development

### Running Tests

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run with coverage
make coverage

# View coverage report (generates coverage.html)
open coverage.html
```

### Test Statistics
- **Total Tests**: 134 (all passing ✅)
- **Coverage**: 87.7% (layout), 87.8% (components), 90.2% (models)
- **Test Files**: 8
- **Test Runtime**: <1s

### Building

```bash
# Build binary
make build

# Clean build artifacts
make clean

# Format code
make fmt

# Run linter
make lint

# Tidy dependencies
make tidy
```

### Project Structure

```
packages/tui-v2/
├── cmd/rycode/              # Main entry point
│   ├── main.go             # CLI flags and initialization
│   └── demo.go             # Theme demo showcase
├── internal/
│   ├── layout/             # Responsive layout system
│   │   ├── types.go        # DeviceClass enum (6 breakpoints)
│   │   └── manager.go      # LayoutManager for adaptation
│   ├── theme/              # Matrix cyberpunk theme
│   │   ├── colors.go       # 20+ color palette
│   │   ├── theme.go        # Theme system & styles
│   │   └── effects.go      # Visual effects (10+)
│   ├── ui/
│   │   ├── components/     # Reusable UI components
│   │   │   ├── message.go  # MessageBubble (markdown, code)
│   │   │   ├── input.go    # InputBar (ghost text, actions)
│   │   │   └── filetree.go # FileTree (navigation, git)
│   │   └── models/         # Application models
│   │       ├── chat.go     # ChatModel (streaming AI)
│   │       └── workspace.go # WorkspaceModel (split-pane)
├── Makefile                # Build automation
├── go.mod                  # Go dependencies
└── README.md               # This file
```

---

## 🎨 Theme System

### Color Palette

**Primary Matrix Colors:**
- `MatrixGreen` (#00ff00) - Primary UI elements
- `MatrixGreenDim` (#00dd00) - Secondary text
- `MatrixGreenDark` (#004400) - Backgrounds

**Neon Accents:**
- `NeonCyan` (#00ffff) - Informational
- `NeonPink` (#ff3366) - Errors
- `NeonPurple` (#cc00ff) - Types
- `NeonYellow` (#ffaa00) - Warnings
- `NeonOrange` (#ff6600) - Modified files
- `NeonBlue` (#0088ff) - Functions

**Code Syntax:**
- Keywords → Neon Pink
- Strings → Neon Yellow
- Numbers → Neon Cyan
- Comments → Dark Green
- Functions → Neon Blue
- Types → Neon Purple

### Visual Effects

- **Gradient Text**: 4 presets (Matrix, Fire, Cool, Warm)
- **Glow Effects**: Intensity-based neon glow
- **Matrix Rain**: Animated digital rain
- **Pulse Animation**: Breathing effect
- **Rainbow Text**: Multi-color cycling
- **Scanlines**: CRT monitor effect

---

## 📊 Architecture

### Design Patterns

**Bubble Tea (Elm Architecture):**
- Model-View-Update pattern
- Immutable state updates
- Command-based side effects
- Type-safe message passing

**Component-Based:**
- Reusable UI components
- Isolated state management
- Composition over inheritance
- Clear interfaces

**Responsive:**
- Device class detection
- Breakpoint-based layouts
- Dynamic dimension updates
- Mobile-first design

**Theme-Driven:**
- Centralized color system
- Consistent styling
- Easy customization
- Semantic colors

---

## 🔮 Roadmap

### Phase 1: Foundation ✅ (Complete)
- [x] Responsive framework (6 device classes)
- [x] Matrix theme system (20+ colors, 10+ effects)
- [x] MessageBubble component (markdown, code blocks)
- [x] InputBar component (ghost text, quick actions)
- [x] ChatModel (streaming AI responses)
- [x] FileTree component (navigation, git status)
- [x] Workspace integration (split-pane)
- [x] 134 tests (100% passing)

### Phase 2: AI Integration ✅ (Complete!)
- [x] Real AI provider integration (Claude Opus 4, GPT-4o)
- [x] Streaming token-by-token responses
- [x] Context-aware multi-turn conversations
- [x] Auto-provider selection (Claude → OpenAI → Mock)
- [x] Conversation history tracking
- [x] Error handling & graceful fallback
- [ ] Token usage tracking & cost monitoring
- [ ] Rate limiting with retry logic

### Phase 3: Code Editor
- [ ] Syntax highlighting (200+ languages)
- [ ] LSP integration (go-to-definition, autocomplete)
- [ ] Multi-file editing (tabs)
- [ ] Search & replace
- [ ] Git integration (commit, diff, blame)
- [ ] Code actions (format, refactor)

### Phase 4: Advanced Features
- [ ] Voice input (speech-to-text)
- [ ] Gesture recognition (swipe, pinch)
- [ ] Plugin system
- [ ] Custom themes
- [ ] Session persistence
- [ ] Multi-workspace support

---

## 🤝 Contributing

We welcome contributions! Here's how you can help:

### Development Setup

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Format code (`make fmt`)
6. Run linter (`make lint`)
7. Commit changes (`git commit -m 'feat: Add amazing feature'`)
8. Push to branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Add godoc comments for exported types and functions
- Write tests for new features (maintain 80%+ coverage)
- Use semantic commit messages (feat:, fix:, docs:, test:, refactor:)

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

### Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [Chroma](https://github.com/alecthomas/chroma) - Syntax highlighting

### Inspiration

- [toolkit-cli.com](https://toolkit-cli.com) - Matrix theme inspiration
- [neovim](https://neovim.io) - Vim-style shortcuts
- [VSCode](https://code.visualstudio.com) - IDE features
- [The Matrix](https://www.imdb.com/title/tt0133093/) - Cyberpunk aesthetic

---

## 📞 Support

- **Documentation**:
  - [README.md](README.md) - Getting started & features
  - [AI_INTEGRATION.md](AI_INTEGRATION.md) - AI providers guide
  - [ARCHITECTURE.md](ARCHITECTURE.md) - Technical architecture
  - [SHOWCASE.md](SHOWCASE.md) - Visual demos
- **Issues**: [GitHub Issues](https://github.com/aaronmrosenthal/RyCode/issues)
- **Discussions**: [GitHub Discussions](https://github.com/aaronmrosenthal/RyCode/discussions)

---

<div align="center">

**Built with ❤️ by the RyCode Team**

*Making terminal coding beautiful, one green pixel at a time* 🟢

</div>
