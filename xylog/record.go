// This file copied and modified comments of python logging.
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
// passed in is in Message. The record also includes information as when the
// record was created or the source line where the logging call was made.
//
// The current attributes are described by:
//
// %(asctime)s         Textual time when the LogRecord was created
// %(created)f         Time when the LogRecord was created (time.Now().Unix()
//                     return value)
// %(filename)s        Filename portion of pathname
// %(funcName)s        Function name logged the record
// %(levelname)s       Text logging level for the message ("DEBUG", "INFO",
//                     "WARNING", "ERROR", "CRITICAL")
// %(levelno)s         Numeric logging level for the message (DEBUG, INFO,
//                     WARNING, ERROR, CRITICAL)
// %(lineno)d          Source line number where the logging call was issued
//                     (if available)
// %(message)s         The logging message
// %(module)s          Module (name portion of filename)
// %(msecs)d           Millisecond portion of the creation time
// %(name)s            Name of the logger
// %(pathname)s        Full pathname of the source file where the logging
//                     call was issued (if available)
// %(process)d         Process ID (if available)
// %(relativeCreated)d Time in milliseconds when the LogRecord was created,
//                     relative to the time the logging module was loaded
//                     (typically at application startup time)
type LogRecord struct {
	Asctime         string `map:"asctime"`
	Created         int64  `map:"created"`
	FileName        string `map:"filename"`
	FuncName        string `map:"funcname"`
	LevelName       string `map:"levelname"`
	LevelNo         int    `map:"levelno"`
	LineNo          int    `map:"lineno"`
	Message         string `map:"message"`
	Module          string `map:"module"`
	Msecs           int    `map:"msecs"`
	Name            string `map:"name"`
	PathName        string `map:"pathname"`
	Process         int    `map:"process"`
	RelativeCreated int64  `map:"relativeCreated"`
}

func MakeRecord(
	name string, level int, pathname string, lineno int, msg string,
	pc uintptr,
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
		RelativeCreated: created.Unix() - startTime,
	}
}

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
