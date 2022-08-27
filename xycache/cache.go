// Package xycache provides cache management.
package xycache

import (
	"time"
	"unsafe"

	"github.com/xybor/xyplatform/xylock"
	"github.com/xybor/xyplatform/xysched"
)

const (
	B = 1 << (10 * iota)
	KB
	MB
	GB
)

// Cache is a key-value, in-memory, and limited-size storage. It supports to
// remove entries which are oldest or expired from storage.
type Cache[kt comparable, vt any] struct {
	lock        xylock.RWLock
	sched       *xysched.Scheduler
	lifetime    time.Duration
	size        int
	cache       map[any]*entry[kt, vt]
	oldestEntry *entry[kt, vt]
	newestEntry *entry[kt, vt]
}

// WithSize sets the maximum number of entries in Cache, no-limited size by
// default. When a new entry is added but no available slot, the oldest entry
// will be deleted.
func (c *Cache[kt, vt]) WithSize(s int) {
	c.lock.WLockFunc(func() {
		c.size = s
	})
}

// WithSizeInBytes estimates the maximum number of entries in Cache based on
// the maximum size.
func (c *Cache[kt, vt]) WithSizeInBytes(s int) {
	var e entry[kt, vt]
	var k kt
	var entrySize = int(unsafe.Sizeof(e)) + int(unsafe.Sizeof(k))
	c.WithSize(s / entrySize)
}

// WithExpiration sets the lifetime of entries. Expired entries will be deleted
// from Cache.
func (c *Cache[kt, vt]) WithExpiration(d time.Duration) {
	c.lock.WLockFunc(func() {
		c.lifetime = d
	})
}

// Sets adds a new entry to Cache. If the entry's key has existed, replace it
// with the new value.
func (c *Cache[kt, vt]) Set(key kt, value vt) {
	c.lock.WLockFunc(func() {
		c.set(key, value)
	})
}

// Get returns the value corresponding to the key in Cache. If the key does not
// exist, ok will be false.
func (c *Cache[kt, vt]) Get(key kt) (value vt, ok bool) {
	c.lock.RLockFunc(func() any {
		value, ok = c.get(key)
		return nil
	})
	return
}

// Delete removes the entry corresponding to the key from Cache.
func (c *Cache[kt, vt]) Delete(key kt) {
	c.lock.WLockFunc(func() {
		c.delete(key)
	})
}

// Replace is similar to Set in case of existed key. This method allows you to
// modify the value of entry with a function.
func (c *Cache[kt, vt]) Replace(key kt, replace func(vt) vt) {
	c.lock.WLockFunc(func() {
		var value, ok = c.get(key)
		if ok {
			value = replace(value)
			c.set(key, value)
		}
	})
}

// entry is a dependent element in cache, it contains key, value, and pointer to
// previous and next elements. Pointers supports to find the first and last
// element in Cache.
type entry[kt comparable, vt any] struct {
	key   kt
	value vt
	prev  *entry[kt, vt]
	next  *entry[kt, vt]
}

// set is the underlying method of adding or replacing an entry. This entry will
// be always the newest entry in Cache after set. If the Cache was full, the
// oldest one will be deleted.
func (c *Cache[kt, vt]) set(key kt, value vt) {
	// Create the cache in first use.
	if c.cache == nil {
		c.cache = make(map[any]*entry[kt, vt])
		c.sched = xysched.NewScheduler()
	}

	// Replace the current entry in case of existed key.
	if _, ok := c.cache[key]; ok {
		c.cache[key].value = value
		c.renew(key)
		return
	}

	// Delete the oldest entry if Cache was full.
	if c.size > 0 && len(c.cache) >= c.size {
		c.delete(c.oldestEntry.key)
	}

	// The current entry will be the newest one in Cache.
	var e = &entry[kt, vt]{
		key:   key,
		value: value,
		prev:  c.newestEntry,
		next:  nil,
	}
	c.cache[key] = e

	if c.oldestEntry == nil {
		c.oldestEntry = e
	}

	if c.newestEntry == nil {
		c.newestEntry = e
	} else {
		c.newestEntry.next = e
		c.newestEntry = e
	}

	// Delete entry from Cache when it expired.
	if c.lifetime > 0 {
		c.sched.After(c.lifetime) <- xysched.NewTask(func() {
			c.lock.WLockFunc(func() { c.delete(key) })
		})
	}
}

// get is the underlying method of get the value corresponding to the key. If
// the key doesn't exist, the latter output will be false.
func (c *Cache[kt, vt]) get(key kt) (vt, bool) {
	var value vt
	var ok bool
	var e *entry[kt, vt]
	e, ok = c.cache[key]
	if ok {
		value = e.value
		c.renew(key)
	}
	return value, ok
}

// delete removes an entry from Cache.
func (c *Cache[kt, vt]) delete(key kt) {
	var e, ok = c.cache[key]
	if !ok {
		return
	}

	if e == c.oldestEntry {
		c.oldestEntry = e.next
	}
	if e == c.newestEntry {
		c.newestEntry = e.prev
	}

	if e.prev != nil {
		e.prev.next = e.next
	}

	if e.next != nil {
		e.next.prev = e.prev
	}

	delete(c.cache, key)
}

// renew transforms an entry to the newest one.
func (c *Cache[kt, vt]) renew(key kt) {
	var e, ok = c.cache[key]
	if !ok {
		return
	}

	if e == c.oldestEntry {
		c.oldestEntry = e.next
	}

	if e.prev != nil {
		e.prev.next = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	}

	e.prev = c.newestEntry
	e.next = nil
}
