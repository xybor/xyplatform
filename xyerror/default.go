package xyerror

var eid = Register("default", 100000)

var (
	UnknownError        = eid.NewClass("UnknownError")
	IOError             = eid.NewClass("IOError")
	FloatingPointError  = eid.NewClass("FloatingPointError")
	IndexError          = eid.NewClass("IndexError")
	KeyError            = eid.NewClass("KeyError")
	NotImplementedError = eid.NewClass("NotImplementedError")
	ValueError          = eid.NewClass("ValueError")
	ParameterError      = eid.NewClass("ParameterError")
	TypeError           = eid.NewClass("TypeError")
)
