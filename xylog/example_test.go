package xylog_test

import (
	"os"

	"github.com/xybor/xyplatform/xylog"
)

func Example() {
	// You can directly use xylog functions to log with the root logger.
	var handler = xylog.NewStreamHandler()
	handler.SetStream(os.Stdout)

	xylog.SetLevel(xylog.DEBUG)
	xylog.AddHandler(handler)
	xylog.Debug("foo")

	// Handlers in the root logger will affect to other logger, so in this
	// example, it should remove this handler from the root logger.
	xylog.RemoveHandler(handler)

	// Output:
	// foo
}

func ExampleGetLogger() {
	var handler = xylog.NewStreamHandler()
	handler.SetStream(os.Stdout)
	handler.SetFormatter(xylog.NewTextFormatter(
		"module=%(name)s level=%(levelname)s %(message)s"))

	var logger = xylog.GetLogger("example")
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Debug("foo %s", "bar")

	// Output:
	// module=example level=DEBUG foo bar
}
