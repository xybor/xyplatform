package xyselect

import "reflect"

// See documentation of R().
type rselector struct {
	cases []reflect.SelectCase
}

func (rs *rselector) recv(c <-chan any) int {
	sc := reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c)}
	rs.cases = append(rs.cases, sc)

	return len(rs.cases) - 2
}

func (rs *rselector) send(c any, v any) int {
	sc := reflect.SelectCase{
		Dir:  reflect.SelectSend,
		Chan: reflect.ValueOf(c), Send: reflect.ValueOf(v),
	}

	rs.cases = append(rs.cases, sc)

	return len(rs.cases) - 2
}

func (rs *rselector) xselect(isDefault bool) (index int, value any, err error) {
	cases := rs.cases
	n := len(cases)
	if isDefault {
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})
	}

	i, v, ok := reflect.Select(cases)
	if i == n {
		return 0, nil, DefaultCaseError.New("default case")
	}

	if !ok {
		return i, nil, ClosedChannelError.New("channel is closed")
	}

	return i, v.Interface(), nil
}
