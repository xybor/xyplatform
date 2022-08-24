package xylog

import (
	"fmt"

	"github.com/xybor/xyplatform/xycond"
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
type textFormatter struct {
	formatstring  string
	attrbuteIndex []int
}

// NewTextFormatter creates a textFormatter which uses LogRecord attributes to
// contribute logging string, e.g. %(message)s or %(levelno)d. See LogRecord for
// more details.
func NewTextFormatter(s string) textFormatter {
	var record = LogRecord{}
	var attributeIndex []int
	var fmtstr = ""
	var i, n = 0, len(s)
	for i < n {
		fmtstr += string(s[i])
		if s[i] == '%' {
			xycond.MustTrue(i+1 < n).Assert("unexpectedly end with %%")
			i++
			switch s[i] {
			case '%':
			case '(':
				i++
				var token = ""
				for {
					xycond.MustTrue(i < n).Assert("uncompleted token %s", token)
					if s[i] == ')' {
						break
					}
					token += string(s[i])
					i++
				}
				attributeIndex = append(attributeIndex, record.mapName(token))
			default:
				xycond.Panic("unexpected token: %s", s[i-2:i])
			}
		}
		i++
	}

	return textFormatter{
		formatstring:  fmtstr,
		attrbuteIndex: attributeIndex,
	}
}

// Format creates a logging string by combine format string and logging record.
func (f textFormatter) Format(record LogRecord) string {
	var attrs = make([]any, len(f.attrbuteIndex))
	for i := range f.attrbuteIndex {
		attrs[i] = record.mapIndex(f.attrbuteIndex[i])
	}
	return fmt.Sprintf(f.formatstring, attrs...)
}
