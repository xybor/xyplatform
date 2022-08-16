package xylock_test

import (
	"fmt"
	"time"

	"github.com/xybor/xyplatform/xylock"
)

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
