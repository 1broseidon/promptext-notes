package context

import (
	"path/filepath"

	"github.com/1broseidon/promptext/pkg/promptext"
	"github.com/bmatcuk/doublestar/v4"
)

// ExtractCodeContext extracts code context from changed files using promptext.
// It focuses on relevant file types (.go, .md, .yml, .yaml) and applies a token budget.
// The excludePatterns parameter allows filtering out specific files using glob patterns.
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
		// Check if file matches any exclude pattern (supports globs)
		excluded := false
		for _, pattern := range excludePatterns {
			// Try glob match first
			matched, err := doublestar.Match(pattern, file)
			if err == nil && matched {
				excluded = true
				break
			}
			// Fallback to basename exact match for backwards compatibility
			if filepath.Base(file) == pattern {
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

	// Extract context with token budget and exclusions
	// NOTE: We extract from the full repo (not just changed files) because:
	// 1. AI needs context from unchanged files (imports, interfaces, related code)
	// 2. Promptext doesn't support extracting from specific file list
	// The relevantFiles list is used to determine token budget and can be logged
	result, err := promptext.Extract(".",
		promptext.WithExtensions(relevantExts...),
		promptext.WithTokenBudget(tokenBudget),
		promptext.WithExcludes(excludePatterns...), // Apply exclude patterns to promptext
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
