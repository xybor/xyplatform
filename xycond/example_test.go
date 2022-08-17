package xycond_test

import "github.com/xybor/xyplatform/xycond"

func ExampleCondition() {
	// Assert 1 == 2
	xycond.MustFalse(1 == 2).Assert("weird")

	// Assert x is 0
	var x int
	xycond.MustZero(x).Assert("%d is not initialized with zero", x)

	// Assert string is empty, panic without any message if the condition fails.
	xycond.MustEmpty("").JustAssert()
}
