package xyselect

import (
	"sync"

	"github.com/xybor/xyplatform/xycond"
)

type chanResult struct {
	index  int
	value  any
	recvOK bool
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
	defer es.mu.Unlock()

	// If this is the only live channel currently, recreate center channel.
	if es.liveCounter == 0 {
		es.center = make(chan chanResult)
	}

	es.counter++
	es.liveCounter++

	go func(i int) {
		// Until the channel is closed, receiving all it values then send them
		// to the center channel.
		for v := range c {
			es.center <- chanResult{i, v, true}
		}
		es.center <- chanResult{i, nil, false}

		es.mu.Lock()
		defer es.mu.Unlock()

		es.liveCounter--
		// If there is no more live channel, closing the center channel.
		if es.liveCounter == 0 {
			close(es.center)
		}
	}(es.counter - 1)

	return es.counter
}

func (es *eselector) send(any, any) int {
	xycond.Panic("Exhausted-selector doesn't support Send")
	return 0 // Never reach, avoid syntax error
}

func (es *eselector) xselect(isDefault bool) (index int, value any, recvOk error) {
	var r chanResult
	var ok bool
	if isDefault {
		select {
		case r, ok = <-es.center:
			if !ok {
				r = chanResult{-1, nil, false}
			}
		default:
			r = chanResult{-1, nil, false}
			ok = true
		}
	} else {
		r, ok = <-es.center
	}

	// The default case.
	if r.index == -1 {
		return -1, nil, nil
	}

	// There is no live channel, the selector is exhausted.
	if !ok {
		return 0, nil, ExhaustedError.New("selector is exhausted")
	}

	// The channel closed.
	if !r.recvOK {
		return r.index, nil, ClosedChannelError.New("channel closed")
	}

	return r.index, r.value, nil
}
