package xylog

import (
	"time"

	"github.com/xybor/xyplatform"
)

const (
	Byte  = int64(1)
	KByte = 1024 * Byte
	MByte = 1024 * KByte
	GByte = 1024 * MByte
)

const (
	Minute = time.Minute
	Hour   = time.Hour
	Day    = 24 * time.Hour
	Week   = 7 * Day
)

var (
	TRACE    = RegisterLevel(10, "TRACE")
	DEBUG    = RegisterLevel(20, "DEBUG")
	INFO     = RegisterLevel(50, "INFO")
	WARN     = RegisterLevel(70, "WARN")
	ERROR    = RegisterLevel(80, "ERROR")
	CRITICAL = RegisterLevel(90, "CRITICAL")
)

var TimeFormat = "01-02-2006.15-04-05"

// Log is a general log function allowing you to print a message with a custom
// log level.
func Log(m xyplatform.Module, level uint, msg string, a ...interface{}) {
	Logger(m).Log(level, msg, a...)
}

// Config is a function allowing you to configure the logger. The returned
// value of this function is meaningless, it only helps run it in the global
// scope.
func Config(m xyplatform.Module, configurators ...configurator) bool {
	return Logger(m).Config(configurators...)
}

// Trace logs very detailed information for debugging.
func Trace(m xyplatform.Module, msg string, a ...interface{}) {
	Logger(m).Trace(msg, a...)
}

// Debug logs helpful and diagnostic information for debugging.
func Debug(m xyplatform.Module, msg string, a ...interface{}) {
	Logger(m).Debug(msg, a...)
}

// Info logs normal actions, such as start and stop a process.
func Info(m xyplatform.Module, msg string, a ...interface{}) {
	Logger(m).Info(msg, a...)
}

// Warn logs errors the program can recover and continue after that, or
// something leads to unexpected behaviors.
func Warn(m xyplatform.Module, msg string, a ...interface{}) {
	Logger(m).Warn(msg, a...)
}

// Error logs errors causing the function or operation to be stopped, but it
// could be fixed later.
func Error(m xyplatform.Module, msg string, a ...interface{}) {
	Logger(m).Error(msg, a...)
}

// Critical logs errors causing the application or program to be stopped, or
// something needs to be fixed immediately.
func Critical(m xyplatform.Module, msg string, a ...interface{}) {
	Logger(m).Critical(msg, a...)
}
