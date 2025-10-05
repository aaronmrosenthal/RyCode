import { describe, test, expect, beforeAll } from "bun:test"
import { FileSecurity } from "../../src/file/security"
import { TestSetup } from "../setup"
import { Instance } from "../../src/project/instance"
import { Project } from "../../src/project/project"
import path from "path"

describe("FileSecurity", () => {
  let tempDir: string
  let worktree: string

  beforeAll(async () => {
    tempDir = await TestSetup.createTempDir()
    worktree = tempDir
  })

  test("allows paths within directory", async () => {
    await Instance.provide({
      directory: tempDir,
      async fn() {
        const validPath = path.join(tempDir, "file.txt")
        await TestSetup.createTestFile("file.txt", "content")

        const result = FileSecurity.validatePath("file.txt")
        expect(result).toBe(validPath)
      },
    })
  })

  test("blocks path traversal attempts", async () => {
    await Instance.provide({
      directory: tempDir,
      async fn() {
        expect(() => {
          FileSecurity.validatePath("../../etc/passwd")
        }).toThrow(FileSecurity.PathTraversalError)
      },
    })
  })

  test("blocks access to .env files", async () => {
    await Instance.provide({
      directory: tempDir,
      async fn() {
        expect(() => {
          FileSecurity.validatePath(".env")
        }).toThrow(FileSecurity.SensitiveFileError)
      },
    })
  })

  test("blocks access to .env.* files", async () => {
    await Instance.provide({
      directory: tempDir,
      async fn() {
        expect(() => {
          FileSecurity.validatePath(".env.local")
        }).toThrow(FileSecurity.SensitiveFileError)

        expect(() => {
          FileSecurity.validatePath(".env.production")
        }).toThrow(FileSecurity.SensitiveFileError)
      },
    })
  })

  test("blocks access to SSH keys", async () => {
    await Instance.provide({
      directory: tempDir,
      async fn() {
        expect(() => {
          FileSecurity.validatePath(".ssh/id_rsa")
        }).toThrow(FileSecurity.SensitiveFileError)

        expect(() => {
          FileSecurity.validatePath(".ssh/id_ed25519")
        }).toThrow(FileSecurity.SensitiveFileError)
      },
    })
  })

  test("blocks access to credential files", async () => {
    await Instance.provide({
      directory: tempDir,
      async fn() {
        expect(() => {
          FileSecurity.validatePath(".aws/credentials")
        }).toThrow(FileSecurity.SensitiveFileError)

        expect(() => {
          FileSecurity.validatePath("credentials.json")
        }).toThrow(FileSecurity.SensitiveFileError)
      },
    })
  })

  test("blocks access to PEM and key files", async () => {
    await Instance.provide({
      directory: tempDir,
      async fn() {
        expect(() => {
          FileSecurity.validatePath("private.pem")
        }).toThrow(FileSecurity.SensitiveFileError)

        expect(() => {
          FileSecurity.validatePath("server.key")
        }).toThrow(FileSecurity.SensitiveFileError)
      },
    })
  })

  test("normalizes paths correctly", async () => {
    await Instance.provide({
      directory: tempDir,
      async fn() {
        await TestSetup.createTestFile("subdir/file.txt", "content")

        const result = FileSecurity.validatePath("./subdir/../subdir/file.txt")
        expect(result).toContain("subdir/file.txt")
      },
    })
  })

  test("validates multiple paths in batch", async () => {
    await Instance.provide({
      directory: tempDir,
      async fn() {
        await TestSetup.createTestFile("file1.txt", "content1")
        await TestSetup.createTestFile("file2.txt", "content2")

        const results = FileSecurity.validatePaths(["file1.txt", "file2.txt"])
        expect(results).toHaveLength(2)
        expect(results[0]).toContain("file1.txt")
        expect(results[1]).toContain("file2.txt")
      },
    })
  })

  test("isPathSafe returns boolean without throwing", async () => {
    await Instance.provide({
      directory: tempDir,
      async fn() {
        expect(FileSecurity.isPathSafe("file.txt")).toBe(true)
        expect(FileSecurity.isPathSafe(".env")).toBe(false)
        expect(FileSecurity.isPathSafe("../../etc/passwd")).toBe(false)
      },
    })
  })

  test("blocks access to system files on Unix", async () => {
    await Instance.provide({
      directory: tempDir,
      async fn() {
        expect(() => {
          FileSecurity.validatePath("/etc/passwd")
        }).toThrow(FileSecurity.PathTraversalError)

        expect(() => {
          FileSecurity.validatePath("/etc/shadow")
        }).toThrow(FileSecurity.PathTraversalError)
      },
    })
  })

  test("allows safe file paths", async () => {
    await Instance.provide({
      directory: tempDir,
      async fn() {
        await TestSetup.createTestFile("src/index.ts", "code")
        await TestSetup.createTestFile("package.json", "{}")
        await TestSetup.createTestFile("README.md", "docs")

        expect(() => {
          FileSecurity.validatePath("src/index.ts")
        }).not.toThrow()

        expect(() => {
          FileSecurity.validatePath("package.json")
        }).not.toThrow()

        expect(() => {
          FileSecurity.validatePath("README.md")
        }).not.toThrow()
      },
    })
  })
})
