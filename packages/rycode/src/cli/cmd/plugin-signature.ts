import { cmd } from "./cmd"
import { PluginSignature } from "../../plugin/signature"
import { UI } from "../ui"
import path from "path"
import { existsSync } from "fs"

/**
 * Plugin Signature CLI Commands
 *
 * Provides commands to manage plugin signatures:
 * - plugin:sign          Sign a plugin file
 * - plugin:verify-sig    Verify a plugin signature
 * - plugin:keygen        Generate signing key pair
 */

export const PluginSignCommand = cmd({
  command: "plugin:sign <plugin-path>",
  describe: "Sign a plugin file",
  builder: (yargs) =>
    yargs
      .positional("plugin-path", {
        describe: "Path to the plugin file",
        type: "string",
        demandOption: true,
      })
      .option("key", {
        describe: "Path to private key file (PEM format)",
        type: "string",
        demandOption: true,
      })
      .option("algorithm", {
        describe: "Signature algorithm",
        type: "string",
        default: "RSA-SHA256",
        choices: ["RSA-SHA256", "RSA-SHA512"],
      })
      .option("output", {
        describe: "Output signature file path",
        type: "string",
      })
      .option("json", {
        describe: "Output as JSON",
        type: "boolean",
        default: false,
      }),
  handler: async (args) => {
    const pluginPath = path.resolve(args.pluginPath as string)
    const keyPath = path.resolve(args.key as string)

    if (!existsSync(pluginPath)) {
      UI.error(`Plugin file not found: ${pluginPath}`)
      process.exit(1)
    }

    if (!existsSync(keyPath)) {
      UI.error(`Private key file not found: ${keyPath}`)
      process.exit(1)
    }

    try {
      // Read private key
      const keyFile = Bun.file(keyPath)
      const privateKey = await keyFile.text()

      // Sign the plugin
      const signature = await PluginSignature.signWithCrypto(
        pluginPath,
        privateKey,
        args.algorithm as string
      )

      // Save signature if output specified
      if (args.output) {
        const outputPath = path.resolve(args.output as string)
        await Bun.write(outputPath, JSON.stringify(signature, null, 2))
      }

      if (args.json) {
        console.log(JSON.stringify(signature, null, 2))
      } else {
        UI.println()
        UI.println(UI.Style.BOLD + "✓ Plugin Signed Successfully" + UI.Style.RESET)
        UI.println()
        UI.println(UI.Style.TEXT_INFO + "File:      " + UI.Style.RESET + pluginPath)
        UI.println(UI.Style.TEXT_INFO + "Algorithm: " + UI.Style.RESET + signature.algorithm)
        UI.println(UI.Style.TEXT_INFO + "Key ID:    " + UI.Style.RESET + signature.keyId)
        UI.println(UI.Style.TEXT_INFO + "Timestamp: " + UI.Style.RESET + new Date(signature.timestamp).toISOString())
        UI.println()

        if (args.output) {
          UI.println(UI.Style.TEXT_SUCCESS + `Signature saved to: ${args.output}` + UI.Style.RESET)
          UI.println()
        }

        UI.println(UI.Style.DIM + "Signature (base64):" + UI.Style.RESET)
        UI.println(UI.Style.DIM + signature.signature.substring(0, 80) + "..." + UI.Style.RESET)
        UI.println()
      }
    } catch (error) {
      UI.error(`Failed to sign plugin: ${error instanceof Error ? error.message : String(error)}`)
      process.exit(1)
    }
  },
})

export const PluginVerifySignatureCommand = cmd({
  command: "plugin:verify-sig <plugin-path> <signature-file>",
  describe: "Verify a plugin signature",
  builder: (yargs) =>
    yargs
      .positional("plugin-path", {
        describe: "Path to the plugin file",
        type: "string",
        demandOption: true,
      })
      .positional("signature-file", {
        describe: "Path to signature file (JSON)",
        type: "string",
        demandOption: true,
      })
      .option("public-key", {
        describe: "Path to public key file (PEM format)",
        type: "string",
        demandOption: true,
      })
      .option("json", {
        describe: "Output as JSON",
        type: "boolean",
        default: false,
      }),
  handler: async (args) => {
    const pluginPath = path.resolve(args.pluginPath as string)
    const sigPath = path.resolve(args.signatureFile as string)
    const pubKeyPath = path.resolve(args.publicKey as string)

    if (!existsSync(pluginPath)) {
      UI.error(`Plugin file not found: ${pluginPath}`)
      process.exit(1)
    }

    if (!existsSync(sigPath)) {
      UI.error(`Signature file not found: ${sigPath}`)
      process.exit(1)
    }

    if (!existsSync(pubKeyPath)) {
      UI.error(`Public key file not found: ${pubKeyPath}`)
      process.exit(1)
    }

    try {
      // Read signature
      const sigFile = Bun.file(sigPath)
      const sigData = await sigFile.json()
      const signature = PluginSignature.Signature.parse(sigData)

      // Read public key
      const pubKeyFile = Bun.file(pubKeyPath)
      const publicKey = await pubKeyFile.text()

      // Verify signature
      const result = await PluginSignature.verifyCryptoSignature(
        pluginPath,
        signature,
        publicKey
      )

      if (args.json) {
        console.log(JSON.stringify({
          valid: result.valid,
          error: result.error,
          signature: signature,
        }, null, 2))
      } else {
        UI.println()
        UI.println(UI.Style.BOLD + "Signature Verification" + UI.Style.RESET)
        UI.println()

        UI.println(UI.Style.TEXT_INFO + "File:      " + UI.Style.RESET + pluginPath)
        UI.println(UI.Style.TEXT_INFO + "Algorithm: " + UI.Style.RESET + signature.algorithm)
        UI.println(UI.Style.TEXT_INFO + "Key ID:    " + UI.Style.RESET + signature.keyId)
        UI.println(UI.Style.TEXT_INFO + "Signed:    " + UI.Style.RESET + new Date(signature.timestamp).toISOString())
        UI.println()

        if (result.valid) {
          UI.println(UI.Style.TEXT_SUCCESS + "✓ Signature is VALID" + UI.Style.RESET)
          UI.println(UI.Style.TEXT_SUCCESS + "  Plugin has been signed by the claimed signer." + UI.Style.RESET)
        } else {
          UI.println(UI.Style.TEXT_DANGER + "✗ Signature is INVALID" + UI.Style.RESET)
          UI.println(UI.Style.TEXT_DANGER + `  ${result.error}` + UI.Style.RESET)
          UI.println()
          UI.println(UI.Style.TEXT_WARNING + "⚠ DO NOT use this plugin." + UI.Style.RESET)
          process.exitCode = 1
        }

        UI.println()
      }
    } catch (error) {
      UI.error(`Failed to verify signature: ${error instanceof Error ? error.message : String(error)}`)
      process.exit(1)
    }
  },
})

export const PluginKeygenCommand = cmd({
  command: "plugin:keygen",
  describe: "Generate RSA key pair for plugin signing",
  builder: (yargs) =>
    yargs
      .option("output", {
        describe: "Output directory for keys",
        type: "string",
        default: ".",
      })
      .option("name", {
        describe: "Key name prefix",
        type: "string",
        default: "plugin-signing",
      })
      .option("json", {
        describe: "Output as JSON",
        type: "boolean",
        default: false,
      }),
  handler: async (args) => {
    try {
      const keyPair = PluginSignature.generateKeyPair()

      const outputDir = path.resolve(args.output as string)
      const keyName = args.name as string

      const privateKeyPath = path.join(outputDir, `${keyName}-private.pem`)
      const publicKeyPath = path.join(outputDir, `${keyName}-public.pem`)

      // Save keys
      await Bun.write(privateKeyPath, keyPair.privateKey)
      await Bun.write(publicKeyPath, keyPair.publicKey)

      if (args.json) {
        console.log(JSON.stringify({
          keyId: keyPair.keyId,
          privateKeyPath,
          publicKeyPath,
        }, null, 2))
      } else {
        UI.println()
        UI.println(UI.Style.BOLD + "✓ RSA Key Pair Generated" + UI.Style.RESET)
        UI.println()
        UI.println(UI.Style.TEXT_INFO + "Key ID:      " + UI.Style.RESET + keyPair.keyId)
        UI.println(UI.Style.TEXT_INFO + "Private Key: " + UI.Style.RESET + privateKeyPath)
        UI.println(UI.Style.TEXT_INFO + "Public Key:  " + UI.Style.RESET + publicKeyPath)
        UI.println()
        UI.println(UI.Style.TEXT_WARNING + "⚠ Keep your private key secure!" + UI.Style.RESET)
        UI.println(UI.Style.DIM + "  Add it to .gitignore and never commit it." + UI.Style.RESET)
        UI.println()
        UI.println(UI.Style.TEXT_INFO + "Next steps:" + UI.Style.RESET)
        UI.println(UI.Style.DIM + `  1. Sign a plugin: rycode plugin:sign <plugin-path> --key ${privateKeyPath}` + UI.Style.RESET)
        UI.println(UI.Style.DIM + `  2. Verify signature: rycode plugin:verify-sig <plugin-path> <sig.json> --public-key ${publicKeyPath}` + UI.Style.RESET)
        UI.println()
      }
    } catch (error) {
      UI.error(`Failed to generate key pair: ${error instanceof Error ? error.message : String(error)}`)
      process.exit(1)
    }
  },
})
