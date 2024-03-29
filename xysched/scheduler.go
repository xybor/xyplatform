package xysched

import (
	"fmt"
	"time"

	"github.com/xybor/xyplatform/xylock"
)

// future interface defines a schedulable object.
type future interface {
	// run calls the future's function.
	run()

	// Return the next time this future will be run. Leave it as nil if you do
	// not want to run anymore.
	next() *time.Time

	// Return callback future objects, they will be sent to scheduler once this
	// future completes.
	callbacks() []future

	// stop returns the channel which is closed if the future early stops.
	stop() <-chan any
}

// Scheduler is used for scheduling future objects.
type Scheduler struct {
	name    string
	futureQ chan future
	stop    chan any
	sem     *xylock.Semaphore
}

// NewScheduler returns Scheduler associated with the name, if it has not yet
// existed, create a new one.
//
// Any Scheduler with a non-empty name will be associated with its name. Calling
// this function twice with the same name gives you the same Scheduler. If you
// want to create different Schedulers each call, use the empty name.
func NewScheduler(name string) *Scheduler {
	var sched *Scheduler
	var ok bool
	lock.RLockFunc(func() any {
		sched, ok = schedulerManager[name]
		return nil
	})
	if ok {
		return sched
	}

	if name == "" {
		lock.WLockFunc(func() {
			name = fmt.Sprintf("scheduler-%d", anonSchedCounter)
			anonSchedCounter++
		})
	}

	sched = &Scheduler{
		name:    name,
		futureQ: make(chan future),
		stop:    make(chan any),
		sem:     nil,
	}

	lock.WLockFunc(func() {
		schedulerManager[name] = sched
	})

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
func (s *Scheduler) After(d time.Duration) chan<- future {
	if d < 0 {
		d = 0
	}

	var c = make(chan future)
	go func() {
		var f future
		var timer *time.Timer
		var done = make(chan any)
		select {
		case <-s.stop:
		case f = <-c:
			timer = time.AfterFunc(d, func() {
				s.futureQ <- f
				close(done)
			})
			logger.Event("prepare-to-schedule").
				Field("scheduler", s.name).Field("future", f).Field("after", d).
				Debug()
		}

		select {
		case <-s.stop:
			if timer != nil {
				timer.Stop()
			}
		case <-f.stop():
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
func (s *Scheduler) At(next time.Time) chan<- future {
	return s.After(time.Until(next))
}

// Now is a shortcut of After(0).
//
// NOTE: You should send ONLY ONE future to this channel because it is designed
// to handle one. If you try sending another, it will be blocked forever. To
// send other futures to scheduler, let call this method again.
func (s *Scheduler) Now() chan<- future {
	return s.After(0)
}

// Stop terminates the scheduler and all pending futures from now on. Running
// futures still run until they complete.
func (s *Scheduler) Stop() {
	logger.Event("signal-stop").Field("scheduler", s.name).Info()
	close(s.stop)
}

// Singleton is a shortcut of Concurrent(1).
func (s *Scheduler) Singleton() {
	s.Concurrent(1)
}

// Concurrent limits the number of running futures at the same time. By default,
// there is no limited.
func (s *Scheduler) Concurrent(n int) {
	s.sem = xylock.NewSemaphore(int64(n))
	logger.Event("set-concurrent").
		Field("scheduler", s.name).Field("futures", n).Debug()
}

// start begins the scheduled loop.
func (s *Scheduler) start() {
	logger.Event("start").Field("scheduler", s.name).Info()
	var isStop = false

	for !isStop {
		select {
		case <-s.stop:
			isStop = true
		case f := <-s.futureQ:
			select {
			case <-f.stop():
			default:
				if next := f.next(); next != nil {
					s.At(*next) <- f
				}
				go s.sem.AcquireFunc(1, func() {
					f.run()
					for _, cb := range f.callbacks() {
						s.Now() <- cb
					}
				})
			}
		}
	}
	logger.Event("stop").Field("scheduler", s.name).Info()
}
