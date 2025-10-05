import z from "zod/v4"
import { EOL } from "os"
import { NamedError } from "../util/error"

export namespace UI {
  const LOGO = [
    [`█▀▀█ █▀▀█ █▀▀ █▀▀▄ `, `█▀▀ █▀▀█ █▀▀▄ █▀▀`],
    [`█░░█ █░░█ █▀▀ █░░█ `, `█░░ █░░█ █░░█ █▀▀`],
    [`▀▀▀▀ █▀▀▀ ▀▀▀ ▀  ▀ `, `▀▀▀ ▀▀▀▀ ▀▀▀  ▀▀▀`],
  ]

  // Cyberpunk/Matrix-inspired ASCII art logo alternative
  const LOGO_CYBERPUNK = [
    `╔═══════════════════════════════════╗`,
    `║  ▓▓▓▓▓  ▓▓▓▓▓  ▓▓▓▓▓  ▓    ▓     ║`,
    `║  ▓   ▓  ▓   ▓  ▓      ▓▓   ▓     ║`,
    `║  ▓   ▓  ▓▓▓▓▓  ▓▓▓    ▓ ▓  ▓     ║`,
    `║  ▓   ▓  ▓      ▓      ▓  ▓ ▓     ║`,
    `║  ▓▓▓▓▓  ▓      ▓▓▓▓▓  ▓   ▓▓     ║`,
    `║                                   ║`,
    `║  ▓▓▓▓▓  ▓▓▓▓▓  ▓▓▓▓▓  ▓▓▓▓▓      ║`,
    `║  ▓      ▓   ▓  ▓   ▓  ▓          ║`,
    `║  ▓      ▓   ▓  ▓   ▓  ▓▓▓        ║`,
    `║  ▓      ▓   ▓  ▓   ▓  ▓          ║`,
    `║  ▓▓▓▓▓  ▓▓▓▓▓  ▓▓▓▓▓  ▓▓▓▓▓  ▓   ║`,
    `╚═══════════════════════════════════╝`,
  ]

  export const CancelledError = NamedError.create("UICancelledError", z.void())

  // Toolkit-CLI.com inspired color palette
  export const Style = {
    // Reset
    RESET: "\x1b[0m",

    // Matrix/Cyberpunk theme colors (256-color mode for better gradients)
    MATRIX_GREEN: "\x1b[38;2;0;255;65m", // Bright matrix green
    MATRIX_GREEN_DIM: "\x1b[38;2;0;180;50m", // Dimmer matrix green
    CLAUDE_BLUE: "\x1b[38;2;102;153;255m", // Claude blue
    GEMINI_GREEN: "\x1b[38;2;26;214;137m", // Gemini green
    NEON_CYAN: "\x1b[38;2;0;255;255m", // Bright cyan
    NEON_MAGENTA: "\x1b[38;2;255;0;255m", // Bright magenta
    CYBER_PURPLE: "\x1b[38;2;147;51;234m", // Purple accent

    // Background colors
    BG_DARK: "\x1b[48;2;26;27;38m", // Dark background (#1a1b26)
    BG_TRANSPARENT: "\x1b[49m",

    // Text styles
    BOLD: "\x1b[1m",
    DIM: "\x1b[2m",
    ITALIC: "\x1b[3m",
    UNDERLINE: "\x1b[4m",
    BLINK: "\x1b[5m",

    // Original styles (maintained for compatibility)
    TEXT_HIGHLIGHT: "\x1b[96m",
    TEXT_HIGHLIGHT_BOLD: "\x1b[96m\x1b[1m",
    TEXT_DIM: "\x1b[90m",
    TEXT_DIM_BOLD: "\x1b[90m\x1b[1m",
    TEXT_NORMAL: "\x1b[0m",
    TEXT_NORMAL_BOLD: "\x1b[1m",
    TEXT_WARNING: "\x1b[93m",
    TEXT_WARNING_BOLD: "\x1b[93m\x1b[1m",
    TEXT_DANGER: "\x1b[91m",
    TEXT_DANGER_BOLD: "\x1b[91m\x1b[1m",
    TEXT_SUCCESS: "\x1b[92m",
    TEXT_SUCCESS_BOLD: "\x1b[92m\x1b[1m",
    TEXT_INFO: "\x1b[94m",
    TEXT_INFO_BOLD: "\x1b[94m\x1b[1m",
  }

  /**
   * Create gradient text effect (matrix green to cyan)
   */
  export function gradient(text: string, colors: string[] = [Style.MATRIX_GREEN, Style.NEON_CYAN]): string {
    if (text.length === 0) return ""
    if (colors.length === 0) return text
    if (colors.length === 1) return colors[0] + text + Style.RESET

    const result: string[] = []
    const step = Math.max(1, Math.floor(text.length / (colors.length - 1)))

    for (let i = 0; i < text.length; i++) {
      const colorIndex = Math.min(colors.length - 1, Math.floor(i / step))
      result.push(colors[colorIndex] + text[i])
    }
    result.push(Style.RESET)

    return result.join("")
  }

  /**
   * Create glowing text effect with matrix green
   */
  export function glow(text: string, color: string = Style.MATRIX_GREEN): string {
    return `${Style.BOLD}${color}${text}${Style.RESET}`
  }

  /**
   * Create boxed text with cyberpunk border
   */
  export function box(text: string, options: { color?: string; padding?: number } = {}): string {
    const color = options.color || Style.CLAUDE_BLUE
    const padding = options.padding ?? 1
    const pad = " ".repeat(padding)
    const width = text.length + padding * 2
    const top = `${color}╔${"═".repeat(width)}╗${Style.RESET}`
    const middle = `${color}║${pad}${Style.RESET}${text}${color}${pad}║${Style.RESET}`
    const bottom = `${color}╚${"═".repeat(width)}╝${Style.RESET}`

    return `${top}\n${middle}\n${bottom}`
  }

  /**
   * Create typing cursor effect (for static display)
   */
  export function withCursor(text: string): string {
    return `${text}${Style.MATRIX_GREEN}${Style.BLINK}▌${Style.RESET}`
  }

  export function println(...message: string[]) {
    print(...message)
    Bun.stderr.write(EOL)
  }

  export function print(...message: string[]) {
    blank = false
    Bun.stderr.write(message.join(" "))
  }

  let blank = false
  export function empty() {
    if (blank) return
    println("" + Style.TEXT_NORMAL)
    blank = true
  }

  export function logo(pad?: string, style: "classic" | "cyberpunk" | "gradient" = "classic") {
    if (style === "cyberpunk") {
      return LOGO_CYBERPUNK.map((line, i) => {
        const color = i === 0 || i === LOGO_CYBERPUNK.length - 1 ? Style.MATRIX_GREEN : Style.CLAUDE_BLUE
        return (pad || "") + color + line + Style.RESET
      }).join(EOL)
    }

    if (style === "gradient") {
      // Apply gradient effect to each line
      return LOGO_CYBERPUNK.map((line, i) => {
        const colors = [Style.MATRIX_GREEN, Style.GEMINI_GREEN, Style.CLAUDE_BLUE, Style.CYBER_PURPLE]
        const colorIndex = Math.floor((i / LOGO_CYBERPUNK.length) * (colors.length - 1))
        const color = colors[colorIndex]
        return (pad || "") + color + line + Style.RESET
      }).join(EOL)
    }

    // Classic logo (original behavior)
    const result = []
    for (const row of LOGO) {
      if (pad) result.push(pad)
      result.push(Bun.color("gray", "ansi"))
      result.push(row[0])
      result.push("\x1b[0m")
      result.push(row[1])
      result.push(EOL)
    }
    return result.join("").trimEnd()
  }

  /**
   * Display an animated banner with toolkit-cli aesthetic
   */
  export function banner(message: string) {
    const boxed = box(glow(message, Style.MATRIX_GREEN), {
      color: Style.CLAUDE_BLUE,
      padding: 2,
    })
    println(boxed)
  }

  /**
   * Display a section header with cyberpunk styling
   */
  export function section(title: string) {
    const line = Style.MATRIX_GREEN_DIM + "─".repeat(40) + Style.RESET
    println()
    println(line)
    println(glow(`▶ ${title}`, Style.CLAUDE_BLUE))
    println(line)
  }

  export async function input(prompt: string): Promise<string> {
    const readline = require("readline")
    const rl = readline.createInterface({
      input: process.stdin,
      output: process.stdout,
    })

    return new Promise((resolve) => {
      rl.question(prompt, (answer: string) => {
        rl.close()
        resolve(answer.trim())
      })
    })
  }

  export function error(message: string) {
    println(Style.TEXT_DANGER_BOLD + "Error: " + Style.TEXT_NORMAL + message)
  }

  export function markdown(text: string): string {
    return text
  }
}
