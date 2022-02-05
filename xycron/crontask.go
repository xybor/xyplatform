package xycron

import (
	"math"
	"strings"
	"sync"
	"time"

	"github.com/xybor/xyplatform/xycond"
)

// cronTask is a task which is mainly used for cyclical task. It can help task
// run at a specific time in hour, day, or more.
type cronTask struct {
	// These are time units at which the task is run.
	second, minute, hour, day, weekday, month unitRange

	// The task will run in t times. Nil represents for infinite times.
	t int64

	// The channel on which task sends if the time has come.
	signalC chan time.Time

	// done is used to notify the next time should be scheduled.
	done chan any

	// canceledFunc is used for stopping the task before out of runs.
	canceledFunc func() bool

	// mu is a locker ensuring there is no unexpected behaviors.
	mu sync.Mutex

	*commonTask
}

func newCronTask(sched *scheduler) *cronTask {
	var t = &cronTask{
		second:       newUnitRange("second", 59, atoui),
		minute:       newUnitRange("minute", 59, atoui),
		hour:         newUnitRange("hour", 23, atoui),
		day:          newUnitRange("day", 30, atoui),
		weekday:      newUnitRange("weekday", 06, atowd),
		month:        newUnitRange("month", 11, atom),
		t:            math.MaxInt64,
		signalC:      make(chan time.Time),
		done:         make(chan any),
		canceledFunc: nil,
		mu:           sync.Mutex{},
	}
	t.commonTask = newCommonTask(sched, t)

	return t
}

func (t *cronTask) signal() <-chan time.Time {
	return t.signalC
}

// run signals for start() method to schedule the next run, then calls the
// function in task.
func (t *cronTask) run() {
	t.mu.Lock()
	t.t -= 1
	t.mu.Unlock()

	// The next run could be scheduled now.
	t.done <- nil

	logger.Info("event=call-func task=%s remain=%d", t.id, t.t)
	rv, err := callFunc(t.f, t.params...)

	rvi := make([]any, len(rv))
	for i, v := range rv {
		rvi[i] = v.Interface()
	}

	if err != nil {
		logger.Error("event=call-func-failed task=%s params=%v return=%v err=%s",
			t.id, t.params, rvi, err)
	} else {
		logger.Debug("event=call-func-done task=%s params=%v return=%v",
			t.id, t.params, rvi)
	}
}

// findNextRun() finds scheduled time of the next run.
func (t *cronTask) findNextRun(now time.Time) timePoint {
	var current = newTimePoint(now)

	// Assume now as next run.
	var next = newTimePoint(now)

	// Avoid the task is executed at the same time when calling findNextRun
	// multiple times and too fast, the task will be scheduled to run at least
	// after one second from now.
	current.sec += 1

	next.sec = t.second.findNext(current.sec)
	// If the next second rotates back, the expected next minute must be at
	// least added by 1.
	if next.sec < current.sec {
		current.min += 1
	}

	next.min = t.minute.findNext(current.min)
	// If the next minute rotates back, the expected next hour must be at least
	// added by 1.
	if next.min < current.min {
		current.hour += 1
	}
	// If the next minute is different, reset the next second.
	if next.min != now.Minute() {
		next.sec = t.second.r[0]
	}

	next.hour = t.hour.findNext(current.hour)
	// If the next hour rotates back, the expected next day must be at least
	// added by 1.
	if next.hour < current.hour {
		current.day += 1
	}
	// If the next hour is different, reset the next second and minute.
	if next.hour != now.Hour() {
		next.sec = t.second.r[0]
		next.min = t.minute.r[0]
	}

	// Maybe the next date is invalid (e.g., 30/2/2000, invalid weekday), it
	// should find the next date until the date is valid.
	var isValidDate = false
	for !isValidDate {
		next.day = t.day.findNext(current.day)
		// If the next day rotates back, the expected next month must be at
		// least added by 1.
		if next.day < current.day {
			current.mon += 1
		}
		// If the next day is different, reset the next clock.
		if next.day != current.day {
			next.sec = t.second.r[0]
			next.min = t.minute.r[0]
			next.hour = t.hour.r[0]
		}

		next.mon = time.Month(t.month.findNext(int(current.mon)))
		// If the next month rotates back, the expected next year must be at
		// least added by 1.
		if next.mon < current.mon {
			next.year += 1
		}
		// If the next month is different, reset the next clock and day.
		if next.mon != current.mon {
			next.sec = t.second.r[0]
			next.min = t.minute.r[0]
			next.hour = t.hour.r[0]
			next.day = t.day.r[0]
		}

		// Cron task does not need to find year, all years are valid.

		// Always find at the next day in case the date is invalid.
		current.day = next.day + 1
		current.mon = next.mon

		var nextWd = next.toTime(t.sched.loc).Weekday()
		isValidDate = t.weekday.contains(int(nextWd)) && next.isValidDay()
	}

	return next
}

// scheduleNextRun() finds the datetime of next run and will send a value to
// task.signal when the countdown is over.
func (t *cronTask) scheduleNextRun() {
	now := time.Now().In(t.sched.loc)
	nextTime := t.findNextRun(now)

	logger.Debug("event=scheduled task=%s time=%s", t.id, nextTime.toTime(t.sched.loc))

	d := nextTime.toTime(t.sched.loc).Sub(now)
	timer := time.AfterFunc(d, func() {
		t.signalC <- time.Now().In(t.sched.loc)
	})

	t.canceledFunc = timer.Stop
}

// start() begins a loop of calling task.scheduleNextRun(), it only stops if
// there is no more run times.
func (t *cronTask) start() {
	logger.Info("event=start task=%s", t.id)

	var isBreak = false

	for {
		t.mu.Lock()
		if t.t <= 0 {
			close(t.signalC)
			t.canceledFunc = nil
			isBreak = true
		} else {
			t.scheduleNextRun()
		}
		t.mu.Unlock()

		if isBreak {
			break
		}

		<-t.done
	}
}

func (t *cronTask) stop() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	logger.Info("event=stop task=%s", t.id)

	if t.canceledFunc == nil && t.t > 0 {
		return StopError.New("do not stop when the task has not started yet")
	}

	t.t = 0

	if t.canceledFunc != nil && t.canceledFunc() {
		logger.Debug("event=cancel-timer task=%s", t.id)
		t.done <- nil
	}

	return nil
}

// Times sets the number of run times, it panics if the parameter n is not
// positive.
func (t *cronTask) Times(n int) *cronTask {
	xycond.Condition(n > 0).Assertf("Expected a positive number, but got %d", n)

	t.mu.Lock()
	defer t.mu.Unlock()

	t.t = int64(n)
	return t
}

// Once is a shortcut of Times(1).
func (t *cronTask) Once() *cronTask {
	return t.Times(1)
}

// Twice is a shortcut of Times(2).
func (t *cronTask) Twice() *cronTask {
	return t.Times(2)
}

// Infinity is a shortcut of Times(math.MaxInt64).
func (t *cronTask) Infinity() *cronTask {
	return t.Times(math.MaxInt64)
}

// Set seconds that the task will be run. The parameter must be a string, see
// linux cron for further details. It panics if the value is wrong.
func (t *cronTask) Seconds(s string) *cronTask {
	err := t.second.set(s)
	xycond.Nil(err).Assertf("%s", err)

	return t
}

// Set seconds that the task will be run. The parameters are int numbers.
func (t *cronTask) Second(s ...int) *cronTask {
	return t.Seconds(strings.Join(intListToString(s...), ","))
}

// Set minutes that the task will be run. The parameter must be a string, see
// linux cron for further details. It panics if the value is wrong.
func (t *cronTask) Minutes(s string) *cronTask {
	err := t.minute.set(s)
	xycond.Nil(err).Assertf("%s", err)

	return t
}

// Set minutes that the task will be run. The parameters are int numbers.
func (t *cronTask) Minute(m ...int) *cronTask {
	return t.Minutes(strings.Join(intListToString(m...), ","))
}

// Set hours that the task will be run. The parameter must be a string, see
// linux cron for further details. It panics if the value is wrong.
func (t *cronTask) Hours(s string) *cronTask {
	err := t.hour.set(s)
	xycond.Nil(err).Assertf("%s", err)

	return t
}

// Set hours that the task will be run. The parameters are int numbers.
func (t *cronTask) Hour(h ...int) *cronTask {
	return t.Hours(strings.Join(intListToString(h...), ","))
}

// Set days that the task will be run. The parameter must be a string, see
// linux cron for further details. It panics if the value is wrong.
func (t *cronTask) Days(s string) *cronTask {
	err := t.day.set(s)
	xycond.Nil(err).Assertf("%s", err)

	return t
}

// Set days that the task will be run. The parameters are int numbers.
func (t *cronTask) Day(d ...int) *cronTask {
	return t.Days(strings.Join(intListToString(d...), ","))
}

// Set weekdays that the task will be run. The parameter must be a string, see
// linux cron for further details. It panics if the value is wrong.
func (t *cronTask) Weekdays(s string) *cronTask {
	err := t.weekday.set(s)
	xycond.Nil(err).Assertf("%s", err)

	return t
}

// Set weekdays that the task will be run. The parameters are int numbers.
func (t *cronTask) Weekday(wd ...time.Weekday) *cronTask {
	return t.Weekdays(strings.Join(wdListToString(wd...), ","))
}

// Set months that the task will be run. The parameter must be a string, see
// linux cron for further details. It panics if the value is wrong.
func (t *cronTask) Months(s string) *cronTask {
	err := t.month.set(s)
	xycond.Nil(err).Assertf("%s", err)

	return t
}

// Set months that the task will be run. The parameters are int numbers.
func (t *cronTask) Month(m ...time.Month) *cronTask {
	return t.Months(strings.Join(monListToString(m...), ","))
}
