//go:build ignore

package main

import (
	"fmt"
	gotoken "go/token"
	"os"

	"github.com/chai2010/ugo/lexer"
)

func main() {
	code := loadCode("../hello.ugo")
	tokens := lexer.Lex("../hello.ugo", code, lexer.Option{})
	for i, tok := range tokens {
		fmt.Printf(
			"%02d: %-12v: %-16q // %s\n",
			i, tok.Type, tok.Literal,
			PosString("../hello.ugo", code, tok.Pos),
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

func PosString(filename string, src string, pos int) string {
	fset := gotoken.NewFileSet()
	fset.AddFile(filename, 1, len(src)).SetLinesForContent([]byte(src))
	return fmt.Sprintf("%v", fset.Position(gotoken.Pos(pos+1)))
}
