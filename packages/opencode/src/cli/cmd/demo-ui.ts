/**
 * Demo command to showcase the new toolkit-cli.com inspired UI theme
 */

import { cmd } from "./cmd"
import { UI } from "../ui"
import { CyberpunkPrompts, symbols, formatters } from "../theme"

export const DemoUICommand = cmd({
  command: "demo-ui",
  describe: "showcase the cyberpunk UI theme",
  async handler() {
    UI.empty()

    // Show all logo variants
    UI.println()
    UI.println(UI.glow("=== CLASSIC LOGO ===", UI.Style.CLAUDE_BLUE))
    UI.println()
    UI.println(UI.logo(undefined, "classic"))
    UI.println()

    UI.println()
    UI.println(UI.glow("=== CYBERPUNK LOGO ===", UI.Style.MATRIX_GREEN))
    UI.println()
    UI.println(UI.logo(undefined, "cyberpunk"))
    UI.println()

    UI.println()
    UI.println(UI.glow("=== GRADIENT LOGO ===", UI.Style.GEMINI_GREEN))
    UI.println()
    UI.println(UI.logo(undefined, "gradient"))
    UI.println()

    // Show gradient text
    UI.println()
    CyberpunkPrompts.divider("Gradient Text Examples")
    UI.println()
    UI.println(
      UI.gradient("Matrix Green to Cyan Gradient", [UI.Style.MATRIX_GREEN, UI.Style.NEON_CYAN]),
    )
    UI.println(
      UI.gradient("Multi-color Gradient Effect", [
        UI.Style.MATRIX_GREEN,
        UI.Style.GEMINI_GREEN,
        UI.Style.CLAUDE_BLUE,
        UI.Style.CYBER_PURPLE,
      ]),
    )
    UI.println()

    // Show glow effects
    CyberpunkPrompts.divider("Glow Effects")
    UI.println()
    UI.println(UI.glow("Matrix Green Glow", UI.Style.MATRIX_GREEN))
    UI.println(UI.glow("Claude Blue Glow", UI.Style.CLAUDE_BLUE))
    UI.println(UI.glow("Neon Cyan Glow", UI.Style.NEON_CYAN))
    UI.println(UI.glow("Cyber Purple Glow", UI.Style.CYBER_PURPLE))
    UI.println()

    // Show boxed text
    CyberpunkPrompts.divider("Boxed Messages")
    UI.println()
    UI.println(
      UI.box("This is a Claude Blue box", {
        color: UI.Style.CLAUDE_BLUE,
        padding: 2,
      }),
    )
    UI.println()
    UI.println(
      UI.box("This is a Matrix Green box", {
        color: UI.Style.MATRIX_GREEN,
        padding: 3,
      }),
    )
    UI.println()

    // Show banner
    UI.banner("🚀 OPENCODE INITIALIZED")
    UI.println()

    // Show section headers
    UI.section("Configuration")
    UI.println("  " + UI.Style.TEXT_DIM + "Loading configuration files..." + UI.Style.RESET)
    UI.println("  " + UI.Style.MATRIX_GREEN + "✓" + UI.Style.RESET + " Config loaded")
    UI.println()

    UI.section("Dependencies")
    UI.println("  " + UI.Style.TEXT_DIM + "Checking dependencies..." + UI.Style.RESET)
    UI.println("  " + UI.Style.MATRIX_GREEN + "✓" + UI.Style.RESET + " All dependencies satisfied")
    UI.println()

    // Show symbols
    CyberpunkPrompts.divider("Theme Symbols")
    UI.println()
    UI.println(`  ${symbols.step_submit} Step submitted`)
    UI.println(`  ${symbols.step_active} Step active`)
    UI.println(`  ${symbols.step_cancel} Step cancelled`)
    UI.println(`  ${symbols.step_error} Step error`)
    UI.println()
    UI.println(`  ${symbols.radio_active} Selected radio option`)
    UI.println(`  ${symbols.radio_inactive} Unselected radio option`)
    UI.println()
    UI.println(`  ${symbols.checkbox_active} Active checkbox`)
    UI.println(`  ${symbols.checkbox_selected} Selected checkbox`)
    UI.println(`  ${symbols.checkbox_inactive} Inactive checkbox`)
    UI.println()
    UI.println(`  ${symbols.info} Info message`)
    UI.println(`  ${symbols.success} Success message`)
    UI.println(`  ${symbols.warning} Warning message`)
    UI.println(`  ${symbols.error} Error message`)
    UI.println()

    // Show formatters
    CyberpunkPrompts.divider("Text Formatters")
    UI.println()
    UI.println(formatters.intro("This is an intro message"))
    UI.println(formatters.prompt("This is a prompt message"))
    UI.println(formatters.hint("This is a hint"))
    UI.println(formatters.selected("This is selected text"))
    UI.println(formatters.active("This is active text"))
    UI.println(formatters.success("This is a success message"))
    UI.println(formatters.info("This is an info message"))
    UI.println(formatters.warning("This is a warning message"))
    UI.println(formatters.error("This is an error message"))
    UI.println(formatters.outro("This is an outro message"))
    UI.println()

    // Show typing cursor effect
    CyberpunkPrompts.divider("Interactive Elements")
    UI.println()
    UI.println(UI.withCursor("Typing text here"))
    UI.println()

    // Final message
    UI.println()
    UI.println(
      UI.box(
        UI.gradient("🎨 UI Demo Complete - Toolkit-CLI Aesthetic Applied! 🎨", [
          UI.Style.MATRIX_GREEN,
          UI.Style.GEMINI_GREEN,
          UI.Style.CLAUDE_BLUE,
          UI.Style.CYBER_PURPLE,
        ]),
        { color: UI.Style.NEON_CYAN, padding: 1 },
      ),
    )
    UI.println()
  },
})
