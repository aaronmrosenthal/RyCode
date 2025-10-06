import { describe, test, expect, beforeEach } from "bun:test"
import { Password } from "../../src/util/password"

describe("Password Utility", () => {
  describe("hash", () => {
    test("should hash a valid password", async () => {
      const password = "MySecureP@ss123"
      const hash = await Password.hash(password)

      expect(hash).toBeString()
      expect(hash).toStartWith("$2b$12$")
      expect(hash.length).toBeGreaterThan(50)
    })

    test("should generate different hashes for same password", async () => {
      const password = "MySecureP@ss123"
      const hash1 = await Password.hash(password)
      const hash2 = await Password.hash(password)

      expect(hash1).not.toBe(hash2) // Different salts
    })

    test("should reject password that's too short", async () => {
      expect(async () => await Password.hash("Short1!")).toThrow()
    })

    test("should reject password missing uppercase", async () => {
      expect(async () => await Password.hash("mypassword1!")).toThrow()
    })

    test("should reject password missing lowercase", async () => {
      expect(async () => await Password.hash("MYPASSWORD1!")).toThrow()
    })

    test("should reject password missing number", async () => {
      expect(async () => await Password.hash("MyPassword!")).toThrow()
    })

    test("should reject password missing special character", async () => {
      expect(async () => await Password.hash("MyPassword123")).toThrow()
    })

    test("should reject password exceeding max length", async () => {
      const longPassword = "A".repeat(129) + "a1!"
      expect(async () => await Password.hash(longPassword)).toThrow()
    })
  })

  describe("verify", () => {
    test("should verify correct password", async () => {
      const password = "MySecureP@ss123"
      const hash = await Password.hash(password)
      const isValid = await Password.verify(password, hash)

      expect(isValid).toBe(true)
    })

    test("should reject incorrect password", async () => {
      const password = "MySecureP@ss123"
      const hash = await Password.hash(password)
      const isValid = await Password.verify("WrongP@ssword123", hash)

      expect(isValid).toBe(false)
    })

    test("should reject empty password", async () => {
      const password = "MySecureP@ss123"
      const hash = await Password.hash(password)
      const isValid = await Password.verify("", hash)

      expect(isValid).toBe(false)
    })

    test("should handle invalid hash gracefully", async () => {
      const isValid = await Password.verify("MySecureP@ss123", "invalid_hash")

      expect(isValid).toBe(false)
    })

    test("should be case-sensitive", async () => {
      const password = "MySecureP@ss123"
      const hash = await Password.hash(password)
      const isValid = await Password.verify("mysecurep@ss123", hash)

      expect(isValid).toBe(false)
    })
  })

  describe("needsRehash", () => {
    test("should return false for hash with current rounds", async () => {
      const password = "MySecureP@ss123"
      const hash = await Password.hash(password)
      const needs = await Password.needsRehash(hash)

      expect(needs).toBe(false)
    })

    test("should return true for hash with different rounds", async () => {
      // Simulate old hash with 10 rounds
      const oldHash = "$2b$10$KSGOsCl3aiSahBmx8z8dleL2pLHvP2wJpJ8xQ7X0Mz3Q4F5e6g7Ii"
      const needs = await Password.needsRehash(oldHash)

      expect(needs).toBe(true)
    })

    test("should return true for invalid hash format", async () => {
      const needs = await Password.needsRehash("invalid_hash")

      expect(needs).toBe(true)
    })
  })

  describe("validate", () => {
    test("should accept valid password", () => {
      const result = Password.validate("MySecureP@ss123")

      expect(result.valid).toBe(true)
      expect(result.errors).toHaveLength(0)
    })

    test("should reject password too short", () => {
      const result = Password.validate("Short1!")

      expect(result.valid).toBe(false)
      expect(result.errors).toContain("Password must be at least 8 characters")
    })

    test("should reject password too long", () => {
      const result = Password.validate("A".repeat(129) + "a1!")

      expect(result.valid).toBe(false)
      expect(result.errors).toContain("Password must not exceed 128 characters")
    })

    test("should list all validation errors", () => {
      const result = Password.validate("weak")

      expect(result.valid).toBe(false)
      expect(result.errors.length).toBeGreaterThan(3)
    })

    test("should accept password with all requirements", () => {
      const passwords = ["MyP@ss1234", "SecureP@ssw0rd", "Complex!Pass9"]

      passwords.forEach((password) => {
        const result = Password.validate(password)
        expect(result.valid).toBe(true)
      })
    })
  })

  describe("strength", () => {
    test("should score very weak password", () => {
      const result = Password.strength("12345")

      expect(result.score).toBeLessThanOrEqual(1)
      expect(result.feedback).toContain("weak")
    })

    test("should score weak password", () => {
      const result = Password.strength("Password1")

      expect(result.score).toBeGreaterThanOrEqual(1)
      expect(result.score).toBeLessThan(3)
    })

    test("should score strong password", () => {
      const result = Password.strength("MySecureP@ss123")

      expect(result.score).toBeGreaterThanOrEqual(3)
      expect(result.feedback).toMatch(/strong/i)
    })

    test("should score very strong password", () => {
      const result = Password.strength("MyVerySecureP@ssword123!")

      expect(result.score).toBe(4)
      expect(result.feedback).toContain("Very strong")
    })

    test("should penalize common patterns", () => {
      const commonPatterns = ["password123!", "qwerty123!", "abc12345!"]

      commonPatterns.forEach((password) => {
        const result = Password.strength(password)
        expect(result.score).toBeLessThan(3)
      })
    })

    test("should penalize repeating characters", () => {
      const result = Password.strength("Aaa111!!!")

      expect(result.score).toBeLessThanOrEqual(3)
      expect(result.feedback).toMatch(/repeating|variety/i)
    })
  })

  describe("schema", () => {
    test("should parse valid password", () => {
      expect(() => Password.schema.parse("MySecureP@ss123")).not.toThrow()
    })

    test("should reject invalid password with detailed errors", () => {
      try {
        Password.schema.parse("weak")
      } catch (error: any) {
        expect(error.issues).toBeDefined()
        expect(error.issues.length).toBeGreaterThan(0)
      }
    })

    test("should provide clear error messages", () => {
      const invalidPasswords = [
        { password: "short", expectedError: "at least 8 characters" },
        { password: "nouppercase1!", expectedError: "uppercase" },
        { password: "NOLOWERCASE1!", expectedError: "lowercase" },
        { password: "NoNumbers!", expectedError: "number" },
        { password: "NoSpecial123", expectedError: "special character" },
      ]

      invalidPasswords.forEach(({ password, expectedError }) => {
        try {
          Password.schema.parse(password)
          expect(true).toBe(false) // Should not reach here
        } catch (error: any) {
          const message = error.issues[0]?.message || ""
          expect(message.toLowerCase()).toContain(expectedError.toLowerCase())
        }
      })
    })
  })

  describe("edge cases", () => {
    test("should handle unicode characters", async () => {
      const password = "MyP@ss123\u00e9" // with accent
      const hash = await Password.hash(password)
      const isValid = await Password.verify(password, hash)

      expect(isValid).toBe(true)
    })

    test("should handle special characters in verification", async () => {
      const password = "Test!@#$%^&*()123Aa"
      const hash = await Password.hash(password)
      const isValid = await Password.verify(password, hash)

      expect(isValid).toBe(true)
    })

    test("should handle whitespace in password", async () => {
      const password = "My Secure P@ss 123"
      const hash = await Password.hash(password)
      const isValid = await Password.verify(password, hash)

      expect(isValid).toBe(true)
    })
  })
})
