package xycron

var globalScheduler = New()

// A scheduler you could use throughout the program without creating a new one.
func Global() *scheduler {
	return globalScheduler
}
