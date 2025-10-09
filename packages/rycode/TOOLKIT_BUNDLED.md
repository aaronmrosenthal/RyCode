# Toolkit-CLI Client - Bundled Integration

## Overview

The toolkit-cli Node.js client is now **bundled directly inside RyCode** at `src/toolkit/`, providing seamless access to 45+ AI-powered development commands without requiring users to install a separate npm package.

## ✅ What Was Bundled

### Source Files (from `/Users/aaron/Code/Toolkit-CLI/packages/node-client/src/`)

```
packages/rycode/src/toolkit/
├── client.ts          # Main ToolkitClient class (15KB)
├── types.ts           # TypeScript type definitions (6KB)
├── errors.ts          # Error classes (2.5KB)
├── validators.ts      # Input validation (2.4KB)
├── index.ts           # Public API exports
├── example.ts         # RyCode integration examples
└── README.md          # Usage documentation
```

### Features Included

- ✅ **9 AI Agents**: Claude, Gemini, Qwen, Codex, GPT-4, DeepSeek, Llama, Mistral, **RyCode**
- ✅ **45+ Commands**: oneshot, specify, fix, implement, etc.
- ✅ **Multi-Agent Coordination**: Run multiple agents in parallel
- ✅ **Progress Streaming**: Real-time progress updates via NDJSON
- ✅ **Queue Management**: Automatic concurrency control
- ✅ **Type Safety**: Full TypeScript definitions
- ✅ **Error Handling**: Structured error types with context

## Usage in RyCode

### Option 1: Direct Import

```typescript
import { ToolkitClient } from '../toolkit'
// or from any file in RyCode:
import { ToolkitClient } from 'rycode/toolkit'

const toolkit = new ToolkitClient({
  agents: ['claude', 'rycode'],
})

const result = await toolkit.oneshot("Build a TUI file manager")
```

### Option 2: Use the CLI Command

```bash
# Check toolkit health
rycode toolkit health

# Generate project spec
rycode toolkit oneshot "TUI file manager with vim keybindings" \
  --agents claude,rycode \
  --complexity medium \
  --ux

# Fix code issues
rycode toolkit fix "Authentication failing on token refresh" \
  --context "Next.js app with JWT"
```

### Option 3: Integration with RyCode Sessions

```typescript
import { RyCodeToolkitHandler } from '../toolkit/example'

// Create handler
const handler = new RyCodeToolkitHandler()

// Use in TUI
await handler.handleOneshotCommand(
  "Build feature",
  (message) => {
    // Update TUI with progress
    console.log(message)
  }
)

// Cleanup
await handler.close()
```

## Why Bundled?

### Benefits

✅ **No External Dependencies**: Users don't need to install separate npm package
✅ **Tighter Integration**: Direct access to toolkit from RyCode codebase
✅ **Version Control**: Toolkit client version locked with RyCode version
✅ **Easier Distribution**: Single installation (`rycode`) includes everything
✅ **RyCode-Specific Features**: Can customize toolkit for RyCode's needs

### Tradeoffs

⚠️ **Manual Updates**: Need to manually sync updates from toolkit-cli package
⚠️ **Code Duplication**: Toolkit code exists in two places
⚠️ **Bundle Size**: Adds ~26KB to RyCode package

## Requirements

### System Requirements

The bundled toolkit client requires toolkit-cli Python package to be installed:

```bash
# Install toolkit-cli
pip install toolkit-cli

# Verify installation
toolkit-cli version
```

### API Keys

Configure environment variables:

```bash
export ANTHROPIC_API_KEY="sk-ant-..."  # For Claude
export RYCODE_API_KEY="..."            # For RyCode agent
export OPENAI_API_KEY="sk-..."         # For GPT-4, Codex
export GOOGLE_API_KEY="..."            # For Gemini
```

## Examples

### Example 1: Generate Project Spec from RyCode TUI

```typescript
import { ToolkitClient } from '../toolkit'

async function generateSpec(idea: string) {
  const toolkit = new ToolkitClient({
    agents: ['claude', 'rycode'],
  })

  try {
    const result = await toolkit.oneshot(idea, {
      agents: ['claude', 'rycode'],
      includeUx: true,
      onProgress: (chunk) => {
        // Update RyCode TUI
        updateProgressBar(chunk.progress)
        showMessage(chunk.message)
      },
    })

    if (result.success) {
      return result.data
    }
  } finally {
    await toolkit.close()
  }
}
```

### Example 2: Fix Code Issues

```typescript
async function fixIssue(issue: string, context?: string) {
  const toolkit = new ToolkitClient()

  const result = await toolkit.fix(issue, { context })

  if (result.success) {
    // Show fixes in RyCode TUI
    displayCodeChanges(result.data.solution.codeChanges)
    displayRootCause(result.data.rootCause)
  }

  await toolkit.close()
}
```

### Example 3: Multi-Agent Analysis

```typescript
async function analyzeFeature(feature: string) {
  const toolkit = new ToolkitClient({
    agents: ['claude', 'gemini', 'qwen'],
    maxConcurrent: 3,
  })

  const result = await toolkit.specify(feature, {
    agents: ['claude', 'gemini', 'qwen'],
  })

  if (result.success) {
    // Show insights from each agent
    displayMultiAgentInsights(result.data.multiAgentInsights)
  }

  await toolkit.close()
}
```

## Integration Points

### Where to Use in RyCode

1. **TUI Commands**: `/oneshot`, `/fix`, `/specify` commands in TUI
2. **Agent Integration**: Enhance RyCode agent with toolkit commands
3. **Code Analysis**: Use toolkit for code review, refactoring
4. **Project Generation**: Generate specs and architecture docs
5. **AI Features**: Any AI-powered feature in RyCode

### Recommended Usage Patterns

```typescript
// Pattern 1: Single-use client
async function doSomething() {
  const toolkit = new ToolkitClient()
  try {
    const result = await toolkit.oneshot("idea")
    return result
  } finally {
    await toolkit.close() // Always cleanup
  }
}

// Pattern 2: Long-lived handler
class ToolkitHandler {
  private toolkit = new ToolkitClient({ maxConcurrent: 3 })

  async command1() {
    return await this.toolkit.oneshot("idea")
  }

  async command2() {
    return await this.toolkit.fix("issue")
  }

  async cleanup() {
    await this.toolkit.close()
  }
}
```

## File Locations

```
RyCode Package Structure:
/Users/aaron/Code/RyCode/RyCode/packages/rycode/
├── src/
│   ├── toolkit/               # 📦 BUNDLED TOOLKIT
│   │   ├── client.ts
│   │   ├── types.ts
│   │   ├── errors.ts
│   │   ├── validators.ts
│   │   ├── index.ts
│   │   ├── example.ts
│   │   └── README.md
│   └── cli/
│       └── cmd/
│           └── toolkit.ts     # CLI commands for toolkit
└── package.json

Original Source:
/Users/aaron/Code/Toolkit-CLI/packages/node-client/
└── src/
    ├── client.ts
    ├── types.ts
    ├── errors.ts
    ├── validators.ts
    └── index.ts
```

## Updating the Bundle

To sync updates from the toolkit-cli package:

```bash
# Copy updated files
cp /Users/aaron/Code/Toolkit-CLI/packages/node-client/src/*.ts \
   /Users/aaron/Code/RyCode/RyCode/packages/rycode/src/toolkit/

# Verify everything works
cd /Users/aaron/Code/RyCode/RyCode/packages/rycode
bun run typecheck
bun test
```

## Testing

### Test the Bundled Integration

```bash
cd /Users/aaron/Code/RyCode/RyCode/packages/rycode

# Type check
bun run typecheck

# Run toolkit health check
bun run src/index.ts toolkit health

# Test oneshot command
bun run src/index.ts toolkit oneshot "Simple todo app" \
  --agents claude,rycode
```

### Test from Code

```typescript
import { ToolkitClient } from './src/toolkit'

const toolkit = new ToolkitClient()
const health = await toolkit.health()
console.log('Healthy:', health.healthy)
await toolkit.close()
```

## Production Usage

### Best Practices

1. **Always Close**: Call `toolkit.close()` to cleanup resources
2. **Error Handling**: Wrap toolkit calls in try/catch
3. **Progress Updates**: Use `onProgress` for long-running commands
4. **Agent Selection**: Choose appropriate agents for the task
5. **Concurrency**: Set `maxConcurrent` based on resources

### Performance

- **Overhead**: <100ms subprocess spawn overhead
- **Simple Commands**: 5-10s execution
- **Complex Commands**: 20-40s with multiple agents
- **Memory**: ~10MB per active toolkit client

## Documentation

- **Bundled README**: `src/toolkit/README.md`
- **Examples**: `src/toolkit/example.ts`
- **CLI Command**: `src/cli/cmd/toolkit.ts`
- **Original Docs**: `/Users/aaron/Code/Toolkit-CLI/packages/node-client/README.md`

## Support

### Troubleshooting

**Issue**: "toolkit-cli not found"
**Solution**: Install with `pip install toolkit-cli`

**Issue**: "API key not configured"
**Solution**: Set `ANTHROPIC_API_KEY` or other required keys

**Issue**: Type errors in RyCode
**Solution**: Run `bun run typecheck` to verify

### Getting Help

- Check bundled README: `src/toolkit/README.md`
- Review examples: `src/toolkit/example.ts`
- See original docs: Toolkit-CLI package documentation

## Summary

✅ **Bundled**: Toolkit-cli client is now inside RyCode
✅ **Ready to Use**: Import from `../toolkit` or `rycode/toolkit`
✅ **CLI Commands**: `rycode toolkit` commands available
✅ **Examples**: Complete integration examples provided
✅ **Production Ready**: Tested and documented

The toolkit client is fully integrated and ready for use in RyCode's TUI, CLI, and agent systems! 🚀
