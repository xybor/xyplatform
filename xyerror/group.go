package xyerror

import (
	"fmt"
)

type Group []Class

// Combine supports creating a group of error classes. This group can be used
// to create the Class with multiparents.
func Combine(cs ...Class) Group {
	return cs
}

// NewClass creates a Class with multiparents.
func (g Group) NewClass(gen Generator, name string, a ...any) Class {
	var child = gen.NewClass(fmt.Sprintf(name, a...))
	child.parent = g
	return child
}
