package xycond

import (
	"fmt"
	"reflect"
)

// tester instances represent testing.T or testing.B.
type tester interface {
	Error(...any)
	Errorf(string, ...any)
}

// Condition is a type of bool which later you must call JustAssert or Assert.
// If the condition is false, the program will be panicked.
type Condition bool

// not returns the reverse condition of the origin.
func not(c Condition) Condition {
	return !c
}

// JustAssert panics the program without a message if condition fails.
func (c Condition) JustAssert() {
	if !c {
		panic("Something was wrong")
	}
}

// Assert prints a formatted message and panics the program if the condition
// fails.
func (c Condition) Assert(format string, args ...any) {
	if !c {
		panic(fmt.Sprintf(format, args...))
	}
}

// Test will call t.Error if condition is false.
func (c Condition) Test(t tester, args ...any) {
	if !c {
		t.Error(args...)
	}
}

// Testf will call t.Errorf if condition is false.
func (c Condition) Testf(t tester, format string, args ...any) {
	if !c {
		t.Errorf(format, args...)
	}
}

type integer interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64
}

type number interface {
	integer | float32 | float64
}

// MustEqual returns true if two values are the same.
func MustEqual(a, b any) Condition {
	MustSameType(a, b).Assert("two parameters must be the same type")
	return a == b
}

// MustNotEqual returns true if two values are not the same.
func MustNotEqual(a, b any) Condition {
	return not(MustEqual(a, b))
}

func MustPanic(f func()) (c Condition) {
	defer func() {
		var r = recover()
		if r == nil {
			c = false
		} else {
			c = true
		}
	}()

	f()
	return
}

func MustNotPanic(f func()) (c Condition) {
	return not(MustPanic(f))
}

// MustZero returns true if a is zero. MustZero only accepts number parameter.
func MustZero[T number](a T) Condition {
	return a == 0
}

// MustNotZero returns true if a is not zero. MustNotZero only accepts number
// parameter.
func MustNotZero[T number](a T) Condition {
	return not(MustZero(a))
}

// MustNil returns true if a is nil.
func MustNil(a any) Condition {
	if a == nil {
		return true
	}
	if MustNotBe(a, reflect.Pointer) {
		return false
	}
	return Condition(reflect.ValueOf(a).IsNil())
}

// MustNotNil returns true if a is not nil.
func MustNotNil(a any) Condition {
	return not(MustNil(a))
}

// MustEmpty returns true if a is an empty string, slice, array, or channel.
func MustEmpty(a any) Condition {
	MustBeLenghtType(a).Assert("parameter must be a length type")
	return MustZero(reflect.ValueOf(a).Len())
}

// MustNotEmpty returns true if a is a not empty string, slice, array, or
// channel.
func MustNotEmpty(a any) Condition {
	MustBeLenghtType(a).Assert("parameter must be a length type")
	return MustNotZero(reflect.ValueOf(a).Len())
}

// MustContainM returns true if map contains the key.
func MustContainM[kt comparable, vt any](m map[kt]vt, k kt) Condition {
	_, ok := m[k]
	return MustTrue(ok)
}

// MustNotContainM returns true if map doesn't contain the key.
func MustNotContainM[kt comparable, vt any](m map[kt]vt, k kt) Condition {
	return not(MustContainM(m, k))
}

// MustContainA returns true if array or slice contains the element.
func MustContainA(a any, e any) Condition {
	var va = reflect.ValueOf(a)
	MustBe(a, reflect.Array, reflect.Slice).
		Assert("expected an array or slice, but got %s", va.Kind())

	for i := 0; i < va.Len(); i++ {
		if va.Index(i).Interface() == e {
			return true
		}
	}
	return false
}

// MustNotContainA returns true if array doesn't contains the element.
func MustNotContainA(a any, e any) Condition {
	return not(MustContainA(a, e))
}

// MustBe returns true if value belongs to one of basic types.
func MustBe(v any, kinds ...reflect.Kind) Condition {
	var kindV = reflect.TypeOf(v).Kind()
	for i := range kinds {
		if kindV == kinds[i] {
			return true
		}
	}
	return false
}

// MustNotBe returns true if value doesn't belong to all of basic types.
func MustNotBe(v any, kinds ...reflect.Kind) Condition {
	return not(MustBe(v, kinds...))
}

// MustBeLengthType returns true if a is a string, slice, array, or chan.
func MustBeLenghtType(a any) Condition {
	return MustBe(a, reflect.String, reflect.Slice, reflect.Array, reflect.Chan)
}

// MustBeElemType returns true if a is LengthType or Pointer.
func MustBeElemType(a any) Condition {
	return MustBeLenghtType(a) || MustBe(a, reflect.Pointer)
}

// MustSameType returns true if values are the same type.
func MustSameType(v ...any) Condition {
	var t0 = reflect.TypeOf(v[0])
	for i := 1; i < len(v); i++ {
		if t0 != reflect.TypeOf(v[i]) {
			return false
		}
	}
	return true
}

// MustWritableChan returns true if channel is writable.
func MustWritableChan(c any) Condition {
	MustBe(c, reflect.Chan).Assert("c must be a channel")
	var dir = reflect.TypeOf(c).ChanDir()
	return dir == reflect.BothDir || dir == reflect.SendDir
}

// MustReadableChan returns true if channel is readable.
func MustReadableChan(c any) Condition {
	MustBe(c, reflect.Chan).Assert("c must be a channel")
	var dir = reflect.TypeOf(c).ChanDir()
	return dir == reflect.BothDir || dir == reflect.RecvDir
}

// MustTrue checks if b is true.
func MustTrue(b bool) Condition {
	return Condition(b)
}

// MustFalse checks if b is false.
func MustFalse(b bool) Condition {
	return not(MustTrue(b))
}

// JustPanic panics immediately.
func JustPanic() {
	MustTrue(false).JustAssert()
}

// Panic panics with a formatted string.
func Panic(msg string, a ...any) {
	MustTrue(false).Assert(msg, a...)
}
