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
    echo "Usage: $0 <version> [since-tag] [api-provider]" >&2
    echo "Example: $0 v1.0.0" >&2
    echo "Example: $0 v1.0.0 v0.9.0" >&2
    echo "Example: $0 v1.0.0 v0.9.0 grok" >&2
    exit 1
fi

# Check if promptext-notes binary exists
if ! command -v promptext-notes &> /dev/null; then
    echo "❌ promptext-notes binary not found. Please install it first:" >&2
    echo "   go install github.com/1broseidon/promptext-notes/cmd/promptext-notes@latest" >&2
    exit 1
fi

echo "🤖 Generating AI prompt for $VERSION..." >&2

# Generate the AI prompt
PROMPT_ARGS="--version $VERSION --ai-prompt"
if [ -n "$SINCE_TAG" ]; then
    PROMPT_ARGS="$PROMPT_ARGS --since $SINCE_TAG"
fi

AI_PROMPT=$(promptext-notes $PROMPT_ARGS 2>/dev/null)

if [ -z "$AI_PROMPT" ]; then
    echo "❌ Failed to generate AI prompt" >&2
    exit 1
fi

echo "✅ AI prompt generated (~$(echo "$AI_PROMPT" | wc -c) characters)" >&2

# Function to call Cerebras API
call_cerebras() {
    local api_key="$1"
    local prompt="$2"

    echo "🧠 Calling Cerebras API (llama-3.3-70b)..." >&2

    # Create JSON payload
    local json_payload=$(jq -n \
        --arg prompt "$prompt" \
        '{
            model: "llama-3.3-70b",
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
        echo "❌ API Error: $(echo "$response" | jq -r '.error.message')" >&2
        echo "Full response: $response" >&2
        return 1
    fi

    # Check if response is empty
    if [ -z "$response" ]; then
        echo "❌ Empty response from API" >&2
        return 1
    fi

    # Extract content
    local content=$(echo "$response" | jq -r '.choices[0].message.content // empty')

    if [ -z "$content" ]; then
        echo "❌ No content in API response" >&2
        echo "Full response: $response" >&2
        return 1
    fi

    echo "$content"
}

# Function to call Grok API (xAI)
call_grok() {
    local api_key="$1"
    local prompt="$2"

    echo "🧠 Calling Grok API (grok-beta)..." >&2

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
        echo "❌ API Error: $(echo "$response" | jq -r '.error.message')" >&2
        echo "Full response: $response" >&2
        return 1
    fi

    # Check if response is empty
    if [ -z "$response" ]; then
        echo "❌ Empty response from API" >&2
        return 1
    fi

    # Extract content
    local content=$(echo "$response" | jq -r '.choices[0].message.content // empty')

    if [ -z "$content" ]; then
        echo "❌ No content in API response" >&2
        echo "Full response: $response" >&2
        return 1
    fi

    echo "$content"
}

# Call the appropriate API
RELEASE_NOTES=""
case "$API_PROVIDER" in
    cerebras)
        if [ -z "$CEREBRAS_API_KEY" ]; then
            echo "❌ CEREBRAS_API_KEY environment variable not set" >&2
            exit 1
        fi
        RELEASE_NOTES=$(call_cerebras "$CEREBRAS_API_KEY" "$AI_PROMPT")
        ;;
    grok)
        if [ -z "$GROK_API_KEY" ]; then
            echo "❌ GROK_API_KEY environment variable not set" >&2
            exit 1
        fi
        RELEASE_NOTES=$(call_grok "$GROK_API_KEY" "$AI_PROMPT")
        ;;
    *)
        echo "❌ Unknown API provider: $API_PROVIDER" >&2
        echo "Supported providers: cerebras, grok" >&2
        exit 1
        ;;
esac

if [ -z "$RELEASE_NOTES" ]; then
    echo "❌ Failed to generate release notes from API" >&2
    exit 1
fi

# Output the release notes to stdout (this is what gets captured)
echo "$RELEASE_NOTES"

# Status message to stderr
echo "" >&2
echo "✅ Release notes generated successfully!" >&2

# Save to file if RELEASE_NOTES_FILE is set
if [ -n "$RELEASE_NOTES_FILE" ]; then
    echo "$RELEASE_NOTES" > "$RELEASE_NOTES_FILE"
    echo "" >&2
    echo "📝 Saved to: $RELEASE_NOTES_FILE" >&2
fi
