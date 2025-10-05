#!/usr/bin/env bun

import { glob } from "glob"
import { readFile } from "fs/promises"

interface Issue {
  file: string
  line: number
  match: string
}

const issues: Issue[] = []

// Files and patterns to ignore
const IGNORE_PATTERNS = [
  "**/node_modules/**",
  "**/dist/**",
  "**/build/**",
  "**/.git/**",
  "scripts/rename-*.ts",
  "scripts/rename-*.sh",
  "scripts/verify-*.ts",
  "OPENCODE_RENAME_*.md",
  "PROJECT_CONTEXT.md", // Contains historical references
  "IMPLEMENTATION_SUMMARY.md", // Contains historical references
  "SECURITY_MIGRATION.md", // May contain migration notes
]

// Patterns that are acceptable (backwards compatibility)
const ACCEPTABLE_PATTERNS = [
  /\/\/ Migration:/,
  /\/\/ TODO: rename/,
  /\/\/ Formerly opencode/,
  /\/\/ Was: opencode/,
  /formerly opencode/i,
  /previously opencode/i,
  /renamed from opencode/i,
  /OPENCODE_\w+ is deprecated/,
  /deprecated.*opencode/i,
  /backwards compatibility.*opencode/i,
  /legacy.*opencode/i,
]

function isAcceptableLine(line: string): boolean {
  return ACCEPTABLE_PATTERNS.some((pattern) => pattern.test(line))
}

async function checkFile(filePath: string) {
  try {
    const content = await readFile(filePath, "utf-8")
    const lines = content.split("\n")

    lines.forEach((line, idx) => {
      const lowerLine = line.toLowerCase()

      // Check for "opencode" references
      if (lowerLine.includes("opencode")) {
        // Skip if it's an acceptable pattern (migration notes, etc.)
        if (isAcceptableLine(line)) {
          return
        }

        // Skip if it's in a URL comment about the old site
        if (line.includes("https://opencode.ai") && line.trim().startsWith("//")) {
          return
        }

        issues.push({
          file: filePath,
          line: idx + 1,
          match: line.trim(),
        })
      }

      // Check for old package references
      if (lowerLine.includes("@opencode-ai")) {
        issues.push({
          file: filePath,
          line: idx + 1,
          match: line.trim(),
        })
      }

      // Check for old GitHub references
      if (lowerLine.includes("github.com/sst/opencode")) {
        issues.push({
          file: filePath,
          line: idx + 1,
          match: line.trim(),
        })
      }
    })
  } catch (e) {
    // Skip binary files and files that can't be read as text
  }
}

async function main() {
  console.log("🔍 Searching for remaining 'opencode' references...\n")

  const files = await glob("**/*", {
    ignore: IGNORE_PATTERNS,
    nodir: true,
  })

  console.log(`Checking ${files.length} files...\n`)

  for (const file of files) {
    await checkFile(file)
  }

  if (issues.length > 0) {
    console.log(`❌ Found ${issues.length} remaining "opencode" references:\n`)

    // Group by file
    const byFile = new Map<string, Issue[]>()
    for (const issue of issues) {
      if (!byFile.has(issue.file)) {
        byFile.set(issue.file, [])
      }
      byFile.get(issue.file)!.push(issue)
    }

    // Print grouped by file
    for (const [file, fileIssues] of byFile.entries()) {
      console.log(`📄 ${file}:`)
      for (const issue of fileIssues) {
        console.log(`   Line ${issue.line}: ${issue.match}`)
      }
      console.log()
    }

    console.log(`\n⚠️  Please review and fix these references manually.`)
    console.log(`   Some may be acceptable (e.g., in comments about migration).`)
    process.exit(1)
  } else {
    console.log("✅ No remaining 'opencode' references found!")
    console.log("   The rename appears to be complete.\n")
    process.exit(0)
  }
}

await main()
