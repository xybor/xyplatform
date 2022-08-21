package xysched_test

import (
	"testing"
	"time"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xysched"
)

func TestNewScheduler(t *testing.T) {
	xycond.MustNotPanic(func() {
		var sched = xysched.NewScheduler()
		sched.Stop()
	}).Test(t, "A panic occurred")
}

func TestSchedulerAfter(t *testing.T) {
	var sched = xysched.NewScheduler()
	defer sched.Stop()

	xycond.MustNotPanic(func() {
		sched.After(time.Nanosecond) <- xysched.NewTask(func() {})
		time.Sleep(time.Millisecond)
	}).Test(t, "A panic occurs")
}

func TestSchedulerAfterButCloseSoon(t *testing.T) {
	var sched = xysched.NewScheduler()

	xycond.MustNotPanic(func() {
		sched.After(time.Second) <- xysched.NewTask(func() {})
		sched.Stop()
		time.Sleep(time.Millisecond)
	}).Test(t, "A panic occurs")
}

func TestSchedulerAfterWithNegativeDuration(t *testing.T) {
	var sched = xysched.NewScheduler()
	defer sched.Stop()

	xycond.MustNotPanic(func() {
		sched.After(-1) <- xysched.NewTask(func() {})
	}).Test(t, "A panic occurs")
}

func TestSchedulerAfterMacros(t *testing.T) {
	var sched = xysched.NewScheduler()
	defer sched.Stop()

	xycond.MustNotPanic(func() {
		sched.Now() <- xysched.NewTask(func() {})
	}).Test(t, "A panic occurs")

	xycond.MustNotPanic(func() {
		sched.At(time.Now()) <- xysched.NewTask(func() {})
	}).Test(t, "A panic occurs")
}

func TestSchedulerConcurrent(t *testing.T) {
	var sched = xysched.NewScheduler()
	defer sched.Stop()

	xycond.MustNotPanic(func() {
		sched.Singleton()
	}).Test(t, "A panic occurs")

	xycond.MustNotPanic(func() {
		sched.Concurrent(10)
	}).Test(t, "A panic occurs")

	xycond.MustNotPanic(func() {
		sched.Now() <- xysched.NewTask(func() {})
		time.Sleep(time.Millisecond)
	}).Test(t, "A panic occurs")
}
