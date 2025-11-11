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

// OllamaProvider implements the Provider interface for local Ollama
type OllamaProvider struct {
	config     *config.Config
	httpClient *http.Client
	baseURL    string
}

// ollamaRequest represents the Ollama API request format
type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// ollamaResponse represents the Ollama API response format
type ollamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider(cfg *config.Config) (*OllamaProvider, error) {
	// Get base URL from config or use default
	baseURL := "http://localhost:11434"
	if url, ok := cfg.AI.Custom["ollama_url"]; ok {
		baseURL = url
	}

	return &OllamaProvider{
		config:  cfg,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: cfg.AI.Timeout,
		},
	}, nil
}

// Name returns the provider name
func (p *OllamaProvider) Name() string {
	return "ollama"
}

// ValidateConfig checks if the configuration is valid
func (p *OllamaProvider) ValidateConfig() error {
	if p.config.AI.Model == "" {
		return fmt.Errorf("Ollama model is not specified")
	}

	return nil
}

// NewRequest creates a request from a prompt using provider's configured defaults
func (p *OllamaProvider) NewRequest(prompt string) *Request {
	return &Request{
		Prompt:      prompt,
		Model:       p.config.AI.Model,
		MaxTokens:   p.config.AI.MaxTokens,
		Temperature: p.config.AI.Temperature,
	}
}

// Generate sends a request to Ollama and returns the response
func (p *OllamaProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
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
func (p *OllamaProvider) generateOnce(ctx context.Context, req *Request) (*Response, error) {
	// Combine system prompt and user prompt
	prompt := req.Prompt
	if req.SystemPrompt != "" {
		prompt = fmt.Sprintf("System: %s\n\nUser: %s", req.SystemPrompt, req.Prompt)
	}

	// Build request payload
	apiReq := ollamaRequest{
		Model:  req.Model,
		Prompt: prompt,
		Stream: false,
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/api/generate", p.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	httpResp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request (is Ollama running?): %w", err)
	}
	defer httpResp.Body.Close()

	// Read response body
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Handle non-200 responses
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama API error (status %d): %s", httpResp.StatusCode, string(body))
	}

	// Parse response
	var apiResp ollamaResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &Response{
		Content:      apiResp.Response,
		TokensUsed:   0, // Ollama doesn't provide token counts in this simple format
		Model:        apiResp.Model,
		Provider:     "ollama",
		CostEstimate: 0.0, // Local model, no cost
		Metadata: map[string]interface{}{
			"created_at": apiResp.CreatedAt,
		},
	}, nil
}
