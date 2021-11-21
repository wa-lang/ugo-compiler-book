package token

import "fmt"

type Reader struct {
	comments []Token
	tokens   []Token
	pos      int
	width    int
}

func NewReader(input []Token) *Reader {
	var p = new(Reader)
	for _, tok := range input {
		if tok.Type == COMMENT {
			p.comments = append(p.comments, tok)
		} else {
			p.tokens = append(p.tokens, tok)
		}
	}
	return p
}

func (p *Reader) PeekToken() Token {
	tok := p.ReadToken()
	p.UnreadToken()
	return tok
}

func (p *Reader) ReadToken() Token {
	if p.pos >= len(p.tokens) {
		p.width = 0
		return Token{Type: EOF}
	}
	tok := p.tokens[p.pos]
	p.width = 1
	p.pos += p.width
	return tok
}

func (p *Reader) UnreadToken() {
	p.pos -= p.width
	return
}

func (p *Reader) AcceptToken(expectTypes ...TokenType) (tok Token, ok bool) {
	tok = p.ReadToken()
	for _, x := range expectTypes {
		if tok.Type == x {
			return tok, true
		}
	}
	p.UnreadToken()
	return tok, false
}

func (p *Reader) AcceptTokenList(expectTypes ...TokenType) (toks []Token, ok bool) {
	for {
		tok, ok := p.AcceptToken(expectTypes...)
		if !ok || tok.Type == EOF {
			return toks, len(toks) != 0
		}
		toks = append(toks, tok)
	}
}

func (p *Reader) MustAcceptToken(expectTypes ...TokenType) (tok Token) {
	tok, ok := p.AcceptToken(expectTypes...)
	if !ok {
		panic(fmt.Errorf("token.Reader.MustAcceptToken(%v) failed", expectTypes))
	}
	return tok
}

func (p *Reader) MustAcceptTokenList(expectTypes ...TokenType) (toks []Token) {
	toks, ok := p.AcceptTokenList(expectTypes...)
	if !ok {
		panic(fmt.Errorf("token.Reader.AcceptTokenList(%v) failed", expectTypes))
	}
	return toks
}

func (p *Reader) Comments() []Token {
	return p.comments
}
