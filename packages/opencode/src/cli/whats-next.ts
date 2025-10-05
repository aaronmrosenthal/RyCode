/**
 * "What's Next?" UI Component
 *
 * Display actionable next steps with clickable links
 */

import { UI } from "./ui"

export interface NextStepOption {
  key: string
  label: string
  time?: string
  command?: string
  description?: string
  file?: string // File path to make clickable
  url?: string // URL to make clickable
}

/**
 * Display a "What's Next?" prompt with clickable links
 */
export function whatsNext(options: NextStepOption[], recommendation?: string) {
  console.log()
  console.log(UI.glow("🎯 What's Next?", UI.Style.MATRIX_GREEN))
  console.log()

  console.log(UI.Style.TEXT_DIM + "Immediate Options:" + UI.Style.RESET)
  console.log()

  for (const opt of options) {
    // Option header: [A] Label (time)
    const header =
      UI.Style.BOLD + UI.Style.CLAUDE_BLUE + `[${opt.key}] ` + UI.Style.RESET +
      UI.Style.BOLD + opt.label + UI.Style.RESET +
      (opt.time ? UI.Style.TEXT_DIM + ` (${opt.time})` + UI.Style.RESET : "")

    console.log(header)

    // Command or description (with clickable links)
    if (opt.command) {
      console.log(opt.command)
    }

    if (opt.description) {
      // Auto-link any file paths or URLs in description
      const linkedDescription = UI.autoLink(opt.description)
      console.log(linkedDescription)
    }

    // If explicit file path provided, add as clickable link
    if (opt.file) {
      console.log("   " + UI.styledFileLink(opt.file))
    }

    // If explicit URL provided, add as clickable link
    if (opt.url) {
      console.log("   " + UI.styledLink(opt.url, opt.url))
    }

    console.log()
  }

  // Recommendation
  if (recommendation) {
    console.log(UI.Style.TEXT_DIM + "───" + UI.Style.RESET)
    console.log(UI.Style.BOLD + "Recommendation: " + UI.Style.RESET + UI.autoLink(recommendation))
    console.log()
  }

  // Footer with instructions
  const keys = options.map(o => UI.Style.BOLD + o.label.toLowerCase().split(" ")[0] + UI.Style.RESET)
  const keysText = keys.slice(0, -1).join(", ") + ", or " + keys[keys.length - 1]

  console.log(
    UI.Style.TEXT_INFO +
    "Type " + keysText + " to proceed, or ask for other suggestions." +
    UI.Style.RESET
  )
  console.log()
}

/**
 * Example usage matching the screenshot
 */
export function exampleWhatsNext() {
  whatsNext(
    [
      {
        key: "A",
        label: "Push to remote",
        time: "30 sec",
        command: "git push origin main",
        description: "Deploy your security hardening to the repository.",
      },
      {
        key: "B",
        label: "Clean up tasks.md",
        time: "10 min",
        description: "Fix the malformed tasks in /Users/aaron/Code/Toolkit-CLI/tasks.md",
      },
      {
        key: "C",
        label: "Start Quick-Win implementation",
        time: "4 weeks",
        description: "Begin T001: Create directory structure for the /quick-win command",
      },
    ],
    "Push to remote first (A), then tackle tasks.md cleanup (B), then start Quick-Win when ready for a focused sprint"
  )
}

// Run example if called directly
if (import.meta.main) {
  exampleWhatsNext()
}
