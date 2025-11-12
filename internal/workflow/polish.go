package workflow

import (
	"context"
	"fmt"

	"github.com/1broseidon/promptext-notes/internal/ai"
	"github.com/1broseidon/promptext-notes/internal/config"
)

// DefaultPolishPrompt is the default prompt for polishing changelogs
const DefaultPolishPrompt = `You are a technical writer. Polish this changelog for public consumption.

CHANGELOG:
%s

Rules:
- Keep the exact same items (do NOT add or remove anything)
- Only improve the wording
- Avoid "we", "we've", "our"
- Use active voice: "Updated X", "Fixed Y"
- Keep it concise

Output only the polished changelog.`

// PolishChangelog takes a draft changelog and polishes it using a second AI model
func PolishChangelog(ctx context.Context, draftChangelog string, cfg *config.Config) (string, error) {
	if !cfg.AI.Polish.Enabled {
		return draftChangelog, nil // Polish not enabled, return draft as-is
	}

	// Determine polish provider
	polishProvider := cfg.GetPolishProvider()
	polishModel := cfg.GetPolishModel()

	// Get polish API key
	polishAPIKey, err := cfg.GetPolishAPIKey()
	if err != nil {
		return "", fmt.Errorf("failed to get polish API key: %w", err)
	}

	// Create polish config
	polishCfg := &config.Config{
		AI: config.AIConfig{
			Provider:    polishProvider,
			Model:       polishModel,
			APIKeyEnv:   cfg.GetPolishAPIKeyEnv(),
			MaxTokens:   cfg.AI.Polish.PolishMaxTokens,
			Temperature: cfg.AI.Polish.PolishTemperature,
			Timeout:     cfg.AI.Timeout,
			Retry:       cfg.AI.Retry,
			Custom:      cfg.AI.Custom,
		},
	}

	// Create polish provider
	var polishAI ai.Provider
	switch polishProvider {
	case "anthropic":
		polishAI, err = ai.NewAnthropicProvider(polishAPIKey, polishCfg)
	case "openai":
		polishAI, err = ai.NewOpenAIProvider(polishAPIKey, polishCfg)
	case "cerebras":
		polishAI, err = ai.NewCerebrasProvider(polishAPIKey, polishCfg)
	case "groq":
		polishAI, err = ai.NewGroqProvider(polishAPIKey, polishCfg)
	case "openrouter":
		polishAI, err = ai.NewOpenRouterProvider(polishAPIKey, polishCfg)
	case "ollama":
		polishAI, err = ai.NewOllamaProvider(polishCfg)
	default:
		return "", fmt.Errorf("unsupported polish provider: %s", polishProvider)
	}

	if err != nil {
		return "", fmt.Errorf("failed to create polish AI provider: %w", err)
	}

	// Prepare polish prompt
	polishPrompt := cfg.AI.Polish.PolishPrompt
	if polishPrompt == "" {
		polishPrompt = fmt.Sprintf(DefaultPolishPrompt, draftChangelog)
	} else {
		polishPrompt = fmt.Sprintf(polishPrompt, draftChangelog)
	}

	// Create polish request
	req := &ai.Request{
		Prompt:      polishPrompt,
		Model:       polishModel,
		MaxTokens:   cfg.AI.Polish.PolishMaxTokens,
		Temperature: cfg.AI.Polish.PolishTemperature,
	}

	// Generate polished changelog
	resp, err := polishAI.Generate(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to polish changelog: %w", err)
	}

	return resp.Content, nil
}
