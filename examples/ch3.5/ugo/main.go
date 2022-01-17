package main

import (
	"os"

	"github.com/wa-lang/ugo/ast"
	"github.com/wa-lang/ugo/compiler"
	"github.com/wa-lang/ugo/parser"
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
