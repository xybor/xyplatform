package xyplatform

import (
	"os"

	"github.com/xybor/xyplatform/xylog"
)

func init() {
	var handler = xylog.StreamHandler()
	handler.SetStream(os.Stdout)
	handler.SetLevel(xylog.WARNING)
	handler.SetFormatter(
		xylog.Formatter(
			"source=%(filename)s.%(funcname)s:%(lineno)d time=%(asctime)s+%(msecs)d " +
				"level=%(levelname)s %(message)s",
		))

	var logger = xylog.GetLogger("xyplatform")
	logger.SetLevel(xylog.WARNING)
	logger.AddHandler(handler)
}
