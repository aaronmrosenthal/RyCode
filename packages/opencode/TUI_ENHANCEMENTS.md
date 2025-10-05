# 🎨 Enhanced TUI Components

Advanced terminal UI components for creating engaging, verbose, and interactive command-line experiences.

## Overview

The Enhanced TUI library provides a comprehensive set of components to make your CLI interactions more verbose, visually appealing, and user-friendly. Built on top of the existing UI system with cyberpunk/matrix aesthetic.

## Features

### 📊 **Progress Indicators**
- Animated progress bars with customizable colors
- Multiple spinner styles (dots, pulse, cyber, matrix, etc.)
- Real-time progress updates
- Percentage display

### 📋 **Data Display**
- Rich formatted tables with borders and alignment
- Key-value pair displays with auto-alignment
- Timeline views for events
- File tree visualization

### 🎯 **Interactive Elements**
- Multi-step wizard indicators
- Interactive menus with descriptions
- Collapsible sections
- Status indicators (online/offline/loading/error)

### 🎨 **Visual Components**
- Syntax-highlighted code blocks
- Git-style diff viewers
- Badges and tags
- Notifications (success/error/warning/info)
- Gradient text effects

## Quick Start

```typescript
import { EnhancedTUI } from "./src/cli/tui-enhanced"
import { UI } from "./src/cli/ui"

// Show a progress bar
console.log(EnhancedTUI.progressBar(75, 100, {
  label: "Loading",
  color: UI.Style.MATRIX_GREEN
}))

// Create a spinner
const spinner = EnhancedTUI.spinner("Processing...")
spinner.start()
// ... do work ...
spinner.stop("Complete!")

// Display a table
console.log(EnhancedTUI.table(
  ["Name", "Status", "Progress"],
  [
    ["Build", "✓ Success", "100%"],
    ["Tests", "⟳ Running", "75%"],
  ]
))
```

## Component Reference

### 1. Progress Bars

Display progress with animated bars.

```typescript
EnhancedTUI.progressBar(current: number, total: number, options?: {
  width?: number          // Bar width (default: 40)
  label?: string         // Optional label
  showPercentage?: boolean  // Show % (default: true)
  color?: string        // Bar color
  bgColor?: string      // Background color
})
```

**Example:**
```typescript
// Basic progress bar
EnhancedTUI.progressBar(75, 100, { label: "Downloading" })
// Output: Downloading [████████████████████████░░░░░] 75%

// Styled progress bar
EnhancedTUI.progressBar(60, 100, {
  label: "Building",
  color: UI.Style.CLAUDE_BLUE,
  width: 30
})
```

### 2. Spinners

Animated loading indicators with multiple frame styles.

```typescript
const spinner = EnhancedTUI.spinner(message: string, options?: {
  frames?: string[]     // Animation frames
  color?: string       // Spinner color
  interval?: number    // Frame interval in ms
})

spinner.start()          // Start animation
spinner.update(message)  // Update message
spinner.stop(message?)   // Stop with success message
spinner.fail(message?)   // Stop with error message
```

**Available frame styles:**
- `spinnerFrames.dots` - Classic dots: ⠋⠙⠹⠸⠼⠴⠦⠧
- `spinnerFrames.pulse` - Pulsing circle: ○◔◑◕●
- `spinnerFrames.cyber` - Vertical bars: ▁▂▃▄▅▆▇█
- `spinnerFrames.matrix` - Matrix style: ⡀⠄⠂⠁⠈⠐⠠⢀
- `spinnerFrames.neon` - Neon corners: ◜◠◝◞◡◟
- `spinnerFrames.arrow` - Rotating arrows: ←↖↑↗→↘↓↙
- `spinnerFrames.line` - Classic line: -\|/

**Example:**
```typescript
const spin = EnhancedTUI.spinner("Installing dependencies...", {
  frames: EnhancedTUI.spinnerFrames.dots,
  color: UI.Style.NEON_CYAN
})

spin.start()
await installPackages()
spin.update("Running tests...")
await runTests()
spin.stop("All tasks completed!")
```

### 3. Tables

Rich formatted tables with styled borders and alignment.

```typescript
EnhancedTUI.table(
  headers: string[],
  rows: string[][],
  options?: {
    headerColor?: string
    borderColor?: string
    align?: ("left" | "center" | "right")[]
    maxWidth?: number
  }
)
```

**Example:**
```typescript
EnhancedTUI.table(
  ["Status", "Task", "Time", "Progress"],
  [
    ["✓", "Build", "2.3s", "100%"],
    ["⟳", "Test", "1.5s", "75%"],
    ["○", "Deploy", "—", "0%"]
  ],
  {
    align: ["center", "left", "right", "right"],
    headerColor: UI.Style.CLAUDE_BLUE,
    borderColor: UI.Style.MATRIX_GREEN
  }
)
```

### 4. Timeline

Display events in chronological order with status indicators.

```typescript
EnhancedTUI.timeline(events: Array<{
  time: string
  title: string
  description?: string
  status?: "success" | "error" | "warning" | "info"
}>, options?: {
  showConnector?: boolean
})
```

**Example:**
```typescript
EnhancedTUI.timeline([
  {
    time: "10:23 AM",
    title: "Build Started",
    description: "Compiling TypeScript...",
    status: "info"
  },
  {
    time: "10:24 AM",
    title: "Tests Passed",
    description: "142 tests completed",
    status: "success"
  },
  {
    time: "10:25 AM",
    title: "Deployment Failed",
    description: "Connection refused",
    status: "error"
  }
])
```

### 5. Wizard Steps

Multi-step process indicator showing progress through a workflow.

```typescript
EnhancedTUI.wizardSteps(
  steps: string[],
  currentStep: number,
  options?: {
    completedColor?: string
    activeColor?: string
    pendingColor?: string
  }
)
```

**Example:**
```typescript
const steps = [
  "Install Dependencies",
  "Configure Project",
  "Run Tests",
  "Deploy"
]

EnhancedTUI.wizardSteps(steps, 2)
// Shows step 3 as active, 1-2 as completed, 4 as pending
```

### 6. Notifications

Toast-style notification boxes.

```typescript
EnhancedTUI.notification(
  message: string,
  type: "success" | "error" | "warning" | "info",
  options?: {
    title?: string
    duration?: number
  }
)
```

**Example:**
```typescript
console.log(EnhancedTUI.notification(
  "Build completed in 2.3 seconds",
  "success",
  { title: "Success" }
))

console.log(EnhancedTUI.notification(
  "Failed to connect to database",
  "error",
  { title: "Connection Error" }
))
```

### 7. Key-Value Display

Aligned key-value pairs for configuration display.

```typescript
EnhancedTUI.keyValue(
  pairs: Record<string, string>,
  options?: {
    keyColor?: string
    valueColor?: string
    separator?: string
    indent?: number
  }
)
```

**Example:**
```typescript
EnhancedTUI.keyValue({
  "Project": "opencode",
  "Version": "1.0.0",
  "Environment": "production",
  "Status": "✓ Running"
}, {
  keyColor: UI.Style.CLAUDE_BLUE,
  separator: " →",
  indent: 2
})
```

### 8. Code Blocks

Syntax-highlighted code with line numbers.

```typescript
EnhancedTUI.codeBlock(
  code: string,
  options?: {
    language?: string
    showLineNumbers?: boolean
    highlightLines?: number[]
    theme?: "dark" | "matrix"
  }
)
```

**Example:**
```typescript
const code = `function hello(name: string) {
  console.log(\`Hello, \${name}!\`)
  return true
}`

EnhancedTUI.codeBlock(code, {
  language: "typescript",
  showLineNumbers: true,
  highlightLines: [2],
  theme: "matrix"
})
```

### 9. Diff Viewer

Git-style diff display for code changes.

```typescript
EnhancedTUI.diff(
  additions: string[],
  deletions: string[],
  options?: {
    addColor?: string
    deleteColor?: string
  }
)
```

**Example:**
```typescript
EnhancedTUI.diff(
  ["const newVersion = '2.0.0'", "console.log('Updated!')"],
  ["const oldVersion = '1.0.0'"],
  {
    addColor: UI.Style.MATRIX_GREEN,
    deleteColor: UI.Style.TEXT_DANGER
  }
)
```

### 10. Badges

Labeled tags for status indicators.

```typescript
EnhancedTUI.badge(
  text: string,
  options?: {
    color?: string
    variant?: "filled" | "outlined"
  }
)
```

**Example:**
```typescript
const badges = [
  EnhancedTUI.badge("NEW", { color: UI.Style.MATRIX_GREEN }),
  EnhancedTUI.badge("BETA", { color: UI.Style.CLAUDE_BLUE }),
  EnhancedTUI.badge("EXPERIMENTAL", {
    color: UI.Style.CYBER_PURPLE,
    variant: "outlined"
  })
]

console.log(badges.join(" "))
```

### 11. Status Indicators

Real-time status with visual indicators.

```typescript
EnhancedTUI.status(
  label: string,
  state: "online" | "offline" | "loading" | "error",
  options?: {
    showLabel?: boolean
  }
)
```

**Example:**
```typescript
console.log([
  EnhancedTUI.status("Server", "online"),
  EnhancedTUI.status("Database", "loading"),
  EnhancedTUI.status("Cache", "error")
].join("  "))
```

### 12. File Tree

Visual file/folder hierarchy.

```typescript
EnhancedTUI.fileTree(
  tree: Record<string, any>
)
```

**Example:**
```typescript
EnhancedTUI.fileTree({
  src: {
    components: {
      "Button.tsx": null,
      "Input.tsx": null
    },
    "index.ts": null
  },
  "package.json": null,
  "README.md": null
})
```

### 13. Interactive Menus

Selection menus with descriptions.

```typescript
EnhancedTUI.menu(
  options: Array<{
    label: string
    description?: string
    selected?: boolean
    disabled?: boolean
  }>,
  options?: {
    selectedColor?: string
    disabledColor?: string
  }
)
```

**Example:**
```typescript
EnhancedTUI.menu([
  {
    label: "Start Dev Server",
    description: "Run with hot reload",
    selected: true
  },
  {
    label: "Run Tests",
    description: "Execute test suite"
  },
  {
    label: "Deploy",
    description: "Deploy to production",
    disabled: true
  }
])
```

### 14. Collapsible Sections

Expandable/collapsible content sections.

```typescript
EnhancedTUI.collapsible(
  title: string,
  content: string,
  isExpanded: boolean,
  options?: {
    expandIcon?: string
    collapseIcon?: string
    titleColor?: string
  }
)
```

**Example:**
```typescript
// Collapsed
EnhancedTUI.collapsible("Details", "Content here...", false)
// Output: ▶ Details

// Expanded
EnhancedTUI.collapsible("Details", "Content here...", true)
// Output: ▼ Details
//           Content here...
```

## Advanced Usage

### Combining Components

Create rich, verbose interfaces by combining multiple components:

```typescript
import { UI } from "./src/cli/ui"
import { EnhancedTUI } from "./src/cli/tui-enhanced"
import { CyberpunkPrompts } from "./src/cli/theme"

// Application startup sequence
UI.banner("🚀 OPENCODE INITIALIZATION")
UI.println()

// Show configuration
CyberpunkPrompts.divider("Configuration")
console.log(EnhancedTUI.keyValue({
  "Version": "1.0.0",
  "Environment": "production",
  "Port": "3000"
}, { indent: 2 }))
UI.println()

// Show wizard progress
CyberpunkPrompts.divider("Setup Progress")
console.log(EnhancedTUI.wizardSteps([
  "Load Configuration",
  "Connect to Database",
  "Start Server",
  "Ready"
], 2))
UI.println()

// Show status indicators
console.log("Services:")
console.log([
  "  " + EnhancedTUI.status("API", "online"),
  EnhancedTUI.status("Database", "online"),
  EnhancedTUI.status("Cache", "loading")
].join("  "))
UI.println()

// Success notification
console.log(EnhancedTUI.notification(
  "System initialization complete!",
  "success"
))
```

### Creating Interactive Flows

Build multi-step interactive experiences:

```typescript
async function deploymentFlow() {
  // Step 1: Show wizard
  console.log(EnhancedTUI.wizardSteps([
    "Build Project",
    "Run Tests",
    "Deploy to Server"
  ], 0))

  // Step 2: Build with progress
  const buildSpinner = EnhancedTUI.spinner("Building project...")
  buildSpinner.start()
  await buildProject()
  buildSpinner.stop("Build complete!")

  // Step 3: Tests with table results
  console.log(EnhancedTUI.table(
    ["Test Suite", "Passed", "Failed"],
    [
      ["Unit Tests", "142", "0"],
      ["Integration", "23", "0"]
    ]
  ))

  // Step 4: Timeline of deployment
  console.log(EnhancedTUI.timeline([
    { time: "14:23", title: "Upload Assets", status: "success" },
    { time: "14:24", title: "Database Migration", status: "success" },
    { time: "14:25", title: "Server Restart", status: "success" }
  ]))

  // Final notification
  console.log(EnhancedTUI.notification(
    "Deployment completed successfully!",
    "success",
    { title: "Deployment Complete" }
  ))
}
```

## Styling Guide

### Colors

Use the built-in UI.Style colors for consistent theming:

```typescript
UI.Style.MATRIX_GREEN      // Bright matrix green #00FF41
UI.Style.MATRIX_GREEN_DIM  // Dimmer matrix green
UI.Style.CLAUDE_BLUE       // Claude blue #6699FF
UI.Style.GEMINI_GREEN      // Gemini green #1AD689
UI.Style.NEON_CYAN         // Bright cyan #00FFFF
UI.Style.NEON_MAGENTA      // Bright magenta #FF00FF
UI.Style.CYBER_PURPLE      // Purple accent #9333EA

// Semantic colors
UI.Style.TEXT_SUCCESS      // Success messages
UI.Style.TEXT_DANGER       // Error messages
UI.Style.TEXT_WARNING      // Warning messages
UI.Style.TEXT_INFO         // Info messages
UI.Style.TEXT_DIM          // Dimmed text
```

### Best Practices

1. **Use progress indicators** for long-running operations
2. **Provide context** with labels and descriptions
3. **Use appropriate colors** - green for success, red for errors, etc.
4. **Show status updates** to keep users informed
5. **Combine components** for rich experiences
6. **Keep it readable** - don't overuse colors or animations
7. **Test in terminals** - verify your output looks good in different terminals

## Examples

### Build Tool Output

```typescript
CyberpunkPrompts.intro("Starting Build")

const buildSpinner = EnhancedTUI.spinner("Compiling...")
buildSpinner.start()
// ... compile ...
buildSpinner.stop("Compilation complete")

console.log(EnhancedTUI.table(
  ["File", "Size", "Gzip"],
  [
    ["index.js", "45 KB", "12 KB"],
    ["vendor.js", "234 KB", "78 KB"]
  ]
))

console.log(EnhancedTUI.notification(
  "Build completed in 2.3s",
  "success"
))
```

### Project Dashboard

```typescript
UI.banner("PROJECT DASHBOARD")

// Status overview
console.log(EnhancedTUI.keyValue({
  "Status": EnhancedTUI.status("System", "online", { showLabel: false }),
  "Uptime": "3 days",
  "Requests": "1.2M",
  "Errors": "0.01%"
}))

// Recent events
console.log(EnhancedTUI.timeline([
  { time: "Just now", title: "API Request", status: "success" },
  { time: "2m ago", title: "Deploy Complete", status: "success" },
  { time: "5m ago", title: "Warning", status: "warning" }
]))
```

## Running the Demo

To see all components in action:

```bash
bun run demo-tui-enhanced.ts
```

This showcases:
- All component types
- Color schemes
- Animation effects
- Layout combinations
- Best practices

## See Also

- [CLICKABLE_LINKS.md](./CLICKABLE_LINKS.md) - Clickable hyperlinks guide
- [src/cli/ui.ts](./src/cli/ui.ts) - Base UI utilities
- [src/cli/theme.ts](./src/cli/theme.ts) - Cyberpunk theme
- [src/cli/whats-next.ts](./src/cli/whats-next.ts) - "What's Next?" component

## Contributing

Want to add more components? Follow these guidelines:

1. **Keep it terminal-friendly** - Test in multiple terminals
2. **Use ANSI escape codes** - Stick to widely-supported codes
3. **Provide options** - Make colors and styles configurable
4. **Document well** - Add clear examples
5. **Match the aesthetic** - Follow the cyberpunk/matrix theme

Happy building! 🎨✨
