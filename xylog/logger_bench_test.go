package xylog_test

import (
	"testing"

	"github.com/xybor/xyplatform/xylog"
)

func BenchmarkLoggerWithoutLog(b *testing.B) {
	xylog.SetLevel(xylog.CRITICAL)
	for i := 0; i < b.N; i++ {
		xylog.Debug("msg")
	}
}

func BenchmarkLoggerWithOneHandler(b *testing.B) {
	var handler = xylog.NewHandler("", &CapturedEmitter{})
	handler.SetLevel(xylog.DEBUG)
	xylog.SetLevel(xylog.DEBUG)
	xylog.AddHandler(handler)
	for i := 0; i < b.N; i++ {
		xylog.Critical("msg")
	}
	xylog.RemoveHandler(handler)
}

func BenchmarkLoggerWithMultiHandler(b *testing.B) {
	for i := 0; i < 100; i++ {
		var handler = xylog.NewHandler("", &CapturedEmitter{})
		handler.SetLevel(xylog.DEBUG)
		xylog.AddHandler(handler)
	}
	xylog.SetLevel(xylog.DEBUG)
	for i := 0; i < b.N; i++ {
		xylog.Critical("msg")
	}
	for _, h := range xylog.GetHandlers() {
		xylog.RemoveHandler(h)
	}
}
