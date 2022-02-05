package xycron

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xybor/xyplatform/xyerror"
)

func callFunc(f any, p ...any) ([]reflect.Value, error) {
	fv := reflect.ValueOf(f)

	var ninput = fv.Type().NumIn()
	if fv.Type().IsVariadic() {
		if len(p) < ninput-1 {
			return nil, xyerror.ParameterError.Newf(
				"expected at least %d, but got %d parameters", ninput-1, len(p))
		}
	} else {
		if len(p) != ninput {
			return nil, xyerror.ParameterError.Newf(
				"expected %d, but got %d parameters", ninput, len(p))
		}
	}

	in := make([]reflect.Value, len(p))
	for k, param := range p {
		in[k] = reflect.ValueOf(param)
	}

	return fv.Call(in), nil
}

func wdListToString(wd ...time.Weekday) []string {
	il := make([]int, len(wd))
	for i, v := range wd {
		il[i] = int(v)
	}

	return intListToString(il...)
}

func monListToString(m ...time.Month) []string {
	il := make([]int, len(m))
	for i, v := range m {
		il[i] = int(v)
	}

	return intListToString(il...)
}

func intListToString(i ...int) []string {
	return strings.Split(strings.Trim(fmt.Sprint(i), "[]"), " ")
}

func funcName(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

var atoiWeekdays = map[string]int{
	"sun": 0, "mon": 1, "tue": 2, "wed": 3, "thu": 4, "fri": 5, "sat": 6,
}

// atowd converts a string to weekday as int.
func atowd(s string) (int, error) {
	if v, err := strconv.Atoi(s); err == nil && v >= 0 && v < 7 {
		return v, nil
	}

	s = strings.ToLower(s)
	if v, ok := atoiWeekdays[s]; ok {
		return v, nil
	}

	return 0, xyerror.ValueError.Newf("cannot convert %s to weekday", s)
}

func itowd(wd time.Weekday) string {
	return strconv.Itoa(int(wd))
}

var atoiMonths = map[string]int{
	"jan": 0, "feb": 1, "mar": 2, "apr": 3, "may": 4, "jun": 5,
	"jul": 6, "aug": 7, "sep": 8, "oct": 9, "nov": 10, "dec": 11,
}

// atom converts a string to month as int.
func atom(s string) (int, error) {
	if v, err := strconv.Atoi(s); err == nil && v >= 0 && v < 12 {
		return v, nil
	}

	s = strings.ToLower(s)
	if v, ok := atoiMonths[s]; ok {
		return v, nil
	}

	return 0, xyerror.ValueError.Newf("cannot convert %s to month", s)
}

func itom(m time.Month) string {
	return strconv.Itoa(int(m))
}

// atoui converts a string to unsigned int.
func atoui(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, xyerror.ValueError.Newf("%v", err)
	}

	if v < 0 {
		return 0, xyerror.ValueError.New("negative number")
	}

	return v, nil
}

func parseListToken(s string, atoi func(string) (int, error)) ([]int, error) {
	if len(s) == 0 {
		return nil, xyerror.ValueError.New("s is empty")
	}

	splitted := strings.Split(s, ",")
	if len(splitted) < 2 {
		return nil, xyerror.ValueError.Newf("%s is not a list token", s)
	}

	var result []int
	result = make([]int, len(splitted))

	var v int
	var err error
	for i, token := range splitted {
		v, err = atoi(token)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}

	return result, nil
}

func parseRange(s string, max int, atoi func(string) (int, error)) (start, end int, err error) {
	splitted := strings.Split(s, "-")
	if len(splitted) > 2 {
		err = xyerror.ValueError.Newf("%s is not a range token", s)
		return
	}

	var e error
	if len(splitted) == 2 {
		start, e = atoi(splitted[0])
		if e != nil {
			err = xyerror.ValueError.Newf("can not convert %s to int: %s", splitted[0], e)
			return
		}

		end, e = atoi(splitted[1])
		if e != nil {
			err = xyerror.ValueError.Newf("can not convert %s to int: %s", splitted[1], e)
			return
		}
	} else if splitted[0] == "*" {
		start = 0
		end = max
	} else {
		start, e = atoi(splitted[0])
		end = start

		if e != nil {
			err = xyerror.ValueError.Newf("can not convert %s to int: %s", splitted[0], e)
			return
		}
	}

	return start, end, nil
}

func parseRangeToken(s string, max int, atoi func(string) (int, error)) ([]int, error) {
	if len(s) == 0 {
		return nil, xyerror.ValueError.New("s is empty")
	}

	splitted := strings.Split(s, "/")
	if len(splitted) > 2 {
		return nil, xyerror.ValueError.Newf("%s is not a range token", s)
	}

	var values string = splitted[0]
	var step int = 1
	var err error

	if len(values) == 0 {
		return nil, xyerror.ValueError.New("values is empty")
	}

	if len(splitted) == 2 {
		step, err = strconv.Atoi(splitted[1])
		if err != nil {
			return nil, xyerror.ValueError.Newf("can not convert %s to int: %s", splitted[1], err)
		}
	}

	var start, end int
	start, end, err = parseRange(values, max, atoi)
	if err != nil {
		return nil, err
	}

	n := math.Round(float64(end-start+1) / float64(step))
	result := make([]int, int(n))
	for i := range result {
		result[i] = start + i*step
	}

	return result, nil
}

type timePoint struct {
	sec, min, hour, day, year int
	mon                       time.Month
	weekday                   time.Weekday
}

func newTimePoint(t time.Time) timePoint {
	tp := timePoint{}
	tp.hour, tp.min, tp.sec = t.Clock()
	tp.year, tp.mon, tp.day = t.Date()
	tp.weekday = t.Weekday()

	return tp
}

func (tp timePoint) toTime(loc *time.Location) time.Time {
	return time.Date(
		tp.year, tp.mon, tp.day,
		tp.hour, tp.min, tp.sec, 0,
		loc)
}

func (tp timePoint) isValidDay() bool {
	y, m, d := tp.toTime(time.UTC).Date()
	return d == tp.day && m == tp.mon && y == tp.year
}

// unitRange is the struct representing for a range of time unit.
type unitRange struct {
	name string
	r    []int
	max  int
	atoi func(string) (int, error)
}

func newUnitRange(name string, max int, atoi func(string) (int, error)) unitRange {
	ur := unitRange{name: name, max: max, atoi: atoi}
	ur.set("*")
	return ur
}

func (ur *unitRange) set(s string) error {
	// Only one parsing function succeeds, the other returns error.
	r1, err1 := parseListToken(s, ur.atoi)
	r2, err2 := parseRangeToken(s, ur.max, ur.atoi)

	ur.r = append(r1, r2...)

	if len(ur.r) == 0 {
		return xyerror.ValueError.Newf("%s, %s", err1, err2)
	}

	for _, v := range ur.r {
		if v > ur.max {
			return xyerror.ValueError.Newf("value %d is larger than %d", v, ur.max)
		}
	}

	sort.Ints(ur.r)

	return nil
}

// findNext finds the element which is greater than or equal to u. If u is
// greater than the greatest element, it returns the smallest one.
func (ur unitRange) findNext(u int) int {
	var l, r int = 0, len(ur.r) - 1
	for l <= r {
		m := (l + r) / 2

		if ur.r[m] < u {
			l = m + 1
		}

		if l <= m {
			r = m - 1
		}
	}

	if r+1 >= len(ur.r) {
		return ur.r[0]
	}

	return ur.r[r+1]
}

func (ur unitRange) contains(u int) bool {
	return ur.findNext(u) == u
}
