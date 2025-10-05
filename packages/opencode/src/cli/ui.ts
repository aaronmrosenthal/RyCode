import z from "zod/v4"
import { EOL } from "os"
import { NamedError } from "../util/error"

export namespace UI {
  // RyCode Logo - Classic (Original OpenCode style for backwards compatibility)
  const LOGO = [
    [`█▀▀█ █▀▀█ █▀▀ █▀▀▄ `, `█▀▀ █▀▀█ █▀▀▄ █▀▀`],
    [`█░░█ █░░█ █▀▀ █░░█ `, `█░░ █░░█ █░░█ █▀▀`],
    [`▀▀▀▀ █▀▀▀ ▀▀▀ ▀  ▀ `, `▀▀▀ ▀▀▀▀ ▀▀▀  ▀▀▀`],
  ]

  // RyCode Logo - Modern & Sleek (RECOMMENDED)
  const LOGO_RYCODE_MODERN = [
    ``,
    `  ██████╗ ██╗   ██╗ ██████╗ ██████╗ ██████╗ ███████╗`,
    `  ██╔══██╗╚██╗ ██╔╝██╔════╝██╔═══██╗██╔══██╗██╔════╝`,
    `  ██████╔╝ ╚████╔╝ ██║     ██║   ██║██║  ██║█████╗  `,
    `  ██╔══██╗  ╚██╔╝  ██║     ██║   ██║██║  ██║██╔══╝  `,
    `  ██║  ██║   ██║   ╚██████╗╚██████╔╝██████╔╝███████╗`,
    `  ╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝`,
    ``,
  ]

  // RyCode Logo - Cyberpunk/Matrix Boxed
  const LOGO_CYBERPUNK = [
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

  // RyCode Logo - Big & Bold
  const LOGO_RYCODE_BIG = [
    ``,
    `  ██████╗  ██╗   ██╗  ██████╗  ██████╗  ██████╗  ███████╗`,
    `  ██╔══██╗ ╚██╗ ██╔╝ ██╔════╝ ██╔═══██╗ ██╔══██╗ ██╔════╝`,
    `  ██████╔╝  ╚████╔╝  ██║      ██║   ██║ ██║  ██║ █████╗  `,
    `  ██╔══██╗   ╚██╔╝   ██║      ██║   ██║ ██║  ██║ ██╔══╝  `,
    `  ██║  ██║    ██║    ╚██████╗ ╚██████╔╝ ██████╔╝ ███████╗`,
    `  ╚═╝  ╚═╝    ╚═╝     ╚═════╝  ╚═════╝  ╚═════╝  ╚══════╝`,
    ``,
  ]

  // RyCode Logo - Slant Style
  const LOGO_RYCODE_SLANT = [
    ``,
    `    ____                 ______          __   `,
    `   / __ \\__  __  _____  / ____/___  ____/ /__ `,
    `  / /_/ / / / / / ___/ / /   / __ \\/ __  / _ \\`,
    ` / _, _/ /_/ / / /__  / /___/ /_/ / /_/ /  __/`,
    `/_/ |_|\\__, /  \\___/  \\____/\\____/\\__,_/\\___/ `,
    `      /____/                                   `,
    ``,
  ]

  export const CancelledError = NamedError.create("UICancelledError", z.void())

  // Toolkit-CLI.com inspired color palette
  export const Style = {
    // Reset
    RESET: "\x1b[0m",

    // Matrix/Cyberpunk theme colors (256-color mode for better gradients)
    MATRIX_GREEN: "\x1b[38;2;0;255;65m", // Bright matrix green
    MATRIX_GREEN_DIM: "\x1b[38;2;0;180;50m", // Dimmer matrix green
    MATRIX_GREEN_DARK: "\x1b[38;2;0;100;30m", // Dark matrix green
    MATRIX_GREEN_LIGHT: "\x1b[38;2;100;255;150m", // Light matrix green
    MATRIX_GREEN_BRIGHT: "\x1b[38;2;150;255;180m", // Brightest matrix green
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

  export function logo(
    pad?: string,
    style: "classic" | "cyberpunk" | "gradient" | "rycode" | "rycode-modern" | "rycode-big" | "rycode-slant" | "matrix" = "rycode"
  ) {
    // RyCode Modern with Matrix Gradient (Default - recommended)
    if (style === "rycode" || style === "rycode-modern" || style === "matrix") {
      return LOGO_RYCODE_MODERN.map((line, i) => {
        // Matrix digital rain gradient: dark green → bright green → light green → bright green
        const colors = [
          Style.MATRIX_GREEN_DARK,    // Start dark
          Style.MATRIX_GREEN_DIM,     // Dim
          Style.MATRIX_GREEN,         // Bright core
          Style.MATRIX_GREEN_BRIGHT,  // Brightest
          Style.MATRIX_GREEN_LIGHT,   // Light
          Style.MATRIX_GREEN,         // Bright
          Style.MATRIX_GREEN_DIM,     // Fade
        ]
        const colorIndex = Math.floor((i / LOGO_RYCODE_MODERN.length) * (colors.length - 1))
        const color = colors[colorIndex]
        return (pad || "") + color + line + Style.RESET
      }).join(EOL)
    }

    // RyCode Big & Bold with Matrix Gradient
    if (style === "rycode-big") {
      return LOGO_RYCODE_BIG.map((line, i) => {
        // Same matrix gradient effect
        const colors = [
          Style.MATRIX_GREEN_DARK,
          Style.MATRIX_GREEN_DIM,
          Style.MATRIX_GREEN,
          Style.MATRIX_GREEN_BRIGHT,
          Style.MATRIX_GREEN_LIGHT,
          Style.MATRIX_GREEN,
          Style.MATRIX_GREEN_DIM,
        ]
        const colorIndex = Math.floor((i / LOGO_RYCODE_BIG.length) * (colors.length - 1))
        const color = colors[colorIndex]
        return (pad || "") + color + line + Style.RESET
      }).join(EOL)
    }

    // RyCode Slant Style
    if (style === "rycode-slant") {
      return LOGO_RYCODE_SLANT.map((line, i) => {
        const colors = [Style.CYBER_PURPLE, Style.NEON_MAGENTA, Style.CYBER_PURPLE, Style.NEON_MAGENTA]
        const colorIndex = Math.floor((i / LOGO_RYCODE_SLANT.length) * (colors.length - 1))
        const color = colors[colorIndex]
        return (pad || "") + color + line + Style.RESET
      }).join(EOL)
    }

    // Cyberpunk Boxed (RyCode Cyberpunk)
    if (style === "cyberpunk") {
      return LOGO_CYBERPUNK.map((line, i) => {
        const color = i === 0 || i === LOGO_CYBERPUNK.length - 1 ? Style.MATRIX_GREEN : Style.CLAUDE_BLUE
        return (pad || "") + color + line + Style.RESET
      }).join(EOL)
    }

    // Gradient (uses cyberpunk logo)
    if (style === "gradient") {
      return LOGO_CYBERPUNK.map((line, i) => {
        const colors = [Style.MATRIX_GREEN, Style.GEMINI_GREEN, Style.CLAUDE_BLUE, Style.CYBER_PURPLE]
        const colorIndex = Math.floor((i / LOGO_CYBERPUNK.length) * (colors.length - 1))
        const color = colors[colorIndex]
        return (pad || "") + color + line + Style.RESET
      }).join(EOL)
    }

    // Classic OpenCode logo (original behavior - for backwards compatibility)
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

  /**
   * Create clickable hyperlink using OSC 8 escape sequences
   * Supported in: iTerm2, Alacritty, Kitty, WezTerm, Windows Terminal, VSCode
   *
   * @param text - Display text
   * @param url - URL or file path to link to
   * @returns Clickable link (falls back to plain text in unsupported terminals)
   *
   * @example
   * // URL link
   * link("Click here", "https://example.com")
   *
   * // File path link
   * link("/path/to/file.md", "file:///path/to/file.md")
   *
   * // Auto-detect file paths
   * fileLink("/Users/aaron/tasks.md")
   */
  export function link(text: string, url: string): string {
    // OSC 8 format: \x1b]8;;URL\x1b\\TEXT\x1b]8;;\x1b\\
    const OSC = "\x1b]"
    const ST = "\x1b\\"

    return `${OSC}8;;${url}${ST}${text}${OSC}8;;${ST}`
  }

  /**
   * Create clickable file path link
   * Automatically converts file paths to file:// URLs
   *
   * @param path - Absolute or relative file path
   * @param displayText - Optional display text (uses path if not provided)
   */
  export function fileLink(path: string, displayText?: string): string {
    // Convert to file:// URL if not already
    const url = path.startsWith("file://")
      ? path
      : "file://" + (path.startsWith("/") ? path : process.cwd() + "/" + path)

    return link(displayText || path, url)
  }

  /**
   * Auto-detect and make file paths/URLs clickable in text
   * Finds patterns like /path/to/file or https://example.com
   *
   * @param text - Text potentially containing paths or URLs
   */
  export function autoLink(text: string): string {
    // Match file paths (absolute Unix/Windows paths)
    text = text.replace(
      /(\/?[a-zA-Z]:)?\/[\w\-./]+/g,
      (match) => fileLink(match)
    )

    // Match URLs
    text = text.replace(
      /(https?:\/\/[^\s]+)/g,
      (match) => link(match, match)
    )

    return text
  }

  /**
   * Create styled clickable link with icon
   */
  export function styledLink(text: string, url: string, icon: string = "🔗"): string {
    return `${icon} ${Style.UNDERLINE}${Style.CLAUDE_BLUE}${link(text, url)}${Style.RESET}`
  }

  /**
   * Create a file path link with file icon
   */
  export function styledFileLink(path: string, displayText?: string): string {
    return styledLink(displayText || path, "file://" + path, "📄")
  }
}
