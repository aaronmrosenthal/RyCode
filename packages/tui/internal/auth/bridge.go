package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// AuthStatus represents the authentication status from TypeScript
type AuthStatus struct {
	IsAuthenticated bool   `json:"isAuthenticated"`
	Provider        string `json:"provider"`
	ModelsCount     int    `json:"modelsCount"`
}

// ProviderHealth represents circuit breaker health from TypeScript
type ProviderHealth struct {
	Provider      string     `json:"provider"`
	Status        string     `json:"status"` // "healthy", "degraded", "down"
	FailureCount  int        `json:"failureCount"`
	NextAttemptAt *time.Time `json:"nextAttemptAt,omitempty"`
}

// CostSummary represents cost tracking data from TypeScript
type CostSummary struct {
	TodayCost  float64 `json:"todayCost"`
	MonthCost  float64 `json:"monthCost"`
	Projection float64 `json:"projection"`
	SavingsTip string  `json:"savingsTip,omitempty"`
}

// ProviderInfo represents provider information
type ProviderInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ModelsCount int    `json:"modelsCount"`
}

// AuthResult represents authentication result
type AuthResult struct {
	Provider    string `json:"provider"`
	ModelsCount int    `json:"modelsCount"`
	Message     string `json:"message"`
}

// AutoDetectResult represents auto-detection result
type AutoDetectResult struct {
	Message     string `json:"message"`
	Found       int    `json:"found"`
	Credentials []struct {
		Provider string `json:"provider"`
		Count    int    `json:"count"`
	} `json:"credentials"`
}

// Recommendation represents a model recommendation
type Recommendation struct {
	Provider  string  `json:"provider"`
	Model     string  `json:"model"`
	Score     float64 `json:"score"`
	Reasoning string  `json:"reasoning"`
}

// RecommendationsResult represents recommendation results
type RecommendationsResult struct {
	Recommendations []Recommendation `json:"recommendations"`
}

// Bridge provides access to the TypeScript authentication system
type Bridge struct {
	cliPath string
}

var debugLog *os.File

func init() {
	var err error
	debugLog, err = os.OpenFile("/tmp/rycode-debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		debugLog = nil
	}
}

func logDebug(format string, args ...interface{}) {
	if debugLog != nil {
		fmt.Fprintf(debugLog, format+"\n", args...)
		debugLog.Sync()
	}
}

// NewBridge creates a new authentication bridge
func NewBridge(projectRoot string) *Bridge {
	// Convert projectRoot to absolute path
	absProjectRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		// Fallback to original if abs fails
		absProjectRoot = projectRoot
	}

	cliPath := filepath.Join(absProjectRoot, "packages", "rycode", "src", "auth", "cli.ts")
	logDebug("DEBUG [NewBridge]: projectRoot=%s, absProjectRoot=%s, cliPath=%s", projectRoot, absProjectRoot, cliPath)
	return &Bridge{
		cliPath: cliPath,
	}
}

// runCLI executes a CLI command and returns the result
func (b *Bridge) runCLI(ctx context.Context, args ...string) ([]byte, error) {
	fullArgs := append([]string{"run", b.cliPath}, args...)
	cmd := exec.CommandContext(ctx, "bun", fullArgs...)

	// CRITICAL FIX: Set working directory to the project root
	// Extract project root from cliPath (remove "/packages/rycode/src/auth/cli.ts")
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(b.cliPath))))
	cmd.Dir = projectRoot

	// DEBUG: Log the command being run with working directory
	logDebug("DEBUG [bridge]: Running: bun %v (from %s)", fullArgs, projectRoot)

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Try to parse error from stderr
			var errorResp struct {
				Success bool   `json:"success"`
				Error   string `json:"error"`
			}
			logDebug("DEBUG [bridge]: Command failed with stderr: %s", string(exitErr.Stderr))
			if jsonErr := json.Unmarshal(exitErr.Stderr, &errorResp); jsonErr == nil {
				return nil, fmt.Errorf("auth CLI error: %s", errorResp.Error)
			}
			return nil, fmt.Errorf("auth CLI failed: %s", string(exitErr.Stderr))
		}
		logDebug("DEBUG [bridge]: Command failed: %v", err)
		return nil, fmt.Errorf("failed to run auth CLI: %w", err)
	}

	logDebug("DEBUG [bridge]: Got output (%d bytes)", len(output))
	return output, nil
}

// CheckAuthStatus checks if a provider is authenticated
func (b *Bridge) CheckAuthStatus(ctx context.Context, provider string) (*AuthStatus, error) {
	output, err := b.runCLI(ctx, "check", provider)
	if err != nil {
		return nil, err
	}

	var status AuthStatus
	if err := json.Unmarshal(output, &status); err != nil {
		return nil, fmt.Errorf("failed to parse auth status: %w", err)
	}

	return &status, nil
}

// Authenticate authenticates with a provider
func (b *Bridge) Authenticate(ctx context.Context, provider, apiKey string) (*AuthResult, error) {
	output, err := b.runCLI(ctx, "auth", provider, apiKey)
	if err != nil {
		return nil, err
	}

	var result AuthResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse auth result: %w", err)
	}

	return &result, nil
}

// GetCostSummary retrieves cost tracking summary
func (b *Bridge) GetCostSummary(ctx context.Context) (*CostSummary, error) {
	output, err := b.runCLI(ctx, "cost")
	if err != nil {
		return nil, err
	}

	var summary CostSummary
	if err := json.Unmarshal(output, &summary); err != nil {
		return nil, fmt.Errorf("failed to parse cost summary: %w", err)
	}

	return &summary, nil
}

// GetProviderHealth retrieves provider health status
func (b *Bridge) GetProviderHealth(ctx context.Context, provider string) (*ProviderHealth, error) {
	output, err := b.runCLI(ctx, "health", provider)
	if err != nil {
		return nil, err
	}

	var health ProviderHealth
	if err := json.Unmarshal(output, &health); err != nil {
		return nil, fmt.Errorf("failed to parse provider health: %w", err)
	}

	return &health, nil
}

// ListAuthenticatedProviders lists all authenticated providers
func (b *Bridge) ListAuthenticatedProviders(ctx context.Context) ([]ProviderInfo, error) {
	output, err := b.runCLI(ctx, "list")
	if err != nil {
		return nil, err
	}

	var result struct {
		Providers []ProviderInfo `json:"providers"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse provider list: %w", err)
	}

	return result.Providers, nil
}

// AutoDetect attempts to auto-detect credentials
func (b *Bridge) AutoDetect(ctx context.Context) (*AutoDetectResult, error) {
	output, err := b.runCLI(ctx, "auto-detect")
	if err != nil {
		return nil, err
	}

	var result AutoDetectResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse auto-detect result: %w", err)
	}

	return &result, nil
}

// AutoDetectProvider attempts to auto-detect credentials for a specific provider
func (b *Bridge) AutoDetectProvider(ctx context.Context, provider string) (*AuthResult, error) {
	// Run auto-detect first
	result, err := b.AutoDetect(ctx)
	if err != nil {
		return nil, err
	}

	// Check if provider was found
	for _, cred := range result.Credentials {
		if cred.Provider == provider {
			// Provider credentials detected, now authenticate
			return &AuthResult{
				Provider:    provider,
				ModelsCount: cred.Count,
				Message:     fmt.Sprintf("Auto-detected credentials for %s", provider),
			}, nil
		}
	}

	return nil, fmt.Errorf("no credentials found for provider: %s", provider)
}

// GetAuthStatus retrieves the full auth status for all providers
func (b *Bridge) GetAuthStatus(ctx context.Context) (*struct {
	Authenticated []ProviderInfo `json:"authenticated"`
}, error) {
	providers, err := b.ListAuthenticatedProviders(ctx)
	if err != nil {
		return nil, err
	}

	return &struct {
		Authenticated []ProviderInfo `json:"authenticated"`
	}{
		Authenticated: providers,
	}, nil
}

// GetRecommendations gets model recommendations for a task
func (b *Bridge) GetRecommendations(ctx context.Context, task string) ([]Recommendation, error) {
	args := []string{"recommendations"}
	if task != "" {
		args = append(args, task)
	}

	output, err := b.runCLI(ctx, args...)
	if err != nil {
		return nil, err
	}

	var result RecommendationsResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse recommendations: %w", err)
	}

	return result.Recommendations, nil
}

// CLIProviderInfo represents a CLI provider with models
type CLIProviderInfo struct {
	Provider string   `json:"provider"`
	Models   []string `json:"models"`
	Source   string   `json:"source"`
}

// GetCLIProviders retrieves available CLI providers with models
func (b *Bridge) GetCLIProviders(ctx context.Context) ([]CLIProviderInfo, error) {
	output, err := b.runCLI(ctx, "cli-providers")
	if err != nil {
		return nil, err
	}

	var result struct {
		Providers []CLIProviderInfo `json:"providers"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse CLI providers: %w", err)
	}

	return result.Providers, nil
}
