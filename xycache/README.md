# Introduction

Package xycache provides cache management.

# Features

Visit [pkg.go.dev](https://pkg.go.dev/github.com/xybor/xyplatform/xycache) for
more public APIs.

Cache is a key-value, in-memory, expirable, and limited-size storage.

After a expirable Cache starts, it schedules a goroutine which removes expired
entries from Cache every time interval.

If you adds an entry to a limited-size storage when it is full, the oldest entry
will be deleted to make slot for the new one.

# Example

1.  Default cache

```golang
var cache = xycache.Cache[string, int]{}
cache.Set("foo", 1)
cache.Set("bar", 2)

cache.Get("foo")        // 1, true
cache.MustGet("bar")    // 2
cache.MustGet("buzzz")  // panic
```

2.  Limited-size cache

```golang
var cache = xycache.Cache[string, int]{}
cache.WithSize(2)       // The cache has the maximum of two entries.
cache.Set("foo", 1)
cache.Set("bar", 2)

cache.Get("bar")        // bar is renewed (foo now is the oldest entry)
cache.Get("foo")        // foo is renewed (bar now is the oldest entry)

cache.Set("buzz", 3)    // Delete the oldest entry before adding the new one.
cache.MustGet("bar")    // panic (bar was deleted)
cache.MustGet("foo")    // 1
cache.MustGet("buzz")   // 3
```

3.  Expirable cache

```golang
var cache = xycache.Cache[string, int]{}
// Every entry will exist for a second before being deleted.
cache.WithExpiration(time.Second)
// Clean the expired entries every milisecond.
cache.WithColectorInterval(time.Millisecond)

cache.Set("foo", 1)
cache.Set("bar", 2)

cache.MustGet("foo")      // 1
time.Sleep(time.Second)   // foo, bar are deleted because they expired.
cache.Set("buzz", 3)

cache.Get("foo")          // 0, false
cache.Get("bar")          // 0, false
cache.Get("buzz")         // 3, true
```

4.  Replace entry

```golang
var cache = xycache.Cache[string, int]{}
cache.WithExpiration(time.Second)
cache.Set("foo", 1)
cache.MustGet("foo")                    // 1, true

cache.Set("foo", 2)
cache.Get("foo")                        // 2, true

cache.Replace("foo", func(v int) int {
    return v*5
})
cache.Get("foo")                        // 10, true
```
