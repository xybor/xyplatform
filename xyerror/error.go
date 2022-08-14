package xyerror

import (
	"errors"
)

// XyError is the error of xyplatform. It supports checking if an error belongs
// to a class or not.
//
// errors.Is(err, cls) returns true if err is created by cls itself or cls's
// child class.
type XyError struct {
	// error class
	c Class

	// error message
	msg string
}

// Error is the method to treat XyError as an error.
func (xerr XyError) Error() string {
	return xerr.msg
}

// Is is the method used to customize errors.Is method.
func (xerr XyError) Is(target error) bool {
	if !errors.As(target, &Class{}) {
		return false
	}

	tc := target.(Class)

	return xerr.c.belongsTo(tc)
}

// Or returns the first not-nil error. If all errors are nil, return nil.
func Or(errs ...error) error {
	for i := range errs {
		if errs[i] != nil {
			return errs[i]
		}
	}

	return nil
}
