package xysched

import (
	"reflect"
	"time"

	"github.com/jinzhu/copier"
	"github.com/xybor/xyplatform/xycond"
)

// task is a future which runs one time.
type task struct {
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
	onsuccess []*task

	// The recovered error in case the function panicked.
	recover error

	// Callback tasks handle the panicked error if task panicked in runtime.
	onfailure []*task

	// Other callback futures.
	cb []future
}

// Task creates a future which runs function f with parameters params. This
// future runs only one time.
func Task(f any, params ...any) *task {
	var fv = reflect.ValueOf(f)
	xycond.True(fv.Kind() == reflect.Func).
		Assertf("Expected a function, but got %s", fv.Kind())

	return &task{
		fv: fv, params: params,
		variadic: false, ret: make([]any, fv.Type().NumOut()),
		cb: make([]future, 0),
	}
}

// Variadic requires the task to convert returned values from []any to ...any.
// It is helpful if you are using a generic function in task (output is []any)
// and want to pass them to another callback function as variadic instead of
// slice or array.
//
// n is the number of returned values of function.
func (t *task) Variadic(n int) *task {
	var ftype = t.fv.Type()
	var nout = ftype.NumOut()
	xycond.True(nout == 1).Assertf(
		"Expected only one output, but %d found", nout)

	var outkind = ftype.Out(0).Kind()
	xycond.True(outkind == reflect.Array || outkind == reflect.Slice).Assertf(
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
func (t *task) Callback(f any, params ...any) *task {
	cb, ok := f.(future)
	if ok {
		xycond.Empty(params).
			Assertf("Do not pass params if f was already a tasker")
	} else {
		cb = Task(f, params...)
	}

	t.cb = append(t.cb, cb)

	if t, ok := cb.(*task); ok {
		return t
	}
	return nil
}

// Then sets a callback task which will be run after this task ran successfully.
// The callback task's input parameters are the output of this task.
//
// It returns the callback task.
func (t *task) Then(f any) *task {
	var cb = Task(f)
	t.onsuccess = append(t.onsuccess, cb)
	return cb
}

// Catch sets a callback task which will be run after this task panicked. The
// only parameter of the callback task is the panicked error.
//
// It returns the callback task.
func (t *task) Catch(f any) *task {
	var cb = Task(f)
	t.onfailure = append(t.onfailure, cb)
	return cb
}

// Required method of future.
func (t *task) run() {
	if len(t.onfailure) > 0 {
		defer func() {
			if r := recover(); r != nil {
				if e, ok := r.(error); ok {
					t.recover = e
				} else {
					t.recover = CallError.Newf("%s", r)
				}
			}
		}()
	}

	v := callFunc(t.fv, t.params, t.variadic)
	copy(t.ret, v)
	t.recover = nil
}

// Required method of future.
func (t *task) next() *time.Time {
	return nil
}

// Required method of future.
func (t *task) callbacks() []future {
	var cb []future
	cb = append(cb, t.cb...)

	if t.recover != nil {
		for _, ft := range t.onfailure {
			ft.params = []any{t.recover}
			cb = append(cb, ft)
		}
	} else {
		for _, fs := range t.onsuccess {
			fs.params = t.ret
			cb = append(cb, fs)
		}
	}

	return cb
}

func (t *task) copy() future {
	var f = new(task)
	copier.Copy(f, t)
	return f
}
