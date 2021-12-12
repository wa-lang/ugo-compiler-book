package parser

import (
	"fmt"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

func (p *Parser) parseExpr_const() ast.Expr {
	return p.parseExpr_const_binary(1)
}

func (p *Parser) parseExpr_const_binary(prec1 int) ast.Expr {
	x := p.parseExpr_const_unary()
	for {
		op := p.PeekToken()
		if op.Type.Precedence() < prec1 {
			return x
		}

		p.MustAcceptToken(op.Type)
		y := p.parseExpr_const_binary(op.Type.Precedence() + 1)
		x = &ast.BinaryExpr{Op: op.Type, X: x, Y: y}
	}
}

func (p *Parser) parseExpr_const_unary() ast.Expr {
	if _, ok := p.AcceptToken(token.ADD); ok {
		return p.parseExpr_const_primary()
	}
	if _, ok := p.AcceptToken(token.SUB); ok {
		return &ast.UnaryExpr{
			X: p.parseExpr_const_primary(),
		}
	}
	return p.parseExpr_const_primary()
}

func (p *Parser) parseExpr_const_primary() ast.Expr {
	peek := p.PeekToken()

	switch peek.Type {
	case token.IDENT:
		tokIdent := p.MustAcceptToken(token.IDENT)
		// const x = int(1.0) + 1
		if lparen, ok := p.AcceptToken(token.LPAREN); ok {
			expr := p.parseExpr_const()
			rparen := p.MustAcceptToken(token.RPAREN)
			return &ast.CallExpr{
				FuncName: &ast.Ident{
					NamePos: tokIdent.Pos,
					Name:    tokIdent.Literal,
				},
				Lparen: lparen.Pos,
				Args:   []ast.Expr{expr},
				Rparen: rparen.Pos,
			}
		}
		return &ast.Ident{
			NamePos: tokIdent.Pos,
			Name:    tokIdent.Literal,
		}

	case token.NUMBER:
		tokInt := p.MustAcceptToken(token.NUMBER)
		return &ast.Number{
			ValuePos: tokInt.Pos,
			Value:    tokInt.IntValue(),
			ValueEnd: tokInt.Pos + token.Pos(len(tokInt.Literal)),
		}

	case token.LPAREN:
		p.MustAcceptToken(token.LPAREN)
		defer p.MustAcceptToken(token.RPAREN)
		return p.parseExpr_const()

	default:
		panic(fmt.Errorf("expr: %v", p.PeekToken()))
	}
}
