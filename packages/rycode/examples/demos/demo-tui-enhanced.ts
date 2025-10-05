#!/usr/bin/env bun

/**
 * Enhanced TUI Demo
 * Showcases all the new verbose and interactive TUI components
 */

import { UI } from "../../src/cli/ui"
import { EnhancedTUI } from "../../src/cli/tui-enhanced"
import { CyberpunkPrompts } from "../../src/cli/theme"

async function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

async function main() {
  UI.empty()

  // Opening banner
  UI.println(UI.logo(undefined, "gradient"))
  UI.println()
  UI.banner("🎨 ENHANCED TUI SHOWCASE 🎨")
  UI.println()

  await sleep(500)

  // ============================================
  // SECTION 1: Progress Bars
  // ============================================
  CyberpunkPrompts.divider("Progress Bars")
  UI.println()

  UI.println(UI.Style.BOLD + "1. Basic Progress Bars:" + UI.Style.RESET)
  UI.println()

  for (let i = 0; i <= 100; i += 20) {
    UI.println(`  ${EnhancedTUI.progressBar(i, 100, { label: "Loading", width: 30 })}`)
  }

  UI.println()
  UI.println(UI.Style.BOLD + "2. Styled Progress Bars:" + UI.Style.RESET)
  UI.println()

  UI.println(
    `  ${EnhancedTUI.progressBar(75, 100, {
      label: "Matrix Green",
      color: UI.Style.MATRIX_GREEN,
      width: 35,
    })}`
  )
  UI.println(
    `  ${EnhancedTUI.progressBar(45, 100, {
      label: "Claude Blue ",
      color: UI.Style.CLAUDE_BLUE,
      width: 35,
    })}`
  )
  UI.println(
    `  ${EnhancedTUI.progressBar(90, 100, {
      label: "Cyber Purple",
      color: UI.Style.CYBER_PURPLE,
      width: 35,
    })}`
  )
  UI.println()

  await sleep(800)

  // ============================================
  // SECTION 2: Tables
  // ============================================
  CyberpunkPrompts.divider("Formatted Tables")
  UI.println()

  const tableData = [
    [
      UI.Style.MATRIX_GREEN + "Active" + UI.Style.RESET,
      "opencode/main",
      "1,234",
      "98.5%",
    ],
    [
      UI.Style.CLAUDE_BLUE + "Running" + UI.Style.RESET,
      "feature/tui",
      "567",
      "87.2%",
    ],
    [
      UI.Style.TEXT_WARNING + "Pending" + UI.Style.RESET,
      "fix/auth",
      "89",
      "45.0%",
    ],
    [UI.Style.TEXT_DIM + "Stopped" + UI.Style.RESET, "dev/test", "12", "12.1%"],
  ]

  UI.println(
    EnhancedTUI.table(
      ["Status", "Branch", "Commits", "Coverage"],
      tableData,
      {
        align: ["center", "left", "right", "right"],
      }
    )
  )
  UI.println()

  await sleep(800)

  // ============================================
  // SECTION 3: Timeline
  // ============================================
  CyberpunkPrompts.divider("Event Timeline")
  UI.println()

  const events = [
    {
      time: "10:23 AM",
      title: "Build Started",
      description: "Compiling TypeScript sources...",
      status: "info" as const,
    },
    {
      time: "10:24 AM",
      title: "Tests Passed",
      description: "All 142 tests completed successfully",
      status: "success" as const,
    },
    {
      time: "10:24 AM",
      title: "Linting Warning",
      description: "3 style warnings found in src/index.ts",
      status: "warning" as const,
    },
    {
      time: "10:25 AM",
      title: "Deployment Failed",
      description: "Could not connect to production server",
      status: "error" as const,
    },
  ]

  UI.println(EnhancedTUI.timeline(events))
  UI.println()

  await sleep(800)

  // ============================================
  // SECTION 4: Wizard Steps
  // ============================================
  CyberpunkPrompts.divider("Multi-Step Wizard")
  UI.println()

  const wizardSteps = [
    "Configure Project",
    "Install Dependencies",
    "Run Tests",
    "Deploy to Production",
  ]

  UI.println(UI.Style.BOLD + "Installation Progress:" + UI.Style.RESET)
  UI.println()
  UI.println(EnhancedTUI.wizardSteps(wizardSteps, 2))
  UI.println()

  await sleep(800)

  // ============================================
  // SECTION 5: Notifications
  // ============================================
  CyberpunkPrompts.divider("Notifications")
  UI.println()

  UI.println(
    EnhancedTUI.notification("Build completed successfully!", "success", {
      title: "Success",
    })
  )
  UI.println()

  UI.println(
    EnhancedTUI.notification("Failed to connect to database", "error", {
      title: "Connection Error",
    })
  )
  UI.println()

  UI.println(
    EnhancedTUI.notification("API rate limit approaching", "warning", {
      title: "Warning",
    })
  )
  UI.println()

  UI.println(
    EnhancedTUI.notification("New version available: v2.0.0", "info", {
      title: "Update Available",
    })
  )
  UI.println()

  await sleep(800)

  // ============================================
  // SECTION 6: Key-Value Pairs
  // ============================================
  CyberpunkPrompts.divider("Configuration Display")
  UI.println()

  const config = {
    Project: UI.Style.MATRIX_GREEN + "opencode" + UI.Style.RESET,
    Version: UI.Style.CLAUDE_BLUE + "1.0.0" + UI.Style.RESET,
    Environment: UI.Style.NEON_CYAN + "production" + UI.Style.RESET,
    "Build Time": "2.3s",
    "Total Size": "1.2 MB",
    "Dependencies": "42",
  }

  UI.println(EnhancedTUI.keyValue(config, { indent: 2 }))
  UI.println()

  await sleep(800)

  // ============================================
  // SECTION 7: Code Blocks
  // ============================================
  CyberpunkPrompts.divider("Syntax Highlighting")
  UI.println()

  const code = `function greet(name: string) {
  console.log(\`Hello, \${name}!\`)
  return true
}

greet("OpenCode")`

  UI.println(
    EnhancedTUI.codeBlock(code, {
      language: "typescript",
      showLineNumbers: true,
      highlightLines: [2, 5],
      theme: "matrix",
    })
  )
  UI.println()

  await sleep(800)

  // ============================================
  // SECTION 8: Diff Viewer
  // ============================================
  CyberpunkPrompts.divider("Git Diff Viewer")
  UI.println()

  const deletions = [
    "const oldFunction = () => {",
    '  console.log("deprecated")',
    "}",
  ]

  const additions = [
    "const newFunction = (name: string) => {",
    "  console.log(`Modern: ${name}`)",
    "  return true",
    "}",
  ]

  UI.println(EnhancedTUI.diff(additions, deletions))
  UI.println()

  await sleep(800)

  // ============================================
  // SECTION 9: Badges & Status
  // ============================================
  CyberpunkPrompts.divider("Badges & Status Indicators")
  UI.println()

  UI.println(
    "  " +
      EnhancedTUI.badge("NEW", { color: UI.Style.MATRIX_GREEN }) +
      " " +
      EnhancedTUI.badge("BETA", { color: UI.Style.CLAUDE_BLUE }) +
      " " +
      EnhancedTUI.badge("DEPRECATED", { color: UI.Style.TEXT_WARNING }) +
      " " +
      EnhancedTUI.badge("EXPERIMENTAL", {
        color: UI.Style.CYBER_PURPLE,
        variant: "outlined",
      })
  )
  UI.println()

  UI.println("  Status Indicators:")
  UI.println(
    "    " +
      EnhancedTUI.status("Server", "online") +
      "  " +
      EnhancedTUI.status("Database", "online") +
      "  " +
      EnhancedTUI.status("Cache", "loading") +
      "  " +
      EnhancedTUI.status("API", "error")
  )
  UI.println()

  await sleep(800)

  // ============================================
  // SECTION 10: File Tree
  // ============================================
  CyberpunkPrompts.divider("File Tree")
  UI.println()

  const fileStructure = {
    src: {
      cli: {
        "ui.ts": null,
        "tui-enhanced.ts": null,
        "theme.ts": null,
      },
      server: {
        "server.ts": null,
        "routes.ts": null,
      },
      "index.ts": null,
    },
    tests: {
      "ui.test.ts": null,
      "integration.test.ts": null,
    },
    "package.json": null,
    "README.md": null,
  }

  UI.println(EnhancedTUI.fileTree(fileStructure))
  UI.println()

  await sleep(800)

  // ============================================
  // SECTION 11: Interactive Menu
  // ============================================
  CyberpunkPrompts.divider("Interactive Menu")
  UI.println()

  const menuOptions = [
    { label: "Start Development Server", description: "Run with hot reload", selected: true },
    { label: "Run Tests", description: "Execute test suite" },
    { label: "Build for Production", description: "Optimize and bundle" },
    { label: "Deploy to Staging", description: "Deploy to test environment", disabled: true },
  ]

  UI.println(EnhancedTUI.menu(menuOptions))
  UI.println()

  await sleep(800)

  // ============================================
  // SECTION 12: Collapsible Sections
  // ============================================
  CyberpunkPrompts.divider("Collapsible Sections")
  UI.println()

  UI.println(
    EnhancedTUI.collapsible(
      "Build Details",
      "Source: TypeScript\nTarget: ES2022\nModules: 156\nWarnings: 0",
      false
    )
  )
  UI.println()

  UI.println(
    EnhancedTUI.collapsible(
      "Test Results",
      UI.Style.MATRIX_GREEN +
        "✓ All tests passed\n" +
        UI.Style.RESET +
        "  Total: 142 tests\n  Duration: 2.3s\n  Coverage: 94.2%",
      true
    )
  )
  UI.println()

  await sleep(800)

  // ============================================
  // SECTION 13: Spinner Demo
  // ============================================
  CyberpunkPrompts.divider("Spinner Animations")
  UI.println()

  UI.println(UI.Style.BOLD + "Spinner Types:" + UI.Style.RESET)
  UI.println()

  const spinnerTypes: Array<[string, string[]]> = [
    ["Dots", EnhancedTUI.spinnerFrames.dots],
    ["Pulse", EnhancedTUI.spinnerFrames.pulse],
    ["Cyber", EnhancedTUI.spinnerFrames.cyber],
    ["Matrix", EnhancedTUI.spinnerFrames.matrix],
  ]

  for (const [name, frames] of spinnerTypes) {
    const preview = frames.slice(0, 8).join(" ")
    UI.println(`  ${name.padEnd(10)}: ${UI.Style.NEON_CYAN}${preview}${UI.Style.RESET}`)
  }
  UI.println()

  // Live spinner demo
  UI.println(UI.Style.BOLD + "Live Spinner Demo:" + UI.Style.RESET)
  UI.println()

  const spin = EnhancedTUI.spinner("Processing files...", {
    frames: EnhancedTUI.spinnerFrames.dots,
    color: UI.Style.NEON_CYAN,
  })

  spin.start()
  await sleep(2000)
  spin.update("Compiling TypeScript...")
  await sleep(1500)
  spin.update("Running tests...")
  await sleep(1500)
  spin.stop("All tasks completed!")
  UI.println()

  await sleep(500)

  // ============================================
  // FINAL BANNER
  // ============================================
  UI.println()
  UI.println(
    UI.box(
      UI.gradient("🚀 Enhanced TUI Demo Complete! 🚀", [
        UI.Style.MATRIX_GREEN,
        UI.Style.GEMINI_GREEN,
        UI.Style.CLAUDE_BLUE,
        UI.Style.CYBER_PURPLE,
      ]),
      { color: UI.Style.NEON_CYAN, padding: 2 }
    )
  )
  UI.println()

  // Summary
  const summary = {
    "Components Shown": "13 categories",
    "Visual Effects": "Gradients, Glow, Animations",
    "Interactive Elements": "Progress, Spinners, Tables, Menus",
    "Theme": "Cyberpunk/Matrix",
  }

  UI.println(EnhancedTUI.keyValue(summary, { indent: 2, keyColor: UI.Style.MATRIX_GREEN }))
  UI.println()

  UI.println(
    UI.Style.TEXT_INFO +
      "💡 Tip: Use these components in your CLI to create engaging user experiences!" +
      UI.Style.RESET
  )
  UI.println()
}

main().catch(console.error)
