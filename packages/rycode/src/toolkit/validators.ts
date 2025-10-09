/**
 * Input validation utilities
 */

import type { AgentType, ComplexityLevel } from './types';
import { ValidationError } from './errors';

const VALID_AGENTS: AgentType[] = [
  'claude',
  'gemini',
  'qwen',
  'codex',
  'gpt4',
  'deepseek',
  'llama',
  'mistral',
  'rycode'
];
const VALID_COMPLEXITY: ComplexityLevel[] = ['low', 'medium', 'high', 'enterprise'];

export class Validators {
  static isValidAgent(agent: string): boolean {
    return VALID_AGENTS.includes(agent as AgentType);
  }

  static isValidComplexity(level: string): boolean {
    return VALID_COMPLEXITY.includes(level as ComplexityLevel);
  }

  static isValidTimeout(timeout: number): boolean {
    return timeout > 0 && timeout <= 600000; // Max 10 minutes
  }

  static isValidMaxConcurrent(max: number): boolean {
    return max >= 1 && max <= 10;
  }

  static sanitizeInput(input: string): string {
    // Remove shell metacharacters to prevent command injection
    return input.replace(/[;&|`$()\n]/g, '');
  }

  static validateProjectIdea(idea: string): void {
    if (idea.length < 10 || idea.length > 5000) {
      throw new ValidationError(
        'project_idea',
        idea.length,
        '10-5000 characters required'
      );
    }
  }

  static validateFeature(feature: string): void {
    if (feature.length < 10 || feature.length > 5000) {
      throw new ValidationError(
        'feature',
        feature.length,
        '10-5000 characters required'
      );
    }
  }

  static validateAgents(agents: string[]): void {
    if (agents.length === 0) {
      throw new ValidationError('agents', agents, 'At least one agent required');
    }

    const invalid = agents.filter((a) => !this.isValidAgent(a));
    if (invalid.length > 0) {
      throw new ValidationError(
        'agents',
        invalid,
        `Invalid agents: ${invalid.join(', ')}. Valid: ${VALID_AGENTS.join(', ')}`
      );
    }
  }

  static validateComplexity(complexity: string): void {
    if (!this.isValidComplexity(complexity)) {
      throw new ValidationError(
        'complexity',
        complexity,
        `Must be one of: ${VALID_COMPLEXITY.join(', ')}`
      );
    }
  }

  static validateTimeout(timeout: number): void {
    if (!this.isValidTimeout(timeout)) {
      throw new ValidationError(
        'timeout',
        timeout,
        'Must be > 0 and <= 600000ms (10 minutes)'
      );
    }
  }
}
