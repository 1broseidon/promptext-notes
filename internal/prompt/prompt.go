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
	prompt.WriteString("Generate release notes in Keep a Changelog format with ONLY these sections.\n")
	prompt.WriteString("**IMPORTANT: Omit any section that has no entries. Do not include placeholder text for empty sections.**\n\n")

	prompt.WriteString("### ⚠️ BREAKING CHANGES (if any)\n")
	prompt.WriteString("- Changes that break existing code or require user action\n")
	prompt.WriteString("- API changes that affect backwards compatibility\n")
	prompt.WriteString("- Maximum 2 sentences per item\n\n")

	prompt.WriteString("### Added\n")
	prompt.WriteString("- **New user-facing features ONLY**\n")
	prompt.WriteString("- What users can now do that they couldn't before\n")
	prompt.WriteString("- Format: Brief, 1-2 sentences maximum per item\n")
	prompt.WriteString("- Focus on capabilities, not implementation\n\n")

	prompt.WriteString("### Changed\n")
	prompt.WriteString("- **User-visible improvements** to existing features\n")
	prompt.WriteString("- What works better or differently for users\n")
	prompt.WriteString("- Format: Brief, 1-2 sentences maximum per item\n")
	prompt.WriteString("- NOT internal refactoring or code reorganization\n\n")

	prompt.WriteString("### Fixed\n")
	prompt.WriteString("- **User-impacting bug fixes ONLY**\n")
	prompt.WriteString("- What problems users will no longer experience\n")
	prompt.WriteString("- Format: Brief, 1 sentence per fix\n")
	prompt.WriteString("- NOT internal bugs users never saw\n\n")

	prompt.WriteString("### Deprecated (if any)\n")
	prompt.WriteString("- Features or APIs that will be removed in future versions\n")
	prompt.WriteString("- Include timeline if known\n\n")

	prompt.WriteString("### Security (if any)\n")
	prompt.WriteString("- Security improvements or vulnerability fixes\n")
	prompt.WriteString("- Be specific but don't reveal exploits\n\n")

	// Requirements
	prompt.WriteString("## Critical Rules - MUST FOLLOW\n\n")
	prompt.WriteString("**MUST OMIT** (these provide ZERO user value):\n")
	prompt.WriteString("- ❌ Documentation updates (README changes, CHANGELOG meta-references)\n")
	prompt.WriteString("- ❌ Internal refactoring or code reorganization\n")
	prompt.WriteString("- ❌ Test coverage, CI/CD, or build system changes\n")
	prompt.WriteString("- ❌ Implementation details (token budgets, internal APIs, data structures)\n")
	prompt.WriteString("- ❌ Statistics sections (file counts, commit counts, token counts)\n")
	prompt.WriteString("- ❌ Empty sections with placeholders like \"*(No fixes)*\"\n")
	prompt.WriteString("- ❌ Meta-references (\"Updated CHANGELOG\", \"Added release notes\")\n\n")

	prompt.WriteString("**MUST FOCUS ON**:\n")
	prompt.WriteString("- ✅ User-facing changes ONLY\n")
	prompt.WriteString("- ✅ What users can do differently\n")
	prompt.WriteString("- ✅ Problems users will no longer experience\n")
	prompt.WriteString("- ✅ Features users need to know about\n")
	prompt.WriteString("- ✅ Breaking changes that require action\n\n")

	prompt.WriteString("**LENGTH LIMITS**:\n")
	prompt.WriteString("- Maximum 2 sentences per item\n")
	prompt.WriteString("- Maximum 3 sub-points per item if needed\n")
	prompt.WriteString("- Be concise and direct - no fluff\n\n")

	prompt.WriteString("**CATEGORIZATION RULES**:\n")
	prompt.WriteString("- \"Added\" = Truly NEW features users didn't have before\n")
	prompt.WriteString("- \"Changed\" = IMPROVEMENTS to existing features\n")
	prompt.WriteString("- \"Fixed\" = BUG FIXES that resolve user problems\n")
	prompt.WriteString("- Do NOT list the same change in multiple sections\n")
	prompt.WriteString("- Model upgrades go in \"Changed\", not \"Added\"\n\n")

	// Example format
	prompt.WriteString("## Example Format\n\n")
	prompt.WriteString("```markdown\n")
	prompt.WriteString("## [" + version + "] - " +
		time.Now().Format("2006-01-02") + "\n\n")
	prompt.WriteString("### ⚠️ BREAKING CHANGES\n")
	prompt.WriteString("- **API endpoint changes** - `/api/v1/users` is now `/api/v2/users`. Update all API calls.\n\n")
	prompt.WriteString("### Added\n")
	prompt.WriteString("- **PDF export** - Export release notes as styled PDF documents\n")
	prompt.WriteString("- **Dark mode** - Toggle dark theme in settings for better readability\n\n")
	prompt.WriteString("### Changed\n")
	prompt.WriteString("- **Performance** - Release note generation is 3x faster for large repositories\n")
	prompt.WriteString("- **Error messages** - More helpful context when git operations fail\n\n")
	prompt.WriteString("### Fixed\n")
	prompt.WriteString("- **Unicode handling** - Fixed crash when commit messages contain emoji\n")
	prompt.WriteString("- **Memory leak** - Resolved issue causing high memory usage on large repos\n\n")
	prompt.WriteString("### Deprecated\n")
	prompt.WriteString("- **Old API format** - Legacy `/api/v1` endpoints deprecated, will be removed in v2.0.0\n")
	prompt.WriteString("```\n\n")

	prompt.WriteString("Generate ONLY the sections with content. Omit empty sections entirely. Be ruthlessly focused on user value.\n")

	return prompt.String()
}
