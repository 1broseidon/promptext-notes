package generator

import (
	"fmt"
	"strings"
	"time"

	"github.com/1broseidon/promptext-notes/internal/analyzer"
	"github.com/1broseidon/promptext/pkg/promptext"
)

// GenerateReleaseNotes generates release notes in Keep a Changelog format.
func GenerateReleaseNotes(version string, categories analyzer.CommitCategories, result *promptext.Result) string {
	var notes strings.Builder

	// Determine version
	if version == "" {
		version = "Unreleased"
	}

	// Header
	notes.WriteString(fmt.Sprintf("## [%s] - %s\n\n",
		version, time.Now().Format("2006-01-02")))

	// Breaking Changes (highest priority)
	if len(categories.Breaking) > 0 {
		notes.WriteString("### ⚠️ Breaking Changes\n")
		for _, item := range categories.Breaking {
			notes.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
		}
		notes.WriteString("\n")
	}

	// Added (features)
	if len(categories.Features) > 0 {
		notes.WriteString("### Added\n")
		for _, item := range categories.Features {
			notes.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
		}
		notes.WriteString("\n")
	}

	// Fixed
	if len(categories.Fixes) > 0 {
		notes.WriteString("### Fixed\n")
		for _, item := range categories.Fixes {
			notes.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
		}
		notes.WriteString("\n")
	}

	// Changed
	if len(categories.Changes) > 0 {
		notes.WriteString("### Changed\n")
		for _, item := range categories.Changes {
			notes.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
		}
		notes.WriteString("\n")
	}

	// Documentation
	if len(categories.Docs) > 0 {
		notes.WriteString("### Documentation\n")
		for _, item := range categories.Docs {
			notes.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
		}
		notes.WriteString("\n")
	}

	// Statistics
	notes.WriteString("### Statistics\n")
	notes.WriteString(fmt.Sprintf("- **Files changed**: %d\n",
		len(result.ProjectOutput.Files)))
	notes.WriteString(fmt.Sprintf("- **Commits**: %d\n",
		categories.CountTotal()))
	notes.WriteString(fmt.Sprintf("- **Context analyzed**: ~%d tokens\n",
		result.TokenCount))
	notes.WriteString("\n")

	notes.WriteString("---\n\n")

	return notes.String()
}
