package xylock_test

import (
	"testing"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xylock"
)

func TestLock(t *testing.T) {
	var tests = []*xylock.Lock{{}, nil}

	for i := range tests {
		xycond.MustNotPanic(func() {
			tests[i].LockFunc(func() {})
			tests[i].RLockFunc(func() any { return nil })
			tests[i].TryLock()
		}).Test(t, "A panic occurred")
	}
}

func TestRWLock(t *testing.T) {
	var tests = []*xylock.RWLock{{}, nil}
	for i := range tests {
		xycond.MustNotPanic(func() {
			tests[i].WLockFunc(func() {})
			tests[i].RLockFunc(func() any { return nil })
			tests[i].RWLockFunc(func() any { return nil })
			tests[i].TryLock()
			tests[i].TryRLock()
			tests[i].RLocker()
		}).Test(t, "A panic occurred")
	}
}

func TestSemaphore(t *testing.T) {
	var tests = []*xylock.Semaphore{xylock.NewSemaphore(1), nil}

	for i := range tests {
		xycond.MustNotPanic(func() {
			tests[i].AcquireFunc(1, func() {})
			tests[i].RAcquireFunc(1, func() any { return nil })
			tests[i].TryAcquire(1)
		}).Test(t, "A panic occurred")
	}
}
