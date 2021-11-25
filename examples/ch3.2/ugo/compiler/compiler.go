package compiler

import (
	"bytes"
	"fmt"
	"io"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/builtin"
)

type Compiler struct{}

func (p *Compiler) Compile(file *ast.File) string {
	var buf bytes.Buffer

	p.genHeader(&buf, file)
	p.compileFile(&buf, file)
	p.genMain(&buf, file)

	return buf.String()
}

func (p *Compiler) genHeader(w io.Writer, file *ast.File) {
	fmt.Fprintf(w, "; package %s\n", file.Pkg.Name)
	fmt.Fprintln(w, builtin.Header)
}

func (p *Compiler) genMain(w io.Writer, file *ast.File) {
	if file.Pkg.Name != "main" {
		return
	}
	for _, fn := range file.Funcs {
		if fn.Name == "main" {
			fmt.Fprintln(w, builtin.MainMain)
			return
		}
	}
}

func (p *Compiler) compileFile(w io.Writer, file *ast.File) {
	for _, fn := range file.Funcs {
		p.compileFunc(w, file, fn)
	}
}

func (p *Compiler) compileFunc(w io.Writer, file *ast.File, fn *ast.Func) {
	if fn.Body == nil {
		fmt.Fprintf(w, "declare i32 @ugo_%s_%s() {\n", file.Pkg.Name, fn.Name)
		return
	}

	fmt.Fprintf(w, "define i32 @ugo_%s_%s() {\n", file.Pkg.Name, fn.Name)
	p.compileStmt(w, fn.Body)

	fmt.Fprintln(w, "\tret i32 0")
	fmt.Fprintln(w, "}")
}

func (p *Compiler) compileStmt(w io.Writer, stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.BlockStmt:
		for _, x := range stmt.List {
			p.compileStmt(w, x)
		}
	case *ast.ExprStmt:
		p.compileExpr(w, stmt.X)

	default:
		panic("unreachable")
	}
}

func (p *Compiler) compileExpr(w io.Writer, expr ast.Expr) {
	// TODO
}
