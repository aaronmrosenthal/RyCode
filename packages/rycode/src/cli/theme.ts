/**
 * Toolkit-CLI.com inspired theme for @clack/prompts
 * Cyberpunk/Matrix aesthetic with vibrant accent colors
 */

import { UI } from "./ui"

/**
 * Custom symbols for @clack/prompts with cyberpunk aesthetic
 */
export const symbols = {
  // Use matrix-style characters
  step_active: UI.Style.MATRIX_GREEN + "▸" + UI.Style.RESET,
  step_cancel: UI.Style.TEXT_DANGER + "✖" + UI.Style.RESET,
  step_error: UI.Style.TEXT_DANGER + "●" + UI.Style.RESET,
  step_submit: UI.Style.CLAUDE_BLUE + "✓" + UI.Style.RESET,

  // Box drawing characters with matrix theme
  bar: UI.Style.MATRIX_GREEN_DIM + "│" + UI.Style.RESET,
  bar_end: UI.Style.MATRIX_GREEN_DIM + "└" + UI.Style.RESET,

  // Radio and checkbox with neon styling
  radio_active: UI.Style.NEON_CYAN + "●" + UI.Style.RESET,
  radio_inactive: UI.Style.TEXT_DIM + "○" + UI.Style.RESET,
  checkbox_active: UI.Style.MATRIX_GREEN + "◉" + UI.Style.RESET,
  checkbox_selected: UI.Style.CLAUDE_BLUE + "◈" + UI.Style.RESET,
  checkbox_inactive: UI.Style.TEXT_DIM + "◯" + UI.Style.RESET,

  // Password field
  password_mask: UI.Style.CYBER_PURPLE + "●" + UI.Style.RESET,

  // Arrows with gradient colors
  arrow_down: UI.Style.GEMINI_GREEN + "▼" + UI.Style.RESET,
  arrow_left: UI.Style.GEMINI_GREEN + "◀" + UI.Style.RESET,
  arrow_right: UI.Style.GEMINI_GREEN + "▶" + UI.Style.RESET,
  arrow_up: UI.Style.GEMINI_GREEN + "▲" + UI.Style.RESET,

  // Indicators
  info: UI.Style.CLAUDE_BLUE + "ℹ" + UI.Style.RESET,
  success: UI.Style.MATRIX_GREEN + "✓" + UI.Style.RESET,
  warning: UI.Style.TEXT_WARNING + "⚠" + UI.Style.RESET,
  error: UI.Style.TEXT_DANGER + "✖" + UI.Style.RESET,
}

/**
 * Format text with cyberpunk styling
 */
export const formatters = {
  /**
   * Style intro messages with matrix green
   */
  intro(message: string): string {
    return `${UI.Style.MATRIX_GREEN}${UI.Style.BOLD}┌${UI.Style.RESET} ${UI.gradient(message, [
      UI.Style.MATRIX_GREEN,
      UI.Style.CLAUDE_BLUE,
    ])}`
  },

  /**
   * Style outro messages with success color
   */
  outro(message: string): string {
    return `${UI.Style.CLAUDE_BLUE}${UI.Style.BOLD}└${UI.Style.RESET} ${UI.glow(message, UI.Style.MATRIX_GREEN)}`
  },

  /**
   * Style prompts with cyberpunk accent
   */
  prompt(message: string): string {
    return `${UI.Style.CLAUDE_BLUE}▸${UI.Style.RESET} ${message}`
  },

  /**
   * Style hints with dim text
   */
  hint(hint: string): string {
    return `${UI.Style.TEXT_DIM}${hint}${UI.Style.RESET}`
  },

  /**
   * Style selected items with matrix green
   */
  selected(text: string): string {
    return `${UI.Style.MATRIX_GREEN}${text}${UI.Style.RESET}`
  },

  /**
   * Style active items with bold cyan
   */
  active(text: string): string {
    return `${UI.Style.BOLD}${UI.Style.NEON_CYAN}${text}${UI.Style.RESET}`
  },

  /**
   * Style errors with danger color
   */
  error(message: string): string {
    return `${UI.Style.TEXT_DANGER}${message}${UI.Style.RESET}`
  },

  /**
   * Style success messages with matrix green
   */
  success(message: string): string {
    return `${UI.Style.MATRIX_GREEN}${message}${UI.Style.RESET}`
  },

  /**
   * Style info messages with blue
   */
  info(message: string): string {
    return `${UI.Style.CLAUDE_BLUE}${message}${UI.Style.RESET}`
  },

  /**
   * Style warnings with yellow
   */
  warning(message: string): string {
    return `${UI.Style.TEXT_WARNING}${message}${UI.Style.RESET}`
  },
}

/**
 * Enhanced prompts wrappers with cyberpunk styling
 * These wrap @clack/prompts functions to apply custom theming
 */
export namespace CyberpunkPrompts {
  /**
   * Display a cyberpunk-styled intro with optional ASCII art
   */
  export function intro(message: string, showArt: boolean = false) {
    if (showArt) {
      UI.println(UI.logo(undefined, "gradient"))
      UI.println()
    }
    UI.println(formatters.intro(message))
  }

  /**
   * Display a cyberpunk-styled outro
   */
  export function outro(message: string) {
    UI.println(formatters.outro(message))
  }

  /**
   * Display a section divider with matrix styling
   */
  export function divider(text?: string) {
    const line = UI.Style.MATRIX_GREEN_DIM + "─".repeat(50) + UI.Style.RESET
    UI.println()
    if (text) {
      const styledText = UI.gradient(text, [UI.Style.MATRIX_GREEN, UI.Style.CLAUDE_BLUE])
      UI.println(`${line}`)
      UI.println(`${UI.Style.MATRIX_GREEN}┤${UI.Style.RESET} ${styledText} ${UI.Style.MATRIX_GREEN}├${UI.Style.RESET}`)
      UI.println(`${line}`)
    } else {
      UI.println(line)
    }
  }

  /**
   * Display a loading spinner with matrix styling
   */
  export function loading(message: string) {
    const frames = ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"]
    let i = 0
    return {
      message: (msg: string) => {
        message = msg
      },
      stop: () => {
        UI.println(`${UI.Style.MATRIX_GREEN}✓${UI.Style.RESET} ${message}`)
      },
      // Note: Actual animation would require setInterval, this is a simplified version
      frame: () => {
        return `${UI.Style.NEON_CYAN}${frames[i++ % frames.length]}${UI.Style.RESET} ${message}`
      },
    }
  }
}
