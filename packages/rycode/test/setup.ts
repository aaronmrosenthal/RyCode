/**
 * Test setup and utilities for OpenCode test suite
 */

import { afterAll, beforeAll } from "bun:test"
import path from "path"
import fs from "fs/promises"
import os from "os"

export namespace TestSetup {
  let tempDir: string

  /**
   * Creates a temporary directory for test files
   */
  export async function createTempDir(): Promise<string> {
    tempDir = await fs.mkdtemp(path.join(os.tmpdir(), "opencode-test-"))
    return tempDir
  }

  /**
   * Cleanup temporary directory after tests
   */
  export async function cleanup() {
    if (tempDir) {
      await fs.rm(tempDir, { recursive: true, force: true }).catch(() => {})
    }
  }

  /**
   * Create a test file in the temp directory
   */
  export async function createTestFile(relativePath: string, content: string): Promise<string> {
    if (!tempDir) {
      tempDir = await createTempDir()
    }
    const filePath = path.join(tempDir, relativePath)
    const dir = path.dirname(filePath)
    await fs.mkdir(dir, { recursive: true })
    await fs.writeFile(filePath, content, "utf-8")
    return filePath
  }

  /**
   * Mock environment variables for tests
   */
  export function mockEnv(vars: Record<string, string>): () => void {
    const original: Record<string, string | undefined> = {}
    for (const [key, value] of Object.entries(vars)) {
      original[key] = process.env[key]
      process.env[key] = value
    }
    return () => {
      for (const [key, value] of Object.entries(original)) {
        if (value === undefined) {
          delete process.env[key]
        } else {
          process.env[key] = value
        }
      }
    }
  }

  /**
   * Helper to create a mock Request object
   */
  export function createMockRequest(options: {
    method?: string
    path?: string
    headers?: Record<string, string>
    query?: Record<string, string>
    body?: any
  }): Request {
    const { method = "GET", path = "/", headers = {}, query = {}, body } = options

    const url = new URL(path, "http://localhost:3000")
    for (const [key, value] of Object.entries(query)) {
      url.searchParams.set(key, value)
    }

    return new Request(url.toString(), {
      method,
      headers: new Headers(headers),
      body: body ? JSON.stringify(body) : undefined,
    })
  }
}

/**
 * Global setup
 */
beforeAll(async () => {
  // Setup test environment
  process.env.NODE_ENV = "test"
})

/**
 * Global teardown
 */
afterAll(async () => {
  await TestSetup.cleanup()
})
