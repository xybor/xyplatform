package xysched

import (
	"reflect"

	"github.com/xybor/xyplatform/xycond"
)

// callFunc calls a function and returns its returned values.
func callFunc(fv reflect.Value, p []any) []any {
	var in = make([]reflect.Value, len(p))
	for k, param := range p {
		in[k] = reflect.ValueOf(param)
	}

	var result = fv.Call(in)
	var iresult = make([]any, 0)
	for i := range result {
		iresult = append(iresult, result[i].Interface())
	}

	return iresult
}

// checkParam panics if function input can not fit with params.
func checkParam(params []reflect.Type, in []reflect.Type, isVariadic bool) {
	var ninput = len(params)
	if isVariadic {
		ninput--
	}
	xycond.AssertNotLessThan(len(in), ninput)
	if !isVariadic {
		xycond.AssertNotGreaterThan(len(in), ninput)
	}

	for i := range in {
		xycond.AssertNotEqual(in[i].Kind(), reflect.Invalid)
	}

	for i := 0; i < ninput; i++ {
		xycond.AssertTrue(in[i].AssignableTo(params[i]))
	}
}

// anyArrayToTypeArray converts an array of any to array of reflect.Type.
func anyArrayToTypeArray(a []any) []reflect.Type {
	var in = make([]reflect.Type, len(a))
	for i := range a {
		in[i] = reflect.TypeOf(a[i])
	}
	return in
}

// getFuncIn returns all function's input parameters under an array of
// reflect.Type.
func getFuncIn(f reflect.Type) []reflect.Type {
	var in = make([]reflect.Type, f.NumIn())
	for i := range in {
		in[i] = f.In(i)
	}
	return in
}

// getFuncOut returns all function's output parameters under an array of
// reflect.Type.
func getFuncOut(f reflect.Type) []reflect.Type {
	var in = make([]reflect.Type, f.NumOut())
	for i := range in {
		in[i] = f.Out(i)
	}
	return in
}
