package xyselect

import (
	"github.com/xybor/xyplatform/xyerror"
)

var egen = xyerror.Register("xyselect", 200000)

// Errors of package xyselect.
var (
	SelectorError      = egen.NewClass("SelectorError")
	ClosedChannelError = SelectorError.NewClass("ClosedChannelError")
	ExhaustedError     = SelectorError.NewClass("ExhaustedError")
)
