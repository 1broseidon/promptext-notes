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
	Polish      PolishConfig      `yaml:"polish"`
}

// PolishConfig defines 2-stage polish workflow configuration
type PolishConfig struct {
	Enabled           bool    `yaml:"enabled"`
	PolishModel       string  `yaml:"polish_model"`       // Model for stage 2 (stage 1 uses ai.model)
	PolishProvider    string  `yaml:"polish_provider"`    // Optional: different provider for polish (defaults to ai.provider)
	PolishAPIKeyEnv   string  `yaml:"polish_api_key_env"` // Optional: API key env var (auto-detected from provider)
	PolishPrompt      string  `yaml:"polish_prompt"`      // Custom polish prompt (optional)
	PolishMaxTokens   int     `yaml:"polish_max_tokens"`  // Max tokens for polish stage
	PolishTemperature float64 `yaml:"polish_temperature"` // Temperature for polish stage
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
	Include         []string `yaml:"include"`
	Exclude         []string `yaml:"exclude"`
	AutoExcludeMeta bool     `yaml:"auto_exclude_meta"` // Auto-exclude CI/config/meta files (default: true)
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

// GetDefaultMetaExclusions returns file patterns that are auto-excluded when AutoExcludeMeta is true
func GetDefaultMetaExclusions() []string {
	return []string{
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
			Polish: PolishConfig{
				Enabled:           false,
				PolishModel:       "", // Uses ai.model if not specified
				PolishProvider:    "", // Uses ai.provider if not specified
				PolishAPIKeyEnv:   "", // Auto-detected from provider
				PolishPrompt:      "", // Uses default prompt
				PolishMaxTokens:   4000,
				PolishTemperature: 0.3,
			},
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
				AutoExcludeMeta: true, // Enable meta file exclusion by default
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

	// Apply auto-exclude-meta (merge with existing exclusions)
	// Note: AutoExcludeMeta defaults to true if not explicitly set to false
	if config.Filters.Files.AutoExcludeMeta {
		config.Filters.Files.Exclude = mergeUnique(config.Filters.Files.Exclude, GetDefaultMetaExclusions())
	}

	if len(config.Filters.Commits.ExcludeAuthors) == 0 {
		config.Filters.Commits.ExcludeAuthors = defaults.Filters.Commits.ExcludeAuthors
	}
	if len(config.Filters.Commits.ExcludePatterns) == 0 {
		config.Filters.Commits.ExcludePatterns = defaults.Filters.Commits.ExcludePatterns
	}
}

// mergeUnique merges two string slices, removing duplicates
func mergeUnique(a, b []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(a)+len(b))

	for _, item := range a {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	for _, item := range b {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// GetDefaultAPIKeyEnv returns the default environment variable for API key
func GetDefaultAPIKeyEnv(provider string) string {
	switch provider {
	case "anthropic":
		return "ANTHROPIC_API_KEY"
	case "openai":
		return "OPENAI_API_KEY"
	case "cerebras":
		return "CEREBRAS_API_KEY"
	case "groq":
		return "GROQ_API_KEY"
	case "openrouter":
		return "OPENROUTER_API_KEY"
	case "ollama":
		return "" // No API key needed for local Ollama
	default:
		return ""
	}
}

// getDefaultAPIKeyEnv is a private wrapper for backwards compatibility
func getDefaultAPIKeyEnv(provider string) string {
	return GetDefaultAPIKeyEnv(provider)
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
	case "openrouter":
		return "openai/gpt-4o-mini" // Cost-effective default
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
		"anthropic":  true,
		"openai":     true,
		"cerebras":   true,
		"groq":       true,
		"openrouter": true,
		"ollama":     true,
	}

	if !validProviders[c.AI.Provider] {
		return fmt.Errorf("invalid AI provider: %s (supported: anthropic, openai, cerebras, groq, openrouter, ollama)", c.AI.Provider)
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

	// Validate polish config if enabled
	if c.AI.Polish.Enabled {
		polishProvider := c.GetPolishProvider()
		if !validProviders[polishProvider] {
			return fmt.Errorf("invalid polish provider: %s (supported: anthropic, openai, cerebras, groq, openrouter, ollama)", polishProvider)
		}
	}

	return nil
}

// GetPolishProvider returns the effective polish provider (defaults to main provider)
func (c *Config) GetPolishProvider() string {
	if c.AI.Polish.PolishProvider != "" {
		return c.AI.Polish.PolishProvider
	}
	return c.AI.Provider
}

// GetPolishModel returns the effective polish model (defaults to main model)
func (c *Config) GetPolishModel() string {
	if c.AI.Polish.PolishModel != "" {
		return c.AI.Polish.PolishModel
	}
	return c.AI.Model
}

// GetPolishAPIKeyEnv returns the API key env var for polish provider
func (c *Config) GetPolishAPIKeyEnv() string {
	if c.AI.Polish.PolishAPIKeyEnv != "" {
		return c.AI.Polish.PolishAPIKeyEnv
	}
	return GetDefaultAPIKeyEnv(c.GetPolishProvider())
}

// GetPolishAPIKey retrieves the polish API key from the environment
func (c *Config) GetPolishAPIKey() (string, error) {
	polishProvider := c.GetPolishProvider()

	// Same provider as main - use main API key
	if polishProvider == c.AI.Provider {
		return c.GetAPIKey()
	}

	// Different provider - get its API key
	apiKeyEnv := c.GetPolishAPIKeyEnv()
	if apiKeyEnv == "" {
		return "", nil // No API key required (e.g., Ollama)
	}

	key := os.Getenv(apiKeyEnv)
	if key == "" {
		return "", fmt.Errorf("polish API key not found in environment variable: %s", apiKeyEnv)
	}

	return key, nil
}
