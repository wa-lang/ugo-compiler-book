package parser

import (
	"fmt"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/logger"
	"github.com/chai2010/ugo/token"
)

func (p *parser) parseExpr() ast.Expr {
	logger.Debugln("parseExpr: peek =", p.peek())

	expr := p.parseExpr_mul()
	for {
		switch p.peekTokenType() {
		case token.ADD, token.SUB:
			tok := p.next()
			expr = &ast.BinaryExpr{
				X:  expr,
				Op: tok,
				Y:  p.parseExpr_mul(),
			}
		default:
			return expr
		}
	}
}

func (p *parser) parseExpr_mul() ast.Expr {
	expr := p.parseExpr_unary()
	for {
		switch p.peekTokenType() {
		case token.MUL, token.QUO:
			tok := p.next()
			expr = &ast.BinaryExpr{
				X:  expr,
				Op: tok,
				Y:  p.parseExpr_unary(),
			}
		default:
			return expr
		}
	}
}

func (p *parser) parseExpr_unary() ast.Expr {
	if p.accept(token.ADD) {
		return p.parseExpr_primary()
	}
	if p.accept(token.SUB) {
		return &ast.UnaryExpr{
			X: p.parseExpr_primary(),
		}
	}
	return p.parseExpr_primary()
}

func (p *parser) parseExpr_primary() ast.Expr {
	peek := p.peek()

	logger.Debugf("parseExpr_primary: peek = %v\n", peek)

	switch peek.Type {
	case token.IDENT:
		ident := p.next()
		if p.accept(token.LPAREN) {
			var args []ast.Expr
			for {
				if p.accept(token.RPAREN) {
					return &ast.CallExpr{
						Fun: &ast.Ident{
							Name: ident.IdentName(),
						},
						Args: args,
					}
				}
				args = append(args, p.parseExpr())
				p.accept(token.COMMA)
			}
		}
		return &ast.Ident{
			Name: ident.IdentName(),
		}
	case token.INT:
		switch tok := p.next(); tok.Type {
		case token.INT:
			return &ast.Number{
				Value: int(tok.IntValue()),
			}
		default:
			p.err = fmt.Errorf("todo")
			panic(p.err)
		}
	case token.LPAREN:
		p.next()
		expr := p.parseExpr()
		if !p.accept(token.RPAREN) {
			p.err = fmt.Errorf("todo")
			panic(p.err)
		}
		p.next()
		return expr
	default:
		p.errorf("todo: peek=%v", peek)
		panic(p.err)
	}
}
