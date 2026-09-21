package helper

import (
	"math/rand"
	"time"
)

func CalculateBackoff(attempt int) time.Duration {
	baseDelay := time.Second
	jitter := time.Duration(rand.Float64() * float64(baseDelay))

	return (baseDelay*time.Duration(1<<uint(attempt)) + jitter)
}
