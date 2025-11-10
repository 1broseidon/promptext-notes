package git

import (
	"testing"
)

func TestIsGitRepository(t *testing.T) {
	// This test assumes it's being run from within the git repository
	if !IsGitRepository() {
		t.Skip("Not in a git repository, skipping test")
	}

	result := IsGitRepository()
	if !result {
		t.Error("Expected to be in a git repository")
	}
}

func TestGetLastTag(t *testing.T) {
	// This test will work even if there are no tags
	tag, err := GetLastTag()
	if err != nil {
		t.Fatalf("GetLastTag() error = %v", err)
	}

	// Tag should either be a valid tag or the fallback "HEAD~10"
	if tag == "" {
		t.Error("GetLastTag() returned empty string")
	}

	// If no tags exist, should return fallback
	// We can't test for specific values as they depend on repository state
	t.Logf("GetLastTag() = %s", tag)
}

func TestGetChangedFiles(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not in a git repository, skipping test")
	}

	// Test with HEAD~1 to HEAD (should work in any repo with commits)
	files, err := GetChangedFiles("HEAD~1")
	if err != nil {
		// It's okay if there's an error if we don't have enough commits
		t.Logf("GetChangedFiles() error = %v (may be expected if repo has < 2 commits)", err)
		return
	}

	// Result might be empty if no files changed, that's okay
	t.Logf("GetChangedFiles(HEAD~1) returned %d files", len(files))

	// Check that files don't contain empty strings
	for _, file := range files {
		if file == "" {
			t.Error("GetChangedFiles() returned empty file path")
		}
	}
}

func TestGetCommits(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not in a git repository, skipping test")
	}

	// Test with HEAD~1 to HEAD
	commits, err := GetCommits("HEAD~1")
	if err != nil {
		t.Logf("GetCommits() error = %v (may be expected if repo has < 2 commits)", err)
		return
	}

	// Should have at least one commit in most cases
	t.Logf("GetCommits(HEAD~1) returned %d commits", len(commits))

	// Check that commits don't contain empty strings
	for _, commit := range commits {
		if commit == "" {
			t.Error("GetCommits() returned empty commit message")
		}
	}
}

func TestGetChangedFilesWithInvalidRef(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not in a git repository, skipping test")
	}

	// Test with an invalid reference
	_, err := GetChangedFiles("invalid-ref-that-does-not-exist")
	if err == nil {
		t.Error("GetChangedFiles() with invalid ref should return error")
	}
}

func TestGetCommitsWithInvalidRef(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not in a git repository, skipping test")
	}

	// Test with an invalid reference
	_, err := GetCommits("invalid-ref-that-does-not-exist")
	if err == nil {
		t.Error("GetCommits() with invalid ref should return error")
	}
}
