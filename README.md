```
██████╗ ██╗   ██╗ ██████╗ ██████╗ ██████╗ ███████╗
██╔══██╗╚██╗ ██╔╝██╔════╝██╔═══██╗██╔══██╗██╔════╝
██████╔╝ ╚████╔╝ ██║     ██║   ██║██║  ██║█████╗
██╔══██╗  ╚██╔╝  ██║     ██║   ██║██║  ██║██╔══╝
██║  ██║   ██║   ╚██████╗╚██████╔╝██████╔╝███████╗
╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝
```

<p align="center">
  <strong>OpenCode, AI-Optimized</strong><br>
  What happens when you run <a href="https://toolkit-cli.com">toolkit-cli</a> on OpenCode
</p>

<p align="center">
  <a href="https://github.com/aaronmrosenthal/RyCode"><img alt="GitHub" src="https://img.shields.io/github/stars/aaronmrosenthal/RyCode?style=flat-square" /></a>
  <a href="https://github.com/aaronmrosenthal/RyCode"><img alt="License" src="https://img.shields.io/badge/license-MIT-blue?style=flat-square" /></a>
</p>

---

## 🤖 The AI Optimization Demo

```bash
$ toolkit-cli optimize opencode/

→ Claude analyzing architecture patterns...
→ Gemini reviewing performance bottlenecks...
→ Codex refactoring critical paths...
→ Qwen optimizing bundle size...

✓ Generated RyCode - Multi-agent optimized IDE
```

---

[![RyCode Terminal UI](packages/web/src/assets/lander/screenshot.png)](https://github.com/aaronmrosenthal/RyCode)

---

## About RyCode

**RyCode is OpenCode enhanced by toolkit-cli's multi-agent AI collaboration.**

This project demonstrates what happens when you run toolkit-cli's optimization commands on an existing codebase. Instead of a single AI making changes, **Claude, Gemini, Codex, and Qwen collaborate** - each bringing specialized expertise to create a superior result.

### The Before → After Story

**Before: OpenCode**
A great terminal-based AI coding agent, but optimized by a single perspective.

**After: RyCode**
The same codebase, enhanced by running toolkit-cli's multi-agent commands:

```bash
/fix      # Claude identifies and fixes architectural issues
/improve  # Gemini suggests performance optimizations
/security # Qwen hardens security vulnerabilities
/optimize # Codex refactors for efficiency
```

### What toolkit-cli Commands Were Run

RyCode showcases these toolkit-cli capabilities:

- **`/fix`** - Multi-agent bug detection and resolution
- **`/improve`** - Collaborative code quality enhancement
- **`/security`** - AI-powered security hardening
- **`/optimize`** - Performance and bundle size optimization
- **`/peer-review`** - Multi-perspective code review
- **`/test`** - Comprehensive testing strategy

### Why Multi-Agent Optimization Works

Instead of one AI making all decisions, toolkit-cli orchestrates specialists:

- **Claude** → Deep reasoning and architecture analysis
- **Gemini** → Fast iteration and performance insights
- **Codex** → Precise code generation and refactoring
- **Qwen** → Efficient optimization and security

**Result:** Better code than any single AI could produce alone.

---

### Installation

```bash
# Package managers
npm i -g rycode@latest             # or bun/pnpm/yarn
brew install aaronmrosenthal/tap/rycode      # macOS and Linux
```

> [!TIP]
> RyCode is actively developed. Check back for updates!

#### Installation Directory

The install script respects the following priority order for the installation path:

1. `$RYCODE_INSTALL_DIR` - Custom installation directory
2. `$XDG_BIN_DIR` - XDG Base Directory Specification compliant path
3. `$HOME/bin` - Standard user binary directory (if exists or can be created)
4. `$HOME/.rycode/bin` - Default fallback

### Documentation

For more info on how to configure RyCode, see the inline help and configuration options.

### Contributing

RyCode is built with toolkit-cli and showcases multi-agent AI development. Contributions are welcome!

We accept PRs for:

- Bug fixes
- Improvements to LLM performance
- Support for new AI providers
- Env-specific fixes
- Documentation improvements
- New toolkit-cli slash command integrations

To run RyCode locally you need:

- Bun
- Golang 1.24.x

And run:

```bash
$ bun install
$ bun dev
```

#### Development Notes

**Built with toolkit-cli**: This project demonstrates collaborative AI development patterns. Changes to core functionality may involve multiple AI agents working together through toolkit's slash command system.

### FAQ

#### What is RyCode?

RyCode is **OpenCode optimized by toolkit-cli's multi-agent AI collaboration**. It's a living demonstration of what happens when you run toolkit's `/fix`, `/improve`, `/security`, and `/optimize` commands on a real codebase.

#### How was RyCode created?

```bash
# Started with OpenCode
git clone https://github.com/sst/opencode

# Ran toolkit-cli optimization commands
toolkit-cli /fix opencode/
toolkit-cli /improve opencode/
toolkit-cli /security opencode/
toolkit-cli /optimize opencode/

# Result: RyCode - The AI-optimized version
```

#### What's different from OpenCode?

- **Multi-agent reviewed** - Every change analyzed by Claude, Gemini, Codex, and Qwen
- **Collaboratively optimized** - Not just one AI's perspective, but collective intelligence
- **Systematically enhanced** - Used toolkit's 26 slash commands for comprehensive improvements
- **Before/After showcase** - Demonstrates the power of multi-agent collaboration

#### Why multi-agent over single-AI?

**Single AI:** One perspective, one set of biases, one approach.

**Multi-agent (toolkit-cli):**
- Claude catches architectural issues
- Gemini spots performance bottlenecks
- Codex suggests better patterns
- Qwen finds security vulnerabilities

**Result:** Code that's been reviewed and optimized from 4 different expert perspectives.

---

**Built with** [toolkit-cli](https://toolkit-cli.com) | **Created by** [Aaron Rosenthal](https://github.com/aaronmrosenthal)
