package parser

import (
	"fmt"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/logger"
	"github.com/chai2010/ugo/token"
)

// x, y :=
func (p *parser) parseExprList() (exprs []ast.Expr) {
	for {
		exprs = append(exprs, p.parseExpr())
		if p.peekTokenType() != token.COMMA {
			break
		}
	}

	p.acceptTokenRun(token.SEMICOLON)
	return
}

func (p *parser) parseExpr() ast.Expr {
	logger.Debugln("peek =", p.peekToken())

	expr := p.parseExpr_mul()

	for {
		switch p.peekTokenType() {
		case token.ADD, token.SUB:
			tok := p.nextToken()
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
	logger.Debugln("peek =", p.peekToken())

	expr := p.parseExpr_unary()
	for {
		switch p.peekTokenType() {
		case token.SEMICOLON:
			return expr
		case token.MUL, token.QUO:
			tok := p.nextToken()
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
	logger.Debugln("peek =", p.peekToken())

	if _, ok := p.acceptToken(token.ADD); ok {
		return p.parseExpr_primary()
	}
	if _, ok := p.acceptToken(token.SUB); ok {
		return &ast.UnaryExpr{
			X: p.parseExpr_primary(),
		}
	}
	return p.parseExpr_primary()
}

func (p *parser) parseExpr_primary() ast.Expr {
	logger.Debugln("peek =", p.peekToken())

	peek := p.peekToken()

	switch peek.Type {
	case token.IDENT:
		ident := p.nextToken()
		if lparen, ok := p.acceptToken(token.LPAREN); ok {
			var args []ast.Expr
			for {
				if rparen, ok := p.acceptToken(token.RPAREN); ok {
					return &ast.CallExpr{
						Fun: &ast.Ident{
							NamePos: ident.Pos,
							Name:    ident.IdentName(),
						},
						Lparen: lparen.Pos,
						Args:   args,
						Rparen: rparen.Pos,
					}
				}
				args = append(args, p.parseExpr())
				p.acceptToken(token.COMMA)
			}
		}
		return &ast.Ident{
			NamePos: ident.Pos,
			Name:    ident.IdentName(),
		}
	case token.INT:
		tok := p.nextToken()
		return &ast.Number{
			ValuePos: tok.Pos,
			Value:    tok.IntValue(),
			ValueEnd: tok.EndPos(),
		}
	case token.FLOAT:
		tok := p.nextToken()
		return &ast.Number{
			ValuePos: tok.Pos,
			Value:    tok.FloatValue(),
			ValueEnd: tok.EndPos(),
		}

	case token.LPAREN:
		p.nextToken()
		expr := p.parseExpr()
		if _, ok := p.acceptToken(token.RPAREN); !ok {
			p.err = fmt.Errorf("todo")
			panic(p.err)
		}
		p.nextToken()
		return expr
	default:
		panic(fmt.Errorf("expr: %v", p.peekToken()))
	}
}
