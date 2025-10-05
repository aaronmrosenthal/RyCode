#!/usr/bin/env bun

/**
 * Demo: Polished Installer Experience
 * Showcases the professional onboarding flow
 */

import { UI } from "../../src/cli/ui"
import { InstallerMessages } from "../../src/cli/installer-messages"
import { EnhancedTUI } from "../../src/cli/tui-enhanced"

async function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

async function main() {
  // 1. Welcome Screen
  InstallerMessages.welcome()
  await sleep(2000)

  // 2. Show features
  InstallerMessages.features()
  await sleep(2000)

  // 3. Auth required
  InstallerMessages.authRequired()
  await sleep(2000)

  // 4. Provider intro
  InstallerMessages.providerIntro()
  await sleep(1500)

  // 5. Show API key help
  InstallerMessages.showApiKeyHelp("anthropic")
  await sleep(2000)

  // 6. Progress simulation
  UI.println()
  UI.println(UI.Style.BOLD + "Simulating Setup Process:" + UI.Style.RESET)
  UI.println()

  InstallerMessages.progress(1, 3, "Configuring provider...")
  await sleep(1000)

  InstallerMessages.progress(2, 3, "Verifying credentials...")
  await sleep(1000)

  InstallerMessages.progress(3, 3, "Initializing AI models...")
  await sleep(1000)

  // 7. Success
  InstallerMessages.authSuccess("Anthropic Claude")
  await sleep(1500)

  // 8. System check
  InstallerMessages.systemCheck([
    { name: "API Connection", status: "ok", message: "Connected to Anthropic" },
    { name: "Model Access", status: "ok", message: "Claude 3.5 Sonnet available" },
    { name: "Cache Setup", status: "ok", message: "Local cache initialized" },
    { name: "Git Integration", status: "ok", message: "Repository detected" },
  ])
  await sleep(1500)

  // 9. Quick tips
  InstallerMessages.quickTips()
  await sleep(1500)

  // 10. Complete
  InstallerMessages.complete()
  await sleep(1000)

  // 11. Starting
  InstallerMessages.starting()
  await sleep(800)

  // 12. Simulate loading
  const spinner = EnhancedTUI.spinner("Loading AI models...", {
    frames: EnhancedTUI.spinnerFrames.dots,
    color: UI.Style.MATRIX_GREEN,
  })
  spinner.start()
  await sleep(2000)
  spinner.stop("Models loaded successfully!")

  UI.println()

  // 13. Ready
  InstallerMessages.ready()

  // Showcase error/warning/info messages
  UI.println()
  UI.println("═".repeat(80))
  UI.println()
  UI.println(UI.Style.BOLD + "MESSAGE VARIANTS:" + UI.Style.RESET)
  UI.println()
  UI.println("═".repeat(80))

  await sleep(500)

  InstallerMessages.info(
    "Additional Feature Available",
    "You can enable GitHub Copilot integration for even more AI assistance."
  )
  await sleep(800)

  InstallerMessages.warning(
    "Rate Limit Approaching",
    "You've used 80% of your API quota. Consider upgrading your plan."
  )
  await sleep(800)

  InstallerMessages.error(
    "Connection Failed",
    "Could not reach the AI service.",
    "Check your internet connection and try again"
  )
  await sleep(1000)

  // Final showcase
  UI.println()
  UI.println("═".repeat(80))
  UI.println()
  UI.println(UI.Style.BOLD + UI.Style.MATRIX_GREEN + "INSTALLER DEMO COMPLETE" + UI.Style.RESET)
  UI.println()
  UI.println("═".repeat(80))
  UI.println()

  UI.println(UI.logo())
  UI.println()

  UI.println(
    UI.box(
      UI.gradient("✨ Polished Installer Experience ✨", [
        UI.Style.MATRIX_GREEN,
        UI.Style.MATRIX_GREEN_BRIGHT,
      ]) +
      "\n\n" +
      UI.Style.TEXT_DIM +
      "Professional messaging throughout the entire onboarding flow" +
      UI.Style.RESET,
      { color: UI.Style.MATRIX_GREEN, padding: 1 }
    )
  )
  UI.println()

  const features = {
    "Welcome Screen": "✓ Matrix aesthetic with animated header",
    "Clear Steps": "✓ Wizard-style progress indicators",
    "Help Links": "✓ Clickable URLs for API key creation",
    "Progress Bars": "✓ Visual feedback during setup",
    "Status Icons": "✓ Clear success/error/warning states",
    "Quick Tips": "✓ Helpful guidance for new users",
    "Professional Tone": "✓ Friendly yet polished messaging",
  }

  UI.println(UI.Style.BOLD + "Features Demonstrated:" + UI.Style.RESET)
  UI.println()
  UI.println(EnhancedTUI.keyValue(features, { indent: 2, keyColor: UI.Style.MATRIX_GREEN }))
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

main().catch(console.error)
