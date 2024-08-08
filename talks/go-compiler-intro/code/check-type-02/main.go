package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	info := &types.Info{
		Types:  make(map[ast.Expr]types.TypeAndValue), // HL
		Defs:   make(map[*ast.Ident]types.Object),     // HL
		Uses:   make(map[*ast.Ident]types.Object),     // HL
		Scopes: make(map[ast.Node]*types.Scope),       // HL
	}
	conf := types.Config{Importer: nil}
	pkg, err := conf.Check("hello.go", fset, []*ast.File{f}, info) // HL
	if err != nil {
		log.Fatal(err)
	}
	_ = pkg
}

const src = `
package main
var s = "hello ssa"
func main() {
	for i := 0; i < 3; i++ {
		println(s)
	}
}
`
