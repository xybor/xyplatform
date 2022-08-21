package xyselect

import (
	"reflect"
	"sync"
)

// See documentation of R().
type rselector struct {
	cases []reflect.SelectCase
	mu    sync.Mutex
}

func (rs *rselector) recv(c <-chan any) int {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	var sc = reflect.SelectCase{
		Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c)}
	rs.cases = append(rs.cases, sc)

	return len(rs.cases) - 2
}

func (rs *rselector) send(c any, v any) int {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	var sc = reflect.SelectCase{
		Dir: reflect.SelectSend, Chan: reflect.ValueOf(c),
		Send: reflect.ValueOf(v),
	}

	rs.cases = append(rs.cases, sc)

	return len(rs.cases) - 2
}

func (rs *rselector) xselect(isDefault bool) (index int, value any, err error) {
	rs.mu.Lock()
	var cases = rs.cases
	rs.mu.Unlock()

	if isDefault {
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})
	}

	var i, v, ok = reflect.Select(cases)
	switch cases[i].Dir {
	case reflect.SelectSend:
		return i, nil, nil
	case reflect.SelectRecv:
		if ok {
			return i, v.Interface(), nil
		}
		return i, nil, ClosedChannelError.New("channel closed")
	default:
		// reflect.SelectDefault
		return -1, nil, nil
	}
}
