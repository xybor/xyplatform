package xyselect

import (
	"sync"

	"github.com/xybor/xyplatform/xycond"
)

type chanResult struct {
	index int
	value any
}

// See documentation of E().
type eselector struct {
	cases      []<-chan any
	center     chan chanResult
	counter    int
	isSelected bool
	mu         sync.Mutex
}

func (es *eselector) recv(c <-chan any) int {
	xycond.False(es.isSelected).Assert("Don't add case after selecting")

	n := len(es.cases)
	es.cases = append(es.cases, c)
	go func(i int) {
		// Until the channel is closed, receiving all it values then send them
		// to the center channel.
		for v := range c {
			es.center <- chanResult{i, v}
		}

		es.mu.Lock()
		defer es.mu.Unlock()

		// When all channels are closed, also closing the center channel.
		es.counter += 1
		if es.counter == len(es.cases) {
			close(es.center)
		}
	}(n)

	return n
}

func (es *eselector) send(any, any) int {
	xycond.Exit(1, "Exhausted-selector doesn't support Send")
	return 0 // Never reach, avoid syntax error
}

func (es *eselector) xselect(isDefault bool) (index int, value any, recvOk error) {
	es.isSelected = true

	var r chanResult
	var ok bool
	if isDefault {
		select {
		case r, ok = <-es.center:
		default:
			r = chanResult{-1, nil}
			ok = true
		}
	} else {
		r, ok = <-es.center
	}

	if r.index == -1 {
		return 0, nil, DefaultCaseError.New("default case")
	}

	if !ok {
		return 0, nil, ExhaustedError.New("selector is exhausted")
	}

	return r.index, r.value, nil
}
