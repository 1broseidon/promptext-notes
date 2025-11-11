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

const openaiAPIURL = "https://api.openai.com/v1/chat/completions"

// OpenAIProvider implements the Provider interface for OpenAI
type OpenAIProvider struct {
	apiKey     string
	config     *config.Config
	httpClient *http.Client
}

// openaiRequest represents the OpenAI API request format
type openaiRequest struct {
	Model       string          `json:"model"`
	Messages    []openaiMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

// openaiMessage represents a message in the conversation
type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openaiResponse represents the OpenAI API response format
type openaiResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []openaiChoice `json:"choices"`
	Usage   openaiUsage    `json:"usage"`
}

// openaiChoice represents a completion choice
type openaiChoice struct {
	Index        int           `json:"index"`
	Message      openaiMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

// openaiUsage represents token usage information
type openaiUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// openaiError represents an error response from OpenAI
type openaiError struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string, cfg *config.Config) (*OpenAIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("openAI API key is required")
	}

	return &OpenAIProvider{
		apiKey: apiKey,
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.AI.Timeout,
		},
	}, nil
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// ValidateConfig checks if the configuration is valid
func (p *OpenAIProvider) ValidateConfig() error {
	if p.apiKey == "" {
		return fmt.Errorf("openAI API key is not set")
	}

	if p.config.AI.Model == "" {
		return fmt.Errorf("openAI model is not specified")
	}

	return nil
}

// NewRequest creates a request from a prompt using provider's configured defaults
func (p *OpenAIProvider) NewRequest(prompt string) *Request {
	return &Request{
		Prompt:      prompt,
		Model:       p.config.AI.Model,
		MaxTokens:   p.config.AI.MaxTokens,
		Temperature: p.config.AI.Temperature,
	}
}

// Generate sends a request to OpenAI and returns the response
func (p *OpenAIProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
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
func (p *OpenAIProvider) generateOnce(ctx context.Context, req *Request) (*Response, error) {
	// Build messages array
	messages := []openaiMessage{}

	// Add system prompt if provided
	if req.SystemPrompt != "" {
		messages = append(messages, openaiMessage{
			Role:    "system",
			Content: req.SystemPrompt,
		})
	}

	// Add user prompt
	messages = append(messages, openaiMessage{
		Role:    "user",
		Content: req.Prompt,
	})

	// Build request payload
	apiReq := openaiRequest{
		Model:       req.Model,
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", openaiAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))

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
		var apiErr openaiError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", httpResp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("OpenAI API error: %s", apiErr.Error.Message)
	}

	// Parse response
	var apiResp openaiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract content from first choice
	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	content := apiResp.Choices[0].Message.Content

	// Calculate cost estimate
	costEstimate := p.estimateCost(apiResp.Model, apiResp.Usage.PromptTokens, apiResp.Usage.CompletionTokens)

	return &Response{
		Content:      content,
		TokensUsed:   apiResp.Usage.TotalTokens,
		Model:        apiResp.Model,
		Provider:     "openai",
		CostEstimate: costEstimate,
		Metadata: map[string]interface{}{
			"prompt_tokens":     apiResp.Usage.PromptTokens,
			"completion_tokens": apiResp.Usage.CompletionTokens,
			"finish_reason":     apiResp.Choices[0].FinishReason,
			"id":                apiResp.ID,
		},
	}, nil
}

// estimateCost calculates approximate cost based on token usage
func (p *OpenAIProvider) estimateCost(model string, inputTokens, outputTokens int) float64 {
	// Pricing as of January 2025 (per million tokens)
	var inputCost, outputCost float64

	switch {
	case strings.Contains(model, "gpt-4o-mini"):
		inputCost = 0.150
		outputCost = 0.600
	case strings.Contains(model, "gpt-4o"):
		inputCost = 2.50
		outputCost = 10.00
	case strings.Contains(model, "gpt-4-turbo"):
		inputCost = 10.00
		outputCost = 30.00
	case strings.Contains(model, "gpt-3.5-turbo"):
		inputCost = 0.50
		outputCost = 1.50
	default:
		// Unknown model, use gpt-4o-mini pricing as baseline
		inputCost = 0.150
		outputCost = 0.600
	}

	// Calculate cost in dollars
	cost := (float64(inputTokens) * inputCost / 1_000_000) +
		(float64(outputTokens) * outputCost / 1_000_000)

	return cost
}
