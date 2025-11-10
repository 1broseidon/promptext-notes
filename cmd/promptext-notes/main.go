package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/1broseidon/promptext-notes/internal/analyzer"
	"github.com/1broseidon/promptext-notes/internal/context"
	"github.com/1broseidon/promptext-notes/internal/generator"
	"github.com/1broseidon/promptext-notes/internal/git"
	"github.com/1broseidon/promptext-notes/internal/prompt"
)

func main() {
	// Parse flags
	version := flag.String("version", "", "Version to generate notes for (e.g., v0.7.4)")
	sinceTag := flag.String("since", "", "Generate notes since this tag (auto-detects if empty)")
	output := flag.String("output", "", "Output file (prints to stdout if empty)")
	aiPrompt := flag.Bool("ai-prompt", false, "Generate prompt for AI to enhance release notes")
	flag.Parse()

	// Check if we're in a git repository
	if !git.IsGitRepository() {
		log.Fatal("Error: Not a git repository. Please run this command from within a git repository.")
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

	fmt.Fprintf(os.Stderr, "📊 Analyzing changes since %s...\n", fromTag)

	// Get changed files
	changedFiles, err := git.GetChangedFiles(fromTag)
	if err != nil {
		log.Fatalf("Failed to get changed files: %v", err)
	}

	if len(changedFiles) == 0 {
		fmt.Fprintln(os.Stderr, "⚠️  No changes detected")
		return
	}

	fmt.Fprintf(os.Stderr, "   Found %d changed files\n", len(changedFiles))

	// Get commits
	commits, err := git.GetCommits(fromTag)
	if err != nil {
		log.Fatalf("Failed to get commits: %v", err)
	}

	fmt.Fprintf(os.Stderr, "   Found %d commits\n", len(commits))

	// Extract code context
	fmt.Fprintln(os.Stderr, "\n🔍 Extracting code context with promptext...")
	result, err := context.ExtractCodeContext(changedFiles)
	if err != nil {
		log.Fatalf("Failed to extract context: %v", err)
	}

	fmt.Fprintf(os.Stderr, "   Extracted context: ~%d tokens from %d files\n",
		result.TokenCount, len(result.ProjectOutput.Files))

	// Categorize commits
	categories := analyzer.CategorizeCommits(commits)

	// Generate output
	var outputText string
	if *aiPrompt {
		fmt.Fprintln(os.Stderr, "\n📝 Generating AI prompt...")
		outputText = prompt.GenerateAIPrompt(*version, fromTag, commits, categories, result)
	} else {
		fmt.Fprintln(os.Stderr, "\n📝 Generating release notes...")
		outputText = generator.GenerateReleaseNotes(*version, categories, result)
	}

	// Write output
	if *output != "" {
		if err := os.WriteFile(*output, []byte(outputText), 0644); err != nil {
			log.Fatalf("Failed to write output: %v", err)
		}
		fmt.Fprintf(os.Stderr, "✅ Written to %s\n", *output)
	} else {
		fmt.Print(outputText)
	}
}
