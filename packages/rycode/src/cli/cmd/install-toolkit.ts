/**
 * Install Toolkit CLI Command
 *
 * Installs the toolkit-cli Python package required by the bundled toolkit client
 */

import type { CommandModule } from 'yargs'
import { spawn } from 'child_process'
import { UI } from '../ui'

export const InstallToolkitCommand: CommandModule = {
  command: 'install-toolkit-cli',
  describe: 'Install toolkit-cli Python package',
  builder: (yargs) =>
    yargs
      .option('upgrade', {
        alias: 'u',
        describe: 'Upgrade to latest version',
        type: 'boolean',
        default: false,
      })
      .option('version', {
        alias: 'v',
        describe: 'Install specific version',
        type: 'string',
      }),
  handler: async (argv) => {
    console.log('')
    UI.header('Toolkit-CLI Installation')
    console.log('')

    // Check if Python is installed
    process.stdout.write('🔍 Checking Python... ')
    const pythonCheck = await checkPython()

    if (!pythonCheck.installed) {
      console.log('❌')
      console.log('')
      UI.error('Python 3.11+ is required but not found')
      console.log('')
      console.log('Install Python 3.11+:')
      console.log('  macOS:   brew install python@3.11')
      console.log('  Linux:   apt-get install python3.11')
      console.log('  Windows: https://python.org/downloads')
      console.log('')
      process.exit(1)
    }

    console.log(`✅ ${pythonCheck.version}`)

    // Build pip install command
    const args = ['install', '--quiet']

    if (argv.upgrade) {
      args.push('--upgrade')
      process.stdout.write('📦 Upgrading toolkit-cli... ')
    } else {
      process.stdout.write('📦 Installing toolkit-cli... ')
    }

    if (argv.version) {
      args.push(`toolkit-cli==${argv.version}`)
    } else {
      args.push('toolkit-cli')
    }

    // Run pip install
    const success = await installToolkit(args)

    if (success) {
      console.log('✅')
      console.log('')

      // Verify installation
      process.stdout.write('🔍 Verifying installation... ')
      const version = await getToolkitVersion()

      if (version) {
        console.log(`✅ v${version}`)
        console.log('')
        console.log('╔═══════════════════════════════════════════╗')
        console.log('║  ✅ Installation Complete                ║')
        console.log('╚═══════════════════════════════════════════╝')
        console.log('')
        console.log('Next Steps:')
        console.log('')
        console.log('1️⃣  Configure API Keys')
        console.log('   export ANTHROPIC_API_KEY="sk-ant-..."')
        console.log('   export RYCODE_API_KEY="..."')
        console.log('')
        console.log('2️⃣  Check Health')
        console.log('   rycode toolkit health')
        console.log('')
        console.log('3️⃣  Try a Command')
        console.log('   rycode toolkit oneshot "Build a todo app"')
        console.log('')
      } else {
        console.log('⚠️')
        console.log('')
        UI.error('Installation completed but verification failed')
        console.log('Try: toolkit-cli version')
        console.log('')
      }
    } else {
      console.log('❌')
      console.log('')
      UI.error('Installation failed')
      console.log('')
      console.log('Troubleshooting:')
      console.log('  • Check pip: pip3 --version')
      console.log('  • Try sudo: sudo pip3 install toolkit-cli')
      console.log('  • Use venv: python3 -m venv venv')
      console.log('')
      process.exit(1)
    }
  },
}

/**
 * Check if Python is installed and get version
 */
async function checkPython(): Promise<{ installed: boolean; version?: string }> {
  return new Promise((resolve) => {
    const proc = spawn('python3', ['--version'])

    let output = ''

    proc.stdout?.on('data', (data) => {
      output += data.toString()
    })

    proc.stderr?.on('data', (data) => {
      output += data.toString()
    })

    proc.on('close', (code) => {
      if (code === 0) {
        const match = output.match(/Python (\d+\.\d+\.\d+)/)
        const version = match ? match[1] : undefined

        // Check if version is 3.11+
        if (version) {
          const [major, minor] = version.split('.').map(Number)
          const valid = major === 3 && minor >= 11

          resolve({ installed: valid, version })
        } else {
          resolve({ installed: false })
        }
      } else {
        resolve({ installed: false })
      }
    })

    proc.on('error', () => {
      resolve({ installed: false })
    })
  })
}

/**
 * Install toolkit-cli using pip
 */
async function installToolkit(args: string[]): Promise<boolean> {
  return new Promise((resolve) => {
    const proc = spawn('pip3', args, {
      stdio: 'pipe', // Capture output for quiet install
    })

    let stdout = ''
    let stderr = ''

    proc.stdout?.on('data', (data) => {
      stdout += data.toString()
    })

    proc.stderr?.on('data', (data) => {
      stderr += data.toString()
    })

    proc.on('close', (code) => {
      if (code !== 0) {
        // Show error details if failed
        console.log('')
        if (stderr) console.error(stderr)
        if (stdout) console.log(stdout)
      }
      resolve(code === 0)
    })

    proc.on('error', (error) => {
      console.log('')
      console.error('Error running pip:', error.message)
      resolve(false)
    })
  })
}

/**
 * Get installed toolkit-cli version
 */
async function getToolkitVersion(): Promise<string | null> {
  return new Promise((resolve) => {
    const proc = spawn('toolkit-cli', ['version'])

    let output = ''

    proc.stdout?.on('data', (data) => {
      output += data.toString()
    })

    proc.on('close', (code) => {
      if (code === 0) {
        // Parse version from output
        const match = output.match(/(\d+\.\d+\.\d+)/)
        resolve(match ? match[1] : null)
      } else {
        resolve(null)
      }
    })

    proc.on('error', () => {
      resolve(null)
    })
  })
}
