# Promptext-Notes: Release Notes Generator - Complete Specification

**Version**: 1.0
**Last Updated**: 2025-11-10
**Status**: Ready for Implementation

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Core Purpose](#core-purpose)
3. [Technical Stack](#technical-stack)
4. [Dependencies](#dependencies)
5. [Architecture](#architecture)
6. [Feature Specifications](#feature-specifications)
7. [Git Integration](#git-integration)
8. [Promptext Integration](#promptext-integration)
9. [CLI Interface](#cli-interface)
10. [Output Formats](#output-formats)
11. [Implementation Guide](#implementation-guide)
12. [Testing Strategy](#testing-strategy)
13. [Installation & Distribution](#installation--distribution)
14. [Usage Examples](#usage-examples)
15. [Future Enhancements](#future-enhancements)

---

## Project Overview

### What is Promptext-Notes?

A Go-based CLI tool that generates intelligent, context-aware release notes by combining git history analysis with code context extraction using the promptext library.

### Key Differentiator

Unlike traditional changelog generators that only parse commit messages, promptext-notes:
- Extracts actual code context from changed files
- Uses token-aware analysis to stay within LLM limits
- Generates AI-ready prompts with full code context
- Understands conventional commits
- Produces Keep a Changelog formatted output

### Target Users

1. **Solo Developers**: Automate changelog creation for releases
2. **Teams**: Standardize release notes across projects
3. **AI-Assisted Workflows**: Generate context-rich prompts for LLMs to write polished notes
4. **CI/CD Pipelines**: Integrate into release automation

---

## Core Purpose

### Problem Statement

Writing release notes is tedious and error-prone:
- Manually reviewing git history is time-consuming
- Easy to miss important changes
- Hard to understand impact without code context
- Inconsistent formatting across releases
- LLMs need full context to write good release notes

### Solution

Promptext-notes automates release note generation by:
1. Analyzing git commits since last tag
2. Extracting code context with promptext (token-aware)
3. Categorizing changes by type (feat, fix, docs, breaking)
4. Generating structured markdown in Keep a Changelog format
5. Optionally creating AI prompts with full context for polished notes

---

## Technical Stack

### Language & Version
- **Go**: 1.22 or higher
- **Reason**: Cross-platform, single binary, excellent git/CLI tooling

### Build System
- Standard `go build`
- No external build tools required

### Runtime Dependencies
- Git must be installed and accessible in PATH
- No other runtime dependencies

---

## Dependencies

### Primary Dependency: Promptext Library

**Package**: `github.com/1broseidon/promptext`

**Purpose**: Token-aware code context extraction

**Key Functions Used**:
```go
import "github.com/1broseidon/promptext/pkg/promptext"

// Extract code context
result, err := promptext.Extract(directory string, opts ...Option)

// Options
promptext.WithExtensions(exts ...string)
promptext.WithTokenBudget(tokens int)
promptext.WithFormat(format Format)
```

**Result Structure**:
```go
type Result struct {
    FormattedOutput string       // Formatted code context
    TokenCount      int           // Estimated token count
    ProjectOutput   ProjectOutput // Structured data
}

type ProjectOutput struct {
    Files    []FileInfo
    Metadata ProjectMetadata
}

type FileInfo struct {
    Path    string
    Content string
    Tokens  int
}
```

### Standard Library
- `flag` - CLI argument parsing
- `os` - File I/O
- `os/exec` - Git command execution
- `fmt` - String formatting
- `log` - Error logging
- `strings` - String manipulation
- `time` - Date formatting
- `path/filepath` - Path handling

### Go Module Setup

```go
// go.mod
module github.com/yourusername/promptext-notes

go 1.22

require github.com/1broseidon/promptext v0.7.3
```

---

## Architecture

### High-Level Flow

```
Input (CLI flags)
    ↓
Git Analysis
    ├─ Get last tag (or use --since)
    ├─ Get changed files (git diff)
    └─ Get commit messages (git log)
    ↓
Context Extraction
    ├─ Filter relevant files (.go, .md, .yml)
    └─ Use promptext.Extract() with 8K budget
    ↓
Processing
    ├─ Categorize commits (feat/fix/docs/chore/breaking)
    └─ Build changelog structure
    ↓
Output Generation
    ├─ Mode 1: Basic release notes (markdown)
    └─ Mode 2: AI prompt (with full context)
    ↓
Output (stdout or file)
```

### Module Structure

```
promptext-notes/
├── main.go                 # Entry point, CLI setup
├── git.go                  # Git operations
├── analyzer.go             # Commit categorization
├── generator.go            # Release notes generation
├── ai_prompt.go            # AI prompt generation
├── go.mod                  # Dependencies
├── go.sum                  # Dependency checksums
├── README.md               # Documentation
├── LICENSE                 # License file
└── .gitignore              # Git ignore patterns
```

### Function Responsibilities

**main.go**:
- Parse CLI flags
- Orchestrate workflow
- Handle errors and logging
- Write output

**git.go**:
- `getLastTag() string` - Get most recent git tag
- `getChangedFiles(since string) ([]string, error)` - Get changed files
- `getCommits(since string) ([]string, error)` - Get commit messages

**analyzer.go**:
- `categorizeCommits(commits []string) CommitCategories` - Sort by type
- `detectBreakingChanges(commits []string) []string` - Find breaking changes

**generator.go**:
- `generateReleaseNotes(version string, data AnalysisData) string` - Create changelog

**ai_prompt.go**:
- `generateAIPrompt(version string, data AnalysisData) string` - Create AI prompt

---

## Feature Specifications

### Feature 1: Git History Analysis

**Input**:
- Version tag (e.g., `v0.7.4`)
- Since tag (optional, auto-detects if empty)

**Process**:
```go
// Get last tag if not specified
func getLastTag() string {
    cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
    output, err := cmd.Output()
    if err != nil {
        return "HEAD~10" // Fallback to last 10 commits
    }
    return strings.TrimSpace(string(output))
}

// Get changed files between tags
func getChangedFiles(since string) ([]string, error) {
    cmd := exec.Command("git", "diff", "--name-only", since+"..HEAD")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    lines := strings.Split(strings.TrimSpace(string(output)), "\n")
    var files []string
    for _, line := range lines {
        if line != "" {
            files = append(files, line)
        }
    }
    return files, nil
}

// Get commits between tags
func getCommits(since string) ([]string, error) {
    cmd := exec.Command("git", "log", since+"..HEAD", "--pretty=format:%s")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    lines := strings.Split(strings.TrimSpace(string(output)), "\n")
    var commits []string
    for _, line := range lines {
        if line != "" {
            commits = append(commits, line)
        }
    }
    return commits, nil
}
```

**Output**: Lists of changed files and commit messages

**Error Handling**:
- No git repository → Fatal error
- No tags found → Fallback to HEAD~10
- Empty result → Warning and exit

---

### Feature 2: Code Context Extraction

**Purpose**: Extract relevant code changes with token awareness

**Implementation**:
```go
func extractCodeContext(changedFiles []string) (*promptext.Result, error) {
    // Focus on code and documentation
    relevantExts := []string{".go", ".md", ".yml", ".yaml"}

    // Filter changed files by extension
    var relevantFiles []string
    for _, file := range changedFiles {
        ext := filepath.Ext(file)
        for _, relevantExt := range relevantExts {
            if ext == relevantExt {
                relevantFiles = append(relevantFiles, file)
                break
            }
        }
    }

    // If no relevant files, get project summary
    if len(relevantFiles) == 0 {
        return promptext.Extract(".",
            promptext.WithExtensions(relevantExts...),
            promptext.WithTokenBudget(4000),
        )
    }

    // Extract context with 8K token budget
    return promptext.Extract(".",
        promptext.WithExtensions(relevantExts...),
        promptext.WithTokenBudget(8000),
    )
}
```

**Token Budget**: 8,000 tokens
- Enough for meaningful context
- Fits within most LLM prompts with room for instructions

**File Filtering**: `.go`, `.md`, `.yml`, `.yaml`
- Skip binary files, lock files, build artifacts

---

### Feature 3: Commit Categorization

**Conventional Commit Types**:
- `feat:` → Added section
- `fix:` → Fixed section
- `docs:` → Documentation section
- `chore:` → Chores/maintenance section
- `refactor:` → Changed section
- `test:` → Testing section
- `breaking:` or `BREAKING CHANGE:` → Breaking Changes section

**Implementation**:
```go
type CommitCategories struct {
    Features []string
    Fixes    []string
    Docs     []string
    Chores   []string
    Changes  []string
    Breaking []string
}

func categorizeCommits(commits []string) CommitCategories {
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

        if strings.HasPrefix(lower, "feat:") || strings.HasPrefix(lower, "feature:") {
            cats.Features = append(cats.Features,
                strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(commit, "feat:"), "feature:")))
        } else if strings.HasPrefix(lower, "fix:") {
            cats.Fixes = append(cats.Fixes,
                strings.TrimSpace(strings.TrimPrefix(commit, "fix:")))
        } else if strings.HasPrefix(lower, "docs:") {
            cats.Docs = append(cats.Docs,
                strings.TrimSpace(strings.TrimPrefix(commit, "docs:")))
        } else if strings.HasPrefix(lower, "chore:") {
            cats.Chores = append(cats.Chores,
                strings.TrimSpace(strings.TrimPrefix(commit, "chore:")))
        } else if strings.Contains(lower, "breaking") {
            cats.Breaking = append(cats.Breaking, commit)
        } else {
            cats.Changes = append(cats.Changes, commit)
        }
    }

    return cats
}
```

---

### Feature 4: Release Notes Generation

**Format**: Keep a Changelog

**Template**:
```markdown
## [VERSION] - YYYY-MM-DD

### ⚠️ Breaking Changes (if any)
- Change description

### Added
- New feature descriptions

### Fixed
- Bug fix descriptions

### Changed
- Modification descriptions

### Documentation
- Doc update descriptions

### Statistics
- **Files changed**: X
- **Commits**: Y
- **Context analyzed**: ~Z tokens
```

**Implementation**:
```go
func generateReleaseNotes(version string, categories CommitCategories,
                          result *promptext.Result) string {
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
        len(categories.Features) + len(categories.Fixes) +
        len(categories.Changes) + len(categories.Docs)))
    notes.WriteString(fmt.Sprintf("- **Context analyzed**: ~%d tokens\n",
        result.TokenCount))
    notes.WriteString("\n")

    notes.WriteString("---\n\n")

    return notes.String()
}
```

---

### Feature 5: AI Prompt Generation

**Purpose**: Generate comprehensive prompts for LLMs to write polished release notes

**Prompt Structure**:
1. Task description
2. Context metadata (version, commits count, tokens)
3. Full commit history
4. Changed files summary
5. Complete code context from promptext
6. Task instructions
7. Requirements list
8. Example format

**Implementation**:
```go
func generateAIPrompt(version string, fromTag string, commits []string,
                       categories CommitCategories, result *promptext.Result) string {
    var prompt strings.Builder

    // Determine version
    if version == "" {
        version = "Unreleased"
    }

    // Header
    prompt.WriteString("# Release Notes Enhancement Request\n\n")
    prompt.WriteString("Please generate comprehensive release notes for version " +
        version + "\n\n")

    // Context metadata
    prompt.WriteString("## Context\n\n")
    prompt.WriteString(fmt.Sprintf("- **Version**: %s\n", version))
    prompt.WriteString(fmt.Sprintf("- **Changes since**: %s\n", fromTag))
    prompt.WriteString(fmt.Sprintf("- **Commits analyzed**: %d\n", len(commits)))
    prompt.WriteString(fmt.Sprintf("- **Files changed**: %d\n",
        len(result.ProjectOutput.Files)))
    prompt.WriteString(fmt.Sprintf("- **Context extracted**: ~%d tokens\n\n",
        result.TokenCount))

    // Commit history
    prompt.WriteString("## Commit History\n\n")
    prompt.WriteString("```\n")
    for _, commit := range commits {
        prompt.WriteString(commit + "\n")
    }
    prompt.WriteString("```\n\n")

    // Changed files summary
    prompt.WriteString("## Changed Files Summary\n\n")
    for _, file := range result.ProjectOutput.Files {
        prompt.WriteString(fmt.Sprintf("- `%s` (~%d tokens)\n",
            file.Path, file.Tokens))
    }
    prompt.WriteString("\n")

    // Full code context
    prompt.WriteString("## Code Context (via promptext)\n\n")
    prompt.WriteString("```\n")
    prompt.WriteString(result.FormattedOutput)
    prompt.WriteString("\n```\n\n")

    // Task instructions
    prompt.WriteString("## Task\n\n")
    prompt.WriteString("Generate release notes in Keep a Changelog format with these sections:\n\n")
    prompt.WriteString("### Added\n")
    prompt.WriteString("- New features (be specific and detailed)\n")
    prompt.WriteString("- Focus on user-facing value\n\n")
    prompt.WriteString("### Changed\n")
    prompt.WriteString("- Improvements and modifications\n\n")
    prompt.WriteString("### Fixed\n")
    prompt.WriteString("- Bug fixes\n\n")
    prompt.WriteString("### Documentation\n")
    prompt.WriteString("- Doc updates\n\n")

    // Requirements
    prompt.WriteString("## Requirements\n\n")
    prompt.WriteString("1. Use the commit history and code context to write detailed, clear descriptions\n")
    prompt.WriteString("2. Group related changes together logically\n")
    prompt.WriteString("3. Focus on user impact, not implementation details\n")
    prompt.WriteString("4. Be specific about what changed and why it matters\n")
    prompt.WriteString("5. Follow Keep a Changelog format\n")
    prompt.WriteString("6. Include markdown formatting for code, paths, etc.\n\n")

    // Example format
    prompt.WriteString("## Example Format\n\n")
    prompt.WriteString("```markdown\n")
    prompt.WriteString("## [" + version + "] - " +
        time.Now().Format("2006-01-02") + "\n\n")
    prompt.WriteString("### Added\n")
    prompt.WriteString("- **Feature Name**: Description of what was added\n")
    prompt.WriteString("  - Sub-detail about the feature\n")
    prompt.WriteString("  - Another aspect of the feature\n\n")
    prompt.WriteString("### Changed\n")
    prompt.WriteString("- **Component**: What changed and why it's better\n\n")
    prompt.WriteString("...\n")
    prompt.WriteString("```\n\n")

    prompt.WriteString("Please generate the complete, polished release notes now.\n")

    return prompt.String()
}
```

---

## Git Integration

### Git Commands Used

**1. Get Last Tag**:
```bash
git describe --tags --abbrev=0
```
Returns: `v0.7.3`

**2. Get Changed Files**:
```bash
git diff --name-only TAG..HEAD
```
Returns: List of file paths

**3. Get Commits**:
```bash
git log TAG..HEAD --pretty=format:%s
```
Returns: Commit messages (subject lines only)

**4. Get Commit Details** (future):
```bash
git log TAG..HEAD --pretty=format:"%h - %s (%an)"
```
Returns: Hash, message, author

### Error Scenarios

| Scenario | Command Exit Code | Handling |
|----------|-------------------|----------|
| Not a git repo | Non-zero | Fatal error with helpful message |
| No tags found | Non-zero | Fallback to `HEAD~10` |
| Tag doesn't exist | Non-zero | Error: "Tag not found" |
| No changes | Zero, empty output | Warning: "No changes" |

---

## Promptext Integration

### Library Documentation Reference

**Package**: `github.com/1broseidon/promptext/pkg/promptext`

**Primary Function**:
```go
func Extract(dir string, opts ...Option) (*Result, error)
```

**Functional Options Pattern**:
```go
type Option func(*Extractor)

func WithExtensions(exts ...string) Option
func WithTokenBudget(tokens int) Option
func WithFormat(format Format) Option
func WithExcludes(patterns ...string) Option
func WithRelevance(keywords ...string) Option
func WithGitIgnore(enabled bool) Option
func WithDefaultRules(enabled bool) Option
func WithVerbose(enabled bool) Option
func WithDebug(enabled bool) Option
```

### Integration Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/1broseidon/promptext/pkg/promptext"
)

func main() {
    // Extract code context
    result, err := promptext.Extract(".",
        promptext.WithExtensions(".go", ".md", ".yml", ".yaml"),
        promptext.WithTokenBudget(8000),
        promptext.WithFormat(promptext.FormatMarkdown),
    )

    if err != nil {
        log.Fatalf("Failed to extract context: %v", err)
    }

    // Access results
    fmt.Printf("Extracted %d tokens from %d files\n",
        result.TokenCount, len(result.ProjectOutput.Files))

    // Use formatted output
    fmt.Println(result.FormattedOutput)

    // Access individual files
    for _, file := range result.ProjectOutput.Files {
        fmt.Printf("File: %s (%d tokens)\n", file.Path, file.Tokens)
    }
}
```

### Token Budget Strategy

**Budget**: 8,000 tokens
- **Reasoning**: Leaves ~120K tokens for instructions in most LLMs
- **Coverage**: Typically covers 10-20 changed files with full context
- **Fallback**: If no relevant files, use 4,000 token budget for project summary

### File Extension Filtering

**Included**: `.go`, `.md`, `.yml`, `.yaml`
- Code files: `.go` for implementation changes
- Documentation: `.md` for README, docs, changelog
- Config: `.yml`, `.yaml` for configuration changes

**Excluded** (by promptext defaults):
- Lock files (`package-lock.json`, `go.sum`)
- Binary files
- Build artifacts (`dist/`, `build/`)
- Dependencies (`node_modules/`, `vendor/`)

---

## CLI Interface

### Command Syntax

```bash
promptext-notes [flags]
```

### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--version` | string | "" | Version to generate notes for (e.g., v0.7.4) |
| `--since` | string | "" | Generate notes since this tag (auto-detects if empty) |
| `--output` | string | "" | Output file path (stdout if empty) |
| `--ai-prompt` | bool | false | Generate AI enhancement prompt instead of basic notes |

### Flag Parsing Implementation

```go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
)

func main() {
    // Define flags
    version := flag.String("version", "", "Version to generate notes for (e.g., v0.7.4)")
    sinceTag := flag.String("since", "", "Generate notes since this tag (auto-detects if empty)")
    output := flag.String("output", "", "Output file (prints to stdout if empty)")
    aiPrompt := flag.Bool("ai-prompt", false, "Generate prompt for AI to enhance release notes")

    flag.Parse()

    // Validation
    if *version == "" {
        log.Println("Warning: No version specified, using 'Unreleased'")
    }

    // Get from tag
    fromTag := *sinceTag
    if fromTag == "" {
        fromTag = getLastTag()
    }

    fmt.Fprintf(os.Stderr, "📊 Analyzing changes since %s...\n", fromTag)

    // ... rest of implementation
}
```

### Progress Output (to stderr)

All progress messages go to `stderr`, actual output to `stdout`:

```go
fmt.Fprintln(os.Stderr, "📊 Analyzing changes since v0.7.3...")
fmt.Fprintf(os.Stderr, "   Found %d changed files\n", len(files))
fmt.Fprintf(os.Stderr, "   Found %d commits\n", len(commits))
fmt.Fprintln(os.Stderr, "\n🔍 Extracting code context with promptext...")
fmt.Fprintf(os.Stderr, "   Extracted context: ~%d tokens from %d files\n",
    tokens, fileCount)
fmt.Fprintln(os.Stderr, "\n📝 Generating release notes...")
```

**Note**: Do NOT use `\n` at the end of `fmt.Fprintln()` - it already adds a newline!

---

## Output Formats

### Mode 1: Basic Release Notes

**Flag**: (default, no flags)

**Output**: Markdown changelog

**Example**:
```markdown
## [0.7.4] - 2025-11-10

### Added
- New feature for code analysis
- Support for additional file types

### Fixed
- Bug in token counting
- Edge case in file filtering

### Statistics
- **Files changed**: 12
- **Commits**: 8
- **Context analyzed**: ~7,850 tokens
```

---

### Mode 2: AI Enhancement Prompt

**Flag**: `--ai-prompt`

**Output**: Comprehensive prompt for LLM

**Example**:
```markdown
# Release Notes Enhancement Request

Please generate comprehensive release notes for version 0.7.4

## Context

- **Version**: 0.7.4
- **Changes since**: v0.7.3
- **Commits analyzed**: 8
- **Files changed**: 12
- **Context extracted**: ~7,850 tokens

## Commit History

```
feat: add code analysis feature
fix: token counting bug
docs: update README
...
```

## Changed Files Summary

- `pkg/analyzer/analyzer.go` (~1,200 tokens)
- `README.md` (~450 tokens)
...

## Code Context (via promptext)

```
[Full formatted code context here]
```

## Task

Generate release notes in Keep a Changelog format with these sections:

### Added
- New features (be specific and detailed)
...

## Requirements

1. Use the commit history and code context to write detailed, clear descriptions
2. Group related changes together logically
...

Please generate the complete, polished release notes now.
```

---

## Implementation Guide

### Step 1: Project Setup

```bash
# Create project directory
mkdir promptext-notes
cd promptext-notes

# Initialize Go module
go mod init github.com/yourusername/promptext-notes

# Add promptext dependency
go get github.com/1broseidon/promptext@latest

# Create main file
touch main.go
```

### Step 2: Main Entry Point

**File**: `main.go`

```go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"

    "github.com/1broseidon/promptext/pkg/promptext"
)

func main() {
    // Parse flags
    version := flag.String("version", "", "Version to generate notes for")
    sinceTag := flag.String("since", "", "Generate notes since this tag")
    output := flag.String("output", "", "Output file")
    aiPrompt := flag.Bool("ai-prompt", false, "Generate AI prompt")
    flag.Parse()

    // Get from tag
    fromTag := *sinceTag
    if fromTag == "" {
        fromTag = getLastTag()
    }

    fmt.Fprintf(os.Stderr, "📊 Analyzing changes since %s...\n", fromTag)

    // Get changed files
    changedFiles, err := getChangedFiles(fromTag)
    if err != nil {
        log.Fatalf("Failed to get changed files: %v", err)
    }

    if len(changedFiles) == 0 {
        fmt.Fprintln(os.Stderr, "⚠️  No changes detected")
        return
    }

    fmt.Fprintf(os.Stderr, "   Found %d changed files\n", len(changedFiles))

    // Get commits
    commits, err := getCommits(fromTag)
    if err != nil {
        log.Fatalf("Failed to get commits: %v", err)
    }

    fmt.Fprintf(os.Stderr, "   Found %d commits\n", len(commits))

    // Extract code context
    fmt.Fprintln(os.Stderr, "\n🔍 Extracting code context with promptext...")
    result, err := extractCodeContext(changedFiles)
    if err != nil {
        log.Fatalf("Failed to extract context: %v", err)
    }

    fmt.Fprintf(os.Stderr, "   Extracted context: ~%d tokens from %d files\n",
        result.TokenCount, len(result.ProjectOutput.Files))

    // Categorize commits
    categories := categorizeCommits(commits)

    // Generate output
    var outputText string
    if *aiPrompt {
        fmt.Fprintln(os.Stderr, "\n📝 Generating AI prompt...")
        outputText = generateAIPrompt(*version, fromTag, commits, categories, result)
    } else {
        fmt.Fprintln(os.Stderr, "\n📝 Generating release notes...")
        outputText = generateReleaseNotes(*version, categories, result)
    }

    // Write output
    if *output != "" {
        if err := os.WriteFile(*output, []byte(outputText), 0644); err != nil {
            log.Fatalf("Failed to write output: %v", err)
        }
        fmt.Fprintf(os.Stderr, "✅ Written to %s\n", *output)
    } else {
        fmt.Println(outputText)
    }
}
```

### Step 3: Implement Helper Functions

Create separate files for each module (git.go, analyzer.go, generator.go, ai_prompt.go) with the functions specified in earlier sections.

### Step 4: Build & Test

```bash
# Build
go build -o promptext-notes

# Test basic mode
./promptext-notes --version v0.7.4

# Test AI prompt mode
./promptext-notes --version v0.7.4 --ai-prompt

# Test with output file
./promptext-notes --version v0.7.4 --output CHANGELOG.md

# Test with custom since tag
./promptext-notes --version v0.7.4 --since v0.7.0
```

---

## Testing Strategy

### Unit Tests

**Test Coverage**:
1. Git command execution
2. Commit categorization
3. File filtering
4. Output formatting

**Example Test**:
```go
package main

import (
    "testing"
)

func TestCategorizeCommits(t *testing.T) {
    commits := []string{
        "feat: add new feature",
        "fix: resolve bug",
        "docs: update README",
        "chore: update dependencies",
        "breaking: remove old API",
    }

    categories := categorizeCommits(commits)

    if len(categories.Features) != 1 {
        t.Errorf("Expected 1 feature, got %d", len(categories.Features))
    }

    if len(categories.Fixes) != 1 {
        t.Errorf("Expected 1 fix, got %d", len(categories.Fixes))
    }

    if len(categories.Breaking) != 1 {
        t.Errorf("Expected 1 breaking change, got %d", len(categories.Breaking))
    }
}

func TestGetLastTag(t *testing.T) {
    // Test in a git repo with tags
    tag := getLastTag()
    if tag == "" {
        t.Error("Expected a tag, got empty string")
    }
}
```

### Integration Tests

**Test Scenarios**:
1. Run in actual git repository
2. Test with different tag ranges
3. Verify output format
4. Test error conditions (no git, no tags, etc.)

### Manual Testing Checklist

- [ ] Run in repository with no tags
- [ ] Run with --version flag
- [ ] Run with --since flag
- [ ] Run with --output flag
- [ ] Run with --ai-prompt flag
- [ ] Test all flag combinations
- [ ] Verify markdown formatting
- [ ] Verify token counts are reasonable
- [ ] Test with repositories of different sizes

---

## Installation & Distribution

### Local Installation

```bash
go install github.com/yourusername/promptext-notes@latest
```

### Manual Build

```bash
# Clone repository
git clone https://github.com/yourusername/promptext-notes
cd promptext-notes

# Build
go build -o promptext-notes

# Install to PATH
sudo mv promptext-notes /usr/local/bin/
```

### Cross-Platform Builds

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o promptext-notes-linux-amd64

# macOS
GOOS=darwin GOARCH=amd64 go build -o promptext-notes-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o promptext-notes-darwin-arm64

# Windows
GOOS=windows GOARCH=amd64 go build -o promptext-notes-windows-amd64.exe
```

### GitHub Releases (Optional)

Use GoReleaser for automated releases:

**.goreleaser.yml**:
```yaml
builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"

checksum:
  name_template: 'checksums.txt'

release:
  github:
    owner: yourusername
    name: promptext-notes
```

---

## Usage Examples

### Example 1: Basic Release Notes

```bash
# Generate notes for v1.0.0 since last tag
promptext-notes --version v1.0.0

# Output:
## [1.0.0] - 2025-11-10

### Added
- Initial release
- Core functionality

### Statistics
- **Files changed**: 25
- **Commits**: 15
- **Context analyzed**: ~6,500 tokens
```

### Example 2: AI-Enhanced Release Notes

```bash
# Generate AI prompt for v1.0.0
promptext-notes --version v1.0.0 --ai-prompt > prompt.txt

# Copy prompt.txt and paste into Claude/ChatGPT/etc.
# AI will generate polished release notes using full code context
```

### Example 3: Custom Date Range

```bash
# Generate notes from v0.5.0 to current
promptext-notes --version v1.0.0 --since v0.5.0
```

### Example 4: Write to File

```bash
# Append to CHANGELOG.md
promptext-notes --version v1.0.0 --output release-notes.md
cat release-notes.md >> CHANGELOG.md
```

### Example 5: CI/CD Integration

```bash
# In GitHub Actions
- name: Generate Release Notes
  run: |
    go install github.com/yourusername/promptext-notes@latest
    promptext-notes --version ${{ github.ref_name }} --output RELEASE_NOTES.md

- name: Create Release
  uses: softprops/action-gh-release@v1
  with:
    body_path: RELEASE_NOTES.md
```

---

## Future Enhancements

### Phase 2 Features

1. **Custom Templates**: Allow users to specify custom output templates
2. **Multiple Output Formats**: JSON, HTML, PDF
3. **Author Attribution**: Group changes by contributor
4. **Link Generation**: Auto-link issues and PRs (e.g., `#123` → GitHub link)
5. **Change Impact Analysis**: Categorize changes by scope (major/minor/patch)
6. **Interactive Mode**: Prompt user to edit/refine notes before output
7. **Changelog File Update**: Automatically prepend to existing CHANGELOG.md
8. **Release Type Detection**: Auto-suggest version based on changes (semver)

### Phase 3 Features

1. **Multi-Repository Support**: Aggregate changes across monorepos
2. **Jira/Linear Integration**: Fetch ticket details for commit messages
3. **Slack/Discord Notifications**: Post release notes to team channels
4. **Historical Analysis**: Compare current release to past releases
5. **Breaking Change Detection**: Parse code for API changes
6. **Migration Guide Generation**: Auto-generate upgrade instructions

---

## Implementation Checklist

### Core Functionality
- [ ] Set up Go project with promptext dependency
- [ ] Implement CLI flag parsing
- [ ] Implement `getLastTag()` function
- [ ] Implement `getChangedFiles()` function
- [ ] Implement `getCommits()` function
- [ ] Implement `extractCodeContext()` function
- [ ] Implement `categorizeCommits()` function
- [ ] Implement `generateReleaseNotes()` function
- [ ] Implement `generateAIPrompt()` function
- [ ] Add progress output to stderr
- [ ] Add output file writing

### Testing
- [ ] Write unit tests for commit categorization
- [ ] Write integration tests
- [ ] Test in real git repositories
- [ ] Test error conditions
- [ ] Test all flag combinations

### Documentation
- [ ] Write README.md
- [ ] Add usage examples
- [ ] Document all CLI flags
- [ ] Add contributing guidelines
- [ ] Create LICENSE file

### Distribution
- [ ] Set up GitHub repository
- [ ] Configure GitHub Actions for testing
- [ ] Set up GoReleaser (optional)
- [ ] Publish first release
- [ ] Add installation instructions

### Polish
- [ ] Add colored output (optional)
- [ ] Improve error messages
- [ ] Add verbose/debug mode
- [ ] Add version flag (`--version`)
- [ ] Add help text (`--help`)

---

## Quick Start for Implementation

```bash
# 1. Create project
mkdir promptext-notes && cd promptext-notes

# 2. Initialize Go module
go mod init github.com/yourusername/promptext-notes

# 3. Add dependency
go get github.com/1broseidon/promptext@latest

# 4. Create main.go with the implementation from this spec

# 5. Build
go build

# 6. Test
./promptext-notes --version v1.0.0

# 7. Success! You now have a working release notes generator
```

---

## Support & Resources

- **Promptext Library**: https://github.com/1broseidon/promptext
- **Promptext Docs**: https://promptext.sh
- **Go Documentation**: https://go.dev/doc/
- **Keep a Changelog**: https://keepachangelog.com/
- **Conventional Commits**: https://www.conventionalcommits.org/

---

**End of Specification**

*This document contains everything needed to build promptext-notes from scratch. No prior context required.*
