package xycond

import (
	"fmt"
	"reflect"
)

// Condition is a type of bool which later you must call JustAssert or Assert.
// If the condition is false, the program will be panicked.
type Condition bool

func Not(c Condition) Condition {
	return !c
}

// JustAssert terminates the program without a message if condition fails.
func (c Condition) JustAssert() {
	if !c {
		panic("Something was wrong")
	}
}

// Assert prints a formatted message and terminates the program if the
// condition fails.
func (c Condition) Assert(format string, args ...any) {
	if !c {
		panic(fmt.Sprintf(format, args...))
	}
}

type integer interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64
}

type number interface {
	integer | float32 | float64
}

// Zero returns true if a is zero. Zero only accepts number parameter.
func Zero[T number](a T) Condition {
	return Condition(a == 0)
}

// NotZero returns true if a is not zero. NotZero only accepts number parameter.
func NotZero[T number](a T) Condition {
	return Not(Zero(a))
}

// Nil returns true if a is nil.
func Nil(a any) Condition {
	return Condition(a == nil)
}

// NotNil returns true if a is not nil.
func NotNil(a any) Condition {
	return Not(Nil(a))
}

// Empty returns true if a is an empty string, slice, array, or channel.
func Empty(a any) Condition {
	return Zero(reflect.ValueOf(a).Len())
}

// NotEmpty returns true if a is a not empty string, slice, array, or channel.
func NotEmpty(a any) Condition {
	return Not(Empty(a))
}

// ContainM returns a condition checking if a map contains the key.
func ContainM[kt comparable, vt any](m map[kt]vt, k kt) Condition {
	_, ok := m[k]
	return True(ok)
}

// NotContainM returns a condition checking if a map doesn't contain the key.
func NotContainM[kt comparable, vt any](m map[kt]vt, k kt) Condition {
	return Not(ContainM(m, k))
}

// ContainA returns a condition checking if an array contains the element.
func ContainA(a any, e any) Condition {
	var v = reflect.ValueOf(a)
	var kind = v.Kind()
	Condition(kind == reflect.Array || kind == reflect.Slice).
		Assert("Expected an array or slice, but got %s", kind)

	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == e {
			return Condition(true)
		}
	}
	return Condition(false)
}

// NotContainA returns a condition checking if an array doesn't contains the
// element.
func NotContainA(a any, e any) Condition {
	return Not(ContainA(a, e))
}

// IsKind returns a condition checking if a value belongs to one of basic types.
func IsKind(v any, kinds ...reflect.Kind) Condition {
	var kindV = reflect.TypeOf(v).Kind()
	for i := range kinds {
		if kindV == kinds[i] {
			return Condition(true)
		}
	}

	return Condition(false)
}

// SameType returns a condition checking if values are the same type.
func SameType(v ...any) Condition {
	var t0 = reflect.TypeOf(v[0])
	for i := 1; i < len(v); i++ {
		if t0 != reflect.TypeOf(v[i]) {
			return Condition(false)
		}
	}
	return Condition(true)
}

// IsWritableChan returns true if channel is writable.
func IsWritableChan(c any) Condition {
	var dir = reflect.TypeOf(c).ChanDir()
	return dir == reflect.BothDir || dir == reflect.SendDir
}

// IsReadableChan returns true if channel is readable.
func IsReadableChan(c any) Condition {
	var dir = reflect.TypeOf(c).ChanDir()
	return dir == reflect.BothDir || dir == reflect.RecvDir
}

// True checks if b is true.
func True(b bool) Condition {
	return Condition(b)
}

// False checks if b is false
func False(b bool) Condition {
	return Not(True(b))
}

// JustPanic panics immediately. It is equivalent to
// Condition(false).JustAssert()
func JustPanic() {
	Condition(false).JustAssert()
}

// Panic panics with a formatted string.
func Panic(msg string, a ...any) {
	c := Condition(false)
	c.Assert(msg, a...)
}
