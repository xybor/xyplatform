package xyerror

import (
	"fmt"
)

type group []class

// Combine supports creating a group of error classes. This group can be used
// to create the error class with multiparents.
//
// For example, xyerror.Combine(C1, C2).NewClass(xyerror.DEFAULT, "New")
func Combine(cs ...class) group {
	return cs
}

// NewClass creates an error class of multiparents with format string.
func (g group) NewClass(eid errorid, name string, a ...any) class {
	var child = eid.NewClass(fmt.Sprintf(name, a...))
	child.parent = g
	return child
}
