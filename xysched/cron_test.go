package xysched_test

import (
	"testing"
	"time"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xysched"
)

func TestNewCron(t *testing.T) {
	xycond.ExpectNotPanic(func() {
		xysched.NewCron(func() {})
	}).Test(t)

	xycond.ExpectNotPanic(func() {
		xysched.NewCron(func(a, b int) {}, 1, 2)
	}).Test(t)
}

func TestNewCronWithMismatchedParameterNumber(t *testing.T) {
	xycond.ExpectPanic(func() {
		xysched.NewCron(func(a, b int) {}, 1, 2, 3)
	}).Test(t)
}

func TestNewCronWithMismatchedParameterType(t *testing.T) {
	xycond.ExpectPanic(func() {
		xysched.NewCron(func(a, b int) {}, 1)
	}).Test(t)

	xycond.ExpectPanic(func() {
		xysched.NewCron(func(a, b int) {}, 1, "a")
	}).Test(t)
}

func TestNewCronWithVariadicFunction(t *testing.T) {
	xycond.ExpectNotPanic(func() {
		xysched.NewCron(func(a ...int) {}, 1, 2, 3, 4)
	}).Test(t)

	xycond.ExpectNotPanic(func() {
		xysched.NewCron(func(a ...int) {})
	}).Test(t)

	xycond.ExpectPanic(func() {
		xysched.NewCron(func(a int, b ...int) {})
	}).Test(t)
}

func TestCronPeriodic(t *testing.T) {
	var c = xysched.NewCron(func() {})

	xycond.ExpectNotPanic(func() {
		c.Every(3 * time.Hour)
	}).Test(t)

	xycond.ExpectPanic(func() {
		c.Every(-3 * time.Hour)
	}).Test(t)
}

func TestCronMacroPeriodic(t *testing.T) {
	var c = xysched.NewCron(func() {})

	xycond.ExpectNotPanic(func() {
		c.Secondly()
	}).Test(t)

	xycond.ExpectNotPanic(func() {
		c.Minutely()
	}).Test(t)

	xycond.ExpectNotPanic(func() {
		c.Hourly()
	}).Test(t)

	xycond.ExpectNotPanic(func() {
		c.Daily()
	}).Test(t)
}

func TestCronTimes(t *testing.T) {
	var c = xysched.NewCron(func() {})

	xycond.ExpectNotPanic(func() {
		c.Times(5)
	}).Test(t)
}

func TestCronMacroTimes(t *testing.T) {
	var c = xysched.NewCron(func() {})

	xycond.ExpectNotPanic(func() {
		c.Once()
	}).Test(t)

	xycond.ExpectNotPanic(func() {
		c.Twice()
	}).Test(t)
}

func TestCronFinish(t *testing.T) {
	var cron = xysched.NewCron(func() {})
	var callback *xysched.Task

	xycond.ExpectNotPanic(func() {
		callback = cron.Finish(func(a, b int) {}, 1, 2)
	}).Test(t)

	xycond.ExpectNotNil(callback).Test(t)

	xycond.ExpectNotPanic(func() {
		cron.Finish(xysched.NewTask(func(a, b int) {}, 1, 2))
	}).Test(t)

	xycond.ExpectNotPanic(func() {
		callback = cron.Finish(xysched.NewCron(func() {}))
	}).Test(t)

	xycond.ExpectNil(callback).Test(t)

	xycond.ExpectPanic(func() {
		cron.Finish(callback, 1, 2)
	}).Test(t)
}

func TestCronStop(t *testing.T) {
	var sched = xysched.NewScheduler("")
	defer sched.Stop()

	var captured int
	var cron = xysched.NewCron(func() { captured++ })
	cron.Every(time.Millisecond)
	cron.Times(10)
	sched.Now() <- cron
	cron.Stop()
	time.Sleep(time.Second)

	xycond.ExpectLessThan(captured, 10).Test(t)
}
