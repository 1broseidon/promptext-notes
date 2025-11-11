# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

---

## [v0.5.0] - 2025-11-11

### Added
- **Multi-provider AI support** - Generate release notes using OpenAI, Anthropic, Cerebras, or Groq with configurable models and automatic provider selection
- **GitHub Actions automation** - Automated workflow that generates AI-enhanced release notes on version tag push with full GitHub release integration
- **Local script generation** - `generate-release-notes.sh` script for generating AI-enhanced notes locally with support for all AI providers
- **Configurable AI models** - GitHub Variables support for customizing which AI model each provider uses (e.g., `OPENAI_MODEL`, `CEREBRAS_MODEL`)

### Changed
- **AI model upgrade** - Switched default AI model to gpt-oss-120b for significantly better quality and more coherent release notes generation
- **Enhanced prompt rules** - Improved AI prompts with explicit rules to omit non-user-value content and stricter categorization for better relevance

### Fixed
- **YAML workflow syntax** - Corrected heredoc syntax in GitHub Actions workflow to properly handle multi-line strings
- **Duplicate headers** - Prevented duplicate section headers in CHANGELOG output

---

---

## [v0.5.1] - 2025-11-11

### Fixed
- **YAML workflow syntax** - Corrected heredoc syntax in GitHub Actions workflow to properly handle multi-line strings

---

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
