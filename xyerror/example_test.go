package xyerror_test

import (
	"errors"
	"fmt"

	"github.com/xybor/xyplatform/xyerror"
)

var exampleGen = xyerror.Register("example", 400000)

func ExampleClass() {
	// To create a root Class, call Generator.NewClass with the name of Class.
	var RootError = exampleGen.NewClass("RootError")

	// You can create a class from another one.
	var ChildError = RootError.NewClass("ChildError")

	fmt.Println(RootError)
	fmt.Println(ChildError)

	// Output:
	// [400001] RootError
	// [400002] ChildError
}

func ExampleXyError() {
	// You can compare a XyError with an Class by using the built-in method
	// errors.Is.
	var NegativeIndexError = xyerror.IndexError.NewClass("NegativeIndexError")

	var err1 = xyerror.ValueError.New("some value error")
	if errors.Is(err1, xyerror.ValueError) {
		fmt.Println("err1 is a ValueError")
	}
	if !errors.Is(err1, NegativeIndexError) {
		fmt.Println("err1 is not a NegativeIndexError")
	}

	var err2 = NegativeIndexError.Newf("some negative index error %d", -1)
	if errors.Is(err2, NegativeIndexError) {
		fmt.Println("err2 is a NegativeIndexError")
	}
	if errors.Is(err2, xyerror.IndexError) {
		fmt.Println("err2 is a IndexError")
	}
	if !errors.Is(err2, xyerror.ValueError) {
		fmt.Println("err2 is not a ValueError")
	}

	// Output:
	// err1 is a ValueError
	// err1 is not a NegativeIndexError
	// err2 is a NegativeIndexError
	// err2 is a IndexError
	// err2 is not a ValueError
}

func ExampleGroup() {
	// Group allows you to create a class with multiparents.
	var KeyValueError = xyerror.
		Combine(xyerror.KeyError, xyerror.ValueError).
		NewClass(exampleGen, "KeyValueError")

	var err = KeyValueError.New("something is wrong")

	if errors.Is(err, xyerror.KeyError) {
		fmt.Println("err is a KeyError")
	}

	if errors.Is(err, xyerror.ValueError) {
		fmt.Println("err is a ValueError")
	}

	// Output:
	// err is a KeyError
	// err is a ValueError
}
