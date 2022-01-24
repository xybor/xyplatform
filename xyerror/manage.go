package xyerror

import (
	"fmt"

	"github.com/xybor/xyplatform"
)

// manager is a map of module as key and errno as value.
var manager = make(map[xyplatform.Module]int)

// extractModule returns a Module given an errno.
func extractModule(errno int) xyplatform.Module {
	var module = xyplatform.NewModule(0, "Temporary")
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

	if module.ID() == 0 {
		panic("Cannot find a module of this error")
	}

	return module
}

// nextErrno returns the next errno given a Module.
func nextErrno(m xyplatform.Module) int {
	if _, ok := manager[m]; !ok {
		msg := fmt.Sprintf("Module %s had not registered yet", m)
		panic(msg)
	}

	manager[m] += 1
	return m.ID() + manager[m]
}

// Register adds a Module to pool for managing new error types. The returned
// value of this function is meaningless, it only helps run it in the global
// scope.
func Register(m xyplatform.Module) bool {
	if m.ID()%xyplatform.Default.ID() != 0 {
		msg := fmt.Sprintf("%s's ID %d is not divisible by %d",
			m.Name(), m.ID(), xyplatform.Default.ID())
		panic(msg)
	}

	if _, ok := manager[m]; ok {
		msg := fmt.Sprintf("Module %s had already registered", m)
		panic(msg)
	}

	manager[m] = 0

	return true
}
