/**
 * Toolkit Command - Access toolkit-cli from RyCode
 *
 * Provides CLI commands for using the bundled toolkit-cli client
 */

import type { CommandModule } from 'yargs'
import { ToolkitClient } from '../../toolkit'
import { UI } from '../ui'

export const ToolkitCommand: CommandModule = {
  command: 'toolkit <subcommand>',
  describe: 'Access toolkit-cli AI commands',
  builder: (yargs) =>
    yargs
      .command({
        command: 'health',
        describe: 'Check toolkit-cli installation status',
        handler: async () => {
          const toolkit = new ToolkitClient()

          UI.header('Toolkit Health Check')

          try {
            const health = await toolkit.health()

            if (health.healthy) {
              UI.success('✅ Toolkit is healthy')
              console.log('')
              console.log('📦 Version:', health.toolkitCliVersion)
              console.log('🐍 Python:', health.pythonVersion)
              console.log('')
              console.log('🤖 Agents:')
              health.agentsAvailable.forEach((agent) => {
                const status = agent.configured ? '✅' : '⚠️ '
                console.log(`   ${status} ${agent.name}`)
              })
            } else {
              UI.error('❌ Toolkit is not healthy')
              console.log('')
              console.log('Issues:')
              health.issues.forEach((issue) => {
                console.log(`   • ${issue}`)
              })
            }
          } catch (error) {
            UI.error('Failed to check toolkit health')
            if (error instanceof Error) {
              console.error(error.message)
            }
          } finally {
            await toolkit.close()
          }
        },
      })
      .command({
        command: 'oneshot <idea>',
        describe: 'Generate complete project specification',
        builder: (yargs) =>
          yargs
            .positional('idea', {
              describe: 'Project idea',
              type: 'string',
            })
            .option('agents', {
              alias: 'a',
              describe: 'AI agents to use (comma-separated)',
              type: 'string',
              default: 'claude,rycode',
            })
            .option('complexity', {
              alias: 'c',
              describe: 'Project complexity',
              type: 'string',
              choices: ['low', 'medium', 'high', 'enterprise'],
              default: 'medium',
            })
            .option('ux', {
              describe: 'Include UX designs',
              type: 'boolean',
              default: true,
            }),
        handler: async (argv) => {
          if (!argv.idea) {
            UI.error('Project idea is required')
            return
          }

          const agents = (argv.agents as string).split(',')
          const toolkit = new ToolkitClient({
            agents: agents as any,
          })

          UI.header('Generating Project Specification')
          console.log('')
          console.log('💡 Idea:', argv.idea)
          console.log('🤖 Agents:', agents.join(', '))
          console.log('📊 Complexity:', argv.complexity)
          console.log('')

          try {
            const result = await toolkit.oneshot(argv.idea as string, {
              agents: agents as any,
              complexity: argv.complexity as any,
              includeUx: argv.ux,
              onProgress: (chunk) => {
                console.log(`[${chunk.progress}%] ${chunk.message}`)
              },
            })

            if (result.success) {
              console.log('')
              UI.success('✅ Specification generated')
              console.log('')
              console.log('═══════════════════════════════════')
              console.log('📋 OVERVIEW')
              console.log('═══════════════════════════════════')
              console.log(result.data?.specification.overview)
              console.log('')
              console.log('═══════════════════════════════════')
              console.log('🏗️  ARCHITECTURE')
              console.log('═══════════════════════════════════')
              console.log(result.data?.architecture.overview)
              console.log('')
              console.log('📊 Metrics:')
              console.log(`   ⏱️  Time: ${result.metadata.executionTime}ms`)
              console.log(`   🤖 Agents: ${result.metadata.agentsUsed.join(', ')}`)
            } else {
              UI.error('❌ Generation failed')
              console.error(result.error?.message)
            }
          } catch (error) {
            UI.error('Failed to generate specification')
            if (error instanceof Error) {
              console.error(error.message)
            }
          } finally {
            await toolkit.close()
          }
        },
      })
      .command({
        command: 'fix <issue>',
        describe: 'Analyze and fix code issues',
        builder: (yargs) =>
          yargs
            .positional('issue', {
              describe: 'Issue description',
              type: 'string',
            })
            .option('context', {
              alias: 'c',
              describe: 'Additional context',
              type: 'string',
            }),
        handler: async (argv) => {
          if (!argv.issue) {
            UI.error('Issue description is required')
            return
          }

          const toolkit = new ToolkitClient()

          UI.header('Analyzing Issue')
          console.log('')
          console.log('🐛 Issue:', argv.issue)
          if (argv.context) {
            console.log('📝 Context:', argv.context)
          }
          console.log('')

          try {
            const result = await toolkit.fix(argv.issue as string, {
              context: argv.context as string | undefined,
              agents: ['claude'],
            })

            if (result.success) {
              console.log('')
              UI.success('✅ Analysis complete')
              console.log('')
              console.log('═══════════════════════════════════')
              console.log('🔍 ROOT CAUSE')
              console.log('═══════════════════════════════════')
              console.log(result.data?.rootCause)
              console.log('')
              console.log('═══════════════════════════════════')
              console.log('💡 SOLUTION')
              console.log('═══════════════════════════════════')
              console.log(result.data?.solution.approach)
              console.log('')
              console.log('📝 Code Changes:')
              result.data?.solution.codeChanges.forEach((change, idx) => {
                console.log(`   ${idx + 1}. ${change.file}`)
                console.log(`      ${change.explanation}`)
              })
            } else {
              UI.error('❌ Analysis failed')
              console.error(result.error?.message)
            }
          } catch (error) {
            UI.error('Failed to analyze issue')
            if (error instanceof Error) {
              console.error(error.message)
            }
          } finally {
            await toolkit.close()
          }
        },
      })
      .demandCommand(1, 'You must specify a subcommand'),
  handler: () => {
    // No-op, subcommands handle everything
  },
}
