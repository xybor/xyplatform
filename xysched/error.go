package xysched

import (
	"github.com/xybor/xyplatform"
	"github.com/xybor/xyplatform/xyerror"
)

var _ = xyerror.Register(xyplatform.XySched)

var (
	CallError      = xyerror.NewClass(xyplatform.XySched, "CallError")
	ParameterError = xyerror.ParameterError.NewClassM(xyplatform.XySched)
)
