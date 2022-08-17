package xyplatform

import (
	"os"

	"github.com/xybor/xyplatform/xylog"
)

func init() {
	var handler = xylog.NewStreamHandler("xyplatform")
	handler.SetStream(os.Stdout)
	handler.SetLevel(xylog.WARNING)
	handler.SetFormatter(xylog.NewTextFormatter(
		"time=%(asctime)s+%(msecs)d " +
			"source=%(filename)s.%(funcname)s:%(lineno)d " +
			"level=%(levelname)s " +
			"module=%(module)s " +
			"%(message)s",
	))

	var logger = xylog.GetLogger("xybor.xyplatform")
	logger.SetLevel(xylog.WARNING)
	logger.AddHandler(handler)
}
