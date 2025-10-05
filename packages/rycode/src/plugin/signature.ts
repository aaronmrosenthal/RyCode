/**
 * Plugin Signature Verification
 *
 * Provides GPG signature verification for plugins to ensure authenticity
 * and integrity. Works in conjunction with the registry system.
 */

import z from "zod/v4"
import { NamedError } from "../util/error"
import { Log } from "../util/log"
import crypto from "crypto"
import { $ } from "bun"
import { existsSync } from "fs"

export namespace PluginSignature {
  const log = Log.create({ service: "plugin.signature" })

  /**
   * Signature verification error
   */
  export const SignatureVerificationError = NamedError.create(
    "SignatureVerificationError",
    z.object({
      plugin: z.string(),
      reason: z.string(),
    })
  )
  export type SignatureVerificationError = InstanceType<typeof SignatureVerificationError>

  /**
   * Signature data structure
   */
  export const Signature = z.object({
    /** Algorithm used (e.g., "RSA-SHA256", "Ed25519") */
    algorithm: z.string(),
    /** Base64-encoded signature */
    signature: z.string(),
    /** Signer's key ID or fingerprint */
    keyId: z.string(),
    /** Timestamp when signed */
    timestamp: z.number(),
    /** Optional: Public key (PEM format) */
    publicKey: z.string().optional(),
  })
  export type Signature = z.infer<typeof Signature>

  /**
   * Trusted signer configuration
   */
  export const TrustedSigner = z.object({
    /** Signer name */
    name: z.string(),
    /** Key ID or fingerprint */
    keyId: z.string(),
    /** Public key (PEM format) */
    publicKey: z.string(),
    /** Organization */
    organization: z.string().optional(),
    /** Email */
    email: z.string().email().optional(),
    /** Trust level */
    trustLevel: z.enum(["full", "marginal", "never"]).default("full"),
  })
  export type TrustedSigner = z.infer<typeof TrustedSigner>

  /**
   * Signature policy
   */
  export const Policy = z.object({
    /** Require signatures for all plugins */
    requireSignatures: z.boolean().default(false),
    /** List of trusted signers */
    trustedSigners: z.array(TrustedSigner).default([]),
    /** Accept self-signed plugins in development */
    allowSelfSigned: z.boolean().default(true),
    /** Signature expiration in days (0 = no expiration) */
    signatureExpiration: z.number().default(365),
  })
  export type Policy = z.infer<typeof Policy>

  /**
   * Default official RyCode signers
   */
  export const OFFICIAL_SIGNERS: TrustedSigner[] = [
    {
      name: "RyCode Release Team",
      keyId: "RYCODE-RELEASE-2025",
      publicKey: "", // Will be populated with actual public key
      organization: "RyCode",
      email: "security@rycode.ai",
      trustLevel: "full",
    },
  ]

  /**
   * Check if GPG is available
   */
  export async function isGPGAvailable(): Promise<boolean> {
    try {
      const result = await $`gpg --version`.quiet().nothrow()
      return result.exitCode === 0
    } catch {
      return false
    }
  }

  /**
   * Generate signature for a file using GPG
   */
  export async function signFile(
    filePath: string,
    keyId: string,
    passphrase?: string
  ): Promise<Signature> {
    if (!existsSync(filePath)) {
      throw new Error(`File not found: ${filePath}`)
    }

    const gpgAvailable = await isGPGAvailable()
    if (!gpgAvailable) {
      throw new Error("GPG is not available. Please install GPG to sign plugins.")
    }

    try {
      // Generate detached signature
      const sigPath = `${filePath}.sig`

      let result

      if (passphrase) {
        // Use passphrase for signing
        result = await $`gpg --batch --yes --pinentry-mode loopback --passphrase ${passphrase} --armor --detach-sign --local-user ${keyId} ${filePath}`.quiet().nothrow()
      } else {
        // Sign without passphrase
        result = await $`gpg --armor --detach-sign --local-user ${keyId} ${filePath}`.quiet().nothrow()
      }

      if (result.exitCode !== 0) {
        throw new Error(`GPG signing failed: ${result.stderr}`)
      }

      // Read the signature
      const sigFile = Bun.file(sigPath)
      const signatureContent = await sigFile.text()

      // Clean up signature file
      await Bun.write(sigPath, "") // Delete sig file

      const signature: Signature = {
        algorithm: "GPG",
        signature: Buffer.from(signatureContent).toString("base64"),
        keyId,
        timestamp: Date.now(),
      }

      log.info("file signed", { path: filePath, keyId })

      return signature
    } catch (error) {
      log.error("signing failed", {
        path: filePath,
        error: error instanceof Error ? error.message : String(error),
      })
      throw error
    }
  }

  /**
   * Verify signature using GPG
   */
  export async function verifySignature(
    filePath: string,
    signature: Signature,
    trustedSigners: TrustedSigner[] = OFFICIAL_SIGNERS
  ): Promise<{
    valid: boolean
    signer?: TrustedSigner
    error?: string
  }> {
    const gpgAvailable = await isGPGAvailable()
    if (!gpgAvailable) {
      return {
        valid: false,
        error: "GPG is not available",
      }
    }

    try {
      // Write signature to temp file
      const sigPath = `${filePath}.sig`
      const sigContent = Buffer.from(signature.signature, "base64").toString()
      await Bun.write(sigPath, sigContent)

      // Verify using GPG
      const result = await $`gpg --verify ${sigPath} ${filePath}`.quiet().nothrow()

      // Clean up
      await Bun.write(sigPath, "")

      if (result.exitCode !== 0) {
        return {
          valid: false,
          error: `GPG verification failed: ${result.stderr}`,
        }
      }

      // Check if signer is trusted
      const signer = trustedSigners.find(s => s.keyId === signature.keyId)

      if (!signer) {
        return {
          valid: false,
          error: `Signer ${signature.keyId} is not in trusted signers list`,
        }
      }

      // Check signature expiration
      const age = Date.now() - signature.timestamp
      const maxAge = 365 * 24 * 60 * 60 * 1000 // 365 days

      if (age > maxAge) {
        return {
          valid: false,
          error: `Signature expired (age: ${Math.floor(age / (24 * 60 * 60 * 1000))} days)`,
        }
      }

      log.info("signature verified", {
        path: filePath,
        signer: signer.name,
        keyId: signature.keyId,
      })

      return {
        valid: true,
        signer,
      }
    } catch (error) {
      log.error("verification failed", {
        path: filePath,
        error: error instanceof Error ? error.message : String(error),
      })

      return {
        valid: false,
        error: error instanceof Error ? error.message : String(error),
      }
    }
  }

  /**
   * Sign using Node.js crypto (alternative to GPG)
   * Useful for development and testing
   */
  export async function signWithCrypto(
    filePath: string,
    privateKeyPEM: string,
    algorithm: string = "RSA-SHA256",
    publicKeyPEM?: string
  ): Promise<Signature> {
    if (!existsSync(filePath)) {
      throw new Error(`File not found: ${filePath}`)
    }

    try {
      // Read file content
      const file = Bun.file(filePath)
      const content = await file.arrayBuffer()

      // Create signature
      const sign = crypto.createSign(algorithm)
      sign.update(Buffer.from(content))
      sign.end()

      const signature = sign.sign(privateKeyPEM, "base64")

      // Generate key ID from public key if provided, otherwise from private key
      const keySource = publicKeyPEM || privateKeyPEM
      const keyId = crypto
        .createHash("sha256")
        .update(keySource)
        .digest("hex")
        .substring(0, 16)
        .toUpperCase()

      log.info("file signed with crypto", { path: filePath, algorithm, keyId })

      return {
        algorithm,
        signature,
        keyId,
        timestamp: Date.now(),
        publicKey: publicKeyPEM,
      }
    } catch (error) {
      log.error("crypto signing failed", {
        path: filePath,
        error: error instanceof Error ? error.message : String(error),
      })
      throw error
    }
  }

  /**
   * Verify signature using Node.js crypto
   */
  export async function verifyCryptoSignature(
    filePath: string,
    signature: Signature,
    publicKeyPEM: string
  ): Promise<{
    valid: boolean
    error?: string
  }> {
    if (!existsSync(filePath)) {
      return {
        valid: false,
        error: `File not found: ${filePath}`,
      }
    }

    try {
      // Read file content
      const file = Bun.file(filePath)
      const content = await file.arrayBuffer()

      // Verify signature
      const verify = crypto.createVerify(signature.algorithm)
      verify.update(Buffer.from(content))
      verify.end()

      const valid = verify.verify(publicKeyPEM, signature.signature, "base64")

      if (valid) {
        log.info("crypto signature verified", { path: filePath })
      } else {
        log.warn("crypto signature invalid", { path: filePath })
      }

      return {
        valid,
        error: valid ? undefined : "Signature verification failed",
      }
    } catch (error) {
      log.error("crypto verification failed", {
        path: filePath,
        error: error instanceof Error ? error.message : String(error),
      })

      return {
        valid: false,
        error: error instanceof Error ? error.message : String(error),
      }
    }
  }

  /**
   * Generate RSA key pair for signing (development/testing)
   */
  export function generateKeyPair(): {
    privateKey: string
    publicKey: string
    keyId: string
  } {
    const { privateKey, publicKey } = crypto.generateKeyPairSync("rsa", {
      modulusLength: 2048,
      publicKeyEncoding: {
        type: "spki",
        format: "pem",
      },
      privateKeyEncoding: {
        type: "pkcs8",
        format: "pem",
      },
    })

    const keyId = crypto
      .createHash("sha256")
      .update(publicKey)
      .digest("hex")
      .substring(0, 16)
      .toUpperCase()

    log.info("key pair generated", { keyId })

    return {
      privateKey,
      publicKey,
      keyId,
    }
  }

  /**
   * Check if a signer is trusted
   */
  export function isTrustedSigner(
    keyId: string,
    trustedSigners: TrustedSigner[] = OFFICIAL_SIGNERS
  ): boolean {
    return trustedSigners.some(
      s => s.keyId === keyId && s.trustLevel === "full"
    )
  }

  /**
   * Get signer information
   */
  export function getSigner(
    keyId: string,
    trustedSigners: TrustedSigner[] = OFFICIAL_SIGNERS
  ): TrustedSigner | null {
    return trustedSigners.find(s => s.keyId === keyId) || null
  }
}
