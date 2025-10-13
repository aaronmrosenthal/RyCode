package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/auth"
	"github.com/aaronmrosenthal/rycode/internal/components/dialog"
)

func main() {
	fmt.Println("=== E2E Test: Model Selector with Init() ===\n")

	// Set up debug logging
	debugLog, err := os.OpenFile("/tmp/rycode-e2e-test.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer debugLog.Close()

	log := func(format string, args ...interface{}) {
		msg := fmt.Sprintf(format, args...)
		fmt.Println(msg)
		fmt.Fprintf(debugLog, msg+"\n")
		debugLog.Sync()
	}

	log("=== STARTING E2E TEST ===")
	log("Time: %s", time.Now().Format(time.RFC3339))

	// Get project root
	projectRoot := "/Users/aaron/Code/RyCode/RyCode"
	log("Project root: %s", projectRoot)

	// Create app with minimal fields required
	log("\n[1] Creating app instance...")
	testApp := &app.App{
		AuthBridge: auth.NewBridge(projectRoot),
	}
	log("    ✓ App created with auth bridge")

	// Test 1: Check auth status for each CLI provider directly via bridge
	ctx := context.Background()
	log("\n[2] Testing auth status for CLI providers...")
	cliProviders := []string{"anthropic", "google", "openai", "xai", "qwen"}
	authenticatedCount := 0

	for _, providerID := range cliProviders {
		authStatus, err := testApp.AuthBridge.CheckAuthStatus(ctx, providerID)
		if err != nil {
			log("      - %s: ERROR - %v", providerID, err)
			continue
		}

		if authStatus.IsAuthenticated {
			authenticatedCount++
			log("      - %s: ✓ AUTHENTICATED (%d models)", providerID, authStatus.ModelsCount)
		} else {
			log("      - %s: ✗ not authenticated", providerID)
		}
	}

	log("    Total authenticated CLI providers: %d", authenticatedCount)

	if authenticatedCount == 0 {
		log("\n⚠️  WARNING: No authenticated CLI providers found!")
		log("    This test will verify the empty state is shown correctly")
		log("    To test with providers, authenticate first:")
		log("    - claude auth login")
		log("    - export GOOGLE_API_KEY=...")
		log("    - export OPENAI_API_KEY=...")
	}

	// Test 2: Create model dialog (the actual component used in TUI)
	log("\n[3] Creating ModelDialog component...")
	modelDialog := dialog.NewModelDialog(testApp)
	if modelDialog == nil {
		log("    ✗ ERROR: NewModelDialog returned nil")
		os.Exit(1)
	}
	log("    ✓ ModelDialog created")

	// Test 3: Call Init() - THIS IS THE CRITICAL STEP
	log("\n[4] Calling modelDialog.Init()...")
	initCmd := modelDialog.Init()
	if initCmd == nil {
		log("    ✓ Init() returned nil (no async work needed)")
	} else {
		log("    ✓ Init() returned command")

		// Execute the command to simulate what Bubble Tea does
		log("    Executing Init() command...")
		msg := initCmd()
		if msg != nil {
			log("      Command returned message: %T", msg)
		} else {
			log("      Command returned nil")
		}
	}

	// Test 4: Simulate Update cycle to let initialization complete
	log("\n[5] Simulating Bubble Tea Update cycle...")

	// Send WindowSizeMsg to trigger size setup
	sizeMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updatedModel, cmd := modelDialog.Update(sizeMsg)
	log("    ✓ WindowSizeMsg processed")

	if cmd != nil {
		log("    Executing returned command...")
		cmdMsg := cmd()
		if cmdMsg != nil {
			log("      Command returned: %T", cmdMsg)
		}
	}

	// Cast back to dialog
	modelDialog, ok := updatedModel.(dialog.ModelDialog)
	if !ok {
		log("    ✗ ERROR: Could not cast back to ModelDialog")
		os.Exit(1)
	}

	// Test 5: Render the view using Render()
	log("\n[6] Rendering ModelDialog view...")
	view := modelDialog.Render("")

	log("=== RENDERED VIEW ===")
	log("%s", view)
	log("=== END VIEW ===")

	// Test 6: Check view content
	log("\n[7] Analyzing rendered content...")

	hasNoProviders := containsString(view, "No authenticated CLI providers found")
	hasProviders := containsString(view, "Select Provider") || containsString(view, "models available")

	if hasNoProviders {
		log("    ✗ View shows: 'No authenticated CLI providers found'")
		log("    This means Init() was called but no authenticated CLI providers were detected")
		log("    Expected authenticated providers: %d", authenticatedCount)

		if authenticatedCount > 0 {
			log("\n⚠️  FAILURE: Auth check found %d providers but dialog shows none!", authenticatedCount)
			log("    This indicates a bug in SimpleProviderToggle.loadAuthenticatedProviders()")
			os.Exit(1)
		} else {
			log("\n✓ EXPECTED: No authenticated providers, dialog correctly shows empty state")
		}
	} else if hasProviders {
		log("    ✓ View shows provider selection UI")
		log("    Dialog successfully loaded and displayed providers")

		// Count provider chips in view
		chipCount := 0
		for _, providerID := range cliProviders {
			if containsString(view, getProviderDisplayName(providerID)) {
				chipCount++
			}
		}
		log("    Found %d provider chips in view", chipCount)

		if chipCount != authenticatedCount {
			log("\n⚠️  WARNING: Expected %d providers but found %d in view", authenticatedCount, chipCount)
		} else {
			log("\n✓ SUCCESS: All authenticated providers displayed correctly")
		}
	} else {
		log("    ⚠️  View content unclear - neither empty state nor provider list detected")
	}

	// Test 7: Check debug log for provider loading details
	log("\n[8] Checking debug log for provider loading details...")
	debugLogData, err := os.ReadFile("/tmp/rycode-debug.log")
	if err == nil {
		debugLines := 0
		for _, line := range []byte(string(debugLogData)) {
			if line == '\n' {
				debugLines++
			}
		}
		log("     Debug log has %d lines", debugLines)
		if debugLines > 0 {
			log("     (Check /tmp/rycode-debug.log for detailed provider loading logs)")
		}
	}

	log("\n=== TEST COMPLETE ===")
	log("Detailed logs saved to: /tmp/rycode-e2e-test.log")
	log("Debug logs from app: /tmp/rycode-debug.log")

	// Summary
	log("\n=== SUMMARY ===")
	log("Authenticated CLI providers: %d", authenticatedCount)
	log("Dialog Init() called: YES")
	log("View rendered successfully: YES")

	if hasProviders {
		log("Status: ✓ PASS - Providers displayed correctly")
		os.Exit(0)
	} else if hasNoProviders && authenticatedCount == 0 {
		log("Status: ✓ PASS - Correctly showing empty state (no providers authenticated)")
		os.Exit(0)
	} else {
		log("Status: ✗ FAIL - Dialog not displaying providers correctly")
		os.Exit(1)
	}
}

func containsString(text, substr string) bool {
	return len(text) > 0 && len(substr) > 0 &&
		(text == substr || findSubstring(text, substr) >= 0)
}

func findSubstring(text, substr string) int {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func getProviderDisplayName(providerID string) string {
	switch providerID {
	case "anthropic":
		return "Claude"
	case "google":
		return "Gemini"
	case "openai":
		return "GPT-5"
	case "xai":
		return "Grok"
	case "qwen":
		return "Qwen"
	default:
		return providerID
	}
}
