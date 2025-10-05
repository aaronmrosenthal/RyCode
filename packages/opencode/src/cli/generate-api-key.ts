import { APIKey } from "../server/middleware/api-key"

async function main() {
  const plaintext = APIKey.generate()
  const hashed = await APIKey.hash(plaintext)

  console.log("\n" + "=".repeat(70))
  console.log("  🔑 NEW API KEY GENERATED")
  console.log("=".repeat(70))
  console.log("\n⚠️  COPY THIS KEY NOW - IT WILL NOT BE SHOWN AGAIN\n")
  console.log("   " + plaintext)
  console.log("\n" + "=".repeat(70))
  console.log("\nAdd this to your opencode.json config:\n")
  console.log('"server": {')
  console.log('  "require_auth": true,')
  console.log('  "api_keys": [')
  console.log('    "' + hashed + '"')
  console.log("  ]")
  console.log("}")
  console.log("\n" + "=".repeat(70))
  console.log("\n✅ Key is hashed using scrypt (secure storage)")
  console.log("✅ 256-bit entropy (cryptographically secure)")
  console.log("✅ Safe to commit hashed value to version control\n")
}

main().catch(console.error)
