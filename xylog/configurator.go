package xylog

import "strings"

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

type allowCfg []string

func (cfg allowCfg) apply(lg *logger) {
	for k := range lg.allow {
		delete(lg.allow, k)
	}

	for _, a := range cfg {
		lg.allow[a] = true
	}
}

// Allow is configurator which indicates log levels to be allowed to print.
// Use "ALL" if you want to print all log levels.
func Allow(levels ...string) allowCfg {
	cfg := make(allowCfg, 0)

	for _, l := range levels {
		cfg = append(cfg, strings.ToUpper(l))
	}

	return cfg
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
