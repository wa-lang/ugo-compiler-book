package parser

import (
	"fmt"

	"github.com/chai2010/ugo/token"
)

type TokenStream struct {
	tokens   []token.Token
	comments []token.Token
	pos      int
	width    int
}

func NewTokenStream(tokens, comments []token.Token) *TokenStream {
	return &TokenStream{
		tokens:   tokens,
		comments: comments,
	}
}

func (p *TokenStream) PeekToken() token.Token {
	tok := p.ReadToken()
	p.UnreadToken()
	return tok
}

func (p *TokenStream) ReadToken() token.Token {
	if p.pos >= len(p.tokens) {
		p.width = 0
		return token.Token{Type: token.EOF}
	}
	tok := p.tokens[p.pos]
	p.width = 1
	p.pos += p.width
	return tok
}

func (p *TokenStream) UnreadToken() {
	p.pos -= p.width
	return
}

func (p *TokenStream) AcceptToken(expectTypes ...token.TokenType) (tok token.Token, ok bool) {
	tok = p.ReadToken()
	for _, x := range expectTypes {
		if tok.Type == x {
			return tok, true
		}
	}
	p.UnreadToken()
	return tok, false
}

func (p *TokenStream) AcceptTokenList(expectTypes ...token.TokenType) (toks []token.Token, ok bool) {
	for {
		tok, ok := p.AcceptToken(expectTypes...)
		if !ok || tok.Type == token.EOF {
			return toks, len(toks) != 0
		}
		toks = append(toks, tok)
	}
}

func (p *TokenStream) MustAcceptToken(expectTypes ...token.TokenType) (tok token.Token) {
	tok, ok := p.AcceptToken(expectTypes...)
	if !ok {
		panic(fmt.Errorf("token.TokenStream.MustAcceptToken(%v) failed", expectTypes))
	}
	return tok
}

func (p *TokenStream) MustAcceptTokenList(expectTypes ...token.TokenType) (toks []token.Token) {
	toks, ok := p.AcceptTokenList(expectTypes...)
	if !ok {
		panic(fmt.Errorf("token.TokenStream.AcceptTokenList(%v) failed", expectTypes))
	}
	return toks
}

func (p *TokenStream) Comments() []token.Token {
	return p.comments
}
