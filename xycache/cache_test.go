package xycache_test

import (
	"testing"

	"github.com/xybor/xyplatform/xycache"
	"github.com/xybor/xyplatform/xycond"
)

func TestCacheSet(t *testing.T) {
	var cache = xycache.Cache[int, int]{}
	cache.Set(1, 1)
	var value, ok = cache.Get(1)
	xycond.ExpectTrue(ok).Test(t)
	xycond.ExpectEqual(value, 1).Test(t)
}

func TestCacheSetDuplicate(t *testing.T) {
	var cache = xycache.Cache[int, string]{}
	cache.Set(1, "foo")
	cache.Set(1, "bar")
	var value, ok = cache.Get(1)
	xycond.ExpectTrue(ok).Test(t)
	xycond.ExpectEqual(value, "bar").Test(t)
}

func TestCacheReplace(t *testing.T) {
	var cache = xycache.Cache[string, int]{}
	cache.Set("foo", 1)
	cache.Replace("foo", func(i int) int { return i + 2 })
	var value, ok = cache.Get("foo")
	xycond.ExpectTrue(ok).Test(t)
	xycond.ExpectEqual(value, 3).Test(t)
}

func TestCacheDelete(t *testing.T) {
	var cache = xycache.Cache[string, string]{}
	cache.Set("foo", "bar")
	cache.Delete("foo")
	var _, ok = cache.Get("foo")
	xycond.ExpectFalse(ok).Test(t)
}

func TestCacheRemoveOldest(t *testing.T) {
	var cache = xycache.Cache[string, any]{}
	cache.LimitSize(2)
	cache.Set("foo", nil)
	cache.Set("bar", nil)
	cache.Set("buzz", nil)
	var _, ok = cache.Get("foo")
	xycond.ExpectFalse(ok).Test(t)
}

func TestCacheRenew(t *testing.T) {
	var cache = xycache.Cache[int, any]{}
	cache.LimitSize(2)
	cache.Set(1, nil)
	cache.Set(2, nil)
	cache.Get(1)
	cache.Set(3, nil)
	var _, ok = cache.Get(1)
	xycond.ExpectTrue(ok).Test(t)
}
