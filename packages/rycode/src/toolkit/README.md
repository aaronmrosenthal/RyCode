# Toolkit-CLI Client (Bundled in RyCode)

This directory contains the toolkit-cli Node.js client, bundled directly into RyCode for seamless integration.

## Overview

The toolkit-cli client provides access to 45+ AI-powered development commands through a TypeScript API, with special support for RyCode's TUI integration.

## Usage in RyCode

```typescript
import { ToolkitClient } from '../toolkit'
// or
import { ToolkitClient } from 'rycode/toolkit'

// Create client with RyCode agent
const toolkit = new ToolkitClient({
  agents: ['claude', 'rycode'],
  apiKeys: {
    anthropic: process.env.ANTHROPIC_API_KEY,
    rycode: process.env.RYCODE_API_KEY
  }
})

// Generate project spec
const result = await toolkit.oneshot(
  "TUI file manager with vim keybindings",
  {
    agents: ['claude', 'rycode'],
    includeUx: true
  }
)
```

## Features

- ✅ **9 AI Agents**: Claude, Gemini, Qwen, Codex, GPT-4, DeepSeek, Llama, Mistral, RyCode
- ✅ **45+ Commands**: oneshot, specify, fix, implement, etc.
- ✅ **Multi-Agent**: Run multiple agents in parallel
- ✅ **Progress Streaming**: Real-time progress updates
- ✅ **Type Safe**: Full TypeScript definitions
- ✅ **Queue Management**: Automatic concurrency control

## Files

- `client.ts` - Main ToolkitClient class
- `types.ts` - TypeScript type definitions
- `errors.ts` - Error classes
- `validators.ts` - Input validation
- `index.ts` - Public API exports

## Configuration

The client requires toolkit-cli to be installed:

```bash
pip install toolkit-cli
```

## API Keys

Configure in environment or programmatically:

```bash
export ANTHROPIC_API_KEY="sk-ant-..."
export RYCODE_API_KEY="..."
```

## Examples

See `examples/` directory for usage examples.

## Source

This is a bundled copy of the toolkit-cli Node.js client from:
`/Users/aaron/Code/Toolkit-CLI/packages/node-client`

To update, copy from source or regenerate the bundle.
