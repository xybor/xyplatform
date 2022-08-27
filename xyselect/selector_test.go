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
		xycond.ExpectZero(i).Test(t)
		xycond.ExpectEqual(v, 10).Test(t)
		xycond.ExpectNil(e).Test(t)

		i, v, e = selector.Select(false)
		xycond.ExpectZero(i).Test(t)
		xycond.ExpectNil(v).Test(t)
		xycond.ExpectError(e, xyselect.ClosedChannelError).Test(t)

		_, _, e = selector.Select(false)
		xycond.ExpectError(
			e, xyselect.ExhaustedError, xyselect.ClosedChannelError).Test(t)
	}
}

func TestSelectorDefault(t *testing.T) {
	var tests = []*xyselect.Selector{xyselect.E(), xyselect.R()}

	for i := range tests {
		var selector = tests[i]

		var i, v, e = selector.Select(true)
		xycond.ExpectEqual(i, -1).Test(t)
		xycond.ExpectNil(v).Test(t)
		xycond.ExpectNil(e).Test(t)
	}
}

func TestRSelectorSend(t *testing.T) {
	var selector = xyselect.R()
	var c = make(chan int)
	var rv int
	go func() {
		rv = <-c
		xycond.ExpectEqual(rv, 10).Test(t)
	}()
	selector.Send(c, 10)
	var i, v, e = selector.Select(false)
	xycond.ExpectZero(i).Test(t)
	xycond.ExpectNil(v).Test(t)
	xycond.ExpectNil(e).Test(t)
}

func TestESelectorSend(t *testing.T) {
	var selector = xyselect.E()
	var c = make(chan int)

	xycond.ExpectPanic(func() {
		selector.Send(c, 10)
	}).Test(t)
}
