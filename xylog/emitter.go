package xylog

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/xybor/xyplatform/xycond"
)

// LogWriter instances define a writer using to log.
type LogWriter interface {
	Write([]byte) (int, error)
	Close() error
}

// Emitter instances dispatch logging events to specific destinations.
type Emitter interface {
	// Emit will be called after a record was decided to log.
	Emit(LogRecord)

	// SetFormatter sets the new formatter to Emitter.
	SetFormatter(Formatter)
}

// StreamEmitter writes logging message to a stream.
type StreamEmitter struct {
	stream    *bufio.Writer
	formatter Formatter
}

// NewStreamEmitter creates a StreamEmitter which writes message to a stream
// (os.Stderr by default).
func NewStreamEmitter(w io.Writer) *StreamEmitter {
	var e = &StreamEmitter{formatter: defaultFormatter}
	e.setStream(w)
	return e
}

// Emit will be called after a record was decided to log.
func (e *StreamEmitter) Emit(record LogRecord) {
	var msg = e.formatter.Format(record)
	var _, err = e.stream.WriteString(msg + "\n")
	if err == nil {
		err = e.stream.Flush()
	}

	if err != nil {
		log.Println("------------ Logging error ------------")
		log.Printf("An error occurs when logging: %s\n", err)
		log.Panic(string(debug.Stack()))
	}
}

// SetFormatter sets the new formatter to Emitter.
func (e *StreamEmitter) SetFormatter(f Formatter) {
	e.formatter = f
}

// setStream sets a new stream to emitter.
func (e *StreamEmitter) setStream(w io.Writer) {
	if e.stream != nil {
		e.stream.Flush()
	}

	if w == nil {
		e.stream = nil
	} else {
		var stream = bufio.NewWriter(w)
		stream.Flush()
		e.stream = stream
	}
}

type FileEmitter struct {
	*StreamEmitter
	rotator     rotator
	filename    string
	writer      LogWriter
	backupCount uint
}

// NewFileEmitter creates a StreamEmitter by providing the file name.
func NewFileEmitter(fn string) *FileEmitter {
	var emitter = &FileEmitter{
		rotator: nil, backupCount: 0,
		filename: fn, writer: nil,
		StreamEmitter: NewStreamEmitter(nil),
	}
	return emitter
}

func NewSizeRotatingFileEmitter(
	fn string, maxBytes uint64, backupCount uint,
) *FileEmitter {
	var emitter = NewFileEmitter(fn)
	emitter.backupCount = backupCount
	emitter.rotator = &sizeRotator{filename: fn, maxBytes: maxBytes}

	return emitter
}

func NewTimeRotatingFileEmitter(
	fn string, interval time.Duration, backupCount uint,
) *FileEmitter {
	var emitter = NewFileEmitter(fn)
	emitter.backupCount = backupCount
	emitter.rotator = &timeRotator{
		d:            interval,
		nextRollover: time.Now().Add(interval),
	}

	return emitter
}

func (e *FileEmitter) Emit(record LogRecord) {
	if e.writer == nil {
		e.open()
	}
	if e.rotator != nil && e.rotator.shouldRollover() {
		e.doRollover()
	}
	e.StreamEmitter.Emit(record)
}

func (e *FileEmitter) open() {
	var f, err = os.OpenFile(e.filename, fileflag, fileperm)
	xycond.MustNil(err).Assert("%v", err)
	e.writer = f
	e.setStream(f)
}

// close stops to write to the log writer.
func (e *FileEmitter) close() {
	if e.writer != nil {
		var err = e.writer.Close()
		xycond.MustNil(err).Assert("%v", err)
		e.writer = nil
		e.setStream(nil)
	}
}

func (e *FileEmitter) doRollover() {
	e.close()

	for i := e.backupCount; i > 0; i-- {
		var sfn = rotationFilename(e.filename, i-1)
		var dfn = rotationFilename(e.filename, i)

		if _, err := os.Stat(sfn); err == nil {
			if _, err := os.Stat(dfn); err == nil {
				os.Remove(dfn)
			}
			os.Rename(sfn, dfn)
		}
	}

	e.open()
}

type rotator interface {
	shouldRollover() bool
}

type sizeRotator struct {
	filename string
	fd       *os.File
	maxBytes uint64
}

func (r *sizeRotator) shouldRollover() bool {
	var err error
	if r.fd == nil {
		r.fd, err = os.Open(r.filename)
		xycond.MustNil(err).Assert("%v", err)
	}

	stat, err := r.fd.Stat()
	xycond.MustNil(err).Assert("%v", err)

	if uint64(stat.Size()) >= r.maxBytes {
		r.fd = nil
		return true
	}
	return false
}

type timeRotator struct {
	d            time.Duration
	nextRollover time.Time
}

func (r *timeRotator) shouldRollover() bool {
	if time.Now().After(r.nextRollover) {
		r.nextRollover = time.Now().Add(r.d)
		return true
	}
	return false
}

func rotationFilename(base string, i uint) string {
	if i == 0 {
		return base
	}
	return fmt.Sprintf("%s.%d", base, i)
}
