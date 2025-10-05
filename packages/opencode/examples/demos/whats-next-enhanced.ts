#!/usr/bin/env bun

/**
 * Enhanced "What's Next?" Demo
 * Shows how to use the new TUI components in a practical workflow
 */

import { UI } from "../../src/cli/ui"
import { EnhancedTUI } from "../../src/cli/tui-enhanced"
import { CyberpunkPrompts } from "../../src/cli/theme"

async function enhancedWhatsNext() {
  UI.empty()

  // Opening with gradient logo
  UI.println(UI.logo(undefined, "gradient"))
  UI.println()

  // Main heading with glow effect
  UI.println(
    UI.box(
      UI.glow("🎯 What's Next? - Enhanced Edition", UI.Style.MATRIX_GREEN),
      { color: UI.Style.CLAUDE_BLUE, padding: 1 }
    )
  )
  UI.println()

  // System status overview
  CyberpunkPrompts.divider("System Status")
  UI.println()

  const statusInfo = {
    "Project": UI.Style.MATRIX_GREEN + "opencode" + UI.Style.RESET,
    "Branch": UI.Style.CLAUDE_BLUE + "feature/tui-enhancements" + UI.Style.RESET,
    "Status": EnhancedTUI.status("All Systems", "online", { showLabel: false }) + " Online",
    "Tests": UI.Style.MATRIX_GREEN + "✓ 142 passing" + UI.Style.RESET,
    "Coverage": EnhancedTUI.progressBar(94, 100, { width: 20, showPercentage: true }),
  }

  UI.println(EnhancedTUI.keyValue(statusInfo, { indent: 2 }))
  UI.println()

  // Recent activity timeline
  CyberpunkPrompts.divider("Recent Activity")
  UI.println()

  UI.println(
    EnhancedTUI.timeline([
      {
        time: "Just now",
        title: "Added Enhanced TUI Components",
        description: "13 new components: progress bars, spinners, tables, and more",
        status: "success",
      },
      {
        time: "2 min ago",
        title: "All Tests Passing",
        description: "142 tests completed successfully",
        status: "success",
      },
      {
        time: "5 min ago",
        title: "Documentation Updated",
        description: "Added comprehensive guide in TUI_ENHANCEMENTS.md",
        status: "info",
      },
    ])
  )
  UI.println()

  // Immediate action items in a fancy table
  CyberpunkPrompts.divider("Immediate Options")
  UI.println()

  const optionsTable = [
    [
      UI.Style.CLAUDE_BLUE + "[A]" + UI.Style.RESET,
      UI.Style.BOLD + "Commit Changes" + UI.Style.RESET,
      UI.Style.TEXT_DIM + "30 sec" + UI.Style.RESET,
      EnhancedTUI.badge("QUICK", { color: UI.Style.MATRIX_GREEN }),
    ],
    [
      UI.Style.CLAUDE_BLUE + "[B]" + UI.Style.RESET,
      UI.Style.BOLD + "Run Full Test Suite" + UI.Style.RESET,
      UI.Style.TEXT_DIM + "2 min" + UI.Style.RESET,
      EnhancedTUI.badge("VERIFY", { color: UI.Style.CLAUDE_BLUE }),
    ],
    [
      UI.Style.CLAUDE_BLUE + "[C]" + UI.Style.RESET,
      UI.Style.BOLD + "Create Pull Request" + UI.Style.RESET,
      UI.Style.TEXT_DIM + "5 min" + UI.Style.RESET,
      EnhancedTUI.badge("DEPLOY", { color: UI.Style.CYBER_PURPLE }),
    ],
  ]

  UI.println(
    EnhancedTUI.table(["Key", "Action", "Time", "Type"], optionsTable, {
      headerColor: UI.Style.MATRIX_GREEN,
      borderColor: UI.Style.MATRIX_GREEN_DIM,
      align: ["center", "left", "right", "center"],
    })
  )
  UI.println()

  // Detailed descriptions with clickable links
  UI.println(UI.Style.BOLD + "Option Details:" + UI.Style.RESET)
  UI.println()

  UI.println(
    EnhancedTUI.collapsible(
      "A: Commit Changes",
      `Stage and commit all changes to version control\n` +
        `  Modified: ${UI.fileLink("/Users/aaron/Code/RyCode/opencode/packages/opencode/src/cli/tui-enhanced.ts")}\n` +
        `  Created: ${UI.fileLink("/Users/aaron/Code/RyCode/opencode/packages/opencode/TUI_ENHANCEMENTS.md")}`,
      true,
      { titleColor: UI.Style.CLAUDE_BLUE }
    )
  )
  UI.println()

  UI.println(
    EnhancedTUI.collapsible(
      "B: Run Full Test Suite",
      `Execute complete test suite including:\n` +
        `  • Unit tests (142 tests)\n` +
        `  • Integration tests (23 tests)\n` +
        `  • Coverage analysis`,
      true,
      { titleColor: UI.Style.CLAUDE_BLUE }
    )
  )
  UI.println()

  UI.println(
    EnhancedTUI.collapsible(
      "C: Create Pull Request",
      `Create PR with title: "feat: Add Enhanced TUI Components"\n` +
        `  Target: main branch\n` +
        `  Review: Required\n` +
        `  Link: ${UI.link("GitHub PR", "https://github.com/anthropics/opencode/pulls")}`,
      true,
      { titleColor: UI.Style.CLAUDE_BLUE }
    )
  )
  UI.println()

  // Smart recommendation with wizard-style progress
  CyberpunkPrompts.divider("Recommended Workflow")
  UI.println()

  UI.println(UI.Style.BOLD + "Suggested Steps:" + UI.Style.RESET)
  UI.println()

  const recommendedSteps = [
    "Commit changes to git",
    "Run test suite to verify",
    "Create pull request",
    "Request code review",
  ]

  UI.println(EnhancedTUI.wizardSteps(recommendedSteps, 0))
  UI.println()

  // Additional context
  UI.println(
    EnhancedTUI.notification(
      "All changes are ready to commit. No conflicts detected.",
      "success",
      { title: "Git Status" }
    )
  )
  UI.println()

  // Interactive menu of actions
  CyberpunkPrompts.divider("Quick Actions")
  UI.println()

  UI.println(
    EnhancedTUI.menu([
      {
        label: "Commit & Push",
        description: "Stage all changes, commit, and push to remote",
        selected: true,
      },
      {
        label: "Run Demo",
        description: "Execute demo-tui-enhanced.ts to see components",
      },
      {
        label: "Review Documentation",
        description: "Open TUI_ENHANCEMENTS.md for full guide",
      },
      {
        label: "Run Benchmarks",
        description: "Performance testing (optional)",
      },
    ])
  )
  UI.println()

  // Footer with file structure
  CyberpunkPrompts.divider("Project Structure")
  UI.println()

  const fileStructure = {
    "src/cli": {
      "tui-enhanced.ts": null,
      "ui.ts": null,
      "theme.ts": null,
      "whats-next.ts": null,
    },
    "docs": {
      "TUI_ENHANCEMENTS.md": null,
      "CLICKABLE_LINKS.md": null,
    },
    "demos": {
      "demo-tui-enhanced.ts": null,
      "whats-next-enhanced.ts": null,
    },
  }

  UI.println(EnhancedTUI.fileTree(fileStructure))
  UI.println()

  // Final prompt
  UI.println()
  UI.println(
    UI.box(
      UI.gradient("🚀 Type a, b, c, or ask for other suggestions 🚀", [
        UI.Style.MATRIX_GREEN,
        UI.Style.CLAUDE_BLUE,
      ]),
      { color: UI.Style.NEON_CYAN, padding: 1 }
    )
  )
  UI.println()

  // Stats summary
  const stats = {
    "Components": "13",
    "Demo Lines": "300+",
    "Documentation": "500+ lines",
    "Theme": "Cyberpunk/Matrix",
  }

  UI.println(UI.Style.TEXT_DIM + "Session Stats:" + UI.Style.RESET)
  UI.println(EnhancedTUI.keyValue(stats, { indent: 2, keyColor: UI.Style.TEXT_DIM }))
  UI.println()
}

enhancedWhatsNext().catch(console.error)
