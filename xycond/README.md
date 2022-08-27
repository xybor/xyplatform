# Introduction

Package xycond supports to assert or expect many conditions.

It makes source code to be shorter and more readable by using inline commands.

# Features

This package has the following features:

-   Assert a condition, panic in case condition is false.
-   Expect a condition to occur and perform actions on this expectation.

Visit [pkg.go.dev](https://pkg.go.dev/github.com/xybor/xyplatform/xycond) for
more details.

# Example

```golang
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
```
