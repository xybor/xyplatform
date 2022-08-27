package xycache_test

import (
	"testing"

	"github.com/xybor/xyplatform/xycache"
	"github.com/xybor/xyplatform/xycond"
)

func TestCacheSet(t *testing.T) {
	var cache = xycache.Cache[int, int]{}
	xycond.MustNotPanic(func() {
		cache.Set(1, 1)
	})
	var value, ok = cache.Get(1)
	xycond.MustTrue(ok).Test(t, "Expected true, but got false")
	xycond.MustEqual(value, 1).Testf(t, "Expected 1, but got %d", value)
}

func TestCacheSetDuplicate(t *testing.T) {
	var cache = xycache.Cache[int, string]{}
	xycond.MustNotPanic(func() {
		cache.Set(1, "foo")
		cache.Set(1, "bar")
	})
	var value, ok = cache.Get(1)
	xycond.MustTrue(ok).Test(t, "Expected true, but got false")
	xycond.MustEqual(value, "bar").Testf(t, "Expecte bar, but got %s", value)
}

func TestCacheReplace(t *testing.T) {
	var cache = xycache.Cache[string, int]{}
	xycond.MustNotPanic(func() {
		cache.Set("foo", 1)
		cache.Replace("foo", func(i int) int { return i + 2 })
	})
	var value, ok = cache.Get("foo")
	xycond.MustTrue(ok).Test(t, "Expected true, but got false")
	xycond.MustEqual(value, 3).Testf(t, "Expected 3, but got %d", value)
}

func TestCacheDelete(t *testing.T) {
	var cache = xycache.Cache[string, string]{}
	xycond.MustNotPanic(func() {
		cache.Set("foo", "bar")
		cache.Delete("foo")
	})
	var _, ok = cache.Get("foo")
	xycond.MustFalse(ok).Test(t, "Expected false, but got true")
}

func TestCacheRemoveOldest(t *testing.T) {
	var cache = xycache.Cache[string, any]{}
	cache.WithSize(2)
	cache.Set("foo", nil)
	cache.Set("bar", nil)
	cache.Set("buzz", nil)
	var _, ok = cache.Get("foo")
	xycond.MustFalse(ok).Test(t, "Expected false, but got true")
}

func TestCacheRenew(t *testing.T) {
	var cache = xycache.Cache[int, any]{}
	cache.WithSize(2)
	cache.Set(1, nil)
	cache.Set(2, nil)
	cache.Get(1)
	cache.Set(3, nil)
	var _, ok = cache.Get(1)
	xycond.MustTrue(ok).Test(t, "Expected true, but got false")
}

func TestCacheWithSizeInBytes(t *testing.T) {
	var cache = xycache.Cache[int, string]{}
	cache.WithSizeInBytes(100 * xycache.B)
	cache.Set(1, "foo")
	cache.Set(2, "bar")
	cache.Set(3, "foobar")
	var _, ok = cache.Get(1)
	xycond.MustFalse(ok).Test(t, "Expected false, but got true")
}
