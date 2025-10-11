# Security Features Migration Guide

This guide helps you adopt the new security features introduced in OpenCode 0.14.x.

## Overview of Changes

### New Features

1. **API Key Authentication** - Secure your server endpoints
2. **Rate Limiting** - Prevent abuse and resource exhaustion
3. **Path Validation** - Protect against directory traversal and sensitive file access

### Breaking Changes

None! All security features are **opt-in** and backward compatible.

## Migration Steps

### Step 1: Update OpenCode

```bash
npm update -g opencode-ai@latest
# or
brew upgrade opencode
```

### Step 2: Review Your Configuration

Check your current `opencode.json` or `opencode.jsonc`:

```bash
cat opencode.json
```

### Step 3: Enable Security Features (Optional)

#### Option A: Development/Local Use (No Changes Required)

If you're using OpenCode locally, you don't need to make any changes. Security features are disabled by default for backward compatibility.

#### Option B: Production/Exposed Server (Recommended)

Add security configuration to your `opencode.json`:

```json
{
  "server": {
    "require_auth": true,
    "api_keys": [
      "GENERATE_A_SECURE_KEY_HERE"
    ],
    "rate_limit": {
      "enabled": true,
      "limit": 100,
      "window_ms": 60000
    }
  }
}
```

**Generate a secure API key:**
```bash
# Using OpenSSL
openssl rand -hex 32

# Using Node.js
node -e "console.log(require('crypto').randomBytes(32).toString('hex'))"

# Using Python
python3 -c "import secrets; print(secrets.token_hex(32))"
```

### Step 4: Update Clients

If you enabled authentication, update your API clients:

**HTTP Requests:**
```typescript
fetch('http://localhost:3000/session', {
  headers: {
    'X-OpenCode-API-Key': 'your-api-key-here'
  }
})
```

**Go SDK:**
```go
import "github.com/sst/opencode-sdk-go"

client := opencode.NewClient(
    option.WithAPIKey("your-api-key-here"),
)
```

**JavaScript SDK:**
```typescript
import { OpenCode } from '@opencode-ai/sdk'

const client = new OpenCode({
  apiKey: 'your-api-key-here'
})
```

### Step 5: Test Your Setup

1. **Without API Key** (should fail if auth enabled):
```bash
curl http://localhost:3000/session
# Expected: 401 Unauthorized
```

2. **With API Key** (should succeed):
```bash
curl -H "X-OpenCode-API-Key: your-api-key-here" http://localhost:3000/session
# Expected: 200 OK
```

3. **Rate Limit Test**:
```bash
# Make multiple rapid requests
for i in {1..110}; do
  curl -H "X-OpenCode-API-Key: your-api-key-here" http://localhost:3000/session
done
# Expected: Last 10 requests return 429 (if limit is 100)
```

4. **Path Validation Test**:
```bash
curl -H "X-OpenCode-API-Key: your-api-key-here" \
  "http://localhost:3000/file/content?path=.env"
# Expected: 403 Forbidden
```

## Configuration Examples

### Basic Security (Recommended for most users)

```json
{
  "server": {
    "require_auth": true,
    "api_keys": ["your-secure-key-here"]
  }
}
```

### High Security (For public-facing servers)

```json
{
  "server": {
    "require_auth": true,
    "api_keys": [
      "primary-key-here",
      "secondary-key-for-rotation"
    ],
    "rate_limit": {
      "enabled": true,
      "limit": 50,
      "window_ms": 60000
    }
  }
}
```

### Development (Default)

```json
{
  "server": {
    "require_auth": false,
    "rate_limit": {
      "enabled": false
    }
  }
}
```

## Environment-Specific Configuration

### Using Environment Variables

Instead of hardcoding API keys, use environment variables:

```json
{
  "server": {
    "require_auth": true
  }
}
```

Then set keys via environment:
```bash
export OPENCODE_API_KEYS="key1,key2,key3"
```

*Note: Environment variable support coming in future release*

### Multi-Environment Setup

**Development (`opencode.dev.json`):**
```json
{
  "server": {
    "require_auth": false
  }
}
```

**Production (`opencode.prod.json`):**
```json
{
  "server": {
    "require_auth": true,
    "api_keys": ["${PROD_API_KEY}"]
  }
}
```

Load the appropriate config:
```bash
# Development
opencode --config opencode.dev.json

# Production
export OPENCODE_CONFIG=opencode.prod.json
opencode
```

## Troubleshooting

### "Missing API key" Error

**Problem**: Getting 401 errors even though you set an API key.

**Solutions**:
1. Check the header name is exactly `X-OpenCode-API-Key`
2. Verify the API key in your config matches what you're sending
3. Ensure `opencode.json` is in the current directory or parent directories

### Rate Limit Constantly Triggered

**Problem**: Hitting rate limits too frequently.

**Solutions**:
1. Increase the limit:
   ```json
   {
     "server": {
       "rate_limit": {
         "limit": 200
       }
     }
   }
   ```
2. Implement exponential backoff in your client
3. Cache responses when possible

### Can't Access Legitimate Files

**Problem**: Getting 403 errors for files that should be accessible.

**Solutions**:
1. Check if the file matches a sensitive pattern (`.env`, `credentials`, etc.)
2. Verify the file is within your project directory
3. Use relative paths from your project root
4. If you need to access a file with a sensitive name pattern, rename it

### Authentication Bypassed on Localhost

**Problem**: Auth works on production but not locally.

**Explanation**: This is by design! Localhost requests bypass authentication in development mode.

**To test auth locally**:
```typescript
AuthMiddleware.middleware(c, next, {
  bypassLocalhost: false
})
```

## Best Practices

### 1. Key Rotation

Rotate API keys regularly:

```bash
# Generate new key
NEW_KEY=$(openssl rand -hex 32)

# Add to config (keep old key temporarily)
{
  "server": {
    "api_keys": [
      "new-key-here",
      "old-key-here"
    ]
  }
}

# Update clients to use new key
# Remove old key after all clients updated
```

### 2. Monitoring

Monitor rate limit headers to optimize limits:

```bash
curl -I http://localhost:3000/session \
  -H "X-OpenCode-API-Key: your-key" \
  | grep X-RateLimit
```

### 3. Security Headers

Use a reverse proxy for additional security:

**nginx example:**
```nginx
server {
    listen 443 ssl;
    server_name opencode.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## Getting Help

If you encounter issues:

1. Check the [Security Documentation](./SECURITY.md)
2. Review the [API Documentation](https://opencode.ai/docs)
3. Ask in [Discord](https://opencode.ai/discord)
4. Open an [issue](https://github.com/sst/opencode/issues) with:
   - Your OpenCode version
   - Configuration (redact API keys!)
   - Error messages
   - Steps to reproduce

## Rollback Instructions

If you need to rollback security features:

```json
{
  "server": {
    "require_auth": false,
    "rate_limit": {
      "enabled": false
    }
  }
}
```

Or remove the `server` section entirely to use defaults.

## Next Steps

- [ ] Review [SECURITY.md](./SECURITY.md) for detailed security information
- [ ] Set up monitoring for security events
- [ ] Configure your reverse proxy/load balancer
- [ ] Test your security setup thoroughly
- [ ] Document your security configuration for your team

## Changelog

- **0.14.0**: Initial release of security features
  - API key authentication
  - Rate limiting
  - Path validation
