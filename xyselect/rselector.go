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

	sc := reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c)}
	rs.cases = append(rs.cases, sc)

	return len(rs.cases) - 2
}

func (rs *rselector) send(c any, v any) int {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	sc := reflect.SelectCase{
		Dir:  reflect.SelectSend,
		Chan: reflect.ValueOf(c), Send: reflect.ValueOf(v),
	}

	rs.cases = append(rs.cases, sc)

	return len(rs.cases) - 2
}

func (rs *rselector) xselect(isDefault bool) (index int, value any, err error) {
	rs.mu.Lock()
	cases := rs.cases
	rs.mu.Unlock()

	if isDefault {
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})
	}

	i, v, ok := reflect.Select(cases)
	if i == len(cases)-1 && isDefault {
		return 0, nil, DefaultCaseError.New("default case")
	}
	if !ok {
		return i, nil, ClosedChannelError.New("channel is closed")
	}

	return i, v.Interface(), nil
}
