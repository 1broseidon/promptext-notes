package prompt

import (
	"strings"
	"testing"

	"github.com/1broseidon/promptext-notes/internal/analyzer"
	"github.com/1broseidon/promptext/pkg/promptext"
)

func TestGenerateAIPrompt(t *testing.T) {
	commits := []string{
		"feat: add new feature",
		"fix: resolve bug",
		"docs: update README",
	}

	categories := analyzer.CommitCategories{
		Features: []string{"add new feature"},
		Fixes:    []string{"resolve bug"},
		Docs:     []string{"update README"},
	}

	result := &promptext.Result{
		TokenCount:      5000,
		FormattedOutput: "# Project Code\n\nThis is the code context",
		ProjectOutput: &promptext.ProjectOutput{
			Files: []promptext.FileInfo{
				{Path: "main.go", Tokens: 3000},
				{Path: "README.md", Tokens: 2000},
			},
		},
	}

	tests := []struct {
		name      string
		version   string
		fromTag   string
		commits   []string
		wantParts []string
	}{
		{
			name:    "full AI prompt with version",
			version: "v1.0.0",
			fromTag: "v0.9.0",
			commits: commits,
			wantParts: []string{
				"# Release Notes Enhancement Request",
				"Please generate comprehensive release notes for version v1.0.0",
				"## Context",
				"**Version**: v1.0.0",
				"**Changes since**: v0.9.0",
				"**Commits analyzed**: 3",
				"**Files changed**: 2",
				"**Context extracted**: ~5000 tokens",
				"## 🎯 Executive Summary",
				"## Commit History",
				"feat: add new feature",
				"fix: resolve bug",
				"docs: update README",
				"## Changed Files Summary",
				"`main.go` (~3000 tokens)",
				"`README.md` (~2000 tokens)",
				"## Code Context (via promptext)",
				"This is the code context",
				"## Task",
				"Generate release notes in Keep a Changelog format",
				"## Critical Rules",
				"PRIMARY SOURCE",
				"USER VALUE ONLY",
				"## Example Format",
				"Generate ONLY the sections with content",
			},
		},
		{
			name:    "unreleased version",
			version: "",
			fromTag: "v0.9.0",
			commits: commits,
			wantParts: []string{
				"version Unreleased",
				"**Version**: Unreleased",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateAIPrompt(tt.version, tt.fromTag, tt.commits, categories, result, "", "")

			// Check that all expected parts are present
			for _, part := range tt.wantParts {
				if !strings.Contains(got, part) {
					t.Errorf("GenerateAIPrompt() missing expected part: %q", part)
				}
			}
		})
	}
}

func TestGenerateAIPromptStructure(t *testing.T) {
	commits := []string{"feat: test"}
	categories := analyzer.CommitCategories{
		Features: []string{"test"},
	}
	result := &promptext.Result{
		TokenCount:      1000,
		FormattedOutput: "code",
		ProjectOutput: &promptext.ProjectOutput{
			Files: []promptext.FileInfo{{Path: "test.go", Tokens: 1000}},
		},
	}

	prompt := GenerateAIPrompt("v1.0.0", "v0.9.0", commits, categories, result, "", "")

	// Verify all major sections are present in order
	sections := []string{
		"# Release Notes Enhancement Request",
		"## Context",
		"## 🎯 Executive Summary",
		"## 📊 Git Diff Summary",
		"## Code Context (via promptext)",
		"## Changed Files Summary",
		"## Commit History",
		"## Task",
		"## Critical Rules",
		"## Example Format",
	}

	lastIndex := -1
	for _, section := range sections {
		index := strings.Index(prompt, section)
		if index == -1 {
			t.Errorf("Missing section: %s", section)
			continue
		}
		if index <= lastIndex {
			t.Errorf("Section %s is out of order", section)
		}
		lastIndex = index
	}
}

func TestGenerateAIPromptCodeBlocks(t *testing.T) {
	commits := []string{"feat: test"}
	categories := analyzer.CommitCategories{Features: []string{"test"}}
	result := &promptext.Result{
		TokenCount:      1000,
		FormattedOutput: "code content",
		ProjectOutput: &promptext.ProjectOutput{
			Files: []promptext.FileInfo{{Path: "test.go", Tokens: 1000}},
		},
	}

	prompt := GenerateAIPrompt("v1.0.0", "v0.9.0", commits, categories, result, "", "")

	// Count code blocks (should have at least 3: commit history, code context, example)
	codeBlockCount := strings.Count(prompt, "```")
	if codeBlockCount < 6 { // Each code block has opening and closing ```
		t.Errorf("Expected at least 6 code block markers (3 blocks), got %d", codeBlockCount)
	}

	// Verify code context is included
	if !strings.Contains(prompt, "code content") {
		t.Error("Code context should be included in the prompt")
	}
}

func TestGenerateAIPromptEmptyCommits(t *testing.T) {
	commits := []string{}
	categories := analyzer.CommitCategories{}
	result := &promptext.Result{
		TokenCount:      0,
		FormattedOutput: "",
		ProjectOutput: &promptext.ProjectOutput{
			Files: []promptext.FileInfo{},
		},
	}

	prompt := GenerateAIPrompt("v1.0.0", "v0.9.0", commits, categories, result, "", "")

	// Should still have the structure
	if !strings.Contains(prompt, "# Release Notes Enhancement Request") {
		t.Error("Should have header even with empty commits")
	}

	// Should show 0 commits
	if !strings.Contains(prompt, "**Commits analyzed**: 0") {
		t.Error("Should show 0 commits analyzed")
	}

	// Should show 0 files
	if !strings.Contains(prompt, "**Files changed**: 0") {
		t.Error("Should show 0 files changed")
	}
}
