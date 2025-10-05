# Plugin Signatures

> **Cryptographic signature verification for plugin authenticity and integrity**

Plugin signatures provide cryptographic proof that a plugin was created by a trusted developer and hasn't been tampered with. This complements hash-based verification with identity verification.

---

## 🎯 Overview

RyCode supports two signature methods:

1. **Crypto Signatures** (Built-in) - RSA signatures using Node.js crypto
2. **GPG Signatures** (Optional) - Traditional GPG/PGP signatures

Both methods provide:
- ✅ **Authenticity** - Proves who created the plugin
- ✅ **Integrity** - Detects any tampering
- ✅ **Trust Chain** - Verifiable chain of trust
- ✅ **Non-repudiation** - Signer cannot deny signing

---

## 🚀 Quick Start

### For Plugin Users

**Verify a signed plugin:**

```bash
# Download plugin and signature
curl -O https://example.com/my-plugin.js
curl -O https://example.com/my-plugin.sig.json

# Download public key
curl -O https://example.com/signing-key-public.pem

# Verify signature
rycode plugin:verify-sig my-plugin.js my-plugin.sig.json --public-key signing-key-public.pem
```

### For Plugin Developers

**Sign your plugin:**

```bash
# 1. Generate signing keys (one-time)
rycode plugin:keygen --output ./keys --name my-plugin

# 2. Sign your plugin
rycode plugin:sign dist/my-plugin.js \
  --key ./keys/my-plugin-private.pem \
  --output dist/my-plugin.sig.json

# 3. Distribute plugin, signature, and public key
```

---

## 📖 Detailed Guide

### Crypto Signatures (Recommended)

Crypto signatures use RSA keys and are built into RyCode. No external dependencies required.

#### Generate Key Pair

```bash
rycode plugin:keygen --output ./keys --name my-company
```

**Output:**
```
✓ RSA Key Pair Generated

Key ID:      A1B2C3D4E5F6G7H8
Private Key: ./keys/my-company-private.pem
Public Key:  ./keys/my-company-public.pem

⚠ Keep your private key secure!
```

**Important:**
- Add `*-private.pem` to `.gitignore`
- Store private keys in secure vault (e.g., 1Password, AWS Secrets Manager)
- Never commit private keys to version control

#### Sign a Plugin

```bash
rycode plugin:sign dist/my-plugin.js \
  --key ./keys/my-company-private.pem \
  --algorithm RSA-SHA256 \
  --output dist/my-plugin.sig.json
```

**Signature file format:**

```json
{
  "algorithm": "RSA-SHA256",
  "signature": "Base64EncodedSignature==",
  "keyId": "A1B2C3D4E5F6G7H8",
  "timestamp": 1696521600000,
  "publicKey": "-----BEGIN PUBLIC KEY-----\n..."
}
```

#### Verify Signature

```bash
rycode plugin:verify-sig dist/my-plugin.js dist/my-plugin.sig.json \
  --public-key ./keys/my-company-public.pem
```

**Output on success:**
```
Signature Verification

File:      dist/my-plugin.js
Algorithm: RSA-SHA256
Key ID:    A1B2C3D4E5F6G7H8
Signed:    2025-10-05T10:00:00.000Z

✓ Signature is VALID
  Plugin has been signed by the claimed signer.
```

**Output on failure:**
```
✗ Signature is INVALID
  Signature verification failed

⚠ DO NOT use this plugin.
```

---

### GPG Signatures (Advanced)

For teams already using GPG, RyCode supports GPG signatures.

#### Prerequisites

Install GPG:
```bash
# macOS
brew install gnupg

# Ubuntu/Debian
sudo apt-get install gnupg

# Verify installation
gpg --version
```

#### Generate GPG Key

```bash
gpg --full-generate-key

# Choose:
# - RSA and RSA (default)
# - 2048 bits minimum (4096 recommended)
# - Valid for 1-2 years
# - Your name and email
```

#### Export Public Key

```bash
# List your keys
gpg --list-keys

# Export public key
gpg --armor --export your@email.com > signing-key-public.asc
```

#### Sign with GPG

Currently, use GPG directly (CLI integration coming soon):

```bash
# Create detached signature
gpg --armor --detach-sign dist/my-plugin.js

# Creates dist/my-plugin.js.asc
```

#### Verify GPG Signature

```bash
# Import signer's public key
gpg --import signing-key-public.asc

# Verify signature
gpg --verify dist/my-plugin.js.asc dist/my-plugin.js
```

---

## 🔐 Trust Levels

RyCode supports three trust levels for signers:

### Full Trust

```typescript
{
  name: "RyCode Official",
  keyId: "RYCODE-2025",
  publicKey: "...",
  trustLevel: "full"  // Fully trusted
}
```

**Used for:**
- Official RyCode plugins
- Your organization's internal plugins
- Well-known trusted developers

### Marginal Trust

```typescript
{
  name: "Community Developer",
  keyId: "DEV-123",
  publicKey: "...",
  trustLevel: "marginal"  // Partially trusted
}
```

**Used for:**
- Community plugins under review
- New contributors
- Transitional trust period

### Never Trust

```typescript
{
  name: "Untrusted Source",
  keyId: "UNKNOWN",
  publicKey: "...",
  trustLevel: "never"  // Explicitly distrusted
}
```

**Used for:**
- Revoked keys
- Compromised signers
- Blacklisted sources

---

## 📋 Signature Policy

Configure signature requirements in `.rycode.json`:

```jsonc
{
  "plugin_security": {
    "signature_policy": {
      // Require signatures for all plugins
      "requireSignatures": false,

      // List of trusted signers
      "trustedSigners": [
        {
          "name": "My Company",
          "keyId": "A1B2C3D4E5F6G7H8",
          "publicKey": "-----BEGIN PUBLIC KEY-----\n...",
          "organization": "ACME Corp",
          "email": "security@acme.com",
          "trustLevel": "full"
        }
      ],

      // Allow self-signed in development
      "allowSelfSigned": true,

      // Signature expiration (days, 0 = never)
      "signatureExpiration": 365
    }
  }
}
```

---

## 🔄 Workflow Examples

### Workflow 1: Sign Official Plugin

```bash
# 1. Build plugin
bun build src/index.ts --outfile dist/plugin.js

# 2. Generate hash
HASH=$(rycode plugin:hash dist/plugin.js --json | jq -r '.hash')

# 3. Sign plugin
rycode plugin:sign dist/plugin.js \
  --key ~/.rycode/official-signing-key.pem \
  --output dist/plugin.sig.json

# 4. Add to registry with signature
rycode plugin:registry:add dist/plugin.js my-plugin 1.0.0 \
  --description "Official plugin" \
  --verified-by official

# 5. Publish
npm publish
```

### Workflow 2: Verify Before Installation

```bash
# 1. Install plugin
npm install opencode-my-plugin

# 2. Locate plugin file
PLUGIN_PATH=$(npm root)/opencode-my-plugin/dist/plugin.js

# 3. Download signature and public key from registry/npm
curl -O https://registry.rycode.ai/opencode-my-plugin/1.0.0/signature.json
curl -O https://registry.rycode.ai/keys/opencode-official.pem

# 4. Verify signature
rycode plugin:verify-sig $PLUGIN_PATH signature.json \
  --public-key opencode-official.pem

# 5. Verify hash against registry
rycode plugin:verify $PLUGIN_PATH --hash <hash-from-registry>

# 6. If both pass, plugin is safe to use
```

### Workflow 3: Corporate Plugin Signing

```bash
# Setup (once)
# 1. Generate company signing key
rycode plugin:keygen --output /secure/vault/ --name acme-plugins

# 2. Store private key in secrets manager
aws secretsmanager create-secret \
  --name acme-plugin-signing-key \
  --secret-string file:///secure/vault/acme-plugins-private.pem

# 3. Distribute public key to team
cp /secure/vault/acme-plugins-public.pem /shared/keys/

# CI/CD Pipeline
# 1. Retrieve signing key from secrets manager
aws secretsmanager get-secret-value \
  --secret-id acme-plugin-signing-key \
  --query SecretString \
  --output text > /tmp/signing-key.pem

# 2. Build and sign
bun run build
rycode plugin:sign dist/plugin.js \
  --key /tmp/signing-key.pem \
  --output dist/plugin.sig.json

# 3. Upload to artifact storage
aws s3 cp dist/plugin.js s3://plugins/my-plugin/1.0.0/
aws s3 cp dist/plugin.sig.json s3://plugins/my-plugin/1.0.0/
aws s3 cp /shared/keys/acme-plugins-public.pem s3://plugins/keys/

# 4. Clean up
rm /tmp/signing-key.pem
```

---

## 🛡️ Security Best Practices

### For Signers

1. **Protect Private Keys**
   - Use hardware security modules (HSM) for production
   - Encrypt key files at rest
   - Use strong passphrases
   - Rotate keys annually

2. **Key Storage**
   ```bash
   # Encrypt private key
   gpg --symmetric --cipher-algo AES256 signing-key-private.pem

   # Store encrypted version only
   rm signing-key-private.pem
   git add signing-key-private.pem.gpg
   ```

3. **Signature Distribution**
   - Include signature files in npm packages
   - Upload to registry alongside plugins
   - Provide public key download link
   - Document signature verification steps

4. **Key Rotation**
   ```bash
   # Generate new key
   rycode plugin:keygen --name my-plugin-2026

   # Sign with both keys (transition period)
   rycode plugin:sign plugin.js --key old-key.pem --output plugin.sig.old.json
   rycode plugin:sign plugin.js --key new-key.pem --output plugin.sig.new.json

   # After transition, revoke old key
   ```

### For Verifiers

1. **Always Verify**
   ```bash
   # Never skip verification
   rycode plugin:verify-sig plugin.js signature.json --public-key public-key.pem
   ```

2. **Trust Verification**
   - Only trust known public keys
   - Verify public key fingerprint out-of-band
   - Use official key distribution channels

3. **Check Expiration**
   ```bash
   # Signatures older than 1 year should be re-verified
   jq '.timestamp' signature.json
   ```

4. **Defense in Depth**
   ```bash
   # Combine hash + signature verification
   rycode plugin:verify plugin.js --hash <hash>
   rycode plugin:verify-sig plugin.js signature.json --public-key key.pem
   ```

---

## 🔧 Programmatic API

### Sign a Plugin

```typescript
import { PluginSignature } from "./src/plugin/signature"

const keyPair = PluginSignature.generateKeyPair()

const signature = await PluginSignature.signWithCrypto(
  "/path/to/plugin.js",
  keyPair.privateKey,
  "RSA-SHA256",
  keyPair.publicKey
)

console.log(signature)
// {
//   algorithm: "RSA-SHA256",
//   signature: "base64...",
//   keyId: "A1B2C3D4...",
//   timestamp: 1696521600000,
//   publicKey: "-----BEGIN PUBLIC KEY-----\n..."
// }
```

### Verify a Signature

```typescript
const result = await PluginSignature.verifyCryptoSignature(
  "/path/to/plugin.js",
  signature,
  publicKey
)

if (result.valid) {
  console.log("✓ Signature verified")
} else {
  console.error("✗ Invalid signature:", result.error)
}
```

### Check Trust

```typescript
const trustedSigners: PluginSignature.TrustedSigner[] = [
  {
    name: "RyCode Official",
    keyId: "RYCODE-2025",
    publicKey: "...",
    trustLevel: "full",
  }
]

const isTrusted = PluginSignature.isTrustedSigner(
  signature.keyId,
  trustedSigners
)

const signer = PluginSignature.getSigner(
  signature.keyId,
  trustedSigners
)
```

---

## 🔍 Troubleshooting

### Signature Verification Failed

**Problem:** Signature doesn't match

**Solutions:**
1. File was modified after signing - re-download
2. Wrong public key - verify key fingerprint
3. Corrupted signature file - re-download signature

### GPG Not Available

**Problem:** `GPG is not available`

**Solution:**
```bash
# Install GPG
brew install gnupg  # macOS
sudo apt install gnupg  # Linux

# Verify
gpg --version
```

### Key Not Found

**Problem:** `Signer X is not in trusted signers list`

**Solution:**
```json
// Add to .rycode.json
{
  "plugin_security": {
    "signature_policy": {
      "trustedSigners": [
        {
          "name": "Signer Name",
          "keyId": "X",
          "publicKey": "...",
          "trustLevel": "full"
        }
      ]
    }
  }
}
```

### Expired Signature

**Problem:** `Signature expired (age: 400 days)`

**Solution:**
- Contact plugin author for re-signed version
- Or adjust expiration policy (not recommended)

---

## 📊 Comparison: Crypto vs GPG

| Feature | Crypto (RSA) | GPG/PGP |
|---------|-------------|---------|
| **Setup** | Simple (built-in) | Requires GPG install |
| **Performance** | Fast | Moderate |
| **Key Management** | Manual | GPG keyring |
| **Web of Trust** | ❌ No | ✅ Yes |
| **Industry Standard** | ✅ Yes (RSA) | ✅ Yes (PGP) |
| **RyCode Integration** | ✅ Full | ⚠️ Partial |
| **CI/CD Friendly** | ✅ Yes | ⚠️ Needs setup |

**Recommendation:**
- **Use Crypto** for most cases (simpler, no dependencies)
- **Use GPG** if already using GPG for code signing

---

## 🔗 Related Documentation

- [Plugin Security Guide](./PLUGIN_SECURITY.md) - Overall security system
- [Plugin Registry](./PLUGIN_REGISTRY.md) - Hash-based verification
- [Security Policy](./SECURITY.md) - Security procedures and policies

---

## 📚 Further Reading

### Standards & Specifications

- [RFC 8017 - PKCS #1: RSA Cryptography](https://tools.ietf.org/html/rfc8017)
- [RFC 4880 - OpenPGP Message Format](https://tools.ietf.org/html/rfc4880)
- [NIST FIPS 186-4 - Digital Signature Standard](https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.186-4.pdf)

### Best Practices

- [OWASP Code Signing Best Practices](https://owasp.org/www-community/controls/Code_Signing)
- [Node.js Crypto Module](https://nodejs.org/api/crypto.html)
- [GPG Handbook](https://www.gnupg.org/gph/en/manual.html)

---

**Last Updated:** October 5, 2025

For questions about plugin signatures, contact [security@rycode.ai](mailto:security@rycode.ai)
