package main

import (
	"fmt"
	"os"

	"github.com/chai2010/ugo/compiler"
	"github.com/chai2010/ugo/parser"
)

func main() {
	code := loadCode("./hello.ugo")
	f, err := parser.ParseFile("./hello.ugo", code)
	if err != nil {
		panic(err)
	}

	ll := new(compiler.Compiler).Compile(f)
	fmt.Print(ll)
}

func loadCode(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(data)
}
