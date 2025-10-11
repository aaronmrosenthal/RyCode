# Phase 2 Manual Testing Guide

**Version:** 1.0
**Date:** October 11, 2024
**Build:** `/tmp/rycode-tui-phase2`
**Status:** Ready for testing

---

## Prerequisites

### Required Setup

1. **RyCode Server Running**
   ```bash
   cd packages/rycode
   bun run dev
   # Should start on default port
   ```

2. **TypeScript CLI Accessible**
   ```bash
   # Test CLI works
   cd packages/rycode
   bun run src/auth/cli.ts check anthropic
   # Should return JSON with auth status
   ```

3. **At Least One Unauthenticated Provider**
   - Remove API keys from environment
   - Or test with fresh provider

4. **Terminal Requirements**
   - Minimum width: 80 columns
   - Minimum height: 24 rows
   - Color support recommended

---

## Test Scenarios

### Test 1: Status Bar Display ✅

**Goal:** Verify model and cost display in status bar

**Steps:**
1. Start TUI: `/tmp/rycode-tui-phase2`
2. Observe status bar at bottom right
3. Verify displays: `[tab Model Name | 💰 $0.00 | tab→]`

**Expected Results:**
- ✅ Model name appears correctly
- ✅ Cost shows "$0.00" initially
- ✅ Tab hint shows "tab→"
- ✅ No "No model" message (unless no model selected)
- ✅ Layout aligned properly

**Pass Criteria:**
- Status bar renders without errors
- Model name matches current model
- Cost format is correct

---

### Test 2: Background Cost Updates ✅

**Goal:** Verify cost updates automatically

**Steps:**
1. Start TUI with existing session
2. Send a prompt that costs money
3. Wait 5-10 seconds
4. Observe cost in status bar

**Expected Results:**
- ✅ Cost updates from $0.00 to actual cost
- ✅ Updates happen automatically (every 5s)
- ✅ No UI freezing during update
- ✅ Cost reflects actual API usage

**Pass Criteria:**
- Cost value changes
- UI remains responsive
- No error messages

---

### Test 3: Tab Key Model Cycling ✅

**Goal:** Verify Tab cycles through recent models

**Steps:**
1. Use 3+ different models (build history)
2. Press Tab key
3. Observe model change + toast
4. Press Tab again
5. Press Shift+Tab

**Expected Results:**
- ✅ Tab cycles to next recent model
- ✅ Toast shows "Switched to {Model} ({Provider})"
- ✅ Status bar updates
- ✅ Shift+Tab cycles backward

**Pass Criteria:**
- Model switches on Tab
- Toast notification appears
- Status bar updates immediately

---

### Test 4: View Locked Models 🔒

**Goal:** Verify locked model display

**Steps:**
1. Open model dialog (Ctrl+X M)
2. Find an unauthenticated provider
3. Observe provider header and models

**Expected Results:**
- ✅ Provider header shows "🔒" icon
- ✅ Models show "[locked]" indicator
- ✅ Models are grayed out (faint color)
- ✅ Cannot navigate to locked models
- ✅ Cannot select locked models

**Pass Criteria:**
- Visual indicators clear
- Selection prevented
- No crashes

---

### Test 5: Keyboard Shortcut - 'a' Key 🔐

**Goal:** Verify 'a' key starts authentication

**Steps:**
1. Open model dialog (Ctrl+X M)
2. Navigate to locked provider (e.g., OpenAI 🔒)
3. Press 'a' key

**Expected Results:**
- ✅ Auth prompt dialog appears
- ✅ Shows "Authenticate with {Provider}"
- ✅ Input field visible (masked password)
- ✅ Hints shown: "Press Enter to submit | Ctrl+D for auto-detect | Esc to cancel"

**Pass Criteria:**
- Prompt appears instantly
- Input field focused
- Can type API key (shows •••)

---

### Test 6: Successful Authentication ✓

**Goal:** Verify authentication flow works

**Steps:**
1. Open model dialog (Ctrl+X M)
2. Press 'a' on locked provider
3. Enter valid API key
4. Press Enter
5. Wait for validation (up to 5 seconds)

**Expected Results:**
- ✅ Prompt closes
- ✅ Toast appears: "✓ Authenticated with {Provider} (X models)"
- ✅ Provider header changes: 🔒 → ✓
- ✅ Models unlock (no more [locked])
- ✅ Models become selectable
- ✅ Can now select a model

**Pass Criteria:**
- Auth succeeds within 5 seconds
- Models immediately unlock
- Success toast shows
- No errors

---

### Test 7: Invalid API Key ✗

**Goal:** Verify error handling

**Steps:**
1. Open model dialog (Ctrl+X M)
2. Press 'a' on locked provider
3. Enter invalid/fake API key
4. Press Enter
5. Wait for validation

**Expected Results:**
- ✅ Error message appears in prompt
- ✅ Shows "✗ {Error message}"
- ✅ Prompt stays open (doesn't close)
- ✅ Can retry with correct key
- ✅ Esc still cancels

**Pass Criteria:**
- Error displays clearly
- Can retry without reopening
- No crashes

---

### Test 8: Cancel Authentication (Esc)

**Goal:** Verify cancel works

**Steps:**
1. Open model dialog (Ctrl+X M)
2. Press 'a' on locked provider
3. Start typing API key
4. Press Esc

**Expected Results:**
- ✅ Auth prompt closes
- ✅ Back to model list
- ✅ Provider still locked
- ✅ No error messages
- ✅ Can press 'a' again

**Pass Criteria:**
- Clean exit
- State unchanged
- No memory leaks

---

### Test 9: Select Locked Model → Auth Prompt

**Goal:** Verify locked model selection triggers auth

**Steps:**
1. Open model dialog (Ctrl+X M)
2. Navigate to locked model
3. Press Enter on locked model

**Expected Results:**
- ✅ Auth prompt appears (same as 'a' key)
- ✅ Shows provider name
- ✅ After auth, model selection continues automatically

**Pass Criteria:**
- Seamless flow
- User doesn't need to reselect model
- Auth success → model selected

---

### Test 10: Auto-Detect from Model Dialog 🔍

**Goal:** Verify 'd' key auto-detect

**Prerequisites:**
```bash
# Set test environment variable
export OPENAI_API_KEY="sk-test123..."
```

**Steps:**
1. Open model dialog (Ctrl+X M)
2. Press 'd' key

**Expected Results:**
- ✅ Dialog stays open (doesn't show prompt)
- ✅ Toast appears: "✓ Auto-detected X credential(s)"
- ✅ Locked providers unlock automatically
- ✅ Models become selectable
- ✅ Provider headers update: 🔒 → ✓

**Pass Criteria:**
- Finds credentials
- Multiple providers can unlock at once
- Fast (<1 second)

---

### Test 11: Auto-Detect from Auth Prompt (Ctrl+D)

**Goal:** Verify Ctrl+D in auth prompt

**Prerequisites:**
```bash
export ANTHROPIC_API_KEY="sk-ant-test123..."
```

**Steps:**
1. Open model dialog (Ctrl+X M)
2. Press 'a' on Anthropic
3. Press Ctrl+D (don't type API key)

**Expected Results:**
- ✅ Auth prompt closes
- ✅ Auto-detect runs
- ✅ Toast: "✓ Auto-detected 1 credential(s)"
- ✅ Anthropic unlocks
- ✅ Models available

**Pass Criteria:**
- Works from inside prompt
- Finds env var correctly
- Unlocks automatically

---

### Test 12: Auto-Detect No Credentials

**Goal:** Verify graceful handling when nothing found

**Prerequisites:**
```bash
# Remove all API keys
unset OPENAI_API_KEY
unset ANTHROPIC_API_KEY
unset GOOGLE_API_KEY
# etc.
```

**Steps:**
1. Open model dialog (Ctrl+X M)
2. Press 'd' key

**Expected Results:**
- ✅ Toast: "No credentials found. Please enter manually."
- ✅ Toast type: Info (not error)
- ✅ Model dialog stays open
- ✅ Can press 'a' to enter manually

**Pass Criteria:**
- Helpful message
- Not treated as error
- User can continue

---

### Test 13: Provider Health Indicators

**Goal:** Verify health status icons

**Prerequisites:**
- Need provider with degraded/down status
- Or mock it in auth manager

**Steps:**
1. Open model dialog (Ctrl+X M)
2. Observe provider headers

**Expected Results:**
- ✅ Healthy: "Provider ✓"
- ✅ Degraded: "Provider ⚠"
- ✅ Down: "Provider ✗"
- ✅ Not authenticated: "Provider 🔒"

**Pass Criteria:**
- Icons match health status
- Clear visual distinction
- Cached (30s)

---

### Test 14: Responsive Auth Prompt

**Goal:** Verify prompt adapts to terminal size

**Steps:**
1. Terminal width: 100+ columns
2. Open model dialog, press 'a'
3. Observe input width (should be 60 chars)
4. Resize terminal to 70 columns
5. Press Esc, then 'a' again
6. Observe input width (should be 50 chars)
7. Resize to 50 columns
8. Press Esc, then 'a' again
9. Observe input width (should be 40 chars)

**Expected Results:**
- ✅ Wide terminal: 60-char input
- ✅ Medium terminal: 50-char input
- ✅ Narrow terminal: 40-char input
- ✅ No text overflow
- ✅ Dialog remains usable

**Pass Criteria:**
- Input width adjusts
- No visual glitches
- Always usable

---

### Test 15: Multiple Auth Sessions

**Goal:** Verify can auth multiple providers

**Steps:**
1. Open model dialog (Ctrl+X M)
2. Press 'a' on OpenAI, enter key, auth succeeds
3. Observe OpenAI unlocks
4. Press 'a' on Anthropic, enter key, auth succeeds
5. Observe Anthropic unlocks
6. Press 'a' on Google, enter key, auth succeeds
7. Observe Google unlocks

**Expected Results:**
- ✅ Can auth multiple providers in sequence
- ✅ Each shows success toast
- ✅ All unlock correctly
- ✅ No interference between auths

**Pass Criteria:**
- Multiple auths work
- No state corruption
- All providers functional

---

### Test 16: Auth Cache Invalidation

**Goal:** Verify cache clears after auth

**Steps:**
1. Open model dialog (Ctrl+X M)
2. Observe OpenAI locked
3. Press 'a', authenticate
4. Observe OpenAI unlocked
5. Close dialog (Esc)
6. Reopen dialog (Ctrl+X M)
7. Observe OpenAI still unlocked

**Expected Results:**
- ✅ Auth status persists
- ✅ No re-authentication needed
- ✅ Cache was invalidated correctly
- ✅ Fresh status from server

**Pass Criteria:**
- Auth persists
- Status accurate
- Cache works correctly

---

### Test 17: Stale Cost Indicator

**Goal:** Verify "$--" shows when cost is stale

**Steps:**
1. Start TUI
2. Observe cost in status bar (e.g., "$0.12")
3. Kill TypeScript server (simulate failure)
4. Wait 10+ seconds
5. Observe status bar

**Expected Results:**
- ✅ Cost changes to "💰 $--"
- ✅ Status bar still renders
- ✅ No crash or error UI
- ✅ Other elements intact

**Pass Criteria:**
- Stale indicator shows
- UI remains functional
- Graceful degradation

---

### Test 18: Auth Timeout

**Goal:** Verify 5-second timeout works

**Steps:**
1. Stop TypeScript server
2. Open model dialog (Ctrl+X M)
3. Press 'a' on any provider
4. Enter fake API key
5. Press Enter
6. Wait and observe

**Expected Results:**
- ✅ After ~5 seconds, error appears
- ✅ Error message: "auth CLI error: ..." or timeout message
- ✅ Prompt stays open
- ✅ Can retry or cancel
- ✅ UI doesn't freeze

**Pass Criteria:**
- Times out appropriately
- Error message shown
- No hang

---

### Test 19: Rapid Key Presses

**Goal:** Verify no race conditions

**Steps:**
1. Open model dialog (Ctrl+X M)
2. Rapidly press 'a', 'a', 'a' (multiple times)
3. Press Esc
4. Rapidly press 'd', 'd', 'd'
5. Navigate models and press Enter quickly

**Expected Results:**
- ✅ No crashes
- ✅ No duplicate prompts
- ✅ No orphaned state
- ✅ Behaves predictably

**Pass Criteria:**
- Stable under rapid input
- No visual glitches
- Correct state

---

### Test 20: Edge Case - <2 Recent Models

**Goal:** Verify Tab with insufficient history

**Steps:**
1. Fresh start (clear history)
2. Select one model
3. Press Tab key

**Expected Results:**
- ✅ Toast: "Need at least 2 recent models to cycle"
- ✅ No crash
- ✅ Model stays the same
- ✅ Helpful message

**Pass Criteria:**
- Graceful handling
- Clear feedback
- No errors

---

## Performance Tests

### Test P1: Auth Response Time

**Goal:** Measure authentication speed

**Steps:**
1. Open model dialog
2. Press 'a'
3. Enter API key
4. Press Enter
5. Measure time to toast

**Target:** <2 seconds
**Pass:** Success within timeout
**Fail:** >5 seconds (timeout)

---

### Test P2: Auto-Detect Speed

**Goal:** Measure auto-detect performance

**Steps:**
1. Set 3 API keys in environment
2. Open model dialog
3. Press 'd'
4. Measure time to toast

**Target:** <1 second
**Pass:** Near-instant feedback
**Fail:** >2 seconds

---

### Test P3: Cost Update Latency

**Goal:** Measure background update speed

**Steps:**
1. Start TUI
2. Wait for cost update cycle
3. Measure time from tick to display

**Target:** <100ms per update
**Pass:** Status bar updates smoothly
**Fail:** Visible lag

---

### Test P4: Auth Prompt Display

**Goal:** Measure prompt appearance speed

**Steps:**
1. Open model dialog
2. Press 'a'
3. Measure time to prompt visible

**Target:** <10ms (instant)
**Pass:** No perceptible delay
**Fail:** Noticeable lag

---

## Integration Tests

### Test I1: Status Bar + Model Dialog

**Goal:** Verify status bar updates when model selected via dialog

**Steps:**
1. Note current model in status bar
2. Open model dialog (Ctrl+X M)
3. Select different model
4. Observe status bar

**Expected:**
- ✅ Status bar updates immediately
- ✅ Shows new model name
- ✅ Cost resets or updates

---

### Test I2: Tab Cycling + Status Bar

**Goal:** Verify status bar updates when cycling with Tab

**Steps:**
1. Note current model in status bar
2. Press Tab key
3. Observe status bar

**Expected:**
- ✅ Status bar updates immediately
- ✅ Shows new model name
- ✅ Cost updates (if different session)

---

### Test I3: Auth + Model Selection

**Goal:** Verify seamless locked model → auth → select flow

**Steps:**
1. Open model dialog (Ctrl+X M)
2. Press Enter on locked model
3. Auth prompt appears
4. Enter valid API key
5. Press Enter

**Expected:**
- ✅ Model becomes available
- ✅ Model gets selected automatically (or can be selected)
- ✅ Dialog closes
- ✅ Status bar shows selected model

---

## Regression Tests

### Test R1: Existing Keybindings

**Goal:** Verify no keybinding conflicts

**Steps:**
1. Test all existing shortcuts:
   - Ctrl+X M (model dialog) ✓
   - Ctrl+X A (agent dialog) ✓
   - Ctrl+X L (session list) ✓
   - F2 (cycle model) ✓
   - Esc (close dialogs) ✓

**Expected:**
- ✅ All work as before
- ✅ No conflicts with 'a' or 'd'

---

### Test R2: Theme Switching

**Goal:** Verify auth prompt respects theme

**Steps:**
1. Open model dialog
2. Press 'a'
3. Observe colors
4. Esc, switch theme
5. Press 'a' again
6. Observe colors changed

**Expected:**
- ✅ Colors match theme
- ✅ No hard-coded colors
- ✅ Readable in all themes

---

## Bug Reporting Template

If you find an issue, report it using this format:

```markdown
## Bug: [Short description]

**Severity:** Critical / High / Medium / Low
**Test:** [Test number, e.g., Test 6]

**Steps to Reproduce:**
1. ...
2. ...
3. ...

**Expected Behavior:**
- ...

**Actual Behavior:**
- ...

**Screenshots/Logs:**
[Attach if available]

**Environment:**
- Terminal: [iTerm2, Kitty, etc.]
- Shell: [bash, zsh, etc.]
- OS: [macOS, Linux, etc.]
- Build: /tmp/rycode-tui-phase2
```

---

## Test Results Tracking

Use this checklist to track progress:

### Core Features
- [ ] Test 1: Status Bar Display
- [ ] Test 2: Background Cost Updates
- [ ] Test 3: Tab Key Model Cycling
- [ ] Test 4: View Locked Models
- [ ] Test 5: Keyboard Shortcut 'a'
- [ ] Test 6: Successful Authentication
- [ ] Test 7: Invalid API Key
- [ ] Test 8: Cancel Authentication
- [ ] Test 9: Select Locked Model
- [ ] Test 10: Auto-Detect (Model Dialog)
- [ ] Test 11: Auto-Detect (Auth Prompt)
- [ ] Test 12: Auto-Detect No Credentials
- [ ] Test 13: Provider Health Indicators

### Edge Cases
- [ ] Test 14: Responsive Auth Prompt
- [ ] Test 15: Multiple Auth Sessions
- [ ] Test 16: Auth Cache Invalidation
- [ ] Test 17: Stale Cost Indicator
- [ ] Test 18: Auth Timeout
- [ ] Test 19: Rapid Key Presses
- [ ] Test 20: <2 Recent Models

### Performance
- [ ] Test P1: Auth Response Time
- [ ] Test P2: Auto-Detect Speed
- [ ] Test P3: Cost Update Latency
- [ ] Test P4: Auth Prompt Display

### Integration
- [ ] Test I1: Status Bar + Model Dialog
- [ ] Test I2: Tab Cycling + Status Bar
- [ ] Test I3: Auth + Model Selection

### Regression
- [ ] Test R1: Existing Keybindings
- [ ] Test R2: Theme Switching

---

## Success Criteria

**Phase 2 is ready for production when:**

✅ All core feature tests pass (Tests 1-13)
✅ No critical bugs found
✅ Edge cases handled gracefully (Tests 14-20)
✅ Performance targets met (Tests P1-P4)
✅ No regressions (Tests R1-R2)
✅ Integration works smoothly (Tests I1-I3)

---

**Happy Testing!** 🧪

Report results in a new document: `MANUAL_TEST_RESULTS.md`
