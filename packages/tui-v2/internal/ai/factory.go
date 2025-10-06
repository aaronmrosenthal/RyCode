package ai

import (
	"fmt"
	"os"
	"strings"
)

// LoadConfigFromEnv loads AI configuration from environment variables
func LoadConfigFromEnv() (*Config, error) {
	config := DefaultConfig()

	// Load API keys from environment
	config.ClaudeAPIKey = os.Getenv("ANTHROPIC_API_KEY")
	config.OpenAIAPIKey = os.Getenv("OPENAI_API_KEY")

	// Allow provider override
	if provider := os.Getenv("RYCODE_AI_PROVIDER"); provider != "" {
		config.Provider = strings.ToLower(provider)
	}

	// Allow model overrides
	if model := os.Getenv("RYCODE_CLAUDE_MODEL"); model != "" {
		config.ClaudeModel = model
	}
	if model := os.Getenv("RYCODE_OPENAI_MODEL"); model != "" {
		config.OpenAIModel = model
	}

	return config, nil
}

// NewProvider creates an AI provider based on configuration
func NewProvider(config *Config) (Provider, error) {
	if config == nil {
		var err error
		config, err = LoadConfigFromEnv()
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}

	// Determine which provider to use
	providerName := strings.ToLower(config.Provider)

	// Auto-select based on available API keys
	if providerName == "auto" {
		if config.ClaudeAPIKey != "" {
			providerName = "claude"
		} else if config.OpenAIAPIKey != "" {
			providerName = "openai"
		} else {
			return nil, fmt.Errorf("no API keys found; set ANTHROPIC_API_KEY or OPENAI_API_KEY")
		}
	}

	// Create the appropriate provider
	switch providerName {
	case "claude", "anthropic":
		if config.ClaudeAPIKey == "" {
			return nil, fmt.Errorf("Claude API key not found; set ANTHROPIC_API_KEY")
		}
		return newClaudeProvider(config.ClaudeAPIKey, config)

	case "openai", "gpt", "gpt-4":
		if config.OpenAIAPIKey == "" {
			return nil, fmt.Errorf("OpenAI API key not found; set OPENAI_API_KEY")
		}
		return newOpenAIProvider(config.OpenAIAPIKey, config)

	default:
		return nil, fmt.Errorf("unknown provider: %s (supported: claude, openai)", providerName)
	}
}

// MockProvider interface allows injecting provider implementations
// This is used by the factory to avoid import cycles
var (
	newClaudeProvider func(apiKey string, config *Config) (Provider, error)
	newOpenAIProvider func(apiKey string, config *Config) (Provider, error)
)

// RegisterProviders sets up the provider constructors (called from providers package)
func RegisterProviders(
	claudeConstructor func(apiKey string, config *Config) (Provider, error),
	openAIConstructor func(apiKey string, config *Config) (Provider, error),
) {
	newClaudeProvider = claudeConstructor
	newOpenAIProvider = openAIConstructor
}
