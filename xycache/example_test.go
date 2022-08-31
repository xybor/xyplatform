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
	fmt.Println(cache.MustGet(1))

	// Output:
	// foo
}

func ExampleCache() {
	var cache = xycache.Cache[int, string]{}
	cache.LimitSize(2)
	cache.Expiration(time.Millisecond)
	cache.Set(1, "foo")
	cache.Set(2, "bar")
	cache.Set(3, "buzz")
	var value, ok = cache.Get(1)
	if !ok {
		fmt.Println("key 1 has removed because the cache is full")
	}
	value, ok = cache.Get(2)
	if ok {
		fmt.Println("key 2 is", value)
	}
	time.Sleep(10 * time.Millisecond)
	value, ok = cache.Get(2)
	if !ok {
		fmt.Println("key 2 has removed because it expired")
	}

	// Output:
	// key 1 has removed because the cache is full
	// key 2 is bar
	// key 2 has removed because it expired
}
