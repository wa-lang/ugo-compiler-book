package parser

import (
	"strconv"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

func (p *Parser) parseExpr() ast.Expr {
	return p.parseExpr_binary(1)
}

func (p *Parser) parseExpr_binary(prec int) ast.Expr {
	x := p.parseExpr_unary()
	for {
		op := p.PeekToken()
		if op.Type.Precedence() < prec {
			return x
		}

		p.MustAcceptToken(op.Type)
		y := p.parseExpr_binary(op.Type.Precedence() + 1)
		x = &ast.BinaryExpr{Op: op, X: x, Y: y}
	}
	return nil
}

func (p *Parser) parseExpr_unary() ast.Expr {
	if _, ok := p.AcceptToken(token.ADD); ok {
		return p.parseExpr_primary()
	}
	if _, ok := p.AcceptToken(token.SUB); ok {
		return &ast.UnaryExpr{
			Op: token.Token{Type: token.SUB},
			X:  p.parseExpr_primary(),
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

	tokNumber := p.MustAcceptToken(token.NUMBER)
	value, _ := strconv.Atoi(tokNumber.Literal)
	return &ast.Number{
		Value: value,
	}
}
