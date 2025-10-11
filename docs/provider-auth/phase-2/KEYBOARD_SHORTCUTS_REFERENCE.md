# RyCode TUI - Keyboard Shortcuts Reference (Phase 2)

**Version:** Phase 2 Complete
**Date:** October 11, 2024

---

## New Shortcuts (Phase 2)

### Model Dialog

| Key | Context | Action | Description |
|-----|---------|--------|-------------|
| **a** | On provider header or model | **Authenticate** | Opens auth prompt for provider |
| **d** | In model list | **Auto-detect** | Scans for credentials automatically |
| **Enter** | On locked model | **Auth & Select** | Shows auth prompt, then selects model |
| **Tab** | Anywhere | **Next Model** | Cycles to next recent model |
| **Shift+Tab** | Anywhere | **Previous Model** | Cycles to previous recent model |

### Auth Prompt

| Key | Context | Action | Description |
|-----|---------|--------|-------------|
| **Enter** | In auth prompt | **Submit** | Submits API key for validation |
| **Ctrl+D** | In auth prompt | **Auto-detect** | Closes prompt and runs auto-detect |
| **Esc** | In auth prompt | **Cancel** | Closes prompt without authenticating |

---

## Existing Shortcuts (Unchanged)

### Global

| Key | Action | Description |
|-----|--------|-------------|
| **Ctrl+C** | Quit | Exit RyCode TUI |
| **Ctrl+L** | Clear | Clear terminal screen |

### Dialogs

| Key | Action | Description |
|-----|--------|-------------|
| **Ctrl+X M** | Model Dialog | Open model selector |
| **Ctrl+X A** | Agent Dialog | Open agent selector |
| **Ctrl+X L** | Session List | Open session manager |
| **Esc** | Close | Close current dialog |

### Navigation

| Key | Action | Description |
|-----|--------|-------------|
| **↑ / Ctrl+P** | Up | Move up in lists |
| **↓ / Ctrl+N** | Down | Move down in lists |
| **Enter** | Select | Select current item |
| **Ctrl+X** | Remove | Remove from recent (in lists) |

### Model Management

| Key | Action | Description |
|-----|--------|-------------|
| **F2** | Cycle Model | Alternative to Tab (same function) |

---

## Quick Start Guide

### Authenticating with a Provider

**Option 1: Using 'a' key**
```
1. Ctrl+X M         → Open model dialog
2. Navigate to provider
3. Press 'a'        → Auth prompt appears
4. Type API key     → Shows as ••••••
5. Press Enter      → Authenticates
6. ✓ Models unlock
```

**Option 2: Auto-detect**
```
1. Set env var:     export OPENAI_API_KEY="sk-..."
2. Ctrl+X M         → Open model dialog
3. Press 'd'        → Scans for credentials
4. ✓ Models unlock  → Toast shows "Auto-detected X credential(s)"
```

**Option 3: Select locked model**
```
1. Ctrl+X M         → Open model dialog
2. Navigate to locked model
3. Press Enter      → Auth prompt appears
4. Type API key
5. Press Enter      → Authenticates
6. ✓ Model selected → Automatically after auth
```

---

## Common Workflows

### Switch Models Quickly
```
Tab → Tab → Tab
(Cycles through recently used models)
```

### Authenticate Multiple Providers
```
Ctrl+X M → Navigate to Provider1 → a → Enter key → Success
         → Navigate to Provider2 → a → Enter key → Success
         → Navigate to Provider3 → a → Enter key → Success
         → Esc
```

### Check What Needs Authentication
```
Ctrl+X M
Look for 🔒 icons next to provider names
```

### Cancel Authentication
```
Ctrl+X M → a → (start typing) → Esc → (back to model list)
```

---

## Visual Indicators

### Provider Status Icons

| Icon | Meaning | Can Select Models |
|------|---------|-------------------|
| ✓ | Authenticated & Healthy | ✅ Yes |
| ⚠ | Authenticated but Degraded | ✅ Yes |
| ✗ | Authenticated but Down | ❌ No |
| 🔒 | Not Authenticated | ❌ No |

### Model Status

| Appearance | Meaning | Selectable |
|------------|---------|------------|
| `Claude 3.5 Sonnet` | Normal (authenticated) | ✅ Yes |
| `~~GPT-4 Turbo [locked]~~` | Locked (not authenticated) | ❌ No |

### Status Bar

```
[RyCode v1.0] [~/project:main]      [tab Claude 3.5 Sonnet | 💰 $0.12 | tab→]
                                     ^    ^                   ^          ^
                                     |    Current model       |          Hint
                                     |                        Today's cost
                                     Tab cycles models
```

---

## Troubleshooting

### "No credentials found"
- Check environment variables: `env | grep API_KEY`
- Check config file: `~/.config/rycode/credentials`
- Try manual entry with 'a' key

### "Need at least 2 recent models to cycle"
- Use more than one model to build history
- Use Ctrl+X M to select different models

### Auth prompt doesn't appear
- Verify you pressed 'a' on a locked provider (🔒)
- Check terminal size (minimum 80x24)
- Try Esc then retry

### Models won't unlock after auth
- Check API key is valid
- Verify provider is not down (✗)
- Try closing and reopening dialog

---

## Tips & Tricks

### Fastest Authentication
1. Set environment variables for all providers
2. Press 'd' in model dialog once
3. All providers unlock automatically

### Keyboard-Only Workflow
```
Ctrl+X M  → Open model dialog
↓ ↓ ↓     → Navigate to provider
a         → Start auth
(type)    → Enter API key
Enter     → Submit
Esc       → Close dialog
```

### Check Authentication Status Anytime
```
Ctrl+X M → Look for icons → Esc
(Quick peek at which providers are authenticated)
```

---

## Cheat Sheet (Printable)

```
┌─────────────────────────────────────────────────────┐
│ RyCode TUI - Phase 2 Keyboard Shortcuts            │
├─────────────────────────────────────────────────────┤
│                                                     │
│ MODEL DIALOG                                        │
│   a         Authenticate provider                   │
│   d         Auto-detect credentials                 │
│   Enter     Select (or auth if locked)              │
│   Tab       Next model                              │
│   Shift+Tab Previous model                          │
│   Esc       Close dialog                            │
│                                                     │
│ AUTH PROMPT                                         │
│   Enter     Submit API key                          │
│   Ctrl+D    Auto-detect                             │
│   Esc       Cancel                                  │
│                                                     │
│ GLOBAL                                              │
│   Ctrl+X M  Model dialog                            │
│   Ctrl+X A  Agent dialog                            │
│   Ctrl+X L  Session list                            │
│   Ctrl+C    Quit                                    │
│                                                     │
│ ICONS                                               │
│   ✓         Authenticated & healthy                 │
│   ⚠         Authenticated but degraded              │
│   ✗         Authenticated but down                  │
│   🔒        Not authenticated                       │
│                                                     │
└─────────────────────────────────────────────────────┘
```

---

**Need Help?**
- Full documentation: `docs/provider-auth/phase-2/`
- Testing guide: `MANUAL_TESTING_GUIDE.md`
- Implementation details: `INLINE_AUTH_PHASE2_COMPLETE.md`
