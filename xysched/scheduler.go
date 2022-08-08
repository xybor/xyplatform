package xysched

import (
	"context"
	"time"

	"golang.org/x/sync/semaphore"
)

// CallbackFunc is a type alias of a generic callback function
type CallbackFunc any

// future interface defines a schedulable object.
type future interface {
	// run calls the future's function.
	run()

	// copy creates a copy of this future to avoid race condition when adding
	// callback futures to scheduler many times.
	copy() future

	// Return the next time this future will be run. Leave it as nil if you do
	// not want to run anymore.
	next() *time.Time

	// Return callback future objects, they will be sent to scheduler once this
	// future completes.
	callbacks() []future
}

// scheduler is used for scheduling future objects.
type scheduler struct {
	futureQ chan future
	stop    chan any
	sem     *semaphore.Weighted
}

// New creates a scheduler and starts it.
func New() *scheduler {
	sched := &scheduler{
		futureQ: make(chan future),
		stop:    make(chan any),
		sem:     nil,
	}
	go sched.start()
	return sched
}

// After creates a send-only channel. Sending a future to this channel will
// add it to scheduler after a duration. If d is negative, After will send the
// future to scheduler immediately.
//
// NOTE: You should send ONLY ONE future to this channel because it is designed
// to handle one. If you try sending another, it will be blocked forever. To
// send other futures to scheduler, let call this method again.
func (s *scheduler) After(d time.Duration) chan<- future {
	if d < 0 {
		d = 0
	}

	var c = make(chan future)
	go func() {
		var timer *time.Timer
		var done = make(chan any)
		select {
		case <-s.stop:
		case t := <-c:
			timer = time.AfterFunc(d, func() {
				s.futureQ <- t
				close(done)
			})
		}

		select {
		case <-s.stop:
			if timer != nil {
				timer.Stop()
			}
		case <-done:
		}
	}()
	return c
}

// At is a shortcut of After(time.Until(next)).
//
// NOTE: You should send ONLY ONE future to this channel because it is designed
// to handle one. If you try sending another, it will be blocked forever. To
// send other futures to scheduler, let call this method again.
func (s *scheduler) At(next time.Time) chan<- future {
	return s.After(time.Until(next))
}

// Now is a shortcut of After(0).
//
// NOTE: You should send ONLY ONE future to this channel because it is designed
// to handle one. If you try sending another, it will be blocked forever. To
// send other futures to scheduler, let call this method again.
func (s *scheduler) Now() chan<- future {
	return s.After(0)
}

// Stop terminates the scheduler and all pending futures from now on. Running
// futures still run until they complete.
func (s *scheduler) Stop() {
	logger.Debug("event=stopping scheduler=%p", s)
	close(s.stop)
}

// Singleton is a shortcut of Concurrent(1).
func (s *scheduler) Singleton() {
	s.Concurrent(1)
}

// Concurrent limits the number of running futures at the same time. By default,
// there is no limited.
func (s *scheduler) Concurrent(n int) {
	logger.Trace("event=set-concurrent scheduler=%p n=%d", s, n)
	s.sem = semaphore.NewWeighted(int64(n))
}

// start begins the scheduled loop.
func (s *scheduler) start() {
	logger.Debug("event=start scheduler=%p", s)
	var isStop = false

	for !isStop {
		select {
		case <-s.stop:
			isStop = true
		case f := <-s.futureQ:
			go func() {
				if s.sem != nil {
					s.sem.Acquire(context.TODO(), 1)
					defer s.sem.Release(1)
				}
				f.run()
				for _, cb := range f.callbacks() {
					s.Now() <- cb.copy()
				}
				if next := f.next(); next != nil {
					s.At(*next) <- f
				}
			}()
		}
	}
	logger.Debug("event=stop scheduler=%p", s)
}
