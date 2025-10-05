# OpenCode → RyCode Implementation Plan

> **AI-Native Renaming Project**: Comprehensive implementation roadmap for renaming "opencode" to "rycode" across the entire codebase.

**Project Duration:** 2 weeks (10 working days)
**Affected Files:** 367+
**Risk Level:** High (breaking changes to public API)
**Team:** Multi-agent collaboration (Claude, Codex, Gemini)

---

## 📋 Executive Summary

### Objectives
1. Rename all "opencode" references to "rycode" across the codebase
2. Update package names from `@opencode-ai/*` to `@rycode-ai/*`
3. Maintain backwards compatibility for 6 months
4. Ensure zero-downtime migration for users
5. Update all documentation and external references

### Success Criteria
- ✅ All builds pass (TypeScript + Go)
- ✅ All tests pass (unit + integration + E2E)
- ✅ No broken imports or module references
- ✅ Packages published to npm/Go module registry
- ✅ Documentation completely updated
- ✅ Backwards compatibility verified

---

## 🗓️ 10-Day Sprint Plan

### **Day 1-2: Preparation & Setup**

#### Day 1 Morning: Environment Setup
```bash
# Create feature branch
git checkout -b refactor/rename-opencode-to-rycode

# Backup current state
git tag backup-pre-rename-$(date +%Y%m%d)

# Reserve npm packages
npm info @rycode-ai/plugin  # Check availability
npm info @rycode-ai/sdk     # Check availability

# Verify Go module availability
# Check github.com/aaronmrosenthal/rycode
```

**Tasks:**
- [ ] Create feature branch
- [ ] Tag backup point
- [ ] Reserve npm package names (`@rycode-ai/plugin`, `@rycode-ai/sdk`)
- [ ] Update GitHub repo settings (if renaming repo)
- [ ] Set up test environment

#### Day 1 Afternoon: Dependency Analysis
```bash
# Find all package.json files
find . -name "package.json" -not -path "*/node_modules/*"

# Find all go.mod files
find . -name "go.mod"

# Find all imports
grep -r "from '@opencode-ai" --include="*.ts" --include="*.tsx"
grep -r "github.com/sst/opencode" --include="*.go"
```

**Tasks:**
- [ ] Map all package dependencies
- [ ] Identify circular dependencies
- [ ] Document import patterns
- [ ] Create dependency graph

#### Day 2: Create Migration Scripts

**Script 1: Rename Package Names** (`scripts/rename-packages.ts`)
```typescript
#!/usr/bin/env bun

import { readdir, readFile, writeFile } from 'fs/promises';
import { join } from 'path';

async function updatePackageJson(filePath: string) {
  const content = await readFile(filePath, 'utf-8');
  const pkg = JSON.parse(content);

  // Update name
  if (pkg.name?.includes('opencode')) {
    pkg.name = pkg.name.replace('opencode', 'rycode');
  }

  // Update dependencies
  for (const depType of ['dependencies', 'devDependencies', 'peerDependencies']) {
    if (pkg[depType]) {
      const deps = pkg[depType];
      for (const [name, version] of Object.entries(deps)) {
        if (name.includes('@opencode-ai')) {
          const newName = name.replace('@opencode-ai', '@rycode-ai');
          delete deps[name];
          deps[newName] = version;
        }
      }
    }
  }

  // Update bin
  if (pkg.bin?.opencode) {
    pkg.bin.rycode = pkg.bin.opencode;
    delete pkg.bin.opencode;
  }

  await writeFile(filePath, JSON.stringify(pkg, null, 2) + '\n');
  console.log(`✅ Updated: ${filePath}`);
}

async function findAndUpdateAllPackageJsons() {
  // Implementation to recursively find and update all package.json files
  const packages = [
    'package.json',
    'packages/opencode/package.json',
    'packages/plugin/package.json',
    'packages/sdk/js/package.json',
    'packages/web/package.json',
  ];

  for (const pkg of packages) {
    await updatePackageJson(pkg);
  }
}

await findAndUpdateAllPackageJsons();
```

**Script 2: Rename Imports** (`scripts/rename-imports.ts`)
```typescript
#!/usr/bin/env bun

import { readFile, writeFile } from 'fs/promises';
import { glob } from 'glob';

async function updateImports(filePath: string) {
  let content = await readFile(filePath, 'utf-8');
  let modified = false;

  // Update package imports
  if (content.includes('@opencode-ai')) {
    content = content.replace(/@opencode-ai\//g, '@rycode-ai/');
    modified = true;
  }

  // Update relative imports to renamed directories
  if (content.includes('packages/opencode/')) {
    content = content.replace(/packages\/opencode\//g, 'packages/rycode/');
    modified = true;
  }

  if (modified) {
    await writeFile(filePath, content);
    console.log(`✅ Updated imports: ${filePath}`);
  }
}

// Find all TypeScript/JavaScript files
const files = await glob('**/*.{ts,tsx,js,jsx,mjs}', {
  ignore: ['node_modules/**', 'dist/**', 'build/**', '.git/**']
});

for (const file of files) {
  await updateImports(file);
}
```

**Script 3: Rename Go Modules** (`scripts/rename-go-modules.sh`)
```bash
#!/bin/bash

# Update go.mod files
find . -name "go.mod" -not -path "*/node_modules/*" | while read -r file; do
  echo "Updating: $file"

  # Update module path
  sed -i.bak 's|github.com/sst/opencode|github.com/aaronmrosenthal/rycode|g' "$file"

  # Clean up backup
  rm "${file}.bak"
done

# Update Go import statements
find . -name "*.go" | while read -r file; do
  echo "Updating imports: $file"

  sed -i.bak 's|github.com/sst/opencode|github.com/aaronmrosenthal/rycode|g' "$file"

  rm "${file}.bak"
done

echo "✅ Go modules updated"
```

**Script 4: Verification** (`scripts/verify-rename.ts`)
```typescript
#!/usr/bin/env bun

import { glob } from 'glob';
import { readFile } from 'fs/promises';

interface Issue {
  file: string;
  line: number;
  match: string;
}

const issues: Issue[] = [];

// Find remaining "opencode" references
const files = await glob('**/*', {
  ignore: [
    'node_modules/**',
    'dist/**',
    '.git/**',
    'scripts/rename-*.ts',  // Ignore migration scripts
    'OPENCODE_RENAME_*.md', // Ignore planning docs
  ]
});

for (const file of files) {
  try {
    const content = await readFile(file, 'utf-8');
    const lines = content.split('\n');

    lines.forEach((line, idx) => {
      // Look for "opencode" but ignore comments about the migration
      if (
        line.toLowerCase().includes('opencode') &&
        !line.includes('// Migration:') &&
        !line.includes('// TODO: rename') &&
        !line.includes('formerly opencode')
      ) {
        issues.push({
          file,
          line: idx + 1,
          match: line.trim()
        });
      }
    });
  } catch (e) {
    // Skip binary files
  }
}

if (issues.length > 0) {
  console.log(`❌ Found ${issues.length} remaining "opencode" references:\n`);
  issues.forEach(({ file, line, match }) => {
    console.log(`${file}:${line}`);
    console.log(`  ${match}\n`);
  });
  process.exit(1);
} else {
  console.log('✅ No remaining "opencode" references found');
  process.exit(0);
}
```

---

### **Day 3-4: Core Package Renaming**

#### Phase 1: Rename Directories
```bash
# Rename main package directory
git mv packages/opencode packages/rycode

# Update bin directory
cd packages/rycode
git mv bin/opencode bin/rycode

# Update test fixtures if needed
find . -path "*/opencode/*" -type d
```

**Tasks:**
- [ ] Rename `packages/opencode` → `packages/rycode`
- [ ] Rename binary `bin/opencode` → `bin/rycode`
- [ ] Update `.gitignore` references
- [ ] Update `tsconfig.json` path mappings

#### Phase 2: Update Package Manifests
```bash
# Run the package rename script
bun run scripts/rename-packages.ts

# Verify changes
git diff packages/*/package.json
```

**Files to update:**
- [ ] `packages/rycode/package.json`
- [ ] `packages/plugin/package.json`
- [ ] `packages/sdk/js/package.json`
- [ ] `packages/web/package.json`
- [ ] Root `package.json`

#### Phase 3: Update Imports
```bash
# Run import update script
bun run scripts/rename-imports.ts

# Manual verification of critical files
code packages/rycode/src/index.ts
code packages/plugin/src/index.ts
code packages/sdk/js/src/client.ts
```

**Critical files requiring manual review:**
- [ ] `packages/rycode/src/index.ts`
- [ ] `packages/plugin/src/index.ts`
- [ ] `packages/sdk/js/src/client.ts`
- [ ] `packages/sdk/js/src/server.ts`

#### Phase 4: Update Go Modules
```bash
# Run Go module rename
bash scripts/rename-go-modules.sh

# Update go.mod in TUI package
cd packages/tui
go mod edit -module github.com/aaronmrosenthal/rycode/packages/tui

# Tidy dependencies
go mod tidy

# Verify build
go build ./...
```

**Tasks:**
- [ ] Update `packages/tui/go.mod`
- [ ] Update all Go import statements
- [ ] Run `go mod tidy`
- [ ] Verify `go build` succeeds

---

### **Day 5-6: Configuration & Theme Updates**

#### Config Files
```bash
# Rename config file
git mv opencode.json rycode.json

# Update config schema reference
# Edit rycode.json
```

**Update `rycode.json`:**
```json
{
  "$schema": "https://rycode.ai/config.json"
}
```

**Create backwards compatibility wrapper** (`packages/rycode/src/config/legacy.ts`):
```typescript
import { existsSync } from 'fs';
import { join } from 'path';

/**
 * Load config with backwards compatibility for opencode.json
 */
export async function loadConfig(cwd: string) {
  const rycodeConfig = join(cwd, 'rycode.json');
  const opencodeConfig = join(cwd, 'opencode.json');

  if (existsSync(rycodeConfig)) {
    return await import(rycodeConfig);
  }

  if (existsSync(opencodeConfig)) {
    console.warn(
      '⚠️  Deprecation Warning: opencode.json is deprecated. ' +
      'Please rename to rycode.json. Support will be removed in v1.0.0'
    );
    return await import(opencodeConfig);
  }

  return null;
}
```

#### Theme Files
```bash
# Rename theme file
git mv packages/tui/internal/theme/themes/opencode.json \
       packages/tui/internal/theme/themes/rycode.json

# Update theme references
grep -r "opencode.json" packages/tui/
```

**Update theme loader** (`packages/tui/internal/theme/loader.go`):
```go
// Update default theme reference
const DefaultTheme = "rycode"  // was "opencode"

// Add backwards compatibility
func LoadTheme(name string) (*Theme, error) {
  // Support legacy "opencode" theme name
  if name == "opencode" {
    log.Warn("Theme 'opencode' is deprecated, use 'rycode' instead")
    name = "rycode"
  }

  // ... rest of loading logic
}
```

#### Environment Variables
```bash
# Update env var references
grep -r "OPENCODE_" packages/
```

**Update all env var references:**
- `OPENCODE_API_KEY` → `RYCODE_API_KEY`
- `OPENCODE_CONFIG` → `RYCODE_CONFIG`
- `OPENCODE_HOME` → `RYCODE_HOME`

**Add backwards compatibility** (`packages/rycode/src/config/env.ts`):
```typescript
export function getApiKey(): string | undefined {
  // Try new name first
  let key = process.env.RYCODE_API_KEY;

  // Fallback to legacy name with warning
  if (!key && process.env.OPENCODE_API_KEY) {
    console.warn(
      '⚠️  OPENCODE_API_KEY is deprecated. Use RYCODE_API_KEY instead.'
    );
    key = process.env.OPENCODE_API_KEY;
  }

  return key;
}
```

---

### **Day 7: Documentation Updates**

#### Update All Markdown Files
```bash
# Find all docs with "opencode"
grep -r "opencode" --include="*.md" --include="*.mdx" .

# Update with sed or manual editing
find . -name "*.md" -o -name "*.mdx" | while read file; do
  sed -i.bak 's/opencode/rycode/g' "$file"
  sed -i.bak 's/OpenCode/RyCode/g' "$file"
  sed -i.bak 's/@opencode-ai/@rycode-ai/g' "$file"
  rm "${file}.bak"
done
```

**Critical documentation files:**
- [ ] `README.md` (root)
- [ ] `packages/rycode/README.md`
- [ ] `SECURITY.md`
- [ ] `SECURITY_GUIDE.md`
- [ ] `PLUGIN_SECURITY.md`
- [ ] All `packages/web/src/content/docs/*.mdx`

#### Update Code Comments
```bash
# Find TODO/FIXME comments mentioning opencode
grep -r "TODO.*opencode" --include="*.ts" --include="*.go"
grep -r "FIXME.*opencode" --include="*.ts" --include="*.go"
```

#### Update Examples
```bash
# Update example code
ls packages/rycode/examples/
```

**Example files to update:**
- [ ] `packages/rycode/examples/README.md`
- [ ] All demo files in `packages/rycode/examples/demos/`

---

### **Day 8: Testing & Validation**

#### Unit Tests
```bash
# Update test files
find . -name "*.test.ts" -o -name "*_test.go" | while read file; do
  # Update imports and references
  sed -i.bak 's/@opencode-ai/@rycode-ai/g' "$file"
  rm "${file}.bak"
done

# Run TypeScript tests
cd packages/rycode
bun test

# Run Go tests
cd packages/tui
go test ./...
```

**Test suites to verify:**
- [ ] `packages/rycode/test/` (all tests)
- [ ] `packages/tui/internal/` (Go tests)
- [ ] `packages/sdk/js/` (SDK tests)
- [ ] `packages/plugin/` (Plugin tests)

#### Integration Tests
```typescript
// Create integration test: scripts/test-integration.ts
#!/usr/bin/env bun

import { spawn } from 'child_process';
import { tmpdir } from 'os';
import { join } from 'path';
import { mkdir, writeFile } from 'fs/promises';

async function testInstallation() {
  console.log('🧪 Testing installation...');

  const testDir = join(tmpdir(), 'rycode-test-' + Date.now());
  await mkdir(testDir, { recursive: true });

  // Create test package.json
  await writeFile(
    join(testDir, 'package.json'),
    JSON.stringify({ dependencies: { '@rycode-ai/sdk': 'latest' } })
  );

  // Test npm install
  const npm = spawn('npm', ['install'], { cwd: testDir });

  return new Promise((resolve, reject) => {
    npm.on('close', (code) => {
      if (code === 0) {
        console.log('✅ Installation test passed');
        resolve(true);
      } else {
        console.error('❌ Installation test failed');
        reject(new Error('Install failed'));
      }
    });
  });
}

async function testCLI() {
  console.log('🧪 Testing CLI...');

  const cli = spawn('./packages/rycode/bin/rycode', ['--version']);

  return new Promise((resolve, reject) => {
    cli.stdout.on('data', (data) => {
      console.log(`Version: ${data}`);
    });

    cli.on('close', (code) => {
      if (code === 0) {
        console.log('✅ CLI test passed');
        resolve(true);
      } else {
        reject(new Error('CLI test failed'));
      }
    });
  });
}

async function testImports() {
  console.log('🧪 Testing imports...');

  try {
    // Test new package names
    await import('@rycode-ai/sdk');
    await import('@rycode-ai/plugin');
    console.log('✅ Import test passed');
  } catch (error) {
    console.error('❌ Import test failed:', error);
    throw error;
  }
}

// Run all tests
await testInstallation();
await testCLI();
await testImports();

console.log('\n✅ All integration tests passed!');
```

#### Build Verification
```bash
# Build all packages
bun run build

# Verify outputs
ls -la packages/rycode/dist/
ls -la packages/sdk/js/dist/

# Build Go TUI
cd packages/tui
go build -o dist/rycode ./cmd/rycode

# Test binary
./dist/rycode --help
```

**Build verification checklist:**
- [ ] TypeScript builds without errors
- [ ] Go builds without errors
- [ ] All dist directories populated correctly
- [ ] Binary is executable
- [ ] No broken imports in build output

#### Run Verification Script
```bash
# Check for remaining "opencode" references
bun run scripts/verify-rename.ts
```

---

### **Day 9: Publishing & Deployment**

#### Pre-Publish Checklist
- [ ] All tests passing
- [ ] All builds successful
- [ ] Documentation complete
- [ ] Changelog updated
- [ ] Version bumped (0.14.1 → 0.15.0)

#### Publish npm Packages
```bash
# Login to npm
npm login

# Publish packages in dependency order
cd packages/plugin
npm publish --access public

cd ../sdk/js
npm publish --access public

cd ../rycode
npm publish --access public
```

#### Publish Go Module
```bash
# Tag release
git tag packages/tui/v0.15.0
git push origin packages/tui/v0.15.0

# Go will automatically index from GitHub
```

#### Update Package Registries
```bash
# Homebrew tap (if applicable)
# Update formula in homebrew-tap repo

# Chocolatey (Windows)
# Update package manifest

# AUR (Arch Linux)
# Update PKGBUILD
```

---

### **Day 10: Verification & Cleanup**

#### Post-Deploy Verification
```bash
# Test fresh install
npx @rycode-ai/sdk@latest

# Test CLI install
npm install -g rycode@latest
rycode --version

# Test backwards compatibility
OPENCODE_API_KEY=test rycode --help  # Should show deprecation warning
```

#### Update External Services
- [ ] Update GitHub repo description
- [ ] Update package registry metadata
- [ ] Update website (if applicable)
- [ ] Update social media references

#### Cleanup Old Packages
```bash
# Deprecate old packages (after 6 months)
npm deprecate @opencode-ai/sdk "Package has been renamed to @rycode-ai/sdk"
npm deprecate @opencode-ai/plugin "Package has been renamed to @rycode-ai/plugin"
```

#### Create Migration Guide
Create `MIGRATION_GUIDE.md`:
```markdown
# Migration Guide: OpenCode → RyCode

## For Users

### Update your dependencies:
\`\`\`bash
npm uninstall @opencode-ai/sdk
npm install @rycode-ai/sdk
\`\`\`

### Update imports:
\`\`\`typescript
// Old
import { Client } from '@opencode-ai/sdk';

// New
import { Client } from '@rycode-ai/sdk';
\`\`\`

### Rename config file:
\`\`\`bash
mv opencode.json rycode.json
\`\`\`

### Update environment variables:
\`\`\`bash
# Old
export OPENCODE_API_KEY=xxx

# New
export RYCODE_API_KEY=xxx
\`\`\`

## Backwards Compatibility

For 6 months (until v1.0.0), the following will continue to work:
- `opencode.json` config file (with deprecation warning)
- `OPENCODE_*` environment variables (with deprecation warning)
- Old package names will redirect to new ones

## Timeline

- **v0.15.0** (Now): Rename complete, backwards compatibility active
- **v0.20.0** (3 months): Deprecation warnings become more prominent
- **v1.0.0** (6 months): Backwards compatibility removed
\`\`\`
```

---

## 🔍 Risk Assessment & Mitigation

### P0 - Critical Risks

#### Risk: Package installation breaks for existing users
**Impact:** Users cannot install or update
**Probability:** High
**Mitigation:**
- Maintain old package names with deprecation notice
- Publish both old and new packages for 6 months
- Clear migration guide
- Automated migration script

**Rollback:**
```bash
# If critical issues found
npm unpublish @rycode-ai/sdk --force
git revert <commit-sha>
git push origin main
```

#### Risk: Build failures in CI/CD
**Impact:** Cannot deploy
**Probability:** Medium
**Mitigation:**
- Update all CI/CD configs before merge
- Test in staging environment
- Run builds locally first

**Files to update:**
- [ ] `.github/workflows/test.yml`
- [ ] `.github/workflows/release.yml`
- [ ] Any deployment scripts

### P1 - High Risks

#### Risk: Broken imports in production code
**Impact:** Runtime errors for users
**Probability:** Medium
**Mitigation:**
- Comprehensive test coverage
- Integration tests
- Canary deployment

**Testing:**
```typescript
// Test all public API entry points
describe('Public API', () => {
  it('exports Client from @rycode-ai/sdk', async () => {
    const { Client } = await import('@rycode-ai/sdk');
    expect(Client).toBeDefined();
  });

  it('exports Plugin from @rycode-ai/plugin', async () => {
    const { Plugin } = await import('@rycode-ai/plugin');
    expect(Plugin).toBeDefined();
  });
});
```

#### Risk: Documentation links broken
**Impact:** Poor user experience
**Probability:** High
**Mitigation:**
- Automated link checking
- Redirect old URLs
- Update all references

**Verification:**
```bash
# Check all markdown links
npx markdown-link-check README.md
npx markdown-link-check packages/*/README.md
```

### P2 - Medium Risks

#### Risk: Theme compatibility issues
**Impact:** Visual inconsistencies
**Probability:** Low
**Mitigation:**
- Keep old theme file as alias
- Update theme loader with fallback

#### Risk: Plugin ecosystem breaks
**Impact:** Third-party plugins fail
**Probability:** Medium
**Mitigation:**
- Maintain `@opencode-ai/plugin` as alias
- Notify plugin authors
- Provide migration guide for plugin developers

### P3 - Low Risks

#### Risk: SEO impact from URL changes
**Impact:** Temporary search ranking drop
**Probability:** Low
**Mitigation:**
- 301 redirects from old URLs
- Update sitemap
- Notify search engines

---

## 🧪 Comprehensive Testing Strategy

### Unit Tests
```bash
# All package tests must pass
bun test                           # TypeScript
go test ./...                      # Go
```

### Integration Tests
```bash
# Test package installation
npm install @rycode-ai/sdk
npm install @rycode-ai/plugin

# Test CLI
rycode --version
rycode --help

# Test with real project
mkdir test-project
cd test-project
npm init -y
npm install @rycode-ai/sdk
```

### E2E Tests
```typescript
// e2e/rename-migration.test.ts
describe('Migration E2E', () => {
  it('should work with new package names', async () => {
    // Test complete workflow with new names
  });

  it('should maintain backwards compatibility', async () => {
    // Test old config files still work
    // Test old env vars still work
  });

  it('should show deprecation warnings', async () => {
    // Verify warnings are shown for old patterns
  });
});
```

### Backwards Compatibility Tests
```typescript
// test/backwards-compat.test.ts
describe('Backwards Compatibility', () => {
  it('loads opencode.json config', async () => {
    const config = await loadConfig(testDir);
    expect(config).toBeDefined();
  });

  it('reads OPENCODE_API_KEY env var', () => {
    process.env.OPENCODE_API_KEY = 'test-key';
    const key = getApiKey();
    expect(key).toBe('test-key');
  });

  it('supports opencode theme name', async () => {
    const theme = await loadTheme('opencode');
    expect(theme).toBeDefined();
  });
});
```

### Platform-Specific Tests
```bash
# macOS
bun test:macos

# Linux
bun test:linux

# Windows
bun test:windows
```

### Manual Testing Checklist
- [ ] Fresh install on clean machine
- [ ] Upgrade from v0.14.1 to v0.15.0
- [ ] CLI commands all work
- [ ] Config file migration
- [ ] Environment variables
- [ ] Theme loading
- [ ] Plugin system
- [ ] Documentation links
- [ ] Code examples in docs

---

## 📊 Success Metrics

### Quantitative Metrics
- **Build Success Rate:** 100% across all platforms
- **Test Pass Rate:** 100% (unit + integration + E2E)
- **Installation Success:** >95% successful installs
- **Migration Rate:** >80% of users migrated within 3 months
- **Support Tickets:** <10% increase related to rename

### Qualitative Metrics
- User feedback sentiment (GitHub issues, discussions)
- Developer experience improvements
- Documentation clarity
- Brand recognition

### Monitoring
```typescript
// Add telemetry for migration tracking
export function trackConfigLoad(source: 'rycode.json' | 'opencode.json') {
  analytics.track('config_loaded', { source });
}

export function trackEnvVarUsage(name: string) {
  analytics.track('env_var_used', { name });
}
```

**Dashboard metrics to monitor:**
- Config file usage (rycode.json vs opencode.json)
- Environment variable usage (RYCODE_* vs OPENCODE_*)
- Package download stats (old vs new names)
- Error rates by package version
- Migration guide page views

---

## 🔄 Rollback Procedure

### If Critical Issues Detected

#### Immediate Actions (< 1 hour)
```bash
# 1. Unpublish broken packages
npm unpublish @rycode-ai/sdk@0.15.0
npm unpublish @rycode-ai/plugin@0.15.0

# 2. Revert Git changes
git revert HEAD~10..HEAD  # Revert last 10 commits
git push origin main --force-with-lease

# 3. Re-publish previous version
git checkout v0.14.1
cd packages/opencode
npm publish
```

#### Communication (< 2 hours)
- Post incident notice on GitHub
- Update package README with known issues
- Notify users via social media
- Email notification to enterprise users

#### Root Cause Analysis (< 24 hours)
- Identify what went wrong
- Document lessons learned
- Update testing strategy
- Create prevention measures

### Partial Rollback
If only specific components are affected:
```bash
# Roll back just the problematic package
npm unpublish @rycode-ai/sdk@0.15.0
git revert <specific-commit>
```

---

## 📝 Changelog Entry

```markdown
## [0.15.0] - 2025-10-05

### 🎉 Major Changes

#### Package Rename: OpenCode → RyCode
- **Breaking:** All packages renamed from `@opencode-ai/*` to `@rycode-ai/*`
- **Breaking:** Binary renamed from `opencode` to `rycode`
- **Breaking:** Config file renamed from `opencode.json` to `rycode.json`
- **Breaking:** Environment variables renamed from `OPENCODE_*` to `RYCODE_*`

### ✨ New Features
- Full backwards compatibility for 6 months
- Automatic config file migration
- Deprecation warnings for old patterns
- Comprehensive migration guide

### 🔧 Improvements
- Updated all documentation
- Improved error messages
- Better CLI help text
- Enhanced theme system

### 📦 Migration Guide
See [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md) for detailed migration instructions.

### ⚠️ Deprecation Notices
- `@opencode-ai/*` packages deprecated (use `@rycode-ai/*`)
- `opencode.json` config deprecated (use `rycode.json`)
- `OPENCODE_*` env vars deprecated (use `RYCODE_*`)

All deprecated features will be removed in v1.0.0 (approximately 6 months).
```

---

## 🤖 Multi-Agent Task Distribution

### Parallel Execution Opportunities

#### Phase 1: Preparation (Parallel)
- **Agent 1 (Claude):** Analyze dependencies, create migration scripts
- **Agent 2 (Codex):** Update package.json files, handle npm logistics
- **Agent 3 (Gemini):** Review Go module structure, plan Go changes

#### Phase 2: Core Renaming (Sequential, then Parallel)
1. **Sequential:** Rename directories (one agent, avoid conflicts)
2. **Parallel:**
   - **Agent 1:** Update TypeScript imports
   - **Agent 2:** Update Go imports
   - **Agent 3:** Update configuration files

#### Phase 3: Documentation (Parallel)
- **Agent 1:** Update main docs (README, guides)
- **Agent 2:** Update code comments and examples
- **Agent 3:** Update test files and specs

#### Phase 4: Testing (Parallel)
- **Agent 1:** Run and verify unit tests
- **Agent 2:** Run and verify integration tests
- **Agent 3:** Run and verify E2E tests

### Sequential Dependencies
Must happen in order:
1. Rename directories → Update imports
2. Update package.json → Run install
3. All code changes → Run tests
4. Tests passing → Publish packages

### Critical Decision Points (Human Review Required)
- [ ] Final package name approval
- [ ] Pre-publish verification
- [ ] Release notes approval
- [ ] Deprecation timeline confirmation
- [ ] Rollback decision (if issues found)

---

## ✅ Final Checklist

### Pre-Implementation
- [ ] Backup created (Git tag)
- [ ] Feature branch created
- [ ] Team aligned on approach
- [ ] npm packages reserved
- [ ] Migration scripts ready

### During Implementation
- [ ] All directories renamed
- [ ] All imports updated
- [ ] All configs updated
- [ ] All docs updated
- [ ] All tests updated
- [ ] Backwards compatibility added

### Pre-Deployment
- [ ] All builds passing
- [ ] All tests passing (100%)
- [ ] Verification script clean
- [ ] Integration tests passing
- [ ] Documentation complete
- [ ] Changelog updated
- [ ] Version bumped
- [ ] Migration guide ready

### Post-Deployment
- [ ] Packages published
- [ ] Git tags created
- [ ] Release notes published
- [ ] Migration guide published
- [ ] Old packages deprecated (message only)
- [ ] Monitoring enabled
- [ ] Success metrics tracked

### After 6 Months
- [ ] Remove backwards compatibility
- [ ] Unpublish old packages
- [ ] Remove deprecation warnings
- [ ] Update documentation
- [ ] Close migration phase

---

## 📞 Communication Plan

### Internal Team
- Daily standup during implementation
- Slack updates on progress
- Blocker escalation process

### Users
- **Pre-Release:** Blog post announcing rename
- **Release:** Detailed release notes with migration guide
- **Post-Release:** Tutorial video on migration
- **Ongoing:** FAQ updates based on support tickets

### Community
- GitHub Discussions announcement
- Social media posts
- Newsletter update
- Documentation updates

---

## 🎯 Next Steps

1. **Review this plan** with team and stakeholders
2. **Get approval** for package names and timeline
3. **Reserve npm packages** `@rycode-ai/*`
4. **Set sprint start date**
5. **Begin Day 1 tasks**

---

**Plan Status:** Ready for Review
**Last Updated:** 2025-10-05
**Owner:** Multi-Agent Team (Claude, Codex, Gemini)
**Reviewers:** Human oversight required at critical decision points
