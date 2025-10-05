#!/usr/bin/env bun

/**
 * Test clickable links - matches the user's "What's Next?" scenario
 */

import { whatsNext } from "../../src/cli/whats-next"
import { UI } from "../../src/cli/ui"

console.log("\n" + "=".repeat(70))
console.log(UI.glow("🧪 TESTING CLICKABLE LINKS", UI.Style.MATRIX_GREEN))
console.log("=".repeat(70) + "\n")

console.log(UI.Style.TEXT_INFO + "Try Cmd+Click (Mac) or Ctrl+Click on the file paths below!" + UI.Style.RESET)
console.log()

// Simulate the exact "What's Next?" from the screenshot
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

console.log("=".repeat(70))
console.log(UI.Style.BOLD + "Individual Link Tests:" + UI.Style.RESET)
console.log("=".repeat(70) + "\n")

// Test 1: Simple file link
console.log(UI.Style.BOLD + "1. Simple file link:" + UI.Style.RESET)
console.log("   " + UI.fileLink("/Users/aaron/Code/Toolkit-CLI/tasks.md"))
console.log()

// Test 2: File link with custom text
console.log(UI.Style.BOLD + "2. File link with custom text:" + UI.Style.RESET)
console.log("   " + UI.fileLink("/Users/aaron/Code/RyCode/opencode/README.md", "📖 Project README"))
console.log()

// Test 3: URL link
console.log(UI.Style.BOLD + "3. URL link:" + UI.Style.RESET)
console.log("   " + UI.link("Visit OpenAI", "https://openai.com"))
console.log()

// Test 4: Styled file link
console.log(UI.Style.BOLD + "4. Styled file link with icon:" + UI.Style.RESET)
console.log("   " + UI.styledFileLink("/Users/aaron/Code/Toolkit-CLI/tasks.md", "My Tasks"))
console.log()

// Test 5: Auto-linking in text
console.log(UI.Style.BOLD + "5. Auto-linking in paragraph:" + UI.Style.RESET)
const paragraph = `
The configuration file is located at /Users/aaron/Code/RyCode/opencode/opencode.json
and the documentation is at https://github.com/anthropics/opencode. You can also
check /Users/aaron/Code/Toolkit-CLI/tasks.md for pending tasks.
`.trim()
console.log("   " + UI.autoLink(paragraph))
console.log()

// Test 6: Multiple styled links in a sentence
console.log(UI.Style.BOLD + "6. Multiple links in one line:" + UI.Style.RESET)
console.log(
  "   Read the " +
  UI.styledFileLink("/Users/aaron/Code/RyCode/opencode/README.md", "README") +
  " or visit " +
  UI.styledLink("GitHub", "https://github.com/anthropics/opencode", "🔗") +
  " for more info."
)
console.log()

console.log("=".repeat(70))
console.log(UI.Style.TEXT_SUCCESS + "✅ All links generated successfully!" + UI.Style.RESET)
console.log(UI.Style.TEXT_DIM + "Try clicking the links above (Cmd+Click on Mac, Ctrl+Click on Windows/Linux)" + UI.Style.RESET)
console.log("=".repeat(70) + "\n")
