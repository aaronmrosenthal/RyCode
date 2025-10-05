#!/usr/bin/env bun

/**
 * RyCode Matrix Gradient Logo Demo
 * Showcase the Matrix digital rain aesthetic
 */

import { UI } from "./src/cli/ui"

async function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

async function main() {
  UI.empty()

  // Opening animation
  UI.println()
  UI.println(UI.Style.MATRIX_GREEN_DARK + "█".repeat(80) + UI.Style.RESET)
  UI.println(UI.Style.MATRIX_GREEN_DIM + "█".repeat(80) + UI.Style.RESET)
  UI.println(UI.Style.MATRIX_GREEN + "█".repeat(80) + UI.Style.RESET)
  UI.println(UI.Style.MATRIX_GREEN_BRIGHT + "█".repeat(80) + UI.Style.RESET)
  UI.println(UI.Style.MATRIX_GREEN_LIGHT + "█".repeat(80) + UI.Style.RESET)
  UI.println()

  await sleep(300)

  // Title
  UI.println(
    UI.box(
      UI.glow("⚡ RYCODE - MATRIX DIGITAL RAIN AESTHETIC ⚡", UI.Style.MATRIX_GREEN_BRIGHT),
      { color: UI.Style.MATRIX_GREEN, padding: 2 }
    )
  )
  UI.println()

  await sleep(500)

  // Gradient showcase
  UI.println(UI.Style.BOLD + UI.Style.MATRIX_GREEN + "MATRIX GRADIENT SPECTRUM:" + UI.Style.RESET)
  UI.println()

  const gradientDemo = [
    { name: "DARK", color: UI.Style.MATRIX_GREEN_DARK, hex: "#00641E" },
    { name: "DIM", color: UI.Style.MATRIX_GREEN_DIM, hex: "#00B432" },
    { name: "STANDARD", color: UI.Style.MATRIX_GREEN, hex: "#00FF41" },
    { name: "BRIGHT", color: UI.Style.MATRIX_GREEN_BRIGHT, hex: "#96FFB4" },
    { name: "LIGHT", color: UI.Style.MATRIX_GREEN_LIGHT, hex: "#64FF96" },
  ]

  gradientDemo.forEach((item) => {
    UI.println(
      `  ${item.color}████████${UI.Style.RESET}  ${UI.Style.BOLD}${item.name.padEnd(12)}${UI.Style.RESET} ${UI.Style.TEXT_DIM}${item.hex}${UI.Style.RESET}`
    )
  })
  UI.println()

  await sleep(800)

  // Main logo with matrix gradient
  UI.println("═".repeat(80))
  UI.println()
  UI.println(
    UI.Style.BOLD +
      UI.Style.MATRIX_GREEN_BRIGHT +
      "        THE MATRIX DIGITAL RAIN LOGO" +
      UI.Style.RESET
  )
  UI.println()
  UI.println("═".repeat(80))
  UI.println()

  await sleep(500)

  // Show the logo
  UI.println(UI.logo())
  UI.println()

  await sleep(800)

  // Tagline with matrix gradient
  const tagline = "AI-POWERED DEVELOPMENT ASSISTANT"
  UI.println("  " + UI.gradient(tagline, [
    UI.Style.MATRIX_GREEN_DARK,
    UI.Style.MATRIX_GREEN,
    UI.Style.MATRIX_GREEN_BRIGHT,
    UI.Style.MATRIX_GREEN_LIGHT,
  ]))
  UI.println()

  UI.println("═".repeat(80))
  UI.println()

  await sleep(500)

  // Show comparison
  UI.println(UI.Style.BOLD + "COMPARISON:" + UI.Style.RESET)
  UI.println()

  UI.println(UI.Style.MATRIX_GREEN + "1. MATRIX GRADIENT (Default):" + UI.Style.RESET)
  UI.println()
  UI.println(UI.logo(undefined, "rycode"))
  UI.println()
  UI.println(UI.Style.TEXT_DIM + "  Dark green at top → bright at center → light fade" + UI.Style.RESET)
  UI.println()

  await sleep(800)

  UI.println(UI.Style.MATRIX_GREEN + "2. BIG MATRIX:" + UI.Style.RESET)
  UI.println()
  UI.println(UI.logo(undefined, "rycode-big"))
  UI.println()
  UI.println(UI.Style.TEXT_DIM + "  Same matrix gradient, bigger spacing" + UI.Style.RESET)
  UI.println()

  await sleep(800)

  // Digital rain effect visualization
  UI.println("═".repeat(80))
  UI.println()
  UI.println(UI.Style.BOLD + UI.Style.MATRIX_GREEN + "DIGITAL RAIN VISUALIZATION:" + UI.Style.RESET)
  UI.println()

  // Create digital rain effect with the gradient
  const rainColumns = 20
  for (let row = 0; row < 10; row++) {
    let line = "  "
    for (let col = 0; col < rainColumns; col++) {
      const offset = (row + col) % gradientDemo.length
      const color = gradientDemo[offset].color
      const char = ["0", "1", "ﾊ", "ﾐ", "ﾋ", "ｰ", "ｳ", "ｼ", "ﾅ", "ﾓ"][Math.floor(Math.random() * 10)]
      line += color + char + " " + UI.Style.RESET
    }
    UI.println(line)
  }
  UI.println()

  await sleep(800)

  // Features
  UI.println("═".repeat(80))
  UI.println()
  UI.println(UI.Style.BOLD + "MATRIX AESTHETIC FEATURES:" + UI.Style.RESET)
  UI.println()

  const features = [
    "⚡ Authentic Matrix digital rain gradient",
    "🌊 Flowing dark → bright → light green cascade",
    "✨ Multiple green shades for depth and dimension",
    "💚 Pure green palette - no rainbow distraction",
    "🎯 Toolkit-CLI.com inspired color scheme",
    "🔥 Killer visual impact in any terminal",
  ]

  features.forEach((feature) => {
    UI.println(`  ${UI.Style.MATRIX_GREEN}${feature}${UI.Style.RESET}`)
  })
  UI.println()

  await sleep(500)

  // Code usage
  UI.println("═".repeat(80))
  UI.println()
  UI.println(UI.Style.BOLD + "USAGE:" + UI.Style.RESET)
  UI.println()

  UI.println(UI.Style.MATRIX_GREEN_DIM + "  // Default - Matrix gradient" + UI.Style.RESET)
  UI.println(UI.Style.MATRIX_GREEN + '  UI.logo()' + UI.Style.RESET)
  UI.println()
  UI.println(UI.Style.MATRIX_GREEN_DIM + "  // Explicit matrix style" + UI.Style.RESET)
  UI.println(UI.Style.MATRIX_GREEN + '  UI.logo(undefined, "matrix")' + UI.Style.RESET)
  UI.println()
  UI.println(UI.Style.MATRIX_GREEN_DIM + "  // Big matrix logo" + UI.Style.RESET)
  UI.println(UI.Style.MATRIX_GREEN + '  UI.logo(undefined, "rycode-big")' + UI.Style.RESET)
  UI.println()

  UI.println("═".repeat(80))
  UI.println()

  // Final showcase
  UI.println()
  UI.println(UI.logo())
  UI.println()
  UI.println(
    "  " +
      UI.gradient(">>> WELCOME TO THE MATRIX <<<", [
        UI.Style.MATRIX_GREEN_DARK,
        UI.Style.MATRIX_GREEN,
        UI.Style.MATRIX_GREEN_BRIGHT,
        UI.Style.MATRIX_GREEN,
        UI.Style.MATRIX_GREEN_DARK,
      ])
  )
  UI.println()

  // Success banner
  UI.println()
  UI.println(
    UI.box(
      UI.glow("🟢 MATRIX GRADIENT ACTIVE - READY TO CODE 🟢", UI.Style.MATRIX_GREEN_BRIGHT),
      { color: UI.Style.MATRIX_GREEN, padding: 1 }
    )
  )
  UI.println()

  // Closing rain effect
  UI.println(UI.Style.MATRIX_GREEN_LIGHT + "█".repeat(80) + UI.Style.RESET)
  UI.println(UI.Style.MATRIX_GREEN_BRIGHT + "█".repeat(80) + UI.Style.RESET)
  UI.println(UI.Style.MATRIX_GREEN + "█".repeat(80) + UI.Style.RESET)
  UI.println(UI.Style.MATRIX_GREEN_DIM + "█".repeat(80) + UI.Style.RESET)
  UI.println(UI.Style.MATRIX_GREEN_DARK + "█".repeat(80) + UI.Style.RESET)
  UI.println()
}

main().catch(console.error)
