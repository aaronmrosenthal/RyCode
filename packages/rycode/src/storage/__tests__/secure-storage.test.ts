import { describe, test, expect } from "bun:test"
import { SecureStorage } from "../secure-storage"

describe("SecureStorage", () => {
  const testKey = SecureStorage.generateKey()
  const testData = JSON.stringify({
    apiKey: "super-secret-key-12345",
    token: "access-token-67890",
  })

  describe("encryption and decryption", () => {
    test("encrypts and decrypts data correctly", async () => {
      const encrypted = await SecureStorage.encrypt(testData, testKey)
      const decrypted = await SecureStorage.decrypt(encrypted, testKey)

      expect(decrypted).toBe(testData)
    })

    test("encrypted data is not plaintext", async () => {
      const encrypted = await SecureStorage.encrypt(testData, testKey)

      expect(encrypted).not.toBe(testData)
      expect(encrypted).not.toContain("super-secret-key")
      expect(encrypted).not.toContain("access-token")
    })

    test("encrypted data has correct format", async () => {
      const encrypted = await SecureStorage.encrypt(testData, testKey)

      // Format: salt:iv:authTag:encryptedData (all hex)
      const parts = encrypted.split(":")
      expect(parts).toHaveLength(4)

      // Check all parts are valid hex
      for (const part of parts) {
        expect(/^[0-9a-fA-F]+$/.test(part)).toBe(true)
      }
    })

    test("different encryptions produce different outputs (randomized IV)", async () => {
      const encrypted1 = await SecureStorage.encrypt(testData, testKey)
      const encrypted2 = await SecureStorage.encrypt(testData, testKey)

      expect(encrypted1).not.toBe(encrypted2)

      // But both decrypt to same plaintext
      const decrypted1 = await SecureStorage.decrypt(encrypted1, testKey)
      const decrypted2 = await SecureStorage.decrypt(encrypted2, testKey)
      expect(decrypted1).toBe(testData)
      expect(decrypted2).toBe(testData)
    })
  })

  describe("key validation", () => {
    test("rejects decryption with wrong key", async () => {
      const encrypted = await SecureStorage.encrypt(testData, testKey)
      const wrongKey = SecureStorage.generateKey()

      await expect(SecureStorage.decrypt(encrypted, wrongKey)).rejects.toThrow("Failed to decrypt")
    })

    test("rejects decryption with no key", async () => {
      const encrypted = await SecureStorage.encrypt(testData, testKey)

      // Clear environment variable for this test
      const originalKey = process.env['RYCODE_ENCRYPTION_KEY']
      delete process.env['RYCODE_ENCRYPTION_KEY']

      await expect(SecureStorage.decrypt(encrypted)).rejects.toThrow("RYCODE_ENCRYPTION_KEY")

      // Restore
      if (originalKey) process.env['RYCODE_ENCRYPTION_KEY'] = originalKey
    })
  })

  describe("isEncrypted", () => {
    test("detects encrypted data", async () => {
      const encrypted = await SecureStorage.encrypt(testData, testKey)

      expect(SecureStorage.isEncrypted(encrypted)).toBe(true)
    })

    test("detects plaintext data", () => {
      expect(SecureStorage.isEncrypted("plaintext:" + testData)).toBe(false)
      expect(SecureStorage.isEncrypted(testData)).toBe(false)
      expect(SecureStorage.isEncrypted("not:encrypted:data")).toBe(false)
    })
  })

  describe("plaintext fallback", () => {
    test("handles plaintext marker for backward compatibility", async () => {
      const plaintext = "plaintext:" + testData

      const decrypted = await SecureStorage.decrypt(plaintext, testKey)

      expect(decrypted).toBe(testData)
    })

    test("encrypts without key creates plaintext", async () => {
      // Clear environment variable
      const originalKey = process.env['RYCODE_ENCRYPTION_KEY']
      delete process.env['RYCODE_ENCRYPTION_KEY']

      const encrypted = await SecureStorage.encrypt(testData)

      expect(encrypted).toStartWith("plaintext:")
      expect(SecureStorage.isEncrypted(encrypted)).toBe(false)

      // Restore
      if (originalKey) process.env['RYCODE_ENCRYPTION_KEY'] = originalKey
    })
  })

  describe("reencrypt", () => {
    test("re-encrypts plaintext data", async () => {
      const plaintext = "plaintext:" + testData

      const encrypted = await SecureStorage.reencrypt(plaintext, testKey)

      expect(SecureStorage.isEncrypted(encrypted)).toBe(true)

      const decrypted = await SecureStorage.decrypt(encrypted, testKey)
      expect(decrypted).toBe(testData)
    })

    test("leaves encrypted data unchanged", async () => {
      const encrypted = await SecureStorage.encrypt(testData, testKey)

      const reencrypted = await SecureStorage.reencrypt(encrypted, testKey)

      expect(reencrypted).toBe(encrypted)
    })
  })

  describe("key generation", () => {
    test("generates valid encryption keys", () => {
      const key = SecureStorage.generateKey()

      expect(key).toBeTruthy()
      expect(SecureStorage.isValidKey(key)).toBe(true)
    })

    test("generated keys are different", () => {
      const key1 = SecureStorage.generateKey()
      const key2 = SecureStorage.generateKey()

      expect(key1).not.toBe(key2)
    })

    test("validates key format", () => {
      const validKey = SecureStorage.generateKey()
      expect(SecureStorage.isValidKey(validKey)).toBe(true)

      expect(SecureStorage.isValidKey("too-short")).toBe(false)
      expect(SecureStorage.isValidKey("not-base64!@#$%")).toBe(false)
    })
  })

  describe("tampering detection", () => {
    test("detects tampering via auth tag", async () => {
      const encrypted = await SecureStorage.encrypt(testData, testKey)

      // Tamper with encrypted data
      const parts = encrypted.split(":")
      parts[3] = "00" + parts[3].substring(2) // Modify encrypted data

      const tampered = parts.join(":")

      await expect(SecureStorage.decrypt(tampered, testKey)).rejects.toThrow("Failed to decrypt")
    })

    test("detects invalid format", async () => {
      const invalidFormat = "not:valid:format"

      await expect(SecureStorage.decrypt(invalidFormat, testKey)).rejects.toThrow()
    })
  })
})
