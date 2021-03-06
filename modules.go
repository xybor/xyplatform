package xyplatform

import "fmt"

type Module struct {
	id   int
	name string
}

func (m Module) ID() int {
	return m.id
}

func (m Module) Name() string {
	return m.name
}

func (m Module) String() string {
	return fmt.Sprintf("[%d]%s", m.id, m.name)
}

func NewModule(id int, name string) Module {
	return Module{id: id, name: name}
}

var (
	Default    = NewModule(10000, "Default")
	XySelector = NewModule(20000, "XySelector")
	XyCron     = NewModule(30000, "XyCron")
)
