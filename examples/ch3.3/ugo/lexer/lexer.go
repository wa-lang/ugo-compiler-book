package lexer

import (
	"fmt"

	"github.com/chai2010/ugo/token"
)

const eof = 0

type Option struct {
	SkipComment    bool
	DontInsertSemi bool
}

func Lex(name, input string, opt Option) []token.Token {
	return newLexer(name, input, opt).run()
}

// lexer holds the state of the scanner.
type lexer struct {
	r     Reader
	opt   Option
	name  string
	items []token.Token
}

func newLexer(name, input string, opt Option) *lexer {
	return &lexer{
		r:    NewReader(input),
		opt:  opt,
		name: name,
	}
}

func (p *lexer) emit(typ token.TokenType) {
	lit, pos := p.r.EmitToken()
	if typ == token.EOF {
		pos = p.r.Pos()
		lit = ""
	}
	p.items = append(p.items, token.Token{
		Type:    typ,
		Literal: lit,
		Pos:     pos,
	})
}

func (p *lexer) errorf(format string, args ...interface{}) {
	tok := token.Token{
		Type:    token.ERROR,
		Literal: fmt.Sprintf(format, args...),
		Pos:     p.r.Pos(),
	}
	p.items = append(p.items, tok)
	panic(tok)
}

func (p *lexer) run() (items []token.Token) {
	defer func() {
		items = p.items
		if r := recover(); r != nil {
			if _, ok := r.(token.Token); !ok {
				panic(r)
			}
		}
	}()

	for {
		r := p.r.Read()
		if r == eof {
			p.emit(token.EOF)
			return
		}

		switch {
		case r == '\n':
			p.r.IgnoreToken()

			if len(p.items) > 0 {
				switch p.items[len(p.items)-1].Type {
				case token.RPAREN, token.RBRACE:
					p.emit(token.SEMICOLON)
				}
			}

		case isSpace(r):
			p.r.IgnoreToken()

		case isAlpha(r):
			p.r.Unread()
			for {
				if r := p.r.Read(); !isAlphaNumeric(r) {
					p.r.Unread()
					p.emit(token.IDENT)
					break
				}
			}

		case ('0' <= r && r <= '9'): // 123, 1.0
			p.r.Unread()

			digits := "0123456789"
			p.r.AcceptRun(digits)
			p.emit(token.NUMBER)

		case r == '+': // +, +=, ++
			p.emit(token.ADD)
		case r == '-': // -, -=, --
			p.emit(token.SUB)
		case r == '*': // *, *=
			p.emit(token.MUL)
		case r == '/': // /, //, /*, /=
			if p.r.Peek() == '/' {
				// line comment
				for {
					t := p.r.Read()
					if t == '\n' {
						p.r.IgnoreToken()
						break
					}
					if t == eof {
						return
					}
				}
			} else {
				p.emit(token.DIV)
			}

		case r == '(':
			p.emit(token.LPAREN)
		case r == '{':
			p.emit(token.LBRACE)

		case r == ')':
			p.emit(token.RPAREN)
		case r == '}':
			p.emit(token.RBRACE)

		case r == ';':
			p.emit(token.SEMICOLON)

		default:
			p.errorf("unrecognized character: %#U", r)
			return
		}
	}
}
