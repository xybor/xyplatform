package xylog

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/xybor/xyplatform"
	"github.com/xybor/xyplatform/xyerror"
)

func parseTemplate(e string, t T) string {
	var s string = fmt.Sprintf("event=%s", e)

	for k, v := range t {
		s += fmt.Sprintf(" %s=%s", k, v)
	}

	return s
}

func makeRecord(f string, m xyplatform.Module, lvl string, msg string, a ...interface{}) string {
	content := f

	content = strings.Replace(content, "$MODULE$", m.Name(), 1)

	content = strings.Replace(content, "$LEVEL$", lvl, 1)

	t := time.Now().Format("01-02-2006 15:04:05")
	content = strings.Replace(content, "$TIME$", t, 1)

	msg = fmt.Sprintf(msg, a...)
	content = strings.Replace(content, "$MESSAGE$", msg, 1)

	return content
}

func appendFile(fn string, data []byte) error {
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return xyerror.IOError.Newf("Cannot open file %s: %s", fn, err)
	}

	defer f.Close()

	data = append(data, '\n')
	_, err = f.Write(data)
	if err != nil {
		return xyerror.IOError.Newf("Cannot write to file %s: %s", fn, err)
	}

	return nil
}
