package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// MatrixLogo is the ASCII art logo for RyCode
const MatrixLogo = `
╦═╗┬ ┬╔═╗┌─┐┌┬┐┌─┐
╠╦╝└┬┘║  │ │ ││├┤
╩╚═ ┴ ╚═╝└─┘─┴┘└─┘`

// MatrixLogoSmall is a compact version for small screens
const MatrixLogoSmall = `
╦═╗┬ ┬╔═╗┌─┐┌┬┐┌─┐
╩╚═┴ ╚═╝└─┘─┴┘└─┘`

// MatrixLogoMini is the smallest version for tiny screens
const MatrixLogoMini = `RyCode`

// RenderLogo renders the ASCII logo with optional animation
func RenderLogo(animated bool, frame int, width int) string {
	// Choose logo size based on width
	logo := MatrixLogo
	if width < 60 {
		logo = MatrixLogoSmall
	}
	if width < 40 {
		logo = MatrixLogoMini
	}

	if animated {
		// Rainbow animation
		return RainbowText(logo)
	}

	// Static gradient
	return GradientTextPreset(logo, GradientMatrix)
}

// RenderLogoWithTagline renders logo with tagline
func RenderLogoWithTagline(animated bool, frame int, width int) string {
	logo := RenderLogo(animated, frame, width)

	// Don't show tagline on very small screens
	if width < 60 {
		return logo
	}

	tagline := "The AI-Native Terminal IDE"
	taglineStyle := lipgloss.NewStyle().
		Foreground(MatrixGreenDim).
		Italic(true)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		taglineStyle.Render(tagline),
	)
}

// RenderLogoBordered renders logo with border
func RenderLogoBordered(animated bool, frame int, width int) string {
	content := RenderLogoWithTagline(animated, frame, width)

	// Calculate border width
	borderWidth := width
	if borderWidth > 80 {
		borderWidth = 80
	}

	style := lipgloss.NewStyle().
		Width(borderWidth).
		Align(lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(MatrixGreen).
		Padding(1, 2)

	return style.Render(content)
}
