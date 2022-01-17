package lexer

import (
	"fmt"

	"github.com/wa-lang/ugo/token"
)

type Lexer struct {
	*SourceStream
	tokens   []token.Token
	comments []token.Token
}

func Lex(name, input string) []token.Token {
	return NewLexer(name, input).Tokens()
}

func NewLexer(name, input string) *Lexer {
	return &Lexer{
		SourceStream: NewSourceStream(name, input),
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
	lit, pos := p.EmitToken()
	if typ == token.EOF {
		pos = p.Pos()
		lit = ""
	}
	if typ == token.IDENT {
		typ = token.Lookup(lit)
	}
	p.tokens = append(p.tokens, token.Token{
		Type:    typ,
		Literal: lit,
		Pos:     pos,
	})
}

func (p *Lexer) emitComment() {
	lit, pos := p.EmitToken()

	p.comments = append(p.comments, token.Token{
		Type:    token.COMMENT,
		Literal: lit,
		Pos:     pos,
	})
}

func (p *Lexer) errorf(format string, args ...interface{}) {
	tok := token.Token{
		Type:    token.ERROR,
		Literal: fmt.Sprintf(format, args...),
		Pos:     p.Pos(),
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
		r := p.Read()
		if r == rune(token.EOF) {
			p.emit(token.EOF)
			return
		}

		switch {
		case r == '\n':
			p.IgnoreToken()

			if len(p.tokens) > 0 {
				switch p.tokens[len(p.tokens)-1].Type {
				case token.RPAREN, token.IDENT, token.NUMBER:
					p.emit(token.SEMICOLON)
				}
			}

		case isSpace(r):
			p.IgnoreToken()

		case isAlpha(r):
			p.Unread()
			for {
				if r := p.Read(); !isAlphaNumeric(r) {
					p.Unread()
					p.emit(token.IDENT)
					break
				}
			}

		case ('0' <= r && r <= '9'): // 123, 1.0
			p.Unread()

			digits := "0123456789"
			p.AcceptRun(digits)
			p.emit(token.NUMBER)

		case r == '+': // +, +=, ++
			p.emit(token.ADD)
		case r == '-': // -, -=, --
			p.emit(token.SUB)
		case r == '*': // *, *=
			p.emit(token.MUL)
		case r == '/': // /, //, /*, /=
			if p.Peek() != '/' {
				p.emit(token.DIV)
			} else {
				// line comment
				for {
					t := p.Read()
					if t == '\n' {
						p.emitComment()
						break
					}
					if t == rune(token.EOF) {
						p.Unread()
						p.emitComment()
						return
					}
				}
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
