//go:build ignore

package main

import (
	"fmt"
	"os"

	"github.com/wa-lang/ugo/parser"
)

func main() {
	code := loadCode("../hello.ugo")
	f, err := parser.ParseFile("../hello.ugo", code)
	if err != nil {
		panic(err)
	}

	fmt.Println(f.JSONString())
	fmt.Println(f.String())
}

func loadCode(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(data)
}
