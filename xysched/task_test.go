package xysched_test

import (
	"testing"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xysched"
)

func TestNewTask(t *testing.T) {
	xycond.MustNotPanic(func() {
		xysched.NewTask(func() {})
	}).Test(t, "A panic occurred")

	xycond.MustNotPanic(func() {
		xysched.NewTask(func(a, b int) {}, 1, 2)
	}).Test(t, "A panic occurred")
}

func TestNewTaskWithMismatchedParameterNumber(t *testing.T) {
	xycond.MustPanic(func() {
		xysched.NewTask(func(a, b int) {}, 1, 2, 3)
	}).Test(t, "Expected a panic, but not found")
}

func TestNewTaskWithMismatchedParameterType(t *testing.T) {
	xycond.MustPanic(func() {
		xysched.NewTask(func(a, b int) {}, 1)
	}).Test(t, "Expected a panic, but not found")

	xycond.MustPanic(func() {
		xysched.NewTask(func(a, b int) {}, 1, "a")
	}).Test(t, "Expected a panic, but not found")
}

func TestNewTaskWithVariadicFunction(t *testing.T) {
	xycond.MustNotPanic(func() {
		xysched.NewTask(func(a ...int) {}, 1, 2, 3, 4)
	}).Test(t, "A panic occurred")

	xycond.MustNotPanic(func() {
		xysched.NewTask(func(a ...int) {})
	}).Test(t, "A panic occurred")

	xycond.MustPanic(func() {
		xysched.NewTask(func(a int, b ...int) {})
	}).Test(t, "Expected a panic, but not found")
}

func TestTaskCallback(t *testing.T) {
	var task = xysched.NewTask(func() {})
	var callback *xysched.Task

	xycond.MustNotPanic(func() {
		callback = task.Callback(func(a, b int) {}, 1, 2)
	}).Test(t, "A panic occurred")

	xycond.MustNotNil(callback).Test(t, "Expected a not nil callback")

	xycond.MustNotPanic(func() {
		task.Callback(xysched.NewTask(func(a, b int) {}, 1, 2))
	}).Test(t, "A panic occurred")

	xycond.MustNotPanic(func() {
		callback = task.Callback(xysched.NewCron(func() {}))
	}).Test(t, "A panic occurred")

	xycond.MustNil(callback).Test(t, "Expected a nil callback")

	xycond.MustPanic(func() {
		task.Callback(callback, 1, 2)
	}).Test(t, "Expected a panic, but not found")
}

func TestTaskThen(t *testing.T) {
	var task = xysched.NewTask(func() string { return "" })

	xycond.MustNotPanic(func() {
		task.Then(func(string) {})
	}).Test(t, "A panic occurred")

	xycond.MustPanic(func() {
		task.Then(func() {})
	}).Test(t, "Expected a panic, but not found")
}

func TestTaskCatch(t *testing.T) {
	var task = xysched.NewTask(func() string { return "" })

	xycond.MustNotPanic(func() {
		task.Catch(func(error) {})
	}).Test(t, "A panic occurred")

	xycond.MustPanic(func() {
		task.Catch(func(string) {})
	}).Test(t, "Expected a panic, but not found")

	xycond.MustPanic(func() {
		task.Catch(func() {})
	}).Test(t, "Expected a panic, but not found")
}
