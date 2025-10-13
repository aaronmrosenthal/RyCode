# RyCode - AI-Powered IDE Context Package

## 🚀 Project Overview

**RyCode** is an AI-generated IDE built with modern web technologies and Go. It provides an intelligent terminal user interface (TUI) for AI-assisted coding with multi-provider support, OAuth authentication, and intelligent model management.

**Repository**: https://github.com/aaronmrosenthal/RyCode
**License**: MIT
**Version**: 0.14.1
**Built With**: toolkit-cli

---

## 📁 Project Architecture

### Monorepo Structure

```
RyCode/
├── packages/
│   ├── tui/                    # Go-based Terminal UI (Primary Interface)
│   ├── rycode/                 # TypeScript Backend Server & CLI
│   ├── desktop/                # SolidJS Desktop Application
│   ├── web/                    # Web Application
│   ├── ry-code-landing/        # Next.js Landing Page
│   ├── plugin/                 # Plugin System
│   ├── sdk/
│   │   ├── go/                 # Go SDK (Stainless-generated)
│   │   └── js/                 # JavaScript SDK
│   ├── console/                # Admin Console
│   └── function/               # Serverless Functions
├── sdks/
│   └── vscode/                 # VS Code Extension
├── docs/                       # Documentation
├── test-tui/                   # TUI Tests
└── infra/                      # Infrastructure as Code
```

---

## 🏗️ Core Components

### 1. TUI (Terminal User Interface)
**Location**: `packages/tui/`
**Language**: Go 1.24.0
**Framework**: Bubble Tea v2, Lip Gloss v2

**Key Features**:
- Epic splash screen with neural cortex animation
- Intelligent model selector with keyboard navigation
- Real-time authentication status
- Background auto-detection of API credentials
- Smart model recommendations based on task type
- Responsive design with lipgloss styling

**Entry Point**: `packages/tui/cmd/rycode/main.go`

**Architecture**:
```
internal/
├── app/                # Core application logic
├── auth/              # Authentication bridge
├── components/        # Reusable UI components
├── splash/            # Epic splash screen
├── tui/               # TUI orchestration
├── intelligence/      # AI-powered features
└── responsive/        # Responsive layout system
```

**Key Technologies**:
- `charmbracelet/bubbletea`: TUI framework
- `charmbracelet/lipgloss`: Terminal styling
- `charmbracelet/bubbles`: UI components
- `aaronmrosenthal/rycode-sdk-go`: API client

### 2. RyCode Server
**Location**: `packages/rycode/`
**Language**: TypeScript 5.8.2
**Runtime**: Bun 1.2.21

**Key Features**:
- REST API server (Hono framework)
- Multi-provider AI model support
- OAuth authentication system
- Plugin system with sandboxing
- LSP (Language Server Protocol) integration
- MCP (Model Context Protocol) support
- File watching and security scanning
- Session management with compaction

**Entry Point**: `packages/rycode/src/index.ts`

**Architecture**:
```
src/
├── server/            # Hono API server
├── auth/              # OAuth & API key auth
├── provider/          # Model provider integrations
├── plugin/            # Plugin system & security
├── lsp/               # LSP client/server
├── mcp/               # Model Context Protocol
├── session/           # Session & message management
├── tool/              # AI tools (read, edit, grep, etc.)
├── file/              # File operations & security
├── cli/               # CLI commands
└── util/              # Utilities
```

**Key Technologies**:
- `hono`: Fast web framework
- `ai`: Vercel AI SDK (5.0.8)
- `zod`: Schema validation (4.1.8)
- `@modelcontextprotocol/sdk`: MCP support
- `vscode-jsonrpc`: LSP communication
- `chokidar`: File watching

### 3. Desktop Application
**Location**: `packages/desktop/`
**Language**: TypeScript
**Framework**: SolidJS 1.9.9

**Key Features**:
- Modern SolidJS UI
- Agent-based workflows
- TailwindCSS 4.x styling
- Custom theme system

**Guidelines** (from `AGENTS.md`):
- Use `@/` import alias for `src/` directory
- Prettier: no semicolons, 120 char line width
- UI primitives in `/ui/`, components in `/components/`
- PascalCase for components, camelCase for variables

---

## 🔐 Authentication System

### Priority Order (from `docs/AUTH_PRIORITY.md`)

1. **Environment Variables** (backup)
2. **CLI Auth** (`~/.local/share/rycode/auth.json`) ✅ **PRIMARY**
3. Custom Provider Loaders
4. Plugin Auth
5. Config File (`opencode.json`)

### OAuth Support (from `docs/OAUTH_AUTHENTICATION.md`)

**Supported Providers**:
- ✅ Anthropic (Claude Pro/Max)
- ✅ GitHub Copilot
- ✅ Google (OAuth or API Key)

**OAuth Flow**:
```
User → rycode auth login
     → Browser opens / Device code shown
     → User authenticates
     → Tokens stored in ~/.local/share/rycode/auth.json
     → Auto-refresh on expiry
```

**Token Storage**:
```json
{
  "anthropic": {
    "type": "oauth",
    "access": "...",
    "refresh": "...",
    "expires": "2025-10-12T12:00:00Z"
  }
}
```

### API Key Support

**Environment Variables**:
```bash
ANTHROPIC_API_KEY / CLAUDE_API_KEY
OPENAI_API_KEY
GOOGLE_API_KEY
QWEN_API_KEY / DASHSCOPE_API_KEY
XAI_API_KEY / GROK_API_KEY
```

**CLI Commands**:
```bash
rycode auth login     # Interactive setup
rycode auth list      # Show configured providers
rycode auth logout    # Remove credentials
rycode models         # List available models
```

---

## 🤖 Intelligent Model Management

**Location**: `packages/tui/INTELLIGENT_MODEL_AUTOMATION.md`

### Auto-Setup on First Run
- Detects credentials in environment
- Authenticates automatically
- Shows success toast
- Zero user interruption

### Background Authentication
When selecting a locked model:
1. Tries auto-detection (3s timeout)
2. Authenticates automatically if found
3. Only prompts for API key if auto-detect fails

### Smart Model Recommendations

**Task Detection**:
- Debugging: `test`, `bug`, `debug`, `fix`
- Refactoring: `refactor`, `clean`, `improve`, `optimize`
- Code Generation: `build`, `create`, `implement`, `add`
- Code Review: `review`, `analyze`, `explain`
- Quick Questions: `quick`, `?`, `how`, `what`

**Recommendation Logic**:
- Analyzes prompt to detect task type
- Gets AI-powered recommendations (>70% confidence)
- Shows non-intrusive toast if better model available
- User maintains full control

**Implementation**: `packages/tui/internal/app/app.go:AnalyzePromptAndRecommendModel()`

---

## 🔌 Provider System

### Supported Providers

**Cloud Providers**:
- OpenAI (GPT-3.5, GPT-4, Codex)
- Anthropic (Claude Sonnet, Opus, Haiku)
- Google (Gemini models)
- xAI (Grok)
- Alibaba (Qwen)
- GitHub Copilot (multi-model access)

**Local Providers**:
- Ollama (OpenAI-compatible)
- LM Studio (OpenAI-compatible)
- Any OpenAI-compatible endpoint

### Auto-Detection System
**File**: `packages/rycode/src/auth/auto-detect.ts`

Detects credentials from:
- Environment variables
- Config files (`~/.anthropic/config.json`, etc.)
- CLI tools (`gcloud`, `anthropic`, `openai`)

**Test Scripts**:
```bash
bun run packages/rycode/test/provider-test.ts        # Detection test
bun run packages/rycode/test/provider-e2e-test.ts    # E2E API test
```

---

## 🛠️ Development Workflow

### Getting Started

```bash
# Install dependencies
bun install

# Start RyCode server
bun run packages/rycode/src/index.ts serve --port 4096

# Start TUI (in another terminal)
rycode

# Or use combined dev script
bun run dev
```

### Building

```bash
# Type check
bun run typecheck

# Build all packages
bun turbo build

# Build TUI
cd packages/tui
go build -o dist/rycode ./cmd/rycode

# Build rycode server
cd packages/rycode
bun run ./script/build.ts
```

### Testing

```bash
# Run tests
bun test

# Test provider detection
bun run packages/rycode/test/provider-test.ts

# Test TUI
cd test-tui
playwright test
```

---

## 📝 Code Style & Conventions

### TypeScript (RyCode Server)
- **Runtime**: Bun
- **Style**: Prettier (no semicolons, 120 char width)
- **Imports**: Absolute paths from `src/`
- **Types**: Zod for validation
- **Error Handling**: `NamedError` class
- **Logging**: `Log.Default` with structured logging

### Go (TUI)
- **Version**: Go 1.24.0
- **Style**: gofmt standard
- **Patterns**: Bubble Tea architecture (Model-Update-View)
- **Error Handling**: Explicit error returns
- **Logging**: slog with context

### SolidJS (Desktop)
- **Framework**: SolidJS 1.9.9
- **Style**: Prettier (no semicolons, 120 char width)
- **Imports**: `@/` alias for `src/`
- **Components**: Function declarations, splitProps
- **Styling**: TailwindCSS 4.x with CSS variables

---

## 🔧 Key Technologies

### Backend
- **Hono**: Fast web framework (4.7.10)
- **AI SDK**: Vercel AI SDK (5.0.8)
- **Zod**: Schema validation (4.1.8)
- **Bun**: JavaScript runtime & bundler
- **LSP**: Language Server Protocol support
- **MCP**: Model Context Protocol support

### Frontend
- **SolidJS**: Reactive UI framework (1.9.9)
- **TailwindCSS**: Utility-first CSS (4.x)
- **Kobalte**: Accessible UI primitives (0.13.11)
- **Vite**: Build tool

### TUI
- **Bubble Tea**: TUI framework (v2)
- **Lip Gloss**: Terminal styling (v2)
- **Bubbles**: UI components (v2)
- **Glamour**: Markdown rendering

---

## 📚 Important Files & Locations

### Configuration
- `opencode.json` - Project-level configuration
- `~/.local/share/rycode/auth.json` - User credentials
- `.env` - Environment variables (gitignored)

### Key Source Files
- `packages/tui/cmd/rycode/main.go` - TUI entry point
- `packages/rycode/src/index.ts` - Server entry point
- `packages/rycode/src/provider/provider.ts` - Provider loading logic
- `packages/tui/internal/app/app.go` - TUI core logic
- `packages/tui/internal/auth/bridge.go` - Auth bridge

### Documentation
- `docs/AUTH_PRIORITY.md` - Authentication priority
- `docs/OAUTH_AUTHENTICATION.md` - OAuth guide
- `docs/DEVELOPER_API_KEYS.md` - API key setup
- `packages/tui/INTELLIGENT_MODEL_AUTOMATION.md` - Model automation
- `docs/PROVIDER_TESTING_SUMMARY.md` - Provider testing
- `packages/desktop/AGENTS.md` - Desktop development guidelines

---

## 🎯 Recent Improvements

### Splash Screen
- Epic neural cortex animation
- Responsive borders with lipgloss
- Model selector with 5 SOTA models
- Keyboard navigation (↑↓, Tab, Enter)

### Authentication
- OAuth support for Claude Pro/Max
- GitHub Copilot integration
- Auto-detection on first run
- Background authentication for locked models

### Model Management
- Smart task-based recommendations
- Auto-detection of credentials
- Non-intrusive toast notifications
- Confidence-based suggestions (>70%)

---

## 🔒 Security

### Plugin System
- Sandboxed execution
- Code signing and verification
- Hash-based integrity checks
- Permission system

### File Security
- Security scanning for sensitive files
- Gitignore pattern matching
- Path validation
- XSS prevention

### Authentication
- Secure token storage (chmod 600)
- Auto-refresh of OAuth tokens
- No plaintext passwords
- CSRF protection

---

## 📦 Package Management

**Package Manager**: Bun 1.2.21
**Monorepo**: Turbo
**Workspaces**: Bun workspaces

**Catalog Dependencies** (shared versions):
```json
{
  "@types/bun": "1.2.21",
  "@hono/zod-validator": "0.4.2",
  "@kobalte/core": "0.13.11",
  "typescript": "5.8.2",
  "zod": "4.1.8",
  "hono": "4.7.10",
  "ai": "5.0.8",
  "solid-js": "1.9.9"
}
```

---

## 🚢 Deployment

### Server
```bash
bun run packages/rycode/src/index.ts serve --port 4096
```

### TUI Binary
```bash
cd packages/tui
go build -ldflags="-s -w" -o releases/v1.0.0/rycode-linux-arm64 ./cmd/rycode
```

### Environment
```bash
export RYCODE_SERVER="http://127.0.0.1:4096"
```

---

## 🎓 Learning Resources

### For New Contributors
1. Read `docs/AUTH_PRIORITY.md` for authentication system
2. Review `packages/desktop/AGENTS.md` for code style
3. Check `docs/PROVIDER_TESTING_SUMMARY.md` for testing

### For AI Context
1. This document provides complete project overview
2. Authentication flows documented in `docs/OAUTH_AUTHENTICATION.md`
3. Model automation explained in `packages/tui/INTELLIGENT_MODEL_AUTOMATION.md`
4. Provider testing in `docs/PROVIDER_TESTING_SUMMARY.md`

---

## 🐛 Common Tasks

### Add New Provider
1. Add to `packages/rycode/src/provider/models.ts`
2. Add detection to `packages/rycode/src/auth/auto-detect.ts`
3. Test with `bun run packages/rycode/test/provider-test.ts`

### Add New CLI Command
1. Create in `packages/rycode/src/cli/cmd/`
2. Register in `packages/rycode/src/index.ts`
3. Follow existing command patterns

### Add TUI Component
1. Create in `packages/tui/internal/components/`
2. Use Bubble Tea patterns (Init/Update/View)
3. Style with Lip Gloss

### Debug Authentication
```bash
# Check environment detection
bun run packages/rycode/test/provider-test.ts

# List configured providers
bun run packages/rycode/src/index.ts auth list

# Test actual API call
bun run packages/rycode/test/provider-e2e-test.ts
```

---

## 🎨 UI/UX Patterns

### TUI Design Principles
- **Responsive**: Adapts to terminal size
- **Accessible**: Keyboard-first navigation
- **Informative**: Clear status indicators
- **Non-blocking**: Background operations
- **Beautiful**: Lip Gloss styling

### Color Scheme
- Primary: Cyan/Blue
- Success: Green
- Warning: Yellow
- Error: Red
- Locked: Gray

### Keyboard Shortcuts
- `↑↓`: Navigate
- `Tab`: Quick switch
- `Enter`: Select
- `/model`: Open model selector
- `Ctrl+C`: Exit

---

## 📊 Performance Considerations

### Timeouts
- First run detection: 2s
- Auto-detect: 10s (background)
- Provider-specific: 3s
- Model recommendations: 3s (background)

### Optimizations
- Lazy loading of providers
- Background authentication
- Session compaction
- LSP caching

---

## 🔮 Future Enhancements

### Planned Features
1. Setup wizard for new users
2. Ollama auto-detection
3. Learning user preferences
4. Cost-aware recommendations
5. Performance-aware model switching

### Potential Improvements
- Multi-agent workflows
- Advanced plugin system
- Enhanced debugging tools
- Team collaboration features

---

## 📞 Support & Contributing

**Issues**: https://github.com/aaronmrosenthal/RyCode/issues
**Documentation**: `docs/` directory
**Tests**: `packages/rycode/test/` and `test-tui/`

### Contribution Guidelines
1. Follow existing code style
2. Add tests for new features
3. Update documentation
4. Run type checking before commit
5. Use conventional commits

---

## 📄 License

MIT License - See repository for full license text

---

**Generated**: 2025-10-12
**Context Version**: 1.0
**For**: AI-assisted development sessions
