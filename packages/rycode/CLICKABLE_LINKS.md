# 🔗 Clickable Hyperlinks in Terminal

## Overview

OpenCode now supports **clickable hyperlinks** in terminal output using OSC 8 escape sequences. Links work in modern terminals and allow users to click file paths and URLs directly.

## Supported Terminals

✅ **Fully Supported:**
- iTerm2 (macOS)
- Alacritty
- Kitty
- WezTerm
- Windows Terminal
- VSCode Terminal
- GNOME Terminal 3.28+
- Konsole

⚠️ **Fallback:** In unsupported terminals, links display as regular text.

## API Reference

### `UI.link(text, url)`

Create a clickable link.

```typescript
import { UI } from "./cli/ui"

// Basic link
UI.link("Click here", "https://example.com")

// File link
UI.link("README", "file:///path/to/README.md")
```

### `UI.fileLink(path, displayText?)`

Create a clickable file path link. Automatically converts paths to `file://` URLs.

```typescript
// Absolute path
UI.fileLink("/Users/aaron/Code/project/README.md")

// Custom display text
UI.fileLink("/Users/aaron/tasks.md", "My Tasks")
```

### `UI.autoLink(text)`

Automatically detect and link file paths and URLs in text.

```typescript
const text = "Check /Users/aaron/tasks.md and visit https://github.com"
console.log(UI.autoLink(text))
// Makes both the file path and URL clickable
```

### `UI.styledLink(text, url, icon?)`

Create a styled clickable link with icon.

```typescript
UI.styledLink("Documentation", "https://docs.example.com", "📚")
// Output: 📚 Documentation (clickable, styled)
```

### `UI.styledFileLink(path, displayText?)`

Create a styled file link with file icon.

```typescript
UI.styledFileLink("/Users/aaron/Code/project/README.md")
// Output: 📄 /Users/aaron/Code/project/README.md (clickable)

UI.styledFileLink("/path/to/file.md", "Config File")
// Output: 📄 Config File (clickable)
```

## Usage Examples

### Example 1: "What's Next?" with Clickable Files

```typescript
import { whatsNext } from "./cli/whats-next"

whatsNext(
  [
    {
      key: "A",
      label: "Push to remote",
      time: "30 sec",
      command: "git push origin main",
      description: "Deploy your changes to the repository.",
    },
    {
      key: "B",
      label: "Clean up tasks.md",
      time: "10 min",
      description: "Fix the malformed tasks in /Users/aaron/Code/Toolkit-CLI/tasks.md",
      // Path in description will be auto-linked!
    },
    {
      key: "C",
      label: "Start Quick-Win implementation",
      time: "4 weeks",
      file: "/Users/aaron/Code/quick-win", // Explicit file link
    },
  ],
  "Push to remote first (A), then tackle tasks.md cleanup (B)"
)
```

**Output:**
```
🎯 What's Next?

Immediate Options:

[A] Push to remote (30 sec)
git push origin main
Deploy your changes to the repository.

[B] Clean up tasks.md (10 min)
Fix the malformed tasks in /Users/aaron/Code/Toolkit-CLI/tasks.md
                              ↑ clickable!

[C] Start Quick-Win implementation (4 weeks)
📄 /Users/aaron/Code/quick-win
   ↑ clickable!

───
Recommendation: Push to remote first (A), then tackle tasks.md cleanup (B)

Type push, clean, or start to proceed, or ask for other suggestions.
```

### Example 2: Error Messages with File Links

```typescript
import { UI } from "./cli/ui"

const errorFile = "/Users/aaron/Code/project/src/index.ts"
const errorLine = 42

UI.error(
  `Parse error in ${UI.fileLink(errorFile)}:${errorLine}\n` +
  `Click the link above to open the file.`
)
```

### Example 3: Documentation References

```typescript
import { UI } from "./cli/ui"

console.log(
  "For more information, see:\n" +
  UI.styledLink("Documentation", "https://docs.opencode.dev", "📖") + "\n" +
  UI.styledFileLink("/Users/aaron/Code/project/README.md", "Local README")
)
```

### Example 4: Auto-Linking in Output

```typescript
import { UI } from "./cli/ui"

const output = `
Task completed!

Modified files:
  /Users/aaron/Code/project/src/index.ts
  /Users/aaron/Code/project/README.md

For details, visit https://github.com/user/repo/pull/123
`

console.log(UI.autoLink(output))
// All paths and URLs become clickable
```

## How It Works

OSC 8 is a terminal escape sequence standard for hyperlinks:

```
\x1b]8;;URL\x1b\\TEXT\x1b]8;;\x1b\\
```

Example:
```typescript
// This:
UI.link("Click me", "https://example.com")

// Generates:
"\x1b]8;;https://example.com\x1b\\Click me\x1b]8;;\x1b\\"

// Terminal renders it as:
Click me  <- clickable link
```

## File Protocol

File paths are converted to `file://` URLs:

```typescript
UI.fileLink("/Users/aaron/tasks.md")
// Becomes: file:///Users/aaron/tasks.md

// Windows paths:
UI.fileLink("C:/Users/aaron/tasks.md")
// Becomes: file:///C:/Users/aaron/tasks.md
```

## Click Behavior

- **macOS:** `Cmd + Click` to open
- **Linux:** `Ctrl + Click` to open
- **Windows:** `Ctrl + Click` to open

**Files:** Open in default editor (usually VSCode if installed)
**URLs:** Open in default browser

## Testing

Run the demo to see clickable links in action:

```bash
bun run src/cli/cmd/demo-links.ts
```

Or test the "What's Next?" component:

```bash
bun run src/cli/whats-next.ts
```

## Tips

### 1. **Always Use Absolute Paths**

Relative paths may not work correctly:
```typescript
// ❌ Avoid
UI.fileLink("./tasks.md")

// ✅ Better
UI.fileLink(process.cwd() + "/tasks.md")

// ✅ Best
import { resolve } from "path"
UI.fileLink(resolve("./tasks.md"))
```

### 2. **Handle Unsupported Terminals Gracefully**

Links degrade gracefully to plain text in unsupported terminals. No special handling needed.

### 3. **Use `autoLink` for User-Generated Content**

If displaying content that might contain paths or URLs:
```typescript
console.log(UI.autoLink(userMessage))
```

### 4. **Combine with Styling**

```typescript
// Link with color and underline
const styled = UI.Style.UNDERLINE + UI.Style.CLAUDE_BLUE +
               UI.link("Click", "https://example.com") +
               UI.Style.RESET

// Or use the helper
UI.styledLink("Click", "https://example.com")
```

## Browser-Based Terminals

For browser-based terminals (like code-server or web SSH), OSC 8 support depends on the terminal emulator implementation (e.g., xterm.js supports OSC 8).

## Security

**Important:** Only link to trusted URLs and file paths. Do not create links from untrusted user input without validation.

```typescript
// ❌ Dangerous
UI.link("Click", userProvidedURL) // Could be javascript: or other malicious scheme

// ✅ Safe
if (userURL.startsWith("http://") || userURL.startsWith("https://")) {
  UI.link("Click", userURL)
}
```

## See Also

- [OSC 8 Specification](https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda)
- [Terminal Hyperlinks Support](https://github.com/Alhadis/OSC8-Adoption)
