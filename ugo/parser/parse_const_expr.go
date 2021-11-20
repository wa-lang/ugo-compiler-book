package parser

import (
	"fmt"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

func (p *parser) parseConstExpr() ast.Expr {
	return p.parseConstExpr_binary(token.LowestPrec + 1)
}

func (p *parser) parseConstExpr_binary(prec1 int) ast.Expr {
	x := p.parseConstExpr_unary()
	for {
		op := p.r.PeekToken()
		if op.Type.Precedence() < prec1 {
			return x
		}

		p.r.MustAcceptToken(op.Type)
		y := p.parseConstExpr_binary(op.Type.Precedence() + 1)
		x = &ast.BinaryExpr{Op: op, X: x, Y: y}
	}
}

func (p *parser) parseConstExpr_unary() ast.Expr {
	if _, ok := p.r.AcceptToken(token.ADD); ok {
		return p.parseConstExpr_primary()
	}
	if _, ok := p.r.AcceptToken(token.SUB); ok {
		return &ast.UnaryExpr{
			X: p.parseConstExpr_primary(),
		}
	}
	return p.parseConstExpr_primary()
}

func (p *parser) parseConstExpr_primary() ast.Expr {
	peek := p.r.PeekToken()

	switch peek.Type {
	case token.IDENT:
		tokIdent := p.r.MustAcceptToken(token.IDENT)
		// const x = int(1.0) + 1
		if lparen, ok := p.r.AcceptToken(token.LPAREN); ok {
			expr := p.parseConstExpr()
			rparen := p.r.MustAcceptToken(token.RPAREN)
			return &ast.CallExpr{
				Fun: &ast.Ident{
					NamePos: tokIdent.Pos,
					Name:    tokIdent.IdentName(),
				},
				Lparen: lparen.Pos,
				Args:   []ast.Expr{expr},
				Rparen: rparen.Pos,
			}
		}
		return &ast.Ident{
			NamePos: tokIdent.Pos,
			Name:    tokIdent.IdentName(),
		}

	case token.INT:
		tokInt := p.r.MustAcceptToken(token.INT)
		return &ast.Number{
			ValuePos: tokInt.Pos,
			Value:    tokInt.IntValue(),
			ValueEnd: tokInt.EndPos(),
		}
	case token.FLOAT:
		tokFloat := p.r.MustAcceptToken(token.FLOAT)
		return &ast.Number{
			ValuePos: tokFloat.Pos,
			Value:    tokFloat.FloatValue(),
			ValueEnd: tokFloat.EndPos(),
		}

	case token.CHAR:
		panic("TODO")
	case token.STRING:
		panic("TODO")

	case token.LPAREN:
		p.r.MustAcceptToken(token.LPAREN)
		defer p.r.MustAcceptToken(token.RPAREN)
		return p.parseConstExpr()

	default:
		panic(fmt.Errorf("expr: %v", p.r.PeekToken()))
	}
}
