package providers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/ai"
)

const (
	openAIAPIURL = "https://api.openai.com/v1/chat/completions"
)

// OpenAIProvider implements the AI Provider interface for OpenAI's GPT models
type OpenAIProvider struct {
	apiKey      *ai.SecureString // Encrypted in memory to prevent key extraction
	model       string
	maxTokens   int
	temperature float64
	topP        float64
	httpClient  *http.Client
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string, config *ai.Config) (*OpenAIProvider, error) {
	if config == nil {
		config = ai.DefaultConfig()
	}

	// Encrypt API key in memory
	secureKey, err := ai.NewSecureString(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to secure API key: %w", err)
	}

	// Zero out the plaintext parameter (best effort)
	ai.ZeroString(apiKey)

	model := config.OpenAIModel
	if model == "" {
		model = "gpt-4o"
	}

	return &OpenAIProvider{
		apiKey:      secureKey,
		model:       model,
		maxTokens:   config.MaxTokens,
		temperature: config.Temperature,
		topP:        config.TopP,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // 2 minute total timeout
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 30 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				IdleConnTimeout:       90 * time.Second,
				MaxIdleConns:          10,
				MaxIdleConnsPerHost:   2,
			},
		},
	}, nil
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
	eventCh := make(chan ai.StreamEvent, 10) // Small buffer with backpressure

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

	// Decrypt API key temporarily for HTTP request
	apiKey, err := o.apiKey.Reveal()
	if err != nil {
		close(eventCh)
		return eventCh, fmt.Errorf("failed to access API key: %w", err)
	}
	// Ensure key is zeroed after this function
	defer ai.ZeroString(apiKey)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", openAIAPIURL, bytes.NewReader(payloadBytes))
	if err != nil {
		close(eventCh)
		return eventCh, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send request
	go func() {
		defer close(eventCh)

		resp, err := o.httpClient.Do(req)
		if err != nil {
			select {
			case eventCh <- ai.StreamEvent{Type: ai.EventTypeError, Error: fmt.Errorf("request failed: %w", err)}:
			case <-ctx.Done():
			}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 10*1024)) // Max 10KB error
			select {
			case eventCh <- ai.StreamEvent{Type: ai.EventTypeError, Error: fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))}:
			case <-ctx.Done():
			}
			return
		}

		// Parse SSE stream with larger buffer
		scanner := bufio.NewScanner(resp.Body)
		buf := make([]byte, 1024*1024) // 1MB max line
		scanner.Buffer(buf, len(buf))

		malformedCount := 0
		for scanner.Scan() {
			// Check if context cancelled
			select {
			case <-ctx.Done():
				return
			default:
			}

			line := scanner.Text()

			// SSE format: "data: {...}"
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")

			// Check for stream end marker
			if data == "[DONE]" {
				select {
				case eventCh <- ai.StreamEvent{Type: ai.EventTypeComplete, Done: true}:
				case <-ctx.Done():
				}
				return
			}

			// Parse JSON event
			var event openAIStreamEvent
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				malformedCount++
				// Log first few errors, then fail if too many
				if malformedCount <= 3 {
					// Skip malformed events but track them
					continue
				}
				select {
				case eventCh <- ai.StreamEvent{Type: ai.EventTypeError, Error: fmt.Errorf("too many malformed events (%d): %w", malformedCount, err)}:
				case <-ctx.Done():
				}
				return
			}

			// Extract content from delta
			if len(event.Choices) > 0 {
				choice := event.Choices[0]

				// Check for finish reason
				if choice.FinishReason != "" && choice.FinishReason != "null" {
					select {
					case eventCh <- ai.StreamEvent{Type: ai.EventTypeComplete, Done: true}:
					case <-ctx.Done():
					}
					return
				}

				// Send content delta
				if choice.Delta.Content != "" {
					select {
					case eventCh <- ai.StreamEvent{Type: ai.EventTypeChunk, Content: choice.Delta.Content}:
					case <-ctx.Done():
						return
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			select {
			case eventCh <- ai.StreamEvent{Type: ai.EventTypeError, Error: fmt.Errorf("stream read error: %w", err)}:
			case <-ctx.Done():
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
