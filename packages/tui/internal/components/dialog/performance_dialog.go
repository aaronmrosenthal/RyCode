package dialog

import (
	"fmt"
	"strings"

	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/layout"
	"github.com/aaronmrosenthal/rycode/internal/performance"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/aaronmrosenthal/rycode/internal/typography"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// PerformanceDialog displays performance metrics and monitoring
type PerformanceDialog interface {
	layout.Modal
}

type performanceDialog struct {
	app     *app.App
	monitor *performance.PerformanceMonitor
	width   int
	height  int
}

// NewPerformanceDialog creates a new performance monitoring dialog
func NewPerformanceDialog(app *app.App) PerformanceDialog {
	return &performanceDialog{
		app:     app,
		monitor: performance.GetMonitor(),
	}
}

func (d *performanceDialog) Init() tea.Cmd {
	return nil
}

func (d *performanceDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height

	case tea.KeyPressMsg:
		switch msg.String() {
		case "r":
			// Reset metrics
			d.monitor.Reset()
		case "t":
			// Toggle monitoring
			if d.monitor.IsEnabled() {
				d.monitor.Disable()
			} else {
				d.monitor.Enable()
			}
		case "c":
			// Clear warnings
			d.monitor.ClearWarnings()
		}
	}

	return d, nil
}

func (d *performanceDialog) View() string {
	t := theme.CurrentTheme()
	typo := typography.New()

	var sections []string

	// Header
	header := styles.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		Render("⚡ Performance Monitor")

	sections = append(sections, header)
	sections = append(sections, "")

	// Health score
	health := d.monitor.GetHealth()
	healthViz := d.renderHealthScore(health)
	sections = append(sections, healthViz)
	sections = append(sections, "")

	// Frame metrics
	frameMetrics := d.monitor.GetFrameMetrics()
	frameSection := d.renderFrameMetrics(frameMetrics)
	sections = append(sections, typo.Subheading.Render("🎬 Frame Performance"))
	sections = append(sections, frameSection)
	sections = append(sections, "")

	// Memory metrics
	memMetrics := d.monitor.GetMemoryMetrics()
	memSection := d.renderMemoryMetrics(memMetrics)
	sections = append(sections, typo.Subheading.Render("💾 Memory Usage"))
	sections = append(sections, memSection)
	sections = append(sections, "")

	// Component metrics
	compMetrics := d.monitor.GetComponentMetrics()
	if len(compMetrics) > 0 {
		compSection := d.renderComponentMetrics(compMetrics)
		sections = append(sections, typo.Subheading.Render("🎨 Component Timings"))
		sections = append(sections, compSection)
		sections = append(sections, "")
	}

	// Warnings
	warnings := d.monitor.GetWarnings()
	if len(warnings) > 0 {
		warnSection := d.renderWarnings(warnings)
		sections = append(sections, typo.Subheading.Render("⚠️  Performance Warnings"))
		sections = append(sections, warnSection)
		sections = append(sections, "")
	}

	// Help footer
	helpStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Faint(true)

	help := helpStyle.Render("[r] Reset  [t] Toggle  [c] Clear warnings  [ESC] Close")
	sections = append(sections, help)

	content := strings.Join(sections, "\n")

	// Wrap in bordered box
	boxStyle := styles.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border()).
		Padding(2, 3).
		Width(d.width - 4)

	return boxStyle.Render(content)
}

// renderHealthScore creates a health score visualization
func (d *performanceDialog) renderHealthScore(health float64) string {
	t := theme.CurrentTheme()

	// Health bar
	barWidth := 40
	filledWidth := int((health / 100.0) * float64(barWidth))

	bar := ""
	healthColor := t.Success()
	if health < 70 {
		healthColor = t.Warning()
	}
	if health < 40 {
		healthColor = t.Error()
	}

	for i := 0; i < filledWidth; i++ {
		bar += styles.NewStyle().Foreground(healthColor).Render("█")
	}
	for i := filledWidth; i < barWidth; i++ {
		bar += styles.NewStyle().Foreground(t.TextMuted()).Faint(true).Render("░")
	}

	// Score label
	scoreStyle := styles.NewStyle().
		Foreground(healthColor).
		Bold(true)

	score := scoreStyle.Render(fmt.Sprintf("%.0f%%", health))

	// Status text
	status := "Excellent"
	if health < 70 {
		status = "Good"
	}
	if health < 40 {
		status = "Poor"
	}

	statusStyle := styles.NewStyle().
		Foreground(t.TextMuted())

	return fmt.Sprintf("Health: %s %s   %s", score, bar, statusStyle.Render(status))
}

// renderFrameMetrics creates frame performance visualization
func (d *performanceDialog) renderFrameMetrics(metrics performance.FrameMetrics) string {
	t := theme.CurrentTheme()

	var lines []string

	// FPS
	fpsColor := t.Success()
	if metrics.FPS < 50 {
		fpsColor = t.Warning()
	}
	if metrics.FPS < 30 {
		fpsColor = t.Error()
	}

	fpsStyle := styles.NewStyle().
		Foreground(fpsColor).
		Bold(true)

	fps := fpsStyle.Render(fmt.Sprintf("%.1f FPS", metrics.FPS))

	labelStyle := styles.NewStyle().
		Foreground(t.TextMuted())

	lines = append(lines, fps+"  "+labelStyle.Render("(target: 60 FPS)"))

	// Frame time
	avgMs := float64(metrics.AverageFrameTime.Microseconds()) / 1000.0
	lastMs := float64(metrics.FrameTime.Microseconds()) / 1000.0

	frameTimeStyle := styles.NewStyle().
		Foreground(t.Text())

	lines = append(lines, frameTimeStyle.Render(
		fmt.Sprintf("Average: %.2fms  Last: %.2fms  Target: 16.67ms", avgMs, lastMs)))

	// Dropped frames
	dropRate := 0.0
	if metrics.TotalFrames > 0 {
		dropRate = float64(metrics.DroppedFrames) / float64(metrics.TotalFrames) * 100.0
	}

	dropColor := t.Success()
	if dropRate > 5 {
		dropColor = t.Warning()
	}
	if dropRate > 10 {
		dropColor = t.Error()
	}

	dropStyle := styles.NewStyle().
		Foreground(dropColor)

	lines = append(lines, dropStyle.Render(
		fmt.Sprintf("Dropped: %d/%d (%.1f%%)", metrics.DroppedFrames, metrics.TotalFrames, dropRate)))

	return strings.Join(lines, "\n")
}

// renderMemoryMetrics creates memory usage visualization
func (d *performanceDialog) renderMemoryMetrics(metrics performance.MemoryMetrics) string {
	t := theme.CurrentTheme()

	var lines []string

	// Allocated memory
	allocMB := float64(metrics.Alloc) / (1024 * 1024)
	sysMB := float64(metrics.Sys) / (1024 * 1024)

	memColor := t.Success()
	if allocMB > 50 {
		memColor = t.Warning()
	}
	if allocMB > 100 {
		memColor = t.Error()
	}

	memStyle := styles.NewStyle().
		Foreground(memColor).
		Bold(true)

	labelStyle := styles.NewStyle().
		Foreground(t.TextMuted())

	lines = append(lines, memStyle.Render(fmt.Sprintf("%.2f MB", allocMB))+
		" "+labelStyle.Render(fmt.Sprintf("allocated (%.2f MB system)", sysMB)))

	// Heap objects
	lines = append(lines, labelStyle.Render(
		fmt.Sprintf("Heap objects: %s", formatNumber(metrics.HeapObjects))))

	// GC info
	gcPauseMs := float64(metrics.LastGCPause.Microseconds()) / 1000.0

	gcColor := t.Success()
	if gcPauseMs > 5 {
		gcColor = t.Warning()
	}
	if gcPauseMs > 10 {
		gcColor = t.Error()
	}

	gcStyle := styles.NewStyle().
		Foreground(gcColor)

	lines = append(lines, gcStyle.Render(
		fmt.Sprintf("GC runs: %d  Last pause: %.2fms", metrics.NumGC, gcPauseMs)))

	return strings.Join(lines, "\n")
}

// renderComponentMetrics creates component timing table
func (d *performanceDialog) renderComponentMetrics(metrics map[string]*performance.RenderMetrics) string {
	t := theme.CurrentTheme()

	var lines []string

	// Sort by average render time
	type sortable struct {
		name    string
		metrics *performance.RenderMetrics
	}

	var sorted []sortable
	for name, m := range metrics {
		sorted = append(sorted, sortable{name, m})
	}

	// Simple bubble sort
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].metrics.AverageRenderTime < sorted[j].metrics.AverageRenderTime {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Show top 5
	limit := 5
	if len(sorted) < limit {
		limit = len(sorted)
	}

	for i := 0; i < limit; i++ {
		item := sorted[i]

		avgMs := float64(item.metrics.AverageRenderTime.Microseconds()) / 1000.0
		lastMs := float64(item.metrics.RenderTime.Microseconds()) / 1000.0

		// Color based on speed
		color := t.Success()
		if avgMs > 5 {
			color = t.Warning()
		}
		if avgMs > 10 {
			color = t.Error()
		}

		nameStyle := styles.NewStyle().
			Foreground(t.Text()).
			Bold(true)

		timeStyle := styles.NewStyle().
			Foreground(color)

		countStyle := styles.NewStyle().
			Foreground(t.TextMuted()).
			Faint(true)

		line := fmt.Sprintf("%-30s %s  %s",
			nameStyle.Render(item.name),
			timeStyle.Render(fmt.Sprintf("%.2fms avg", avgMs)),
			countStyle.Render(fmt.Sprintf("(%.2fms last, %d renders)", lastMs, item.metrics.RenderCount)))

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// renderWarnings creates warnings list
func (d *performanceDialog) renderWarnings(warnings []string) string {
	t := theme.CurrentTheme()

	var lines []string

	for _, warning := range warnings {
		warnStyle := styles.NewStyle().
			Foreground(t.Warning())

		lines = append(lines, warnStyle.Render("⚠ "+warning))
	}

	return strings.Join(lines, "\n")
}

func (d *performanceDialog) Render(background string) string {
	return d.View()
}

func (d *performanceDialog) Close() tea.Cmd {
	return nil
}

// formatNumber formats large numbers with commas
func formatNumber(n uint64) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}

	var result strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(c)
	}

	return result.String()
}
