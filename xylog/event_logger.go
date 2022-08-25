package xylog

import "fmt"

// EventLogger is a logger wrapper supporting to compose logging message with
// key-value pair.
type EventLogger struct {
	fields map[string]any
	lg     *Logger
}

// Field adds a key-value pair to logging message.
func (e *EventLogger) Field(key string, value any) *EventLogger {
	e.fields[key] = value
	return e
}

// compose creates the final logging messages from fields.
func (e *EventLogger) compose() string {
	var s = ""
	for k, v := range e.fields {
		s = prefixMessage(s, k+"="+fmt.Sprint(v))
	}
	return s
}

// Debug calls Log with DEBUG level.
func (e *EventLogger) Debug() {
	if e.lg.isEnabledFor(DEBUG) {
		e.lg.log(DEBUG, e.compose())
	}
}

// Info calls Log with INFO level.
func (e *EventLogger) Info() {
	if e.lg.isEnabledFor(INFO) {
		e.lg.log(INFO, e.compose())
	}
}

// Warn calls Log with WARN level.
func (e *EventLogger) Warn() {
	if e.lg.isEnabledFor(WARN) {
		e.lg.log(WARN, e.compose())
	}
}

// Warning calls Log with WARNING level.
func (e *EventLogger) Warning() {
	if e.lg.isEnabledFor(WARNING) {
		e.lg.log(WARNING, e.compose())
	}
}

// Error calls Log with ERROR level.
func (e *EventLogger) Error() {
	if e.lg.isEnabledFor(ERROR) {
		e.lg.log(ERROR, e.compose())
	}
}

// Fatal calls Log with FATAL level.
func (e *EventLogger) Fatal() {
	if e.lg.isEnabledFor(FATAL) {
		e.lg.log(FATAL, e.compose())
	}
}

// Critical calls Log with CRITICAL level.
func (e *EventLogger) Critical() {
	if e.lg.isEnabledFor(CRITICAL) {
		e.lg.log(CRITICAL, e.compose())
	}
}

// Log logs with a custom level.
func (e *EventLogger) Log(level int) {
	level = checkLevel(level)
	if e.lg.isEnabledFor(level) {
		e.lg.log(level, e.compose())
	}
}
