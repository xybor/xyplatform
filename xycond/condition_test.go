package xycond_test

import (
	"reflect"
	"testing"

	"github.com/xybor/xyplatform/xycond"
)

func TestMustEqual(t *testing.T) {
	if xycond.MustEqual(1, 2) {
		t.Error("MustEqual failed")
	}

	if xycond.MustNotEqual(1, 1) {
		t.Error("MustNotEqual failed")
	}
}

func TestMustEqualDiffType(t *testing.T) {
	defer func() {
		var r = recover()
		if r == nil {
			t.Error("Expect a panic, but not found")
		}
	}()

	xycond.MustEqual(1, true)
}

func TestMustPanic(t *testing.T) {
	if xycond.MustPanic(func() {}) {
		t.Error("MustPanic failed")
	}

	if xycond.MustNotPanic(func() { panic("") }) {
		t.Error("MustNotPanic failed")
	}
}

func TestMustZero(t *testing.T) {
	if xycond.MustZero(1) {
		t.Error("MustZero failed")
	}

	if xycond.MustNotZero(0) {
		t.Error("MustNotZero failed")
	}
}

func TestMustNil(t *testing.T) {
	if xycond.MustNil("") {
		t.Error("MustNil failed with empty string")
	}

	if xycond.MustNil(new(int)) {
		t.Error("MustNil failed with new(int)")
	}

	if xycond.MustNil(0) {
		t.Error("MustNil failed with zero")
	}

	if xycond.MustNotNil(nil) {
		t.Error("MustNotNil failed with nil")
	}

	var x *int = nil
	if xycond.MustNotNil(x) {
		t.Error("MustNotNil failed with nil pointer of int")
	}
}

func TestMustEmpty(t *testing.T) {
	if xycond.MustEmpty("a") {
		t.Error("MustEmpty failed with non-empty string")
	}

	if xycond.MustEmpty([]int{1}) {
		t.Error("MustEmpty failed with non-empty slice")
	}

	if xycond.MustEmpty([1]int{1}) {
		t.Error("MustEmpty failed with non-empty array")
	}

	if xycond.MustNotEmpty("") {
		t.Error("MustNotEmpty failed with empty string")
	}

	if xycond.MustNotEmpty([]int{}) {
		t.Error("MustNotEmpty failed with empty slice")
	}

	if xycond.MustNotEmpty([]int{1, 2, 3}[0:0]) {
		t.Error("MustNotEmpty failed with empty array")
	}
}

func TestMustEmptyWithNonLengthType(t *testing.T) {
	defer func() {
		var r = recover()
		if r == nil {
			t.Error("Expect a panic, but not found")
		}
	}()

	xycond.MustEmpty(1)
}

func TestMustContainM(t *testing.T) {
	var m = map[int]string{
		1: "foo",
		2: "bar",
	}

	if xycond.MustContainM(m, 0) {
		t.Error("MustContainM failed with not-existed key")
	}

	if xycond.MustNotContainM(m, 1) {
		t.Error("MustNotContainM failed with an existed key")
	}
}

func TestMustContainASlice(t *testing.T) {
	var a = []int{1, 2, 3, 4, 5}

	if xycond.MustContainA(a, 0) {
		t.Error("MustContainA failed with not-existed element")
	}

	if xycond.MustNotContainA(a, 1) {
		t.Error("MustNotContainA failed with an existed element")
	}
}

func TestMustContainAArray(t *testing.T) {
	var a = [5]int{1, 2, 3, 4, 5}

	if xycond.MustContainA(a, 0) {
		t.Error("MustContainA failed with not-existed element")
	}

	if xycond.MustNotContainA(a, 1) {
		t.Error("MustNotContainA failed with an existed element")
	}
}

func TestMustContainAErrorType(t *testing.T) {
	defer func() {
		var r = recover()
		if r == nil {
			t.Error("Expect a panic, but not found")
		}
	}()

	if xycond.MustContainA("abc", 0) {
		t.Error("MustContainA failed with not-existed element")
	}
}

func TestMustBe(t *testing.T) {
	var tests = map[any]reflect.Kind{
		1:     reflect.Int,
		"foo": reflect.String,
		1.1:   reflect.Float64,
		true:  reflect.Bool,
		'c':   reflect.Int32,
	}

	for value, kind := range tests {
		if xycond.MustBe(value, kind+1) {
			t.Errorf("MustBe failed with %v", value)
		}

		if xycond.MustNotBe(value, kind) {
			t.Errorf("MustNotBe failed with %v", value)
		}
	}
}

func TestMustBeLengthType(t *testing.T) {
	if xycond.MustBeLenghtType(123) {
		t.Error("MustBeLengthType failed with int")
	}

	if xycond.MustBeLenghtType(true) {
		t.Error("MustBeLengthType failed with bool")
	}

	if xycond.MustBeLenghtType(new(int)) {
		t.Error("MustBeLengthType failed with pointer of int")
	}

	if !xycond.MustBeLenghtType([]int{}) {
		t.Error("MustBeLengthType failed with int slice")
	}

	if !xycond.MustBeLenghtType([1]float32{}) {
		t.Error("MustBeLengthType failed with float32 array")
	}

	if !xycond.MustBeLenghtType("foo") {
		t.Error("MustBeLengthType failed with string")
	}

	if !xycond.MustBeLenghtType(make(chan rune)) {
		t.Error("MustBeLengthType failed with chan rune")
	}
}

func TestMustBeElemType(t *testing.T) {
	if xycond.MustBeElemType(1) {
		t.Error("MustBeElemType failed with int")
	}

	if !xycond.MustBeElemType(new(int)) {
		t.Error("MustBeElemType failed with pointer of int")
	}

	if !xycond.MustBeElemType("bar") {
		t.Error("MustBeElemType failed with string")
	}
}

func TestMustSameType(t *testing.T) {
	if xycond.MustSameType(1, "a") {
		t.Error("MustSameType failed with int and string")
	}

	if xycond.MustSameType(1, '3') {
		t.Error("MustSameType failed with int and rune")
	}

	if xycond.MustSameType("a", 1) {
		t.Error("MustSameType failed with string and int")
	}

	if xycond.MustSameType(1, 2, 3, "a") {
		t.Error("MustSameType failed with multiple int and string")
	}

	if xycond.MustSameType([]int{1}, [1]int{1}) {
		t.Error("MustSameType failed with slice and array")
	}

	if !xycond.MustSameType(1, 2) {
		t.Error("MustSameType failed with only int")
	}
	if !xycond.MustSameType(1, 2, 3, 4, 5) {
		t.Error("MustSameType failed with multiple int")
	}

	if !xycond.MustSameType(make(chan int), make(chan int)) {
		t.Error("MustSameType failed with only chan int")
	}
}

func TestMustWritableChan(t *testing.T) {
	var receive = make(<-chan int)
	var both = make(chan int)
	var send = make(chan<- int)

	if xycond.MustWritableChan(receive) {
		t.Error("MustWritableChan failed with receive-only (read-only)")
	}

	if !xycond.MustWritableChan(both) {
		t.Error("MustWritableChan failed with both direction channel")
	}

	if !xycond.MustWritableChan(send) {
		t.Error("MustWritableChan failed with send-only (write-only)")
	}
}

func TestMustReadableChan(t *testing.T) {
	var send = make(chan<- int)
	var both = make(chan int)
	var receive = make(<-chan int)

	if xycond.MustReadableChan(send) {
		t.Error("MustReadableChan failed with send-only (write-only)")
	}

	if !xycond.MustReadableChan(both) {
		t.Error("MustReadableChan failed with both direction channel")
	}

	if !xycond.MustReadableChan(receive) {
		t.Error("MustReadableChan failed with receive-only (read-only)")
	}
}

func TestConditionAssert(t *testing.T) {
	defer func() {
		var r = recover()
		if r == nil {
			t.Error("Expect a panic, but not found")
		}
	}()
	xycond.MustTrue(false).Assert("should panic")
}

func TestConditionAssertNotOccur(t *testing.T) {
	xycond.MustTrue(true).Assert("should not panic")
}

func TestConditionJustAssert(t *testing.T) {
	defer func() {
		var r = recover()
		if r == nil {
			t.Error("Expect a panic, but not found")
		}
	}()
	xycond.MustTrue(false).JustAssert()
}

func TestConditionJustAssertNotOccur(t *testing.T) {
	xycond.MustTrue(true).JustAssert()
}

type panictester struct{}

func (*panictester) Error(a ...any) {
	xycond.JustPanic()
}

func (*panictester) Errorf(m string, a ...any) {
	xycond.Panic(m, a...)
}

func TestConditionTest(t *testing.T) {
	defer func() {
		var r = recover()
		if r == nil {
			t.Error("Expect a panic, but not found")
		}
	}()
	xycond.MustTrue(false).Test(&panictester{}, "should panic")
}

func TestConditionTestNotOccur(t *testing.T) {
	xycond.MustTrue(true).Test(&panictester{}, "should not panic")
}

func TestConditionTestf(t *testing.T) {
	defer func() {
		var r = recover()
		if r == nil {
			t.Error("Expect a panic, but not found")
		}
	}()
	xycond.MustTrue(false).Testf(&panictester{}, "should panic")
}

func TestConditionTestfNotOccur(t *testing.T) {
	xycond.MustFalse(false).Testf(&panictester{}, "should not panic")
}
