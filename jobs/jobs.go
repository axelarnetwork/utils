package jobs

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/semaphore"
)

// Job represents a (long-running) process that can be spawned on a separate go-routine.
// When encountering an error, a Job should send that error to the given channel and continue to run.
type Job func(ctx context.Context) error

// JobManager manages multiple concurrent jobs and handles their errors. Can wait for all jobs and error handling to finish.
type JobManager struct {
	wgJobs      *sync.WaitGroup
	done        chan struct{}
	errChan     chan error
	ctx         context.Context
	once        *sync.Once
	jobCapacity *capacityMgr
}

type capacityMgr struct {
	semaphore *semaphore.Weighted
}

func (mgr *capacityMgr) Acquire(ctx context.Context, n int64) error {
	if mgr.semaphore != nil {
		return mgr.semaphore.Acquire(ctx, n)
	}
	return ctx.Err()
}

func (mgr *capacityMgr) Release(n int64) {
	if mgr.semaphore != nil {
		mgr.semaphore.Release(n)
	}
}

// MgrOptions modify the behaviour of the JobManager
type MgrOptions func(*JobManager) *JobManager

// WithMaxCapacity defines how many jobs will be run in parallel
func WithMaxCapacity(cap int64) MgrOptions {
	return func(mgr *JobManager) *JobManager {
		mgr.jobCapacity = &capacityMgr{semaphore.NewWeighted(cap)}
		return mgr
	}
}

// WithErrorCacheCapacity defines the size of the error cache. If the cache is full new errors from jobs will be ignored. Default is 1000
func WithErrorCacheCapacity(cap int64) MgrOptions {
	return func(mgr *JobManager) *JobManager {
		mgr.errChan = make(chan error, cap)
		return mgr
	}
}

// NewMgr returns a new JobManager
func NewMgr(ctx context.Context, opts ...MgrOptions) *JobManager {
	mgr := &JobManager{
		ctx:         ctx,
		jobCapacity: &capacityMgr{},
		done:        make(chan struct{}),
		once:        &sync.Once{},
		errChan:     make(chan error, 1000),
		wgJobs:      &sync.WaitGroup{},
	}

	for _, opt := range opts {
		mgr = opt(mgr)
	}

	return mgr
}

// AddJobs calls AddJob for each of the given jobs
func (mgr *JobManager) AddJobs(jobs ...Job) {
	for _, j := range jobs {
		mgr.AddJob(j)
	}
}

// AddJob spawns a new goroutine for the given job, manages its lifetime and handles its errors
func (mgr *JobManager) AddJob(j Job) {
	mgr.wgJobs.Add(1)
	go func() {
		if err := mgr.jobCapacity.Acquire(mgr.ctx, 1); err != nil {
			mgr.tryCacheError(err)
			mgr.wgJobs.Done()
			return
		}
		go func() {
			defer mgr.wgJobs.Done()
			defer mgr.jobCapacity.Release(1)
			defer mgr.recovery()
			if err := j(mgr.ctx); err != nil {
				mgr.tryCacheError(err)
			}
		}()
	}()
}

func (mgr *JobManager) recovery() {
	if r := recover(); r != nil {
		mgr.tryCacheError(fmt.Errorf("job panicked: %s", r))
	}
}

func (mgr *JobManager) tryCacheError(err error) {
	// do not block if the error queue is already full
	select {
	case mgr.errChan <- err:
		break
	default:
		break
	}
}

// Done returns a channel that gets closed when all jobs finished
func (mgr *JobManager) Done() <-chan struct{} {
	go func() {
		mgr.once.Do(func() {
			mgr.wgJobs.Wait()
			close(mgr.errChan)
			close(mgr.done)
		})
	}()

	return mgr.done
}

// Errs returns errors encountered during job execution
func (mgr *JobManager) Errs() <-chan error {
	return mgr.errChan
}
