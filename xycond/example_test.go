package xycond_test

import (
	"fmt"
	"testing"

	"github.com/xybor/xyplatform/xycond"
)

func Example() {
	xycond.AssertFalse(1 == 2)

	var x int
	xycond.AssertZero(x)

	// Test a condition with *testing.T or *testing.B.
	var t = &testing.T{}
	xycond.ExpectEmpty("").Test(t)

	// Perform actions on an expectation.
	xycond.ExpectEqual(1, 2).
		True(func() {
			fmt.Printf("1 == 2")
		}).
		False(func() {
			fmt.Printf("1 != 2")
		})

	// Output:
	// 1 != 2
}
