package providers

import (
	"context"
	"testing"

	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/ai"
)

func TestOpenAIProvider_Name(t *testing.T) {
	config := ai.DefaultConfig()
	provider, err := NewOpenAIProvider("test-key", config)
	if err != nil {
		t.Fatalf("NewOpenAIProvider() error = %v", err)
	}

	if provider.Name() != "OpenAI" {
		t.Errorf("Name() = %v, want OpenAI", provider.Name())
	}
}

func TestOpenAIProvider_Model(t *testing.T) {
	t.Run("Default model", func(t *testing.T) {
		config := ai.DefaultConfig()
		provider, err := NewOpenAIProvider("test-key", config)
		if err != nil {
			t.Fatalf("NewOpenAIProvider() error = %v", err)
		}

		if provider.Model() != "gpt-4o" {
			t.Errorf("Model() = %v, want gpt-4o", provider.Model())
		}
	})

	t.Run("Custom model", func(t *testing.T) {
		config := ai.DefaultConfig()
		config.OpenAIModel = "gpt-4-turbo"
		provider, err := NewOpenAIProvider("test-key", config)
		if err != nil {
			t.Fatalf("NewOpenAIProvider() error = %v", err)
		}

		if provider.Model() != "gpt-4-turbo" {
			t.Errorf("Model() = %v, want gpt-4-turbo", provider.Model())
		}
	})

	t.Run("Empty model uses default", func(t *testing.T) {
		config := ai.DefaultConfig()
		config.OpenAIModel = ""
		provider, err := NewOpenAIProvider("test-key", config)
		if err != nil {
			t.Fatalf("NewOpenAIProvider() error = %v", err)
		}

		if provider.Model() != "gpt-4o" {
			t.Errorf("Model() = %v, want gpt-4o (default)", provider.Model())
		}
	})
}

func TestOpenAIProvider_NilConfig(t *testing.T) {
	provider, err := NewOpenAIProvider("test-key", nil)
	if err != nil {
		t.Fatalf("NewOpenAIProvider() error = %v", err)
	}

	if provider == nil {
		t.Fatal("NewOpenAIProvider() returned nil with nil config")
	}

	// Should use default config values
	if provider.model != "gpt-4o" {
		t.Errorf("model = %v, want gpt-4o", provider.model)
	}
	if provider.maxTokens != 4096 {
		t.Errorf("maxTokens = %v, want 4096", provider.maxTokens)
	}
	if provider.temperature != 0.7 {
		t.Errorf("temperature = %v, want 0.7", provider.temperature)
	}
	if provider.topP != 0.9 {
		t.Errorf("topP = %v, want 0.9", provider.topP)
	}
}

func TestOpenAIProvider_Configuration(t *testing.T) {
	config := &ai.Config{
		OpenAIModel: "gpt-4-turbo",
		MaxTokens:   8192,
		Temperature: 0.5,
		TopP:        0.95,
	}

	provider, err := NewOpenAIProvider("test-key", config)
	if err != nil {
		t.Fatalf("NewOpenAIProvider() error = %v", err)
	}

	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"model", provider.model, "gpt-4-turbo"},
		{"maxTokens", provider.maxTokens, 8192},
		{"temperature", provider.temperature, 0.5},
		{"topP", provider.topP, 0.95},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}

	// Verify apiKey is not empty (can't check value since it's encrypted)
	if provider.apiKey == nil || provider.apiKey.IsEmpty() {
		t.Error("apiKey should not be empty")
	}
}

func TestOpenAIProvider_Stream_Context(t *testing.T) {
	config := ai.DefaultConfig()
	provider, err := NewOpenAIProvider("test-key", config)
	if err != nil {
		t.Fatalf("NewOpenAIProvider() error = %v", err)
	}

	// Test with canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	messages := []ai.Message{
		{Role: ai.RoleUser, Content: "Hello"},
	}

	// Note: This will fail to connect since we canceled context
	// In a real test with dependency injection, we'd verify the context is passed through
	_, streamErr := provider.Stream(ctx, "Test", messages)

	// We expect an error due to canceled context
	if streamErr != nil {
		// This is expected - context was canceled
		t.Logf("Expected error with canceled context: %v", streamErr)
	}
}

func TestOpenAIProvider_HTTPClient(t *testing.T) {
	config := ai.DefaultConfig()
	provider, err := NewOpenAIProvider("test-key", config)
	if err != nil {
		t.Fatalf("NewOpenAIProvider() error = %v", err)
	}

	if provider.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
}

func TestOpenAIProvider_EmptyAPIKey(t *testing.T) {
	config := ai.DefaultConfig()
	_, err := NewOpenAIProvider("", config)

	// Should fail to create SecureString from empty key
	if err == nil {
		t.Error("Expected error when creating provider with empty API key")
	}
}
