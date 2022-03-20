package xylog

import (
	"fmt"
	"sync"
	"time"
)

type writer interface {
	write(string, ...interface{})
}

type stdout struct {
	mutex sync.Mutex
}

// Stdout is a writer which prints log to standard output.
var Stdout = &stdout{mutex: sync.Mutex{}}

func (w *stdout) write(msg string, a ...interface{}) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	content := fmt.Sprintf(msg, a...)
	fmt.Println(content)
}

type file struct {
	fn    string
	mutex sync.Mutex
}

// File is a writer which prints log to a specified file.
func File(fn string) *file {
	sep := "------------------------------------------------------------------"
	err := appendFile(fn, []byte(sep))

	if err != nil {
		msg := fmt.Sprintf("Test writing to %s failed: %s", fn, err)
		panic(msg)
	}

	return &file{fn: fn, mutex: sync.Mutex{}}
}

func (f *file) write(msg string, a ...interface{}) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	content := fmt.Sprintf(msg, a...)
	err := appendFile(f.fn, []byte(content))

	if err != nil {
		fmt.Printf("WARNING: Cannot write this line to file %s: %s", f.fn, err)
		fmt.Printf(msg, a...)
	}
}

type sfile struct {
	format string
	s      stopper
	f      *file
	mutex  sync.Mutex
}

// SFile is a writer which prints log to another file if the stop contidion is
// reach.
func SFile(format string, s stopper) *sfile {
	t := time.Now().Format(TimeFormat)
	fn := fmt.Sprintf(format, t)

	return &sfile{
		format: format,
		f:      File(fn),
		s:      s,
		mutex:  sync.Mutex{},
	}
}

func (sf *sfile) write(msg string, a ...interface{}) {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	if sf.s.isStop(sf.f.fn, sf.format) {
		t := time.Now().Format(TimeFormat)
		fn := fmt.Sprintf(sf.format, t)
		sf.f = File(fn)
	}

	sf.f.write(msg, a...)
}
