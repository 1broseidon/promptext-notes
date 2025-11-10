package prompt

import (
	"fmt"
	"strings"
	"time"

	"github.com/1broseidon/promptext-notes/internal/analyzer"
	"github.com/1broseidon/promptext/pkg/promptext"
)

// GenerateAIPrompt generates a comprehensive prompt for LLMs to write polished release notes.
func GenerateAIPrompt(version, fromTag string, commits []string, categories analyzer.CommitCategories, result *promptext.Result) string {
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
