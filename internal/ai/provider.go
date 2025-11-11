package ai

import (
	"context"
	"fmt"

	"github.com/1broseidon/promptext-notes/internal/config"
)

// Provider defines the interface for AI service providers
type Provider interface {
	// Generate calls the AI provider with the given request and returns the response
	Generate(ctx context.Context, req *Request) (*Response, error)

	// Name returns the provider name (anthropic, openai, etc.)
	Name() string

	// ValidateConfig checks if required credentials and configuration are present
	ValidateConfig() error
}

// Request represents an AI generation request
type Request struct {
	// Prompt is the main content to send to the AI
	Prompt string

	// SystemPrompt is an optional system prompt (supported by some providers)
	SystemPrompt string

	// Model specifies which model to use (provider-specific)
	Model string

	// MaxTokens is the maximum number of tokens to generate
	MaxTokens int

	// Temperature controls randomness (0.0 to 1.0)
	Temperature float64
}

// Response represents an AI generation response
type Response struct {
	// Content is the generated text
	Content string

	// TokensUsed is the number of tokens consumed (if available)
	TokensUsed int

	// Model is the actual model used
	Model string

	// Provider is the provider name
	Provider string

	// CostEstimate is an optional cost estimate in USD
	CostEstimate float64

	// Metadata contains provider-specific information
	Metadata map[string]interface{}
}

// NewProvider creates a new AI provider based on the configuration
func NewProvider(cfg *config.Config) (Provider, error) {
	apiKey, err := cfg.GetAPIKey()
	if err != nil && cfg.AI.Provider != "ollama" {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	switch cfg.AI.Provider {
	case "anthropic":
		return NewAnthropicProvider(apiKey, cfg)
	case "openai":
		return NewOpenAIProvider(apiKey, cfg)
	case "cerebras":
		return NewCerebrasProvider(apiKey, cfg)
	case "groq":
		return NewGroqProvider(apiKey, cfg)
	case "ollama":
		return NewOllamaProvider(cfg)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", cfg.AI.Provider)
	}
}

// RequestFromConfig creates a Request from configuration and prompt
func RequestFromConfig(cfg *config.Config, prompt string) *Request {
	return &Request{
		Prompt:      prompt,
		Model:       cfg.AI.Model,
		MaxTokens:   cfg.AI.MaxTokens,
		Temperature: cfg.AI.Temperature,
	}
}
