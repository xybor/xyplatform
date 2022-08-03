package xysched

import (
	"github.com/xybor/xyplatform"
	"github.com/xybor/xyplatform/xylog"
)

var _ = xylog.Register(xyplatform.XySched)
var _ = xylog.Config(xyplatform.XySched, xylog.Allow(xylog.WARN))
var logger = xylog.Logger(xyplatform.XySched)
