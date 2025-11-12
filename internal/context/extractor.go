package context

import (
	"path/filepath"

	"github.com/1broseidon/promptext/pkg/promptext"
)

// ExtractCodeContext extracts code context from changed files using promptext.
// It focuses on relevant file types (.go, .md, .yml, .yaml) and applies a token budget.
// The excludePatterns parameter allows filtering out specific files by filename.
func ExtractCodeContext(changedFiles []string, excludePatterns []string) (*promptext.Result, error) {
	// Focus on code and documentation files
	relevantExts := []string{".go", ".md", ".yml", ".yaml"}

	// Default exclusions if none provided
	if len(excludePatterns) == 0 {
		excludePatterns = []string{"CHANGELOG.md", "README.md"}
	}

	// Filter changed files by extension and exclude patterns
	var relevantFiles []string
	for _, file := range changedFiles {
		// Check if file matches any exclude pattern (by basename)
		excluded := false
		base := filepath.Base(file)
		for _, pattern := range excludePatterns {
			if base == pattern {
				excluded = true
				break
			}
		}
		if excluded {
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
