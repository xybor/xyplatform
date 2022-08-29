package xyerror_test

import (
	"testing"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xyerror"
)

func TestXyError(t *testing.T) {
	var id = nextid()
	var egen = xyerror.Register("", id)
	var c = egen.NewClass("class")
	var xerr1 = c.Newf("error-%d", 1)
	var xerr2 = c.New("error-2")

	xycond.ExpectEqual(xerr1.Error(), "class: error-1").Test(t)
	xycond.ExpectEqual(xerr2.Error(), "class: error-2").Test(t)
}

func TestXyErrorIs(t *testing.T) {
	var err1 = xyerror.ValueError.New("err1")
	var err2 = xyerror.TypeError.New("err2")

	xycond.ExpectError(err1, xyerror.ValueError).Test(t)
	xycond.ExpectError(err2, xyerror.TypeError).Test(t)

	xycond.ExpectErrorNot(err1, err1).Test(t)
	xycond.ExpectErrorNot(err1, err2).Test(t)
	xycond.ExpectErrorNot(err1, xyerror.TypeError).Test(t)
}

func TestOr(t *testing.T) {
	var err1 = xyerror.ValueError.New("err1")
	var err2 = xyerror.TypeError.New("err2")
	var err3 error

	xycond.ExpectError(xyerror.Or(err1, err2), xyerror.ValueError).Test(t)
	xycond.ExpectError(xyerror.Or(err2, err1), xyerror.TypeError).Test(t)
	xycond.ExpectNil(xyerror.Or(nil, err3)).Test(t)
}

func TestCombine(t *testing.T) {
	var gen = xyerror.Register("", nextid())
	var c = xyerror.Combine(xyerror.ValueError, xyerror.TypeError).
		NewClass(gen, "class")
	var xerr = c.New("error")

	xycond.ExpectError(xerr, xyerror.ValueError).Test(t)
	xycond.ExpectError(xerr, xyerror.TypeError).Test(t)
	xycond.ExpectErrorNot(xerr, xyerror.IndexError).Test(t)
}
