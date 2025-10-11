# RyCode Phase 2 - Quick Start Guide

**Ready in 2 minutes!** ⚡

---

## What's New in Phase 2?

Three game-changing features:

1. 🎯 **Status Bar** - See your current model and today's cost at a glance
2. ⚡ **Tab Cycling** - Switch between recent models instantly
3. 🔐 **Inline Auth** - Authenticate with any provider without leaving the model dialog

---

## Try It Now!

### 1. Launch RyCode

```bash
/tmp/rycode-tui-phase2
```

### 2. Check Your Status Bar

Look at the bottom right:
```
[tab Claude 3.5 Sonnet | 💰 $0.12 | tab→]
```

**What you see:**
- Current model name
- Today's cost (updates every 5 seconds)
- Tab hint (press to cycle models)

### 3. Cycle Through Models

Press **Tab** key repeatedly:
```
Tab → GPT-4 Turbo
Tab → Claude 3 Opus
Tab → Claude 3.5 Sonnet
```

Each press switches to your next recent model. **Shift+Tab** goes backward.

### 4. Authenticate a Provider

**Option A: Keyboard shortcut (fastest)**
```
1. Ctrl+X M         (Open model dialog)
2. Navigate to a provider with 🔒
3. Press 'a'        (Auth prompt appears)
4. Type your API key
5. Press Enter      (Authenticates)
6. ✓ Done!          (Models unlock)
```

**Option B: Auto-detect (easiest)**
```
1. Set environment variable:
   export OPENAI_API_KEY="sk-..."
2. Ctrl+X M         (Open model dialog)
3. Press 'd'        (Auto-detect)
4. ✓ Done!          (All found providers unlock)
```

**Option C: Select locked model (most intuitive)**
```
1. Ctrl+X M         (Open model dialog)
2. Navigate to a locked model
3. Press Enter      (Auth prompt appears)
4. Type your API key
5. Press Enter      (Authenticates and selects model)
6. ✓ Done!
```

---

## Visual Guide

### Before Authentication
```
┌─────────────────────────────────────┐
│ Search models...                    │
├─────────────────────────────────────┤
│ Anthropic ✓                         │
│   Claude 3.5 Sonnet                 │
│   Claude 3 Opus                     │
│                                     │
│ OpenAI 🔒                           │
│   GPT-4 Turbo          [locked]     │  ← Can't select
│   GPT-3.5 Turbo        [locked]     │
└─────────────────────────────────────┘
```

### After Pressing 'a' on OpenAI
```
┌─────────────────────────────────────┐
│ Authenticate with OpenAI            │
│                                     │
│ ••••••••••••••••••                  │  ← Password input
│                                     │
│ Press Enter to submit | Ctrl+D for  │
│ auto-detect | Esc to cancel         │
└─────────────────────────────────────┘
```

### After Authentication
```
┌─────────────────────────────────────┐
│ Search models...                    │
├─────────────────────────────────────┤
│ Anthropic ✓                         │
│   Claude 3.5 Sonnet                 │
│   Claude 3 Opus                     │
│                                     │
│ OpenAI ✓                            │  ← Now authenticated!
│   GPT-4 Turbo                       │  ← Can select
│   GPT-3.5 Turbo                     │
└─────────────────────────────────────┘

Toast: ✓ Authenticated with OpenAI (8 models)
```

---

## Keyboard Shortcuts

### Essential Shortcuts

| Key | What It Does |
|-----|--------------|
| **Tab** | Switch to next recent model |
| **Ctrl+X M** | Open model dialog |
| **a** | Authenticate provider (in model dialog) |
| **d** | Auto-detect credentials (in model dialog) |
| **Esc** | Close dialog or cancel auth |

### All Available Shortcuts

See [KEYBOARD_SHORTCUTS_REFERENCE.md](./KEYBOARD_SHORTCUTS_REFERENCE.md) for complete list.

---

## Common Tasks

### Authenticate with OpenAI
```bash
# Set your API key
export OPENAI_API_KEY="sk-proj-..."

# Launch RyCode
/tmp/rycode-tui-phase2

# Auto-detect
Ctrl+X M → d

# Or manually
Ctrl+X M → navigate to OpenAI → a → paste key → Enter
```

### Authenticate with Anthropic
```bash
export ANTHROPIC_API_KEY="sk-ant-..."
/tmp/rycode-tui-phase2
Ctrl+X M → d
```

### Authenticate with Google
```bash
export GOOGLE_API_KEY="..."
/tmp/rycode-tui-phase2
Ctrl+X M → d
```

### Check Which Providers Need Auth
```bash
/tmp/rycode-tui-phase2
Ctrl+X M
# Look for 🔒 icons
Esc
```

### Quickly Switch Between 3 Models
```bash
# Use any 3 models first to build history
# Then:
Tab → Tab → Tab
# Cycles through your 3 recent models
```

---

## Provider Status Icons

| Icon | Meaning | What To Do |
|------|---------|------------|
| ✓ | Authenticated & working | Nothing - use normally |
| 🔒 | Not authenticated | Press 'a' to authenticate |
| ⚠ | Having issues | Wait a bit, may recover |
| ✗ | Currently down | Check provider status page |

---

## Troubleshooting

### "No credentials found"

**Problem:** Pressed 'd' but nothing happened

**Solution:**
```bash
# Check your environment
env | grep API_KEY

# If nothing, set keys:
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."

# Try again
Ctrl+X M → d
```

### Can't select a model

**Problem:** Model is grayed out with "[locked]"

**Solution:** Provider needs authentication
```
Ctrl+X M → navigate to model → a → enter key → Enter
```

### Auth prompt won't appear

**Problem:** Pressed 'a' but nothing happened

**Solution:** Make sure you're on a locked provider (🔒)
```
Ctrl+X M → find provider with 🔒 → a
```

### Models won't unlock after auth

**Problem:** Entered key but still locked

**Solutions:**
1. Check if key is valid
2. Look at provider icon (✗ = provider down)
3. Try closing and reopening dialog (Esc then Ctrl+X M)
4. Check server logs for errors

---

## Environment Variables

### Supported Providers

```bash
# Anthropic
export ANTHROPIC_API_KEY="sk-ant-..."

# OpenAI
export OPENAI_API_KEY="sk-proj-..."

# Google
export GOOGLE_API_KEY="..."

# Grok (X.AI)
export GROK_API_KEY="..."

# Qwen (Alibaba)
export QWEN_API_KEY="..."
```

### Quick Setup Script

Create `~/.rycode_keys.sh`:
```bash
#!/bin/bash
export ANTHROPIC_API_KEY="sk-ant-..."
export OPENAI_API_KEY="sk-proj-..."
export GOOGLE_API_KEY="..."
```

Then:
```bash
source ~/.rycode_keys.sh
/tmp/rycode-tui-phase2
Ctrl+X M → d
# All providers unlock automatically!
```

---

## Tips for Power Users

### 1. Set All Keys Once
```bash
# Add to ~/.bashrc or ~/.zshrc
export ANTHROPIC_API_KEY="sk-ant-..."
export OPENAI_API_KEY="sk-proj-..."
export GOOGLE_API_KEY="..."

# Then auto-detect on first launch
```

### 2. Keyboard-Only Workflow
```bash
Ctrl+X M    # Open dialog
d           # Auto-detect
Esc         # Close
Tab         # Cycle models
Enter       # Start chatting
```

### 3. Quick Model Check
```bash
Ctrl+X M    # Peek at status
Esc         # Close immediately
# Just checking which models are available
```

### 4. Fastest Model Switch
```bash
Tab         # Instant switch to last model
# No dialog needed!
```

---

## What's Next?

After you're comfortable with Phase 2:

1. **Read Full Docs**
   - [INLINE_AUTH_PHASE2_COMPLETE.md](./INLINE_AUTH_PHASE2_COMPLETE.md) - Complete feature guide
   - [KEYBOARD_SHORTCUTS_REFERENCE.md](./KEYBOARD_SHORTCUTS_REFERENCE.md) - All shortcuts

2. **Explore Advanced Features**
   - Multiple provider auth
   - Health monitoring
   - Cost tracking

3. **Provide Feedback**
   - Report bugs
   - Suggest improvements
   - Share workflows

---

## Getting Help

### Documentation

- **Quick Start** - This document
- **Full Guide** - [INLINE_AUTH_PHASE2_COMPLETE.md](./INLINE_AUTH_PHASE2_COMPLETE.md)
- **Shortcuts** - [KEYBOARD_SHORTCUTS_REFERENCE.md](./KEYBOARD_SHORTCUTS_REFERENCE.md)
- **Testing** - [MANUAL_TESTING_GUIDE.md](./MANUAL_TESTING_GUIDE.md)

### Common Issues

See "Troubleshooting" section above or check full documentation.

---

## Summary

**Three things to remember:**

1. **Tab** = Switch models quickly
2. **Ctrl+X M → a** = Authenticate
3. **Ctrl+X M → d** = Auto-detect credentials

That's it! You're ready to use Phase 2.

---

**Enjoy the improved workflow!** 🚀

*Built with ❤️ by the RyCode team*
