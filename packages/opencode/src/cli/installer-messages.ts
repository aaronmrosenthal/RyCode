/**
 * Polished Installer Messages
 * Professional, engaging messaging for the RyCode installation flow
 */

import { UI } from "./ui"
import { EnhancedTUI } from "./tui-enhanced"

export namespace InstallerMessages {
  /**
   * Welcome screen - First impression
   */
  export function welcome() {
    UI.empty()

    // Matrix rain effect header
    UI.println(UI.Style.MATRIX_GREEN_DARK + "█".repeat(80) + UI.Style.RESET)
    UI.println(UI.Style.MATRIX_GREEN + "█".repeat(80) + UI.Style.RESET)
    UI.println(UI.Style.MATRIX_GREEN_BRIGHT + "█".repeat(80) + UI.Style.RESET)
    UI.println()

    // Logo with tagline
    UI.println(UI.logo())
    UI.println()
    UI.println(
      "  " +
        UI.gradient("AI-POWERED DEVELOPMENT ASSISTANT", [
          UI.Style.MATRIX_GREEN_DARK,
          UI.Style.MATRIX_GREEN,
          UI.Style.MATRIX_GREEN_BRIGHT,
        ])
    )
    UI.println()

    // Welcome message
    UI.println(
      UI.box(
        UI.glow("✨ Welcome to RyCode ✨", UI.Style.MATRIX_GREEN_BRIGHT) +
          "\n\n" +
          UI.Style.TEXT_DIM +
          "Let's get you set up in just a few steps..." +
          UI.Style.RESET,
        { color: UI.Style.MATRIX_GREEN, padding: 1 }
      )
    )
    UI.println()
  }

  /**
   * Authentication required screen
   */
  export function authRequired() {
    UI.empty()
    UI.println()

    UI.println(
      UI.box(
        UI.Style.MATRIX_GREEN + "🔐 Authentication Required" + UI.Style.RESET + "\n\n" +
        UI.Style.TEXT_DIM + "To use RyCode, you'll need to authenticate with an AI provider.\n" +
        "This enables access to powerful AI models for code assistance." + UI.Style.RESET,
        { color: UI.Style.MATRIX_GREEN, padding: 1 }
      )
    )
    UI.println()

    // Show steps
    const steps = [
      "Select your AI provider",
      "Enter your API credentials",
      "Start coding with AI assistance",
    ]

    UI.println(UI.Style.BOLD + "Setup Steps:" + UI.Style.RESET)
    UI.println()
    UI.println(EnhancedTUI.wizardSteps(steps, 0))
    UI.println()

    UI.println(
      UI.Style.TEXT_INFO +
        "📝 Don't have an API key? We'll show you where to get one." +
        UI.Style.RESET
    )
    UI.println()
  }

  /**
   * Provider selection intro
   */
  export function providerIntro() {
    UI.println()
    UI.println(UI.Style.MATRIX_GREEN + "═".repeat(70) + UI.Style.RESET)
    UI.println()
    UI.println(UI.Style.BOLD + UI.Style.MATRIX_GREEN + "🤖 AI Provider Setup" + UI.Style.RESET)
    UI.println()
    UI.println(
      UI.Style.TEXT_DIM +
        "Choose your AI provider. We recommend starting with Anthropic or OpenAI." +
        UI.Style.RESET
    )
    UI.println()
  }

  /**
   * Success message after auth
   */
  export function authSuccess(providerName: string) {
    UI.println()
    UI.println(
      UI.box(
        UI.glow("✓ Authentication Successful", UI.Style.MATRIX_GREEN_BRIGHT) +
          "\n\n" +
          UI.Style.TEXT_DIM +
          `Connected to ${providerName}` +
          UI.Style.RESET,
        { color: UI.Style.MATRIX_GREEN, padding: 1 }
      )
    )
    UI.println()

    // Completion status
    const steps = [
      "Select your AI provider",
      "Enter your API credentials",
      "Start coding with AI assistance",
    ]

    UI.println(EnhancedTUI.wizardSteps(steps, 2))
    UI.println()
  }

  /**
   * Starting RyCode message
   */
  export function starting() {
    UI.empty()
    UI.println()

    // Animated bars
    UI.println(UI.Style.MATRIX_GREEN_DARK + "█".repeat(80) + UI.Style.RESET)
    UI.println(UI.Style.MATRIX_GREEN + "█".repeat(80) + UI.Style.RESET)
    UI.println(UI.Style.MATRIX_GREEN_BRIGHT + "█".repeat(80) + UI.Style.RESET)
    UI.println()

    UI.println(UI.logo())
    UI.println()

    const tagline = ">>> INITIALIZING AI DEVELOPMENT ENVIRONMENT <<<"
    UI.println(
      "  " +
        UI.gradient(tagline, [
          UI.Style.MATRIX_GREEN_DARK,
          UI.Style.MATRIX_GREEN,
          UI.Style.MATRIX_GREEN_BRIGHT,
          UI.Style.MATRIX_GREEN,
          UI.Style.MATRIX_GREEN_DARK,
        ])
    )
    UI.println()
  }

  /**
   * Ready to code message
   */
  export function ready() {
    UI.println()
    UI.println(
      UI.box(
        UI.glow("🚀 READY TO CODE 🚀", UI.Style.MATRIX_GREEN_BRIGHT) +
          "\n\n" +
          UI.Style.TEXT_DIM +
          "RyCode is now running. Start chatting to build amazing things!" +
          UI.Style.RESET,
        { color: UI.Style.MATRIX_GREEN, padding: 1 }
      )
    )
    UI.println()

    UI.println(
      UI.Style.MATRIX_GREEN_BRIGHT + "█".repeat(80) + UI.Style.RESET
    )
    UI.println(
      UI.Style.MATRIX_GREEN + "█".repeat(80) + UI.Style.RESET
    )
    UI.println(
      UI.Style.MATRIX_GREEN_DARK + "█".repeat(80) + UI.Style.RESET
    )
    UI.println()
  }

  /**
   * API key helper messages for different providers
   */
  export const apiKeyHelp = {
    anthropic: {
      message: "Get your Anthropic API key at:",
      url: "https://console.anthropic.com/settings/keys",
      hint: "You'll need to create an account if you don't have one",
    },
    openai: {
      message: "Get your OpenAI API key at:",
      url: "https://platform.openai.com/api-keys",
      hint: "Click '+ Create new secret key' to generate one",
    },
    google: {
      message: "Get your Google AI API key at:",
      url: "https://aistudio.google.com/app/apikey",
      hint: "Free tier available with generous limits",
    },
    openrouter: {
      message: "Get your OpenRouter API key at:",
      url: "https://openrouter.ai/keys",
      hint: "Access multiple AI models through one API",
    },
    opencode: {
      message: "Create an OpenCode API key at:",
      url: "https://opencode.ai/auth",
      hint: "Official RyCode hosting service",
    },
  }

  /**
   * Display API key help for a provider
   */
  export function showApiKeyHelp(providerId: string) {
    const help = apiKeyHelp[providerId as keyof typeof apiKeyHelp]
    if (!help) return

    UI.println()
    UI.println(UI.Style.MATRIX_GREEN + "═".repeat(70) + UI.Style.RESET)
    UI.println()
    UI.println(
      UI.Style.MATRIX_GREEN_BRIGHT + "🔑 " + help.message + UI.Style.RESET
    )
    UI.println()
    UI.println("  " + UI.link(help.url, help.url))
    UI.println()
    UI.println(UI.Style.TEXT_DIM + "  💡 " + help.hint + UI.Style.RESET)
    UI.println()
    UI.println(UI.Style.MATRIX_GREEN + "═".repeat(70) + UI.Style.RESET)
    UI.println()
  }

  /**
   * Error message styling
   */
  export function error(title: string, message: string, suggestion?: string) {
    UI.println()
    UI.println(
      UI.box(
        UI.Style.TEXT_DANGER + "✖ " + title + UI.Style.RESET + "\n\n" +
        UI.Style.TEXT_DIM + message + UI.Style.RESET +
        (suggestion ? "\n\n" + UI.Style.TEXT_INFO + "💡 " + suggestion + UI.Style.RESET : ""),
        { color: UI.Style.TEXT_DANGER, padding: 1 }
      )
    )
    UI.println()
  }

  /**
   * Warning message styling
   */
  export function warning(title: string, message: string) {
    UI.println()
    UI.println(
      UI.box(
        UI.Style.TEXT_WARNING + "⚠ " + title + UI.Style.RESET + "\n\n" +
        UI.Style.TEXT_DIM + message + UI.Style.RESET,
        { color: UI.Style.TEXT_WARNING, padding: 1 }
      )
    )
    UI.println()
  }

  /**
   * Info message styling
   */
  export function info(title: string, message: string) {
    UI.println()
    UI.println(
      UI.box(
        UI.Style.CLAUDE_BLUE + "ℹ " + title + UI.Style.RESET + "\n\n" +
        UI.Style.TEXT_DIM + message + UI.Style.RESET,
        { color: UI.Style.CLAUDE_BLUE, padding: 1 }
      )
    )
    UI.println()
  }

  /**
   * Progress indicator
   */
  export function progress(step: number, total: number, message: string) {
    UI.println()
    UI.println(
      EnhancedTUI.progressBar(step, total, {
        label: `Step ${step}/${total}`,
        color: UI.Style.MATRIX_GREEN,
        width: 50,
      })
    )
    UI.println()
    UI.println(
      UI.Style.TEXT_DIM + "  " + message + UI.Style.RESET
    )
    UI.println()
  }

  /**
   * Feature showcase
   */
  export function features() {
    UI.println()
    UI.println(UI.Style.BOLD + UI.Style.MATRIX_GREEN + "✨ What You Get:" + UI.Style.RESET)
    UI.println()

    const features = [
      { icon: "🤖", title: "AI Pair Programming", desc: "Chat with AI to build features faster" },
      { icon: "🔍", title: "Intelligent Code Search", desc: "Find anything in your codebase instantly" },
      { icon: "✨", title: "Smart Completions", desc: "Context-aware code suggestions" },
      { icon: "🐛", title: "Bug Detection", desc: "Catch issues before they happen" },
      { icon: "📚", title: "Code Explanations", desc: "Understand complex code quickly" },
      { icon: "🚀", title: "Productivity Boost", desc: "Ship features 10x faster" },
    ]

    features.forEach((feature) => {
      UI.println(
        `  ${feature.icon}  ${UI.Style.BOLD}${feature.title}${UI.Style.RESET}`
      )
      UI.println(
        `     ${UI.Style.TEXT_DIM}${feature.desc}${UI.Style.RESET}`
      )
      UI.println()
    })
  }

  /**
   * Quick tips
   */
  export function quickTips() {
    UI.println()
    UI.println(UI.Style.BOLD + UI.Style.MATRIX_GREEN + "💡 Quick Tips:" + UI.Style.RESET)
    UI.println()

    const tips = [
      "Use natural language - just describe what you want to build",
      "RyCode can read your entire codebase for context",
      "Ask for explanations, refactoring, or new features",
      "Use /commands for quick actions and workflows",
    ]

    tips.forEach((tip, i) => {
      UI.println(
        `  ${UI.Style.MATRIX_GREEN}${i + 1}.${UI.Style.RESET} ${UI.Style.TEXT_DIM}${tip}${UI.Style.RESET}`
      )
    })
    UI.println()
  }

  /**
   * Completion message
   */
  export function complete() {
    UI.println()
    UI.println(
      UI.box(
        UI.gradient("✓ Setup Complete!", [
          UI.Style.MATRIX_GREEN,
          UI.Style.MATRIX_GREEN_BRIGHT,
        ]) +
        "\n\n" +
        UI.Style.TEXT_DIM +
        "You're all set! RyCode is ready to supercharge your development." +
        UI.Style.RESET,
        { color: UI.Style.MATRIX_GREEN, padding: 1 }
      )
    )
    UI.println()
  }

  /**
   * System check result
   */
  export function systemCheck(checks: Array<{ name: string; status: "ok" | "warn" | "error"; message?: string }>) {
    UI.println()
    UI.println(UI.Style.BOLD + "System Check:" + UI.Style.RESET)
    UI.println()

    const statuses = {
      ok: EnhancedTUI.status("", "online", { showLabel: false }),
      warn: EnhancedTUI.status("", "loading", { showLabel: false }),
      error: EnhancedTUI.status("", "error", { showLabel: false }),
    }

    checks.forEach((check) => {
      UI.println(`  ${statuses[check.status]} ${check.name}`)
      if (check.message) {
        UI.println(`     ${UI.Style.TEXT_DIM}${check.message}${UI.Style.RESET}`)
      }
    })

    UI.println()
  }
}
