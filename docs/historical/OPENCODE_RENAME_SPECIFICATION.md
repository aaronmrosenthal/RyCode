# OpenCode to RyCode Renaming Specification

## Executive Summary

This document provides a comprehensive specification for renaming "opencode" throughout the RyCode codebase to a more semantic name that better reflects its identity as an AI-native development platform. This is a critical refactoring that will affect package names, import paths, binary names, configuration files, URLs, and documentation across the entire monorepo.

**Current State**: The codebase uses "opencode" in 367+ files across TypeScript/JavaScript packages, Go packages, documentation, themes, and configuration files.

**Proposed Change**: Rename to **"rycode"** to establish a unique brand identity for the AI-native multi-agent development platform while maintaining backwards compatibility where critical.

---

## Name Options Analysis

### Option 1: **rycode** (RECOMMENDED)
**Semantic Meaning**: Already established as the project name, representing the "Ry" prefix (suggesting collaboration/flow) combined with "code"

**Rationale**:
- Already the monorepo root name and GitHub repository name
- Unique, memorable, and brandable
- Natural evolution from the "toolkit-cli enhanced OpenCode" origin story
- Maintains connection to "code" for developer familiarity
- No namespace conflicts with existing npm packages or GitHub repos

**Pros**:
- Establishes unique brand identity
- Already recognized in the README and repository structure
- Short, memorable, CLI-friendly
- Clean separation from upstream OpenCode project

**Cons**:
- Requires updating all references to establish new identity
- Need to ensure npm package availability

---

### Option 2: **aiflow**
**Semantic Meaning**: Emphasizes the flow-based multi-agent AI collaboration aspect

**Rationale**:
- Descriptive of the multi-agent workflow nature
- Terminal/CLI context appropriate
- Developer-friendly naming

**Pros**:
- Clearly communicates AI-native purpose
- Professional, enterprise-ready name
- Good for marketing and positioning

**Cons**:
- Generic, may conflict with existing tools
- Loses connection to the RyCode brand
- Less unique in the crowded AI dev tools space

---

### Option 3: **agentkit**
**Semantic Meaning**: Emphasizes the toolkit/SDK nature for AI agent collaboration

**Rationale**:
- Aligns with "toolkit-cli" heritage
- Clearly indicates purpose (agent development kit)
- Developer-centric naming

**Pros**:
- Self-documenting name
- Strong connection to multi-agent architecture
- Appeals to technical audience

**Cons**:
- Doesn't leverage existing RyCode brand
- "kit" suffix may feel generic
- Potential conflicts with other "kit" libraries

---

### Option 4: **matrixcode**
**Semantic Meaning**: Emphasizes the Matrix cyberpunk aesthetic and collaborative coding

**Rationale**:
- Aligns with Matrix theme in the codebase
- Unique and memorable
- Appeals to the cyberpunk/terminal aesthetic

**Pros**:
- Strong thematic connection
- Memorable brand identity
- Differentiates from competition

**Cons**:
- May be too niche/themed
- Potential trademark concerns with "Matrix"
- Doesn't clearly communicate AI-native purpose

---

### Option 5: **codeweave**
**Semantic Meaning**: Represents weaving together multiple AI agents to create code

**Rationale**:
- Poetic metaphor for multi-agent collaboration
- Unique in the development tools space
- Professional yet creative

**Pros**:
- Beautiful metaphor for collaboration
- Unique namespace
- Memorable

**Cons**:
- Less immediately clear purpose
- Doesn't leverage existing RyCode identity
- May feel abstract to some users

---

## Recommended Choice: **rycode**

### Justification

**rycode** is the recommended choice for the following strategic reasons:

1. **Brand Continuity**: Already established as the project name, repository name, and in documentation
2. **Unique Identity**: Clearly separates this project from upstream OpenCode while maintaining connection to code/development
3. **Marketing Story**: Natural evolution narrative: "OpenCode → toolkit-cli optimization → RyCode"
4. **Technical Clarity**: Short, memorable, CLI-friendly name
5. **Low Risk**: Already partially adopted, reducing transition complexity
6. **Community Recognition**: Users already associate the project with RyCode
7. **SEO & Discovery**: Unique enough for strong search presence

### Semantic Alignment

- **AI-native**: The "Ry" prefix suggests flow, collaboration, and modern development
- **Multi-agent**: Implies coordination and intelligent orchestration
- **Developer-focused**: "code" suffix maintains clear connection to development
- **Terminal/CLI**: Short, memorable command name

---

## Complete Renaming Map

### 1. Package Names

#### NPM Packages
| Current | New | Impact |
|---------|-----|--------|
| `opencode` | `rycode` | BREAKING - Main CLI package |
| `@opencode-ai/sdk` | `@rycode-ai/sdk` | BREAKING - SDK package |
| `@opencode-ai/plugin` | `@rycode-ai/plugin` | BREAKING - Plugin package |
| `@opencode-ai/console-app` | `@rycode-ai/console-app` | Internal |
| `@opencode-ai/console-core` | `@rycode-ai/console-core` | Internal |
| `@opencode-ai/console-function` | `@rycode-ai/console-function` | Internal |
| `@opencode-ai/console-mail` | `@rycode-ai/console-mail` | Internal |
| `@opencode-ai/console-resource` | `@rycode-ai/console-resource` | Internal |
| `@opencode-ai/console-scripts` | `@rycode-ai/console-scripts` | Internal |
| `@opencode-ai/desktop` | `@rycode-ai/desktop` | Internal |
| `@opencode-ai/function` | `@rycode-ai/function` | Internal |
| `@opencode-ai/web` | `@rycode-ai/web` | Internal |

#### Go Packages
| Current | New | Impact |
|---------|-----|--------|
| `github.com/sst/opencode` | `github.com/sst/rycode` | BREAKING - Main TUI module |
| `github.com/sst/opencode-sdk-go` | `github.com/sst/rycode-sdk-go` | BREAKING - Go SDK |
| `github.com/sst/opencode/internal/*` | `github.com/sst/rycode/internal/*` | Internal |

### 2. Directory Structure

| Current | New | Notes |
|---------|-----|-------|
| `/packages/opencode/` | `/packages/rycode/` | Main CLI package |
| `/packages/opencode/bin/opencode` | `/packages/rycode/bin/rycode` | Binary wrapper script |
| `/packages/opencode/bin/opencode.cmd` | `/packages/rycode/bin/rycode.cmd` | Windows wrapper |
| `/packages/opencode/dist/opencode-*` | `/packages/rycode/dist/rycode-*` | Distribution folders |
| `/packages/tui/cmd/opencode/` | `/packages/tui/cmd/rycode/` | TUI entry point |
| `/.opencode/` | `/.rycode/` | Hidden config directory |
| `/packages/console/app/.opencode/` | `/packages/console/app/.rycode/` | Console config |
| `/node_modules/@opencode-ai/` | `/node_modules/@rycode-ai/` | Installed packages |

### 3. Binary & Executable Names

| Current | New | Platforms |
|---------|-----|-----------|
| `opencode` | `rycode` | CLI command (all platforms) |
| `opencode.exe` | `rycode.exe` | Windows executable |
| `opencode-darwin-x64` | `rycode-darwin-x64` | macOS Intel binary |
| `opencode-darwin-arm64` | `rycode-darwin-arm64` | macOS Apple Silicon binary |
| `opencode-linux-x64` | `rycode-linux-x64` | Linux x64 binary |
| `opencode-linux-arm64` | `rycode-linux-arm64` | Linux ARM64 binary |
| `opencode-windows-x64` | `rycode-windows-x64` | Windows x64 binary |

### 4. Configuration Files

| Current | New | Location |
|---------|-----|----------|
| `opencode.json` | `rycode.json` | Project root config |
| `$HOME/.opencode/config.json` | `$HOME/.rycode/config.json` | User config (legacy support) |
| `$XDG_CONFIG_HOME/opencode/` | `$XDG_CONFIG_HOME/rycode/` | XDG config directory |

**Config Schema URLs**:
- `https://opencode.ai/config.json` → `https://rycode.ai/config.json`
- `https://opencode.ai/theme.json` → `https://rycode.ai/theme.json`

### 5. Environment Variables

| Current | New | Backwards Compatibility |
|---------|-----|-------------------------|
| `OPENCODE_SERVER` | `RYCODE_SERVER` | Check both, prefer new |
| `OPENCODE_BIN_PATH` | `RYCODE_BIN_PATH` | Check both, prefer new |
| `OPENCODE_INSTALL_DIR` | `RYCODE_INSTALL_DIR` | Check both, prefer new |
| `OPENCODE_API_KEY` | `RYCODE_API_KEY` | Check both, prefer new |
| `OPENCODE_LOG_LEVEL` | `RYCODE_LOG_LEVEL` | Check both, prefer new |
| `OPENCODE_CONFIG_PATH` | `RYCODE_CONFIG_PATH` | Check both, prefer new |
| `OPENCODE_THEME` | `RYCODE_THEME` | Check both, prefer new |

### 6. Import Statements

#### TypeScript/JavaScript
```typescript
// Before
import { createOpencodeClient, createOpencodeServer } from "@opencode-ai/sdk"
import "@opencode-ai/plugin"
import type { OpencodeConfig } from "@opencode-ai/sdk"

// After
import { createRycodeClient, createRycodeServer } from "@rycode-ai/sdk"
import "@rycode-ai/plugin"
import type { RycodeConfig } from "@rycode-ai/sdk"
```

#### Go
```go
// Before
import (
    "github.com/sst/opencode-sdk-go"
    "github.com/sst/opencode/internal/api"
)

// After
import (
    "github.com/sst/rycode-sdk-go"
    "github.com/sst/rycode/internal/api"
)
```

### 7. URLs & Domains

| Current | New | Type |
|---------|-----|------|
| `https://opencode.ai` | `https://rycode.ai` | Main website |
| `https://opencode.ai/docs` | `https://rycode.ai/docs` | Documentation |
| `https://opencode.ai/install` | `https://rycode.ai/install` | Install script |
| `https://opencode.ai/s/[id]` | `https://rycode.ai/s/[id]` | Share links |
| `https://opencode.ai/config.json` | `https://rycode.ai/config.json` | Config schema |
| `https://opencode.ai/theme.json` | `https://rycode.ai/theme.json` | Theme schema |

### 8. Theme Files

| Current | New | Location |
|---------|-----|----------|
| `themes/opencode.json` | `themes/rycode.json` | Default theme |
| Theme schema references | Update schema URLs | All theme files |

### 9. Documentation References

**File Types to Update**:
- All `.md` and `.mdx` files
- All `README.md` files
- `SECURITY.md` guides
- API documentation
- Code comments
- JSDoc/GoDoc comments
- Example code snippets
- Tutorial content

**Search Patterns**:
- Case-insensitive "opencode"
- "OpenCode" (capitalized)
- "OPENCODE" (all caps)
- "@opencode-ai"
- "opencode.ai"
- "sst/opencode"

### 10. Package Manager Configs

#### Homebrew
- Tap name: `sst/tap/opencode` → `sst/tap/rycode`
- Formula name: `opencode.rb` → `rycode.rb`

#### Chocolatey
- Package ID: `opencode` → `rycode`

#### WinGet
- Package ID: `opencode` → `rycode`

#### Scoop
- Manifest: `extras/opencode` → `extras/rycode`

#### AUR (Arch Linux)
- Package: `opencode-bin` → `rycode-bin`

#### NPM
- Package: `opencode-ai` → `rycode-ai` (consider)
- Or: `opencode` → `rycode`

---

## Migration Strategy

### Phase 1: Preparation (Week 1)
**Goal**: Set up infrastructure and backwards compatibility

1. **Reserve Package Names**
   - Register `@rycode-ai` org on npm
   - Reserve `rycode` package on npm
   - Check Go module availability

2. **Domain & Infrastructure**
   - Secure `rycode.ai` domain
   - Set up redirects from `opencode.ai` → `rycode.ai`
   - Configure DNS and SSL certificates

3. **Create Compatibility Layer**
   - Implement environment variable fallbacks
   - Add config file migration logic
   - Create symlink support for old binary names

4. **Update Package Managers**
   - Prepare Homebrew formula
   - Submit Chocolatey package
   - Update AUR package
   - Prepare WinGet manifest
   - Update Scoop manifest

### Phase 2: Core Renaming (Week 2)
**Goal**: Rename internal packages and modules

1. **Rename Go Packages**
   - Update `go.mod` module path
   - Update all Go import statements
   - Update TUI cmd directory
   - Run tests to verify

2. **Rename TypeScript/JavaScript Packages**
   - Rename `packages/opencode` → `packages/rycode`
   - Update `package.json` files
   - Update workspace references
   - Update import statements

3. **Rename SDK Packages**
   - `@opencode-ai/sdk` → `@rycode-ai/sdk`
   - `@opencode-ai/plugin` → `@rycode-ai/plugin`
   - Update generated types
   - Update client/server exports

4. **Update Binary Names**
   - Rename binary wrapper scripts
   - Update platform-specific binaries
   - Update build scripts
   - Update distribution configs

### Phase 3: Configuration & Theming (Week 2-3)
**Goal**: Update config files and themes

1. **Config Migration**
   - Add `rycode.json` config support
   - Implement migration from `opencode.json`
   - Update schema URLs
   - Add deprecation warnings for old config

2. **Theme Updates**
   - Rename `opencode.json` theme → `rycode.json`
   - Update theme schema references
   - Update theme loader logic
   - Test all themes

3. **Environment Variables**
   - Add support for `RYCODE_*` variables
   - Maintain fallback to `OPENCODE_*`
   - Add deprecation warnings
   - Update documentation

### Phase 4: Documentation (Week 3)
**Goal**: Update all documentation and external references

1. **Core Documentation**
   - Update README.md files
   - Update SECURITY.md guides
   - Update API documentation
   - Update example code

2. **Website & Guides**
   - Update `/packages/web/` content
   - Update all `.mdx` documentation
   - Update tutorial content
   - Update troubleshooting guides

3. **External References**
   - Update GitHub repository description
   - Update npm package descriptions
   - Update social media links
   - Update community forums

### Phase 5: Publishing (Week 4)
**Goal**: Release updated packages

1. **Publish New Packages**
   - Publish `@rycode-ai/*` to npm
   - Publish Go modules
   - Update Homebrew tap
   - Submit to package managers

2. **Deprecate Old Packages**
   - Mark `@opencode-ai/*` as deprecated
   - Add upgrade instructions
   - Set up automatic redirects
   - Monitor migration

3. **Update Install Scripts**
   - Update `install.sh` script
   - Test installation on all platforms
   - Update installation docs
   - Monitor for issues

### Phase 6: Cleanup (Week 5-6)
**Goal**: Remove deprecated code and finalize migration

1. **Remove Old References**
   - Clean up old package references
   - Remove deprecated configs
   - Clean up old environment variables
   - Update tests

2. **Final Verification**
   - Run full test suite
   - Verify all platforms
   - Check all package managers
   - Validate documentation

3. **Community Communication**
   - Announce migration completion
   - Provide migration guide
   - Answer community questions
   - Monitor for issues

---

## Breaking Changes

### Critical Breaking Changes

1. **Package Names** (BREAKING)
   - `opencode` → `rycode` CLI command
   - `@opencode-ai/*` → `@rycode-ai/*` npm packages
   - `github.com/sst/opencode*` → `github.com/sst/rycode*` Go modules

2. **Binary Names** (BREAKING with compatibility)
   - CLI command changes from `opencode` to `rycode`
   - **Mitigation**: Provide `opencode` symlink with deprecation warning

3. **Import Paths** (BREAKING)
   - All TypeScript/JavaScript imports change
   - All Go imports change
   - **Mitigation**: Publish one final version of old packages with upgrade instructions

4. **Configuration Files** (BREAKING with migration)
   - `opencode.json` → `rycode.json`
   - **Mitigation**: Auto-migrate on first run, support both temporarily

5. **Environment Variables** (BREAKING with fallback)
   - `OPENCODE_*` → `RYCODE_*`
   - **Mitigation**: Check both, prefer new, warn on old

### Non-Breaking Changes (with deprecation)

1. **URLs** (Redirects)
   - `opencode.ai` → `rycode.ai`
   - **Mitigation**: Set up permanent redirects (301)

2. **Theme Names** (Aliased)
   - `opencode` theme → `rycode` theme
   - **Mitigation**: Support both names, prefer new

3. **Config Schemas** (Versioned)
   - Update schema URLs
   - **Mitigation**: Maintain old schema URLs with redirects

---

## Backwards Compatibility Strategy

### Compatibility Layer

Implement a compatibility layer that:

1. **Environment Variables**
   ```typescript
   function getEnvVar(name: string): string | undefined {
     // Try new name first
     const newValue = process.env[`RYCODE_${name}`]
     if (newValue) return newValue

     // Fall back to old name with warning
     const oldValue = process.env[`OPENCODE_${name}`]
     if (oldValue) {
       console.warn(`DEPRECATION: OPENCODE_${name} is deprecated. Use RYCODE_${name} instead.`)
       return oldValue
     }

     return undefined
   }
   ```

2. **Config Files**
   ```typescript
   async function loadConfig(): Promise<Config> {
     // Try new config first
     const newConfig = await loadFile('rycode.json')
     if (newConfig) return newConfig

     // Fall back to old config with migration
     const oldConfig = await loadFile('opencode.json')
     if (oldConfig) {
       console.warn('DEPRECATION: opencode.json is deprecated. Migrating to rycode.json...')
       await migrateConfig(oldConfig, 'rycode.json')
       return oldConfig
     }

     return defaultConfig
   }
   ```

3. **Binary Symlinks**
   - Install `rycode` as primary binary
   - Create `opencode` symlink that warns about deprecation
   - Remove symlink after 6 months

4. **Package Aliases**
   - Publish final version of `@opencode-ai/*` packages
   - Add postinstall script with upgrade instructions
   - Point to new `@rycode-ai/*` packages

### Deprecation Timeline

- **Month 1-2**: Both names supported, warnings issued
- **Month 3-4**: Increase warning frequency, update docs
- **Month 5-6**: Prepare for removal, final migration push
- **Month 7+**: Remove `opencode` compatibility (major version bump)

---

## Testing Strategy

### 1. Unit Tests
- Update all test imports
- Update test fixtures
- Update mock configs
- Verify all tests pass

### 2. Integration Tests
- Test CLI installation on all platforms
- Test package manager installations
- Test config migration
- Test environment variable fallbacks

### 3. End-to-End Tests
- Install via npm, verify `rycode` command
- Install via Homebrew, verify binary
- Test config file migration
- Test theme loading with new names

### 4. Backwards Compatibility Tests
- Verify old environment variables work
- Verify old config files migrate
- Verify `opencode` symlink works (with warning)
- Test old import paths are properly deprecated

### 5. Platform-Specific Tests

| Platform | Test Scenarios |
|----------|---------------|
| macOS | Homebrew install, binary execution, config migration |
| Linux | Package manager install, permissions, XDG paths |
| Windows | Chocolatey/WinGet/Scoop install, .cmd wrapper |
| All | NPM install, Bun install, global vs local |

### 6. Documentation Tests
- Verify all links work
- Check schema URLs resolve
- Test install scripts
- Validate example code

---

## Rollout Plan

### Pre-Release (2 weeks before)

1. **Communication**
   - Blog post announcing the rename
   - GitHub discussions thread
   - Social media announcements
   - Email to users (if available)

2. **Infrastructure**
   - Set up rycode.ai domain
   - Configure redirects
   - Prepare package manager submissions
   - Reserve package names

3. **Beta Testing**
   - Release beta versions with new names
   - Gather feedback from early adopters
   - Fix critical issues
   - Refine migration scripts

### Release Day (v1.0.0)

1. **Package Publishing**
   - Publish `@rycode-ai/*` to npm
   - Publish Go modules to GitHub
   - Submit to Homebrew
   - Submit to other package managers

2. **Documentation**
   - Deploy updated website
   - Publish migration guide
   - Update GitHub README
   - Update all external links

3. **Deprecation**
   - Mark old packages as deprecated
   - Add upgrade instructions
   - Enable compatibility warnings
   - Monitor for issues

### Post-Release (First week)

1. **Monitoring**
   - Track installation metrics
   - Monitor error reports
   - Answer community questions
   - Fix critical bugs quickly

2. **Support**
   - Provide migration assistance
   - Update troubleshooting docs
   - Create FAQ for common issues
   - Engage with community

### Long-term (First 6 months)

1. **Gradual Deprecation**
   - Month 1-2: Soft warnings
   - Month 3-4: Stronger warnings
   - Month 5-6: Removal preparation
   - Month 7+: Remove compatibility (v2.0.0)

2. **Community Engagement**
   - Regular migration status updates
   - Celebrate migration milestones
   - Highlight success stories
   - Thank early adopters

---

## Risk Mitigation

### High-Risk Areas

1. **Package Installation Failures**
   - **Risk**: Users unable to install new packages
   - **Mitigation**: Thorough testing on all platforms, maintain old packages temporarily
   - **Rollback**: Keep old package versions available

2. **Breaking Existing Workflows**
   - **Risk**: Users' automation breaks
   - **Mitigation**: Provide `opencode` compatibility symlink, clear migration guide
   - **Rollback**: Document how to use old versions

3. **Lost SEO/Discovery**
   - **Risk**: Users can't find the project
   - **Mitigation**: Proper redirects, update all external links, SEO optimization
   - **Rollback**: Maintain content on both domains temporarily

4. **Community Confusion**
   - **Risk**: Users confused about the rename
   - **Mitigation**: Clear communication, comprehensive docs, active support
   - **Rollback**: Extend compatibility period

### Medium-Risk Areas

1. **Theme Compatibility**
   - **Risk**: Custom themes break
   - **Mitigation**: Support both theme names, auto-migration
   - **Rollback**: Keep old theme loader

2. **Plugin Ecosystem**
   - **Risk**: Third-party plugins break
   - **Mitigation**: Plugin compatibility layer, update docs
   - **Rollback**: Maintain old plugin API

3. **IDE Integration**
   - **Risk**: VS Code extension breaks
   - **Mitigation**: Update extension simultaneously, test thoroughly
   - **Rollback**: Keep extension backwards compatible

### Low-Risk Areas

1. **Documentation Links**
   - **Risk**: Some links 404
   - **Mitigation**: Comprehensive link audit, redirects
   - **Rollback**: Update links as discovered

2. **Internal References**
   - **Risk**: Missed renames in code
   - **Mitigation**: Automated search/replace, code review
   - **Rollback**: Quick patches for discoveries

---

## Success Metrics

### Quantitative Metrics

1. **Installation Success Rate**
   - Target: >95% successful installs across all platforms
   - Track: NPM downloads, Homebrew installs, error reports

2. **Migration Rate**
   - Target: >80% of active users migrated within 3 months
   - Track: Package version analytics, config file usage

3. **Breaking Issues**
   - Target: <10 critical issues in first week
   - Track: GitHub issues, error reporting, support tickets

4. **Documentation Coverage**
   - Target: 100% of docs updated
   - Track: Automated link checker, manual review

### Qualitative Metrics

1. **Community Sentiment**
   - Track: GitHub discussions, social media, user feedback
   - Goal: Positive reception of new brand

2. **Developer Experience**
   - Track: Migration ease, documentation clarity
   - Goal: Smooth transition with minimal friction

3. **Brand Recognition**
   - Track: SEO rankings, community mentions
   - Goal: Strong association with AI-native development

---

## Update Checklist

### Package Files
- [ ] All `package.json` files renamed
- [ ] All `go.mod` files updated
- [ ] Workspace references updated
- [ ] Lock files regenerated

### Source Code
- [ ] All TypeScript/JavaScript imports updated
- [ ] All Go imports updated
- [ ] All binary references updated
- [ ] All config file references updated

### Configuration
- [ ] Environment variable names updated
- [ ] Config file schema updated
- [ ] Theme files renamed
- [ ] Migration scripts created

### Documentation
- [ ] All README files updated
- [ ] All .mdx documentation updated
- [ ] API documentation updated
- [ ] Example code updated
- [ ] Tutorial content updated
- [ ] Troubleshooting guides updated

### Infrastructure
- [ ] Domain acquired and configured
- [ ] Redirects set up
- [ ] SSL certificates configured
- [ ] Package names reserved

### Distribution
- [ ] NPM packages published
- [ ] Go modules published
- [ ] Homebrew formula updated
- [ ] Chocolatey package updated
- [ ] WinGet manifest updated
- [ ] Scoop manifest updated
- [ ] AUR package updated

### Testing
- [ ] All unit tests passing
- [ ] All integration tests passing
- [ ] Platform-specific tests passing
- [ ] Backwards compatibility verified
- [ ] Migration scripts tested

### Communication
- [ ] Blog post published
- [ ] GitHub announcement posted
- [ ] Social media updated
- [ ] Migration guide published
- [ ] FAQ created

### Post-Launch
- [ ] Monitor installation metrics
- [ ] Track error reports
- [ ] Provide migration support
- [ ] Update based on feedback
- [ ] Plan deprecation timeline

---

## Appendix A: File Count Analysis

### Files Affected by Category

| Category | Count | Examples |
|----------|-------|----------|
| TypeScript/JavaScript | 150+ | src/**/*.ts, src/**/*.tsx |
| Go | 80+ | internal/**/*.go, cmd/**/*.go |
| Documentation | 60+ | *.md, *.mdx, docs/** |
| Configuration | 30+ | package.json, go.mod, *.json |
| Build/Deploy | 20+ | scripts/**, .github/workflows/** |
| Tests | 40+ | test/**, **/*.test.ts |
| Binaries | 10+ | bin/**, dist/** |
| Themes | 25+ | themes/**/*.json |

**Total**: 367+ files directly affected

### Search Patterns for Complete Coverage

```bash
# Case-insensitive search for all variants
rg -i "opencode" --type-add 'config:*.{json,yaml,yml,toml}' -t config
rg -i "opencode" --type ts --type js
rg -i "opencode" --type go
rg -i "opencode" --type md
rg -i "@opencode-ai"
rg "github\.com/sst/opencode"
rg "https://opencode\.ai"
rg "OPENCODE_"
```

---

## Appendix B: Migration Scripts

### Package.json Renamer
```bash
#!/bin/bash
# Rename all package.json references

find . -name "package.json" -type f | while read -r file; do
  sed -i.bak 's/"@opencode-ai\//"@rycode-ai\//g' "$file"
  sed -i.bak 's/"opencode"/"rycode"/g' "$file"
  sed -i.bak 's/opencode-ai/rycode-ai/g' "$file"
  rm "${file}.bak"
done
```

### Go Module Renamer
```bash
#!/bin/bash
# Update Go import paths

find . -name "*.go" -type f | while read -r file; do
  sed -i.bak 's|github.com/sst/opencode-sdk-go|github.com/sst/rycode-sdk-go|g' "$file"
  sed -i.bak 's|github.com/sst/opencode/|github.com/sst/rycode/|g' "$file"
  rm "${file}.bak"
done

# Update go.mod
sed -i.bak 's|module github.com/sst/opencode|module github.com/sst/rycode|g' go.mod
sed -i.bak 's|github.com/sst/opencode-sdk-go|github.com/sst/rycode-sdk-go|g' go.mod
```

### Config Migration Script
```typescript
// Auto-migrate opencode.json to rycode.json
import { readFile, writeFile, exists } from 'fs/promises'

async function migrateConfig() {
  const oldConfig = 'opencode.json'
  const newConfig = 'rycode.json'

  if (await exists(oldConfig) && !await exists(newConfig)) {
    const config = JSON.parse(await readFile(oldConfig, 'utf-8'))

    // Update schema URL if present
    if (config.$schema) {
      config.$schema = config.$schema.replace('opencode.ai', 'rycode.ai')
    }

    await writeFile(newConfig, JSON.stringify(config, null, 2))
    console.log('✓ Migrated opencode.json → rycode.json')
  }
}
```

---

## Appendix C: Communication Templates

### Blog Post Template
```markdown
# Announcing RyCode: OpenCode Reimagined

We're excited to announce that OpenCode is evolving into **RyCode** -
a more focused brand identity for our AI-native development platform.

## What's Changing?

- Command name: `opencode` → `rycode`
- Package names: `@opencode-ai/*` → `@rycode-ai/*`
- Domain: opencode.ai → rycode.ai

## Why the Change?

RyCode better represents our vision as a unique AI-native platform
built on multi-agent collaboration, while honoring our toolkit-cli heritage.

## How to Migrate

[Migration steps...]

## Questions?

[Support information...]
```

### Migration Guide Template
```markdown
# RyCode Migration Guide

## Quick Start

1. Update your installation:
   ```bash
   npm uninstall -g opencode-ai
   npm install -g rycode-ai
   ```

2. Update your imports:
   ```typescript
   // Before
   import { createOpencodeClient } from '@opencode-ai/sdk'

   // After
   import { createRycodeClient } from '@rycode-ai/sdk'
   ```

3. Rename your config:
   ```bash
   mv opencode.json rycode.json
   ```

## Backwards Compatibility

[Compatibility information...]
```

---

## Conclusion

This specification provides a comprehensive roadmap for renaming "opencode" to "rycode" throughout the codebase. The migration is significant but manageable with proper planning, phased execution, and robust backwards compatibility measures.

**Key Success Factors**:
1. Thorough testing across all platforms
2. Clear communication with the community
3. Robust backwards compatibility
4. Phased rollout with monitoring
5. Responsive support during transition

**Estimated Timeline**: 6-8 weeks for complete migration

**Risk Level**: Medium (mitigated by compatibility layer and phased approach)

**Recommendation**: Proceed with "rycode" as the new name, following the phased migration strategy outlined above.
