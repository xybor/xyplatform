package xylog_test

import (
	"os"
	"testing"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xyerror"
	"github.com/xybor/xyplatform/xylog"
)

type ErrorWriter struct{}

func (ew *ErrorWriter) Write(p []byte) (n int, err error) {
	return 0, xyerror.UnknownError.New("unknown")
}

func (ew *ErrorWriter) Close() error {
	return nil
}

func TestNewStreamEmitterWithNil(t *testing.T) {
	xycond.MustNotPanic(func() {
		xylog.NewStreamEmitter(nil)
	}).Test(t, "A panic occurred")
}

func TestStreamEmitterEmit(t *testing.T) {
	var emitter = xylog.NewStreamEmitter(os.Stderr)
	xycond.MustNotPanic(func() {
		emitter.Emit(xylog.LogRecord{})
	}).Test(t, "A panic occurred")
}

func TestStreamEmitterEmitError(t *testing.T) {
	var emitter = xylog.NewStreamEmitter(&ErrorWriter{})
	xycond.MustPanic(func() {
		emitter.Emit(xylog.LogRecord{})
	}).Test(t, "Expect a panic, but not found")
}

func TestFileEmitter(t *testing.T) {
	var emitter = xylog.NewFileEmitter("a.log")
	xycond.MustNotPanic(func() {
		emitter.Emit(xylog.LogRecord{})
	}).Test(t, "A panic occurred")
}

func TestSizeRotatingFileEmitter(t *testing.T) {
	var emitter = xylog.NewSizeRotatingFileEmitter("a.log", 100, 1)
	xycond.MustNotPanic(func() {
		emitter.Emit(xylog.LogRecord{})
	}).Test(t, "A panic occurred")
}

func TestTimeRotatingFileEmitter(t *testing.T) {
	var emitter = xylog.NewTimeRotatingFileEmitter("a.log", 100, 1)
	xycond.MustNotPanic(func() {
		emitter.Emit(xylog.LogRecord{})
	}).Test(t, "A panic occurred")
}
