package lexer

import (
	"fmt"
	gotoken "go/token"
	"unicode"
)

func PosString(filename string, src string, pos int) string {
	fset := gotoken.NewFileSet()
	fset.AddFile(filename, 1, len(src)).SetLinesForContent([]byte(src))
	return fmt.Sprintf("%v", fset.Position(gotoken.Pos(pos+1)))
}

func isSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\r':
		return true
	}
	return false
}

func isAlpha(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(rune(r)) || unicode.IsDigit(r)
}
