package xylog_test

import (
	"testing"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xylog"
)

func testRootLogger(t *testing.T, f func(int)) {
	var handler = xylog.NewHandler("", &CapturedEmitter{})
	handler.SetLevel(xylog.DEBUG)
	xylog.AddHandler(handler)
	defer xylog.RemoveHandler(handler)

	var loggerLevel = xylog.INFO
	xylog.SetLevel(loggerLevel)

	f(loggerLevel)

}

func TestRootHandler(t *testing.T) {
	var handler = xylog.NewHandler("", xylog.StdoutEmitter)
	xycond.ExpectNotPanic(func() {
		xylog.AddHandler(handler)
		xylog.RemoveHandler(handler)
	}).Test(t)
}

func TestRootFilter(t *testing.T) {
	var filter = &NameFilter{}
	xycond.ExpectNotPanic(func() {
		xylog.AddFilter(filter)
		xylog.RemoveFilter(filter)
	}).Test(t)
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
		xycond.ExpectNotPanic(func() {
			xylog.SetLevel(levels[i])
		}).Test(t)
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

	testRootLogger(t, func(loggerLevel int) {
		for i := range levels {
			checkLogOutput(t, func() { xylog.Logf(levels[i], "foo") }, "foo",
				levels[i], loggerLevel)
			checkLogOutput(t, func() { xylog.Log(levels[i], "foo") }, "foo",
				levels[i], loggerLevel)
		}
	})
}

func TestRootLogfMethods(t *testing.T) {
	var methods = map[int]func(string, ...any){
		xylog.DEBUG:    xylog.Debugf,
		xylog.INFO:     xylog.Infof,
		xylog.WARN:     xylog.Warnf,
		xylog.ERROR:    xylog.Errorf,
		xylog.CRITICAL: xylog.Criticalf,
	}

	testRootLogger(t, func(loggerLevel int) {
		for level, method := range methods {
			checkLogOutput(t, func() { method("foo") }, "foo", level, loggerLevel)
		}
	})
}

func TestRootLogMethods(t *testing.T) {
	var methods = map[int]func(...any){
		xylog.DEBUG:    xylog.Debug,
		xylog.INFO:     xylog.Info,
		xylog.WARN:     xylog.Warn,
		xylog.ERROR:    xylog.Error,
		xylog.CRITICAL: xylog.Critical,
	}

	testRootLogger(t, func(loggerLevel int) {
		for level, method := range methods {
			checkLogOutput(t, func() { method("foo") }, "foo", level, loggerLevel)
		}
	})
}
