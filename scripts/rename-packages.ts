#!/usr/bin/env bun

import { readdir, readFile, writeFile } from "fs/promises"
import { join } from "path"
import { glob } from "glob"

interface PackageJson {
  name?: string
  bin?: Record<string, string> | string
  dependencies?: Record<string, string>
  devDependencies?: Record<string, string>
  peerDependencies?: Record<string, string>
  [key: string]: any
}

async function updatePackageJson(filePath: string) {
  const content = await readFile(filePath, "utf-8")
  const pkg: PackageJson = JSON.parse(content)
  let modified = false

  // Update name
  if (pkg.name?.includes("opencode")) {
    pkg.name = pkg.name.replace(/opencode/g, "rycode").replace(/@opencode-ai/g, "@rycode-ai")
    modified = true
    console.log(`  ✓ Updated name: ${pkg.name}`)
  }

  // Update bin
  if (pkg.bin) {
    if (typeof pkg.bin === "string") {
      if (pkg.bin.includes("opencode")) {
        pkg.bin = pkg.bin.replace(/opencode/g, "rycode")
        modified = true
        console.log(`  ✓ Updated bin: ${pkg.bin}`)
      }
    } else {
      const newBin: Record<string, string> = {}
      for (const [name, path] of Object.entries(pkg.bin)) {
        const newName = name.replace(/opencode/g, "rycode")
        const newPath = path.replace(/opencode/g, "rycode")
        newBin[newName] = newPath
        if (newName !== name || newPath !== path) {
          modified = true
          console.log(`  ✓ Updated bin: ${newName} -> ${newPath}`)
        }
      }
      pkg.bin = newBin
    }
  }

  // Update dependencies
  for (const depType of ["dependencies", "devDependencies", "peerDependencies"]) {
    if (pkg[depType]) {
      const deps = pkg[depType] as Record<string, string>
      const newDeps: Record<string, string> = {}
      for (const [name, version] of Object.entries(deps)) {
        if (name.includes("@opencode-ai") || name === "opencode") {
          const newName = name.replace(/@opencode-ai/g, "@rycode-ai").replace(/^opencode$/, "rycode")
          newDeps[newName] = version
          modified = true
          console.log(`  ✓ Updated ${depType}: ${name} -> ${newName}`)
        } else {
          newDeps[name] = version
        }
      }
      pkg[depType] = newDeps
    }
  }

  if (modified) {
    await writeFile(filePath, JSON.stringify(pkg, null, 2) + "\n")
    console.log(`✅ Updated: ${filePath}\n`)
    return true
  }

  console.log(`⏭️  Skipped (no changes): ${filePath}\n`)
  return false
}

async function findAndUpdateAllPackageJsons() {
  console.log("🔍 Finding all package.json files...\n")

  const files = await glob("**/package.json", {
    ignore: ["**/node_modules/**", "**/dist/**", "**/build/**"],
  })

  console.log(`Found ${files.length} package.json files\n`)

  let updatedCount = 0
  for (const file of files) {
    console.log(`Processing: ${file}`)
    const updated = await updatePackageJson(file)
    if (updated) updatedCount++
  }

  console.log(`\n📊 Summary:`)
  console.log(`   Total files: ${files.length}`)
  console.log(`   Updated: ${updatedCount}`)
  console.log(`   Skipped: ${files.length - updatedCount}`)
}

await findAndUpdateAllPackageJsons()
