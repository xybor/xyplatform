package xysched

import "sync"

var global *Scheduler = nil
var globalOnce sync.Once

// A scheduler you could use throughout the program without creating a new one.
func Global() *Scheduler {
	globalOnce.Do(func() { global = NewScheduler() })
	return global
}
