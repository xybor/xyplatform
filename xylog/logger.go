package xylog

import (
	"strings"

	"github.com/xybor/xyplatform"
)

type logger struct {
	module xyplatform.Module
	format string
	allow  map[string]bool
	output writer
}

// Log is a general log function allowing you to print a message with a custom
// log level.
func (lg logger) Log(level string, msg string, a ...interface{}) {
	level = strings.ToUpper(level)

	_, all := lg.allow["ALL"]
	_, ok := lg.allow[level]
	if !ok && !all {
		return
	}

	log := replaceLog(lg.format, lg.module, level, msg, a...)

	lg.output.write(log)
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

// Log normal actions, such as start or stop a process.
func (lg logger) Info(msg string, a ...interface{}) {
	lg.Log("INFO", msg, a...)
}

// Log errors but the program can recover, or something leads to unexpected
// behaviors.
func (lg logger) Warn(msg string, a ...interface{}) {
	lg.Log("WARN", msg, a...)
}

// Log errors crashing the program, but errors come from external affects.
// For example, FileNotFound or DatabaseError.
func (lg logger) Error(msg string, a ...interface{}) {
	lg.Log("ERROR", msg, a...)
}

// Log errors crashing the program, and errors is internal issues.
func (lg logger) Critical(msg string, a ...interface{}) {
	lg.Log("CRITICAL", msg, a...)
}

// Log helpful or diagnostic information for debugging.
func (lg logger) Debug(msg string, a ...interface{}) {
	lg.Log("DEBUG", msg, a...)
}

// Log information helps determine where the problem occurs.
func (lg logger) Trace(msg string, a ...interface{}) {
	lg.Log("TRACE", msg, a...)
}
