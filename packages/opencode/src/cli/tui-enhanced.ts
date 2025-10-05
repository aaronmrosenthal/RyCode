/**
 * Enhanced TUI Components
 * Advanced terminal UI components with animations, progress bars, tables, and more
 */

import { UI } from "./ui"

export namespace EnhancedTUI {
  /**
   * Progress bar with percentage and animation
   */
  export function progressBar(
    current: number,
    total: number,
    options: {
      width?: number
      label?: string
      showPercentage?: boolean
      color?: string
      bgColor?: string
    } = {}
  ): string {
    const width = options.width || 40
    const label = options.label || ""
    const showPercentage = options.showPercentage ?? true
    const color = options.color || UI.Style.MATRIX_GREEN
    const bgColor = options.bgColor || UI.Style.TEXT_DIM

    const percentage = Math.min(100, Math.max(0, Math.round((current / total) * 100)))
    const filledWidth = Math.round((width * percentage) / 100)
    const emptyWidth = width - filledWidth

    const filled = "█".repeat(filledWidth)
    const empty = "░".repeat(emptyWidth)

    const bar = `${color}${filled}${UI.Style.RESET}${bgColor}${empty}${UI.Style.RESET}`
    const percentText = showPercentage ? ` ${percentage}%` : ""
    const labelText = label ? `${label} ` : ""

    return `${labelText}[${bar}]${percentText}`
  }

  /**
   * Animated spinner frames
   */
  export const spinnerFrames = {
    dots: ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"],
    line: ["-", "\\", "|", "/"],
    dots2: ["⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"],
    matrix: ["⡀", "⠄", "⠂", "⠁", "⠈", "⠐", "⠠", "⢀"],
    pulse: ["○", "◔", "◑", "◕", "●", "◕", "◑", "◔"],
    neon: ["◜", "◠", "◝", "◞", "◡", "◟"],
    cyber: ["▁", "▂", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃", "▂"],
    arrow: ["←", "↖", "↑", "↗", "→", "↘", "↓", "↙"],
  }

  /**
   * Create a spinner instance
   */
  export function spinner(
    message: string,
    options: {
      frames?: string[]
      color?: string
      interval?: number
    } = {}
  ) {
    const frames = options.frames || spinnerFrames.dots
    const color = options.color || UI.Style.NEON_CYAN
    let frame = 0
    let intervalId: Timer | null = null

    return {
      start() {
        intervalId = setInterval(() => {
          const spinChar = frames[frame % frames.length]
          process.stderr.write(`\r${color}${spinChar}${UI.Style.RESET} ${message}`)
          frame++
        }, options.interval || 80)
      },
      update(newMessage: string) {
        message = newMessage
      },
      stop(finalMessage?: string) {
        if (intervalId) clearInterval(intervalId)
        process.stderr.write(`\r${UI.Style.MATRIX_GREEN}✓${UI.Style.RESET} ${finalMessage || message}\n`)
      },
      fail(errorMessage?: string) {
        if (intervalId) clearInterval(intervalId)
        process.stderr.write(`\r${UI.Style.TEXT_DANGER}✖${UI.Style.RESET} ${errorMessage || message}\n`)
      },
    }
  }

  /**
   * Table renderer with styled columns
   */
  export function table(
    headers: string[],
    rows: string[][],
    options: {
      headerColor?: string
      borderColor?: string
      align?: ("left" | "center" | "right")[]
      maxWidth?: number
    } = {}
  ): string {
    const headerColor = options.headerColor || UI.Style.CLAUDE_BLUE
    const borderColor = options.borderColor || UI.Style.MATRIX_GREEN_DIM
    const align = options.align || headers.map(() => "left")

    // Calculate column widths
    const widths = headers.map((header, i) => {
      const maxContentWidth = Math.max(
        header.length,
        ...rows.map((row) => (row[i] || "").length)
      )
      return options.maxWidth ? Math.min(maxContentWidth, options.maxWidth) : maxContentWidth
    })

    // Helper to align text
    const alignText = (text: string, width: number, alignment: "left" | "center" | "right") => {
      const stripped = text.replace(/\x1b\[[0-9;]*m/g, "") // Remove ANSI codes for length calculation
      const padding = width - stripped.length
      if (padding <= 0) return text.slice(0, width)

      switch (alignment) {
        case "center":
          const leftPad = Math.floor(padding / 2)
          const rightPad = padding - leftPad
          return " ".repeat(leftPad) + text + " ".repeat(rightPad)
        case "right":
          return " ".repeat(padding) + text
        default:
          return text + " ".repeat(padding)
      }
    }

    const lines: string[] = []

    // Top border
    const topBorder =
      borderColor + "╔" + widths.map((w) => "═".repeat(w + 2)).join("╦") + "╗" + UI.Style.RESET
    lines.push(topBorder)

    // Header row
    const headerRow =
      borderColor + "║ " + UI.Style.RESET +
      headers
        .map((h, i) => headerColor + UI.Style.BOLD + alignText(h, widths[i], align[i]) + UI.Style.RESET)
        .join(borderColor + " ║ " + UI.Style.RESET) +
      borderColor + " ║" + UI.Style.RESET
    lines.push(headerRow)

    // Separator
    const separator =
      borderColor + "╠" + widths.map((w) => "═".repeat(w + 2)).join("╬") + "╣" + UI.Style.RESET
    lines.push(separator)

    // Data rows
    for (const row of rows) {
      const dataRow =
        borderColor + "║ " + UI.Style.RESET +
        row
          .map((cell, i) => alignText(cell || "", widths[i], align[i]))
          .join(borderColor + " ║ " + UI.Style.RESET) +
        borderColor + " ║" + UI.Style.RESET
      lines.push(dataRow)
    }

    // Bottom border
    const bottomBorder =
      borderColor + "╚" + widths.map((w) => "═".repeat(w + 2)).join("╩") + "╝" + UI.Style.RESET
    lines.push(bottomBorder)

    return lines.join("\n")
  }

  /**
   * Timeline display for events
   */
  export function timeline(
    events: Array<{
      time: string
      title: string
      description?: string
      status?: "success" | "error" | "warning" | "info"
    }>,
    options: {
      showConnector?: boolean
      colors?: Record<string, string>
    } = {}
  ): string {
    const showConnector = options.showConnector ?? true
    const statusIcons = {
      success: UI.Style.MATRIX_GREEN + "✓" + UI.Style.RESET,
      error: UI.Style.TEXT_DANGER + "✖" + UI.Style.RESET,
      warning: UI.Style.TEXT_WARNING + "⚠" + UI.Style.RESET,
      info: UI.Style.CLAUDE_BLUE + "ℹ" + UI.Style.RESET,
    }

    const lines: string[] = []

    events.forEach((event, i) => {
      const icon = event.status ? statusIcons[event.status] : UI.Style.MATRIX_GREEN + "●" + UI.Style.RESET
      const time = UI.Style.TEXT_DIM + event.time + UI.Style.RESET
      const title = UI.Style.BOLD + event.title + UI.Style.RESET

      lines.push(`  ${icon} ${time} ${title}`)

      if (event.description) {
        const connector = showConnector ? UI.Style.MATRIX_GREEN_DIM + "│" + UI.Style.RESET : " "
        lines.push(`  ${connector}   ${UI.Style.TEXT_DIM}${event.description}${UI.Style.RESET}`)
      }

      if (showConnector && i < events.length - 1) {
        lines.push(`  ${UI.Style.MATRIX_GREEN_DIM}│${UI.Style.RESET}`)
      }
    })

    return lines.join("\n")
  }

  /**
   * Notification box (toast-style)
   */
  export function notification(
    message: string,
    type: "success" | "error" | "warning" | "info" = "info",
    options: {
      title?: string
      duration?: number
    } = {}
  ): string {
    const icons = {
      success: "✓",
      error: "✖",
      warning: "⚠",
      info: "ℹ",
    }

    const colors = {
      success: UI.Style.MATRIX_GREEN,
      error: UI.Style.TEXT_DANGER,
      warning: UI.Style.TEXT_WARNING,
      info: UI.Style.CLAUDE_BLUE,
    }

    const icon = icons[type]
    const color = colors[type]
    const title = options.title || type.toUpperCase()

    const content = options.title
      ? `${icon} ${UI.Style.BOLD}${title}${UI.Style.RESET}\n${message}`
      : `${icon} ${message}`

    return UI.box(content, { color, padding: 1 })
  }

  /**
   * Multi-step wizard progress indicator
   */
  export function wizardSteps(
    steps: string[],
    currentStep: number,
    options: {
      completedColor?: string
      activeColor?: string
      pendingColor?: string
    } = {}
  ): string {
    const completedColor = options.completedColor || UI.Style.MATRIX_GREEN
    const activeColor = options.activeColor || UI.Style.CLAUDE_BLUE
    const pendingColor = options.pendingColor || UI.Style.TEXT_DIM

    const lines: string[] = []

    steps.forEach((step, i) => {
      const stepNum = i + 1
      let icon: string
      let color: string
      let connector = ""

      if (i < currentStep) {
        // Completed step
        icon = "✓"
        color = completedColor
      } else if (i === currentStep) {
        // Current step
        icon = "▶"
        color = activeColor
      } else {
        // Pending step
        icon = "○"
        color = pendingColor
      }

      const stepText = `${color}${UI.Style.BOLD}${icon}${UI.Style.RESET} ${color}Step ${stepNum}:${UI.Style.RESET} ${
        i === currentStep ? UI.Style.BOLD + color + step + UI.Style.RESET : color + step + UI.Style.RESET
      }`

      lines.push(stepText)

      // Add connector line between steps
      if (i < steps.length - 1) {
        connector = i < currentStep ? completedColor + "  │" + UI.Style.RESET : pendingColor + "  │" + UI.Style.RESET
        lines.push(connector)
      }
    })

    return lines.join("\n")
  }

  /**
   * Key-value pair display with alignment
   */
  export function keyValue(
    pairs: Record<string, string>,
    options: {
      keyColor?: string
      valueColor?: string
      separator?: string
      indent?: number
    } = {}
  ): string {
    const keyColor = options.keyColor || UI.Style.CLAUDE_BLUE
    const valueColor = options.valueColor || UI.Style.RESET
    const separator = options.separator || ":"
    const indent = " ".repeat(options.indent || 0)

    const maxKeyLength = Math.max(...Object.keys(pairs).map((k) => k.length))

    return Object.entries(pairs)
      .map(([key, value]) => {
        const paddedKey = key.padEnd(maxKeyLength)
        return `${indent}${keyColor}${paddedKey}${UI.Style.RESET} ${separator} ${valueColor}${value}${UI.Style.RESET}`
      })
      .join("\n")
  }

  /**
   * Code block with syntax highlighting (basic)
   */
  export function codeBlock(
    code: string,
    options: {
      language?: string
      showLineNumbers?: boolean
      highlightLines?: number[]
      theme?: "dark" | "matrix"
    } = {}
  ): string {
    const showLineNumbers = options.showLineNumbers ?? true
    const highlightLines = options.highlightLines || []
    const theme = options.theme || "matrix"

    const lines = code.split("\n")
    const maxLineNumWidth = lines.length.toString().length

    const themeColors = {
      dark: {
        lineNumber: UI.Style.TEXT_DIM,
        highlight: UI.Style.BG_DARK,
        code: UI.Style.RESET,
      },
      matrix: {
        lineNumber: UI.Style.MATRIX_GREEN_DIM,
        highlight: UI.Style.MATRIX_GREEN,
        code: UI.Style.CLAUDE_BLUE,
      },
    }

    const colors = themeColors[theme]

    return lines
      .map((line, i) => {
        const lineNum = (i + 1).toString().padStart(maxLineNumWidth, " ")
        const lineNumText = showLineNumbers ? `${colors.lineNumber}${lineNum}${UI.Style.RESET} │ ` : ""
        const highlight = highlightLines.includes(i + 1) ? colors.highlight : ""
        return `${lineNumText}${highlight}${colors.code}${line}${UI.Style.RESET}`
      })
      .join("\n")
  }

  /**
   * Animated typing effect (returns array of frames)
   */
  export function* typingAnimation(text: string, options: { cursor?: string } = {}): Generator<string> {
    const cursor = options.cursor || UI.Style.MATRIX_GREEN + "▌" + UI.Style.RESET
    for (let i = 0; i <= text.length; i++) {
      yield text.slice(0, i) + cursor
    }
    yield text
  }

  /**
   * Diff viewer for git-style diffs
   */
  export function diff(
    additions: string[],
    deletions: string[],
    options: {
      addColor?: string
      deleteColor?: string
      context?: number
    } = {}
  ): string {
    const addColor = options.addColor || UI.Style.MATRIX_GREEN
    const deleteColor = options.deleteColor || UI.Style.TEXT_DANGER

    const lines: string[] = []

    deletions.forEach((line) => {
      lines.push(`${deleteColor}- ${line}${UI.Style.RESET}`)
    })

    additions.forEach((line) => {
      lines.push(`${addColor}+ ${line}${UI.Style.RESET}`)
    })

    return lines.join("\n")
  }

  /**
   * Badge/tag component
   */
  export function badge(
    text: string,
    options: {
      color?: string
      variant?: "filled" | "outlined"
    } = {}
  ): string {
    const color = options.color || UI.Style.CLAUDE_BLUE
    const variant = options.variant || "filled"

    if (variant === "outlined") {
      return `${color}[${text}]${UI.Style.RESET}`
    }

    return `${color}${UI.Style.BOLD} ${text} ${UI.Style.RESET}`
  }

  /**
   * Collapsible section (returns formatted text, actual collapse needs interaction)
   */
  export function collapsible(
    title: string,
    content: string,
    isExpanded: boolean = false,
    options: {
      expandIcon?: string
      collapseIcon?: string
      titleColor?: string
    } = {}
  ): string {
    const expandIcon = options.expandIcon || "▶"
    const collapseIcon = options.collapseIcon || "▼"
    const titleColor = options.titleColor || UI.Style.CLAUDE_BLUE

    const icon = isExpanded ? collapseIcon : expandIcon
    const titleText = `${titleColor}${icon} ${UI.Style.BOLD}${title}${UI.Style.RESET}`

    if (isExpanded) {
      const indentedContent = content
        .split("\n")
        .map((line) => `  ${line}`)
        .join("\n")
      return `${titleText}\n${indentedContent}`
    }

    return titleText
  }

  /**
   * Interactive menu selector (static representation)
   */
  export function menu(
    options: Array<{ label: string; description?: string; selected?: boolean; disabled?: boolean }>,
    opts: {
      selectedColor?: string
      disabledColor?: string
    } = {}
  ): string {
    const selectedColor = opts.selectedColor || UI.Style.MATRIX_GREEN
    const disabledColor = opts.disabledColor || UI.Style.TEXT_DIM

    return options
      .map((option) => {
        const icon = option.selected ? "◉" : "○"
        const color = option.disabled ? disabledColor : option.selected ? selectedColor : UI.Style.RESET
        const label = option.disabled ? `${option.label} (disabled)` : option.label
        const description = option.description ? `\n  ${UI.Style.TEXT_DIM}${option.description}${UI.Style.RESET}` : ""

        return `${color}${icon} ${label}${UI.Style.RESET}${description}`
      })
      .join("\n")
  }

  /**
   * Status indicator with pulse effect (static representation)
   */
  export function status(
    label: string,
    state: "online" | "offline" | "loading" | "error",
    options: {
      showLabel?: boolean
    } = {}
  ): string {
    const showLabel = options.showLabel ?? true

    const states = {
      online: { icon: "●", color: UI.Style.MATRIX_GREEN },
      offline: { icon: "●", color: UI.Style.TEXT_DIM },
      loading: { icon: "◔", color: UI.Style.NEON_CYAN },
      error: { icon: "●", color: UI.Style.TEXT_DANGER },
    }

    const { icon, color } = states[state]
    const labelText = showLabel ? ` ${label}` : ""

    return `${color}${icon}${UI.Style.RESET}${labelText}`
  }

  /**
   * File tree display
   */
  export function fileTree(
    tree: Record<string, any>,
    options: {
      indent?: number
      level?: number
      prefix?: string
      isLast?: boolean
    } = {}
  ): string {
    const indent = options.indent || 0
    const level = options.level || 0
    const prefix = options.prefix || ""
    const lines: string[] = []

    const entries = Object.entries(tree)
    entries.forEach(([name, value], index) => {
      const isLast = index === entries.length - 1
      const connector = isLast ? "└─" : "├─"
      const icon = typeof value === "object" ? "📁" : "📄"
      const color = typeof value === "object" ? UI.Style.CLAUDE_BLUE : UI.Style.TEXT_DIM

      lines.push(`${prefix}${UI.Style.MATRIX_GREEN_DIM}${connector}${UI.Style.RESET} ${icon} ${color}${name}${UI.Style.RESET}`)

      if (typeof value === "object" && value !== null) {
        const newPrefix = prefix + (isLast ? "   " : UI.Style.MATRIX_GREEN_DIM + "│  " + UI.Style.RESET)
        lines.push(fileTree(value, { ...options, prefix: newPrefix, level: level + 1 }))
      }
    })

    return lines.join("\n")
  }
}
