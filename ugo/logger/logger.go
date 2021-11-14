package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var DebugMode = false

func callerInfo(skip int) (fn, filename string, line int) {
	pc, filename, line, _ := runtime.Caller(skip + 1)
	fn = runtime.FuncForPC(pc).Name()
	if idx := strings.LastIndex(fn, "/"); idx >= 0 {
		fn = fn[idx+1:]
	}
	if wd, _ := os.Getwd(); wd != "" {
		if rel, err := filepath.Rel(wd, filename); err == nil {
			filename = rel
		}
	}
	return
}

func Debug(a ...interface{}) {
	if DebugMode {
		fn, filename, line := callerInfo(1)
		msg := fmt.Sprint(a...)
		fmt.Printf("%s:%d: %s: %s", filename, line, fn, msg)
	}
}

func Debugln(a ...interface{}) {
	if DebugMode {
		fn, filename, line := callerInfo(1)
		msg := fmt.Sprintln(a...)
		fmt.Printf("%s:%d: %s: %s", filename, line, fn, msg)
	}
}

func Debugf(format string, a ...interface{}) {
	if DebugMode {
		fn, filename, line := callerInfo(1)
		msg := fmt.Sprintf(format, a...)
		fmt.Printf("%s:%d: %s: %s", filename, line, fn, msg)
	}
}

func Print(a ...interface{}) {
	_, filename, line := callerInfo(1)
	msg := fmt.Sprint(a...)
	fmt.Printf("%s:%d: %s", filename, line, msg)
}

func Println(a ...interface{}) {
	_, filename, line := callerInfo(1)
	msg := fmt.Sprintln(a...)
	fmt.Printf("%s:%d: %s", filename, line, msg)
}

func Printf(format string, a ...interface{}) {
	_, filename, line := callerInfo(1)
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s:%d: %s", filename, line, msg)
}
