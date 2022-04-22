package xylog

import (
	"github.com/xybor/xyplatform"
)

type logger struct {
	module xyplatform.Module
	format string
	level  uint
	output writer
}

// Log is a general log function allowing you to print a message with a custom
// log level.
func (lg logger) Log(level uint, msg string, a ...interface{}) {
	if level < lg.level {
		return
	}

	checkLevel(level)
	levelName := levelManager[level]

	log := makeRecord(lg.format, lg.module, levelName, msg, a...)

	lg.output.write(log)
}

// LogT is a template log function allowing you to print a template message
// with a custom log level.
func (lg logger) LogT(level uint, event string, t T) {
	if level < lg.level {
		return
	}

	checkLevel(level)
	levelName := levelManager[level]

	record := makeRecord(lg.format, lg.module, levelName, parseTemplate(event, t))

	lg.output.write(record)
}

// Config is a function allowing you to configure the logger. The returned
// value of this function is meaningless, it only helps run it in the global
// scope.
func (lg *logger) Config(configurators ...configurator) bool {
	for _, c := range configurators {
		c.apply(lg)
	}

	return true
}

// Trace logs very detailed information for debugging.
func (lg logger) Trace(msg string, a ...interface{}) {
	lg.Log(TRACE, msg, a...)
}

func (lg logger) TraceT(event string, t T) {
	lg.LogT(TRACE, event, t)
}

// Debug logs helpful and diagnostic information for debugging.
func (lg logger) Debug(msg string, a ...interface{}) {
	lg.Log(DEBUG, msg, a...)
}

func (lg logger) DebugT(event string, t T) {
	lg.LogT(DEBUG, event, t)
}

// Info logs normal actions, such as start and stop a process.
func (lg logger) Info(msg string, a ...interface{}) {
	lg.Log(INFO, msg, a...)
}

func (lg logger) InfoT(event string, t T) {
	lg.LogT(INFO, event, t)
}

// Warn logs errors the program can recover and continue after that, or
// something leads to unexpected behaviors.
func (lg logger) Warn(msg string, a ...interface{}) {
	lg.Log(WARN, msg, a...)
}

func (lg logger) WarnT(event string, t T) {
	lg.LogT(WARN, event, t)
}

// Error logs errors causing the function or operation to be stopped, but it
// could be fixed later.
func (lg logger) Error(msg string, a ...interface{}) {
	lg.Log(ERROR, msg, a...)
}

func (lg logger) ErrorT(event string, t T) {
	lg.LogT(ERROR, event, t)
}

// Critical logs errors causing the application or program to be stopped, or
// something needs to be fixed immediately.
func (lg logger) Critical(msg string, a ...interface{}) {
	lg.Log(CRITICAL, msg, a...)
}

func (lg logger) CriticalT(event string, t T) {
	lg.LogT(CRITICAL, event, t)
}
