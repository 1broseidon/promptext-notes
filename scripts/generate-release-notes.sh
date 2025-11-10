#!/bin/bash

# Script to generate AI-enhanced release notes using multiple AI providers
# Usage: ./generate-release-notes.sh <version> [since-tag] [api-provider] [model]
#
# Supported providers: openai, anthropic, cerebras, groq
#
# Environment variables:
#   API Keys (required for chosen provider):
#     OPENAI_API_KEY - API key for OpenAI
#     ANTHROPIC_API_KEY - API key for Anthropic
#     CEREBRAS_API_KEY - API key for Cerebras
#     GROQ_API_KEY - API key for Groq
#
#   Models (optional, uses defaults if not specified):
#     OPENAI_MODEL - OpenAI model (default: gpt-5-nano)
#     ANTHROPIC_MODEL - Anthropic model (default: claude-haiku-4-5)
#     CEREBRAS_MODEL - Cerebras model (default: gpt-oss-120b)
#     GROQ_MODEL - Groq model (default: llama-3.3-70b-versatile)

set -e

VERSION="${1:-}"
SINCE_TAG="${2:-}"
API_PROVIDER="${3:-cerebras}"
MODEL_OVERRIDE="${4:-}"

# Set default models
OPENAI_MODEL="${OPENAI_MODEL:-gpt-5-nano}"
ANTHROPIC_MODEL="${ANTHROPIC_MODEL:-claude-haiku-4-5}"
CEREBRAS_MODEL="${CEREBRAS_MODEL:-gpt-oss-120b}"
GROQ_MODEL="${GROQ_MODEL:-llama-3.3-70b-versatile}"

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version> [since-tag] [api-provider] [model]" >&2
    echo "" >&2
    echo "Examples:" >&2
    echo "  $0 v1.0.0" >&2
    echo "  $0 v1.0.0 v0.9.0 openai" >&2
    echo "  $0 v1.0.0 v0.9.0 anthropic claude-opus-4" >&2
    echo "" >&2
    echo "Supported providers: openai, anthropic, cerebras, groq" >&2
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

# Validation function to check provider configuration
validate_provider_config() {
    local provider="$1"
    local required_key=""
    local setup_url=""

    case "$provider" in
        openai)
            required_key="OPENAI_API_KEY"
            setup_url="https://platform.openai.com/api-keys"
            ;;
        anthropic)
            required_key="ANTHROPIC_API_KEY"
            setup_url="https://console.anthropic.com/settings/keys"
            ;;
        cerebras)
            required_key="CEREBRAS_API_KEY"
            setup_url="https://cerebras.ai"
            ;;
        groq)
            required_key="GROQ_API_KEY"
            setup_url="https://console.groq.com/keys"
            ;;
        *)
            echo "❌ Unknown provider: $provider" >&2
            echo "Supported providers: openai, anthropic, cerebras, groq" >&2
            return 1
            ;;
    esac

    # Check if the required API key is set
    if [ -z "${!required_key}" ]; then
        echo "❌ ${required_key} environment variable not set" >&2
        echo "" >&2
        echo "To use the $provider provider:" >&2
        echo "  1. Get an API key from: $setup_url" >&2
        echo "  2. Set the environment variable: export ${required_key}=your-key-here" >&2
        echo "  3. Or add it to GitHub Secrets if running in CI" >&2
        echo "" >&2
        return 1
    fi

    echo "✅ Provider validation passed for $provider" >&2
    return 0
}

# Function to call OpenAI API
call_openai() {
    local api_key="$1"
    local prompt="$2"
    local model="${OPENAI_MODEL}"

    if [ -n "$MODEL_OVERRIDE" ]; then
        model="$MODEL_OVERRIDE"
    fi

    echo "🧠 Calling OpenAI API ($model)..." >&2

    # Create JSON payload
    local json_payload=$(jq -n \
        --arg prompt "$prompt" \
        --arg model "$model" \
        '{
            model: $model,
            messages: [
                {
                    role: "user",
                    content: $prompt
                }
            ],
            temperature: 0.6,
            max_tokens: 4096,
            top_p: 0.95
        }')

    # Call API
    local response=$(curl -s --location 'https://api.openai.com/v1/chat/completions' \
        --header 'Content-Type: application/json' \
        --header "Authorization: Bearer ${api_key}" \
        --data "$json_payload")

    # Check for errors
    if echo "$response" | jq -e '.error' > /dev/null 2>&1; then
        local error_msg=$(echo "$response" | jq -r '.error.message')
        local error_type=$(echo "$response" | jq -r '.error.type // "unknown"')
        echo "❌ OpenAI API Error ($error_type): $error_msg" >&2

        if [ "$error_type" = "invalid_api_key" ]; then
            echo "💡 Check your OPENAI_API_KEY is correct" >&2
        fi
        return 1
    fi

    # Check if response is empty
    if [ -z "$response" ]; then
        echo "❌ Empty response from OpenAI API" >&2
        return 1
    fi

    # Extract content
    local content=$(echo "$response" | jq -r '.choices[0].message.content // empty')

    if [ -z "$content" ]; then
        echo "❌ No content in OpenAI API response" >&2
        echo "Full response: $response" >&2
        return 1
    fi

    echo "$content"
}

# Function to call Anthropic API
call_anthropic() {
    local api_key="$1"
    local prompt="$2"
    local model="${ANTHROPIC_MODEL}"

    if [ -n "$MODEL_OVERRIDE" ]; then
        model="$MODEL_OVERRIDE"
    fi

    echo "🧠 Calling Anthropic API ($model)..." >&2

    # Create JSON payload (Anthropic uses slightly different structure)
    local json_payload=$(jq -n \
        --arg prompt "$prompt" \
        --arg model "$model" \
        '{
            model: $model,
            max_tokens: 4096,
            temperature: 0.6,
            messages: [
                {
                    role: "user",
                    content: $prompt
                }
            ]
        }')

    # Call API (note: Anthropic uses x-api-key header)
    local response=$(curl -s --location 'https://api.anthropic.com/v1/messages' \
        --header 'Content-Type: application/json' \
        --header "x-api-key: ${api_key}" \
        --header 'anthropic-version: 2023-06-01' \
        --data "$json_payload")

    # Check for errors
    if echo "$response" | jq -e '.error' > /dev/null 2>&1; then
        local error_msg=$(echo "$response" | jq -r '.error.message')
        local error_type=$(echo "$response" | jq -r '.error.type // "unknown"')
        echo "❌ Anthropic API Error ($error_type): $error_msg" >&2

        if [ "$error_type" = "authentication_error" ]; then
            echo "💡 Check your ANTHROPIC_API_KEY is correct" >&2
        fi
        return 1
    fi

    # Check if response is empty
    if [ -z "$response" ]; then
        echo "❌ Empty response from Anthropic API" >&2
        return 1
    fi

    # Extract content (Anthropic returns content in different structure)
    local content=$(echo "$response" | jq -r '.content[0].text // empty')

    if [ -z "$content" ]; then
        echo "❌ No content in Anthropic API response" >&2
        echo "Full response: $response" >&2
        return 1
    fi

    echo "$content"
}

# Function to call Cerebras API
call_cerebras() {
    local api_key="$1"
    local prompt="$2"
    local model="${CEREBRAS_MODEL}"

    if [ -n "$MODEL_OVERRIDE" ]; then
        model="$MODEL_OVERRIDE"
    fi

    echo "🧠 Calling Cerebras API ($model)..." >&2

    # Create JSON payload
    local json_payload=$(jq -n \
        --arg prompt "$prompt" \
        --arg model "$model" \
        '{
            model: $model,
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

    # Check for errors (Cerebras uses both .error and direct error fields)
    if echo "$response" | jq -e '.error' > /dev/null 2>&1; then
        local error_msg=$(echo "$response" | jq -r '.error.message')
        local error_code=$(echo "$response" | jq -r '.error.code // "unknown"')
        echo "❌ Cerebras API Error ($error_code): $error_msg" >&2

        if [[ "$error_msg" == *"not found"* ]] || [[ "$error_msg" == *"does not exist"* ]]; then
            echo "💡 Model not found. Available models: gpt-oss-120b, llama-3.3-70b, zai-glm-4.6" >&2
        fi
        return 1
    elif echo "$response" | jq -e '.code' > /dev/null 2>&1; then
        # Cerebras sometimes returns errors without .error wrapper
        local error_msg=$(echo "$response" | jq -r '.message')
        local error_code=$(echo "$response" | jq -r '.code')
        echo "❌ Cerebras API Error ($error_code): $error_msg" >&2

        if [ "$error_code" = "wrong_api_key" ] || [ "$error_code" = "invalid_api_key" ]; then
            echo "💡 Check your CEREBRAS_API_KEY is correct" >&2
        fi
        return 1
    fi

    # Check if response is empty
    if [ -z "$response" ]; then
        echo "❌ Empty response from Cerebras API" >&2
        return 1
    fi

    # Extract content
    local content=$(echo "$response" | jq -r '.choices[0].message.content // empty')

    if [ -z "$content" ]; then
        echo "❌ No content in Cerebras API response" >&2
        echo "Full response: $response" >&2
        return 1
    fi

    echo "$content"
}

# Function to call Groq API
call_groq() {
    local api_key="$1"
    local prompt="$2"
    local model="${GROQ_MODEL}"

    if [ -n "$MODEL_OVERRIDE" ]; then
        model="$MODEL_OVERRIDE"
    fi

    echo "🧠 Calling Groq API ($model)..." >&2

    # Create JSON payload
    local json_payload=$(jq -n \
        --arg prompt "$prompt" \
        --arg model "$model" \
        '{
            model: $model,
            messages: [
                {
                    role: "user",
                    content: $prompt
                }
            ],
            temperature: 0.6,
            max_tokens: 4096,
            top_p: 0.95
        }')

    # Call API
    local response=$(curl -s --location 'https://api.groq.com/openai/v1/chat/completions' \
        --header 'Content-Type: application/json' \
        --header "Authorization: Bearer ${api_key}" \
        --data "$json_payload")

    # Check for errors
    if echo "$response" | jq -e '.error' > /dev/null 2>&1; then
        local error_msg=$(echo "$response" | jq -r '.error.message')
        local error_type=$(echo "$response" | jq -r '.error.type // "unknown"')
        echo "❌ Groq API Error ($error_type): $error_msg" >&2

        if [[ "$error_msg" == *"model"* ]] && [[ "$error_msg" == *"not found"* ]]; then
            echo "💡 Model not found. Common models: llama-3.3-70b-versatile, mixtral-8x7b-32768" >&2
        fi
        return 1
    fi

    # Check if response is empty
    if [ -z "$response" ]; then
        echo "❌ Empty response from Groq API" >&2
        return 1
    fi

    # Extract content
    local content=$(echo "$response" | jq -r '.choices[0].message.content // empty')

    if [ -z "$content" ]; then
        echo "❌ No content in Groq API response" >&2
        echo "Full response: $response" >&2
        return 1
    fi

    echo "$content"
}

# Validate provider configuration before making API call
if ! validate_provider_config "$API_PROVIDER"; then
    exit 1
fi

# Call the appropriate API
RELEASE_NOTES=""
case "$API_PROVIDER" in
    openai)
        RELEASE_NOTES=$(call_openai "$OPENAI_API_KEY" "$AI_PROMPT")
        ;;
    anthropic)
        RELEASE_NOTES=$(call_anthropic "$ANTHROPIC_API_KEY" "$AI_PROMPT")
        ;;
    cerebras)
        RELEASE_NOTES=$(call_cerebras "$CEREBRAS_API_KEY" "$AI_PROMPT")
        ;;
    groq)
        RELEASE_NOTES=$(call_groq "$GROQ_API_KEY" "$AI_PROMPT")
        ;;
    *)
        echo "❌ Unknown API provider: $API_PROVIDER" >&2
        echo "Supported providers: openai, anthropic, cerebras, groq" >&2
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
