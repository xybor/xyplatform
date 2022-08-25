// Package xyerror defines error type used in xyplatform.
package xyerror

import (
	"log"
)

// Generator is used to generate root Class for every module. It is determined
// by the identifier of module.
type Generator struct {
	// The identifier of module.
	id int
}

// erroinfo includes the name and the number of created errors of an error id.
type errorinfo struct {
	name  string
	count int
}

// The minimum and default id of module
var minid = 100000

// manager is a map of errorid as key and errorinfo as value.
var manager = make(map[Generator]*errorinfo)

// getGenerator returns the Generator with the given errno.
func getGenerator(errno int) Generator {
	for gen := range manager {
		var d = errno - gen.id
		if d < 0 || d > gen.id {
			continue
		}

		if d < minid {
			return gen
		}
	}

	return Generator{0}
}

// Register adds a Module with its identifier to managing pool for creating new
// Classes.
func Register(name string, id int) Generator {
	if id%minid != 0 {
		log.Panicf("cannot register, %d is not divisible by %d", id, minid)
	}
	var gen = Generator{id}
	if _, ok := manager[gen]; ok {
		log.Panicf("id %d had already registered", id)
	}

	manager[gen] = &errorinfo{name: name, count: 0}
	return gen
}
