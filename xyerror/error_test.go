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
	var xerr = c.New("error")
	xycond.MustEqual(xerr.Error(), "class: error").
		Testf(t, "%s != %s", xerr.Error(), "class: error")
}

func TestXyErrorIs(t *testing.T) {
	var err1 = xyerror.ValueError.New("err1")
	var err2 = xyerror.TypeError.New("err2")

	xycond.ErrorMustBe(err1, xyerror.ValueError).
		Testf(t, "err1 should be xyerror.ValueError")
	xycond.ErrorMustBe(err2, xyerror.TypeError).
		Testf(t, "err2 should be xyerror.TypeError")

	xycond.ErrorMustNotBe(err1, err1).
		Testf(t, "err1 should not be err1")
	xycond.ErrorMustNotBe(err1, err2).
		Testf(t, "err1 should not be err2")
	xycond.ErrorMustNotBe(err1, xyerror.TypeError).
		Testf(t, "err1 should not be xyerror.TypeError")
}

func TestOr(t *testing.T) {
	var err1 = xyerror.ValueError.New("err1")
	var err2 = xyerror.TypeError.New("err2")
	var err3 error

	xycond.ErrorMustBe(xyerror.Or(err1, err2), xyerror.ValueError).
		Testf(t, "err1 or err2 should be the ValueError")
	xycond.ErrorMustBe(xyerror.Or(err2, err1), xyerror.TypeError).
		Testf(t, "err2 or err1 should be the TypeError")
	xycond.MustNil(xyerror.Or(nil, err3)).
		Testf(t, "Or with all nil errors should return nil")
}

func TestCombine(t *testing.T) {
	var gen = xyerror.Register("", nextid())
	var c = xyerror.Combine(xyerror.ValueError, xyerror.TypeError).
		NewClass(gen, "class")
	var xerr = c.New("error")

	xycond.ErrorMustBe(xerr, xyerror.ValueError).
		Testf(t, "xerr should be ValueError")

	xycond.ErrorMustBe(xerr, xyerror.TypeError).
		Testf(t, "xerr should be TypeError")

	xycond.ErrorMustNotBe(xerr, xyerror.IndexError).
		Testf(t, "xerr should not be IndexError")
}
