package utils

import (
	"math/rand"
	"time"
)

// BackOff computes the next back-off duration
type BackOff func(currentRetryCount int) time.Duration

// ExponentialBackOff computes an exponential back-off
func ExponentialBackOff(minTimeout time.Duration) BackOff {
	return func(currentRetryCount int) time.Duration {
		jitter := rand.Float64()
		jitterMax := 200 * time.Millisecond
		if currentRetryCount < 1 {
			currentRetryCount = 1
		}
		strategy := 1 << (currentRetryCount - 1)
		backoff := (float64(strategy) * float64(minTimeout.Nanoseconds())) + (jitter * float64(jitterMax.Nanoseconds()))

		return time.Duration(backoff)
	}
}

// LinearBackOff computes a linear back-off
func LinearBackOff(minTimeout time.Duration) BackOff {
	return func(currentRetryCount int) time.Duration {
		jitter := rand.Float64()
		jitterMax := 200 * time.Millisecond
		if currentRetryCount < 1 {
			currentRetryCount = 1
		}
		backoff := (float64(currentRetryCount) * float64(minTimeout.Nanoseconds())) + (jitter * float64(jitterMax.Nanoseconds()))
		return time.Duration(backoff)
	}
}
