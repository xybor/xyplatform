package xysched_test

import (
	"testing"
	"time"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xysched"
)

func TestNewScheduler(t *testing.T) {
	var sched = xysched.NewScheduler(t.Name())
	xycond.ExpectEqual(sched, xysched.GetScheduler(t.Name())).Test(t)
	defer sched.Stop()
}

func TestSchedulerAfter(t *testing.T) {
	var sched = xysched.NewScheduler("")
	defer sched.Stop()

	xycond.ExpectNotPanic(func() {
		sched.After(time.Nanosecond) <- xysched.NewTask(func() {})
		time.Sleep(time.Millisecond)
	}).Test(t)
}

func TestSchedulerAfterButCloseSoon(t *testing.T) {
	var sched = xysched.NewScheduler("")

	xycond.ExpectNotPanic(func() {
		sched.After(time.Second) <- xysched.NewTask(func() {})
		sched.Stop()
		time.Sleep(time.Millisecond)
	}).Test(t)
}

func TestSchedulerAfterWithNegativeDuration(t *testing.T) {
	var sched = xysched.NewScheduler("")
	defer sched.Stop()

	xycond.ExpectNotPanic(func() {
		sched.After(-1) <- xysched.NewTask(func() {})
	}).Test(t)
}

func TestSchedulerAfterMacros(t *testing.T) {
	var sched = xysched.NewScheduler("")
	defer sched.Stop()

	xycond.ExpectNotPanic(func() {
		sched.Now() <- xysched.NewTask(func() {})
	}).Test(t)

	xycond.ExpectNotPanic(func() {
		sched.At(time.Now()) <- xysched.NewTask(func() {})
	}).Test(t)
}

func TestSchedulerConcurrent(t *testing.T) {
	var sched = xysched.NewScheduler("")
	defer sched.Stop()

	xycond.ExpectNotPanic(func() {
		sched.Singleton()
	}).Test(t)

	xycond.ExpectNotPanic(func() {
		sched.Concurrent(10)
	}).Test(t)

	xycond.ExpectNotPanic(func() {
		sched.Now() <- xysched.NewTask(func() {})
		time.Sleep(time.Millisecond)
	}).Test(t)
}
