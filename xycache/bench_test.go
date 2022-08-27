package xycache_test

import (
	"testing"
	"time"

	"github.com/xybor/xyplatform/xycache"
)

var cache xycache.Cache[int, string]

func init() {
	cache = xycache.Cache[int, string]{}
	cache.WithSize(1000000)
	cache.WithExpiration(10 * time.Second)
}

func BenchmarkCache1Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cache.Set(i, "many thing to show")
	}
}

func BenchmarkCache2Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cache.Get(i % 1000000)
	}
}

func BenchmarkCache3Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cache.Set(i%1000000, "a new value to set")
	}
}
