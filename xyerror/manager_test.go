package xyerror_test

import (
	"fmt"
	"testing"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xyerror"
)

var autoid = 1000000

func nextid() int {
	autoid += 100000
	return autoid
}

func classmsg(id int, msg string) string {
	return fmt.Sprintf("[%d] %s", id, msg)
}

func TestInitiateGeneratorClassDirectly(t *testing.T) {
	var errorGen = xyerror.Generator{}
	xycond.ExpectPanic(func() { errorGen.NewClass("foo") }).Test(t)
}

func TestRegister(t *testing.T) {
	var id = nextid()
	var egen = xyerror.Register(t.Name(), id)
	var c = egen.NewClass("bar")
	xycond.ExpectEqual(c.Error(), classmsg(id+1, "bar")).Test(t)
}

func TestRegisterDuplicate(t *testing.T) {
	var id = nextid()
	xyerror.Register(t.Name(), id)
	xycond.ExpectPanic(func() { xyerror.Register("foobar", id) }).Test(t)
}
