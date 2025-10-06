package providers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/ai"
)

const (
	claudeAPIURL     = "https://api.anthropic.com/v1/messages"
	claudeAPIVersion = "2023-06-01"
)

// ClaudeProvider implements the AI Provider interface for Anthropic's Claude API
type ClaudeProvider struct {
	apiKey      string
	model       string
	maxTokens   int
	temperature float64
	topP        float64
	httpClient  *http.Client
}

// NewClaudeProvider creates a new Claude AI provider
func NewClaudeProvider(apiKey string, config *ai.Config) *ClaudeProvider {
	if config == nil {
		config = ai.DefaultConfig()
	}

	model := config.ClaudeModel
	if model == "" {
		model = "claude-opus-4-20250514"
	}

	return &ClaudeProvider{
		apiKey:      apiKey,
		model:       model,
		maxTokens:   config.MaxTokens,
		temperature: config.Temperature,
		topP:        config.TopP,
		httpClient:  &http.Client{},
	}
}

// Name returns the provider name
func (c *ClaudeProvider) Name() string {
	return "Claude"
}

// Model returns the model identifier
func (c *ClaudeProvider) Model() string {
	return c.model
}

// Stream sends a prompt and streams back response tokens
func (c *ClaudeProvider) Stream(ctx context.Context, prompt string, messages []ai.Message) (<-chan ai.StreamEvent, error) {
	eventCh := make(chan ai.StreamEvent, 100)

	// Build request payload
	reqMessages := make([]map[string]string, 0, len(messages)+1)

	// Add conversation history
	for _, msg := range messages {
		reqMessages = append(reqMessages, map[string]string{
			"role":    string(msg.Role),
			"content": msg.Content,
		})
	}

	// Add current prompt
	reqMessages = append(reqMessages, map[string]string{
		"role":    "user",
		"content": prompt,
	})

	payload := map[string]interface{}{
		"model":      c.model,
		"messages":   reqMessages,
		"max_tokens": c.maxTokens,
		"stream":     true,
	}

	if c.temperature > 0 {
		payload["temperature"] = c.temperature
	}
	if c.topP > 0 && c.topP < 1 {
		payload["top_p"] = c.topP
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		close(eventCh)
		return eventCh, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", claudeAPIURL, bytes.NewReader(payloadBytes))
	if err != nil {
		close(eventCh)
		return eventCh, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", claudeAPIVersion)
	req.Header.Set("x-api-key", c.apiKey)

	// Send request
	go func() {
		defer close(eventCh)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			eventCh <- ai.StreamEvent{
				Type:  ai.EventTypeError,
				Error: fmt.Errorf("request failed: %w", err),
			}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			eventCh <- ai.StreamEvent{
				Type:  ai.EventTypeError,
				Error: fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes)),
			}
			return
		}

		// Parse SSE stream
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			// SSE format: "data: {...}"
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")

			// Check for stream end marker
			if data == "[DONE]" {
				eventCh <- ai.StreamEvent{
					Type: ai.EventTypeComplete,
					Done: true,
				}
				return
			}

			// Parse JSON event
			var event claudeStreamEvent
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				// Skip malformed events
				continue
			}

			// Handle different event types
			switch event.Type {
			case "content_block_delta":
				if event.Delta.Type == "text_delta" && event.Delta.Text != "" {
					eventCh <- ai.StreamEvent{
						Type:    ai.EventTypeChunk,
						Content: event.Delta.Text,
					}
				}

			case "message_stop":
				eventCh <- ai.StreamEvent{
					Type: ai.EventTypeComplete,
					Done: true,
				}
				return

			case "error":
				eventCh <- ai.StreamEvent{
					Type:  ai.EventTypeError,
					Error: fmt.Errorf("stream error: %s", event.Error.Message),
				}
				return
			}
		}

		if err := scanner.Err(); err != nil {
			eventCh <- ai.StreamEvent{
				Type:  ai.EventTypeError,
				Error: fmt.Errorf("stream read error: %w", err),
			}
		}
	}()

	return eventCh, nil
}

// claudeStreamEvent represents a Claude SSE event
type claudeStreamEvent struct {
	Type  string `json:"type"`
	Delta struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}
