package xylog

import (
	"fmt"
	"strings"

	"github.com/xybor/xyplatform"
)

var manager = make(map[xyplatform.Module]*logger)

var maxLevel = uint(100)
var levelManager = make([]string, maxLevel)

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
		level:  0,
	}

	return true
}

// RegisterLevel adds a log level to the manager. The returned value of this
// function is equal to the value parameter, and it also helps run the function
// in the global scope.
func RegisterLevel(value uint, name string) uint {
	if value >= maxLevel {
		msg := fmt.Sprintf("The level value is expected less than %d, but "+
			"got %d", maxLevel, value)
		panic(msg)
	}

	if levelManager[value] != "" {
		msg := fmt.Sprintf("Level value %d has already registered", value)
		panic(msg)
	}

	levelManager[value] = strings.ToUpper(name)

	return value
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

func checkLevel(level uint) {
	if levelManager[level] == "" {
		msg := fmt.Sprintf("Level %d has not registered yet", level)
		panic(msg)
	}
}
