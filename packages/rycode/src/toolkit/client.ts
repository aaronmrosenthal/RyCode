/**
 * Main ToolkitClient class for invoking toolkit-cli commands
 */

import { spawn, ChildProcess } from 'child_process';
import {
  ToolkitConfig,
  CommandResult,
  HealthStatus,
  VersionInfo,
  QueueStatus,
  OneshotResult,
  OneshotOptions,
  SpecifyResult,
  SpecifyOptions,
  FixResult,
  FixOptions,
  PlanResult,
  PlanOptions,
  TasksResult,
  TasksOptions,
  ImplementResult,
  ImplementOptions,
  MakeResult,
  MakeOptions,
  ProgressChunk,
} from './types';
import { Validators } from './validators';
import {
  ToolkitError,
  NotFoundError,
  TimeoutError,
  PythonError,
} from './errors';

interface QueuedCommand {
  execute: () => Promise<any>;
  resolve: (value: any) => void;
  reject: (error: any) => void;
}

const DEFAULT_CONFIG: Required<ToolkitConfig> = {
  pythonPath: 'python3',
  toolkitCliPath: 'toolkit-cli',
  agents: ['claude'],
  apiKeys: {},
  timeout: 120000, // 2 minutes
  maxConcurrent: 5,
  logLevel: 'info',
};

export class ToolkitClient {
  private config: Required<ToolkitConfig>;
  private activeProcesses: number = 0;
  private queue: QueuedCommand[] = [];

  constructor(config?: ToolkitConfig) {
    this.config = { ...DEFAULT_CONFIG, ...config };
    this.validateConfig();
  }

  // ==================== Configuration ====================

  public updateConfig(updates: Partial<ToolkitConfig>): void {
    this.config = { ...this.config, ...updates };
    this.validateConfig();
  }

  public getConfig(): ToolkitConfig {
    return JSON.parse(JSON.stringify(this.config));
  }

  private validateConfig(): void {
    if (this.config.agents) {
      Validators.validateAgents(this.config.agents);
    }
    if (this.config.timeout) {
      Validators.validateTimeout(this.config.timeout);
    }
    if (this.config.maxConcurrent) {
      if (!Validators.isValidMaxConcurrent(this.config.maxConcurrent)) {
        throw new Error('maxConcurrent must be between 1 and 10');
      }
    }
  }

  // ==================== Health & Status ====================

  public async health(): Promise<HealthStatus> {
    try {
      const result = await this.executeCommand('version', ['--format', 'json']);
      const versionData = JSON.parse(result.stdout);

      return {
        healthy: true,
        toolkitCliInstalled: true,
        toolkitCliVersion: versionData.version,
        pythonVersion: versionData.python_version,
        agentsAvailable: await this.checkAgentAvailability(),
        issues: [],
      };
    } catch (error) {
      if (error instanceof NotFoundError) {
        return {
          healthy: false,
          toolkitCliInstalled: false,
          agentsAvailable: [],
          issues: [error.message],
        };
      }

      return {
        healthy: false,
        toolkitCliInstalled: false,
        agentsAvailable: [],
        issues: [(error as Error).message],
      };
    }
  }

  public async version(): Promise<VersionInfo> {
    const packageJson = require('../package.json');
    const nodePackageVersion = packageJson.version;

    try {
      const result = await this.executeCommand('version', ['--format', 'json']);
      const versionData = JSON.parse(result.stdout);

      const toolkitVersion = versionData.version;
      const minimumRequired = packageJson.peerDependencies?.['toolkit-cli'] || '1.3.0';

      const compatible = this.compareVersions(toolkitVersion, minimumRequired) >= 0;

      return {
        nodePackage: nodePackageVersion,
        toolkitCli: toolkitVersion,
        compatible,
        minimumRequired,
        error: compatible
          ? undefined
          : `toolkit-cli version ${minimumRequired}+ required, found ${toolkitVersion}`,
      };
    } catch (error) {
      throw new NotFoundError('toolkit-cli');
    }
  }

  private async checkAgentAvailability() {
    const agentEnvMap: Record<string, string> = {
      claude: 'ANTHROPIC_API_KEY',
      gemini: 'GOOGLE_API_KEY',
      qwen: 'QWEN_API_KEY',
      codex: 'OPENAI_API_KEY',
      gpt4: 'OPENAI_API_KEY',
      deepseek: 'DEEPSEEK_API_KEY',
      llama: 'TOGETHER_API_KEY',
      mistral: 'TOGETHER_API_KEY',
      rycode: 'RYCODE_API_KEY',
    };

    const statuses = [];
    for (const [agent, envVar] of Object.entries(agentEnvMap)) {
      statuses.push({
        name: agent,
        available: true,
        configured: !!process.env[envVar] || !!this.config.apiKeys?.[agent as keyof typeof this.config.apiKeys],
      });
    }

    return statuses;
  }

  private compareVersions(v1: string, v2: string): number {
    const parts1 = v1.split('.').map(Number);
    const parts2 = v2.split('.').map(Number);

    for (let i = 0; i < Math.max(parts1.length, parts2.length); i++) {
      const part1 = parts1[i] || 0;
      const part2 = parts2[i] || 0;

      if (part1 > part2) return 1;
      if (part1 < part2) return -1;
    }

    return 0;
  }

  // ==================== Queue Management ====================

  public getQueueStatus(): QueueStatus {
    return {
      active: this.activeProcesses,
      queued: this.queue.length,
    };
  }

  private async executeWithQueue<T>(fn: () => Promise<T>): Promise<T> {
    if (this.activeProcesses < this.config.maxConcurrent) {
      this.activeProcesses++;
      try {
        return await fn();
      } finally {
        this.activeProcesses--;
        this.processQueue();
      }
    }

    return new Promise<T>((resolve, reject) => {
      this.queue.push({
        execute: fn,
        resolve,
        reject,
      });
    });
  }

  private async processQueue(): Promise<void> {
    if (this.queue.length === 0 || this.activeProcesses >= this.config.maxConcurrent) {
      return;
    }

    const item = this.queue.shift();
    if (!item) return;

    this.activeProcesses++;
    try {
      const result = await item.execute();
      item.resolve(result);
    } catch (error) {
      item.reject(error);
    } finally {
      this.activeProcesses--;
      this.processQueue();
    }
  }

  // ==================== Command Execution ====================

  private async executeCommand(
    command: string,
    args: string[],
    options?: { timeout?: number; onProgress?: (chunk: ProgressChunk) => void; signal?: AbortSignal }
  ): Promise<{ stdout: string; stderr: string }> {
    return new Promise((resolve, reject) => {
      const timeout = options?.timeout || this.config.timeout;
      const proc: ChildProcess = spawn(this.config.toolkitCliPath, [command, ...args], {
        env: { ...process.env, ...this.buildEnv() },
      });

      let stdout = '';
      let stderr = '';
      let timedOut = false;
      let aborted = false;

      const timeoutId = setTimeout(() => {
        timedOut = true;
        proc.kill();
        reject(new TimeoutError(command, timeout));
      }, timeout);

      const abortHandler = () => {
        aborted = true;
        proc.kill();
        reject(new Error('Command aborted'));
      };

      options?.signal?.addEventListener('abort', abortHandler);

      proc.stdout?.on('data', (data) => {
        stdout += data.toString();

        // Handle NDJSON progress updates
        if (options?.onProgress) {
          const lines = data.toString().split('\n');
          for (const line of lines) {
            if (line.trim() && line.startsWith('{')) {
              try {
                const chunk = JSON.parse(line);
                if (chunk.phase && chunk.message !== undefined && chunk.progress !== undefined) {
                  options.onProgress(chunk);
                }
              } catch {
                // Not a progress chunk, ignore
              }
            }
          }
        }
      });

      proc.stderr?.on('data', (data) => {
        stderr += data.toString();
      });

      proc.on('error', (error) => {
        clearTimeout(timeoutId);
        options?.signal?.removeEventListener('abort', abortHandler);

        if (error.message.includes('ENOENT')) {
          reject(new NotFoundError('toolkit-cli'));
        } else {
          reject(error);
        }
      });

      proc.on('close', (code) => {
        clearTimeout(timeoutId);
        options?.signal?.removeEventListener('abort', abortHandler);

        if (timedOut || aborted) return;

        if (code === 0) {
          resolve({ stdout, stderr });
        } else {
          reject(new PythonError(`Command failed with exit code ${code}`, stderr, code || 1));
        }
      });
    });
  }

  private buildEnv(): Record<string, string> {
    const env: Record<string, string> = {};

    // Map config API keys to environment variables
    if (this.config.apiKeys?.anthropic) {
      env.ANTHROPIC_API_KEY = this.config.apiKeys.anthropic;
    }
    if (this.config.apiKeys?.openai) {
      env.OPENAI_API_KEY = this.config.apiKeys.openai;
    }
    if (this.config.apiKeys?.google) {
      env.GOOGLE_API_KEY = this.config.apiKeys.google;
    }
    if (this.config.apiKeys?.qwen) {
      env.QWEN_API_KEY = this.config.apiKeys.qwen;
    }
    if (this.config.apiKeys?.deepseek) {
      env.DEEPSEEK_API_KEY = this.config.apiKeys.deepseek;
    }
    if (this.config.apiKeys?.together) {
      env.TOGETHER_API_KEY = this.config.apiKeys.together;
    }
    if (this.config.apiKeys?.rycode) {
      env.RYCODE_API_KEY = this.config.apiKeys.rycode;
    }

    return env;
  }

  private parseResult<T>(stdout: string): CommandResult<T> {
    try {
      // Filter out progress chunks (NDJSON lines)
      const lines = stdout.split('\n');
      let resultJson = '';

      for (const line of lines) {
        if (line.trim() && line.startsWith('{')) {
          const parsed = JSON.parse(line);
          // Final result has success field
          if (parsed.success !== undefined) {
            resultJson = line;
            break;
          }
        }
      }

      if (!resultJson) {
        // Fallback: try parsing entire stdout
        return JSON.parse(stdout);
      }

      return JSON.parse(resultJson);
    } catch (error) {
      throw new Error(`Failed to parse command result: ${error}`);
    }
  }

  // ==================== Command Methods ====================

  public async oneshot(
    projectIdea: string,
    options?: OneshotOptions
  ): Promise<CommandResult<OneshotResult>> {
    Validators.validateProjectIdea(projectIdea);

    const args = [Validators.sanitizeInput(projectIdea), '--format', 'json'];

    if (options?.agents) {
      Validators.validateAgents(options.agents);
      args.push('--ai', options.agents.join(' '));
    }

    if (options?.complexity) {
      Validators.validateComplexity(options.complexity);
      args.push('--complexity', options.complexity);
    }

    if (options?.includeUx) {
      args.push('--ux');
    }

    if (options?.flags) {
      for (const flag of options.flags) {
        args.push(`--${flag}`);
      }
    }

    return this.executeWithQueue(async () => {
      const result = await this.executeCommand('oneshot', args, {
        timeout: options?.timeout,
        onProgress: options?.onProgress,
        signal: options?.signal,
      });
      return this.parseResult<OneshotResult>(result.stdout);
    });
  }

  public async specify(
    feature: string,
    options?: SpecifyOptions
  ): Promise<CommandResult<SpecifyResult>> {
    Validators.validateFeature(feature);

    const args = [Validators.sanitizeInput(feature), '--format', 'json'];

    if (options?.agents) {
      Validators.validateAgents(options.agents);
      args.push('--ai', options.agents.join(' '));
    }

    if (options?.flags) {
      for (const flag of options.flags) {
        args.push(`--${flag}`);
      }
    }

    return this.executeWithQueue(async () => {
      const result = await this.executeCommand('specify', args, {
        timeout: options?.timeout,
        signal: options?.signal,
      });
      return this.parseResult<SpecifyResult>(result.stdout);
    });
  }

  public async fix(issue: string, options?: FixOptions): Promise<CommandResult<FixResult>> {
    if (issue.length < 10 || issue.length > 5000) {
      throw new Error('Issue description must be 10-5000 characters');
    }

    const args = [Validators.sanitizeInput(issue), '--format', 'json'];

    if (options?.agents) {
      Validators.validateAgents(options.agents);
      args.push('--ai', options.agents.join(' '));
    }

    if (options?.context) {
      args.push('--context', Validators.sanitizeInput(options.context));
    }

    if (options?.autoApply) {
      args.push('--auto-apply');
    }

    return this.executeWithQueue(async () => {
      const result = await this.executeCommand('fix', args, {
        timeout: options?.timeout,
        signal: options?.signal,
      });
      return this.parseResult<FixResult>(result.stdout);
    });
  }

  public async plan(objective: string, options?: PlanOptions): Promise<CommandResult<PlanResult>> {
    const args = [Validators.sanitizeInput(objective), '--format', 'json'];

    if (options?.agents) {
      Validators.validateAgents(options.agents);
      args.push('--ai', options.agents.join(' '));
    }

    return this.executeWithQueue(async () => {
      const result = await this.executeCommand('plan', args, {
        timeout: options?.timeout,
        signal: options?.signal,
      });
      return this.parseResult<PlanResult>(result.stdout);
    });
  }

  public async tasks(input: string, options?: TasksOptions): Promise<CommandResult<TasksResult>> {
    const args = [Validators.sanitizeInput(input), '--format', 'json'];

    if (options?.agents) {
      Validators.validateAgents(options.agents);
      args.push('--ai', options.agents.join(' '));
    }

    return this.executeWithQueue(async () => {
      const result = await this.executeCommand('tasks', args, {
        timeout: options?.timeout,
        signal: options?.signal,
      });
      return this.parseResult<TasksResult>(result.stdout);
    });
  }

  public async implement(
    task: string,
    options?: ImplementOptions
  ): Promise<CommandResult<ImplementResult>> {
    const args = [Validators.sanitizeInput(task), '--format', 'json'];

    if (options?.agents) {
      Validators.validateAgents(options.agents);
      args.push('--ai', options.agents.join(' '));
    }

    return this.executeWithQueue(async () => {
      const result = await this.executeCommand('implement', args, {
        timeout: options?.timeout,
        signal: options?.signal,
      });
      return this.parseResult<ImplementResult>(result.stdout);
    });
  }

  public async make(description: string, options?: MakeOptions): Promise<CommandResult<MakeResult>> {
    const args = [Validators.sanitizeInput(description), '--format', 'json'];

    if (options?.agents) {
      Validators.validateAgents(options.agents);
      args.push('--ai', options.agents.join(' '));
    }

    return this.executeWithQueue(async () => {
      const result = await this.executeCommand('make', args, {
        timeout: options?.timeout,
        signal: options?.signal,
      });
      return this.parseResult<MakeResult>(result.stdout);
    });
  }

  // ==================== Resource Cleanup ====================

  public async close(): Promise<void> {
    // Clear queue
    for (const item of this.queue) {
      item.reject(new Error('Client closed'));
    }
    this.queue = [];

    // Note: Active processes will complete naturally
    // In a production version, we'd track process PIDs and kill them
  }
}
