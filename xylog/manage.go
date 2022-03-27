package xylog

import (
	"strings"

	"github.com/xybor/xyplatform"
	"github.com/xybor/xyplatform/xycond"
)

var manager = make(map[xyplatform.Module]*logger)

const maxLevel = uint(100)

var levelManager = make([]string, maxLevel)

// Register adds a module to log manager and creates its logger. The returned
// value of this function is meaningless, it only helps run the function in the
// global scope.
func Register(m xyplatform.Module) bool {
	xycond.NotContains(manager, m).Assertf("Module %s had already registered", m)

	manager[m] = &logger{
		module: m,
		output: Stdout,
		format: "$TIME$ -- $MODULE$ [$LEVEL$] $MESSAGE$",
		level:  0,
	}

	return true
}

// RegisterLevel adds a log level to the manager. The returned value of this
// function is equal to the value parameter, and it also helps run the function
// in the global scope.
func RegisterLevel(value uint, name string) uint {
	xycond.Condition(value < maxLevel).
		Assertf("Level value must less than %d, but got %d", maxLevel, value)

	xycond.StringEmpty(levelManager[value]).
		Assertf("Level value %d has already registered", value)

	levelManager[value] = strings.ToUpper(name)

	return value
}

// Logger gets the logger of a registered module.
func Logger(m xyplatform.Module) *logger {
	xycond.Contains(manager, m).Assertf("Module %s had not registered yet", m)
	return manager[m]
}

func checkLevel(level uint) {
	xycond.Condition(levelManager[level] != "").
		Assertf("Level %d has not registered yet", level)
}
