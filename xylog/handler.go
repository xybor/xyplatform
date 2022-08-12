// This file copied and modified comments of python logging.
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

// baseHandler acts as a placeholder which defines the handler interface.
// Handlers can optionally use Formatter instances to format records as desired.
// By default, no formatter is specified; in this case, the 'raw' message as
// determined by record.message is logged.
type baseHandler struct {
	filterer

	base      emiter
	level     int
	formatter formatter
	lock      xylock.RWLock
}

// newBaseHandler creates the baseHandler with a specified emiter as the base
// object of baseHandler. This function is useful in create a new concrete
// handler (rather than baseHandler).
func newBaseHandler(e emiter) *baseHandler {
	return &baseHandler{
		filterer:  newfilterer(),
		base:      e,
		level:     NOTSET,
		formatter: defaultFormatter,
		lock:      xylock.RWLock{},
	}
}

// SetLevel sets the new logging level of handler. It is NOTSET by default.
func (h *baseHandler) SetLevel(level int) {
	h.lock.WLockFunc(func() { h.level = checkLevel(level) })
}

// SetFormatter sets the new formatter of handler. It is defaultFormatter by
// default.
func (h *baseHandler) SetFormatter(f formatter) {
	h.lock.WLockFunc(func() { h.formatter = f })
}

// format uses formatter to format the record.
func (h *baseHandler) format(record LogRecord) string {
	return h.formatter.Format(record)
}

// handle handles a new record, it will check if the record should be logged or
// not, then call emit if it is.
func (h *baseHandler) handle(record LogRecord) {
	if h.filter(record) && record.LevelNo >= h.level {
		h.lock.WLockFunc(func() { h.base.emit(record) })
	}
}

// streamHandler writes logging records, appropriately formatted, to a stream.
// Note that this class does not close the stream, as os.stdout or os.stderr
// may be used.
type streamHandler struct {
	*baseHandler
	stream *bufio.Writer
}

// StreamHandler returns a streamHandler, the handler writes records to a stream
// (os.Stderr by default).
func StreamHandler() *streamHandler {
	var hdr = &streamHandler{
		stream: bufio.NewWriter(os.Stderr),
	}
	hdr.baseHandler = newBaseHandler(hdr)

	return hdr
}

// SetStream sets a new stream for this handler. Note that this stream will not
// be closed, so it may use os.Stderr or os.Stdout.
func (hdr *streamHandler) SetStream(f *os.File) {
	var stream = bufio.NewWriter(f)
	stream.Flush()
	hdr.lock.WLockFunc(func() { hdr.stream = stream })
}

// emit will be called after a record was decided to log.
func (hdr *streamHandler) emit(record LogRecord) {
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
