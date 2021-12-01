//go:build ignore

package main

import (
	"fmt"
	"os"

	"github.com/chai2010/ugo/lexer"
)

func main() {
	code := loadCode("../hello.ugo")
	tokens := lexer.Lex("../hello.ugo", code)
	for i, tok := range tokens {
		fmt.Printf(
			"%02d: %-12v: %-16q // %s\n",
			i, tok.Type, tok.Literal,
			lexer.PosString("../hello.ugo", code, tok.Pos),
		)
	}
}

func loadCode(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(data)
}
