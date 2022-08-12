// This file copied and modified comments of python logging.
package xylog

import (
	"fmt"
	"reflect"
	"strings"
)

// formatter instances are used to convert a LogRecord to text.
type formatter interface {
	Format(LogRecord) string
}

// Formatter need to know how a LogRecord is constructed. They are responsible
// for converting a LogRecord to (usually) a string which can be interpreted by
// either a human or an external system. The base Formatter allows a formatting
// string to be specified.
//
// The Formatter can be initialized with a format string which makes use of
// knowledge of the LogRecord attributes - e.g. %(message)s or %(levelname).
// See LogRecord for more details about attributes.
type Formatter string

// Format creates a logging string by combine format string and logging record.
func (f Formatter) Format(record LogRecord) string {
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
