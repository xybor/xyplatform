package xylock

import (
	"context"

	"golang.org/x/sync/semaphore"
)

// Semaphore is a wrapper of semaphore.Weighted. All methods of Semaphore
// support to do nothing if Semaphore pointer is nil.
type Semaphore struct {
	w *semaphore.Weighted
}

// NewSemaphore creates a new semaphore with the given maximum combined weight
// for concurrent access. All method of Semaphore support to do nothing if
// Semaphore pointer is nil.
func NewSemaphore(n int64) *Semaphore {
	return &Semaphore{w: semaphore.NewWeighted(n)}
}

// AcquireCtx acquires the semaphore with a weight of n, blocking until
// resources are available or ctx is done. On success, returns nil. On failure,
// returns ctx.Err() and leaves the semaphore unchanged.
//
// If ctx is already done, AcquireCtx may still succeed without blocking.
func (s *Semaphore) AcquireCtx(ctx context.Context, n int64) {
	if s != nil {
		s.w.Acquire(ctx, n)
	}
}

// Acquire is a shortcut of AcquireCtx with context.TODO().
func (s *Semaphore) Acquire(n int64) {
	if s != nil {
		s.AcquireCtx(context.TODO(), n)
	}
}

// TryAcquire acquires the semaphore with a weight of n without blocking.
// On success, returns true. On failure, returns false and leaves the semaphore
// unchanged.
func (s *Semaphore) TryAcquire(n int64) bool {
	if s != nil {
		return s.w.TryAcquire(n)
	}
	return true
}

// Release releases the semaphore with a weight of n.
func (s *Semaphore) Release(n int64) {
	if s != nil {
		s.w.Release(n)
	}
}

// AcquireFunc acquires the semaphore with a weight of n and blocks to call a
// function until it completes, then releases with the same weight.
func (s *Semaphore) AcquireFunc(n int64, f func()) {
	s.Acquire(n)
	defer s.Release(n)
	f()
}

// RAcquireFunc acquires the semaphore with a weight of n and blocks to call a
// function for reading until it completes, then releases with the same weight
// and returns the reading value.
func (s *Semaphore) RAcquireFunc(n int64, f func() any) any {
	s.Acquire(n)
	defer s.Release(n)
	return f()
}
