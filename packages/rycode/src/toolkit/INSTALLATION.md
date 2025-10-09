# Installing Toolkit-CLI for RyCode

## Quick Install

RyCode now includes a built-in command to install toolkit-cli:

```bash
rycode install-toolkit-cli
```

That's it! The command will:
- ✅ Check if Python 3.11+ is installed
- ✅ Install toolkit-cli via pip
- ✅ Verify the installation
- ✅ Show next steps for configuration

## Installation Options

### Basic Installation

```bash
rycode install-toolkit-cli
```

### Upgrade to Latest Version

```bash
rycode install-toolkit-cli --upgrade
```

or

```bash
rycode install-toolkit-cli -u
```

### Install Specific Version

```bash
rycode install-toolkit-cli --version 1.3.3
```

or

```bash
rycode install-toolkit-cli -v 1.3.3
```

## Prerequisites

### Python 3.11+

Toolkit-cli requires Python 3.11 or higher.

**macOS**:
```bash
brew install python@3.11
```

**Ubuntu/Debian**:
```bash
sudo apt-get update
sudo apt-get install python3.11 python3-pip
```

**Windows**:
Download from [python.org/downloads](https://python.org/downloads)

**Verify Python Version**:
```bash
python3 --version
# Should show: Python 3.11.x or higher
```

## Manual Installation

If you prefer to install manually:

```bash
# Using pip
pip3 install toolkit-cli

# Upgrade
pip3 install --upgrade toolkit-cli

# Specific version
pip3 install toolkit-cli==1.3.3
```

## Verify Installation

After installation, verify toolkit-cli is working:

```bash
# Check version
toolkit-cli version

# Or use RyCode's health check
rycode toolkit health
```

Expected output:
```
✅ Toolkit is healthy
📦 Version: 1.3.3
🐍 Python: 3.11.5

🤖 Agents:
   ✅ claude
   ⚠️  gemini
   ⚠️  qwen
   ...
```

## Configure API Keys

After installation, configure your AI provider API keys:

### Using Environment Variables

```bash
# Required for Claude
export ANTHROPIC_API_KEY="sk-ant-..."

# Optional - Add as needed
export OPENAI_API_KEY="sk-..."
export GOOGLE_API_KEY="..."
export QWEN_API_KEY="..."
export RYCODE_API_KEY="..."
```

### Using .env File

Create `.env` file in your project:

```env
ANTHROPIC_API_KEY=sk-ant-...
OPENAI_API_KEY=sk-...
GOOGLE_API_KEY=...
RYCODE_API_KEY=...
```

### Permanent Configuration

Add to your shell profile (`~/.zshrc` or `~/.bashrc`):

```bash
# Toolkit-CLI API Keys
export ANTHROPIC_API_KEY="sk-ant-..."
export OPENAI_API_KEY="sk-..."
export GOOGLE_API_KEY="..."
export RYCODE_API_KEY="..."
```

Then reload:
```bash
source ~/.zshrc
# or
source ~/.bashrc
```

## Test the Installation

### Check Health

```bash
rycode toolkit health
```

This shows:
- ✅ toolkit-cli installation status
- ✅ Python version
- ✅ Available AI agents
- ✅ Configured API keys
- ⚠️  Any issues

### Try a Simple Command

```bash
rycode toolkit oneshot "Simple todo app" --agents claude
```

This will:
1. Generate a complete project specification
2. Show progress updates
3. Display the results
4. Demonstrate that everything is working

## Troubleshooting

### Error: "Python 3.11+ is required but not found"

**Problem**: Python is not installed or version is too old

**Solution**:
1. Install Python 3.11+:
   ```bash
   # macOS
   brew install python@3.11

   # Ubuntu
   sudo apt-get install python3.11
   ```

2. Verify:
   ```bash
   python3 --version
   ```

3. Retry installation:
   ```bash
   rycode install-toolkit-cli
   ```

### Error: "pip3: command not found"

**Problem**: pip is not installed

**Solution**:
```bash
# macOS
python3 -m ensurepip --upgrade

# Ubuntu
sudo apt-get install python3-pip

# Verify
pip3 --version
```

### Error: "Permission denied"

**Problem**: Need sudo for system-wide install

**Solution**:
```bash
sudo pip3 install toolkit-cli
```

Or use a virtual environment (recommended):
```bash
python3 -m venv ~/venv-toolkit
source ~/venv-toolkit/bin/activate
pip install toolkit-cli
```

### Error: "toolkit-cli not found" after installation

**Problem**: toolkit-cli not in PATH

**Solution**:
1. Find where it was installed:
   ```bash
   pip3 show toolkit-cli
   ```

2. Add to PATH in `~/.zshrc` or `~/.bashrc`:
   ```bash
   export PATH="$PATH:$HOME/.local/bin"
   ```

3. Reload shell:
   ```bash
   source ~/.zshrc
   ```

### Error: "API key not configured"

**Problem**: Missing API keys

**Solution**:
1. Get API key from provider:
   - Anthropic: https://console.anthropic.com
   - OpenAI: https://platform.openai.com
   - Google: https://aistudio.google.com

2. Set environment variable:
   ```bash
   export ANTHROPIC_API_KEY="sk-ant-..."
   ```

3. Verify:
   ```bash
   rycode toolkit health
   ```

### Installation hangs or fails

**Problem**: Network issues or pip cache

**Solution**:
```bash
# Clear pip cache
pip3 cache purge

# Retry with verbose output
pip3 install toolkit-cli --verbose

# Use alternative index
pip3 install toolkit-cli --index-url https://pypi.org/simple
```

## Virtual Environments (Recommended)

Using a virtual environment isolates toolkit-cli from system Python:

```bash
# Create virtual environment
python3 -m venv ~/venv-toolkit

# Activate
source ~/venv-toolkit/bin/activate

# Install toolkit-cli
pip install toolkit-cli

# Use with RyCode
rycode toolkit health

# Deactivate when done
deactivate
```

**Note**: Remember to activate the virtual environment before using toolkit commands.

## Upgrade Toolkit-CLI

Keep toolkit-cli updated for latest features and fixes:

```bash
# Using RyCode command
rycode install-toolkit-cli --upgrade

# Or manually
pip3 install --upgrade toolkit-cli

# Check new version
toolkit-cli version
```

## Uninstall

If you need to remove toolkit-cli:

```bash
pip3 uninstall toolkit-cli
```

## Complete Installation Workflow

Here's the complete workflow from scratch:

```bash
# 1. Verify Python
python3 --version
# Should be 3.11+

# 2. Install toolkit-cli
rycode install-toolkit-cli

# 3. Configure API keys
export ANTHROPIC_API_KEY="sk-ant-..."
export RYCODE_API_KEY="..."

# 4. Verify installation
rycode toolkit health

# 5. Test with a command
rycode toolkit oneshot "Test app" --agents claude,rycode

# 6. Start using in RyCode!
```

## Next Steps

After successful installation:

1. ✅ Configure all API keys you plan to use
2. ✅ Run `rycode toolkit health` to verify
3. ✅ Try example commands from `src/toolkit/example.ts`
4. ✅ Integrate toolkit into your RyCode workflows
5. ✅ Read full documentation in `TOOLKIT_BUNDLED.md`

## Support

If you encounter issues:

1. Check this troubleshooting guide
2. Run `rycode toolkit health` for diagnostics
3. Check toolkit-cli logs
4. Open an issue with error details

---

**Quick Reference**:
- Install: `rycode install-toolkit-cli`
- Upgrade: `rycode install-toolkit-cli -u`
- Health: `rycode toolkit health`
- Help: `rycode install-toolkit-cli --help`
