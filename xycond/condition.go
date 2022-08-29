// Package xycond supports to assert or expect many conditions.
package xycond

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/xybor/xyplatform/xyerror"
)

type integer interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64
}

type number interface {
	integer | float32 | float64
}

// failer instances may be *testing.T or *testing.B.
type failer interface {
	Fail()
}

// ExpectEqual returns a true Condition if the two values are equal.
func ExpectEqual(a, b any) Condition {
	return Condition{
		result:   a == b,
		trueMsg:  fmt.Sprintf("%v == %v", a, b),
		falseMsg: fmt.Sprintf("%v != %v", a, b),
	}
}

// ExpectNotEqual returns a true Condition if the two values are not equal.
func ExpectNotEqual(a, b any) Condition {
	return ExpectEqual(a, b).revert()
}

// ExpectLessThan returns a true Condition if the first parameter is less than
// the second.
func ExpectLessThan[t number](a, b t) Condition {
	return Condition{
		result:   a < b,
		trueMsg:  fmt.Sprintf("%v is less than %v", a, b),
		falseMsg: fmt.Sprintf("%v is not less than %v", a, b),
	}
}

// ExpectNotLessThan returns a true Condition if the first parameter is not less
// than the second.
func ExpectNotLessThan[t number](a, b t) Condition {
	return ExpectLessThan(a, b).revert()
}

// ExpectGreaterThan returns a true Condition if the first parameter is greater
// than the second.
func ExpectGreaterThan[t number](a, b t) Condition {
	return Condition{
		result:   a > b,
		trueMsg:  fmt.Sprintf("%v is greater than %v", a, b),
		falseMsg: fmt.Sprintf("%v is not greater than %v", a, b),
	}
}

// ExpectNotGreaterThan returns a true Condition if the first parameter is not
// greater than the second.
func ExpectNotGreaterThan[t number](a, b t) Condition {
	return ExpectGreaterThan(a, b).revert()
}

// ExpectPanic returns a true Condition if it found a panic after calling
// function.
func ExpectPanic(f func()) (c Condition) {
	defer func() {
		var r = recover()
		if r == nil {
			c.result = false
			c.falseMsg = "no panic found"
		} else {
			c.result = true
			c.trueMsg = fmt.Sprintf("got a panic: %v", r)
		}
	}()

	f()
	return
}

// ExpectPanic returns a true Condition if it doesn't found any panic after
// calling function.
func ExpectNotPanic(f func()) (c Condition) {
	return ExpectPanic(f).revert()
}

// ExpectZero returns a true Condition if the parameter is zero.
func ExpectZero[T number](a T) Condition {
	var zero T = 0
	return ExpectEqual(a, zero)
}

// ExpectNotZero returns a true Condition if the parameter is not zero.
func ExpectNotZero[T number](a T) Condition {
	return ExpectZero(a).revert()
}

// ExpectNil returns a true Condition if the parameter is nil.
func ExpectNil(a any) Condition {
	var va = reflect.ValueOf(a)
	return Condition{
		result:   a == nil || va.IsNil(),
		trueMsg:  "got a nil value",
		falseMsg: fmt.Sprintf("got a not-nil value: %v", a),
	}
}

// ExpectNotNil returns a true Condition if the parameter is not nil.
func ExpectNotNil(a any) Condition {
	return ExpectNil(a).revert()
}

// ExpectEmpty returns a true Condition if the parameter is an empty string,
// slice, array, or channel.
func ExpectEmpty(a any) Condition {
	var va = reflect.ValueOf(a)
	return Condition{
		result:   va.Len() == 0,
		trueMsg:  fmt.Sprintf("got an empty %s", va.Kind()),
		falseMsg: fmt.Sprintf("got %s with %d element(s)", va.Kind(), va.Len()),
	}
}

// ExpectNotEmpty returns a true Condition if the parameter is not an empty
// string, slice, array, or channel.
func ExpectNotEmpty(a any) Condition {
	return ExpectEmpty(a).revert()
}

// ExpectIs returns a true Condition if value belongs to one of passed kinds.
func ExpectIs(v any, kinds ...reflect.Kind) Condition {
	var kindV = reflect.TypeOf(v).Kind()
	var result = false
	for i := range kinds {
		if kindV == kinds[i] {
			result = true
		}
	}
	return Condition{
		result:   result,
		trueMsg:  fmt.Sprintf("the value is %s", kindV),
		falseMsg: fmt.Sprintf("the value is %s", kindV),
	}
}

// ExpectIsNot returns a true Condition if value doesn't belong to any passed
// kinds.
func ExpectIsNot(v any, kinds ...reflect.Kind) Condition {
	return ExpectIs(v, kinds...).revert()
}

// ExpectSame returns a true Condition if parameters are the same type.
func ExpectSame(v ...any) Condition {
	var t0 = reflect.TypeOf(v[0])
	var result = true
	var types = []string{fmt.Sprint(t0)}
	for i := 1; i < len(v); i++ {
		var ti = reflect.TypeOf(v[i])
		if t0 != ti {
			result = false
		}
		types = append(types, fmt.Sprint(ti))
	}
	return Condition{
		result:   result,
		trueMsg:  fmt.Sprintf("all value are the same type (%s)", t0),
		falseMsg: strings.Join(types, "-"),
	}
}

// ExpectNotSameType returns a true Condition if there is at least one value
// whose type is different from the rest.
func ExpectNotSame(v ...any) Condition {
	return ExpectSame(v...).revert()
}

// ExpectWritable returns a true Condition if the channel is writable.
func ExpectWritable(c any) Condition {
	AssertIs(c, reflect.Chan)
	var dir = reflect.TypeOf(c).ChanDir()
	return Condition{
		result:   dir == reflect.BothDir || dir == reflect.SendDir,
		trueMsg:  "channel is writable",
		falseMsg: "channel is not writable",
	}
}

// ExpectNotWritable returns a true Condition if the channel is not writable.
func ExpectNotWritable(c any) Condition {
	return ExpectWritable(c).revert()
}

// ExpectReadable returns a true Condition if the channel is readable.
func ExpectReadable(c any) Condition {
	AssertIs(c, reflect.Chan)
	var dir = reflect.TypeOf(c).ChanDir()
	return Condition{
		result:   dir == reflect.BothDir || dir == reflect.RecvDir,
		trueMsg:  "channel is readable",
		falseMsg: "channel is not readable",
	}
}

// ExpectNotReadable returns a true Condition if the channel is not readable.
func ExpectNotReadable(c any) Condition {
	return ExpectReadable(c).revert()
}

// ExpectError returns a true Condition if err belongs to one of the passed
// targets.
func ExpectError(err error, targets ...error) Condition {
	var result = false
	var trueTarget error
	for i := range targets {
		if errors.Is(err, targets[i]) {
			result = true
			trueTarget = targets[i]
		}
	}
	return Condition{
		result:   result,
		trueMsg:  fmt.Sprintf("err is %v", trueTarget),
		falseMsg: fmt.Sprintf("err doesn't belong to any targets (%s)", err),
	}
}

// ExpectErrorNot returns a true Condition if the err doesn't belong to any
// targets.
func ExpectErrorNot(err error, targets ...error) Condition {
	return ExpectError(err, targets...).revert()
}

// ExpectTrue returns true if the the parameter is true.
func ExpectTrue(b bool) Condition {
	return Condition{
		result:   b,
		trueMsg:  "the condition is true",
		falseMsg: "the condition is false",
	}
}

// ExpectFalse returns a true Condition if the parameter is false.
func ExpectFalse(b bool) Condition {
	return ExpectTrue(b).revert()
}

// Panic panics with a formatted string.
func Panic(msg string, a ...any) {
	panic(xyerror.AssertionError.Newf(msg, a...))
}

// JustPanic panics immediately.
func JustPanic() {
	Panic("")
}

// Condition supports to perform actions on expectation.
type Condition struct {
	result   bool
	trueMsg  string
	falseMsg string
}

// Test will call Fail method if it is a false Condition. It is used while
// testing, with *testing.T or *testing.B.
func (c Condition) Test(f failer) {
	if !c.result {
		var _, fn, ln, ok = runtime.Caller(1)
		if ok {
			fmt.Printf("%s:%d: ", fn, ln)
		}
		fmt.Println(c.falseMsg)
		f.Fail()
	}
}

// True is performed when Condition is true.
func (c Condition) True(f func()) Condition {
	if c.result {
		f()
	}
	return c
}

// False is performed when Condition is false.
func (c Condition) False(f func()) Condition {
	if !c.result {
		f()
	}
	return c
}

// assert prints the false message and panics if it is a false Condition.
func (c Condition) assert() {
	if !c.result {
		Panic(c.falseMsg)
	}
}

// revert returns the reverse Condition.
func (c Condition) revert() Condition {
	return Condition{
		result:   !c.result,
		trueMsg:  c.falseMsg,
		falseMsg: c.trueMsg,
	}
}
