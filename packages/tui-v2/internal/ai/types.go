package ai

import "context"

// Provider defines the interface for AI service providers (Claude, GPT-4, etc.)
type Provider interface {
	// Stream sends a prompt and streams back response tokens via the channel
	// Returns error if request fails
	Stream(ctx context.Context, prompt string, messages []Message) (<-chan StreamEvent, error)

	// Name returns the provider name (e.g., "Claude", "GPT-4")
	Name() string

	// Model returns the model identifier (e.g., "claude-opus-4", "gpt-4")
	Model() string
}

// Message represents a single message in the conversation
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

// Role represents the message sender role
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

// StreamEvent represents a single event in the streaming response
type StreamEvent struct {
	Type         EventType
	Content      string
	Error        error
	Done         bool
	TokensUsed   int // Tokens used in this chunk (if available)
	TotalTokens  int // Total tokens used so far (cumulative)
	PromptTokens int // Tokens in the prompt (set on first event)
}

// EventType represents the type of streaming event
type EventType string

const (
	EventTypeChunk    EventType = "chunk"    // Token/text chunk
	EventTypeComplete EventType = "complete" // Stream completed
	EventTypeError    EventType = "error"    // Error occurred
)

// Config holds configuration for AI providers
type Config struct {
	// Provider selection ("claude", "openai", "auto")
	Provider string

	// API keys
	ClaudeAPIKey string
	OpenAIAPIKey string

	// Model selection
	ClaudeModel string // Default: "claude-opus-4-20250514"
	OpenAIModel string // Default: "gpt-4o"

	// Request parameters
	MaxTokens   int     // Maximum tokens to generate (default: 4096)
	Temperature float64 // Sampling temperature 0-1 (default: 0.7)
	TopP        float64 // Nucleus sampling (default: 0.9)

	// Rate limiting
	RequestsPerMinute int // Max requests per minute (default: 50)
	TokensPerMinute   int // Max tokens per minute (default: 100000)
}

// DefaultConfig returns the default AI configuration
func DefaultConfig() *Config {
	return &Config{
		Provider:          "auto", // Auto-select based on available API keys
		ClaudeModel:       "claude-opus-4-20250514",
		OpenAIModel:       "gpt-4o",
		MaxTokens:         4096,
		Temperature:       0.7,
		TopP:              0.9,
		RequestsPerMinute: 50,
		TokensPerMinute:   100000,
	}
}
