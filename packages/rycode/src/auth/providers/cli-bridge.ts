/**
 * CLI Provider Bridge
 *
 * Enables RyCode to communicate with LLM providers through their CLI tools
 * instead of requiring API keys. This allows users to leverage their existing
 * authenticated CLI sessions without wasting tokens via API keys.
 *
 * Supported CLI tools:
 * - claude (Anthropic Claude Code CLI)
 * - qwen (Qwen/Alibaba CLI)
 * - codex (OpenAI Codex CLI)
 * - gemini (Google Gemini CLI, if available)
 */

import { exec } from 'child_process'
import { promisify } from 'util'

const execAsync = promisify(exec)

export interface CLIProviderInfo {
  name: string
  cliCommand: string
  testCommand: string
  available: boolean
  version?: string
}

export interface CLIRequest {
  provider: string
  prompt: string
  model?: string
  options?: Record<string, any>
}

export interface CLIResponse {
  content: string
  model: string
  usage?: {
    inputTokens: number
    outputTokens: number
  }
  raw?: any
}

/**
 * CLI Provider Bridge
 * Manages communication with LLM providers through their CLI interfaces
 */
export class CLIProviderBridge {
  private availableProviders: Map<string, CLIProviderInfo> = new Map()

  constructor() {
    this.initializeProviders()
  }

  /**
   * Initialize and detect available CLI providers
   */
  private initializeProviders(): void {
    // Define supported CLI providers
    const providers = [
      {
        name: 'claude',
        cliCommand: 'claude',
        testCommand: 'claude --version',
      },
      {
        name: 'qwen',
        cliCommand: 'qwen',
        testCommand: 'qwen --version',
      },
      {
        name: 'codex',
        cliCommand: 'codex',
        testCommand: 'codex --version',
      },
      {
        name: 'gemini',
        cliCommand: 'gemini',
        testCommand: 'gemini --version',
      },
    ]

    // Check which providers are available
    providers.forEach(provider => {
      this.availableProviders.set(provider.name, {
        ...provider,
        available: false, // Will be updated by detectAvailableProviders()
      })
    })
  }

  /**
   * Detect which CLI providers are available and authenticated
   */
  async detectAvailableProviders(): Promise<CLIProviderInfo[]> {
    const results: CLIProviderInfo[] = []

    for (const [name, info] of this.availableProviders) {
      try {
        const { stdout } = await execAsync(`which ${info.cliCommand}`)
        if (stdout.trim()) {
          // CLI tool exists, try to get version
          try {
            const { stdout: versionOut } = await execAsync(info.testCommand)
            info.available = true
            info.version = versionOut.trim()
            this.availableProviders.set(name, info)
            results.push(info)
          } catch {
            // Version check failed, but CLI exists
            info.available = true
            this.availableProviders.set(name, info)
            results.push(info)
          }
        }
      } catch {
        // CLI tool not found
        info.available = false
        this.availableProviders.set(name, info)
      }
    }

    return results
  }

  /**
   * Check if a specific provider is available
   */
  isProviderAvailable(provider: string): boolean {
    const info = this.availableProviders.get(provider)
    return info?.available || false
  }

  /**
   * Send a request to Claude CLI
   */
  async sendToClaudeCLI(request: CLIRequest): Promise<CLIResponse> {
    if (!this.isProviderAvailable('claude')) {
      throw new Error('Claude CLI is not available. Install it with: npm install -g @anthropic-ai/claude-code')
    }

    try {
      // Use claude CLI in non-interactive mode with --print flag and JSON output
      const modelFlag = request.model ? `--model ${request.model}` : ''
      const escapedPrompt = request.prompt.replace(/"/g, '\\"').replace(/\n/g, ' ')
      const command = `claude --print ${modelFlag} --output-format json "${escapedPrompt}"`

      const { stdout, stderr } = await execAsync(command, {
        timeout: 120000, // 2 minute timeout for longer responses
        maxBuffer: 50 * 1024 * 1024, // 50MB buffer
      })

      if (stderr && !stdout) {
        throw new Error(`Claude CLI error: ${stderr}`)
      }

      // Parse JSON output from claude CLI
      try {
        const response = JSON.parse(stdout)
        return {
          content: response.content || response.text || stdout,
          model: response.model || request.model || 'claude-sonnet-4-5',
          usage: response.usage,
          raw: response,
        }
      } catch {
        // If JSON parsing fails, return raw text
        return {
          content: stdout,
          model: request.model || 'claude-sonnet-4-5',
          raw: { stdout, stderr },
        }
      }
    } catch (error: any) {
      throw new Error(`Failed to communicate with Claude CLI: ${error.message}`)
    }
  }

  /**
   * Send a request to Qwen CLI
   */
  async sendToQwenCLI(request: CLIRequest): Promise<CLIResponse> {
    if (!this.isProviderAvailable('qwen')) {
      throw new Error('Qwen CLI is not available. Install it with: npm install -g @qwen/cli')
    }

    try {
      // Qwen CLI command format (adjust based on actual CLI)
      const modelFlag = request.model ? `--model ${request.model}` : ''
      const command = `qwen ${modelFlag} "${request.prompt.replace(/"/g, '\\"')}"`

      const { stdout, stderr } = await execAsync(command, {
        timeout: 60000,
        maxBuffer: 10 * 1024 * 1024,
      })

      if (stderr && !stdout) {
        throw new Error(`Qwen CLI error: ${stderr}`)
      }

      return {
        content: stdout.trim(),
        model: request.model || 'qwen3-max',
        raw: { stdout, stderr },
      }
    } catch (error: any) {
      throw new Error(`Failed to communicate with Qwen CLI: ${error.message}`)
    }
  }

  /**
   * Send a request to Codex CLI (OpenAI GPT)
   */
  async sendToCodexCLI(request: CLIRequest): Promise<CLIResponse> {
    if (!this.isProviderAvailable('codex')) {
      throw new Error('Codex CLI is not available. Make sure you have the OpenAI CLI installed and authenticated')
    }

    try {
      // Codex/GPT CLI command format (using OpenAI's CLI interface)
      const modelFlag = request.model ? `--model ${request.model}` : '--model gpt-5'
      const escapedPrompt = request.prompt.replace(/"/g, '\\"').replace(/\n/g, ' ')
      const command = `codex ${modelFlag} "${escapedPrompt}"`

      const { stdout, stderr } = await execAsync(command, {
        timeout: 120000, // 2 minute timeout
        maxBuffer: 50 * 1024 * 1024, // 50MB buffer
      })

      if (stderr && !stdout) {
        throw new Error(`Codex CLI error: ${stderr}`)
      }

      // Try to parse JSON if available, otherwise return raw text
      try {
        const response = JSON.parse(stdout)
        return {
          content: response.content || response.choices?.[0]?.message?.content || stdout,
          model: response.model || request.model || 'gpt-5',
          usage: response.usage ? {
            inputTokens: response.usage.prompt_tokens || 0,
            outputTokens: response.usage.completion_tokens || 0,
          } : undefined,
          raw: response,
        }
      } catch {
        return {
          content: stdout.trim(),
          model: request.model || 'gpt-5',
          raw: { stdout, stderr },
        }
      }
    } catch (error: any) {
      throw new Error(`Failed to communicate with Codex CLI: ${error.message}`)
    }
  }

  /**
   * Send a request to Gemini CLI (Google)
   */
  async sendToGeminiCLI(request: CLIRequest): Promise<CLIResponse> {
    if (!this.isProviderAvailable('gemini')) {
      throw new Error('Gemini CLI is not available. Make sure you have gcloud CLI installed and authenticated')
    }

    try {
      // Gemini CLI command format (using gcloud or gemini CLI)
      const modelFlag = request.model ? `--model ${request.model}` : '--model gemini-2.5-pro'
      const escapedPrompt = request.prompt.replace(/"/g, '\\"').replace(/\n/g, ' ')
      const command = `gemini ${modelFlag} "${escapedPrompt}"`

      const { stdout, stderr } = await execAsync(command, {
        timeout: 120000, // 2 minute timeout
        maxBuffer: 50 * 1024 * 1024, // 50MB buffer
      })

      if (stderr && !stdout) {
        throw new Error(`Gemini CLI error: ${stderr}`)
      }

      // Try to parse JSON if available
      try {
        const response = JSON.parse(stdout)
        return {
          content: response.content || response.candidates?.[0]?.content?.parts?.[0]?.text || stdout,
          model: response.model || request.model || 'gemini-2.5-pro',
          usage: response.usageMetadata ? {
            inputTokens: response.usageMetadata.promptTokenCount || 0,
            outputTokens: response.usageMetadata.candidatesTokenCount || 0,
          } : undefined,
          raw: response,
        }
      } catch {
        return {
          content: stdout.trim(),
          model: request.model || 'gemini-2.5-pro',
          raw: { stdout, stderr },
        }
      }
    } catch (error: any) {
      throw new Error(`Failed to communicate with Gemini CLI: ${error.message}`)
    }
  }

  /**
   * Send a request to any provider
   */
  async sendRequest(request: CLIRequest): Promise<CLIResponse> {
    switch (request.provider.toLowerCase()) {
      case 'claude':
      case 'anthropic':
        return this.sendToClaudeCLI(request)

      case 'qwen':
        return this.sendToQwenCLI(request)

      case 'codex':
      case 'openai':
        return this.sendToCodexCLI(request)

      case 'gemini':
      case 'google':
        return this.sendToGeminiCLI(request)

      default:
        throw new Error(`Unknown provider: ${request.provider}`)
    }
  }

  /**
   * Get list of available providers with their models
   */
  async getAvailableProvidersWithModels(): Promise<Array<{
    provider: string
    models: string[]
    source: 'cli'
  }>> {
    const available = await this.detectAvailableProviders()

    return available.map(provider => {
      // Define latest models for each CLI provider (as of October 2025)
      const models: Record<string, string[]> = {
        claude: [
          'claude-sonnet-4-5',                 // Latest Sonnet 4.5 (Sep 2025) - best coding
          'claude-opus-4-1',                   // Opus 4.1 (Aug 2025) - strongest reasoning
          'claude-sonnet-4',                   // Sonnet 4 (May 2025)
          'claude-3-7-sonnet',                 // Claude 3.7 Sonnet (Feb 2025) - hybrid reasoning
          'claude-3-5-sonnet-20241022',        // Claude 3.5 Sonnet
          'claude-3-5-haiku-20241022',         // Latest Haiku
        ],
        qwen: [
          'qwen3-max',                         // Qwen3-Max (Sep 2025) - most capable
          'qwen3-next',                        // Qwen3-Next (Sep 2025)
          'qwen3-omni',                        // Qwen3-Omni (Sep 2025) - multimodal
          'qwen3-thinking-2507',               // Thinking model (Jul 2025)
          'qwen3-instruct-2507',               // Instruct model (Jul 2025)
          'qwen3-235b',                        // 235B MoE model
          'qwen3-32b',                         // 32B dense model
        ],
        codex: [
          'gpt-5',                             // GPT-5 (Aug 2025) - latest flagship
          'gpt-5-mini',                        // GPT-5 Mini - more usage
          'gpt-5-nano',                        // GPT-5 Nano - fastest
          'gpt-4-5',                           // GPT-4.5 (Orion) - transitional
          'gpt-4o',                            // GPT-4 Omni
          'gpt-4o-mini',                       // GPT-4 Omni Mini
          'o3',                                // O3 reasoning model
          'o3-mini',                           // O3 Mini
        ],
        gemini: [
          'gemini-2.5-pro',                    // Gemini 2.5 Pro - #1 on LMArena
          'gemini-2.5-flash',                  // Gemini 2.5 Flash (Sep 2025) - best price/perf
          'gemini-2.5-flash-lite',             // Flash-Lite - fastest, lowest cost
          'gemini-2.5-flash-image',            // Flash Image (Aug 2025) - image gen
          'gemini-2.5-computer-use',           // Computer Use model
          'gemini-2.5-deep-think',             // Deep Think - advanced reasoning
          'gemini-exp-1206',                   // Experimental
        ],
      }

      return {
        provider: provider.name,
        models: models[provider.name] || [],
        source: 'cli' as const,
      }
    })
  }

  /**
   * Test if a CLI provider is working
   */
  async testProvider(provider: string): Promise<boolean> {
    try {
      const response = await this.sendRequest({
        provider,
        prompt: 'Say "OK" if you can read this.',
      })
      return response.content.length > 0
    } catch {
      return false
    }
  }
}

/**
 * Global CLI provider bridge instance
 */
export const cliProviderBridge = new CLIProviderBridge()
