package xylog

import "fmt"

type configurator interface {
	apply(*logger)
}

type formatCfg string

func (cfg formatCfg) apply(lg *logger) {
	lg.format = string(cfg)
}

// Format is a configurator which adjusts the log string.
// There are some macros used to format:
//     $TIME$    - time in format dd-mm-yy hh:mm:ss when logging
//     $LEVEL$   - the log level
//     $MODULE$  - module name
//     $MESSAGE$ - log message
// Example: "$TIME$ -- $MODULE$ [$LEVEL$]  $MESSAGE$".
func Format(f string) formatCfg {
	return formatCfg(f)
}

type allowCfg uint

func (cfg allowCfg) apply(lg *logger) {
	lg.level = uint(cfg)
}

// Allow is configurator which indicates log levels to be allowed to print.
func Allow(level uint) allowCfg {
	if level >= maxLevel {
		msg := fmt.Sprintf("The level value is expected less than %d, but "+
			"got %d", maxLevel, level)
		panic(msg)
	}

	return allowCfg(level)
}

// AllowAll help you print all log levels.
func AllowAll() allowCfg {
	return allowCfg(0)
}

// NoAllow help you print no log.
func NoAllow() allowCfg {
	return allowCfg(maxLevel)
}

type writerCfg struct {
	writer
}

func (cfg writerCfg) apply(lg *logger) {
	lg.output = cfg
}

// Writer is a configurator which chooses the type of output you want to print
// the log.
func Writer(w writer) writerCfg {
	return writerCfg{w}
}

// StdWriter is a shortcut of Writer(Stdout)
func StdWriter() writerCfg {
	return Writer(Stdout)
}
