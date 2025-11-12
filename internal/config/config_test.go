package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	config := Default()

	if config.Version != "1" {
		t.Errorf("Expected version 1, got %s", config.Version)
	}

	if config.AI.Provider != "anthropic" {
		t.Errorf("Expected anthropic provider, got %s", config.AI.Provider)
	}

	if config.AI.Model != "claude-haiku-4-5" {
		t.Errorf("Expected claude-haiku-4-5 model, got %s", config.AI.Model)
	}

	if config.AI.MaxTokens != 8000 {
		t.Errorf("Expected 8000 max tokens, got %d", config.AI.MaxTokens)
	}

	if config.Output.Format != "keepachangelog" {
		t.Errorf("Expected keepachangelog format, got %s", config.Output.Format)
	}
}

func TestLoad(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yml")

	configContent := `version: "1"
ai:
  provider: openai
  model: gpt-4o
  api_key_env: OPENAI_API_KEY
  max_tokens: 4000
  temperature: 0.5
  timeout: 60s
  retry:
    attempts: 5
    backoff: linear
    initial_delay: 1s

output:
  format: keepachangelog
  sections:
    - breaking
    - added
    - fixed

filters:
  files:
    include:
      - "*.go"
    exclude:
      - "*_test.go"
  commits:
    exclude_authors:
      - "bot"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Validate loaded values
	if config.AI.Provider != "openai" {
		t.Errorf("Expected openai provider, got %s", config.AI.Provider)
	}

	if config.AI.Model != "gpt-4o" {
		t.Errorf("Expected gpt-4o model, got %s", config.AI.Model)
	}

	if config.AI.MaxTokens != 4000 {
		t.Errorf("Expected 4000 max tokens, got %d", config.AI.MaxTokens)
	}

	if config.AI.Temperature != 0.5 {
		t.Errorf("Expected 0.5 temperature, got %.2f", config.AI.Temperature)
	}

	if config.AI.Timeout != 60*time.Second {
		t.Errorf("Expected 60s timeout, got %v", config.AI.Timeout)
	}

	if config.AI.Retry.Attempts != 5 {
		t.Errorf("Expected 5 retry attempts, got %d", config.AI.Retry.Attempts)
	}

	if config.AI.Retry.Backoff != "linear" {
		t.Errorf("Expected linear backoff, got %s", config.AI.Retry.Backoff)
	}

	if len(config.Output.Sections) != 3 {
		t.Errorf("Expected 3 output sections, got %d", len(config.Output.Sections))
	}

	if len(config.Filters.Files.Include) != 1 {
		t.Errorf("Expected 1 include filter, got %d", len(config.Filters.Files.Include))
	}
}

func TestLoadNonExistent(t *testing.T) {
	_, err := Load("/nonexistent/config.yml")
	if err == nil {
		t.Error("Expected error when loading nonexistent file")
	}
}

func TestLoadOrDefault(t *testing.T) {
	// Test with nonexistent file - should return defaults
	config := LoadOrDefault("/nonexistent/config.yml")
	if config.AI.Provider != "anthropic" {
		t.Errorf("Expected default anthropic provider, got %s", config.AI.Provider)
	}

	// Test with existing file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yml")

	configContent := `version: "1"
ai:
  provider: groq
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	config = LoadOrDefault(configPath)
	if config.AI.Provider != "groq" {
		t.Errorf("Expected groq provider, got %s", config.AI.Provider)
	}
}

func TestApplyDefaults(t *testing.T) {
	config := &Config{
		AI: AIConfig{
			Provider: "openai",
			// Other fields left empty
		},
	}

	applyDefaults(config)

	// Should have defaults for empty fields
	if config.AI.MaxTokens != 8000 {
		t.Errorf("Expected default max tokens, got %d", config.AI.MaxTokens)
	}

	if config.AI.Model != "gpt-4o-mini" {
		t.Errorf("Expected default OpenAI model, got %s", config.AI.Model)
	}

	if config.AI.APIKeyEnv != "OPENAI_API_KEY" {
		t.Errorf("Expected default OpenAI API key env, got %s", config.AI.APIKeyEnv)
	}
}

func TestGetDefaultModel(t *testing.T) {
	tests := []struct {
		provider string
		expected string
	}{
		{"anthropic", "claude-haiku-4-5"},
		{"openai", "gpt-4o-mini"},
		{"cerebras", "llama-3.3-70b"},
		{"groq", "llama-3.3-70b-versatile"},
		{"ollama", "llama3.2"},
		{"unknown", ""},
	}

	for _, tt := range tests {
		result := getDefaultModel(tt.provider)
		if result != tt.expected {
			t.Errorf("getDefaultModel(%s) = %s, expected %s", tt.provider, result, tt.expected)
		}
	}
}

func TestGetDefaultAPIKeyEnv(t *testing.T) {
	tests := []struct {
		provider string
		expected string
	}{
		{"anthropic", "ANTHROPIC_API_KEY"},
		{"openai", "OPENAI_API_KEY"},
		{"cerebras", "CEREBRAS_API_KEY"},
		{"groq", "GROQ_API_KEY"},
		{"ollama", ""},
		{"unknown", ""},
	}

	for _, tt := range tests {
		result := getDefaultAPIKeyEnv(tt.provider)
		if result != tt.expected {
			t.Errorf("getDefaultAPIKeyEnv(%s) = %s, expected %s", tt.provider, result, tt.expected)
		}
	}
}

func TestGetAPIKey(t *testing.T) {
	config := &Config{
		AI: AIConfig{
			APIKeyEnv: "TEST_API_KEY",
		},
	}

	// Test with missing key
	_, err := config.GetAPIKey()
	if err == nil {
		t.Error("Expected error when API key not set")
	}

	// Test with key present
	os.Setenv("TEST_API_KEY", "test-key-value")
	defer os.Unsetenv("TEST_API_KEY")

	key, err := config.GetAPIKey()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if key != "test-key-value" {
		t.Errorf("Expected test-key-value, got %s", key)
	}

	// Test with no API key required (Ollama)
	config.AI.APIKeyEnv = ""
	key, err = config.GetAPIKey()
	if err != nil {
		t.Errorf("Unexpected error for no API key: %v", err)
	}
	if key != "" {
		t.Errorf("Expected empty key for Ollama, got %s", key)
	}
}

func TestAutoExcludeMeta(t *testing.T) {
	// Test that meta files are auto-excluded when AutoExcludeMeta is true
	config := &Config{
		Filters: FiltersConfig{
			Files: FileFilters{
				AutoExcludeMeta: true,
				Exclude: []string{
					"*_test.go",
					"vendor/*",
				},
			},
		},
	}

	applyDefaults(config)

	// Check that meta exclusions were merged
	expectedExclusions := []string{
		"*_test.go",
		"vendor/*",
		"CHANGELOG.md",
		"README.md",
		".github/**",
		".vscode/**",
		".idea/**",
		"*.example.*",
		".promptext-notes*.yml",
		"**/.gitignore",
		"**/.*ignore",
	}

	if len(config.Filters.Files.Exclude) != len(expectedExclusions) {
		t.Errorf("Expected %d exclusions, got %d", len(expectedExclusions), len(config.Filters.Files.Exclude))
	}

	// Check all expected exclusions are present
	exclusionMap := make(map[string]bool)
	for _, excl := range config.Filters.Files.Exclude {
		exclusionMap[excl] = true
	}

	for _, expected := range expectedExclusions {
		if !exclusionMap[expected] {
			t.Errorf("Expected exclusion not found: %s", expected)
		}
	}
}

func TestAutoExcludeMetaDisabled(t *testing.T) {
	// Test that meta files are NOT auto-excluded when AutoExcludeMeta is false
	config := &Config{
		Filters: FiltersConfig{
			Files: FileFilters{
				AutoExcludeMeta: false,
				Exclude: []string{
					"*_test.go",
					"vendor/*",
				},
			},
		},
	}

	applyDefaults(config)

	// Check that ONLY the original exclusions remain (no meta exclusions added)
	expectedExclusions := []string{
		"*_test.go",
		"vendor/*",
	}

	if len(config.Filters.Files.Exclude) != len(expectedExclusions) {
		t.Errorf("Expected %d exclusions, got %d (meta exclusions should not be added)", len(expectedExclusions), len(config.Filters.Files.Exclude))
	}
}

func TestMergeUnique(t *testing.T) {
	tests := []struct {
		name     string
		a        []string
		b        []string
		expected []string
	}{
		{
			name:     "No duplicates",
			a:        []string{"a", "b"},
			b:        []string{"c", "d"},
			expected: []string{"a", "b", "c", "d"},
		},
		{
			name:     "With duplicates",
			a:        []string{"a", "b", "c"},
			b:        []string{"b", "c", "d"},
			expected: []string{"a", "b", "c", "d"},
		},
		{
			name:     "Empty slices",
			a:        []string{},
			b:        []string{},
			expected: []string{},
		},
		{
			name:     "First empty",
			a:        []string{},
			b:        []string{"a", "b"},
			expected: []string{"a", "b"},
		},
		{
			name:     "Second empty",
			a:        []string{"a", "b"},
			b:        []string{},
			expected: []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeUnique(tt.a, tt.b)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
			}

			// Convert to map for easy comparison
			resultMap := make(map[string]bool)
			for _, item := range result {
				resultMap[item] = true
			}

			for _, expected := range tt.expected {
				if !resultMap[expected] {
					t.Errorf("Expected item %s not found in result", expected)
				}
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		expectErr bool
	}{
		{
			name:      "Valid config",
			config:    Default(),
			expectErr: false,
		},
		{
			name: "Invalid provider",
			config: &Config{
				AI: AIConfig{
					Provider:    "invalid",
					MaxTokens:   8000,
					Temperature: 0.3,
					Retry: RetryConfig{
						Backoff: "exponential",
					},
				},
			},
			expectErr: true,
		},
		{
			name: "Invalid max tokens",
			config: &Config{
				AI: AIConfig{
					Provider:    "anthropic",
					MaxTokens:   -100,
					Temperature: 0.3,
					Retry: RetryConfig{
						Backoff: "exponential",
					},
				},
			},
			expectErr: true,
		},
		{
			name: "Invalid temperature",
			config: &Config{
				AI: AIConfig{
					Provider:    "anthropic",
					MaxTokens:   8000,
					Temperature: 1.5,
					Retry: RetryConfig{
						Backoff: "exponential",
					},
				},
			},
			expectErr: true,
		},
		{
			name: "Invalid backoff",
			config: &Config{
				AI: AIConfig{
					Provider:    "anthropic",
					MaxTokens:   8000,
					Temperature: 0.3,
					Retry: RetryConfig{
						Backoff: "invalid",
					},
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectErr && err == nil {
				t.Error("Expected validation error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}
