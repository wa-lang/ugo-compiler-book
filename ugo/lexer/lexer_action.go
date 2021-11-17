package lexer

import (
	"github.com/chai2010/ugo/token"
)

func (l *lexer) run() {
	defer l.adjustItems()

	for {
		switch r := l.next(); {
		case r == eof:
			l.emit(token.EOF)
			return

		case r == '\n':
			if len(l.items) > 0 {
				lastTok := l.items[len(l.items)-1].Type
				if lastTok.IsShouldAppendSemi() && !l.opt.DontInsertSemi {
					l.emit(token.SEMICOLON)
				}
			}

		case isSpace(r):
			l.ignore()

		case isAlpha(r):
			l.backup()
			l.lexIdentifier()

		case ('0' <= r && r <= '9'): // 123, 1.0
			l.backup()
			l.lexNumber()

		case r == '\'': // 'a', '\n'
			l.lexChar()

		case r == '"': // "abc\n"
			l.lexQuote()

		case r == '`': // `abc`
			l.lexRawQuote()

		case r == '+': // +, +=, ++
			switch l.next() {
			case '=':
				l.emit(token.ADD_ASSIGN)
			case '+':
				l.emit(token.INC)
			default:
				l.backup()
				l.emit(token.ADD)
			}
		case r == '-': // -, -=, --
			switch l.next() {
			case '=':
				l.emit(token.SUB_ASSIGN)
			case '-':
				l.emit(token.DEC)
			default:
				l.backup()
				l.emit(token.SUB)
			}
		case r == '*': // *, *=
			switch l.next() {
			case '=':
				l.emit(token.MUL_ASSIGN)
			default:
				l.backup()
				l.emit(token.MUL)
			}
		case r == '/': // /, //, /*, /=
			switch l.next() {
			case '=':
				l.emit(token.QUO_ASSIGN)
			case '/':
				l.lexLineComment()
			case '*':
				l.lexMultiLineComment()
			default:
				l.backup()
				l.emit(token.QUO)
			}
		case r == '%': // %, %=
			switch l.next() {
			case '=':
				l.emit(token.REM_ASSIGN)
			default:
				l.backup()
				l.emit(token.REM)
			}

		case r == '&': // &, &=, &^, &^=, &&
			switch l.next() {
			case '=':
				l.emit(token.AND_ASSIGN)
			case '^':
				l.next()
				if l.peek() == '=' {
					l.emit(token.AND_NOT_ASSIGN)
				} else {
					l.emit(token.AND_NOT)
				}
			case '&':
				l.emit(token.LAND)
			default:
				l.emit(token.AND)
			}
		case r == '|': // |, |=, ||
			switch l.next() {
			case '=':
				l.emit(token.OR_ASSIGN)
			case '|':
				l.emit(token.LOR)
			default:
				l.backup()
				l.emit(token.OR)
			}
		case r == '^': // ^, ^=
			switch l.next() {
			case '=':
				l.emit(token.XOR_ASSIGN)
			default:
				l.backup()
				l.emit(token.XOR)
			}
		case r == '<': // <, <=, <<, <<=
			switch l.next() {
			case '=':
				l.emit(token.LSS)
			case '<':
				switch l.next() {
				case '=':
					l.emit(token.SHL_ASSIGN)
				default:
					l.backup()
					l.emit(token.SHL)
				}
			default:
				l.backup()
				l.emit(token.LSS)
			}
		case r == '>': // >, >=, >>, >>=
			switch l.next() {
			case '=':
				l.emit(token.GEQ)
			case '>':
				switch l.next() {
				case '=':
					l.emit(token.SHR_ASSIGN)
				default:
					l.backup()
					l.emit(token.SHR)
				}
			default:
				l.backup()
				l.emit(token.GTR)
			}

		case r == '=': // =, ==
			switch l.next() {
			case '=':
				l.emit(token.EQL)
			default:
				l.backup()
				l.emit(token.ASSIGN)
			}
		case r == '!': // !, !=
			switch l.next() {
			case '=':
				l.emit(token.NEQ)
			default:
				l.backup()
				l.emit(token.NOT)
			}
		case r == ':': // :, :=
			switch l.next() {
			case '=':
				l.emit(token.DEFINE)
			default:
				l.backup()
				l.emit(token.COLON)
			}

		case r == '(':
			l.emit(token.LPAREN)
		case r == '[':
			l.emit(token.LBRACK)
		case r == '{':
			l.emit(token.LBRACE)

		case r == ')':
			l.emit(token.RPAREN)
		case r == ']':
			l.emit(token.RBRACK)
		case r == '}':
			l.emit(token.RBRACE)

		case r == ',':
			l.emit(token.COMMA)
		case r == '.':
			l.emit(token.PERIOD)
		case r == ';':
			l.emit(token.SEMICOLON)

		default:
			l.errorf("unrecognized character in action: %#U", r)
			return
		}
	}
}

func (l *lexer) adjustItems() {
	if l.opt.SkipComment {
		items := l.items[:0]
		for _, t := range l.items {
			if t.Type != token.COMMENT {
				items = append(items, t)
			}
		}
		l.items = items
	}
}

func (l *lexer) lexLineComment() {
	for {
		if r := l.next(); r == '\n' || r == eof {
			l.emit(token.COMMENT)
			return
		}
	}
}

func (l *lexer) lexMultiLineComment() {
	for {
		if r := l.next(); r == '*' && l.peek() == '/' {
			l.next()
			l.emit(token.COMMENT)
			return
		}
	}
}

func (l *lexer) lexIdentifier() {
	l.start = l.pos
	for {
		if r := l.next(); !isAlphaNumeric(r) {
			l.backup()
			l.emit(token.IDENT)
			return
		}
	}
}

func (l *lexer) lexNumber() {
	isFloat := false
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		isFloat = true
		l.acceptRun(digits)
	}

	if isFloat {
		l.emit(token.FLOAT)
	} else {
		l.emit(token.INT)
	}
}

func (l *lexer) lexChar() {
	for {
		switch l.next() {
		case '\\':
			l.next()
		case eof:
			l.errorf("unterminated quoted string")
			return
		case '\'':
			l.emit(token.CHAR)
			return
		}
	}
}

func (l *lexer) lexQuote() {
	for {
		switch l.next() {
		case '\\':
			l.next()
		case eof:
			l.errorf("unterminated quoted string")
			return
		case '"':
			l.emit(token.STRING)
			return
		}
	}
}

func (l *lexer) lexRawQuote() {
	for {
		switch l.next() {
		case eof:
			l.errorf("unterminated raw quoted string")
			return
		case '`':
			l.emit(token.STRING)
			return
		}
	}
}
