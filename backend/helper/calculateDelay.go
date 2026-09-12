package helper

import "time"

func CalculateBackoff(attempt int) time.Duration {
	baseDelay := time.Second

	return baseDelay * time.Duration(1<<uint(attempt))
}
