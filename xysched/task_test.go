package xysched_test

import (
	"testing"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xysched"
)

func TestNewTask(t *testing.T) {
	xycond.ExpectNotPanic(func() {
		xysched.NewTask(func() {})
	}).Test(t)

	xycond.ExpectNotPanic(func() {
		xysched.NewTask(func(a, b int) {}, 1, 2)
	}).Test(t)
}

func TestNewTaskWithMismatchedParameterNumber(t *testing.T) {
	xycond.ExpectPanic(func() {
		xysched.NewTask(func(a, b int) {}, 1, 2, 3)
	}).Test(t)
}

func TestNewTaskWithMismatchedParameterType(t *testing.T) {
	xycond.ExpectPanic(func() {
		xysched.NewTask(func(a, b int) {}, 1)
	}).Test(t)

	xycond.ExpectPanic(func() {
		xysched.NewTask(func(a, b int) {}, 1, "a")
	}).Test(t)
}

func TestNewTaskWithVariadicFunction(t *testing.T) {
	xycond.ExpectNotPanic(func() {
		xysched.NewTask(func(a ...int) {}, 1, 2, 3, 4)
	}).Test(t)

	xycond.ExpectNotPanic(func() {
		xysched.NewTask(func(a ...int) {})
	}).Test(t)

	xycond.ExpectPanic(func() {
		xysched.NewTask(func(a int, b ...int) {})
	}).Test(t)
}

func TestTaskCallback(t *testing.T) {
	var task = xysched.NewTask(func() {})
	var callback *xysched.Task

	xycond.ExpectNotPanic(func() {
		callback = task.Callback(func(a, b int) {}, 1, 2)
	}).Test(t)

	xycond.ExpectNotNil(callback).Test(t)

	xycond.ExpectNotPanic(func() {
		task.Callback(xysched.NewTask(func(a, b int) {}, 1, 2))
	}).Test(t)

	xycond.ExpectNotPanic(func() {
		callback = task.Callback(xysched.NewCron(func() {}))
	}).Test(t)

	xycond.ExpectNil(callback).Test(t)

	xycond.ExpectPanic(func() {
		task.Callback(callback, 1, 2)
	}).Test(t)
}

func TestTaskThen(t *testing.T) {
	var task = xysched.NewTask(func() string { return "" })

	xycond.ExpectNotPanic(func() {
		task.Then(func(string) {})
	}).Test(t)

	xycond.ExpectPanic(func() {
		task.Then(func() {})
	}).Test(t)
}

func TestTaskCatch(t *testing.T) {
	var task = xysched.NewTask(func() string { return "" })

	xycond.ExpectNotPanic(func() {
		task.Catch(func(error) {})
	}).Test(t)

	xycond.ExpectPanic(func() {
		task.Catch(func(string) {})
	}).Test(t)

	xycond.ExpectPanic(func() {
		task.Catch(func() {})
	}).Test(t)
}
