package xycond_test

import (
	"reflect"
	"testing"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xyerror"
)

func TestAssert(t *testing.T) {
	defer func() {
		var r = recover()
		if r == nil {
			t.Fail()
		}
	}()

	xycond.AssertTrue(false)
}

func TestAssertEqual(t *testing.T) {
	xycond.AssertEqual(1, 1)
	xycond.AssertNotEqual(1, 2)
}

func TestAssertLessThan(t *testing.T) {
	xycond.AssertLessThan(1, 2)
	xycond.AssertNotLessThan(1, 0)
}

func TestAssertGreaterThan(t *testing.T) {
	xycond.AssertGreaterThan(1, 0)
	xycond.AssertNotGreaterThan(1, 1)
}

func TestAssertPanic(t *testing.T) {
	xycond.AssertPanic(func() { panic("") })
	xycond.AssertNotPanic(func() {})
}

func TestAssertZero(t *testing.T) {
	xycond.AssertZero(0)
	xycond.AssertNotZero(1)
}

func TestAssertNil(t *testing.T) {
	var x *int
	xycond.AssertNil(x)
	xycond.AssertNil(nil)

	var a = make([]int, 0)
	xycond.AssertNotNil(a)
	xycond.AssertNotNil(new(int))
}

func TestAssertEmpty(t *testing.T) {
	xycond.AssertEmpty("")
	xycond.AssertEmpty([]int{})
	xycond.AssertEmpty([]int{1, 2, 3}[0:0])

	xycond.AssertNotEmpty("a")
	xycond.AssertNotEmpty([]int{1})
	xycond.AssertNotEmpty([1]int{1})
}

func TestAssertIs(t *testing.T) {
	xycond.AssertIs(1, reflect.Int)
	xycond.AssertIsNot(1, reflect.String)
}

func TestAssertSame(t *testing.T) {
	xycond.AssertSame(1, 2)
	xycond.AssertSame(1, 2, 3, 4, 5)
	xycond.AssertSame(make(chan int), make(chan int))

	xycond.AssertNotSame(1, "a")
	xycond.AssertNotSame(1, '3')
	xycond.AssertNotSame("a", 1)
	xycond.AssertNotSame(1, 2, 3, "a")
	xycond.AssertNotSame([]int{1}, [1]int{1})
}

func TestAssertWritable(t *testing.T) {
	var receive = make(<-chan int)
	var both = make(chan int)
	var send = make(chan<- int)

	xycond.AssertWritable(both)
	xycond.AssertWritable(send)
	xycond.AssertNotWritable(receive)
}

func TestAssertReadable(t *testing.T) {
	var send = make(chan<- int)
	var both = make(chan int)
	var receive = make(<-chan int)

	xycond.AssertReadable(both)
	xycond.AssertReadable(receive)
	xycond.AssertNotReadable(send)
}

func TestAssertError(t *testing.T) {
	var err = xyerror.ValueError.New("")
	xycond.AssertError(err, xyerror.ValueError)
	xycond.AssertErrorNot(err, xyerror.AssertionError)
}

func TestAssertTrue(t *testing.T) {
	xycond.AssertTrue(true)
	xycond.AssertFalse(false)
}
