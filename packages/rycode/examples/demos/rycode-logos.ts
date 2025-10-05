#!/usr/bin/env bun

/**
 * RyCode ASCII Art Logo Designs
 * Multiple killer logo variants for RyCode
 */

import { UI } from "../../src/cli/ui"

// Logo Variant 1: Bold Blocky (Classic)
const RYCODE_CLASSIC = [
  `██████╗ ██╗   ██╗ ██████╗ ██████╗ ██████╗ ███████╗`,
  `██╔══██╗╚██╗ ██╔╝██╔════╝██╔═══██╗██╔══██╗██╔════╝`,
  `██████╔╝ ╚████╔╝ ██║     ██║   ██║██║  ██║█████╗  `,
  `██╔══██╗  ╚██╔╝  ██║     ██║   ██║██║  ██║██╔══╝  `,
  `██║  ██║   ██║   ╚██████╗╚██████╔╝██████╔╝███████╗`,
  `╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝`,
]

// Logo Variant 2: Cyberpunk Boxed
const RYCODE_CYBERPUNK = [
  `╔════════════════════════════════════════════════════════════════╗`,
  `║                                                                ║`,
  `║  ▓█████▄   ▓██   ██▓    ▄████▄   ▒█████  ▓█████▄ ▓████▄      ║`,
  `║  ▒██▀ ██▌   ▒██  ██▒   ▒██▀ ▀█  ▒██▒  ██▒▒██▀ ██▌ ██▀ ██▌    ║`,
  `║  ░██   █▌   ░██ ██░    ▒▓█    ▄ ▒██░  ██▒░██   █▌░██   █▌    ║`,
  `║  ░▓█▄   ▌   ░ ▐██▓░    ▒▓▓▄ ▄██▒▒██   ██░░▓█▄   ▌░▓█▄   ▌    ║`,
  `║  ░▒████▓    ░ ██▒▓░    ▒ ▓███▀ ░░ ████▓▒░░▒████▓ ░▒████▓     ║`,
  `║   ▒▒▓  ▒     ██▒▒▒     ░ ░▒ ▒  ░░ ▒░▒░▒░  ▒▒▓  ▒  ▒▒▓  ▒     ║`,
  `║   ░ ▒  ▒   ▓██ ░▒░       ░  ▒     ░ ▒ ▒░  ░ ▒  ▒  ░ ▒  ▒     ║`,
  `║   ░ ░  ░   ▒ ▒ ░░      ░        ░ ░ ░ ▒   ░ ░  ░  ░ ░  ░     ║`,
  `║     ░      ░ ░         ░ ░          ░ ░     ░       ░         ║`,
  `║   ░        ░ ░         ░                  ░       ░           ║`,
  `║                                                                ║`,
  `╚════════════════════════════════════════════════════════════════╝`,
]

// Logo Variant 3: Neon Style
const RYCODE_NEON = [
  ``,
  `  ██▀███ ▓██   ██▓ ▄████▄   ▒█████  ▓█████▄ ▓█████ `,
  ` ▓██ ▒ ██▒▒██  ██▒▒██▀ ▀█  ▒██▒  ██▒▒██▀ ██▌▓█   ▀ `,
  ` ▓██ ░▄█ ▒ ▒██ ██░▒▓█    ▄ ▒██░  ██▒░██   █▌▒███   `,
  ` ▒██▀▀█▄   ░ ▐██▓░▒▓▓▄ ▄██▒▒██   ██░░▓█▄   ▌▒▓█  ▄ `,
  ` ░██▓ ▒██▒ ░ ██▒▓░▒ ▓███▀ ░░ ████▓▒░░▒████▓ ░▒████▒`,
  ` ░ ▒▓ ░▒▓░  ██▒▒▒ ░ ░▒ ▒  ░░ ▒░▒░▒░  ▒▒▓  ▒ ░░ ▒░ ░`,
  `   ░▒ ░ ▒░▓██ ░▒░   ░  ▒     ░ ▒ ▒░  ░ ▒  ▒  ░ ░  ░`,
  `   ░░   ░ ▒ ▒ ░░  ░        ░ ░ ░ ▒   ░ ░  ░    ░   `,
  `    ░     ░ ░     ░ ░          ░ ░     ░       ░  ░`,
  `          ░ ░     ░                  ░              `,
  ``,
]

// Logo Variant 4: Compact & Modern
const RYCODE_MODERN = [
  ``,
  `  ██████╗ ██╗   ██╗ ██████╗ ██████╗ ██████╗ ███████╗`,
  `  ██╔══██╗╚██╗ ██╔╝██╔════╝██╔═══██╗██╔══██╗██╔════╝`,
  `  ██████╔╝ ╚████╔╝ ██║     ██║   ██║██║  ██║█████╗  `,
  `  ██╔══██╗  ╚██╔╝  ██║     ██║   ██║██║  ██║██╔══╝  `,
  `  ██║  ██║   ██║   ╚██████╗╚██████╔╝██████╔╝███████╗`,
  `  ╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝`,
  ``,
]

// Logo Variant 5: Slant Style (Super Stylish)
const RYCODE_SLANT = [
  ``,
  `    ____                 ______          __   `,
  `   / __ \\__  __  _____  / ____/___  ____/ /__ `,
  `  / /_/ / / / / / ___/ / /   / __ \\/ __  / _ \\`,
  ` / _, _/ /_/ / / /__  / /___/ /_/ / /_/ /  __/`,
  `/_/ |_|\\__, /  \\___/  \\____/\\____/\\__,_/\\___/ `,
  `      /____/                                   `,
  ``,
]

// Logo Variant 6: Big & Bold
const RYCODE_BIG = [
  ``,
  `  ██████╗  ██╗   ██╗  ██████╗  ██████╗  ██████╗  ███████╗`,
  `  ██╔══██╗ ╚██╗ ██╔╝ ██╔════╝ ██╔═══██╗ ██╔══██╗ ██╔════╝`,
  `  ██████╔╝  ╚████╔╝  ██║      ██║   ██║ ██║  ██║ █████╗  `,
  `  ██╔══██╗   ╚██╔╝   ██║      ██║   ██║ ██║  ██║ ██╔══╝  `,
  `  ██║  ██║    ██║    ╚██████╗ ╚██████╔╝ ██████╔╝ ███████╗`,
  `  ╚═╝  ╚═╝    ╚═╝     ╚═════╝  ╚═════╝  ╚═════╝  ╚══════╝`,
  ``,
]

// Logo Variant 7: Matrix Digital Rain Style
const RYCODE_MATRIX = [
  `╔═══════════════════════════════════════════════════════════╗`,
  `║  ██▀███   ▓██   ██▓ ▄████▄   ▒█████  ▓█████▄ ▓█████     ║`,
  `║ ▓██ ▒ ██▒  ▒██  ██▒▒██▀ ▀█  ▒██▒  ██▒▒██▀ ██▌▓█   ▀     ║`,
  `║ ▓██ ░▄█ ▒   ▒██ ██░▒▓█    ▄ ▒██░  ██▒░██   █▌▒███       ║`,
  `║ ▒██▀▀█▄     ░ ▐██▓░▒▓▓▄ ▄██▒▒██   ██░░▓█▄   ▌▒▓█  ▄     ║`,
  `║ ░██▓ ▒██▒   ░ ██▒▓░▒ ▓███▀ ░░ ████▓▒░░▒████▓ ░▒████▒    ║`,
  `║ ░ ▒▓ ░▒▓░    ██▒▒▒ ░ ░▒ ▒  ░░ ▒░▒░▒░  ▒▒▓  ▒ ░░ ▒░ ░    ║`,
  `║   ░▒ ░ ▒░  ▓██ ░▒░   ░  ▒     ░ ▒ ▒░  ░ ▒  ▒  ░ ░  ░    ║`,
  `║   ░░   ░   ▒ ▒ ░░  ░        ░ ░ ░ ▒   ░ ░  ░    ░       ║`,
  `║    ░       ░ ░     ░ ░          ░ ░     ░       ░  ░    ║`,
  `╚═══════════════════════════════════════════════════════════╝`,
]

// Logo Variant 8: Double-Strike (Most Impactful)
const RYCODE_DOUBLE = [
  ``,
  `  ██████╗ ██╗   ██╗ ██████╗███████╗██████╗ ██████╗ ███████╗`,
  `  ██╔══██╗╚██╗ ██╔╝██╔════╝██╔════╝██╔══██╗██╔══██╗██╔════╝`,
  `  ██████╔╝ ╚████╔╝ ██║     █████╗  ██████╔╝██║  ██║█████╗  `,
  `  ██╔══██╗  ╚██╔╝  ██║     ██╔══╝  ██╔══██╗██║  ██║██╔══╝  `,
  `  ██║  ██║   ██║   ╚██████╗███████╗██║  ██║██████╔╝███████╗`,
  `  ╚═╝  ╚═╝   ╚═╝    ╚═════╝╚══════╝╚═╝  ╚═╝╚═════╝ ╚══════╝`,
  `  ═════════════════════════════════════════════════════════`,
  ``,
]

// Helper function to print logo with colors
function printLogo(
  lines: string[],
  name: string,
  colorStyle: "gradient" | "matrix" | "neon" | "classic" | "cyber" = "gradient"
) {
  UI.println()
  UI.println(UI.Style.BOLD + UI.Style.NEON_CYAN + `${name}:` + UI.Style.RESET)
  UI.println()

  if (colorStyle === "gradient") {
    // Multi-color gradient
    const colors = [
      UI.Style.MATRIX_GREEN,
      UI.Style.GEMINI_GREEN,
      UI.Style.NEON_CYAN,
      UI.Style.CLAUDE_BLUE,
      UI.Style.CYBER_PURPLE,
    ]
    lines.forEach((line, i) => {
      const colorIndex = Math.floor((i / lines.length) * (colors.length - 1))
      UI.println(colors[colorIndex] + line + UI.Style.RESET)
    })
  } else if (colorStyle === "matrix") {
    // Matrix green
    lines.forEach((line) => {
      UI.println(UI.Style.MATRIX_GREEN + line + UI.Style.RESET)
    })
  } else if (colorStyle === "neon") {
    // Neon cyan
    lines.forEach((line) => {
      UI.println(UI.Style.NEON_CYAN + line + UI.Style.RESET)
    })
  } else if (colorStyle === "cyber") {
    // Alternating cyber colors
    lines.forEach((line, i) => {
      const color = i % 2 === 0 ? UI.Style.CYBER_PURPLE : UI.Style.NEON_MAGENTA
      UI.println(color + line + UI.Style.RESET)
    })
  } else {
    // Classic - alternating matrix green and blue
    lines.forEach((line, i) => {
      const color = i % 2 === 0 ? UI.Style.MATRIX_GREEN : UI.Style.CLAUDE_BLUE
      UI.println(color + line + UI.Style.RESET)
    })
  }

  UI.println()
}

// Main showcase
UI.empty()
UI.println()
UI.println(
  UI.box(UI.glow("🎨 RYCODE ASCII ART SHOWCASE 🎨", UI.Style.MATRIX_GREEN), {
    color: UI.Style.CYBER_PURPLE,
    padding: 2,
  })
)
UI.println()

printLogo(RYCODE_MODERN, "1. MODERN (Recommended)", "gradient")
printLogo(RYCODE_MATRIX, "2. MATRIX BOXED", "matrix")
printLogo(RYCODE_NEON, "3. NEON STYLE", "neon")
printLogo(RYCODE_BIG, "4. BIG & BOLD", "classic")
printLogo(RYCODE_SLANT, "5. SLANT STYLE", "cyber")
printLogo(RYCODE_CYBERPUNK, "6. CYBERPUNK BOXED", "gradient")
printLogo(RYCODE_CLASSIC, "7. CLASSIC BLOCKY", "gradient")
printLogo(RYCODE_DOUBLE, "8. DOUBLE-STRIKE", "neon")

// Show the recommended one with special treatment
UI.println()
UI.println("═".repeat(70))
UI.println()
UI.println(UI.Style.BOLD + UI.Style.MATRIX_GREEN + "🌟 RECOMMENDED LOGO:" + UI.Style.RESET)
UI.println()

// Gradient version of the modern logo
RYCODE_MODERN.forEach((line, i) => {
  const colors = [
    UI.Style.MATRIX_GREEN,
    UI.Style.GEMINI_GREEN,
    UI.Style.NEON_CYAN,
    UI.Style.CLAUDE_BLUE,
    UI.Style.CYBER_PURPLE,
  ]
  const colorIndex = Math.floor((i / RYCODE_MODERN.length) * (colors.length - 1))
  UI.println("  " + colors[colorIndex] + line + UI.Style.RESET)
})

UI.println()
UI.println(
  UI.Style.TEXT_DIM +
    "  AI-Powered Development Assistant • Cyberpunk Aesthetic • Next-Gen CLI" +
    UI.Style.RESET
)
UI.println()
UI.println("═".repeat(70))
UI.println()

// Export the logos
export {
  RYCODE_CLASSIC,
  RYCODE_CYBERPUNK,
  RYCODE_NEON,
  RYCODE_MODERN,
  RYCODE_SLANT,
  RYCODE_BIG,
  RYCODE_MATRIX,
  RYCODE_DOUBLE,
}
