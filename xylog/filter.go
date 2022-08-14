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
func (fr *filterer) RemoveFilter(f Filter) {
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
