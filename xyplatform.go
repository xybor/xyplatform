package xyplatform

import (
	"github.com/xybor/xyplatform/xylog"
)

func init() {
	var handler = xylog.NewHandler("xybor.xyplatform", xylog.StderrEmitter)
	handler.SetLevel(xylog.WARNING)
	handler.SetFormatter(xylog.NewTextFormatter(
		"time=%(asctime)s " +
			"source=%(filename)s.%(funcname)s:%(lineno)d " +
			"level=%(levelname)s " +
			"module=%(module)s " +
			"%(message)s",
	))

	var logger = xylog.GetLogger("xybor.xyplatform")
	logger.SetLevel(xylog.WARNING)
	logger.AddHandler(handler)
}
