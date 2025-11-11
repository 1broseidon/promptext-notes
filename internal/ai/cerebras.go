package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/1broseidon/promptext-notes/internal/config"
)

const cerebrasAPIURL = "https://api.cerebras.ai/v1/chat/completions"

// CerebrasProvider implements the Provider interface for Cerebras
type CerebrasProvider struct {
	apiKey     string
	config     *config.Config
	httpClient *http.Client
}

// NewCerebrasProvider creates a new Cerebras provider
func NewCerebrasProvider(apiKey string, cfg *config.Config) (*CerebrasProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("cerebras API key is required")
	}

	return &CerebrasProvider{
		apiKey: apiKey,
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.AI.Timeout,
		},
	}, nil
}

// Name returns the provider name
func (p *CerebrasProvider) Name() string {
	return "cerebras"
}

// ValidateConfig checks if the configuration is valid
func (p *CerebrasProvider) ValidateConfig() error {
	if p.apiKey == "" {
		return fmt.Errorf("cerebras API key is not set")
	}

	if p.config.AI.Model == "" {
		return fmt.Errorf("cerebras model is not specified")
	}

	return nil
}

// NewRequest creates a request from a prompt using provider's configured defaults
func (p *CerebrasProvider) NewRequest(prompt string) *Request {
	return &Request{
		Prompt:      prompt,
		Model:       p.config.AI.Model,
		MaxTokens:   p.config.AI.MaxTokens,
		Temperature: p.config.AI.Temperature,
	}
}

// Generate sends a request to Cerebras and returns the response
func (p *CerebrasProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
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
func (p *CerebrasProvider) generateOnce(ctx context.Context, req *Request) (*Response, error) {
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

	// Build request payload (OpenAI-compatible format)
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
	httpReq, err := http.NewRequestWithContext(ctx, "POST", cerebrasAPIURL, bytes.NewBuffer(jsonData))
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
		if apiErr.Error.Message == "" {
			return nil, fmt.Errorf("cerebras API error (status %d): %s", httpResp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("cerebras API error: %s", apiErr.Error.Message)
	}

	// Parse response (OpenAI-compatible format)
	var apiResp openaiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract content from first choice
	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	content := apiResp.Choices[0].Message.Content

	return &Response{
		Content:      content,
		TokensUsed:   apiResp.Usage.TotalTokens,
		Model:        apiResp.Model,
		Provider:     "cerebras",
		CostEstimate: 0.0, // Cerebras offers free tier
		Metadata: map[string]interface{}{
			"prompt_tokens":     apiResp.Usage.PromptTokens,
			"completion_tokens": apiResp.Usage.CompletionTokens,
			"finish_reason":     apiResp.Choices[0].FinishReason,
			"id":                apiResp.ID,
		},
	}, nil
}
