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
func NewTask(f any, params ...any) *Task {
	var fv = reflect.ValueOf(f)
	xycond.True(fv.Kind() == reflect.Func).
		Assert("Expected a function, but got %s", fv.Kind())

	return &Task{
		fv: fv, params: params,
		variadic: false, ret: make([]any, fv.Type().NumOut()),
		cb: make([]future, 0), lock: xylock.Lock{},
	}
}

// Variadic requires the task to convert returned values from []any to ...any.
// It is helpful if you are using a generic function in task (output is []any)
// and want to pass them to another callback function as variadic instead of
// slice or array.
//
// n is the number of returned values of function.
func (t *Task) Variadic(n int) *Task {
	var ftype = t.fv.Type()
	var nout = ftype.NumOut()
	xycond.True(nout == 1).Assert(
		"Expected only one output, but %d found", nout)

	var outkind = ftype.Out(0).Kind()
	xycond.True(outkind == reflect.Array || outkind == reflect.Slice).Assert(
		"Expected output as []any, but got %s", outkind)

	t.variadic = true
	t.ret = make([]any, n)
	return t
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
		xycond.Empty(params).
			Assert("Do not pass params if f was already a tasker")
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
	var cb = NewTask(f)
	t.onsuccess = append(t.onsuccess, cb)
	return cb
}

// Catch sets a callback task which will be run after this task panicked. The
// only parameter of the callback task is the panicked error.
//
// It returns the callback task.
func (t *Task) Catch(f any) *Task {
	var cb = NewTask(f)
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
		v := callFunc(t.fv, t.params, t.variadic)
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
