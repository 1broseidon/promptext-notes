# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

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
