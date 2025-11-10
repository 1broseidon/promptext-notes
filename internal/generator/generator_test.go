package generator

import (
	"strings"
	"testing"
	"time"

	"github.com/1broseidon/promptext-notes/internal/analyzer"
	"github.com/1broseidon/promptext/pkg/promptext"
)

func TestGenerateReleaseNotes(t *testing.T) {
	// Create test data
	categories := analyzer.CommitCategories{
		Features: []string{"add new feature", "add another feature"},
		Fixes:    []string{"fix critical bug"},
		Docs:     []string{"update README"},
		Changes:  []string{"refactor code"},
		Breaking: []string{"remove deprecated API"},
	}

	result := &promptext.Result{
		TokenCount: 5000,
		ProjectOutput: &promptext.ProjectOutput{
			Files: []promptext.FileInfo{
				{Path: "main.go", Tokens: 1000},
				{Path: "README.md", Tokens: 500},
			},
		},
	}

	tests := []struct {
		name       string
		version    string
		categories analyzer.CommitCategories
		result     *promptext.Result
		wantParts  []string
	}{
		{
			name:       "full release notes with version",
			version:    "v1.0.0",
			categories: categories,
			result:     result,
			wantParts: []string{
				"## [v1.0.0]",
				"### ⚠️ Breaking Changes",
				"remove deprecated API",
				"### Added",
				"add new feature",
				"add another feature",
				"### Fixed",
				"fix critical bug",
				"### Changed",
				"refactor code",
				"### Documentation",
				"update README",
				"### Statistics",
				"**Files changed**: 2",
				"**Commits**: 6",
				"**Context analyzed**: ~5000 tokens",
			},
		},
		{
			name:       "unreleased version",
			version:    "",
			categories: categories,
			result:     result,
			wantParts: []string{
				"## [Unreleased]",
			},
		},
		{
			name:    "empty categories",
			version: "v0.1.0",
			categories: analyzer.CommitCategories{
				Features: []string{},
				Fixes:    []string{},
				Docs:     []string{},
				Changes:  []string{},
				Breaking: []string{},
			},
			result: &promptext.Result{
				TokenCount: 100,
				ProjectOutput: &promptext.ProjectOutput{
					Files: []promptext.FileInfo{},
				},
			},
			wantParts: []string{
				"## [v0.1.0]",
				"### Statistics",
				"**Files changed**: 0",
				"**Commits**: 0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateReleaseNotes(tt.version, tt.categories, tt.result)

			// Check that all expected parts are present
			for _, part := range tt.wantParts {
				if !strings.Contains(got, part) {
					t.Errorf("GenerateReleaseNotes() missing expected part: %q\nGot:\n%s", part, got)
				}
			}

			// Check date format is valid
			today := time.Now().Format("2006-01-02")
			if !strings.Contains(got, today) {
				t.Errorf("GenerateReleaseNotes() should contain today's date %s", today)
			}

			// Should end with separator
			if !strings.Contains(got, "---") {
				t.Error("GenerateReleaseNotes() should end with separator")
			}
		})
	}
}

func TestGenerateReleaseNotesFormat(t *testing.T) {
	categories := analyzer.CommitCategories{
		Features: []string{"feature 1"},
	}
	result := &promptext.Result{
		TokenCount: 1000,
		ProjectOutput: &promptext.ProjectOutput{
			Files: []promptext.FileInfo{
				{Path: "test.go", Tokens: 1000},
			},
		},
	}

	notes := GenerateReleaseNotes("v1.0.0", categories, result)

	// Verify markdown structure
	if !strings.HasPrefix(notes, "##") {
		t.Error("Release notes should start with ## header")
	}

	// Verify sections have proper markdown headers
	expectedHeaders := []string{"### Added", "### Statistics"}
	for _, header := range expectedHeaders {
		if !strings.Contains(notes, header) {
			t.Errorf("Missing expected header: %s", header)
		}
	}

	// Verify bullet points
	if !strings.Contains(notes, "- ") {
		t.Error("Release notes should contain bullet points")
	}
}

func TestGenerateReleaseNotesOnlyBreaking(t *testing.T) {
	categories := analyzer.CommitCategories{
		Breaking: []string{"major breaking change"},
	}
	result := &promptext.Result{
		TokenCount: 100,
		ProjectOutput: &promptext.ProjectOutput{
			Files: []promptext.FileInfo{},
		},
	}

	notes := GenerateReleaseNotes("v2.0.0", categories, result)

	// Should have breaking changes section
	if !strings.Contains(notes, "### ⚠️ Breaking Changes") {
		t.Error("Should contain breaking changes section")
	}

	// Should not have other sections (except Statistics)
	unwantedSections := []string{"### Added", "### Fixed", "### Changed", "### Documentation"}
	for _, section := range unwantedSections {
		if strings.Contains(notes, section) {
			t.Errorf("Should not contain section: %s", section)
		}
	}
}
