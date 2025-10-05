/**
 * Enhanced Auth Command with Polished Messaging
 * Professional installer experience
 */

import { Auth } from "../../auth"
import { cmd } from "./cmd"
import * as prompts from "@clack/prompts"
import { UI } from "../ui"
import { InstallerMessages } from "../installer-messages"
import { EnhancedTUI } from "../tui-enhanced"
import { ModelsDev } from "../../provider/models"
import { map, pipe, sortBy, values } from "remeda"
import path from "path"
import os from "os"
import { Global } from "../../global"
import { Plugin } from "../../plugin"
import { Instance } from "../../project/instance"

/**
 * Enhanced login command with polished UI
 */
export const AuthLoginEnhancedCommand = cmd({
  command: "login [url]",
  describe: "log in to a provider with enhanced UI",
  builder: (yargs) =>
    yargs.positional("url", {
      describe: "opencode auth provider",
      type: "string",
    }),
  async handler(args) {
    await Instance.provide({
      directory: process.cwd(),
      async fn() {
        // Welcome screen
        InstallerMessages.authRequired()

        await prompts.intro(
          UI.Style.MATRIX_GREEN + "🔐 Provider Authentication" + UI.Style.RESET
        )

        // Handle wellknown URL if provided
        if (args.url) {
          try {
            prompts.log.info(`Connecting to ${args.url}...`)
            const wellknown = await fetch(`${args.url}/.well-known/opencode`).then((x) => x.json())

            const spinner = prompts.spinner()
            spinner.start(`Running authentication command...`)

            const proc = Bun.spawn({
              cmd: wellknown.auth.command,
              stdout: "pipe",
            })

            const exit = await proc.exited
            if (exit !== 0) {
              spinner.stop("Authentication failed", 1)
              InstallerMessages.error(
                "Authentication Failed",
                `The authentication command exited with code ${exit}`,
                "Check your credentials and try again"
              )
              prompts.outro("Setup cancelled")
              return
            }

            const token = await new Response(proc.stdout).text()
            await Auth.set(args.url, {
              type: "wellknown",
              key: wellknown.auth.env,
              token: token.trim(),
            })

            spinner.stop("Successfully authenticated")
            InstallerMessages.authSuccess(args.url)
            prompts.outro(UI.Style.MATRIX_GREEN + "✓ Ready to go!" + UI.Style.RESET)
            return
          } catch (error) {
            InstallerMessages.error(
              "Connection Failed",
              `Could not connect to ${args.url}`,
              "Verify the URL is correct and try again"
            )
            prompts.outro("Setup cancelled")
            return
          }
        }

        // Show provider intro
        InstallerMessages.providerIntro()

        // Refresh provider database
        await ModelsDev.refresh().catch(() => {})
        const providers = await ModelsDev.get()

        const priority: Record<string, number> = {
          anthropic: 0,
          openai: 1,
          opencode: 2,
          google: 3,
          "github-copilot": 4,
          openrouter: 5,
          vercel: 6,
        }

        // Provider selection with enhanced UI
        let provider = await prompts.autocomplete({
          message: "Select your AI provider",
          maxItems: 8,
          options: [
            ...pipe(
              providers,
              values(),
              sortBy(
                (x) => priority[x.id] ?? 99,
                (x) => x.name ?? x.id
              ),
              map((x) => ({
                label: x.name,
                value: x.id,
                hint: priority[x.id] <= 1 ? "⭐ recommended" : undefined,
              }))
            ),
            {
              value: "other",
              label: "Other Provider",
            },
          ],
        })

        if (prompts.isCancel(provider)) {
          prompts.outro("Setup cancelled")
          throw new UI.CancelledError()
        }

        // Handle plugin-based auth
        const plugin = await Plugin.list().then((x) => x.find((x) => x.auth?.provider === provider))

        if (plugin && plugin.auth) {
          let index = 0
          if (plugin.auth.methods.length > 1) {
            const method = await prompts.select({
              message: "Choose login method",
              options: [
                ...plugin.auth.methods.map((x, index) => ({
                  label: x.label,
                  value: index.toString(),
                })),
              ],
            })
            if (prompts.isCancel(method)) {
              prompts.outro("Setup cancelled")
              throw new UI.CancelledError()
            }
            index = parseInt(method)
          }

          const method = plugin.auth.methods[index]
          if (method.type === "oauth") {
            await new Promise((resolve) => setTimeout(resolve, 10))
            const authorize = await method.authorize()

            if (authorize.url) {
              UI.println()
              UI.println(UI.Style.MATRIX_GREEN + "═".repeat(70) + UI.Style.RESET)
              UI.println()
              UI.println(UI.Style.BOLD + "🌐 Authorization Required" + UI.Style.RESET)
              UI.println()
              UI.println("  Open this URL in your browser:")
              UI.println()
              UI.println("  " + UI.link(authorize.url, authorize.url))
              UI.println()
              UI.println(UI.Style.MATRIX_GREEN + "═".repeat(70) + UI.Style.RESET)
              UI.println()
            }

            if (authorize.method === "auto") {
              if (authorize.instructions) {
                prompts.log.info(authorize.instructions)
              }

              const spinner = prompts.spinner()
              spinner.start("Waiting for authorization...")

              const result = await authorize.callback()

              if (result.type === "failed") {
                spinner.stop("Authorization failed", 1)
                InstallerMessages.error(
                  "Authorization Failed",
                  "Could not complete the OAuth flow",
                  "Try again or use a different provider"
                )
                prompts.outro("Setup cancelled")
                return
              }

              if (result.type === "success") {
                if ("refresh" in result) {
                  await Auth.set(provider, {
                    type: "oauth",
                    refresh: result.refresh,
                    access: result.access,
                    expires: result.expires,
                  })
                }
                if ("key" in result) {
                  await Auth.set(provider, {
                    type: "api",
                    key: result.key,
                  })
                }
                spinner.stop("Authorization successful")

                const providerName = providers[provider]?.name || provider
                InstallerMessages.authSuccess(providerName)
              }
            }

            if (authorize.method === "code") {
              const code = await prompts.text({
                message: "Paste the authorization code here:",
                validate: (x) => (x && x.length > 0 ? undefined : "Authorization code is required"),
              })

              if (prompts.isCancel(code)) {
                prompts.outro("Setup cancelled")
                throw new UI.CancelledError()
              }

              const spinner = prompts.spinner()
              spinner.start("Verifying code...")

              const result = await authorize.callback(code)

              if (result.type === "failed") {
                spinner.stop("Verification failed", 1)
                InstallerMessages.error(
                  "Invalid Code",
                  "The authorization code was rejected",
                  "Make sure you copied the entire code and try again"
                )
                prompts.outro("Setup cancelled")
                return
              }

              if (result.type === "success") {
                if ("refresh" in result) {
                  await Auth.set(provider, {
                    type: "oauth",
                    refresh: result.refresh,
                    access: result.access,
                    expires: result.expires,
                  })
                }
                if ("key" in result) {
                  await Auth.set(provider, {
                    type: "api",
                    key: result.key,
                  })
                }

                spinner.stop("Code verified successfully")

                const providerName = providers[provider]?.name || provider
                InstallerMessages.authSuccess(providerName)
              }
            }

            prompts.outro(UI.Style.MATRIX_GREEN + "✓ Authentication complete!" + UI.Style.RESET)
            return
          }
        }

        // Handle "other" provider
        if (provider === "other") {
          provider = await prompts.text({
            message: "Enter provider ID",
            validate: (x) => (x && x.match(/^[0-9a-z-]+$/) ? undefined : "Use lowercase letters, numbers, and hyphens only"),
          })

          if (prompts.isCancel(provider)) {
            prompts.outro("Setup cancelled")
            throw new UI.CancelledError()
          }

          provider = provider.replace(/^@ai-sdk\//, "")

          InstallerMessages.warning(
            "Custom Provider",
            `This will store credentials for "${provider}".\nYou'll need to configure it in opencode.json separately.\nCheck the documentation for examples.`
          )
        }

        // Handle cloud providers with env vars
        if (provider === "amazon-bedrock") {
          InstallerMessages.info(
            "AWS Bedrock Configuration",
            "Amazon Bedrock uses AWS credentials.\n\n" +
            "Configure using environment variables:\n" +
            "• AWS_BEARER_TOKEN_BEDROCK\n" +
            "• AWS_PROFILE\n" +
            "• AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY"
          )
          prompts.outro("Configuration info displayed")
          return
        }

        if (provider === "google-vertex") {
          InstallerMessages.info(
            "Google Vertex AI Configuration",
            "Vertex AI uses Application Default Credentials.\n\n" +
            "Setup options:\n" +
            "• Set GOOGLE_APPLICATION_CREDENTIALS environment variable\n" +
            "• Run: gcloud auth application-default login\n" +
            "• Optionally set GOOGLE_CLOUD_PROJECT and GOOGLE_CLOUD_LOCATION"
          )
          prompts.outro("Configuration info displayed")
          return
        }

        // Show API key help
        InstallerMessages.showApiKeyHelp(provider)

        // Get API key
        const key = await prompts.password({
          message: "Enter your API key",
          validate: (x) => (x && x.length > 0 ? undefined : "API key is required"),
        })

        if (prompts.isCancel(key)) {
          prompts.outro("Setup cancelled")
          throw new UI.CancelledError()
        }

        // Save credentials
        const spinner = prompts.spinner()
        spinner.start("Saving credentials...")

        await Auth.set(provider, {
          type: "api",
          key,
        })

        spinner.stop("Credentials saved")

        // Success message
        const providerName = providers[provider]?.name || provider
        InstallerMessages.authSuccess(providerName)

        // Show next steps
        UI.println(UI.Style.BOLD + "Next Steps:" + UI.Style.RESET)
        UI.println()
        UI.println(`  ${UI.Style.MATRIX_GREEN}1.${UI.Style.RESET} Start RyCode with: ${UI.Style.BOLD}rycode${UI.Style.RESET}`)
        UI.println(`  ${UI.Style.MATRIX_GREEN}2.${UI.Style.RESET} Ask AI to help build features`)
        UI.println(`  ${UI.Style.MATRIX_GREEN}3.${UI.Style.RESET} Enjoy 10x faster development!`)
        UI.println()

        prompts.outro(UI.Style.MATRIX_GREEN + "✓ Setup complete - ready to code!" + UI.Style.RESET)
      },
    })
  },
})
