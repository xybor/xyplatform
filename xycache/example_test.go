package xycache_test

import (
	"fmt"
	"time"

	"github.com/xybor/xyplatform/xycache"
)

func Example() {
	var cache = xycache.Cache[int, string]{}
	cache.Set(1, "foo")
	cache.Set(2, "bar")
	var value, ok = cache.Get(1)
	if ok {
		fmt.Println(value)
	}

	// Output:
	// foo
}

func ExampleCache() {
	var cache = xycache.Cache[int, string]{}
	cache.WithSize(2)
	cache.WithExpiration(time.Millisecond)
	cache.Set(1, "foo")
	cache.Set(2, "bar")
	cache.Set(3, "buzz")
	var value, ok = cache.Get(1)
	if !ok {
		fmt.Println("key 1 has removed because cache is full")
	}
	value, ok = cache.Get(2)
	if ok {
		fmt.Println(value)
	}
	time.Sleep(10 * time.Millisecond)
	value, ok = cache.Get(2)
	if !ok {
		fmt.Println("key 2 has removed because it expired")
	}

	// Output:
	// key 1 has removed because cache is full
	// bar
	// key 2 has removed because it expired
}
