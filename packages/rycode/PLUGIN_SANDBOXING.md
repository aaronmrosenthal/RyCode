# Plugin Sandboxing

> **Process-level isolation using worker threads for maximum security**

Plugin sandboxing provides true process isolation, preventing malicious or buggy plugins from compromising your system or other plugins. Each plugin runs in its own isolated worker thread with strict resource limits and capability restrictions.

---

## 🎯 Overview

RyCode's sandboxing system provides:

- ✅ **Process Isolation** - Plugins run in separate worker threads
- ✅ **Resource Limits** - Memory, CPU, and execution time constraints
- ✅ **Timeout Protection** - Automatic termination of runaway plugins
- ✅ **Capability Enforcement** - Fine-grained permission control
- ✅ **Resource Monitoring** - Real-time tracking of resource usage
- ✅ **Graceful Termination** - Clean shutdown and error handling

---

## 🚀 Quick Start

### Basic Usage

```typescript
import { PluginSandbox } from "./src/plugin/sandbox"
import { PluginSecurity } from "./src/plugin/security"

// Create sandbox
const sandbox = await PluginSandbox.createSandbox({
  pluginName: "my-plugin",
  pluginVersion: "1.0.0",
  capabilities: {
    fileSystemRead: true,
    fileSystemWrite: false,
    network: true,
    shell: false,
    env: false,
    projectMetadata: true,
    aiClient: true,
  },
  resourceLimits: {
    maxMemoryMB: 512,
    maxExecutionTime: 30000, // 30 seconds
    maxCPUTime: 10000,       // 10 seconds
    maxFileSystemOps: 1000,
    maxNetworkRequests: 100,
  },
  strictMode: true,
})

// Execute plugin
const result = await sandbox.execute({ input: "data" })

// Get resource usage
const usage = sandbox.getResourceUsage()
console.log(`Memory: ${usage.memoryMB}MB, CPU: ${usage.cpuTime}ms`)

// Terminate sandbox
await sandbox.terminate()
```

---

## 📖 Architecture

### Components

```
┌─────────────────────────────────────────┐
│         Main Thread                     │
│  ┌──────────────────────────────────┐  │
│  │   Sandbox Manager                │  │
│  │  - Resource monitoring           │  │
│  │  - Timeout enforcement            │  │
│  │  - Capability checking            │  │
│  └──────────────┬───────────────────┘  │
│                 │ IPC                   │
│                 ▼                       │
│  ┌──────────────────────────────────┐  │
│  │   Worker Thread                  │  │
│  │  ┌────────────────────────────┐  │  │
│  │  │  Sandboxed Environment     │  │  │
│  │  │  - Plugin execution        │  │  │
│  │  │  - Wrapped APIs            │  │  │
│  │  │  - Resource tracking       │  │  │
│  │  └────────────────────────────┘  │  │
│  └──────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

### Security Layers

1. **Worker Thread Isolation** - Separate process space
2. **Capability Restrictions** - API access control
3. **Resource Limits** - Memory and CPU constraints
4. **Timeout Enforcement** - Execution time limits
5. **Monitored APIs** - Tracked file system and network access

---

## ⚙️ Configuration

### Resource Limits

```typescript
const resourceLimits = {
  // Maximum memory in MB (default: 512)
  maxMemoryMB: 512,

  // Maximum execution time in ms (default: 30000)
  maxExecutionTime: 30000,

  // Maximum CPU time in ms (default: 10000)
  maxCPUTime: 10000,

  // Maximum file system operations (default: 1000)
  maxFileSystemOps: 1000,

  // Maximum network requests (default: 100)
  maxNetworkRequests: 100,
}
```

### Capabilities

```typescript
const capabilities = {
  // Read files
  fileSystemRead: true,

  // Write files (HIGH RISK - only grant if needed)
  fileSystemWrite: false,

  // Make network requests
  network: true,

  // Execute shell commands (VERY HIGH RISK)
  shell: false,

  // Access environment variables
  env: false,

  // Access project metadata
  projectMetadata: true,

  // Access AI client
  aiClient: true,
}
```

### Strict Mode

```typescript
{
  // Disable eval(), Function(), and dynamic require
  strictMode: true, // Recommended for production
}
```

---

## 🛡️ Security Guarantees

### What Sandboxing Protects Against

✅ **Memory Leaks** - Automatic cleanup on termination
✅ **Infinite Loops** - Timeout enforcement
✅ **Resource Exhaustion** - Memory and CPU limits
✅ **Unauthorized File Access** - Capability-based restrictions
✅ **Uncontrolled Network Requests** - Request counting and limits
✅ **Process Interference** - Isolated execution environment

### What Sandboxing Does NOT Protect Against

⚠️ **Side-channel attacks** - Worker threads share the same process
⚠️ **Timing attacks** - No timing protections
⚠️ **Spectre/Meltdown** - CPU-level vulnerabilities
⚠️ **Supply chain attacks** - Malicious dependencies

**Mitigation:** Combine sandboxing with registry verification and signatures.

---

## 📊 Resource Monitoring

### Real-time Tracking

```typescript
// Monitor resource usage
setInterval(() => {
  const usage = sandbox.getResourceUsage()
  console.log(`
    Memory: ${usage.memoryMB}MB
    CPU: ${usage.cpuTime}ms
    File Ops: ${usage.fileSystemOps}
    Network: ${usage.networkRequests}
  `)
}, 1000)
```

### Global Statistics

```typescript
const stats = PluginSandbox.getStatistics()

console.log(`
  Active Sandboxes: ${stats.activeSandboxes}
  Total Memory: ${stats.totalMemoryMB}MB
  Average Memory: ${stats.averageMemoryMB}MB
  Total CPU: ${stats.totalCPUTime}ms
  Average CPU: ${stats.averageCPUTime}ms
`)
```

---

## 🔄 Lifecycle Management

### Creation

```typescript
const sandbox = await PluginSandbox.createSandbox(config)
```

### Execution

```typescript
try {
  const result = await sandbox.execute(input)
  console.log("Success:", result)
} catch (error) {
  if (error instanceof PluginSandbox.SandboxTimeoutError) {
    console.error("Plugin timed out")
  } else if (error instanceof PluginSandbox.SandboxResourceError) {
    console.error("Resource limit exceeded:", error.data.resource)
  } else {
    console.error("Execution failed:", error)
  }
}
```

### Termination

```typescript
// Terminate single sandbox
await sandbox.terminate()

// Terminate all active sandboxes
await PluginSandbox.terminateAll()
```

---

## 💡 Best Practices

### For Plugin Users

1. **Use Minimal Capabilities**
   ```typescript
   // Only grant what's needed
   capabilities: {
     fileSystemRead: true,  // Only if plugin needs to read
     fileSystemWrite: false, // Deny unless absolutely necessary
     network: false,        // Deny unless needed
     shell: false,          // Almost never grant this
   }
   ```

2. **Set Conservative Limits**
   ```typescript
   resourceLimits: {
     maxMemoryMB: 256,      // Start small
     maxExecutionTime: 10000, // 10 seconds
     maxFileSystemOps: 100,
     maxNetworkRequests: 10,
   }
   ```

3. **Always Use Strict Mode**
   ```typescript
   strictMode: true // Prevents eval() and dynamic code
   ```

4. **Monitor Resource Usage**
   ```typescript
   const usage = sandbox.getResourceUsage()
   if (usage.memoryMB > 400) {
     console.warn("Plugin using high memory")
   }
   ```

### For Plugin Developers

1. **Design for Limits**
   - Assume low resource limits
   - Implement efficient algorithms
   - Clean up resources properly

2. **Handle Timeouts Gracefully**
   - Checkpoint progress periodically
   - Provide meaningful progress updates
   - Support resumable operations

3. **Request Minimal Capabilities**
   - Only request what you need
   - Document why each capability is required
   - Provide fallback for denied capabilities

4. **Test with Limits**
   ```typescript
   // Test with production-like limits
   const testLimits = {
     maxMemoryMB: 128,
     maxExecutionTime: 5000,
   }
   ```

---

## 🔍 Troubleshooting

### Timeout Errors

**Problem:** `SandboxTimeoutError: Plugin execution exceeded timeout`

**Solutions:**
1. Increase `maxExecutionTime` if legitimate
2. Optimize plugin code
3. Break work into smaller chunks
4. Use async operations properly

### Memory Limit Exceeded

**Problem:** `SandboxResourceError: Memory limit exceeded`

**Solutions:**
1. Increase `maxMemoryMB` if needed
2. Profile memory usage
3. Fix memory leaks
4. Process data in streams instead of loading all at once

### File System Operations Exceeded

**Problem:** Too many file operations

**Solutions:**
1. Increase `maxFileSystemOps`
2. Cache file reads
3. Batch file operations
4. Use in-memory data structures

### Worker Fails to Start

**Problem:** Worker thread creation fails

**Solutions:**
1. Check Node.js/Bun version compatibility
2. Verify worker file path is correct
3. Check system resource availability
4. Review error logs

---

## 📋 Configuration Examples

### Development Mode

```typescript
{
  capabilities: {
    fileSystemRead: true,
    fileSystemWrite: true,  // Allow for development
    network: true,
    shell: false,
    env: true,
    projectMetadata: true,
    aiClient: true,
  },
  resourceLimits: {
    maxMemoryMB: 1024,      // Generous limits
    maxExecutionTime: 60000,
    maxCPUTime: 30000,
    maxFileSystemOps: 5000,
    maxNetworkRequests: 500,
  },
  strictMode: false, // Allow dynamic code for debugging
}
```

### Production Mode

```typescript
{
  capabilities: {
    fileSystemRead: true,
    fileSystemWrite: false,  // Strict
    network: true,
    shell: false,             // Never in production
    env: false,
    projectMetadata: true,
    aiClient: true,
  },
  resourceLimits: {
    maxMemoryMB: 256,         // Conservative
    maxExecutionTime: 10000,
    maxCPUTime: 5000,
    maxFileSystemOps: 500,
    maxNetworkRequests: 50,
  },
  strictMode: true, // Always in production
}
```

### Untrusted Plugin

```typescript
{
  capabilities: {
    fileSystemRead: false,    // No file access
    fileSystemWrite: false,
    network: false,            // No network
    shell: false,
    env: false,
    projectMetadata: false,    // No project info
    aiClient: false,
  },
  resourceLimits: {
    maxMemoryMB: 64,          // Minimal
    maxExecutionTime: 3000,    // 3 seconds
    maxCPUTime: 1000,
    maxFileSystemOps: 0,
    maxNetworkRequests: 0,
  },
  strictMode: true,
}
```

---

## 🔗 Related Documentation

- [Plugin Security Guide](./PLUGIN_SECURITY.md) - Overall security system
- [Plugin Registry](./PLUGIN_REGISTRY.md) - Hash verification
- [Plugin Signatures](./PLUGIN_SIGNATURES.md) - Cryptographic verification
- [Security Policy](./SECURITY.md) - Security procedures

---

## 🚀 Future Enhancements

### Planned Features

- [ ] Native resource limits (when Bun adds support)
- [ ] Memory snapshots for debugging
- [ ] CPU profiling in sandbox
- [ ] Network request interception and filtering
- [ ] Syscall filtering (seccomp on Linux)
- [ ] Container-based isolation (Docker)

### Experimental Features

- WebAssembly sandboxing for near-native performance
- V8 isolates for stronger separation
- Capability delegation between plugins

---

**Last Updated:** October 5, 2025

For questions about plugin sandboxing, contact [security@rycode.ai](mailto:security@rycode.ai)
