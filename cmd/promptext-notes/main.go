package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/1broseidon/promptext-notes/internal/ai"
	"github.com/1broseidon/promptext-notes/internal/config"
	"github.com/1broseidon/promptext-notes/internal/git"
	"github.com/1broseidon/promptext-notes/internal/workflow"
)

func main() {
	// Parse flags
	version := flag.String("version", "", "Version to generate notes for (e.g., v0.7.4)")
	sinceTag := flag.String("since", "", "Generate notes since this tag (auto-detects if empty)")
	output := flag.String("output", "", "Output file (prints to stdout if empty)")
	configFile := flag.String("config", ".promptext-notes.yml", "Configuration file path")

	// AI flags
	generate := flag.Bool("generate", false, "Generate AI-enhanced changelog (requires AI provider)")
	aiPrompt := flag.Bool("ai-prompt", false, "Generate prompt for AI to enhance release notes (legacy mode)")
	providerFlag := flag.String("provider", "", "AI provider (anthropic, openai, cerebras, groq, ollama)")
	modelFlag := flag.String("model", "", "AI model to use")

	// Other flags
	quiet := flag.Bool("quiet", false, "Suppress progress messages")

	flag.Parse()

	// Check if we're in a git repository
	if !git.IsGitRepository() {
		log.Fatal("Error: Not a git repository. Please run this command from within a git repository.")
	}

	// Load configuration
	cfg := config.LoadOrDefault(*configFile)

	// Override config with CLI flags
	if *providerFlag != "" {
		cfg.AI.Provider = *providerFlag
		// Update API key env var to match the new provider
		cfg.AI.APIKeyEnv = config.GetDefaultAPIKeyEnv(*providerFlag)
	}
	if *modelFlag != "" {
		cfg.AI.Model = *modelFlag
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Get from tag
	fromTag := *sinceTag
	if fromTag == "" {
		var err error
		fromTag, err = git.GetLastTag()
		if err != nil {
			log.Fatalf("Failed to get last tag: %v", err)
		}
	}

	// Create AI provider if needed
	var provider ai.Provider
	var err error

	if *generate {
		provider, err = ai.NewProvider(cfg)
		if err != nil {
			log.Fatalf("Failed to create AI provider: %v", err)
		}
	}

	// Set up workflow options
	opts := workflow.GenerateOptions{
		Version:      *version,
		SinceTag:     fromTag,
		Output:       *output,
		UseAI:        *generate,
		AIPromptOnly: *aiPrompt,
		Verbose:      !*quiet,
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.AI.Timeout*2)
	defer cancel()

	// Generate release notes
	outputText, err := workflow.GenerateReleaseNotes(ctx, opts, provider)
	if err != nil {
		log.Fatalf("Failed to generate release notes: %v", err)
	}

	// Write output
	if *output != "" {
		if err := os.WriteFile(*output, []byte(outputText), 0644); err != nil {
			log.Fatalf("Failed to write output: %v", err)
		}
		if !*quiet {
			fmt.Fprintf(os.Stderr, "✅ Written to %s\n", *output)
		}
	} else {
		fmt.Print(outputText)
	}
}
