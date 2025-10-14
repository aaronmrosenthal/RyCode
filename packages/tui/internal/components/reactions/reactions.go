package reactions

import (
	"time"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/aaronmrosenthal/rycode/internal/theme"
)

// Reaction represents an emoji reaction to an AI message
type Reaction string

const (
	ReactionThumbsUp   Reaction = "👍" // Great response
	ReactionThumbsDown Reaction = "👎" // Not helpful
	ReactionThinking   Reaction = "🤔" // Confusing
	ReactionBulb       Reaction = "💡" // Aha moment
	ReactionRocket     Reaction = "🚀" // Excellent
	ReactionBug        Reaction = "🐛" // Found a bug
	ReactionParty      Reaction = "🎉" // Celebration
)

// MessageReaction stores a reaction for a message
type MessageReaction struct {
	MessageID string
	Reaction  Reaction
	Timestamp time.Time
}

// ReactionManager handles reactions and learning
type ReactionManager struct {
	reactions map[string]MessageReaction // messageID -> reaction
	stats     map[Reaction]int            // reaction type -> count
}

// NewReactionManager creates a new reaction manager
func NewReactionManager() *ReactionManager {
	return &ReactionManager{
		reactions: make(map[string]MessageReaction),
		stats:     make(map[Reaction]int),
	}
}

// Add adds a reaction to a message
func (rm *ReactionManager) Add(messageID string, reaction Reaction) {
	mr := MessageReaction{
		MessageID: messageID,
		Reaction:  reaction,
		Timestamp: time.Now(),
	}

	rm.reactions[messageID] = mr
	rm.stats[reaction]++
}

// Get retrieves a reaction for a message
func (rm *ReactionManager) Get(messageID string) *MessageReaction {
	if mr, exists := rm.reactions[messageID]; exists {
		return &mr
	}
	return nil
}

// Remove removes a reaction from a message
func (rm *ReactionManager) Remove(messageID string) {
	if mr, exists := rm.reactions[messageID]; exists {
		rm.stats[mr.Reaction]--
		delete(rm.reactions, messageID)
	}
}

// GetStats returns reaction statistics
func (rm *ReactionManager) GetStats() map[Reaction]int {
	return rm.stats
}

// Render renders a reaction with styling
func Render(reaction Reaction, theme *theme.Theme) string {
	style := lipgloss.NewStyle().
		Foreground((*theme).Primary()).
		Bold(true)

	return style.Render(string(reaction))
}

// RenderPicker renders the reaction picker UI
func RenderPicker(theme *theme.Theme) string {
	reactions := []Reaction{
		ReactionThumbsUp,
		ReactionThumbsDown,
		ReactionThinking,
		ReactionBulb,
		ReactionRocket,
		ReactionBug,
		ReactionParty,
	}

	descriptions := map[Reaction]string{
		ReactionThumbsUp:   "Great response!",
		ReactionThumbsDown: "Not helpful",
		ReactionThinking:   "Confusing",
		ReactionBulb:       "Aha moment!",
		ReactionRocket:     "Excellent!",
		ReactionBug:        "Found a bug",
		ReactionParty:      "Celebration",
	}

	// Build picker
	items := []string{}
	for i, reaction := range reactions {
		desc := descriptions[reaction]
		style := lipgloss.NewStyle().
			Foreground((*theme).Text())

		item := lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().
				Foreground((*theme).Primary()).
				Bold(true).
				Render(string(reaction)),
			style.Render(" "+desc),
		)

		// Add number for keyboard selection
		numberStyle := lipgloss.NewStyle().
			Foreground((*theme).TextMuted()).
			PaddingRight(1)

		items = append(items,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				numberStyle.Render(string(rune('1'+i))),
				item,
			),
		)
	}

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground((*theme).Primary()).
		Bold(true).
		MarginBottom(1)

	title := titleStyle.Render("React to this message")

	// Instruction
	instructionStyle := lipgloss.NewStyle().
		Foreground((*theme).TextMuted()).
		MarginTop(1)

	instruction := instructionStyle.Render("Press 1-7 to react, ESC to cancel")

	// Combine
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		lipgloss.JoinVertical(lipgloss.Left, items...),
		instruction,
	)
}

// GetLearningFeedback converts reactions to AI learning signals
func GetLearningFeedback(reaction Reaction) map[string]interface{} {
	feedback := make(map[string]interface{})

	switch reaction {
	case ReactionThumbsUp, ReactionBulb, ReactionRocket:
		feedback["quality"] = "high"
		feedback["helpful"] = true
		feedback["clarity"] = "good"

	case ReactionThumbsDown:
		feedback["quality"] = "low"
		feedback["helpful"] = false
		feedback["needs_improvement"] = true

	case ReactionThinking:
		feedback["clarity"] = "poor"
		feedback["needs_elaboration"] = true
		feedback["confusing"] = true

	case ReactionBug:
		feedback["contains_error"] = true
		feedback["needs_correction"] = true

	case ReactionParty:
		feedback["celebration"] = true
		feedback["milestone"] = true
		feedback["excellent"] = true
	}

	return feedback
}

// GetSuggestionFromReaction provides AI improvement suggestions
func GetSuggestionFromReaction(reaction Reaction) string {
	suggestions := map[Reaction]string{
		ReactionThumbsDown: "Try rephrasing your question or ask for more details",
		ReactionThinking:   "Ask AI to explain in simpler terms or provide examples",
		ReactionBug:        "Report the error and ask for a correction",
		ReactionThumbsUp:   "Great! You can ask AI to save this explanation for later",
		ReactionBulb:       "Nice! Consider documenting this insight",
		ReactionRocket:     "Excellent! This pattern worked well - AI will remember it",
		ReactionParty:      "🎉 Celebrate and commit your changes!",
	}

	if suggestion, exists := suggestions[reaction]; exists {
		return suggestion
	}

	return ""
}
