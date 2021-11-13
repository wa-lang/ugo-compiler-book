package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var DebugMode = false

func callerInfo(skip int) (filename string, line int) {
	_, filename, line, _ = runtime.Caller(skip + 1)
	if wd, _ := os.Getwd(); wd != "" {
		if rel, err := filepath.Rel(wd, filename); err == nil {
			filename = rel
		}
	}
	return
}

func Debug(a ...interface{}) {
	if DebugMode {
		filename, line := callerInfo(1)
		msg := fmt.Sprint(a...)
		fmt.Printf("%s:%d: %s", filename, line, msg)
	}
}

func Debugln(a ...interface{}) {
	if DebugMode {
		filename, line := callerInfo(1)
		msg := fmt.Sprintln(a...)
		fmt.Printf("%s:%d: %s", filename, line, msg)
	}
}

func Debugf(format string, a ...interface{}) {
	if DebugMode {
		filename, line := callerInfo(1)
		msg := fmt.Sprintf(format, a...)
		fmt.Printf("%s:%d: %s", filename, line, msg)
	}
}

func Print(a ...interface{}) {
	filename, line := callerInfo(1)
	msg := fmt.Sprint(a...)
	fmt.Printf("%s:%d: %s", filename, line, msg)
}

func Println(a ...interface{}) {
	filename, line := callerInfo(1)
	msg := fmt.Sprintln(a...)
	fmt.Printf("%s:%d: %s", filename, line, msg)
}

func Printf(format string, a ...interface{}) {
	filename, line := callerInfo(1)
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s:%d: %s", filename, line, msg)
}
