# RyCode Matrix TUI: Actionable Implementation Tasks

## Overview

This document breaks down the Matrix TUI implementation into specific, executable tasks with priorities, dependencies, time estimates, and acceptance criteria.

**Total Estimated Time:** 12 weeks (480 hours)
**Team Size:** 4 developers
**Tasks:** 180+ discrete tasks

---

## Task Priority System

- **P0 (Critical):** Blocking, must complete before proceeding
- **P1 (High):** Core features, important for MVP
- **P2 (Medium):** Nice-to-have features
- **P3 (Low):** Future enhancements

---

## Phase 1: Foundation & Matrix Theme (Weeks 1-2)

### Sprint 1.1: Project Setup (Days 1-2)

#### TASK-001: Initialize Project Structure
- **Priority:** P0
- **Assignee:** Tech Lead
- **Estimate:** 2 hours
- **Dependencies:** None
- **Description:** Create directory structure for new TUI package

**Acceptance Criteria:**
```bash
✓ Directory structure created:
  packages/tui-v2/
  ├── cmd/rycode/
  ├── internal/{ui,input,ai,theme,layout,animation}/
  ├── pkg/api/
  └── test/{unit,integration,e2e}/
✓ README.md with setup instructions
✓ .gitignore configured
```

**Commands:**
```bash
mkdir -p packages/tui-v2/{cmd/rycode,internal/{ui/{models,components,views},input,ai/{providers},theme,layout,animation,util},pkg/api,test/{unit,integration,e2e}}
cd packages/tui-v2
touch README.md .gitignore
```

---

#### TASK-002: Initialize Go Module
- **Priority:** P0
- **Assignee:** Tech Lead
- **Estimate:** 1 hour
- **Dependencies:** TASK-001
- **Description:** Setup Go module and install core dependencies

**Acceptance Criteria:**
```bash
✓ go.mod created with correct module path
✓ Core dependencies installed:
  - github.com/charmbracelet/bubbletea
  - github.com/charmbracelet/lipgloss
  - github.com/charmbracelet/glamour
  - github.com/alecthomas/chroma/v2
✓ Test dependencies installed:
  - github.com/stretchr/testify
✓ go.sum generated
```

**Commands:**
```bash
go mod init github.com/aaronmrosenthal/rycode/packages/tui-v2
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/charmbracelet/glamour@latest
go get github.com/alecthomas/chroma/v2@latest
go get github.com/stretchr/testify@latest
```

---

#### TASK-003: Create Makefile
- **Priority:** P0
- **Assignee:** Tech Lead
- **Estimate:** 1 hour
- **Dependencies:** TASK-002
- **Description:** Build automation for common tasks

**Acceptance Criteria:**
```makefile
✓ Makefile includes:
  - build: Build binary
  - test: Run all tests
  - test-unit: Run unit tests only
  - test-integration: Run integration tests
  - coverage: Generate coverage report
  - lint: Run linter
  - clean: Clean build artifacts
  - install: Install binary
```

**File:** `packages/tui-v2/Makefile`
```makefile
.PHONY: build test test-unit test-integration coverage lint clean install

BINARY_NAME=rycode
BUILD_DIR=../../packages/rycode/dist

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/rycode

test:
	go test -v ./...

test-unit:
	go test -v ./internal/... ./pkg/...

test-integration:
	go test -v -tags=integration ./test/integration/...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint:
	golangci-lint run

clean:
	rm -rf $(BUILD_DIR)/$(BINARY_NAME)
	rm -f coverage.out

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(HOME)/bin/
```

---

#### TASK-004: Setup CI/CD Pipeline
- **Priority:** P0
- **Assignee:** QA/DevOps
- **Estimate:** 3 hours
- **Dependencies:** TASK-003
- **Description:** GitHub Actions workflow for automated testing

**Acceptance Criteria:**
```yaml
✓ .github/workflows/tui-test.yml created
✓ Runs on push and PR to dev/main
✓ Tests all Go versions (1.21, 1.22)
✓ Runs linter
✓ Generates coverage report
✓ Uploads coverage to Codecov
```

**File:** `.github/workflows/tui-test.yml`
```yaml
name: TUI Tests

on:
  push:
    branches: [main, dev]
    paths:
      - 'packages/tui-v2/**'
  pull_request:
    branches: [main, dev]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.21, 1.22]

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}

      - name: Install dependencies
        working-directory: packages/tui-v2
        run: go mod download

      - name: Run tests
        working-directory: packages/tui-v2
        run: make test

      - name: Run linter
        working-directory: packages/tui-v2
        run: make lint

      - name: Generate coverage
        working-directory: packages/tui-v2
        run: make coverage

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./packages/tui-v2/coverage.out
```

---

#### TASK-005: Create Hello World TUI
- **Priority:** P0
- **Assignee:** Frontend Developer
- **Estimate:** 2 hours
- **Dependencies:** TASK-002
- **Description:** Minimal Bubble Tea app to verify setup

**Acceptance Criteria:**
```
✓ cmd/rycode/main.go created
✓ Displays "Hello, RyCode!" in Matrix green
✓ Responds to 'q' to quit
✓ Runs without errors
✓ Binary builds successfully
```

**File:** `packages/tui-v2/cmd/rycode/main.go`
```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var matrixGreen = lipgloss.Color("#00ff00")

type model struct {
	message string
}

func initialModel() model {
	return model{message: "Hello, RyCode!"}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	style := lipgloss.NewStyle().
		Foreground(matrixGreen).
		Bold(true)

	return style.Render(m.message) + "\n\nPress 'q' to quit"
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
```

**Test:**
```bash
cd packages/tui-v2
make build
../../packages/rycode/dist/rycode
# Should display green text, quit with 'q'
```

---

### Sprint 1.2: Responsive Framework (Days 3-4)

#### TASK-006: Implement DeviceClass Enum
- **Priority:** P0
- **Assignee:** Frontend Developer
- **Estimate:** 1 hour
- **Dependencies:** TASK-005
- **Description:** Define device classes based on terminal dimensions

**Acceptance Criteria:**
```
✓ internal/layout/types.go created
✓ DeviceClass enum defined (6 variants)
✓ String() method for debugging
✓ Unit tests for all variants
```

**File:** `packages/tui-v2/internal/layout/types.go`
```go
package layout

type DeviceClass int

const (
	PhonePortrait DeviceClass = iota
	PhoneLandscape
	TabletPortrait
	TabletLandscape
	DesktopSmall
	DesktopLarge
)

func (dc DeviceClass) String() string {
	return [...]string{
		"PhonePortrait",
		"PhoneLandscape",
		"TabletPortrait",
		"TabletLandscape",
		"DesktopSmall",
		"DesktopLarge",
	}[dc]
}

func (dc DeviceClass) MinWidth() int {
	return [...]int{0, 60, 80, 100, 120, 160}[dc]
}

func (dc DeviceClass) IsMobile() bool {
	return dc == PhonePortrait || dc == PhoneLandscape
}

func (dc DeviceClass) IsTablet() bool {
	return dc == TabletPortrait || dc == TabletLandscape
}

func (dc DeviceClass) IsDesktop() bool {
	return dc == DesktopSmall || dc == DesktopLarge
}
```

---

#### TASK-007: Implement LayoutManager
- **Priority:** P0
- **Assignee:** Frontend Developer
- **Estimate:** 3 hours
- **Dependencies:** TASK-006
- **Description:** Detect device class and manage layout transitions

**Acceptance Criteria:**
```
✓ internal/layout/manager.go created
✓ Auto-detects device class from terminal size
✓ Updates on window resize
✓ Provides appropriate layout for each device class
✓ Unit tests for all breakpoints
```

**File:** `packages/tui-v2/internal/layout/manager.go`
```go
package layout

import (
	"time"
)

type LayoutManager struct {
	width      int
	height     int
	class      DeviceClass
	lastUpdate time.Time
}

func NewLayoutManager(width, height int) *LayoutManager {
	lm := &LayoutManager{
		width:      width,
		height:     height,
		lastUpdate: time.Now(),
	}
	lm.detectDevice()
	return lm
}

func (lm *LayoutManager) detectDevice() {
	switch {
	case lm.width >= 160:
		lm.class = DesktopLarge
	case lm.width >= 120:
		lm.class = DesktopSmall
	case lm.width >= 100:
		lm.class = TabletLandscape
	case lm.width >= 80:
		lm.class = TabletPortrait
	case lm.width >= 60:
		lm.class = PhoneLandscape
	default:
		lm.class = PhonePortrait
	}
}

func (lm *LayoutManager) Update(width, height int) {
	oldClass := lm.class
	lm.width = width
	lm.height = height
	lm.detectDevice()
	lm.lastUpdate = time.Now()

	if oldClass != lm.class {
		// Device class changed, trigger layout transition
	}
}

func (lm *LayoutManager) GetDeviceClass() DeviceClass {
	return lm.class
}

func (lm *LayoutManager) GetDimensions() (int, int) {
	return lm.width, lm.height
}
```

**Test File:** `packages/tui-v2/internal/layout/manager_test.go`
```go
package layout

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLayoutManager_DetectDevice(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		want   DeviceClass
	}{
		{"Phone portrait", 40, 80, PhonePortrait},
		{"Phone landscape", 80, 40, PhoneLandscape},
		{"Tablet portrait", 90, 100, TabletPortrait},
		{"Tablet landscape", 120, 60, TabletLandscape},
		{"Desktop small", 140, 60, DesktopSmall},
		{"Desktop large", 180, 60, DesktopLarge},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm := NewLayoutManager(tt.width, tt.height)
			assert.Equal(t, tt.want, lm.GetDeviceClass())
		})
	}
}

func TestLayoutManager_Update(t *testing.T) {
	lm := NewLayoutManager(40, 80)
	assert.Equal(t, PhonePortrait, lm.GetDeviceClass())

	lm.Update(180, 60)
	assert.Equal(t, DesktopLarge, lm.GetDeviceClass())
}
```

---

#### TASK-008: Implement Layout Interfaces
- **Priority:** P0
- **Assignee:** Frontend Developer
- **Estimate:** 2 hours
- **Dependencies:** TASK-007
- **Description:** Define layout interface and concrete implementations

**Acceptance Criteria:**
```
✓ internal/layout/layout.go with Layout interface
✓ StackLayout (phone)
✓ SplitLayout (tablet)
✓ MultiPaneLayout (desktop)
✓ Unit tests for each layout type
```

**File:** `packages/tui-v2/internal/layout/layout.go`
```go
package layout

import tea "github.com/charmbracelet/bubbletea"

type Layout interface {
	AddView(name string, view tea.Model)
	RemoveView(name string)
	SetActive(name string)
	Render(width, height int) string
	Update(msg tea.Msg) tea.Cmd
}

// StackLayout - Single pane, stack navigation (phone)
type StackLayout struct {
	views   map[string]tea.Model
	stack   []string
	current string
}

func NewStackLayout() *StackLayout {
	return &StackLayout{
		views: make(map[string]tea.Model),
		stack: []string{},
	}
}

func (sl *StackLayout) AddView(name string, view tea.Model) {
	sl.views[name] = view
}

func (sl *StackLayout) SetActive(name string) {
	if _, exists := sl.views[name]; exists {
		sl.stack = append(sl.stack, sl.current)
		sl.current = name
	}
}

func (sl *StackLayout) Render(width, height int) string {
	if view, exists := sl.views[sl.current]; exists {
		return view.View()
	}
	return "No active view"
}

// SplitLayout - Two panes side-by-side (tablet)
type SplitLayout struct {
	left     tea.Model
	right    tea.Model
	ratio    float64 // 0.0 - 1.0
}

func NewSplitLayout(ratio float64) *SplitLayout {
	return &SplitLayout{
		ratio: ratio,
	}
}

// MultiPaneLayout - 4 panes (desktop)
type MultiPaneLayout struct {
	tree    tea.Model
	editor  tea.Model
	chat    tea.Model
	metrics tea.Model
}

func NewMultiPaneLayout() *MultiPaneLayout {
	return &MultiPaneLayout{}
}
```

---

### Sprint 1.3: Matrix Theme (Day 5)

#### TASK-009: Define Matrix Color Palette
- **Priority:** P0
- **Assignee:** Frontend Developer
- **Estimate:** 1 hour
- **Dependencies:** TASK-005
- **Description:** Define all colors for Matrix theme

**Acceptance Criteria:**
```
✓ internal/theme/colors.go created
✓ Primary colors defined (Matrix green variants)
✓ Accent colors defined (neon pink, cyan, yellow)
✓ Background colors defined (black, dark green)
✓ Color constants exported
```

**File:** `packages/tui-v2/internal/theme/colors.go`
```go
package theme

import "github.com/charmbracelet/lipgloss"

// Primary Matrix colors
var (
	MatrixGreen      = lipgloss.Color("#00ff00")
	MatrixGreenDim   = lipgloss.Color("#00dd00")
	MatrixGreenDark  = lipgloss.Color("#004400")
	MatrixGreenVDark = lipgloss.Color("#002200")
)

// Neon accents (cyberpunk)
var (
	NeonCyan   = lipgloss.Color("#00ffff")
	NeonPink   = lipgloss.Color("#ff3366")
	NeonYellow = lipgloss.Color("#ffaa00")
	NeonPurple = lipgloss.Color("#cc00ff")
)

// Backgrounds
var (
	Black       = lipgloss.Color("#000000")
	DarkGreen   = lipgloss.Color("#001100")
	DarkerGreen = lipgloss.Color("#000800")
)

// Semantic colors
var (
	ColorError   = NeonPink
	ColorWarning = NeonYellow
	ColorSuccess = MatrixGreen
	ColorInfo    = NeonCyan
)
```

---

#### TASK-010: Implement Theme System
- **Priority:** P0
- **Assignee:** Frontend Developer
- **Estimate:** 2 hours
- **Dependencies:** TASK-009
- **Description:** Create theme struct and styling helpers

**Acceptance Criteria:**
```
✓ internal/theme/theme.go created
✓ Theme struct with all styles
✓ MatrixTheme instance created
✓ Helper functions for common styles
✓ Unit tests for styling
```

**File:** `packages/tui-v2/internal/theme/theme.go`
```go
package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name string

	// Text styles
	Primary   lipgloss.Style
	Secondary lipgloss.Style
	Dim       lipgloss.Style
	Error     lipgloss.Style
	Success   lipgloss.Style
	Warning   lipgloss.Style
	Info      lipgloss.Style

	// UI elements
	Border    lipgloss.Style
	Highlight lipgloss.Style
	Selected  lipgloss.Style

	// Components
	Button       lipgloss.Style
	Input        lipgloss.Style
	CodeBlock    lipgloss.Style
	MessageUser  lipgloss.Style
	MessageAI    lipgloss.Style
}

var MatrixTheme = Theme{
	Name: "Matrix",

	Primary:   lipgloss.NewStyle().Foreground(MatrixGreen),
	Secondary: lipgloss.NewStyle().Foreground(MatrixGreenDim),
	Dim:       lipgloss.NewStyle().Foreground(MatrixGreenDark),
	Error:     lipgloss.NewStyle().Foreground(ColorError).Bold(true),
	Success:   lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true),
	Warning:   lipgloss.NewStyle().Foreground(ColorWarning),
	Info:      lipgloss.NewStyle().Foreground(ColorInfo),

	Border: lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(MatrixGreen),

	Highlight: lipgloss.NewStyle().
		Background(DarkGreen),

	Selected: lipgloss.NewStyle().
		Background(DarkGreen).
		Foreground(MatrixGreen).
		Bold(true),

	Button: lipgloss.NewStyle().
		Foreground(Black).
		Background(MatrixGreen).
		Padding(0, 2).
		Bold(true),

	Input: lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(MatrixGreen).
		Padding(0, 1),

	CodeBlock: lipgloss.NewStyle().
		Background(DarkerGreen).
		Padding(1).
		MarginTop(1).
		MarginBottom(1),

	MessageUser: lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(NeonCyan).
		Padding(1),

	MessageAI: lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(MatrixGreen).
		Padding(1),
}

// Helper functions
func (t Theme) Title(text string) string {
	return t.Primary.Bold(true).Render(text)
}

func (t Theme) Subtitle(text string) string {
	return t.Secondary.Render(text)
}

func (t Theme) Hint(text string) string {
	return t.Dim.Italic(true).Render(text)
}
```

---

#### TASK-011: Implement Gradient Text Helper
- **Priority:** P1
- **Assignee:** Frontend Developer
- **Estimate:** 2 hours
- **Dependencies:** TASK-010
- **Description:** Create gradient text effect

**Acceptance Criteria:**
```
✓ internal/theme/effects.go created
✓ GradientText function implemented
✓ Color interpolation working
✓ Supports horizontal gradients
✓ Unit tests with visual validation
```

**File:** `packages/tui-v2/internal/theme/effects.go`
```go
package theme

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

// GradientText creates gradient effect across text
func GradientText(text string, from, to lipgloss.Color) string {
	if len(text) == 0 {
		return ""
	}

	fromRGB := hexToRGB(string(from))
	toRGB := hexToRGB(string(to))

	result := ""
	for i, char := range text {
		progress := float64(i) / float64(len(text)-1)
		color := interpolateRGB(fromRGB, toRGB, progress)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(rgbToHex(color)))
		result += style.Render(string(char))
	}

	return result
}

type RGB struct {
	R, G, B int
}

func hexToRGB(hex string) RGB {
	var r, g, b int
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return RGB{r, g, b}
}

func rgbToHex(rgb RGB) string {
	return fmt.Sprintf("#%02x%02x%02x", rgb.R, rgb.G, rgb.B)
}

func interpolateRGB(from, to RGB, progress float64) RGB {
	return RGB{
		R: int(float64(from.R) + (float64(to.R)-float64(from.R))*progress),
		G: int(float64(from.G) + (float64(to.G)-float64(from.G))*progress),
		B: int(float64(from.B) + (float64(to.B)-float64(from.B))*progress),
	}
}

// GlowText simulates glow effect (limited in terminal)
func GlowText(text string, color lipgloss.Color, intensity float64) string {
	style := lipgloss.NewStyle().Foreground(color)
	if intensity > 0.5 {
		style = style.Bold(true)
	}
	return style.Render(text)
}
```

---

#### TASK-012: Create Theme Demo Command
- **Priority:** P1
- **Assignee:** Frontend Developer
- **Estimate:** 2 hours
- **Dependencies:** TASK-011
- **Description:** Showcase all theme elements

**Acceptance Criteria:**
```
✓ cmd/rycode/main.go updated with demo flag
✓ Displays all colors
✓ Shows all text styles
✓ Demonstrates gradient effect
✓ Shows UI components (buttons, borders)
✓ Runs with `rycode --demo`
```

---

### Sprint 1.4: Core Components (Week 2, Days 1-4)

#### TASK-013: Implement MessageBubble Component
- **Priority:** P0
- **Assignee:** Frontend Developer
- **Estimate:** 4 hours
- **Dependencies:** TASK-010
- **Description:** Chat message display component

**Acceptance Criteria:**
```
✓ internal/ui/components/message.go created
✓ MessageBubble struct defined
✓ Render method with responsive wrapping
✓ Markdown rendering via Glamour
✓ Syntax highlighting for code blocks
✓ Timestamp formatting
✓ Reaction support
✓ Unit tests with sample messages
```

[Continue with remaining 167 tasks...]

---

## Summary Statistics

### By Phase
- **Phase 1 (Foundation):** 35 tasks, 80 hours
- **Phase 2 (Mobile UX):** 40 tasks, 96 hours
- **Phase 3 (Visual):** 30 tasks, 72 hours
- **Phase 4 (AI Features):** 35 tasks, 88 hours
- **Phase 5 (Desktop):** 25 tasks, 64 hours
- **Phase 6 (Polish & Launch):** 15 tasks, 80 hours

### By Priority
- **P0 (Critical):** 45 tasks (25%)
- **P1 (High):** 75 tasks (42%)
- **P2 (Medium):** 45 tasks (25%)
- **P3 (Low):** 15 tasks (8%)

### By Role
- **Tech Lead:** 40 tasks
- **Frontend Developer:** 70 tasks
- **Backend Developer:** 50 tasks
- **QA/DevOps:** 20 tasks

---

## Task Tracking

**Recommended Tools:**
- GitHub Projects (Kanban board)
- Linear (for sprint planning)
- Notion (for documentation)

**Status Labels:**
- 🔴 Blocked
- 🟡 In Progress
- 🟢 Completed
- ⚪ Not Started

**Burn-down Chart Target:**
- Week 1: 15 tasks completed
- Week 2: 20 tasks completed
- Week 3-4: 40 tasks completed
- Week 5-6: 30 tasks completed
- Week 7-8: 35 tasks completed
- Week 9-10: 25 tasks completed
- Week 11-12: 15 tasks completed

---

## Next Steps

1. Import tasks into project management tool
2. Assign initial tasks to team members
3. Set up daily standup rhythm
4. Begin Sprint 1.1 (Project Setup)
5. Track progress and adjust estimates

**Let's execute!** 🚀
