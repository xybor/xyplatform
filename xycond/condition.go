package xycond

import (
	"fmt"
	"reflect"
)

type condition bool

// This function returns a condition which later you must call JustAssert,
// Assert, or Assertf to check it. If the condition is false, the program will
// exit with exit code 1.
func Condition(cond bool) condition {
	return condition(cond)
}

func Not(c condition) condition {
	return !c
}

// JustAssert terminates the program without a message if condition fails.
func (c condition) JustAssert() {
	if !c {
		panic(c)
	}
}

// Assert prints a formatted message and terminates the program if the
// condition fails.
func (c condition) Assert(format string, args ...any) {
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
func Zero[T number](a T) condition {
	return Condition(a == 0)
}

// NotZero returns true if a is not zero. NotZero only accepts number parameter.
func NotZero[T number](a T) condition {
	return Not(Zero(a))
}

// Nil returns true if a is nil.
func Nil(a any) condition {
	return Condition(a == nil)
}

// NotNil returns true if a is not nil.
func NotNil(a any) condition {
	return Not(Nil(a))
}

// Empty returns true if a is an empty string, slice, array, or channel.
func Empty(a any) condition {
	return Zero(reflect.ValueOf(a).Len())
}

// NotEmpty returns true if a is a not empty string, slice, array, or channel.
func NotEmpty(a any) condition {
	return Not(Empty(a))
}

// ContainM returns a condition checking if a map contains the key.
func ContainM[kt comparable, vt any](m map[kt]vt, k kt) condition {
	_, ok := m[k]
	return True(ok)
}

// NotContainM returns a condition checking if a map doesn't contain the key.
func NotContainM[kt comparable, vt any](m map[kt]vt, k kt) condition {
	return Not(ContainM(m, k))
}

// ContainA returns a condition checking if an array contains the element.
func ContainA(a any, e any) condition {
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
func NotContainA(a any, e any) condition {
	return Not(ContainA(a, e))
}

// Divisible returns a condition checking if a is divisible by b.
func Divisible[t integer](a, b t) condition {
	return Condition(a%b == 0)
}

// IsKind returns a condition checking if a value belongs to one of basic types.
func IsKind(v any, kinds ...reflect.Kind) condition {
	var kindV = reflect.TypeOf(v).Kind()
	for i := range kinds {
		if kindV == kinds[i] {
			return Condition(true)
		}
	}

	return Condition(false)
}

// SameType returns a condition checking if values are the same type.
func SameType(v ...any) condition {
	var k0 = reflect.TypeOf(v[0]).Kind()
	for i := 1; i < len(v); i++ {
		if k0 != reflect.TypeOf(v[i]).Kind() {
			return Condition(false)
		}
	}

	return Condition(true)
}

func True(b bool) condition {
	return Condition(b)
}

func False(b bool) condition {
	return Not(True(b))
}

func JustPanic() {
	c := Condition(false)
	c.JustAssert()
}

func Panic(msg string, a ...any) {
	c := Condition(false)
	c.Assert(msg, a...)
}
