# Introduction

Xylock contains wrapper structs of built-in `sync` library, such as `sync.Mutex`
or `semaphore.Weighted`.

# Features

Xylock wrapper structs have fully methods of origin structs. For example, `Lock`
is the wrapper struct of `sync.Mutex`, and it has the following methods:

```golang
func (l *Lock) Lock()
func (l *Lock) Unlock()
```

Methods of wrapper structs have an additional features, that is to do nothing if
the receiver pointer is nil. This is helpful when lock is an optional
development.

```golang
// These following commands will not cause a panic. They just do nothing.
var lock *xylock.Lock = nil
lock.Lock()
lock.Unlock()
```

Xylock structs allows to run a function in thread-safe area (with Lock and
Unlock cover the function).

```golang
var lock = xylock.Lock{}
lock.LockFunc(func() {
    // thread-safe area
})
```

Thread-safe methods of Xylog structs with `R` as prefix will support to read
data.

```golang
var foo int
var lock = xylock.Lock{}
var result = lock.RLockFunc(func() any {
    return foo
}).(int)
```

Visit [pkg.go.dev](https://pkg.go.dev/github.com/xybor/xyplatform/xylock) for
more details.

# Example

```golang
func Example() {
	var x int
	var lock = xylock.Lock{}
	for i := 0; i < 10000; i++ {
		go lock.LockFunc(func() {
			x = x + 1
		})
	}

	time.Sleep(time.Millisecond)
	fmt.Println(lock.RLockFunc(func() any { return x }))

	// Output:
	// 10000
}
```
