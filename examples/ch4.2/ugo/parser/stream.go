package parser

import (
	"fmt"

	"github.com/wa-lang/ugo/token"
)

type TokenStream struct {
	filename string
	src      string
	tokens   []token.Token
	comments []token.Token
	pos      int
	width    int
}

func NewTokenStream(filename, src string, tokens, comments []token.Token) *TokenStream {
	return &TokenStream{
		filename: filename,
		src:      src,
		tokens:   tokens,
		comments: comments,
	}
}

func (p *TokenStream) Filename() string {
	return p.filename
}

func (p *TokenStream) Source() string {
	return p.src
}

func (p *TokenStream) Tokens() []token.Token {
	return p.tokens
}

func (p *TokenStream) Comments() []token.Token {
	return p.comments
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
		panic(fmt.Errorf("parser.TokenStream.MustAcceptToken(%v) failed", expectTypes))
	}
	return tok
}

func (p *TokenStream) MustAcceptTokenList(expectTypes ...token.TokenType) (toks []token.Token) {
	toks, ok := p.AcceptTokenList(expectTypes...)
	if !ok {
		panic(fmt.Errorf("parser.TokenStream.AcceptTokenList(%v) failed", expectTypes))
	}
	return toks
}

func (p *TokenStream) PrintTokens() {
	for i, tok := range p.Tokens() {
		fmt.Printf(
			"%02d: %-12v: %-20q // %s\n",
			i, tok.Type, tok.Literal,
			tok.Pos.Position(p.filename, p.src),
		)
	}

	fmt.Println("----")

	for i, tok := range p.Comments() {
		fmt.Printf(
			"%02d: %-12v: %-20q // %s\n",
			i, tok.Type, tok.Literal,
			tok.Pos.Position(p.filename, p.src),
		)
	}
}
