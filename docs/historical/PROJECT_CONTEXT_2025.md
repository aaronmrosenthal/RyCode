# RyCode Project Context

**Generated:** 2025-10-10
**Purpose:** Comprehensive context package for AI sessions

---

## Executive Summary

**RyCode** is an AI-native evolution of OpenCode - a terminal-based AI coding assistant built with multi-agent collaboration at its core. The project demonstrates what's possible when multiple AI agents (Claude, Gemini, Codex, Qwen) collaborate on software architecture, performance, security, and code quality from day one.

**Key Statistics:**
- **Primary Language:** TypeScript (~28,673 lines) + Go (~36,457 lines)
- **Architecture:** Monorepo with 5+ packages managed by Turborepo
- **Runtime:** Bun 1.2.21
- **Version:** 0.14.1
- **License:** MIT

---

## Project Architecture

### High-Level Structure

```
RyCode/
├── packages/
│   ├── rycode/          # Core CLI and server (TypeScript/Bun)
│   ├── tui/             # Terminal UI (Go/Bubble Tea)
│   ├── web/             # Marketing website (Astro + SolidJS)
│   ├── desktop/         # Desktop app (SolidJS)
│   ├── plugin/          # Plugin SDK
│   ├── sdk/             # Go & JS SDKs
│   └── console/         # Admin console
├── .claude/             # 45+ toolkit-cli slash commands
├── infra/               # Infrastructure as code
└── scripts/             # Build and deployment scripts
```

### Monorepo Architecture

**Package Manager:** Bun with Turborepo for task orchestration

**Workspaces:**
- `packages/*` - Main packages
- `packages/console/*` - Console sub-packages
- `packages/sdk/js` - JavaScript SDK

**Build System:**
- Turborepo task pipeline with dependency management
- Tasks: `typecheck`, `build`, `test`
- Shared catalog dependencies for version consistency

---

## Core Technologies

### Backend Stack

**Runtime & Language:**
- **Bun 1.2.21** - JavaScript runtime and package manager
- **TypeScript 5.8.2** - Primary language for backend/CLI
- **Go 1.24.x** - TUI and performance-critical components

**Core Frameworks:**
- **Hono 4.7.10** - Fast web framework for API server
- **hono-openapi** - OpenAPI integration
- **Zod 4.1.8** - Schema validation and type inference
- **Yargs** - CLI argument parsing

**AI & LLM Integration:**
- **ai** (Vercel AI SDK 5.0.8) - Multi-provider LLM integration
- **@ai-sdk/anthropic** - Claude integration
- **@ai-sdk/google-vertex** - Gemini integration
- **@ai-sdk/amazon-bedrock** - Bedrock integration
- **@ai-sdk/openai-compatible** - OpenAI-compatible providers
- **@modelcontextprotocol/sdk** - MCP protocol implementation

**Development Tools:**
- **LSP Integration:** vscode-jsonrpc, vscode-languageserver-types
- **File Operations:** tree-sitter, web-tree-sitter, chokidar (file watching)
- **Utilities:** ignore (gitignore parsing), minimatch (glob patterns), fuzzysort

### Frontend Stack

**Web (Marketing/Docs):**
- **Astro 5.7.13** - Static site generator
- **SolidJS 1.9.9** - Reactive UI framework
- **Shiki 3.4.2** - Syntax highlighting
- **Marked** - Markdown processing
- **TailwindCSS** - Styling (via toolbeam-docs-theme)

**Desktop App:**
- **SolidJS 1.9.9** - Primary UI framework
- **@kobalte/core 0.13.11** - Accessible UI primitives
- **TailwindCSS 4.x** - Utility-first CSS
- **Vite** - Build tool and dev server

**TUI (Terminal UI):**
- **Go + Bubble Tea** - Terminal UI framework
- **Charm libraries** (lipgloss, bubbletea) - Terminal styling and interaction
- Custom responsive system with breakpoints and gesture support

---

## Component Responsibilities

### 1. rycode (Core Package)

**Location:** `packages/rycode/`
**Language:** TypeScript
**Size:** ~28,673 lines

**Key Modules:**

**CLI Commands** (`src/cli/cmd/`)
- `run.ts` - Main execution command
- `tui.ts` - Launch terminal UI
- `auth.ts` - Authentication management
- `agent.ts` - AI agent configuration
- `mcp.ts` - Model Context Protocol management
- `plugin.ts` - Plugin system commands
- `serve.ts` - Server management
- `debug/` - Debugging utilities

**Project Management** (`src/project/`)
- `project.ts` - Git-based project detection and management
- `instance.ts` - Project instance lifecycle
- `state.ts` - Project state management
- `bootstrap.ts` - Project initialization

**Session Management** (`src/session/`)
- `message.ts` - Message schemas (user, assistant, tool invocations)
- `message-v2.ts` - Enhanced message system
- `system.ts` - System message generation
- `compaction.ts` - Context window management
- `revert.ts` - Session revert/undo functionality

**LSP Integration** (`src/lsp/`)
- `client.ts` - Language Server Protocol client
- `server.ts` - LSP server management
- `language.ts` - Language detection and mapping
- Supports TypeScript, Python, Rust, Go, and more

**Plugin System** (`src/plugin/`)
- `index.ts` - Plugin loader with security controls
- `security.ts` - Plugin sandboxing and capability management
- `sandbox.ts` - Isolated plugin execution
- `registry.ts` - Plugin discovery and installation
- `signature.ts` - Plugin signing and verification

**Tool System** (`src/tool/`)
- `glob.ts` - File pattern matching
- `grep.ts` - Code search
- `read.ts` - File reading
- `multiedit.ts` - Batch file editing
- `lsp-diagnostics.ts` - LSP diagnostics integration
- `lsp-hover.ts` - LSP hover information
- `task.ts` - Task execution
- `webfetch.ts` - Web content fetching

**MCP Integration** (`src/mcp/`)
- Support for local and remote MCP servers
- Transport layer: StreamableHTTP, SSE, Stdio
- Tool aggregation from multiple MCP sources

**Server/API** (`src/server/`)
- `project.ts` - Project API routes (Hono)
- `tui.ts` - TUI server integration
- Middleware: authentication, rate limiting, security monitoring

**Authentication** (`src/auth/`)
- `github-copilot.ts` - GitHub Copilot integration
- Multi-provider auth support

**File Operations** (`src/file/`)
- `ignore.ts` - .gitignore parsing
- `ripgrep.ts` - Fast code search
- `fzf.ts` - Fuzzy file finding
- `watcher.ts` - File system watching
- `security.ts` - File security checks

**Configuration** (`src/config/`)
- `config.ts` - Configuration management
- `markdown.ts` - Markdown processing
- Support for project-specific and global configs

### 2. tui (Terminal UI)

**Location:** `packages/tui/`
**Language:** Go
**Size:** ~36,457 lines

**Key Components:**

**Main Application** (`cmd/rycode/main.go`)
- Entry point for TUI
- Connects to rycode server via SDK
- Handles piped input and command-line arguments
- Event streaming from server

**Core App** (`internal/app/`)
- `app.go` - Main application state and logic
- `state.go` - State management
- `prompt.go` - Prompt handling

**UI Components** (`internal/components/`)
- `chat/` - Message display and rendering
- `textarea/` - Input handling with memoization
- `diff/` - Diff viewer with syntax highlighting
- `dialog/` - Modal dialogs (agents, models, help, search)
- `list/` - Reusable list component
- `status/` - Status bar
- `timeline/` - Session timeline view
- `toast/` - Notifications
- `qr/` - QR code display
- `ghost/` - Ghost text suggestions
- `reactions/` - Message reactions
- `smarthistory/` - Command history

**Responsive System** (`internal/responsive/`)
- `breakpoints.go` - Terminal size breakpoints
- `gestures.go` - Touch/mouse gesture detection
- `contrast.go` - Adaptive color contrast
- `accessibility.go` - Accessibility features
- `platform.go` - Platform-specific adaptations

**Styling** (`internal/styles/`)
- `styles.go` - Style definitions
- `background.go` - Background rendering
- `utilities.go` - Style utilities
- Theme integration with color system

**API Client** (`internal/api/`)
- `api.go` - Server communication
- HTTP client for rycode server

**Input Handling** (`input/`)
- Custom input driver supporting Kitty protocol
- Mouse and keyboard event parsing
- Clipboard integration
- Focus management

### 3. web (Marketing Site)

**Location:** `packages/web/`
**Framework:** Astro + SolidJS
**Purpose:** Marketing website and documentation

**Key Features:**
- Static site generation with Astro
- Interactive components with SolidJS
- Documentation with Starlight integration
- Syntax highlighting with Shiki
- Custom toolbeam-docs-theme integration

### 4. desktop (Desktop App)

**Location:** `packages/desktop/`
**Framework:** SolidJS
**Purpose:** Native-like desktop experience

**Architecture:**
- SolidJS with TypeScript
- TailwindCSS with CSS variables theme
- @kobalte/core for accessible primitives
- Vite for build and dev server
- File structure: `/ui/` primitives, `/components/` higher-level, `/pages/`, `/providers/`

**Code Conventions:**
- Function declarations for components
- splitProps for component props
- No semicolons (Prettier config)
- 120 char line width
- PascalCase components, camelCase variables, snake_case files

### 5. plugin (Plugin SDK)

**Location:** `packages/plugin/`
**Purpose:** Plugin development SDK

**Exports:**
- Plugin types and interfaces
- Hook definitions
- Tool integration

### 6. sdk (Client SDKs)

**Location:** `packages/sdk/`
**Languages:** Go + JavaScript

**Go SDK:**
- Type-safe client for rycode server
- Auto-generated from OpenAPI spec
- Used by TUI

**JavaScript SDK:**
- Type-safe TypeScript client
- Used by plugins and web integrations

---

## Development Patterns & Conventions

### Code Organization

**Namespace Pattern:**
```typescript
// Heavily used throughout rycode package
export namespace Project {
  export const Info = z.object({...})
  export type Info = z.infer<typeof Info>

  export async function fromDirectory(directory: string) {...}
}
```

**State Management with Instance:**
```typescript
// Pattern for project-scoped state
const state = Instance.state(
  async () => {
    // Initialize state
    return { ... }
  },
  async (state) => {
    // Cleanup on project change
  }
)
```

**Schema-First Development:**
```typescript
// Zod schemas define both runtime validation and TypeScript types
export const Message = z.object({
  id: z.string(),
  role: z.enum(["user", "assistant"]),
  parts: z.array(MessagePart)
}).meta({ ref: "Message" })

export type Message = z.infer<typeof Message>
```

### Error Handling

**Named Errors:**
```typescript
// Pattern from util/error.ts
export const InitializeError = NamedError.create(
  "LSPInitializeError",
  z.object({
    serverID: z.string()
  })
)

// Usage
throw new InitializeError({ serverID: input.serverID }, { cause: err })
```

### Event Bus Pattern

```typescript
// Event definition
export const Event = {
  Diagnostics: Bus.event(
    "lsp.client.diagnostics",
    z.object({
      serverID: z.string(),
      path: z.string()
    })
  )
}

// Publishing
Bus.publish(Event.Diagnostics, { path, serverID })

// Subscribing
Bus.subscribe(Event.Diagnostics, (event) => {...})
```

### Logging

```typescript
const log = Log.create({ service: "component-name" })
log.info("action", { data: "..." })
log.error("failure", { error: err })
```

### Configuration

**Location:** User config at `~/.config/opencode/config.json` (transitioning to rycode)

**Schema-driven:**
```typescript
export const Config = z.object({
  provider: z.record(z.string(), ProviderConfig),
  plugin: z.array(z.string()).optional(),
  plugin_security: PluginSecurityPolicy.optional(),
  mcp: z.record(z.string(), MCPConfig).optional()
})
```

### Security Patterns

**Plugin Security:**
- Capability-based permissions
- Sandboxed execution
- Trust allowlists
- Integrity verification (hashing)
- User approval for untrusted plugins

**File Security:**
- Gitignore parsing and respect
- Security file detection (.env, credentials)
- Path validation and sanitization

---

## Key Features & Functionality

### Multi-Agent AI Collaboration

**Providers Supported:**
- Anthropic (Claude)
- Google Vertex (Gemini)
- Amazon Bedrock
- OpenAI-compatible providers

**Agent System:**
- Multiple agents per session
- Specialized roles (architecture, security, performance)
- Collaborative decision-making

### Tool System

**Built-in Tools:**
- File operations (glob, grep, read, write, edit)
- LSP integration (diagnostics, hover, definitions)
- Task execution
- Web fetching
- Multiedit (batch file changes)

**MCP Integration:**
- Extensible tool system via Model Context Protocol
- Support for local and remote MCP servers
- Dynamic tool discovery and registration

### Plugin System

**Features:**
- NPM-based plugin distribution
- Security policies (strict/warn/permissive modes)
- Capability-based permissions
- Plugin signing and verification
- Sandboxed execution
- Default plugins: copilot-auth, anthropic-auth

**Hooks:**
- `config` - Configuration access
- `auth` - Authentication provider
- `event` - Event subscription
- `tool` - Custom tool registration

### Session Management

**Capabilities:**
- Persistent sessions across restarts
- Context window management with compaction
- Message history with tool invocations
- Session revert/undo
- Timeline view of session history

### LSP Integration

**Supported Languages:**
- TypeScript/JavaScript (via tsserver)
- Python (via pylsp/pyright)
- Rust (via rust-analyzer)
- Go (via gopls)
- And more...

**Features:**
- Real-time diagnostics
- Hover information
- Type-aware code understanding
- Automatic language server detection

### Terminal UI Features

**Responsive Design:**
- Breakpoint system for terminal sizes
- Adaptive layouts
- Touch/gesture support
- Platform-specific optimizations

**Accessibility:**
- High contrast support
- Screen reader considerations
- Keyboard navigation
- Configurable color schemes

**Components:**
- Syntax-highlighted code display
- Interactive diff viewer
- Fuzzy search dialogs
- Command history
- Multi-pane layouts

---

## Build & Development

### Prerequisites

- Bun 1.2.21+
- Go 1.24.x
- Git

### Development Commands

```bash
# Install dependencies
bun install

# Run development server
bun dev

# Type checking
bun run typecheck

# Build all packages
bun run build
```

### Package-Specific Commands

**rycode:**
```bash
cd packages/rycode
bun run dev                    # Run CLI
bun run build                  # Build for production
bun run test                   # Run tests
```

**tui:**
```bash
cd packages/tui
go build ./cmd/rycode          # Build TUI
go test ./...                  # Run tests
```

**web:**
```bash
cd packages/web
bun run dev                    # Dev server
bun run build                  # Build static site
```

**desktop:**
```bash
cd packages/desktop
bun run dev                    # Vite dev server on port 3000
bun run build                  # Production build
bun run typecheck              # Type validation only
```

### Testing Strategy

**Validation:**
- Type checking with TypeScript
- No automated test suite (manual testing focused)
- Integration testing via development usage

---

## Notable Implementation Details

### Project Detection

Projects are identified by git repository root commit hash:
```typescript
// packages/rycode/src/project/project.ts
const [id] = await $`git rev-list --max-parents=0 --all`
```

This provides stable project IDs across file system moves.

### Context Window Management

Implements intelligent message compaction to fit within LLM context limits:
- Tool results summarization
- Old message pruning
- Preservation of critical context

### Security Model

**Plugin Security Modes:**
1. **Strict:** Only trusted plugins allowed
2. **Warn:** Untrusted plugins allowed with warnings
3. **Permissive:** All plugins allowed

**Capabilities:**
- `file:read` - Read file system
- `file:write` - Write to file system
- `network` - Network access
- `process` - Process execution
- `env` - Environment variable access

### Message Format

Messages use a parts-based structure supporting:
- Text content
- Reasoning (extended thinking)
- Tool invocations (call, partial-call, result states)
- Source URLs (for citations)
- File attachments
- Step boundaries

---

## Configuration & Customization

### Config File Location

- User: `~/.config/opencode/config.json` (transitioning to `~/.config/rycode/`)
- Project: `.opencode/` or `.claude/` directories

### Slash Commands

45+ custom slash commands in `.claude/commands/`:

**Development:**
- `/fix` - Bug detection and resolution
- `/improve` - Code quality enhancement
- `/make` - Implementation guidance
- `/implement` - Interactive task implementation

**Quality:**
- `/test` - Testing strategy
- `/debug` - Debugging assistance
- `/peer-review` - Multi-agent code review
- `/reflect` - Constructive critique

**Planning:**
- `/plan` - Implementation planning
- `/specify` - Feature specification
- `/tasks` - Generate task lists
- `/next` - Next step recommendations

**Security:**
- `/security` - Security analysis
- `/verify` - Verification and validation

**UX/Design:**
- `/ux` - UX design and accessibility
- `/mock-ups` - Mockup management
- `/responsive` - Responsive design testing

**Deployment:**
- `/ship` - Pre-deployment checklist
- `/deploy` - Deployment automation
- `/optimize` - Performance optimization

**Utilities:**
- `/re-context` - Generate context package (this command!)
- `/errors` - Screenshot OCR + log parsing
- `/diff` - Git diff analysis
- `/keys` - .env file management
- `/undo` - Checkpoint-based rollback

---

## Recent Development Focus

Based on git history (last 5 commits):

1. **Type Safety & Error Handling** (13fb274c)
   - Enhanced error handling throughout codebase
   - Type safety improvements

2. **Environment & Path Updates** (cfc6bc81)
   - Migration from "opencode" to "rycode" branding
   - Updated environment variables and paths

3. **Process.env Access Fix** (4edea357)
   - Fixed environment variable access patterns

4. **Build Error Resolution** (590d486e)
   - Resolved debugger and toolkit client build issues
   - Fixed Go component compatibility

5. **Branding Update** (a4b046bc)
   - Systematic update of references from opencode to rycode

---

## Architecture Insights

### Monorepo Benefits

**Shared Dependencies:**
- Catalog system ensures version consistency
- TypeScript/Bun/Zod versions unified across packages

**Cross-Package Communication:**
- rycode server provides HTTP API
- TUI consumes API via Go SDK
- Desktop/web can integrate via JS SDK

**Build Pipeline:**
- Turborepo manages inter-package dependencies
- Parallel builds where possible
- Cached builds for efficiency

### Multi-Runtime Strategy

**Why TypeScript + Go?**
- **TypeScript (Bun):** Fast development, AI SDK ecosystem, plugin system
- **Go:** Performance-critical TUI, binary distribution, platform compatibility

**Communication:**
- HTTP API boundary between runtimes
- Auto-generated SDKs from OpenAPI specs
- Event streaming for real-time updates

### Security-First Design

**Defense in Depth:**
1. Plugin sandboxing and capability restrictions
2. File security checks (gitignore, secret detection)
3. Rate limiting on API routes
4. Security monitoring middleware
5. Plugin signature verification
6. User approval for untrusted operations

---

## Future Considerations

Based on project documentation:

**Planned Features:**
- Enhanced mobile-first TUI (see MOBILE_FIRST_UX_ARCHITECTURE.md)
- Matrix-style TUI visualization (see MATRIX_TUI_SPECIFICATION.md)
- Improved concurrency (see CONCURRENCY_IMPROVEMENTS.md)
- Database migrations (see DATABASE_MIGRATIONS.md)

**Evolution from OpenCode:**
- Systematic rename from "opencode" to "rycode"
- Enhanced multi-agent collaboration
- Improved security model
- Modern UI/UX across all interfaces

---

## Quick Reference

### File Locations

| Component | Path | Purpose |
|-----------|------|---------|
| Main CLI | `packages/rycode/src/index.ts` | Entry point |
| TUI Main | `packages/tui/cmd/rycode/main.go` | Terminal UI entry |
| Config | `~/.config/opencode/config.json` | User configuration |
| Logs | Platform-specific | Debug logs |
| Plugins | `~/.bun/install/cache/` | Installed plugins |
| Sessions | Storage (platform-specific) | Session persistence |

### Key Concepts

**Project:** A git repository identified by root commit hash
**Session:** A conversation with AI, potentially spanning multiple messages
**Agent:** An AI model with specific role/configuration
**Tool:** An action the AI can perform (file ops, LSP, etc.)
**Plugin:** User-installed extension providing hooks/tools
**MCP:** Model Context Protocol - standard for tool integration
**LSP:** Language Server Protocol - code intelligence

### Important Entry Points

Reading the codebase? Start here:

1. **CLI Setup:** `packages/rycode/src/index.ts`
2. **Project Detection:** `packages/rycode/src/project/project.ts`
3. **Session Logic:** `packages/rycode/src/session/message.ts`
4. **Tool System:** `packages/rycode/src/tool/tool.ts`
5. **Plugin System:** `packages/rycode/src/plugin/index.ts`
6. **TUI App:** `packages/tui/internal/app/app.go`
7. **Server Routes:** `packages/rycode/src/server/project.ts`

---

## Documentation Files

The project includes extensive documentation:

- `AGENTS.md` - Agent system overview
- `SECURITY*.md` - Multiple security documentation files
- `FEATURE_SPECIFICATION.md` - Feature specs
- `TESTING_STRATEGY.md` - Testing approach
- `PROJECT_CONTEXT.md` - Original context doc (now superseded)
- `PRODUCTION_DEPLOYMENT_CHECKLIST.md` - Deployment guide
- Various specification documents for planned features

---

## Contributing

**Accepted PRs:**
- Bug fixes
- LLM performance improvements
- New AI provider support
- Environment-specific fixes
- Documentation improvements
- New toolkit-cli slash command integrations

**Development Notes:**
- Built with toolkit-cli demonstrating collaborative AI patterns
- Changes to core functionality may involve multiple AI agents
- Use slash commands for complex tasks

---

## Links & Resources

- **Repository:** https://github.com/aaronmrosenthal/RyCode
- **toolkit-cli:** https://toolkit-cli.com
- **Author:** Aaron Rosenthal

---

## Context Usage Tips

When starting a new AI session with this context:

1. **For feature work:** Reference the specific package (`rycode`, `tui`, `web`, `desktop`)
2. **For debugging:** Check the component responsibilities section first
3. **For API work:** Look at server routes in `packages/rycode/src/server/`
4. **For UI work:** TUI (Go) vs Desktop (SolidJS) - check the appropriate section
5. **For security:** Review the security patterns and plugin system sections
6. **For architecture:** Start with the high-level structure and monorepo benefits

**This document is comprehensive but not exhaustive.** Use it as a map, then dive into specific files as needed.

---

*Generated by RyCode `/re-context` command - 2025-10-10*
