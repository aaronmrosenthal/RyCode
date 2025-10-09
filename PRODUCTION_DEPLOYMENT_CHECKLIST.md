# Production Deployment Checklist

Complete this checklist before deploying RyCode to production with security features enabled.

## Pre-Deployment

### 1. Security Configuration

- [ ] **Generate encryption key**
  ```bash
  bun run packages/rycode/src/index.ts generate-key
  ```

- [ ] **Set encryption key in environment**
  ```bash
  export RYCODE_ENCRYPTION_KEY="<generated-key>"
  ```

- [ ] **Store key securely**
  - [ ] Add to secrets manager (AWS Secrets Manager, HashiCorp Vault, etc.)
  - [ ] Never commit to git
  - [ ] Document key rotation procedure
  - [ ] Set up key backup strategy

- [ ] **Configure HTTPS**
  - [ ] SSL/TLS certificates installed
  - [ ] Redirect HTTP → HTTPS
  - [ ] HSTS headers will auto-enable on HTTPS

### 2. File Permissions

- [ ] **Set restrictive permissions on sensitive files**
  ```bash
  chmod 600 ~/.local/share/rycode/auth.json
  chmod 600 ~/.local/share/rycode/*.json
  ```

- [ ] **Verify ownership**
  ```bash
  chown $USER:$USER ~/.local/share/rycode/*
  ```

### 3. Data Migration

- [ ] **Backup existing data**
  ```bash
  cp ~/.local/share/rycode/auth.json ~/.local/share/rycode/auth.json.backup
  ```

- [ ] **Migrate to encrypted storage**
  ```typescript
  import { Auth } from "./auth"
  const count = await Auth.migrateToEncrypted()
  console.log(`Migrated ${count} credentials`)
  ```

- [ ] **Verify migration**
  - [ ] Test authentication with encrypted credentials
  - [ ] Confirm no plaintext warnings in logs

### 4. Security Headers

- [ ] **HSTS enabled** (automatic on HTTPS)
  ```bash
  curl -I https://your-domain.com | grep Strict-Transport-Security
  # Expected: Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
  ```

- [ ] **CSP configured**
  ```bash
  curl -I https://your-domain.com | grep Content-Security-Policy
  ```

- [ ] **Other security headers present**
  - [ ] X-Content-Type-Options: nosniff
  - [ ] X-Frame-Options: DENY
  - [ ] X-XSS-Protection: 1; mode=block

### 5. Rate Limiting

- [ ] **Configure rate limits** in `.rycode/config.json`
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

- [ ] **Test rate limiting**
  - [ ] Verify 429 responses when limit exceeded
  - [ ] Confirm legitimate traffic not blocked

### 6. Plugin Security

- [ ] **Review plugin permissions**
  ```json
  {
    "plugins": {
      "*": {
        "sandbox": {
          "timeout": 30000,
          "maxMemoryMB": 256,
          "allowNetwork": false,
          "strictMode": true
        }
      }
    }
  }
  ```

- [ ] **Verify plugin signatures** (if using signed plugins)

### 7. Logging & Monitoring

- [ ] **Configure log levels**
  - [ ] Production: `info` or `warn`
  - [ ] Staging: `debug`

- [ ] **Set up monitoring**
  - [ ] Integrity check failures
  - [ ] Authentication failures
  - [ ] Rate limit violations
  - [ ] Unusual activity patterns

- [ ] **Log aggregation** (optional but recommended)
  - [ ] Centralized logging (Datadog, Splunk, ELK)
  - [ ] Log retention policy
  - [ ] Security event alerts

### 8. Testing

- [ ] **Run full test suite**
  ```bash
  bun test --timeout 60000
  # Expected: 302 tests passing
  ```

- [ ] **Integration tests**
  - [ ] Authentication flow
  - [ ] API endpoints
  - [ ] Plugin execution
  - [ ] File operations

- [ ] **Security tests**
  - [ ] Encryption/decryption
  - [ ] Integrity verification
  - [ ] Tamper detection
  - [ ] Key validation

## Deployment

### 9. Deploy Application

- [ ] **Build for production**
  ```bash
  bun run build
  ```

- [ ] **Deploy to environment**
  - [ ] Update environment variables
  - [ ] Deploy application code
  - [ ] Run database migrations (if applicable)

- [ ] **Verify deployment**
  - [ ] Health check endpoint responding
  - [ ] HTTPS working
  - [ ] Authentication working
  - [ ] No errors in logs

### 10. Post-Deployment Verification

- [ ] **Smoke tests**
  - [ ] User can authenticate
  - [ ] API requests work
  - [ ] Data encrypted at rest
  - [ ] Integrity checks passing

- [ ] **Security headers verification**
  ```bash
  curl -I https://your-domain.com
  ```

- [ ] **Log review**
  - [ ] No security warnings
  - [ ] No integrity failures
  - [ ] Authentication events logged

## Post-Deployment

### 11. Documentation

- [ ] **Update runbooks**
  - [ ] Key rotation procedure
  - [ ] Incident response plan
  - [ ] Recovery procedures

- [ ] **Document configuration**
  - [ ] Environment variables
  - [ ] Security settings
  - [ ] Plugin configurations

### 12. Backup & Recovery

- [ ] **Set up automated backups**
  - [ ] Auth data backup schedule
  - [ ] Encryption key backup (secure location)
  - [ ] Configuration backups

- [ ] **Test recovery procedure**
  - [ ] Restore from backup
  - [ ] Verify decryption with backup key
  - [ ] Document recovery time

### 13. Monitoring Setup

- [ ] **Configure alerts**
  - [ ] Integrity check failures
  - [ ] Authentication failures (threshold)
  - [ ] Rate limit violations
  - [ ] Server errors

- [ ] **Dashboard setup**
  - [ ] Security event metrics
  - [ ] Authentication success/failure rates
  - [ ] API usage patterns

### 14. Key Rotation Plan

- [ ] **Document rotation schedule** (recommended: quarterly)

- [ ] **Rotation procedure**
  1. Generate new key
  2. Decrypt all data with old key
  3. Re-encrypt with new key
  4. Update environment variable
  5. Securely delete old key
  6. Verify all services working

- [ ] **Set calendar reminder** for next rotation

### 15. Compliance & Audit

- [ ] **Document compliance**
  - [ ] OWASP compliance: 90%
  - [ ] Encryption algorithm: AES-256-GCM
  - [ ] Key derivation: PBKDF2 100K iterations
  - [ ] Integrity: SHA-256 checksums

- [ ] **Prepare for audit**
  - [ ] Security assessment document ready
  - [ ] Implementation summary available
  - [ ] Test results documented

## Emergency Procedures

### Integrity Check Failure

If integrity check fails in production:

1. **Immediate action**
   - Isolate affected system
   - Do not use potentially tampered data
   - Alert security team

2. **Investigation**
   - Review logs for suspicious activity
   - Check file modification times
   - Verify backup integrity

3. **Recovery**
   - Restore from verified backup
   - Force re-authentication
   - Update credentials if compromised

### Key Compromise

If encryption key is compromised:

1. **Immediate action**
   - Rotate key immediately
   - Force all users to re-authenticate
   - Review access logs

2. **Investigation**
   - Determine scope of compromise
   - Identify affected data
   - Check for unauthorized access

3. **Recovery**
   - Generate new key
   - Re-encrypt all data
   - Update all environments
   - Document incident

## Production Configuration Reference

### Recommended Settings

```json
{
  "server": {
    "require_auth": true,
    "rate_limit": {
      "enabled": true,
      "limit": 100,
      "window_ms": 60000
    }
  },
  "security": {
    "hsts": {
      "enabled": true,
      "maxAge": 31536000,
      "includeSubDomains": true,
      "preload": true
    },
    "csp": {
      "enabled": true,
      "useNonces": true
    }
  },
  "plugins": {
    "*": {
      "sandbox": {
        "timeout": 30000,
        "maxMemoryMB": 256,
        "allowNetwork": false,
        "strictMode": true
      }
    }
  },
  "logging": {
    "level": "info",
    "security_events": true
  }
}
```

### Environment Variables

```bash
# Required
RYCODE_ENCRYPTION_KEY="<base64-encoded-key>"

# Optional
NODE_ENV="production"
LOG_LEVEL="info"
```

## Support

- **Security Documentation:** [SECURITY_ASSESSMENT.md](./SECURITY_ASSESSMENT.md)
- **Migration Guide:** [SECURITY_MIGRATION_GUIDE.md](./SECURITY_MIGRATION_GUIDE.md)
- **API Documentation:** [packages/rycode/src/storage/README.md](./packages/rycode/src/storage/README.md)

## Security Contacts

- **Security Issues:** Create private security advisory (GitHub)
- **Emergency:** Document emergency contact procedure

---

**Version:** 1.0.0
**Last Updated:** 2025-10-08
**Next Review:** Quarterly (2026-01-08)

## Sign-Off

- [ ] **Technical Lead:** Security implementation reviewed and approved
- [ ] **Security Team:** Configuration audited and approved
- [ ] **DevOps:** Deployment procedure validated
- [ ] **Compliance:** Standards compliance verified

---

✅ **All items checked** → Ready for production deployment
