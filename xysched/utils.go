package xysched

import (
	"reflect"

	"github.com/xybor/xyplatform/xycond"
)

// callFunc calls a function and returns its returned values. This function
// return a non-nil error if parameter is mismatched or function panics.
//
// if variadic is true, this function will convert the returned values from
// []any to ...any. It is helpful to pass them as parameters of other functions.
func callFunc(fv reflect.Value, p []any, variadic bool) []any {
	var ftype = fv.Type()
	var ninput = ftype.NumIn()
	if ftype.IsVariadic() {
		xycond.Condition(len(p) >= ninput-1).Assert(
			"expected at least %d, but got %d parameters", ninput-1, len(p))
	} else {
		xycond.Condition(len(p) == ninput).Assert(
			"expected %d, but got %d parameters", ninput, len(p))
	}

	in := make([]reflect.Value, len(p))
	for k, param := range p {
		in[k] = reflect.ValueOf(param)
	}

	var result = fv.Call(in)
	var iresult = make([]any, 0)
	for i := range result {
		iresult = append(iresult, result[i].Interface())
	}

	if variadic {
		iresult = iresult[0].([]any)
	}
	return iresult
}
