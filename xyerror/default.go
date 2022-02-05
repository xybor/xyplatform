package xyerror

import "github.com/xybor/xyplatform"

var _ = Register(xyplatform.Default)

var (
	UnknownError        = NewClass(xyplatform.Default, "UnknownError")
	IOError             = NewClass(xyplatform.Default, "IOError")
	FloatingPointError  = NewClass(xyplatform.Default, "FloatingPointError")
	IndexError          = NewClass(xyplatform.Default, "IndexError")
	KeyError            = NewClass(xyplatform.Default, "KeyError")
	NotImplementedError = NewClass(xyplatform.Default, "NotImplementedError")
	ValueError          = NewClass(xyplatform.Default, "ValueError")
	ParameterError      = NewClass(xyplatform.Default, "ParameterError")
	TypeError           = NewClass(xyplatform.Default, "TypeError")
)
