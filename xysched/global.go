// Package xysched supports to schedule future tasks with a simple syntax.
package xysched

import (
	"time"

	_ "github.com/xybor/xyplatform" // This import will init xyplatform logger.
	"github.com/xybor/xyplatform/xylock"
	"github.com/xybor/xyplatform/xylog"
)

func init() {
	logger.AddExtra("module", "xysched")
}

var logger = xylog.GetLogger("xybor.xyplatform.xysched")

// schedulerManager stores Schedulers with their names.
var schedulerManager = map[string]*Scheduler{}
var anonSchedCounter = 0
var futureCounter = 0
var lock = xylock.RWLock{}

// A default scheduler.
var global = NewScheduler("")

// After creates a send-only channel. Sending a future to this channel will
// add it to scheduler after a duration. If d is negative, After will send the
// future to scheduler immediately.
//
// NOTE: You should send ONLY ONE future to this channel because it is designed
// to handle one. If you try sending another, it will be blocked forever. To
// send other futures to scheduler, let call this method again.
func After(d time.Duration) chan<- future {
	return global.After(d)
}

// At is a shortcut of After(time.Until(next)).
//
// NOTE: You should send ONLY ONE future to this channel because it is designed
// to handle one. If you try sending another, it will be blocked forever. To
// send other futures to scheduler, let call this method again.
func At(next time.Time) chan<- future {
	return global.At(next)
}

// Now is a shortcut of After(0).
//
// NOTE: You should send ONLY ONE future to this channel because it is designed
// to handle one. If you try sending another, it will be blocked forever. To
// send other futures to scheduler, let call this method again.
func Now() chan<- future {
	return global.Now()
}
