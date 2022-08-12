package xylog_test

import (
	"os"

	"github.com/xybor/xyplatform/xylog"
)

func Example() {
	var handler = xylog.StreamHandler()
	handler.SetStream(os.Stdout)
	handler.SetFormatter(xylog.Formatter("module=%(name)s level=%(levelname)s %(message)s"))

	var logger = xylog.GetLogger("example")
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	logger.Debug("foo %s", "bar")

	// Output:
	// module=example level=DEBUG foo bar
}
