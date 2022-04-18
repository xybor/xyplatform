package xyerror

import (
	"errors"
	"fmt"
)

// xyerror is the error of xybor projects. It supports checking if an error
// belongs to a class or not.
//
// errors.Is(err, cls) returns true if err is created by:
// 	- cls.
// 	- another class created by cls.
type xyerror struct {
	// error class
	c class

	// error message
	msg string
}

func (xerr xyerror) Error() string {
	return fmt.Sprintf("%s %s", xerr.c.Error(), xerr.msg)
}

func (xerr xyerror) Is(target error) bool {
	if !errors.As(target, &class{}) {
		return false
	}

	tc := target.(class)

	return xerr.c.belongsTo(tc)
}

// Or returns the first not-nil error. If all errors are nil, return nil.
func Or(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
