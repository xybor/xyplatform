package xylog

import (
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// A LogRecord instance represents an event being logged.
//
// LogRecord instances are created every time something is logged. They contain
// all the information pertinent to the event being logged. The main information
// passed in is Message. The record also includes information as when the record
// was created or the source line where the logging call was made.
type LogRecord struct {
	// Textual time when the LogRecord was created.
	Asctime string `map:"asctime"`

	//Time when the LogRecord was created (time.Now().Unix() return value).
	Created int64 `map:"created"`

	// Filename portion of pathname.
	FileName string `map:"filename"`

	// Function name logged the record.
	FuncName string `map:"funcname"`

	// Text logging level for the message ("DEBUG", "INFO", "WARNING", "ERROR",
	// "CRITICAL").
	LevelName string `map:"levelname"`

	// Numeric logging level for the message (DEBUG, INFO, WARNING, ERROR,
	// CRITICAL).
	LevelNo int `map:"levelno"`

	// Source line number where the logging call was issued.
	LineNo int `map:"lineno"`

	// The logging message.
	Message string `map:"message"`

	// The module called log method.
	Module string `map:"module"`

	// Millisecond portion of the creation time.
	Msecs int `map:"msecs"`

	// Name of the logger.
	Name string `map:"name"`

	// Full pathname of the source file where the logging call was issued.
	PathName string `map:"pathname"`

	// Process ID.
	Process int `map:"process"`

	// Time in milliseconds when the LogRecord was created, relative to the time
	// the logging module was loaded (typically at application startup time).
	RelativeCreated int64 `map:"relativeCreated"`
}

// makeRecord creates specialized LogRecords.
func makeRecord(
	name string, level int, pathname string, lineno int, msg string, pc uintptr,
) LogRecord {
	var created = time.Now()
	var module, funcname = extractFromPC(pc)

	return LogRecord{
		Asctime:         created.Format(timeLayout),
		Created:         created.Unix(),
		FileName:        filepath.Base(pathname),
		FuncName:        funcname,
		LevelName:       getLevelName(level),
		LevelNo:         level,
		LineNo:          lineno,
		Message:         msg,
		Module:          module,
		Msecs:           created.Nanosecond() / int(time.Millisecond),
		Name:            name,
		PathName:        pathname,
		Process:         processid,
		RelativeCreated: created.UnixMilli() - startTime,
	}
}

// extractFromPC returns module name and function name from program counter.
func extractFromPC(pc uintptr) (module, fname string) {
	var s = runtime.FuncForPC(pc).Name()

	// Split the funcname in the form of func with receiver.
	// E.g. module.(receiver).func
	var parts []string
	var sep = ".("
	parts = strings.Split(s, ".(")

	// If it is not the form of func with receiver, split it with normal func.
	// E.g. module.func
	if len(parts) <= 1 {
		sep = "."
		parts = strings.Split(s, ".")
	}

	// In case one of form is valid, remove the funcname from string.
	if len(parts) > 1 {
		var funcname = parts[len(parts)-1]
		module = strings.TrimSuffix(s, sep+funcname)
		fname = strings.TrimPrefix(sep+funcname, ".")
		return
	}

	// Otherwise, the string contains only funcname.
	return "unknown", s
}
