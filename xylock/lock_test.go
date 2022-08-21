package xylock_test

import (
	"testing"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xylock"
)

func TestLock(t *testing.T) {
	var l = xylock.Lock{}
	xycond.MustNotPanic(func() {
		l.LockFunc(func() {})
		l.RLockFunc(func() any { return nil })
		l.TryLock()
	}).Test(t, "A panic occurred")
}

func TestLockNil(t *testing.T) {
	var l *xylock.Lock = nil
	xycond.MustNotPanic(func() {
		l.LockFunc(func() {})
		l.RLockFunc(func() any { return nil })
		l.TryLock()
	}).Test(t, "A panic occurred")
}

func TestRWLock(t *testing.T) {
	var l = xylock.RWLock{}
	xycond.MustNotPanic(func() {
		l.WLockFunc(func() {})
		l.RLockFunc(func() any { return nil })
		l.RWLockFunc(func() any { return nil })
		l.TryLock()
		l.TryRLock()
		l.RLocker()
	}).Test(t, "A panic occurred")
}

func TestRWLockNil(t *testing.T) {
	var l *xylock.RWLock = nil
	xycond.MustNotPanic(func() {
		l.WLockFunc(func() {})
		l.RLockFunc(func() any { return nil })
		l.RWLockFunc(func() any { return nil })
		l.TryLock()
		l.TryRLock()
		l.RLocker()
	}).Test(t, "A panic occurred")
}

func TestSemaphore(t *testing.T) {
	var s = xylock.NewSemaphore(1)
	xycond.MustNotPanic(func() {
		s.AcquireFunc(1, func() {})
		s.RAcquireFunc(1, func() any { return nil })
		s.TryAcquire(1)
	}).Test(t, "A panic occurred")
}

func TestSemaphoreNil(t *testing.T) {
	var s *xylock.Semaphore = nil
	xycond.MustNotPanic(func() {
		s.AcquireFunc(1, func() {})
		s.RAcquireFunc(1, func() any { return nil })
		s.TryAcquire(1)
	}).Test(t, "A panic occurred")
}
