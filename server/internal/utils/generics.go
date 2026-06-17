package utils

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func StripJSONFences(s string) string {
	t := strings.TrimSpace(s)
	t = strings.TrimPrefix(t, "```json")
	t = strings.TrimPrefix(t, "```")
	t = strings.TrimSuffix(t, "```")
	return strings.TrimSpace(t)
}

// ExecuteWithRetry runs operation up to maxRetries times. Successive attempts
// wait baseDelay, 2·baseDelay, 4·baseDelay … (pass 0 to retry immediately).
// Stops early if ctx is cancelled.
func ExecuteWithRetry[T any](
	ctx context.Context,
	operationName string,
	maxRetries int,
	baseDelay time.Duration,
	operation func(ctx context.Context) (T, error),
) (T, error) {
	var fallback T
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		result, err := operation(ctx)
		if err == nil {
			return result, nil
		}
		lastErr = err

		if attempt < maxRetries && baseDelay > 0 {
			wait := baseDelay * time.Duration(1<<(attempt-1))
			select {
			case <-time.After(wait):
			case <-ctx.Done():
				return fallback, ctx.Err()
			}
		}
	}

	return fallback, fmt.Errorf(
		"%s failed after %d attempts: %w",
		operationName,
		maxRetries,
		lastErr,
	)
}

func AppendUnique(s []string, v string) []string {
	for _, x := range s {
		if x == v {
			return s
		}
	}
	return append(s, v)
}
