package xysched

import "sync"

var global *scheduler = nil
var globalOnce sync.Once

// A scheduler you could use throughout the program without creating a new one.
func Global() *scheduler {
	globalOnce.Do(func() { global = New() })
	return global
}
