// Package xyplatform contains many utilities which make developer easy to code.
package xyplatform

import (
	"github.com/xybor/xyplatform/xylog"
)

func init() {
	var handler = xylog.NewHandler("xybor.xyplatform", xylog.StderrEmitter)
	handler.SetLevel(xylog.WARNING)
	handler.SetFormatter(xylog.NewTextFormatter(
		"time=%(asctime)-30s " +
			"level=%(levelname)-8s " +
			"%(message)s",
	))

	var logger = xylog.GetLogger("xybor.xyplatform")
	logger.SetLevel(xylog.WARNING)
	logger.AddHandler(handler)
}
