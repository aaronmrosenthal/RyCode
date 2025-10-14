package dialog

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/aaronmrosenthal/rycode-sdk-go"
	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/components/modal"
	"github.com/aaronmrosenthal/rycode/internal/layout"
)

// ModelDialog interface for the model selection dialog
type ModelDialog interface {
	layout.Modal
}

type modelDialog struct {
	app          *app.App
	width        int
	height       int
	modal        *modal.Modal
	simpleToggle *SimpleProviderToggle
}

type ModelWithProvider struct {
	Model    opencode.Model
	Provider opencode.Provider
}

func (m *modelDialog) Init() tea.Cmd {
	return m.simpleToggle.Init()
}

func (m *modelDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.simpleToggle.SetSize(msg.Width, msg.Height)
		return m, nil
	}

	// Pass all messages to simple toggle
	updatedToggle, cmd := m.simpleToggle.Update(msg)
	if toggle, ok := updatedToggle.(*SimpleProviderToggle); ok {
		m.simpleToggle = toggle
	}
	return m, cmd
}

func (m *modelDialog) View() string {
	return m.simpleToggle.View()
}

func (m *modelDialog) Render(background string) string {
	// Show cortex overlay FIRST (not in modal) while loading providers
	// Once loaded, fade in the modal content in place of the cortex
	if m.simpleToggle.IsLoading() {
		// Render cortex directly on background (no modal box)
		// This provides instant visual feedback while data loads
		return m.simpleToggle.View()
	}

	// Once loaded, render the full modal with provider selection content
	return m.modal.Render(m.View(), background)
}

func (s *modelDialog) Close() tea.Cmd {
	return nil
}

func NewModelDialog(app *app.App) ModelDialog {
	simpleToggle := NewSimpleProviderToggle(app)

	// Calculate dynamic width based on provider count and names
	// Will be recalculated after providers load, but start with reasonable default
	dialogWidth := 80

	dialog := &modelDialog{
		app:          app,
		simpleToggle: simpleToggle,
	}

	dialog.modal = modal.New(
		modal.WithTitle(""),
		modal.WithMaxWidth(dialogWidth+4), // Add padding for modal border
	)

	return dialog
}
