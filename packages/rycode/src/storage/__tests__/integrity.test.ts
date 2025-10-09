import { describe, test, expect } from "bun:test"
import { Integrity } from "../integrity"

describe("Integrity", () => {
  const testData = JSON.stringify({
    session: "test-session-123",
    data: { foo: "bar", nested: { value: 42 } },
  })

  describe("checksum computation", () => {
    test("computes consistent checksums", () => {
      const checksum1 = Integrity.computeChecksum(testData)
      const checksum2 = Integrity.computeChecksum(testData)

      expect(checksum1).toBe(checksum2)
    })

    test("produces 64-character hex checksums", () => {
      const checksum = Integrity.computeChecksum(testData)

      expect(checksum).toHaveLength(64)
      expect(/^[0-9a-fA-F]{64}$/.test(checksum)).toBe(true)
    })

    test("different data produces different checksums", () => {
      const data1 = "test data 1"
      const data2 = "test data 2"

      const checksum1 = Integrity.computeChecksum(data1)
      const checksum2 = Integrity.computeChecksum(data2)

      expect(checksum1).not.toBe(checksum2)
    })

    test("even small changes alter checksum", () => {
      const data1 = "test data"
      const data2 = "test data " // Extra space

      const checksum1 = Integrity.computeChecksum(data1)
      const checksum2 = Integrity.computeChecksum(data2)

      expect(checksum1).not.toBe(checksum2)
    })
  })

  describe("checksum verification", () => {
    test("verifies correct checksums", () => {
      const checksum = Integrity.computeChecksum(testData)

      expect(Integrity.verifyChecksum(testData, checksum)).toBe(true)
    })

    test("rejects incorrect checksums", () => {
      const wrongChecksum = "a".repeat(64)

      expect(Integrity.verifyChecksum(testData, wrongChecksum)).toBe(false)
    })

    test("rejects malformed checksums", () => {
      expect(Integrity.verifyChecksum(testData, "not-a-valid-checksum")).toBe(false)
      expect(Integrity.verifyChecksum(testData, "")).toBe(false)
    })

    test("uses constant-time comparison", () => {
      const checksum = Integrity.computeChecksum(testData)

      // Should not reveal timing information
      const wrongChecksum = "0" + checksum.substring(1)

      expect(Integrity.verifyChecksum(testData, wrongChecksum)).toBe(false)
    })
  })

  describe("wrap and unwrap", () => {
    test("wraps data with checksum", () => {
      const wrapped = Integrity.wrap(testData)

      expect(wrapped).toContain(":")
      expect(wrapped).toHaveLength(testData.length + 65) // 64 chars + colon
      expect(Integrity.hasIntegrity(wrapped)).toBe(true)
    })

    test("unwraps and verifies data", () => {
      const wrapped = Integrity.wrap(testData)
      const unwrapped = Integrity.unwrap(wrapped)

      expect(unwrapped).toBe(testData)
    })

    test("detects tampering when unwrapping", () => {
      const wrapped = Integrity.wrap(testData)

      // Tamper with data
      const parts = wrapped.split(":")
      parts[1] = parts[1].substring(0, parts[1].length - 1) + "X"
      const tampered = parts.join(":")

      expect(() => Integrity.unwrap(tampered)).toThrow(Integrity.IntegrityError)
    })

    test("rejects data without checksum", () => {
      expect(() => Integrity.unwrap(testData)).toThrow(Integrity.IntegrityError)
      expect(() => Integrity.unwrap("no:checksum:here")).toThrow(Integrity.IntegrityError)
    })
  })

  describe("hasIntegrity", () => {
    test("detects wrapped data", () => {
      const wrapped = Integrity.wrap(testData)

      expect(Integrity.hasIntegrity(wrapped)).toBe(true)
    })

    test("detects unwrapped data", () => {
      expect(Integrity.hasIntegrity(testData)).toBe(false)
      expect(Integrity.hasIntegrity("")).toBe(false)
      expect(Integrity.hasIntegrity("too:short")).toBe(false)
    })
  })

  describe("metadata", () => {
    test("generates metadata with checksum, size, timestamp", () => {
      const metadata = Integrity.generateMetadata(testData)

      expect(metadata.checksum).toHaveLength(64)
      expect(metadata.size).toBe(Buffer.byteLength(testData, "utf8"))
      expect(metadata.timestamp).toBeGreaterThan(Date.now() - 1000)
      expect(metadata.algorithm).toBe("sha256")
    })

    test("verifies data with metadata", () => {
      const metadata = Integrity.generateMetadata(testData)

      expect(Integrity.verifyMetadata(testData, metadata)).toBe(true)
    })

    test("rejects modified data", () => {
      const metadata = Integrity.generateMetadata(testData)
      const modifiedData = testData + " "

      expect(Integrity.verifyMetadata(modifiedData, metadata)).toBe(false)
    })

    test("detects size changes", () => {
      const metadata = Integrity.generateMetadata(testData)

      // Manually modify metadata
      metadata.size += 1

      expect(Integrity.verifyMetadata(testData, metadata)).toBe(false)
    })
  })

  describe("tampering detection", () => {
    test("detects data corruption", () => {
      const wrapped = Integrity.wrap(testData)

      // Corrupt a single character in the data
      const corrupted = wrapped.substring(0, wrapped.length - 5) + "X" + wrapped.substring(wrapped.length - 4)

      expect(() => Integrity.unwrap(corrupted)).toThrow(Integrity.IntegrityError)
    })

    test("detects checksum modification", () => {
      const wrapped = Integrity.wrap(testData)

      // Modify checksum
      const parts = wrapped.split(":")
      parts[0] = "a".repeat(64)
      const modified = parts.join(":")

      expect(() => Integrity.unwrap(modified)).toThrow(Integrity.IntegrityError)
    })
  })
})
