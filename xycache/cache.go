// Package xycache provides cache management.
package xycache

import (
	"time"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xylock"
	"github.com/xybor/xyplatform/xysched"
)

// entry is a dependent element in cache, it contains key, value, expired time.
type entry[kt comparable, vt any] struct {
	value     vt
	exipireAt time.Time
	unode     *node[kt]
	enode     *node[kt]
}

// Cache is a key-value, in-memory, expirable, and limited-size storage.
//
// Cache stores entries in three data structures:
//  - map: supports to quickly access to an entry via its key.
//  - ulist (or use list): supports to delete the oldest entry when the cache is
//    full. This list stores entries in ascending order of last use.
//  - elist (or expiration list): supports to collect expired entries without
//    iterating over total cache. This list stores entries according to
//    ascending order of expired time.
type Cache[kt comparable, vt any] struct {
	lock        xylock.RWLock
	expiredTime time.Duration
	colector    *xysched.Cron
	size        int
	cache       map[any]*entry[kt, vt]
	ulist       list[kt]
	elist       list[kt]
}

// LimitSize sets the maximum number of entries in Cache, no-limited size by
// default. When a new entry is added but no available slot, the oldest entry
// will be deleted.
func (c *Cache[kt, vt]) LimitSize(s int) {
	c.lock.WLockFunc(func() {
		c.size = s
	})
}

// Expiration sets the lifetime of entries. Expired entries will be deleted
// from Cache.
func (c *Cache[kt, vt]) Expiration(d time.Duration) {
	c.lock.WLockFunc(func() {
		xycond.AssertZero(int64(c.expiredTime))
		c.expiredTime = d
	})
	c.SetCollectorInterval(d)
}

// SetCollectorInterval sets the interval time between cleanings of expired
// entries. This method is only available if expiration was set.
func (c *Cache[kt, vt]) SetCollectorInterval(interval time.Duration) {
	c.lock.WLockFunc(func() {
		xycond.AssertNotZero(int64(c.expiredTime))

		if c.colector != nil {
			c.colector.Stop()
		}

		c.colector = xysched.NewCron(func() {
			c.lock.WLockFunc(func() {
				for enode := c.elist.first; enode != nil; enode = enode.next {
					var exipireAt = c.cache[enode.key].exipireAt
					if exipireAt.After(time.Now()) {
						// If the current entry in elist has not expired yet,
						// neither will the next item.
						break
					}
					c.delete(enode.key)
				}
			})
		})
		c.colector.Every(interval)
		xysched.Now() <- c.colector
	})
}

// Sets adds a new entry to Cache. If the entry's key has existed, replace it
// with the new value. If the Cache is full, delete the oldest one.
func (c *Cache[kt, vt]) Set(key kt, value vt) {
	c.lock.WLockFunc(func() {
		c.set(key, value)
	})
}

// Get returns the value corresponding to the key in Cache. If the key does not
// exist Cache or it expired, ok will be false.
func (c *Cache[kt, vt]) Get(key kt) (value vt, ok bool) {
	c.lock.RLockFunc(func() any {
		value, ok = c.get(key)
		return nil
	})
	return
}

// MustGet returns the value corresponding to the key in Cache. It panics if the
// key doesn't exist in Cache or entry expired.
func (c *Cache[kt, vt]) MustGet(key kt) vt {
	var value, ok = c.Get(key)
	xycond.AssertTrue(ok)
	return value
}

// Delete removes the entry corresponding to the key from Cache.
func (c *Cache[kt, vt]) Delete(key kt) {
	c.lock.WLockFunc(func() {
		c.delete(key)
	})
}

// Replace is similar to Set in case of existed key. This method allows you to
// modify the value of entry within a function.
func (c *Cache[kt, vt]) Replace(key kt, replace func(vt) vt) {
	c.lock.WLockFunc(func() {
		var value, ok = c.get(key)
		if ok {
			value = replace(value)
			c.set(key, value)
		}
	})
}

// set adds or replaces value of an entry. This entry will be the newest entry
// in Cache. If the Cache was full, the oldest one will be deleted.
func (c *Cache[kt, vt]) set(key kt, value vt) {
	// Create the cache in the first use.
	if c.cache == nil {
		c.cache = make(map[any]*entry[kt, vt])
	}

	// Replace the current entry in case of existed key.
	if _, ok := c.cache[key]; ok {
		c.cache[key].value = value
		c.renew(key)
		return
	}

	// Delete the oldest entry if Cache was full.
	if c.size > 0 && len(c.cache) >= c.size {
		c.delete(c.ulist.first.key)
	}

	// The current entry will be the newest one in Cache.
	var e = &entry[kt, vt]{
		value: value,
	}
	e.unode = &node[kt]{key: key}
	if c.expiredTime > 0 {
		e.exipireAt = time.Now().Add(c.expiredTime)
		e.enode = &node[kt]{key: key}
	}

	c.cache[key] = e
	c.ulist.append(e.unode)
	c.elist.append(e.enode)
}

// get is the underlying method of get the value corresponding to the key. If
// the key doesn't exist or it expired, the latter output will be false,
// otherwise, renew the entry.
func (c *Cache[kt, vt]) get(key kt) (vt, bool) {
	var value vt
	var e, ok = c.cache[key]
	if ok {
		value = e.value
		if c.expiredTime > 0 && e.exipireAt.Before(time.Now()) {
			ok = false
		} else {
			c.renew(key)
		}
	}

	return value, ok
}

// delete removes an entry from Cache.
func (c *Cache[kt, vt]) delete(key kt) {
	var e, ok = c.cache[key]
	if !ok {
		return
	}

	c.ulist.remove(e.unode)
	c.elist.remove(e.enode)
	delete(c.cache, key)
}

// renew transforms an entry to the newest one.
func (c *Cache[kt, vt]) renew(key kt) {
	var e, ok = c.cache[key]
	if !ok {
		return
	}

	c.ulist.remove(e.unode)
	c.ulist.append(e.unode)
}
