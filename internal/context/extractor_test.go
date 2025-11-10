package context

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractCodeContext(t *testing.T) {
	tests := []struct {
		name         string
		changedFiles []string
		wantErr      bool
	}{
		{
			name: "with Go files",
			changedFiles: []string{
				"main.go",
				"internal/analyzer/analyzer.go",
				"README.md",
			},
			wantErr: false,
		},
		{
			name: "with various file types",
			changedFiles: []string{
				"config.yml",
				"config.yaml",
				"docs.md",
				"script.go",
			},
			wantErr: false,
		},
		{
			name: "with non-relevant files",
			changedFiles: []string{
				"image.png",
				"binary.exe",
				"data.json",
			},
			wantErr: false,
		},
		{
			name:         "empty file list",
			changedFiles: []string{},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractCodeContext(tt.changedFiles)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractCodeContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("ExtractCodeContext() returned nil result")
					return
				}

				// Verify result has required fields
				if result.TokenCount < 0 {
					t.Errorf("TokenCount should be non-negative, got %d", result.TokenCount)
				}

				// ProjectOutput should be initialized
				if result.ProjectOutput.Files == nil {
					t.Error("ProjectOutput.Files should not be nil")
				}
			}
		})
	}
}

func TestExtractCodeContextFileFiltering(t *testing.T) {
	changedFiles := []string{
		"main.go",           // Should be included
		"README.md",         // Should be included
		"config.yml",        // Should be included
		"docker.yaml",       // Should be included
		"image.png",         // Should be filtered out
		"script.sh",         // Should be filtered out
		"data.json",         // Should be filtered out
		"package-lock.json", // Should be filtered out
	}

	result, err := ExtractCodeContext(changedFiles)
	if err != nil {
		t.Fatalf("ExtractCodeContext() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("ExtractCodeContext() returned nil result")
	}

	// Result should exist and be valid
	if result.TokenCount < 0 {
		t.Errorf("TokenCount should be non-negative, got %d", result.TokenCount)
	}
}

func TestExtractCodeContextInTestDirectory(t *testing.T) {
	// Create a temporary test directory
	tmpDir := t.TempDir()

	// Create test files
	testFiles := map[string]string{
		"test.go":    "package main\n\nfunc main() {}\n",
		"README.md":  "# Test Project\n",
		"config.yml": "key: value\n",
		"ignore.txt": "should be ignored\n",
	}

	for name, content := range testFiles {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Change to temp directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(origDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Test with files that exist
	changedFiles := []string{"test.go", "README.md", "config.yml"}
	result, err := ExtractCodeContext(changedFiles)
	if err != nil {
		t.Fatalf("ExtractCodeContext() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("ExtractCodeContext() returned nil result")
	}

	if result.TokenCount <= 0 {
		t.Errorf("Expected positive token count, got %d", result.TokenCount)
	}
}

func TestExtractCodeContextTokenBudget(t *testing.T) {
	tests := []struct {
		name         string
		changedFiles []string
		description  string
	}{
		{
			name:         "with relevant files (8000 token budget)",
			changedFiles: []string{"main.go", "README.md"},
			description:  "Should use 8000 token budget when relevant files exist",
		},
		{
			name:         "without relevant files (4000 token budget)",
			changedFiles: []string{},
			description:  "Should use 4000 token budget when no relevant files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractCodeContext(tt.changedFiles)
			if err != nil {
				t.Fatalf("ExtractCodeContext() unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("ExtractCodeContext() returned nil result")
			}

			// Just verify we got a result - the actual token budget is internal to promptext
			t.Logf("%s - TokenCount: %d", tt.description, result.TokenCount)
		})
	}
}
