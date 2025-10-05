# RyCode Security Guide

> **Complete guide to RyCode's enterprise-grade security features**

This guide provides a comprehensive overview of RyCode's security system, helping you understand and effectively use all security features to protect your development environment.

---

## 🎯 Quick Navigation

- **[Getting Started](#-getting-started)** - Essential security setup
- **[Security Features](#-security-features)** - What protects you
- **[Configuration Guide](#%EF%B8%8F-configuration-guide)** - How to configure
- **[Best Practices](#-best-practices)** - Security recommendations
- **[Detailed Documentation](#-detailed-documentation)** - Deep dives

---

## 🚀 Getting Started

### Essential Security Setup (5 minutes)

**1. Enable Strict Mode**

```jsonc
// .rycode.json
{
  "plugin_security": {
    "mode": "strict",  // Reject untrusted plugins
    "requireApproval": true
  }
}
```

**2. Verify Official Plugins**

```bash
# Check if a plugin is trusted
rycode plugin:check opencode-auth 1.0.0
```

**3. Enable Integrity Verification**

```jsonc
{
  "plugin_security": {
    "verifyIntegrity": true  // Check hashes
  }
}
```

**That's it!** You're now protected by RyCode's security system.

---

## 🔒 Security Features

RyCode implements a **defense-in-depth** strategy with three layers:

### Layer 1: Plugin Registry (Hash Verification)

**What it does:** Verifies plugins haven't been tampered with

```bash
# Add plugin to registry
rycode plugin:registry:add my-plugin.js my-plugin 1.0.0

# Verify against registry
rycode plugin:verify my-plugin.js --hash <hash>
```

📖 **[Plugin Registry Guide](./PLUGIN_REGISTRY.md)**

### Layer 2: Cryptographic Signatures

**What it does:** Proves who created the plugin

```bash
# Sign your plugin
rycode plugin:sign my-plugin.js --key private-key.pem

# Verify signature
rycode plugin:verify-sig my-plugin.js signature.json --public-key public-key.pem
```

📖 **[Plugin Signatures Guide](./PLUGIN_SIGNATURES.md)**

### Layer 3: Worker Thread Sandboxing

**What it does:** Isolates plugin execution with resource limits

```typescript
const sandbox = await PluginSandbox.createSandbox({
  pluginName: "my-plugin",
  pluginVersion: "1.0.0",
  capabilities: {
    fileSystemWrite: false,  // Read-only
    network: true,
    shell: false,  // No shell access
  },
  resourceLimits: {
    maxMemoryMB: 256,
    maxExecutionTime: 10000,
  },
})
```

📖 **[Plugin Sandboxing Guide](./PLUGIN_SANDBOXING.md)**

---

## ⚙️ Configuration Guide

### Security Levels

Choose based on your risk tolerance:

#### 🔴 **Strict (Recommended for Production)**

```jsonc
{
  "plugin_security": {
    "mode": "strict",
    "verifyIntegrity": true,
    "requireApproval": true,
    "trustedPlugins": [
      // Only list explicitly trusted plugins
    ],
    "signature_policy": {
      "requireSignatures": true,
      "allowSelfSigned": false
    }
  }
}
```

**Use when:** Production environments, handling sensitive data

#### 🟡 **Warn (Default)**

```jsonc
{
  "plugin_security": {
    "mode": "warn",  // Warn but allow
    "verifyIntegrity": true,
    "requireApproval": true
  }
}
```

**Use when:** Development, testing, learning

#### 🟢 **Permissive (Development Only)**

```jsonc
{
  "plugin_security": {
    "mode": "permissive",
    "verifyIntegrity": false,
    "requireApproval": false
  }
}
```

**Use when:** Local development, plugin creation

⚠️ **Never use permissive mode in production!**

---

## 🛡️ Best Practices

### For All Users

1. **Keep RyCode Updated**
   ```bash
   rycode upgrade
   ```

2. **Review Plugin Permissions**
   ```bash
   rycode plugin:check <plugin-name> <version>
   ```

3. **Monitor Security Audit Log**
   ```bash
   rycode plugin:audit
   ```

4. **Use Official Plugins When Possible**
   - Look for ✓ (official) verification
   - Check plugin:check output for "Official: Yes"

### For Plugin Users

1. **Verify Before Installing**
   ```bash
   # Check trust status
   rycode plugin:check my-plugin 1.0.0

   # Verify hash
   rycode plugin:verify plugin.js --hash <expected-hash>

   # Verify signature
   rycode plugin:verify-sig plugin.js sig.json --public-key key.pem
   ```

2. **Pin Plugin Versions**
   ```jsonc
   {
     "plugin_security": {
       "trustedPlugins": [{
         "name": "my-plugin",
         "versions": ["1.2.3"],  // Exact version
         "hash": "abc123..."
       }]
     }
   }
   ```

3. **Grant Minimal Capabilities**
   ```jsonc
   {
     "trustedPlugins": [{
       "name": "my-plugin",
       "capabilities": {
         "fileSystemRead": true,
         "fileSystemWrite": false,  // Deny write
         "shell": false,  // Deny shell
       }
     }]
   }
   ```

### For Plugin Developers

1. **Sign Your Releases**
   ```bash
   # Generate signing key
   rycode plugin:keygen

   # Sign plugin
   rycode plugin:sign dist/plugin.js --key private-key.pem

   # Include signature in npm package
   cp plugin.sig.json dist/
   ```

2. **Document Capabilities**
   ```markdown
   ## Required Capabilities
   - `fileSystemRead`: Read project files
   - `network`: Fetch external resources
   - `aiClient`: Generate AI completions
   ```

3. **Test with Strict Limits**
   ```typescript
   // Test with production-like sandbox
   const sandbox = await PluginSandbox.createSandbox({
     resourceLimits: {
       maxMemoryMB: 128,  // Low limit
       maxExecutionTime: 5000,  // 5 seconds
     }
   })
   ```

4. **Submit to Registry**
   ```bash
   rycode plugin:registry:add dist/plugin.js my-plugin 1.0.0 \
     --description "My awesome plugin" \
     --verified-by community
   ```

### For Organizations

1. **Corporate Registry**
   ```bash
   # Set up corporate registry
   export RYCODE_REGISTRY="https://registry.company.com"

   # Sync all workstations
   rycode plugin:registry:sync --url $RYCODE_REGISTRY
   ```

2. **Centralized Signing**
   ```bash
   # Store company signing key in secrets manager
   # Sign all internal plugins with company key
   # Distribute public key to all developers
   ```

3. **Security Policy Enforcement**
   ```jsonc
   // Enforced via .rycode.json in repos
   {
     "plugin_security": {
       "mode": "strict",
       "requireSignatures": true,
       "trustedPlugins": [
         // Only company-approved plugins
       ]
     }
   }
   ```

---

## 🎯 Decision Trees

### Should I Trust This Plugin?

```
Is it an official RyCode plugin?
├─ Yes → ✅ Safe to use
└─ No
    ├─ Is it in the registry?
    │   ├─ Yes, verified by "official" → ✅ Safe
    │   ├─ Yes, verified by "community" → ⚠️ Review first
    │   └─ No → 🔴 High risk
    └─ Does it have a valid signature?
        ├─ Yes, from trusted signer → ✅ Safe
        ├─ Yes, but unknown signer → ⚠️ Verify signer
        └─ No signature → 🔴 Don't use
```

### What Security Level Should I Use?

```
What environment?
├─ Production → Use "strict" mode
├─ Staging → Use "strict" or "warn"
├─ Development → Use "warn"
└─ Local plugin development → Use "permissive"

Handling sensitive data?
├─ Yes → Use "strict" + requireSignatures
└─ No → Use "warn"

Team size?
├─ Enterprise → Corporate registry + strict mode
├─ Small team → Shared registry + warn mode
└─ Solo developer → Local registry + warn mode
```

---

## 🔍 Troubleshooting

### Common Issues

**Problem:** "Plugin is not trusted"

**Solution:**
```bash
# Check why
rycode plugin:check my-plugin 1.0.0

# Add to trusted list
# Edit .rycode.json and add to trustedPlugins
```

**Problem:** "Integrity check failed"

**Solution:**
```bash
# Verify hash
rycode plugin:hash plugin.js

# Compare with expected
# If different, plugin was modified - re-download
```

**Problem:** "Signature verification failed"

**Solution:**
```bash
# Verify you have the correct public key
# Re-download signature file
# Check for file corruption
```

**Problem:** "Sandbox timeout"

**Solution:**
```typescript
// Increase timeout
resourceLimits: {
  maxExecutionTime: 60000,  // 1 minute
}
```

---

## 📊 Security Checklist

### Before Installing a Plugin

- [ ] Check if plugin is in registry
- [ ] Verify hash matches registry
- [ ] Check signature if available
- [ ] Review required capabilities
- [ ] Check plugin reputation/reviews
- [ ] Read plugin source code (if open source)

### Before Deploying to Production

- [ ] Use strict security mode
- [ ] Enable signature verification
- [ ] Pin exact plugin versions
- [ ] Set conservative resource limits
- [ ] Enable audit logging
- [ ] Review all trusted plugins
- [ ] Test security policies
- [ ] Document security decisions

### Regular Maintenance

- [ ] Review audit logs weekly
- [ ] Update RyCode monthly
- [ ] Rotate signing keys annually
- [ ] Review trusted plugins quarterly
- [ ] Update security policies as needed
- [ ] Monitor for security advisories

---

## 📚 Detailed Documentation

### Core Security Docs

- **[SECURITY.md](./SECURITY.md)** - Security policy and procedures
- **[SECURITY_ASSESSMENT.md](./SECURITY_ASSESSMENT.md)** - Security audit report

### Feature Guides

- **[PLUGIN_SECURITY.md](./PLUGIN_SECURITY.md)** - Plugin security system overview
- **[PLUGIN_REGISTRY.md](./PLUGIN_REGISTRY.md)** - Registry and hash verification
- **[PLUGIN_SIGNATURES.md](./PLUGIN_SIGNATURES.md)** - Cryptographic signatures
- **[PLUGIN_SANDBOXING.md](./PLUGIN_SANDBOXING.md)** - Worker thread isolation

### Technical Details

- **[PLUGIN_SECURITY_IMPLEMENTATION.md](./PLUGIN_SECURITY_IMPLEMENTATION.md)** - Implementation details

---

## 🆘 Getting Help

### Security Issues

**Found a security vulnerability?**

🔒 **DO NOT** create a public GitHub issue

✅ **DO** email: [security@rycode.ai](mailto:security@rycode.ai)

See [SECURITY.md](./SECURITY.md#-reporting-security-vulnerabilities) for details.

### General Questions

- **GitHub Discussions:** Ask questions
- **Documentation:** Check docs first
- **Examples:** See `examples/` directory

---

## 📈 Security Roadmap

### ✅ Implemented (October 2025)

- [x] Plugin registry with hash verification
- [x] Cryptographic signature verification (GPG/RSA)
- [x] Worker thread sandboxing

### 🔄 In Progress

- [ ] Automated CVE scanning
- [ ] Rate limiting for API requests
- [ ] Network request filtering

### 🔮 Planned

- [ ] Container-based isolation (Docker)
- [ ] WebAssembly sandboxing
- [ ] Hardware security module (HSM) support
- [ ] Plugin certification program

---

## 📊 Compliance

RyCode security features help meet:

- ✅ **OWASP Top 10** - Addressed 8/10
- ✅ **CWE Top 25** - Mitigated common weaknesses
- ✅ **NIST Cybersecurity Framework** - Supply chain risk management
- ✅ **SOC 2 Type II** - In progress
- ✅ **ISO 27001** - Aligned controls

---

## 🏆 Summary

RyCode provides **enterprise-grade plugin security** through:

1. **Hash Verification** - Detect tampering
2. **Cryptographic Signatures** - Verify authenticity
3. **Process Sandboxing** - Isolate execution

**Result:** Multiple layers of protection ensuring safe plugin usage.

**Recommendation:** Start with strict mode, verify all plugins, and monitor audit logs.

---

**Last Updated:** October 5, 2025

**Version:** 1.0.0

For questions or feedback, contact [security@rycode.ai](mailto:security@rycode.ai)
