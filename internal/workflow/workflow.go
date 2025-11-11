package workflow

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/1broseidon/promptext-notes/internal/ai"
	"github.com/1broseidon/promptext-notes/internal/analyzer"
	aicontext "github.com/1broseidon/promptext-notes/internal/context"
	"github.com/1broseidon/promptext-notes/internal/generator"
	"github.com/1broseidon/promptext-notes/internal/git"
	"github.com/1broseidon/promptext-notes/internal/prompt"
)

// GenerateOptions contains options for release notes generation
type GenerateOptions struct {
	Version      string
	SinceTag     string
	Output       string
	UseAI        bool
	AIPromptOnly bool
	Verbose      bool
}

// GenerateReleaseNotes orchestrates the full release notes generation process
func GenerateReleaseNotes(ctx context.Context, opts GenerateOptions, provider ai.Provider) (string, error) {
	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "📊 Analyzing changes since %s...\n", opts.SinceTag)
	}

	// Get changed files
	changedFiles, err := git.GetChangedFiles(opts.SinceTag)
	if err != nil {
		return "", fmt.Errorf("failed to get changed files: %w", err)
	}

	if len(changedFiles) == 0 {
		if opts.Verbose {
			fmt.Fprintln(os.Stderr, "⚠️  No changes detected")
		}
		return "", fmt.Errorf("no changes detected since %s", opts.SinceTag)
	}

	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "   Found %d changed files\n", len(changedFiles))
	}

	// Get commits
	commits, err := git.GetCommits(opts.SinceTag)
	if err != nil {
		return "", fmt.Errorf("failed to get commits: %w", err)
	}

	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "   Found %d commits\n", len(commits))
	}

	// Extract code context
	if opts.Verbose {
		fmt.Fprintln(os.Stderr, "\n🔍 Extracting code context with promptext...")
	}

	result, err := aicontext.ExtractCodeContext(changedFiles)
	if err != nil {
		return "", fmt.Errorf("failed to extract context: %w", err)
	}

	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "   Extracted context: ~%d tokens from %d files\n",
			result.TokenCount, len(result.ProjectOutput.Files))
	}

	// Categorize commits
	categories := analyzer.CategorizeCommits(commits)

	// Generate AI prompt
	promptText := prompt.GenerateAIPrompt(opts.Version, opts.SinceTag, commits, categories, result)

	// If only prompt is requested, return it
	if opts.AIPromptOnly {
		if opts.Verbose {
			fmt.Fprintln(os.Stderr, "\n📝 Generated AI prompt (see stdout)")
		}
		return promptText, nil
	}

	// If AI enhancement is requested, call the AI provider
	if opts.UseAI && provider != nil {
		if opts.Verbose {
			fmt.Fprintf(os.Stderr, "\n🤖 Generating AI-enhanced changelog using %s...\n", provider.Name())
		}

		// Create AI request
		req := &ai.Request{
			Prompt:      promptText,
			Model:       "", // Will be set from config
			MaxTokens:   0,  // Will be set from config
			Temperature: 0,  // Will be set from config
		}

		// Call AI provider
		response, err := provider.Generate(ctx, req)
		if err != nil {
			return "", fmt.Errorf("failed to generate AI response: %w", err)
		}

		if opts.Verbose {
			fmt.Fprintf(os.Stderr, "   ✓ Generated %d tokens", response.TokensUsed)
			if response.CostEstimate > 0 {
				fmt.Fprintf(os.Stderr, " (estimated cost: $%.4f)", response.CostEstimate)
			}
			fmt.Fprintln(os.Stderr)
		}

		// Post-process the AI response to remove any extra headers
		content := stripAIHeaders(response.Content)

		return content, nil
	}

	// Otherwise, generate basic release notes
	if opts.Verbose {
		fmt.Fprintln(os.Stderr, "\n📝 Generating release notes...")
	}

	return generator.GenerateReleaseNotes(opts.Version, categories, result), nil
}

// stripAIHeaders removes common AI-generated headers from the response
func stripAIHeaders(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	skipNext := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip lines that look like AI headers
		if strings.HasPrefix(trimmed, "# Release Notes for") ||
			strings.HasPrefix(trimmed, "# Changelog for") ||
			strings.HasPrefix(trimmed, "Here are the") ||
			strings.HasPrefix(trimmed, "Here is the") {
			skipNext = true
			continue
		}

		// Skip empty lines immediately after headers
		if skipNext && trimmed == "" {
			skipNext = false
			continue
		}

		skipNext = false

		// Keep everything else
		if i > 0 || trimmed != "" {
			result = append(result, line)
		}
	}

	return strings.TrimSpace(strings.Join(result, "\n"))
}
