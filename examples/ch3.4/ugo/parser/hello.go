//go:build ignore

package main

import (
	"encoding/json"
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

	fmt.Println(JSONString(f))
}

func JSONString(x interface{}) string {
	b, err := json.MarshalIndent(x, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func loadCode(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(data)
}
