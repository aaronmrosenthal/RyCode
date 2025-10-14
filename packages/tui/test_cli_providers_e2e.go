package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/auth"
)

func main() {
	fmt.Println("=== E2E Test: All CLI Providers Authenticated ===\n")

	// Set up logging
	logFile, err := os.OpenFile("/tmp/rycode-e2e-cli-providers.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
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

	log("=== STARTING CLI PROVIDERS E2E TEST ===")
	log("Time: %s", time.Now().Format(time.RFC3339))
	log("Purpose: Validate ALL SOTA CLI providers are authenticated and accessible")
	log("")

	// Get project root dynamically
	projectRoot := getProjectRoot()
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
	log("\n[3] Validating authentication for all providers...")

	expectedProviders := map[string]int{
		"claude":  6,  // Expected model count
		"qwen":    7,
		"codex":   8,
		"gemini":  7,
	}

	authenticatedCount := 0
	failedProviders := []string{}
	missingProviders := []string{}

	// Check each expected provider
	for providerID, expectedModelCount := range expectedProviders {
		found := false
		for _, cliProv := range cliProviders {
			if cliProv.Provider == providerID {
				found = true
				authStatus, err := testApp.AuthBridge.CheckAuthStatus(ctx, providerID)
				if err != nil {
					log("    ✗ %s: Authentication check failed: %v", providerID, err)
					log("       → Check if %s API key is set", getProviderEnvVar(providerID))
					failedProviders = append(failedProviders, providerID)
					continue
				}

				if !authStatus.IsAuthenticated {
					log("    ✗ %s: NOT AUTHENTICATED", providerID)
					log("       → Set %s environment variable", getProviderEnvVar(providerID))
					log("       → Or run: rycode /auth to authenticate providers")
					failedProviders = append(failedProviders, providerID)
					continue
				}

				if authStatus.ModelsCount != expectedModelCount {
					log("    ⚠️  %s: AUTHENTICATED but model count mismatch (expected %d, got %d)",
						providerID, expectedModelCount, authStatus.ModelsCount)
				} else {
					log("    ✓ %s: AUTHENTICATED (%d models)", providerID, authStatus.ModelsCount)
				}

				authenticatedCount++
				break
			}
		}

		if !found {
			log("    ✗ %s: NOT FOUND in CLI providers", providerID)
			missingProviders = append(missingProviders, providerID)
		}
	}

	// Print summary
	log("\n=== TEST SUMMARY ===")
	log("Expected providers: %d", len(expectedProviders))
	log("Authenticated: %d", authenticatedCount)
	log("Failed: %d", len(failedProviders))
	log("Missing: %d", len(missingProviders))
	log("")

	if len(failedProviders) > 0 {
		log("✗ Failed providers:")
		for _, p := range failedProviders {
			log("  - %s", p)
		}
		log("")
	}

	if len(missingProviders) > 0 {
		log("✗ Missing providers:")
		for _, p := range missingProviders {
			log("  - %s", p)
		}
		log("")
	}

	log("Test logs saved to: /tmp/rycode-e2e-cli-providers.log")

	// Exit with appropriate code
	if len(failedProviders) > 0 || len(missingProviders) > 0 {
		log("\n❌ TEST FAILED: %d provider(s) not properly configured", len(failedProviders)+len(missingProviders))
		os.Exit(1)
	} else {
		log("\n✅ TEST PASSED: All %d SOTA providers are authenticated and ready!", authenticatedCount)
		os.Exit(0)
	}
}

// getProjectRoot dynamically discovers the project root directory
func getProjectRoot() string {
	// Try environment variable first
	if root := os.Getenv("PROJECT_ROOT"); root != "" {
		return root
	}

	// Try git root
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output))
	}

	// Fallback: try to find go.work or package.json going up from current directory
	dir, err := os.Getwd()
	if err == nil {
		for {
			// Check for go.work
			if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
				return dir
			}
			// Check for package.json with workspaces
			if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
				return dir
			}

			parent := filepath.Dir(dir)
			if parent == dir {
				break // Reached root
			}
			dir = parent
		}
	}

	// Last resort: hardcoded fallback
	return "/Users/aaron/Code/RyCode/RyCode"
}

// getProviderEnvVar returns the expected environment variable name for a provider
func getProviderEnvVar(providerID string) string {
	envVars := map[string]string{
		"claude":  "ANTHROPIC_API_KEY",
		"qwen":    "DASHSCOPE_API_KEY or QWEN_API_KEY",
		"codex":   "OPENAI_API_KEY",
		"gemini":  "GOOGLE_API_KEY or GEMINI_API_KEY",
	}

	if envVar, ok := envVars[providerID]; ok {
		return envVar
	}

	// Generic fallback
	return fmt.Sprintf("%s_API_KEY", strings.ToUpper(providerID))
}
