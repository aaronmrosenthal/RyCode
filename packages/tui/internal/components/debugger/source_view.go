package debugger

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

type SourceViewModel struct {
	width       int
	height      int
	currentFile string
	currentLine int
	lines       []string
	theme       theme.Theme
}

func NewSourceView(width, height int) SourceViewModel {
	return SourceViewModel{
		width:  width,
		height: height,
		theme:  theme.CurrentTheme(),
		lines:  []string{},
	}
}

func (m SourceViewModel) UpdateCurrentLine(file string, line int) SourceViewModel {
	m.currentFile = file
	m.currentLine = line

	// Load the file content
	content, err := os.ReadFile(file)
	if err != nil {
		m.lines = []string{fmt.Sprintf("Error reading file: %v", err)}
		return m
	}

	m.lines = strings.Split(string(content), "\n")
	return m
}

func (m SourceViewModel) UpdateSize(width, height int) SourceViewModel {
	m.width = width
	m.height = height
	return m
}

func (m SourceViewModel) View() string {
	if len(m.lines) == 0 {
		return styles.NewStyle().
			Foreground(m.theme.TextMuted()).
			Render("No source loaded")
	}

	t := m.theme

	// Calculate visible range
	contextLines := (m.height - 2) / 2
	startLine := max(1, m.currentLine-contextLines)
	endLine := min(len(m.lines), m.currentLine+contextLines)

	var output strings.Builder

	for i := startLine; i <= endLine; i++ {
		lineNum := i
		lineContent := ""
		if lineNum <= len(m.lines) {
			lineContent = m.lines[lineNum-1]
		}

		// Line number style
		lineNumStyle := styles.NewStyle().
			Foreground(t.TextMuted()).
			Width(4).
			Align(lipgloss.Right)

		// Is this the current line?
		if lineNum == m.currentLine {
			// Highlight current execution line
			currentLineStyle := styles.NewStyle().
				Background(t.Primary()).
				Foreground(t.Background()).
				Bold(true)

			arrowStyle := styles.NewStyle().
				Foreground(t.Primary()).
				Bold(true)

			output.WriteString(arrowStyle.Render("►"))
			output.WriteString(currentLineStyle.Render(fmt.Sprintf("%3d", lineNum)))
			output.WriteString(" ")
			output.WriteString(currentLineStyle.Render(lineContent))
		} else {
			// Normal line
			lineStyle := styles.NewStyle().
				Foreground(t.Text())

			output.WriteString(" ")
			output.WriteString(lineNumStyle.Render(fmt.Sprintf("%3d", lineNum)))
			output.WriteString(" ")
			output.WriteString(lineStyle.Render(lineContent))
		}

		if i < endLine {
			output.WriteString("\n")
		}
	}

	return output.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
