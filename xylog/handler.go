package xylog

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xylock"
)

// Emitter instances dispatch logging events to specific destinations.
type Emitter interface {
	Emit(string)
}

// Handler handles logging events. Do NOT instantiated directly this struct.
//
// Any Handler with a not-empty name will be associated with its name.
type Handler struct {
	f *filterer
	e Emitter

	level     int
	formatter Formatter
	lock      xylock.RWLock
}

// NewHandler creates a Handler with a specified Emitter.
//
// Any Handler with a not-empty name will be associated with its name. Calling
// NewHandler twice with the same name will cause a panic. If you want to create
// an anonymous Handler, call this function with an empty name.
func NewHandler(name string, e Emitter) *Handler {
	var handler = GetHandler(name)
	xycond.MustNil(handler).Assert(
		"The handler with name %s is associated with another Emitter", name)

	handler = &Handler{
		f:         newfilterer(),
		e:         e,
		level:     NOTSET,
		formatter: defaultFormatter,
		lock:      xylock.RWLock{},
	}

	if name != "" {
		mapHandler(name, handler)
	}

	return handler
}

// SetLevel sets the new logging level of handler. It is NOTSET by default.
func (h *Handler) SetLevel(level int) {
	h.lock.WLockFunc(func() { h.level = checkLevel(level) })
}

// SetFormatter sets the new formatter of handler.
func (h *Handler) SetFormatter(f Formatter) {
	h.lock.WLockFunc(func() { h.formatter = f })
}

// AddFilter adds a specified filter.
func (h *Handler) AddFilter(f Filter) {
	h.f.AddFilter(f)
}

// RemoveFilter removes an existed filter.
func (h *Handler) RemoveFilter(f Filter) {
	h.f.RemoveFilter(f)
}

// filter checks all filters in filterer, if there is any failed filter, it will
// returns false.
func (h *Handler) filter(r LogRecord) bool {
	return h.f.filter(r)
}

// format uses formatter to format the record.
func (h *Handler) format(record LogRecord) string {
	var f = h.lock.RLockFunc(func() any { return h.formatter }).(Formatter)
	return f.Format(record)
}

// handle handles a new record, it will check if the record should be logged or
// not, then call emit if it is.
func (h *Handler) handle(record LogRecord) {
	var level = h.lock.RLockFunc(func() any { return h.level }).(int)
	if h.filter(record) && record.LevelNo >= level {
		h.e.Emit(h.format(record))
	}
}

// StreamEmitter writes logging message to a stream. Note that this class does
// not close the stream, as os.Stdout or os.Stderr may be used.
type StreamEmitter struct {
	stream *bufio.Writer
}

// NewStreamEmitter creates a StreamEmitter which writes message to a stream
// (os.Stderr by default).
func NewStreamEmitter(f *os.File) *StreamEmitter {
	var stream = bufio.NewWriter(f)
	stream.Flush()
	return &StreamEmitter{stream: stream}
}

// Emit will be called after a record was decided to log.
func (e *StreamEmitter) Emit(msg string) {
	var _, err = e.stream.WriteString(msg + "\n")
	if err == nil {
		err = e.stream.Flush()
	}

	if err != nil {
		os.Stderr.Write([]byte("------------ Logging error ------------\n"))
		os.Stderr.Write([]byte(
			fmt.Sprintf("An error occurs when logging: %s\n", err)))
		debug.PrintStack()
	}
}
