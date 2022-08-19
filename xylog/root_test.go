package xylog_test

import (
	"testing"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xylog"
)

func TestRootHandler(t *testing.T) {
	var handler = xylog.NewHandler("", xylog.StdoutEmitter)
	xycond.MustNotPanic(func() {
		xylog.AddHandler(handler)
		xylog.GetHandlers()
		xylog.RemoveHandler(handler)
	}).Test(t, "A panic occurred")
}

func TestRootFilter(t *testing.T) {
	var filter = &NameFilter{}
	xycond.MustNotPanic(func() {
		xylog.AddFilter(filter)
		xylog.GetFilters()
		xylog.RemoveFilter(filter)
	}).Test(t, "A panic occurred")
}

func TestRootSetLevel(t *testing.T) {
	var levels = []int{
		xylog.NOTSET,
		xylog.DEBUG,
		xylog.INFO,
		xylog.WARN,
		xylog.WARNING,
		xylog.ERROR,
		xylog.FATAL,
		xylog.CRITICAL,
	}

	for i := range levels {
		xycond.MustNotPanic(func() {
			xylog.SetLevel(levels[i])
		}).Test(t, "A panic occurred")
	}
}

func TestRootLog(t *testing.T) {
	var levels = []int{
		xylog.NOTSET,
		xylog.DEBUG,
		xylog.INFO,
		xylog.WARN,
		xylog.WARNING,
		xylog.ERROR,
		xylog.FATAL,
		xylog.CRITICAL,
	}

	var handler = xylog.NewHandler("", &CapturedEmitter{})
	handler.SetLevel(xylog.DEBUG)
	xylog.AddHandler(handler)

	var loggerLevel = xylog.INFO
	xylog.SetLevel(loggerLevel)

	var expectedMessage = "foo"

	for i := range levels {
		xycond.MustNotPanic(func() {
			capturedOutput = ""
			xylog.Log(levels[i], expectedMessage)
			if levels[i] < loggerLevel {
				xycond.MustEmpty(capturedOutput).Testf(t,
					"Expect an empty output, but got %s", capturedOutput)
			} else {
				xycond.MustEqual(capturedOutput, expectedMessage).
					Testf(t, "%s != %s", capturedOutput, expectedMessage)
			}
		}).Test(t, "A panic occurred")
	}

	xylog.RemoveHandler(handler)
}

func TestRootLogMethods(t *testing.T) {
	var methods = map[int]func(string, ...any){
		xylog.DEBUG:    xylog.Debug,
		xylog.INFO:     xylog.Info,
		xylog.WARN:     xylog.Warn,
		xylog.ERROR:    xylog.Error,
		xylog.CRITICAL: xylog.Critical,
	}

	var handler = xylog.NewHandler("", &CapturedEmitter{})
	handler.SetLevel(xylog.DEBUG)
	xylog.AddHandler(handler)

	var loggerLevel = xylog.INFO
	xylog.SetLevel(loggerLevel)

	var expectedMessage = "foo"

	for level, method := range methods {
		xycond.MustNotPanic(func() {
			capturedOutput = ""
			method(expectedMessage)
			if level < loggerLevel {
				xycond.MustEmpty(capturedOutput).Testf(t,
					"Expect an empty output, but got %s", capturedOutput)
			} else {
				xycond.MustEqual(capturedOutput, expectedMessage).
					Testf(t, "%s != %s", capturedOutput, expectedMessage)
			}
		}).Test(t, "A panic occurred")
	}

	xylog.RemoveHandler(handler)
}
