# Introduction
Xycond supports to check many types of condition and panic if the condition
fails.

It makes source code to be shorter and more readable by using inline assertion
commands.

# Features
This package has only one struct, `Condition`, a type alias of `bool`.

The package defines two methods for this struct:
```golang
// Panic without message if the Condition fails.
func (Condition) JustAssert()

// Panic with a formatted messsage if the Condition fails.
func (Condition) Assert(msg string, a ...any)
```

There are many functions to create `Condition` instances. Example:
```golang
// Check condition directly.
func True(bool) Condition
func False(bool) Condition

// Assert a number is zero.
func Zero(number) Condition

// Assert an object is nil.
func Nil(any) Condition

// Assert a slice contains the element.
func ContainA(a any, e any) Condition

// Assert a map contains the key.
func ContainM(m map, k any) Condition
```

Besides, Xycond also has two functions to panic without using a `Condition`,
they are `Panic` and `JustPanic`.

Visit [pkg.go.dev](https://pkg.go.dev/github.com/xybor/xyplatform/xycond) for
more details.

# Example
```golang
// Assert 1 == 2
xycond.False(1 == 2).Assert("Weird")

// Assert x is 0
var x int
xycond.Zero(x).Assert("%d is not initialized with zero", x)

// Assert string is empty, panic without any message if the condition fails.
xycond.Empty("").JustAssert()
```
