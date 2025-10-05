```
██████╗ ██╗   ██╗ ██████╗ ██████╗ ██████╗ ███████╗
██╔══██╗╚██╗ ██╔╝██╔════╝██╔═══██╗██╔══██╗██╔════╝
██████╔╝ ╚████╔╝ ██║     ██║   ██║██║  ██║█████╗
██╔══██╗  ╚██╔╝  ██║     ██║   ██║██║  ██║██╔══╝
██║  ██║   ██║   ╚██████╗╚██████╔╝██████╔╝███████╗
╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝
```

<p align="center">
  <strong>The AI-Native Evolution of OpenCode</strong><br>
  Multi-agent development powered by <a href="https://toolkit-cli.com">toolkit-cli</a>
</p>

<p align="center">
  <a href="https://github.com/aaronmrosenthal/RyCode"><img alt="GitHub" src="https://img.shields.io/github/stars/aaronmrosenthal/RyCode?style=flat-square" /></a>
  <a href="https://github.com/aaronmrosenthal/RyCode"><img alt="License" src="https://img.shields.io/badge/license-MIT-blue?style=flat-square" /></a>
</p>

---

## 🤖 The AI-Native Difference

```bash
$ toolkit-cli fix opencode/ --ai "claude gemini qwen codex"

● Claude analyzing architecture patterns...
● Gemini reviewing UX implications...
● Qwen checking i18n compliance...
● Codex generating optimized solution...

✓ Synthesized fix applied. 4 agents collaborated.
```

---

[![RyCode Terminal UI](packages/web/src/assets/lander/screenshot.png)](https://github.com/aaronmrosenthal/RyCode)

---

## About RyCode

**RyCode is the AI-native version of OpenCode** - rebuilt from the ground up with multi-agent collaboration at its core.

While OpenCode pioneered AI-powered coding assistants, RyCode takes it further by making **every aspect of development AI-native**. Built using toolkit-cli's multi-agent architecture, RyCode doesn't just use AI - it's designed by AI collaboration.

Instead of a single AI perspective, **Claude, Gemini, Codex, and Qwen work together** - each bringing specialized expertise to architecture, performance, security, and optimization.

### The Evolution Story

**🔷 OpenCode** - The Foundation
A powerful terminal-based AI coding assistant, built with a single AI perspective.

**✨ RyCode** - The AI-Native Evolution
OpenCode's codebase transformed through toolkit-cli's multi-agent collaboration:

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

RyCode is **the AI-native version of OpenCode** - rebuilt from the ground up using toolkit-cli's multi-agent collaboration.

It's not just optimized code; it's a complete evolution that demonstrates what happens when you design software with AI-first principles, where multiple AI agents collaborate on architecture, security, performance, and code quality from day one.

#### How was RyCode created?

```bash
# Started with OpenCode
git clone https://github.com/aaronmrosenthal/rycode

# Ran toolkit-cli optimization commands
toolkit-cli /fix opencode/
toolkit-cli /improve opencode/
toolkit-cli /security opencode/
toolkit-cli /optimize opencode/

# Result: RyCode - The AI-native version
```

#### What's different from OpenCode?

- **AI-Native Architecture** - Designed from the ground up with multi-agent collaboration
- **Multi-Agent Foundation** - Every feature built with Claude, Gemini, Codex, and Qwen working together
- **Collective Intelligence** - Not just one AI's perspective, but 4 specialized agents collaborating
- **Systematically Evolved** - 26 toolkit-cli commands applied across the entire codebase
- **Living Demonstration** - Showcases what's possible when AI agents collaborate on software design

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
