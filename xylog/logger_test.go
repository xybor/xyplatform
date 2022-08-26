package xylog_test

import (
	"testing"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xylog"
)

type CapturedEmitter struct{}

func (h *CapturedEmitter) Emit(record xylog.LogRecord) {
	capturedOutput = record.Message
}

func (h *CapturedEmitter) SetFormatter(xylog.Formatter) {}

type NameFilter struct {
	name string
}

func (f *NameFilter) Filter(r xylog.LogRecord) bool {
	return f.name == r.Name
}

// capturedOutput is the output which CapturedHandler printed.
var capturedOutput string

// validCustomLevels will be added to xylog's level system.
var validCustomLevels = []int{-1, 25, 100}

// invalidCustomLevels will not be added to xylog's level system.
var invalidCustomLevels = []int{-10, 35, 75}

func init() {
	for i := range validCustomLevels {
		xylog.AddLevel(validCustomLevels[i], "")
	}
}

func TestLoggerValidCustomLevel(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())

	// Test not panic.
	for i := range validCustomLevels {
		logger.SetLevel(validCustomLevels[i])
	}
}

func TestLoggerInvalidCustomLevel(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())

	for i := range invalidCustomLevels {
		xycond.MustPanic(func() {
			logger.SetLevel(invalidCustomLevels[i])
		}).Testf(t, "Expected a panic, but not found")
	}
}

func TestLoggerHandler(t *testing.T) {
	var expectedHandler = xylog.NewHandler("", &CapturedEmitter{})
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(expectedHandler)
	var handlers = logger.GetHandlers()
	xycond.MustEqual(len(handlers), 1).
		Testf(t, "Expected one handler, but got %d", len(handlers))
	xycond.MustEqual(handlers[0], expectedHandler).
		Testf(t, "Expected the same handler, but got different ones")

	logger.RemoveHandler(expectedHandler)
	handlers = logger.GetHandlers()
	xycond.MustEmpty(handlers).
		Testf(t, "Expected no handler, but got %d", len(handlers))
}

func TestLoggerAddHandlerNil(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())
	xycond.MustPanic(func() {
		logger.AddHandler(nil)
	}).Test(t, "Expected a panic, but not found")
}

func TestLoggerRemoveInvalidHandler(t *testing.T) {
	var expectedHandler = xylog.NewHandler("", &CapturedEmitter{})
	var logger = xylog.GetLogger(t.Name())
	logger.RemoveHandler(expectedHandler)
	var handlers = logger.GetHandlers()
	xycond.MustEmpty(handlers).
		Testf(t, "Expected no handler, but got %d", len(handlers))
}

func TestLoggerLogMethods(t *testing.T) {
	var expectedMessage = "foo"
	var loggerLevel = xylog.WARN
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(xylog.NewHandler("", &CapturedEmitter{}))
	logger.SetLevel(loggerLevel)

	var loggerMethods = map[int]func(string, ...any){
		xylog.DEBUG:    logger.Debug,
		xylog.INFO:     logger.Info,
		xylog.WARN:     logger.Warn,
		xylog.ERROR:    logger.Error,
		xylog.CRITICAL: logger.Critical,
	}

	for level, method := range loggerMethods {
		capturedOutput = ""
		method(expectedMessage)
		if level < loggerLevel {
			xycond.MustEmpty(capturedOutput).
				Testf(t, "Expected an empty output, but got %s", capturedOutput)
		} else {
			xycond.MustEqual(capturedOutput, expectedMessage).
				Testf(t, "Expected output %s, but got empty", expectedMessage)
		}
	}
}

func TestLoggerCallHandlerHierarchy(t *testing.T) {
	var expectedMessage = "foo"
	var handler = xylog.NewHandler("", &CapturedEmitter{})
	var logger = xylog.GetLogger(t.Name())
	logger.SetLevel(xylog.DEBUG)
	logger.AddHandler(handler)

	logger = xylog.GetLogger(t.Name() + ".main")
	capturedOutput = ""
	logger.Info(expectedMessage)
	xycond.MustEqual(capturedOutput, expectedMessage).
		Testf(t, "%s != %s", capturedOutput, expectedMessage)
}

func TestLoggerLogNoHandler(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())
	logger.SetLevel(xylog.DEBUG)

	xycond.MustNotPanic(func() {
		logger.Info("foo")
	}).Test(t, "A panic occurred")
}

func TestLoggerLogNotSetLevel(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())

	xycond.MustNotPanic(func() {
		logger.Info("foo")
	}).Test(t, "A panic occurred")
}

func TestLoggerLogInvalidCustomLevel(t *testing.T) {
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(xylog.NewHandler("", &CapturedEmitter{}))
	logger.SetLevel(xylog.DEBUG)

	for i := range invalidCustomLevels {
		xycond.MustPanic(func() {
			logger.Log(invalidCustomLevels[i], "msg")
		}).Test(t, "Expected a panic, but not found")
	}
}

func TestLoggerLogValidCustomLevel(t *testing.T) {
	var expectedMessage = "foo"
	var loggerLevel = xylog.DEBUG
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(xylog.NewHandler("", &CapturedEmitter{}))
	logger.SetLevel(loggerLevel)

	for i := range validCustomLevels {
		capturedOutput = ""
		logger.Log(validCustomLevels[i], expectedMessage)
		if validCustomLevels[i] < loggerLevel {
			xycond.MustEmpty(capturedOutput).
				Testf(t, "Expected an empty output, but got %s", capturedOutput)
		} else {
			xycond.MustEqual(capturedOutput, expectedMessage).
				Testf(t, "Expected output %s, but got empty", expectedMessage)
		}
	}
}

func TestLoggerFilter(t *testing.T) {
	var expectedFilter = &NameFilter{}
	var logger = xylog.GetLogger(t.Name())
	logger.AddFilter(expectedFilter)
	var filters = logger.GetFilters()
	xycond.MustTrue(len(filters) == 1).
		Testf(t, "Expected one elements, but got %d", len(filters))
	xycond.MustEqual(filters[0], expectedFilter).
		Test(t, "Expected the same filter, but got different")
	xycond.MustNotPanic(func() {
		logger.RemoveFilter(expectedFilter)
	}).Test(t, "A panic occurred")
}

func TestLoggerFilterLog(t *testing.T) {
	var expectedMessage = "foo"
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(xylog.NewHandler("", &CapturedEmitter{}))
	logger.SetLevel(xylog.DEBUG)

	capturedOutput = ""
	logger.AddFilter(&NameFilter{t.Name()})
	logger.Debug(expectedMessage)
	xycond.MustEqual(capturedOutput, expectedMessage).Testf(t,
		"Expected output %s, but got %s", expectedMessage, capturedOutput)

	capturedOutput = ""
	logger.AddFilter(&NameFilter{"bar name"})
	logger.Debug(expectedMessage)
	xycond.MustEmpty(capturedOutput).
		Testf(t, "Expected an empty output, but got %s", capturedOutput)
}

func TestLoggerAddExtra(t *testing.T) {
	var handler = xylog.NewHandler("", &CapturedEmitter{})
	var logger = xylog.GetLogger(t.Name())
	logger.SetLevel(xylog.DEBUG)
	logger.AddHandler(handler)
	logger.AddExtra("bar", "something")

	capturedOutput = ""
	logger.Info("foo")
	xycond.MustEqual(capturedOutput, "bar=something foo").
		Testf(t, "%s != %s", capturedOutput, "bar=something foo")
}
