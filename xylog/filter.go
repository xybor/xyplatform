package xylog

import (
	"reflect"

	"github.com/xybor/xyplatform/xycond"
	"github.com/xybor/xyplatform/xylock"
)

// Filter instances are used to perform arbitrary filtering of LogRecord.
type Filter interface {
	Filter(record LogRecord) bool
}

// A base class for loggers and handlers which allows them to share common code.
type filterer struct {
	filters []Filter
	lock    xylock.RWLock
}

func newfilterer() filterer {
	return filterer{
		filters: nil,
		lock:    xylock.RWLock{},
	}
}

// AddFilter adds a specified filter. Passed filter must be a pointer.
func (fr *filterer) AddFilter(f Filter) {
	xycond.IsKind(f, reflect.Pointer).Assert("Expected a pointer of filter")

	fr.lock.WLockFunc(func() {
		if xycond.NotContainA(fr.filters, f) {
			fr.filters = append(fr.filters, f)
		}
	})
}

// RemoveFilter removes a specified filter. Passed filter must be a pointer.
func (fr *filterer) RemoveFilter(f Filter) {
	xycond.IsKind(f, reflect.Pointer).Assert("Expected a pointer of filter")

	fr.lock.WLockFunc(func() {
		for i := range fr.filters {
			if fr.filters[i] == f {
				fr.filters = append(fr.filters[:i], fr.filters[i+1:]...)
				break
			}
		}
	})
}

// filter checks all filters in filterer, if there is any failed filter, it will
// returns false.
func (fr *filterer) filter(record LogRecord) bool {
	// Avoid calling locks.
	if len(fr.filters) == 0 {
		return true
	}

	return fr.lock.RLockFunc(func() any {
		for i := range fr.filters {
			if !fr.filters[i].Filter(record) {
				return false
			}
		}
		return true
	}).(bool)
}
