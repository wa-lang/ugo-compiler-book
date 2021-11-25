package compiler

import (
	"bytes"
	"fmt"
	"io"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/builtin"
)

type Compiler struct{}

func (p *Compiler) Compile(f *ast.File) string {
	var buf bytes.Buffer

	p.genBuiltin(&buf)
	p.compileFile(&buf, f)
	p.genMainMain(&buf)

	return buf.String()
}

func (p *Compiler) genBuiltin(w io.Writer) {
	fmt.Fprintln(w, builtin.Header)
}
func (p *Compiler) genMainMain(w io.Writer) {
	fmt.Fprintln(w, builtin.MainMain)
}

func (p *Compiler) compileFile(w io.Writer, file *ast.File) {
	fmt.Fprintln(w, "; TODO")
}

func (p *Compiler) compileFunc(w io.Writer, fn *ast.Func)   {}
func (p *Compiler) compileStmt(w io.Writer, stmt ast.Stmt)  {}
func (p *Compiler) compileExpr(w io.Writer, expr *ast.Expr) {}
