package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the complete configuration for promptext-notes
type Config struct {
	Version string        `yaml:"version"`
	AI      AIConfig      `yaml:"ai"`
	Output  OutputConfig  `yaml:"output"`
	Filters FiltersConfig `yaml:"filters"`
}

// AIConfig holds AI provider configuration
type AIConfig struct {
	Provider    string            `yaml:"provider"`
	Model       string            `yaml:"model"`
	APIKeyEnv   string            `yaml:"api_key_env"`
	MaxTokens   int               `yaml:"max_tokens"`
	Temperature float64           `yaml:"temperature"`
	Timeout     time.Duration     `yaml:"timeout"`
	Retry       RetryConfig       `yaml:"retry"`
	Custom      map[string]string `yaml:"custom"`
}

// RetryConfig defines retry behavior
type RetryConfig struct {
	Attempts     int           `yaml:"attempts"`
	Backoff      string        `yaml:"backoff"`
	InitialDelay time.Duration `yaml:"initial_delay"`
}

// OutputConfig defines output formatting
type OutputConfig struct {
	Format   string   `yaml:"format"`
	Sections []string `yaml:"sections"`
	Template string   `yaml:"template"`
}

// FiltersConfig defines filtering rules
type FiltersConfig struct {
	Files   FileFilters   `yaml:"files"`
	Commits CommitFilters `yaml:"commits"`
}

// FileFilters defines file inclusion/exclusion patterns
type FileFilters struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}

// CommitFilters defines commit filtering rules
type CommitFilters struct {
	ExcludeAuthors  []string `yaml:"exclude_authors"`
	ExcludePatterns []string `yaml:"exclude_patterns"`
}

// Load reads and parses a configuration file
func Load(path string) (*Config, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults for missing values
	applyDefaults(&config)

	return &config, nil
}

// LoadOrDefault attempts to load config from path, returns defaults if not found
func LoadOrDefault(path string) *Config {
	config, err := Load(path)
	if err != nil {
		return Default()
	}
	return config
}

// Default returns a configuration with sensible defaults
func Default() *Config {
	return &Config{
		Version: "1",
		AI: AIConfig{
			Provider:    "anthropic",
			Model:       "claude-haiku-4-5",
			APIKeyEnv:   "ANTHROPIC_API_KEY",
			MaxTokens:   8000,
			Temperature: 0.3,
			Timeout:     30 * time.Second,
			Retry: RetryConfig{
				Attempts:     3,
				Backoff:      "exponential",
				InitialDelay: 2 * time.Second,
			},
			Custom: make(map[string]string),
		},
		Output: OutputConfig{
			Format: "keepachangelog",
			Sections: []string{
				"breaking",
				"added",
				"changed",
				"fixed",
				"docs",
			},
		},
		Filters: FiltersConfig{
			Files: FileFilters{
				Include: []string{
					"*.go",
					"*.md",
					"*.yml",
					"*.yaml",
					"*.json",
				},
				Exclude: []string{
					"*_test.go",
					"vendor/*",
					"node_modules/*",
					".git/*",
				},
			},
			Commits: CommitFilters{
				ExcludeAuthors: []string{
					"dependabot[bot]",
					"renovate[bot]",
				},
				ExcludePatterns: []string{
					"^Merge pull request",
					"^Merge branch",
				},
			},
		},
	}
}

// applyDefaults fills in missing values with defaults
func applyDefaults(config *Config) {
	defaults := Default()

	if config.Version == "" {
		config.Version = defaults.Version
	}

	// AI defaults
	if config.AI.Provider == "" {
		config.AI.Provider = defaults.AI.Provider
	}
	if config.AI.MaxTokens == 0 {
		config.AI.MaxTokens = defaults.AI.MaxTokens
	}
	if config.AI.Temperature == 0 {
		config.AI.Temperature = defaults.AI.Temperature
	}
	if config.AI.Timeout == 0 {
		config.AI.Timeout = defaults.AI.Timeout
	}
	if config.AI.Retry.Attempts == 0 {
		config.AI.Retry = defaults.AI.Retry
	}
	if config.AI.Custom == nil {
		config.AI.Custom = make(map[string]string)
	}

	// Set default API key env var based on provider
	if config.AI.APIKeyEnv == "" {
		config.AI.APIKeyEnv = getDefaultAPIKeyEnv(config.AI.Provider)
	}

	// Set default model based on provider
	if config.AI.Model == "" {
		config.AI.Model = getDefaultModel(config.AI.Provider)
	}

	// Output defaults
	if config.Output.Format == "" {
		config.Output.Format = defaults.Output.Format
	}
	if len(config.Output.Sections) == 0 {
		config.Output.Sections = defaults.Output.Sections
	}

	// Filters defaults
	if len(config.Filters.Files.Include) == 0 {
		config.Filters.Files.Include = defaults.Filters.Files.Include
	}
	if len(config.Filters.Files.Exclude) == 0 {
		config.Filters.Files.Exclude = defaults.Filters.Files.Exclude
	}
	if len(config.Filters.Commits.ExcludeAuthors) == 0 {
		config.Filters.Commits.ExcludeAuthors = defaults.Filters.Commits.ExcludeAuthors
	}
	if len(config.Filters.Commits.ExcludePatterns) == 0 {
		config.Filters.Commits.ExcludePatterns = defaults.Filters.Commits.ExcludePatterns
	}
}

// getDefaultAPIKeyEnv returns the default environment variable for API key
func getDefaultAPIKeyEnv(provider string) string {
	switch provider {
	case "anthropic":
		return "ANTHROPIC_API_KEY"
	case "openai":
		return "OPENAI_API_KEY"
	case "cerebras":
		return "CEREBRAS_API_KEY"
	case "groq":
		return "GROQ_API_KEY"
	case "ollama":
		return "" // No API key needed for local Ollama
	default:
		return ""
	}
}

// getDefaultModel returns the default model for a provider
func getDefaultModel(provider string) string {
	switch provider {
	case "anthropic":
		return "claude-haiku-4-5"
	case "openai":
		return "gpt-4o-mini"
	case "cerebras":
		return "llama-3.3-70b"
	case "groq":
		return "llama-3.3-70b-versatile"
	case "ollama":
		return "llama3.2"
	default:
		return ""
	}
}

// GetAPIKey retrieves the API key from the environment
func (c *Config) GetAPIKey() (string, error) {
	if c.AI.APIKeyEnv == "" {
		return "", nil // No API key required (e.g., Ollama)
	}

	key := os.Getenv(c.AI.APIKeyEnv)
	if key == "" {
		return "", fmt.Errorf("API key not found in environment variable: %s", c.AI.APIKeyEnv)
	}

	return key, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	validProviders := map[string]bool{
		"anthropic": true,
		"openai":    true,
		"cerebras":  true,
		"groq":      true,
		"ollama":    true,
	}

	if !validProviders[c.AI.Provider] {
		return fmt.Errorf("invalid AI provider: %s (supported: anthropic, openai, cerebras, groq, ollama)", c.AI.Provider)
	}

	if c.AI.MaxTokens <= 0 {
		return fmt.Errorf("max_tokens must be positive, got: %d", c.AI.MaxTokens)
	}

	if c.AI.Temperature < 0 || c.AI.Temperature > 1 {
		return fmt.Errorf("temperature must be between 0 and 1, got: %.2f", c.AI.Temperature)
	}

	validBackoffs := map[string]bool{
		"exponential": true,
		"linear":      true,
		"constant":    true,
	}

	if !validBackoffs[c.AI.Retry.Backoff] {
		return fmt.Errorf("invalid backoff strategy: %s (supported: exponential, linear, constant)", c.AI.Retry.Backoff)
	}

	return nil
}
