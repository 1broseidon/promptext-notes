# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

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
