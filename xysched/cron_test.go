package xysched_test

import (
	"testing"
	"time"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xysched"
)

func TestNewCron(t *testing.T) {
	xycond.MustNotPanic(func() {
		xysched.NewCron(func() {})
	}).Test(t, "A panic occurred")

	xycond.MustNotPanic(func() {
		xysched.NewCron(func(a, b int) {}, 1, 2)
	}).Test(t, "A panic occurred")
}

func TestNewCronWithMismatchedParameterNumber(t *testing.T) {
	xycond.MustPanic(func() {
		xysched.NewCron(func(a, b int) {}, 1, 2, 3)
	}).Test(t, "Expected a panic, but not found")
}

func TestNewCronWithMismatchedParameterType(t *testing.T) {
	xycond.MustPanic(func() {
		xysched.NewCron(func(a, b int) {}, 1)
	}).Test(t, "Expected a panic, but not found")

	xycond.MustPanic(func() {
		xysched.NewCron(func(a, b int) {}, 1, "a")
	}).Test(t, "Expected a panic, but not found")
}

func TestNewCronWithVariadicFunction(t *testing.T) {
	xycond.MustNotPanic(func() {
		xysched.NewCron(func(a ...int) {}, 1, 2, 3, 4)
	}).Test(t, "A panic occurred")

	xycond.MustNotPanic(func() {
		xysched.NewCron(func(a ...int) {})
	}).Test(t, "A panic occurred")

	xycond.MustPanic(func() {
		xysched.NewCron(func(a int, b ...int) {})
	}).Test(t, "Expected a panic, but not found")
}

func TestCronPeriodic(t *testing.T) {
	var c = xysched.NewCron(func() {})

	xycond.MustNotPanic(func() {
		c.Every(3 * time.Hour)
	}).Test(t, "A panic occured")

	xycond.MustPanic(func() {
		c.Every(-3 * time.Hour)
	}).Test(t, "Expected a panic, but not found")
}

func TestCronMacroPeriodic(t *testing.T) {
	var c = xysched.NewCron(func() {})

	xycond.MustNotPanic(func() {
		c.Secondly()
	}).Test(t, "A panic occured")

	xycond.MustNotPanic(func() {
		c.Minutely()
	}).Test(t, "A panic occured")

	xycond.MustNotPanic(func() {
		c.Hourly()
	}).Test(t, "A panic occured")

	xycond.MustNotPanic(func() {
		c.Daily()
	}).Test(t, "A panic occured")
}

func TestCronTimes(t *testing.T) {
	var c = xysched.NewCron(func() {})

	xycond.MustNotPanic(func() {
		c.Times(5)
	}).Test(t, "A panic occured")
}

func TestCronMacroTimes(t *testing.T) {
	var c = xysched.NewCron(func() {})

	xycond.MustNotPanic(func() {
		c.Once()
	}).Test(t, "A panic occured")

	xycond.MustNotPanic(func() {
		c.Twice()
	}).Test(t, "A panic occured")
}

func TestCronFinish(t *testing.T) {
	var cron = xysched.NewCron(func() {})
	var callback *xysched.Task

	xycond.MustNotPanic(func() {
		callback = cron.Finish(func(a, b int) {}, 1, 2)
	}).Test(t, "A panic occurred")

	xycond.MustNotNil(callback).Test(t, "Expected a not nil callback")

	xycond.MustNotPanic(func() {
		cron.Finish(xysched.NewTask(func(a, b int) {}, 1, 2))
	}).Test(t, "A panic occurred")

	xycond.MustNotPanic(func() {
		callback = cron.Finish(xysched.NewCron(func() {}))
	}).Test(t, "A panic occurred")

	xycond.MustNil(callback).Test(t, "Expected a nil callback")

	xycond.MustPanic(func() {
		cron.Finish(callback, 1, 2)
	}).Test(t, "Expected a panic, but not found")
}
