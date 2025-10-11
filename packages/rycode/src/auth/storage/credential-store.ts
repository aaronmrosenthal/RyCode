/**
 * Credential Store - Storage Adapter
 *
 * Bridges the new provider authentication system with existing RyCode
 * storage infrastructure. Provides a clean interface for storing and
 * retrieving provider credentials with encryption and integrity checks.
 */

import { Auth } from '../index'
import { inputValidator } from '../security/input-validator'
import { StorageError } from '../errors'
import type { ProviderAuthResult } from '../provider-registry'

export interface StoredCredential {
  provider: string
  method: 'api-key' | 'oauth' | 'cli'
  createdAt: Date
  updatedAt: Date
  expiresAt?: Date
  metadata?: Record<string, any>
}

export class CredentialStore {
  /**
   * Store provider credentials
   */
  async store(
    provider: string,
    authResult: ProviderAuthResult
  ): Promise<void> {
    try {
      // Validate provider
      if (!provider || typeof provider !== 'string') {
        throw new Error('Invalid provider name')
      }

      // Convert to Auth.Info format based on method
      let authInfo: Auth.Info

      if (authResult.method === 'api-key' || authResult.method === 'api') {
        // API key authentication
        const apiKey = authResult['apiKey'] || authResult['credential']

        if (!apiKey) {
          throw new Error('API key is required for api-key authentication')
        }

        // Validate before storing
        const sanitized = inputValidator.sanitizeAPIKey(apiKey)

        authInfo = {
          type: 'api',
          key: sanitized
        }
      } else if (authResult.method === 'oauth') {
        // OAuth authentication
        const { oauthToken, refreshToken } = authResult

        if (!oauthToken && !refreshToken) {
          throw new Error('OAuth token or refresh token is required')
        }

        authInfo = {
          type: 'oauth',
          access: oauthToken || '',
          refresh: refreshToken || '',
          expires: authResult.expiresAt
            ? authResult.expiresAt.getTime()
            : Date.now() + 3600000 // 1 hour default
        }
      } else {
        throw new Error(`Unsupported authentication method: ${authResult.method}`)
      }

      // Store using existing Auth namespace
      await Auth.set(provider, authInfo)
    } catch (error) {
      throw new StorageError(
        provider,
        'store',
        {
          provider,
          timestamp: new Date(),
          originalError: error instanceof Error ? error : undefined
        }
      )
    }
  }

  /**
   * Retrieve provider credentials
   */
  async retrieve(provider: string): Promise<StoredCredential | null> {
    try {
      const authInfo = await Auth.get(provider)

      if (!authInfo) {
        return null
      }

      // Convert from Auth.Info to StoredCredential
      const credential: StoredCredential = {
        provider,
        method: authInfo.type === 'oauth' ? 'oauth' : 'api-key',
        createdAt: new Date(), // We don't track this yet
        updatedAt: new Date()
      }

      if (authInfo.type === 'oauth') {
        credential.expiresAt = new Date(authInfo.expires)
      }

      return credential
    } catch (error) {
      throw new StorageError(
        provider,
        'retrieve',
        {
          provider,
          timestamp: new Date(),
          originalError: error instanceof Error ? error : undefined
        }
      )
    }
  }

  /**
   * Get raw auth info (for internal use)
   */
  async getAuthInfo(provider: string): Promise<Auth.Info | null> {
    try {
      const authInfo = await Auth.get(provider)
      return authInfo || null
    } catch (error) {
      throw new StorageError(
        provider,
        'retrieve',
        {
          provider,
          timestamp: new Date(),
          originalError: error instanceof Error ? error : undefined
        }
      )
    }
  }

  /**
   * Check if provider has stored credentials
   */
  async has(provider: string): Promise<boolean> {
    try {
      const authInfo = await Auth.get(provider)
      return authInfo !== undefined
    } catch {
      return false
    }
  }

  /**
   * Remove provider credentials
   */
  async remove(provider: string): Promise<boolean> {
    try {
      return await Auth.remove(provider)
    } catch (error) {
      throw new StorageError(
        provider,
        'remove',
        {
          provider,
          timestamp: new Date(),
          originalError: error instanceof Error ? error : undefined
        }
      )
    }
  }

  /**
   * List all stored providers
   */
  async list(): Promise<string[]> {
    try {
      const allAuth = await Auth.all()
      return Object.keys(allAuth)
    } catch (error) {
      throw new StorageError(
        'unknown',
        'list',
        {
          timestamp: new Date(),
          originalError: error instanceof Error ? error : undefined
        }
      )
    }
  }

  /**
   * Get all stored credentials
   */
  async listAll(): Promise<StoredCredential[]> {
    try {
      const providers = await this.list()
      const credentials: StoredCredential[] = []

      for (const provider of providers) {
        const credential = await this.retrieve(provider)
        if (credential) {
          credentials.push(credential)
        }
      }

      return credentials
    } catch (error) {
      throw new StorageError(
        'unknown',
        'list',
        {
          timestamp: new Date(),
          originalError: error instanceof Error ? error : undefined
        }
      )
    }
  }

  /**
   * Check if credentials are expired (for OAuth)
   */
  async isExpired(provider: string): Promise<boolean> {
    try {
      const authInfo = await Auth.get(provider)

      if (!authInfo || authInfo.type !== 'oauth') {
        return false
      }

      return authInfo.expires < Date.now()
    } catch {
      return true
    }
  }

  /**
   * Update OAuth tokens (refresh)
   */
  async updateOAuthTokens(
    provider: string,
    accessToken: string,
    refreshToken?: string,
    expiresAt?: Date
  ): Promise<void> {
    try {
      const authInfo = await Auth.get(provider)

      if (!authInfo || authInfo.type !== 'oauth') {
        throw new Error('Provider does not use OAuth authentication')
      }

      // Update tokens
      const updated: Auth.Info = {
        type: 'oauth',
        access: accessToken,
        refresh: refreshToken || authInfo.refresh,
        expires: expiresAt ? expiresAt.getTime() : Date.now() + 3600000
      }

      await Auth.set(provider, updated)
    } catch (error) {
      throw new StorageError(
        provider,
        'update',
        {
          provider,
          timestamp: new Date(),
          originalError: error instanceof Error ? error : undefined
        }
      )
    }
  }

  /**
   * Clear all credentials (use with caution!)
   */
  async clear(): Promise<number> {
    try {
      const providers = await this.list()
      let count = 0

      for (const provider of providers) {
        const removed = await this.remove(provider)
        if (removed) count++
      }

      return count
    } catch (error) {
      throw new StorageError(
        'unknown',
        'clear',
        {
          timestamp: new Date(),
          originalError: error instanceof Error ? error : undefined
        }
      )
    }
  }

  /**
   * Migrate to encrypted storage
   */
  async migrateToEncrypted(): Promise<number> {
    try {
      return await Auth.migrateToEncrypted()
    } catch (error) {
      throw new StorageError(
        'unknown',
        'migrate',
        {
          timestamp: new Date(),
          originalError: error instanceof Error ? error : undefined
        }
      )
    }
  }

  /**
   * Export credentials (for backup)
   */
  async export(): Promise<Record<string, StoredCredential>> {
    try {
      const credentials = await this.listAll()
      const exported: Record<string, StoredCredential> = {}

      for (const cred of credentials) {
        exported[cred.provider] = cred
      }

      return exported
    } catch (error) {
      throw new StorageError(
        'unknown',
        'export',
        {
          timestamp: new Date(),
          originalError: error instanceof Error ? error : undefined
        }
      )
    }
  }
}

/**
 * Global credential store instance
 */
export const credentialStore = new CredentialStore()
