package xyselect

import (
	"github.com/xybor/xyplatform/xyerror"
)

var egen = xyerror.Register("xyselect", 200000)

var (
	SelectorError      = egen.NewClass("SelectorError")
	ClosedChannelError = SelectorError.NewClass("ClosedChannelError")
	DefaultCaseError   = SelectorError.NewClass("DefaultCaseError")
	ExhaustedError     = SelectorError.NewClass("ExhaustedError")
)
