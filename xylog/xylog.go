package xylog

import (
	"fmt"
	"time"

	"github.com/xybor/xyplatform"
)

var (
	Byte  = int64(1)
	KByte = 1024 * Byte
	MByte = 1024 * KByte
	GByte = 1024 * MByte
)

var (
	Minute = time.Minute
	Hour   = time.Hour
	Day    = 24 * time.Hour
	Week   = 7 * Day
)

var TimeFormat = "01-02-2006.15-04-05"

// Log is a general log function allowing you to print a message with a custom
// log level.
func Log(m xyplatform.Module, level string, msg string, a ...interface{}) {
	if _, ok := manager[m]; !ok {
		msg := fmt.Sprintf("Module %s had not registered yet", m)
		panic(msg)
	}

	manager[m].Log(level, msg, a...)
}

// Config is a function allowing you to configure the logger. The returned
// value of this function is meaningless, it only helps run it in the global
// scope.
func Config(m xyplatform.Module, configurators ...configurator) bool {
	if _, ok := manager[m]; !ok {
		msg := fmt.Sprintf("Module %s had not registered yet", m)
		panic(msg)
	}

	return manager[m].Config(configurators...)
}

// Log errors but the program can recover, or something leads to unexpected
// behaviors.
func Info(m xyplatform.Module, msg string, a ...interface{}) {
	Log(m, "INFO", msg, a...)
}

// Log errors crashing the program, but errors come from external affects.
// For example, FileNotFound or DatabaseError.
func Warn(m xyplatform.Module, msg string, a ...interface{}) {
	Log(m, "WARN", msg, a...)
}

// Log errors crashing the program, but errors come from external affects.
// For example, FileNotFound or DatabaseError.
func Error(m xyplatform.Module, msg string, a ...interface{}) {
	Log(m, "ERROR", msg, a...)
}

// Log errors crashing the program, and errors is internal issues.
func Critical(m xyplatform.Module, msg string, a ...interface{}) {
	Log(m, "CRITICAL", msg, a...)
}

// Log helpful or diagnostic information for debugging.
func Debug(m xyplatform.Module, msg string, a ...interface{}) {
	Log(m, "DEBUG", msg, a...)
}

// Log information helps determine where the problem occurs.
func Trace(m xyplatform.Module, msg string, a ...interface{}) {
	Log(m, "TRACE", msg, a...)
}
