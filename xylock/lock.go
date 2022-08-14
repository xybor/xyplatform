// This module defines wrapper types of sync mutex, rwmutex, and semaphore.
package xylock

import (
	"sync"
)

// Lock is a wrapper struct of sync.Mutex. All methods of Lock support to do
// nothing if Lock pointer is nil.
type Lock struct {
	m sync.Mutex
}

// Lock locks l.
// If the lock is already in use, the calling goroutine blocks until the mutex
// is available.
func (l *Lock) Lock() {
	if l != nil {
		l.m.Lock()
	}
}

// TryLock tries to lock l and reports whether it succeeded.
//
// Note that while correct uses of TryLock do exist, they are rare, and use of
// TryLock is often a sign of a deeper problem in a particular use of mutexes.
func (l *Lock) TryLock() bool {
	if l != nil {
		return l.m.TryLock()
	}
	return true
}

// Unlock unlocks l.
// It is a run-time error if l is not locked on entry to Unlock.
//
// A locked Mutex is not associated with a particular goroutine. It is allowed
// for one goroutine to lock a Mutex and then arrange for another goroutine to
// unlock it.
func (l *Lock) Unlock() {
	if l != nil {
		l.m.Unlock()
	}
}

// LockFunc blocks to call a function until it completes.
func (l *Lock) LockFunc(f func()) {
	l.Lock()
	defer l.Unlock()
	f()
}

// RLockFunc blocks to call a function until it completes, then returns the
// reading value.
func (l *Lock) RLockFunc(f func() any) any {
	l.Lock()
	defer l.Unlock()
	return f()
}
