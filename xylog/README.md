# Introduction
Xylog provides flexible logging methods to the program.

This package is inspired by idea of Python `logging`.

# Feature
The basic structs defined by the module, together with their functions, are
listed below.
1. Loggers expose the interface that application code directly uses.
2. Handlers send the log records (created by loggers) to the appropriate
destination.
3. Filters provide a finer grained facility for determining which log records to
output.
4. Formatters specify the layout of log records in the final output.

Visit [pkg.go.dev](https://pkg.go.dev/github.com/xybor/xyplatform/xylog) for
more details.

## Logger
Loggers should NEVER be instantiated directly, but always through the
module-level function `xylog.GetLogger(name)`. Multiple calls to `GetLogger()`
with the same name will always return a reference to the same Logger object.

You can logs a message by using one of the following methods:
```golang
func (*Logger) Critical(msg string, a ...any)
func (*Logger) Error(msg string, a ...any)
func (*Logger) Fatal(msg string, a ...any)
func (*Logger) Warn(msg string, a ...any)
func (*Logger) Warning(msg string, a ...any)
func (*Logger) Info(msg string, a ...any)
func (*Logger) Debug(msg string, a ...any)
func (*Logger) Log(level int, msg string, a ...any)
```

To add `Handler` or `Filter` instances to the `logger`, call `AddHandler` or
`AddFilter` methods.

To adjust the level, using `SetLevel` method.

## Logging level
The numeric values of logging levels are given in the following table. These are
primarily of interest if you want to define your own levels, and need them to
have specific values relative to the predefined levels. If you define a level
with the same numeric value, it overwrites the predefined value; the predefined
name is lost.
|      Level     | Numeric value  |
|----------------|----------------|
| CRITICAL       |              50|
| ERROR/FATAL    |              40|
| WARN/WARNING   |              30|
| INFO           |              20|
| DEBUG          |              10|
| NOTSET         |               0|

## Handler
`Handler` instances handle logging events.

Any `Handler` created with a not-empty name will be associated with this name.
Later calls of this function with the same name will return the same `Handler`.
If a name is associated with a `Handler` type, do not reuse this name for other
types.

`Handler` can use `SetFormatter` method to format the logging message.

Each `Handler` has its own way and place to log the message. Such as
`SteamHandler` is used to log into stdout or stderr.

## Formatter
`Formatter` instances are used to convert a `LogRecord` to text.

`Formatter` need to know how a `LogRecord` is constructed. They are responsible
for converting a `LogRecord` to a string which can be interpreted by either a
human or an external system.

`TextFormatter` is a built-in `Formatter` which uses logging macros to format
the message.

| MACROS            |  DESCRIPTION                                    |
|-------------------|-------------------------------------------------|
|`asctime`          |Textual time when the LogRecord was created.|
|`created`          |Time when the LogRecord was created (time.Now().Unix() return value).|
|`filename`         |Filename portion of pathname.|
|`funcname`         |Function name logged the record.|
|`levelname`        |Text logging level for the message ("DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL").|
|`levelno`          |Numeric logging level for the message (DEBUG, INFO, WARNING, ERROR, CRITICAL).|
|`lineno`           |Source line number where the logging call was issued.|
|`message`          |The logging message.|
|`module`           |The module called log method.|
|`msecs`            |Millisecond portion of the creation time.|
|`name`             |Name of the logger.|
|`pathname`         |Full pathname of the source file where the logging call was issued.|
|`process`          |Process ID.|
|`relativeCreated`  |Time in milliseconds when the LogRecord was created, relative to the time the logging module was loaded (typically at application startup time).|

## Filter
`Filter` instances are used to perform arbitrary filtering of `LogRecord`.

A `Filter` struct needs to define `Format(LogRecord)` method, which return true
if it allows to log the `LogRecord`, and vice versa.

`Filter` can be used in both `Handler` and `Logger`.

# Example
## Simple usage
```golang
var handler = xylog.StreamHandler("xybor")
handler.SetFormatter(xylog.NewTextFormmater("%(level)s %(message)s"))

var logger = xylog.GetLogger("xybor.service")
logger.AddHandler(handler)
logger.SetLevel(xylog.DEBUG)

logger.Debug("foo")

// Output:
// DEBUG foo
```

## Filter definition
```golang
// LoggerNameFilter only logs out records belongs to a specified logger.
type LoggerNameFilter struct {
    name string
}

func (f *LoggerNameFilter) Filter(r xylog.LogRecord) bool {
    return f.name == r.name
}

// Get the logger of above example
var logger = xylog.GetLogger("xybor.service")
logger.AddFilter(&LoggerNameFilter{"xybor.service.chat"})

logger.Debug("foo")
xylog.GetLogger("xybor.service.auth").Debug("auth foo")
xylog.GetLogger("xybor.service.chat").Debug("chat foo")

// Output:
// chat foo
```

## Root logger
```golang
// A simple program with only one application area could use directly the root
// logger.
var handler = xylog.StreamHandler("")
xylog.SetLevel(xylog.DEBUG)
xylog.AddHandler(handler)

xylog.Debug("bar")

// Output:
// bar
```

## Xyplatform log
```golang
// example.go
package foo

// To apply xyplatform standard logging, you must import xyplatform.
import (
    _ "github.com/xybor/xyplatform"
    "github.com/xybor/xyplatform/xylog"
)

// All loggers need to be started with the prefix "xybor.xyplatform."
var logger = xylog.GetLogger("xybor.xyplatform.foo")

func main() {
    logger.Debug("message=bar")
}

// Output:
// time=[time] source=example.go.main:13 level=DEBUG module=foo message=bar
```
