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

func newfilterer() *filterer {
	return &filterer{
		filters: nil,
		lock:    xylock.RWLock{},
	}
}

// AddFilter adds a specified filter.
func (ftr *filterer) AddFilter(f Filter) {
	ftr.lock.WLockFunc(func() {
		if xycond.MustNotContainA(ftr.filters, f) {
			ftr.filters = append(ftr.filters, f)
		}
	})
}

// RemoveFilter removes an existed filter.
func (ftr *filterer) RemoveFilter(f Filter) {
	ftr.lock.WLockFunc(func() {
		for i := range ftr.filters {
			if ftr.filters[i] == f {
				ftr.filters = append(ftr.filters[:i], ftr.filters[i+1:]...)
				break
			}
		}
	})
}

// GetFilters returns all filters of filterer.
func (ftr *filterer) GetFilters() []Filter {
	return ftr.lock.RLockFunc(func() any { return ftr.filters }).([]Filter)
}

// filter checks all filters in filterer, if there is any failed filter, it will
// returns false.
func (ftr *filterer) filter(record LogRecord) bool {
	// Avoid calling locks.
	if len(ftr.filters) == 0 {
		return true
	}

	return ftr.lock.RLockFunc(func() any {
		for i := range ftr.filters {
			if !ftr.filters[i].Filter(record) {
				return false
			}
		}
		return true
	}).(bool)
}
