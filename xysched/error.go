package xysched

import (
	"github.com/xybor/xyplatform/xyerror"
)

var eid = xyerror.Register("xysched", 300000)

var (
	CallError      = eid.NewClass("CallError")
	ParameterError = xyerror.ParameterError.NewClassM(eid)
)
