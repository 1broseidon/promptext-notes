# Configuration Reference

Complete reference for `.promptext-notes.yml` configuration file.

## Table of Contents

- [Overview](#overview)
- [AI Configuration](#ai-configuration)
- [Output Configuration](#output-configuration)
- [Filters Configuration](#filters-configuration)
- [Complete Example](#complete-example)
- [Use Cases](#use-cases)

---

## Overview

Configuration file: `.promptext-notes.yml` (YAML format)

**Location:** Repository root

**Validation:** Validated on load with helpful error messages

**Priority:** Config file < CLI flags (CLI flags override config)

---

## AI Configuration

### Basic AI Settings

```yaml
ai:
  # Provider (required)
  provider: cerebras  # cerebras, openai, anthropic, groq, openrouter, ollama

  # Model (required)
  model: zai-glm-4.6  # Exact API model name

  # API Key Environment Variable (required, except ollama)
  api_key_env: CEREBRAS_API_KEY

  # Max Tokens (optional, default: 8000)
  max_tokens: 8000

  # Temperature (optional, default: 0.3, range: 0.0-1.0)
  temperature: 0.3

  # Timeout (optional, default: 30s)
  timeout: 30s
```

### Supported Providers

| Provider | Models | API Key Env | Free? |
|----------|--------|-------------|-------|
| `cerebras` | `zai-glm-4.6`, `llama-3.3-70b` | `CEREBRAS_API_KEY` | ✅ Yes |
| `groq` | `llama-3.3-70b-versatile`, `mixtral-8x7b-32768` | `GROQ_API_KEY` | ✅ Yes |
| `openai` | `gpt-4o`, `gpt-4o-mini` | `OPENAI_API_KEY` | ❌ Paid |
| `anthropic` | `claude-sonnet-4.5`, `claude-haiku-4.5`, `claude-opus-4-20250514` | `ANTHROPIC_API_KEY` | ❌ Paid |
| `openrouter` | 100+ models (e.g., `anthropic/claude-sonnet-4.5`, `google/gemini-2.5-flash`) | `OPENROUTER_API_KEY` | ❌ Paid |
| `ollama` | Any local model | N/A | ✅ Free (local) |

### Recommended Models

**Best Free (Discovery):**
```yaml
ai:
  provider: cerebras
  model: zai-glm-4.6  # 10/10 accuracy, completely free
```

**Best Paid (Polish):**
```yaml
ai:
  polish:
    polish_provider: openrouter
    polish_model: "anthropic/claude-sonnet-4.5"  # 8/10 accuracy, ~$0.004/run
```

### Retry Configuration

```yaml
ai:
  retry:
    # Number of retry attempts (default: 3)
    attempts: 3

    # Backoff strategy: exponential, linear, constant (default: exponential)
    backoff: exponential

    # Initial delay before first retry (default: 2s)
    initial_delay: 2s
```

**Backoff Strategies:**
- `exponential`: 2s, 4s, 8s, 16s... (recommended)
- `linear`: 2s, 4s, 6s, 8s...
- `constant`: 2s, 2s, 2s, 2s...

### Custom Provider Options

```yaml
ai:
  custom:
    # Anthropic API version (optional)
    anthropic_version: "2023-06-01"

    # OpenRouter optional headers
    http_referer: "https://your-app-url.com"
    x_title: "Your App Name"

    # Ollama base URL (default: http://localhost:11434)
    ollama_url: "http://localhost:11434"
```

### 2-Stage Polish Workflow

**Overview:** Stage 1 (Discovery) + Stage 2 (Polish) = Premium quality release notes

```yaml
ai:
  # Stage 1: Discovery (main AI config above)
  provider: cerebras
  model: zai-glm-4.6

  # Stage 2: Polish (optional refinement)
  polish:
    # Enable polish workflow (default: false)
    # Can also enable via --polish CLI flag
    enabled: false

    # Polish model (required if enabled)
    polish_model: "anthropic/claude-sonnet-4.5"

    # Polish provider (optional, defaults to main provider)
    polish_provider: "openrouter"

    # Polish API key env (optional, auto-detected from provider)
    polish_api_key_env: "OPENROUTER_API_KEY"

    # Polish max tokens (default: 4000)
    polish_max_tokens: 4000

    # Polish temperature (default: 0.3)
    polish_temperature: 0.3

    # Custom polish prompt (optional, uses default if not specified)
    polish_prompt: ""
```

**How It Works:**
1. Discovery model analyzes code changes → generates technical changelog
2. Polish model refines language → produces user-friendly release notes

**Cost:**
- Stage 1 (cerebras/zai-glm-4.6): **FREE**
- Stage 2 (claude-sonnet-4.5): **~$0.004/run**
- **Total:** ~$0.004/run

**When to Use:**
- ✅ Production releases
- ✅ Public-facing changelogs
- ✅ When quality > cost
- ❌ Internal releases (discovery only is fine)
- ❌ Frequent pre-releases (free discovery only)

---

## Output Configuration

```yaml
output:
  # Format: keepachangelog or conventional (default: keepachangelog)
  format: keepachangelog

  # Sections to include (default: all)
  sections:
    - breaking      # ⚠️ BREAKING CHANGES
    - added         # ✨ New features
    - changed       # 🔄 Changes to existing features
    - fixed         # 🐛 Bug fixes
    - deprecated    # 🗑️ Deprecated features
    - removed       # ❌ Removed features
    - security      # 🔒 Security fixes
    - docs          # 📚 Documentation updates

  # Custom template path (optional)
  template: ./templates/custom-changelog.tmpl
```

### Format Options

**keepachangelog** (default):
```markdown
## [v1.0.0] - 2025-11-12

### ⚠️ BREAKING CHANGES
- Removed deprecated API

### Added
- New feature X

### Changed
- Updated feature Y

### Fixed
- Bug fix Z
```

**conventional**:
```markdown
# Release v1.0.0 (2025-11-12)

**Features**
- New feature X

**Changes**
- Updated feature Y

**Fixes**
- Bug fix Z
```

### Section Filtering

Only include specific sections:

```yaml
output:
  sections:
    - added
    - fixed
    # Excludes: changed, deprecated, removed, security, docs
```

---

## Filters Configuration

### File Filtering

```yaml
filters:
  files:
    # Auto-exclude meta files (default: true) - NEW in v0.8.0
    auto_exclude_meta: true

    # Include patterns (glob format)
    include:
      - "*.go"
      - "*.md"
      - "*.yml"
      - "*.yaml"
      - "*.json"
      - "*.js"
      - "*.ts"
      - "*.tsx"
      - "*.py"

    # Exclude patterns (glob format)
    exclude:
      - "*_test.go"         # Test files
      - "vendor/*"          # Vendor dependencies
      - "node_modules/*"    # Node dependencies
      - ".git/*"            # Git metadata
      - "dist/*"            # Build output
      - "build/*"           # Build artifacts
```

### Auto-Exclude-Meta (v0.8.0+)

**When `auto_exclude_meta: true` (default):**

Automatically excludes these files from AI context:
- `CHANGELOG.md` - Existing changelog
- `README.md` - Project documentation
- `.github/**` - GitHub Actions, workflows, issue templates
- `.vscode/**` - VS Code configs
- `.idea/**` - JetBrains configs
- `*.example.*` - Example files
- `.promptext-notes*.yml` - Tool config files
- `**/.gitignore` - Git ignore files
- `**/.*ignore` - All ignore files

**Purpose:** Prevents internal tooling changes from appearing in user-facing release notes.

**Example Impact:**

Without auto-exclude-meta:
```markdown
### Changed
- Updated GitHub Actions workflow to use new model
- Fixed typo in README
- Improved configuration defaults
```

With auto-exclude-meta:
```markdown
### Changed
- **Improved configuration defaults** - Better error handling and validation
```

**Disable:**
```yaml
filters:
  files:
    auto_exclude_meta: false  # Include all files
```

### Commit Filtering

```yaml
filters:
  commits:
    # Exclude commits from these authors
    exclude_authors:
      - "dependabot[bot]"
      - "renovate[bot]"
      - "github-actions[bot]"

    # Exclude commits matching these regex patterns
    exclude_patterns:
      - "^Merge pull request"    # PR merge commits
      - "^Merge branch"           # Branch merge commits
      - "^chore\\(deps\\):"       # Dependency updates
```

---

## Complete Example

Full production configuration:

```yaml
version: "1"

# AI Provider Configuration
ai:
  # Discovery stage (FREE)
  provider: cerebras
  model: zai-glm-4.6
  api_key_env: CEREBRAS_API_KEY
  max_tokens: 8000
  temperature: 0.3
  timeout: 30s

  # Retry configuration
  retry:
    attempts: 3
    backoff: exponential
    initial_delay: 2s

  # Custom options (optional)
  custom: {}

  # Polish workflow (premium quality)
  polish:
    enabled: false  # Enable via --polish flag or set true here
    polish_model: "anthropic/claude-sonnet-4.5"
    polish_provider: "openrouter"
    polish_api_key_env: "OPENROUTER_API_KEY"
    polish_max_tokens: 4000
    polish_temperature: 0.3

# Output Configuration
output:
  format: keepachangelog
  sections:
    - breaking
    - added
    - changed
    - fixed
    - deprecated
    - security
    - docs

# File & Commit Filtering
filters:
  files:
    # Auto-exclude meta files (CI, configs, CHANGELOG, README)
    auto_exclude_meta: true

    include:
      - "*.go"
      - "*.md"
      - "*.yml"
      - "*.yaml"
      - "*.json"

    exclude:
      - "*_test.go"
      - "vendor/*"
      - "node_modules/*"
      - ".git/*"

  commits:
    exclude_authors:
      - "dependabot[bot]"
      - "renovate[bot]"
      - "github-actions[bot]"

    exclude_patterns:
      - "^Merge pull request"
      - "^Merge branch"
```

---

## Use Cases

### Use Case 1: Startup (Free)

**Goal:** Free, automated release notes

```yaml
version: "1"

ai:
  provider: cerebras
  model: zai-glm-4.6
  api_key_env: CEREBRAS_API_KEY

filters:
  files:
    auto_exclude_meta: true
```

**Cost:** $0 (FREE)
**Quality:** 10/10 accuracy

---

### Use Case 2: Production (Premium)

**Goal:** Highest quality release notes for public releases

```yaml
version: "1"

ai:
  provider: cerebras
  model: zai-glm-4.6
  api_key_env: CEREBRAS_API_KEY

  polish:
    enabled: true  # Or use --polish flag
    polish_model: "anthropic/claude-sonnet-4.5"
    polish_provider: "openrouter"
    polish_api_key_env: "OPENROUTER_API_KEY"

filters:
  files:
    auto_exclude_meta: true
```

**Cost:** ~$0.004/run
**Quality:** 8/10 with premium polish

---

### Use Case 3: Public Library

**Goal:** Show all changes, including CI/tooling improvements

```yaml
version: "1"

ai:
  provider: cerebras
  model: zai-glm-4.6
  api_key_env: CEREBRAS_API_KEY

filters:
  files:
    auto_exclude_meta: false  # Include CI/config changes

output:
  sections:
    - breaking
    - added
    - changed
    - fixed
    - docs  # Show documentation updates
```

---

### Use Case 4: Internal Tools

**Goal:** Fast, simple changelogs for internal use

```yaml
version: "1"

ai:
  provider: cerebras
  model: llama-3.3-70b  # Fast model
  api_key_env: CEREBRAS_API_KEY
  max_tokens: 4000      # Smaller context
  temperature: 0.1      # More deterministic

filters:
  files:
    auto_exclude_meta: true
  commits:
    exclude_patterns:
      - "^WIP:"           # Work in progress commits
      - "^temp:"          # Temporary commits
```

---

### Use Case 5: Monorepo

**Goal:** Separate changelogs per service

```yaml
version: "1"

ai:
  provider: cerebras
  model: zai-glm-4.6
  api_key_env: CEREBRAS_API_KEY

filters:
  files:
    include:
      - "services/api/**"  # Only API service files
    auto_exclude_meta: true
```

---

## Environment Variables

Set these in your shell or CI/CD:

```bash
# Required (choose one)
export CEREBRAS_API_KEY="your-key"
export OPENAI_API_KEY="your-key"
export ANTHROPIC_API_KEY="your-key"
export GROQ_API_KEY="your-key"
export OPENROUTER_API_KEY="your-key"

# Optional (for 2-stage polish if using different providers)
export CEREBRAS_API_KEY="discovery-key"
export OPENROUTER_API_KEY="polish-key"
```

---

## CLI Overrides

CLI flags override config file:

```bash
# Override provider/model
promptext-notes --generate --version v1.0.0 \
  --provider openai \
  --model gpt-4o-mini

# Override exclude files
promptext-notes --generate --version v1.0.0 \
  --exclude-files "*.test.js,*.spec.ts"

# Enable polish (overrides config)
promptext-notes --generate --version v1.0.0 --polish
```

---

## Validation

Config is validated on load. Common errors:

### Invalid Provider

```
Error: invalid AI provider: invalid (supported: anthropic, openai, cerebras, groq, openrouter, ollama)
```

**Fix:** Use a supported provider name

### Missing API Key

```
Error: API key not found in environment variable: CEREBRAS_API_KEY
```

**Fix:**
```bash
export CEREBRAS_API_KEY="your-key"
```

### Invalid Temperature

```
Error: temperature must be between 0 and 1, got: 1.5
```

**Fix:** Set temperature between 0.0 and 1.0

### Invalid Max Tokens

```
Error: max_tokens must be positive, got: -100
```

**Fix:** Set max_tokens > 0

---

## Next Steps

- **Usage Examples:** [USAGE.md](USAGE.md)
- **GitHub Actions Setup:** [USAGE.md#github-actions](USAGE.md#github-actions)
- **Git Hooks:** [USAGE.md#git-hooks](USAGE.md#git-hooks)
