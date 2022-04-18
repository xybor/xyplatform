package xyselect

import (
	"reflect"
	"sync"

	"github.com/xybor/xyplatform/xycond"
)

type selector interface {
	xselect(isDefault bool) (index int, v any, err error)
	recv(c <-chan any) int
	send(c any, v any) int
}

// safeSelector is the struct supporting thread-safe, type-safe selector.
type safeSelector struct {
	selector
}

// IMPORTANT: ONLY use this selector if (1) exhausted-select, a type of select
// which is only stopped when all channels are closed, if the selector is
// exhausted, the error returned in Select() method will be StoppedError; (2)
// only receiving cases.
//
// NOTE: Don't add channel to this selector after calling Select() method,
// otherwise, it panics.
//
// ESelector is the exhausted-version selector. Its workflow is that all
// channels will send its received values to a center channel.
// Instead of calling select statement with all channels, you only need to
// receive on the center one.
//
// When ESelector adds a case of channel, it creates a goroutine which receives
// values of that channel until it is closed.
// For a received value, the goroutine will send that value to the center
// channel.
//
// The center channel is only closed if all channel-cases are closed. This is
// the reason why you must call Select until there is no any alive channel.
func E() *safeSelector {
	// Create a closed channel
	var center = make(chan chanResult)
	close(center)

	return &safeSelector{
		selector: &eselector{
			counter:     0,
			liveCounter: 0,
			center:      center,
			mu:          sync.Mutex{},
		},
	}
}

// RSelector is the reflect-version selector. It uses reflect module to handle
// customized select statment.
func R() *safeSelector {
	return &safeSelector{
		selector: &rselector{
			cases: nil,
			mu:    sync.Mutex{},
		},
	}
}

// C creates a read-only chan any from read-only chan T. This channel is only
// closed when c channel is closed.
func C[T any](c <-chan T) <-chan any {
	r := make(chan any)
	go func() {
		for v := range c {
			r <- v
		}
		close(r)
	}()

	return r
}

// Recv adds a receiving case to selector. If the channel is not the type of
// chan any, using xyselector.C to cast it. This method returns the index of
// the added case. The received value is returned by the second parameter of
// Select() method.
//
// For example:
//     c1 := make(chan any)
//     c2 := make(chan int)
//     selector.Recv(c1)
//     selector.Recv(xyselector.C(c2))
func (s *safeSelector) Recv(c <-chan any) int {
	return s.selector.recv(c)
}

// Send adds a sending case to selector. The first parameter c must be a
// writable channel. This method returns the index of the added case.
func (s *safeSelector) Send(c any, v any) int {
	cType := reflect.TypeOf(c)
	xycond.Condition(cType.Kind() == reflect.Chan).
		Assert("The first parameter must be a channel.")

	dir := cType.ChanDir()
	xycond.Condition(dir == reflect.BothDir || dir == reflect.SendDir).
		Assert("The first parameter of Send must be a writable channel.")

	cKind := cType.Elem().Kind()
	vKind := reflect.ValueOf(v).Kind()
	xycond.Condition(cKind == vKind).
		Assertf("channel and value must be the same type, but got chan %s and %s.", cKind, vKind)

	return s.selector.send(c, v)
}

// Select executes a select operation described by the list of cases in
// selector.
//
// Like the Go select statement, it blocks until at least one of the cases
// can proceed and then executes that case. If isDefault is true, it will be
// the non-blocking select.
//
// It returns the index of the chosen case, the value received, and a error of
// selector. Nil for the case of receiving is not closed. DefaultCaseError for
// the case of default. ExhaustedError if there is no more available channel in
// exhausted-selector.
func (s *safeSelector) Select(isDefault bool) (index int, v any, err error) {
	return s.selector.xselect(isDefault)
}
