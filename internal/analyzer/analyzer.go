package analyzer

import (
	"strings"
)

// CommitCategories holds categorized commit messages.
type CommitCategories struct {
	Features []string
	Fixes    []string
	Docs     []string
	Chores   []string
	Changes  []string
	Breaking []string
}

// CategorizeCommits categorizes commit messages based on conventional commit format.
func CategorizeCommits(commits []string) CommitCategories {
	cats := CommitCategories{
		Features: []string{},
		Fixes:    []string{},
		Docs:     []string{},
		Chores:   []string{},
		Changes:  []string{},
		Breaking: []string{},
	}

	for _, commit := range commits {
		lower := strings.ToLower(commit)

		// Check for breaking changes first
		if strings.Contains(lower, "breaking") || strings.Contains(commit, "BREAKING CHANGE") {
			cats.Breaking = append(cats.Breaking, commit)
			continue
		}

		// Categorize by conventional commit prefix
		if strings.HasPrefix(lower, "feat:") || strings.HasPrefix(lower, "feature:") {
			message := extractMessage(commit, "feat:", "feature:")
			cats.Features = append(cats.Features, message)
		} else if strings.HasPrefix(lower, "fix:") {
			message := extractMessage(commit, "fix:")
			cats.Fixes = append(cats.Fixes, message)
		} else if strings.HasPrefix(lower, "docs:") {
			message := extractMessage(commit, "docs:")
			cats.Docs = append(cats.Docs, message)
		} else if strings.HasPrefix(lower, "chore:") {
			message := extractMessage(commit, "chore:")
			cats.Chores = append(cats.Chores, message)
		} else if strings.HasPrefix(lower, "refactor:") {
			message := extractMessage(commit, "refactor:")
			cats.Changes = append(cats.Changes, message)
		} else if strings.HasPrefix(lower, "test:") {
			// Skip test commits or include in chores
			message := extractMessage(commit, "test:")
			cats.Chores = append(cats.Chores, message)
		} else {
			// Uncategorized commits go to Changes
			cats.Changes = append(cats.Changes, commit)
		}
	}

	return cats
}

// extractMessage removes the prefix from a commit message and trims whitespace.
func extractMessage(commit string, prefixes ...string) string {
	lower := strings.ToLower(commit)
	for _, prefix := range prefixes {
		if strings.HasPrefix(lower, prefix) {
			// Find the prefix in the original commit (preserving case)
			idx := len(prefix)
			if idx < len(commit) {
				return strings.TrimSpace(commit[idx:])
			}
		}
	}
	return strings.TrimSpace(commit)
}

// CountTotal returns the total number of commits across all categories.
func (c *CommitCategories) CountTotal() int {
	return len(c.Features) + len(c.Fixes) + len(c.Docs) +
		len(c.Chores) + len(c.Changes) + len(c.Breaking)
}
