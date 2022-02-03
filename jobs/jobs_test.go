package jobs_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/axelarnetwork/utils/jobs"
	"github.com/stretchr/testify/assert"

	"github.com/axelarnetwork/utils/test/rand"
)

func TestJobManager_Errs(t *testing.T) {
	// setup
	jobCount := rand.I64Between(0, 100)
	mgr := jobs.NewMgr(context.Background())
	for i := int64(0); i < jobCount; i++ {
		job := func(ctx context.Context) error {
			return fmt.Errorf("error by job %d", i)
		}
		mgr.AddJob(job)
	}
	// test
	<-mgr.Done()
	assert.Len(t, mgr.Errs(), int(jobCount))
}

func TestJobManager_Errs_WithSmallCachSize(t *testing.T) {
	jobCount := rand.I64Between(20, 100)
	errorCacheSize := int64(10)
	mgr := jobs.NewMgr(context.Background(), jobs.WithErrorCacheCapacity(errorCacheSize))
	for i := int64(0); i < jobCount; i++ {
		job := func(ctx context.Context) error {
			return fmt.Errorf("error by job %d", i)
		}
		mgr.AddJob(job)
	}
	// test
	<-mgr.Done()
	assert.Len(t, mgr.Errs(), int(errorCacheSize))
}

func TestJobManager_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	mgr := jobs.NewMgr(ctx)

	jobCount := rand.I64Between(0, 100)
	for i := 0; i < int(jobCount); i++ {
		mgr.AddJob(func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		})
	}

	select {
	case <-mgr.Done():
		assert.Fail(t, "it should be impossible for the mgr to be done here")
	default:
		break
	}

	cancel()

	timeout, timeoutCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer timeoutCancel()

	select {
	case <-mgr.Done():
		break
	case <-timeout.Done():
		assert.Fail(t, "timed out")
	}
}

func TestJobManager_MoreJobsThanCap(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	cap := rand.I64Between(1, 20)
	mgr := jobs.NewMgr(ctx, jobs.WithMaxCapacity(cap))

	jobCount := rand.I64Between(cap, 100)
	blockingCtx, cancel := context.WithCancel(context.Background())
	jobsStarted := int64(0)
	for i := 0; i < int(jobCount); i++ {
		mgr.AddJob(func(context.Context) error {
			atomic.AddInt64(&jobsStarted, 1)
			<-blockingCtx.Done()
			return nil
		})
	}

	// only jobs up to the cap have started because no jobs have finished yet
	assert.Equal(t, cap, jobsStarted)
	cancel()

	timeout, timeoutCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer timeoutCancel()

	select {
	case <-mgr.Done():
		break
	case <-timeout.Done():
		assert.Fail(t, "timed out")
	}

	// now all jobs must have finished
	assert.Equal(t, jobCount, jobsStarted)
}

// this extracted function is needed to close over the loop counter i
