/**
 * Smart Model Recommender
 *
 * Intelligently recommends the best model based on task context,
 * user preferences, cost considerations, and performance requirements.
 */

export interface TaskContext {
  task: 'code_generation' | 'code_review' | 'quick_question' | 'documentation' | 'debugging' | 'refactoring' | 'general'
  complexity?: 'simple' | 'medium' | 'complex'
  contextSize?: number // Number of tokens in context
  needsVision?: boolean
  needsRealTimeInfo?: boolean
  speedPreference?: 'fastest' | 'balanced' | 'quality'
  costPreference?: 'cheapest' | 'balanced' | 'premium'
}

export interface ModelRecommendation {
  provider: string
  model: string
  reason: string
  pros: string[]
  cons: string[]
  estimatedCost: string
  speed: 'fast' | 'medium' | 'slow'
  quality: number // 1-5 stars
  confidence: number // 0-1
}

export interface ModelCapabilities {
  provider: string
  model: string
  contextWindow: number
  supportsVision: boolean
  supportsRealTime: boolean
  speed: number // 1-10, higher is faster
  quality: number // 1-5 stars
  costEfficiency: number // 1-10, higher is more efficient
  bestFor: string[]
}

/**
 * Model capabilities database
 */
const MODEL_CAPABILITIES: ModelCapabilities[] = [
  // Anthropic
  {
    provider: 'anthropic',
    model: 'claude-3-5-sonnet-20241022',
    contextWindow: 200000,
    supportsVision: true,
    supportsRealTime: false,
    speed: 8,
    quality: 5,
    costEfficiency: 8,
    bestFor: ['code_generation', 'code_review', 'refactoring', 'complex reasoning']
  },
  {
    provider: 'anthropic',
    model: 'claude-3-5-haiku-20241022',
    contextWindow: 200000,
    supportsVision: false,
    supportsRealTime: false,
    speed: 10,
    quality: 4,
    costEfficiency: 10,
    bestFor: ['quick_question', 'documentation', 'simple tasks']
  },
  {
    provider: 'anthropic',
    model: 'claude-3-opus-20240229',
    contextWindow: 200000,
    supportsVision: true,
    supportsRealTime: false,
    speed: 5,
    quality: 5,
    costEfficiency: 4,
    bestFor: ['complex reasoning', 'creative writing', 'analysis']
  },

  // OpenAI
  {
    provider: 'openai',
    model: 'gpt-4-turbo-preview',
    contextWindow: 128000,
    supportsVision: true,
    supportsRealTime: false,
    speed: 7,
    quality: 5,
    costEfficiency: 6,
    bestFor: ['code_generation', 'creative writing', 'analysis']
  },
  {
    provider: 'openai',
    model: 'gpt-4',
    contextWindow: 8192,
    supportsVision: false,
    supportsRealTime: false,
    speed: 6,
    quality: 5,
    costEfficiency: 3,
    bestFor: ['complex reasoning', 'analysis']
  },
  {
    provider: 'openai',
    model: 'gpt-3.5-turbo',
    contextWindow: 16385,
    supportsVision: false,
    supportsRealTime: false,
    speed: 10,
    quality: 3,
    costEfficiency: 10,
    bestFor: ['quick_question', 'simple tasks']
  },

  // Grok (xAI)
  {
    provider: 'grok',
    model: 'grok-2',
    contextWindow: 128000,
    supportsVision: false,
    supportsRealTime: true,
    speed: 9,
    quality: 4,
    costEfficiency: 8,
    bestFor: ['code_generation', 'real-time info', 'humor']
  },
  {
    provider: 'grok',
    model: 'grok-2-mini',
    contextWindow: 128000,
    supportsVision: false,
    supportsRealTime: true,
    speed: 10,
    quality: 3,
    costEfficiency: 10,
    bestFor: ['quick_question', 'real-time info']
  },

  // Qwen (Alibaba)
  {
    provider: 'qwen',
    model: 'qwen-max',
    contextWindow: 8192,
    supportsVision: false,
    supportsRealTime: false,
    speed: 7,
    quality: 4,
    costEfficiency: 8,
    bestFor: ['code_generation', 'analysis']
  },
  {
    provider: 'qwen',
    model: 'qwen-plus',
    contextWindow: 32768,
    supportsVision: false,
    supportsRealTime: false,
    speed: 8,
    quality: 4,
    costEfficiency: 9,
    bestFor: ['code_generation', 'documentation']
  },
  {
    provider: 'qwen',
    model: 'qwen-turbo',
    contextWindow: 8192,
    supportsVision: false,
    supportsRealTime: false,
    speed: 10,
    quality: 3,
    costEfficiency: 10,
    bestFor: ['quick_question', 'simple tasks']
  },

  // Google
  {
    provider: 'google',
    model: 'gemini-1.5-pro',
    contextWindow: 1000000,
    supportsVision: true,
    supportsRealTime: false,
    speed: 6,
    quality: 5,
    costEfficiency: 7,
    bestFor: ['large context', 'analysis', 'vision tasks']
  },
  {
    provider: 'google',
    model: 'gemini-1.5-flash',
    contextWindow: 1000000,
    supportsVision: true,
    supportsRealTime: false,
    speed: 10,
    quality: 4,
    costEfficiency: 10,
    bestFor: ['quick_question', 'large context', 'vision tasks']
  }
]

export class ModelRecommender {
  /**
   * Get model recommendations based on context
   */
  recommend(
    context: TaskContext,
    availableModels: Array<{ provider: string; model: string }>
  ): ModelRecommendation[] {
    // Filter to available models
    const available = MODEL_CAPABILITIES.filter(cap =>
      availableModels.some(m => m.provider === cap.provider && m.model === cap.model)
    )

    // Score each model
    const scored = available.map(cap => ({
      model: cap,
      score: this.scoreModel(cap, context)
    }))

    // Sort by score (descending)
    scored.sort((a, b) => b.score - a.score)

    // Generate recommendations
    return scored.slice(0, 3).map(({ model, score }) =>
      this.generateRecommendation(model, context, score)
    )
  }

  /**
   * Score a model based on context
   */
  private scoreModel(model: ModelCapabilities, context: TaskContext): number {
    let score = 0

    // Task-specific scoring
    if (model.bestFor.includes(context.task)) {
      score += 30
    }

    // Vision requirement
    if (context.needsVision && model.supportsVision) {
      score += 20
    } else if (context.needsVision && !model.supportsVision) {
      return 0 // Disqualify if vision needed but not supported
    }

    // Real-time info requirement
    if (context.needsRealTimeInfo && model.supportsRealTime) {
      score += 20
    }

    // Context size
    if (context.contextSize) {
      if (context.contextSize <= model.contextWindow) {
        score += 10
      } else {
        return 0 // Disqualify if context too large
      }

      // Bonus for large context windows when needed
      if (context.contextSize > 100000 && model.contextWindow >= 200000) {
        score += 10
      }
    }

    // Speed preference
    if (context.speedPreference === 'fastest') {
      score += model.speed * 2
    } else if (context.speedPreference === 'balanced') {
      score += model.speed
    }

    // Quality preference
    if (context.speedPreference === 'quality') {
      score += model.quality * 5
    } else {
      score += model.quality * 2
    }

    // Cost preference
    if (context.costPreference === 'cheapest') {
      score += model.costEfficiency * 2
    } else if (context.costPreference === 'balanced') {
      score += model.costEfficiency
    }

    // Complexity matching
    if (context.complexity === 'simple' && model.quality >= 4) {
      score -= 5 // Slight penalty for overengineering
    } else if (context.complexity === 'complex' && model.quality < 4) {
      score -= 10 // Penalty for using weak model on hard task
    }

    return score
  }

  /**
   * Generate a recommendation object
   */
  private generateRecommendation(
    model: ModelCapabilities,
    context: TaskContext,
    score: number
  ): ModelRecommendation {
    const pros: string[] = []
    const cons: string[] = []

    // Determine pros
    if (model.speed >= 9) pros.push('Very fast response times')
    if (model.quality === 5) pros.push('Highest quality outputs')
    if (model.costEfficiency >= 9) pros.push('Extremely cost-efficient')
    if (model.contextWindow >= 100000) pros.push(`Large ${this.formatContextWindow(model.contextWindow)} context`)
    if (model.supportsVision) pros.push('Supports image understanding')
    if (model.supportsRealTime) pros.push('Real-time web access')

    // Determine cons
    if (model.speed < 7) cons.push('Slower response times')
    if (model.costEfficiency < 6) cons.push('Higher cost')
    if (model.contextWindow < 32000) cons.push('Limited context window')
    if (!model.supportsVision && context.needsVision) cons.push('No vision support')

    // Generate reason
    const reason = this.generateReason(model, context)

    // Estimate cost
    const estimatedCost = this.estimateCost(model, context)

    // Determine speed category
    const speed = model.speed >= 9 ? 'fast' : model.speed >= 7 ? 'medium' : 'slow'

    // Confidence based on score
    const confidence = Math.min(score / 100, 1)

    return {
      provider: model.provider,
      model: model.model,
      reason,
      pros,
      cons,
      estimatedCost,
      speed,
      quality: model.quality,
      confidence
    }
  }

  /**
   * Generate human-readable reason
   */
  private generateReason(model: ModelCapabilities, context: TaskContext): string {
    const reasons: string[] = []

    if (context.task === 'quick_question' && model.speed >= 9) {
      reasons.push('Lightning fast for quick questions')
    } else if (context.task === 'code_generation' && model.bestFor.includes('code_generation')) {
      reasons.push('Excellent for code generation')
    } else if (context.task === 'code_review' && model.quality >= 4) {
      reasons.push('High-quality code analysis')
    }

    if (model.costEfficiency >= 9 && context.costPreference === 'cheapest') {
      reasons.push('Most cost-effective option')
    }

    if (model.supportsRealTime && context.needsRealTimeInfo) {
      reasons.push('Provides real-time information')
    }

    if (reasons.length === 0) {
      reasons.push(`Good balance of ${model.quality === 5 ? 'quality' : 'performance'} and cost`)
    }

    return reasons.join(', ')
  }

  /**
   * Estimate cost for typical usage
   */
  private estimateCost(model: ModelCapabilities, _context: TaskContext): string {
    // Rough estimates per request
    const estimates: Record<string, string> = {
      'anthropic/claude-3-5-sonnet-20241022': '$0.01-0.05 per request',
      'anthropic/claude-3-5-haiku-20241022': '$0.001-0.01 per request',
      'anthropic/claude-3-opus-20240229': '$0.05-0.20 per request',
      'openai/gpt-4-turbo-preview': '$0.02-0.10 per request',
      'openai/gpt-4': '$0.05-0.20 per request',
      'openai/gpt-3.5-turbo': '$0.001-0.005 per request',
      'grok/grok-2': '$0.01-0.03 per request',
      'grok/grok-2-mini': '$0.001-0.01 per request',
      'qwen/qwen-max': '$0.01-0.03 per request',
      'qwen/qwen-plus': '$0.002-0.01 per request',
      'qwen/qwen-turbo': '$0.0005-0.002 per request',
      'google/gemini-1.5-pro': '$0.01-0.05 per request',
      'google/gemini-1.5-flash': '$0.0005-0.002 per request'
    }

    const key = `${model.provider}/${model.model}`
    return estimates[key] || '$0.01-0.05 per request'
  }

  /**
   * Format context window for display
   */
  private formatContextWindow(tokens: number): string {
    if (tokens >= 1000000) {
      return `${(tokens / 1000000).toFixed(1)}M token`
    }
    if (tokens >= 1000) {
      return `${(tokens / 1000).toFixed(0)}K token`
    }
    return `${tokens} token`
  }

  /**
   * Get default model recommendation
   */
  getDefaultRecommendation(): ModelRecommendation {
    return {
      provider: 'anthropic',
      model: 'claude-3-5-sonnet-20241022',
      reason: 'Best all-around model for coding tasks',
      pros: [
        'Excellent code generation',
        'Large 200K context window',
        'Good balance of speed and quality',
        'Cost-efficient'
      ],
      cons: [
        'Premium features require paid tier'
      ],
      estimatedCost: '$0.01-0.05 per request',
      speed: 'fast',
      quality: 5,
      confidence: 0.95
    }
  }

  /**
   * Compare models side-by-side
   */
  compareModels(models: Array<{ provider: string; model: string }>): Array<{
    provider: string
    model: string
    speed: number
    quality: number
    costEfficiency: number
    contextWindow: string
    capabilities: string[]
  }> {
    return models
      .map(({ provider, model }) => {
        const cap = MODEL_CAPABILITIES.find(
          c => c.provider === provider && c.model === model
        )

        if (!cap) return null

        const capabilities: string[] = []
        if (cap.supportsVision) capabilities.push('Vision')
        if (cap.supportsRealTime) capabilities.push('Real-time')
        capabilities.push(...cap.bestFor)

        return {
          provider: cap.provider,
          model: cap.model,
          speed: cap.speed,
          quality: cap.quality,
          costEfficiency: cap.costEfficiency,
          contextWindow: this.formatContextWindow(cap.contextWindow),
          capabilities
        }
      })
      .filter((c): c is NonNullable<typeof c> => c !== null)
  }
}

/**
 * Global model recommender instance
 */
export const modelRecommender = new ModelRecommender()
