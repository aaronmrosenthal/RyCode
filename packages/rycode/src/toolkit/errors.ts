/**
 * Error classes for toolkit-cli client
 */

import { ErrorCode } from './types';

export class ToolkitError extends Error {
  public readonly code: ErrorCode;

  constructor(message: string, code: ErrorCode = 'UNKNOWN_ERROR') {
    super(message);
    this.name = this.constructor.name;
    this.code = code;
    Error.captureStackTrace(this, this.constructor);
  }
}

export class ValidationError extends ToolkitError {
  public readonly details: {
    field: string;
    value: any;
    constraint: string;
  };

  constructor(field: string, value: any, constraint: string) {
    super(
      `Validation failed for ${field}: ${constraint}`,
      'VALIDATION_ERROR'
    );
    this.details = { field, value, constraint };
  }
}

export class ApiError extends ToolkitError {
  public readonly provider: string;
  public readonly statusCode?: number;
  public readonly retryable: boolean;

  constructor(
    message: string,
    provider: string,
    statusCode?: number,
    retryable: boolean = false
  ) {
    super(message, 'API_ERROR');
    this.provider = provider;
    this.statusCode = statusCode;
    this.retryable = retryable;
  }
}

export class TimeoutError extends ToolkitError {
  public readonly command: string;
  public readonly timeout: number;

  constructor(command: string, timeout: number) {
    super(
      `Command '${command}' exceeded timeout of ${timeout}ms`,
      'TIMEOUT_ERROR'
    );
    this.command = command;
    this.timeout = timeout;
  }
}

export class RateLimitError extends ToolkitError {
  public readonly retryAfter: number;

  constructor(retryAfter: number) {
    super(
      `Rate limit exceeded. Retry after ${retryAfter} seconds`,
      'RATE_LIMIT_ERROR'
    );
    this.retryAfter = retryAfter;
  }
}

export class NotFoundError extends ToolkitError {
  public readonly missing: 'python' | 'toolkit-cli';
  public readonly installInstructions: string;

  constructor(missing: 'python' | 'toolkit-cli') {
    const instructions =
      missing === 'python'
        ? 'Install Python 3.11+ from https://python.org'
        : 'Install toolkit-cli with: pip install toolkit-cli';

    super(`${missing} not found. ${instructions}`, 'NOT_FOUND_ERROR');
    this.missing = missing;
    this.installInstructions = instructions;
  }
}

export class PythonError extends ToolkitError {
  public readonly stderr: string;
  public readonly exitCode: number;

  constructor(message: string, stderr: string, exitCode: number) {
    super(message, 'PYTHON_ERROR');
    this.stderr = stderr;
    this.exitCode = exitCode;
  }
}
