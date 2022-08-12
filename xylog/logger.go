// This file copied and modified comments of python logging.
package xylog

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xylock"
)

// logger represent a single logging channel. A "logging channel" indicates an
// area of an application. Exactly how an "area" is defined is up to the
// application developer. Since an application can have any number of areas,
// logging channels are identified by a unique string. Application areas can be
// nested (e.g. an area of "input processing" might include sub-areas "read CSV
// files", "read XLS files" and "read Gnumeric files"). To cater for this
// natural nesting, channel names are organized into a namespace hierarchy where
// levels are separated by periods. So in the instance given above, channel
// names might be "input" for the upper level, and "input.csv", "input.xls" and
// "input.gnu" for the sub-levels. There is no arbitrary limit to the depth of
// nesting.
type logger struct {
	filterer

	fullname string
	children map[string]*logger
	parent   *logger
	level    int
	handlers []handler
	lock     xylock.RWLock
	cache    map[int]bool
}

// newlogger creates a new logger with a name and parent. The fullname of logger
// will be concatenated by the parent's fullname. This logger will not be
// automatically added to logger hierarchy. The returned logger has no child,
// no handler, and NOTSET level.
func newlogger(name string, parent *logger) *logger {
	var c = parent
	if c != nil && c != rootLogger {
		name = c.fullname + "." + name
	}

	return &logger{
		filterer: newfilterer(),
		fullname: name,
		children: make(map[string]*logger),
		parent:   parent,
		level:    NOTSET,
		handlers: nil,
		lock:     xylock.RWLock{},
		cache:    make(map[int]bool),
	}
}

// SetLevel sets the new logging level. It also clears logging level cache of
// all loggers in program.
func (lg *logger) SetLevel(level int) {
	lg.lock.WLockFunc(func() { lg.level = checkLevel(level) })
	rootLogger.clearCache()
}

// AddHandler adds a new handler. The passed handler must be a pointer.
func (lg *logger) AddHandler(h handler) {
	xycond.IsKind(h, reflect.Pointer).Assert("h must be a pointer of handler")
	lg.lock.WLockFunc(func() { lg.handlers = append(lg.handlers, h) })
}

// RemoveHandler removes a existed handler. The passed handler must be a
// pointer.
func (lg *logger) RemoveHandler(h handler) {
	xycond.IsKind(h, reflect.Pointer).Assert("h must be a pointer of handler")
	lg.lock.WLockFunc(func() {
		for i := range lg.handlers {
			if lg.handlers[i] == h {
				lg.handlers = append(lg.handlers[:i], lg.handlers[i+1:]...)
				break
			}
		}
	})
}

// Debug calls Log with DEBUG level.
func (lg *logger) Debug(msg string, a ...any) {
	if lg.isEnabledFor(DEBUG) {
		lg.log(DEBUG, fmt.Sprintf(msg, a...))
	}
}

// Info calls Log with INFO level.
func (lg *logger) Info(msg string, a ...any) {
	if lg.isEnabledFor(INFO) {
		lg.log(INFO, fmt.Sprintf(msg, a...))
	}
}

// Warn calls Log with WARN level.
func (lg *logger) Warn(msg string, a ...any) {
	if lg.isEnabledFor(WARN) {
		lg.log(WARN, fmt.Sprintf(msg, a...))
	}
}

// Warning calls Log with WARNING level.
func (lg *logger) Warning(msg string, a ...any) {
	if lg.isEnabledFor(WARNING) {
		lg.log(WARNING, fmt.Sprintf(msg, a...))
	}
}

// Error calls Log with ERROR level.
func (lg *logger) Error(msg string, a ...any) {
	if lg.isEnabledFor(ERROR) {
		lg.log(ERROR, fmt.Sprintf(msg, a...))
	}
}

// Fatal calls Log with FATAL level.
func (lg *logger) Fatal(msg string, a ...any) {
	if lg.isEnabledFor(FATAL) {
		lg.log(FATAL, fmt.Sprintf(msg, a...))
	}
}

// Critical calls Log with CRITICAL level.
func (lg *logger) Critical(msg string, a ...any) {
	if lg.isEnabledFor(CRITICAL) {
		lg.log(CRITICAL, fmt.Sprintf(msg, a...))
	}
}

// Log logs a message with a custom level.
func (lg *logger) Log(level int, msg string, a ...any) {
	if lg.isEnabledFor(level) {
		lg.log(level, fmt.Sprintf(msg, a...))
	}
}

// log is a low-level logging method which creates a LogRecord and then calls
// all the handlers of this logger to handle the record.
func (lg *logger) log(level int, msg string) {
	pc, filename, lineno, ok := runtime.Caller(2)
	if !ok {
		filename = "unknown"
		lineno = -1
	}

	var record = MakeRecord(lg.fullname, level, filename, lineno, msg, pc)

	lg.handle(record)
}

// handle calls the handlers for the specified record.
func (lg *logger) handle(record LogRecord) {
	if lg.filter(record) {
		lg.callHandlers(record)
	}
}

// callHandlers passes a record to all relevant handlers.
//
// Loop through all handlers for this logger and its parents in the logger
// hierarchy. If no handler was found, output a one-off error message to
// os.Stderr.
func (lg *logger) callHandlers(record LogRecord) {
	var c = lg
	var found = 0
	for c != nil {
		for i := range c.handlers {
			lg.handlers[i].handle(record)
			found += 1
		}
		c = c.parent
	}

	if found == 0 {
		lastHandler.handle(record)
	}
}

// isEnabledFor checks if a logging level should be logged in this logger.
func (lg *logger) isEnabledFor(level int) bool {
	var isEnabled, isCached bool
	var _ = lg.lock.RLockFunc(func() any {
		isEnabled, isCached = lg.cache[level]
		return nil
	})

	if !isCached {
		isEnabled = level >= lg.getEffectiveLevel()
		lg.lock.WLockFunc(func() { lg.cache[level] = isEnabled })
	}
	return isEnabled
}

// getEffectiveLevel gets the effective level for this logger.
//
// Loop through this logger and its parents in the logger hierarchy,
// looking for a non-zero logging level. Return the first one found.
func (lg *logger) getEffectiveLevel() int {
	var level = lg.lock.RLockFunc(func() any { return lg.level }).(int)
	if level == NOTSET && lg.parent != nil {
		return lg.parent.getEffectiveLevel()
	}
	return level
}

// clearCache clears logging level cache of this logger and all its children.
func (lg *logger) clearCache() {
	lg.lock.WLockFunc(func() {
		for k := range lg.cache {
			delete(lg.cache, k)
		}
	})
	for i := range lg.children {
		lg.children[i].clearCache()
	}
}
