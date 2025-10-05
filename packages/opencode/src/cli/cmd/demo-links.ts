#!/usr/bin/env bun

/**
 * Demo: Clickable Hyperlinks in Terminal
 *
 * Shows how to create clickable links using OSC 8 escape sequences
 * Works in: iTerm2, Alacritty, Kitty, WezTerm, Windows Terminal, VSCode Terminal
 */

import { UI } from "../ui"

console.log("\n" + UI.glow("🔗 Clickable Hyperlinks Demo", UI.Style.MATRIX_GREEN))
console.log(UI.Style.TEXT_DIM + "═".repeat(60) + UI.Style.RESET + "\n")

// Example 1: Simple URL link
console.log(UI.Style.BOLD + "1. URL Links:" + UI.Style.RESET)
console.log("   " + UI.link("Visit OpenAI", "https://openai.com"))
console.log("   " + UI.link("Claude AI", "https://claude.ai"))
console.log()

// Example 2: File path links
console.log(UI.Style.BOLD + "2. File Path Links:" + UI.Style.RESET)
console.log("   " + UI.fileLink("/Users/aaron/Code/Toolkit-CLI/tasks.md"))
console.log("   " + UI.fileLink("/Users/aaron/Code/RyCode/opencode/README.md", "Project README"))
console.log()

// Example 3: Styled links with icons
console.log(UI.Style.BOLD + "3. Styled Links:" + UI.Style.RESET)
console.log("   " + UI.styledLink("Documentation", "https://docs.example.com", "📚"))
console.log("   " + UI.styledFileLink("/Users/aaron/Code/Toolkit-CLI/tasks.md"))
console.log()

// Example 4: Auto-detect and link paths/URLs in text
console.log(UI.Style.BOLD + "4. Auto-Linking:" + UI.Style.RESET)
const text = "Check the file at /Users/aaron/Code/Toolkit-CLI/tasks.md or visit https://github.com"
console.log("   " + UI.autoLink(text))
console.log()

// Example 5: "What's Next?" style with clickable options
console.log(UI.Style.BOLD + "5. What's Next? (Interactive Example):" + UI.Style.RESET)
console.log()

const options = [
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
    description: "Fix the malformed tasks in " + UI.fileLink("/Users/aaron/Code/Toolkit-CLI/tasks.md"),
  },
  {
    key: "C",
    label: "Start Quick-Win implementation",
    time: "4 weeks",
    description: "Begin T001: Create directory structure for the " + UI.fileLink("/quick-win", "file:///Users/aaron/Code/Toolkit-CLI/quick-win") + " command",
  },
]

for (const opt of options) {
  console.log(
    UI.Style.BOLD + UI.Style.CLAUDE_BLUE + `[${opt.key}] ` + UI.Style.RESET +
    UI.Style.BOLD + opt.label + UI.Style.RESET +
    UI.Style.TEXT_DIM + ` (${opt.time})` + UI.Style.RESET
  )
  console.log(opt.command ? opt.command : opt.description)
  console.log()
}

console.log(UI.Style.TEXT_DIM + "─".repeat(60) + UI.Style.RESET)
console.log(UI.Style.BOLD + "Recommendation: " + UI.Style.RESET + "Push to remote first (A), then tackle tasks.md cleanup (B), then start Quick-Win when ready for a focused sprint")
console.log()
console.log(UI.Style.TEXT_INFO + "Type " + UI.Style.BOLD + "push" + UI.Style.RESET + UI.Style.TEXT_INFO + ", " + UI.Style.BOLD + "clean" + UI.Style.RESET + UI.Style.TEXT_INFO + ", or " + UI.Style.BOLD + "quick-win" + UI.Style.RESET + UI.Style.TEXT_INFO + " to proceed, or ask for other suggestions." + UI.Style.RESET)
console.log()

// Example 6: Multiple links in a sentence
console.log(UI.Style.BOLD + "6. Links in Context:" + UI.Style.RESET)
console.log(
  "   Read the " + UI.styledFileLink("/Users/aaron/Code/RyCode/opencode/README.md", "README") +
  " or check " + UI.styledLink("issues", "https://github.com/anthropics/opencode/issues", "🐛")
)
console.log()

// Terminal compatibility note
console.log(UI.Style.TEXT_DIM + "Note: Clickable links work in modern terminals (iTerm2, Alacritty, Kitty, WezTerm, etc.)")
console.log("      Cmd+Click (Mac) or Ctrl+Click (Linux/Windows) to open links" + UI.Style.RESET)
console.log()
