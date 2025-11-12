# Usage Guide

Complete guide to using `promptext-notes` for automated, AI-enhanced release notes.

## Table of Contents

- [Quick Start](#quick-start)
- [Installation](#installation)
- [Configuration](#configuration)
- [Local Usage](#local-usage)
- [GitHub Actions](#github-actions)
- [Git Hooks](#git-hooks)
- [Examples](#examples)

---

## Quick Start

**TL;DR:** Install → Configure → Use

```bash
# 1. Install
go install github.com/1broseidon/promptext-notes/cmd/promptext-notes@latest

# 2. Configure (optional but recommended)
cat > .promptext-notes.yml <<EOF
version: "1"
ai:
  provider: cerebras
  model: zai-glm-4.6
  api_key_env: CEREBRAS_API_KEY
EOF

# 3. Use
export CEREBRAS_API_KEY="your-key-here"
promptext-notes --generate --version v1.0.0
```

---

## Installation

### Prerequisites

- **Go 1.22+**: `go version`
- **Git repository**: `git status`
- **API key**: From Cerebras, OpenAI, Anthropic, Groq, or OpenRouter

### Install via Go

```bash
go install github.com/1broseidon/promptext-notes/cmd/promptext-notes@latest
```

### Verify Installation

```bash
promptext-notes --help
```

### Get an API Key

**Recommended: Cerebras (Free)**
1. Visit [cerebras.ai](https://cerebras.ai)
2. Sign up and get API key
3. `export CEREBRAS_API_KEY="your-key"`

**Alternatives:**
- **OpenRouter** (access 100+ models): [openrouter.ai](https://openrouter.ai)
- **Anthropic** (best coding): [console.anthropic.com](https://console.anthropic.com/settings/keys)
- **OpenAI**: [platform.openai.com](https://platform.openai.com/api-keys)
- **Groq** (free, fast): [console.groq.com](https://console.groq.com/keys)

---

## Configuration

### Create Config File

Create `.promptext-notes.yml` in your repository root:

```yaml
version: "1"

# AI Provider Configuration
ai:
  provider: cerebras           # cerebras, openai, anthropic, groq, openrouter, ollama
  model: zai-glm-4.6           # Best free model (10/10 accuracy)
  api_key_env: CEREBRAS_API_KEY
  max_tokens: 8000
  temperature: 0.3
  timeout: 30s

  retry:
    attempts: 3
    backoff: exponential
    initial_delay: 2s

  # 2-Stage Polish Workflow (optional)
  polish:
    enabled: false             # Enable with --polish flag
    polish_model: "anthropic/claude-sonnet-4.5"
    polish_provider: "openrouter"
    polish_api_key_env: "OPENROUTER_API_KEY"
    polish_max_tokens: 4000
    polish_temperature: 0.3

# Output Configuration
output:
  format: keepachangelog       # keepachangelog or conventional
  sections:
    - breaking
    - added
    - changed
    - fixed
    - docs

# File Filtering
filters:
  files:
    # Auto-exclude meta files (CI, configs, CHANGELOG, README)
    auto_exclude_meta: true    # NEW in v0.8.0 (prevents internal changes in release notes)

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

  commits:
    exclude_authors:
      - "dependabot[bot]"
      - "renovate[bot]"
      - "github-actions[bot]"

    exclude_patterns:
      - "^Merge pull request"
      - "^Merge branch"
```

See [CONFIGURATION.md](CONFIGURATION.md) for full reference.

---

## Local Usage

### Basic Release Notes

Generate release notes without AI:

```bash
promptext-notes --version v1.0.0
```

Output:
```markdown
# Release Notes for v1.0.0

## Changed (3)
- Update configuration system
- Improve error handling
- Refactor prompt generation

## Fixed (1)
- Fix API key environment variable override
```

### AI-Enhanced Release Notes

Generate with AI for better quality:

```bash
# Set API key
export CEREBRAS_API_KEY="your-key-here"

# Generate
promptext-notes --generate --version v1.0.0

# With specific version range
promptext-notes --generate --version v1.0.0 --since v0.9.0

# Save to file
promptext-notes --generate --version v1.0.0 --output CHANGELOG.md
```

Output (AI-enhanced):
```markdown
## [v1.0.0] - 2025-11-12

### Changed
- **Enhanced configuration system** - Improved config file handling with better defaults and validation
- **Error handling improvements** - Added retry logic with exponential backoff for more reliable AI requests
- **Refactored prompt generation** - Streamlined prompt structure for better AI comprehension

### Fixed
- **API key environment override** - CLI provider override now correctly updates corresponding API key
```

### 2-Stage Polish Workflow

For premium quality release notes:

```bash
# Set API keys (discovery + polish providers)
export CEREBRAS_API_KEY="your-cerebras-key"      # Stage 1: Discovery
export OPENROUTER_API_KEY="your-openrouter-key"  # Stage 2: Polish

# Generate with polish
promptext-notes --generate --version v1.0.0 --polish
```

**How it works:**
1. **Stage 1 (Discovery):** Uses `zai-glm-4.6` to analyze code changes (FREE, 10/10 accuracy)
2. **Stage 2 (Polish):** Uses `claude-sonnet-4.5` to polish language (~$0.004/run)

### Manual Mode (AI Prompt Only)

Generate just the AI prompt without calling the API:

```bash
promptext-notes --ai-prompt --version v1.0.0 > prompt.txt

# Use with any AI tool
cat prompt.txt | pbcopy  # Copy to clipboard
# Paste into ChatGPT, Claude.ai, etc.
```

### Override Provider/Model

```bash
# Use different provider
promptext-notes --generate --version v1.0.0 --provider openai --model gpt-4o-mini

# Use OpenRouter (access 100+ models)
export OPENROUTER_API_KEY="your-key"
promptext-notes --generate --version v1.0.0 \
  --provider openrouter \
  --model anthropic/claude-sonnet-4.5
```

---

## GitHub Actions

### Setup

1. **Create workflow file** `.github/workflows/release-notes.yml`:

```yaml
name: Auto-Generate Release Notes

on:
  push:
    tags:
      - 'v*'

jobs:
  generate-release-notes:
    name: Generate AI Release Notes
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Full history for git operations

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true

      - name: Install promptext-notes from source
        run: |
          go install ./cmd/promptext-notes
          echo "$HOME/go/bin" >> $GITHUB_PATH
          promptext-notes --help

      - name: Determine version and previous tag
        id: version
        run: |
          VERSION="${GITHUB_REF#refs/tags/}"
          echo "version=$VERSION" >> $GITHUB_OUTPUT

          PREV_TAG=$(git describe --tags --abbrev=0 ${VERSION}^ 2>/dev/null || echo "")
          echo "since_tag=$PREV_TAG" >> $GITHUB_OUTPUT

          echo "🏷️  Version: $VERSION"
          echo "📍 Previous tag: ${PREV_TAG:-auto-detect}"

      - name: Generate release notes with AI
        id: generate
        env:
          CEREBRAS_API_KEY: ${{ secrets.CEREBRAS_API_KEY }}
          OPENROUTER_API_KEY: ${{ secrets.OPENROUTER_API_KEY }}
        run: |
          VERSION="${{ steps.version.outputs.version }}"
          SINCE_TAG="${{ steps.version.outputs.since_tag }}"

          # Build command using config file
          CMD="promptext-notes --generate --output release-notes.md --version $VERSION"

          if [ -n "$SINCE_TAG" ]; then
            CMD="$CMD --since $SINCE_TAG"
          fi

          # Enable 2-stage polish workflow
          CMD="$CMD --polish"
          echo "✨ 2-stage polish workflow enabled"
          echo "🚫 Auto-exclude-meta enabled: CI configs, CHANGELOG, README excluded"

          echo "🚀 Running: $CMD"
          $CMD

          if [ -f release-notes.md ]; then
            echo "✅ Release notes generated successfully"
          else
            echo "❌ Failed to generate release notes"
            exit 1
          fi

      - name: Create GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          tag_name: ${{ steps.version.outputs.version }}
          name: Release ${{ steps.version.outputs.version }}
          body_path: release-notes.md
          draft: false
          prerelease: false
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Update CHANGELOG.md
        run: |
          VERSION="${{ steps.version.outputs.version }}"

          # Check if version already exists
          if [ -f CHANGELOG.md ] && grep -q "^## \\[$VERSION\\]" CHANGELOG.md; then
            echo "⚠️  Version $VERSION already exists in CHANGELOG.md, skipping"
            exit 0
          fi

          # Strip release notes header and prepend to CHANGELOG
          if grep -q "^# Release Notes for" release-notes.md; then
            sed -n '/^## \\[/,$p' release-notes.md > release-notes-clean.md
          else
            cp release-notes.md release-notes-clean.md
          fi

          cat release-notes-clean.md > temp-changelog.md
          echo "" >> temp-changelog.md
          echo "---" >> temp-changelog.md
          echo "" >> temp-changelog.md

          if [ -f CHANGELOG.md ]; then
            sed -n '/^## \\[/,$p' CHANGELOG.md >> temp-changelog.md
          fi

          {
            echo "# Changelog"
            echo ""
            echo "All notable changes to this project will be documented in this file."
            echo ""
            echo "The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)."
            echo ""
            echo "---"
            echo ""
            cat temp-changelog.md
          } > CHANGELOG.md

          rm -f release-notes-clean.md temp-changelog.md
          echo "📝 Updated CHANGELOG.md"

      - name: Commit CHANGELOG.md
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"

          git add CHANGELOG.md

          if git diff --staged --quiet; then
            echo "No changes to commit"
          else
            git commit -m "docs: update CHANGELOG for ${{ steps.version.outputs.version }}"
            git push origin HEAD:main
            echo "✅ Pushed CHANGELOG.md to main branch"
          fi

      - name: Upload release notes artifact
        uses: actions/upload-artifact@v4
        with:
          name: release-notes-${{ steps.version.outputs.version }}
          path: release-notes.md
          retention-days: 90
```

2. **Add API keys to GitHub Secrets:**
   - Go to **Settings → Secrets and variables → Actions**
   - Add `CEREBRAS_API_KEY` (required)
   - Add `OPENROUTER_API_KEY` (optional, for polish workflow)

3. **Commit config file** (`.promptext-notes.yml`)

4. **Push a tag:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

5. **Check the release:** Go to your GitHub releases page

### Workflow Features

✅ **Automatic trigger** on version tags (`v*`)
✅ **2-stage polish** (discovery + refinement)
✅ **Auto-exclude-meta** (no CI/config changes in release notes)
✅ **CHANGELOG update** with deduplication
✅ **GitHub release** creation
✅ **Artifact upload** for archival

---

## Git Hooks

### Pre-Tag Hook

Generate release notes before creating a tag:

Create `.git/hooks/pre-tag` (or use a tool like [Husky](https://typicode.github.io/husky/)):

```bash
#!/bin/bash
# .git/hooks/pre-tag

TAG_NAME="$1"

echo "🚀 Generating release notes for $TAG_NAME..."

export CEREBRAS_API_KEY="your-key-here"  # Or read from env

promptext-notes --generate --version "$TAG_NAME" --output CHANGELOG-preview.md

echo "✅ Preview saved to CHANGELOG-preview.md"
echo ""
cat CHANGELOG-preview.md
echo ""
echo "Press Enter to continue, Ctrl+C to abort"
read
```

Usage:
```bash
chmod +x .git/hooks/pre-tag
git tag -a v1.0.0 -m "v1.0.0"  # Hook runs automatically
```

### Post-Commit Hook

Auto-categorize commits (without AI):

```bash
#!/bin/bash
# .git/hooks/post-commit

LAST_COMMIT=$(git log -1 --pretty=%B)

if [[ $LAST_COMMIT =~ ^(feat|fix|docs|refactor|perf|test|chore): ]]; then
  echo "✅ Conventional commit detected: $LAST_COMMIT"
else
  echo "⚠️  Non-conventional commit. Use: feat|fix|docs|refactor|perf|test|chore: message"
fi
```

---

## Examples

### Example 1: Local Development

```bash
# During development
git tag v0.1.0-alpha

# Generate preview
promptext-notes --generate --version v0.1.0-alpha

# Review, then push
git push origin v0.1.0-alpha
```

### Example 2: Production Release

```bash
# 1. Config file with polish workflow
cat .promptext-notes.yml
# ai:
#   provider: cerebras
#   model: zai-glm-4.6
#   polish:
#     enabled: false
#     polish_model: "anthropic/claude-sonnet-4.5"
#     polish_provider: "openrouter"

# 2. Generate with polish
export CEREBRAS_API_KEY="..."
export OPENROUTER_API_KEY="..."

promptext-notes --generate --version v1.0.0 --polish

# 3. Commit and tag
git add CHANGELOG.md
git commit -m "docs: v1.0.0 release notes"
git tag v1.0.0
git push origin main v1.0.0
```

### Example 3: Manual Review Workflow

```bash
# 1. Generate draft
promptext-notes --generate --version v1.0.0 --output draft.md

# 2. Review and edit
vim draft.md

# 3. Commit manually
cat draft.md >> CHANGELOG.md
git add CHANGELOG.md
git commit -m "docs: add v1.0.0 release notes"
```

### Example 4: Compare Providers

```bash
# Test different models
for provider in cerebras openai anthropic; do
  echo "=== Testing $provider ==="
  promptext-notes --generate --version v1.0.0 \
    --provider $provider \
    --output "release-notes-$provider.md"
done

# Compare outputs
diff release-notes-{cerebras,openai}.md
```

### Example 5: Monorepo

```bash
# Generate notes for specific subdirectory
cd services/api
promptext-notes --generate --version v1.0.0 --output CHANGELOG.md

cd ../web
promptext-notes --generate --version v1.0.0 --output CHANGELOG.md
```

---

## Tips & Best Practices

### 1. Use Config File

**Don't:**
```bash
promptext-notes --generate --version v1.0.0 \
  --provider cerebras \
  --model zai-glm-4.6 \
  --exclude-files "CHANGELOG.md,README.md,.github/**"
```

**Do:**
```yaml
# .promptext-notes.yml
ai:
  provider: cerebras
  model: zai-glm-4.6
filters:
  files:
    auto_exclude_meta: true
```

```bash
promptext-notes --generate --version v1.0.0
```

### 2. Enable Auto-Exclude-Meta

Keeps changelogs focused on user-facing changes:

```yaml
filters:
  files:
    auto_exclude_meta: true  # Excludes CI, configs, CHANGELOG, README
```

Before:
```
### Changed
- Updated GitHub Actions workflow
- Fixed typo in README
- Improved configuration defaults
```

After:
```
### Changed
- **Improved configuration defaults** - Better error handling and validation
```

### 3. Use 2-Stage Polish for Production

Free discovery + cheap polish = best quality:

```yaml
ai:
  provider: cerebras
  model: zai-glm-4.6        # FREE
  polish:
    enabled: true
    polish_model: "anthropic/claude-sonnet-4.5"  # ~$0.004/run
    polish_provider: "openrouter"
```

### 4. Version Ranges

Always specify `--since` for better context:

```bash
# Good
promptext-notes --generate --version v1.0.0 --since v0.9.0

# Less accurate (auto-detects previous tag)
promptext-notes --generate --version v1.0.0
```

### 5. Test Locally First

Before pushing tags, test locally:

```bash
# Dry run
promptext-notes --generate --version v1.0.0-test

# Review
cat release-notes.md

# Clean up test tag
git tag -d v1.0.0-test
```

---

## Troubleshooting

### API Key Not Found

```
Error: API key not found in environment variable: CEREBRAS_API_KEY
```

**Fix:**
```bash
export CEREBRAS_API_KEY="your-key-here"
# Or add to ~/.bashrc or ~/.zshrc
```

### Config File Not Loaded

```bash
# Check if file exists
ls -la .promptext-notes.yml

# Validate YAML syntax
cat .promptext-notes.yml | grep -v "^#" | head -20
```

### No Changes Detected

```
Error: no commits found in range v1.0.0..v0.9.0
```

**Fix:**
```bash
# Check tags exist
git tag -l

# Check commit range
git log v0.9.0..v1.0.0 --oneline
```

---

## Next Steps

- **Configuration Reference:** [CONFIGURATION.md](CONFIGURATION.md)
- **Source Code:** [github.com/1broseidon/promptext-notes](https://github.com/1broseidon/promptext-notes)
- **Report Issues:** [GitHub Issues](https://github.com/1broseidon/promptext-notes/issues)
