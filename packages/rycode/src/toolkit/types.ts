/**
 * TypeScript type definitions for toolkit-cli Node.js client
 * Auto-generated from Python Pydantic models
 */

// ==================== Core Types ====================

export type LogLevel = 'silent' | 'error' | 'warn' | 'info' | 'debug';
export type ComplexityLevel = 'low' | 'medium' | 'high' | 'enterprise';
export type AgentType =
  | 'claude'      // Anthropic Claude (Opus, Sonnet, Haiku)
  | 'gemini'      // Google Gemini
  | 'qwen'        // Alibaba Qwen
  | 'codex'       // OpenAI Codex
  | 'gpt4'        // OpenAI GPT-4
  | 'deepseek'    // DeepSeek Coder
  | 'llama'       // Meta Llama
  | 'mistral'     // Mistral AI
  | 'rycode';     // RyCode AI Agent
export type ErrorCode =
  | 'VALIDATION_ERROR'
  | 'API_ERROR'
  | 'TIMEOUT_ERROR'
  | 'RATE_LIMIT_ERROR'
  | 'NOT_FOUND_ERROR'
  | 'PYTHON_ERROR'
  | 'UNKNOWN_ERROR';

export type MonitorMode =
  | 'start'
  | 'stop'
  | 'status'
  | 'dashboard'
  | 'metrics'
  | 'logs'
  | 'errors'
  | 'profile'
  | 'analyze'
  | 'fix';

// ==================== Configuration ====================

export interface ApiKeys {
  anthropic?: string;    // Claude
  openai?: string;       // GPT-4, Codex
  google?: string;       // Gemini
  qwen?: string;         // Qwen
  deepseek?: string;     // DeepSeek
  together?: string;     // Llama, Mistral (via Together AI)
  rycode?: string;       // RyCode API key
}

export interface ToolkitConfig {
  pythonPath?: string;
  toolkitCliPath?: string;
  agents?: AgentType[];
  apiKeys?: ApiKeys;
  timeout?: number;
  maxConcurrent?: number;
  logLevel?: LogLevel;
}

// ==================== Command Results ====================

export interface CommandMetadata {
  executionTime: number;
  agentsUsed: string[];
  tokensConsumed?: number;
  cost?: number;
  timestamp: string;
}

export interface CommandError {
  code: ErrorCode;
  message: string;
  details?: Record<string, any>;
  stack?: string;
}

export interface CommandResult<T> {
  success: boolean;
  data?: T;
  error?: CommandError;
  metadata: CommandMetadata;
}

// ==================== Oneshot Models ====================

export interface AcceptanceScenario {
  id: string;
  given: string;
  when: string;
  then: string;
}

export interface Requirement {
  id: string;
  priority: 'critical' | 'high' | 'medium' | 'low';
  description: string;
  acceptanceCriteria: string[];
}

export interface Specification {
  overview: string;
  userStory: string;
  acceptanceScenarios: AcceptanceScenario[];
  functionalRequirements: Requirement[];
  technicalRequirements: Requirement[];
  uxRequirements: Requirement[];
}

export interface Component {
  name: string;
  responsibility: string;
  interfaces: string[];
}

export interface Architecture {
  overview: string;
  components: Component[];
  dataFlow: string;
  techStack: string[];
}

export interface Phase {
  id: number;
  name: string;
  description: string;
  tasks: string[];
  estimatedWeeks: number;
}

export interface Milestone {
  name: string;
  description: string;
  criteria: string[];
}

export interface Roadmap {
  phases: Phase[];
  milestones: Milestone[];
  timeline: string;
}

export interface UxDesign {
  screenName: string;
  asciiArt: string;
  description: string;
}

export interface OneshotMetadata extends CommandMetadata {
  complexity: ComplexityLevel;
  includesUx: boolean;
}

export interface OneshotResult {
  specification: Specification;
  architecture: Architecture;
  roadmap: Roadmap;
  uxDesigns?: UxDesign[];
  metadata: OneshotMetadata;
}

// ==================== Specify Models ====================

export interface MultiAgentInsight {
  agent: string;
  perspective: string;
  keyFindings: string[];
  recommendations: string[];
}

export interface SpecifyResult {
  specification: Specification;
  multiAgentInsights?: MultiAgentInsight[];
}

// ==================== Fix Models ====================

export interface CodeChange {
  file: string;
  before: string;
  after: string;
  explanation: string;
}

export interface Solution {
  approach: string;
  codeChanges: CodeChange[];
  testCases: string[];
  verificationSteps: string[];
}

export interface FixResult {
  issue: string;
  rootCause: string;
  solution: Solution;
  preventionStrategy: string;
}

// ==================== Command Options ====================

export interface ProgressChunk {
  phase: string;
  message: string;
  progress: number;
}

export interface OneshotOptions {
  agents?: AgentType[];
  complexity?: ComplexityLevel;
  includeUx?: boolean;
  flags?: ('creative' | 'deep' | 'fast')[];
  timeout?: number;
  onProgress?: (chunk: ProgressChunk) => void;
  signal?: AbortSignal;
}

export interface SpecifyOptions {
  agents?: AgentType[];
  flags?: ('creative' | 'deep' | 'validation' | 'research')[];
  timeout?: number;
  signal?: AbortSignal;
}

export interface FixOptions {
  agents?: AgentType[];
  context?: string;
  autoApply?: boolean;
  timeout?: number;
  signal?: AbortSignal;
}

export interface PlanOptions {
  agents?: AgentType[];
  timeout?: number;
  signal?: AbortSignal;
}

export interface TasksOptions {
  agents?: AgentType[];
  timeout?: number;
  signal?: AbortSignal;
}

export interface ImplementOptions {
  agents?: AgentType[];
  timeout?: number;
  signal?: AbortSignal;
}

export interface MakeOptions {
  agents?: AgentType[];
  timeout?: number;
  signal?: AbortSignal;
}

// ==================== Health & Status ====================

export interface AgentStatus {
  name: string;
  available: boolean;
  configured: boolean;
  error?: string;
}

export interface HealthStatus {
  healthy: boolean;
  toolkitCliInstalled: boolean;
  toolkitCliVersion?: string;
  pythonVersion?: string;
  agentsAvailable: AgentStatus[];
  issues: string[];
}

export interface VersionInfo {
  nodePackage: string;
  toolkitCli: string;
  compatible: boolean;
  minimumRequired: string;
  error?: string;
}

// ==================== Queue Status ====================

export interface QueueStatus {
  active: number;
  queued: number;
}

// ==================== Generic Result Types ====================

export interface PlanResult {
  plan: string;
  phases: Phase[];
  risks: string[];
}

export interface TasksResult {
  tasks: string[];
  prioritized: boolean;
}

export interface ImplementResult {
  filesChanged: string[];
  summary: string;
}

export interface MakeResult {
  guidance: string;
  bestPractices: string[];
}
