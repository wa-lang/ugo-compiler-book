package logger

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var DebugMode = false

func CallerInfo(skip int) (fn, filename string, line int) {
	return callerInfo(skip + 1)
}

func CallStack() string {
	var buf bytes.Buffer
	for skip := 1; ; skip += 1 {
		fn, filename, line := callerInfo(skip)
		if filename == "" {
			break
		}
		fmt.Fprintf(&buf, "%s:%d: %s", filename, line, fn)
	}
	return buf.String()
}

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

func Assert(ok bool, a ...interface{}) {
	if ok {
		return
	}
	fn, filename, line := callerInfo(1)
	msg := fmt.Sprint(a...)
	if msg != "" {
		panic(fmt.Sprintf("%s:%d: %s: assert failed: %v", filename, line, fn, msg))
	} else {
		panic(fmt.Sprintf("%s:%d: %s: assert failed", filename, line, fn))
	}
}

func Assertf(ok bool, format string, a ...interface{}) {
	if ok {
		return
	}
	fn, filename, line := callerInfo(1)
	msg := fmt.Sprintf(format, a...)
	if msg != "" {
		panic(fmt.Sprintf("%s:%d: %s: assert failed: %v", filename, line, fn, msg))
	} else {
		panic(fmt.Sprintf("%s:%d: %s: assert failed", filename, line, fn))
	}
}

func AssertEQ(a, b interface{}) {
	if a == b {
		return
	}
	fn, filename, line := callerInfo(1)
	panic(fmt.Sprintf("%s:%d: %s: AssertEQ: %v != %v", filename, line, fn, a, b))
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
