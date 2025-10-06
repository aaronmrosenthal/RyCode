package providers

import "github.com/aaronmrosenthal/rycode/packages/tui-v2/internal/ai"

func init() {
	// Register provider constructors with the factory
	ai.RegisterProviders(
		func(apiKey string, config *ai.Config) (ai.Provider, error) {
			return NewClaudeProvider(apiKey, config)
		},
		func(apiKey string, config *ai.Config) (ai.Provider, error) {
			return NewOpenAIProvider(apiKey, config)
		},
	)
}
