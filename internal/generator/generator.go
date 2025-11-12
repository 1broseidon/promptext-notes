package generator

import (
	"fmt"
	"strings"
	"time"

	"github.com/1broseidon/promptext-notes/internal/analyzer"
	"github.com/1broseidon/promptext-notes/internal/config"
	"github.com/1broseidon/promptext/pkg/promptext"
)

// OutputConfig holds configuration for output generation
type OutputConfig struct {
	Sections []string // List of sections to include (breaking, added, changed, fixed, docs, etc.)
}

// GenerateReleaseNotes generates release notes with configurable format and sections.
// If cfg is nil, uses default Keep a Changelog format with all sections.
func GenerateReleaseNotes(version string, categories analyzer.CommitCategories, result *promptext.Result, cfg *config.Config) string {
	var notes strings.Builder

	// Determine version
	if version == "" {
		version = "Unreleased"
	}

	// Header
	notes.WriteString(fmt.Sprintf("## [%s] - %s\n\n",
		version, time.Now().Format("2006-01-02")))

	// Determine which sections to include
	sections := []string{"breaking", "added", "fixed", "changed", "docs"}
	if cfg != nil && len(cfg.Output.Sections) > 0 {
		sections = cfg.Output.Sections
	}

	// Generate sections based on config
	for _, section := range sections {
		switch strings.ToLower(section) {
		case "breaking":
			if len(categories.Breaking) > 0 {
				notes.WriteString("### ⚠️ Breaking Changes\n")
				for _, item := range categories.Breaking {
					notes.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
				}
				notes.WriteString("\n")
			}

		case "added":
			if len(categories.Features) > 0 {
				notes.WriteString("### Added\n")
				for _, item := range categories.Features {
					notes.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
				}
				notes.WriteString("\n")
			}

		case "fixed":
			if len(categories.Fixes) > 0 {
				notes.WriteString("### Fixed\n")
				for _, item := range categories.Fixes {
					notes.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
				}
				notes.WriteString("\n")
			}

		case "changed":
			if len(categories.Changes) > 0 {
				notes.WriteString("### Changed\n")
				for _, item := range categories.Changes {
					notes.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
				}
				notes.WriteString("\n")
			}

		case "docs", "documentation":
			if len(categories.Docs) > 0 {
				notes.WriteString("### Documentation\n")
				for _, item := range categories.Docs {
					notes.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
				}
				notes.WriteString("\n")
			}

			// Future: add support for deprecated, removed, security sections
		}
	}

	// Statistics (always include unless explicitly disabled)
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
