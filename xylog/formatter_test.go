package xylog_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xylog"
)

func TestNewTextFormatter(t *testing.T) {
	var f = xylog.NewTextFormatter(
		"time=%(asctime)s %(levelno).3d %(module)s something")
	xycond.MustTrue(strings.Contains(
		fmt.Sprint(f),
		"time=%s %.3d %s something",
	)).Testf(t, "Got unepected formatter %v", f)
}

func TestNewTextFormatterWithPercentageSign(t *testing.T) {
	var f = xylog.NewTextFormatter(
		"%%abc)s")
	xycond.MustTrue(strings.Contains(fmt.Sprint(f), "%abc)s")).
		Testf(t, "Got unepected formatter %v", f)
}

func TestTextFormatter(t *testing.T) {
	var formatter = xylog.NewTextFormatter(
		"%(asctime)s %(created)d %(filename)s %(funcname)s %(levelname)s " +
			"%(levelno)d %(lineno)d %(message)s %(module)s %(msecs)d " +
			"%(name)s %(pathname)s %(process)d %(relativeCreated)d")

	var s = formatter.Format(xylog.LogRecord{
		Asctime:         "ASCTIME",
		Created:         1,
		FileName:        "FILENAME",
		FuncName:        "FUNCNAME",
		LevelName:       "LEVELNAME",
		LevelNo:         2,
		LineNo:          3,
		Message:         "MESSAGE",
		Module:          "MODULE",
		Msecs:           4,
		Name:            "NAME",
		PathName:        "PATHNAME",
		Process:         5,
		RelativeCreated: 6,
	})

	xycond.MustEqual(s, "ASCTIME 1 FILENAME FUNCNAME LEVELNAME 2 3 MESSAGE "+
		"MODULE 4 NAME PATHNAME 5 6",
	).Testf(t, "Got unexpected formatted string: %s", s)
}
