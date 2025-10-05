# 🔒 Plugin Security Guide

## Overview

RyCode implements a **comprehensive plugin security system** to protect against malicious or compromised plugins. The system provides:

- ✅ **Plugin Allowlist** - Trust only approved plugins
- ✅ **Capability-Based Permissions** - Restrict plugin access to system resources
- ✅ **Integrity Verification** - Verify plugin code hasn't been tampered with
- ✅ **User Approval Flow** - Require confirmation before loading untrusted plugins
- ✅ **Security Audit Log** - Track all plugin loading and security events

---

## Security Modes

RyCode supports three security modes:

### 1. **Strict Mode** (Recommended for Production)

```jsonc
{
  "plugin_security": {
    "mode": "strict"
  }
}
```

- ❌ **Blocks** all plugins not in the trusted allowlist
- ✅ **Requires** explicit configuration for each plugin
- 🔒 **Best for:** Production environments, security-sensitive projects

### 2. **Warn Mode** (Default)

```jsonc
{
  "plugin_security": {
    "mode": "warn"
  }
}
```

- ⚠️ **Warns** about untrusted plugins but allows them
- 🔒 **Restricts** untrusted plugins to limited capabilities
- ✅ **Requires** user approval before installation
- 🔧 **Best for:** Development, testing new plugins

### 3. **Permissive Mode** (Not Recommended)

```jsonc
{
  "plugin_security": {
    "mode": "permissive"
  }
}
```

- ✅ **Allows** all plugins without warnings
- ⚠️ **No restrictions** on plugin capabilities
- 🚨 **Risk:** Vulnerable to malicious plugins
- 🔧 **Best for:** Testing only, never use in production

---

## Configuring Plugin Security

### Default Configuration

By default, RyCode trusts official plugins and restricts untrusted plugins:

```typescript
// Official plugins automatically trusted:
- opencode-copilot-auth@0.0.3
- opencode-anthropic-auth@0.0.2

// Untrusted plugins get limited capabilities:
- fileSystemRead: true
- fileSystemWrite: false  // ❌ No writes
- network: false          // ❌ No network access
- shell: false            // ❌ No shell commands
- env: false              // ❌ No environment variables
- aiClient: false         // ❌ No AI client access
```

### Adding Trusted Plugins

To trust a third-party plugin, add it to your `opencode.json`:

```jsonc
{
  "plugin_security": {
    "mode": "warn",  // or "strict"
    "trustedPlugins": [
      {
        "name": "my-custom-plugin",
        "versions": ["1.2.3", "^1.2.0"],  // Specific versions or semver
        "capabilities": {
          "fileSystemRead": true,
          "fileSystemWrite": true,   // Grant write access
          "network": true,            // Grant network access
          "shell": false,             // Still no shell
          "env": false,
          "projectMetadata": true,
          "aiClient": true
        }
      }
    ]
  }
}
```

### Version Matching

Supported version patterns:

```json
"versions": ["latest"]           // Any version
"versions": ["1.2.3"]            // Exact version
"versions": ["^1.2.0"]           // >= 1.2.0, < 2.0.0
"versions": ["~1.2.0"]           // >= 1.2.0, < 1.3.0
"versions": ["*"]                // Any version (not recommended)
```

---

## Plugin Capabilities

### Available Capabilities

| Capability | Description | Risk Level |
|------------|-------------|------------|
| `fileSystemRead` | Read files from project | 🟡 Medium |
| `fileSystemWrite` | Write/modify files | 🔴 High |
| `network` | Make HTTP requests | 🟡 Medium |
| `shell` | Execute shell commands | 🔴 Critical |
| `env` | Access environment variables | 🔴 High |
| `projectMetadata` | Read project info | 🟢 Low |
| `aiClient` | Access AI client | 🟡 Medium |

### Capability Examples

#### Safe Plugin (Read-Only)

```json
{
  "name": "code-analyzer",
  "capabilities": {
    "fileSystemRead": true,
    "projectMetadata": true,
    "aiClient": false,
    "network": false,
    "shell": false,
    "env": false
  }
}
```

#### Moderate Risk (Network Access)

```json
{
  "name": "api-integration",
  "capabilities": {
    "fileSystemRead": true,
    "network": true,           // Can make external requests
    "aiClient": true,
    "fileSystemWrite": false,  // Cannot modify files
    "shell": false,
    "env": false
  }
}
```

#### High Risk (Full Access)

```json
{
  "name": "build-tool",
  "capabilities": {
    "fileSystemRead": true,
    "fileSystemWrite": true,   // Can modify files
    "network": true,
    "shell": true,             // Can run commands
    "env": true,               // Can access secrets
    "aiClient": true
  }
}
```

⚠️ **Warning:** Only grant full capabilities to plugins you completely trust!

---

## Integrity Verification

### Enabling Hash Verification

To verify a plugin hasn't been tampered with, add its SHA-256 hash:

```jsonc
{
  "plugin_security": {
    "verifyIntegrity": true,  // Enable integrity checks
    "trustedPlugins": [
      {
        "name": "my-plugin",
        "versions": ["1.0.0"],
        "hash": "a3c5f8d9e2b1... (SHA-256 hash)",
        "capabilities": { /* ... */ }
      }
    ]
  }
}
```

### Generating Plugin Hashes

Use RyCode CLI to generate hashes:

```bash
# Install plugin first
bun add my-plugin

# Generate hash
rycode plugin:hash my-plugin

# Output:
# SHA-256: a3c5f8d9e2b1...
```

Or programmatically:

```typescript
import { Plugin } from "./src/plugin"

const hash = await Plugin.generatePluginHash("/path/to/plugin")
console.log("SHA-256:", hash)
```

---

## User Approval Flow

When a user installs an untrusted plugin, RyCode prompts for approval:

```
╔═══════════════════════════════════════════════════════════╗
║  ⚠️  Install Untrusted Plugin                             ║
╠═══════════════════════════════════════════════════════════╣
║                                                           ║
║  Plugin: my-custom-plugin@1.2.3                          ║
║  Status: NOT IN TRUSTED ALLOWLIST                        ║
║                                                           ║
║  This plugin will have LIMITED capabilities:             ║
║    ✅ Read project files                                  ║
║    ✅ Access project metadata                             ║
║    ❌ Cannot write files                                  ║
║    ❌ Cannot access network                               ║
║    ❌ Cannot execute shell commands                       ║
║    ❌ Cannot access environment variables                 ║
║                                                           ║
║  Do you want to proceed?                                 ║
║                                                           ║
║  [Yes]  [No]  [Show Details]                             ║
╚═══════════════════════════════════════════════════════════╝
```

To bypass approval (not recommended):

```json
{
  "plugin_security": {
    "requireApproval": false
  }
}
```

---

## Security Audit Log

RyCode logs all plugin security events for monitoring:

```typescript
import { Plugin } from "./src/plugin"

// View audit log
const auditLog = Plugin.getSecurityAuditLog()

console.log(auditLog)
// Output:
// [
//   {
//     timestamp: 1704099600000,
//     plugin: "opencode-copilot-auth",
//     version: "0.0.3",
//     action: "loaded",
//     trusted: true,
//     capabilities: { ... }
//   },
//   {
//     timestamp: 1704099601000,
//     plugin: "untrusted-plugin",
//     version: "1.0.0",
//     action: "denied",
//     trusted: false,
//     reason: "user_denied"
//   }
// ]
```

### Audit Events

| Action | Description |
|--------|-------------|
| `loaded` | Plugin successfully loaded |
| `denied` | Plugin blocked or user denied |
| `capability_check` | Capability permission checked |

---

## Best Practices

### 1. **Use Strict Mode in Production**

```json
{
  "plugin_security": {
    "mode": "strict",
    "verifyIntegrity": true
  }
}
```

### 2. **Minimize Plugin Capabilities**

Only grant the minimum permissions needed:

```json
// ❌ Don't do this
{
  "capabilities": {
    "shell": true,  // Too permissive!
    "env": true
  }
}

// ✅ Do this
{
  "capabilities": {
    "fileSystemRead": true,
    "projectMetadata": true
  }
}
```

### 3. **Verify Plugin Sources**

Before adding a plugin to your allowlist:

- ✅ Check npm download counts
- ✅ Review GitHub repository
- ✅ Read source code
- ✅ Check for known vulnerabilities
- ✅ Verify maintainer identity

### 4. **Pin Plugin Versions**

```json
// ❌ Risky
"versions": ["latest"]

// ✅ Safe
"versions": ["1.2.3"]
```

### 5. **Enable Integrity Verification**

For production, always verify plugin integrity:

```json
{
  "plugin_security": {
    "verifyIntegrity": true,
    "trustedPlugins": [
      {
        "name": "my-plugin",
        "versions": ["1.0.0"],
        "hash": "a3c5f8d9e2b1..."  // Required!
      }
    ]
  }
}
```

### 6. **Monitor Audit Logs**

Regularly review security events:

```bash
# Export audit log
rycode plugin:audit > plugin-audit.json

# Look for:
# - Denied plugins
# - Failed integrity checks
# - Unexpected capability requests
```

---

## Advanced Configuration

### Per-Environment Configuration

```jsonc
// .opencode/opencode.json (production)
{
  "plugin_security": {
    "mode": "strict",
    "verifyIntegrity": true,
    "requireApproval": true
  }
}

// opencode.json (development)
{
  "plugin_security": {
    "mode": "warn",
    "verifyIntegrity": false,
    "requireApproval": true
  }
}
```

### Default Capabilities for Untrusted Plugins

Override default capabilities:

```json
{
  "plugin_security": {
    "defaultCapabilities": {
      "fileSystemRead": true,
      "projectMetadata": true,
      "aiClient": true,        // Grant AI access by default
      "network": false,
      "shell": false,
      "env": false,
      "fileSystemWrite": false
    }
  }
}
```

---

## Sandboxing Implementation

RyCode uses capability-based sandboxing:

```typescript
// Plugin receives sandboxed input
const sandboxedInput = {
  // Available based on capabilities
  client: capabilities.aiClient ? client : undefined,
  project: capabilities.projectMetadata ? project : undefined,
  worktree: capabilities.fileSystemRead ? worktree : undefined,

  // Shell access throws error if not permitted
  $: capabilities.shell ? Bun.$ : new Proxy({}, {
    get() {
      throw new Error("Plugin does not have shell permission")
    }
  })
}
```

### Capability Enforcement

When a plugin attempts restricted operations:

```typescript
// Plugin tries to execute shell command without permission
try {
  await $.`echo "hello"`
} catch (error) {
  // CapabilityDeniedError: Plugin "my-plugin" does not have shell permission
}
```

---

## Threat Scenarios & Mitigations

### 1. Malicious Plugin Upload

**Scenario:** Attacker uploads malicious package to npm

**Mitigations:**
- ✅ Allowlist prevents loading
- ✅ User approval required
- ✅ Restricted capabilities limit damage
- ✅ Audit log tracks attempts

### 2. Dependency Confusion

**Scenario:** Attacker publishes package with same name

**Mitigations:**
- ✅ Version pinning prevents unexpected updates
- ✅ Integrity verification detects tampering
- ✅ Strict mode blocks unknown packages

### 3. Supply Chain Attack

**Scenario:** Legitimate plugin compromised

**Mitigations:**
- ✅ Hash verification detects modifications
- ✅ Version pinning prevents auto-updates
- ✅ Capability limits reduce blast radius

### 4. Credential Theft

**Scenario:** Plugin tries to steal API keys

**Mitigations:**
- ✅ `env: false` blocks environment access
- ✅ Sandboxing prevents `process.env` access
- ✅ Network restrictions prevent exfiltration

---

## Migrating Existing Configurations

### From No Security to Warn Mode

```diff
{
+  "plugin_security": {
+    "mode": "warn",
+    "requireApproval": true
+  }
}
```

### From Warn to Strict Mode

1. **List current plugins:**

```bash
rycode plugin:list
```

2. **Add to allowlist:**

```json
{
  "plugin_security": {
    "mode": "strict",
    "trustedPlugins": [
      {
        "name": "opencode-copilot-auth",
        "versions": ["0.0.3"],
        "official": true
      }
      // Add all plugins you use
    ]
  }
}
```

3. **Generate hashes:**

```bash
rycode plugin:hash opencode-copilot-auth
```

4. **Test in development first!**

---

## Troubleshooting

### Plugin Blocked in Strict Mode

**Error:** `UntrustedPluginError: Plugin "my-plugin" is not in the trusted allowlist`

**Solution:** Add plugin to `plugin_security.trustedPlugins`

### Integrity Check Failed

**Error:** `IntegrityCheckFailedError: Hash mismatch for "my-plugin"`

**Cause:** Plugin code was modified or version changed

**Solutions:**
- Regenerate hash: `rycode plugin:hash my-plugin`
- Verify plugin wasn't tampered with
- Check version matches exactly

### Capability Denied

**Error:** `CapabilityDeniedError: Plugin does not have shell permission`

**Solution:** Grant capability in plugin configuration:

```json
{
  "capabilities": {
    "shell": true
  }
}
```

### User Approval Required

Plugins prompt for approval on first run. To avoid:

```json
{
  "plugin_security": {
    "requireApproval": false  // Not recommended
  }
}
```

Or add to allowlist.

---

## API Reference

### Configuration Schema

```typescript
interface PluginSecurityPolicy {
  /** Enforcement mode */
  mode: "strict" | "warn" | "permissive"

  /** Trusted plugins */
  trustedPlugins: TrustedPlugin[]

  /** Default capabilities for untrusted plugins */
  defaultCapabilities: Capabilities

  /** Require user approval */
  requireApproval: boolean

  /** Verify plugin integrity */
  verifyIntegrity: boolean
}

interface TrustedPlugin {
  /** Package name */
  name: string

  /** Allowed versions */
  versions: string[]

  /** Plugin capabilities */
  capabilities: Capabilities

  /** SHA-256 hash (optional) */
  hash?: string

  /** Is official plugin */
  official?: boolean
}

interface Capabilities {
  fileSystemRead: boolean
  fileSystemWrite: boolean
  network: boolean
  shell: boolean
  env: boolean
  projectMetadata: boolean
  aiClient: boolean
}
```

### CLI Commands

```bash
# List loaded plugins
rycode plugin:list

# Generate plugin hash
rycode plugin:hash <plugin-name>

# View security audit log
rycode plugin:audit

# Check if plugin is trusted
rycode plugin:check <plugin-name> <version>
```

---

## Security Contact

Found a security vulnerability in a plugin or RyCode itself?

**Report to:** security@rycode.ai

**DO NOT** create public issues for security vulnerabilities.

---

## Changelog

### v1.0.0 (2025-01-05)

- ✅ Initial plugin security system
- ✅ Capability-based permissions
- ✅ Allowlist and verification
- ✅ Integrity verification
- ✅ User approval flow
- ✅ Security audit logging

---

**Stay secure! 🔒**
