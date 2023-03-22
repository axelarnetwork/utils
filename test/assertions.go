package testutils

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// FailOnTimeout blocks until `ctx` is done or until specified timeout has elapsed.
// In the latter case it calls require.FailNow(t, "test timed out").
func FailOnTimeout(ctx context.Context, t *testing.T, timeout time.Duration) {
	select {
	case <-ctx.Done():
		// context got cancelled in time, nothing to do
	case <-time.After(timeout):
		require.FailNow(t, "test timed out")
	}
}
