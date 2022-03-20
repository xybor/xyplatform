package xyerror

import (
	"fmt"

	"github.com/xybor/xyplatform"
)

type group []class

// Combine supports creating a group of error classes. This group can be used
// to create the error class with multiparents.
//
// For example, xyerror.Combine(C1, C2).NewClass(xyplatform.DEFAULT, "New")
func Combine(cs ...class) group {
	return cs
}

// NewClass creates an error class of multiparents.
func (g group) NewClass(m xyplatform.Module, name string) class {
	return class{
		errno:  nextErrno(m),
		name:   name,
		parent: g,
	}
}

// NewClassf creates an error class of multiparents with format string.
func (g group) NewClassf(m xyplatform.Module, name string, a ...interface{}) class {
	return class{
		errno:  nextErrno(m),
		name:   fmt.Sprintf(name, a...),
		parent: g,
	}
}
