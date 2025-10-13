package dialog

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/aaronmrosenthal/rycode-sdk-go"
	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/components/modal"
	"github.com/aaronmrosenthal/rycode/internal/layout"
)

var modelsDebugLog *os.File

func init() {
	var err error
	// Use owner-only permissions (0600) for security
	modelsDebugLog, err = os.OpenFile("/tmp/rycode-debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		modelsDebugLog = nil
	}
}

func logModelsDebug(format string, args ...interface{}) {
	if modelsDebugLog != nil {
		fmt.Fprintf(modelsDebugLog, format+"\n", args...)
		modelsDebugLog.Sync()
	}
}

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
	dialogWidth  int
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
	return m.modal.Render(m.View(), background)
}

func (s *modelDialog) Close() tea.Cmd {
	return nil
}

func NewModelDialog(app *app.App) ModelDialog {
	logModelsDebug("=== NewModelDialog() called ===")

	dialog := &modelDialog{
		app:          app,
		simpleToggle: NewSimpleProviderToggle(app),
		dialogWidth:  60,
	}

	dialog.modal = modal.New(
		modal.WithTitle(""),
		modal.WithMaxWidth(dialog.dialogWidth+4),
	)

	return dialog
}
