package xylog

import (
	"fmt"
	"os"
	"time"
)

type stopper interface {
	isStop(fn string, format string) bool
}

type timeAfter struct {
	d time.Duration
}

// TimeAfter is a stopper causing SFile to create another file after a
// specified time.
func TimeAfter(duration time.Duration) timeAfter {
	return timeAfter{d: duration}
}

func (a timeAfter) isStop(fn string, format string) bool {
	var timeStr string
	fmt.Sscanf(fn, format, timeStr)
	start, err := time.Parse(TimeFormat, timeStr)
	if err != nil {
		fmt.Printf("WARNING: Something was wrong when parse time in "+
			"file %s: %s\n", fn, err)
		return false
	}

	return start.Add(a.d).After(time.Now())
}

type timePeriod struct {
	p time.Duration
}

// TimePeriod is a stopper causing SFile to create another file when a new time
// period comes.
func TimePeriod(p time.Duration) timePeriod {
	if p != Minute && p != Hour && p != Day && p != Week {
		panic("Only support minute, hour, day, and week for time period")
	}

	return timePeriod{p: p}
}

func (p timePeriod) isStop(fn string, format string) bool {
	var timeStr string
	fmt.Sscanf(fn, format, timeStr)
	last, err := time.Parse(TimeFormat, timeStr)
	if err != nil {
		fmt.Printf("WARNING: Something wrong when parse time in file %s: %s\n", fn, err)
		return false
	}

	switch p.p {
	case Minute:
		if time.Now().Minute() != last.Minute() {
			return true
		}
	case Hour:
		if time.Now().Hour() != last.Hour() {
			return true
		}
	case Day:
		if time.Now().Day() != last.Day() {
			return true
		}
	case Week:
		_, nowWeek := time.Now().ISOWeek()
		_, lastWeek := last.ISOWeek()

		if nowWeek != lastWeek {
			return true
		}
	}

	return false
}

type limitSize struct {
	sz int64
}

// LimitSize is a stopper causing SFile to create another file if the file size
// is exceed the limit.
func LimitSize(sz int64) limitSize {
	return limitSize{sz: sz}
}

func (ls limitSize) isStop(fn string, format string) bool {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Printf("WARNING: Cannot open file %s: %s\n", fn, err)
		return false
	}

	fs, err := f.Stat()
	if err != nil {
		fmt.Printf("WARNING: Cannot read size of file %s: %s\n", fn, err)
		return false
	}

	return fs.Size() > ls.sz
}
