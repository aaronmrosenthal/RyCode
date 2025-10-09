package debugger

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

type StackFrame struct {
	Name   string
	File   string
	Line   int
	Active bool
}

type CallStackViewModel struct {
	width  int
	height int
	frames []StackFrame
	theme  theme.Theme
}

func NewCallStackView(width, height int) CallStackViewModel {
	return CallStackViewModel{
		width:  width,
		height: height,
		theme:  theme.CurrentTheme(),
		frames: []StackFrame{},
	}
}

func (m CallStackViewModel) UpdateSize(width, height int) CallStackViewModel {
	m.width = width
	m.height = height
	return m
}

func (m CallStackViewModel) UpdateFrames(frames []StackFrame) CallStackViewModel {
	m.frames = frames
	return m
}

func (m CallStackViewModel) View() string {
	if len(m.frames) == 0 {
		// Show placeholder
		return styles.NewStyle().
			Foreground(m.theme.TextMuted()).
			Render("Waiting for stack trace...\n\nExample:\n  › calculateTotal() L67\n    processOrder() L45\n    main() L120")
	}

	t := m.theme

	var output strings.Builder

	for i, frame := range m.frames {
		// Active frame indicator
		if frame.Active {
			indicatorStyle := styles.NewStyle().
				Foreground(t.Primary()).
				Bold(true)
			output.WriteString(indicatorStyle.Render("› "))
		} else {
			output.WriteString("  ")
		}

		// Function name
		nameStyle := styles.NewStyle().
			Foreground(t.Text())
		if frame.Active {
			nameStyle = nameStyle.Bold(true)
		}

		output.WriteString(nameStyle.Render(frame.Name))
		output.WriteString(" ")

		// Location
		locationStyle := styles.NewStyle().
			Foreground(t.TextMuted())

		if frame.File != "" {
			output.WriteString(locationStyle.Render(fmt.Sprintf("L%d", frame.Line)))
		}

		if i < len(m.frames)-1 {
			output.WriteString("\n")
		}
	}

	return output.String()
}
