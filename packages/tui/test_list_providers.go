package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aaronmrosenthal/rycode/internal/auth"
)

func main() {
	// Get project root (two levels up from packages/tui)
	projectRoot := "/Users/aaron/Code/RyCode/RyCode"

	// Create auth bridge
	authBridge := auth.NewBridge(projectRoot)

	// Call GetCLIProviders directly
	cliProviders, err := authBridge.GetCLIProviders(context.Background())

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("SUCCESS: Got %d CLI providers\n", len(cliProviders))
	for _, p := range cliProviders {
		fmt.Printf("  Provider: %s (%s) - %d models\n", p.Provider, p.Source, len(p.Models))
		for _, model := range p.Models {
			fmt.Printf("    - %s\n", model)
		}
	}
}
