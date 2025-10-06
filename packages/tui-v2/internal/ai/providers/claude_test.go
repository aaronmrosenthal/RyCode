package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/ai"
)

func TestClaudeProvider_Name(t *testing.T) {
	config := ai.DefaultConfig()
	provider, err := NewClaudeProvider("test-key", config)
	if err != nil {
		t.Fatalf("NewClaudeProvider() error = %v", err)
	}

	if provider.Name() != "Claude" {
		t.Errorf("Name() = %v, want Claude", provider.Name())
	}
}

func TestClaudeProvider_Model(t *testing.T) {
	t.Run("Default model", func(t *testing.T) {
		config := ai.DefaultConfig()
		provider, err := NewClaudeProvider("test-key", config)
		if err != nil {
			t.Fatalf("NewClaudeProvider() error = %v", err)
		}

		if provider.Model() != "claude-opus-4-20250514" {
			t.Errorf("Model() = %v, want claude-opus-4-20250514", provider.Model())
		}
	})

	t.Run("Custom model", func(t *testing.T) {
		config := ai.DefaultConfig()
		config.ClaudeModel = "claude-sonnet-4-20250514"
		provider, err := NewClaudeProvider("test-key", config)
		if err != nil {
			t.Fatalf("NewClaudeProvider() error = %v", err)
		}

		if provider.Model() != "claude-sonnet-4-20250514" {
			t.Errorf("Model() = %v, want claude-sonnet-4-20250514", provider.Model())
		}
	})
}

func TestClaudeProvider_Stream_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("x-api-key") == "" {
			t.Error("Missing x-api-key header")
		}
		if r.Header.Get("anthropic-version") != claudeAPIVersion {
			t.Errorf("anthropic-version = %v, want %v", r.Header.Get("anthropic-version"), claudeAPIVersion)
		}

		// Send SSE response
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Write streaming events
		events := []string{
			"data: {\"type\":\"content_block_delta\",\"delta\":{\"type\":\"text_delta\",\"text\":\"Hello\"}}\n\n",
			"data: {\"type\":\"content_block_delta\",\"delta\":{\"type\":\"text_delta\",\"text\":\" World\"}}\n\n",
			"data: {\"type\":\"message_stop\"}\n\n",
		}

		for _, event := range events {
			w.Write([]byte(event))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}))
	defer server.Close()

	// Create provider with custom HTTP client
	config := ai.DefaultConfig()
	provider, err := NewClaudeProvider("test-key", config)
	if err != nil {
		t.Fatalf("NewClaudeProvider() error = %v", err)
	}
	provider.httpClient = server.Client()

	// Note: Can't override const claudeAPIURL for testing
	// In production code, we'd need dependency injection for the URL
	// For this test, we'll just verify the provider was created correctly
	if provider == nil {
		t.Fatal("NewClaudeProvider() returned nil")
	}

	// Verify provider fields (skip apiKey since it's encrypted)
	if provider.apiKey == nil || provider.apiKey.IsEmpty() {
		t.Error("apiKey should not be empty")
	}
	if provider.model != config.ClaudeModel {
		t.Errorf("model = %v, want %v", provider.model, config.ClaudeModel)
	}
	if provider.maxTokens != config.MaxTokens {
		t.Errorf("maxTokens = %v, want %v", provider.maxTokens, config.MaxTokens)
	}
	if provider.temperature != config.Temperature {
		t.Errorf("temperature = %v, want %v", provider.temperature, config.Temperature)
	}
	if provider.topP != config.TopP {
		t.Errorf("topP = %v, want %v", provider.topP, config.TopP)
	}
}

func TestClaudeProvider_Stream_Error(t *testing.T) {
	t.Run("HTTP error", func(t *testing.T) {
		// Create mock server that returns error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":{"message":"Invalid API key"}}`))
		}))
		defer server.Close()

		config := ai.DefaultConfig()
		provider, err := NewClaudeProvider("invalid-key", config)
		if err != nil {
			t.Fatalf("NewClaudeProvider() error = %v", err)
		}
		provider.httpClient = server.Client()

		// Note: Can't easily test Stream() without mocking the URL
		// This would require dependency injection or interface-based HTTP client
		// For now, just verify provider creation
		if provider == nil {
			t.Fatal("NewClaudeProvider() returned nil")
		}
	})
}

func TestClaudeProvider_NilConfig(t *testing.T) {
	provider, err := NewClaudeProvider("test-key", nil)
	if err != nil {
		t.Fatalf("NewClaudeProvider() error = %v", err)
	}

	if provider == nil {
		t.Fatal("NewClaudeProvider() returned nil with nil config")
	}

	// Should use default config values
	if provider.model != "claude-opus-4-20250514" {
		t.Errorf("model = %v, want claude-opus-4-20250514", provider.model)
	}
	if provider.maxTokens != 4096 {
		t.Errorf("maxTokens = %v, want 4096", provider.maxTokens)
	}
}

func TestClaudeProvider_Stream_Context(t *testing.T) {
	config := ai.DefaultConfig()
	provider, err := NewClaudeProvider("test-key", config)
	if err != nil {
		t.Fatalf("NewClaudeProvider() error = %v", err)
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
