package main

import "fmt"

type TokenReader struct {
	tokens []Token
	pos    int
	width  int
}

func NewTokenReader(input []Token) *TokenReader {
	return &TokenReader{tokens: input}
}

func (p *TokenReader) PeekToken() Token {
	tok := p.ReadToken()
	p.UnreadToken()
	return tok
}

func (p *TokenReader) ReadToken() Token {
	if p.pos >= len(p.tokens) {
		p.width = 0
		return Token{Type: EOF}
	}
	tok := p.tokens[p.pos]
	p.width = 1
	p.pos += p.width
	return tok
}

func (p *TokenReader) UnreadToken() {
	p.pos -= p.width
	return
}

func (p *TokenReader) AcceptToken(expectTypes ...TokenType) (tok Token, ok bool) {
	tok = p.ReadToken()
	for _, x := range expectTypes {
		if tok.Type == x {
			return tok, true
		}
	}
	p.UnreadToken()
	return tok, false
}

func (p *TokenReader) MustAcceptToken(expectTypes ...TokenType) (tok Token) {
	tok, ok := p.AcceptToken(expectTypes...)
	if !ok {
		panic(fmt.Errorf("token.Reader.MustAcceptToken(%v) failed", expectTypes))
	}
	return tok
}
