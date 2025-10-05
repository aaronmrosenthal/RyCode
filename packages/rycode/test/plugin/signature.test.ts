import { describe, expect, test } from "bun:test"
import { PluginSignature } from "../../src/plugin/signature"
import path from "path"
import { tmpdir } from "../fixture/fixture"

/**
 * Tests for Plugin Signature System
 */

describe("PluginSignature", () => {
  describe("Key Generation", () => {
    test("should generate RSA key pair", () => {
      const keyPair = PluginSignature.generateKeyPair()

      expect(keyPair.privateKey).toContain("BEGIN PRIVATE KEY")
      expect(keyPair.publicKey).toContain("BEGIN PUBLIC KEY")
      expect(keyPair.keyId).toMatch(/^[A-F0-9]{16}$/)
    })

    test("should generate unique keys", () => {
      const keyPair1 = PluginSignature.generateKeyPair()
      const keyPair2 = PluginSignature.generateKeyPair()

      expect(keyPair1.keyId).not.toBe(keyPair2.keyId)
      expect(keyPair1.privateKey).not.toBe(keyPair2.privateKey)
    })
  })

  describe("Crypto Signing", () => {
    test("should sign a file with crypto", async () => {
      await using tmp = await tmpdir()
      const testFile = path.join(tmp.path, "test-plugin.js")
      await Bun.write(testFile, "console.log('test')")

      const keyPair = PluginSignature.generateKeyPair()

      const signature = await PluginSignature.signWithCrypto(
        testFile,
        keyPair.privateKey,
        "RSA-SHA256",
        keyPair.publicKey
      )

      expect(signature.algorithm).toBe("RSA-SHA256")
      expect(signature.signature).toBeTruthy()
      expect(signature.keyId).toBe(keyPair.keyId)
      expect(signature.timestamp).toBeGreaterThan(0)
    })

    test("should verify valid signature", async () => {
      await using tmp = await tmpdir()
      const testFile = path.join(tmp.path, "test-plugin.js")
      await Bun.write(testFile, "console.log('verify me')")

      const keyPair = PluginSignature.generateKeyPair()

      const signature = await PluginSignature.signWithCrypto(
        testFile,
        keyPair.privateKey,
        "RSA-SHA256",
        keyPair.publicKey
      )

      const result = await PluginSignature.verifyCryptoSignature(
        testFile,
        signature,
        keyPair.publicKey
      )

      expect(result.valid).toBe(true)
      expect(result.error).toBeUndefined()
    })

    test("should fail verification with wrong public key", async () => {
      await using tmp = await tmpdir()
      const testFile = path.join(tmp.path, "test-plugin.js")
      await Bun.write(testFile, "console.log('test')")

      const keyPair1 = PluginSignature.generateKeyPair()
      const keyPair2 = PluginSignature.generateKeyPair()

      const signature = await PluginSignature.signWithCrypto(
        testFile,
        keyPair1.privateKey
      )

      const result = await PluginSignature.verifyCryptoSignature(
        testFile,
        signature,
        keyPair2.publicKey // Wrong key!
      )

      expect(result.valid).toBe(false)
      expect(result.error).toBeTruthy()
    })

    test("should fail verification if file is modified", async () => {
      await using tmp = await tmpdir()
      const testFile = path.join(tmp.path, "test-plugin.js")
      await Bun.write(testFile, "console.log('original')")

      const keyPair = PluginSignature.generateKeyPair()

      const signature = await PluginSignature.signWithCrypto(
        testFile,
        keyPair.privateKey,
        "RSA-SHA256",
        keyPair.publicKey
      )

      // Modify the file
      await Bun.write(testFile, "console.log('tampered')")

      const result = await PluginSignature.verifyCryptoSignature(
        testFile,
        signature,
        keyPair.publicKey
      )

      expect(result.valid).toBe(false)
    })

    test("should support different algorithms", async () => {
      await using tmp = await tmpdir()
      const testFile = path.join(tmp.path, "test-plugin.js")
      await Bun.write(testFile, "console.log('algorithm test')")

      const keyPair = PluginSignature.generateKeyPair()

      const sig256 = await PluginSignature.signWithCrypto(
        testFile,
        keyPair.privateKey,
        "RSA-SHA256",
        keyPair.publicKey
      )

      const sig512 = await PluginSignature.signWithCrypto(
        testFile,
        keyPair.privateKey,
        "RSA-SHA512",
        keyPair.publicKey
      )

      expect(sig256.algorithm).toBe("RSA-SHA256")
      expect(sig512.algorithm).toBe("RSA-SHA512")

      // Both should verify successfully
      const result256 = await PluginSignature.verifyCryptoSignature(
        testFile,
        sig256,
        keyPair.publicKey
      )

      const result512 = await PluginSignature.verifyCryptoSignature(
        testFile,
        sig512,
        keyPair.publicKey
      )

      expect(result256.valid).toBe(true)
      expect(result512.valid).toBe(true)
    })
  })

  describe("Trusted Signers", () => {
    test("should check if signer is trusted", () => {
      const trustedSigners: PluginSignature.TrustedSigner[] = [
        {
          name: "Test Signer",
          keyId: "TEST123",
          publicKey: "test-key",
          trustLevel: "full",
        },
      ]

      expect(PluginSignature.isTrustedSigner("TEST123", trustedSigners)).toBe(true)
      expect(PluginSignature.isTrustedSigner("UNKNOWN", trustedSigners)).toBe(false)
    })

    test("should get signer information", () => {
      const trustedSigners: PluginSignature.TrustedSigner[] = [
        {
          name: "Test Signer",
          keyId: "TEST123",
          publicKey: "test-key",
          organization: "Test Org",
          trustLevel: "full",
        },
      ]

      const signer = PluginSignature.getSigner("TEST123", trustedSigners)

      expect(signer).not.toBeNull()
      expect(signer?.name).toBe("Test Signer")
      expect(signer?.organization).toBe("Test Org")
    })

    test("should not trust marginal or never signers", () => {
      const trustedSigners: PluginSignature.TrustedSigner[] = [
        {
          name: "Marginal Signer",
          keyId: "MARGINAL",
          publicKey: "test-key",
          trustLevel: "marginal",
        },
        {
          name: "Never Signer",
          keyId: "NEVER",
          publicKey: "test-key",
          trustLevel: "never",
        },
      ]

      expect(PluginSignature.isTrustedSigner("MARGINAL", trustedSigners)).toBe(false)
      expect(PluginSignature.isTrustedSigner("NEVER", trustedSigners)).toBe(false)
    })
  })

  describe("GPG Availability", () => {
    test("should check if GPG is available", async () => {
      const available = await PluginSignature.isGPGAvailable()
      // Just check it returns a boolean - actual availability depends on system
      expect(typeof available).toBe("boolean")
    })
  })
})
