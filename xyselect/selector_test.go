package xyselect_test

import (
	"testing"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xyselect"
)

func TestESelectorRecv(t *testing.T) {
	var tests = []*xyselect.Selector{xyselect.E(), xyselect.R()}

	for i := range tests {
		var selector = tests[i]
		var c = make(chan int)
		var rc = xyselect.C(c)
		selector.Recv(rc)

		go func() {
			c <- 10
			close(c)
		}()

		var i, v, e = selector.Select(false)
		xycond.MustZero(i).Testf(t, "Expected a zero index, but got %d", i)
		xycond.MustEqual(v, 10).Testf(t, "Expected value of 10, but got %v", v)
		xycond.MustNil(e).Testf(t, "Expected no error, but got %s", e)

		i, v, e = selector.Select(false)
		xycond.MustZero(i).Testf(t, "Expected a zero index, but got %d", i)
		xycond.MustNil(v).Testf(t, "Expected a nil value, but got %v", v)
		xycond.ErrorMustBe(e, xyselect.ClosedChannelError).
			Testf(t, "Expected ClosedChannelError, but got %s", e)

		_, _, e = selector.Select(false)
		xycond.ErrorMustBe(
			e, xyselect.ExhaustedError, xyselect.ClosedChannelError).Testf(t,
			"Expected ExhaustedError or ClosedChannelError, but got %s", e)
	}
}

func TestSelectorDefault(t *testing.T) {
	var tests = []*xyselect.Selector{xyselect.E(), xyselect.R()}

	for i := range tests {
		var selector = tests[i]

		var i, v, e = selector.Select(true)
		xycond.MustEqual(i, -1).Testf(t, "Expected index of -1, but got %d", i)
		xycond.MustNil(v).Testf(t, "Expected a nil value, but got %v", v)
		xycond.MustNil(e).Testf(t, "Expected no error, but got %v", e)
	}
}

func TestRSelectorSend(t *testing.T) {
	var selector = xyselect.R()
	var c = make(chan int)
	var rv int
	go func() {
		rv = <-c
		xycond.MustEqual(rv, 10).
			Testf(t, "Expected received value of 10, but got %v", rv)
	}()
	selector.Send(c, 10)
	var i, v, e = selector.Select(false)
	xycond.MustZero(i).Testf(t, "Expected a zero index, but got %d", i)
	xycond.MustNil(v).Testf(t, "Expected a nil value, but got %v", v)
	xycond.MustNil(e).Testf(t, "Expected no error, but got %s", e)
}

func TestESelectorSend(t *testing.T) {
	var selector = xyselect.E()
	var c = make(chan int)

	xycond.MustPanic(func() {
		selector.Send(c, 10)
	}).Test(t, "Expected a panic, but not found")
}
