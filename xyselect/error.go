package xyselect

import (
	"github.com/xybor/xyplatform"
	"github.com/xybor/xyplatform/xyerror"
)

var _ = xyerror.Register(xyplatform.XySelector)

var (
	SelectorError = xyerror.NewClass(xyplatform.XySelector, "SelectorError")

	ClosedChannelError = SelectorError.NewClass("ClosedChannelError")
	DefaultCaseError   = SelectorError.NewClass("DefaultCaseError")
	ExhaustedError     = SelectorError.NewClass("ExhaustedError")
)
