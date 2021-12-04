package ast

func Walk(node Node, fn func(node Node) bool) {
	walk(node, fn)
}

func walk(n Node, fn func(n Node) bool) {
	if n == nil || !fn(n) {
		return
	}

	switch n := n.(type) {
	case *File:
		walk(n.Pkg, fn)
		for _, f := range n.Funcs {
			walk(f, fn)
		}

	case *Package:
		walk(n, fn)

	case *Func:
		walk(n, fn)

	case *BlockStmt:
		walk(n, fn)
		for _, stmt := range n.List {
			walk(stmt, fn)
		}

	case *ExprStmt:
		walk(n.X, fn)

	case *BinaryExpr:
		walk(n.X, fn)
		walk(n.Y, fn)

	case *UnaryExpr:
		walk(n.X, fn)

	case *ParenExpr:
		walk(n.X, fn)

	case *CallExpr:
		for _, arg := range n.Args {
			walk(arg, fn)
		}

	case *Ident:
		return
	case *Number:
		return

	default:
		panic("unreachable")
	}
}
