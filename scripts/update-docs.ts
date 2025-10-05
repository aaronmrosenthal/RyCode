#!/usr/bin/env bun

import { readFile, writeFile } from "fs/promises"
import { glob } from "glob"

async function updateDocFile(filePath: string) {
  let content = await readFile(filePath, "utf-8")
  let modified = false
  const changes: string[] = []

  // Update GitHub URLs
  if (content.includes("github.com/sst/opencode")) {
    content = content.replace(/github\.com\/sst\/opencode/g, "github.com/aaronmrosenthal/rycode")
    modified = true
    changes.push("GitHub URLs")
  }

  // Update opencode.ai references (but preserve in historical context)
  const lines = content.split("\n")
  const updatedLines = lines.map((line) => {
    // Skip lines that are clearly historical or migration-related
    if (
      line.includes("formerly") ||
      line.includes("previously") ||
      line.includes("was:") ||
      line.includes("renamed from") ||
      line.toLowerCase().includes("migration")
    ) {
      return line
    }

    // Update opencode.ai to rycode.ai
    if (line.includes("opencode.ai")) {
      modified = true
      changes.push("Domain URLs")
      return line.replace(/opencode\.ai/g, "rycode.ai")
    }

    return line
  })
  content = updatedLines.join("\n")

  // Update "OpenCode" brand name to "RyCode" (preserve in code examples and historical references)
  const safeReplacements = [
    // Product name in titles and headers
    [/^# OpenCode\b/gm, "# RyCode"],
    [/^## OpenCode\b/gm, "## RyCode"],
    [/^### OpenCode\b/gm, "### RyCode"],

    // In descriptive text (but not in code or quotes)
    [/\bOpenCode takes\b/g, "RyCode takes"],
    [/\bOpenCode supports\b/g, "RyCode supports"],
    [/\bOpenCode includes\b/g, "RyCode includes"],
    [/\bOpenCode validates\b/g, "RyCode validates"],
    [/\busing OpenCode\b/g, "using RyCode"],
    [/\bKeep OpenCode\b/g, "Keep RyCode"],
    [/\bRun OpenCode\b/g, "Run RyCode"],

    // Package and command references
    [/`opencode`/g, "`rycode`"],
    [/opencode\\.json/g, "rycode.json"],
    [/opencode\\.jsonc/g, "rycode.jsonc"],

    // CLI commands (outside of git clone examples)
    [/\$ opencode\b/g, "$ rycode"],
  ]

  for (const [pattern, replacement] of safeReplacements) {
    if (pattern.test(content)) {
      content = content.replace(pattern, replacement)
      modified = true
      if (!changes.includes("Brand names")) changes.push("Brand names")
    }
  }

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
  console.log("🔍 Finding documentation files to update...\n")

  const files = [
    "README.md",
    "SECURITY.md",
    "AGENTS.md",
    "SECURITY_ASSESSMENT.md",
    "WEAKNESSES_ANALYSIS.md",
    "UX_IMPROVEMENT_SUMMARY.md",
    "BUG_FIXES_SUMMARY.md",
    "FINAL_BUG_FIXES.md",
    "COMPLETE_IMPLEMENTATION_SUMMARY.md",
    "CONCURRENCY_IMPROVEMENTS.md",
  ]

  console.log(`Found ${files.length} documentation files\n`)

  let updatedCount = 0
  for (const file of files) {
    try {
      const updated = await updateDocFile(file)
      if (updated) updatedCount++
    } catch (e) {
      console.log(`⏭️  Skipped (not found): ${file}`)
    }
  }

  console.log(`\n📊 Summary:`)
  console.log(`   Total files: ${files.length}`)
  console.log(`   Updated: ${updatedCount}`)
  console.log(`   Skipped: ${files.length - updatedCount}`)
}

await main()
