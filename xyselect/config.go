package xyselect

import (
	_ "github.com/xybor/xyplatform" // This import will init xyplatform logger.
	"github.com/xybor/xyplatform/xyerror"
	"github.com/xybor/xyplatform/xylog"
)

var egen = xyerror.Register("xyselect", 200000)

// Errors of package xyselect.
var (
	SelectorError      = egen.NewClass("SelectorError")
	ClosedChannelError = SelectorError.NewClass("ClosedChannelError")
	ExhaustedError     = SelectorError.NewClass("ExhaustedError")
)

var logger = xylog.GetLogger("xybor.xyplatform.xyselect")

func init() {
	logger.AddExtra("module", "xyselect")
}
