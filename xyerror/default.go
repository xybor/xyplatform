package xyerror

var defaultGen = Register("default", 100000)

var (
	UnknownError        = defaultGen.NewClass("UnknownError")
	IOError             = defaultGen.NewClass("IOError")
	FloatingPointError  = defaultGen.NewClass("FloatingPointError")
	IndexError          = defaultGen.NewClass("IndexError")
	KeyError            = defaultGen.NewClass("KeyError")
	NotImplementedError = defaultGen.NewClass("NotImplementedError")
	ValueError          = defaultGen.NewClass("ValueError")
	ParameterError      = defaultGen.NewClass("ParameterError")
	TypeError           = defaultGen.NewClass("TypeError")
	AssertionError      = defaultGen.NewClass("AssertionError")
)
