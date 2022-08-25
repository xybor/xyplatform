package xyerror

import (
	"fmt"
)

// Class is a special error with error number and error name. Error number is a
// unique number of Class and helps to determine which module the error Class
// belongs to.
//
// The main purpose of error Class is creating a XyError, so that it should not
// be used for returning.
//
// A Class can be created by one or many parent Classes. The Class without
// parent is called root Class..
type Class struct {
	// The unique number of each Class.
	errno int

	// The error name
	name string

	// The parent classes
	parent []Class
}

// NewClass creates a root Class with error number will be determined by
// module's id in Generator.
func (gen Generator) NewClass(name string, args ...any) Class {
	manager[gen].count++
	return Class{
		errno:  manager[gen].count + gen.id,
		name:   fmt.Sprintf(name, args...),
		parent: nil,
	}
}

// NewClass creates a new Class with called Class as parent.
func (c Class) NewClass(name string, args ...any) Class {
	var gen = getGenerator(c.errno)
	var class = gen.NewClass(name, args...)
	class.parent = []Class{c}
	return class
}

// NewClassM creates a new error class with this class as parent. It has another
// errorid and the same name.
func (c Class) NewClassM(gen Generator) Class {
	var class = gen.NewClass(c.name)
	class.parent = []Class{c}
	return class
}

// New creates a XyError with an error message.
func (c Class) New(msg string, a ...any) XyError {
	return XyError{c: c, msg: fmt.Sprintf(msg, a...)}
}

// belongsTo checks if a Class is inherited from a target class. A class belongs
// to the target Class if it is created by the target itself or target's child.
func (c Class) belongsTo(t Class) bool {
	if c.errno == t.errno {
		return true
	}

	for i := range c.parent {
		if c.parent[i].belongsTo(t) {
			return true
		}
	}

	return false
}

// Error is the method to treat Class as an error.
func (c Class) Error() string {
	return fmt.Sprintf("[%d] %s", c.errno, c.name)
}
