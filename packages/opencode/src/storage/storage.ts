import { Log } from "../util/log"
import path from "path"
import fs from "fs/promises"
import { Global } from "../global"
import { lazy } from "../util/lazy"
import { Lock } from "../util/lock"
import { $ } from "bun"

export namespace Storage {
  const log = Log.create({ service: "storage" })

  type Migration = (dir: string) => Promise<void>

  const MIGRATIONS: Migration[] = [
    async (dir) => {
      const project = path.resolve(dir, "../project")
      for await (const projectDir of new Bun.Glob("*").scan({
        cwd: project,
        onlyFiles: false,
      })) {
        log.info(`migrating project ${projectDir}`)
        let projectID = projectDir
        const fullProjectDir = path.join(project, projectDir)
        let worktree = "/"

        if (projectID !== "global") {
          for await (const msgFile of new Bun.Glob("storage/session/message/*/*.json").scan({
            cwd: path.join(project, projectDir),
            absolute: true,
          })) {
            const json = await Bun.file(msgFile).json()
            worktree = json.path?.root
            if (worktree) break
          }
          if (!worktree) continue
          if (!(await fs.exists(worktree))) continue
          const [id] = await $`git rev-list --max-parents=0 --all`
            .quiet()
            .nothrow()
            .cwd(worktree)
            .text()
            .then((x) =>
              x
                .split("\n")
                .filter(Boolean)
                .map((x) => x.trim())
                .toSorted(),
            )
          if (!id) continue
          projectID = id

          await Bun.write(
            path.join(dir, "project", projectID + ".json"),
            JSON.stringify({
              id,
              vcs: "git",
              worktree,
              time: {
                created: Date.now(),
                initialized: Date.now(),
              },
            }),
          )

          log.info(`migrating sessions for project ${projectID}`)
          for await (const sessionFile of new Bun.Glob("storage/session/info/*.json").scan({
            cwd: fullProjectDir,
            absolute: true,
          })) {
            const dest = path.join(dir, "session", projectID, path.basename(sessionFile))
            log.info("copying", {
              sessionFile,
              dest,
            })
            const session = await Bun.file(sessionFile).json()
            await Bun.write(dest, JSON.stringify(session))
            log.info(`migrating messages for session ${session.id}`)
            for await (const msgFile of new Bun.Glob(`storage/session/message/${session.id}/*.json`).scan({
              cwd: fullProjectDir,
              absolute: true,
            })) {
              const dest = path.join(dir, "message", session.id, path.basename(msgFile))
              log.info("copying", {
                msgFile,
                dest,
              })
              const message = await Bun.file(msgFile).json()
              await Bun.write(dest, JSON.stringify(message))

              log.info(`migrating parts for message ${message.id}`)
              for await (const partFile of new Bun.Glob(`storage/session/part/${session.id}/${message.id}/*.json`).scan(
                {
                  cwd: fullProjectDir,
                  absolute: true,
                },
              )) {
                const dest = path.join(dir, "part", message.id, path.basename(partFile))
                const part = await Bun.file(partFile).json()
                log.info("copying", {
                  partFile,
                  dest,
                })
                await Bun.write(dest, JSON.stringify(part))
              }
            }
          }
        }
      }
    },
  ]

  const state = lazy(async () => {
    const dir = path.join(Global.Path.data, "storage")
    const migration = await Bun.file(path.join(dir, "migration"))
      .json()
      .then((x) => parseInt(x))
      .catch(() => 0)
    for (let index = migration; index < MIGRATIONS.length; index++) {
      log.info("running migration", { index })
      const migration = MIGRATIONS[index]
      await migration(dir).catch((e) => {
        log.error("failed to run migration", { error: e, index })
      })
      await Bun.write(path.join(dir, "migration"), (index + 1).toString())
    }
    return {
      dir,
    }
  })

  // BUG FIX: Validate storage keys
  function validateKey(key: string[]): void {
    if (key.length === 0) {
      throw new Error("Storage key cannot be empty")
    }

    for (const segment of key) {
      if (!segment || typeof segment !== "string") {
        throw new Error(`Invalid key segment: ${segment}`)
      }
      // Prevent directory traversal attacks
      if (segment.includes("..") || segment.includes("/") || segment.includes("\\")) {
        throw new Error(`Invalid characters in key segment: ${segment}`)
      }
      // Prevent hidden files
      if (segment.startsWith(".")) {
        throw new Error(`Key segments cannot start with dot: ${segment}`)
      }
    }
  }

  export async function remove(key: string[]) {
    validateKey(key)
    const dir = await state().then((x) => x.dir)
    const target = path.join(dir, ...key) + ".json"
    await fs.unlink(target).catch((error) => {
      log.debug("File delete failed (may not exist)", { target, error: error.message })
    })
  }

  export async function read<T>(key: string[]) {
    validateKey(key)
    const dir = await state().then((x) => x.dir)
    const target = path.join(dir, ...key) + ".json"
    // Use file-specific read lock
    using _ = await Lock.read(target)
    return Bun.file(target).json() as Promise<T>
  }

  export async function update<T>(key: string[], fn: (draft: T) => void) {
    validateKey(key)
    const dir = await state().then((x) => x.dir)
    const target = path.join(dir, ...key) + ".json"
    // Ensure parent directory exists
    await fs.mkdir(path.dirname(target), { recursive: true })
    // FIXED: Use file-specific write lock instead of global "storage" lock
    using _ = await Lock.write(target)
    const content = await Bun.file(target).json()
    fn(content)
    await Bun.write(target, JSON.stringify(content, null, 2))
    return content as T
  }

  // BUG FIX: Validate content size (max 10MB)
  const MAX_CONTENT_SIZE = 10 * 1024 * 1024 // 10MB

  export async function write<T>(key: string[], content: T) {
    validateKey(key)

    // Validate content size
    const json = JSON.stringify(content, null, 2)
    if (json.length > MAX_CONTENT_SIZE) {
      throw new Error(
        `Content too large: ${json.length} bytes (max ${MAX_CONTENT_SIZE} bytes)`,
      )
    }

    const dir = await state().then((x) => x.dir)
    const target = path.join(dir, ...key) + ".json"
    // Ensure parent directory exists
    await fs.mkdir(path.dirname(target), { recursive: true })
    // FIXED: Use file-specific write lock instead of global "storage" lock
    using _ = await Lock.write(target)
    await Bun.write(target, json)
  }

  const glob = new Bun.Glob("**/*")
  export async function list(prefix: string[]) {
    const dir = await state().then((x) => x.dir)
    try {
      const result = await Array.fromAsync(
        glob.scan({
          cwd: path.join(dir, ...prefix),
          onlyFiles: true,
        }),
      ).then((results) => results.map((x) => [...prefix, ...x.slice(0, -5).split(path.sep)]))
      result.sort()
      return result
    } catch {
      return []
    }
  }

  /**
   * Transaction support for atomic multi-file operations
   */
  export class Transaction {
    private operations: Array<{ type: "write" | "remove"; key: string[]; content?: any }> = []
    private locks: Disposable[] = []
    private committed = false

    async write<T>(key: string[], content: T) {
      if (this.committed) throw new Error("Transaction already committed")
      this.operations.push({ type: "write", key, content })
    }

    async remove(key: string[]) {
      if (this.committed) throw new Error("Transaction already committed")
      this.operations.push({ type: "remove", key })
    }

    async commit() {
      if (this.committed) throw new Error("Transaction already committed")
      this.committed = true

      const dir = await state().then((x) => x.dir)

      try {
        // Acquire all locks first (sorted to prevent deadlocks)
        const lockKeys = this.operations
          .map((op) => path.join(dir, ...op.key) + ".json")
          .sort()
          .filter((key, index, arr) => arr.indexOf(key) === index) // Unique

        for (const key of lockKeys) {
          const lock = await Lock.write(key)
          this.locks.push(lock)
        }

        // Execute all operations
        for (const op of this.operations) {
          const target = path.join(dir, ...op.key) + ".json"
          if (op.type === "write") {
            // Ensure parent directory exists
            await fs.mkdir(path.dirname(target), { recursive: true })
            await Bun.write(target, JSON.stringify(op.content, null, 2))
          } else if (op.type === "remove") {
            await fs.unlink(target).catch((error) => {
              // File may not exist, log for debugging
              log.debug("File delete failed (may not exist)", { target, error: error.message })
            })
          }
        }
      } finally {
        // Release all locks
        for (const lock of this.locks) {
          lock[Symbol.dispose]()
        }
      }
    }

    async rollback() {
      if (this.committed) throw new Error("Transaction already committed or rolled back")
      this.committed = true

      // Release locks without executing operations
      for (const lock of this.locks) {
        lock[Symbol.dispose]()
      }
      this.operations = []
    }
  }

  export function transaction() {
    return new Transaction()
  }
}
