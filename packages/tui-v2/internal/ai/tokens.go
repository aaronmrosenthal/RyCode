package ai

import "strings"

// EstimateTokens provides a rough estimate of token count for text
// This is a simple approximation: ~4 characters per token on average
// Real token counting would require provider-specific tokenizers
func EstimateTokens(text string) int {
	if text == "" {
		return 0
	}

	// Rough approximation: 1 token ≈ 4 characters
	// This matches OpenAI's rule of thumb
	chars := len(text)
	tokens := chars / 4

	// Adjust for whitespace and punctuation (slightly more efficient)
	words := len(strings.Fields(text))
	if words > 0 {
		// Average English word is ~5 chars, ~1.3 tokens
		tokens = int(float64(words) * 1.3)
	}

	// Minimum 1 token for non-empty text
	if tokens == 0 && text != "" {
		tokens = 1
	}

	return tokens
}

// EstimateConversationTokens estimates total tokens for a conversation
func EstimateConversationTokens(messages []Message) int {
	total := 0

	// Add overhead for conversation structure
	// Each message has ~4 tokens of overhead (role markers, etc.)
	overhead := len(messages) * 4

	for _, msg := range messages {
		total += EstimateTokens(msg.Content)
	}

	return total + overhead
}
