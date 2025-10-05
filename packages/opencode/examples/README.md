# RyCode Examples & Demos

This directory contains examples and demonstrations of RyCode features.

## 📂 Directory Structure

```
examples/
├── demos/           # UI/TUI feature demonstrations
└── README.md        # This file
```

## 🎨 Demos

### TUI & UI Demonstrations

Located in `demos/`:

- **`demo-installer.ts`** - Polished installer experience showcase
- **`demo-tui-enhanced.ts`** - Enhanced TUI components demonstration
- **`demo-matrix-logo.ts`** - Matrix-themed logo animations
- **`demo-rycode-logo.ts`** - RyCode logo variants
- **`rycode-logos.ts`** - Collection of ASCII art logos
- **`test-clickable.ts`** - Clickable terminal links testing
- **`whats-next-enhanced.ts`** - "What's Next?" workflow demo

### Running Demos

All demos are executable with Bun:

```bash
# From the opencode package root
bun examples/demos/demo-installer.ts
bun examples/demos/demo-tui-enhanced.ts
bun examples/demos/test-clickable.ts
# ... etc
```

## 🔧 Usage

These demos showcase:

- ✨ Terminal UI components and styling
- 🎨 Logo and branding animations
- 🔗 Clickable terminal links (OSC 8 sequences)
- 📊 Progress indicators and status displays
- 🎯 Professional installer flows
- 💫 Enhanced user experience patterns

## 📝 Notes

- All demos are standalone and don't require additional setup
- Demos use the actual RyCode UI components from `src/cli/`
- Some demos include `sleep()` delays for visual effect
- Terminal support for features like clickable links may vary by terminal emulator

## 🚀 Adding New Examples

To add a new demo:

1. Create a new `.ts` file in `demos/`
2. Add shebang: `#!/usr/bin/env bun`
3. Make it executable: `chmod +x examples/demos/your-demo.ts`
4. Import necessary components from `../../src/`
5. Document it in this README

---

**Last Updated:** October 5, 2025
