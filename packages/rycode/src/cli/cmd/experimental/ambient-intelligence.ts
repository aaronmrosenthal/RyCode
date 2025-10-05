/**
 * Ambient Intelligence Features
 *
 * Invisible helpers that make the development experience magical:
 * - Mood detection and adaptive UI
 * - Background task orchestration
 * - Contextual command palette
 * - Celebration moments
 */

import { UI } from "../../ui"

export namespace AmbientIntelligence {
  /**
   * Mood Detection System
   */
  export class MoodDetector {
    private recentActions: Action[] = []
    private readonly windowSize = 10
    private currentMood: Mood = 'neutral'

    addAction(action: Action): void {
      this.recentActions.push(action)
      if (this.recentActions.length > this.windowSize) {
        this.recentActions.shift()
      }

      this.currentMood = this.detectMood()
    }

    private detectMood(): Mood {
      const actions = this.recentActions

      // Frustrated: repeated errors, same commands, long pauses
      const errorRate = actions.filter(a => a.type === 'error').length / actions.length
      const repeatRate = this.calculateRepeatRate()
      const avgPause = this.calculateAveragePause()

      if (errorRate > 0.5 || repeatRate > 0.6 || avgPause > 30000) {
        return 'frustrated'
      }

      // Productive: steady progress, variety of commands, quick responses
      const commandVariety = new Set(actions.map(a => a.command)).size
      const quickResponses = actions.filter(a => a.duration && a.duration < 5000).length

      if (commandVariety > 5 && quickResponses > actions.length * 0.6) {
        return 'productive'
      }

      // Exploring: many different files, read-heavy, documentation access
      const fileChanges = new Set(actions.map(a => a.file)).size
      const readActions = actions.filter(a => a.type === 'read').length

      if (fileChanges > 5 && readActions > actions.length * 0.5) {
        return 'exploring'
      }

      // Debugging: same file, errors, test runs
      const sameFile = actions.filter(a => a.file === actions[0]?.file).length
      const testRuns = actions.filter(a => a.command?.includes('test')).length

      if (sameFile > actions.length * 0.7 && (errorRate > 0.3 || testRuns > 2)) {
        return 'debugging'
      }

      // Celebrating: tests passed, build succeeded, commits
      const successActions = actions.filter(a => a.type === 'success').length
      if (successActions > 3) {
        return 'celebrating'
      }

      return 'neutral'
    }

    private calculateRepeatRate(): number {
      const commands = this.recentActions.map(a => a.command)
      const unique = new Set(commands)
      return 1 - (unique.size / commands.length)
    }

    private calculateAveragePause(): number {
      let totalPause = 0
      for (let i = 1; i < this.recentActions.length; i++) {
        const pause = this.recentActions[i].timestamp - this.recentActions[i - 1].timestamp
        totalPause += pause
      }
      return totalPause / (this.recentActions.length - 1)
    }

    getMood(): Mood {
      return this.currentMood
    }

    getSuggestions(): string[] {
      switch (this.currentMood) {
        case 'frustrated':
          return [
            '💡 Take a step back? Try /explain to understand the issue',
            '🔍 Debug mode might help: /debug',
            '📚 Check the docs: /docs',
            '🤝 Ask for help: describe what you\'re trying to do',
          ]
        case 'debugging':
          return [
            '🔎 Add logging to trace the issue',
            '🧪 Try a minimal reproduction',
            '📊 Check the stack trace carefully',
            '⏮️ Undo recent changes: /undo',
          ]
        case 'celebrating':
          return [
            '🎉 Great work! Time to commit?',
            '✅ Run the full test suite to be sure',
            '📝 Document what you learned',
            '🚀 Ready to deploy?',
          ]
        case 'exploring':
          return [
            '🗺️ Create a diagram of the architecture',
            '📝 Take notes on key findings',
            '🔗 Map out dependencies',
            '💾 Bookmark important files',
          ]
        case 'productive':
          return [
            '⚡ Keep the momentum going!',
            '💾 Maybe time for a checkpoint?',
            '🧹 Quick refactor while things are fresh?',
          ]
        default:
          return []
      }
    }

    adaptUI(): UIAdjustments {
      switch (this.currentMood) {
        case 'frustrated':
          return {
            simplify: true,
            showHelp: true,
            colorScheme: 'calming',
            verbosity: 'high',
          }
        case 'productive':
          return {
            simplify: false,
            showHelp: false,
            colorScheme: 'energetic',
            verbosity: 'low',
          }
        case 'debugging':
          return {
            simplify: false,
            showHelp: true,
            colorScheme: 'focused',
            verbosity: 'medium',
          }
        case 'celebrating':
          return {
            simplify: false,
            showHelp: false,
            colorScheme: 'celebratory',
            verbosity: 'low',
            animations: true,
          }
        default:
          return {
            simplify: false,
            showHelp: false,
            colorScheme: 'default',
            verbosity: 'medium',
          }
      }
    }
  }

  type Mood = 'frustrated' | 'productive' | 'exploring' | 'debugging' | 'celebrating' | 'neutral'

  interface Action {
    type: 'error' | 'success' | 'read' | 'write' | 'command'
    timestamp: number
    command?: string
    file?: string
    duration?: number
  }

  interface UIAdjustments {
    simplify: boolean
    showHelp: boolean
    colorScheme: string
    verbosity: 'low' | 'medium' | 'high'
    animations?: boolean
  }

  /**
   * Background Task Orchestrator
   */
  export class BackgroundOrchestrator {
    private tasks: Map<string, BackgroundTask> = new Map()

    register(task: BackgroundTask): void {
      this.tasks.set(task.id, task)
    }

    async start(taskId: string): Promise<void> {
      const task = this.tasks.get(taskId)
      if (!task || task.running) return

      task.running = true
      task.startTime = Date.now()

      try {
        await task.execute()
        task.status = 'completed'
      } catch (error) {
        task.status = 'failed'
        task.error = error instanceof Error ? error.message : 'Unknown error'
      } finally {
        task.running = false
        task.endTime = Date.now()
      }
    }

    getStatus(): string {
      const lines: string[] = []

      lines.push(UI.glow('Background Tasks', UI.Style.CLAUDE_BLUE))
      lines.push('')

      for (const task of this.tasks.values()) {
        const icon =
          task.running ? UI.Style.NEON_CYAN + '◉' :
          task.status === 'completed' ? UI.Style.MATRIX_GREEN + '✓' :
          task.status === 'failed' ? UI.Style.TEXT_DANGER + '✖' :
          UI.Style.TEXT_DIM + '○'

        lines.push(`${icon} ${task.name}${UI.Style.RESET}`)

        if (task.running && task.progress !== undefined) {
          const bar = this.renderProgressBar(task.progress)
          lines.push(`  ${bar}`)
        }

        if (task.error) {
          lines.push(`  ${UI.Style.TEXT_DANGER}Error: ${task.error}${UI.Style.RESET}`)
        }
      }

      return lines.join('\n')
    }

    private renderProgressBar(progress: number): string {
      const width = 30
      const filled = Math.floor(width * progress)
      const empty = width - filled

      return (
        UI.Style.CLAUDE_BLUE +
        '[' +
        '█'.repeat(filled) +
        '░'.repeat(empty) +
        '] ' +
        Math.floor(progress * 100) +
        '%' +
        UI.Style.RESET
      )
    }
  }

  interface BackgroundTask {
    id: string
    name: string
    description: string
    execute: () => Promise<void>
    running: boolean
    status: 'pending' | 'running' | 'completed' | 'failed'
    progress?: number
    startTime?: number
    endTime?: number
    error?: string
  }

  /**
   * Celebration System
   */
  export class CelebrationEngine {
    celebrate(event: CelebrationType, context?: any): void {
      const message = this.getCelebrationMessage(event, context)
      const animation = this.getCelebrationAnimation(event)

      console.log('\n')
      console.log(animation)
      console.log(message)
      console.log('\n')
    }

    private getCelebrationMessage(event: CelebrationType, context?: any): string {
      switch (event) {
        case 'tests_passing':
          return UI.gradient(
            `🎉 All tests passing! ${context?.count || ''} tests succeeded`,
            [UI.Style.MATRIX_GREEN, UI.Style.GEMINI_GREEN]
          )
        case 'build_success':
          return UI.gradient(
            '✨ Build successful! Ready to deploy',
            [UI.Style.CLAUDE_BLUE, UI.Style.NEON_CYAN]
          )
        case 'bug_fixed':
          return UI.gradient(
            '🐛 Bug squashed! Great debugging',
            [UI.Style.CYBER_PURPLE, UI.Style.NEON_MAGENTA]
          )
        case 'milestone':
          return UI.gradient(
            `🚀 Milestone achieved: ${context?.name || 'unnamed'}`,
            [UI.Style.MATRIX_GREEN, UI.Style.CLAUDE_BLUE, UI.Style.CYBER_PURPLE]
          )
        case 'first_contribution':
          return UI.gradient(
            '🌟 First contribution! Welcome to the project',
            [UI.Style.GEMINI_GREEN, UI.Style.NEON_CYAN]
          )
        default:
          return UI.glow('🎊 Success!', UI.Style.MATRIX_GREEN)
      }
    }

    private getCelebrationAnimation(event: CelebrationType): string {
      const animations = {
        tests_passing: `
${UI.Style.MATRIX_GREEN}
   ✓  ✓  ✓
  ✓ ✓ ✓ ✓ ✓
   ✓  ✓  ✓
${UI.Style.RESET}`,
        build_success: `
${UI.Style.CLAUDE_BLUE}
    ╭━━━━━━╮
    │ ✨ ✨ │
    ╰━━━━━━╯
${UI.Style.RESET}`,
        bug_fixed: `
${UI.Style.CYBER_PURPLE}
    🐛 → 💥 → ✨
${UI.Style.RESET}`,
        milestone: `
${UI.Style.NEON_CYAN}
       🏆
    ╱╲╱╲╱╲
   ╱  ╲  ╲
${UI.Style.RESET}`,
        first_contribution: `
${UI.Style.GEMINI_GREEN}
      ★
     ★ ★
    ★ ★ ★
${UI.Style.RESET}`,
      }

      return animations[event] || ''
    }
  }

  type CelebrationType = 'tests_passing' | 'build_success' | 'bug_fixed' | 'milestone' | 'first_contribution'

  /**
   * Contextual Command Palette
   */
  export class ContextualPalette {
    private context: Context = {}
    private commandHistory: string[] = []

    updateContext(newContext: Partial<Context>): void {
      this.context = { ...this.context, ...newContext }
    }

    getSuggestions(): CommandSuggestion[] {
      const suggestions: CommandSuggestion[] = []

      // Based on current file type
      if (this.context.currentFile?.endsWith('.test.ts')) {
        suggestions.push(
          { command: '/test', relevance: 0.95, reason: 'You\'re in a test file' },
          { command: '/coverage', relevance: 0.8, reason: 'Check test coverage' }
        )
      }

      if (this.context.currentFile?.endsWith('.tsx') || this.context.currentFile?.endsWith('.jsx')) {
        suggestions.push(
          { command: '/preview', relevance: 0.9, reason: 'Preview React component' },
          { command: '/test components', relevance: 0.8, reason: 'Test React components' }
        )
      }

      // Based on recent errors
      if (this.context.hasErrors) {
        suggestions.push(
          { command: '/debug', relevance: 0.95, reason: 'Debug recent errors' },
          { command: '/fix', relevance: 0.9, reason: 'Auto-fix issues' },
          { command: '/lint --fix', relevance: 0.8, reason: 'Fix linting errors' }
        )
      }

      // Based on uncommitted changes
      if (this.context.hasUncommittedChanges) {
        suggestions.push(
          { command: '/commit', relevance: 0.85, reason: 'Commit your changes' },
          { command: '/review', relevance: 0.75, reason: 'Review changes before commit' }
        )
      }

      // Based on command history patterns
      const lastCommands = this.commandHistory.slice(-3)
      if (lastCommands.filter(c => c === '/test').length >= 2) {
        suggestions.push(
          { command: '/watch test', relevance: 0.9, reason: 'Auto-run tests on changes' }
        )
      }

      // Sort by relevance
      return suggestions.sort((a, b) => b.relevance - a.relevance)
    }

    addToHistory(command: string): void {
      this.commandHistory.push(command)
      if (this.commandHistory.length > 100) {
        this.commandHistory.shift()
      }
    }

    render(): string {
      const suggestions = this.getSuggestions().slice(0, 5)
      const lines: string[] = []

      lines.push(UI.glow('Suggested Commands', UI.Style.CLAUDE_BLUE))
      lines.push('')

      for (const [index, suggestion] of suggestions.entries()) {
        const key = (index + 1).toString()
        const relevanceBar = '█'.repeat(Math.floor(suggestion.relevance * 10))

        lines.push(
          `${UI.Style.NEON_CYAN}${key}${UI.Style.RESET} ${UI.Style.MATRIX_GREEN}${suggestion.command}${UI.Style.RESET}`
        )
        lines.push(
          `  ${UI.Style.TEXT_DIM}${suggestion.reason} ${UI.Style.CLAUDE_BLUE}${relevanceBar}${UI.Style.RESET}`
        )
      }

      return lines.join('\n')
    }
  }

  interface Context {
    currentFile?: string
    hasErrors?: boolean
    hasUncommittedChanges?: boolean
    currentBranch?: string
    openFiles?: string[]
  }

  interface CommandSuggestion {
    command: string
    relevance: number // 0-1
    reason: string
  }
}
