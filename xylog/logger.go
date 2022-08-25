package xylog

import (
	"fmt"
	"runtime"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xylock"
)

// Logger represents a single logging channel. A "logging channel" indicates an
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
type Logger struct {
	f *filterer

	fullname string
	children map[string]*Logger
	parent   *Logger
	level    int
	handlers []*Handler
	lock     xylock.RWLock
	cache    map[int]bool
	extra    string
}

// newlogger creates a new logger with a name and parent. The fullname of logger
// will be concatenated by the parent's fullname. This logger will not be
// automatically added to logger hierarchy. The returned logger has no child,
// no handler, and NOTSET level.
func newlogger(name string, parent *Logger) *Logger {
	var c = parent
	if c != nil && c != rootLogger {
		name = c.fullname + "." + name
	}

	return &Logger{
		f:        newfilterer(),
		fullname: name,
		children: make(map[string]*Logger),
		parent:   parent,
		level:    NOTSET,
		handlers: nil,
		lock:     xylock.RWLock{},
		cache:    make(map[int]bool),
		extra:    "",
	}
}

// SetLevel sets the new logging level. It also clears logging level cache of
// all loggers in program.
func (lg *Logger) SetLevel(level int) {
	lg.lock.WLockFunc(func() { lg.level = checkLevel(level) })
	rootLogger.clearCache()
}

// AddHandler adds a new handler.
func (lg *Logger) AddHandler(h *Handler) {
	xycond.MustNotNil(h).Assert("expected a not-nil Handler")
	lg.lock.WLockFunc(func() { lg.handlers = append(lg.handlers, h) })
}

// RemoveHandler removes an existed handler.
func (lg *Logger) RemoveHandler(h *Handler) {
	lg.lock.WLockFunc(func() {
		for i := range lg.handlers {
			if lg.handlers[i] == h {
				lg.handlers = append(lg.handlers[:i], lg.handlers[i+1:]...)
				break
			}
		}
	})
}

// GetHandlers returns all handlers of logger.
func (lg *Logger) GetHandlers() []*Handler {
	return lg.lock.RLockFunc(func() any { return lg.handlers }).([]*Handler)
}

// AddFilter adds a specified filter.
func (lg *Logger) AddFilter(f Filter) {
	lg.f.AddFilter(f)
}

// RemoveFilter removes an existed filter.
func (lg *Logger) RemoveFilter(f Filter) {
	lg.f.RemoveFilter(f)
}

// GetFilters returns all filters of filterer.
func (lg *Logger) GetFilters() []Filter {
	return lg.f.GetFilters()
}

func (lg *Logger) AddExtra(key string, value any) {
	var extra = key + "=" + fmt.Sprint(value)
	lg.extra = prefixMessage(lg.extra, extra)
}

// filter checks all filters in filterer, if there is any failed filter, it will
// returns false.
func (lg *Logger) filter(r LogRecord) bool {
	return lg.f.filter(r)
}

// Debug calls Log with DEBUG level.
func (lg *Logger) Debug(s string, a ...any) {
	if lg.isEnabledFor(DEBUG) {
		lg.log(DEBUG, fmt.Sprintf(s, a...))
	}
}

// Info calls Log with INFO level.
func (lg *Logger) Info(s string, a ...any) {
	if lg.isEnabledFor(INFO) {
		lg.log(INFO, s, a...)
	}
}

// Warn calls Log with WARN level.
func (lg *Logger) Warn(s string, a ...any) {
	if lg.isEnabledFor(WARN) {
		lg.log(WARN, s, a...)
	}
}

// Warning calls Log with WARNING level.
func (lg *Logger) Warning(s string, a ...any) {
	if lg.isEnabledFor(WARNING) {
		lg.log(WARNING, s, a...)
	}
}

// Error calls Log with ERROR level.
func (lg *Logger) Error(s string, a ...any) {
	if lg.isEnabledFor(ERROR) {
		lg.log(ERROR, s, a...)
	}
}

// Fatal calls Log with FATAL level.
func (lg *Logger) Fatal(s string, a ...any) {
	if lg.isEnabledFor(FATAL) {
		lg.log(FATAL, s, a...)
	}
}

// Critical calls Log with CRITICAL level.
func (lg *Logger) Critical(s string, a ...any) {
	if lg.isEnabledFor(CRITICAL) {
		lg.log(CRITICAL, s, a...)
	}
}

// Log logs a message with a custom level.
func (lg *Logger) Log(level int, s string, a ...any) {
	level = checkLevel(level)
	if lg.isEnabledFor(level) {
		lg.log(level, s, a...)
	}
}

// Event creates an eventLogger which logs key-value pairs.
func (lg *Logger) Event(e string) *EventLogger {
	return &EventLogger{
		lg:     lg,
		fields: map[string]any{"event": e},
	}
}

// log is a low-level logging method which creates a LogRecord and then calls
// all the handlers of this logger to handle the record.
func (lg *Logger) log(level int, s string, a ...any) {
	var msg = prefixMessage(lg.extra, fmt.Sprintf(s, a...))
	var pc, filename, lineno, ok = runtime.Caller(skipCall)
	if !ok {
		filename = "unknown"
		lineno = -1
	}

	var record = makeRecord(lg.fullname, level, filename, lineno, msg, pc)

	lg.handle(record)
}

// handle calls the handlers for the specified record.
func (lg *Logger) handle(record LogRecord) {
	if lg.filter(record) {
		lg.callHandlers(record)
	}
}

// callHandlers passes a record to all relevant handlers.
//
// Loop through all handlers for this logger and its parents in the logger
// hierarchy. If no handler was found, output a one-off error message to
// os.Stderr.
func (lg *Logger) callHandlers(record LogRecord) {
	var c = lg
	var found = 0
	for c != nil {
		for i := range c.handlers {
			c.handlers[i].handle(record)
			found += 1
		}
		c = c.parent
	}

	if found == 0 {
		lastHandler.handle(record)
	}
}

// isEnabledFor checks if a logging level should be logged in this logger.
func (lg *Logger) isEnabledFor(level int) bool {
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
func (lg *Logger) getEffectiveLevel() int {
	var level = lg.lock.RLockFunc(func() any { return lg.level }).(int)
	if level == NOTSET && lg.parent != nil {
		return lg.parent.getEffectiveLevel()
	}
	return level
}

// clearCache clears logging level cache of this logger and all its children.
func (lg *Logger) clearCache() {
	lg.lock.WLockFunc(func() {
		for k := range lg.cache {
			delete(lg.cache, k)
		}
	})
	for i := range lg.children {
		lg.children[i].clearCache()
	}
}

// prefixMessage adds a prefix to origin message if the prefix is not empty.
func prefixMessage(prefix, msg string) string {
	if prefix != "" {
		msg = fmt.Sprintf("%s %s", prefix, msg)
	}
	return msg
}
