package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetLastTag retrieves the most recent git tag.
// Returns "HEAD~10" as fallback if no tags are found.
func GetLastTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		// No tags found, return fallback
		return "HEAD~10", nil
	}
	return strings.TrimSpace(string(output)), nil
}

// GetChangedFiles returns a list of files changed between the given tag/commit and HEAD.
func GetChangedFiles(since string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", since+"..HEAD")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var files []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			files = append(files, trimmed)
		}
	}
	return files, nil
}

// GetCommits returns a list of commit messages between the given tag/commit and HEAD.
func GetCommits(since string) ([]string, error) {
	cmd := exec.Command("git", "log", since+"..HEAD", "--pretty=format:%s")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var commits []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			commits = append(commits, trimmed)
		}
	}
	return commits, nil
}

// IsGitRepository checks if the current directory is a git repository.
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// GetDiffStats returns git diff --stat output between the given tag/commit and HEAD.
func GetDiffStats(since string) (string, error) {
	cmd := exec.Command("git", "diff", "--stat", since+"..HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff stats: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetDiff returns git diff output between the given tag/commit and HEAD.
// Use --unified=3 for standard context.
func GetDiff(since string) (string, error) {
	cmd := exec.Command("git", "diff", "--unified=3", since+"..HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
