package main

import (
	"os"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/compiler"
	"github.com/chai2010/ugo/parser"
)

func main() {
	code := loadCode("./hello.ugo")
	f, err := parser.ParseFile("./hello.ugo", code)
	if err != nil {
		panic(err)
	}

	ast.Print(f)

	ll := new(compiler.Compiler).Compile(f)
	//fmt.Print(ll)
	_ = ll
}

func loadCode(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(data)
}
