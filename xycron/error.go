package xycron

import (
	"github.com/xybor/xyplatform"
	"github.com/xybor/xyplatform/xyerror"
)

var _ = xyerror.Register(xyplatform.XyCron)

var (
	StopError       = xyerror.NewClass(xyplatform.XyCron, "StopError")
	InProgressError = xyerror.NewClass(xyplatform.XyCron, "InProgressError")
)
