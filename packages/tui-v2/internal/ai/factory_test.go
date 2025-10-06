package ai

import (
	"context"
	"os"
	"testing"
)

func TestLoadConfigFromEnv(t *testing.T) {
	// Save original env vars
	origClaudeKey := os.Getenv("ANTHROPIC_API_KEY")
	origOpenAIKey := os.Getenv("OPENAI_API_KEY")
	origProvider := os.Getenv("RYCODE_AI_PROVIDER")
	origClaudeModel := os.Getenv("RYCODE_CLAUDE_MODEL")
	origOpenAIModel := os.Getenv("RYCODE_OPENAI_MODEL")

	// Restore env vars after test
	defer func() {
		os.Setenv("ANTHROPIC_API_KEY", origClaudeKey)
		os.Setenv("OPENAI_API_KEY", origOpenAIKey)
		os.Setenv("RYCODE_AI_PROVIDER", origProvider)
		os.Setenv("RYCODE_CLAUDE_MODEL", origClaudeModel)
		os.Setenv("RYCODE_OPENAI_MODEL", origOpenAIModel)
	}()

	t.Run("Default config", func(t *testing.T) {
		// Clear all env vars
		os.Unsetenv("ANTHROPIC_API_KEY")
		os.Unsetenv("OPENAI_API_KEY")
		os.Unsetenv("RYCODE_AI_PROVIDER")
		os.Unsetenv("RYCODE_CLAUDE_MODEL")
		os.Unsetenv("RYCODE_OPENAI_MODEL")

		config, err := LoadConfigFromEnv()
		if err != nil {
			t.Fatalf("LoadConfigFromEnv() error = %v", err)
		}

		if config.Provider != "auto" {
			t.Errorf("config.Provider = %v, want auto", config.Provider)
		}
		if config.ClaudeAPIKey != "" {
			t.Errorf("config.ClaudeAPIKey = %v, want empty", config.ClaudeAPIKey)
		}
		if config.OpenAIAPIKey != "" {
			t.Errorf("config.OpenAIAPIKey = %v, want empty", config.OpenAIAPIKey)
		}
	})

	t.Run("Claude API key from env", func(t *testing.T) {
		os.Setenv("ANTHROPIC_API_KEY", "sk-ant-test-key")
		os.Unsetenv("OPENAI_API_KEY")

		config, err := LoadConfigFromEnv()
		if err != nil {
			t.Fatalf("LoadConfigFromEnv() error = %v", err)
		}

		if config.ClaudeAPIKey != "sk-ant-test-key" {
			t.Errorf("config.ClaudeAPIKey = %v, want sk-ant-test-key", config.ClaudeAPIKey)
		}
	})

	t.Run("OpenAI API key from env", func(t *testing.T) {
		os.Unsetenv("ANTHROPIC_API_KEY")
		os.Setenv("OPENAI_API_KEY", "sk-test-key")

		config, err := LoadConfigFromEnv()
		if err != nil {
			t.Fatalf("LoadConfigFromEnv() error = %v", err)
		}

		if config.OpenAIAPIKey != "sk-test-key" {
			t.Errorf("config.OpenAIAPIKey = %v, want sk-test-key", config.OpenAIAPIKey)
		}
	})

	t.Run("Provider override", func(t *testing.T) {
		os.Setenv("RYCODE_AI_PROVIDER", "claude")

		config, err := LoadConfigFromEnv()
		if err != nil {
			t.Fatalf("LoadConfigFromEnv() error = %v", err)
		}

		if config.Provider != "claude" {
			t.Errorf("config.Provider = %v, want claude", config.Provider)
		}
	})

	t.Run("Model overrides", func(t *testing.T) {
		os.Setenv("RYCODE_CLAUDE_MODEL", "claude-sonnet-4")
		os.Setenv("RYCODE_OPENAI_MODEL", "gpt-4-turbo")

		config, err := LoadConfigFromEnv()
		if err != nil {
			t.Fatalf("LoadConfigFromEnv() error = %v", err)
		}

		if config.ClaudeModel != "claude-sonnet-4" {
			t.Errorf("config.ClaudeModel = %v, want claude-sonnet-4", config.ClaudeModel)
		}
		if config.OpenAIModel != "gpt-4-turbo" {
			t.Errorf("config.OpenAIModel = %v, want gpt-4-turbo", config.OpenAIModel)
		}
	})
}

func TestNewProvider(t *testing.T) {
	// Register a mock provider for testing
	RegisterProviders(
		func(apiKey string, config *Config) (Provider, error) {
			return &mockProvider{name: "MockClaude", model: config.ClaudeModel}, nil
		},
		func(apiKey string, config *Config) (Provider, error) {
			return &mockProvider{name: "MockOpenAI", model: config.OpenAIModel}, nil
		},
	)

	t.Run("No API keys", func(t *testing.T) {
		config := &Config{
			Provider: "auto",
		}

		_, err := NewProvider(config)
		if err == nil {
			t.Error("NewProvider() should return error when no API keys")
		}
	})

	t.Run("Auto-select Claude", func(t *testing.T) {
		config := &Config{
			Provider:     "auto",
			ClaudeAPIKey: "sk-ant-test",
			ClaudeModel:  "claude-opus-4-20250514",
		}

		provider, err := NewProvider(config)
		if err != nil {
			t.Fatalf("NewProvider() error = %v", err)
		}

		if provider.Name() != "MockClaude" {
			t.Errorf("provider.Name() = %v, want MockClaude", provider.Name())
		}
	})

	t.Run("Auto-select OpenAI", func(t *testing.T) {
		config := &Config{
			Provider:     "auto",
			OpenAIAPIKey: "sk-test",
			OpenAIModel:  "gpt-4o",
		}

		provider, err := NewProvider(config)
		if err != nil {
			t.Fatalf("NewProvider() error = %v", err)
		}

		if provider.Name() != "MockOpenAI" {
			t.Errorf("provider.Name() = %v, want MockOpenAI", provider.Name())
		}
	})

	t.Run("Force Claude", func(t *testing.T) {
		config := &Config{
			Provider:     "claude",
			ClaudeAPIKey: "sk-ant-test",
			ClaudeModel:  "claude-opus-4-20250514",
		}

		provider, err := NewProvider(config)
		if err != nil {
			t.Fatalf("NewProvider() error = %v", err)
		}

		if provider.Name() != "MockClaude" {
			t.Errorf("provider.Name() = %v, want MockClaude", provider.Name())
		}
	})

	t.Run("Force OpenAI", func(t *testing.T) {
		config := &Config{
			Provider:     "openai",
			OpenAIAPIKey: "sk-test",
			OpenAIModel:  "gpt-4o",
		}

		provider, err := NewProvider(config)
		if err != nil {
			t.Fatalf("NewProvider() error = %v", err)
		}

		if provider.Name() != "MockOpenAI" {
			t.Errorf("provider.Name() = %v, want MockOpenAI", provider.Name())
		}
	})

	t.Run("Unknown provider", func(t *testing.T) {
		config := &Config{
			Provider:     "unknown",
			ClaudeAPIKey: "sk-ant-test",
		}

		_, err := NewProvider(config)
		if err == nil {
			t.Error("NewProvider() should return error for unknown provider")
		}
	})

	t.Run("Missing Claude key", func(t *testing.T) {
		config := &Config{
			Provider: "claude",
		}

		_, err := NewProvider(config)
		if err == nil {
			t.Error("NewProvider() should return error when Claude key missing")
		}
	})

	t.Run("Missing OpenAI key", func(t *testing.T) {
		config := &Config{
			Provider: "openai",
		}

		_, err := NewProvider(config)
		if err == nil {
			t.Error("NewProvider() should return error when OpenAI key missing")
		}
	})
}

// mockProvider is a test implementation of Provider
type mockProvider struct {
	name  string
	model string
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) Model() string {
	return m.model
}

func (m *mockProvider) Stream(ctx context.Context, prompt string, messages []Message) (<-chan StreamEvent, error) {
	ch := make(chan StreamEvent)
	go func() {
		defer close(ch)
		ch <- StreamEvent{Type: EventTypeChunk, Content: "test"}
		ch <- StreamEvent{Type: EventTypeComplete, Done: true}
	}()
	return ch, nil
}
