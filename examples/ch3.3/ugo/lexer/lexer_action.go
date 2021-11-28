package lexer

import (
	"github.com/chai2010/ugo/token"
)

func (l *lexer) run() {
	for {
		switch r := l.next(); {
		case r == eof:
			l.emit(token.EOF)
			return

		case r == '\n':
			if len(l.items) > 0 {
				//lastTok := l.items[len(l.items)-1].Type
				//if lastTok.IsShouldAppendSemi() && !l.opt.DontInsertSemi {
				//	l.emit(token.SEMICOLON)
				//}
			}

		case isSpace(r):
			l.ignore()

		case isAlpha(r):
			l.backup()

			l.lexIdentifier()

		case ('0' <= r && r <= '9'): // 123, 1.0
			l.backup()

			digits := "0123456789"
			l.acceptRun(digits)
			l.emit(token.NUMBER)

		case r == '+': // +, +=, ++
			l.emit(token.ADD)
		case r == '-': // -, -=, --
			l.emit(token.SUB)
		case r == '*': // *, *=
			l.emit(token.MUL)
		case r == '/': // /, //, /*, /=
			l.emit(token.DIV)

		case r == '(':
			l.emit(token.LPAREN)
		case r == '{':
			l.emit(token.LBRACE)

		case r == ')':
			l.emit(token.RPAREN)
		case r == '}':
			l.emit(token.RBRACE)

		case r == ';':
			l.emit(token.SEMICOLON)

		default:
			l.errorf("unrecognized character in action: %#U", r)
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
