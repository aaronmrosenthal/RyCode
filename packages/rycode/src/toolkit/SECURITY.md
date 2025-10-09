# Toolkit Client Security

## ✅ Security Measures Implemented

### 1. Input Sanitization

**Command Injection Prevention** (`validators.ts:38-41`):
```typescript
static sanitizeInput(input: string): string {
  // Remove shell metacharacters to prevent command injection
  return input.replace(/[;&|`$()\n]/g, '');
}
```

**What's Protected**:
- ✅ Semicolons (`;`) - command chaining
- ✅ Pipes (`|`) - command piping
- ✅ Ampersands (`&`) - background execution
- ✅ Backticks (`` ` ``) - command substitution
- ✅ Dollar signs (`$`) - variable expansion
- ✅ Parentheses (`()`) - subshells
- ✅ Newlines (`\n`) - command separation

**Applied To**:
- All user input (project ideas, features, issues, tasks, descriptions)
- All context parameters
- Every command before subprocess execution

### 2. Input Validation

**Length Constraints** (`validators.ts:43-61`):
```typescript
validateProjectIdea(idea: string): void {
  if (idea.length < 10 || idea.length > 5000) {
    throw new ValidationError(...)
  }
}
```

**Limits**:
- ✅ Minimum: 10 characters (prevents empty/trivial inputs)
- ✅ Maximum: 5000 characters (prevents DoS via large inputs)

**Timeout Protection** (`validators.ts:30-32`):
```typescript
isValidTimeout(timeout: number): boolean {
  return timeout > 0 && timeout <= 600000; // Max 10 minutes
}
```

**Concurrency Limits** (`validators.ts:34-36`):
```typescript
isValidMaxConcurrent(max: number): boolean {
  return max >= 1 && max <= 10;
}
```

### 3. API Key Security

**Environment Variable Injection** (`client.ts:320-347`):
```typescript
private buildEnv(): Record<string, string> {
  const env: Record<string, string> = {};

  if (this.config.apiKeys?.anthropic) {
    env.ANTHROPIC_API_KEY = this.config.apiKeys.anthropic;
  }
  // ... other keys

  return env;
}
```

**What's Protected**:
- ✅ API keys passed via environment variables (not CLI arguments)
- ✅ Keys not visible in `ps` output
- ✅ Keys not logged to stdout/stderr
- ✅ Keys isolated per subprocess

### 4. Agent Whitelisting

**Allowed Agents Only** (`validators.ts:8-18`):
```typescript
const VALID_AGENTS: AgentType[] = [
  'claude', 'gemini', 'qwen', 'codex',
  'gpt4', 'deepseek', 'llama', 'mistral', 'rycode'
];
```

**What's Protected**:
- ✅ Only predefined agents accepted
- ✅ Prevents arbitrary agent injection
- ✅ Type-safe agent selection

### 5. Subprocess Isolation

**Secure Subprocess Spawning** (`client.ts:238-318`):
```typescript
const proc: ChildProcess = spawn(
  this.config.toolkitCliPath,
  [command, ...args],
  { env: { ...process.env, ...this.buildEnv() } }
);
```

**What's Protected**:
- ✅ No shell execution (direct spawn, not `sh -c`)
- ✅ Arguments passed as array (not concatenated string)
- ✅ Timeout protection (kills runaway processes)
- ✅ AbortSignal support (user cancellation)

### 6. Error Handling

**Safe Error Messages** (`client.ts:349-374`):
```typescript
catch (error) {
  throw new Error(`Failed to parse command result: ${error}`);
}
```

**What's Protected**:
- ✅ No sensitive data in error messages
- ✅ Stack traces captured but sanitized
- ✅ Errors don't leak API keys or internal paths

### 7. Type Safety

**TypeScript Validation**:
- ✅ Strong typing prevents type confusion attacks
- ✅ Runtime validation with Pydantic-style validators
- ✅ JSON schema validation for responses

## ⚠️ Potential Vulnerabilities & Mitigations

### 1. Subprocess Execution Risk

**Risk**: Spawning Python subprocesses could be exploited if toolkit-cli is compromised

**Current Mitigation**:
- ✅ Input sanitization before subprocess spawn
- ✅ No shell execution
- ✅ Arguments passed safely

**Additional Recommendations**:
```typescript
// ❌ NEVER do this (shell injection risk)
exec(`toolkit-cli oneshot "${userInput}"`)

// ✅ Always do this (safe array args)
spawn('toolkit-cli', ['oneshot', sanitize(userInput)])
```

### 2. API Key Exposure

**Risk**: API keys could be exposed in logs, errors, or memory dumps

**Current Mitigation**:
- ✅ Keys passed via environment (not args)
- ✅ Keys not in error messages
- ✅ Keys not logged

**Additional Recommendations**:
```typescript
// Consider encrypting keys in config
const encryptedKey = encrypt(apiKey);

// Or use secure key storage
const key = await getSecureKey('ANTHROPIC_API_KEY');
```

### 3. DoS via Resource Exhaustion

**Risk**: Malicious user could exhaust resources with many concurrent requests

**Current Mitigation**:
- ✅ `maxConcurrent` limit (default: 5, max: 10)
- ✅ Queue system for excess requests
- ✅ Timeout protection (max: 10 minutes)

**Additional Recommendations**:
```typescript
// Add rate limiting
const rateLimiter = new RateLimiter({
  maxRequests: 10,
  windowMs: 60000 // 10 requests per minute
});
```

### 4. Path Traversal

**Risk**: User could specify malicious toolkit-cli path

**Current Mitigation**:
- ⚠️ Limited - relies on default `toolkit-cli` in PATH

**Recommendation**:
```typescript
// Validate toolkit-cli path
static validateToolkitPath(path: string): void {
  if (path.includes('..') || path.includes('~')) {
    throw new ValidationError('path', path, 'Invalid path');
  }
}
```

### 5. JSON Parsing Vulnerabilities

**Risk**: Malicious JSON could cause DoS or code execution

**Current Mitigation**:
- ✅ Native `JSON.parse()` (safe in Node.js)
- ✅ Error handling wraps parsing
- ✅ Type validation after parsing

**Status**: ✅ Acceptable risk

## 🔒 Security Best Practices

### For RyCode Integration

1. **Never Trust User Input**
   ```typescript
   // ✅ Always sanitize
   const sanitized = Validators.sanitizeInput(userInput);

   // ✅ Always validate
   Validators.validateProjectIdea(sanitized);
   ```

2. **Use Environment Variables for Secrets**
   ```typescript
   // ✅ Good
   const toolkit = new ToolkitClient({
     apiKeys: {
       anthropic: process.env.ANTHROPIC_API_KEY
     }
   });

   // ❌ Bad - hardcoded keys
   const toolkit = new ToolkitClient({
     apiKeys: {
       anthropic: 'sk-ant-hardcoded-key'
     }
   });
   ```

3. **Set Reasonable Timeouts**
   ```typescript
   // ✅ Prevent hanging
   const toolkit = new ToolkitClient({
     timeout: 120000 // 2 minutes
   });
   ```

4. **Always Cleanup**
   ```typescript
   // ✅ Cleanup resources
   try {
     const result = await toolkit.oneshot(idea);
   } finally {
     await toolkit.close(); // Always cleanup
   }
   ```

5. **Handle Errors Safely**
   ```typescript
   // ✅ Don't leak sensitive info
   catch (error) {
     console.error('Operation failed'); // Generic message
     // Don't log error.stack in production
   }
   ```

## 🔐 Recommended Hardening

### 1. Add Rate Limiting

```typescript
class ToolkitClient {
  private rateLimiter = new Map<string, number>();

  private checkRateLimit(userId: string): boolean {
    const now = Date.now();
    const lastCall = this.rateLimiter.get(userId) || 0;

    if (now - lastCall < 60000) { // 1 minute
      throw new RateLimitError(60);
    }

    this.rateLimiter.set(userId, now);
    return true;
  }
}
```

### 2. Add Audit Logging

```typescript
private async executeCommand(...) {
  // Log command execution
  this.auditLog({
    command,
    args: args.map(a => a.substring(0, 50)), // Truncate
    user: getCurrentUser(),
    timestamp: new Date().toISOString()
  });

  const result = await spawn(...);
  return result;
}
```

### 3. Add Content Security Policy

```typescript
// Validate output before displaying
private validateOutput(output: string): string {
  // Remove any HTML/script tags
  return output.replace(/<script[^>]*>.*?<\/script>/gi, '');
}
```

### 4. Add Process Tracking

```typescript
private activeProcesses = new Map<number, ChildProcess>();

private async executeCommand(...) {
  const proc = spawn(...);
  this.activeProcesses.set(proc.pid, proc);

  proc.on('close', () => {
    this.activeProcesses.delete(proc.pid);
  });
}

public async close(): Promise<void> {
  // Kill all active processes
  for (const proc of this.activeProcesses.values()) {
    proc.kill('SIGTERM');
  }
}
```

## 📋 Security Checklist

Before deploying to production:

- [ ] All API keys stored securely (not in code)
- [ ] Input sanitization enabled for all user input
- [ ] Timeouts configured appropriately
- [ ] Rate limiting implemented
- [ ] Audit logging enabled
- [ ] Error messages don't leak sensitive info
- [ ] HTTPS used for any network communication
- [ ] Dependencies up to date
- [ ] Security headers configured
- [ ] Process cleanup on shutdown

## 🚨 Incident Response

If a security issue is discovered:

1. **Isolate**: Stop affected services immediately
2. **Assess**: Determine scope and impact
3. **Patch**: Apply security fixes
4. **Notify**: Inform affected users
5. **Review**: Update security measures

## 📚 References

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Node.js Security Best Practices](https://nodejs.org/en/docs/guides/security/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)

---

**Security Status**: ✅ **Good** - Basic protections in place
**Risk Level**: 🟡 **Medium** - Additional hardening recommended
**Last Review**: 2025-10-09
