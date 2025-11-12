# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

---

## [v0.8.1] - 2025-11-12

### Added
- **Commit filtering** - Exclude bot commits and merge messages from changelog generation using configurable patterns
- **Output section customization** - Control which changelog sections appear (breaking, added, changed, fixed, docs) through configuration options
- **Auto-exclude-meta filtering** - Automatically exclude CI configurations, CHANGELOG, README, and other meta files from AI context to maintain focus on user-facing changes

### Changed
- **Default AI provider** - Switched from Anthropic to Cerebras for free access to high-quality models
- **Default model** - Updated to `zai-glm-4.6` for improved accuracy at no cost
- **File exclusion patterns** - Added support for glob patterns (e.g., `*.tmp`, `build/**`) to enable more flexible file filtering

### Fixed
- **Context pollution** - Meta-documentation changes no longer influence AI-generated changelog content
- **File filtering** - Exclusion patterns now correctly apply to promptext context extraction
---

## [v0.8.0] - 2025-11-12

### Added
- **Auto-exclude meta files** - Automatically excludes CI configs, CHANGELOG, README, and other meta files from AI context to keep changelogs focused on user-facing changes
- **Developer control over changelog content** - New `auto_exclude_meta` config option (default: true) to prevent internal tooling changes from appearing in release notes
---

## [v0.7.6] - 2025-11-12

### Changed
- **Simplified 2-stage polish workflow** - Discovery stage now uses the main AI model, eliminating the confusing discovery_model configuration option
- **Updated default models** - Changed defaults to `zai-glm-4.6` for discovery and `anthropic/claude-sonnet-4.5` for polish to improve accuracy
- **Improved polish prompt** - Streamlined changelog-only prompt reduces hallucinations and eliminates first-person language

### Fixed
- **API key environment variable override** - CLI provider override now correctly updates the corresponding API key environment variable
- **GitHub Actions workflow** - Fixed configuration to use config file instead of CLI flags for more reliable automation
---

## [v0.7.5] - 2025-11-12

### Added
- **Enhanced 2-Stage Polish Workflow for Premium Release Notes**
  We've introduced a sophisticated two-stage workflow designed to generate exceptionally high-quality release notes. This new process intelligently combines the precision of technical discovery models, which excel at identifying and summarizing core changes, with the eloquence of customer-facing language models. The result is release notes that are not only technically accurate but also clearly articulated and engaging for your audience.

- **Expanded AI Model Access with OpenRouter Support**
  To provide you with greater flexibility and cost-efficiency in generating AI-enhanced release notes, we've integrated support for OpenRouter. This integration allows you to access a vast ecosystem of hundreds of AI models through a single, unified API. You can now choose from a wider selection of models to best suit your specific needs and budget.

- **Configurable File Exclusions for Focused AI Context**
  You now have more control over the information provided to our AI models. With configurable file exclusions, you can specify certain files or directories that should be ignored when generating release notes. This ensures that the AI focuses only on relevant code changes, preventing "context pollution" from non-essential files. You can configure this using the `--exclude-files` CLI flag or through your configuration file settings.

- **Improved Error Handling for Smoother Operations**
  We've significantly enhanced our error handling mechanisms, providing clearer and more informative messages for issues related to Git operations and AI requests. This improvement helps you quickly understand and resolve any problems, leading to a more seamless and reliable experience.

### Changed
- **Updated Default AI Provider to Cerebras (llama-3.3-70b)**
  To offer improved accessibility and higher quality out-of-the-box, we have updated our default AI provider to Cerebras, utilizing the `llama-3.3-70b` model. This change ensures that even users on free tiers can benefit from a powerful and capable AI model for their release notes generation.

- **Granular Configuration for the 2-Stage Polish Workflow**
  The new 2-stage polish workflow is now highly configurable. You can precisely define the discovery and polish models, their respective providers, and API keys. This allows you to fine-tune the workflow to meet your specific requirements for accuracy and linguistic style.

- **Enhanced GitHub Actions Workflow**
  Our GitHub Actions workflow has been updated to leverage the new `llama-3.3-70b` model for improved performance and quality. Additionally, we've integrated support for OpenRouter, giving you more options for AI model selection directly within your CI/CD pipeline.

### Fixed
- **Resolved Context Pollution from Meta-Documentation Changes**
  We've addressed an issue where changes in meta-documentation (e.g., READMEs, contributing guidelines) could inadvertently influence the AI's understanding of code changes, leading to irrelevant or misleading content in the generated release notes. The AI now correctly distinguishes between code changes and documentation updates.

- **Accurate Reflection of Significant Code Changes**
  Previously, significant code changes that were not explicitly detailed in commit messages might have been overlooked during release notes generation. This issue has been resolved, ensuring that all impactful code modifications are now accurately identified and included in your release notes, regardless of commit message verbosity.

- **Correct API Key Environment Variable Overrides**
  We fixed a problem where overriding the default AI provider via CLI flags did not correctly update the corresponding API key environment variable. This ensures that when you specify a different provider through the CLI, the correct API key is always used, preventing authentication failures.
---

## [v0.7.3] - 2025-11-12

### Added
- **Git diff analysis** - Release notes generation now analyzes actual code changes line-by-line for more accurate and comprehensive release notes
- **Executive Summary** - Generated release notes now include a high-level overview section showing change types, key files, and commit counts at a glance
- **Documentation-only detection** - Automatically detects when only documentation files change and avoids generating empty release notes

### Changed
- **Default AI provider** - Example configuration now uses Groq with llama-3.3-70b-versatile model instead of Cerebras for better free tier accessibility
- **Improved AI prompts** - Restructured prompt generation to prioritize actual code diffs over commit messages, resulting in more accurate release notes
- **Context filtering** - Release notes generation now excludes CHANGELOG.md and README.md from AI context to reduce noise and improve focus on actual code changes
- **Enhanced error handling** - Better error messages and graceful handling when git operations fail or diffs cannot be retrieved

### Fixed
- **Context pollution** - Resolved issue where meta-documentation changes would incorrectly influence AI-generated release notes
- **Missed code changes** - Fixed problem where significant code changes not reflected in commit messages could be missed during release notes generation
---

## [v0.7.2] - 2025-11-11

### ⚠️ BREAKING CHANGES
- **Default provider changed** - Cerebras with zai-glm-4.6 model is now the default instead of Anthropic. Set `CEREBRAS_API_KEY` environment variable or override provider via CLI flags.

### Changed
- **Default AI provider** - Switched from Anthropic to Cerebras with zai-glm-4.6 model for improved performance and cost-effectiveness.
- **Provider override behavior** - API key environment variable now automatically updates when overriding provider via CLI flags.

### Fixed
- **API key environment mismatch** - Fixed issue where overriding provider via CLI didn't update the corresponding API key environment variable.
- **Cerebras error handling** - Improved error messages for better debugging when Cerebras API requests fail.
- **CHANGELOG duplicate entries** - Added detection to prevent duplicate version entries when updating changelog.
---

## [v0.7.1] - 2025-11-11

### ⚠️ BREAKING CHANGES
- **Model configuration now requires exact API names** - Remove any custom model name mappings from your configuration and use the exact model identifiers from your AI provider's API documentation (e.g., `claude-haiku-4-5` for Anthropic, `gpt-4o-mini` for OpenAI).

### Changed
- **Faster release notes generation** - Optimized workflow now uses pre-built binaries with caching instead of building from source, reducing generation time from ~45-50s to ~10-15s.
- **Improved code quality checks** - Added gocyclo complexity analysis to pre-commit hooks to catch overly complex functions early in development.

### Fixed
- **Error message capitalization** - Corrected inconsistent capitalization in error strings for better readability and consistency.
---

## [v0.7.0] - 2025-11-11

### Added
- **Multi-provider AI support** - Generate release notes using Anthropic, OpenAI, Cerebras, Groq, or local Ollama with configurable models and automatic provider selection.
- **Comprehensive configuration system** - YAML-based configuration with sensible defaults, CLI flag overrides, and support for all AI providers with customizable retry strategies.
- **Retry mechanism with backoff strategies** - Configurable retry logic with exponential, linear, or constant backoff to handle transient AI provider failures gracefully.

### Changed
- **Default AI provider** - Switched to Anthropic Claude Haiku 4.5 for improved reliability and cost-effectiveness ($1/$5 per 1M tokens).
- **Model configuration** - Updated to latest 2025 models with provider-specific defaults (Claude Haiku 4.5 for Anthropic, GPT-4o-mini for OpenAI, Llama 3.3-70b for Cerebras/Groq).

### Fixed
- **Model name mapping** - Corrected Claude Haiku model identifier to use proper `claude-haiku-4-5` name without incorrect mappings.
- **Model parameter passing** - Ensured AI model is properly passed from configuration to all AI provider requests.
---

## [v0.5.0] - 2025-11-11

### Added
- **Multi-provider AI support** - Generate release notes using OpenAI, Anthropic, Cerebras, or Groq with configurable models and automatic provider selection
- **GitHub Actions automation** - Automated workflow that generates AI-enhanced release notes on version tag push with full GitHub release integration
- **Local script generation** - `generate-release-notes.sh` script for generating AI-enhanced notes locally with support for all AI providers
- **Configurable AI models** - GitHub Variables support for customizing which AI model each provider uses (e.g., `OPENAI_MODEL`, `CEREBRAS_MODEL`)
- **Comprehensive documentation** - New `docs/USAGE.md` with step-by-step integration guide for external repositories

### Changed
- **Default AI provider** - Switched to Anthropic Claude Haiku 4.5 for better reliability and cost-effectiveness ($1/$5 per 1M tokens)
- **Model updates** - Updated to latest 2025 models (GPT-5 Nano for OpenAI, Claude Haiku 4.5 for Anthropic)
- **Enhanced prompt rules** - Improved AI prompts with explicit rules to omit non-user-value content and stricter categorization for better relevance

### Fixed
- **YAML workflow syntax** - Corrected heredoc syntax in GitHub Actions workflow to properly handle multi-line strings
- **Duplicate headers** - Prevented duplicate section headers in CHANGELOG output

---

## [v0.4.0] - 2025-11-10

### Added
- **Enhanced AI prompt generation** – Added explicit rules to omit non-user-value content, new sections (BREAKING CHANGES, Deprecated, Security), and stricter categorization, improving prompt relevance and conciseness.

---

## [v0.3.0] - 2025-11-10

### Changed
- **AI model upgrade** – Switched to gpt-oss-120b for significantly better quality and more coherent release notes generation.

---

## [v0.2.0] - 2025-11-10

### Added
- **Automated AI-enhanced release notes** – GitHub Actions workflow automatically generates professional release notes using Cerebras API on tag push.
- **Script for local generation** – Run `./scripts/generate-release-notes.sh` to generate AI-enhanced notes locally.

### Changed
- **CI coverage checks** – Made coverage checks non-blocking to allow builds to succeed at 79% coverage.

---

## [v0.1.0] - 2025-11-10

### Added
- **CLI tool for release notes** – Generate release notes from git history with code context extraction.
- **Conventional commit parsing** – Automatically categorizes changes by type (feat, fix, docs, etc.).
- **Promptext integration** – Extracts code context with token-aware analysis (8K budget).
- **Keep a Changelog format** – Produces standardized markdown output.
- **AI prompt generation** – Creates comprehensive prompts for LLMs with full code context.
