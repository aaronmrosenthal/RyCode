# Security Policy

## 🔒 Reporting Security Vulnerabilities

We take the security of RyCode seriously. If you discover a security vulnerability, please follow responsible disclosure practices.

### **DO NOT** Create Public Issues

**Please do not** report security vulnerabilities through public GitHub issues, discussions, or pull requests.

### How to Report

**Send vulnerability reports to:** [security@rycode.ai](mailto:security@rycode.ai)

**Include in your report:**

1. **Description** - Clear explanation of the vulnerability
2. **Impact** - Potential security impact and affected components
3. **Steps to Reproduce** - Detailed steps to reproduce the issue
4. **Proof of Concept** - Code or commands demonstrating the vulnerability (if applicable)
5. **Suggested Fix** - Your recommendation for addressing the issue (optional)
6. **Disclosure Timeline** - Your preferred timeline for public disclosure

### What to Expect

- **Initial Response:** Within 48 hours
- **Status Update:** Within 7 days with assessment and planned fix timeline
- **Fix Timeline:** Critical issues addressed within 30 days, others within 90 days
- **Credit:** We will acknowledge your contribution in our security advisories (unless you prefer to remain anonymous)

---

## 🛡️ Supported Versions

We release security updates for the following versions:

| Version | Status | Security Support |
|---------|--------|------------------|
| 1.x.x   | ✅ Stable | Full support |
| 0.x.x   | ⚠️ Beta | Best effort |
| < 0.1.x | ❌ Legacy | No support |

**Recommendation:** Always use the latest stable version for the best security posture.

---

## 🔐 Security Features

RyCode implements enterprise-grade security controls:

### Plugin Security

- **✅ Allowlist System** - Only trusted plugins can load
- **✅ Capability-Based Permissions** - 7 granular permission types
- **✅ Integrity Verification** - SHA-256 hash checking
- **✅ Security Audit Logging** - Complete activity tracking
- **✅ User Approval Flow** - Confirm untrusted plugin installations

See [PLUGIN_SECURITY.md](./PLUGIN_SECURITY.md) for complete documentation.

### Authentication & Authorization

- **✅ API Key Authentication** - Secure server access
- **✅ Timing-Attack Prevention** - Constant-time comparisons
- **✅ scrypt-based Hashing** - Strong key derivation
- **✅ Header-Only Auth** - Prevents credential logging

### File System Security

- **✅ Path Traversal Prevention** - Validates all file paths
- **✅ Sensitive File Protection** - Blocks access to credentials, SSH keys, etc.
- **✅ Directory Restrictions** - Access limited to project directories

### Command Execution Security

- **✅ Permission System** - User approval for dangerous commands
- **✅ Path Validation** - Prevents operations outside project
- **✅ Timeout Protection** - Prevents runaway processes
- **✅ Output Limiting** - Prevents memory exhaustion

---

## 🎯 Security Best Practices

### For Users

1. **Use Strict Mode in Production**
   ```jsonc
   {
     "plugin_security": {
       "mode": "strict",
       "verifyIntegrity": true
     }
   }
   ```

2. **Use CLI Commands for Plugin Management**
   ```bash
   # Check if a plugin is trusted
   rycode plugin:check my-plugin 1.0.0

   # Generate hash for a plugin
   rycode plugin:hash /path/to/plugin.js

   # Verify plugin integrity
   rycode plugin:verify /path/to/plugin.js --hash <expected-hash>

   # View security audit log
   rycode plugin:audit
   ```

2. **Pin Plugin Versions**
   ```jsonc
   {
     "plugin_security": {
       "trustedPlugins": [{
         "name": "my-plugin",
         "versions": ["1.2.3"]  // Exact version
       }]
     }
   }
   ```

3. **Enable Integrity Verification**
   ```jsonc
   {
     "plugin_security": {
       "verifyIntegrity": true,
       "trustedPlugins": [{
         "hash": "sha256-hash-here"
       }]
     }
   }
   ```

4. **Review Plugin Capabilities**
   - Only grant minimum required permissions
   - Avoid `shell: true` and `env: true` unless absolutely necessary

5. **Monitor Audit Logs**
   ```bash
   rycode plugin:audit
   ```

### For Plugin Developers

1. **Minimize Capabilities**
   - Request only the permissions your plugin needs
   - Document why each capability is required

2. **Validate All Inputs**
   - Never trust user-provided data
   - Sanitize file paths, URLs, and commands

3. **Follow Secure Coding Practices**
   - Avoid eval(), Function(), or dynamic code execution
   - Use parameterized commands, not string concatenation
   - Validate and sanitize all outputs

4. **Publish Source Code**
   - Make your plugin open source for transparency
   - Accept security audits and bug reports

5. **Sign Your Releases**
   - Provide GPG signatures for your packages
   - Include SHA-256 hashes in release notes

---

## 🚨 Known Security Considerations

### Current Limitations

1. **Plugin Sandboxing**
   - Current implementation uses capability-based restrictions
   - Full process isolation not yet implemented
   - **Mitigation:** Use strict mode and only trust verified plugins

2. **Dependency Security**
   - Plugins may include vulnerable dependencies
   - **Mitigation:** Regularly run `bun audit` and update dependencies

3. **AI Prompt Injection**
   - Malicious files could contain prompts that trick the AI
   - **Mitigation:** Review AI-generated commands before execution

### Planned Improvements

- [x] **Process-level plugin sandboxing using worker threads** ✅ **IMPLEMENTED** (October 2025)
- [ ] Automated CVE scanning for plugin dependencies
- [x] **GPG signature verification for plugins** ✅ **IMPLEMENTED** (October 2025)
- [x] **Plugin registry with verified hashes** ✅ **IMPLEMENTED** (October 2025)
- [ ] Rate limiting for API requests
- [ ] Network request filtering and monitoring

**Progress: 3 of 6 completed (50%)**

---

## 📋 Security Compliance

RyCode is designed with compliance in mind:

### OWASP Top 10 (2021)

| Risk | Status | Implementation |
|------|--------|----------------|
| A01: Broken Access Control | ✅ Addressed | Capability-based permissions, path validation |
| A02: Cryptographic Failures | ✅ Addressed | scrypt hashing, secure credential storage |
| A03: Injection | ⚠️ Partial | Command validation, needs further hardening |
| A04: Insecure Design | ✅ Addressed | Security-first architecture |
| A05: Security Misconfiguration | ✅ Addressed | Secure defaults, strict mode |
| A06: Vulnerable Components | ⚠️ Monitor | Dependency auditing recommended |
| A07: Authentication Failures | ✅ Addressed | Timing-safe comparisons, strong hashing |
| A08: Software & Data Integrity | ✅ Addressed | Plugin integrity verification |
| A09: Logging Failures | ✅ Addressed | Comprehensive security audit logs |
| A10: SSRF | ✅ Addressed | No user-controlled URL fetching |

### Standards & Frameworks

- **NIST Cybersecurity Framework** - Supply chain risk management
- **CWE Top 25** - Addressed common weaknesses
- **SANS Top 25** - Implemented security controls

### Future Compliance Goals

- SOC 2 Type II certification
- ISO 27001 compliance
- GDPR compliance for data handling
- PCI DSS for payment processing (if applicable)

---

## 🔍 Security Audits

### Internal Audits

- **Last Audit:** October 5, 2025
- **Findings:** Comprehensive security assessment completed
- **Status:** 4/5 security rating (Good → Enterprise-grade)
- **Report:** [SECURITY_ASSESSMENT.md](./SECURITY_ASSESSMENT.md)

### External Audits

- **Status:** Not yet conducted
- **Planned:** Q2 2025
- **Scope:** Full penetration testing and code review

---

## 🏆 Security Hall of Fame

We recognize and thank security researchers who help keep RyCode secure:

<!--
Add researchers who responsibly disclose vulnerabilities:

### 2025
- **[Researcher Name]** - Discovered [vulnerability type]
- **[Researcher Name]** - Reported [security issue]
-->

*Be the first to contribute!*

---

## 📚 Security Resources

### Documentation

- **[Security Guide](./SECURITY_GUIDE.md)** - Complete user-facing security guide
- **[Plugin Security](./PLUGIN_SECURITY.md)** - Plugin security system overview
- **[Plugin Registry](./PLUGIN_REGISTRY.md)** - Hash verification and discovery
- **[Plugin Signatures](./PLUGIN_SIGNATURES.md)** - Cryptographic signatures (GPG/RSA)
- **[Plugin Sandboxing](./PLUGIN_SANDBOXING.md)** - Worker thread isolation
- **[Security Assessment](./SECURITY_ASSESSMENT.md)** - Security audit report
- **[Plugin Implementation](./PLUGIN_SECURITY_IMPLEMENTATION.md)** - Technical details

### External Resources

- [OWASP Application Security](https://owasp.org/)
- [Node.js Security Best Practices](https://nodejs.org/en/docs/guides/security/)
- [Bun Security Model](https://bun.sh/docs/runtime/security)
- [CWE Top 25](https://cwe.mitre.org/top25/)

### Tools

- **Dependency Auditing:** `bun audit`
- **Static Analysis:** TypeScript strict mode
- **Linting:** ESLint with security plugins
- **Testing:** Comprehensive security test suite

---

## 🔐 Security Contact

**Primary Contact:** [security@rycode.ai](mailto:security@rycode.ai)

**PGP Key:** Coming soon

**Response Time:**
- Critical vulnerabilities: 24-48 hours
- High severity: 2-7 days
- Medium/Low severity: 7-14 days

### Alternative Channels

If you cannot reach us via email:

1. **Encrypted Communication:** Use our PGP key (coming soon)
2. **GitHub Security Advisory:** [Create a private security advisory](https://github.com/rycode/opencode/security/advisories/new)
3. **Emergency Contact:** For critical zero-day vulnerabilities, contact maintainers directly

---

## ⚖️ Disclosure Policy

We follow a **coordinated disclosure** model:

### Timeline

1. **Day 0:** Receive vulnerability report
2. **Day 0-2:** Acknowledge receipt and initial assessment
3. **Day 2-7:** Validate and reproduce the vulnerability
4. **Day 7-30:** Develop and test fix
5. **Day 30-90:** Release patch and publish advisory
6. **Day 90+:** Public disclosure (if not resolved)

### Public Disclosure

- We will coordinate with you on disclosure timing
- Credit will be given (unless anonymous preferred)
- CVE IDs will be assigned for significant vulnerabilities
- Security advisories published on GitHub

---

## 📜 Legal

### Responsible Disclosure

We commit to:
- Not pursuing legal action against security researchers who follow responsible disclosure
- Working with you to understand and address the issue
- Keeping you informed throughout the remediation process

### Safe Harbor

Security research conducted in good faith will not be considered:
- A violation of our terms of service
- A violation of computer fraud laws
- A breach of contract

**Conditions:**
- You must report the vulnerability privately
- You must not exploit the vulnerability beyond demonstration
- You must not access, modify, or delete user data
- You must give us reasonable time to fix the issue before public disclosure

---

## 🙏 Thank You

Thank you for helping keep RyCode and our users safe!

Your efforts to responsibly disclose security vulnerabilities are greatly appreciated and help make RyCode more secure for everyone.

---

## 📝 Version History

- **v1.0.0** (2025-01-05) - Initial security policy
  - Added responsible disclosure process
  - Documented security features
  - Established supported versions
  - Created security hall of fame

---

**Last Updated:** October 5, 2025

For questions about this security policy, contact [security@rycode.ai](mailto:security@rycode.ai)
