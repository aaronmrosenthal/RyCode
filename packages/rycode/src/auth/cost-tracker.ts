/**
 * Cost Tracking Dashboard
 *
 * Real-time cost tracking, projections, and smart recommendations
 * to help users understand and optimize their AI spending.
 */

export interface ModelPricing {
  provider: string
  model: string
  inputTokenPrice: number // Price per 1K tokens
  outputTokenPrice: number // Price per 1K tokens
  contextWindow: number
}

export interface UsageRecord {
  timestamp: Date
  provider: string
  model: string
  inputTokens: number
  outputTokens: number
  cost: number
}

export interface CostSummary {
  today: number
  yesterday: number
  thisWeek: number
  thisMonth: number
  lastMonth: number
  projection: {
    dailyAverage: number
    monthlyProjection: number
    yearlyProjection: number
  }
}

export interface CostBreakdown {
  byProvider: Record<string, number>
  byModel: Record<string, number>
  byDay: Array<{ date: string; cost: number }>
}

export interface CostSavingTip {
  type: 'model_suggestion' | 'usage_pattern' | 'provider_choice'
  message: string
  potentialSaving: number
  action?: string
}

/**
 * Model pricing database (per 1K tokens)
 */
const MODEL_PRICING: ModelPricing[] = [
  // Anthropic
  { provider: 'anthropic', model: 'claude-3-5-sonnet-20241022', inputTokenPrice: 0.003, outputTokenPrice: 0.015, contextWindow: 200000 },
  { provider: 'anthropic', model: 'claude-3-5-haiku-20241022', inputTokenPrice: 0.001, outputTokenPrice: 0.005, contextWindow: 200000 },
  { provider: 'anthropic', model: 'claude-3-opus-20240229', inputTokenPrice: 0.015, outputTokenPrice: 0.075, contextWindow: 200000 },

  // OpenAI
  { provider: 'openai', model: 'gpt-4-turbo-preview', inputTokenPrice: 0.01, outputTokenPrice: 0.03, contextWindow: 128000 },
  { provider: 'openai', model: 'gpt-4', inputTokenPrice: 0.03, outputTokenPrice: 0.06, contextWindow: 8192 },
  { provider: 'openai', model: 'gpt-3.5-turbo', inputTokenPrice: 0.0005, outputTokenPrice: 0.0015, contextWindow: 16385 },

  // Grok (xAI)
  { provider: 'grok', model: 'grok-2', inputTokenPrice: 0.002, outputTokenPrice: 0.01, contextWindow: 128000 },
  { provider: 'grok', model: 'grok-2-mini', inputTokenPrice: 0.0005, outputTokenPrice: 0.002, contextWindow: 128000 },

  // Qwen (Alibaba)
  { provider: 'qwen', model: 'qwen-turbo', inputTokenPrice: 0.0002, outputTokenPrice: 0.0006, contextWindow: 8192 },
  { provider: 'qwen', model: 'qwen-plus', inputTokenPrice: 0.0004, outputTokenPrice: 0.0012, contextWindow: 32768 },
  { provider: 'qwen', model: 'qwen-max', inputTokenPrice: 0.002, outputTokenPrice: 0.006, contextWindow: 8192 },

  // Google
  { provider: 'google', model: 'gemini-1.5-pro', inputTokenPrice: 0.00125, outputTokenPrice: 0.005, contextWindow: 1000000 },
  { provider: 'google', model: 'gemini-1.5-flash', inputTokenPrice: 0.000075, outputTokenPrice: 0.0003, contextWindow: 1000000 }
]

export class CostTracker {
  private usage: UsageRecord[] = []
  private readonly storageKey = 'rycode_usage_history'

  constructor() {
    this.loadFromStorage()
  }

  /**
   * Record a usage event
   */
  recordUsage(
    provider: string,
    model: string,
    inputTokens: number,
    outputTokens: number
  ): UsageRecord {
    const pricing = this.getPricing(provider, model)
    const cost = this.calculateCost(inputTokens, outputTokens, pricing)

    const record: UsageRecord = {
      timestamp: new Date(),
      provider,
      model,
      inputTokens,
      outputTokens,
      cost
    }

    this.usage.push(record)
    this.saveToStorage()

    return record
  }

  /**
   * Calculate cost for a usage event
   */
  calculateCost(
    inputTokens: number,
    outputTokens: number,
    pricing: ModelPricing
  ): number {
    const inputCost = (inputTokens / 1000) * pricing.inputTokenPrice
    const outputCost = (outputTokens / 1000) * pricing.outputTokenPrice
    return inputCost + outputCost
  }

  /**
   * Get pricing for a model
   */
  getPricing(provider: string, model: string): ModelPricing {
    const pricing = MODEL_PRICING.find(
      p => p.provider === provider && p.model === model
    )

    // Return default pricing if not found
    if (!pricing) {
      return {
        provider,
        model,
        inputTokenPrice: 0.001,
        outputTokenPrice: 0.002,
        contextWindow: 8192
      }
    }

    return pricing
  }

  /**
   * Get cost summary
   */
  getCostSummary(): CostSummary {
    const now = new Date()
    const today = this.getStartOfDay(now)
    const yesterday = new Date(today.getTime() - 24 * 60 * 60 * 1000)
    const weekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)
    const monthStart = new Date(now.getFullYear(), now.getMonth(), 1)
    const lastMonthStart = new Date(now.getFullYear(), now.getMonth() - 1, 1)
    const lastMonthEnd = new Date(now.getFullYear(), now.getMonth(), 0)

    const todayCost = this.getCostInRange(today, now)
    const yesterdayCost = this.getCostInRange(yesterday, today)
    const weekCost = this.getCostInRange(weekAgo, now)
    const monthCost = this.getCostInRange(monthStart, now)
    const lastMonthCost = this.getCostInRange(lastMonthStart, lastMonthEnd)

    // Calculate projection
    const daysThisMonth = now.getDate()
    const dailyAverage = monthCost / daysThisMonth
    const daysInMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0).getDate()
    const monthlyProjection = dailyAverage * daysInMonth

    return {
      today: todayCost,
      yesterday: yesterdayCost,
      thisWeek: weekCost,
      thisMonth: monthCost,
      lastMonth: lastMonthCost,
      projection: {
        dailyAverage,
        monthlyProjection,
        yearlyProjection: monthlyProjection * 12
      }
    }
  }

  /**
   * Get cost breakdown
   */
  getCostBreakdown(days: number = 30): CostBreakdown {
    const since = new Date(Date.now() - days * 24 * 60 * 60 * 1000)
    const relevantUsage = this.usage.filter(u => u.timestamp >= since)

    // By provider
    const byProvider: Record<string, number> = {}
    for (const record of relevantUsage) {
      byProvider[record.provider] = (byProvider[record.provider] || 0) + record.cost
    }

    // By model
    const byModel: Record<string, number> = {}
    for (const record of relevantUsage) {
      const key = `${record.provider}/${record.model}`
      byModel[key] = (byModel[key] || 0) + record.cost
    }

    // By day
    const byDay: Array<{ date: string; cost: number }> = []
    const dayMap = new Map<string, number>()

    for (const record of relevantUsage) {
      const dateStr = record.timestamp.toISOString().split('T')[0]
      dayMap.set(dateStr, (dayMap.get(dateStr) || 0) + record.cost)
    }

    for (const [date, cost] of dayMap.entries()) {
      byDay.push({ date, cost })
    }

    byDay.sort((a, b) => a.date.localeCompare(b.date))

    return { byProvider, byModel, byDay }
  }

  /**
   * Get smart cost-saving tips
   */
  getCostSavingTips(): CostSavingTip[] {
    const tips: CostSavingTip[] = []
    const last7Days = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)
    const recentUsage = this.usage.filter(u => u.timestamp >= last7Days)

    // Check if using expensive models for simple tasks
    const expensiveModels = ['gpt-4', 'claude-3-opus-20240229']
    const expensiveUsage = recentUsage.filter(u =>
      expensiveModels.some(m => u.model.includes(m))
    )

    if (expensiveUsage.length > 10) {
      const expensiveCost = expensiveUsage.reduce((sum, u) => sum + u.cost, 0)
      const potentialSaving = expensiveCost * 0.7 // Could save ~70% with cheaper models

      tips.push({
        type: 'model_suggestion',
        message: 'You\'re using premium models frequently. Consider Claude Haiku or GPT-3.5 for simpler tasks.',
        potentialSaving,
        action: 'Switch to cost-efficient models for quick questions'
      })
    }

    // Check for inefficient provider usage
    const breakdown = this.getCostBreakdown(7)
    const providers = Object.keys(breakdown.byProvider)

    if (providers.length === 1 && providers[0] === 'openai') {
      const openaiCost = breakdown.byProvider['openai']
      const potentialSaving = openaiCost * 0.5

      tips.push({
        type: 'provider_choice',
        message: 'Claude (Anthropic) offers similar quality at lower cost for many tasks.',
        potentialSaving,
        action: 'Try Claude 3.5 Haiku as your default model'
      })
    }

    // Check for high-volume usage
    const tokensUsed = recentUsage.reduce((sum, u) => sum + u.inputTokens + u.outputTokens, 0)
    if (tokensUsed > 1000000) { // 1M tokens in 7 days
      tips.push({
        type: 'usage_pattern',
        message: 'High token usage detected. Consider using more specific prompts to reduce context.',
        potentialSaving: 0,
        action: 'Optimize your prompts and reduce context size'
      })
    }

    return tips
  }

  /**
   * Format cost for display
   */
  formatCost(cost: number): string {
    if (cost < 0.01) {
      return '< $0.01'
    }
    return `$${cost.toFixed(2)}`
  }

  /**
   * Get status bar display
   */
  getStatusBarDisplay(currentModel: { provider: string; name: string }): string {
    const summary = this.getCostSummary()
    const todayCost = this.formatCost(summary.today)

    return `${currentModel.name} | ⚡ Fast | 💰 ${todayCost} today | [tab→]`
  }

  /**
   * Helper: Get cost in date range
   */
  private getCostInRange(start: Date, end: Date): number {
    return this.usage
      .filter(u => u.timestamp >= start && u.timestamp < end)
      .reduce((sum, u) => sum + u.cost, 0)
  }

  /**
   * Helper: Get start of day
   */
  private getStartOfDay(date: Date): Date {
    return new Date(date.getFullYear(), date.getMonth(), date.getDate())
  }

  /**
   * Save usage to storage
   */
  private saveToStorage(): void {
    try {
      // Keep only last 90 days
      const ninetyDaysAgo = new Date(Date.now() - 90 * 24 * 60 * 60 * 1000)
      this.usage = this.usage.filter(u => u.timestamp >= ninetyDaysAgo)

      // In browser, use localStorage; in Node, would use a file
      if (typeof localStorage !== 'undefined') {
        localStorage.setItem(this.storageKey, JSON.stringify(this.usage))
      }
    } catch (error) {
      console.warn('Failed to save usage history:', error)
    }
  }

  /**
   * Load usage from storage
   */
  private loadFromStorage(): void {
    try {
      if (typeof localStorage !== 'undefined') {
        const stored = localStorage.getItem(this.storageKey)
        if (stored) {
          const parsed = JSON.parse(stored)
          // Restore Date objects
          this.usage = parsed.map((u: any) => ({
            ...u,
            timestamp: new Date(u.timestamp)
          }))
        }
      }
    } catch (error) {
      console.warn('Failed to load usage history:', error)
      this.usage = []
    }
  }

  /**
   * Export usage data
   */
  exportData(): {
    summary: CostSummary
    breakdown: CostBreakdown
    usage: UsageRecord[]
  } {
    return {
      summary: this.getCostSummary(),
      breakdown: this.getCostBreakdown(30),
      usage: [...this.usage]
    }
  }

  /**
   * Clear all usage data
   */
  clearData(): void {
    this.usage = []
    this.saveToStorage()
  }
}

/**
 * Global cost tracker instance
 */
export const costTracker = new CostTracker()
