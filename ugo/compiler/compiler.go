package compiler

import (
	"bytes"
	"fmt"
	"io"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

type Compiler struct {
	nextId int
}

func (p *Compiler) CompileFile(f *ast.File) string {
	var buf bytes.Buffer

	fmt.Fprintln(&buf, builtin_llir)
	fmt.Fprintln(&buf)

	fmt.Fprintf(&buf, "define i32 @main() {\n")
	fmt.Fprintf(&buf, "\tret i32 0; TODO\n")
	fmt.Fprintf(&buf, "}\n")

	return buf.String()
}

func (p *Compiler) CompileExpr(node ast.Expr) string {
	var buf bytes.Buffer

	fmt.Fprintln(&buf, builtin_llir)
	fmt.Fprintln(&buf)

	fmt.Fprintf(&buf, "define i32 @main() {\n")
	fmt.Fprintf(&buf, "\tret i32 %s\n", p.genValue(&buf, node))
	fmt.Fprintf(&buf, "}\n")

	return buf.String()
}

func (p *Compiler) genValue(w io.Writer, node ast.Expr) (id string) {
	if node == nil {
		return ""
	}
	id = p.genId()
	switch node := node.(type) {
	case *ast.Ident:
		//fmt.Fprintf(w, "\t%[1]s = add i32 0, %[2]v; %[1]s = %[2]v\n",
		//	id, node.Value,
		//)
	case *ast.Number:
		fmt.Fprintf(w, "\t%[1]s = add i32 0, %[2]v; %[1]s = %[2]v\n",
			id, node.Value,
		)
	case *ast.UnaryExpr:
		switch node.Op.Type {
		case token.SUB:
			fmt.Fprintf(w, "\t%s = sub i32 0, %s\n",
				id, p.genValue(w, node.X),
			)
		}
	case *ast.ParenExpr:
		fmt.Fprintf(w, "\t%[1]s = add i32 0, %[2]s; %[1]s = %[2]s\n",
			id, p.genValue(w, node.X),
		)
	case *ast.CallExpr:
		fmt.Fprintf(w, "\t%[1]s = call i32(i32) @%[2]s(i32 %[3]s);\n",
			id, node.Fun.Name, p.genValue(w, node.Args[0]),
		)
	case *ast.BinaryExpr:
		switch node.Op.Type {
		case token.ADD:
			fmt.Fprintf(w, "\t%s = add i32 %s, %s\n",
				id, p.genValue(w, node.X), p.genValue(w, node.Y),
			)
		case token.SUB:
			fmt.Fprintf(w, "\t%s = sub i32 %s, %s\n",
				id, p.genValue(w, node.X), p.genValue(w, node.Y),
			)
		case token.MUL:
			fmt.Fprintf(w, "\t%s = mul i32 %s, %s\n",
				id, p.genValue(w, node.X), p.genValue(w, node.Y),
			)
		case token.QUO:
			fmt.Fprintf(w, "\t%s = div i32 %s, %s\n",
				id, p.genValue(w, node.X), p.genValue(w, node.Y),
			)
		}
	}
	return
}

func (p *Compiler) genId() string {
	id := fmt.Sprintf("%%t%d", p.nextId)
	p.nextId++
	return id
}
