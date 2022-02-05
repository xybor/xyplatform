package xycron

import (
	"fmt"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/xybor/xyplatform/xycond"
)

type task interface {
	signal() <-chan time.Time
	start()
	stop() error
	run()
}

var taskID int64 = 1000000000

type commonTask struct {
	sched  *scheduler
	parent task
	f      any
	params []any
	id     string
}

func newCommonTask(sched *scheduler, parent task) *commonTask {
	return &commonTask{
		sched: sched, parent: parent,
		f: nil, params: nil, id: "",
	}
}

// Set parameters for the function in task.
func (ct *commonTask) Params(p ...any) *commonTask {
	ct.params = p
	return ct
}

// Do receives a function and assigns it for task. Scheduler only handles tasks
// that called Do.
func (ct *commonTask) Do(f any) {
	k := reflect.ValueOf(f).Kind()
	xycond.Condition(k == reflect.Func).Assertf("Expected a function, but got %s", k)

	ct.f = f
	ct.id = fmt.Sprintf("%s-%d", funcName(f), atomic.AddInt64(&taskID, 1))

	ct.sched.add(ct.parent)
}

func (ct commonTask) String() string {
	return ct.id
}
