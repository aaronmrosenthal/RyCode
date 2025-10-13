# TUI Peer Review - Multi-Agent Analysis

## Executive Summary
Comprehensive review of TUI cursor, splash, and model selector improvements from 4 specialized perspectives.

---

## 🏗️ **ARCHITECT REVIEW**

### Structural Assessment

**Strengths:**
- ✅ Clean separation: splash → session → chat flow
- ✅ Cursor positioning abstracted across layers (textarea → editor → tui)
- ✅ Fallback pattern for SOTA models (API first, curated fallback)

**Concerns:**
- ⚠️ **CRITICAL: Cursor offset stacking issue**
  - Line editor.go:437: `+2` for prompt
  - Line tui.go:1047: `+editorX` for layout
  - Line tui.go:1207/1259: Returns `editorX + 2`
  - **Result:** Potential double-counting leading to misalignment

- ⚠️ **Model data source fragility**
  - Hardcoded SOTA models in models.go:553-640
  - No sync mechanism with actual provider APIs
  - Model IDs may become stale (e.g., "claude-4-5-sonnet-20250929")

- ⚠️ **Splash timing hardcoded**
  - 4.5s duration may feel long/short depending on system
  - No adaptive timing based on terminal speed

**Recommendations:**
1. Create cursor positioning tests to verify offset calculations
2. Consider dynamic model fetching from actual provider APIs
3. Add telemetry to measure optimal splash duration

---

## 👨‍💻 **SENIOR ENGINEER REVIEW**

### Code Quality Analysis

**Good Practices:**
- ✅ Descriptive comments explaining cursor offset logic
- ✅ Color constants named semantically (brightCyan, neonGreen)
- ✅ Proper error handling in session creation
- ✅ Immutable state updates in Bubble Tea pattern

**Code Smells:**
- 🔴 **Magic numbers everywhere**
  - `+2`, `+5`, `+3` offsets without clear rationale
  - Should be named constants: `PROMPT_WIDTH`, `BORDER_WIDTH`, etc.

- 🔴 **Cursor blink still an issue**
  - Set `Blink: false` in 3 places but user reports still blinking
  - Terminal may override - need DECSET escape codes
  - Missing: `\x1b[?12l` to explicitly disable terminal blink

- 🟡 **Placeholder cursor logic duplicated**
  - Virtual vs Real cursor paths in placeholderView()
  - Could be extracted to helper function

**Bug Risks:**
- Layout offset calculation is fragile - any change to border/padding breaks cursor
- No validation that SOTA model IDs actually exist in backends
- Splash colors not tested in light themes

**Recommendations:**
1. Extract constants for all offset values
2. Add integration test that verifies cursor position
3. Send explicit DECSET codes for cursor control
4. Add unit tests for offset calculations

---

## 📦 **PRODUCT OWNER REVIEW**

### User Experience Assessment

**Wins:**
- ✅ 4.5s splash creates anticipation without annoyance
- ✅ Matrix green/cyan matches user's terminal aesthetic
- ✅ SOTA models (Claude 4.5, Gemini 3.0) show cutting-edge positioning
- ✅ Auto-session creation removes friction

**User Pain Points:**
- 🔴 **BLOCKER: Cursor still blinking** (user explicitly complained)
  - This breaks the "polish" promise
  - Makes UI feel unfinished

- 🔴 **BLOCKER: Cursor positioning still wrong** (user screenshot)
  - Cursor appears in placeholder text, not at start
  - Breaks typing experience

- 🟡 **Input box height** - user wanted "exact" match to reference
  - Current implementation may still have extra pixels

- 🟡 **Model selector shows "locked" for unauthenticated**
  - Good UX but may frustrate if auth is complex
  - Need smooth auth flow

**Missing Features:**
- No way to skip splash after first time (always shows)
- No visual feedback during 4.5s splash (just animation)
- Model selector doesn't show which model is currently selected

**Recommendations:**
1. **P0:** Fix cursor blinking (terminal escape codes)
2. **P0:** Fix cursor positioning (verify offset math)
3. **P1:** Add "Initializing..." text to splash
4. **P2:** Show current model with ● indicator in selector
5. **P3:** Add splash skip option after first view

---

## 🔒 **SECURITY SPECIALIST REVIEW**

### Security Assessment

**Concerns:**
- 🟡 **API key handling in model selector**
  - Line models.go:340: `apiKey` passed to authentication
  - Need to verify secure transmission
  - Should be masked in logs

- 🟡 **Provider authentication**
  - Auto-detect scans for API keys
  - Could expose keys in error messages
  - Need audit of auth bridge error handling

**Good Practices:**
- ✅ Context timeouts prevent hanging auth calls
- ✅ Authentication status cached (30s TTL)
- ✅ No API keys in source code

**Recommendations:**
1. Audit auth bridge for key exposure in logs
2. Add rate limiting to auth attempts
3. Sanitize all error messages before display

---

## 🎯 **CRITICAL ACTION ITEMS**

### Must Fix Before Ship:

1. **Cursor Blinking** - P0 BLOCKER
   ```go
   // Add to tui Init():
   fmt.Print("\x1b[?12l")  // Disable cursor blinking (DECSET)
   ```

2. **Cursor Position** - P0 BLOCKER
   - Current: textareaOffset(1) + editorOffset(2) + layoutOffset(editorX + 2)
   - Debug actual values being returned
   - Likely fix: Remove one layer of offsetting

3. **Model IDs Validation** - P1
   - Verify SOTA model IDs match actual provider APIs
   - Add fallback if model ID not found

4. **Offset Constants** - P1
   ```go
   const (
       TEXTAREA_PROMPT_WIDTH = 1
       EXTERNAL_PROMPT_WIDTH = 2
       BORDER_WIDTH = 1
       // etc.
   )
   ```

### Nice to Have:

5. **Splash Skip State** - P2
   - Add `State.HasSeenSplash` flag
   - Show only on first run or with `--splash` flag

6. **Current Model Indicator** - P2
   - Add ● or ✓ next to active model in selector

---

## 📊 **METRICS & TESTING**

### Test Coverage Gaps:
- ❌ No cursor positioning tests
- ❌ No splash animation tests
- ❌ No model selector integration tests
- ✅ Created manual verification checklist

### Performance Concerns:
- Splash animation: 90 FPS (50ms ticks) - may be overkill
- Model selector auth check: 1s timeout per provider (could batch)
- No lazy loading of provider lists

---

## 💡 **INNOVATION SCORE: 8/10**

**What Works:**
- Matrix color cascade is genuinely impressive
- 4-phase timing creates professional feel
- SOTA model curation shows product maturity

**What Could Be Better:**
- Cursor issues undermine the polish
- Hardcoded models limit flexibility
- No personalization/customization options

---

## ✅ **APPROVAL STATUS**

- Architecture: ⚠️ **CONDITIONAL** - Fix offset stacking
- Code Quality: ⚠️ **CONDITIONAL** - Add constants, fix cursor
- Product/UX: ❌ **BLOCKED** - Cursor must work correctly
- Security: ✅ **APPROVED** - Minor improvements needed

**Overall: BLOCKED on cursor fixes. Once cursor works, this is production-ready.**

---

## 🚀 **NEXT STEPS**

1. Add cursor blink disable escape code
2. Debug and fix cursor X offset calculation
3. Extract magic numbers to constants
4. Add automated tests for cursor positioning
5. Verify SOTA model IDs against actual APIs
6. Add metrics/telemetry to track splash engagement

**Timeline:** 2-4 hours to address P0 blockers, ship-ready after that.
