#!/usr/bin/env bun

/**
 * Verify clickable links are working
 * Shows both raw escape sequences and rendered output
 */

import { UI } from "./src/cli/ui"

console.log("\n" + UI.box(UI.glow("Clickable Links Verification", UI.Style.CLAUDE_BLUE)))
console.log()

// Test case 1: File link
const filePath = "/Users/aaron/Code/Toolkit-CLI/tasks.md"
const fileLink = UI.fileLink(filePath)

console.log(UI.Style.BOLD + "Test 1: File Link" + UI.Style.RESET)
console.log(UI.Style.TEXT_DIM + "Path: " + filePath + UI.Style.RESET)
console.log(UI.Style.TEXT_DIM + "Raw: " + JSON.stringify(fileLink) + UI.Style.RESET)
console.log(UI.Style.TEXT_HIGHLIGHT + "Rendered: " + fileLink + UI.Style.RESET)
console.log(UI.Style.TEXT_DIM + "Expected: Clickable link that opens file in editor" + UI.Style.RESET)
console.log()

// Test case 2: URL link
const url = "https://github.com/anthropics/opencode"
const urlLink = UI.link("GitHub Repo", url)

console.log(UI.Style.BOLD + "Test 2: URL Link" + UI.Style.RESET)
console.log(UI.Style.TEXT_DIM + "URL: " + url + UI.Style.RESET)
console.log(UI.Style.TEXT_DIM + "Raw: " + JSON.stringify(urlLink) + UI.Style.RESET)
console.log(UI.Style.TEXT_HIGHLIGHT + "Rendered: " + urlLink + UI.Style.RESET)
console.log(UI.Style.TEXT_DIM + "Expected: Clickable link that opens in browser" + UI.Style.RESET)
console.log()

// Test case 3: Styled file link
const styledLink = UI.styledFileLink("/Users/aaron/Code/RyCode/opencode/README.md", "Project README")

console.log(UI.Style.BOLD + "Test 3: Styled File Link" + UI.Style.RESET)
console.log(UI.Style.TEXT_DIM + "Raw: " + JSON.stringify(styledLink) + UI.Style.RESET)
console.log(UI.Style.TEXT_HIGHLIGHT + "Rendered: " + styledLink + UI.Style.RESET)
console.log(UI.Style.TEXT_DIM + "Expected: Styled clickable link with file icon" + UI.Style.RESET)
console.log()

// Test case 4: Auto-linking
const textWithPaths = "Check /Users/aaron/tasks.md and https://example.com"
const autoLinked = UI.autoLink(textWithPaths)

console.log(UI.Style.BOLD + "Test 4: Auto-Linking" + UI.Style.RESET)
console.log(UI.Style.TEXT_DIM + "Input: " + textWithPaths + UI.Style.RESET)
console.log(UI.Style.TEXT_HIGHLIGHT + "Output: " + autoLinked + UI.Style.RESET)
console.log(UI.Style.TEXT_DIM + "Expected: Both path and URL are clickable" + UI.Style.RESET)
console.log()

// Visual test
console.log("=".repeat(70))
console.log(UI.Style.BOLD + UI.Style.CLAUDE_BLUE + "Visual Test - Try Clicking These:" + UI.Style.RESET)
console.log("=".repeat(70))
console.log()

console.log("📁 Files:")
console.log("  • " + UI.fileLink("/Users/aaron/Code/Toolkit-CLI/tasks.md"))
console.log("  • " + UI.styledFileLink("/Users/aaron/Code/RyCode/opencode/README.md"))
console.log("  • " + UI.styledFileLink("/Users/aaron/Code/RyCode/opencode/CLICKABLE_LINKS.md", "Link Documentation"))
console.log()

console.log("🌐 URLs:")
console.log("  • " + UI.link("OpenAI", "https://openai.com"))
console.log("  • " + UI.link("Claude", "https://claude.ai"))
console.log("  • " + UI.styledLink("GitHub", "https://github.com/anthropics/opencode", "⭐"))
console.log()

console.log("✅ " + UI.Style.TEXT_SUCCESS + "If you can Cmd+Click (Mac) or Ctrl+Click the links above, it's working!" + UI.Style.RESET)
console.log()

// Terminal compatibility check
console.log("=".repeat(70))
console.log(UI.Style.BOLD + "Terminal Compatibility" + UI.Style.RESET)
console.log("=".repeat(70))
console.log()

const termProgram = process.env["TERM_PROGRAM"] || "unknown"
const term = process.env["TERM"] || "unknown"

console.log(UI.Style.TEXT_INFO + "TERM_PROGRAM: " + termProgram + UI.Style.RESET)
console.log(UI.Style.TEXT_INFO + "TERM: " + term + UI.Style.RESET)
console.log()

const supported = [
  "iTerm.app",
  "Apple_Terminal",
  "WezTerm",
  "Alacritty",
  "kitty",
  "vscode"
].includes(termProgram) || term.includes("kitty")

if (supported) {
  console.log(UI.Style.TEXT_SUCCESS + "✅ Your terminal supports clickable links!" + UI.Style.RESET)
} else {
  console.log(UI.Style.TEXT_WARNING + "⚠️  Your terminal may not support clickable links." + UI.Style.RESET)
  console.log(UI.Style.TEXT_DIM + "   Supported terminals: iTerm2, Alacritty, Kitty, WezTerm, Windows Terminal, VSCode" + UI.Style.RESET)
}
console.log()
