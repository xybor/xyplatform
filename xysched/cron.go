package xysched

import (
	"math"
	"time"

	"github.com/xybor/xyplatform/xycond"
)

// Cron is a future which runs a task periodically.
type Cron struct {
	// cron struct is a wrapper of task
	Task

	// The maximum times which this future runs.
	n uint

	// The periodic duration.
	d time.Duration

	// Callback futures will be run when cron ran out of times.
	onfinish []future
}

// NewCron creates a future which calls function f with parameter params
// periodically. By default, it runs the function forever secondly.
func NewCron(f any, params ...any) *Cron {
	return &Cron{
		Task:     *NewTask(f, params...),
		n:        math.MaxInt,
		d:        time.Second,
		onfinish: make([]future, 0),
	}
}

// Secondly requires the cron to run once per second.
func (c *Cron) Secondly() *Cron {
	return c.Every(time.Second)
}

// Minutely requires the cron to run the cron once per minute.
func (c *Cron) Minutely() *Cron {
	return c.Every(time.Minute)
}

// Hourly requires the cron to run once per hour.
func (c *Cron) Hourly() *Cron {
	return c.Every(time.Hour)
}

// Daily requires the cron to run once per day.
func (c *Cron) Daily() *Cron {
	return c.Every(24 * time.Hour)
}

// Every requires the cron to run with a custom periodic duration.
func (c *Cron) Every(d time.Duration) *Cron {
	xycond.AssertNotLessThan(int(d), 0)
	c.d = d
	return c
}

// Times sets the maximum times which the cron will run.
func (c *Cron) Times(n uint) *Cron {
	c.n = n
	return c
}

// Once is a shortcut of Times(1)
func (c *Cron) Once() *Cron {
	return c.Times(1)
}

// Twice is a shortcut of Times(2)
func (c *Cron) Twice() *Cron {
	return c.Times(2)
}

// Finish sets a callback future which will run after the cron ran out of times.
// See task.Callback for further details.
func (c *Cron) Finish(f any, params ...any) *Task {
	var cb = toFuture(f, params...)
	c.onfinish = append(c.onfinish, cb)

	if t, ok := cb.(*Task); ok {
		return t
	}
	return nil
}

// Required method of future. This method overrides the one of task.
func (c *Cron) next() *time.Time {
	var n = c.lock.RLockFunc(func() any {
		c.n--
		return c.n
	}).(uint)
	if n > 0 {
		var t = time.Now().Add(c.d)
		return &t
	}
	return nil
}

// Required method of future. This method overrides the one of task.
func (c *Cron) callbacks() []future {
	var cb []future
	cb = append(cb, c.Task.callbacks()...)
	if c.lock.RLockFunc(func() any { return c.n }).(uint) == 0 {
		cb = append(cb, c.onfinish...)
	}
	return cb
}
