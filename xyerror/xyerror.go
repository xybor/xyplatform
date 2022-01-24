package xyerror

import (
	"fmt"

	"github.com/xybor/xyplatform"
)

// XyError is a special error with determined error number and module.
type XyError struct {
	errno  int
	errmsg string
	parent *XyError
}

// Errno returns the error number.
func (xerr XyError) Errno() int {
	return xerr.errno
}

// Errmsg returns the error message.
func (xerr XyError) Errmsg() string {
	return xerr.errmsg
}

// Error returns the error string including errno and errmsg.
func (xerr XyError) Error() string {
	if xerr.parent != nil {
		return fmt.Sprintf("[%d][%s] %s",
			xerr.errno, xerr.parent.errmsg, xerr.errmsg)
	}

	return fmt.Sprintf("[%d] %s", xerr.errno, xerr.errmsg)
}

// New creates an end-error from an error type.
func (xerr XyError) New(msg string, args ...interface{}) XyError {
	return XyError{
		errno:  xerr.errno,
		parent: &xerr,
		errmsg: fmt.Sprintf(msg, args...),
	}
}

// NewType creates a new error type by inheriting a parent error type. The new
// error type now belongs to the error chain of the parent one.
func (xerr XyError) NewType(name string, args ...interface{}) XyError {
	var m = extractModule(xerr.errno)

	return XyError{
		errno:  nextErrno(m),
		errmsg: fmt.Sprintf(name, args...),
		parent: &xerr,
	}
}

// IsA checks if an end-error is created by the error type t or not.
func (xerr XyError) IsA(t XyError) bool {
	return xerr.errno == t.errno
}

// IsNotA is a negative clause of IsA
func (xerr XyError) IsNotA(t XyError) bool {
	return !xerr.IsA(t)
}

// BelongsTo checks if an end-error belongs to an error chain of a error type
// or not.
func (xerr XyError) BelongsTo(t XyError) bool {
	for {
		if xerr.errno == t.errno {
			return true
		}

		if xerr.parent == nil {
			break
		}

		xerr = *xerr.parent
	}

	return false
}

// NotBelongTo is a negative clause of BelongsTo
func (xerr XyError) NotBelongTo(t XyError) bool {
	return !xerr.BelongsTo(t)
}

// NewType creates a root error type of a module.
func NewType(m xyplatform.Module, name string, args ...interface{}) XyError {
	return XyError{
		errno:  nextErrno(m),
		errmsg: fmt.Sprintf(name, args...),
		parent: nil,
	}
}
