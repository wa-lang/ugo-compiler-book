package main

import (
	"bytes"
	"fmt"
	"io"
)

type Compiler struct {
	nextId int
}

func (p *Compiler) GenLLIR(node *ExprNode) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "define i32 @main() {\n")
	fmt.Fprintf(&buf, "\tret i32 %s\n", p.genValue(&buf, node))
	fmt.Fprintf(&buf, "}\n")

	return buf.String()
}

func (p *Compiler) genValue(w io.Writer, node *ExprNode) (id string) {
	if node == nil {
		return ""
	}
	id = p.genId()
	switch node.Value {
	case "+":
		fmt.Fprintf(w, "\t%s = add i32 %s, %s\n",
			id, p.genValue(w, node.Left), p.genValue(w, node.Right),
		)
	case "-":
		fmt.Fprintf(w, "\t%s = sub i32 %s, %s\n",
			id, p.genValue(w, node.Left), p.genValue(w, node.Right),
		)
	case "*":
		fmt.Fprintf(w, "\t%s = mul i32 %s, %s\n",
			id, p.genValue(w, node.Left), p.genValue(w, node.Right),
		)
	case "/":
		fmt.Fprintf(w, "\t%s = sdiv i32 %s, %s\n",
			id, p.genValue(w, node.Left), p.genValue(w, node.Right),
		)
	default:
		fmt.Fprintf(w, "\t%[1]s = add i32 0, %[2]s; %[1]s = %[2]s\n",
			id, node.Value,
		)
	}
	return
}

func (p *Compiler) genId() string {
	id := fmt.Sprintf("%%t%d", p.nextId)
	p.nextId++
	return id
}
