import z from "zod/v4"
import { Bus } from "../bus"
import { $ } from "bun"
import { formatPatch, structuredPatch } from "diff"
import path from "path"
import fs from "fs"
import ignore from "ignore"
import { Log } from "../util/log"
import { Instance } from "../project/instance"
import { Ripgrep } from "./ripgrep"
import fuzzysort from "fuzzysort"

export namespace File {
  const log = Log.create({ service: "file" })

  export const Info = z
    .object({
      path: z.string(),
      added: z.number().int(),
      removed: z.number().int(),
      status: z.enum(["added", "deleted", "modified"]),
    })
    .meta({
      ref: "File",
    })

  export type Info = z.infer<typeof Info>

  export const Node = z
    .object({
      name: z.string(),
      path: z.string(),
      absolute: z.string(),
      type: z.enum(["file", "directory"]),
      ignored: z.boolean(),
    })
    .meta({
      ref: "FileNode",
    })
  export type Node = z.infer<typeof Node>

  export const Content = z
    .object({
      content: z.string(),
      diff: z.string().optional(),
      patch: z
        .object({
          oldFileName: z.string(),
          newFileName: z.string(),
          oldHeader: z.string().optional(),
          newHeader: z.string().optional(),
          hunks: z.array(
            z.object({
              oldStart: z.number(),
              oldLines: z.number(),
              newStart: z.number(),
              newLines: z.number(),
              lines: z.array(z.string()),
            }),
          ),
          index: z.string().optional(),
        })
        .optional(),
    })
    .meta({
      ref: "FileContent",
    })
  export type Content = z.infer<typeof Content>

  export const Event = {
    Edited: Bus.event(
      "file.edited",
      z.object({
        file: z.string(),
      }),
    ),
  }

  const state = Instance.state(async () => {
    type Entry = { files: string[]; dirs: string[] }
    let cache: Entry = { files: [], dirs: [] }
    let fetching = false
    const fn = async (result: Entry) => {
      fetching = true
      const set = new Set<string>()
      for await (const file of Ripgrep.files({ cwd: Instance.directory })) {
        result.files.push(file)
        let current = file
        while (true) {
          const dir = path.dirname(current)
          if (dir === ".") break
          if (dir === current) break
          current = dir
          if (set.has(dir)) continue
          set.add(dir)
          result.dirs.push(dir + "/")
        }
      }
      cache = result
      fetching = false
    }
    fn(cache)

    return {
      async files() {
        if (!fetching) {
          fn({
            files: [],
            dirs: [],
          })
        }
        return cache
      },
    }
  })

  /**
   * Initializes the file system cache.
   *
   * Begins background indexing of files in the project directory.
   * Called automatically during system startup.
   */
  export function init(): void {
    state()
  }

  /**
   * Retrieves git status information for the project.
   *
   * Returns information about modified, added, and deleted files.
   * Only works for git-based projects.
   *
   * @returns Array of file changes with line counts and status
   */
  export async function status(): Promise<Info[]> {
    const project = Instance.project
    if (project.vcs !== "git") return []

    const diffOutput = await $`git diff --numstat HEAD`.cwd(Instance.directory).quiet().nothrow().text()

    const changedFiles: Info[] = []

    if (diffOutput.trim()) {
      const lines = diffOutput.trim().split("\n")
      for (const line of lines) {
        const [added, removed, filepath] = line.split("\t")
        changedFiles.push({
          path: filepath,
          added: added === "-" ? 0 : parseInt(added, 10),
          removed: removed === "-" ? 0 : parseInt(removed, 10),
          status: "modified",
        })
      }
    }

    const untrackedOutput = await $`git ls-files --others --exclude-standard`
      .cwd(Instance.directory)
      .quiet()
      .nothrow()
      .text()

    if (untrackedOutput.trim()) {
      const untrackedFiles = untrackedOutput.trim().split("\n")
      for (const filepath of untrackedFiles) {
        try {
          const content = await Bun.file(path.join(Instance.directory, filepath)).text()
          const lines = content.split("\n").length
          changedFiles.push({
            path: filepath,
            added: lines,
            removed: 0,
            status: "added",
          })
        } catch {
          continue
        }
      }
    }

    // Get deleted files
    const deletedOutput = await $`git diff --name-only --diff-filter=D HEAD`
      .cwd(Instance.directory)
      .quiet()
      .nothrow()
      .text()

    if (deletedOutput.trim()) {
      const deletedFiles = deletedOutput.trim().split("\n")
      for (const filepath of deletedFiles) {
        changedFiles.push({
          path: filepath,
          added: 0,
          removed: 0, // Could get original line count but would require another git command
          status: "deleted",
        })
      }
    }

    return changedFiles.map((x) => ({
      ...x,
      path: path.relative(Instance.directory, x.path),
    }))
  }

  /**
   * Reads a file with optional git diff information.
   *
   * For modified files in git repositories, includes structured patch data
   * showing the differences from the committed version.
   *
   * @param file - Relative path to file from project root
   * @returns Object with file content and optional diff/patch data
   */
  export async function read(file: string): Promise<Content> {
    using _ = log.time("read", { file })
    const project = Instance.project
    const full = path.join(Instance.directory, file)
    const content = await Bun.file(full)
      .text()
      .catch(() => "")
      .then((x) => x.trim())
    if (project.vcs === "git") {
      let diff = await $`git diff ${file}`.cwd(Instance.directory).quiet().nothrow().text()
      if (!diff.trim()) diff = await $`git diff --staged ${file}`.cwd(Instance.directory).quiet().nothrow().text()
      if (diff.trim()) {
        const original = await $`git show HEAD:${file}`.cwd(Instance.directory).quiet().nothrow().text()
        const patch = structuredPatch(file, file, original, content, "old", "new", {
          context: Infinity,
          ignoreWhitespace: true,
        })
        const diff = formatPatch(patch)
        return { content, patch, diff }
      }
    }
    return { content }
  }

  /**
   * Lists files and directories at a given path.
   *
   * Respects .gitignore rules for git projects. Excludes .git and .DS_Store.
   * Results are sorted with directories first, then files alphabetically.
   *
   * @param dir - Optional directory path relative to project root (defaults to root)
   * @returns Array of file/directory nodes with metadata
   */
  export async function list(dir?: string): Promise<Node[]> {
    const exclude = [".git", ".DS_Store"]
    const project = Instance.project
    let ignored = (_: string) => false
    if (project.vcs === "git") {
      const gitignore = Bun.file(path.join(Instance.worktree, ".gitignore"))
      if (await gitignore.exists()) {
        const ig = ignore().add(await gitignore.text())
        ignored = ig.ignores.bind(ig)
      }
    }
    const resolved = dir ? path.join(Instance.directory, dir) : Instance.directory
    const nodes: Node[] = []
    for (const entry of await fs.promises.readdir(resolved, { withFileTypes: true })) {
      if (exclude.includes(entry.name)) continue
      const fullPath = path.join(resolved, entry.name)
      const relativePath = path.relative(Instance.directory, fullPath)
      const type = entry.isDirectory() ? "directory" : "file"
      nodes.push({
        name: entry.name,
        path: relativePath,
        absolute: fullPath,
        type,
        ignored: ignored(type === "directory" ? relativePath + "/" : relativePath),
      })
    }
    return nodes.sort((a, b) => {
      if (a.type !== b.type) {
        return a.type === "directory" ? -1 : 1
      }
      return a.name.localeCompare(b.name)
    })
  }

  /**
   * Fuzzy searches for files and directories by name.
   *
   * Uses fuzzysort for intelligent matching. Results are ranked by relevance.
   *
   * @param input - Search parameters (query string and optional result limit)
   * @returns Array of matching file/directory paths, limited to specified count
   */
  export async function search(input: { query: string; limit?: number }): Promise<string[]> {
    log.info("search", { query: input.query })
    const limit = input.limit ?? 100
    const result = await state().then((x) => x.files())
    if (!input.query) return result.dirs.toSorted().slice(0, limit)
    const items = [...result.files, ...result.dirs]
    const sorted = fuzzysort.go(input.query, items, { limit: limit }).map((r) => r.target)
    log.info("search", { query: input.query, results: sorted.length })
    return sorted
  }
}
