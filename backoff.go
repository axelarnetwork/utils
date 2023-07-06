package utils

import (
	"github.com/axelarnetwork/utils/math"
	"math/rand"
	"time"
)

// BackOff computes the next back-off duration
type BackOff func(currentRetryCount int) time.Duration

// ExponentialBackOff computes an exponential back-off
func ExponentialBackOff(minTimeout time.Duration) BackOff {
	return func(currentRetryCount int) time.Duration {
		jitter := rand.Float64() + 0.5
		strategy := 1 << currentRetryCount
		backoff := math.Max(jitter*float64(strategy)*float64(minTimeout.Nanoseconds()), float64(minTimeout.Nanoseconds()))

		return time.Duration(backoff)
	}
}

// LinearBackOff computes a linear back-off
func LinearBackOff(minTimeout time.Duration) BackOff {
	return func(currentRetryCount int) time.Duration {
		jitter := rand.Float64() + 0.5
		backoff := math.Max(jitter*float64(currentRetryCount)*float64(minTimeout.Nanoseconds()), float64(minTimeout.Nanoseconds()))
		return time.Duration(backoff)
	}
}
