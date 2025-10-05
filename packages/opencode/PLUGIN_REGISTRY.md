# Plugin Registry

> **Centralized registry of verified plugin hashes and metadata**

The Plugin Registry provides a trusted source for plugin verification and discovery, complementing the plugin security system with hash verification and automated trust management.

---

## 🎯 Overview

The Plugin Registry is a centralized database of verified plugin information including:

- **SHA-256 hashes** for integrity verification
- **Version information** for compatibility checks
- **Verification levels** (official, community, user)
- **Metadata** (description, author, homepage, repository)
- **Capabilities** required by each plugin

### Key Features

✅ **Local & Remote Registries** - Support for both local and remote registry sources
✅ **Automatic Synchronization** - Auto-sync with remote registry
✅ **Hash Verification** - Verify plugin integrity against registry
✅ **Search & Discovery** - Find plugins by name or pattern
✅ **Multiple Verification Levels** - Official, community, and user-verified plugins
✅ **CLI Management** - Full command-line interface for registry operations

---

## 📦 Installation

The registry is built into RyCode. No additional installation required.

---

## 🚀 Quick Start

### Add a Plugin to Registry

```bash
# Generate hash and add to registry
rycode plugin:hash /path/to/plugin.js
rycode plugin:registry:add /path/to/plugin.js my-plugin 1.0.0 \
  --description "My awesome plugin" \
  --author "Your Name"
```

### Verify Plugin Against Registry

```bash
# Check if plugin is in registry and hash matches
rycode plugin:verify /path/to/plugin.js --hash <expected-hash>
```

### List Registry Contents

```bash
# List all plugins in registry
rycode plugin:registry:list

# List specific plugin versions
rycode plugin:registry:list my-plugin
```

### Search for Plugins

```bash
# Search by name or description
rycode plugin:registry:search "authentication"
```

### Sync with Remote Registry

```bash
# Sync with default remote registry
rycode plugin:registry:sync

# Sync with custom URL
rycode plugin:registry:sync --url https://my-registry.com/plugins.json
```

### View Statistics

```bash
# Show registry statistics
rycode plugin:registry:stats
```

---

## 📖 CLI Reference

### `plugin:registry:add`

Add a plugin to the registry.

```bash
rycode plugin:registry:add <plugin-path> <name> <version> [options]
```

**Arguments:**
- `plugin-path` - Path to the plugin file
- `name` - Plugin package name (e.g., `opencode-auth`)
- `version` - Plugin version (e.g., `1.0.0`)

**Options:**
- `--description` - Plugin description
- `--author` - Plugin author
- `--homepage` - Homepage URL
- `--repository` - Repository URL
- `--verified-by` - Verification level: `official`, `community`, or `user` (default: `user`)
- `--json` - Output as JSON

**Example:**

```bash
rycode plugin:registry:add ./my-plugin.js opencode-my-plugin 1.2.3 \
  --description "My custom authentication plugin" \
  --author "John Doe" \
  --homepage "https://example.com" \
  --repository "https://github.com/user/plugin" \
  --verified-by community
```

---

### `plugin:registry:remove`

Remove a plugin from the registry.

```bash
rycode plugin:registry:remove <name> <version> [options]
```

**Arguments:**
- `name` - Plugin package name
- `version` - Plugin version

**Options:**
- `--json` - Output as JSON

**Example:**

```bash
rycode plugin:registry:remove opencode-my-plugin 1.2.3
```

---

### `plugin:registry:list`

List plugins in the registry.

```bash
rycode plugin:registry:list [name] [options]
```

**Arguments:**
- `name` - Optional: Filter by plugin name

**Options:**
- `--json` - Output as JSON

**Example:**

```bash
# List all plugins
rycode plugin:registry:list

# List specific plugin (all versions)
rycode plugin:registry:list opencode-auth

# Output as JSON
rycode plugin:registry:list --json
```

**Output:**

```
Plugin Registry
4 entries

✓ opencode-auth@1.0.0
  Official authentication plugin
  Hash: 3a5f8d9e2b1c4f7a...

~ community-formatter@2.1.0
  Code formatting plugin
  Hash: 8b4e2d9f1a3c7e5b...
  Author: Jane Smith

- my-custom-plugin@1.0.0
  Hash: 5c9f2e1a4b7d3e8f...
```

**Verification Icons:**
- `✓` Official plugin
- `~` Community-verified plugin
- `-` User-added plugin

---

### `plugin:registry:search`

Search for plugins by name or description.

```bash
rycode plugin:registry:search <pattern> [options]
```

**Arguments:**
- `pattern` - Search pattern (regex)

**Options:**
- `--json` - Output as JSON

**Example:**

```bash
# Search for authentication plugins
rycode plugin:registry:search "auth"

# Search for OpenCode plugins
rycode plugin:registry:search "^opencode-"
```

---

### `plugin:registry:sync`

Sync with remote registry.

```bash
rycode plugin:registry:sync [options]
```

**Options:**
- `--url` - Remote registry URL (overrides default)
- `--json` - Output as JSON

**Example:**

```bash
# Sync with default registry
rycode plugin:registry:sync

# Sync with custom registry
rycode plugin:registry:sync --url https://my-company.com/registry.json
```

---

### `plugin:registry:stats`

Show registry statistics.

```bash
rycode plugin:registry:stats [options]
```

**Options:**
- `--json` - Output as JSON

**Example:**

```bash
rycode plugin:registry:stats
```

**Output:**

```
Registry Statistics

Total Entries:   156
Unique Plugins:  42

Official:        12
Community:       85
User:            59
```

---

## 💻 Programmatic API

### Load Registry

```typescript
import { PluginRegistry } from "./src/plugin/registry"

// Load registry with default config
const registry = await PluginRegistry.load()

// Load with custom config
const registry = await PluginRegistry.load({
  localPath: "/custom/path/registry.json",
  remoteUrl: "https://my-registry.com/plugins.json",
  autoUpdate: true,
  cacheTTL: 3600000, // 1 hour
})
```

### Add Entry

```typescript
await PluginRegistry.add({
  name: "my-plugin",
  version: "1.0.0",
  hash: "abc123...", // SHA-256 hash
  description: "My plugin",
  author: "Your Name",
  homepage: "https://example.com",
  repository: "https://github.com/user/repo",
  verifiedBy: "user",
})
```

### Find Entry

```typescript
// Find specific version
const entry = await PluginRegistry.find("my-plugin", "1.0.0")

// Find all versions
const entries = await PluginRegistry.findAll("my-plugin")
```

### Verify Hash

```typescript
const { verified, entry } = await PluginRegistry.verify(
  "my-plugin",
  "1.0.0",
  "abc123..." // SHA-256 hash to verify
)

if (verified) {
  console.log(`✓ Plugin verified as ${entry.verifiedBy}`)
} else {
  console.log("✗ Verification failed")
}
```

### Search

```typescript
const results = await PluginRegistry.search("authentication")

for (const plugin of results) {
  console.log(`${plugin.name}@${plugin.version}`)
}
```

### Remove Entry

```typescript
const removed = await PluginRegistry.remove("my-plugin", "1.0.0")
```

### Statistics

```typescript
const stats = await PluginRegistry.stats()

console.log(`Total plugins: ${stats.uniquePlugins}`)
console.log(`Official: ${stats.officialCount}`)
console.log(`Community: ${stats.communityCount}`)
console.log(`User: ${stats.userCount}`)
```

---

## 🔐 Integration with Security System

The registry integrates seamlessly with the plugin security system:

### Verify Plugin Before Loading

```typescript
import { PluginSecurity } from "./src/plugin/security"

const result = await PluginSecurity.verifyWithRegistry(
  "my-plugin",
  "1.0.0",
  "/path/to/plugin.js"
)

if (result.verified) {
  console.log("✓ Plugin verified")
  console.log(`  Verified by: ${result.entry.verifiedBy}`)
  // Safe to load plugin
} else {
  console.warn("⚠ Plugin not in registry or hash mismatch")
  // Prompt user for approval
}
```

### Automatic Registry Check in Configuration

```jsonc
// .rycode.json
{
  "plugin_security": {
    "mode": "strict",
    "verifyIntegrity": true,
    "useRegistry": true,  // Enable registry verification
    "trustedPlugins": [
      {
        "name": "my-plugin",
        "versions": ["1.0.0"],
        // Hash automatically verified against registry
      }
    ]
  }
}
```

---

## 📁 Registry File Format

The registry is stored as JSON:

```json
{
  "version": "1.0.0",
  "lastUpdated": 1696521600000,
  "entries": [
    {
      "name": "opencode-auth",
      "version": "1.0.0",
      "hash": "3a5f8d9e2b1c4f7a8e5d9c2b3f6a4e7b1c8d5f9e2a6b3c7f4e8d1a5b9c2e6f3a",
      "description": "Official authentication plugin",
      "author": "RyCode Team",
      "homepage": "https://rycode.ai/plugins/auth",
      "repository": "https://github.com/rycode/opencode-auth",
      "verifiedBy": "official",
      "timestamp": 1696521600000,
      "capabilities": {
        "fileSystemRead": true,
        "fileSystemWrite": false,
        "network": true,
        "shell": false,
        "env": true,
        "projectMetadata": true,
        "aiClient": true
      }
    }
  ]
}
```

### Default Locations

- **Local Registry:** `~/.rycode/plugin-registry.json`
- **Remote Registry:** `https://registry.rycode.ai/plugins.json`

---

## 🌐 Remote Registry

### Publishing to Remote Registry

For official RyCode plugins:

1. **Generate hash:**
   ```bash
   rycode plugin:hash /path/to/plugin.js
   ```

2. **Submit to registry:** Create PR to https://github.com/rycode/plugin-registry

3. **Verification process:**
   - Code review
   - Security audit
   - Automated tests
   - Manual approval

4. **Publication:** Merged to `main` → Auto-deployed to `registry.rycode.ai`

### Self-Hosted Registry

Host your own registry:

```typescript
// registry-server.ts
import { serve } from "bun"

const registry = {
  version: "1.0.0",
  lastUpdated: Date.now(),
  entries: [
    // Your plugins here
  ]
}

serve({
  port: 3000,
  fetch(req) {
    if (req.url.endsWith("/plugins.json")) {
      return new Response(JSON.stringify(registry), {
        headers: { "Content-Type": "application/json" }
      })
    }
    return new Response("Not Found", { status: 404 })
  }
})
```

Configure clients:

```bash
rycode plugin:registry:sync --url https://your-company.com:3000/plugins.json
```

---

## 🔒 Security Best Practices

### For Plugin Users

1. **Always verify plugins:**
   ```bash
   rycode plugin:verify /path/to/plugin.js --hash <expected-hash>
   ```

2. **Prefer official plugins:**
   - Look for `✓` (official) verification icon
   - Official plugins are audited by RyCode team

3. **Check verification level:**
   - `official` - Fully audited and maintained by RyCode
   - `community` - Reviewed and verified by community
   - `user` - User-submitted, use with caution

4. **Keep registry synced:**
   ```bash
   rycode plugin:registry:sync
   ```

### For Plugin Publishers

1. **Publish source code:**
   - Make your plugin open source
   - Link to public GitHub repository

2. **Document capabilities:**
   - Clearly state what permissions your plugin needs
   - Explain why each capability is required

3. **Sign your releases:**
   - Provide GPG signatures
   - Include SHA-256 hashes in release notes

4. **Follow semantic versioning:**
   - Use semver for version numbers
   - Document breaking changes

---

## 🛠️ Troubleshooting

### Registry Not Syncing

```bash
# Clear cache and force reload
PluginRegistry.clearCache()
rycode plugin:registry:sync
```

### Hash Mismatch

If verification fails:

1. **Check file integrity:**
   ```bash
   rycode plugin:hash /path/to/plugin.js
   ```

2. **Compare with registry:**
   ```bash
   rycode plugin:registry:list plugin-name --json
   ```

3. **Re-download plugin if tampered**

### Registry File Corrupted

```bash
# Backup current registry
cp ~/.rycode/plugin-registry.json ~/.rycode/plugin-registry.backup.json

# Re-sync from remote
rm ~/.rycode/plugin-registry.json
rycode plugin:registry:sync
```

---

## 📊 Example Workflows

### Workflow 1: Install Verified Plugin

```bash
# 1. Search for plugin
rycode plugin:registry:search "formatter"

# 2. Check trust status
rycode plugin:check opencode-formatter 2.1.0

# 3. Install via npm/bun
bun add opencode-formatter

# 4. Verify installed plugin
PLUGIN_PATH=$(bun pm ls opencode-formatter | grep "opencode-formatter" | awk '{print $NF}')
rycode plugin:verify $PLUGIN_PATH/index.js --hash <hash-from-registry>

# 5. Add to trusted plugins
# Edit .rycode.json to add plugin
```

### Workflow 2: Publish Custom Plugin

```bash
# 1. Build your plugin
bun build ./src/index.ts --outfile ./dist/plugin.js

# 2. Generate hash
rycode plugin:hash ./dist/plugin.js

# 3. Add to local registry
rycode plugin:registry:add ./dist/plugin.js my-company-plugin 1.0.0 \
  --description "Internal company plugin" \
  --author "Engineering Team" \
  --verified-by user

# 4. Share hash with team
rycode plugin:registry:list my-company-plugin --json > plugin-info.json

# 5. Team members can verify
rycode plugin:verify ./their-copy/plugin.js --hash <your-hash>
```

### Workflow 3: Corporate Registry

```bash
# 1. Set up corporate registry
export CORPORATE_REGISTRY="https://registry.mycorp.com/plugins.json"

# 2. Configure all developer machines
rycode plugin:registry:sync --url $CORPORATE_REGISTRY

# 3. Add corporate plugins
rycode plugin:registry:add ./corp-plugin.js corp-auth 1.0.0 \
  --verified-by official \
  --repository "https://github.mycorp.com/plugins/auth"

# 4. Auto-verify in CI/CD
npm run verify-plugins  # Custom script using registry API
```

---

## 🔗 Related Documentation

- [Plugin Security Guide](./PLUGIN_SECURITY.md) - Complete security system documentation
- [Security Policy](./SECURITY.md) - Overall security policies and procedures
- [Plugin Development](./docs/PLUGIN_DEVELOPMENT.md) - How to create plugins

---

## 📝 FAQ

### Q: Is the registry required?

A: No, the registry is optional. You can still use plugins without it, but the registry provides automated hash verification and discovery.

### Q: Can I use multiple registries?

A: Currently, one registry at a time. You can switch by changing the `remoteUrl` configuration.

### Q: How often is the remote registry updated?

A: The official registry is updated continuously. Local cache refreshes every 1 hour by default (configurable with `cacheTTL`).

### Q: Can I trust community-verified plugins?

A: Community plugins are reviewed by trusted community members but not officially audited. Review the plugin code yourself for critical applications.

### Q: What happens if a plugin is removed from the registry?

A: Your local cache will retain the entry until next sync. After sync, the plugin will no longer be verified but can still be used if already trusted in your config.

---

**Last Updated:** October 5, 2025

For questions about the plugin registry, contact [security@rycode.ai](mailto:security@rycode.ai)
