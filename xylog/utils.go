package xylog

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/xybor/xyplatform"
	"github.com/xybor/xyplatform/xyerror"
)

func replaceLog(
	f string, m xyplatform.Module, l string, msg string, a ...interface{},
) string {
	content := f

	content = strings.Replace(content, "$MODULE$", m.Name(), 1)

	content = strings.Replace(content, "$LEVEL$", l, 1)

	t := time.Now().Format("01-02-2006 15:04:05")
	content = strings.Replace(content, "$TIME$", t, 1)

	msg = fmt.Sprintf(msg, a...)
	content = strings.Replace(content, "$MESSAGE$", msg, 1)

	return content
}

func appendFile(fn string, data []byte) xyerror.XyError {
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return xyerror.IOError.New("Cannot open file %s: %s", fn, err)
	}

	defer f.Close()

	data = append(data, '\n')
	_, err = f.Write(data)
	if err != nil {
		return xyerror.IOError.New("Cannot write to file %s: %s", fn, err)
	}

	return xyerror.Success
}
