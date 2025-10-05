```
██████╗ ██╗   ██╗ ██████╗ ██████╗ ██████╗ ███████╗
██╔══██╗╚██╗ ██╔╝██╔════╝██╔═══██╗██╔══██╗██╔════╝
██████╔╝ ╚████╔╝ ██║     ██║   ██║██║  ██║█████╗
██╔══██╗  ╚██╔╝  ██║     ██║   ██║██║  ██║██╔══╝
██║  ██║   ██║   ╚██████╗╚██████╔╝██████╔╝███████╗
╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝
```

<p align="center">
  <strong>The AI-Generated IDE</strong><br>
  Built with <a href="https://toolkit-cli.com">toolkit-cli</a> - Where LLMs Collaborate, Not Compete
</p>

<p align="center">
  <a href="https://github.com/aaronmrosenthal/RyCode"><img alt="GitHub" src="https://img.shields.io/github/stars/aaronmrosenthal/RyCode?style=flat-square" /></a>
  <a href="https://github.com/aaronmrosenthal/RyCode"><img alt="License" src="https://img.shields.io/badge/license-MIT-blue?style=flat-square" /></a>
</p>

---

[![RyCode Terminal UI](packages/web/src/assets/lander/screenshot.png)](https://github.com/aaronmrosenthal/RyCode)

---

## About RyCode

RyCode is a next-generation AI coding agent that leverages **multi-agent collaboration** through [toolkit-cli](https://toolkit-cli.com). Unlike traditional single-model approaches, RyCode harnesses the collective intelligence of multiple LLMs (Claude, Gemini, Codex, Qwen) working together to deliver superior development experiences.

### Powered by toolkit-cli

This project showcases the power of toolkit-cli's innovative approach:

- **26 AI-Powered Slash Commands** - Streamlined workflows for every development task
- **Multi-Agent Architecture** - Different AI models collaborate on complex problems
- **Spec-Context Preservation** - Maintains project intent across all changes
- **Provider-Agnostic** - Not locked into any single AI vendor
- **Terminal-First** - Built by developers, for developers

### Why Multi-Agent?

RyCode demonstrates that multiple specialized AI models working together outperform any single model alone. Each agent brings unique strengths:

- **Claude** - Deep reasoning and code analysis
- **Gemini** - Fast iteration and broad knowledge
- **Codex** - Precise code generation
- **Qwen** - Efficient problem-solving

Together, they create an IDE that thinks, adapts, and evolves with your codebase.

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

#### What makes RyCode different?

RyCode is built with **multi-agent AI collaboration** at its core:

- **Powered by toolkit-cli** - Leverages 26 AI slash commands and multi-model orchestration
- **Multiple LLMs working together** - Claude, Gemini, Codex, and Qwen collaborate on your code
- **Provider-agnostic** - Not locked to any single AI vendor
- **Terminal-first** - Built for developers who live in the terminal
- **Spec-context preservation** - Maintains project intent across all AI-generated changes
- **Client/server architecture** - Run locally, control remotely

#### How does multi-agent collaboration work?

Different AI models excel at different tasks. RyCode uses toolkit-cli to orchestrate specialized agents:

- Complex reasoning? Claude leads
- Fast iteration? Gemini steps in
- Code generation? Codex delivers
- Efficient solutions? Qwen optimizes

The result: Better code, faster development, maintained context.

#### Is this better than single-model tools?

Yes! Multi-agent systems consistently outperform single models because:
- Each model contributes its unique strengths
- Models verify and improve each other's outputs
- Specialized agents handle specialized tasks
- The collective intelligence exceeds individual capabilities

---

**Built with** [toolkit-cli](https://toolkit-cli.com) | **Created by** [Aaron Rosenthal](https://github.com/aaronmrosenthal)
