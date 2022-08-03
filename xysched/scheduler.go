package xysched

import (
	"context"
	"time"

	"golang.org/x/sync/semaphore"
)

// Type alias of a generic callback function
type CallbackFunc any

// future interface defines a schedulable object.
type future interface {
	// This method calls the future's function.
	run()

	// Copy to another future, avoid race condition.
	copy() future

	// Return the next time this future will be run. Leave it as nil if you do
	// not want to run anymore.
	next() *time.Time

	// Return callback future objects, they will be sent to scheduler once this
	// future completes.
	callbacks() []future
}

// event interface defines an object which scheduler will run callback functions
// once event is triggered.
type event interface {
	// Wait returns a receive-only channel receiving the parameters passed into
	// callback functions. Once receive parameters, the scheduler will schedule
	// callback functions immediately. Close this channel to stop event waiting.
	Wait() <-chan []any

	// Callbacks return all callback functions of the event. Parameters of these
	// functions need to be the same with the one in wait method.
	Callbacks() []CallbackFunc
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
// add it to scheduler after a duration.
//
// NOTE: You should send ONLY ONE future to this channel because it is designed
// to handle one. If you try sending another, it will be blocked forever. To
// send other futures to scheduler, let call this method again.
func (s *scheduler) After(d time.Duration) chan<- future {
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

// At is a shortcut of After(time.Until(next)).In case next time is in the past,
// At will send the future to scheduler immediately.
//
// NOTE: You should send ONLY ONE future to this channel because it is designed
// to handle one. If you try sending another, it will be blocked forever. To
// send other futures to scheduler, let call this method again.
func (s *scheduler) At(next time.Time) chan<- future {
	var d = time.Until(next)
	if d < 0 {
		d = 0
	}
	return s.After(d)
}

// Now is a shortcut of After(0).
//
// NOTE: You should send ONLY ONE future to this channel because it is designed
// to handle one. If you try sending another, it will be blocked forever. To
// send other futures to scheduler, let call this method again.
func (s *scheduler) Now() chan<- future {
	return s.After(0)
}

// Register will send callback functions of event everytime it triggers.
func (s *scheduler) Register(e event) {
	go func() {
		var isStop = false
		var trigger = e.Wait()
		for !isStop {
			select {
			case <-s.stop:
				isStop = true
			case params, ok := <-trigger:
				if !ok {
					isStop = true
					break
				}
				var t = Task(func() []any { return params })
				t.Variadic(len(params))
				for _, cb := range e.Callbacks() {
					t.Then(cb)
				}
				s.Now() <- t
			}
		}
	}()
}

// Stop the scheduler, all not-yet-run futures will not run forever from now on.
// Running futures still run until it completes.
func (s *scheduler) Stop() {
	logger.Debug("event=stopping scheduler=%p", s)
	close(s.stop)
}

// Singleton is a shortcut of Concurrent(1).
func (s *scheduler) Singleton() {
	s.Concurrent(1)
}

// Singleton limits the number of running futures at the same time. By default,
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
					s.Now() <- cb
				}
				if next := f.next(); next != nil {
					s.At(*next) <- f
				}
			}()
		}
	}
	logger.Debug("event=stop scheduler=%p", s)
}
