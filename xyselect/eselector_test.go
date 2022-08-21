package xyselect_test

import (
	"errors"
	"testing"
	"time"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xyselect"
)

func TestESelectorRecv(t *testing.T) {
	var selector = xyselect.E()
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
	xycond.MustEqual(v, nil).Testf(t, "Expected a nil value, but got %v", v)
	xycond.MustTrue(errors.Is(e, xyselect.ClosedChannelError)).
		Testf(t, "Expected ClosedChannelError, but got %s", e)

	i, v, e = selector.Select(false)
	xycond.MustZero(i).Testf(t, "Expected a zero index, but got %d", i)
	xycond.MustEqual(v, nil).Testf(t, "Expected a nil value, but got %v", v)
	xycond.MustTrue(errors.Is(e, xyselect.ExhaustedError)).
		Testf(t, "Expected ExhaustedError, but got %s", e)
}

func TestESelectorSend(t *testing.T) {
	var selector = xyselect.E()
	var c = make(chan int)

	xycond.MustPanic(func() {
		selector.Send(c, 10)
	}).Test(t, "Expected a panic, but not found")
}

func TestESelectorRecvDefault(t *testing.T) {
	var selector = xyselect.E()
	var c = make(chan int)
	var rc = xyselect.C(c)
	selector.Recv(rc)

	go func() {
		time.Sleep(time.Millisecond)
		c <- 10
		close(c)
	}()

	// Default case
	var i, v, e = selector.Select(true)
	xycond.MustEqual(i, -1).Testf(t, "Expected index of -1, but got %d", i)
	xycond.MustNil(v).Testf(t, "Expected a nil value, but got %v", v)
	xycond.MustNil(e).Testf(t, "Expected no error, but got %v", e)

	// Normal case
	time.Sleep(10 * time.Millisecond)
	i, v, e = selector.Select(true)
	xycond.MustZero(i).Testf(t, "Expected a zero index, but got %d", i)
	xycond.MustEqual(v, 10).Testf(t, "Expected value of 10, but got %v", v)
	xycond.MustNil(e).Testf(t, "Expected no error, but got %s", e)

	// Closed case
	time.Sleep(10 * time.Millisecond)
	i, v, e = selector.Select(true)
	xycond.MustZero(i).Testf(t, "Expected a zero index, but got %d", i)
	xycond.MustEqual(v, nil).Testf(t, "Expected a nil value, but got %v", v)
	xycond.MustTrue(errors.Is(e, xyselect.ClosedChannelError)).
		Testf(t, "Expected ClosedChannelError, but got %s", e)
}
