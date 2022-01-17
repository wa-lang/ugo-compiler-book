package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/wa-lang/ugo/token"
)

// Print 打印语法树到 stdout
func Print(node Node) {
	fprint(os.Stdout, "", "", node)
}

// Fprint 打印语法树到指定目标
func Fprint(w io.Writer, filename, source string, node Node) {
	fprint(w, filename, source, node)
}

func fprint(w io.Writer, filename, source string, x interface{}) (err error) {
	p := printer{
		output:   w,
		filename: filename,
		source:   source,
		ptrmap:   make(map[interface{}]int),
		last:     '\n', // force printing of line number on first line
	}

	if f, ok := x.(*File); ok {
		if p.filename == "" && p.source == "" {
			p.filename = f.Filename
			p.source = f.Source
		}

		file := *f
		if len(file.Source) > 8 {
			file.Source = file.Source[:8] + "..."
		}
		x = file
	}

	// install error handler
	defer func() {
		if e := recover(); e != nil {
			err = e.(localError).err // re-panics if it's not a localError
		}
	}()

	// print x
	if x == nil {
		p.printf("nil\n")
		return
	}
	p.print(reflect.ValueOf(x))
	p.printf("\n")

	return
}

type printer struct {
	output   io.Writer
	filename string
	source   string
	ptrmap   map[interface{}]int
	indent   int
	last     byte
	line     int
}

func (p *printer) notNilFilter(_name string, v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return !v.IsNil()
	}
	return true
}

func (p *printer) Write(data []byte) (n int, err error) {
	var m int
	for i, b := range data {
		// invariant: data[0:n] has been written
		if b == '\n' {
			m, err = p.output.Write(data[n : i+1])
			n += m
			if err != nil {
				return
			}
			p.line++
		} else if p.last == '\n' {
			_, err = fmt.Fprintf(p.output, "%6d  ", p.line)
			if err != nil {
				return
			}
			for j := p.indent; j > 0; j-- {
				_, err = p.output.Write([]byte(".  "))
				if err != nil {
					return
				}
			}
		}
		p.last = b
	}
	if len(data) > n {
		m, err = p.output.Write(data[n:])
		n += m
	}
	return
}

// localError wraps locally caught errors so we can distinguish
// them from genuine panics which we don't want to return as errors.
type localError struct {
	err error
}

func (p *printer) printf(format string, args ...interface{}) {
	if _, err := fmt.Fprintf(p, format, args...); err != nil {
		panic(localError{err})
	}
}

func (p *printer) print(x reflect.Value) {
	if !p.notNilFilter("", x) {
		p.printf("nil")
		return
	}

	switch x.Kind() {
	case reflect.Interface:
		p.print(x.Elem())

	case reflect.Map:
		p.printf("%s (len = %d) {", x.Type(), x.Len())
		if x.Len() > 0 {
			p.indent++
			p.printf("\n")
			for _, key := range x.MapKeys() {
				p.print(key)
				p.printf(": ")
				p.print(x.MapIndex(key))
				p.printf("\n")
			}
			p.indent--
		}
		p.printf("}")

	case reflect.Ptr:
		p.printf("*")
		// type-checked ASTs may contain cycles - use ptrmap
		// to keep track of objects that have been printed
		// already and print the respective line number instead
		ptr := x.Interface()
		if line, exists := p.ptrmap[ptr]; exists {
			p.printf("(obj @ %d)", line)
		} else {
			p.ptrmap[ptr] = p.line
			p.print(x.Elem())
		}

	case reflect.Array:
		p.printf("%s {", x.Type())
		if x.Len() > 0 {
			p.indent++
			p.printf("\n")
			for i, n := 0, x.Len(); i < n; i++ {
				p.printf("%d: ", i)
				p.print(x.Index(i))
				p.printf("\n")
			}
			p.indent--
		}
		p.printf("}")

	case reflect.Slice:
		if s, ok := x.Interface().([]byte); ok {
			p.printf("%#q", s)
			return
		}
		p.printf("%s (len = %d) {", x.Type(), x.Len())
		if x.Len() > 0 {
			p.indent++
			p.printf("\n")
			for i, n := 0, x.Len(); i < n; i++ {
				p.printf("%d: ", i)
				p.print(x.Index(i))
				p.printf("\n")
			}
			p.indent--
		}
		p.printf("}")

	case reflect.Struct:
		t := x.Type()
		p.printf("%s {", t)
		p.indent++
		first := true
		for i, n := 0, t.NumField(); i < n; i++ {
			name := t.Field(i).Name
			value := x.Field(i)
			if p.notNilFilter(name, value) {
				if first {
					p.printf("\n")
					first = false
				}
				p.printf("%s: ", name)
				p.print(value)
				p.printf("\n")
			}
		}
		p.indent--
		p.printf("}")

	default:
		v := x.Interface()
		switch v := v.(type) {
		case string:
			// print strings in quotes
			p.printf("%q", v)
			return
		case token.Pos:
			if p.filename != "" && p.source != "" {
				p.printf("%s", v.Position(p.filename, p.source))
				return
			}
		}
		// default
		p.printf("%v", v)
	}
}

func (p *File) JSONString() string {
	file := *p

	if len(file.Source) > 8 {
		file.Source = file.Source[:8] + "..."
	}
	d, _ := json.MarshalIndent(&file, "", "    ")
	return string(d)
}

func (p *File) String() string {
	var buf bytes.Buffer
	Fprint(&buf, p.Filename, p.Source, p)
	return buf.String()
}
