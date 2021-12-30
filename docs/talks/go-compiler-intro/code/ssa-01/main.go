package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"

	"golang.org/x/tools/go/ssa"
)

func main() {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)

	info := &types.Info{}
	conf := types.Config{Importer: nil}
	pkg, _ := conf.Check("hello.go", fset, []*ast.File{f}, info) // HL

	var ssaProg = ssa.NewProgram(fset, ssa.SanityCheckFunctions)        // HL
	var ssaPkg = ssaProg.CreatePackage(pkg, []*ast.File{f}, info, true) // HL

	ssaPkg.Build() // HL

	ssaPkg.WriteTo(os.Stdout)
	ssaPkg.Func("init").WriteTo(os.Stdout) // HL
	ssaPkg.Func("main").WriteTo(os.Stdout) // HL
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
