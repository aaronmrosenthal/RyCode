/**
 * RyCode Provider Authentication - Exports
 *
 * Complete provider authentication system with security, smart features,
 * and all provider implementations.
 */

// Security
export { RateLimiter, authRateLimiter, apiRateLimiter } from './security/rate-limiter'
export { InputValidator, inputValidator } from './security/input-validator'
export {
  CircuitBreaker,
  CircuitBreakerRegistry,
  circuitBreakerRegistry,
  CircuitBreakerError
} from './security/circuit-breaker'

export type { RateLimitConfig, RateLimitResult } from './security/rate-limiter'
export type { ValidationResult } from './security/input-validator'
export type { CircuitState, CircuitBreakerConfig, CircuitBreakerStats } from './security/circuit-breaker'

// Errors
export {
  AuthenticationError,
  InvalidAPIKeyError,
  ExpiredCredentialsError,
  RateLimitError,
  NetworkError,
  ValidationError,
  StorageError,
  CompromisedKeyError,
  parseHTTPError,
  parseNetworkError,
  ErrorHandler
} from './errors'

export type { ErrorReason, ErrorContext } from './errors'

// Smart Features
export { SmartProviderSetup, smartSetup } from './auto-detect'
export { CostTracker, costTracker } from './cost-tracker'
export { ModelRecommender, modelRecommender } from './model-recommender'

export type {
  DetectedCredential,
  AutoDetectResult
} from './auto-detect'

export type {
  ModelPricing,
  UsageRecord,
  CostSummary,
  CostBreakdown,
  CostSavingTip
} from './cost-tracker'

export type {
  TaskContext,
  ModelRecommendation,
  ModelCapabilities
} from './model-recommender'

// Providers
export { anthropicProvider } from './providers/anthropic'
export { openaiProvider } from './providers/openai'
export { grokProvider } from './providers/grok'
export { qwenProvider } from './providers/qwen'
export { googleProvider } from './providers/google'

export type {
  AnthropicAuthConfig,
  AnthropicAuthResult
} from './providers/anthropic'

export type {
  OpenAIAuthConfig,
  OpenAIAuthResult
} from './providers/openai'

export type {
  GrokAuthConfig,
  GrokAuthResult
} from './providers/grok'

export type {
  QwenAuthConfig,
  QwenAuthResult
} from './providers/qwen'

export type {
  GoogleAuthConfig,
  GoogleAuthResult,
  CSRFToken
} from './providers/google'

// Provider Registry
export { ProviderRegistry, providerRegistry } from './provider-registry'

export type {
  ProviderAuthConfig,
  ProviderAuthResult,
  ProviderInfo
} from './provider-registry'

// Storage
export { CredentialStore, credentialStore } from './storage/credential-store'
export { AuditLog, auditLog } from './storage/audit-log'

export type { StoredCredential } from './storage/credential-store'
export type {
  AuditEventType,
  AuditEvent,
  AuditQuery,
  AuditSummary
} from './storage/audit-log'

// Unified Auth Manager
export { AuthManager, authManager } from './auth-manager'

export type {
  AuthManagerConfig,
  AuthenticateOptions,
  AuthStatus
} from './auth-manager'
