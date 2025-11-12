package workflow

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/1broseidon/promptext-notes/internal/ai"
	"github.com/1broseidon/promptext-notes/internal/analyzer"
	"github.com/1broseidon/promptext-notes/internal/config"
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
	ExcludeFiles []string // Files to exclude from AI context (e.g., CHANGELOG.md)
}

// gitData holds git-related data for release notes
type gitData struct {
	changedFiles []string
	commits      []string
	diffStats    string
	diff         string
}

// fetchGitData retrieves all git-related information needed for release notes
func fetchGitData(sinceTag string, verbose bool) (*gitData, error) {
	data := &gitData{}
	var err error

	// Get changed files
	data.changedFiles, err = git.GetChangedFiles(sinceTag)
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	if len(data.changedFiles) == 0 {
		if verbose {
			fmt.Fprintln(os.Stderr, "⚠️  No changes detected")
		}
		return nil, fmt.Errorf("no changes detected since %s", sinceTag)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "   Found %d changed files\n", len(data.changedFiles))
	}

	// Get commits
	data.commits, err = git.GetCommits(sinceTag)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "   Found %d commits\n", len(data.commits))
	}

	// Get git diff stats (non-fatal)
	data.diffStats, err = git.GetDiffStats(sinceTag)
	if err != nil && verbose {
		fmt.Fprintf(os.Stderr, "   Warning: could not get diff stats: %v\n", err)
	}

	// Get git diff (non-fatal)
	data.diff, err = git.GetDiff(sinceTag)
	if err != nil && verbose {
		fmt.Fprintf(os.Stderr, "   Warning: could not get diff: %v\n", err)
	}

	return data, nil
}

// GenerateReleaseNotes orchestrates the full release notes generation process
func GenerateReleaseNotes(ctx context.Context, opts GenerateOptions, provider ai.Provider, cfg *config.Config) (string, error) {
	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "📊 Analyzing changes since %s...\n", opts.SinceTag)
	}

	// Fetch all git data
	gitData, err := fetchGitData(opts.SinceTag, opts.Verbose)
	if err != nil {
		return "", err
	}

	// Extract code context
	if opts.Verbose {
		fmt.Fprintln(os.Stderr, "\n🔍 Extracting code context with promptext...")
	}

	result, err := aicontext.ExtractCodeContext(gitData.changedFiles, opts.ExcludeFiles)
	if err != nil {
		return "", fmt.Errorf("failed to extract context: %w", err)
	}

	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "   Extracted context: ~%d tokens from %d files\n",
			result.TokenCount, len(result.ProjectOutput.Files))
	}

	// Categorize commits
	categories := analyzer.CategorizeCommits(gitData.commits)

	// Generate AI prompt
	promptText := prompt.GenerateAIPrompt(opts.Version, opts.SinceTag, gitData.commits,
		categories, result, gitData.diffStats, gitData.diff)

	// If only prompt is requested, return it
	if opts.AIPromptOnly {
		if opts.Verbose {
			fmt.Fprintln(os.Stderr, "\n📝 Generated AI prompt (see stdout)")
		}
		return promptText, nil
	}

	// If AI enhancement is requested, call the AI provider
	if opts.UseAI && provider != nil {
		content, err := generateAIContent(ctx, provider, promptText, opts.Verbose)
		if err != nil {
			return "", err
		}

		// Stage 2: Polish if enabled
		if cfg != nil && cfg.AI.Polish.Enabled {
			if opts.Verbose {
				polishProvider := cfg.GetPolishProvider()
				polishModel := cfg.GetPolishModel()
				fmt.Fprintf(os.Stderr, "\n✨ Polishing changelog with %s (%s)...\n", polishProvider, polishModel)
			}

			polishedContent, err := PolishChangelog(ctx, content, gitData.diff, cfg)
			if err != nil {
				return "", fmt.Errorf("failed to polish changelog: %w", err)
			}

			if opts.Verbose {
				fmt.Fprintln(os.Stderr, "   ✓ Polish complete")
			}

			return polishedContent, nil
		}

		return content, nil
	}

	// Otherwise, generate basic release notes
	if opts.Verbose {
		fmt.Fprintln(os.Stderr, "\n📝 Generating release notes...")
	}

	return generator.GenerateReleaseNotes(opts.Version, categories, result), nil
}

// generateAIContent calls the AI provider and processes the response
func generateAIContent(ctx context.Context, provider ai.Provider, promptText string, verbose bool) (string, error) {
	if verbose {
		fmt.Fprintf(os.Stderr, "\n🤖 Generating AI-enhanced changelog using %s...\n", provider.Name())
	}

	// Create AI request using provider's configured defaults
	req := provider.NewRequest(promptText)

	// Call AI provider (stage 1: discovery)
	response, err := provider.Generate(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to generate AI response: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "   ✓ Generated %d tokens", response.TokensUsed)
		if response.CostEstimate > 0 {
			fmt.Fprintf(os.Stderr, " (estimated cost: $%.4f)", response.CostEstimate)
		}
		fmt.Fprintln(os.Stderr)
	}

	// Post-process the AI response to remove any extra headers
	return stripAIHeaders(response.Content), nil
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
