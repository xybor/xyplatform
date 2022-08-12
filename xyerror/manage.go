package xyerror

import (
	"github.com/xybor/xyplatform/xycond"
)

// errorid is the type of a group of errors in a module.
type errorid int

// erroinfo includes the name and the number of created errors of an error id.
type errorinfo struct {
	name  string
	count int
}

// The minimum and default id of module
const minId errorid = 100000

// manager is a map of errorid as key and errorinfo as value.
var manager = make(map[errorid]*errorinfo)

// getErrorId returns the error id with the given errno.
func getErrorId(errno int) errorid {
	for id := range manager {
		d := errno - int(id)
		if d < 0 || d > int(eid) {
			continue
		}

		if d < int(minId) {
			return id
		}
	}

	return 0
}

// Register adds a Module to pool for managing new error types.
func Register(name string, id int) errorid {
	xycond.Divisible(id, int(minId)).Assert(
		"Cannot register: %d is not divisible by %d", id, minId)
	var eid = errorid(id)
	xycond.NotContainM(manager, eid).Assert("Id %d had already registered", id)

	manager[eid] = &errorinfo{name: name, count: 0}
	return eid
}
