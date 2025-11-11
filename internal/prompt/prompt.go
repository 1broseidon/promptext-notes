package prompt

import (
	"fmt"
	"strings"
	"time"

	"github.com/1broseidon/promptext-notes/internal/analyzer"
	"github.com/1broseidon/promptext/pkg/promptext"
)

// GenerateAIPrompt generates a comprehensive prompt for LLMs to write polished release notes.
func GenerateAIPrompt(version, fromTag string, commits []string, categories analyzer.CommitCategories, result *promptext.Result, diffStats, diff string) string {
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

	// IMPROVEMENT #1: Executive Summary
	prompt.WriteString("## 🎯 Executive Summary\n\n")
	prompt.WriteString("**Quick Overview**: Read this section first to understand what changed at a high level.\n\n")

	// Determine change type based on categories
	changeTypes := []string{}
	if len(categories.Breaking) > 0 {
		changeTypes = append(changeTypes, "breaking changes")
	}
	if len(categories.Features) > 0 {
		changeTypes = append(changeTypes, "new features")
	}
	if len(categories.Fixes) > 0 {
		changeTypes = append(changeTypes, "bug fixes")
	}
	if len(categories.Changes) > 0 {
		changeTypes = append(changeTypes, "improvements")
	}

	changeTypeStr := "miscellaneous updates"
	if len(changeTypes) > 0 {
		changeTypeStr = strings.Join(changeTypes, ", ")
	}

	prompt.WriteString(fmt.Sprintf("- **Change Type**: %s\n", changeTypeStr))
	prompt.WriteString(fmt.Sprintf("- **Files Modified**: %d files changed\n", len(result.ProjectOutput.Files)))

	// List key files (top 3 by token count)
	if len(result.ProjectOutput.Files) > 0 {
		prompt.WriteString("- **Key Files**: ")
		numFiles := len(result.ProjectOutput.Files)
		if numFiles > 3 {
			numFiles = 3
		}
		for i := 0; i < numFiles; i++ {
			if i > 0 {
				prompt.WriteString(", ")
			}
			prompt.WriteString("`" + result.ProjectOutput.Files[i].Path + "`")
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString(fmt.Sprintf("- **Commit Count**: %d commit(s)\n", len(commits)))
	prompt.WriteString("\n")
	prompt.WriteString("**Focus Areas**: Analyze the diff below as your PRIMARY source of truth.\n\n")

	// IMPROVEMENT #2: Git Diff Stats and Diff View - NOW MANDATORY AND FIRST
	prompt.WriteString("## 📊 Git Diff Summary (PRIMARY SOURCE)\n\n")
	prompt.WriteString("**CRITICAL**: This is your PRIMARY source. The diff shows EXACTLY what changed line-by-line.\n\n")

	if diffStats != "" {
		prompt.WriteString("### Change Magnitude\n\n")
		prompt.WriteString("```\n")
		prompt.WriteString(diffStats)
		prompt.WriteString("\n```\n\n")
	}

	// Special handling for CHANGELOG-only changes
	if strings.Contains(diffStats, "CHANGELOG.md") && !strings.Contains(diffStats, ".go") && !strings.Contains(diffStats, ".yml") {
		prompt.WriteString("⚠️ **WARNING**: Only CHANGELOG.md changed. This means:\n")
		prompt.WriteString("- No actual code changes for users\n")
		prompt.WriteString("- This is likely an automated documentation update\n")
		prompt.WriteString("- Correct response: \"No user-facing changes in this version\"\n\n")
	}

	// Add full diff for small changes (< 200 lines)
	diffLines := strings.Count(diff, "\n")
	if diff != "" && diffLines > 0 {
		prompt.WriteString("### Detailed Line-by-Line Diff\n\n")
		prompt.WriteString("**Use this as ground truth**: Every + is an addition, every - is a deletion.\n\n")
		prompt.WriteString("```diff\n")
		prompt.WriteString(diff)
		prompt.WriteString("\n```\n\n")
	}

	// Full code context - SECONDARY SOURCE
	prompt.WriteString("## Code Context (via promptext) - SECONDARY SOURCE\n\n")
	prompt.WriteString("**Note**: This shows full file contents for context. The DIFF above is more accurate for what actually changed.\n\n")
	prompt.WriteString("```\n")
	prompt.WriteString(result.FormattedOutput)
	prompt.WriteString("\n```\n\n")

	// Changed files summary
	prompt.WriteString("## Changed Files Summary\n\n")
	for _, file := range result.ProjectOutput.Files {
		prompt.WriteString(fmt.Sprintf("- `%s` (~%d tokens)\n",
			file.Path, file.Tokens))
	}
	prompt.WriteString("\n")

	// Commit history - REFERENCE ONLY
	prompt.WriteString("## Commit History (Reference Only)\n\n")
	prompt.WriteString("**NOTE**: Commit messages may be incomplete or misleading. Rely on the actual code changes above to understand the true nature of changes.\n\n")
	prompt.WriteString("```\n")
	for _, commit := range commits {
		prompt.WriteString(commit + "\n")
	}
	prompt.WriteString("```\n\n")

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

	// IMPROVEMENT #3: Consolidated Requirements
	prompt.WriteString("## Critical Rules\n\n")
	prompt.WriteString("**PRIMARY SOURCE**: Code changes (diff/context above), NOT commit messages\n\n")
	prompt.WriteString("**USER VALUE ONLY**: Omit internal changes (refactoring, tests, CI/CD, docs, meta-references)\n\n")
	prompt.WriteString("**FORMAT**: 1-2 sentences max | Omit empty sections | No placeholders\n\n")
	prompt.WriteString("**CATEGORIES**:\n")
	prompt.WriteString("- Added = NEW capabilities users didn't have\n")
	prompt.WriteString("- Changed = IMPROVEMENTS to existing features\n")
	prompt.WriteString("- Fixed = BUG FIXES solving user problems\n")
	prompt.WriteString("- Breaking = Changes requiring user action\n\n")

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
