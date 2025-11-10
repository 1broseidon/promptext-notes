# Documentation

Welcome to the promptext-notes documentation!

## For External Users

**[📚 Complete Integration Guide](USAGE.md)**

Step-by-step guide for adding automated AI-enhanced release notes to your repository:

- Quick start (5 minutes)
- Detailed setup instructions
- Provider-specific guides (OpenAI, Anthropic, Cerebras, Groq)
- Advanced configuration
- Troubleshooting
- Examples

## For Contributors

### Project Structure

```
promptext-notes/
├── cmd/
│   └── promptext-notes/     # CLI entry point
├── internal/
│   ├── analyzer/            # Commit categorization
│   ├── context/             # Code context extraction
│   ├── generator/           # Release notes generation
│   ├── git/                 # Git operations
│   └── prompt/              # AI prompt generation
├── scripts/
│   └── generate-release-notes.sh  # AI provider integration
├── .github/workflows/
│   ├── ci.yml               # CI/CD pipeline
│   └── auto-docs.yml        # Automated release notes
└── docs/
    ├── README.md            # This file
    └── USAGE.md             # External integration guide
```

### Development

See the main [README.md](../README.md) for development setup, testing, and contribution guidelines.

### Adding a New AI Provider

To add support for a new AI provider:

1. Add provider function in `scripts/generate-release-notes.sh`:
   ```bash
   call_newprovider() {
       local api_key="$1"
       local prompt="$2"
       # Implementation
   }
   ```

2. Add provider to validation function
3. Add to case statement in main logic
4. Update workflow `.github/workflows/auto-docs.yml`
5. Update documentation in `docs/USAGE.md`
6. Add secret/variable names to README

### Improving Prompt Quality

The AI prompt is generated in `internal/prompt/prompt.go`. Key sections:

- **Context metadata**: Version, commits, files changed
- **Critical Rules**: What to omit, what to focus on
- **Categorization rules**: Added vs Changed vs Fixed
- **Example format**: Shows expected output structure

When improving prompts:
1. Test with multiple AI providers
2. Verify no "noise" content (docs updates, stats, etc.)
3. Update tests in `internal/prompt/prompt_test.go`
4. Document changes in CHANGELOG

## Need Help?

- **Bug Reports**: [GitHub Issues](https://github.com/1broseidon/promptext-notes/issues)
- **Questions**: Open a [Discussion](https://github.com/1broseidon/promptext-notes/discussions)
- **Pull Requests**: Always welcome!
