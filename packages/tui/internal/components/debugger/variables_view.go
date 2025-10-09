package debugger

import (
	"fmt"
	"strings"

	"github.com/aaronmrosenthal/rycode/internal/styles"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

type Variable struct {
	Name  string
	Value string
	Type  string
}

type VariablesViewModel struct {
	width     int
	height    int
	variables []Variable
	theme     theme.Theme
}

func NewVariablesView(width, height int) VariablesViewModel {
	return VariablesViewModel{
		width:     width,
		height:    height,
		theme:     theme.CurrentTheme(),
		variables: []Variable{},
	}
}

func (m VariablesViewModel) UpdateSize(width, height int) VariablesViewModel {
	m.width = width
	m.height = height
	return m
}

func (m VariablesViewModel) UpdateVariables(vars []Variable) VariablesViewModel {
	m.variables = vars
	return m
}

func (m VariablesViewModel) View() string {
	if len(m.variables) == 0 {
		// Show placeholder with example
		return styles.NewStyle().
			Foreground(m.theme.TextMuted()).
			Render("Waiting for variable data from debugger...\n\nExample:\n  user: { id: 123, name: \"John\" }\n  total: undefined ⚠️\n  items: Array(3)")
	}

	t := m.theme

	var output strings.Builder

	for i, v := range m.variables {
		// Variable name
		nameStyle := styles.NewStyle().
			Foreground(t.Primary()).
			Bold(true)

		// Value style depends on type
		valueStyle := styles.NewStyle().
			Foreground(t.Text())

		// Special styling for undefined/null
		if v.Value == "undefined" || v.Value == "null" {
			valueStyle = valueStyle.
				Foreground(t.Warning()).
				Bold(true)
			v.Value += " ⚠️"
		}

		// Type hint
		typeStyle := styles.NewStyle().
			Foreground(t.TextMuted()).
			Italic(true)

		output.WriteString(nameStyle.Render(v.Name))
		output.WriteString(": ")
		output.WriteString(valueStyle.Render(v.Value))

		if v.Type != "" {
			output.WriteString(" ")
			output.WriteString(typeStyle.Render(fmt.Sprintf("(%s)", v.Type)))
		}

		if i < len(m.variables)-1 {
			output.WriteString("\n")
		}
	}

	return output.String()
}
