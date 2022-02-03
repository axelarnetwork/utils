package jobs_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/axelarnetwork/utils/jobs"
	. "github.com/axelarnetwork/utils/test"
	"github.com/stretchr/testify/assert"

	"github.com/axelarnetwork/utils/test/rand"
)

func TestJobManager_Errs(t *testing.T) {
	var (
		mgr      *jobs.JobManager
		jobCount int64
	)

	Given("a job manager", func(t *testing.T) {
		mgr = jobs.NewMgr(context.Background())
	}).
		When("adding jobs that fail", func(t *testing.T) {
			jobCount = rand.I64Between(0, 100)
			for i := int64(0); i < jobCount; i++ {
				job := func(ctx context.Context) error {
					return fmt.Errorf("error by job %d", i)
				}
				mgr.AddJob(job)
			}
		}).
		Then("find all errors in cache after jobs are done", func(t *testing.T) {
			<-mgr.Done()
			assert.Len(t, mgr.Errs(), int(jobCount))
		}).Run(t, 20)

	var (
		errorCacheSize int64
	)

	Given("a job manager with small error cache", func(t *testing.T) {
		errorCacheSize = rand.I64Between(1, 20)
		mgr = jobs.NewMgr(context.Background(), jobs.WithErrorCacheCapacity(errorCacheSize))
	}).
		When("more jobs are managed", func(t *testing.T) {
			jobCount = rand.I64Between(errorCacheSize, 100)
			for i := int64(0); i < jobCount; i++ {
				job := func(ctx context.Context) error {
					return fmt.Errorf("error by job %d", i)
				}
				mgr.AddJob(job)
			}
		}).
		Then("ignore errors that exceed the cache", func(t *testing.T) {
			<-mgr.Done()
			assert.Len(t, mgr.Errs(), int(errorCacheSize))
		}).Run(t, 20)

	var (
		cancel context.CancelFunc
	)

	Given("a job manager with cancellable context", func(t *testing.T) {
		var ctx context.Context
		ctx, cancel = context.WithCancel(context.Background())
		mgr = jobs.NewMgr(ctx)
	}).Branch(
		When("blocking jobs are added", func(t *testing.T) {
			jobCount = rand.I64Between(0, 100)
			for i := 0; i < int(jobCount); i++ {
				mgr.AddJob(func(ctx context.Context) error {
					<-ctx.Done()
					return nil
				})
			}
		}).
			Then("block until the context is cancelled", func(t *testing.T) {
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
			}),
		When("jobs are added after the context is cancelled", func(t *testing.T) {
			cancel()

			jobCount = rand.I64Between(0, 100)
			for i := 0; i < int(jobCount); i++ {
				mgr.AddJob(func(ctx context.Context) error {
					assert.Fail(t, "should not have been called")
					<-ctx.Done()
					return nil
				})
			}
		}).Then("do not execute the jobs", func(t *testing.T) {
			timeout, timeoutCancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer timeoutCancel()

			select {
			case <-mgr.Done():
				break
			case <-timeout.Done():
				assert.Fail(t, "timed out")
			}
		})).Run(t, 20)

	var (
		capacity    int64
		jobsStarted int64
		unblockJobs context.CancelFunc
	)

	Given("a capacity limited job manager", func(t *testing.T) {
		var ctx context.Context
		ctx, cancel = context.WithCancel(context.Background())
		capacity = rand.I64Between(1, 20)
		mgr = jobs.NewMgr(ctx, jobs.WithMaxCapacity(capacity))
	}).
		When("more blocking jobs than capacity are added", func(t *testing.T) {
			jobCount = rand.I64Between(capacity, 100)

			// use this context to prevent the jobs from immediately completing
			var blockingCtx context.Context
			blockingCtx, unblockJobs = context.WithCancel(context.Background())
			jobsStarted = 0
			for i := 0; i < int(jobCount); i++ {
				mgr.AddJob(func(context.Context) error {
					atomic.AddInt64(&jobsStarted, 1)
					<-blockingCtx.Done()
					return nil
				})
			}
		}).
		Then("block until all jobs are done", func(t *testing.T) {
			select {
			case <-mgr.Done():
				assert.Fail(t, "it should be impossible for the mgr to be done here")
			default:
				break
			}

			timeout, timeoutCancel := context.WithTimeout(context.Background(), 1*time.Second)

			for jobsStarted < capacity {
				select {
				case <-timeout.Done():
					assert.Fail(t, "timed out", "jobs started: %d, capacity: %d", jobsStarted, capacity)
				default:
					time.Sleep(5 * time.Millisecond)
				}
			}
			timeoutCancel()

			// only jobs up to the cap have started because no jobs have finished yet
			assert.Equal(t, capacity, jobsStarted)

			unblockJobs()

			timeout, timeoutCancel = context.WithTimeout(context.Background(), 1*time.Second)
			defer timeoutCancel()

			select {
			case <-mgr.Done():
				break
			case <-timeout.Done():
				assert.Fail(t, "timed out")
			}

			// now all jobs must have finished
			assert.Equal(t, jobCount, jobsStarted)
		}).Run(t, 20)
}
