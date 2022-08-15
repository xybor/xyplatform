package xycond_test

import "github.com/xybor/xyplatform/xycond"

func ExampleCondition() {
	// Assert 1 == 2
	xycond.False(1 == 2).Assert("Weird")

	// Assert x is 0
	var x int
	xycond.Zero(x).Assert("%d is not initialized with zero", x)

	// Assert string is empty, panic without any message if the condition fails.
	xycond.Empty("").JustAssert()
}
