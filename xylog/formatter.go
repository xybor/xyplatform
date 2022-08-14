// This file copied and modified comments of python logging.
package xylog

import (
	"fmt"
	"reflect"
	"strings"
)

// Formatter instances are used to convert a LogRecord to text.
//
// Formatter need to know how a LogRecord is constructed. They are responsible
// for converting a LogRecord to a string which can be interpreted by either a
// human or an external system.
type Formatter interface {
	Format(LogRecord) string
}

// The TextFormatter can be initialized with a format string which makes use of
// knowledge of the LogRecord attributes - e.g. %(message)s or %(levelno)d. See
// LogRecord for more details.
type textFormatter string

// NewTextFormatter creates a textFormatter which uses LogRecord attributes to
// contribute logging string, e.g. %(message)s or %(levelno)d. See LogRecord for
// more details.
func NewTextFormatter(f string) textFormatter {
	return textFormatter(f)
}

// Format creates a logging string by combine format string and logging record.
func (f textFormatter) Format(record LogRecord) string {
	var m = tomap(record)
	var s = string(f)
	for k, v := range m {
		var token = "%(" + k + ")" // %(foo)
		for {
			var index = strings.Index(s, token)
			if index == -1 {
				break
			}
			token = s[index : index+len(token)+1]    // %(foo)s
			var verb = "%" + token[len(token)-1:]    // %s
			var value = fmt.Sprintf(verb, v)         // bar
			s = strings.Replace(s, token, value, -1) // %(foo)s -> bar
		}
	}

	return s
}

// tomap converts a struct to map[string]any.
func tomap(a any) map[string]any {
	var result = make(map[string]any)
	var v = reflect.ValueOf(a)
	for i := 0; i < v.NumField(); i++ {
		var f = v.Type().Field(i)
		var tag = f.Tag.Get("map")
		if tag != "" {
			result[tag] = v.Field(i).Interface()
		} else {
			result[strings.ToLower(f.Name)] = v.Field(i).Interface()
		}
	}
	return result
}
