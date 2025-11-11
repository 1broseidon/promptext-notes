#!/bin/bash

# Script to install Git hooks for pre-commit checks

set -e

HOOKS_DIR=".git/hooks"
HOOK_FILE="$HOOKS_DIR/pre-commit"

if [ ! -d ".git" ]; then
    echo "❌ This script must be run from the repository root"
    exit 1
fi

echo "📦 Installing Git hooks..."

# Create pre-commit hook
cat > "$HOOK_FILE" << 'EOF'
#!/bin/bash

set -e

echo "🔍 Running pre-commit checks..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    exit 1
fi

# Run go fmt
echo "  → Running go fmt..."
UNFORMATTED=$(gofmt -l .)
if [ -n "$UNFORMATTED" ]; then
    echo "❌ The following files need formatting:"
    echo "$UNFORMATTED"
    echo ""
    echo "Run: go fmt ./..."
    exit 1
fi

# Run go vet
echo "  → Running go vet..."
if ! go vet ./...; then
    echo "❌ go vet failed"
    exit 1
fi

# Run staticcheck if installed
if command -v staticcheck &> /dev/null; then
    echo "  → Running staticcheck..."
    if ! staticcheck ./...; then
        echo "❌ staticcheck failed"
        exit 1
    fi
else
    echo "  ⚠️  staticcheck not installed (skipping)"
    echo "     Install: go install honnef.co/go/tools/cmd/staticcheck@latest"
fi

# Run gocyclo if installed
if command -v gocyclo &> /dev/null; then
    echo "  → Running gocyclo..."
    if ! gocyclo -over 20 .; then
        echo "❌ gocyclo failed (complexity > 20)"
        exit 1
    fi
else
    echo "  ⚠️  gocyclo not installed (skipping)"
    echo "     Install: go install github.com/fzipp/gocyclo/cmd/gocyclo@latest"
fi

# Run tests
echo "  → Running tests..."
if ! go test ./... -short; then
    echo "❌ Tests failed"
    exit 1
fi

echo "✅ All pre-commit checks passed!"
EOF

# Make executable
chmod +x "$HOOK_FILE"

echo "✅ Pre-commit hook installed successfully!"
echo ""
echo "To skip hooks for a specific commit, use: git commit --no-verify"
