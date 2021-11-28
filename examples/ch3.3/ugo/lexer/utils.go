package lexer

import (
	"unicode"
)

// isSpace reports whether r is a space character.
func isSpace(r int) bool {
	switch r {
	case ' ', '\t', '\r':
		return true
	}
	return false
}

func isAlpha(r int) bool {
	return r == '_' || unicode.IsLetter(rune(r))
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r int) bool {
	return r == '_' || unicode.IsLetter(rune(r)) || unicode.IsDigit(rune(r))
}
