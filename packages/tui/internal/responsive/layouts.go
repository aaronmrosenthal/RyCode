package responsive

import (
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

// Message represents a chat message for rendering
type Message struct {
	ID        string
	Role      string
	Content   string
	Timestamp string
	AI        string // "claude", "codex", "gemini"
	Tools     []string
	Reaction  string
}

// PhoneLayout renders mobile-optimized layout
type PhoneLayout struct {
	theme  *theme.Theme
	config LayoutConfig
}

// NewPhoneLayout creates a phone-optimized layout
func NewPhoneLayout(theme *theme.Theme, config LayoutConfig) *PhoneLayout {
	return &PhoneLayout{
		theme:  theme,
		config: config,
	}
}

// RenderMessage renders a single message in phone-optimized bubble style
func (pl *PhoneLayout) RenderMessage(msg Message, isActive bool) string {
	// Chat bubble style for phone
	bubbleStyle := lipgloss.NewStyle().
		Padding(1, 2).
		MarginBottom(1).
		MaxWidth(pl.config.Width - 4)

	// User messages: right-aligned, accent color
	if msg.Role == "user" {
		bubbleStyle = bubbleStyle.
			Align(lipgloss.Right).
			Background(pl.theme.AccentPrimary).
			Foreground(pl.theme.BackgroundPrimary).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(pl.theme.AccentPrimary)

		content := bubbleStyle.Render(msg.Content)

		// Add timestamp below (right-aligned)
		timeStyle := lipgloss.NewStyle().
			Foreground(pl.theme.TextDim).
			Align(lipgloss.Right).
			Width(pl.config.Width - 4)

		timestamp := timeStyle.Render(msg.Timestamp)

		return lipgloss.JoinVertical(lipgloss.Right, content, timestamp)
	}

	// AI messages: left-aligned, secondary background
	bubbleStyle = bubbleStyle.
		Align(lipgloss.Left).
		Background(pl.theme.BackgroundSecondary).
		Foreground(pl.theme.TextPrimary).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(pl.theme.Border)

	if isActive {
		bubbleStyle = bubbleStyle.
			BorderForeground(pl.theme.AccentPrimary).
			BorderStyle(lipgloss.ThickBorder())
	}

	// AI indicator with icon
	aiIcon := getAIIcon(msg.AI)
	aiLabel := lipgloss.NewStyle().
		Foreground(pl.theme.AccentSecondary).
		Bold(true).
		Render(aiIcon + " " + msg.AI)

	// Content
	content := msg.Content
	if len(content) > 200 && pl.config.Width < 50 {
		// Truncate for very small screens
		content = content[:197] + "..."
	}

	// Combine AI label and content
	messageContent := lipgloss.JoinVertical(
		lipgloss.Left,
		aiLabel,
		"",
		content,
	)

	bubble := bubbleStyle.Render(messageContent)

	// Reaction emoji (if present)
	if msg.Reaction != "" {
		reactionStyle := lipgloss.NewStyle().
			Foreground(pl.theme.AccentPrimary).
			MarginLeft(2)

		bubble = lipgloss.JoinHorizontal(
			lipgloss.Left,
			bubble,
			reactionStyle.Render(msg.Reaction),
		)
	}

	// Timestamp
	timeStyle := lipgloss.NewStyle().
		Foreground(pl.theme.TextDim).
		MarginLeft(2)

	timestamp := timeStyle.Render(msg.Timestamp)

	return lipgloss.JoinVertical(lipgloss.Left, bubble, timestamp)
}

// RenderInput renders phone-optimized input
func (pl *PhoneLayout) RenderInput(value string, placeholder string) string {
	inputStyle := lipgloss.NewStyle().
		Width(pl.config.Width - 4).
		Padding(1, 2).
		Background(pl.theme.BackgroundSecondary).
		Foreground(pl.theme.TextPrimary).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(pl.theme.AccentPrimary)

	// Show placeholder if empty
	displayValue := value
	if displayValue == "" {
		displayValue = lipgloss.NewStyle().
			Foreground(pl.theme.TextDim).
			Render(placeholder)
	}

	// Voice button (phone-specific)
	voiceButton := lipgloss.NewStyle().
		Foreground(pl.theme.AccentSecondary).
		Background(pl.theme.BackgroundSecondary).
		Padding(0, 1).
		Bold(true).
		Render("🎤")

	input := inputStyle.Render(displayValue)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		input,
		voiceButton,
	)
}

// RenderQuickActions renders phone-optimized quick action bar
func (pl *PhoneLayout) RenderQuickActions() string {
	buttonStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Background(pl.theme.BackgroundSecondary).
		Foreground(pl.theme.TextPrimary).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(pl.theme.Border)

	actions := []string{
		buttonStyle.Copy().Background(pl.theme.AccentPrimary).Render("💬 Chat"),
		buttonStyle.Render("📜 History"),
		buttonStyle.Render("⚙️ Settings"),
		buttonStyle.Render("🤖 AI"),
	}

	return lipgloss.NewStyle().
		Width(pl.config.Width).
		Align(lipgloss.Center).
		Render(lipgloss.JoinHorizontal(lipgloss.Center, actions...))
}

// TabletLayout renders tablet-optimized layout
type TabletLayout struct {
	theme  *theme.Theme
	config LayoutConfig
}

// NewTabletLayout creates a tablet-optimized layout
func NewTabletLayout(theme *theme.Theme, config LayoutConfig) *TabletLayout {
	return &TabletLayout{
		theme:  theme,
		config: config,
	}
}

// RenderSplitView renders tablet split view (chat + preview)
func (tl *TabletLayout) RenderSplitView(messages []Message, preview string) string {
	// Split into two columns
	leftWidth := tl.config.Width * 50 / 100
	rightWidth := tl.config.Width - leftWidth - 2

	// Left: Messages
	leftStyle := lipgloss.NewStyle().
		Width(leftWidth).
		Height(tl.config.Height).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(tl.theme.Border).
		BorderRight(true).
		Padding(1)

	messageViews := []string{}
	for _, msg := range messages {
		messageViews = append(messageViews, tl.renderCompactMessage(msg))
	}

	leftPane := leftStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, messageViews...),
	)

	// Right: Preview/Context
	rightStyle := lipgloss.NewStyle().
		Width(rightWidth).
		Height(tl.config.Height).
		Padding(1)

	rightPane := rightStyle.Render(preview)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
}

// renderCompactMessage renders a compact message for tablet
func (tl *TabletLayout) renderCompactMessage(msg Message) string {
	// More compact than phone, less verbose than desktop
	roleStyle := lipgloss.NewStyle().
		Foreground(tl.theme.AccentSecondary).
		Bold(true)

	contentStyle := lipgloss.NewStyle().
		Foreground(tl.theme.TextPrimary).
		Width(tl.config.Width*50/100 - 4)

	role := roleStyle.Render(getAIIcon(msg.AI) + " ")
	content := contentStyle.Render(msg.Content)

	return lipgloss.NewStyle().
		MarginBottom(1).
		Padding(1).
		Background(tl.theme.BackgroundSecondary).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(tl.theme.Border).
		Render(lipgloss.JoinVertical(lipgloss.Left, role, content))
}

// DesktopLayout renders traditional desktop layout
type DesktopLayout struct {
	theme  *theme.Theme
	config LayoutConfig
}

// NewDesktopLayout creates a desktop layout
func NewDesktopLayout(theme *theme.Theme, config LayoutConfig) *DesktopLayout {
	return &DesktopLayout{
		theme:  theme,
		config: config,
	}
}

// RenderThreeColumn renders desktop three-column layout
func (dl *DesktopLayout) RenderThreeColumn(sidebar string, messages string, context string) string {
	sidebarWidth := 30
	contextWidth := 40
	messagesWidth := dl.config.Width - sidebarWidth - contextWidth - 4

	sidebarStyle := lipgloss.NewStyle().
		Width(sidebarWidth).
		Height(dl.config.Height).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(dl.theme.Border).
		BorderRight(true).
		Padding(1)

	messagesStyle := lipgloss.NewStyle().
		Width(messagesWidth).
		Height(dl.config.Height).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(dl.theme.Border).
		BorderRight(true).
		Padding(1)

	contextStyle := lipgloss.NewStyle().
		Width(contextWidth).
		Height(dl.config.Height).
		Padding(1)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebarStyle.Render(sidebar),
		messagesStyle.Render(messages),
		contextStyle.Render(context),
	)
}

// AI Icons for different providers
func getAIIcon(ai string) string {
	switch strings.ToLower(ai) {
	case "claude":
		return "🧠"
	case "codex", "openai":
		return "⚡"
	case "gemini":
		return "💎"
	default:
		return "🤖"
	}
}

// ThumbZoneIndicator shows the thumb-reachable zone on phones
func ThumbZoneIndicator(config LayoutConfig, theme *theme.Theme) string {
	if config.Device != DevicePhone {
		return ""
	}

	// Show visual indicator of thumb-friendly zone
	style := lipgloss.NewStyle().
		Foreground(theme.Success).
		Faint(true)

	if config.InputPosition == InputTop {
		return style.Render("👍 Thumb zone (top)")
	}

	return style.Render("👍 Thumb zone (bottom)")
}

// SwipeIndicator shows swipe gesture hints
func SwipeIndicator(direction GestureType, theme *theme.Theme) string {
	style := lipgloss.NewStyle().
		Foreground(theme.AccentPrimary).
		Bold(true)

	switch direction {
	case GestureSwipeLeft:
		return style.Render("← Swipe")
	case GestureSwipeRight:
		return style.Render("Swipe →")
	case GestureSwipeUp:
		return style.Render("↑ Swipe")
	case GestureSwipeDown:
		return style.Render("↓ Swipe")
	default:
		return ""
	}
}

// VoiceInputButton renders voice input button for phone
func VoiceInputButton(active bool, theme *theme.Theme) string {
	style := lipgloss.NewStyle().
		Padding(1, 3).
		Background(theme.AccentPrimary).
		Foreground(theme.BackgroundPrimary).
		BorderStyle(lipgloss.RoundedBorder()).
		Bold(true)

	if active {
		// Pulsing effect when recording
		style = style.
			Background(theme.Error).
			Render("🎤 Recording...")
	}

	return style.Render("🎤 Voice")
}

// AIProviderPicker renders AI provider picker
func AIProviderPicker(current string, theme *theme.Theme, width int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.AccentPrimary).
		Bold(true).
		Width(width).
		Align(lipgloss.Center).
		MarginBottom(1)

	title := titleStyle.Render("🤖 Choose AI")

	providers := []struct {
		name string
		icon string
		desc string
	}{
		{"claude", "🧠", "Claude (Anthropic) - Best for coding"},
		{"codex", "⚡", "Codex (OpenAI) - Fast & efficient"},
		{"gemini", "💎", "Gemini (Google) - Multimodal"},
	}

	items := []string{}
	for i, p := range providers {
		buttonStyle := lipgloss.NewStyle().
			Width(width - 4).
			Padding(1, 2).
			MarginBottom(1).
			Background(theme.BackgroundSecondary).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border)

		if strings.ToLower(current) == p.name {
			buttonStyle = buttonStyle.
				Background(theme.AccentPrimary).
				Foreground(theme.BackgroundPrimary).
				BorderForeground(theme.AccentPrimary)
		}

		label := lipgloss.NewStyle().Bold(true).Render(p.icon + " " + p.name)
		desc := lipgloss.NewStyle().
			Foreground(theme.TextDim).
			Render(p.desc)

		numberStyle := lipgloss.NewStyle().
			Foreground(theme.TextDim).
			Render(string(rune('1' + i)))

		item := lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Left, numberStyle, " ", label),
			desc,
		)

		items = append(items, buttonStyle.Render(item))
	}

	hint := lipgloss.NewStyle().
		Foreground(theme.TextDim).
		Width(width).
		Align(lipgloss.Center).
		MarginTop(1).
		Render("Press 1-3 to switch • ESC to cancel")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		lipgloss.JoinVertical(lipgloss.Left, items...),
		hint,
	)
}
