# Promptext-Notes

[![CI](https://github.com/1broseidon/promptext-notes/actions/workflows/ci.yml/badge.svg)](https://github.com/1broseidon/promptext-notes/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/1broseidon/promptext-notes)](https://goreportcard.com/report/github.com/1broseidon/promptext-notes)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go-based CLI tool that generates intelligent, context-aware release notes by combining git history analysis with code context extraction using the [promptext](https://github.com/1broseidon/promptext) library.

## Features

- 📊 **Git History Analysis**: Automatically analyzes commits since the last tag
- 🔍 **Code Context Extraction**: Uses promptext to extract relevant code changes with token-aware analysis
- 📝 **Conventional Commits**: Categorizes changes by type (feat, fix, docs, breaking, etc.)
- 🤖 **AI-Ready Prompts**: Generates comprehensive prompts for LLMs to write polished release notes
- 📋 **Keep a Changelog Format**: Produces standardized markdown output
- ⚡ **Fast & Lightweight**: Single binary with no runtime dependencies (except Git)

## Installation

### Using `go install`

```bash
go install github.com/1broseidon/promptext-notes/cmd/promptext-notes@latest
```

### From Source

```bash
git clone https://github.com/1broseidon/promptext-notes.git
cd promptext-notes
go build -o promptext-notes ./cmd/promptext-notes
sudo mv promptext-notes /usr/local/bin/
```

### Download Pre-built Binary

Download the latest release from the [releases page](https://github.com/1broseidon/promptext-notes/releases).

## Usage

### Basic Release Notes

Generate release notes for a specific version:

```bash
promptext-notes --version v1.0.0
```

Output:
```markdown
## [v1.0.0] - 2025-11-10

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

### AI-Enhanced Release Notes

Generate a comprehensive prompt for an LLM:

```bash
promptext-notes --version v1.0.0 --ai-prompt > prompt.txt
```

Then paste the contents of `prompt.txt` into Claude, ChatGPT, or your preferred LLM to get polished, detailed release notes.

### Custom Date Range

Specify a starting tag/commit:

```bash
promptext-notes --version v1.0.0 --since v0.5.0
```

### Output to File

Write release notes to a file:

```bash
promptext-notes --version v1.0.0 --output RELEASE_NOTES.md
```

### Append to CHANGELOG

```bash
promptext-notes --version v1.0.0 --output release-notes.md
cat release-notes.md >> CHANGELOG.md
```

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--version` | string | "" | Version to generate notes for (e.g., v0.7.4) |
| `--since` | string | "" | Generate notes since this tag (auto-detects if empty) |
| `--output` | string | "" | Output file path (stdout if empty) |
| `--ai-prompt` | bool | false | Generate AI enhancement prompt instead of basic notes |

## How It Works

1. **Git Analysis**: Retrieves changed files and commit messages since the last tag (or specified tag)
2. **Context Extraction**: Uses [promptext](https://github.com/1broseidon/promptext) to extract code context from changed files (.go, .md, .yml, .yaml)
3. **Categorization**: Parses commit messages using conventional commit format
4. **Generation**: Produces either:
   - **Basic Mode**: Keep a Changelog formatted release notes
   - **AI Mode**: Comprehensive prompt with full code context for LLM enhancement

## Automated Release Notes (AI-Enhanced)

This project includes an automated workflow that generates AI-enhanced release notes using fast, free inference APIs (Cerebras or Grok).

### How It Works

When you push a version tag (e.g., `v1.0.0`), the workflow automatically:

1. ✅ Builds the promptext-notes binary
2. 🔍 Analyzes git history and extracts code context
3. 🤖 Sends the prompt to Cerebras/Grok API for AI enhancement
4. 📝 Creates a GitHub release with polished notes
5. 📋 Updates CHANGELOG.md in the repository

### Setup

1. **Get a free API key** from [Cerebras](https://cerebras.ai) or [Grok](https://x.ai)

2. **Add the API key to GitHub Secrets**:
   - Go to your repository → Settings → Secrets and variables → Actions
   - Add `CEREBRAS_API_KEY` or `GROK_API_KEY`

3. **Push a version tag**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

The workflow will automatically generate and publish AI-enhanced release notes!

### Manual Script Usage

You can also use the script locally:

```bash
# Set API key
export CEREBRAS_API_KEY="your-key-here"

# Generate release notes
./scripts/generate-release-notes.sh v1.0.0

# With custom previous tag and provider
./scripts/generate-release-notes.sh v1.0.0 v0.9.0 cerebras

# Save to file
RELEASE_NOTES_FILE=notes.md ./scripts/generate-release-notes.sh v1.0.0
```

### Supported AI Providers

| Provider | Model | Context | Speed | Free Tier |
|----------|-------|---------|-------|-----------|
| **Cerebras** | gpt-oss-120b | 65K | ⚡ Ultra-fast | ✅ Yes (Default) |
| **Cerebras** | llama-3.3-70b | 65K | ⚡ Ultra-fast | ✅ Yes |
| **Cerebras** | zai-glm-4.6 | 64K | ⚡ Ultra-fast | ✅ Yes |
| **Grok** | grok-beta | - | ⚡ Fast | ✅ Yes |

## CI/CD Integration

### GitHub Actions (Basic)

```yaml
- name: Generate Release Notes
  run: |
    go install github.com/1broseidon/promptext-notes/cmd/promptext-notes@latest
    promptext-notes --version ${{ github.ref_name }} --output RELEASE_NOTES.md

- name: Create Release
  uses: softprops/action-gh-release@v1
  with:
    body_path: RELEASE_NOTES.md
```

### GitHub Actions (With AI Enhancement)

The repository includes a complete automated workflow. See `.github/workflows/auto-docs.yml`.

### GitLab CI

```yaml
release:
  script:
    - go install github.com/1broseidon/promptext-notes/cmd/promptext-notes@latest
    - promptext-notes --version $CI_COMMIT_TAG --output RELEASE_NOTES.md
```

## Development

### Prerequisites

- Go 1.22 or higher
- Git

### Build

```bash
go build -o promptext-notes ./cmd/promptext-notes
```

### Test

```bash
go test ./... -v
```

### Test with Coverage

```bash
go test ./... -cover
```

Current coverage: **88.66%**

### Quality Checks

```bash
# Format code
go fmt ./...

# Run staticcheck
staticcheck ./...

# Check cyclomatic complexity
gocyclo -over 15 .

# Run go vet
go vet ./...
```

### Project Structure

```
promptext-notes/
├── cmd/
│   └── promptext-notes/     # CLI entry point
│       └── main.go
├── internal/
│   ├── analyzer/            # Commit categorization
│   │   ├── analyzer.go
│   │   └── analyzer_test.go
│   ├── context/             # Code context extraction
│   │   ├── extractor.go
│   │   └── extractor_test.go
│   ├── generator/           # Release notes generation
│   │   ├── generator.go
│   │   └── generator_test.go
│   ├── git/                 # Git operations
│   │   ├── git.go
│   │   └── git_test.go
│   └── prompt/              # AI prompt generation
│       ├── prompt.go
│       └── prompt_test.go
├── .github/
│   └── workflows/
│       └── ci.yml           # CI/CD pipeline
├── go.mod
├── go.sum
├── README.md
├── LICENSE
└── .gitignore
```

## Examples

### Example 1: Quick Release Notes

```bash
$ promptext-notes --version v0.7.4

## [v0.7.4] - 2025-11-10

### Added
- Token budget support for code extraction
- File filtering by extension

### Fixed
- Panic when no git tags exist

### Statistics
- **Files changed**: 5
- **Commits**: 3
- **Context analyzed**: ~2,150 tokens

---
```

### Example 2: AI Prompt Generation

```bash
$ promptext-notes --version v0.7.4 --ai-prompt

# Release Notes Enhancement Request

Please generate comprehensive release notes for version v0.7.4

## Context

- **Version**: v0.7.4
- **Changes since**: v0.7.3
- **Commits analyzed**: 3
- **Files changed**: 5
- **Context extracted**: ~2,150 tokens

## Commit History

```
feat: add token budget support
fix: handle missing git tags
docs: update README examples
```

... (full prompt with code context)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes using conventional commits (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [promptext](https://github.com/1broseidon/promptext) - Token-aware code context extraction
- [Keep a Changelog](https://keepachangelog.com/) - Changelog format
- [Conventional Commits](https://www.conventionalcommits.org/) - Commit message convention

## Related Projects

- [promptext](https://github.com/1broseidon/promptext) - Extract code context with token awareness
- [conventional-changelog](https://github.com/conventional-changelog/conventional-changelog) - Generate changelogs from git metadata

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/1broseidon/promptext-notes/issues).
