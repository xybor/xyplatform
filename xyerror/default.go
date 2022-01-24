package xyerror

import "github.com/xybor/xyplatform"

var _ = Register(xyplatform.Default)

var (
	Success = XyError{errno: 0, errmsg: "Success", parent: nil}

	UnknownError        = NewType(xyplatform.Default, "UnknownError")
	IOError             = NewType(xyplatform.Default, "IOError")
	FloatingPointError  = NewType(xyplatform.Default, "FloatingPointError")
	IndexError          = NewType(xyplatform.Default, "IndexError")
	KeyError            = NewType(xyplatform.Default, "KeyError")
	NotImplementedError = NewType(xyplatform.Default, "NotImplementedError")
	ValueError          = NewType(xyplatform.Default, "ValueError")
)
