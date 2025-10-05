import { describe, test, expect } from "bun:test"
import { Lock } from "../../src/util/lock"

describe("Lock", () => {
  describe("basic read locks", () => {
    test("should allow multiple concurrent readers", async () => {
      const results: number[] = []

      const reader1 = async () => {
        using _ = await Lock.read("test-key")
        results.push(1)
        await new Promise((resolve) => setTimeout(resolve, 50))
        results.push(2)
      }

      const reader2 = async () => {
        using _ = await Lock.read("test-key")
        results.push(3)
        await new Promise((resolve) => setTimeout(resolve, 50))
        results.push(4)
      }

      await Promise.all([reader1(), reader2()])

      // Both readers should have executed concurrently
      // Results should be interleaved, not sequential
      expect(results.length).toBe(4)
      expect(results.includes(1)).toBe(true)
      expect(results.includes(2)).toBe(true)
      expect(results.includes(3)).toBe(true)
      expect(results.includes(4)).toBe(true)
    })

    test("should release read lock properly", async () => {
      let locked = false

      {
        using _ = await Lock.read("test-key")
        locked = true
      }

      // Lock should be released
      const diagnostics = Lock.diagnostics()
      expect(diagnostics["test-key"]).toBeUndefined() // Cleaned up
    })
  })

  describe("basic write locks", () => {
    test("should enforce exclusive write access", async () => {
      const results: string[] = []

      const writer1 = async () => {
        using _ = await Lock.write("test-key")
        results.push("w1-start")
        await new Promise((resolve) => setTimeout(resolve, 50))
        results.push("w1-end")
      }

      const writer2 = async () => {
        using _ = await Lock.write("test-key")
        results.push("w2-start")
        await new Promise((resolve) => setTimeout(resolve, 50))
        results.push("w2-end")
      }

      await Promise.all([writer1(), writer2()])

      // Writers should be sequential, not concurrent
      expect(results).toEqual(["w1-start", "w1-end", "w2-start", "w2-end"])
    })

    test("should block readers while writer holds lock", async () => {
      const results: string[] = []

      const writer = async () => {
        using _ = await Lock.write("test-key")
        results.push("writer-start")
        await new Promise((resolve) => setTimeout(resolve, 100))
        results.push("writer-end")
      }

      const reader = async () => {
        await new Promise((resolve) => setTimeout(resolve, 20)) // Start after writer
        using _ = await Lock.read("test-key")
        results.push("reader")
      }

      await Promise.all([writer(), reader()])

      // Reader should wait for writer
      expect(results).toEqual(["writer-start", "writer-end", "reader"])
    })

    test("should block writer while readers hold lock", async () => {
      const results: string[] = []

      const reader = async () => {
        using _ = await Lock.read("test-key")
        results.push("reader-start")
        await new Promise((resolve) => setTimeout(resolve, 100))
        results.push("reader-end")
      }

      const writer = async () => {
        await new Promise((resolve) => setTimeout(resolve, 20)) // Start after reader
        using _ = await Lock.write("test-key")
        results.push("writer")
      }

      await Promise.all([reader(), writer()])

      // Writer should wait for reader
      expect(results).toEqual(["reader-start", "reader-end", "writer"])
    })
  })

  describe("timeout support", () => {
    test("should timeout read lock after specified duration", async () => {
      // Hold a write lock
      const writeLock = await Lock.write("test-key")

      try {
        // Try to acquire read lock with short timeout
        await Lock.read("test-key", 100)
        expect(true).toBe(false) // Should not reach here
      } catch (error) {
        expect(error).toBeInstanceOf(Lock.LockTimeoutError)
        expect((error as Lock.LockTimeoutError).key).toBe("test-key")
        expect((error as Lock.LockTimeoutError).timeoutMs).toBe(100)
      } finally {
        writeLock[Symbol.dispose]()
      }
    })

    test("should timeout write lock after specified duration", async () => {
      // Hold a read lock
      const readLock = await Lock.read("test-key")

      try {
        // Try to acquire write lock with short timeout
        await Lock.write("test-key", 100)
        expect(true).toBe(false) // Should not reach here
      } catch (error) {
        expect(error).toBeInstanceOf(Lock.LockTimeoutError)
        expect((error as Lock.LockTimeoutError).key).toBe("test-key")
      } finally {
        readLock[Symbol.dispose]()
      }
    })

    test("should not timeout if lock acquired before timeout", async () => {
      const writeLock = await Lock.write("test-key")

      // Release after 50ms
      setTimeout(() => writeLock[Symbol.dispose](), 50)

      // Try to acquire with 200ms timeout (should succeed)
      using _ = await Lock.read("test-key", 200)

      expect(true).toBe(true) // Should reach here
    })

    test("should clean up timed out waiters from queue", async () => {
      const writeLock = await Lock.write("test-key")

      try {
        await Lock.read("test-key", 50)
      } catch (error) {
        expect(error).toBeInstanceOf(Lock.LockTimeoutError)
      }

      // Check diagnostics - timed out waiter should be removed
      const diag = Lock.diagnostics()
      expect(diag["test-key"].waitingReaders).toBe(0)

      writeLock[Symbol.dispose]()
    })
  })

  describe("concurrent operations", () => {
    test("should handle many concurrent readers without deadlock", async () => {
      const readers = Array(100)
        .fill(0)
        .map(async (__, i) => {
          using _lock = await Lock.read("test-key")
          await new Promise((resolve) => setTimeout(resolve, 1))
          return i
        })

      const results = await Promise.all(readers)
      expect(results.length).toBe(100)
    })

    test("should handle many concurrent writers without deadlock", async () => {
      const counter = { value: 0 }

      const writers = Array(50)
        .fill(0)
        .map(async () => {
          using _ = await Lock.write("test-key")
          counter.value++
          await new Promise((resolve) => setTimeout(resolve, 1))
        })

      await Promise.all(writers)
      expect(counter.value).toBe(50)
    })

    test("should handle mixed readers and writers", async () => {
      const results: string[] = []

      await Promise.all([
        (async () => {
          using _ = await Lock.read("test-key")
          results.push("r1")
        })(),
        (async () => {
          using _ = await Lock.write("test-key")
          results.push("w1")
        })(),
        (async () => {
          using _ = await Lock.read("test-key")
          results.push("r2")
        })(),
        (async () => {
          using _ = await Lock.write("test-key")
          results.push("w2")
        })(),
        (async () => {
          using _ = await Lock.read("test-key")
          results.push("r3")
        })(),
      ])

      expect(results.length).toBe(5)
      // All operations should complete
      expect(results.includes("r1")).toBe(true)
      expect(results.includes("w1")).toBe(true)
      expect(results.includes("r2")).toBe(true)
      expect(results.includes("w2")).toBe(true)
      expect(results.includes("r3")).toBe(true)
    })
  })

  describe("writer priority", () => {
    test("should prioritize waiting writers over waiting readers", async () => {
      const results: string[] = []

      // Start with a writer holding the lock
      const initialWriter = await Lock.write("test-key")

      // Queue up a reader and a writer (writer queued after reader)
      const readerPromise = (async () => {
        using _ = await Lock.read("test-key")
        results.push("reader")
      })()

      await new Promise((resolve) => setTimeout(resolve, 10))

      const writerPromise = (async () => {
        using _ = await Lock.write("test-key")
        results.push("writer")
      })()

      await new Promise((resolve) => setTimeout(resolve, 10))

      // Release initial writer
      initialWriter[Symbol.dispose]()

      await Promise.all([readerPromise, writerPromise])

      // Writer should go first due to priority
      expect(results[0]).toBe("writer")
      expect(results[1]).toBe("reader")
    })
  })

  describe("diagnostics", () => {
    test("should report current lock state", async () => {
      using w = await Lock.write("key1")
      using r1 = await Lock.read("key2")
      using r2 = await Lock.read("key2")

      const diag = Lock.diagnostics()

      expect(diag["key1"]).toBeDefined()
      expect(diag["key1"].writer).toBe(true)
      expect(diag["key1"].readers).toBe(0)

      expect(diag["key2"]).toBeDefined()
      expect(diag["key2"].writer).toBe(false)
      expect(diag["key2"].readers).toBe(2)
    })

    test("should report lock hold time", async () => {
      using _ = await Lock.write("test-key")

      await new Promise((resolve) => setTimeout(resolve, 100))

      const diag = Lock.diagnostics()
      expect(diag["test-key"].heldFor).toBeGreaterThanOrEqual(100)
    })

    test("should report waiting counts", async () => {
      const writeLock = await Lock.write("test-key")

      // Queue up some waiters (don't await)
      Lock.read("test-key", 5000).catch(() => {})
      Lock.write("test-key", 5000).catch(() => {})

      await new Promise((resolve) => setTimeout(resolve, 10))

      const diag = Lock.diagnostics()
      expect(diag["test-key"].waitingReaders).toBe(1)
      expect(diag["test-key"].waitingWriters).toBe(1)

      writeLock[Symbol.dispose]()
    })
  })

  describe("edge cases", () => {
    test("should handle rapid lock acquire/release cycles", async () => {
      for (let i = 0; i < 100; i++) {
        {
          using _ = await Lock.write("test-key-rapid")
        } // Explicit scope to ensure dispose
      }

      // Should not leak locks
      const diag = Lock.diagnostics()
      expect(diag["test-key-rapid"]).toBeUndefined()
    })

    test("should handle disposing same lock multiple times safely", async () => {
      const lock = await Lock.write("test-key-dispose")

      lock[Symbol.dispose]()
      // Second dispose should not crash (dispose is idempotent)
      lock[Symbol.dispose]()

      expect(true).toBe(true)
    })
  })
})
