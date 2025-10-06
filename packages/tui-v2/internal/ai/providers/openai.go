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
	openAIAPIURL = "https://api.openai.com/v1/chat/completions"
)

// OpenAIProvider implements the AI Provider interface for OpenAI's GPT models
type OpenAIProvider struct {
	apiKey      string
	model       string
	maxTokens   int
	temperature float64
	topP        float64
	httpClient  *http.Client
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string, config *ai.Config) *OpenAIProvider {
	if config == nil {
		config = ai.DefaultConfig()
	}

	model := config.OpenAIModel
	if model == "" {
		model = "gpt-4o"
	}

	return &OpenAIProvider{
		apiKey:      apiKey,
		model:       model,
		maxTokens:   config.MaxTokens,
		temperature: config.Temperature,
		topP:        config.TopP,
		httpClient:  &http.Client{},
	}
}

// Name returns the provider name
func (o *OpenAIProvider) Name() string {
	return "OpenAI"
}

// Model returns the model identifier
func (o *OpenAIProvider) Model() string {
	return o.model
}

// Stream sends a prompt and streams back response tokens
func (o *OpenAIProvider) Stream(ctx context.Context, prompt string, messages []ai.Message) (<-chan ai.StreamEvent, error) {
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
		"model":    o.model,
		"messages": reqMessages,
		"stream":   true,
	}

	if o.maxTokens > 0 {
		payload["max_tokens"] = o.maxTokens
	}
	if o.temperature >= 0 {
		payload["temperature"] = o.temperature
	}
	if o.topP > 0 && o.topP <= 1 {
		payload["top_p"] = o.topP
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		close(eventCh)
		return eventCh, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", openAIAPIURL, bytes.NewReader(payloadBytes))
	if err != nil {
		close(eventCh)
		return eventCh, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	// Send request
	go func() {
		defer close(eventCh)

		resp, err := o.httpClient.Do(req)
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
			var event openAIStreamEvent
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				// Skip malformed events
				continue
			}

			// Extract content from delta
			if len(event.Choices) > 0 {
				choice := event.Choices[0]

				// Check for finish reason
				if choice.FinishReason != "" && choice.FinishReason != "null" {
					eventCh <- ai.StreamEvent{
						Type: ai.EventTypeComplete,
						Done: true,
					}
					return
				}

				// Send content delta
				if choice.Delta.Content != "" {
					eventCh <- ai.StreamEvent{
						Type:    ai.EventTypeChunk,
						Content: choice.Delta.Content,
					}
				}
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

// openAIStreamEvent represents an OpenAI SSE event
type openAIStreamEvent struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}
