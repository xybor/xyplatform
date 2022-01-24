package xyerror_test

import (
	"fmt"

	"github.com/xybor/xyplatform"
	"github.com/xybor/xyplatform/xyerror"
)

// Summary:
//        Error type is an error without using for returning. It
// contains a unique error number (errno). Using that errno, you
// can determine which error is returned by using IsA() method.
//
//        You also create an error chain by inheriting an error type
// from another. To check if an error is a child of another error,
// you use BelongsTo() method.
//
//        End-error (or error, simply) is an error using for returning.
// It is created from an error type by New() method.
//
//        See ExampleNewError for creating error types and error chain.
// See ExampleXyError for creating an end-error and check if an error
// is a or belongs to another error.

// To create a root error type, you must have a module first.
var ExampleModule = xyplatform.NewModule(40000, "ExampleModule")

// Next, register the created module with xyerror.
var _ = xyerror.Register(ExampleModule)

func ExampleNewType() {
	// To create a root error type, call xyerror.NewError with that module and
	// the name of error type.
	var RootError = xyerror.NewType(ExampleModule, "RootError")

	// To create a child error type from the parent one, you call
	// NewError with the parent as receiver and name as parameter.
	var ChildError = RootError.NewType("ChildError")

	fmt.Println(RootError)
	fmt.Println(ChildError)

	// Output:
	// [40001] RootError
	// [40002][RootError] ChildError
}

func ExampleXyError() {
	// Create an end-error by using New() method with an error type.
	// You should write `return xyerror.ValueError.New("something")`.
	// Declaring a variable for this error is only for demonstration.
	var LessThan10Error = xyerror.ValueError.New("input value must be less than 10")

	// Create an error chain and another end-error for demonstration.
	var NegativeIndexError = xyerror.IndexError.NewType("NegativeIndexError")
	var UsePositiveError = NegativeIndexError.New("you must use an positive index")

	fmt.Println(LessThan10Error)
	fmt.Println(UsePositiveError)

	if LessThan10Error.IsA(xyerror.ValueError) {
		fmt.Println("LessThan10Error is a ValueError")
	}

	if LessThan10Error.BelongsTo(xyerror.ValueError) {
		fmt.Println("LessThan10Error belongs to ValueError")
	}

	if !LessThan10Error.IsA(NegativeIndexError) {
		fmt.Println("LessThan10Error is not a NegativeIndexError")
	}

	if UsePositiveError.IsA(NegativeIndexError) {
		fmt.Println("UsePositiveError is a NegativeIndexError")
	}

	if UsePositiveError.BelongsTo(xyerror.IndexError) {
		fmt.Println("UsePositiveError belongs to IndexError")
	}

	if !UsePositiveError.IsA(xyerror.IndexError) {
		fmt.Println("UsePositiveError is not an IndexError")
	}

	// Output:
	// [10008][ValueError] input value must be less than 10
	// [10009][NegativeIndexError] you must use an positive index
	// LessThan10Error is a ValueError
	// LessThan10Error belongs to ValueError
	// LessThan10Error is not a NegativeIndexError
	// UsePositiveError is a NegativeIndexError
	// UsePositiveError belongs to IndexError
	// UsePositiveError is not an IndexError
}
