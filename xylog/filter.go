package xylog

import (
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

// AddFilter adds a specified filter.
func (base *filterer) AddFilter(f Filter) {
	base.lock.WLockFunc(func() {
		if xycond.NotContainA(base.filters, f) {
			base.filters = append(base.filters, f)
		}
	})
}

// RemoveFilter removes an existed filter.
func (base *filterer) RemoveFilter(f Filter) {
	base.lock.WLockFunc(func() {
		for i := range base.filters {
			if base.filters[i] == f {
				base.filters = append(base.filters[:i], base.filters[i+1:]...)
				break
			}
		}
	})
}

// filter checks all filters in filterer, if there is any failed filter, it will
// returns false.
func (f *filterer) filter(record LogRecord) bool {
	// Avoid calling locks.
	if len(f.filters) == 0 {
		return true
	}

	return f.lock.RLockFunc(func() any {
		for i := range f.filters {
			if !f.filters[i].Filter(record) {
				return false
			}
		}
		return true
	}).(bool)
}
