# RyCode Documentation

## 📚 Documentation Structure

RyCode documentation is organized by feature and development phase.

### 🔐 Provider Authentication System

**Location:** `docs/provider-auth/`

The Provider Authentication System allows users to authenticate directly with AI providers (Anthropic, OpenAI, Google, Grok, Qwen) from within RyCode.

#### Quick Links
- [Executive Summary](docs/provider-auth/phase-1/EXECUTIVE_SUMMARY.md) - Business overview and impact
- [Quick Reference](docs/provider-auth/phase-1/QUICK_REFERENCE.md) - Code examples and API usage
- [Integration Guide](packages/rycode/src/auth/INTEGRATION_GUIDE.md) - Developer integration guide
- [TUI Integration Plan](docs/provider-auth/phase-2/TUI_INTEGRATION_PLAN.md) - UI integration design

#### Status
- ✅ **Phase 1 Complete** - Core infrastructure (5,045 lines, 16 files)
- 🔄 **Phase 2 In Progress** - TUI integration
- ⏳ **Phase 3 Pending** - Migration wizard
- ⏳ **Phase 4 Pending** - Testing & launch

#### Features
- 5 AI provider integrations
- Enterprise-grade security (rate limiting, circuit breakers)
- Auto-detection of credentials (12+ sources)
- Real-time cost tracking with projections
- Smart model recommendations
- Inline authentication in model selector
- Tab key model switching

---

### 🎨 TUI (Terminal User Interface)

**Location:** `packages/tui/`

The TUI provides a rich terminal-based interface for interacting with RyCode.

#### Components
- Status bar - Shows current model, cost, and git info
- Model selector - Browse and select models with inline auth
- Chat interface - Conversation with AI models
- Session management - Manage and switch between sessions

---

### 🧰 Core Packages

#### rycode (TypeScript)
**Location:** `packages/rycode/`

Core RyCode functionality including:
- Authentication system (`src/auth/`)
- API client
- Configuration management
- Storage utilities

#### TUI (Go)
**Location:** `packages/tui/`

Terminal interface built with Bubble Tea:
- Components (status bar, dialogs, editor)
- App state management
- Command system
- Keybindings

#### Web
**Location:** `packages/web/`

Web-based interface and landing page.

---

## 📖 Getting Started

### For Users
1. [Installation](README.md#installation) - Install RyCode
2. [Quick Reference](docs/provider-auth/phase-1/QUICK_REFERENCE.md) - Common tasks

### For Developers
1. [Integration Guide](packages/rycode/src/auth/INTEGRATION_GUIDE.md) - Integrate auth system
2. [TUI Integration Plan](docs/provider-auth/phase-2/TUI_INTEGRATION_PLAN.md) - Extend TUI
3. [Contributing](README.md#contributing) - Contribution guidelines

---

## 🔍 Find Documentation

### By Topic
- **Authentication** → `docs/provider-auth/`
- **TUI Components** → `packages/tui/internal/components/`
- **API Reference** → `packages/rycode/src/auth/README.md`
- **Architecture** → `docs/provider-auth/phase-1/ARCHITECTURE_DIAGRAM.md`

### By Role
- **Product Managers** → Start with [Executive Summary](docs/provider-auth/phase-1/EXECUTIVE_SUMMARY.md)
- **Engineers** → Start with [Quick Reference](docs/provider-auth/phase-1/QUICK_REFERENCE.md)
- **QA/Test** → Start with [Launch Checklist](docs/provider-auth/phase-1/LAUNCH_CHECKLIST.md)
- **Security** → See [Security section](docs/provider-auth/phase-1/EXECUTIVE_SUMMARY.md#-security-posture)
- **Support** → Start with [Quick Reference](docs/provider-auth/phase-1/QUICK_REFERENCE.md)

---

## 🚀 Current Development

### Active Work
- Phase 2: TUI Integration (status bar, model selector, Tab cycling)

### Recently Completed
- Phase 1: Provider Authentication Infrastructure (Oct 2024)
  - 5 provider integrations
  - Security layer (rate limiting, circuit breakers, validation)
  - Smart features (auto-detect, cost tracking, recommendations)
  - Complete documentation (16 files)

### Upcoming
- Phase 3: Migration wizard for legacy agent system
- Phase 4: Testing and gradual rollout

---

## 📝 Documentation Standards

### File Organization
- Feature docs → `docs/{feature-name}/`
- Phase-based organization → `docs/{feature}/phase-{n}/`
- Code docs → Alongside source code

### File Types
- `README.md` - Overview and quick start
- `*_GUIDE.md` - How-to guides
- `*_PLAN.md` - Planning documents
- `*_CHECKLIST.md` - Task lists
- `*_SUMMARY.md` - Executive summaries
- `*_COMPLETE.md` - Completion reports

---

**Last Updated:** October 2024
**Maintainer:** Development Team
