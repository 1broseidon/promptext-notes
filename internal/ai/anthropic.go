package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/1broseidon/promptext-notes/internal/config"
)

const anthropicAPIURL = "https://api.anthropic.com/v1/messages"

// AnthropicProvider implements the Provider interface for Anthropic Claude
type AnthropicProvider struct {
	apiKey     string
	config     *config.Config
	httpClient *http.Client
}

// anthropicRequest represents the Anthropic API request format
type anthropicRequest struct {
	Model       string              `json:"model"`
	MaxTokens   int                 `json:"max_tokens"`
	Temperature float64             `json:"temperature,omitempty"`
	Messages    []anthropicMessage  `json:"messages"`
	System      string              `json:"system,omitempty"`
}

// anthropicMessage represents a message in the conversation
type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// anthropicResponse represents the Anthropic API response format
type anthropicResponse struct {
	ID      string              `json:"id"`
	Type    string              `json:"type"`
	Role    string              `json:"role"`
	Content []anthropicContent  `json:"content"`
	Model   string              `json:"model"`
	Usage   anthropicUsage      `json:"usage"`
}

// anthropicContent represents content blocks in the response
type anthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// anthropicUsage represents token usage information
type anthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// anthropicError represents an error response from Anthropic
type anthropicError struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

// NewAnthropicProvider creates a new Anthropic provider
func NewAnthropicProvider(apiKey string, cfg *config.Config) (*AnthropicProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Anthropic API key is required")
	}

	return &AnthropicProvider{
		apiKey: apiKey,
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.AI.Timeout,
		},
	}, nil
}

// Name returns the provider name
func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

// ValidateConfig checks if the configuration is valid
func (p *AnthropicProvider) ValidateConfig() error {
	if p.apiKey == "" {
		return fmt.Errorf("Anthropic API key is not set")
	}

	if p.config.AI.Model == "" {
		return fmt.Errorf("Anthropic model is not specified")
	}

	return nil
}

// Generate sends a request to Anthropic and returns the response
func (p *AnthropicProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
	if err := p.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	var response *Response
	var generateErr error

	// Use retry with backoff
	err := RetryWithBackoff(ctx, p.config, func(ctx context.Context) error {
		response, generateErr = p.generateOnce(ctx, req)
		return generateErr
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

// generateOnce performs a single generation attempt
func (p *AnthropicProvider) generateOnce(ctx context.Context, req *Request) (*Response, error) {
	// Normalize model name
	model := p.normalizeModel(req.Model)

	// Build request payload
	apiReq := anthropicRequest{
		Model:       model,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Messages: []anthropicMessage{
			{
				Role:    "user",
				Content: req.Prompt,
			},
		},
	}

	// Add system prompt if provided
	if req.SystemPrompt != "" {
		apiReq.System = req.SystemPrompt
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", anthropicAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	// Allow custom anthropic version from config
	if version, ok := p.config.AI.Custom["anthropic_version"]; ok {
		httpReq.Header.Set("anthropic-version", version)
	}

	// Send request
	httpResp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	// Read response body
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Handle non-200 responses
	if httpResp.StatusCode != http.StatusOK {
		var apiErr anthropicError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", httpResp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("Anthropic API error: %s", apiErr.Error.Message)
	}

	// Parse response
	var apiResp anthropicResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract text from content blocks
	var content strings.Builder
	for _, block := range apiResp.Content {
		if block.Type == "text" {
			content.WriteString(block.Text)
		}
	}

	// Calculate cost estimate (approximate pricing as of 2024)
	costEstimate := p.estimateCost(model, apiResp.Usage.InputTokens, apiResp.Usage.OutputTokens)

	return &Response{
		Content:      content.String(),
		TokensUsed:   apiResp.Usage.InputTokens + apiResp.Usage.OutputTokens,
		Model:        apiResp.Model,
		Provider:     "anthropic",
		CostEstimate: costEstimate,
		Metadata: map[string]interface{}{
			"input_tokens":  apiResp.Usage.InputTokens,
			"output_tokens": apiResp.Usage.OutputTokens,
			"id":            apiResp.ID,
		},
	}, nil
}

// normalizeModel converts common model names to Anthropic's format
func (p *AnthropicProvider) normalizeModel(model string) string {
	// Map friendly names to API model names
	modelMap := map[string]string{
		"claude-haiku-4-5":  "claude-3-5-haiku-20241022",
		"claude-sonnet-4-5": "claude-3-5-sonnet-20241022",
		"claude-opus-4":     "claude-opus-4-20250514",
		"haiku":             "claude-3-5-haiku-20241022",
		"sonnet":            "claude-3-5-sonnet-20241022",
		"opus":              "claude-opus-4-20250514",
	}

	if normalized, ok := modelMap[model]; ok {
		return normalized
	}

	// Return as-is if not in the map (assume it's already a valid model ID)
	return model
}

// estimateCost calculates approximate cost based on token usage
func (p *AnthropicProvider) estimateCost(model string, inputTokens, outputTokens int) float64 {
	// Pricing as of January 2025 (per million tokens)
	var inputCost, outputCost float64

	switch {
	case strings.Contains(model, "haiku"):
		inputCost = 0.80
		outputCost = 4.00
	case strings.Contains(model, "sonnet"):
		inputCost = 3.00
		outputCost = 15.00
	case strings.Contains(model, "opus"):
		inputCost = 15.00
		outputCost = 75.00
	default:
		// Unknown model, use Haiku pricing as baseline
		inputCost = 0.80
		outputCost = 4.00
	}

	// Calculate cost in dollars
	cost := (float64(inputTokens) * inputCost / 1_000_000) +
		(float64(outputTokens) * outputCost / 1_000_000)

	return cost
}
