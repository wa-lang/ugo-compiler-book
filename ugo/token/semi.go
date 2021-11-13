package token

// https://golang.google.cn/ref/spec#Semicolons
func (tok Token) IsShouldInsertSemi() bool {
	// an identifier
	// an integer, floating-point, imaginary, rune, or string literal
	// one of the keywords break, continue, fallthrough, or return
	// one of the operators and punctuation ++, --, ), ], or }

	if tok.IsLiteral() {
		return true
	}

	switch tok {
	case BREAK, CONTINUE, RETURN:
		return true
	case INC, DEC:
		return true
	case RPAREN, RBRACK, RBRACE:
		return true
	}

	return false
}
