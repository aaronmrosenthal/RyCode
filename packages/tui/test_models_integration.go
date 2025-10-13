package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aaronmrosenthal/rycode/internal/app"
	"github.com/aaronmrosenthal/rycode/internal/components/dialog"
)

func main() {
	fmt.Println("=== Testing Model Dialog Integration ===\n")

	// Set up debug logging
	debugLog, err := os.OpenFile("/tmp/rycode-models-test.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer debugLog.Close()

	// Get project root (assuming we're in packages/tui)
	projectRoot := "/Users/aaron/Code/RyCode/RyCode"

	fmt.Fprintf(debugLog, "=== MODELS INTEGRATION TEST ===\n")
	fmt.Fprintf(debugLog, "Project root: %s\n", projectRoot)

	// Create minimal app struct
	testApp := &app.App{}

	// Initialize auth bridge
	fmt.Println("Initializing auth bridge...")
	testApp.InitAuthBridge(projectRoot)
	fmt.Fprintf(debugLog, "Auth bridge initialized\n")

	// List all providers through the app
	fmt.Println("Listing providers...")
	providers, err := testApp.ListProviders(context.Background())
	if err != nil {
		fmt.Fprintf(debugLog, "ERROR: %v\n", err)
		fmt.Printf("ERROR listing providers: %v\n", err)
		return
	}

	fmt.Fprintf(debugLog, "Found %d providers:\n", len(providers))
	fmt.Printf("Found %d providers:\n", len(providers))

	totalModels := 0
	for _, provider := range providers {
		modelCount := len(provider.Models)
		totalModels += modelCount
		fmt.Fprintf(debugLog, "  - %s (%s): %d models\n", provider.Name, provider.ID, modelCount)
		fmt.Printf("  - %s (%s): %d models\n", provider.Name, provider.ID, modelCount)

		for modelID, model := range provider.Models {
			fmt.Fprintf(debugLog, "      * %s (%s)\n", model.Name, modelID)
		}
	}

	fmt.Fprintf(debugLog, "\nTotal models: %d\n", totalModels)
	fmt.Printf("\nTotal models: %d\n", totalModels)

	// Create model dialog
	fmt.Println("\nCreating model dialog...")
	modelDialog := dialog.NewModelDialog(testApp)
	fmt.Fprintf(debugLog, "Model dialog created: %v\n", modelDialog != nil)
	fmt.Printf("Model dialog created successfully!\n")

	fmt.Println("\n=== Test Complete ===")
	fmt.Printf("Check /tmp/rycode-models-test.log for details\n")
}
