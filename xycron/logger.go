package xycron

import (
	"github.com/xybor/xyplatform"
	"github.com/xybor/xyplatform/xylog"
)

var _ = xylog.Register(xyplatform.XyCron)
var _ = xylog.Config(xyplatform.XyCron, xylog.Allow(xylog.WARN))
var logger = xylog.Logger(xyplatform.XyCron)
