package xyerror

import (
	"github.com/xybor/xyplatform"
	"github.com/xybor/xyplatform/xycond"
)

// manager is a map of module as key and errno as value.
var manager = make(map[xyplatform.Module]int)

// extractModule returns a Module given an errno.
func extractModule(errno int) xyplatform.Module {
	var module = xyplatform.NewModule(0, "")
	var minD int
	for m := range manager {
		d := errno - m.ID()
		if d < 0 || d > xyplatform.Default.ID() {
			continue
		}

		if module.ID() == 0 || d < minD {
			minD = d
			module = m
		}
	}

	xycond.NotZero(module.ID()).Assert("Cannot find the module of this error")

	return module
}

// nextErrno returns the next errno given a Module.
func nextErrno(m xyplatform.Module) int {
	xycond.Contains(manager, m).Assertf("Module %s had not registered yet", m)

	manager[m] += 1
	return m.ID() + manager[m]
}

// Register adds a Module to pool for managing new error types. The returned
// value of this function is meaningless, it only helps run it in the global
// scope.
func Register(m xyplatform.Module) bool {
	defaultID := xyplatform.Default.ID()
	xycond.Divisible(m.ID(), defaultID).
		Assertf("%s's ID %d is not divisible by %d", m.Name(), m.ID(), defaultID)

	xycond.NotContains(manager, m).Assertf("Module %s had already registered", m)

	manager[m] = 0

	return true
}
