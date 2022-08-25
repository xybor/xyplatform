# Introduction

Xycond supports to check many types of condition and panic if the condition
fails.

It makes source code to be shorter and more readable by using inline assertion
commands.

# Features

This package has only one struct, `Condition`, a type alias of `bool`.

The package defines the following methods for this struct:

```golang
// Panic without message if the Condition fails.
func (Condition) JustAssert()

// Panic with a formatted messsage if the Condition fails.
func (Condition) Assert(msg string, a ...any)
```

```golang
// Test will call t.Error if condition is false.
func (c Condition) Test(t tester, args ...any)

// Testf will call t.Errorf if condition is false.
func (c Condition) Testf(t tester, format string, args ...any)
```

There are many functions to create `Condition` instances. Example:

```golang
// Check condition directly.
func MustTrue(bool) Condition
func MustFalse(bool) Condition

// Assert a number to be zero.
func MustZero(number) Condition

// Assert an object to be nil.
func MustNil(any) Condition

// Assert a slice or array to contain the element.
func MustContainA(a any, e any) Condition

// Assert a map to contain the key.
func MustContainM(m map, k any) Condition

// Assert a snippet of code to cause a panic.
func MustPanic(f func()) Condition
```

Besides, Xycond also has two functions to panic without using a `Condition`,
they are `Panic` and `JustPanic`.

Visit [pkg.go.dev](https://pkg.go.dev/github.com/xybor/xyplatform/xycond) for
more details.

# Example

```golang
// Assert 1 == 2
xycond.MustFalse(1 == 2).Assert("weird")

// Assert x is 0
var x int
xycond.MustZero(x).Assert("%d is not initialized with zero", x)

// Assert string is empty, panic without any message if the condition fails.
xycond.MustEmpty("").JustAssert()
```
