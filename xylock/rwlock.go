// This file copies comments from sync/rwmutex.go (go1.18.4).
package xylock

import "sync"

// RWLock is a wrapper struct of sync.RWMutex. All methods of RWLock support to
// do nothing if RWLock pointer is nil.
type RWLock struct {
	m sync.RWMutex
}

// Lock locks rw for writing. If the lock is already locked for reading or
// writing, Lock blocks until the lock is available.
func (rw *RWLock) Lock() {
	if rw != nil {
		rw.m.Lock()
	}
}

// RLock locks rw for reading.
//
// It should not be used for recursive read locking; a blocked Lock call
// excludes new readers from acquiring the lock. See the documentation on the
// sync.RWMutex type.
func (rw *RWLock) RLock() {
	if rw != nil {
		rw.m.RLock()
	}
}

// RLocker returns a Locker interface that implements the Lock and Unlock
// methods by calling rw.RLock and rw.RUnlock.
func (rw *RWLock) RLocker() sync.Locker {
	if rw != nil {
		return rw.m.RLocker()
	}
	return nil
}

// RUnlock undoes a single RLock call; it does not affect other simultaneous
// readers. It is a run-time error if rw is not locked for reading on entry to
// RUnlock.
func (rw *RWLock) RUnlock() {
	if rw != nil {
		rw.m.RUnlock()
	}
}

// TryLock tries to lock rw for writing and reports whether it succeeded.
//
// Note that while correct uses of TryLock do exist, they are rare, and use of
// TryLock is often a sign of a deeper problem in a particular use of mutexes.
func (rw *RWLock) TryLock() bool {
	if rw != nil {
		return rw.m.TryLock()
	}
	return true
}

// TryRLock tries to lock rw for reading and reports whether it succeeded.
//
// Note that while correct uses of TryRLock do exist, they are rare, and use of
// TryRLock is often a sign of a deeper problem in a particular use of mutexes.
func (rw *RWLock) TryRLock() bool {
	if rw != nil {
		return rw.m.TryRLock()
	}
	return true
}

// Unlock unlocks rw for writing. It is a run-time error if rw is not locked for
// writing on entry to Unlock.
//
// As with Mutexes, a locked RWMutex is not associated with a particular
// goroutine. One goroutine may RLock (Lock) a RWMutex and then arrange for
// another goroutine to RUnlock (Unlock) it.
func (rw *RWLock) Unlock() {
	if rw != nil {
		rw.m.Unlock()
	}
}

// WLockFunc blocks to call a function for writing until it completes.
func (rw *RWLock) WLockFunc(f func()) {
	rw.Lock()
	defer rw.Unlock()
	f()
}

// RLockFunc blocks to call a function for reading until it completes, then
// returns reading value.
func (rw *RWLock) RLockFunc(f func() any) any {
	rw.RLock()
	defer rw.RUnlock()
	return f()
}

// RWLockFunc blocks to call a function for reading and writing until it
// completes, then returns reading value.
func (rw *RWLock) RWLockFunc(f func() any) any {
	rw.Lock()
	defer rw.Unlock()
	return f()
}
