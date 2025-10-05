#!/usr/bin/env bun

/**
 * RyCode Logo Demo
 * Showcase the new killer RyCode ASCII art
 */

import { UI } from "../../src/cli/ui"

async function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

async function main() {
  UI.empty()

  // Title card
  UI.println()
  UI.println(
    UI.box(UI.glow("🚀 RYCODE - AI-POWERED DEV ASSISTANT 🚀", UI.Style.MATRIX_GREEN), {
      color: UI.Style.CYBER_PURPLE,
      padding: 2,
    })
  )
  UI.println()

  await sleep(500)

  // Show all logo variants
  UI.println("═".repeat(80))
  UI.println()
  UI.println(UI.Style.BOLD + UI.Style.NEON_CYAN + "LOGO VARIANTS:" + UI.Style.RESET)
  UI.println()
  UI.println("═".repeat(80))

  // 1. Modern (Default & Recommended)
  UI.println()
  UI.println(UI.Style.BOLD + "1. MODERN (Default - Recommended):" + UI.Style.RESET)
  UI.println()
  UI.println(UI.logo(undefined, "rycode-modern"))
  UI.println(
    UI.Style.TEXT_DIM +
      "  Clean, professional gradient from matrix green → cyan → blue → purple" +
      UI.Style.RESET
  )
  UI.println()

  await sleep(800)

  // 2. Big & Bold
  UI.println()
  UI.println(UI.Style.BOLD + "2. BIG & BOLD:" + UI.Style.RESET)
  UI.println()
  UI.println(UI.logo(undefined, "rycode-big"))
  UI.println(
    UI.Style.TEXT_DIM + "  Extra spacing for maximum impact - perfect for splash screens" + UI.Style.RESET
  )
  UI.println()

  await sleep(800)

  // 3. Slant Style
  UI.println()
  UI.println(UI.Style.BOLD + "3. SLANT (Italic/Modern):" + UI.Style.RESET)
  UI.println()
  UI.println(UI.logo(undefined, "rycode-slant"))
  UI.println(UI.Style.TEXT_DIM + "  Sleek italic design with cyber purple/magenta theme" + UI.Style.RESET)
  UI.println()

  await sleep(800)

  // 4. Cyberpunk Boxed
  UI.println()
  UI.println(UI.Style.BOLD + "4. CYBERPUNK BOXED:" + UI.Style.RESET)
  UI.println()
  UI.println(UI.logo(undefined, "cyberpunk"))
  UI.println(
    UI.Style.TEXT_DIM + "  Full cyberpunk aesthetic with border box - maximum presence" + UI.Style.RESET
  )
  UI.println()

  await sleep(800)

  // 5. Gradient Cyberpunk
  UI.println()
  UI.println(UI.Style.BOLD + "5. GRADIENT BOXED:" + UI.Style.RESET)
  UI.println()
  UI.println(UI.logo(undefined, "gradient"))
  UI.println(UI.Style.TEXT_DIM + "  Boxed version with rainbow gradient effect" + UI.Style.RESET)
  UI.println()

  await sleep(800)

  // 6. Classic (Backwards Compatibility)
  UI.println()
  UI.println(UI.Style.BOLD + "6. CLASSIC (OpenCode Legacy):" + UI.Style.RESET)
  UI.println()
  UI.println(UI.logo(undefined, "classic"))
  UI.println(UI.Style.TEXT_DIM + "  Original OpenCode logo for backwards compatibility" + UI.Style.RESET)
  UI.println()

  await sleep(800)

  // Recommended Usage Section
  UI.println()
  UI.println("═".repeat(80))
  UI.println()
  UI.println(UI.Style.BOLD + UI.Style.MATRIX_GREEN + "🌟 RECOMMENDED USAGE:" + UI.Style.RESET)
  UI.println()
  UI.println("═".repeat(80))
  UI.println()

  // Show the default in action
  UI.println(UI.logo()) // Default is now "rycode"
  UI.println()
  UI.println(
    UI.Style.TEXT_INFO + "  Simply call " + UI.Style.BOLD + "UI.logo()" + UI.Style.RESET + UI.Style.TEXT_INFO + " to get the modern RyCode logo!" + UI.Style.RESET
  )
  UI.println()

  // Code examples
  UI.println()
  UI.println(UI.Style.BOLD + "CODE EXAMPLES:" + UI.Style.RESET)
  UI.println()
  UI.println(UI.Style.TEXT_DIM + "  // Use the default modern logo" + UI.Style.RESET)
  UI.println(UI.Style.NEON_CYAN + '  UI.println(UI.logo())' + UI.Style.RESET)
  UI.println()
  UI.println(UI.Style.TEXT_DIM + "  // Use a specific variant" + UI.Style.RESET)
  UI.println(UI.Style.NEON_CYAN + '  UI.println(UI.logo(undefined, "rycode-big"))' + UI.Style.RESET)
  UI.println()
  UI.println(UI.Style.TEXT_DIM + "  // With padding" + UI.Style.RESET)
  UI.println(UI.Style.NEON_CYAN + '  UI.println(UI.logo("  ", "cyberpunk"))' + UI.Style.RESET)
  UI.println()

  // Feature callouts
  UI.println()
  UI.println("═".repeat(80))
  UI.println()
  UI.println(UI.Style.BOLD + "KEY FEATURES:" + UI.Style.RESET)
  UI.println()

  const features = [
    { icon: "✨", text: "Multiple logo variants for different contexts" },
    { icon: "🎨", text: "Beautiful gradient effects using cyberpunk color palette" },
    { icon: "🔥", text: "Killer visual impact - stands out in any terminal" },
    { icon: "⚡", text: "Lightweight - just ASCII art, no dependencies" },
    { icon: "🎯", text: "Professional yet edgy - perfect for dev tools" },
    { icon: "🌈", text: "Fully customizable colors and styles" },
  ]

  features.forEach((feature) => {
    UI.println(`  ${feature.icon}  ${feature.text}`)
  })

  UI.println()

  // Color palette showcase
  UI.println()
  UI.println("═".repeat(80))
  UI.println()
  UI.println(UI.Style.BOLD + "COLOR PALETTE:" + UI.Style.RESET)
  UI.println()

  UI.println(
    `  ${UI.Style.MATRIX_GREEN}█████${UI.Style.RESET} Matrix Green     ${UI.Style.MATRIX_GREEN}#00FF41${UI.Style.RESET}`
  )
  UI.println(
    `  ${UI.Style.GEMINI_GREEN}█████${UI.Style.RESET} Gemini Green     ${UI.Style.GEMINI_GREEN}#1AD689${UI.Style.RESET}`
  )
  UI.println(
    `  ${UI.Style.NEON_CYAN}█████${UI.Style.RESET} Neon Cyan        ${UI.Style.NEON_CYAN}#00FFFF${UI.Style.RESET}`
  )
  UI.println(
    `  ${UI.Style.CLAUDE_BLUE}█████${UI.Style.RESET} Claude Blue      ${UI.Style.CLAUDE_BLUE}#6699FF${UI.Style.RESET}`
  )
  UI.println(
    `  ${UI.Style.CYBER_PURPLE}█████${UI.Style.RESET} Cyber Purple     ${UI.Style.CYBER_PURPLE}#9333EA${UI.Style.RESET}`
  )
  UI.println(
    `  ${UI.Style.NEON_MAGENTA}█████${UI.Style.RESET} Neon Magenta     ${UI.Style.NEON_MAGENTA}#FF00FF${UI.Style.RESET}`
  )
  UI.println()

  // Final showcase with tagline
  UI.println()
  UI.println("═".repeat(80))
  UI.println()

  UI.println(UI.logo(undefined, "rycode"))
  UI.println()

  const tagline = "AI-Powered Development Assistant • Cyberpunk Aesthetic • Next-Gen CLI"
  UI.println("  " + UI.gradient(tagline, [UI.Style.MATRIX_GREEN, UI.Style.NEON_CYAN, UI.Style.CYBER_PURPLE]))
  UI.println()

  UI.println("═".repeat(80))
  UI.println()

  // Success message
  UI.println(
    UI.box(
      UI.gradient("🎨 Killer ASCII Art Loaded! Ready to ship! 🚀", [
        UI.Style.MATRIX_GREEN,
        UI.Style.GEMINI_GREEN,
        UI.Style.NEON_CYAN,
        UI.Style.CLAUDE_BLUE,
        UI.Style.CYBER_PURPLE,
      ]),
      { color: UI.Style.NEON_CYAN, padding: 1 }
    )
  )
  UI.println()
}

main().catch(console.error)
