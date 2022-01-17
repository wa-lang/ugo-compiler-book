package parser

import (
	"strconv"

	"github.com/wa-lang/ugo/ast"
	"github.com/wa-lang/ugo/token"
)

func (p *Parser) parseExpr() ast.Expr {
	return p.parseExpr_binary(1)
}

func (p *Parser) parseExpr_binary(prec int) ast.Expr {
	x := p.parseExpr_unary()
	for {
		switch tok := p.PeekToken(); tok.Type {
		case token.EOF:
			return x
		case token.SEMICOLON: // ;
			return x
		}

		op := p.PeekToken()
		if op.Type.Precedence() < prec {
			return x
		}

		p.MustAcceptToken(op.Type)
		y := p.parseExpr_binary(op.Type.Precedence() + 1)
		x = &ast.BinaryExpr{OpPos: op.Pos, Op: op.Type, X: x, Y: y}
	}
}

func (p *Parser) parseExpr_unary() ast.Expr {
	if _, ok := p.AcceptToken(token.ADD); ok {
		return p.parseExpr_primary()
	}
	if tok, ok := p.AcceptToken(token.SUB); ok {
		return &ast.UnaryExpr{
			OpPos: tok.Pos,
			Op:    tok.Type,
			X:     p.parseExpr_primary(),
		}
	}
	return p.parseExpr_primary()
}
func (p *Parser) parseExpr_primary() ast.Expr {
	if _, ok := p.AcceptToken(token.LPAREN); ok {
		expr := p.parseExpr()
		p.MustAcceptToken(token.RPAREN)
		return expr
	}

	switch tok := p.PeekToken(); tok.Type {
	case token.IDENT: // call
		return p.parseExpr_call()
	case token.NUMBER:
		tokNumber := p.MustAcceptToken(token.NUMBER)
		value, _ := strconv.Atoi(tokNumber.Literal)
		return &ast.Number{
			ValuePos: tokNumber.Pos,
			ValueEnd: tokNumber.Pos + token.Pos(len(tokNumber.Literal)),
			Value:    value,
		}
	default:
		p.errorf(tok.Pos, "unknown tok: type=%v, lit=%q", tok.Type, tok.Literal)
		panic("unreachable")
	}
}

func (p *Parser) parseExpr_call() *ast.CallExpr {
	tokIdent := p.MustAcceptToken(token.IDENT)
	tokLparen := p.MustAcceptToken(token.LPAREN)
	arg0 := p.parseExpr()
	tokRparen := p.MustAcceptToken(token.RPAREN)

	return &ast.CallExpr{
		FuncName: &ast.Ident{
			NamePos: tokIdent.Pos,
			Name:    tokIdent.Literal,
		},
		Lparen: tokLparen.Pos,
		Args:   []ast.Expr{arg0},
		Rparen: tokRparen.Pos,
	}
}
