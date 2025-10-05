# Security Policy

## Reporting Security Issues

OpenCode takes security seriously. If you discover a security vulnerability, please report it to:

**Email**: support@sst.dev
**Subject**: [SECURITY] OpenCode Vulnerability Report

Please include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

We aim to respond to security reports within 48 hours and will provide regular updates on the resolution progress.

## Security Features

### Authentication

OpenCode supports API key authentication for server endpoints to prevent unauthorized access.

#### Enabling Authentication

Add to your `opencode.json` or `opencode.jsonc`:

```json
{
  "server": {
    "require_auth": true,
    "api_keys": [
      "your-secret-api-key-here"
    ]
  }
}
```

#### Using API Keys

**Via Header** (recommended):
```bash
curl -H "X-OpenCode-API-Key: your-secret-api-key-here" http://localhost:3000/session
```

**Via Query Parameter**:
```bash
curl http://localhost:3000/session?api_key=your-secret-api-key-here
```

#### Localhost Bypass

By default, authentication is bypassed for localhost connections in development mode. This can be configured via the middleware options.

### Rate Limiting

OpenCode includes built-in rate limiting to prevent abuse and resource exhaustion.

#### Configuration

```json
{
  "server": {
    "rate_limit": {
      "enabled": true,
      "limit": 100,
      "window_ms": 60000
    }
  }
}
```

**Options:**
- `enabled` (boolean, default: `true`): Enable/disable rate limiting
- `limit` (number, default: `100`): Maximum requests per window
- `window_ms` (number, default: `60000`): Time window in milliseconds (1 minute)

#### Rate Limit Headers

Responses include rate limit information:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 2025-10-04T22:30:00.000Z
```

When rate limited, you'll receive a `429 Too Many Requests` response with a `Retry-After` header.

### Path Validation

OpenCode validates all file paths to prevent directory traversal attacks and access to sensitive files.

#### Protected Files

The following file patterns are automatically blocked:

**Credentials and Secrets:**
- `.env`, `.env.*`
- `*.pem`, `*.key`, `*.p12`, `*.pfx`
- `*credentials*`, `*secret*`, `*password*`
- SSH keys: `id_rsa`, `id_dsa`, `id_ed25519`

**System Files (Unix/Linux):**
- `/etc/passwd`, `/etc/shadow`, `/etc/hosts`
- `/etc/ssh/*`, `/root/*`

**System Files (macOS):**
- `/System/*`, `/Library/Keychains/*`

**Cloud Provider Credentials:**
- `.aws/credentials`, `.azure/credentials`, `.gcp/credentials`
- `.ssh/*`, `.git-credentials`, `.netrc`
- `kubeconfig`, `.kube/config`

**Database Files:**
- `*.sqlite`, `*.db`

#### Path Traversal Prevention

All file operations are restricted to:
1. The current project directory (`Instance.directory`)
2. The git worktree root (`Instance.worktree`)

Attempts to access files outside these boundaries will be rejected with a `403 Forbidden` error.

#### Error Responses

**Path Traversal Attempt:**
```json
{
  "name": "PathTraversalError",
  "data": {
    "requestedPath": "../../etc/passwd",
    "message": "Path '../../etc/passwd' is outside allowed directories"
  }
}
```

**Sensitive File Access:**
```json
{
  "name": "SensitiveFileError",
  "data": {
    "requestedPath": ".env",
    "message": "Access to sensitive file '.env' is not allowed"
  }
}
```

## Best Practices

### 1. API Key Management

- **Generate Strong Keys**: Use cryptographically secure random strings
  ```bash
  openssl rand -hex 32
  ```
- **Rotate Regularly**: Change API keys periodically
- **Store Securely**: Use environment variables or secret management systems
- **Never Commit**: Add `opencode.json` to `.gitignore` if it contains keys

### 2. Network Security

- **Use HTTPS**: Always use HTTPS in production
- **Firewall**: Restrict server access to trusted networks
- **Reverse Proxy**: Use nginx/Caddy for SSL termination and additional security

### 3. Deployment Security

- **Disable Localhost Bypass**: Set `bypassLocalhost: false` in production
- **Enable Auth**: Always enable authentication for exposed servers
- **Monitor Logs**: Review logs for suspicious activity
- **Update Regularly**: Keep OpenCode updated for security patches

### 4. File System Security

- **Principle of Least Privilege**: Run OpenCode with minimal required permissions
- **Isolated Directories**: Use dedicated project directories
- **Read-Only Where Possible**: Mount sensitive directories as read-only

## Security Checklist for Production

- [ ] Authentication enabled (`require_auth: true`)
- [ ] Strong API keys configured
- [ ] Rate limiting enabled
- [ ] Running behind HTTPS
- [ ] Localhost bypass disabled for public deployments
- [ ] Logs monitored for security events
- [ ] Regular security updates applied
- [ ] Sensitive files excluded from project directories
- [ ] Firewall rules configured
- [ ] Regular security audits performed

## Vulnerability Disclosure Timeline

1. **Day 0**: Security vulnerability reported
2. **Day 2**: Initial response and acknowledgment
3. **Day 7**: Investigation and fix development
4. **Day 14**: Patch released (if possible)
5. **Day 30**: Public disclosure (after users have time to update)

## Security Updates

Security patches are released as soon as possible after discovery. Critical vulnerabilities may receive out-of-band releases.

Subscribe to security announcements:
- GitHub Security Advisories: https://github.com/sst/opencode/security/advisories
- Discord: https://opencode.ai/discord

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.14.x  | :white_check_mark: |
| 0.13.x  | :white_check_mark: |
| < 0.13  | :x:                |

## Attribution

We appreciate the security research community's efforts. Security researchers who responsibly disclose vulnerabilities will be credited in our security advisories (unless they prefer to remain anonymous).

## Contact

For non-security issues, please use:
- GitHub Issues: https://github.com/sst/opencode/issues
- Discord: https://opencode.ai/discord

For security-related inquiries: support@sst.dev
