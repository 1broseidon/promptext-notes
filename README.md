# Promptext-Notes

[![CI](https://github.com/1broseidon/promptext-notes/actions/workflows/ci.yml/badge.svg)](https://github.com/1broseidon/promptext-notes/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/1broseidon/promptext-notes?v=v0.7.0)](https://goreportcard.com/report/github.com/1broseidon/promptext-notes)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go-based CLI tool that generates intelligent, context-aware release notes by combining git history analysis with code context extraction using the [promptext](https://github.com/1broseidon/promptext) library.

## Features

- 📊 **Git History Analysis**: Automatically analyzes commits since the last tag
- 🔍 **Code Context Extraction**: Uses promptext to extract relevant code changes with token-aware analysis
- 📝 **Conventional Commits**: Categorizes changes by type (feat, fix, docs, breaking, etc.)
- 🤖 **Integrated AI Generation**: Generate AI-enhanced changelogs directly with `--generate` flag
- ✨ **2-Stage Polish Workflow**: Combine accurate discovery with customer-friendly polish for premium quality
- 🚫 **Auto-Exclude-Meta** (v0.8.0): Automatically excludes CI configs, CHANGELOG, README from AI context
- 🌐 **Multi-Provider Support**: Works with OpenRouter (200+ models), Anthropic, OpenAI, Cerebras, Groq, and local Ollama
- ⚙️ **YAML Configuration**: Customize behavior with `.promptext-notes.yml` config file
- 📋 **Keep a Changelog Format**: Produces standardized markdown output
- ⚡ **Fast & Lightweight**: Single binary with no runtime dependencies (except Git)
- 🔌 **Easy Integration**: Add to any repository with GitHub Actions ([See Guide](docs/USAGE.md))
- 🆓 **Free Options**: Use Cerebras, Groq, or local Ollama (no API cost)

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

### AI-Enhanced Release Notes (Integrated)

**NEW!** Generate AI-enhanced changelog directly with a single command:

```bash
# Using Cerebras (default with config file - free tier, best for large prompts)
export CEREBRAS_API_KEY="your-key-here"
promptext-notes --generate --version v1.0.0

# Or specify provider inline
promptext-notes --generate --provider groq --model llama-3.3-70b-versatile --version v1.0.0
```

The `--generate` flag will:
1. Analyze git history and extract code context
2. Send the comprehensive prompt to your AI provider
3. Return polished, production-ready release notes

**Legacy Method:** Generate a prompt to paste into an LLM manually:

```bash
promptext-notes --version v1.0.0 --ai-prompt > prompt.txt
```

Then paste the contents of `prompt.txt` into Claude, ChatGPT, or your preferred LLM.

### 2-Stage Polish Workflow ✨

**NEW!** Combine the accuracy of technical discovery models with the polish of customer-facing language models for premium quality changelogs.

**Quick Start:**
```bash
# Enable polish with CLI flag (uses config from .promptext-notes.yml)
promptext-notes --generate --polish --version v1.0.0

# Or configure in .promptext-notes.yml
```

**How it works:**
1. **Stage 1 (Discovery)**: Uses a model optimized for code understanding to analyze all changes
2. **Stage 2 (Polish)**: Refines the output into polished, customer-friendly language

**Recommended Setup (from benchmarks):**
```yaml
ai:
  provider: cerebras
  model: zai-glm-4.6  # Stage 1 (Discovery) - FREE, 10/10 accuracy

  polish:
    enabled: true  # Or use --polish CLI flag
    polish_model: "anthropic/claude-sonnet-4.5"  # Stage 2 - Premium polish
    polish_provider: "openrouter"  # Different provider for polish stage
    polish_api_key_env: "OPENROUTER_API_KEY"
    polish_max_tokens: 4000
    polish_temperature: 0.3
```

**Benefits:**
- ✅ GLM-4.6: Best free model (10/10 accuracy, catches all technical details)
- ✅ Claude Sonnet 4.5: Premium polish (8/10 quality, minimal hallucination)
- ✅ **Total cost: ~$0.004/run** (discovery is FREE, polish is cheap)
- ✅ Can mix FREE models (Cerebras) with paid polish (OpenRouter)
- ✅ Auto-exclude-meta (v0.8.0): Keeps changelogs focused on user-facing changes

**Cost Analysis:**
- Single-stage (GLM-4.6): $0 (FREE)
- 2-stage (GLM + Claude Sonnet): ~$0.004/run
- 2-stage (GLM + Haiku): ~$0.001/run (cheaper but less accurate)

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

## Configuration

You can configure promptext-notes using a YAML configuration file. Copy `.promptext-notes.example.yml` to `.promptext-notes.yml` and customize:

```yaml
version: "1"

ai:
  provider: cerebras      # cerebras, anthropic, openai, groq, openrouter, ollama
  model: zai-glm-4.6      # Best free model (10/10 accuracy)
  api_key_env: CEREBRAS_API_KEY
  max_tokens: 8000
  temperature: 0.3
  timeout: 30s

  polish:
    enabled: false  # Enable with --polish flag
    polish_model: "anthropic/claude-sonnet-4.5"
    polish_provider: "openrouter"
    polish_api_key_env: "OPENROUTER_API_KEY"

output:
  format: keepachangelog
  sections: [breaking, added, changed, fixed, docs]

filters:
  files:
    auto_exclude_meta: true  # NEW in v0.8.0 - excludes CI, CHANGELOG, README
    include: ["*.go", "*.md", "*.yml"]
    exclude: ["*_test.go", "vendor/*"]
```

See [CONFIGURATION.md](docs/CONFIGURATION.md) for full configuration options and [USAGE.md](docs/USAGE.md) for usage examples.

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--version` | string | "" | Version to generate notes for (e.g., v0.7.4) |
| `--since` | string | "" | Generate notes since this tag (auto-detects if empty) |
| `--output` | string | "" | Output file path (stdout if empty) |
| `--generate` | bool | false | **NEW!** Generate AI-enhanced changelog directly |
| `--polish` | bool | false | **NEW!** Enable 2-stage polish workflow (discovery + refinement) |
| `--provider` | string | "" | AI provider (anthropic, openai, cerebras, groq, openrouter, ollama) |
| `--model` | string | "" | AI model to use (overrides config) |
| `--exclude-files` | string | "" | Comma-separated files to exclude from AI context (e.g., CHANGELOG.md,README.md) |
| `--config` | string | ".promptext-notes.yml" | Configuration file path |
| `--quiet` | bool | false | Suppress progress messages |
| `--ai-prompt` | bool | false | Generate AI prompt only (legacy mode) |

## How It Works

1. **Git Analysis**: Retrieves changed files and commit messages since the last tag (or specified tag)
2. **Context Extraction**: Uses [promptext](https://github.com/1broseidon/promptext) to extract code context from changed files (.go, .md, .yml, .yaml)
3. **Categorization**: Parses commit messages using conventional commit format
4. **Generation**: Produces either:
   - **Basic Mode**: Keep a Changelog formatted release notes
   - **AI Mode**: Comprehensive prompt with full code context for LLM enhancement

## Automated Release Notes (AI-Enhanced)

This project includes an automated workflow that generates AI-enhanced release notes using multiple AI providers: **OpenAI**, **Anthropic**, **Cerebras**, or **Groq**.

> **📚 Want to use this in your own repository?**
> See the **[Complete Integration Guide](docs/USAGE.md)** for step-by-step instructions on adding automated AI-enhanced release notes to any project.

### How It Works

When you push a version tag (e.g., `v1.0.0`), the workflow automatically:

1. ✅ Builds the promptext-notes binary
2. 🔍 Analyzes git history and extracts code context
3. 🤖 Sends the prompt to your chosen AI provider for enhancement
4. 📝 Creates a GitHub release with polished notes
5. 📋 Updates CHANGELOG.md in the repository

### Supported AI Providers

| Provider | Default Model | Context Limit | Cost | Setup URL |
|----------|---------------|---------------|------|-----------|
| **OpenRouter** 🆕 | openai/gpt-4o-mini | Varies | 💰 Pay-as-you-go | [openrouter.ai/keys](https://openrouter.ai/keys) |
| **Cerebras** | zai-glm-4.6 | 65K tokens | ✅ Free | [cerebras.ai](https://cerebras.ai) |
| **Groq** | llama-3.3-70b-versatile | 32K tokens | ✅ Free | [console.groq.com](https://console.groq.com/keys) |
| **Ollama** | llama3.2 | Varies | ✅ Free (Local) | [ollama.com](https://ollama.com) |
| **OpenAI** | gpt-4o-mini | 128K tokens | 💰 $0.15/$0.60 per 1M | [platform.openai.com](https://platform.openai.com/api-keys) |
| **Anthropic** | claude-haiku-4-5 | 200K tokens | 💰 $0.80/$4.00 per 1M | [console.anthropic.com](https://console.anthropic.com/settings/keys) |

### Setup

1. **Get an API key** from your chosen provider (see Setup URL column above)

2. **Add API key(s) to GitHub Secrets**:
   - Go to your repository → **Settings** → **Secrets and variables** → **Actions**
   - Click **"New repository secret"**
   - Add one or more of these secrets:
     - `OPENROUTER_API_KEY` - For OpenRouter (access 200+ models through one API)
     - `CEREBRAS_API_KEY` - For Cerebras (free, best for large prompts)
     - `GROQ_API_KEY` - For Groq (free, good for smaller prompts)
     - `OPENAI_API_KEY` - For OpenAI
     - `ANTHROPIC_API_KEY` - For Anthropic

3. **(Optional) Configure models via GitHub Variables**:
   - Go to your repository → **Settings** → **Secrets and variables** → **Actions** → **Variables** tab
   - Add variables to customize models (otherwise defaults are used):
     - `OPENAI_MODEL` (default: `gpt-5-nano`)
     - `ANTHROPIC_MODEL` (default: `claude-haiku-4-5`)
     - `CEREBRAS_MODEL` (default: `gpt-oss-120b`)
     - `GROQ_MODEL` (default: `llama-3.3-70b-versatile`)

4. **Push a version tag**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

The workflow will automatically generate and publish AI-enhanced release notes using Cerebras (default, free) or your configured provider!

### Local CLI Usage (Recommended)

**NEW!** Use the integrated `--generate` flag for one-command AI-enhanced changelogs:

```bash
# Using Anthropic (create .promptext-notes.yml config first)
export ANTHROPIC_API_KEY="your-key-here"
promptext-notes --generate --version v1.0.0

# Or specify provider inline
export OPENAI_API_KEY="your-key-here"
promptext-notes --generate --provider openai --model gpt-4o-mini --version v1.0.0

# Using OpenRouter (access 200+ models through one API)
export OPENROUTER_API_KEY="your-key-here"
promptext-notes --generate --provider openrouter --model anthropic/claude-sonnet-4 --version v1.0.0

# Using Cerebras (default, free tier, best for large prompts)
export CEREBRAS_API_KEY="your-key-here"
promptext-notes --generate --version v1.0.0

# Using Groq (free tier, good for smaller prompts)
export GROQ_API_KEY="your-key-here"
promptext-notes --generate --provider groq --version v1.0.0

# Using Ollama (local, free, no API key needed!)
# First: ollama pull llama3.2
promptext-notes --generate --provider ollama --model llama3.2 --version v1.0.0
```

### Legacy Script Method

You can also use the shell script (will be deprecated in future versions):

```bash
# Using Cerebras (default)
export CEREBRAS_API_KEY="your-key-here"
./scripts/generate-release-notes.sh v1.0.0

# Using OpenAI
export OPENAI_API_KEY="your-key-here"
./scripts/generate-release-notes.sh v1.0.0 v0.9.0 openai
```

### Available Models by Provider

**OpenRouter** (pay-as-you-go, access 200+ models):
- `openai/gpt-4o-mini` (default) - Cost-effective OpenAI model
- `anthropic/claude-sonnet-4` - Latest Claude model
- `google/gemini-pro-1.5` - Google's Gemini
- `meta-llama/llama-3.3-70b-instruct` - Open source Llama
- `openai/gpt-4o` - Premium OpenAI model
- And 200+ more models! See [openrouter.ai/models](https://openrouter.ai/models)

**Cerebras** (free, ultra-fast, best for large prompts):
- `zai-glm-4.6` (default) - Multilingual support, best for large context
- `gpt-oss-120b` - 120B params, best free quality
- `llama-3.3-70b` - 70B params, good balance

**Groq** (free, fast, best for smaller prompts):
- `llama-3.3-70b-versatile` (default) - Best for general use, 32K context
- `mixtral-8x7b-32768` - Good for technical content
- `llama-3.1-70b-versatile` - Alternative option
- `moonshotai/kimi-k2-instruct-0905` - Kimi K2 model with 128K context (requires paid tier for large prompts)

**OpenAI** (paid, 2025 models):
- `gpt-5-nano` (default) - **Most economical** ($0.05/$0.40 per 1M tokens)
- `gpt-5-mini` - Good balance ($0.25/$2.00 per 1M tokens)
- `gpt-5` - **Best quality** ($1.25/$10 per 1M tokens)

**Anthropic** (paid, 2025 models):
- `claude-haiku-4-5` (default) - **Best value** ($1/$5 per 1M, 73.3% SWE-bench)
- `claude-sonnet-4-5` - **Best coding model** (frontier performance)
- `claude-opus-4-1` - Highest reasoning capability

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
- [staticcheck](https://staticcheck.dev/) (optional but recommended): `go install honnef.co/go/tools/cmd/staticcheck@latest`
- [gocyclo](https://github.com/fzipp/gocyclo) (optional but recommended): `go install github.com/fzipp/gocyclo/cmd/gocyclo@latest`

### Setup Pre-commit Hooks

Install Git hooks to automatically run quality checks before each commit:

```bash
./scripts/install-hooks.sh
```

This will run `go fmt`, `go vet`, `staticcheck`, `gocyclo`, and tests before allowing commits. To skip hooks for a specific commit:

```bash
git commit --no-verify
```

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
gocyclo -over 20 .

# Run go vet
go vet ./...
```

### Project Structure

```
promptext-notes/
├── cmd/
│   └── promptext-notes/           # CLI entry point
│       └── main.go
├── internal/
│   ├── ai/                        # AI provider integrations (NEW!)
│   │   ├── provider.go            # Provider interface
│   │   ├── anthropic.go           # Anthropic (Claude)
│   │   ├── openai.go              # OpenAI (GPT)
│   │   ├── cerebras.go            # Cerebras (free)
│   │   ├── groq.go                # Groq (free)
│   │   ├── ollama.go              # Local Ollama
│   │   └── retry.go               # Retry logic
│   ├── config/                    # Configuration (NEW!)
│   │   ├── config.go              # YAML config support
│   │   └── config_test.go
│   ├── workflow/                  # Orchestration (NEW!)
│   │   └── workflow.go            # End-to-end workflow
│   ├── analyzer/                  # Commit categorization
│   │   ├── analyzer.go
│   │   └── analyzer_test.go
│   ├── context/                   # Code context extraction
│   │   ├── extractor.go
│   │   └── extractor_test.go
│   ├── generator/                 # Release notes generation
│   │   ├── generator.go
│   │   └── generator_test.go
│   ├── git/                       # Git operations
│   │   ├── git.go
│   │   └── git_test.go
│   └── prompt/                    # AI prompt generation
│       ├── prompt.go
│       └── prompt_test.go
├── .github/
│   └── workflows/
│       ├── ci.yml                 # CI/CD pipeline
│       └── auto-docs.yml          # Automated release notes
├── scripts/
│   └── generate-release-notes.sh  # Shell script (legacy)
├── .promptext-notes.example.yml   # Example config (NEW!)
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

## Troubleshooting

### API Key Issues

**Problem**: `❌ OPENAI_API_KEY environment variable not set`

**Solution**:
- Make sure you've added the API key to GitHub Secrets (for CI) or set it as an environment variable (for local use)
- Check the secret name matches exactly: `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, `CEREBRAS_API_KEY`, or `GROQ_API_KEY`

### Invalid API Key

**Problem**: `❌ OpenAI API Error (invalid_api_key): Invalid API key`

**Solution**:
- Verify your API key is correct and hasn't expired
- For OpenAI/Groq: Key format is `sk-...` or similar
- For Anthropic: Key format is `sk-ant-...`
- For Cerebras: Check [cerebras.ai](https://cerebras.ai) for correct key format

### Model Not Found

**Problem**: `❌ Cerebras API Error: Model llama-3.1-70b does not exist`

**Solution**:
- Check the "Available Models by Provider" section above for valid model names
- Update the `CEREBRAS_MODEL` GitHub Variable or use a different model in the command
- Common mistake: `llama3.1-70b` (no dash) vs `llama-3.1-70b` (with dashes)

### Context Length Exceeded

**Problem**: `❌ API Error: Current length is 8950 while limit is 8192`

**Solution**:
- Your code changes are too large for the model's context window
- Switch to a provider with a larger context limit (see "Supported AI Providers" table)
- Recommended: Anthropic (200K), OpenAI (128K), or Cerebras (65K)

### Rate Limiting

**Problem**: `❌ API Error (429): Rate limit exceeded`

**Solution**:
- Wait a few minutes and try again
- Consider upgrading to a paid tier for higher rate limits
- Switch to a different provider (free tiers: Cerebras, Groq)

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
