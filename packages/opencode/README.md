# 🟢 RyCode

> **AI-Powered Development Assistant • Matrix Cyberpunk Aesthetic • Next-Gen CLI**

<div align="center">

```
  ██████╗ ██╗   ██╗ ██████╗ ██████╗ ██████╗ ███████╗
  ██╔══██╗╚██╗ ██╔╝██╔════╝██╔═══██╗██╔══██╗██╔════╝
  ██████╔╝ ╚████╔╝ ██║     ██║   ██║██║  ██║█████╗
  ██╔══██╗  ╚██╔╝  ██║     ██║   ██║██║  ██║██╔══╝
  ██║  ██║   ██║   ╚██████╗╚██████╔╝██████╔╝███████╗
  ╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝
```

**The most powerful AI coding assistant with a killer terminal experience**

[Features](#-features) • [Installation](#-installation) • [Quick Start](#-quick-start) • [Documentation](#-documentation)

</div>

---

## ⚡ Why RyCode?

Ship features **10x faster** with AI pair programming in a **stunning cyberpunk terminal** that actually makes coding fun again.

### 🎯 What Makes RyCode Different

- **🤖 Full Codebase Understanding** - AI reads your entire project for context-aware suggestions
- **💬 Natural Language Interface** - Just describe what you want to build
- **🎨 Killer Terminal UI** - Matrix digital rain aesthetic with modern TUI components
- **🔗 Clickable Everything** - File paths and URLs are clickable in modern terminals
- **✨ Professional Polish** - Every interaction is smooth, beautiful, and informative
- **🚀 Blazing Fast** - Built on Bun for maximum performance

---

## 🌟 Features

### 🎨 **Stunning Visual Experience**

<details>
<summary><b>Matrix Digital Rain Aesthetic</b></summary>

- **Killer ASCII Logo** with authentic Matrix green gradient (dark → bright → light cascade)
- **Animated Headers** with flowing digital rain effect
- **4 Logo Variants**: Modern, Big, Slant, Cyberpunk Boxed
- **Professional Color Palette** inspired by toolkit-cli.com
- Multiple shades of Matrix green for depth and dimension

</details>

<details>
<summary><b>Advanced TUI Components</b></summary>

**14 Professional Terminal UI Components:**

- 📊 **Progress Bars** - Customizable colors, labels, and animations
- ⚙️ **Spinners** - 7 styles (dots, pulse, cyber, matrix, neon, arrow, line)
- 📋 **Rich Tables** - Formatted with borders, colors, and alignment
- 📅 **Timeline Views** - Event history with status indicators
- 🧙 **Wizard Steps** - Multi-step process indicators
- 🔔 **Notifications** - Toast-style success/error/warning/info boxes
- 🔑 **Key-Value Displays** - Auto-aligned configuration views
- 💻 **Code Blocks** - Syntax highlighting with line numbers
- 🔀 **Diff Viewers** - Git-style code comparison
- 🏷️ **Badges & Tags** - Status indicators and labels
- 📊 **Status Indicators** - Online/offline/loading/error states
- 📁 **File Trees** - Visual directory structures
- 📝 **Interactive Menus** - Selection lists with descriptions
- 📦 **Collapsible Sections** - Expandable content blocks

</details>

<details>
<summary><b>Clickable Terminal Links</b></summary>

**OSC 8 Hyperlinks** - Modern terminal support for clickable links:

- 🔗 **URL Links** - Click to open in browser
- 📄 **File Links** - Click to open in editor
- ✨ **Auto-Detection** - Automatically makes paths/URLs clickable
- 🎨 **Styled Links** - Icons and colors for visual clarity
- ✅ **Wide Support** - Works in iTerm2, Alacritty, Kitty, WezTerm, Windows Terminal, VSCode

```typescript
// Clickable file path
UI.fileLink("/path/to/file.ts")

// Clickable URL
UI.link("Documentation", "https://docs.rycode.ai")

// Auto-link everything in text
UI.autoLink("Check /README.md and https://github.com")
```

</details>

### 🤖 **AI Superpowers**

<details>
<summary><b>Intelligent Code Generation</b></summary>

- **Natural Language → Code** - Describe features in plain English
- **Context-Aware** - AI understands your entire codebase
- **Multiple Providers** - Anthropic Claude, OpenAI, Google, OpenRouter
- **Smart Suggestions** - Proactive improvements and optimizations
- **Error Detection** - Catch bugs before they happen

</details>

<details>
<summary><b>AI Pair Programming</b></summary>

- **Interactive Chat** - Build features through conversation
- **Code Explanations** - Understand complex code instantly
- **Refactoring Assistant** - Clean up code with AI guidance
- **Bug Fixing** - AI helps diagnose and fix issues
- **Documentation Generation** - Auto-generate docs and comments

</details>

### 💎 **Professional Installer**

<details>
<summary><b>Polished Onboarding Experience</b></summary>

**Every step is beautifully crafted:**

- ✨ **Welcome Screen** - Animated Matrix header with logo
- 🎯 **Step-by-Step Wizard** - Clear progress through setup
- 🔑 **API Key Help** - Clickable links to provider signups
- 📊 **Progress Indicators** - Visual feedback during installation
- ✅ **System Checks** - Verify all connections and models
- 💡 **Quick Tips** - Helpful guidance for new users
- 🎉 **Success Celebration** - Engaging completion messages

**Professional messaging throughout:**
- Clear error messages with actionable solutions
- Informative warnings with context
- Encouraging success messages
- Consistent Matrix cyberpunk aesthetic

</details>

### 🛠️ **Developer Experience**

<details>
<summary><b>Productivity Features</b></summary>

- **🔍 Intelligent Search** - Find anything in your codebase instantly
- **📝 Smart Completions** - Context-aware code suggestions
- **🔧 Refactoring Tools** - Rename, extract, optimize
- **🐛 Debugging Assistant** - AI-powered error analysis
- **📚 Documentation** - Auto-generate and update docs
- **⚡ Slash Commands** - Quick actions for common tasks

</details>

<details>
<summary><b>What's Next Component</b></summary>

**Actionable next steps after every operation:**

- 📋 **Clear Options** - See what you can do next
- 🔗 **Clickable Files** - Jump directly to relevant code
- 📅 **Activity Timeline** - Recent changes and actions
- 🎯 **Smart Recommendations** - AI suggests next steps
- 📊 **Progress Tracking** - Visual task completion

</details>

---

## 🚀 Installation

### Prerequisites

- **Bun** v1.2.12 or later ([Install Bun](https://bun.sh))
- **Modern Terminal** (iTerm2, Alacritty, Kitty, WezTerm, Windows Terminal, or VSCode)
- **AI Provider API Key** (Anthropic, OpenAI, Google, etc.)

### Quick Install

```bash
# Clone the repository
git clone https://github.com/rycode/opencode.git
cd opencode/packages/opencode

# Install dependencies
bun install

# Run RyCode
bun run index.ts
```

### First Run

RyCode will guide you through:

1. **🤖 Provider Selection** - Choose your AI provider (we recommend Anthropic or OpenAI)
2. **🔑 API Key Setup** - Enter your credentials (clickable signup links provided)
3. **✅ System Check** - Verify all connections and models
4. **🎉 Ready to Code** - Start building with AI!

---

## 💡 Quick Start

### Basic Usage

```bash
# Start RyCode
rycode

# Chat with AI
> "Create a React component for a user profile card"

# Get help
> "Explain this function"

# Refactor code
> "Refactor this to use async/await"

# Fix bugs
> "Why is this test failing?"
```

### Clickable Links

All file paths and URLs are clickable in modern terminals:

```
Modified: src/components/UserProfile.tsx  ← Click to open
Documentation: https://docs.rycode.ai     ← Click to visit
```

Just **Cmd+Click** (Mac) or **Ctrl+Click** (Windows/Linux) on any link!

### Visual Feedback

Every operation shows beautiful progress:

```
Step 2/3 [████████████████░░░░░░░░] 67%
  Verifying credentials...

✓ Authentication Successful
  Connected to Anthropic Claude
```

---

## 📚 Documentation

### Comprehensive Guides

- **[Clickable Links Guide](./CLICKABLE_LINKS.md)** - Complete OSC 8 implementation
- **[TUI Enhancements](./TUI_ENHANCEMENTS.md)** - All 14 components documented
- **[Installer Experience](./INSTALLER_EXPERIENCE.md)** - Professional onboarding flow
- **[Plugin Security](./PLUGIN_SECURITY.md)** - Enterprise-grade plugin security
- **[Plugin Registry](./PLUGIN_REGISTRY.md)** - Verified plugin registry system
- **[Plugin Signatures](./PLUGIN_SIGNATURES.md)** - GPG and crypto signature verification

### Component Examples

```typescript
import { UI } from "./src/cli/ui"
import { EnhancedTUI } from "./src/cli/tui-enhanced"
import { InstallerMessages } from "./src/cli/installer-messages"

// Show Matrix logo
console.log(UI.logo())

// Progress bar
console.log(EnhancedTUI.progressBar(75, 100, {
  label: "Building",
  color: UI.Style.MATRIX_GREEN
}))

// Professional notification
console.log(EnhancedTUI.notification(
  "Build completed successfully!",
  "success"
))

// Formatted table
console.log(EnhancedTUI.table(
  ["Task", "Status", "Progress"],
  [
    ["Build", "✓", "100%"],
    ["Tests", "⟳", "75%"]
  ]
))
```

### Demos

Run the demos to see everything in action:

```bash
# See all TUI components
bun run demo-tui-enhanced.ts

# See Matrix logo variants
bun run demo-matrix-logo.ts

# See installer flow
bun run demo-installer.ts

# See clickable links
bun run verify-links.ts
```

---

## 🎨 Visual Showcase

### Matrix Color Palette

```
DARK    #00641E  ████  ← Shadows
DIM     #00B432  ████
BRIGHT  #00FF41  ████  ← Classic Matrix
LIGHT   #64FF96  ████
BRIGHT  #96FFB4  ████  ← Highlight
```

### Component Gallery

**Progress Indicators**
```
Loading [████████████░░░░░░░░] 60%
```

**Status Icons**
```
● Online    ◔ Loading    ● Error
```

**Timeline View**
```
  ✓ 10:23 Build Started
  │   Compiling TypeScript...
  │
  ✓ 10:24 Tests Passed
      All 142 tests completed
```

**Wizard Steps**
```
✓ Step 1: Configure Project
│
✓ Step 2: Install Dependencies
│
▶ Step 3: Run Tests
│
○ Step 4: Deploy
```

---

## 🤝 Contributing

We welcome contributions! RyCode is built with:

- **TypeScript** - Type-safe codebase
- **Bun** - Fast runtime and package manager
- **Modern Terminal APIs** - OSC 8, ANSI colors, Unicode
- **Professional UX** - Every detail matters

### Development Setup

```bash
# Clone and install
git clone https://github.com/rycode/opencode.git
cd opencode/packages/opencode
bun install

# Run in development
bun run index.ts

# Run tests
bun test

# Build
bun run build
```

---

## 📄 License

MIT License - See [LICENSE](./LICENSE) for details

---

## 🌟 Star History

If RyCode makes your development faster and more enjoyable, give us a star! ⭐

---

<div align="center">

### 🟢 Built with Matrix Digital Rain Aesthetic

**Made with 💚 by developers, for developers**

[Get Started](#-installation) • [View Demos](./demos) • [Read Docs](./docs)

</div>
