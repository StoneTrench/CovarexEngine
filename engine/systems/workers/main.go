package workers

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrNilJob     = errors.New("a nil job was received")
	ErrPoolClosed = errors.New("the pool is closed")
)

type JobError struct {
	Value any
}

func (j *JobError) Error() string {
	return fmt.Sprintf("worker recovered: %v", j.Value)
}

type Job func(ctc context.Context)

type Pool struct {
	ctx    context.Context
	cancel func()

	jobs      chan Job
	is_closed bool

	mu        sync.Mutex
	err_mu    sync.Mutex
	worker_wg sync.WaitGroup
	submit_wg sync.WaitGroup
	once      sync.Once

	on_error func(e error)
}

// NewPool creates and initializes a new pool.
//
// timeout - when duration is 0, then there is no timeout, jobs get an unlimited amount of execution time.
//
// on_error - if it is not nil, then workers will recover job panics, and emit them as error in the callback.
// This is not exactly thread safe, cause the goroutines call it directly. The callsite has a mutex lock, so that only one worker can call it at a time.
// This callback must not panic.
func NewPool(worker_count uint, job_queue_size uint, timeout time.Duration, on_error func(error)) *Pool {
	p := new(Pool)

	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.jobs = make(chan Job, job_queue_size)
	p.worker_wg.Add(int(worker_count))

	p.on_error = on_error

	for range worker_count {
		go p.worker(timeout)
	}

	return p
}

func (p *Pool) execute_job(job Job, ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			p.err_mu.Lock()
			defer p.err_mu.Unlock()

			if p.on_error != nil {
				p.on_error(&JobError{Value: r})
			}
		}
	}()

	job(ctx)
}
func (p *Pool) worker(timeout time.Duration) {
	defer p.worker_wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case job, ok := <-p.jobs:
			if !ok { // we only get ok: if (channel closed) or (channel open && has Job)
				return
			}
			if p.ctx.Err() != nil {
				return
			}
			if timeout > 0 {
				ctx, cancel := context.WithTimeout(p.ctx, timeout)
				p.execute_job(job, ctx)
				cancel()
			} else {
				p.execute_job(job, p.ctx)
			}
		}
	}
}

// Submit submits a job to the queue.
func (p *Pool) Submit(j Job) error {
	if j == nil {
		return ErrNilJob
	}
	p.mu.Lock()
	if p.is_closed {
		p.mu.Unlock()
		return ErrPoolClosed
	}

	p.submit_wg.Add(1)
	p.mu.Unlock()

	defer p.submit_wg.Done()

	select {
	case p.jobs <- j:
		return nil
	case <-p.ctx.Done():
		return ErrPoolClosed
	}
}

// Close is graceful synchronous shutdown, new jobs cannot be submitted, but existing jobs will get done. This will block the thread until all jobs finish.
func (p *Pool) Close() {
	p.once.Do(func() {
		p.mu.Lock()
		p.is_closed = true
		p.mu.Unlock()

		p.submit_wg.Wait()
		close(p.jobs)
		p.worker_wg.Wait()
		p.cancel()
	})
}

// CloseImmediate shuts down the workers and pool immediately, any submissions underway fail, and any jobs in the queue get dropped. This will block the thread until all workers die.
func (p *Pool) CloseImmediate() {
	p.once.Do(func() {
		p.mu.Lock()
		p.is_closed = true
		p.mu.Unlock()

		p.cancel()
		p.submit_wg.Wait()
		close(p.jobs)
		p.worker_wg.Wait()
	})
}
