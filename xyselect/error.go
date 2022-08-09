package xyselect

import (
	"github.com/xybor/xyplatform/xyerror"
)

var eid = xyerror.Register("xyselect", 200000)

var (
	SelectorError      = eid.NewClass("SelectorError")
	ClosedChannelError = SelectorError.NewClass("ClosedChannelError")
	DefaultCaseError   = SelectorError.NewClass("DefaultCaseError")
	ExhaustedError     = SelectorError.NewClass("ExhaustedError")
)
