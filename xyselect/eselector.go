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
	// The total number of channels.
	counter int

	// The number of live (not closed yet) channels.
	liveCounter int

	// The locker using for accessing to liveCounter between gorountines.
	mu sync.Mutex

	// The channel is using for aggregating other channels.
	center chan chanResult
}

func (es *eselector) recv(c <-chan any) int {
	es.mu.Lock()
	es.mu.Unlock()

	// If this is the only live channel currently, recreate center channel.
	if es.liveCounter == 0 {
		es.center = make(chan chanResult)
	}

	es.counter += 1
	es.liveCounter += 1

	go func(i int) {
		// Until the channel is closed, receiving all it values then send them
		// to the center channel.
		for v := range c {
			es.center <- chanResult{i, v}
		}

		es.mu.Lock()
		es.mu.Unlock()

		es.liveCounter -= 1
		// If there is no more live channel, closing the center channel.
		if es.liveCounter == 0 {
			close(es.center)
		}
	}(es.counter - 1)

	return es.counter
}

func (es *eselector) send(any, any) int {
	xycond.Exit(1, "Exhausted-selector doesn't support Send")
	return 0 // Never reach, avoid syntax error
}

func (es *eselector) xselect(isDefault bool) (index int, value any, recvOk error) {
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
