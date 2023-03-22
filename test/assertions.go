package testutils

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// FailOnTimeout blocks until `done` is closed or until specified timeout has elapsed.
// In the latter case it calls require.FailNow(t, "test timed out").
func FailOnTimeout(t *testing.T, done <-chan struct{}, timeout time.Duration) {
	select {
	case <-done:
		// test is done, nothing to do
	case <-time.After(timeout):
		require.FailNow(t, "test timed out")
	}
}
