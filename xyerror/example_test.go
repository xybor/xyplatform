package xyerror_test

import (
	"errors"
	"fmt"

	"github.com/xybor/xyplatform"
	"github.com/xybor/xyplatform/xyerror"
)

// Create example module and register it with xyerror.
var ExampleModule = xyplatform.NewModule(40000, "ExampleModule")
var _ = xyerror.Register(ExampleModule)

func ExampleClass() {
	// To create a root error class, call xyerror.NewClass with that module.
	var RootError = xyerror.NewClass(ExampleModule, "RootError")

	// You can create an error class from another error class.
	var ChildError = RootError.NewClass("ChildError")

	fmt.Println(RootError)
	fmt.Println(ChildError)

	// Output:
	// [40001] RootError
	// [40002] ChildError
}

// You can compare a xyerror with an error class by using built-in method
// errors.Is(...).
func ExampleXyError() {
	var NegativeIndexError = xyerror.IndexError.NewClass("NegativeIndexError")

	err1 := xyerror.ValueError.New("some value error")
	if errors.Is(err1, xyerror.ValueError) {
		fmt.Println("err1 is a ValueError")
	}

	if !errors.Is(err1, NegativeIndexError) {
		fmt.Println("err1 is not a NegativeIndexError")
	}

	err2 := NegativeIndexError.New("some negative index error")
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

// Group allows you to create an error class with multiparents.
func ExampleGroup() {
	CombinedErrorClass := xyerror.
		Combine(xyerror.KeyError, xyerror.ValueError).
		NewClass(xyplatform.Default, "CombinedErrorClass")

	err := CombinedErrorClass.New("something is wrong")

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
