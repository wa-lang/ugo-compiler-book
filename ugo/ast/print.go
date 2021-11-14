package ast

import (
	"bytes"
	goast "go/ast"
	gotoken "go/token"
	"io"
	"os"
)

// Print 打印语法树到 stdout
func Print(node Node) {
	Fprint(os.Stdout, node)
}

// Fprint 打印语法树到指定目标
func Fprint(w io.Writer, node Node) {
	fset := gotoken.NewFileSet()

	if f, _ := node.(*File); f != nil {
		file := *f
		if len(file.Data) > 0 {
			fset.AddFile(f.Name, 1, len(f.Data)).SetLinesForContent(f.Data)
			file.Data = nil
		}
		node = &file
	}

	goast.Fprint(w, fset, node, goast.NotNilFilter)
}

func (p *File) String() string {
	var buf bytes.Buffer
	Fprint(&buf, p)
	return buf.String()
}
