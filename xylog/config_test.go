package xylog_test

import (
	"os"
	"testing"
	"time"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xylog"
)

var registeredHandlerNames = []string{"foo", "bar"}
var notRegisteredHandlerNames = []string{"foobar", "barfoo"}

func init() {
	for i := range registeredHandlerNames {
		xylog.NewHandler(registeredHandlerNames[i], xylog.StdoutEmitter)
	}
}

func TestGetLogger(t *testing.T) {
	var names = []string{"", "foo", "foo.bar"}
	for i := range names {
		var logger1 = xylog.GetLogger(names[i])
		var logger2 = xylog.GetLogger(names[i])
		xycond.MustEqual(logger1, logger2).
			Test(t, "Given a name, expect the same logger, but got different")
	}
}
func TestGetHandler(t *testing.T) {
	for i := range registeredHandlerNames {
		var handlerA = xylog.GetHandler(registeredHandlerNames[i])
		var handlerB = xylog.GetHandler(registeredHandlerNames[i])
		xycond.MustEqual(handlerA, handlerB).
			Test(t, "Expected the same handler, but got different")
	}
}

func TestGetHandlerDiff(t *testing.T) {
	var handlerA = xylog.GetHandler(registeredHandlerNames[0])
	var handlerB = xylog.GetHandler(registeredHandlerNames[1])
	xycond.MustNotEqual(handlerA, handlerB).
		Test(t, "Expected different handlers, but got one")
}

func TestGetHandlerNotRegisterBefore(t *testing.T) {
	for i := range notRegisteredHandlerNames {
		var handler = xylog.GetHandler(notRegisteredHandlerNames[i])
		xycond.MustNil(handler).
			Test(t, "Expected a nil handler, but got not-nil")
	}
}

func TestSetTimeLayout(t *testing.T) {
	xycond.MustNotPanic(func() {
		xylog.SetTimeLayout("123")
		xylog.SetTimeLayout(time.RFC3339Nano)
	}).Test(t, "A panic occurred")
}

func TestSetFileFlag(t *testing.T) {
	xycond.MustNotPanic(func() {
		xylog.SetFileFlag(os.O_WRONLY | os.O_APPEND | os.O_CREATE)
	}).Test(t, "A panic occurred")
}

func TestSetFilePerm(t *testing.T) {
	xycond.MustNotPanic(func() {
		xylog.SetFilePerm(0666)
	}).Test(t, "A panic occurred")
}

func TestSetSkipCall(t *testing.T) {
	xycond.MustNotPanic(func() {
		xylog.SetSkipCall(2)
	}).Test(t, "A panic occurred")
}
