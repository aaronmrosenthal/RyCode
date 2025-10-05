import { cmd } from "./cmd"
import { PluginRegistry } from "../../plugin/registry"
import { PluginSecurity } from "../../plugin/security"
import { UI } from "../ui"
import path from "path"
import { existsSync } from "fs"

/**
 * Plugin Registry CLI Commands
 *
 * Provides commands to manage the plugin registry:
 * - plugin:registry:add     Add a plugin to the registry
 * - plugin:registry:remove  Remove a plugin from the registry
 * - plugin:registry:list    List all plugins in registry
 * - plugin:registry:search  Search for plugins
 * - plugin:registry:sync    Sync with remote registry
 * - plugin:registry:stats   Show registry statistics
 */

export const PluginRegistryAddCommand = cmd({
  command: "plugin:registry:add <plugin-path> <name> <version>",
  describe: "Add a plugin to the registry",
  builder: (yargs) =>
    yargs
      .positional("plugin-path", {
        describe: "Path to the plugin file",
        type: "string",
        demandOption: true,
      })
      .positional("name", {
        describe: "Plugin package name",
        type: "string",
        demandOption: true,
      })
      .positional("version", {
        describe: "Plugin version",
        type: "string",
        demandOption: true,
      })
      .option("description", {
        describe: "Plugin description",
        type: "string",
      })
      .option("author", {
        describe: "Plugin author",
        type: "string",
      })
      .option("homepage", {
        describe: "Homepage URL",
        type: "string",
      })
      .option("repository", {
        describe: "Repository URL",
        type: "string",
      })
      .option("verified-by", {
        describe: "Verification level",
        type: "string",
        choices: ["official", "community", "user"],
        default: "user",
      })
      .option("json", {
        describe: "Output as JSON",
        type: "boolean",
        default: false,
      }),
  handler: async (args) => {
    const pluginPath = path.resolve(args.pluginPath as string)
    const name = args.name as string
    const version = args.version as string

    if (!existsSync(pluginPath)) {
      UI.error(`Plugin file not found: ${pluginPath}`)
      process.exit(1)
    }

    try {
      // Generate hash
      const hash = await PluginSecurity.generateHash(pluginPath)

      // Add to registry
      await PluginRegistry.add({
        name,
        version,
        hash,
        description: args.description as string | undefined,
        author: args.author as string | undefined,
        homepage: args.homepage as string | undefined,
        repository: args.repository as string | undefined,
        verifiedBy: (args.verifiedBy as "official" | "community" | "user") || "user",
      })

      if (args.json) {
        console.log(JSON.stringify({
          success: true,
          plugin: name,
          version,
          hash,
        }, null, 2))
      } else {
        UI.println()
        UI.println(UI.Style.BOLD + "✓ Plugin Added to Registry" + UI.Style.RESET)
        UI.println()
        UI.println(UI.Style.TEXT_INFO + "Plugin:  " + UI.Style.RESET + name)
        UI.println(UI.Style.TEXT_INFO + "Version: " + UI.Style.RESET + version)
        UI.println(UI.Style.TEXT_INFO + "Hash:    " + UI.Style.RESET + UI.Style.TEXT_SUCCESS + hash + UI.Style.RESET)
        if (args.description) {
          UI.println(UI.Style.TEXT_INFO + "Description: " + UI.Style.RESET + args.description)
        }
        UI.println()
      }
    } catch (error) {
      UI.error(`Failed to add plugin: ${error instanceof Error ? error.message : String(error)}`)
      process.exit(1)
    }
  },
})

export const PluginRegistryRemoveCommand = cmd({
  command: "plugin:registry:remove <name> <version>",
  describe: "Remove a plugin from the registry",
  builder: (yargs) =>
    yargs
      .positional("name", {
        describe: "Plugin package name",
        type: "string",
        demandOption: true,
      })
      .positional("version", {
        describe: "Plugin version",
        type: "string",
        demandOption: true,
      })
      .option("json", {
        describe: "Output as JSON",
        type: "boolean",
        default: false,
      }),
  handler: async (args) => {
    const name = args.name as string
    const version = args.version as string

    try {
      const removed = await PluginRegistry.remove(name, version)

      if (args.json) {
        console.log(JSON.stringify({ success: removed, plugin: name, version }, null, 2))
      } else {
        if (removed) {
          UI.println()
          UI.println(UI.Style.TEXT_SUCCESS + `✓ Removed ${name}@${version} from registry` + UI.Style.RESET)
          UI.println()
        } else {
          UI.println()
          UI.println(UI.Style.TEXT_WARNING + `Plugin ${name}@${version} not found in registry` + UI.Style.RESET)
          UI.println()
        }
      }
    } catch (error) {
      UI.error(`Failed to remove plugin: ${error instanceof Error ? error.message : String(error)}`)
      process.exit(1)
    }
  },
})

export const PluginRegistryListCommand = cmd({
  command: "plugin:registry:list [name]",
  describe: "List plugins in the registry",
  builder: (yargs) =>
    yargs
      .positional("name", {
        describe: "Filter by plugin name",
        type: "string",
      })
      .option("json", {
        describe: "Output as JSON",
        type: "boolean",
        default: false,
      }),
  handler: async (args) => {
    try {
      const registry = await PluginRegistry.load()
      let entries = registry.entries

      if (args.name) {
        entries = entries.filter(e => e.name === args.name)
      }

      if (args.json) {
        console.log(JSON.stringify(entries, null, 2))
      } else {
        if (entries.length === 0) {
          UI.println()
          UI.println(UI.Style.TEXT_INFO + "No plugins found in registry" + UI.Style.RESET)
          UI.println()
          return
        }

        UI.println()
        UI.println(UI.Style.BOLD + "Plugin Registry" + UI.Style.RESET)
        UI.println(UI.Style.DIM + `${entries.length} ${entries.length === 1 ? 'entry' : 'entries'}` + UI.Style.RESET)
        UI.println()

        for (const entry of entries) {
          const verifiedIcon = entry.verifiedBy === "official" ? "✓" :
                              entry.verifiedBy === "community" ? "~" : "-"
          const verifiedColor = entry.verifiedBy === "official" ? UI.Style.TEXT_SUCCESS :
                               entry.verifiedBy === "community" ? UI.Style.TEXT_INFO :
                               UI.Style.DIM

          UI.println(`${verifiedColor}${verifiedIcon}${UI.Style.RESET} ${UI.Style.BOLD}${entry.name}${UI.Style.RESET}@${entry.version}`)

          if (entry.description) {
            UI.println(`  ${UI.Style.DIM}${entry.description}${UI.Style.RESET}`)
          }

          UI.println(`  ${UI.Style.DIM}Hash: ${entry.hash.substring(0, 16)}...${UI.Style.RESET}`)

          if (entry.author) {
            UI.println(`  ${UI.Style.DIM}Author: ${entry.author}${UI.Style.RESET}`)
          }

          UI.println()
        }
      }
    } catch (error) {
      UI.error(`Failed to list plugins: ${error instanceof Error ? error.message : String(error)}`)
      process.exit(1)
    }
  },
})

export const PluginRegistrySearchCommand = cmd({
  command: "plugin:registry:search <pattern>",
  describe: "Search for plugins in the registry",
  builder: (yargs) =>
    yargs
      .positional("pattern", {
        describe: "Search pattern (regex)",
        type: "string",
        demandOption: true,
      })
      .option("json", {
        describe: "Output as JSON",
        type: "boolean",
        default: false,
      }),
  handler: async (args) => {
    const pattern = args.pattern as string

    try {
      const results = await PluginRegistry.search(pattern)

      if (args.json) {
        console.log(JSON.stringify(results, null, 2))
      } else {
        if (results.length === 0) {
          UI.println()
          UI.println(UI.Style.TEXT_INFO + `No plugins found matching "${pattern}"` + UI.Style.RESET)
          UI.println()
          return
        }

        UI.println()
        UI.println(UI.Style.BOLD + `Search Results for "${pattern}"` + UI.Style.RESET)
        UI.println(UI.Style.DIM + `${results.length} ${results.length === 1 ? 'match' : 'matches'}` + UI.Style.RESET)
        UI.println()

        for (const entry of results) {
          UI.println(`${UI.Style.BOLD}${entry.name}${UI.Style.RESET}@${entry.version}`)
          if (entry.description) {
            UI.println(`  ${UI.Style.DIM}${entry.description}${UI.Style.RESET}`)
          }
          UI.println()
        }
      }
    } catch (error) {
      UI.error(`Search failed: ${error instanceof Error ? error.message : String(error)}`)
      process.exit(1)
    }
  },
})

export const PluginRegistrySyncCommand = cmd({
  command: "plugin:registry:sync",
  describe: "Sync with remote registry",
  builder: (yargs) =>
    yargs
      .option("url", {
        describe: "Remote registry URL",
        type: "string",
      })
      .option("json", {
        describe: "Output as JSON",
        type: "boolean",
        default: false,
      }),
  handler: async (args) => {
    try {
      const config: PluginRegistry.RegistryConfig = {
        autoUpdate: true,
      }

      if (args.url) {
        config.remoteUrl = args.url as string
      }

      // Force reload from remote
      PluginRegistry.clearCache()
      const registry = await PluginRegistry.load(config)

      if (args.json) {
        console.log(JSON.stringify({
          success: true,
          entries: registry.entries.length,
          lastUpdated: registry.lastUpdated,
        }, null, 2))
      } else {
        UI.println()
        UI.println(UI.Style.TEXT_SUCCESS + "✓ Registry synced successfully" + UI.Style.RESET)
        UI.println()
        UI.println(UI.Style.TEXT_INFO + "Entries: " + UI.Style.RESET + registry.entries.length)
        UI.println(UI.Style.TEXT_INFO + "Last Updated: " + UI.Style.RESET + new Date(registry.lastUpdated).toISOString())
        UI.println()
      }
    } catch (error) {
      UI.error(`Sync failed: ${error instanceof Error ? error.message : String(error)}`)
      process.exit(1)
    }
  },
})

export const PluginRegistryStatsCommand = cmd({
  command: "plugin:registry:stats",
  describe: "Show registry statistics",
  builder: (yargs) =>
    yargs.option("json", {
      describe: "Output as JSON",
      type: "boolean",
      default: false,
    }),
  handler: async (args) => {
    try {
      const stats = await PluginRegistry.stats()

      if (args.json) {
        console.log(JSON.stringify(stats, null, 2))
      } else {
        UI.println()
        UI.println(UI.Style.BOLD + "Registry Statistics" + UI.Style.RESET)
        UI.println()
        UI.println(UI.Style.TEXT_INFO + "Total Entries:   " + UI.Style.RESET + stats.totalEntries)
        UI.println(UI.Style.TEXT_INFO + "Unique Plugins:  " + UI.Style.RESET + stats.uniquePlugins)
        UI.println()
        UI.println(UI.Style.TEXT_SUCCESS + "Official:        " + UI.Style.RESET + stats.officialCount)
        UI.println(UI.Style.TEXT_INFO + "Community:       " + UI.Style.RESET + stats.communityCount)
        UI.println(UI.Style.DIM + "User:            " + UI.Style.RESET + stats.userCount)
        UI.println()
      }
    } catch (error) {
      UI.error(`Failed to get stats: ${error instanceof Error ? error.message : String(error)}`)
      process.exit(1)
    }
  },
})
