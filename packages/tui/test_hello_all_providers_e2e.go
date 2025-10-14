package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/auth"
)

// SessionResponse is the response from creating a session
type SessionResponse struct {
	ID         string `json:"id"`
	ProviderID string `json:"providerId"`
	ModelID    string `json:"modelId"`
	Error      string `json:"error,omitempty"`
	Message    string `json:"message,omitempty"`
}

// MessageRequest is a request to send a message
type MessageRequest struct {
	SessionID string `json:"sessionId"`
	Message   string `json:"message"`
}

// MessageResponse is the response from sending a message
type MessageResponse struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

// ProviderInfo holds information about a provider for testing
type ProviderInfo struct {
	ID           string
	DisplayName  string
	Models       []string
	DefaultModel string
}

func main() {
	fmt.Println("=== E2E Test: Hello to All SOTA Providers ===\n")

	// Set up logging
	logFile, err := os.OpenFile("/tmp/rycode-e2e-hello-all.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	log := func(format string, args ...interface{}) {
		msg := fmt.Sprintf(format, args...)
		fmt.Println(msg)
		fmt.Fprintf(logFile, msg+"\n")
		logFile.Sync()
	}

	log("=== STARTING HELLO ALL PROVIDERS E2E TEST ===")
	log("Time: %s", time.Now().Format(time.RFC3339))
	log("Purpose: Validate ALL SOTA models respond to messages")
	log("")

	// Get project root
	projectRoot := "/Users/aaron/Code/RyCode/RyCode"
	log("Project root: %s", projectRoot)

	// Create app with auth bridge
	log("\n[1] Creating app instance...")
	testApp := &app.App{
		AuthBridge: auth.NewBridge(projectRoot),
	}
	log("    ✓ App created with auth bridge")

	// Get all CLI providers
	ctx := context.Background()
	log("\n[2] Loading CLI providers...")
	cliProviders, err := testApp.AuthBridge.GetCLIProviders(ctx)
	if err != nil {
		log("    ✗ ERROR: Failed to get CLI providers: %v", err)
		os.Exit(1)
	}
	log("    ✓ Found %d CLI provider configs", len(cliProviders))

	// Check authentication for each provider
	log("\n[3] Checking authentication status...")
	authenticatedProviders := []ProviderInfo{}

	for _, cliProv := range cliProviders {
		authStatus, err := testApp.AuthBridge.CheckAuthStatus(ctx, cliProv.Provider)
		if err != nil {
			log("    - %s: ERROR - %v", cliProv.Provider, err)
			continue
		}

		if authStatus.IsAuthenticated {
			displayName := getProviderDisplayName(cliProv.Provider)
			defaultModel := getDefaultModelForProvider(cliProv.Provider, cliProv.Models)

			providerInfo := ProviderInfo{
				ID:           cliProv.Provider,
				DisplayName:  displayName,
				Models:       cliProv.Models,
				DefaultModel: defaultModel,
			}
			authenticatedProviders = append(authenticatedProviders, providerInfo)

			log("    - %s: ✓ AUTHENTICATED (%d models, default: %s)",
				displayName, authStatus.ModelsCount, defaultModel)
		} else {
			log("    - %s: ✗ not authenticated", cliProv.Provider)
		}
	}

	if len(authenticatedProviders) == 0 {
		log("\n⚠️  FATAL: No authenticated CLI providers found!")
		log("    Cannot test without authenticated providers")
		log("    Run: rycode /auth to authenticate providers")
		os.Exit(1)
	}

	log("\n    Total authenticated providers: %d", len(authenticatedProviders))
	log("    Providers to test: %s",
		strings.Join(getProviderNames(authenticatedProviders), ", "))

	// Test each provider by sending "hello" message
	log("\n[4] Testing message responses from each provider...")
	log("    Test message: \"hello\"")
	log("")

	passedProviders := []string{}
	failedProviders := []string{}

	apiBaseURL := "http://127.0.0.1:4096"

	for i, provider := range authenticatedProviders {
		log("  [%d/%d] Testing %s (model: %s)...",
			i+1, len(authenticatedProviders), provider.DisplayName, provider.DefaultModel)

		// Step 1: Create session
		sessionID, err := createSession(apiBaseURL, provider.ID, provider.DefaultModel)
		if err != nil {
			log("      ✗ FAILED to create session: %v", err)
			failedProviders = append(failedProviders, provider.DisplayName)
			continue
		}
		log("      ✓ Session created: %s", sessionID)

		// Step 2: Send "hello" message
		// Don't specify model - let API use default provider resolution
		// The API will automatically pick up authenticated providers
		response, err := sendMessage(apiBaseURL, sessionID, "hello", "", "")
		if err != nil {
			log("      ✗ FAILED to send message: %v", err)
			failedProviders = append(failedProviders, provider.DisplayName)

			// Clean up session
			deleteSession(apiBaseURL, sessionID)
			continue
		}

		// Step 3: Validate response
		if len(response) == 0 {
			log("      ✗ FAILED: Empty response")
			failedProviders = append(failedProviders, provider.DisplayName)
		} else {
			// Truncate long responses for display
			displayResponse := response
			if len(displayResponse) > 100 {
				displayResponse = displayResponse[:100] + "..."
			}
			log("      ✓ SUCCESS: Got response (%d chars)", len(response))
			log("      Response preview: %s", displayResponse)
			passedProviders = append(passedProviders, provider.DisplayName)
		}

		// Step 4: Clean up session
		err = deleteSession(apiBaseURL, sessionID)
		if err != nil {
			log("      ⚠️  Warning: Failed to delete session: %v", err)
		} else {
			log("      ✓ Session cleaned up")
		}

		log("")
	}

	// Print summary
	log("\n=== TEST SUMMARY ===")
	log("Total providers tested: %d", len(authenticatedProviders))
	log("Passed: %d", len(passedProviders))
	log("Failed: %d", len(failedProviders))
	log("")

	if len(passedProviders) > 0 {
		log("✓ Passed providers:")
		for _, p := range passedProviders {
			log("  - %s", p)
		}
		log("")
	}

	if len(failedProviders) > 0 {
		log("✗ Failed providers:")
		for _, p := range failedProviders {
			log("  - %s", p)
		}
		log("")
	}

	log("Test logs saved to: /tmp/rycode-e2e-hello-all.log")

	// Exit with appropriate code
	if len(failedProviders) > 0 {
		log("\n❌ TEST FAILED: %d provider(s) did not respond correctly", len(failedProviders))
		os.Exit(1)
	} else {
		log("\n✅ TEST PASSED: All %d providers responded successfully!", len(passedProviders))
		os.Exit(0)
	}
}

// createSession creates a new session with the API
func createSession(baseURL, providerID, modelID string) (string, error) {
	// Session creation doesn't require provider/model - those are specified when sending messages
	reqBody := map[string]interface{}{}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(
		baseURL+"/session",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var sessionResp struct {
		ID    string `json:"id"`
		Error string `json:"error,omitempty"`
	}

	if err := json.Unmarshal(body, &sessionResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w (body: %s)", err, string(body))
	}

	if sessionResp.Error != "" {
		return "", fmt.Errorf("API error: %s", sessionResp.Error)
	}

	if sessionResp.ID == "" {
		return "", fmt.Errorf("no session ID in response")
	}

	return sessionResp.ID, nil
}

// sendMessage sends a message to a session and returns the response
func sendMessage(baseURL, sessionID, message, providerID, modelID string) (string, error) {
	reqBody := map[string]interface{}{
		"parts": []map[string]string{
			{
				"type": "text",
				"text": message,
			},
		},
	}

	// Only include model if specified (let API use default otherwise)
	if providerID != "" && modelID != "" {
		reqBody["model"] = map[string]string{
			"providerID": providerID,
			"modelID":    modelID,
		}
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/session/%s/message", baseURL, sessionID),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// The response is a message object with parts array
	// Extract the assistant's text response
	var response struct {
		Parts []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"parts"`
		Info struct {
			Error string `json:"error,omitempty"`
		} `json:"info"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w (body: %s)", err, string(body))
	}

	if response.Info.Error != "" {
		return "", fmt.Errorf("API error: %s", response.Info.Error)
	}

	// Collect text from all text parts
	var textParts []string
	for _, part := range response.Parts {
		if part.Type == "text" && part.Text != "" {
			textParts = append(textParts, part.Text)
		}
	}

	if len(textParts) == 0 {
		// Debug: Show what we actually received
		return "", fmt.Errorf("no text content in response (parts: %d, body preview: %s)",
			len(response.Parts), truncate(string(body), 200))
	}

	return strings.Join(textParts, "\n"), nil
}

// deleteSession deletes a session
func deleteSession(baseURL, sessionID string) error {
	req, err := http.NewRequest("DELETE", baseURL+"/session/"+sessionID, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d (body: %s)", resp.StatusCode, string(body))
	}

	return nil
}

// getProviderDisplayName returns the display name for a provider
func getProviderDisplayName(providerID string) string {
	names := map[string]string{
		"claude":    "Claude",
		"anthropic": "Claude",
		"codex":     "Codex",
		"openai":    "Codex",
		"gemini":    "Gemini",
		"google":    "Gemini",
		"grok":      "Grok",
		"xai":       "Grok",
		"qwen":      "Qwen",
	}
	if name, ok := names[providerID]; ok {
		return name
	}
	return providerID
}

// getDefaultModelForProvider returns the best default model for a provider
func getDefaultModelForProvider(providerID string, availableModels []string) string {
	// Priority order for each provider (latest SOTA models as of 2025)
	priorities := map[string][]string{
		"claude": {
			"claude-sonnet-4-5",
			"claude-opus-4-1",
			"claude-sonnet-4",
			"claude-3-7-sonnet",
			"claude-3-5-sonnet-20241022",
		},
		"codex": {
			"gpt-5",
			"o3",
			"gpt-5-mini",
			"gpt-4-5",
			"gpt-4o",
		},
		"gemini": {
			"gemini-2.5-pro",
			"gemini-2.5-flash",
			"gemini-2.5-flash-lite",
			"gemini-exp-1206",
		},
		"grok": {
			"grok-beta",
			"grok-2-1212",
		},
		"qwen": {
			"qwen3-max",
			"qwen3-thinking-2507",
			"qwen3-next",
			"qwen3-omni",
		},
	}

	// Check if we have priorities for this provider
	if prefs, ok := priorities[providerID]; ok {
		for _, modelID := range prefs {
			if contains(availableModels, modelID) {
				return modelID
			}
		}
	}

	// Fallback: return first available model
	if len(availableModels) > 0 {
		return availableModels[0]
	}

	return ""
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// getProviderNames extracts display names from provider info slice
func getProviderNames(providers []ProviderInfo) []string {
	names := make([]string, len(providers))
	for i, p := range providers {
		names[i] = p.DisplayName
	}
	return names
}

// truncate truncates a string to a maximum length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// mapCLIProviderToAPIProvider maps CLI provider IDs to actual API provider IDs
func mapCLIProviderToAPIProvider(cliProviderID string) string {
	// CLI uses friendly names, but API uses canonical provider names
	mapping := map[string]string{
		"claude":  "anthropic",
		"codex":   "openai",
		"gemini":  "google",
		"grok":    "xai",
		"qwen":    "qwen",
	}

	if apiProviderID, ok := mapping[cliProviderID]; ok {
		return apiProviderID
	}

	// Return as-is if no mapping exists (already canonical)
	return cliProviderID
}
