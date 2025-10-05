import crypto from "crypto"
import { promisify } from "util"

const scryptAsync = promisify(crypto.scrypt)

export namespace APIKey {
  export function generate(): string {
    return crypto.randomBytes(32).toString("base64url")
  }

  export async function hash(apiKey: string): Promise<string> {
    const salt = crypto.randomBytes(16).toString("hex")
    const derivedKey = (await scryptAsync(apiKey, salt, 64)) as Buffer
    return salt + ":" + derivedKey.toString("hex")
  }

  export async function verify(apiKey: string, storedHash: string): Promise<boolean> {
    const parts = storedHash.split(":")
    if (parts.length !== 2) {
      return false
    }

    const salt = parts[0]
    const hash = parts[1]
    const hashBuffer = Buffer.from(hash, "hex")
    const derivedKey = (await scryptAsync(apiKey, salt, 64)) as Buffer

    try {
      return crypto.timingSafeEqual(hashBuffer, derivedKey)
    } catch {
      return false
    }
  }

  export function isHashed(key: string): boolean {
    return key.includes(":") && key.split(":").length === 2 && key.split(":")[0].length === 32
  }

  export async function migrate(plaintextKey: string): Promise<string> {
    return await hash(plaintextKey)
  }
}
