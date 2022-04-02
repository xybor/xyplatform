package xycond

import (
	"fmt"
	"os"
	"runtime"
)

type condition struct {
	cond  bool
	depth int
	code  int
}

// This function returns a condition which later you must call JustAssert,
// Assert, or Assertf to check it. If the condition is false, the program will
// exit with exit code 1.
func Condition(cond bool) condition {
	return condition{cond: cond, depth: 1, code: 1}
}

func Not(c condition) condition {
	c.cond = !c.cond
	return c
}

func (c condition) Code(code int) condition {
	c.code = code
	return c
}

// JustAssert terminates the program without a message if condition fails.
func (c condition) JustAssert() {
	if !c.cond {
		_, fn, ln, ok := runtime.Caller(c.depth)
		if !ok {
			fn = "???"
			ln = 0
		}

		os.Stderr.Write([]byte(fmt.Sprintf("The program exits at %s:%d\n", fn, ln)))
		os.Exit(c.code)
	}
}

// Assert prints a message and terminates the program if the condition fails.
func (c condition) Assert(msg string) {
	if !c.cond {
		_, fn, ln, ok := runtime.Caller(c.depth)
		if !ok {
			fn = "???"
			ln = 0
		}

		os.Stderr.Write([]byte(fmt.Sprintf("The program exits at %s:%d\n", fn, ln)))
		os.Stderr.Write([]byte(msg + "\n"))
		os.Exit(c.code)
	}
}

// Assertf prints a formatted message and terminates the program if the
// condition fails.
func (c condition) Assertf(format string, args ...any) {
	c.depth += 1
	c.Assert(fmt.Sprintf(format, args...) + "\n")
}

type number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
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

// StringEmpty returns true if s is an empty string.
func StringEmpty(s string) condition {
	return Condition(len(s) == 0)
}

// StringNotEmpty returns true if s is not an empty string.
func StringNotEmpty(s string) condition {
	return Not(StringEmpty(s))
}

// SliceEmpty returns true if s is an empty slice.
func SliceEmpty[T any](s []T) condition {
	return Condition(len(s) == 0)
}

// SliceNotEmpty returns true if s is not an empty slice.
func SliceNotEmpty[T any](s []T) condition {
	return Not(SliceEmpty(s))
}

// MapEmpty returns true if s is an empty map.
func MapEmpty[K comparable, T any](m map[K]T) condition {
	return Condition(len(m) == 0)
}

// MapNotEmpty returns true if s is not an empty map.
func MapNotEmpty[K comparable, T any](m map[K]T) condition {
	return Not(MapEmpty(m))
}

// Contains returns a condition checking whether or not the map m contains the
// key.
func Contains[kt comparable, vt any](m map[kt]vt, key kt) condition {
	_, ok := m[key]
	return True(ok)
}

// NotContains returns a condition checking whether or not the map m contains
// the key.
func NotContains[kt comparable, vt any](m map[kt]vt, key kt) condition {
	_, ok := m[key]
	return False(ok)
}

// Divisible returns a condition checking if a is divisible by b.
func Divisible(a, b int) condition {
	return Condition(a%b == 0)
}

func True(b bool) condition {
	return Condition(b)
}

func False(b bool) condition {
	return Not(True(b))
}

func JustExit(code int) {
	c := True(false).Code(code)
	c.depth += 1
	c.JustAssert()
}

func Exit(code int, msg string) {
	c := True(false).Code(code)
	c.depth += 1
	c.Assert(msg)
}

func Exitf(code int, msg string, a ...any) {
	c := True(false).Code(code)
	c.depth += 1
	c.Assertf(msg, a...)
}
