package parser

import (
	"fmt"
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

	switch tok := p.PeekToken(); tok.Type {
	case token.IDENT: // call
		return p.parseExpr_call()
	case token.NUMBER:
		tokNumber := p.MustAcceptToken(token.NUMBER)
		value, _ := strconv.Atoi(tokNumber.Literal)
		return &ast.Number{
			Value: value,
		}
	default:
		s := fmt.Sprint(tok)
		panic(s)
	}
}

func (p *Parser) parseExpr_call() *ast.CallExpr {
	tokIdent := p.MustAcceptToken(token.IDENT)
	p.MustAcceptToken(token.LPAREN)
	arg0 := p.parseExpr()
	p.MustAcceptToken(token.RPAREN)

	return &ast.CallExpr{
		FuncName: tokIdent.Literal,
		Args:     []ast.Expr{arg0},
	}
}
