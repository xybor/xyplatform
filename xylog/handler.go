package xylog

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/xybor/xyplatform/xylock"
)

// handler instances handle logging events.
type handler interface {
	handle(LogRecord)
}

// emiter instances dispatch logging events to specific destinations.
type emiter interface {
	emit(LogRecord)
}

// BaseHandler acts as a placeholder which defines the handler interface.
// Handlers can optionally use Formatter instances to format records as desired.
// By default, no formatter is specified; in this case, the 'raw' message as
// determined by record.message is logged.
type BaseHandler struct {
	filterer

	base      emiter
	level     int
	formatter Formatter
	lock      xylock.RWLock
}

// newBaseHandler creates the baseHandler with a specified emiter as the base
// object of baseHandler. This function is useful in create a new concrete
// handler (rather than baseHandler).
func newBaseHandler(e emiter) *BaseHandler {
	return &BaseHandler{
		filterer:  newfilterer(),
		base:      e,
		level:     NOTSET,
		formatter: defaultFormatter,
		lock:      xylock.RWLock{},
	}
}

// SetLevel sets the new logging level of handler. It is NOTSET by default.
func (h *BaseHandler) SetLevel(level int) {
	h.lock.WLockFunc(func() { h.level = checkLevel(level) })
}

// SetFormatter sets the new formatter of handler. It is defaultFormatter by
// default.
func (h *BaseHandler) SetFormatter(f Formatter) {
	h.lock.WLockFunc(func() { h.formatter = f })
}

// format uses formatter to format the record.
func (h *BaseHandler) format(record LogRecord) string {
	return h.formatter.Format(record)
}

// handle handles a new record, it will check if the record should be logged or
// not, then call emit if it is.
func (h *BaseHandler) handle(record LogRecord) {
	if h.filter(record) && record.LevelNo >= h.level {
		h.lock.WLockFunc(func() { h.base.emit(record) })
	}
}

// StreamHandler writes logging records, appropriately formatted, to a stream.
// Note that this class does not close the stream, as os.Stdout or os.Stderr may
// be used.
type StreamHandler struct {
	BaseHandler
	stream *bufio.Writer
}

// NewStreamHandler returns a StreamHandler, the handler writes records to a
// stream (os.Stderr by default).
func NewStreamHandler() *StreamHandler {
	var hdr = &StreamHandler{
		stream: bufio.NewWriter(os.Stderr),
	}
	hdr.BaseHandler = *newBaseHandler(hdr)

	return hdr
}

// SetStream sets a new stream for this handler. Note that this stream will not
// be closed, so it may use os.Stderr or os.Stdout.
func (hdr *StreamHandler) SetStream(f *os.File) {
	var stream = bufio.NewWriter(f)
	stream.Flush()
	hdr.lock.WLockFunc(func() { hdr.stream = stream })
}

// emit will be called after a record was decided to log.
func (hdr *StreamHandler) emit(record LogRecord) {
	var msg = hdr.format(record)
	if msg == "" {
		return
	}

	var err error
	_, err = hdr.stream.WriteString(hdr.format(record) + "\n")
	if err == nil {
		err = hdr.stream.Flush()
	}

	if err != nil {
		os.Stderr.Write([]byte("------------ Logging error ------------\n"))
		os.Stderr.Write([]byte(
			fmt.Sprintf("An error occurs when logging: %s\n", err)))
		debug.PrintStack()
	}
}
