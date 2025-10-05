#!/usr/bin/env bun

import { readFile, writeFile } from "fs/promises"
import { glob } from "glob"

async function updateImports(filePath: string) {
  let content = await readFile(filePath, "utf-8")
  let modified = false
  const changes: string[] = []

  // Update package imports (@opencode-ai -> @rycode-ai)
  if (content.includes("@opencode-ai")) {
    content = content.replace(/@opencode-ai\//g, "@rycode-ai/")
    modified = true
    changes.push("@opencode-ai -> @rycode-ai")
  }

  // Update relative imports to renamed directories
  if (content.includes("packages/opencode/")) {
    content = content.replace(/packages\/opencode\//g, "packages/rycode/")
    modified = true
    changes.push("packages/opencode -> packages/rycode")
  }

  // Update import from 'opencode' (if used as package name)
  const importPattern = /from ['"]opencode['"]/g
  if (importPattern.test(content)) {
    content = content.replace(importPattern, (match) => match.replace("opencode", "rycode"))
    modified = true
    changes.push("from 'opencode' -> from 'rycode'")
  }

  // Update require('opencode') or require("opencode")
  const requirePattern = /require\(['"]opencode['"]\)/g
  if (requirePattern.test(content)) {
    content = content.replace(requirePattern, (match) => match.replace("opencode", "rycode"))
    modified = true
    changes.push("require('opencode') -> require('rycode')")
  }

  // Update scriptName in CLI files
  if (content.includes('.scriptName("opencode")')) {
    content = content.replace(/\.scriptName\("opencode"\)/g, '.scriptName("rycode")')
    modified = true
    changes.push('scriptName("opencode") -> scriptName("rycode")')
  }

  // Update environment variable references in code comments and strings
  // BUT be careful not to change actual env var names yet (that's separate)
  // We'll just update the documentation strings for now

  if (modified) {
    await writeFile(filePath, content)
    console.log(`✅ Updated: ${filePath}`)
    changes.forEach((change) => console.log(`   - ${change}`))
    console.log()
    return true
  }

  return false
}

async function main() {
  console.log("🔍 Finding all TypeScript/JavaScript files...\n")

  const files = await glob("**/*.{ts,tsx,js,jsx,mjs,cjs}", {
    ignore: [
      "**/node_modules/**",
      "**/dist/**",
      "**/build/**",
      "**/.git/**",
      "scripts/rename-*.ts", // Don't update the migration scripts themselves
      "OPENCODE_RENAME_*.md", // Don't update planning docs
    ],
  })

  console.log(`Found ${files.length} files\n`)

  let updatedCount = 0
  for (const file of files) {
    const updated = await updateImports(file)
    if (updated) updatedCount++
  }

  console.log(`\n📊 Summary:`)
  console.log(`   Total files: ${files.length}`)
  console.log(`   Updated: ${updatedCount}`)
  console.log(`   Skipped: ${files.length - updatedCount}`)
}

await main()
