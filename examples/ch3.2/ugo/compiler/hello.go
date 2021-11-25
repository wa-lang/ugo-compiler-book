//go:build ignore

package main

import (
	"fmt"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/compiler"
)

func main() {
	ll := new(compiler.Compiler).Compile(ugoProg)
	fmt.Print(ll)
}

var ugoProg = &ast.File{}
