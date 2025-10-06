package theme

import "github.com/charmbracelet/lipgloss"

// Matrix Theme Primary Colors
var (
	MatrixGreen       = lipgloss.Color("#00ff00") // Primary Matrix green
	MatrixGreenBright = lipgloss.Color("#00ff88") // Brighter variant
	MatrixGreenDim    = lipgloss.Color("#00dd00") // Dimmed variant
	MatrixGreenDark   = lipgloss.Color("#004400") // Dark variant
	MatrixGreenVDark  = lipgloss.Color("#002200") // Very dark variant
)

// Neon Cyberpunk Accent Colors
var (
	NeonCyan   = lipgloss.Color("#00ffff") // Bright cyan
	NeonPink   = lipgloss.Color("#ff3366") // Hot pink
	NeonPurple = lipgloss.Color("#cc00ff") // Electric purple
	NeonYellow = lipgloss.Color("#ffaa00") // Amber yellow
	NeonOrange = lipgloss.Color("#ff6600") // Vibrant orange
	NeonBlue   = lipgloss.Color("#0088ff") // Electric blue
)

// Background Colors
var (
	Black        = lipgloss.Color("#000000") // Pure black
	DarkGreen    = lipgloss.Color("#001100") // Very dark green
	DarkerGreen  = lipgloss.Color("#000800") // Nearly black green
	DarkestGreen = lipgloss.Color("#000400") // Almost black
)

// Semantic Colors
var (
	ColorError     = NeonPink       // Error states
	ColorWarning   = NeonYellow     // Warning states
	ColorSuccess   = MatrixGreen    // Success states
	ColorInfo      = NeonCyan       // Informational
	ColorPrimary   = MatrixGreen    // Primary actions
	ColorSecondary = MatrixGreenDim // Secondary actions
)

// UI Element Colors
var (
	ColorBorder     = MatrixGreen     // Borders
	ColorText       = MatrixGreen     // Primary text
	ColorTextDim    = MatrixGreenDim  // Secondary text
	ColorTextHint   = MatrixGreenDark // Hint text
	ColorBackground = Black           // Main background
	ColorSurface    = DarkGreen       // Surface/panels
	ColorHighlight  = DarkGreen       // Highlighted elements
	ColorSelection  = DarkerGreen     // Selected items
)

// Code Syntax Colors (for syntax highlighting)
var (
	SyntaxKeyword  = NeonPink        // Keywords
	SyntaxString   = NeonYellow      // Strings
	SyntaxNumber   = NeonCyan        // Numbers
	SyntaxComment  = MatrixGreenDark // Comments
	SyntaxFunction = NeonBlue        // Functions
	SyntaxType     = NeonPurple      // Types
	SyntaxOperator = MatrixGreen     // Operators
)

// Gradient Presets
type GradientPreset struct {
	From lipgloss.Color
	To   lipgloss.Color
}

var (
	GradientMatrix = GradientPreset{
		From: MatrixGreen,
		To:   NeonCyan,
	}
	GradientFire = GradientPreset{
		From: NeonOrange,
		To:   NeonPink,
	}
	GradientCool = GradientPreset{
		From: NeonCyan,
		To:   NeonPurple,
	}
	GradientWarm = GradientPreset{
		From: NeonYellow,
		To:   NeonOrange,
	}
)

// AI Provider Brand Colors
var (
	ClaudeBlue    = lipgloss.Color("#5B8DEF") // Claude brand blue
	ClaudeCyan    = lipgloss.Color("#00D4FF") // Claude accent cyan
	OpenAIMagenta = lipgloss.Color("#FF006E") // OpenAI brand magenta
	OpenAIGreen   = lipgloss.Color("#10A37F") // OpenAI accent green
)

// ProviderColors maps provider names to their brand colors
var ProviderColors = map[string]lipgloss.Color{
	"claude":        ClaudeBlue,
	"claude-opus-4": ClaudeBlue,
	"anthropic":     ClaudeBlue,
	"openai":        OpenAIMagenta,
	"gpt-4":         OpenAIMagenta,
	"gpt-4o":        OpenAIMagenta,
}

// ProviderIcons maps provider names to their display icons
var ProviderIcons = map[string]string{
	"claude":        "🤖",
	"claude-opus-4": "🤖",
	"anthropic":     "🤖",
	"openai":        "🧠",
	"gpt-4":         "🧠",
	"gpt-4o":        "🧠",
}

// GetProviderColor returns the brand color for a provider
func GetProviderColor(provider string) lipgloss.Color {
	if color, ok := ProviderColors[provider]; ok {
		return color
	}
	return MatrixGreen // Default fallback
}

// GetProviderIcon returns the icon for a provider
func GetProviderIcon(provider string) string {
	if icon, ok := ProviderIcons[provider]; ok {
		return icon
	}
	return "🤖" // Default fallback
}
