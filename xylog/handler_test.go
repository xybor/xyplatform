package xylog_test

import (
	"testing"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xylog"
)

func TestNewHandlerWithEmptyName(t *testing.T) {
	var handlerA = xylog.NewHandler("", xylog.StdoutEmitter)
	var handlerB = xylog.NewHandler("", xylog.StdoutEmitter)
	xycond.MustNotEqual(handlerA, handlerB).
		Test(t, "Expected different handlers, but got one")
}

func TestHandlerSetFormatter(t *testing.T) {
	var handler = xylog.NewHandler(t.Name(), xylog.StdoutEmitter)
	xycond.MustNotPanic(func() {
		handler.SetFormatter(xylog.NewTextFormatter(""))
	}).Test(t, "A panic occurred")
}

func TestHandlerSetFormatterNil(t *testing.T) {
	var handler = xylog.NewHandler(t.Name(), xylog.StdoutEmitter)
	xycond.MustPanic(func() {
		handler.SetFormatter(nil)
	}).Test(t, "Expected a panic, but not found")
}

func TestHandlerFilter(t *testing.T) {
	var expectedFilter = &NameFilter{}
	var handler = xylog.NewHandler(t.Name(), xylog.StdoutEmitter)
	handler.AddFilter(expectedFilter)
	var filters = handler.GetFilters()
	xycond.MustTrue(len(filters) == 1).
		Testf(t, "Expected one elements, but got %d", len(filters))
	xycond.MustEqual(filters[0], expectedFilter).
		Test(t, "Expected the same filter, but got different")
	xycond.MustNotPanic(func() {
		handler.RemoveFilter(expectedFilter)
	}).Test(t, "A panic occurred")
}

func TestHandlerFilterLog(t *testing.T) {
	var expectedMessage = "foo foo"
	var tests = []struct {
		handlerName string
		filterName  string
	}{
		{t.Name() + "1", t.Name()},
		{t.Name() + "2", "foobar"},
	}

	for i := range tests {
		var handler = xylog.NewHandler(tests[i].handlerName, &CapturedEmitter{})
		handler.AddFilter(&NameFilter{tests[i].filterName})
		handler.SetLevel(xylog.DEBUG)

		var logger = xylog.GetLogger(t.Name())
		logger.SetLevel(xylog.DEBUG)
		logger.AddHandler(handler)
		capturedOutput = ""
		logger.Info(expectedMessage)
		if tests[i].filterName != t.Name() {
			xycond.MustEmpty(capturedOutput).Test(t, "Expected a empty string")
		} else {
			xycond.MustEqual(capturedOutput, expectedMessage).
				Testf(t, "%s != %s", capturedOutput, expectedMessage)
		}
		logger.RemoveHandler(handler)
	}
}

func TestHandlerLevel(t *testing.T) {
	var expectedMessage = "foo foo"
	var loggerLevel = xylog.INFO
	var tests = []struct {
		handlerName string
		level       int
	}{
		{t.Name() + "1", xylog.DEBUG},
		{t.Name() + "2", xylog.ERROR},
	}

	for i := range tests {
		var handler = xylog.NewHandler(tests[i].handlerName, &CapturedEmitter{})
		handler.SetLevel(tests[i].level)

		var logger = xylog.GetLogger(t.Name())
		logger.SetLevel(xylog.DEBUG)
		logger.AddHandler(handler)
		capturedOutput = ""
		logger.Log(loggerLevel, expectedMessage)
		if loggerLevel < tests[i].level {
			xycond.MustEmpty(capturedOutput).Test(t, "Expected a empty string")
		} else {
			xycond.MustEqual(capturedOutput, expectedMessage).
				Testf(t, "%s != %s", capturedOutput, expectedMessage)
		}
		logger.RemoveHandler(handler)
	}
}
