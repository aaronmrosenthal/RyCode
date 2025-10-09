/**
 * Example: Using Bundled Toolkit Client in RyCode
 *
 * This demonstrates how RyCode can use the toolkit-cli client
 * for AI-powered development features.
 */

import { ToolkitClient } from './index'

/**
 * Example 1: Generate project spec from RyCode TUI
 */
export async function generateProjectSpec(idea: string) {
  const toolkit = new ToolkitClient({
    agents: ['claude', 'rycode'],
    timeout: 120000,
  })

  try {
    const result = await toolkit.oneshot(idea, {
      agents: ['claude', 'rycode'],
      complexity: 'medium',
      includeUx: true,
      onProgress: (chunk) => {
        // Update RyCode TUI with progress
        console.log(`[${chunk.progress}%] ${chunk.message}`)
      },
    })

    if (result.success) {
      return {
        specification: result.data?.specification,
        architecture: result.data?.architecture,
        uxDesigns: result.data?.uxDesigns,
      }
    } else {
      throw new Error(result.error?.message || 'Generation failed')
    }
  } finally {
    await toolkit.close()
  }
}

/**
 * Example 2: Fix code issues
 */
export async function fixCodeIssue(issue: string, context?: string) {
  const toolkit = new ToolkitClient({
    agents: ['claude'],
  })

  try {
    const result = await toolkit.fix(issue, {
      context,
      autoApply: false, // Let user review changes
    })

    if (result.success) {
      return {
        rootCause: result.data?.rootCause,
        solution: result.data?.solution,
        codeChanges: result.data?.solution.codeChanges,
      }
    }
  } finally {
    await toolkit.close()
  }
}

/**
 * Example 3: Multi-agent analysis
 */
export async function analyzeWithMultipleAgents(feature: string) {
  const toolkit = new ToolkitClient({
    agents: ['claude', 'gemini', 'qwen'],
    maxConcurrent: 3,
  })

  try {
    const result = await toolkit.specify(feature, {
      agents: ['claude', 'gemini', 'qwen'],
      flags: ['deep', 'validation'],
    })

    if (result.success) {
      return {
        specification: result.data?.specification,
        insights: result.data?.multiAgentInsights,
      }
    }
  } finally {
    await toolkit.close()
  }
}

/**
 * Example 4: Check toolkit availability
 */
export async function checkToolkitStatus() {
  const toolkit = new ToolkitClient()

  const health = await toolkit.health()

  await toolkit.close()

  return {
    healthy: health.healthy,
    toolkitVersion: health.toolkitCliVersion,
    pythonVersion: health.pythonVersion,
    agents: health.agentsAvailable.map((a) => ({
      name: a.name,
      configured: a.configured,
    })),
    issues: health.issues,
  }
}

/**
 * Example 5: RyCode AI Command Handler
 *
 * This shows how RyCode's TUI can integrate toolkit commands
 */
export class RyCodeToolkitHandler {
  private toolkit: ToolkitClient

  constructor() {
    this.toolkit = new ToolkitClient({
      agents: ['claude', 'rycode'],
      maxConcurrent: 3,
      logLevel: 'info',
    })
  }

  /**
   * Handle /oneshot command from RyCode TUI
   */
  async handleOneshotCommand(
    idea: string,
    onProgress?: (message: string) => void
  ) {
    const result = await this.toolkit.oneshot(idea, {
      agents: ['claude', 'rycode'],
      complexity: 'medium',
      includeUx: true,
      onProgress: (chunk) => {
        onProgress?.(`[${chunk.progress}%] ${chunk.message}`)
      },
    })

    if (result.success) {
      return {
        success: true,
        data: result.data,
        executionTime: result.metadata.executionTime,
      }
    } else {
      return {
        success: false,
        error: result.error?.message,
      }
    }
  }

  /**
   * Handle /fix command from RyCode TUI
   */
  async handleFixCommand(issue: string, context?: string) {
    return await this.toolkit.fix(issue, {
      context,
      agents: ['claude'],
      autoApply: false,
    })
  }

  /**
   * Handle /specify command from RyCode TUI
   */
  async handleSpecifyCommand(feature: string) {
    return await this.toolkit.specify(feature, {
      agents: ['claude', 'gemini'],
      flags: ['deep', 'validation'],
    })
  }

  /**
   * Cleanup
   */
  async close() {
    await this.toolkit.close()
  }
}

/**
 * Example 6: Integration with RyCode session
 *
 * This shows how to integrate with RyCode's session management
 */
export async function integrateWithRyCodeSession(sessionId: string) {
  const toolkit = new ToolkitClient({
    agents: ['claude', 'rycode'],
  })

  // Get session context
  // const session = await getSession(sessionId)

  // Use toolkit with session context
  const result = await toolkit.oneshot('Build feature based on session', {
    // context: session.currentFile,
    agents: ['claude', 'rycode'],
  })

  await toolkit.close()
  return result
}

// Export utilities
export { ToolkitClient }
export * from './types'
export * from './errors'
