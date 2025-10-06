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
	claudeAPIURL     = "https://api.anthropic.com/v1/messages"
	claudeAPIVersion = "2023-06-01"
)

// ClaudeProvider implements the AI Provider interface for Anthropic's Claude API
type ClaudeProvider struct {
	apiKey      *ai.SecureString // Encrypted in memory to prevent key extraction
	model       string
	maxTokens   int
	temperature float64
	topP        float64
	httpClient  *http.Client
}

// NewClaudeProvider creates a new Claude AI provider
func NewClaudeProvider(apiKey string, config *ai.Config) (*ClaudeProvider, error) {
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

	model := config.ClaudeModel
	if model == "" {
		model = "claude-opus-4-20250514"
	}

	return &ClaudeProvider{
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
func (c *ClaudeProvider) Name() string {
	return "Claude"
}

// Model returns the model identifier
func (c *ClaudeProvider) Model() string {
	return c.model
}

// Stream sends a prompt and streams back response tokens
func (c *ClaudeProvider) Stream(ctx context.Context, prompt string, messages []ai.Message) (<-chan ai.StreamEvent, error) {
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

	// Decrypt API key temporarily for HTTP request
	apiKey, err := c.apiKey.Reveal()
	if err != nil {
		close(eventCh)
		return eventCh, fmt.Errorf("failed to access API key: %w", err)
	}
	// Ensure key is zeroed after this function
	defer ai.ZeroString(apiKey)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", claudeAPIURL, bytes.NewReader(payloadBytes))
	if err != nil {
		close(eventCh)
		return eventCh, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", claudeAPIVersion)
	req.Header.Set("x-api-key", apiKey)

	// Send request
	go func() {
		defer close(eventCh)

		resp, err := c.httpClient.Do(req)
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
			var event claudeStreamEvent
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

			// Handle different event types
			switch event.Type {
			case "content_block_delta":
				if event.Delta.Type == "text_delta" && event.Delta.Text != "" {
					select {
					case eventCh <- ai.StreamEvent{Type: ai.EventTypeChunk, Content: event.Delta.Text}:
					case <-ctx.Done():
						return
					}
				}

			case "message_stop":
				select {
				case eventCh <- ai.StreamEvent{Type: ai.EventTypeComplete, Done: true}:
				case <-ctx.Done():
				}
				return

			case "error":
				select {
				case eventCh <- ai.StreamEvent{Type: ai.EventTypeError, Error: fmt.Errorf("stream error: %s", event.Error.Message)}:
				case <-ctx.Done():
				}
				return
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
