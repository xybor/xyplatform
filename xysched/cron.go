package xysched

import (
	"math"
	"time"

	"github.com/xybor/xyplatform/xycond"
)

// cron is a future which runs a task periodically.
type cron struct {
	// cron struct is a wrapper of task
	*task

	// The maximum times which this future runs.
	n int

	// The periodic duration.
	d time.Duration

	// Callback futures will be run when cron ran out of times.
	onfinish []future
}

// Cron creates a future which calls function f with parameter params
// periodically. By default, it runs the function forever secondly.
func Cron(f any, params ...any) *cron {
	return &cron{
		task:     Task(f, params...),
		n:        math.MaxInt,
		d:        time.Second,
		onfinish: make([]future, 0),
	}
}

// Secondly requires the cron to run once per second.
func (c *cron) Secondly() *cron {
	return c.Every(time.Second)
}

// Minutely requires the cron to run the cron once per minute.
func (c *cron) Minutely() *cron {
	return c.Every(time.Minute)
}

// Hourly requires the cron to run once per hour.
func (c *cron) Hourly() *cron {
	return c.Every(time.Hour)
}

// Daily requires the cron to run once per day.
func (c *cron) Daily() *cron {
	return c.Every(24 * time.Hour)
}

// Every requires the cron to run with a custom periodic duration.
func (c *cron) Every(d time.Duration) *cron {
	c.d = d
	return c
}

// Times sets the maximum times which the cron will run.
func (c *cron) Times(n int) *cron {
	c.n = n
	return c
}

// Once is a shortcut of Times(1)
func (c *cron) Once() *cron {
	return c.Times(1)
}

// Twice is a shortcut of Times(2)
func (c *cron) Twice() *cron {
	return c.Times(2)
}

// Finish sets a callback future which will run after the cron ran out of times.
// See task.Callback for further details.
func (c *cron) Finish(f any, params ...any) *task {
	cb, ok := f.(future)
	if ok {
		xycond.Empty(params).
			Assertf("Do not pass params if f was already a tasker")
	} else {
		cb = Task(f, params...)
	}

	c.onfinish = append(c.onfinish, cb)

	if t, ok := cb.(*task); ok {
		return t
	}
	return nil
}

// Required method of future. This method overrides the one of task.
func (c *cron) next() *time.Time {
	var n = c.lock.RLockFunc(func() any {
		c.n -= 1
		return c.n
	}).(int)
	if n > 0 {
		var t = time.Now().Add(c.d)
		return &t
	}
	return nil
}

// Required method of future. This method overrides the one of task.
func (c *cron) callbacks() []future {
	var cb []future
	cb = append(cb, c.task.callbacks()...)
	if c.lock.RLockFunc(func() any { return c.n }).(int) == 0 {
		cb = append(cb, c.onfinish...)
	}
	return cb
}
