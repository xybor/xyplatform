package xyerror_test

import (
	"testing"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xyerror"
)

func TestClassNewClass(t *testing.T) {
	var id = nextid()
	var egen = xyerror.Register("gen", id)
	var c1 = egen.NewClass("class1")
	var c2 = c1.NewClass("class2")
	xycond.ExpectEqual(c1.Error(), classmsg(id+1, "class1")).Test(t)
	xycond.ExpectEqual(c2.Error(), classmsg(id+2, "class2")).Test(t)
}

func TestClassNewClassM(t *testing.T) {
	var id1 = nextid()
	var id2 = nextid()
	var egen1 = xyerror.Register("gen", id1)
	var egen2 = xyerror.Register("gen", id2)
	var c1 = egen1.NewClass("class")
	var c2 = c1.NewClassM(egen2)
	xycond.ExpectEqual(c2.Error(), classmsg(id2+1, "class")).Test(t)
}
