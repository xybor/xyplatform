package xysched

import (
	"reflect"
	"time"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xylock"
)

// Task is a future which runs one time.
type Task struct {
	// The function's reflect value.
	fv reflect.Value

	// The parameters of function.
	params []any

	// In case output is []any, it will be converted to ...any if variadic is
	// true. It is helpful to pass them as parameters of another function.
	variadic bool

	// The returned values in case the function ran successfully.
	ret []any

	// Callback tasks handle the returned values if task ran successfully.
	onsuccess []*Task

	// The recovered error in case the function panicked.
	recover error

	// Callback tasks handle the panicked error if task panicked in runtime.
	onfailure []*Task

	// Other callback futures.
	cb []future

	// Avoid task to be called simultaneously.
	lock xylock.Lock
}

// NewTask creates a future which runs function f with parameters params. This
// future runs only one time.
func NewTask(f any, p ...any) *Task {
	var fv = reflect.ValueOf(f)
	var ft = fv.Type()
	var params = anyArrayToTypeArray(p)
	checkParam(getFuncIn(ft), params, ft.IsVariadic())

	return &Task{
		fv: fv, params: p,
		variadic: false, ret: make([]any, ft.NumOut()),
		cb: make([]future, 0), lock: xylock.Lock{},
	}
}

// newPlaceholderTask creates a Task whose parameters is determined later.
func newPlaceholderTask(f any) *Task {
	var fv = reflect.ValueOf(f)
	return &Task{
		fv: fv, params: nil,
		variadic: false, ret: make([]any, fv.Type().NumOut()),
		cb: make([]future, 0), lock: xylock.Lock{},
	}
}

// Callback sets a callback future which will run after the task completed. The
// parameter params only should be used if f is a function. In case f was
// already a future, DO NOT use.
//
// It returns the callback task if you passed a function or task, otherwise,
// nil.
func (t *Task) Callback(f any, params ...any) *Task {
	cb, ok := f.(future)
	if ok {
		xycond.MustEmpty(params).
			Assert("do not pass params if f was already a future")
	} else {
		cb = NewTask(f, params...)
	}

	t.cb = append(t.cb, cb)

	if t, ok := cb.(*Task); ok {
		return t
	}
	return nil
}

// Then sets a callback task which will be run after this task ran successfully.
// The callback task's input parameters are the output of this task.
//
// It returns the callback task.
func (t *Task) Then(f any) *Task {
	var ft = reflect.TypeOf(f)
	checkParam(getFuncIn(ft), getFuncOut(t.fv.Type()), ft.IsVariadic())
	var cb = newPlaceholderTask(f)
	t.onsuccess = append(t.onsuccess, cb)
	return cb
}

// Catch sets a callback task which will be run after this task panicked. The
// only parameter of the callback task is the panicked error.
//
// It returns the callback task.
func (t *Task) Catch(f any) *Task {
	var ft = reflect.TypeOf(f)
	xycond.MustTrue(ft.NumIn() == 1).
		Assert("catch function must have exact one parameter")

	var errtype = reflect.TypeOf((*error)(nil)).Elem()
	xycond.MustTrue(ft.In(0).AssignableTo(errtype)).
		Assert("catch function's parameter must be an error")

	var cb = newPlaceholderTask(f)
	t.onfailure = append(t.onfailure, cb)
	return cb
}

// Required method of future.
func (t *Task) run() {
	if len(t.onfailure) > 0 {
		defer func() {
			if r := recover(); r != nil {
				var e, ok = r.(error)
				if !ok {
					e = CallError.New("%s", r)
				}
				t.lock.LockFunc(func() { t.recover = e })
			}
		}()
	}

	t.lock.LockFunc(func() {
		v := callFunc(t.fv, t.params)
		copy(t.ret, v)
		t.recover = nil
	})
}

// Required method of future.
func (t *Task) next() *time.Time {
	return nil
}

// Required method of future.
func (t *Task) callbacks() []future {
	var cb []future
	cb = append(cb, t.cb...)

	var rdata = t.lock.RLockFunc(func() any { return t.recover })
	if rdata != nil {
		for i := range t.onfailure {
			t.onfailure[i].params = []any{rdata}
			cb = append(cb, t.onfailure[i])
		}
	} else {
		for i := range t.onsuccess {
			t.onsuccess[i].params = t.ret
			cb = append(cb, t.onsuccess[i])
		}
	}

	return cb
}
