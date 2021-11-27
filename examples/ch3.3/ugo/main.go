package main

import (
	"fmt"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/compiler"
	"github.com/chai2010/ugo/token"
)

func main() {
	ll := new(compiler.Compiler).Compile(ugoProg)
	fmt.Print(ll)
}

var ugoProg = &ast.File{
	Pkg: &ast.Package{
		Name: "main",
	},
	Funcs: []*ast.Func{
		{
			Name: "main",
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							FuncName: "exit",
							Args: []ast.Expr{
								&ast.BinaryExpr{
									Op: token.Token{Type: token.ADD},
									X:  &ast.Number{Value: 40},
									Y:  &ast.Number{Value: 2},
								},
							},
						},
					},
				},
			},
		},
	},
}
