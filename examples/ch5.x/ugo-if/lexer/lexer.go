package lexer

import (
	"fmt"

	"github.com/chai2010/ugo/token"
)

func Lex(name, input string) (tokens, comments []token.Token) {
	l := NewLexer(name, input)
	tokens = l.Tokens()
	comments = l.Comments()
	return
}

type Lexer struct {
	src      *SourceStream
	tokens   []token.Token
	comments []token.Token
}

func NewLexer(name, input string) *Lexer {
	return &Lexer{
		src: NewSourceStream(name, input),
	}
}

func (p *Lexer) Tokens() []token.Token {
	if len(p.tokens) == 0 {
		p.run()
	}
	return p.tokens
}

func (p *Lexer) Comments() []token.Token {
	if len(p.tokens) == 0 {
		p.run()
	}
	return p.comments
}

func (p *Lexer) emit(typ token.TokenType) {
	lit, pos := p.src.EmitToken()
	if typ == token.EOF {
		pos = p.src.Pos()
		lit = ""
	}
	if typ == token.IDENT {
		typ = token.Lookup(lit)
	}
	p.tokens = append(p.tokens, token.Token{
		Type:    typ,
		Literal: lit,
		Pos:     token.Pos(pos + 1),
	})
}

func (p *Lexer) emitComment() {
	lit, pos := p.src.EmitToken()

	p.comments = append(p.comments, token.Token{
		Type:    token.COMMENT,
		Literal: lit,
		Pos:     token.Pos(pos + 1),
	})
}

func (p *Lexer) errorf(format string, args ...interface{}) {
	tok := token.Token{
		Type:    token.ERROR,
		Literal: fmt.Sprintf(format, args...),
		Pos:     token.Pos(p.src.Pos() + 1),
	}
	p.tokens = append(p.tokens, tok)
	panic(tok)
}

func (p *Lexer) run() (tokens []token.Token) {
	defer func() {
		tokens = p.tokens
		if r := recover(); r != nil {
			if _, ok := r.(token.Token); !ok {
				panic(r)
			}
		}
	}()

	for {
		r := p.src.Read()
		if r == rune(token.EOF) {
			p.emit(token.EOF)
			return
		}

		switch {
		case r == '\n':
			p.src.IgnoreToken()

			if len(p.tokens) > 0 {
				switch p.tokens[len(p.tokens)-1].Type {
				case token.RPAREN, token.IDENT, token.NUMBER:
					p.emit(token.SEMICOLON)
				}
			}

		case isSpace(r):
			p.src.IgnoreToken()

		case isAlpha(r):
			p.src.Unread()
			for {
				if r := p.src.Read(); !isAlphaNumeric(r) {
					p.src.Unread()
					p.emit(token.IDENT)
					break
				}
			}

		case ('0' <= r && r <= '9'): // 123, 1.0
			p.src.Unread()

			digits := "0123456789"
			p.src.AcceptRun(digits)
			p.emit(token.NUMBER)

		case r == '+': // +, +=, ++
			p.emit(token.ADD)
		case r == '-': // -, -=, --
			p.emit(token.SUB)
		case r == '*': // *, *=
			p.emit(token.MUL)
		case r == '/': // /, //, /*, /=
			if p.src.Peek() != '/' {
				p.emit(token.DIV)
			} else {
				// line comment
				for {
					t := p.src.Read()
					if t == '\n' {
						p.src.Unread()
						p.emitComment()
						break
					}
					if t == rune(token.EOF) {
						p.emitComment()
						return
					}
				}
			}
		case r == '%': // %
			p.emit(token.MOD)

		case r == '=': // =, ==
			switch p.src.Read() {
			case '=':
				p.emit(token.EQL)
			default:
				p.src.Unread()
				p.emit(token.ASSIGN)
			}

		case r == '!': // !=
			switch p.src.Read() {
			case '=':
				p.emit(token.NEQ)
			default:
				p.errorf("unrecognized character: %#U", r)
			}

		case r == '<': // <, <=
			switch p.src.Read() {
			case '=':
				p.emit(token.LEQ)
			default:
				p.src.Unread()
				p.emit(token.LSS)
			}

		case r == '>': // >, >=
			switch p.src.Read() {
			case '=':
				p.emit(token.GEQ)
			default:
				p.src.Unread()
				p.emit(token.GTR)
			}

		case r == ':': // :, :=
			switch p.src.Read() {
			case '=':
				p.emit(token.DEFINE)
			default:
				p.errorf("unrecognized character: %#U", r)
			}

		case r == '(':
			p.emit(token.LPAREN)
		case r == '{':
			p.emit(token.LBRACE)

		case r == ')':
			p.emit(token.RPAREN)
		case r == '}':
			p.emit(token.RBRACE)

		case r == ',':
			p.emit(token.COMMA)
		case r == ';':
			p.emit(token.SEMICOLON)

		default:
			p.errorf("unrecognized character: %#U", r)
			return
		}
	}
}
