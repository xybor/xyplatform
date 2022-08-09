package xyerror

import (
	"fmt"
)

// xyerror.class (also called error class) is a special error with error number
// and error name. Error number is a unique number for each error class and
// helps to determine which module the error class belongs to.
//
// The purpose of error class is creating a xyerror, so that it should not be
// used for returning.
//
// A error class can be created by one or many parent classes. The error class
// without parent is called ROOT error class.
//
// The xyerror belongs to an error class C if this error is created by:
//   + class C.
//   + another error class which is created by C.
type class struct {
	// The unique number of each error class.
	// If the error number is 100001, the module of this error class is DEFAULT
	// (100000).
	errno int

	// The error name of this class
	name string

	// The parent classes
	parent []class
}

// NewClass creates a ROOT error class.
func (eid errorid) NewClass(name string) class {
	manager[eid].count += 1
	return class{
		errno:  manager[eid].count + int(eid),
		name:   name,
		parent: nil,
	}
}

// NewClassf creates a ROOT error class with string format.
func (eid errorid) NewClassf(name string, args ...interface{}) class {
	return eid.NewClass(fmt.Sprintf(name, args...))
}

// NewClass creates a new error class from this class as parent.
func (c class) NewClass(name string) class {
	var eid = getErrorId(c.errno)
	var child = eid.NewClass(name)
	child.parent = []class{c}
	return child
}

// NewClassf creates a new error class from this class as parent with string
// format.
func (c class) NewClassf(name string, args ...interface{}) class {
	return c.NewClass(fmt.Sprintf(name, args...))
}

// NewClassM creates a new error class from this class as parent with another
// errorid and same name.
func (c class) NewClassM(eid errorid) class {
	var child = eid.NewClass(c.name)
	child.parent = []class{c}
	return child
}

// New creates a xyerror.
func (c class) New(msg string) xyerror {
	return xyerror{c: c, msg: msg}
}

// Newf creates a xyerror with string format.
func (c class) Newf(msg string, a ...interface{}) xyerror {
	return xyerror{c: c, msg: fmt.Sprintf(msg, a...)}
}

// belongsTo checks if an error class is inherited from a target class. A class
// belongs to thr target class if:
//   + it is created by the target class.
//   + it is created by the class which belongs to the target class.
func (c class) belongsTo(t class) bool {
	if c.errno == t.errno {
		return true
	}

	for _, p := range c.parent {
		if p.belongsTo(t) {
			return true
		}
	}

	return false
}

func (c class) Error() string {
	return fmt.Sprintf("[%d] %s", c.errno, c.name)
}
