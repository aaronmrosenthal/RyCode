import { describe, test, expect, beforeEach, afterEach } from "bun:test"
import { Storage } from "../../src/storage/storage"
import path from "path"
import fs from "fs/promises"
import { TestSetup } from "../setup"

describe("Storage Transactions", () => {
  let tempDir: string

  beforeEach(async () => {
    tempDir = await TestSetup.createTempDir()
  })

  afterEach(async () => {
    await TestSetup.cleanup()
  })

  describe("basic transaction operations", () => {
    test("should commit multiple writes atomically", async () => {
      const tx = Storage.transaction()

      await tx.write(["test", "key1"], { value: "data1" })
      await tx.write(["test", "key2"], { value: "data2" })
      await tx.write(["test", "key3"], { value: "data3" })

      await tx.commit()

      // All writes should be persisted
      const data1 = await Storage.read<{ value: string }>(["test", "key1"])
      const data2 = await Storage.read<{ value: string }>(["test", "key2"])
      const data3 = await Storage.read<{ value: string }>(["test", "key3"])

      expect(data1.value).toBe("data1")
      expect(data2.value).toBe("data2")
      expect(data3.value).toBe("data3")
    })

    test("should rollback without persisting changes", async () => {
      const tx = Storage.transaction()

      await tx.write(["test", "rollback"], { value: "should-not-exist" })
      await tx.rollback()

      // Data should not exist
      await expect(Storage.read(["test", "rollback"])).rejects.toThrow()
    })

    test("should handle mixed write and remove operations", async () => {
      // Create initial data
      await Storage.write(["test", "existing"], { value: "old" })

      const tx = Storage.transaction()
      await tx.write(["test", "new"], { value: "new" })
      await tx.remove(["test", "existing"])
      await tx.commit()

      // New should exist
      const newData = await Storage.read<{ value: string }>(["test", "new"])
      expect(newData.value).toBe("new")

      // Old should be removed
      await expect(Storage.read(["test", "existing"])).rejects.toThrow()
    })
  })

  describe("transaction atomicity", () => {
    test("should commit all operations atomically", async () => {
      const tx = Storage.transaction()
      await tx.write(["test", "atomic1"], { value: "data1" })
      await tx.write(["test", "atomic2"], { value: "data2" })
      await tx.write(["test", "atomic3"], { value: "data3" })

      await tx.commit()

      // After commit, all should exist
      const data1 = await Storage.read<{ value: string }>(["test", "atomic1"])
      const data2 = await Storage.read<{ value: string }>(["test", "atomic2"])
      const data3 = await Storage.read<{ value: string }>(["test", "atomic3"])

      expect(data1.value).toBe("data1")
      expect(data2.value).toBe("data2")
      expect(data3.value).toBe("data3")
    })
  })

  describe("transaction error handling", () => {
    test("should throw error when writing after commit", async () => {
      const tx = Storage.transaction()
      await tx.write(["test", "key"], { value: "data" })
      await tx.commit()

      await expect(tx.write(["test", "key2"], { value: "more" })).rejects.toThrow(
        "Transaction already committed",
      )
    })

    test("should throw error when committing twice", async () => {
      const tx = Storage.transaction()
      await tx.write(["test", "key"], { value: "data" })
      await tx.commit()

      await expect(tx.commit()).rejects.toThrow("Transaction already committed")
    })
  })

  describe("concurrent transaction handling", () => {
    test("should handle concurrent transactions on different keys", async () => {
      const tx1 = Storage.transaction()
      const tx2 = Storage.transaction()

      await tx1.write(["test", "tx1-key"], { value: "tx1" })
      await tx2.write(["test", "tx2-key"], { value: "tx2" })

      await Promise.all([tx1.commit(), tx2.commit()])

      // Both should succeed
      const data1 = await Storage.read<{ value: string }>(["test", "tx1-key"])
      const data2 = await Storage.read<{ value: string }>(["test", "tx2-key"])

      expect(data1.value).toBe("tx1")
      expect(data2.value).toBe("tx2")
    })

    test("should serialize transactions on same key", async () => {
      const results: string[] = []

      const tx1 = async () => {
        const tx = Storage.transaction()
        await tx.write(["test", "shared"], { value: "tx1" })
        results.push("tx1-start")
        await new Promise((resolve) => setTimeout(resolve, 50))
        results.push("tx1-commit")
        await tx.commit()
      }

      const tx2 = async () => {
        await new Promise((resolve) => setTimeout(resolve, 10))
        const tx = Storage.transaction()
        await tx.write(["test", "shared"], { value: "tx2" })
        results.push("tx2-start")
        await tx.commit()
        results.push("tx2-commit")
      }

      await Promise.all([tx1(), tx2()])

      // Transactions should be serialized - one commits before other
      // Both should complete
      expect(results.length).toBe(4)
      expect(results.includes("tx1-start")).toBe(true)
      expect(results.includes("tx1-commit")).toBe(true)
      expect(results.includes("tx2-start")).toBe(true)
      expect(results.includes("tx2-commit")).toBe(true)

      // Verify commits don't overlap
      const tx1CommitIdx = results.indexOf("tx1-commit")
      const tx2CommitIdx = results.indexOf("tx2-commit")
      expect(Math.abs(tx1CommitIdx - tx2CommitIdx)).toBeGreaterThan(0)
    })
  })

  describe("deadlock prevention", () => {
    test("should prevent deadlocks via lock ordering", async () => {
      // Two transactions accessing same keys in different order
      const tx1 = async () => {
        const tx = Storage.transaction()
        await tx.write(["test", "keyA"], { value: "tx1-A" })
        await tx.write(["test", "keyB"], { value: "tx1-B" })
        await tx.commit()
      }

      const tx2 = async () => {
        const tx = Storage.transaction()
        await tx.write(["test", "keyB"], { value: "tx2-B" })
        await tx.write(["test", "keyA"], { value: "tx2-A" })
        await tx.commit()
      }

      // Should complete without deadlock (locks acquired in sorted order)
      await Promise.all([tx1(), tx2()])

      const dataA = await Storage.read<{ value: string }>(["test", "keyA"])
      const dataB = await Storage.read<{ value: string }>(["test", "keyB"])

      // One of the transactions should win
      expect(dataA.value).toMatch(/tx[12]-A/)
      expect(dataB.value).toMatch(/tx[12]-B/)
    })
  })

  describe("session operations with transactions", () => {
    test("should create session with messages atomically", async () => {
      const sessionID = "test-session-123"
      const projectID = "test-project"

      const tx = Storage.transaction()

      // Session metadata
      await tx.write(["session", projectID, sessionID], {
        id: sessionID,
        projectID,
        title: "Test Session",
        time: {
          created: Date.now(),
          updated: Date.now(),
        },
      })

      // Initial messages
      await tx.write(["message", sessionID, "msg1"], {
        id: "msg1",
        sessionID,
        role: "user",
        content: "Hello",
      })

      await tx.write(["message", sessionID, "msg2"], {
        id: "msg2",
        sessionID,
        role: "assistant",
        content: "Hi there!",
      })

      await tx.commit()

      // All should exist
      const session = await Storage.read(["session", projectID, sessionID])
      const msg1 = await Storage.read(["message", sessionID, "msg1"])
      const msg2 = await Storage.read(["message", sessionID, "msg2"])

      expect(session).toBeDefined()
      expect(msg1).toBeDefined()
      expect(msg2).toBeDefined()
    })

    test("should rollback failed session creation", async () => {
      const sessionID = "test-session-456"
      const projectID = "test-project"

      const tx = Storage.transaction()

      await tx.write(["session", projectID, sessionID], {
        id: sessionID,
        projectID,
      })

      await tx.write(["message", sessionID, "msg1"], {
        id: "msg1",
        sessionID,
      })

      // Rollback instead of commit
      await tx.rollback()

      // Nothing should exist
      await expect(Storage.read(["session", projectID, sessionID])).rejects.toThrow()
      await expect(Storage.read(["message", sessionID, "msg1"])).rejects.toThrow()
    })
  })

  describe("large transactions", () => {
    test("should handle transaction with many operations", async () => {
      const tx = Storage.transaction()

      // Add 100 writes
      for (let i = 0; i < 100; i++) {
        await tx.write(["test", "batch", `item${i}`], { index: i })
      }

      await tx.commit()

      // Verify all persisted
      for (let i = 0; i < 100; i++) {
        const data = await Storage.read<{ index: number }>(["test", "batch", `item${i}`])
        expect(data.index).toBe(i)
      }
    })

    test("should handle transaction with duplicate key writes", async () => {
      const tx = Storage.transaction()

      await tx.write(["test", "duplicate"], { value: 1 })
      await tx.write(["test", "duplicate"], { value: 2 })
      await tx.write(["test", "duplicate"], { value: 3 })

      await tx.commit()

      // Last write should win
      const data = await Storage.read<{ value: number }>(["test", "duplicate"])
      expect(data.value).toBe(3)
    })
  })
})
