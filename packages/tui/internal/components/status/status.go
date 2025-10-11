package status

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
	"github.com/fsnotify/fsnotify"
	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/commands"
	"github.com/aaronmrosenthal/rycode/internal/layout"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
	"github.com/aaronmrosenthal/rycode/internal/util"
)

type GitBranchUpdatedMsg struct {
	Branch string
}

type StatusComponent interface {
	tea.Model
	tea.ViewModel
	Cleanup()
}

type statusComponent struct {
	app        *app.App
	width      int
	cwd        string
	branch     string
	watcher    *fsnotify.Watcher
	done       chan struct{}
	lastUpdate time.Time
}

func (m *statusComponent) Init() tea.Cmd {
	return m.startGitWatcher()
}

func (m *statusComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case GitBranchUpdatedMsg:
		if m.branch != msg.Branch {
			m.branch = msg.Branch
		}
		// Continue watching for changes (persistent watcher)
		return m, m.watchForGitChanges()
	}
	return m, nil
}

func (m *statusComponent) logo() string {
	t := theme.CurrentTheme()

	// Bright neon green for "Ry"
	ryStyle := styles.NewStyle().
		Foreground(compat.AdaptiveColor{
			Dark:  lipgloss.Color("#00FFAA"),
			Light: lipgloss.Color("#00CC88"),
		}).
		Background(t.BackgroundElement()).
		Bold(true).
		Render

	// Medium neon green for "Code"
	codeStyle := styles.NewStyle().
		Foreground(compat.AdaptiveColor{
			Dark:  lipgloss.Color("#00CC88"),
			Light: lipgloss.Color("#008866"),
		}).
		Background(t.BackgroundElement()).
		Bold(true).
		Render

	versionStyle := styles.NewStyle().
		Foreground(t.TextMuted()).
		Background(t.BackgroundElement()).
		Render

	ry := ryStyle("Ry")
	code := codeStyle("Code")
	version := versionStyle(" " + m.app.Version)

	content := ry + code
	if m.width > 40 {
		content += version
	}
	return styles.NewStyle().
		Background(t.BackgroundElement()).
		Padding(0, 1).
		Render(content)
}

func (m *statusComponent) collapsePath(path string, maxWidth int) string {
	if lipgloss.Width(path) <= maxWidth {
		return path
	}

	const ellipsis = ".."
	ellipsisLen := len(ellipsis)

	if maxWidth <= ellipsisLen {
		if maxWidth > 0 {
			return "..."[:maxWidth]
		}
		return ""
	}

	separator := string(filepath.Separator)
	parts := strings.Split(path, separator)

	if len(parts) == 1 {
		return path[:maxWidth-ellipsisLen] + ellipsis
	}

	truncatedPath := parts[len(parts)-1]
	for i := len(parts) - 2; i >= 0; i-- {
		part := parts[i]
		if len(truncatedPath)+len(separator)+len(part)+ellipsisLen > maxWidth {
			return ellipsis + separator + truncatedPath
		}
		truncatedPath = part + separator + truncatedPath
	}
	return truncatedPath
}

// getProviderBrandColor returns the brand color for a given provider
func getProviderBrandColor(providerName string) compat.AdaptiveColor {
	// Normalize provider name to lowercase for comparison
	provider := strings.ToLower(providerName)

	switch {
	case strings.Contains(provider, "anthropic") || strings.Contains(provider, "claude"):
		// Anthropic Claude - Orange/Coral
		return compat.AdaptiveColor{
			Dark:  lipgloss.Color("#D97757"), // Warm coral
			Light: lipgloss.Color("#B85C3C"),
		}
	case strings.Contains(provider, "openai") || strings.Contains(provider, "gpt"):
		// OpenAI - Teal/Turquoise (like GPT branding)
		return compat.AdaptiveColor{
			Dark:  lipgloss.Color("#10A37F"), // OpenAI green
			Light: lipgloss.Color("#0E8C6E"),
		}
	case strings.Contains(provider, "google") || strings.Contains(provider, "gemini"):
		// Google Gemini - Blue/Purple gradient (using blue)
		return compat.AdaptiveColor{
			Dark:  lipgloss.Color("#4285F4"), // Google blue
			Light: lipgloss.Color("#3367D6"),
		}
	case strings.Contains(provider, "grok") || strings.Contains(provider, "x.ai"):
		// Grok (X.AI) - Dark gray/black (X branding)
		return compat.AdaptiveColor{
			Dark:  lipgloss.Color("#71767B"), // X gray
			Light: lipgloss.Color("#536471"),
		}
	case strings.Contains(provider, "qwen") || strings.Contains(provider, "alibaba"):
		// Qwen (Alibaba) - Orange
		return compat.AdaptiveColor{
			Dark:  lipgloss.Color("#FF6A00"), // Alibaba orange
			Light: lipgloss.Color("#E65C00"),
		}
	default:
		// Default - neutral gray
		return compat.AdaptiveColor{
			Dark:  lipgloss.Color("#6B7280"),
			Light: lipgloss.Color("#4B5563"),
		}
	}
}

func (m *statusComponent) buildModelDisplay() string {
	t := theme.CurrentTheme()

	// Check if we have a model selected
	if m.app.Model == nil || m.app.Provider == nil {
		faintStyle := styles.NewStyle().
			Faint(true).
			Background(t.BackgroundPanel()).
			Foreground(t.TextMuted())
		noModelStyle := styles.NewStyle().
			Background(t.BackgroundElement()).
			Padding(0, 1).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(t.BackgroundElement()).
			BorderBackground(t.BackgroundPanel())
		return faintStyle.Render("  ") + noModelStyle.Render("No model")
	}

	// Get provider brand color
	brandColor := getProviderBrandColor(m.app.Provider.Name)

	// Get cost (from cached value)
	costStr := fmt.Sprintf("💰 $%.2f", m.app.CurrentCost)

	// Check if cost data is stale (>10 seconds old)
	if time.Since(m.app.LastCostUpdate) > 10*time.Second {
		costStr = "💰 $--"
	}

	// Style definitions with brand color background
	modelNameStyle := styles.NewStyle().
		Background(brandColor).
		Foreground(compat.AdaptiveColor{
			Dark:  lipgloss.Color("#FFFFFF"), // White text on dark bg
			Light: lipgloss.Color("#FFFFFF"), // White text on light bg
		}).
		Bold(true).
		Render

	costStyle := styles.NewStyle().
		Background(brandColor).
		Foreground(compat.AdaptiveColor{
			Dark:  lipgloss.Color("#FFFFFF"),
			Light: lipgloss.Color("#FFFFFF"),
		}).
		Render

	hintStyle := styles.NewStyle().
		Background(brandColor).
		Foreground(compat.AdaptiveColor{
			Dark:  lipgloss.Color("#E5E5E5"),
			Light: lipgloss.Color("#F0F0F0"),
		}).
		Faint(true).
		Render

	separatorStyle := styles.NewStyle().
		Background(brandColor).
		Foreground(compat.AdaptiveColor{
			Dark:  lipgloss.Color("#E5E5E5"),
			Light: lipgloss.Color("#F0F0F0"),
		}).
		Faint(true).
		Render

	// Get keybinding for cycling
	command := m.app.Commands[commands.AgentCycleCommand]
	kb := command.Keybindings[0]
	key := kb.Key
	if kb.RequiresLeader {
		key = m.app.Config.Keybinds.Leader + " " + kb.Key
	}

	// Build display: "Model Name | 💰 $0.12 | tab→"
	modelName := modelNameStyle(m.app.Model.Name)
	cost := costStyle(costStr)
	hint := hintStyle(key + "→")
	separator := separatorStyle(" | ")

	var content string
	if m.width > 80 {
		// Full display with all info
		content = modelName + separator + cost + separator + hint
	} else if m.width > 60 {
		// Hide hint, show model and cost
		content = modelName + separator + cost
	} else {
		// Just show model name
		content = modelName
	}

	// Add border and padding with brand color
	faintStyle := styles.NewStyle().
		Faint(true).
		Background(t.BackgroundPanel()).
		Foreground(t.TextMuted())

	displayStyle := styles.NewStyle().
		Background(brandColor).
		Padding(0, 1).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(brandColor).
		BorderBackground(t.BackgroundPanel())

	return faintStyle.Render(key+" ") + displayStyle.Render(content)
}

func (m *statusComponent) View() string {
	t := theme.CurrentTheme()
	logo := m.logo()
	logoWidth := lipgloss.Width(logo)

	// Build model display instead of agent display
	modelDisplay := m.buildModelDisplay()
	modelWidth := lipgloss.Width(modelDisplay)

	availableWidth := m.width - logoWidth - modelWidth
	branchSuffix := ""
	if m.branch != "" {
		branchSuffix = ":" + m.branch
	}

	maxCwdWidth := availableWidth - lipgloss.Width(branchSuffix)
	cwdDisplay := m.collapsePath(m.cwd, maxCwdWidth)

	faintStyle := styles.NewStyle().
		Faint(true).
		Background(t.BackgroundPanel()).
		Foreground(t.TextMuted())

	if m.branch != "" && availableWidth > lipgloss.Width(cwdDisplay)+lipgloss.Width(branchSuffix) {
		cwdDisplay += faintStyle.Render(branchSuffix)
	}

	cwd := styles.NewStyle().
		Foreground(t.TextMuted()).
		Background(t.BackgroundPanel()).
		Padding(0, 1).
		Render(cwdDisplay)

	background := t.BackgroundPanel()
	status := layout.Render(
		layout.FlexOptions{
			Background: &background,
			Direction:  layout.Row,
			Justify:    layout.JustifySpaceBetween,
			Align:      layout.AlignStretch,
			Width:      m.width,
		},
		layout.FlexItem{
			View: logo + cwd,
		},
		layout.FlexItem{
			View: modelDisplay,
		},
	)

	blank := styles.NewStyle().Background(t.Background()).Width(m.width).Render("")
	return blank + "\n" + status
}

func (m *statusComponent) startGitWatcher() tea.Cmd {
	cmd := util.CmdHandler(
		GitBranchUpdatedMsg{Branch: getCurrentGitBranch(util.CwdPath)},
	)
	if err := m.initWatcher(); err != nil {
		return cmd
	}
	return tea.Batch(cmd, m.watchForGitChanges())
}

func (m *statusComponent) initWatcher() error {
	gitDir := filepath.Join(util.CwdPath, ".git")
	headFile := filepath.Join(gitDir, "HEAD")
	if info, err := os.Stat(gitDir); err != nil || !info.IsDir() {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err := watcher.Add(headFile); err != nil {
		watcher.Close()
		return err
	}

	// Also watch the ref file if HEAD points to a ref
	refFile := getGitRefFile(util.CwdPath)
	if refFile != headFile && refFile != "" {
		if _, err := os.Stat(refFile); err == nil {
			watcher.Add(refFile) // Ignore error, HEAD watching is sufficient
		}
	}

	m.watcher = watcher
	m.done = make(chan struct{})
	return nil
}

func (m *statusComponent) watchForGitChanges() tea.Cmd {
	if m.watcher == nil {
		return nil
	}

	return tea.Cmd(func() tea.Msg {
		for {
			select {
			case event, ok := <-m.watcher.Events:
				branch := getCurrentGitBranch(util.CwdPath)
				if !ok {
					return GitBranchUpdatedMsg{Branch: branch}
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					// Debounce updates to prevent excessive refreshes
					now := time.Now()
					if now.Sub(m.lastUpdate) < 100*time.Millisecond {
						continue
					}
					m.lastUpdate = now
					if strings.HasSuffix(event.Name, "HEAD") {
						m.updateWatchedFiles()
					}
					return GitBranchUpdatedMsg{Branch: branch}
				}
			case <-m.watcher.Errors:
				// Continue watching even on errors
			case <-m.done:
				return GitBranchUpdatedMsg{Branch: ""}
			}
		}
	})
}

func (m *statusComponent) updateWatchedFiles() {
	if m.watcher == nil {
		return
	}
	refFile := getGitRefFile(util.CwdPath)
	headFile := filepath.Join(util.CwdPath, ".git", "HEAD")
	if refFile != headFile && refFile != "" {
		if _, err := os.Stat(refFile); err == nil {
			// Try to add the new ref file (ignore error if already watching)
			m.watcher.Add(refFile)
		}
	}
}

func getCurrentGitBranch(cwd string) string {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = cwd
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func getGitRefFile(cwd string) string {
	headFile := filepath.Join(cwd, ".git", "HEAD")
	content, err := os.ReadFile(headFile)
	if err != nil {
		return ""
	}

	headContent := strings.TrimSpace(string(content))
	if after, ok := strings.CutPrefix(headContent, "ref: "); ok {
		// HEAD points to a ref file
		refPath := after
		return filepath.Join(cwd, ".git", refPath)
	}

	// HEAD contains a direct commit hash
	return headFile
}

func (m *statusComponent) Cleanup() {
	if m.done != nil {
		close(m.done)
	}
	if m.watcher != nil {
		m.watcher.Close()
	}
}

func NewStatusCmp(app *app.App) StatusComponent {
	statusComponent := &statusComponent{
		app:        app,
		lastUpdate: time.Now(),
	}

	homePath, err := os.UserHomeDir()
	cwdPath := util.CwdPath
	if err == nil && homePath != "" && strings.HasPrefix(cwdPath, homePath) {
		cwdPath = "~" + cwdPath[len(homePath):]
	}
	statusComponent.cwd = cwdPath

	return statusComponent
}
