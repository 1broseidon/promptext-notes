# Using AI-Enhanced Release Notes in Your Repository

This guide shows you how to integrate automated AI-enhanced release notes into **your own repository** using the `promptext-notes` tool and GitHub Actions.

## Table of Contents

- [Quick Start](#quick-start)
- [Step-by-Step Setup](#step-by-step-setup)
- [Provider Setup Guides](#provider-setup-guides)
- [Advanced Configuration](#advanced-configuration)
- [Local Usage](#local-usage)
- [Troubleshooting](#troubleshooting)
- [Examples](#examples)

---

## Quick Start

**For the impatient**: Copy our workflow, add one secret, push a tag.

```bash
# 1. Copy the workflow to your repo
mkdir -p .github/workflows
curl -o .github/workflows/release-notes.yml \
  https://raw.githubusercontent.com/1broseidon/promptext-notes/main/.github/workflows/auto-docs.yml

# 2. Get a free API key from Cerebras: https://cerebras.ai
# 3. Add CEREBRAS_API_KEY to GitHub Secrets (Settings → Secrets and variables → Actions)
# 4. Push a tag
git tag v1.0.0
git push origin v1.0.0

# Done! Check your GitHub releases.
```

---

## Step-by-Step Setup

### Step 1: Copy the GitHub Action Workflow

Create `.github/workflows/release-notes.yml` in your repository:

```yaml
name: Auto-Generate Release Notes

on:
  push:
    tags:
      - 'v*'
  workflow_dispatch:
    inputs:
      version:
        description: 'Version tag (e.g., v1.0.0)'
        required: true
        type: string
      since_tag:
        description: 'Previous tag to compare from (optional)'
        required: false
        type: string
      api_provider:
        description: 'API provider to use'
        required: true
        type: choice
        options:
          - cerebras
          - openai
          - anthropic
          - groq
        default: cerebras

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
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Install dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y jq

      - name: Install promptext-notes
        run: |
          go install github.com/1broseidon/promptext-notes/cmd/promptext-notes@latest
          promptext-notes --help

      - name: Determine version and previous tag
        id: version
        run: |
          if [ "${{ github.event_name }}" = "push" ]; then
            VERSION="${GITHUB_REF#refs/tags/}"
            echo "version=$VERSION" >> $GITHUB_OUTPUT

            PREV_TAG=$(git describe --tags --abbrev=0 ${VERSION}^ 2>/dev/null || echo "")
            echo "since_tag=$PREV_TAG" >> $GITHUB_OUTPUT
          else
            VERSION="${{ inputs.version }}"
            SINCE_TAG="${{ inputs.since_tag }}"
            echo "version=$VERSION" >> $GITHUB_OUTPUT
            echo "since_tag=$SINCE_TAG" >> $GITHUB_OUTPUT
          fi

          echo "🏷️  Version: $VERSION"
          echo "📍 Previous tag: ${SINCE_TAG:-auto-detect}"

      - name: Select API provider
        id: api
        run: |
          if [ "${{ github.event_name }}" = "push" ]; then
            API_PROVIDER="cerebras"
          else
            API_PROVIDER="${{ inputs.api_provider }}"
          fi
          echo "provider=$API_PROVIDER" >> $GITHUB_OUTPUT
          echo "🤖 Using API provider: $API_PROVIDER"

      - name: Generate release notes with AI
        id: generate
        env:
          OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
          CEREBRAS_API_KEY: ${{ secrets.CEREBRAS_API_KEY }}
          GROQ_API_KEY: ${{ secrets.GROQ_API_KEY }}
          OPENAI_MODEL: ${{ vars.OPENAI_MODEL || 'gpt-5-nano' }}
          ANTHROPIC_MODEL: ${{ vars.ANTHROPIC_MODEL || 'claude-haiku-4-5' }}
          CEREBRAS_MODEL: ${{ vars.CEREBRAS_MODEL || 'gpt-oss-120b' }}
          GROQ_MODEL: ${{ vars.GROQ_MODEL || 'llama-3.3-70b-versatile' }}
          RELEASE_NOTES_FILE: release-notes.md
        run: |
          VERSION="${{ steps.version.outputs.version }}"
          SINCE_TAG="${{ steps.version.outputs.since_tag }}"
          API_PROVIDER="${{ steps.api.outputs.provider }}"

          curl -o generate-release-notes.sh \
            https://raw.githubusercontent.com/1broseidon/promptext-notes/main/scripts/generate-release-notes.sh
          chmod +x generate-release-notes.sh

          CMD="./generate-release-notes.sh $VERSION"
          if [ -n "$SINCE_TAG" ]; then
            CMD="$CMD $SINCE_TAG"
          fi
          CMD="$CMD $API_PROVIDER"

          echo "Running: $CMD"
          $CMD

          if [ -f release-notes.md ]; then
            echo "notes_file=release-notes.md" >> $GITHUB_OUTPUT
            echo "✅ Release notes generated successfully"
          else
            echo "❌ Failed to generate release notes"
            exit 1
          fi

      - name: Create GitHub Release
        if: github.event_name == 'push'
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
        if: github.event_name == 'push'
        run: |
          VERSION="${{ steps.version.outputs.version }}"

          cat release-notes.md > temp-changelog.md
          echo "" >> temp-changelog.md

          if [ -f CHANGELOG.md ]; then
            cat CHANGELOG.md >> temp-changelog.md
          fi

          mv temp-changelog.md CHANGELOG.md
          echo "📝 Updated CHANGELOG.md with release notes"

      - name: Commit CHANGELOG.md
        if: github.event_name == 'push'
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

**Note**: Adjust the `git push origin HEAD:main` line if your default branch is not `main`.

### Step 2: Choose Your AI Provider

Pick **one** provider to start with. Free options are best for getting started.

| Provider | Free? | Setup URL | Best For |
|----------|-------|-----------|----------|
| **Cerebras** | ✅ Yes | [cerebras.ai](https://cerebras.ai) | Fast, free, high-quality (recommended) |
| **Groq** | ✅ Yes | [console.groq.com](https://console.groq.com/keys) | Fast, free, good quality |
| **OpenAI** | 💰 Paid | [platform.openai.com](https://platform.openai.com/api-keys) | GPT-5 models, very economical with Nano |
| **Anthropic** | 💰 Paid | [console.anthropic.com](https://console.anthropic.com/settings/keys) | Best coding (Sonnet 4.5), great value (Haiku 4.5) |

### Step 3: Add API Key to GitHub Secrets

1. Get an API key from your chosen provider (see Setup URL above)
2. Go to your repository on GitHub
3. Navigate to **Settings** → **Secrets and variables** → **Actions**
4. Click **"New repository secret"**
5. Add the secret:
   - For Cerebras: Name = `CEREBRAS_API_KEY`, Value = your API key
   - For OpenAI: Name = `OPENAI_API_KEY`, Value = your API key
   - For Anthropic: Name = `ANTHROPIC_API_KEY`, Value = your API key
   - For Groq: Name = `GROQ_API_KEY`, Value = your API key

### Step 4: Test It!

```bash
# Create a version tag
git tag v1.0.0
git push origin v1.0.0
```

The workflow will automatically:
- Analyze your git history since the last tag
- Extract code context from changed files
- Generate professional release notes using AI
- Create a GitHub release
- Update your CHANGELOG.md

---

## Provider Setup Guides

### Cerebras (Recommended - Free)

**Why**: Ultra-fast inference, 65K token context, completely free.

**Setup**:
1. Visit [cerebras.ai](https://cerebras.ai)
2. Sign up for an account
3. Navigate to API Keys section
4. Create a new API key
5. Add to GitHub Secrets as `CEREBRAS_API_KEY`

**Default Model**: `gpt-oss-120b` (120B parameters)

**Available Models**:
- `gpt-oss-120b` - Best quality, recommended
- `llama-3.3-70b` - Good balance
- `zai-glm-4.6` - Multilingual support

### Groq (Free Alternative)

**Why**: Fast inference, free tier, good quality.

**Setup**:
1. Visit [console.groq.com](https://console.groq.com/keys)
2. Sign up for an account
3. Create an API key
4. Add to GitHub Secrets as `GROQ_API_KEY`

**Default Model**: `llama-3.3-70b-versatile`

**Available Models**:
- `llama-3.3-70b-versatile` - Best for general use
- `mixtral-8x7b-32768` - Good for technical content
- `llama-3.1-70b-versatile` - Alternative option

### OpenAI (Paid)

**Why**: GPT-5 models with excellent performance. GPT-5 Nano is extremely economical.

**Setup**:
1. Visit [platform.openai.com](https://platform.openai.com/api-keys)
2. Sign up and add payment method
3. Create an API key
4. Add to GitHub Secrets as `OPENAI_API_KEY`

**Default Model**: `gpt-5-nano` (272K context)

**Available Models** (2025):
- `gpt-5-nano` - **Most economical** ($0.05 input / $0.40 output per 1M tokens)
- `gpt-5-mini` - Good balance ($0.25 / $2.00 per 1M)
- `gpt-5` - **Best quality** ($1.25 / $10 per 1M)

**Cost**: ~$0.02-0.20 per release with Nano (extremely economical!)

### Anthropic (Paid)

**Why**: Best coding model (Sonnet 4.5). Haiku 4.5 offers excellent value at $1/$5.

**Setup**:
1. Visit [console.anthropic.com](https://console.anthropic.com/settings/keys)
2. Sign up and add payment method
3. Create an API key
4. Add to GitHub Secrets as `ANTHROPIC_API_KEY`

**Default Model**: `claude-haiku-4-5` (200K context)

**Available Models** (2025):
- `claude-haiku-4-5` - **Best value** ($1 / $5 per 1M, 73.3% SWE-bench Verified)
- `claude-sonnet-4-5` - **Best coding model in the world** (frontier performance)
- `claude-opus-4-1` - Highest reasoning capability

**Cost**: ~$0.05-0.25 per release with Haiku (great value!)

---

## Advanced Configuration

### Custom Models

Override default models using GitHub Variables:

1. Go to **Settings** → **Secrets and variables** → **Actions** → **Variables** tab
2. Add variables:
   - `CEREBRAS_MODEL` - Default: `gpt-oss-120b`
   - `OPENAI_MODEL` - Default: `gpt-5-nano`
   - `ANTHROPIC_MODEL` - Default: `claude-haiku-4-5`
   - `GROQ_MODEL` - Default: `llama-3.3-70b-versatile`

### Multiple Providers

You can add API keys for multiple providers and switch between them:

```bash
# Add all secrets
CEREBRAS_API_KEY=xxx
OPENAI_API_KEY=yyy
ANTHROPIC_API_KEY=zzz
GROQ_API_KEY=www

# Switch providers manually via workflow_dispatch
# Go to Actions → Auto-Generate Release Notes → Run workflow
# Select your preferred provider from dropdown
```

### Manual Trigger

Trigger the workflow manually:

1. Go to **Actions** tab
2. Select **"Auto-Generate Release Notes"**
3. Click **"Run workflow"**
4. Fill in:
   - Version: `v1.0.0`
   - Previous tag: `v0.9.0` (optional)
   - Provider: Choose from dropdown
5. Click **"Run workflow"**

### CHANGELOG.md Only (No Release)

If you only want CHANGELOG updates without creating GitHub releases, modify the workflow:

```yaml
# Remove or comment out this step:
- name: Create GitHub Release
  if: github.event_name == 'push'
  uses: softprops/action-gh-release@v1
  # ... rest of step
```

### Change Default Provider

Edit the workflow and change:

```yaml
- name: Select API provider
  id: api
  run: |
    if [ "${{ github.event_name }}" = "push" ]; then
      API_PROVIDER="openai"  # Change this line
    else
      API_PROVIDER="${{ inputs.api_provider }}"
    fi
```

---

## Local Usage

Use the tool locally without GitHub Actions:

### Install

```bash
go install github.com/1broseidon/promptext-notes/cmd/promptext-notes@latest
```

### Generate Basic Release Notes

```bash
promptext-notes --version v1.0.0
```

### Generate AI-Enhanced Release Notes

```bash
# Set API key
export CEREBRAS_API_KEY="your-key-here"

# Download the script
curl -O https://raw.githubusercontent.com/1broseidon/promptext-notes/main/scripts/generate-release-notes.sh
chmod +x generate-release-notes.sh

# Generate notes
./generate-release-notes.sh v1.0.0

# With custom provider
./generate-release-notes.sh v1.0.0 v0.9.0 openai

# Save to file
RELEASE_NOTES_FILE=notes.md ./generate-release-notes.sh v1.0.0
```

---

## Troubleshooting

### ❌ API Key Not Found

**Error**: `CEREBRAS_API_KEY environment variable not set`

**Solution**:
- Verify the secret name matches exactly (case-sensitive)
- Make sure you added it to **Secrets**, not Variables
- Secret names must match: `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, `CEREBRAS_API_KEY`, `GROQ_API_KEY`

### ❌ Wrong API Key

**Error**: `Wrong API Key` or `Invalid API key`

**Solution**:
- Double-check you copied the API key correctly (no extra spaces)
- Verify the key hasn't expired
- Make sure you're using the correct provider's key

### ❌ Model Not Found

**Error**: `Model llama-3.1-70b does not exist`

**Solution**:
- Check the [Available Models](#provider-setup-guides) section
- Common mistake: `llama3.1-70b` (no dashes) vs `llama-3.1-70b` (with dashes)
- Update `CEREBRAS_MODEL` variable with correct model name

### ❌ Context Length Exceeded

**Error**: `Current length is 8950 while limit is 8192`

**Solution**:
- Your code changes are too large for the model's context window
- Switch to a provider with larger context:
  - OpenAI GPT-5: 272K tokens
  - Anthropic: 200K tokens
  - Cerebras: 65K tokens
  - Groq: 32K tokens

### ❌ Workflow Fails to Push CHANGELOG

**Error**: `failed to push some refs`

**Solution**:
- Make sure the workflow has `contents: write` permission (already in provided YAML)
- Check branch protection rules don't block bot commits
- Verify default branch name is correct (change `HEAD:main` if needed)

### ❌ promptext-notes Install Fails

**Error**: `go: github.com/1broseidon/promptext-notes@latest: not found`

**Solution**:
- Make sure Go 1.22+ is installed: `go version`
- Check your internet connection
- Try with explicit version: `go install github.com/1broseidon/promptext-notes/cmd/promptext-notes@v0.4.0`

---

## Examples

### Example 1: Simple Setup (Cerebras)

```yaml
# .github/workflows/release-notes.yml
# ... (use the workflow from Step 1)
```

```bash
# Add secret
# Settings → Secrets → New secret
# Name: CEREBRAS_API_KEY
# Value: <your-key>

# Push tag
git tag v1.0.0
git push origin v1.0.0
```

### Example 2: OpenAI with Best Model

```yaml
# Add to GitHub Variables to use best model:
# OPENAI_MODEL = gpt-5
```

```bash
# Add secret
# Name: OPENAI_API_KEY
# Value: sk-...

# Edit workflow to use openai by default:
API_PROVIDER="openai"  # Instead of "cerebras"

# Push tag
git tag v2.0.0
git push origin v2.0.0
```

### Example 3: Multiple Providers for Flexibility

```bash
# Add ALL secrets:
CEREBRAS_API_KEY=xxx
OPENAI_API_KEY=yyy
ANTHROPIC_API_KEY=zzz
GROQ_API_KEY=www

# Use Cerebras by default (automatic on tag push)
git tag v1.0.0
git push origin v1.0.0

# Manually trigger with OpenAI for important releases
# Actions → Run workflow → Select "openai"
```

### Example 4: Monorepo with Multiple Changelogs

Modify the workflow to update specific changelogs:

```yaml
- name: Update CHANGELOG.md
  run: |
    # Update service-specific changelog
    cat release-notes.md > services/api/CHANGELOG.md
    git add services/api/CHANGELOG.md

    # Also update root changelog
    cat release-notes.md > CHANGELOG.md
    git add CHANGELOG.md
```

---

## Need Help?

- **Issues**: [GitHub Issues](https://github.com/1broseidon/promptext-notes/issues)
- **Examples**: Check out this repository's [releases](https://github.com/1broseidon/promptext-notes/releases)
- **Source Code**: [promptext-notes](https://github.com/1broseidon/promptext-notes)

---

## License

This tool is MIT licensed. Use it freely in your open source and commercial projects.
