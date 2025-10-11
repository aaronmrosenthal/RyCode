# RyCode Splash Screen - Implementation Tasks
## Actionable Task Breakdown with Dependencies

> **Generated from:** SPLASH_IMPLEMENTATION_PLAN.md
> **Total Duration:** 5 weeks (35 days)
> **Total Tasks:** 87 tasks across 5 phases
> **Status:** Ready for execution

---

## Task Priority System

- 🔴 **P0 (Critical):** Blocking tasks - must complete before dependent tasks
- 🟠 **P1 (High):** Core functionality - required for MVP
- 🟡 **P2 (Medium):** Important features - enhance user experience
- 🟢 **P3 (Low):** Nice-to-have - polish and extras

---

## Week 1: Foundation (Core Engine) - 18 Tasks

### Phase 1.1: Project Setup (Day 1) - 4 Tasks

**TASK-001** 🔴 P0 - Create splash module directory structure
- **Duration:** 30 min
- **Dependencies:** None
- **Assignee:** Developer
- **Acceptance Criteria:**
  ```bash
  mkdir -p packages/tui/internal/splash
  touch packages/tui/internal/splash/splash.go
  touch packages/tui/internal/splash/cortex.go
  touch packages/tui/internal/splash/bootsequence.go
  touch packages/tui/internal/splash/closer.go
  touch packages/tui/internal/splash/ansi.go
  touch packages/tui/internal/splash/config.go
  ```
- **Verification:** Directory exists with 6 empty files

**TASK-002** 🔴 P0 - Initialize package with basic types
- **Duration:** 1 hour
- **Dependencies:** TASK-001
- **Files:** `splash.go`
- **Implementation:**
  ```go
  package splash

  import tea "github.com/charmbracelet/bubbletea"
  import "time"

  type Model struct {
      act       int           // 1=boot, 2=cortex, 3=closer
      frame     int
      renderer  *CortexRenderer
      done      bool
      width     int
      height    int
  }

  type tickMsg time.Time
  ```
- **Verification:** `go build` succeeds

**TASK-003** 🟠 P1 - Create test files
- **Duration:** 30 min
- **Dependencies:** TASK-001
- **Files:** Create `*_test.go` for each module
  ```bash
  touch packages/tui/internal/splash/cortex_test.go
  touch packages/tui/internal/splash/ansi_test.go
  touch packages/tui/internal/splash/bootsequence_test.go
  ```
- **Verification:** `go test ./internal/splash` runs (even with no tests)

**TASK-004** 🟡 P2 - Set up benchmark framework
- **Duration:** 30 min
- **Dependencies:** TASK-003
- **Files:** `cortex_test.go`
- **Implementation:**
  ```go
  func BenchmarkRenderFrame(b *testing.B) {
      r := NewCortexRenderer(80, 24)
      b.ResetTimer()
      for i := 0; i < b.N; i++ {
          r.RenderFrame()
      }
  }
  ```
- **Verification:** `go test -bench=.` runs

---

### Phase 1.2: Torus Mathematics (Day 1-2) - 6 Tasks

**TASK-005** 🔴 P0 - Implement torus parametric equations
- **Duration:** 3 hours
- **Dependencies:** TASK-002
- **Files:** `cortex.go`
- **Implementation:**
  ```go
  type CortexRenderer struct {
      width, height int
      A, B         float64  // Rotation angles
      screen       []rune
      zbuffer      []float64
  }

  func NewCortexRenderer(width, height int) *CortexRenderer {
      size := width * height
      return &CortexRenderer{
          width:   width,
          height:  height,
          screen:  make([]rune, size),
          zbuffer: make([]float64, size),
      }
  }

  // Calculate point on torus surface
  func (r *CortexRenderer) torusPoint(theta, phi float64) (x, y, z float64) {
      const R = 2.0  // Major radius
      const r = 1.0  // Minor radius

      sinTheta, cosTheta := math.Sin(theta), math.Cos(theta)
      sinPhi, cosPhi := math.Sin(phi), math.Cos(phi)

      circleX := R + r*cosPhi
      circleY := r * sinPhi

      // Before rotation
      x = circleX * cosTheta
      y = circleX * sinTheta
      z = circleY

      return x, y, z
  }
  ```
- **Acceptance Criteria:**
  - Function returns valid 3D coordinates
  - Distance from origin is between 1.0 and 3.0
- **Verification:** Unit test passes

**TASK-006** 🔴 P0 - Implement rotation matrices
- **Duration:** 2 hours
- **Dependencies:** TASK-005
- **Files:** `cortex.go`
- **Implementation:**
  ```go
  // Apply rotation matrices Rx(A) and Rz(B)
  func (r *CortexRenderer) rotate(x, y, z float64) (float64, float64, float64) {
      sinA, cosA := math.Sin(r.A), math.Cos(r.A)
      sinB, cosB := math.Sin(r.B), math.Cos(r.B)

      // Rotate around X-axis
      y1 := y*cosA - z*sinA
      z1 := y*sinA + z*cosA

      // Rotate around Z-axis
      x2 := x*cosB - y1*sinB
      y2 := x*sinB + y1*cosB
      z2 := z1

      return x2, y2, z2
  }
  ```
- **Acceptance Criteria:** Rotation preserves distance from origin
- **Verification:** Unit test with known rotations

**TASK-007** 🔴 P0 - Implement perspective projection
- **Duration:** 2 hours
- **Dependencies:** TASK-006
- **Files:** `cortex.go`
- **Implementation:**
  ```go
  // Project 3D point to 2D screen
  func (r *CortexRenderer) project(x, y, z float64) (int, int, float64) {
      // Move away from camera
      z = z + 5.0

      // Perspective projection
      ooz := 1.0 / z  // "one over z"
      xp := int(float64(r.width)*0.5 + 30.0*ooz*x)
      yp := int(float64(r.height)*0.5 - 15.0*ooz*y)

      return xp, yp, ooz
  }
  ```
- **Acceptance Criteria:** Points map to screen coordinates
- **Verification:** Visual inspection (render single point)

**TASK-008** 🟠 P1 - Implement luminance calculation
- **Duration:** 2 hours
- **Dependencies:** TASK-006
- **Files:** `cortex.go`
- **Implementation:**
  ```go
  // Calculate surface luminance (Phong shading)
  func (r *CortexRenderer) luminance(theta, phi float64) float64 {
      sinA, cosA := math.Sin(r.A), math.Cos(r.A)
      sinB, cosB := math.Sin(r.B), math.Cos(r.B)
      sinTheta, cosTheta := math.Sin(theta), math.Cos(theta)
      sinPhi, cosPhi := math.Sin(phi), math.Cos(phi)

      // Light direction calculation
      L := cosPhi*cosTheta*sinB - cosA*cosTheta*sinPhi -
           sinA*sinTheta + cosB*(cosA*sinPhi - cosTheta*sinA*sinTheta)

      return L
  }
  ```
- **Acceptance Criteria:** Luminance in range [-1, 1]
- **Verification:** Unit test

**TASK-009** 🟠 P1 - Implement Z-buffer algorithm
- **Duration:** 2 hours
- **Dependencies:** TASK-007
- **Files:** `cortex.go`
- **Implementation:**
  ```go
  func (r *CortexRenderer) plotPoint(x, y int, depth float64, char rune) {
      if x < 0 || x >= r.width || y < 0 || y >= r.height {
          return
      }

      idx := y*r.width + x

      // Z-buffer test
      if depth > r.zbuffer[idx] {
          r.zbuffer[idx] = depth
          r.screen[idx] = char
      }
  }
  ```
- **Acceptance Criteria:** Front surfaces occlude back surfaces
- **Verification:** Visual test with sphere

**TASK-010** 🟠 P1 - Implement complete render loop
- **Duration:** 3 hours
- **Dependencies:** TASK-005, TASK-006, TASK-007, TASK-008, TASK-009
- **Files:** `cortex.go`
- **Implementation:**
  ```go
  func (r *CortexRenderer) RenderFrame() {
      // Clear buffers
      for i := range r.screen {
          r.screen[i] = ' '
          r.zbuffer[i] = 0
      }

      // Render torus
      const thetaStep = 0.07
      const phiStep = 0.02

      for theta := 0.0; theta < 6.28; theta += thetaStep {
          for phi := 0.0; phi < 6.28; phi += phiStep {
              // Calculate 3D point
              x, y, z := r.torusPoint(theta, phi)

              // Apply rotation
              x, y, z = r.rotate(x, y, z)

              // Project to 2D
              xp, yp, depth := r.project(x, y, z)

              // Calculate luminance
              L := r.luminance(theta, phi)

              // Map to character
              lumIdx := int((L + 1.0) * 3.5)
              if lumIdx < 0 { lumIdx = 0 }
              if lumIdx > 7 { lumIdx = 7 }
              chars := []rune{' ', '.', '·', ':', '*', '◉', '◎', '⚡'}

              // Plot with z-buffer
              r.plotPoint(xp, yp, depth, chars[lumIdx])
          }
      }

      // Update rotation
      r.A += 0.04
      r.B += 0.02
  }

  func (r *CortexRenderer) String() string {
      var buf strings.Builder
      for y := 0; y < r.height; y++ {
          for x := 0; x < r.width; x++ {
              buf.WriteRune(r.screen[y*r.width+x])
          }
          if y < r.height-1 {
              buf.WriteRune('\n')
          }
      }
      return buf.String()
  }
  ```
- **Acceptance Criteria:** Full torus renders correctly
- **Verification:** Visual inspection - should see rotating donut

---

### Phase 1.3: Color System (Day 3) - 5 Tasks

**TASK-011** 🟠 P1 - Implement ANSI color utilities
- **Duration:** 2 hours
- **Dependencies:** TASK-002
- **Files:** `ansi.go`
- **Implementation:**
  ```go
  package splash

  import "fmt"

  // RGB color
  type RGB struct {
      R, G, B uint8
  }

  // Convert RGB to ANSI truecolor escape sequence
  func (c RGB) ANSI() string {
      return fmt.Sprintf("\033[38;2;%d;%d;%dm", c.R, c.G, c.B)
  }

  // Reset color
  func ResetColor() string {
      return "\033[0m"
  }

  // Colorize text
  func Colorize(text string, color RGB) string {
      return color.ANSI() + text + ResetColor()
  }
  ```
- **Verification:** Print colored text to terminal

**TASK-012** 🟠 P1 - Implement cyan-to-magenta gradient
- **Duration:** 2 hours
- **Dependencies:** TASK-011
- **Files:** `ansi.go`
- **Implementation:**
  ```go
  // Interpolate between two colors
  func lerp(a, b uint8, t float64) uint8 {
      return uint8(float64(a)*(1.0-t) + float64(b)*t)
  }

  // Cyan to magenta gradient based on angle (0 to 2π)
  func GradientColor(angle float64) RGB {
      cyan := RGB{0, 255, 255}      // #00FFFF
      magenta := RGB{255, 0, 255}   // #FF00FF

      // Normalize angle to [0, 1]
      t := math.Mod(angle, 2*math.Pi) / (2 * math.Pi)

      return RGB{
          R: lerp(cyan.R, magenta.R, t),
          G: lerp(cyan.G, magenta.G, t),
          B: lerp(cyan.B, magenta.B, t),
      }
  }
  ```
- **Verification:** Unit test with known angles

**TASK-013** 🟡 P2 - Implement 256-color fallback
- **Duration:** 2 hours
- **Dependencies:** TASK-011
- **Files:** `ansi.go`
- **Implementation:**
  ```go
  // Convert RGB to nearest 256-color ANSI code
  func (c RGB) ANSI256() string {
      // 256-color cube: 16 + 36*r + 6*g + b
      r := int(c.R) * 6 / 256
      g := int(c.G) * 6 / 256
      b := int(c.B) * 6 / 256
      code := 16 + 36*r + 6*g + b
      return fmt.Sprintf("\033[38;5;%dm", code)
  }
  ```
- **Verification:** Compare truecolor vs 256-color side-by-side

**TASK-014** 🟡 P2 - Detect terminal color capabilities
- **Duration:** 1 hour
- **Dependencies:** TASK-011
- **Files:** `ansi.go`
- **Implementation:**
  ```go
  type ColorMode int

  const (
      Colors16 ColorMode = iota
      Colors256
      Truecolor
  )

  func DetectColorMode() ColorMode {
      colorterm := os.Getenv("COLORTERM")
      if colorterm == "truecolor" || colorterm == "24bit" {
          return Truecolor
      }

      term := os.Getenv("TERM")
      if strings.Contains(term, "256color") {
          return Colors256
      }

      return Colors16
  }
  ```
- **Verification:** Test on different terminals

**TASK-015** 🟠 P1 - Integrate colors into cortex renderer
- **Duration:** 2 hours
- **Dependencies:** TASK-010, TASK-012
- **Files:** `cortex.go`
- **Implementation:**
  ```go
  func (r *CortexRenderer) Render() string {
      r.RenderFrame()

      var buf strings.Builder
      for y := 0; y < r.height; y++ {
          for x := 0; x < r.width; x++ {
              idx := y*r.width + x
              char := r.screen[idx]

              if char != ' ' {
                  // Color based on position (angle around torus)
                  angle := math.Atan2(float64(y-r.height/2), float64(x-r.width/2))
                  color := GradientColor(angle + r.B)  // Rotate with torus
                  buf.WriteString(Colorize(string(char), color))
              } else {
                  buf.WriteRune(' ')
              }
          }
          if y < r.height-1 {
              buf.WriteRune('\n')
          }
      }
      return buf.String()
  }
  ```
- **Acceptance Criteria:** Torus displays with cyan-magenta gradient
- **Verification:** Visual inspection

---

### Phase 1.4: Testing & Optimization (Day 4-5) - 3 Tasks

**TASK-016** 🔴 P0 - Write unit tests for torus math
- **Duration:** 3 hours
- **Dependencies:** TASK-010
- **Files:** `cortex_test.go`
- **Implementation:**
  ```go
  func TestTorusGeometry(t *testing.T) {
      r := NewCortexRenderer(80, 24)

      for theta := 0.0; theta < 6.28; theta += 0.1 {
          for phi := 0.0; phi < 6.28; phi += 0.1 {
              x, y, z := r.torusPoint(theta, phi)

              // Distance should be R ± r (2 ± 1 = 1 to 3)
              dist := math.Sqrt(x*x + y*y + z*z)
              if dist < 0.5 || dist > 3.5 {
                  t.Errorf("Invalid distance: %.2f", dist)
              }
          }
      }
  }

  func TestZBufferOcclusion(t *testing.T) {
      r := NewCortexRenderer(80, 24)
      r.RenderFrame()

      // Center should have depth (not empty)
      centerIdx := (r.height/2)*r.width + (r.width/2)
      if r.zbuffer[centerIdx] == 0 {
          t.Error("Z-buffer not working - center is empty")
      }
  }

  func TestRotationPreservesDistance(t *testing.T) {
      r := NewCortexRenderer(80, 24)
      x, y, z := 1.0, 2.0, 3.0
      distBefore := math.Sqrt(x*x + y*y + z*z)

      x2, y2, z2 := r.rotate(x, y, z)
      distAfter := math.Sqrt(x2*x2 + y2*y2 + z2*z2)

      if math.Abs(distBefore-distAfter) > 0.001 {
          t.Errorf("Rotation changed distance: %.3f -> %.3f", distBefore, distAfter)
      }
  }
  ```
- **Acceptance Criteria:** All tests pass
- **Verification:** `go test -v ./internal/splash`

**TASK-017** 🟠 P1 - Performance benchmarking
- **Duration:** 2 hours
- **Dependencies:** TASK-010
- **Files:** `cortex_test.go`
- **Implementation:**
  ```go
  func BenchmarkRenderFrame(b *testing.B) {
      r := NewCortexRenderer(80, 24)
      b.ResetTimer()
      for i := 0; i < b.N; i++ {
          r.RenderFrame()
      }
      // Target: <1ms (1,000,000 ns)
  }

  func BenchmarkFullRender(b *testing.B) {
      r := NewCortexRenderer(80, 24)
      b.ResetTimer()
      for i := 0; i < b.N; i++ {
          _ = r.Render()  // Include string building + colors
      }
      // Target: <10ms (terminal I/O is slow)
  }
  ```
- **Acceptance Criteria:** <1ms per RenderFrame()
- **Verification:** `go test -bench=. -benchtime=5s`

**TASK-018** 🟡 P2 - Profile and optimize hotspots
- **Duration:** 3 hours
- **Dependencies:** TASK-017
- **Implementation:**
  ```bash
  # Generate CPU profile
  go test -cpuprofile=cpu.prof -bench=BenchmarkRenderFrame

  # Analyze
  go tool pprof cpu.prof
  # (pprof) top10
  # (pprof) list RenderFrame

  # Optimize based on findings
  # Common optimizations:
  # - Precompute sin/cos tables
  # - Reduce allocations (use strings.Builder)
  # - Inline small functions
  ```
- **Acceptance Criteria:** No single function >20% CPU time
- **Verification:** pprof output

---

## Week 2: Animations (3-Act Sequence) - 15 Tasks

### Phase 2.1: Boot Sequence (Day 1-2) - 5 Tasks

**TASK-019** 🟠 P1 - Create model status data structure
- **Duration:** 1 hour
- **Dependencies:** TASK-002
- **Files:** `bootsequence.go`
- **Implementation:**
  ```go
  package splash

  type ModelInfo struct {
      Name       string
      Role       string
      Icon       string
      Color      RGB
      Delay      time.Duration
  }

  var models = []ModelInfo{
      {"Claude", "Logical Reasoning", "🧩", RGB{10, 255, 10}, 100 * time.Millisecond},
      {"Gemini", "System Architecture", "⚙️", RGB{10, 255, 10}, 100 * time.Millisecond},
      {"Codex", "Code Generation", "💻", RGB{10, 255, 10}, 100 * time.Millisecond},
      {"Qwen", "Research Pipeline", "🔎", RGB{10, 255, 10}, 100 * time.Millisecond},
      {"Grok", "Humor & Chaos Engine", "🤖", RGB{10, 255, 10}, 100 * time.Millisecond},
      {"GPT", "Language Core", "✅", RGB{10, 255, 10}, 100 * time.Millisecond},
  }
  ```
- **Verification:** Data structure compiles

**TASK-020** 🟠 P1 - Implement line-by-line reveal animation
- **Duration:** 2 hours
- **Dependencies:** TASK-019
- **Files:** `bootsequence.go`
- **Implementation:**
  ```go
  type BootSequence struct {
      frame       int
      linesShown  int
  }

  func NewBootSequence() *BootSequence {
      return &BootSequence{}
  }

  func (b *BootSequence) Update(frame int) {
      // Show 1 line every 3 frames (100ms at 30 FPS)
      b.linesShown = frame / 3
      if b.linesShown > len(models) {
          b.linesShown = len(models)
      }
  }

  func (b *BootSequence) Render() string {
      var buf strings.Builder

      buf.WriteString(Colorize("> [RYCODE NEURAL CORTEX v1.0.0]\n", RGB{0, 255, 255}))
      buf.WriteString(">\n")

      for i := 0; i < b.linesShown && i < len(models); i++ {
          model := models[i]

          prefix := "├─"
          if i == len(models)-1 {
              prefix = "└─"
          }

          line := fmt.Sprintf("> %s %s ▸ %s: ONLINE %s\n",
              prefix, model.Name, model.Role, model.Icon)

          buf.WriteString(Colorize(line, model.Color))
      }

      // Final message after all models loaded
      if b.linesShown >= len(models) {
          buf.WriteString(">\n")
          buf.WriteString(Colorize("> ⚡ SIX MINDS. ONE COMMAND LINE.\n", RGB{255, 174, 0}))
      }

      return buf.String()
  }
  ```
- **Verification:** Visual test - lines appear one by one

**TASK-021** 🟡 P2 - Add typing effect for dramatic flair
- **Duration:** 2 hours
- **Dependencies:** TASK-020
- **Files:** `bootsequence.go`
- **Implementation:**
  ```go
  func (b *BootSequence) Render() string {
      var buf strings.Builder

      // ... existing code ...

      // For current line being revealed, show partial text
      if b.frame%3 < 3 && b.linesShown < len(models) {
          model := models[b.linesShown]
          partialFrame := b.frame % 3

          prefix := "├─"
          if b.linesShown == len(models)-1 {
              prefix = "└─"
          }

          fullLine := fmt.Sprintf("> %s %s ▸ %s: ONLINE %s",
              prefix, model.Name, model.Role, model.Icon)

          // Show partial text based on frame
          charsToShow := len(fullLine) * partialFrame / 3
          partialLine := fullLine[:charsToShow]

          buf.WriteString(Colorize(partialLine, RGB{100, 100, 100}))  // Dimmed
      }

      return buf.String()
  }
  ```
- **Verification:** Lines "type" onto screen

**TASK-022** 🟢 P3 - Add ASCII art header (optional)
- **Duration:** 1 hour
- **Dependencies:** TASK-020
- **Files:** `bootsequence.go`
- **Implementation:**
  ```go
  const asciiHeader = `
   ____       ____          _
  |  _ \ _   / ___|___   __| | ___
  | |_) | | | |   / _ \ / _' |/ _ \
  |  _ <| |_| |__| (_) | (_| |  __/
  |_| \_\\__, \____\___/ \__,_|\___|
         |___/
  `

  func (b *BootSequence) Render() string {
      var buf strings.Builder

      if b.frame < 10 {
          buf.WriteString(Colorize(asciiHeader, RGB{0, 255, 255}))
      }

      // ... rest of code ...
  }
  ```
- **Verification:** Header displays briefly

**TASK-023** 🟠 P1 - Write unit tests for boot sequence
- **Duration:** 1 hour
- **Dependencies:** TASK-020
- **Files:** `bootsequence_test.go`
- **Implementation:**
  ```go
  func TestBootSequenceProgression(t *testing.T) {
      bs := NewBootSequence()

      // Frame 0: No models shown
      bs.Update(0)
      if bs.linesShown != 0 {
          t.Errorf("Expected 0 lines, got %d", bs.linesShown)
      }

      // Frame 3: 1 model shown
      bs.Update(3)
      if bs.linesShown != 1 {
          t.Errorf("Expected 1 line, got %d", bs.linesShown)
      }

      // Frame 18: All 6 models shown
      bs.Update(18)
      if bs.linesShown != 6 {
          t.Errorf("Expected 6 lines, got %d", bs.linesShown)
      }
  }
  ```
- **Verification:** Tests pass

---

### Phase 2.2: Closer Screen (Day 3-4) - 5 Tasks

**TASK-024** 🟠 P1 - Design closer screen layout
- **Duration:** 1 hour
- **Dependencies:** TASK-002
- **Files:** `closer.go`
- **Implementation:**
  ```go
  package splash

  const closerText = `
  ╔═══════════════════════════════════════════════════════════════════════╗
  ║                                                                       ║
  ║                  🌀 RYCODE NEURAL CORTEX ACTIVE 🌀                   ║
  ║                                                                       ║
  ║         "Every LLM fused. Every edge case covered.                   ║
  ║          You're not just coding—you're orchestrating                 ║
  ║          intelligence."                                              ║
  ║                                                                       ║
  ║                                                                       ║
  ║                   Press any key to begin...                          ║
  ║                                                                       ║
  ╚═══════════════════════════════════════════════════════════════════════╝
  `
  ```
- **Verification:** Manually check alignment

**TASK-025** 🟠 P1 - Implement centered text rendering
- **Duration:** 2 hours
- **Dependencies:** TASK-024
- **Files:** `closer.go`
- **Implementation:**
  ```go
  type Closer struct {
      width  int
      height int
  }

  func NewCloser(width, height int) *Closer {
      return &Closer{width: width, height: height}
  }

  func (c *Closer) Render() string {
      lines := strings.Split(closerText, "\n")

      // Calculate vertical centering
      startY := (c.height - len(lines)) / 2
      if startY < 0 {
          startY = 0
      }

      var buf strings.Builder

      // Add top padding
      for i := 0; i < startY; i++ {
          buf.WriteRune('\n')
      }

      // Render centered lines
      for _, line := range lines {
          // Horizontal centering
          padding := (c.width - len(line)) / 2
          if padding > 0 {
              buf.WriteString(strings.Repeat(" ", padding))
          }
          buf.WriteString(line)
          buf.WriteRune('\n')
      }

      return buf.String()
  }
  ```
- **Verification:** Test on different terminal sizes

**TASK-026** 🟡 P2 - Add subtle color pulse effect
- **Duration:** 2 hours
- **Dependencies:** TASK-025
- **Files:** `closer.go`
- **Implementation:**
  ```go
  func (c *Closer) RenderWithPulse(frame int) string {
      // Pulse intensity: 0.7 to 1.0
      intensity := 0.7 + 0.3*math.Sin(float64(frame)*0.1)

      cyan := RGB{
          R: uint8(float64(0) * intensity),
          G: uint8(float64(255) * intensity),
          B: uint8(float64(255) * intensity),
      }

      lines := strings.Split(closerText, "\n")
      var buf strings.Builder

      for _, line := range lines {
          if strings.Contains(line, "🌀") || strings.Contains(line, "CORTEX") {
              buf.WriteString(Colorize(line, cyan))
          } else {
              buf.WriteString(line)
          }
          buf.WriteRune('\n')
      }

      return buf.String()
  }
  ```
- **Verification:** Visual test - title should pulse

**TASK-027** 🟡 P2 - Responsive layout for different terminal sizes
- **Duration:** 2 hours
- **Dependencies:** TASK-025
- **Files:** `closer.go`
- **Implementation:**
  ```go
  func (c *Closer) Render() string {
      if c.width < 80 || c.height < 24 {
          // Simplified version for small terminals
          return c.renderCompact()
      }

      return c.renderFull()
  }

  func (c *Closer) renderCompact() string {
      return `
  ╔═══════════════════════════════════╗
  ║   🌀 RYCODE CORTEX ACTIVE 🌀     ║
  ║                                   ║
  ║   Six minds. One command line.    ║
  ║                                   ║
  ║   Press any key...                ║
  ╚═══════════════════════════════════╝
  `
  }
  ```
- **Verification:** Test on small terminal (60×20)

**TASK-028** 🟢 P3 - Add easter egg messages (random quotes)
- **Duration:** 1 hour
- **Dependencies:** TASK-025
- **Files:** `closer.go`
- **Implementation:**
  ```go
  var easterEggQuotes = []string{
      "Built by Claude. No humans were harmed.",
      "100% AI-designed. 0% compromises.",
      "The age of typing is over.",
      "Warning: May cause productivity addiction.",
  }

  func (c *Closer) Render() string {
      // 5% chance of easter egg
      if rand.Float64() < 0.05 {
          quote := easterEggQuotes[rand.Intn(len(easterEggQuotes))]
          // ... inject quote into render ...
      }

      return c.renderFull()
  }
  ```
- **Verification:** Run 100 times, should see ~5 easter eggs

---

### Phase 2.3: Orchestrator (Day 5-6) - 5 Tasks

**TASK-029** 🔴 P0 - Implement Bubble Tea Model structure
- **Duration:** 2 hours
- **Dependencies:** TASK-002, TASK-020, TASK-025
- **Files:** `splash.go`
- **Implementation:**
  ```go
  type Model struct {
      act          int  // 1=boot, 2=cortex, 3=closer
      frame        int
      bootSeq      *BootSequence
      cortex       *CortexRenderer
      closer       *Closer
      done         bool
      width        int
      height       int
  }

  func New() Model {
      return Model{
          act:     1,
          frame:   0,
          bootSeq: NewBootSequence(),
          cortex:  NewCortexRenderer(80, 24),
          closer:  NewCloser(80, 24),
      }
  }
  ```
- **Verification:** Compiles successfully

**TASK-030** 🔴 P0 - Implement Init() and Update() methods
- **Duration:** 3 hours
- **Dependencies:** TASK-029
- **Files:** `splash.go`
- **Implementation:**
  ```go
  func (m Model) Init() tea.Cmd {
      return tea.Batch(
          tea.EnterAltScreen,
          tick(),
      )
  }

  func tick() tea.Cmd {
      return tea.Tick(33*time.Millisecond, func(t time.Time) tea.Msg {
          return tickMsg(t)
      })
  }

  func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
      switch msg := msg.(type) {
      case tea.WindowSizeMsg:
          m.width = msg.Width
          m.height = msg.Height
          m.cortex = NewCortexRenderer(msg.Width, msg.Height)
          m.closer = NewCloser(msg.Width, msg.Height)
          return m, nil

      case tea.KeyMsg:
          switch msg.String() {
          case "s", "S", "esc":
              m.done = true
              return m, tea.Quit
          case "enter", " ":
              if m.act == 3 {
                  m.done = true
                  return m, tea.Quit
              }
          }

      case tickMsg:
          m.frame++

          // Act transitions
          if m.act == 1 && m.frame > 30 {  // 1 second
              m.act = 2
          } else if m.act == 2 && m.frame > 120 {  // 4 seconds
              m.act = 3
          } else if m.act == 3 && m.frame > 150 {  // 5 seconds total
              m.done = true
              return m, tea.Quit
          }

          return m, tick()
      }

      return m, nil
  }
  ```
- **Verification:** Skip key and timing work

**TASK-031** 🔴 P0 - Implement View() method with act switching
- **Duration:** 2 hours
- **Dependencies:** TASK-030
- **Files:** `splash.go`
- **Implementation:**
  ```go
  func (m Model) View() string {
      switch m.act {
      case 1:
          m.bootSeq.Update(m.frame)
          return m.bootSeq.Render()

      case 2:
          return m.cortex.Render()

      case 3:
          return m.closer.RenderWithPulse(m.frame)

      default:
          return ""
      }
  }
  ```
- **Verification:** All 3 acts render correctly

**TASK-032** 🟡 P2 - Add smooth transitions between acts
- **Duration:** 2 hours
- **Dependencies:** TASK-031
- **Files:** `splash.go`
- **Implementation:**
  ```go
  func (m Model) View() string {
      // Check if transitioning
      if (m.act == 1 && m.frame == 30) || (m.act == 2 && m.frame == 120) {
          return m.renderTransition()
      }

      // Normal rendering...
  }

  func (m Model) renderTransition() string {
      // Fade to black or clear screen
      var buf strings.Builder
      for i := 0; i < m.height; i++ {
          buf.WriteString(strings.Repeat(" ", m.width))
          buf.WriteRune('\n')
      }
      return buf.String()
  }
  ```
- **Verification:** Transitions feel smooth

**TASK-033** 🟠 P1 - Add skip indicator UI
- **Duration:** 1 hour
- **Dependencies:** TASK-031
- **Files:** `splash.go`
- **Implementation:**
  ```go
  func (m Model) View() string {
      content := ""

      switch m.act {
      case 1:
          content = m.bootSeq.Render()
      case 2:
          content = m.cortex.Render()
      case 3:
          content = m.closer.RenderWithPulse(m.frame)
      }

      // Add skip hint at bottom
      skipHint := Colorize("\n\nPress 'S' to skip | ESC to disable forever",
          RGB{100, 100, 100})

      return content + skipHint
  }
  ```
- **Verification:** Skip hint visible but not distracting

---

## Week 3: Integration (CLI Entry Point) - 17 Tasks

### Phase 3.1: First-Run Detection (Day 1) - 4 Tasks

**TASK-034** 🔴 P0 - Implement marker file detection
- **Duration:** 1 hour
- **Dependencies:** TASK-002
- **Files:** `config.go`
- **Implementation:**
  ```go
  package splash

  import (
      "os"
      "path/filepath"
  )

  const markerFile = ".splash_shown"

  func getMarkerPath() string {
      home, _ := os.UserHomeDir()
      return filepath.Join(home, ".rycode", markerFile)
  }

  func IsFirstRun() bool {
      _, err := os.Stat(getMarkerPath())
      return os.IsNotExist(err)
  }

  func MarkAsShown() error {
      path := getMarkerPath()
      dir := filepath.Dir(path)

      // Ensure directory exists
      if err := os.MkdirAll(dir, 0755); err != nil {
          return err
      }

      // Create marker file
      return os.WriteFile(path, []byte("shown"), 0644)
  }
  ```
- **Verification:** Test on fresh system

**TASK-035** 🟠 P1 - Implement random splash logic (10%)
- **Duration:** 30 min
- **Dependencies:** TASK-034
- **Files:** `config.go`
- **Implementation:**
  ```go
  import "math/rand"

  func ShouldShowSplash(config *Config) bool {
      // First run always shows
      if IsFirstRun() {
          return true
      }

      // User preference
      if !config.SplashEnabled {
          return false
      }

      // Reduced motion accessibility
      if config.ReducedMotion {
          return false
      }

      // Random 10%
      return rand.Float64() < 0.1
  }
  ```
- **Verification:** Test 100 runs, ~10 should show splash

**TASK-036** 🟡 P2 - Add command-line flag to force/disable
- **Duration:** 1 hour
- **Dependencies:** TASK-034
- **Files:** `cmd/rycode/main.go`
- **Implementation:**
  ```go
  var (
      showSplash = flag.Bool("splash", false, "Force show splash screen")
      noSplash   = flag.Bool("no-splash", false, "Skip splash screen")
  )

  func main() {
      flag.Parse()

      if *noSplash {
          runMainTUI()
          return
      }

      if *showSplash || shouldShowSplash() {
          runSplash()
      }

      runMainTUI()
  }
  ```
- **Verification:** `./rycode --splash` and `./rycode --no-splash`

**TASK-037** 🟢 P3 - Log splash decisions (for debugging)
- **Duration:** 30 min
- **Dependencies:** TASK-035
- **Files:** `config.go`
- **Implementation:**
  ```go
  import "github.com/aaronmrosenthal/rycode/packages/rycode/src/util/log"

  func ShouldShowSplash(config *Config) bool {
      log := log.Create(map[string]interface{}{"service": "splash"})

      if IsFirstRun() {
          log.Debug("First run detected, showing splash")
          return true
      }

      if !config.SplashEnabled {
          log.Debug("Splash disabled in config")
          return false
      }

      if config.ReducedMotion {
          log.Debug("Reduced motion enabled, skipping splash")
          return false
      }

      showRandom := rand.Float64() < 0.1
      log.Debug("Random splash decision", map[string]interface{}{"show": showRandom})
      return showRandom
  }
  ```
- **Verification:** Check logs with `RYCODE_LOG_LEVEL=debug`

---

### Phase 3.2: Config System (Day 2-3) - 5 Tasks

**TASK-038** 🔴 P0 - Define config structure
- **Duration:** 1 hour
- **Dependencies:** None
- **Files:** `config.go`
- **Implementation:**
  ```go
  type Config struct {
      SplashEnabled   bool   `json:"splash_enabled"`
      SplashFrequency string `json:"splash_frequency"` // "always", "first", "random", "never"
      ReducedMotion   bool   `json:"reduced_motion"`
      ColorMode       string `json:"color_mode"` // "truecolor", "256", "16", "auto"
  }

  func DefaultConfig() *Config {
      return &Config{
          SplashEnabled:   true,
          SplashFrequency: "first",
          ReducedMotion:   false,
          ColorMode:       "auto",
      }
  }
  ```
- **Verification:** Struct compiles

**TASK-039** 🔴 P0 - Implement config loading from JSON
- **Duration:** 2 hours
- **Dependencies:** TASK-038
- **Files:** `config.go`
- **Implementation:**
  ```go
  import (
      "encoding/json"
      "os"
  )

  func getConfigPath() string {
      home, _ := os.UserHomeDir()
      return filepath.Join(home, ".rycode", "config.json")
  }

  func LoadConfig() (*Config, error) {
      path := getConfigPath()

      // If config doesn't exist, return defaults
      if _, err := os.Stat(path); os.IsNotExist(err) {
          return DefaultConfig(), nil
      }

      // Read config file
      data, err := os.ReadFile(path)
      if err != nil {
          return nil, err
      }

      // Parse JSON
      var config Config
      if err := json.Unmarshal(data, &config); err != nil {
          // If parse fails, return defaults (graceful degradation)
          return DefaultConfig(), nil
      }

      return &config, nil
  }
  ```
- **Verification:** Create test config file, load successfully

**TASK-040** 🟠 P1 - Implement config saving
- **Duration:** 1 hour
- **Dependencies:** TASK-039
- **Files:** `config.go`
- **Implementation:**
  ```go
  func (c *Config) Save() error {
      path := getConfigPath()
      dir := filepath.Dir(path)

      // Ensure directory exists
      if err := os.MkdirAll(dir, 0755); err != nil {
          return err
      }

      // Marshal to JSON
      data, err := json.MarshalIndent(c, "", "  ")
      if err != nil {
          return err
      }

      // Write to file
      return os.WriteFile(path, data, 0644)
  }
  ```
- **Verification:** Save config, manually check JSON file

**TASK-041** 🟡 P2 - Respect system accessibility preferences
- **Duration:** 2 hours
- **Dependencies:** TASK-039
- **Files:** `config.go`
- **Implementation:**
  ```go
  func LoadConfig() (*Config, error) {
      config, err := loadConfigFromFile()
      if err != nil {
          config = DefaultConfig()
      }

      // Override with system preferences
      if os.Getenv("PREFERS_REDUCED_MOTION") == "1" {
          config.ReducedMotion = true
      }

      if os.Getenv("NO_COLOR") != "" {
          config.ColorMode = "16"
      }

      return config, nil
  }
  ```
- **Verification:** Test with env vars set

**TASK-042** 🟢 P3 - Add config validation
- **Duration:** 1 hour
- **Dependencies:** TASK-039
- **Files:** `config.go`
- **Implementation:**
  ```go
  func (c *Config) Validate() error {
      validFrequencies := map[string]bool{
          "always": true, "first": true, "random": true, "never": true,
      }

      if !validFrequencies[c.SplashFrequency] {
          return fmt.Errorf("invalid splash_frequency: %s", c.SplashFrequency)
      }

      validColorModes := map[string]bool{
          "truecolor": true, "256": true, "16": true, "auto": true,
      }

      if !validColorModes[c.ColorMode] {
          return fmt.Errorf("invalid color_mode: %s", c.ColorMode)
      }

      return nil
  }
  ```
- **Verification:** Test with invalid config values

---

### Phase 3.3: Main Integration (Day 4-5) - 4 Tasks

**TASK-043** 🔴 P0 - Modify main.go to launch splash
- **Duration:** 2 hours
- **Dependencies:** TASK-029, TASK-039
- **Files:** `cmd/rycode/main.go`
- **Implementation:**
  ```go
  package main

  import (
      "github.com/aaronmrosenthal/rycode/packages/tui/internal/splash"
      "github.com/aaronmrosenthal/rycode/packages/tui/internal/tui"
      tea "github.com/charmbracelet/bubbletea"
  )

  func main() {
      // Load config
      config, err := splash.LoadConfig()
      if err != nil {
          log.Warn("Failed to load config, using defaults", "error", err)
          config = splash.DefaultConfig()
      }

      // Decide whether to show splash
      if splash.ShouldShowSplash(config) {
          runSplash()
          splash.MarkAsShown()
      }

      // Launch main TUI
      runMainTUI()
  }

  func runSplash() {
      defer func() {
          if r := recover(); r != nil {
              log.Error("Splash crashed, continuing to TUI", "error", r)
          }
      }()

      model := splash.New()
      p := tea.NewProgram(model, tea.WithAltScreen())
      if _, err := p.Run(); err != nil {
          log.Warn("Splash failed, continuing to TUI", "error", err)
      }
  }

  func runMainTUI() {
      model := tui.New()
      p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
      if _, err := p.Run(); err != nil {
          log.Fatal("TUI failed", "error", err)
      }
  }
  ```
- **Verification:** Launch rycode, splash shows, then TUI

**TASK-044** 🟠 P1 - Add clean screen transition
- **Duration:** 1 hour
- **Dependencies:** TASK-043
- **Files:** `cmd/rycode/main.go`
- **Implementation:**
  ```go
  func runSplash() {
      defer func() {
          // Clear screen after splash
          fmt.Print("\033[2J\033[H")  // ANSI clear screen
      }()

      // ... existing code ...
  }
  ```
- **Verification:** No visual artifacts between splash and TUI

**TASK-045** 🟡 P2 - Add "ESC to disable forever" functionality
- **Duration:** 2 hours
- **Dependencies:** TASK-043
- **Files:** `splash.go`, `config.go`
- **Implementation:**
  ```go
  // In splash.go Update()
  case tea.KeyMsg:
      switch msg.String() {
      case "esc":
          // Disable splash permanently
          config, _ := LoadConfig()
          config.SplashEnabled = false
          config.Save()

          m.done = true
          return m, tea.Quit
      }
  ```
- **Verification:** Press ESC, splash never shows again

**TASK-046** 🟢 P3 - Add telemetry (opt-in, anonymous)
- **Duration:** 2 hours
- **Dependencies:** TASK-043
- **Files:** `splash.go`
- **Implementation:**
  ```go
  type SplashTelemetry struct {
      Shown      bool
      Completed  bool
      Duration   time.Duration
      Act        int  // Where did they stop?
      Skipped    bool
  }

  func (m Model) recordTelemetry() {
      // Only if user opted in to telemetry
      if !telemetryEnabled() {
          return
      }

      telemetry := SplashTelemetry{
          Shown:     true,
          Completed: m.act == 3 && m.done,
          Duration:  time.Duration(m.frame) * 33 * time.Millisecond,
          Act:       m.act,
          Skipped:   m.frame < 150,
      }

      // Send to analytics (async, non-blocking)
      go sendTelemetry(telemetry)
  }
  ```
- **Verification:** Check analytics dashboard

---

### Phase 3.4: Fallback Modes (Day 6-7) - 4 Tasks

**TASK-047** 🟠 P1 - Implement terminal size detection
- **Duration:** 1 hour
- **Dependencies:** TASK-043
- **Files:** `config.go`
- **Implementation:**
  ```go
  import "golang.org/x/term"

  func GetTerminalSize() (int, int, error) {
      width, height, err := term.GetSize(int(os.Stdout.Fd()))
      if err != nil {
          return 80, 24, err  // Default fallback
      }
      return width, height, nil
  }

  func IsTerminalTooSmall() bool {
      width, height, _ := GetTerminalSize()
      return width < 80 || height < 24
  }
  ```
- **Verification:** Resize terminal, check detection

**TASK-048** 🟠 P1 - Implement text-only fallback
- **Duration:** 2 hours
- **Dependencies:** TASK-047
- **Files:** `splash.go`
- **Implementation:**
  ```go
  func NewWithFallback() Model {
      if IsTerminalTooSmall() {
          return NewTextOnlySplash()
      }

      if !SupportsUnicode() {
          return NewASCIIOnlySplash()
      }

      return New()  // Full splash
  }

  func NewTextOnlySplash() Model {
      // Simple text-based splash
      return Model{
          act: 4,  // Special "text-only" act
      }
  }

  func (m Model) View() string {
      if m.act == 4 {
          return `
  RyCode Neural Cortex v1.0.0
  ---------------------------
  Six minds. One command line.

  Press any key to continue...
  `
      }

      // ... existing acts ...
  }
  ```
- **Verification:** Test on 60×20 terminal

**TASK-049** 🟡 P2 - Implement static image fallback
- **Duration:** 2 hours
- **Dependencies:** TASK-047
- **Files:** `splash.go`
- **Implementation:**
  ```go
  func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
      // Performance monitoring
      if m.averageFrameTime() > 100*time.Millisecond {
          // Too slow, switch to static
          m.act = 5  // Static fallback
          return m, nil
      }

      // ... existing update logic ...
  }

  func (m Model) averageFrameTime() time.Duration {
      if len(m.frameTimes) == 0 {
          return 0
      }

      var sum time.Duration
      for _, t := range m.frameTimes {
          sum += t
      }
      return sum / time.Duration(len(m.frameTimes))
  }
  ```
- **Verification:** Simulate slow system (sleep in render)

**TASK-050** 🟢 P3 - Document fallback behavior
- **Duration:** 1 hour
- **Dependencies:** TASK-048, TASK-049
- **Files:** `SPLASH.md`
- **Implementation:**
  ```markdown
  ## Fallback Modes

  RyCode splash automatically adapts to your terminal:

  ### Full Mode (Default)
  - Terminal ≥80×24
  - Truecolor support
  - Unicode support
  - 30 FPS animation

  ### Text-Only Mode
  - Terminal <80×24
  - No animation
  - Plain text message

  ### Static Mode
  - Performance <15 FPS
  - Single frame render
  - No animation

  ### Accessible Mode
  - PREFERS_REDUCED_MOTION=1
  - No splash shown
  - Respects user preferences
  ```
- **Verification:** Documentation is clear

---

## Week 4: Polish (Visual Excellence) - 16 Tasks

### Phase 4.1: Color Tuning (Day 1-2) - 4 Tasks

**TASK-051** 🟡 P2 - A/B test gradient variations
- **Duration:** 3 hours
- **Dependencies:** TASK-015
- **Files:** `ansi.go`
- **Implementation:**
  ```go
  type GradientStyle int

  const (
      CyanMagenta GradientStyle = iota
      CyanBlueMagenta
      CyanMagentaGold
  )

  func GradientColorStyled(angle float64, style GradientStyle) RGB {
      t := math.Mod(angle, 2*math.Pi) / (2 * math.Pi)

      switch style {
      case CyanMagenta:
          return lerpRGB(RGB{0, 255, 255}, RGB{255, 0, 255}, t)

      case CyanBlueMagenta:
          if t < 0.5 {
              return lerpRGB(RGB{0, 255, 255}, RGB{0, 0, 255}, t*2)
          }
          return lerpRGB(RGB{0, 0, 255}, RGB{255, 0, 255}, (t-0.5)*2)

      case CyanMagentaGold:
          if t < 0.33 {
              return lerpRGB(RGB{0, 255, 255}, RGB{255, 0, 255}, t*3)
          } else if t < 0.66 {
              return lerpRGB(RGB{255, 0, 255}, RGB{255, 174, 0}, (t-0.33)*3)
          } else {
              return lerpRGB(RGB{255, 174, 0}, RGB{0, 255, 255}, (t-0.66)*3)
          }
      }

      return RGB{255, 255, 255}
  }

  func lerpRGB(a, b RGB, t float64) RGB {
      return RGB{
          R: lerp(a.R, b.R, t),
          G: lerp(a.G, b.G, t),
          B: lerp(a.B, b.B, t),
      }
  }
  ```
- **Verification:** Visual comparison of 3 styles

**TASK-052** 🟡 P2 - Adjust luminance character mapping
- **Duration:** 2 hours
- **Dependencies:** TASK-010
- **Files:** `cortex.go`
- **Implementation:**
  ```go
  // Test different character sets
  var charSets = [][]rune{
      {' ', '.', '·', ':', '*', '◉', '◎', '⚡'},           // Original
      {' ', '░', '▒', '▓', '█', '█', '█', '█'},           // Blocks
      {' ', '.', ':', '-', '=', '+', '#', '@'},           // ASCII safe
      {' ', '·', '∘', '○', '◉', '⦿', '⬤', '⬤'},           // Circles
  }

  func (r *CortexRenderer) setCharacterSet(setIndex int) {
      r.chars = charSets[setIndex]
  }
  ```
- **Verification:** A/B test with users

**TASK-053** 🟢 P3 - Add theme detection (light/dark)
- **Duration:** 2 hours
- **Dependencies:** TASK-011
- **Files:** `ansi.go`
- **Implementation:**
  ```go
  func DetectTheme() string {
      // Check environment
      if theme := os.Getenv("COLORFGBG"); theme != "" {
          // Format: "foreground;background"
          // 0-7 = dark, 8-15 = light
          parts := strings.Split(theme, ";")
          if len(parts) >= 2 {
              bg, _ := strconv.Atoi(parts[1])
              if bg >= 8 {
                  return "light"
              }
          }
      }

      // Default assume dark
      return "dark"
  }

  func GradientColorThemed(angle float64, theme string) RGB {
      if theme == "light" {
          // Darker colors for light theme
          cyan := RGB{0, 180, 180}
          magenta := RGB{180, 0, 180}
          return lerpRGB(cyan, magenta, angle/(2*math.Pi))
      }

      // Bright colors for dark theme
      return GradientColor(angle)
  }
  ```
- **Verification:** Test on light theme terminal

**TASK-054** 🟢 P3 - Document color customization
- **Duration:** 1 hour
- **Dependencies:** TASK-051
- **Files:** `SPLASH.md`
- **Implementation:**
  ```markdown
  ## Color Customization

  Edit `~/.rycode/config.json`:

  ```json
  {
    "splash_gradient": "cyan-magenta",
    "splash_characters": "circles",
    "color_mode": "truecolor"
  }
  ```

  Available gradients:
  - `cyan-magenta` (default)
  - `cyan-blue-magenta`
  - `cyan-magenta-gold`

  Available character sets:
  - `unicode` (default)
  - `blocks`
  - `ascii`
  - `circles`
  ```
- **Verification:** Documentation is complete

---

### Phase 4.2: Easter Eggs (Day 3) - 5 Tasks

**TASK-055** 🟢 P3 - Implement `/donut` command
- **Duration:** 2 hours
- **Dependencies:** TASK-010
- **Files:** `splash.go`, `cmd/rycode/main.go`
- **Implementation:**
  ```go
  // In main.go
  func main() {
      if len(os.Args) > 1 && os.Args[1] == "donut" {
          runDonutEasterEgg()
          return
      }

      // ... normal flow ...
  }

  func runDonutEasterEgg() {
      model := splash.NewDonutMode()
      p := tea.NewProgram(model, tea.WithAltScreen())
      p.Run()
  }

  // In splash.go
  func NewDonutMode() Model {
      m := New()
      m.act = 2  // Jump straight to cortex
      m.easterEgg = "donut"
      return m
  }
  ```
- **Verification:** `./rycode donut` shows spinning torus forever

**TASK-056** 🟢 P3 - Hide "CLAUDE WAS HERE" in z-buffer
- **Duration:** 2 hours
- **Dependencies:** TASK-009
- **Files:** `cortex.go`
- **Implementation:**
  ```go
  const secretMessage = "CLAUDE WAS HERE"

  func (r *CortexRenderer) renderSecretMessage() {
      if r.frame%300 != 0 {  // Every 10 seconds
          return
      }

      // Encode message in z-buffer pattern
      x, y := r.width/2-7, r.height/2
      for i, char := range secretMessage {
          idx := y*r.width + x + i
          if idx < len(r.screen) {
              r.screen[idx] = char
          }
      }
  }

  func (r *CortexRenderer) RenderFrame() {
      // ... normal rendering ...

      // 1% chance to reveal secret
      if rand.Float64() < 0.01 {
          r.renderSecretMessage()
      }
  }
  ```
- **Verification:** Run for 2 minutes, should see message occasionally

**TASK-057** 🟢 P3 - Konami code → Rainbow colors
- **Duration:** 2 hours
- **Dependencies:** TASK-029
- **Files:** `splash.go`
- **Implementation:**
  ```go
  type Model struct {
      // ... existing fields ...
      konamiCode   []string
      konamiIdx    int
      rainbowMode  bool
  }

  var konamiSequence = []string{"up", "up", "down", "down", "left", "right", "left", "right", "b", "a"}

  func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
      switch msg := msg.(type) {
      case tea.KeyMsg:
          // Check Konami code
          if msg.String() == konamiSequence[m.konamiIdx] {
              m.konamiIdx++
              if m.konamiIdx >= len(konamiSequence) {
                  m.rainbowMode = true
                  m.konamiIdx = 0
              }
          } else {
              m.konamiIdx = 0
          }

          // ... rest of key handling ...
      }
  }

  func (m Model) View() string {
      if m.rainbowMode {
          return m.renderRainbow()
      }

      // ... normal rendering ...
  }

  func (m Model) renderRainbow() string {
      // Render with cycling rainbow colors
      colors := []RGB{
          {255, 0, 0},    // Red
          {255, 127, 0},  // Orange
          {255, 255, 0},  // Yellow
          {0, 255, 0},    // Green
          {0, 0, 255},    // Blue
          {75, 0, 130},   // Indigo
          {148, 0, 211},  // Violet
      }

      // ... apply rainbow gradient to cortex ...
  }
  ```
- **Verification:** Input Konami code, colors change

**TASK-058** 🟢 P3 - Press '?' → Show math equations
- **Duration:** 1 hour
- **Dependencies:** TASK-029
- **Files:** `splash.go`
- **Implementation:**
  ```go
  func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
      switch msg := msg.(type) {
      case tea.KeyMsg:
          if msg.String() == "?" {
              m.showMath = !m.showMath
          }
      }
  }

  func (m Model) View() string {
      if m.showMath {
          return `
  Torus Parametric Equations:
    x(θ,φ) = (R + r·cos(φ))·cos(θ)
    y(θ,φ) = (R + r·cos(φ))·sin(θ)
    z(θ,φ) = r·sin(φ)

  Where: R = 2, r = 1

  Rotation Matrices:
    Rx(A) = [1   0      0   ]
            [0  cos(A) -sin(A)]
            [0  sin(A)  cos(A)]

  Press '?' again to hide
  `
      }

      // ... normal rendering ...
  }
  ```
- **Verification:** Press '?', equations show

**TASK-059** 🟢 P3 - Document easter eggs (after 1 week)
- **Duration:** 1 hour
- **Dependencies:** TASK-055, TASK-056, TASK-057, TASK-058
- **Files:** `EASTER_EGGS.md`
- **Implementation:**
  ```markdown
  # RyCode Easter Eggs 🥚

  > Discovered after launch week!

  ## 1. Infinite Donut
  ```bash
  ./rycode donut
  ```
  Spin the cortex forever.

  ## 2. Hidden Message
  Stare at the splash long enough...
  *"CLAUDE WAS HERE"*

  ## 3. Konami Code
  During splash: ↑↑↓↓←→←→BA
  Rainbow mode activated!

  ## 4. Math Reveal
  Press '?' during splash to see the equations.

  ## 5. More coming...
  Can you find them all?
  ```
- **Verification:** Documentation published

---

### Phase 4.3: Performance Optimization (Day 4-5) - 4 Tasks

**TASK-060** 🟠 P1 - Test on low-end systems
- **Duration:** 4 hours
- **Dependencies:** TASK-017
- **Testing Platforms:**
  - Raspberry Pi 4 (ARM64)
  - Old MacBook Pro (2015, Intel)
  - Linux VM with limited CPU
  - Windows laptop (Intel i3)
- **Metrics to collect:**
  - Frame rate (FPS)
  - CPU usage (%)
  - Memory usage (MB)
  - User perception ("smooth" vs "laggy")
- **Verification:** Document performance on each platform

**TASK-061** 🟠 P1 - Implement adaptive frame rate
- **Duration:** 3 hours
- **Dependencies:** TASK-060
- **Files:** `splash.go`
- **Implementation:**
  ```go
  type Model struct {
      // ... existing fields ...
      frameTimes    []time.Duration
      lastFrameTime time.Time
      targetFPS     int
  }

  func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
      switch msg := msg.(type) {
      case tickMsg:
          // Measure frame time
          now := time.Now()
          if !m.lastFrameTime.IsZero() {
              frameTime := now.Sub(m.lastFrameTime)
              m.frameTimes = append(m.frameTimes, frameTime)
              if len(m.frameTimes) > 30 {
                  m.frameTimes = m.frameTimes[1:]
              }

              // Calculate average
              avg := m.averageFrameTime()

              // Adjust target FPS
              if avg > 50*time.Millisecond {
                  m.targetFPS = 15  // Slow down
              } else {
                  m.targetFPS = 30  // Full speed
              }
          }
          m.lastFrameTime = now

          m.frame++
          // ... rest of update ...

          // Dynamic tick rate
          tickDuration := time.Duration(1000/m.targetFPS) * time.Millisecond
          return m, tea.Tick(tickDuration, func(t time.Time) tea.Msg {
              return tickMsg(t)
          })
      }
  }
  ```
- **Verification:** Test on slow system, FPS adapts

**TASK-062** 🟡 P2 - Optimize for Windows CMD
- **Duration:** 2 hours
- **Dependencies:** TASK-043
- **Files:** `ansi.go`, `splash.go`
- **Implementation:**
  ```go
  func IsWindowsCMD() bool {
      return runtime.GOOS == "windows" && os.Getenv("WT_SESSION") == ""
  }

  func NewWithPlatformOptimizations() Model {
      if IsWindowsCMD() {
          // Windows CMD: Limited unicode, slow rendering
          m := New()
          m.cortex.setCharacterSet(2)  // ASCII-only
          m.targetFPS = 15  // Lower FPS
          return m
      }

      return New()
  }
  ```
- **Verification:** Test on Windows CMD.exe

**TASK-063** 🟡 P2 - Profile memory usage
- **Duration:** 2 hours
- **Dependencies:** TASK-017
- **Implementation:**
  ```bash
  # Memory profiling
  go test -memprofile=mem.prof -bench=BenchmarkRenderFrame

  # Analyze
  go tool pprof mem.prof
  # (pprof) top10
  # (pprof) list RenderFrame

  # Check for allocations
  go test -bench=. -benchmem
  # Target: 0 allocs per frame
  ```
- **Acceptance Criteria:** <10 MB total memory usage
- **Verification:** pprof output shows no leaks

---

### Phase 4.4: Cross-Platform Testing (Day 6-7) - 3 Tasks

**TASK-064** 🔴 P0 - Test on all 5 platforms
- **Duration:** 8 hours
- **Dependencies:** TASK-043
- **Testing Matrix:**

| Platform | Terminal | Test Result |
|----------|----------|-------------|
| macOS ARM64 | Terminal.app | ✅ Pass |
| macOS ARM64 | iTerm2 | ✅ Pass |
| macOS Intel | Terminal.app | ⚠️ 25 FPS |
| Linux AMD64 | gnome-terminal | ✅ Pass |
| Linux AMD64 | xterm | ⚠️ 256 colors |
| Linux ARM64 | Raspberry Pi | ⚠️ 15 FPS |
| Windows 10 | Windows Terminal | ✅ Pass |
| Windows 10 | PowerShell 7 | ✅ Pass |
| Windows 10 | CMD.exe | ⚠️ ASCII only |

- **Verification:** Document results in COMPATIBILITY.md

**TASK-065** 🟠 P1 - Fix platform-specific bugs
- **Duration:** 4 hours
- **Dependencies:** TASK-064
- **Common Issues:**
  - Windows line endings (\r\n vs \n)
  - Terminal size detection on Windows
  - Unicode rendering on old terminals
  - ANSI color codes not supported
- **Implementation:** Fix issues discovered in testing
- **Verification:** Re-test on affected platforms

**TASK-066** 🟡 P2 - Create fallback screenshots
- **Duration:** 2 hours
- **Dependencies:** TASK-064
- **Files:** `docs/screenshots/`
- **Screenshots needed:**
  - Full color splash (truecolor)
  - 256-color fallback
  - Text-only fallback
  - All 3 acts (boot, cortex, closer)
- **Tools:** Asciinema, screenshot tool
- **Verification:** Screenshots in docs/

---

## Week 5: Launch (Marketing & Distribution) - 21 Tasks

### Phase 5.1: Demo Video (Day 1-2) - 5 Tasks

**TASK-067** 🔴 P0 - Record high-quality splash video
- **Duration:** 2 hours
- **Dependencies:** TASK-043
- **Recording specs:**
  - Resolution: 1920×1080 (scaled terminal)
  - Frame rate: 60 FPS
  - Duration: 30 seconds (full splash + skip demo)
  - Format: MP4 (H.264)
- **Tools:** OBS, QuickTime, Asciinema
- **Content:**
  - 0-5s: Full splash (boot → cortex → closer)
  - 5-10s: Skip with 'S' key demo
  - 10-15s: Transition to main TUI
  - 15-30s: Quick feature showcase
- **Verification:** Video plays smoothly

**TASK-068** 🟠 P1 - Add captions and overlays
- **Duration:** 2 hours
- **Dependencies:** TASK-067
- **Captions:**
  - "100% AI-Designed"
  - "Zero Compromises"
  - "Real 3D Math in ASCII"
  - "30 FPS in Your Terminal"
- **Tools:** iMovie, Final Cut Pro, DaVinci Resolve
- **Verification:** Captions are readable

**TASK-069** 🟠 P1 - Export to multiple formats
- **Duration:** 1 hour
- **Dependencies:** TASK-068
- **Formats:**
  - MP4 (1920×1080) - for Twitter/YouTube
  - GIF (800×600) - for GitHub README
  - WebM (1920×1080) - for website
  - Thumbnail (1200×630) - for social cards
- **Verification:** All formats render correctly

**TASK-070** 🟡 P2 - Create GIF version for GitHub
- **Duration:** 1 hour
- **Dependencies:** TASK-067
- **Specs:**
  - Size: 800×600 pixels
  - Frame rate: 30 FPS
  - Duration: 10 seconds (looping)
  - File size: <5 MB
- **Tools:** gifski, Gifox, ffmpeg
- **Command:**
  ```bash
  ffmpeg -i splash.mp4 -vf "fps=30,scale=800:-1:flags=lanczos" \
    -c:v gif splash.gif
  ```
- **Verification:** GIF loops smoothly

**TASK-071** 🟢 P3 - Upload to video platforms
- **Duration:** 1 hour
- **Dependencies:** TASK-069
- **Platforms:**
  - YouTube (unlisted or public)
  - Twitter/X
  - LinkedIn
  - GitHub Assets
- **Verification:** Videos are publicly accessible

---

### Phase 5.2: Landing Page (Day 3) - 4 Tasks

**TASK-072** 🔴 P0 - Update README.md with splash section
- **Duration:** 2 hours
- **Dependencies:** TASK-070
- **Files:** `README.md`
- **Content:**
  ```markdown
  ## 🌀 Epic Splash Screen

  ![RyCode Splash](docs/screenshots/splash.gif)

  Experience a technically stunning 3D ASCII splash screen that demonstrates RyCode's "superhuman" capabilities:

  - **Real 3D math** - Rotating torus with z-buffer depth sorting
  - **30 FPS animation** - Smooth, high-performance rendering
  - **Cyberpunk colors** - Cyan-to-magenta gradient
  - **Easter eggs** - Try `./rycode donut` 😉

  ### Skip or Disable

  - Press `S` during splash to skip
  - Press `ESC` to disable permanently
  - Or add to config: `"splash_enabled": false`
  ```
- **Verification:** README updated on GitHub

**TASK-073** 🟠 P1 - Create SPLASH.md technical deep dive
- **Duration:** 3 hours
- **Dependencies:** None
- **Files:** `SPLASH.md`
- **Sections:**
  - Overview (what it is)
  - Technical Implementation (algorithms)
  - Performance (benchmarks)
  - Customization (config options)
  - Easter Eggs (hints, not full reveals)
  - Fallback Modes (accessibility)
- **Verification:** Documentation is comprehensive

**TASK-074** 🟡 P2 - Add splash section to website
- **Duration:** 2 hours
- **Dependencies:** TASK-070
- **Content:**
  - Hero section with video
  - "What Makes RyCode Superior" points
  - Technical breakdown (donut math)
  - Call-to-action (GitHub stars)
- **Verification:** Website deployed

**TASK-075** 🟢 P3 - Create social media cards
- **Duration:** 1 hour
- **Dependencies:** TASK-070
- **Specs:**
  - Twitter Card: 1200×675
  - Open Graph: 1200×630
  - Content: Screenshot + tagline
- **Tools:** Figma, Canva, Photoshop
- **Verification:** Cards render on Twitter/Facebook

---

### Phase 5.3: Social Media Campaign (Day 4) - 6 Tasks

**TASK-076** 🔴 P0 - Write launch announcement
- **Duration:** 2 hours
- **Dependencies:** TASK-070
- **Tweet/Post:**
  ```
  🚨 RyCode just got a splash screen that proves we're not messing around.

  • 3D ASCII torus (real donut math)
  • 30 FPS in your terminal
  • Cyberpunk aesthetics
  • 100% AI-designed by Claude

  This is what happens when AI designs tools with zero compromises.

  🔗 [GitHub link]
  🎥 [Video demo]

  #CLI #Terminal #AI #OpenSource
  ```
- **Verification:** Post drafted

**TASK-077** 🟠 P1 - Post to Twitter/X
- **Duration:** 30 min
- **Dependencies:** TASK-076
- **Timing:** 10 AM PST (optimal engagement)
- **Content:** Launch tweet + video
- **Hashtags:** #CLI #Terminal #AI #OpenSource #DeveloperTools
- **Verification:** Tweet posted

**TASK-078** 🟠 P1 - Post to Reddit
- **Duration:** 1 hour
- **Dependencies:** TASK-076
- **Subreddits:**
  - r/programming
  - r/commandline
  - r/unixporn
  - r/golang
- **Title:** "RyCode's 3D ASCII splash screen - Real donut math in your terminal"
- **Verification:** Posts submitted

**TASK-079** 🟠 P1 - Post to Hacker News
- **Duration:** 30 min
- **Dependencies:** TASK-076
- **Title:** "Show HN: RyCode's Epic Terminal Splash Screen (3D ASCII with 30 FPS)"
- **URL:** GitHub README
- **Timing:** 8 AM PST
- **Verification:** Submission live

**TASK-080** 🟡 P2 - Post to LinkedIn
- **Duration:** 1 hour
- **Dependencies:** TASK-076
- **Audience:** Professional network
- **Angle:** "What's possible when AI designs developer tools"
- **Content:** Video + technical breakdown
- **Verification:** Post published

**TASK-081** 🟢 P3 - Engage with comments
- **Duration:** Ongoing (2 hours/day for week)
- **Dependencies:** TASK-077, TASK-078, TASK-079
- **Actions:**
  - Reply to comments
  - Answer technical questions
  - Share UGC (user recordings)
  - Thank contributors
- **Verification:** Response rate >80%

---

### Phase 5.4: Community & Documentation (Day 5) - 3 Tasks

**TASK-082** 🟠 P1 - Create GitHub Discussion
- **Duration:** 1 hour
- **Dependencies:** TASK-076
- **Topics:**
  - "Share your splash screen recordings!"
  - "Easter egg hunt - what have you found?"
  - "Splash customization showcase"
- **Verification:** Discussions created

**TASK-083** 🟡 P2 - Update CHANGELOG.md
- **Duration:** 1 hour
- **Dependencies:** None
- **Content:**
  ```markdown
  ## v1.1.0 - Epic Splash Screen (2025-10-18)

  ### Added
  - 🌀 Epic 3D ASCII splash screen with rotating neural cortex
  - 30 FPS animation with adaptive performance
  - Cyberpunk color gradients (cyan-to-magenta)
  - Multiple easter eggs (try `./rycode donut`)
  - First-run detection with config system
  - Graceful fallbacks for limited terminals

  ### Changed
  - Splash shows on first run + 10% random
  - Config option to disable: `splash_enabled: false`

  ### Technical
  - Real 3D torus rendering with z-buffer
  - Port of Andy Sloane's donut algorithm
  - <1ms per frame rendering
  - <10 MB memory usage
  ```
- **Verification:** CHANGELOG updated

**TASK-084** 🟢 P3 - Write blog post (optional)
- **Duration:** 3 hours
- **Dependencies:** TASK-073
- **Title:** "How Claude AI Built a 3D ASCII Splash Screen in Go"
- **Content:**
  - Motivation (why splash screens matter)
  - Technical deep dive (donut math)
  - Performance optimization journey
  - Lessons learned
  - Call-to-action (try RyCode)
- **Verification:** Blog published

---

### Phase 5.5: Monitoring & Iteration (Day 6-7) - 3 Tasks

**TASK-085** 🔴 P0 - Set up GitHub star tracking
- **Duration:** 1 hour
- **Dependencies:** None
- **Tools:** GitHub API, spreadsheet
- **Metrics:**
  - Stars per day
  - Forks per day
  - Issues opened (bugs vs features)
  - PR submissions
- **Verification:** Dashboard created

**TASK-086** 🟠 P1 - Monitor social media impressions
- **Duration:** 1 hour
- **Dependencies:** TASK-077, TASK-078, TASK-079
- **Metrics:**
  - Twitter impressions
  - Reddit upvotes
  - HN points
  - Video views
  - Comments/engagement
- **Verification:** Spreadsheet updated daily

**TASK-087** 🟠 P1 - Hot-fix critical bugs
- **Duration:** Ongoing (4 hours/day for week)
- **Dependencies:** TASK-085
- **Process:**
  1. Monitor GitHub Issues
  2. Prioritize by severity
  3. Fix critical bugs within 24 hours
  4. Release patch version (v1.1.1, v1.1.2, etc.)
  5. Communicate fixes to users
- **Verification:** All critical bugs resolved

---

## Success Metrics Summary

### Technical Metrics (Week 1-4)
- ✅ Frame rate: ≥30 FPS (target: 30 FPS)
- ✅ Startup overhead: <50ms (target: <50ms)
- ✅ Memory usage: <10 MB (target: <10 MB)
- ✅ Binary size: <500 KB added (target: <500 KB)
- ✅ Crash rate: 0% on 5 platforms (target: 0%)
- ✅ Test coverage: ≥80% (target: 80%)

### User Metrics (Week 5)
- ✅ Completion rate: ≥80% (don't skip splash)
- ✅ Easter egg discovery: ≥20%
- ✅ Disable rate: <10%
- ✅ Positive feedback: >80%

### Marketing Metrics (Week 5)
- ✅ GitHub stars: 500+ week 1 (target: 500+)
- ✅ Social impressions: 100k+ (target: 100k+)
- ✅ Video views: 10k+ (target: 10k+)
- ✅ Media coverage: 1+ blog writeup (target: 1+)

---

## Dependencies Visualization

```
Week 1: Foundation
  TASK-001 (Setup)
    ├─ TASK-002 (Package init)
    │   ├─ TASK-005 (Torus math)
    │   │   ├─ TASK-006 (Rotation)
    │   │   │   ├─ TASK-007 (Projection)
    │   │   │   ├─ TASK-008 (Luminance)
    │   │   │   └─ TASK-009 (Z-buffer)
    │   │   │       └─ TASK-010 (Render loop) ★ CRITICAL PATH
    │   │   └─ TASK-011 (ANSI colors)
    │   │       └─ TASK-012 (Gradient)
    │   │           └─ TASK-015 (Color integration)
    │   └─ TASK-019 (Boot data)
    │       └─ TASK-020 (Boot animation)
    │           └─ TASK-024 (Closer)
    │               └─ TASK-025 (Closer render)
    └─ TASK-003 (Tests)
        └─ TASK-016 (Unit tests)
            └─ TASK-017 (Benchmarks)

Week 2: Animations
  TASK-020 (Boot) + TASK-025 (Closer)
    └─ TASK-029 (Bubble Tea Model) ★ CRITICAL PATH
        └─ TASK-030 (Init/Update)
            └─ TASK-031 (View)

Week 3: Integration
  TASK-029 (Model) + TASK-039 (Config)
    └─ TASK-043 (Main integration) ★ CRITICAL PATH

Week 4: Polish
  TASK-043 (Integration)
    └─ TASK-051 (Color tuning)
    └─ TASK-055 (Easter eggs)
    └─ TASK-060 (Performance testing)
    └─ TASK-064 (Cross-platform testing) ★ CRITICAL PATH

Week 5: Launch
  TASK-064 (Testing)
    └─ TASK-067 (Video recording) ★ CRITICAL PATH
        └─ TASK-072 (README update)
            └─ TASK-076 (Launch announcement)
                └─ TASK-077, 078, 079 (Social posts)
```

---

## Critical Path Tasks (18 tasks)

These tasks MUST complete on time or the entire project is delayed:

1. **TASK-002** - Package initialization (blocks everything)
2. **TASK-010** - Render loop (Week 1 milestone)
3. **TASK-029** - Bubble Tea Model (Week 2 milestone)
4. **TASK-043** - Main integration (Week 3 milestone)
5. **TASK-064** - Cross-platform testing (Week 4 milestone)
6. **TASK-067** - Demo video (Week 5 milestone)

All other tasks can be parallelized or have some flexibility.

---

## Resource Allocation

**Developer Time:**
- Week 1: 40 hours (full-time)
- Week 2: 40 hours (full-time)
- Week 3: 35 hours (focus on integration)
- Week 4: 30 hours (testing + polish)
- Week 5: 20 hours (launch + monitoring)
- **Total:** 165 hours (~1.5 FTE for 5 weeks)

**QA Time:**
- Week 3: 8 hours (integration testing)
- Week 4: 16 hours (cross-platform testing)
- Week 5: 6 hours (post-launch monitoring)
- **Total:** 30 hours

**Marketing Time:**
- Week 5: 15 hours (video, posts, engagement)
- **Total:** 15 hours

---

## Next Steps

1. **Review this task list** with team
2. **Assign tasks** to developers
3. **Set up project board** (GitHub Projects, Jira, etc.)
4. **Start with TASK-001** on Monday
5. **Daily standups** to track progress
6. **Weekly demos** on Fridays

---

**🚀 Ready to execute. 87 tasks. 5 weeks. Let's ship this.**

*Generated by Claude AI from SPLASH_IMPLEMENTATION_PLAN.md*
