package auth

import (
	"context"
	"testing"
	"time"
)

func TestBridge_CheckAuthStatus(t *testing.T) {
	bridge := NewBridge("../../../..")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test with a provider (may or may not be authenticated)
	status, err := bridge.CheckAuthStatus(ctx, "anthropic")
	if err != nil {
		t.Logf("Note: Auth check returned error (expected if not set up): %v", err)
		return
	}

	if status.Provider != "anthropic" {
		t.Errorf("Expected provider 'anthropic', got '%s'", status.Provider)
	}

	t.Logf("Auth status for anthropic: authenticated=%v, models=%d",
		status.IsAuthenticated, status.ModelsCount)
}

func TestBridge_GetCostSummary(t *testing.T) {
	bridge := NewBridge("../../../..")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	summary, err := bridge.GetCostSummary(ctx)
	if err != nil {
		t.Fatalf("Failed to get cost summary: %v", err)
	}

	if summary.TodayCost < 0 {
		t.Errorf("Invalid today's cost: %f", summary.TodayCost)
	}

	t.Logf("Cost summary: today=$%.2f, month=$%.2f, projection=$%.2f",
		summary.TodayCost, summary.MonthCost, summary.Projection)
}

func TestBridge_ListAuthenticatedProviders(t *testing.T) {
	bridge := NewBridge("../../../..")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	providers, err := bridge.ListAuthenticatedProviders(ctx)
	if err != nil {
		t.Fatalf("Failed to list providers: %v", err)
	}

	t.Logf("Found %d authenticated provider(s)", len(providers))
	for _, p := range providers {
		t.Logf("  - %s (%s): %d models", p.Name, p.ID, p.ModelsCount)
	}
}

func TestBridge_GetProviderHealth(t *testing.T) {
	bridge := NewBridge("../../../..")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	health, err := bridge.GetProviderHealth(ctx, "anthropic")
	if err != nil {
		t.Logf("Note: Health check returned error (expected if not set up): %v", err)
		return
	}

	if health.Status != "healthy" && health.Status != "degraded" && health.Status != "down" {
		t.Errorf("Invalid health status: %s", health.Status)
	}

	t.Logf("Health for anthropic: status=%s, failures=%d",
		health.Status, health.FailureCount)
}

func TestBridge_AutoDetect(t *testing.T) {
	bridge := NewBridge("../../../..")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := bridge.AutoDetect(ctx)
	if err != nil {
		t.Fatalf("Failed to auto-detect: %v", err)
	}

	t.Logf("Auto-detect: %s", result.Message)
	t.Logf("Found %d credential source(s)", result.Found)
	for _, cred := range result.Credentials {
		t.Logf("  - %s: %d key(s)", cred.Provider, cred.Count)
	}
}

func TestBridge_GetRecommendations(t *testing.T) {
	bridge := NewBridge("../../../..")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	recs, err := bridge.GetRecommendations(ctx, "code_generation")
	if err != nil {
		t.Fatalf("Failed to get recommendations: %v", err)
	}

	if len(recs) == 0 {
		t.Log("No recommendations available (may be expected if no providers authenticated)")
		return
	}

	t.Logf("Got %d recommendation(s) for code_generation:", len(recs))
	for i, rec := range recs {
		t.Logf("  %d. %s/%s (score: %.2f)", i+1, rec.Provider, rec.Model, rec.Score)
		t.Logf("     %s", rec.Reasoning)
	}
}
