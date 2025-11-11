package context

import (
	"path/filepath"

	"github.com/1broseidon/promptext/pkg/promptext"
)

// ExtractCodeContext extracts code context from changed files using promptext.
// It focuses on relevant file types (.go, .md, .yml, .yaml) and applies a token budget.
func ExtractCodeContext(changedFiles []string) (*promptext.Result, error) {
	// Focus on code and documentation files
	relevantExts := []string{".go", ".md", ".yml", ".yaml"}

	// Filter changed files by extension and exclude meta files
	var relevantFiles []string
	for _, file := range changedFiles {
		// Skip CHANGELOG and README - these are meta-documentation that pollute context
		base := filepath.Base(file)
		if base == "CHANGELOG.md" || base == "README.md" {
			continue
		}

		ext := filepath.Ext(file)
		for _, relevantExt := range relevantExts {
			if ext == relevantExt {
				relevantFiles = append(relevantFiles, file)
				break
			}
		}
	}

	// Determine token budget based on number of relevant files
	tokenBudget := 8000
	if len(relevantFiles) == 0 {
		// If no relevant files, get project summary with smaller budget
		tokenBudget = 4000
	}

	// Extract context with token budget
	result, err := promptext.Extract(".",
		promptext.WithExtensions(relevantExts...),
		promptext.WithTokenBudget(tokenBudget),
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
