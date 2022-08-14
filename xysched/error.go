package xysched

import (
	"github.com/xybor/xyplatform/xyerror"
)

var egen = xyerror.Register("xysched", 300000)

var (
	CallError      = egen.NewClass("CallError")
	ParameterError = xyerror.ParameterError.NewClassM(egen)
)
