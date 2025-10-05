/**
 * Temporal Navigation & Time Travel System
 *
 * Allows users to scrub through conversation history,
 * create state snapshots, and explore parallel branches.
 */

import { UI } from "../../ui"

export namespace TimeTravel {
  export interface Snapshot {
    id: string
    timestamp: number
    label?: string
    conversationState: ConversationState
    fileStates: Map<string, string> // file path -> content
    metadata: {
      messageCount: number
      currentAgent: string
      currentModel: string
      tokenCount: number
    }
  }

  export interface ConversationState {
    messages: Message[]
    currentMessageIndex: number
    branches: Branch[]
    currentBranch: string
  }

  export interface Message {
    id: string
    role: 'user' | 'assistant'
    content: string
    timestamp: number
    files?: string[]
    tools?: string[]
  }

  export interface Branch {
    id: string
    parentId?: string
    name: string
    createdAt: number
    divergencePoint: number // message index where branch started
    description?: string
  }

  export interface TimelineEvent {
    type: 'message' | 'branch' | 'snapshot' | 'error' | 'success'
    timestamp: number
    significance: number // 0-1, how important this event is
    label?: string
  }

  /**
   * Timeline Scrubber UI
   */
  export class Timeline {
    private events: TimelineEvent[] = []
    private currentPosition: number = 0
    private readonly maxWidth = 80

    addEvent(event: TimelineEvent): void {
      this.events.push(event)
      this.events.sort((a, b) => a.timestamp - b.timestamp)
    }

    render(): string {
      const timeline: string[] = []
      const step = this.events.length / this.maxWidth

      // Top border
      timeline.push(UI.Style.MATRIX_GREEN + '┌' + '─'.repeat(this.maxWidth) + '┐' + UI.Style.RESET)

      // Timeline bar
      let bar = UI.Style.MATRIX_GREEN + '│ '
      for (let i = 0; i < this.maxWidth; i++) {
        const eventIndex = Math.floor(i * step)
        const event = this.events[eventIndex]

        if (Math.abs(i - this.currentPosition) < 1) {
          bar += UI.Style.NEON_CYAN + '●' + UI.Style.MATRIX_GREEN
        } else if (event?.type === 'error') {
          bar += UI.Style.TEXT_DANGER + '✖' + UI.Style.MATRIX_GREEN
        } else if (event?.type === 'success') {
          bar += UI.Style.MATRIX_GREEN + '✓' + UI.Style.MATRIX_GREEN
        } else if (event?.type === 'branch') {
          bar += UI.Style.CYBER_PURPLE + '⎇' + UI.Style.MATRIX_GREEN
        } else if (event?.type === 'snapshot') {
          bar += UI.Style.CLAUDE_BLUE + '◆' + UI.Style.MATRIX_GREEN
        } else {
          bar += '═'
        }
      }
      bar += ' │' + UI.Style.RESET
      timeline.push(bar)

      // Bottom border with legend
      timeline.push(UI.Style.MATRIX_GREEN + '└' + '─'.repeat(this.maxWidth) + '┘' + UI.Style.RESET)

      // Legend
      timeline.push('')
      timeline.push(
        UI.Style.TEXT_DIM +
        `  ${UI.Style.NEON_CYAN}●${UI.Style.TEXT_DIM} Current  ` +
        `${UI.Style.MATRIX_GREEN}✓${UI.Style.TEXT_DIM} Success  ` +
        `${UI.Style.TEXT_DANGER}✖${UI.Style.TEXT_DIM} Error  ` +
        `${UI.Style.CYBER_PURPLE}⎇${UI.Style.TEXT_DIM} Branch  ` +
        `${UI.Style.CLAUDE_BLUE}◆${UI.Style.TEXT_DIM} Snapshot` +
        UI.Style.RESET
      )

      return timeline.join('\n')
    }

    scrubTo(position: number): void {
      this.currentPosition = Math.max(0, Math.min(this.maxWidth - 1, position))
    }

    getCurrentEvent(): TimelineEvent | undefined {
      const step = this.events.length / this.maxWidth
      const eventIndex = Math.floor(this.currentPosition * step)
      return this.events[eventIndex]
    }
  }

  /**
   * Snapshot Manager
   */
  export class SnapshotManager {
    private snapshots: Map<string, Snapshot> = new Map()
    private autoSnapshotEnabled = true

    createSnapshot(state: ConversationState, label?: string): Snapshot {
      const snapshot: Snapshot = {
        id: this.generateId(),
        timestamp: Date.now(),
        label: label || `Snapshot ${this.snapshots.size + 1}`,
        conversationState: this.deepClone(state),
        fileStates: new Map(),
        metadata: {
          messageCount: state.messages.length,
          currentAgent: 'build', // TODO: get from actual state
          currentModel: 'claude', // TODO: get from actual state
          tokenCount: 0, // TODO: calculate
        },
      }

      this.snapshots.set(snapshot.id, snapshot)
      return snapshot
    }

    restoreSnapshot(id: string): Snapshot | null {
      return this.snapshots.get(id) || null
    }

    deleteSnapshot(id: string): boolean {
      return this.snapshots.delete(id)
    }

    listSnapshots(): Snapshot[] {
      return Array.from(this.snapshots.values()).sort((a, b) => b.timestamp - a.timestamp)
    }

    autoSnapshot(state: ConversationState): void {
      if (!this.autoSnapshotEnabled) return

      // Auto-snapshot before major operations
      const lastMessage = state.messages[state.messages.length - 1]
      const shouldSnapshot =
        lastMessage?.tools?.includes('write') ||
        lastMessage?.tools?.includes('edit') ||
        state.messages.length % 10 === 0 // Every 10 messages

      if (shouldSnapshot) {
        this.createSnapshot(state, `Auto: ${new Date().toLocaleTimeString()}`)
      }
    }

    private generateId(): string {
      return `snapshot_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    }

    private deepClone<T>(obj: T): T {
      return JSON.parse(JSON.stringify(obj))
    }
  }

  /**
   * Branch Manager for parallel exploration
   */
  export class BranchManager {
    private branches: Map<string, Branch> = new Map()
    private currentBranch: string = 'main'

    createBranch(name: string, divergencePoint: number, parentId?: string): Branch {
      const branch: Branch = {
        id: this.generateBranchId(name),
        parentId: parentId || this.currentBranch,
        name,
        createdAt: Date.now(),
        divergencePoint,
        description: `Branched from ${parentId || this.currentBranch} at message ${divergencePoint}`,
      }

      this.branches.set(branch.id, branch)
      return branch
    }

    switchBranch(branchId: string): Branch | null {
      const branch = this.branches.get(branchId)
      if (!branch) return null

      this.currentBranch = branchId
      return branch
    }

    mergeBranch(sourceBranchId: string, targetBranchId: string): boolean {
      // TODO: Implement merge logic
      return true
    }

    visualizeBranches(): string {
      const tree: string[] = []
      const sortedBranches = Array.from(this.branches.values()).sort(
        (a, b) => a.createdAt - b.createdAt
      )

      tree.push(UI.glow('Branch Tree', UI.Style.CLAUDE_BLUE))
      tree.push('')

      for (const branch of sortedBranches) {
        const isCurrent = branch.id === this.currentBranch
        const indent = branch.parentId ? '  ' : ''
        const icon = isCurrent ? UI.Style.NEON_CYAN + '●' : UI.Style.TEXT_DIM + '○'
        const name = isCurrent ? UI.glow(branch.name, UI.Style.NEON_CYAN) : branch.name

        tree.push(`${indent}${icon} ${name}${UI.Style.RESET}`)
        if (branch.description) {
          tree.push(`${indent}  ${UI.Style.TEXT_DIM}${branch.description}${UI.Style.RESET}`)
        }
      }

      return tree.join('\n')
    }

    private generateBranchId(name: string): string {
      return `branch_${name.toLowerCase().replace(/\s+/g, '_')}_${Date.now()}`
    }
  }

  /**
   * Undo/Redo with Context
   */
  export class ContextualUndoRedo {
    private history: Array<{
      action: string
      state: any
      reason: string
      timestamp: number
    }> = []
    private currentIndex: number = -1

    record(action: string, state: any, reason: string): void {
      // Remove any future history if we're not at the end
      this.history = this.history.slice(0, this.currentIndex + 1)

      this.history.push({
        action,
        state: JSON.parse(JSON.stringify(state)),
        reason,
        timestamp: Date.now(),
      })

      this.currentIndex++
    }

    undo(): { state: any; context: string } | null {
      if (this.currentIndex <= 0) return null

      this.currentIndex--
      const entry = this.history[this.currentIndex]

      return {
        state: entry.state,
        context: `Undoing: ${entry.reason}`,
      }
    }

    redo(): { state: any; context: string } | null {
      if (this.currentIndex >= this.history.length - 1) return null

      this.currentIndex++
      const entry = this.history[this.currentIndex]

      return {
        state: entry.state,
        context: `Redoing: ${entry.reason}`,
      }
    }

    getHistory(): string {
      const lines: string[] = []

      lines.push(UI.glow('Action History', UI.Style.CLAUDE_BLUE))
      lines.push('')

      this.history.forEach((entry, index) => {
        const isCurrent = index === this.currentIndex
        const marker = isCurrent ? UI.Style.NEON_CYAN + '→' : ' '
        const time = new Date(entry.timestamp).toLocaleTimeString()

        lines.push(
          `${marker} ${UI.Style.TEXT_DIM}${time}${UI.Style.RESET} ${entry.action} - ${UI.Style.TEXT_DIM}${entry.reason}${UI.Style.RESET}`
        )
      })

      return lines.join('\n')
    }
  }

  /**
   * Parallel Universe Mode
   * Try multiple approaches simultaneously
   */
  export class ParallelExplorer {
    private universes: Map<string, {
      id: string
      name: string
      divergencePoint: number
      messages: Message[]
      result?: 'success' | 'failure' | 'pending'
    }> = new Map()

    createUniverse(name: string, divergencePoint: number): string {
      const id = `universe_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`

      this.universes.set(id, {
        id,
        name,
        divergencePoint,
        messages: [],
        result: 'pending',
      })

      return id
    }

    compareUniverses(ids: string[]): string {
      const comparison: string[] = []

      comparison.push(UI.glow('Universe Comparison', UI.Style.CLAUDE_BLUE))
      comparison.push('')

      for (const id of ids) {
        const universe = this.universes.get(id)
        if (!universe) continue

        const resultIcon =
          universe.result === 'success' ? UI.Style.MATRIX_GREEN + '✓' :
          universe.result === 'failure' ? UI.Style.TEXT_DANGER + '✖' :
          UI.Style.TEXT_DIM + '○'

        comparison.push(`${resultIcon} ${universe.name}${UI.Style.RESET}`)
        comparison.push(`  Messages: ${universe.messages.length}`)
        comparison.push(`  Diverged at: #${universe.divergencePoint}`)
        comparison.push('')
      }

      return comparison.join('\n')
    }

    mergeBest(): string | null {
      // Find universe with best result
      const successful = Array.from(this.universes.values()).filter(u => u.result === 'success')

      if (successful.length === 0) return null

      // Return the one with fewest messages (most efficient)
      const best = successful.sort((a, b) => a.messages.length - b.messages.length)[0]
      return best.id
    }
  }
}
