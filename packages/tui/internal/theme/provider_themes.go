package theme

import (
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
)

// ProviderTheme extends BaseTheme with provider-specific branding
type ProviderTheme struct {
	BaseTheme

	// Provider metadata
	ProviderID   string
	ProviderName string

	// Provider-specific UI elements
	LogoASCII       string
	LoadingSpinner  string
	WelcomeMessage  string
	TypingIndicator TypingIndicatorStyle
}

// TypingIndicatorStyle defines how the "thinking" indicator appears
type TypingIndicatorStyle struct {
	Text      string // "Thinking..." or "Processing..."
	Animation string // "dots", "gradient", "pulse", "wave"
	UseGradient bool // Use gradient animation for Gemini
}

// Name returns the theme name (implements Theme interface)
func (t *ProviderTheme) Name() string {
	return t.ProviderName + " Theme"
}

// NewClaudeTheme creates the Claude (Anthropic) provider theme
// Based on Claude Code's warm, friendly aesthetic
func NewClaudeTheme() *ProviderTheme {
	return &ProviderTheme{
		ProviderID:   "claude",
		ProviderName: "Claude",

		BaseTheme: BaseTheme{
			// Primary accents - warm copper/orange
			PrimaryColor:   adaptiveColor("#D4754C", "#D4754C"), // Signature copper orange
			SecondaryColor: adaptiveColor("#B85C3C", "#B85C3C"), // Darker copper
			AccentColor:    adaptiveColor("#F08C5C", "#F08C5C"), // Lighter warm orange

			// Background - warm dark tones
			BackgroundColor:        adaptiveColor("#1A1816", "#F5F2EE"), // Warm dark / warm light
			BackgroundPanelColor:   adaptiveColor("#2C2622", "#FFFFFF"), // Panel bg
			BackgroundElementColor: adaptiveColor("#3A3330", "#F5F2EE"), // Element bg

			// Borders - distinctive orange
			BorderSubtleColor: adaptiveColor("#4A3F38", "#D4BBA8"),
			BorderColor:       adaptiveColor("#D4754C", "#D4754C"), // Signature copper
			BorderActiveColor: adaptiveColor("#F08C5C", "#C96A42"), // Bright on focus

			// Text - warm tones
			TextColor:      adaptiveColor("#E8D5C4", "#2C2622"), // Warm cream / dark
			TextMutedColor: adaptiveColor("#9C8373", "#6B5E52"), // Muted warm gray

			// Status colors
			ErrorColor:   adaptiveColor("#D47C7C", "#B84444"),
			WarningColor: adaptiveColor("#E8A968", "#C98840"),
			SuccessColor: adaptiveColor("#6FA86F", "#4A8A4A"),
			InfoColor:    adaptiveColor("#D4754C", "#D4754C"), // Use primary

			// Diff colors
			DiffAddedColor:            adaptiveColor("#6FA86F", "#E6F4E6"),
			DiffRemovedColor:          adaptiveColor("#D47C7C", "#FFEAEA"),
			DiffContextColor:          adaptiveColor("#9C8373", "#6B5E52"),
			DiffHunkHeaderColor:       adaptiveColor("#D4754C", "#D4754C"),
			DiffHighlightAddedColor:   adaptiveColor("#4A8A4A", "#2D6B2D"),
			DiffHighlightRemovedColor: adaptiveColor("#B84444", "#8B2222"),
			DiffAddedBgColor:          adaptiveColor("#2A3A2A", "#E6F4E6"),
			DiffRemovedBgColor:        adaptiveColor("#3A2A2A", "#FFEAEA"),
			DiffContextBgColor:        adaptiveColor("#1A1816", "#F5F2EE"),
			DiffLineNumberColor:       adaptiveColor("#9C8373", "#6B5E52"),
			DiffAddedLineNumberBgColor:   adaptiveColor("#2A3A2A", "#C8E6C8"),
			DiffRemovedLineNumberBgColor: adaptiveColor("#3A2A2A", "#FFCCCC"),

			// Markdown colors
			MarkdownTextColor:            adaptiveColor("#E8D5C4", "#2C2622"),
			MarkdownHeadingColor:         adaptiveColor("#F08C5C", "#C96A42"),
			MarkdownLinkColor:            adaptiveColor("#D4754C", "#B85C3C"),
			MarkdownLinkTextColor:        adaptiveColor("#F08C5C", "#D4754C"),
			MarkdownCodeColor:            adaptiveColor("#E8A968", "#C98840"),
			MarkdownBlockQuoteColor:      adaptiveColor("#9C8373", "#6B5E52"),
			MarkdownEmphColor:            adaptiveColor("#E8D5C4", "#2C2622"),
			MarkdownStrongColor:          adaptiveColor("#F08C5C", "#C96A42"),
			MarkdownHorizontalRuleColor:  adaptiveColor("#4A3F38", "#D4BBA8"),
			MarkdownListItemColor:        adaptiveColor("#D4754C", "#B85C3C"),
			MarkdownListEnumerationColor: adaptiveColor("#9C8373", "#6B5E52"),
			MarkdownImageColor:           adaptiveColor("#F08C5C", "#D4754C"),
			MarkdownImageTextColor:       adaptiveColor("#9C8373", "#6B5E52"),
			MarkdownCodeBlockColor:       adaptiveColor("#E8A968", "#C98840"),

			// Syntax highlighting - warm tones
			SyntaxCommentColor:     adaptiveColor("#9C8373", "#6B5E52"),
			SyntaxKeywordColor:     adaptiveColor("#F08C5C", "#C96A42"),
			SyntaxFunctionColor:    adaptiveColor("#E8A968", "#C98840"),
			SyntaxVariableColor:    adaptiveColor("#E8D5C4", "#2C2622"),
			SyntaxStringColor:      adaptiveColor("#6FA86F", "#4A8A4A"),
			SyntaxNumberColor:      adaptiveColor("#D47C7C", "#B84444"),
			SyntaxTypeColor:        adaptiveColor("#D4754C", "#B85C3C"),
			SyntaxOperatorColor:    adaptiveColor("#F08C5C", "#D4754C"),
			SyntaxPunctuationColor: adaptiveColor("#9C8373", "#6B5E52"),
		},

		LogoASCII: `
     ▄████▄   ██▓    ▄▄▄       █    ██ ▓█████▄ ▓█████
    ▒██▀ ▀█  ▓██▒   ▒████▄     ██  ▓██▒▒██▀ ██▌▓█   ▀
    ▒▓█    ▄ ▒██░   ▒██  ▀█▄  ▓██  ▒██░░██   █▌▒███
    ▒▓▓▄ ▄██▒▒██░   ░██▄▄▄▄██ ▓▓█  ░██░░▓█▄   ▌▒▓█  ▄
    ▒ ▓███▀ ░░██████▒▓█   ▓██▒▒▒█████▓ ░▒████▓ ░▒████▒
    ░ ░▒ ▒  ░░ ▒░▓  ░▒▒   ▓▒█░░▒▓▒ ▒ ▒  ▒▒▓  ▒ ░░ ▒░ ░
      ░  ▒   ░ ░ ▒  ░ ▒   ▒▒ ░░░▒░ ░ ░  ░ ▒  ▒  ░ ░  ░
    ░          ░ ░    ░   ▒    ░░░ ░ ░  ░ ░  ░    ░
    ░ ░          ░  ░     ░  ░   ░        ░       ░  ░
    ░                                    ░            `,

		WelcomeMessage: "Welcome to Claude! I'm here to help you build amazing things.",
		LoadingSpinner:  "⣾⣽⣻⢿⡿⣟⣯⣷", // Braille spinner

		TypingIndicator: TypingIndicatorStyle{
			Text:      "Thinking",
			Animation: "dots",
		},
	}
}

// NewGeminiTheme creates the Gemini (Google) provider theme
// Based on Gemini CLI's vibrant, modern aesthetic with blue-pink gradient
func NewGeminiTheme() *ProviderTheme {
	return &ProviderTheme{
		ProviderID:   "gemini",
		ProviderName: "Gemini",

		BaseTheme: BaseTheme{
			// Primary accents - Google blue to pink gradient
			PrimaryColor:   adaptiveColor("#4285F4", "#4285F4"), // Google blue
			SecondaryColor: adaptiveColor("#9B72F2", "#9B72F2"), // Purple midpoint
			AccentColor:    adaptiveColor("#EA4335", "#EA4335"), // Google red/pink

			// Background - cool dark
			BackgroundColor:        adaptiveColor("#0D0D0D", "#FFFFFF"),
			BackgroundPanelColor:   adaptiveColor("#1A1A1A", "#F8F9FA"),
			BackgroundElementColor: adaptiveColor("#2A2A2A", "#F1F3F4"),

			// Borders - gradient inspired
			BorderSubtleColor: adaptiveColor("#2A2A45", "#C5CAE9"),
			BorderColor:       adaptiveColor("#4285F4", "#4285F4"), // Blue primary
			BorderActiveColor: adaptiveColor("#9B72F2", "#7B52D2"), // Purple on focus

			// Text - cool tones
			TextColor:      adaptiveColor("#E8EAED", "#202124"), // Light gray / dark
			TextMutedColor: adaptiveColor("#9AA0A6", "#5F6368"), // Medium gray

			// Status colors - Google palette
			ErrorColor:   adaptiveColor("#EA4335", "#C5221F"),
			WarningColor: adaptiveColor("#FBBC04", "#F9AB00"),
			SuccessColor: adaptiveColor("#34A853", "#1E8E3E"),
			InfoColor:    adaptiveColor("#4285F4", "#1A73E8"),

			// Diff colors
			DiffAddedColor:            adaptiveColor("#34A853", "#E6F4EA"),
			DiffRemovedColor:          adaptiveColor("#EA4335", "#FCE8E6"),
			DiffContextColor:          adaptiveColor("#9AA0A6", "#5F6368"),
			DiffHunkHeaderColor:       adaptiveColor("#4285F4", "#4285F4"),
			DiffHighlightAddedColor:   adaptiveColor("#1E8E3E", "#137333"),
			DiffHighlightRemovedColor: adaptiveColor("#C5221F", "#A50E0E"),
			DiffAddedBgColor:          adaptiveColor("#1A2A1A", "#E6F4EA"),
			DiffRemovedBgColor:        adaptiveColor("#2A1A1A", "#FCE8E6"),
			DiffContextBgColor:        adaptiveColor("#0D0D0D", "#FFFFFF"),
			DiffLineNumberColor:       adaptiveColor("#9AA0A6", "#5F6368"),
			DiffAddedLineNumberBgColor:   adaptiveColor("#1A2A1A", "#CEEAD6"),
			DiffRemovedLineNumberBgColor: adaptiveColor("#2A1A1A", "#F9DEDC"),

			// Markdown colors
			MarkdownTextColor:            adaptiveColor("#E8EAED", "#202124"),
			MarkdownHeadingColor:         adaptiveColor("#4285F4", "#1A73E8"),
			MarkdownLinkColor:            adaptiveColor("#4285F4", "#1A73E8"),
			MarkdownLinkTextColor:        adaptiveColor("#9B72F2", "#7B52D2"),
			MarkdownCodeColor:            adaptiveColor("#FBBC04", "#F9AB00"),
			MarkdownBlockQuoteColor:      adaptiveColor("#9AA0A6", "#5F6368"),
			MarkdownEmphColor:            adaptiveColor("#E8EAED", "#202124"),
			MarkdownStrongColor:          adaptiveColor("#4285F4", "#1A73E8"),
			MarkdownHorizontalRuleColor:  adaptiveColor("#2A2A45", "#C5CAE9"),
			MarkdownListItemColor:        adaptiveColor("#4285F4", "#1A73E8"),
			MarkdownListEnumerationColor: adaptiveColor("#9AA0A6", "#5F6368"),
			MarkdownImageColor:           adaptiveColor("#EA4335", "#C5221F"),
			MarkdownImageTextColor:       adaptiveColor("#9AA0A6", "#5F6368"),
			MarkdownCodeBlockColor:       adaptiveColor("#FBBC04", "#F9AB00"),

			// Syntax highlighting - vibrant colors
			SyntaxCommentColor:     adaptiveColor("#9AA0A6", "#5F6368"),
			SyntaxKeywordColor:     adaptiveColor("#EA4335", "#C5221F"),
			SyntaxFunctionColor:    adaptiveColor("#4285F4", "#1A73E8"),
			SyntaxVariableColor:    adaptiveColor("#E8EAED", "#202124"),
			SyntaxStringColor:      adaptiveColor("#34A853", "#1E8E3E"),
			SyntaxNumberColor:      adaptiveColor("#FBBC04", "#F9AB00"),
			SyntaxTypeColor:        adaptiveColor("#9B72F2", "#7B52D2"),
			SyntaxOperatorColor:    adaptiveColor("#EA4335", "#C5221F"),
			SyntaxPunctuationColor: adaptiveColor("#9AA0A6", "#5F6368"),
		},

		LogoASCII: `
     ██████╗ ███████╗███╗   ███╗██╗███╗   ██╗██╗
    ██╔════╝ ██╔════╝████╗ ████║██║████╗  ██║██║
    ██║  ███╗█████╗  ██╔████╔██║██║██╔██╗ ██║██║
    ██║   ██║██╔══╝  ██║╚██╔╝██║██║██║╚██╗██║██║
    ╚██████╔╝███████╗██║ ╚═╝ ██║██║██║ ╚████║██║
     ╚═════╝ ╚══════╝╚═╝     ╚═╝╚═╝╚═╝  ╚═══╝╚═╝`,

		WelcomeMessage: "Welcome to Gemini! Let's explore possibilities together.",
		LoadingSpinner:  "◐◓◑◒", // Circle spinner

		TypingIndicator: TypingIndicatorStyle{
			Text:        "Thinking",
			Animation:   "gradient",
			UseGradient: true, // Enable gradient animation
		},
	}
}

// NewCodexTheme creates the Codex (OpenAI) provider theme
// Based on OpenAI's professional, technical aesthetic
func NewCodexTheme() *ProviderTheme {
	return &ProviderTheme{
		ProviderID:   "codex",
		ProviderName: "Codex",

		BaseTheme: BaseTheme{
			// Primary accents - OpenAI teal
			PrimaryColor:   adaptiveColor("#10A37F", "#10A37F"), // OpenAI teal
			SecondaryColor: adaptiveColor("#0D8569", "#0D8569"), // Darker teal
			AccentColor:    adaptiveColor("#1FC2AA", "#1FC2AA"), // Lighter teal

			// Background - neutral dark
			BackgroundColor:        adaptiveColor("#0E0E0E", "#FFFFFF"),
			BackgroundPanelColor:   adaptiveColor("#1C1C1C", "#F7F7F8"),
			BackgroundElementColor: adaptiveColor("#2D2D2D", "#ECECF1"),

			// Borders - teal accent
			BorderSubtleColor: adaptiveColor("#2D3D38", "#C5E0D8"),
			BorderColor:       adaptiveColor("#10A37F", "#10A37F"), // Teal
			BorderActiveColor: adaptiveColor("#1FC2AA", "#0D8569"), // Bright teal

			// Text - clean neutrals
			TextColor:      adaptiveColor("#ECECEC", "#353740"), // Off-white / dark gray
			TextMutedColor: adaptiveColor("#8E8E8E", "#6E6E80"), // Medium gray

			// Status colors
			ErrorColor:   adaptiveColor("#EF4444", "#DC2626"),
			WarningColor: adaptiveColor("#F59E0B", "#D97706"),
			SuccessColor: adaptiveColor("#10A37F", "#0D8569"), // Use primary
			InfoColor:    adaptiveColor("#3B82F6", "#2563EB"),

			// Diff colors
			DiffAddedColor:            adaptiveColor("#10A37F", "#D1FAE5"),
			DiffRemovedColor:          adaptiveColor("#EF4444", "#FEE2E2"),
			DiffContextColor:          adaptiveColor("#8E8E8E", "#6E6E80"),
			DiffHunkHeaderColor:       adaptiveColor("#10A37F", "#10A37F"),
			DiffHighlightAddedColor:   adaptiveColor("#0D8569", "#059669"),
			DiffHighlightRemovedColor: adaptiveColor("#DC2626", "#B91C1C"),
			DiffAddedBgColor:          adaptiveColor("#1A2D28", "#D1FAE5"),
			DiffRemovedBgColor:        adaptiveColor("#2D1A1A", "#FEE2E2"),
			DiffContextBgColor:        adaptiveColor("#0E0E0E", "#FFFFFF"),
			DiffLineNumberColor:       adaptiveColor("#8E8E8E", "#6E6E80"),
			DiffAddedLineNumberBgColor:   adaptiveColor("#1A2D28", "#A7F3D0"),
			DiffRemovedLineNumberBgColor: adaptiveColor("#2D1A1A", "#FECACA"),

			// Markdown colors
			MarkdownTextColor:            adaptiveColor("#ECECEC", "#353740"),
			MarkdownHeadingColor:         adaptiveColor("#1FC2AA", "#0D8569"),
			MarkdownLinkColor:            adaptiveColor("#10A37F", "#10A37F"),
			MarkdownLinkTextColor:        adaptiveColor("#1FC2AA", "#0D8569"),
			MarkdownCodeColor:            adaptiveColor("#F59E0B", "#D97706"),
			MarkdownBlockQuoteColor:      adaptiveColor("#8E8E8E", "#6E6E80"),
			MarkdownEmphColor:            adaptiveColor("#ECECEC", "#353740"),
			MarkdownStrongColor:          adaptiveColor("#1FC2AA", "#0D8569"),
			MarkdownHorizontalRuleColor:  adaptiveColor("#2D3D38", "#C5E0D8"),
			MarkdownListItemColor:        adaptiveColor("#10A37F", "#10A37F"),
			MarkdownListEnumerationColor: adaptiveColor("#8E8E8E", "#6E6E80"),
			MarkdownImageColor:           adaptiveColor("#3B82F6", "#2563EB"),
			MarkdownImageTextColor:       adaptiveColor("#8E8E8E", "#6E6E80"),
			MarkdownCodeBlockColor:       adaptiveColor("#F59E0B", "#D97706"),

			// Syntax highlighting - technical colors
			SyntaxCommentColor:     adaptiveColor("#8E8E8E", "#6E6E80"),
			SyntaxKeywordColor:     adaptiveColor("#EF4444", "#DC2626"),
			SyntaxFunctionColor:    adaptiveColor("#3B82F6", "#2563EB"),
			SyntaxVariableColor:    adaptiveColor("#ECECEC", "#353740"),
			SyntaxStringColor:      adaptiveColor("#10A37F", "#0D8569"),
			SyntaxNumberColor:      adaptiveColor("#F59E0B", "#D97706"),
			SyntaxTypeColor:        adaptiveColor("#1FC2AA", "#0D8569"),
			SyntaxOperatorColor:    adaptiveColor("#ECECEC", "#353740"),
			SyntaxPunctuationColor: adaptiveColor("#8E8E8E", "#6E6E80"),
		},

		LogoASCII: `
     ██████╗ ██████╗ ██████╗ ███████╗██╗  ██╗
    ██╔════╝██╔═══██╗██╔══██╗██╔════╝╚██╗██╔╝
    ██║     ██║   ██║██║  ██║█████╗   ╚███╔╝
    ██║     ██║   ██║██║  ██║██╔══╝   ██╔██╗
    ╚██████╗╚██████╔╝██████╔╝███████╗██╔╝ ██╗
     ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝`,

		WelcomeMessage: "Welcome to Codex. Let's build something extraordinary.",
		LoadingSpinner:  "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏", // Line spinner

		TypingIndicator: TypingIndicatorStyle{
			Text:      "Processing",
			Animation: "pulse",
		},
	}
}

// NewQwenTheme creates the Qwen (Alibaba) provider theme
// Based on Alibaba's modern, international aesthetic
func NewQwenTheme() *ProviderTheme {
	return &ProviderTheme{
		ProviderID:   "qwen",
		ProviderName: "Qwen",

		BaseTheme: BaseTheme{
			// Primary accents - Alibaba orange
			PrimaryColor:   adaptiveColor("#FF6A00", "#FF6A00"), // Alibaba orange
			SecondaryColor: adaptiveColor("#E55D00", "#E55D00"), // Darker orange
			AccentColor:    adaptiveColor("#FF8533", "#FF8533"), // Lighter orange

			// Background - warm dark
			BackgroundColor:        adaptiveColor("#161410", "#FFFBF5"),
			BackgroundPanelColor:   adaptiveColor("#221E18", "#FFF8ED"),
			BackgroundElementColor: adaptiveColor("#2F2A22", "#FFF0DC"),

			// Borders - orange accent
			BorderSubtleColor: adaptiveColor("#3A352C", "#FFDDB3"),
			BorderColor:       adaptiveColor("#FF6A00", "#FF6A00"), // Orange
			BorderActiveColor: adaptiveColor("#FF8533", "#E55D00"), // Bright orange

			// Text - neutral with warm tint
			TextColor:      adaptiveColor("#F0E8DC", "#2F2A22"), // Warm off-white / dark
			TextMutedColor: adaptiveColor("#A0947C", "#6B5E52"), // Warm gray

			// Status colors - Chinese palette
			ErrorColor:   adaptiveColor("#FF4D4F", "#CF1322"),
			WarningColor: adaptiveColor("#FAAD14", "#D48806"),
			SuccessColor: adaptiveColor("#52C41A", "#389E0D"),
			InfoColor:    adaptiveColor("#1890FF", "#096DD9"),

			// Diff colors
			DiffAddedColor:            adaptiveColor("#52C41A", "#F6FFED"),
			DiffRemovedColor:          adaptiveColor("#FF4D4F", "#FFF1F0"),
			DiffContextColor:          adaptiveColor("#A0947C", "#6B5E52"),
			DiffHunkHeaderColor:       adaptiveColor("#FF6A00", "#FF6A00"),
			DiffHighlightAddedColor:   adaptiveColor("#389E0D", "#237804"),
			DiffHighlightRemovedColor: adaptiveColor("#CF1322", "#A8071A"),
			DiffAddedBgColor:          adaptiveColor("#1A2A1A", "#F6FFED"),
			DiffRemovedBgColor:        adaptiveColor("#2A1A1A", "#FFF1F0"),
			DiffContextBgColor:        adaptiveColor("#161410", "#FFFBF5"),
			DiffLineNumberColor:       adaptiveColor("#A0947C", "#6B5E52"),
			DiffAddedLineNumberBgColor:   adaptiveColor("#1A2A1A", "#D9F7BE"),
			DiffRemovedLineNumberBgColor: adaptiveColor("#2A1A1A", "#FFCCC7"),

			// Markdown colors
			MarkdownTextColor:            adaptiveColor("#F0E8DC", "#2F2A22"),
			MarkdownHeadingColor:         adaptiveColor("#FF6A00", "#E55D00"),
			MarkdownLinkColor:            adaptiveColor("#1890FF", "#096DD9"),
			MarkdownLinkTextColor:        adaptiveColor("#FF8533", "#FF6A00"),
			MarkdownCodeColor:            adaptiveColor("#FAAD14", "#D48806"),
			MarkdownBlockQuoteColor:      adaptiveColor("#A0947C", "#6B5E52"),
			MarkdownEmphColor:            adaptiveColor("#F0E8DC", "#2F2A22"),
			MarkdownStrongColor:          adaptiveColor("#FF6A00", "#E55D00"),
			MarkdownHorizontalRuleColor:  adaptiveColor("#3A352C", "#FFDDB3"),
			MarkdownListItemColor:        adaptiveColor("#FF6A00", "#E55D00"),
			MarkdownListEnumerationColor: adaptiveColor("#A0947C", "#6B5E52"),
			MarkdownImageColor:           adaptiveColor("#1890FF", "#096DD9"),
			MarkdownImageTextColor:       adaptiveColor("#A0947C", "#6B5E52"),
			MarkdownCodeBlockColor:       adaptiveColor("#FAAD14", "#D48806"),

			// Syntax highlighting - international colors
			SyntaxCommentColor:     adaptiveColor("#A0947C", "#6B5E52"),
			SyntaxKeywordColor:     adaptiveColor("#FF4D4F", "#CF1322"),
			SyntaxFunctionColor:    adaptiveColor("#1890FF", "#096DD9"),
			SyntaxVariableColor:    adaptiveColor("#F0E8DC", "#2F2A22"),
			SyntaxStringColor:      adaptiveColor("#52C41A", "#389E0D"),
			SyntaxNumberColor:      adaptiveColor("#FAAD14", "#D48806"),
			SyntaxTypeColor:        adaptiveColor("#FF8533", "#E55D00"),
			SyntaxOperatorColor:    adaptiveColor("#FF6A00", "#E55D00"),
			SyntaxPunctuationColor: adaptiveColor("#A0947C", "#6B5E52"),
		},

		LogoASCII: `
     ██████╗ ██╗    ██╗███████╗███╗   ██╗
    ██╔═══██╗██║    ██║██╔════╝████╗  ██║
    ██║   ██║██║ █╗ ██║█████╗  ██╔██╗ ██║
    ██║▄▄ ██║██║███╗██║██╔══╝  ██║╚██╗██║
    ╚██████╔╝╚███╔███╔╝███████╗██║ ╚████║
     ╚══▀▀═╝  ╚══╝╚══╝ ╚══════╝╚═╝  ╚═══╝`,

		WelcomeMessage: "Welcome to Qwen! Ready to innovate together.",
		LoadingSpinner:  "⣾⣽⣻⢿⡿⣟⣯⣷", // Braille spinner

		TypingIndicator: TypingIndicatorStyle{
			Text:      "Thinking",
			Animation: "wave",
		},
	}
}

// adaptiveColor creates an AdaptiveColor from light/dark hex values
func adaptiveColor(dark, light string) compat.AdaptiveColor {
	return compat.AdaptiveColor{
		Light: lipgloss.Color(light),
		Dark:  lipgloss.Color(dark),
	}
}
