package xylog

import (
	"fmt"

	"github.com/xybor/xyplatform"
)

var manager = make(map[xyplatform.Module]*logger)

// Register adds a module to log manager and creates its logger. The returned
// value of this function is meaningless, it only helps run the function in the
// global scope.
func Register(m xyplatform.Module) bool {
	if _, ok := manager[m]; ok {
		msg := fmt.Sprintf("Module %s had already registered", m)
		panic(msg)
	}

	manager[m] = &logger{
		module: m,
		output: Stdout,
		format: "$TIME$ -- $MODULE$ [$LEVEL$] $MESSAGE$",
		allow:  map[string]bool{"ALL": true},
	}

	return true
}

// Logger gets the logger of a registered module.
func Logger(m xyplatform.Module) *logger {
	lg, ok := manager[m]
	if !ok {
		msg := fmt.Sprintf("Module %s had not registered yet", m)
		panic(msg)
	}

	return lg
}
