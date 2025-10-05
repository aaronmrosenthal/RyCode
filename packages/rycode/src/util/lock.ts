export namespace Lock {
  // Default timeout: 30 seconds
  const DEFAULT_TIMEOUT_MS = 30_000

  const locks = new Map<
    string,
    {
      readers: number
      writer: boolean
      waitingReaders: Array<{ resolve: () => void; timeout: NodeJS.Timeout }>
      waitingWriters: Array<{ resolve: () => void; timeout: NodeJS.Timeout }>
      acquiredAt?: number
    }
  >()

  export class LockTimeoutError extends Error {
    constructor(
      public key: string,
      public timeoutMs: number,
    ) {
      super(`Lock timeout after ${timeoutMs}ms for key: ${key}`)
      this.name = "LockTimeoutError"
    }
  }

  function get(key: string) {
    if (!locks.has(key)) {
      locks.set(key, {
        readers: 0,
        writer: false,
        waitingReaders: [],
        waitingWriters: [],
        acquiredAt: undefined,
      })
    }
    return locks.get(key)!
  }

  function process(key: string) {
    const lock = locks.get(key)
    if (!lock || lock.writer || lock.readers > 0) return

    // Prioritize writers to prevent starvation
    if (lock.waitingWriters.length > 0) {
      const next = lock.waitingWriters.shift()!
      clearTimeout(next.timeout) // Clear timeout
      lock.acquiredAt = Date.now()
      next.resolve()
      return
    }

    // Wake up all waiting readers
    while (lock.waitingReaders.length > 0) {
      const next = lock.waitingReaders.shift()!
      clearTimeout(next.timeout) // Clear timeout
      lock.acquiredAt = Date.now()
      next.resolve()
    }

    // Clean up empty locks
    if (lock.readers === 0 && !lock.writer && lock.waitingReaders.length === 0 && lock.waitingWriters.length === 0) {
      locks.delete(key)
    }
  }

  export async function read(key: string, timeoutMs: number = DEFAULT_TIMEOUT_MS): Promise<Disposable> {
    const lock = get(key)

    return new Promise((resolve, reject) => {
      if (!lock.writer && lock.waitingWriters.length === 0) {
        lock.readers++
        lock.acquiredAt = Date.now()
        resolve({
          [Symbol.dispose]: () => {
            lock.readers--
            process(key)
          },
        })
      } else {
        // Set timeout for waiting readers
        const timeout = setTimeout(() => {
          // Remove from waiting queue
          const index = lock.waitingReaders.findIndex((w) => w.timeout === timeout)
          if (index !== -1) {
            lock.waitingReaders.splice(index, 1)
          }
          reject(new LockTimeoutError(key, timeoutMs))
        }, timeoutMs)

        lock.waitingReaders.push({
          resolve: () => {
            lock.readers++
            resolve({
              [Symbol.dispose]: () => {
                lock.readers--
                process(key)
              },
            })
          },
          timeout,
        })
      }
    })
  }

  export async function write(key: string, timeoutMs: number = DEFAULT_TIMEOUT_MS): Promise<Disposable> {
    const lock = get(key)

    return new Promise((resolve, reject) => {
      if (!lock.writer && lock.readers === 0) {
        lock.writer = true
        lock.acquiredAt = Date.now()
        resolve({
          [Symbol.dispose]: () => {
            lock.writer = false
            process(key)
          },
        })
      } else {
        // Set timeout for waiting writers
        const timeout = setTimeout(() => {
          // Remove from waiting queue
          const index = lock.waitingWriters.findIndex((w) => w.timeout === timeout)
          if (index !== -1) {
            lock.waitingWriters.splice(index, 1)
          }
          reject(new LockTimeoutError(key, timeoutMs))
        }, timeoutMs)

        lock.waitingWriters.push({
          resolve: () => {
            lock.writer = true
            resolve({
              [Symbol.dispose]: () => {
                lock.writer = false
                process(key)
              },
            })
          },
          timeout,
        })
      }
    })
  }

  /**
   * Get diagnostic info about current locks
   */
  export function diagnostics() {
    const result: Record<
      string,
      {
        readers: number
        writer: boolean
        waitingReaders: number
        waitingWriters: number
        acquiredAt?: number
        heldFor?: number
      }
    > = {}

    for (const [key, lock] of locks.entries()) {
      result[key] = {
        readers: lock.readers,
        writer: lock.writer,
        waitingReaders: lock.waitingReaders.length,
        waitingWriters: lock.waitingWriters.length,
        acquiredAt: lock.acquiredAt,
        heldFor: lock.acquiredAt ? Date.now() - lock.acquiredAt : undefined,
      }
    }

    return result
  }
}
