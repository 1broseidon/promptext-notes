#!/bin/bash

# Script to generate AI-enhanced release notes using Cerebras API
# Usage: ./generate-release-notes.sh <version> [since-tag] [api-provider]
#
# Environment variables required:
#   CEREBRAS_API_KEY - API key for Cerebras
#   GROK_API_KEY - API key for Grok (alternative)

set -e

VERSION="${1:-}"
SINCE_TAG="${2:-}"
API_PROVIDER="${3:-cerebras}"  # cerebras or grok

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version> [since-tag] [api-provider]"
    echo "Example: $0 v1.0.0"
    echo "Example: $0 v1.0.0 v0.9.0"
    echo "Example: $0 v1.0.0 v0.9.0 grok"
    exit 1
fi

# Check if promptext-notes binary exists
if ! command -v promptext-notes &> /dev/null; then
    echo "❌ promptext-notes binary not found. Please install it first:"
    echo "   go install github.com/1broseidon/promptext-notes/cmd/promptext-notes@latest"
    exit 1
fi

echo "🤖 Generating AI prompt for $VERSION..."

# Generate the AI prompt
PROMPT_ARGS="--version $VERSION --ai-prompt"
if [ -n "$SINCE_TAG" ]; then
    PROMPT_ARGS="$PROMPT_ARGS --since $SINCE_TAG"
fi

AI_PROMPT=$(promptext-notes $PROMPT_ARGS 2>/dev/null)

if [ -z "$AI_PROMPT" ]; then
    echo "❌ Failed to generate AI prompt"
    exit 1
fi

echo "✅ AI prompt generated (~$(echo "$AI_PROMPT" | wc -c) characters)"

# Function to call Cerebras API
call_cerebras() {
    local api_key="$1"
    local prompt="$2"

    echo "🧠 Calling Cerebras API (llama3.1-70b)..."

    # Create JSON payload
    local json_payload=$(jq -n \
        --arg prompt "$prompt" \
        '{
            model: "llama3.1-70b",
            stream: false,
            max_tokens: 4096,
            temperature: 0.6,
            top_p: 0.95,
            messages: [
                {
                    role: "user",
                    content: $prompt
                }
            ]
        }')

    # Call API
    local response=$(curl -s --location 'https://api.cerebras.ai/v1/chat/completions' \
        --header 'Content-Type: application/json' \
        --header "Authorization: Bearer ${api_key}" \
        --data "$json_payload")

    # Check for errors
    if echo "$response" | jq -e '.error' > /dev/null 2>&1; then
        echo "❌ API Error: $(echo "$response" | jq -r '.error.message')"
        return 1
    fi

    # Extract content
    echo "$response" | jq -r '.choices[0].message.content // empty'
}

# Function to call Grok API (xAI)
call_grok() {
    local api_key="$1"
    local prompt="$2"

    echo "🧠 Calling Grok API (grok-beta)..."

    # Create JSON payload
    local json_payload=$(jq -n \
        --arg prompt "$prompt" \
        '{
            model: "grok-beta",
            stream: false,
            max_tokens: 4096,
            temperature: 0.6,
            top_p: 0.95,
            messages: [
                {
                    role: "user",
                    content: $prompt
                }
            ]
        }')

    # Call API
    local response=$(curl -s --location 'https://api.x.ai/v1/chat/completions' \
        --header 'Content-Type: application/json' \
        --header "Authorization: Bearer ${api_key}" \
        --data "$json_payload")

    # Check for errors
    if echo "$response" | jq -e '.error' > /dev/null 2>&1; then
        echo "❌ API Error: $(echo "$response" | jq -r '.error.message')"
        return 1
    fi

    # Extract content
    echo "$response" | jq -r '.choices[0].message.content // empty'
}

# Call the appropriate API
RELEASE_NOTES=""
case "$API_PROVIDER" in
    cerebras)
        if [ -z "$CEREBRAS_API_KEY" ]; then
            echo "❌ CEREBRAS_API_KEY environment variable not set"
            exit 1
        fi
        RELEASE_NOTES=$(call_cerebras "$CEREBRAS_API_KEY" "$AI_PROMPT")
        ;;
    grok)
        if [ -z "$GROK_API_KEY" ]; then
            echo "❌ GROK_API_KEY environment variable not set"
            exit 1
        fi
        RELEASE_NOTES=$(call_grok "$GROK_API_KEY" "$AI_PROMPT")
        ;;
    *)
        echo "❌ Unknown API provider: $API_PROVIDER"
        echo "Supported providers: cerebras, grok"
        exit 1
        ;;
esac

if [ -z "$RELEASE_NOTES" ]; then
    echo "❌ Failed to generate release notes from API"
    exit 1
fi

# Output the release notes
echo ""
echo "✅ Release notes generated successfully!"
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "$RELEASE_NOTES"
echo "═══════════════════════════════════════════════════════════════"

# Save to file if RELEASE_NOTES_FILE is set
if [ -n "$RELEASE_NOTES_FILE" ]; then
    echo "$RELEASE_NOTES" > "$RELEASE_NOTES_FILE"
    echo ""
    echo "📝 Saved to: $RELEASE_NOTES_FILE"
fi
