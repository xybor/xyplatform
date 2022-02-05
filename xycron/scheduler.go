package xycron

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/xybor/xyplatform/xyerror"
	"github.com/xybor/xyplatform/xyselect"
	"golang.org/x/sync/semaphore"
)

type scheduler struct {
	// The task list, it includes both stopped and live tasks.
	t []task

	// The semaphore object is used for concurrency mode in scheduler. Nil is
	// for free mode.
	sem *semaphore.Weighted

	loc *time.Location

	// The selector is using for selecting a dynamical list of task channels.
	selector xyselect.Selector

	// The channel is used for preventing the selector returning an error of
	// ExhaustedError when all current tasks stopped.
	cancel chan any

	// inProgress indicates if the scheduler is in progress or not.
	inProgress bool

	mu sync.Mutex
}

// New creates a new scheduler. It is used to schedule and run tasks.
//
// Call Schedule() with the scheduler to create a new task, then using task's
// method to schedule when the task will be run, how many times it runs. The
// default task will run every second.
//
// Use Every(.), or EveryX(.), ..., or After(.) to create tasks with
// predefined time.
//
// There are three modes in scheduler: Free, Concurrent, and Singleton.
//   Free - default: There is no limit about how many tasks run simultaneously.
//   Concurrent: You can set the number of tasks run simultaneously.
//   Singleton: There is only one task run at the same time.
// In the case the number of running tasks is exceed the limit, other tasks
// will not be scheduled until one of running tasks completes.
func New() *scheduler {
	return &scheduler{
		t:          make([]task, 0),
		sem:        nil,
		loc:        time.Local,
		selector:   *xyselect.E(),
		cancel:     nil,
		inProgress: false,
		mu:         sync.Mutex{},
	}
}

// add adds a new task to scheduler.
func (sched *scheduler) add(t task) {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	sched.t = append(sched.t, t)

	if sched.inProgress {
		sched.selector.Recv(xyselect.C(t.signal()))
		go t.start()
	}
}

// Location sets the new time location to scheduler, time.Local by default.
func (sched *scheduler) Location(loc *time.Location) *scheduler {
	sched.loc = loc
	return sched
}

// ConcurrentMode sets the limit number of task run at the same time.
func (sched *scheduler) ConcurrentMode(n uint) *scheduler {
	sched.sem = semaphore.NewWeighted(int64(n))
	return sched
}

// SingletonMode is a shortcut of ConcurrentMode(1)
func (sched *scheduler) SingletonMode() *scheduler {
	return sched.ConcurrentMode(1)
}

// FreeMode sets no limit number of tasks run simultaneously.
func (sched *scheduler) FreeMode() *scheduler {
	sched.sem = nil
	return sched
}

// execute checks how many tasks is running, if it exceeds the limit, returns
// immediately, otherwise, runs the task.
func (sched *scheduler) execute(t task) {
	if sched.sem != nil && !sched.sem.TryAcquire(1) {
		return
	}

	t.run()

	if sched.sem != nil {
		sched.sem.Release(1)
	}
}

// loop creates a infinite loop for selecting the next task by using xyselector
// and signal channel of tasks, then executes the task. This method only stops
// when there is no more live signal channel in selector, it means the cancel
// channel is also closed.
func (sched *scheduler) loop() {
	for {
		i, _, err := sched.selector.Select(false)

		if errors.Is(err, xyselect.ExhaustedError) {
			break
		}

		// Decrease the index by 1 because of the cancel channel at zero index.
		logger.Debug("event=execute task=%s index=%d", sched.t[i-1], i-1)
		go sched.execute(sched.t[i-1])
	}
}

// Start begins scheduling tasks in scheduler. The recommendation is to call
// Start() before creating tasks. All tasks still work even if they was added
// to scheduler before it starts, but the time could be wrong. A panic occurs
// if you call Start while the scheduler has started and still been in
// progress.
func (sched *scheduler) Start() error {
	sched.mu.Lock()
	if sched.inProgress {
		sched.mu.Unlock()
		return InProgressError.New("scheduler is in progress")
	}

	sched.cancel = make(chan any)
	sched.selector.Recv(sched.cancel)

	for _, t := range sched.t {
		sched.selector.Recv(xyselect.C(t.signal()))
		go t.start()
	}

	sched.inProgress = true
	sched.mu.Unlock()

	logger.Info("event=start")

	sched.loop()

	sched.mu.Lock()
	sched.inProgress = false
	sched.mu.Unlock()
	logger.Info("event=stopped")

	return nil
}

// Stop signals all tasks need to terminate their timer and not schedule any
// more. It will return the first error if it occurs.
func (sched *scheduler) Stop() error {
	logger.Info("event=stopping")

	sched.mu.Lock()
	defer sched.mu.Unlock()

	var err error = nil

	if sched.cancel == nil {
		err = StopError.New("do not stop before scheduler starts")
	} else {
		close(sched.cancel)
	}

	var il []int
	for i, t := range sched.t {
		e := t.stop()
		if e != nil {
			il = append(il, i)
			logger.Error("An error occurs at %s: %s", t, e)
		}
	}

	if len(il) != 0 {
		var e = StopError.Newf("can not stop these following tasks: %v", il)
		err = xyerror.Or(err, e)
	}

	return err
}

// Cron creates a new cron task with default setting (it will run every second
// and infinitely).
func (sched *scheduler) Cron() *cronTask {
	return newCronTask(sched)
}

// After creates a new task which runs after a duration. You should only call
// this function if you have called scheduler.Start() before. Otherwise, it
// could never run for many years because the time is over.
//
// After is scheduled to run once. You should not call other methods except for
// Params() and Do().
func (sched *scheduler) After(d time.Duration) *cronTask {
	// The minimum duration which scheduler handles is one second.
	if d < time.Second {
		d = time.Second
		logger.Warn("The duration is rounded up to one second.")
	}

	var tp = newTimePoint(time.Now().Add(d))
	var t = newCronTask(sched)
	t.Second(tp.sec).Minute(tp.min).Hour(tp.hour).Day(tp.day).Month(tp.mon)
	t.Once()

	return t
}

type wrapperCronTask struct {
	interval int
	t        *cronTask
}

// Every generates a wrapper of cron task to create a task running in cycles.
// The unit (second, minute, ..., or month) will be choosen later.
func (sched *scheduler) Every(interval int) wrapperCronTask {
	return wrapperCronTask{interval: interval, t: sched.Cron()}
}

// Set the unit cycle as second.
func (wt wrapperCronTask) Seconds() *cronTask {
	wt.t.Minutes("*").Hours("*").Days("*").Weekdays("*").Months("*")
	wt.t.Seconds("*/" + strconv.Itoa(wt.interval))
	return wt.t
}

// Set the unit cycle as minute.
func (wt wrapperCronTask) Minutes() *cronTask {
	wt.t.Hours("*").Days("*").Weekdays("*").Months("*")
	wt.t.Minutes("*/" + strconv.Itoa(wt.interval))
	return wt.t
}

// Set the unit cycle as hour.
func (wt wrapperCronTask) Hours() *cronTask {
	wt.t.Days("*").Weekdays("*").Months("*")
	wt.t.Hours("*/" + strconv.Itoa(wt.interval))
	return wt.t
}

// Set the unit cycle as day.
func (wt wrapperCronTask) Days() *cronTask {
	wt.t.Weekdays("*").Months("*")
	wt.t.Days("*/" + strconv.Itoa(wt.interval))
	return wt.t
}

// Set the unit cycle as month.
func (wt wrapperCronTask) Months() *cronTask {
	wt.t.Weekdays("*")
	wt.t.Months("*/" + strconv.Itoa(wt.interval))
	return wt.t
}

// EverySecond is a shortcut of Every(1).Seconds().
func (sched *scheduler) EverySecond() *cronTask {
	return sched.Every(1).Seconds()
}

// EveryMinute is a shortcut of Every(1).Minutes().
func (sched *scheduler) EveryMinute() *cronTask {
	return sched.Every(1).Minutes()
}

// EveryHour is a shortcut of Every(1).Hours().
func (sched *scheduler) EveryHour() *cronTask {
	return sched.Every(1).Hours()
}

// EveryDay is a shortcut of Every(1).Days().
func (sched *scheduler) EveryDay() *cronTask {
	return sched.Every(1).Days()
}

// EveryMonth is a shortcut of Every(1).Months().
func (sched *scheduler) EveryMonth() *cronTask {
	return sched.Every(1).Months()
}
