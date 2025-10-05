package main

import (
	"fmt"
	"os"

	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/ui/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Matrix theme colors
var (
	matrixGreen    = lipgloss.Color("#00ff00")
	matrixGreenDim = lipgloss.Color("#00dd00")
	black          = lipgloss.Color("#000000")
	darkGreen      = lipgloss.Color("#001100")
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(matrixGreen).
			Bold(true).
			MarginTop(1).
			MarginBottom(1)

	logoStyle = lipgloss.NewStyle().
			Foreground(matrixGreen).
			Bold(true)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(matrixGreenDim)

	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(matrixGreen).
			Padding(1, 2)

	hintStyle = lipgloss.NewStyle().
			Foreground(matrixGreenDim).
			Italic(true)
)

type model struct {
	width   int
	height  int
	message string
}

func initialModel() model {
	return model{
		message: "Hello, RyCode!",
		width:   80,
		height:  24,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	// ASCII logo
	logo := `
██████╗ ██╗   ██╗ ██████╗ ██████╗ ██████╗ ███████╗
██╔══██╗╚██╗ ██╔╝██╔════╝██╔═══██╗██╔══██╗██╔════╝
██████╔╝ ╚████╔╝ ██║     ██║   ██║██║  ██║█████╗
██╔══██╗  ╚██╔╝  ██║     ██║   ██║██║  ██║██╔══╝
██║  ██║   ██║   ╚██████╗╚██████╔╝██████╔╝███████╗
╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝`

	renderedLogo := logoStyle.Render(logo)

	// Title and subtitle
	title := titleStyle.Render("Matrix TUI v2.0")
	subtitle := subtitleStyle.Render("The AI-Native Evolution • Mobile-First Design")

	// Welcome message
	welcome := borderStyle.Render(
		fmt.Sprintf("Welcome to RyCode!\n\n%s\n\nTerminal Size: %d × %d", m.message, m.width, m.height),
	)

	// Features list
	features := lipgloss.NewStyle().
		Foreground(matrixGreen).
		Render(`
Features:
  🎨 Matrix Theme with Neon Effects
  📱 Mobile-First Responsive Design
  👆 Gesture-Based Navigation
  🎤 Voice Input Integration
  🤖 Multi-Agent AI Collaboration
  ⚡ 60 FPS Smooth Animations
`)

	// Hints
	hints := hintStyle.Render("\nPress 'q' or Ctrl+C to quit")

	// Compose the view
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		renderedLogo,
		title,
		subtitle,
		"",
		welcome,
		features,
		hints,
	)

	// Center content
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func main() {
	// Check for flags
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--demo":
			p := tea.NewProgram(
				demoModel(),
				tea.WithAltScreen(),
				tea.WithMouseCellMotion(),
			)
			if _, err := p.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return

		case "--chat":
			p := tea.NewProgram(
				chatModel(),
				tea.WithAltScreen(),
				tea.WithMouseCellMotion(),
			)
			if _, err := p.Run(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return

		case "--help", "-h":
			fmt.Println("RyCode Matrix TUI v2")
			fmt.Println("\nUsage:")
			fmt.Println("  rycode           Run default interface")
			fmt.Println("  rycode --demo    Show theme demo")
			fmt.Println("  rycode --chat    Interactive chat interface")
			fmt.Println("  rycode --help    Show this help")
			return
		}
	}

	// Create Bubble Tea program (default: chat)
	p := tea.NewProgram(
		chatModel(),
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// Demo mode model
type demoModelType struct {
	width  int
	height int
}

func demoModel() demoModelType {
	return demoModelType{
		width:  80,
		height: 24,
	}
}

func (m demoModelType) Init() tea.Cmd {
	return nil
}

func (m demoModelType) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" || msg.String() == "esc" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m demoModelType) View() string {
	return ShowDemo(m.width, m.height)
}

// Chat mode - returns the chat model
func chatModel() models.ChatModel {
	return models.NewChatModel()
}
