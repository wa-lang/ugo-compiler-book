//go:build ignore

package main

import (
	"fmt"
	"os"

	lexpkg "github.com/wa-lang/ugo/lexer"
)

func main() {
	code := loadCode("../hello.ugo")
	lexer := lexpkg.NewLexer("../hello.ugo", code)

	for i, tok := range lexer.Tokens() {
		fmt.Printf(
			"%02d: %-12v: %-20q // %s\n",
			i, tok.Type, tok.Literal,
			lexpkg.PosString("../hello.ugo", code, tok.Pos),
		)
	}

	fmt.Println("----")

	for i, tok := range lexer.Comments() {
		fmt.Printf(
			"%02d: %-12v: %-20q // %s\n",
			i, tok.Type, tok.Literal,
			lexpkg.PosString("../hello.ugo", code, tok.Pos),
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
