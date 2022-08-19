package xylog_test

import (
	"fmt"
	"os"

	"github.com/xybor/xyplatform/xylog"
)

// NOTE: In example_test.go, xylog.StdoutEmitter is not accepted as the compared
// output. For this reason, in all examples, we must create a new one.
// In reality, you should use xylog.StdoutEmitter or xylog.StderrEmitter
// instead.

func Example() {
	// You can directly use xylog functions to log with the root logger.
	var handler = xylog.NewHandler("", xylog.NewStreamEmitter(os.Stdout))

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
	var handler = xylog.NewHandler("", xylog.NewStreamEmitter(os.Stdout))
	handler.SetFormatter(xylog.NewTextFormatter(
		"module=%(name)s level=%(levelname)s %(message)s"))

	var logger = xylog.GetLogger("example")
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Debug("foo %s", "bar")

	// Output:
	// module=example level=DEBUG foo bar
}

func ExampleHandler() {
	// You can use a handler throughout program without storing it in global
	// scope. All handlers can be identified by their names.
	var handlerA = xylog.NewHandler("example", xylog.StdoutEmitter)
	var handlerB = xylog.GetHandler("example")
	if handlerA == handlerB {
		fmt.Println("handlerA == handlerB")
	} else {
		fmt.Println("handlerA != handlerB")
	}

	// In case name is an empty string, it totally is a fresh handler.
	var handlerC = xylog.NewHandler("", xylog.StdoutEmitter)
	var handlerD = xylog.NewHandler("", xylog.StdoutEmitter)
	if handlerC == handlerD {
		fmt.Println("handlerC == handlerD")
	} else {
		fmt.Println("handlerC != handlerD")
	}

	// Output:
	// handlerA == handlerB
	// handlerC != handlerD
}
