package xyerror

// Default is the default error generator.
var Default = Register("default", 100000)

// Default predefined errors.
var (
	Error               = Default.NewClass("Error")
	IOError             = Default.NewClass("IOError")
	FloatingPointError  = Default.NewClass("FloatingPointError")
	IndexError          = Default.NewClass("IndexError")
	KeyError            = Default.NewClass("KeyError")
	NotImplementedError = Default.NewClass("NotImplementedError")
	ValueError          = Default.NewClass("ValueError")
	ParameterError      = Default.NewClass("ParameterError")
	TypeError           = Default.NewClass("TypeError")
	AssertionError      = Default.NewClass("AssertionError")
)
