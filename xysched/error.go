package xysched

import (
	"github.com/xybor/xyplatform/xyerror"
)

var egen = xyerror.Register("xysched", 300000)

// Errors of package xysched.
var (
	CallError = egen.NewClass("CallError")
)
