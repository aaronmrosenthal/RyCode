/**
 * Ghost Text Prediction System
 *
 * Shows dimmed predictions of what the AI might suggest next based on:
 * - Conversation context
 * - Recent patterns
 * - User behavior
 * - Code structure
 */

import { UI } from "../../ui"

export namespace GhostText {
  export interface Suggestion {
    text: string
    confidence: number // 0-1
    trigger: 'tab' | 'enter' | 'auto'
    source: 'pattern' | 'history' | 'ml'
    context: {
      lastMessages: string[]
      currentFile?: string
      recentEdits?: string[]
    }
  }

  export interface PredictionEngine {
    predict(input: string, context: any): Promise<Suggestion | null>
    learn(accepted: boolean, suggestion: Suggestion): void
    getConfidence(pattern: string): number
  }

  /**
   * Pattern-based prediction engine
   * Learns from user's accepted suggestions
   */
  export class PatternPredictor implements PredictionEngine {
    private patterns: Map<string, { count: number; accepted: number }> = new Map()
    private readonly minConfidence = 0.6

    async predict(input: string, context: any): Promise<Suggestion | null> {
      // Look for matching patterns
      const matches = this.findPatterns(input)
      if (matches.length === 0) return null

      const bestMatch = matches[0]
      const confidence = this.getConfidence(bestMatch.pattern)

      if (confidence < this.minConfidence) return null

      return {
        text: bestMatch.completion,
        confidence,
        trigger: confidence > 0.8 ? 'auto' : 'tab',
        source: 'pattern',
        context: {
          lastMessages: context.lastMessages || [],
          currentFile: context.currentFile,
        },
      }
    }

    learn(accepted: boolean, suggestion: Suggestion): void {
      const key = this.extractPattern(suggestion.text)
      const existing = this.patterns.get(key) || { count: 0, accepted: 0 }

      existing.count++
      if (accepted) existing.accepted++

      this.patterns.set(key, existing)
    }

    getConfidence(pattern: string): number {
      const stats = this.patterns.get(pattern)
      if (!stats || stats.count < 3) return 0

      return stats.accepted / stats.count
    }

    private findPatterns(_input: string): Array<{ pattern: string; completion: string }> {
      // Implementation: fuzzy match against known patterns
      return []
    }

    private extractPattern(text: string): string {
      // Implementation: extract key pattern from text
      return text.slice(0, 50) // Simplified
    }
  }

  /**
   * Render ghost text in terminal
   */
  export function renderGhost(text: string, _position: number): string {
    return UI.Style.TEXT_DIM + text + UI.Style.RESET
  }

  /**
   * Display ghost suggestion with acceptance hints
   */
  export function displaySuggestion(suggestion: Suggestion): string {
    const confidenceColor =
      suggestion.confidence > 0.8 ? UI.Style.MATRIX_GREEN :
      suggestion.confidence > 0.6 ? UI.Style.CLAUDE_BLUE :
      UI.Style.TEXT_DIM

    const triggerHint =
      suggestion.trigger === 'auto' ? '↵' :
      suggestion.trigger === 'tab' ? '⇥' : '⏎'

    return `${confidenceColor}${suggestion.text}${UI.Style.RESET} ${UI.Style.TEXT_DIM}${triggerHint}${UI.Style.RESET}`
  }

  /**
   * Context-aware command suggestions
   */
  export const CommandSuggestions = {
    afterError: [
      { text: '/debug', confidence: 0.9 },
      { text: '/fix', confidence: 0.85 },
      { text: '/explain', confidence: 0.7 },
    ],
    afterSuccess: [
      { text: '/commit', confidence: 0.8 },
      { text: '/test', confidence: 0.75 },
      { text: '/review', confidence: 0.7 },
    ],
    editingReact: [
      { text: '/preview', confidence: 0.85 },
      { text: '/test components', confidence: 0.8 },
      { text: '/lint --fix', confidence: 0.75 },
    ],
    editingTests: [
      { text: '/test', confidence: 0.9 },
      { text: '/coverage', confidence: 0.7 },
    ],
  }

  /**
   * Smart clipboard with AI categorization
   */
  export interface SmartClipboard {
    items: ClipboardItem[]
    categories: Map<string, ClipboardItem[]>

    add(text: string, context?: any): void
    search(query: string): ClipboardItem[]
    categorize(item: ClipboardItem): string
  }

  export interface ClipboardItem {
    text: string
    timestamp: number
    category: 'function' | 'config' | 'example' | 'snippet' | 'error' | 'other'
    language?: string
    context?: {
      file?: string
      conversation?: string
    }
  }

  export class AIClipboard implements SmartClipboard {
    items: ClipboardItem[] = []
    categories: Map<string, ClipboardItem[]> = new Map()
    private maxItems = 100

    add(text: string, context?: any): void {
      const category = this.detectCategory(text)
      const language = this.detectLanguage(text)

      const item: ClipboardItem = {
        text,
        timestamp: Date.now(),
        category,
        language,
        context,
      }

      this.items.unshift(item)
      if (this.items.length > this.maxItems) {
        this.items.pop()
      }

      // Update categories
      const categoryItems = this.categories.get(category) || []
      categoryItems.unshift(item)
      this.categories.set(category, categoryItems)
    }

    search(query: string): ClipboardItem[] {
      const lowerQuery = query.toLowerCase()
      return this.items.filter(item =>
        item.text.toLowerCase().includes(lowerQuery) ||
        item.category.includes(lowerQuery) ||
        item.language?.includes(lowerQuery)
      )
    }

    categorize(item: ClipboardItem): string {
      return this.detectCategory(item.text)
    }

    private detectCategory(text: string): ClipboardItem['category'] {
      // Simple heuristics - could be enhanced with ML
      if (text.includes('function') || text.includes('=>')) return 'function'
      if (text.includes('{') && text.includes('}') && !text.includes('(')) return 'config'
      if (text.includes('Error:') || text.includes('Exception')) return 'error'
      if (text.includes('//') || text.includes('/*')) return 'example'
      return 'snippet'
    }

    private detectLanguage(text: string): string | undefined {
      // Simple detection - could use tree-sitter or other tools
      if (text.includes('function') || text.includes('=>')) return 'typescript'
      if (text.includes('def ') || text.includes('import ')) return 'python'
      if (text.includes('package ') || text.includes('func ')) return 'go'
      return undefined
    }
  }

  /**
   * Conversation state tracker for better predictions
   */
  export class ConversationTracker {
    private state: {
      lastIntent?: 'debug' | 'implement' | 'refactor' | 'learn' | 'fix'
      currentFiles: Set<string>
      recentPatterns: string[]
      errorCount: number
      successCount: number
    } = {
      currentFiles: new Set(),
      recentPatterns: [],
      errorCount: 0,
      successCount: 0,
    }

    updateIntent(message: string): void {
      // Detect intent from message
      if (message.includes('bug') || message.includes('error')) {
        this.state.lastIntent = 'fix'
        this.state.errorCount++
      } else if (message.includes('add') || message.includes('create')) {
        this.state.lastIntent = 'implement'
      } else if (message.includes('refactor') || message.includes('improve')) {
        this.state.lastIntent = 'refactor'
      } else if (message.includes('explain') || message.includes('how')) {
        this.state.lastIntent = 'learn'
      } else if (message.includes('debug') || message.includes('why')) {
        this.state.lastIntent = 'debug'
      }
    }

    trackSuccess(): void {
      this.state.successCount++
      this.state.errorCount = 0 // Reset error streak
    }

    getSuggestedCommands(): string[] {
      const { lastIntent, errorCount, successCount } = this.state

      // Frustrated user pattern
      if (errorCount > 2) {
        return ['/help', '/explain', '/debug']
      }

      // Success streak
      if (successCount > 3) {
        return ['/commit', '/test', '/review']
      }

      // Intent-based suggestions
      switch (lastIntent) {
        case 'fix':
          return ['/test', '/debug', '/lint']
        case 'implement':
          return ['/test', '/preview', '/lint']
        case 'refactor':
          return ['/test', '/review', '/commit']
        case 'learn':
          return ['/docs', '/example', '/explain']
        case 'debug':
          return ['/logs', '/trace', '/fix']
        default:
          return ['/help', '/models', '/stats']
      }
    }
  }
}
