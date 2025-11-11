package ai

import (
	"context"
	"fmt"
	"time"

	"github.com/1broseidon/promptext-notes/internal/config"
)

// RetryableFunc is a function that can be retried
type RetryableFunc func(ctx context.Context) error

// RetryWithBackoff retries a function with configurable backoff strategy
func RetryWithBackoff(ctx context.Context, cfg *config.Config, fn RetryableFunc) error {
	var lastErr error

	for attempt := 1; attempt <= cfg.AI.Retry.Attempts; attempt++ {
		// Try the operation
		err := fn(ctx)
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't sleep after the last attempt
		if attempt == cfg.AI.Retry.Attempts {
			break
		}

		// Calculate delay based on backoff strategy
		delay := calculateDelay(cfg.AI.Retry, attempt)

		// Wait with context cancellation support
		select {
		case <-time.After(delay):
			// Continue to next attempt
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		}
	}

	return fmt.Errorf("failed after %d attempts: %w", cfg.AI.Retry.Attempts, lastErr)
}

// calculateDelay calculates the delay before the next retry
func calculateDelay(retry config.RetryConfig, attempt int) time.Duration {
	switch retry.Backoff {
	case "exponential":
		// 2s, 4s, 8s, 16s...
		multiplier := 1 << uint(attempt-1) // 2^(attempt-1)
		return retry.InitialDelay * time.Duration(multiplier)

	case "linear":
		// 2s, 4s, 6s, 8s...
		return retry.InitialDelay * time.Duration(attempt)

	case "constant":
		// 2s, 2s, 2s, 2s...
		return retry.InitialDelay

	default:
		// Default to exponential
		multiplier := 1 << uint(attempt-1)
		return retry.InitialDelay * time.Duration(multiplier)
	}
}
