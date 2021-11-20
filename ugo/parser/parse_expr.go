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
		if p.r.PeekToken().Type != token.COMMA {
			break
		}
	}

	p.r.AcceptTokenList(token.SEMICOLON)
	return
}

func (p *parser) parseExpr() ast.Expr {
	logger.Debugln("peek =", p.r.PeekToken())

	expr := p.parseExpr_mul()

	for {
		switch p.r.PeekToken().Type {
		case token.ADD, token.SUB:
			tok := p.r.ReadToken()
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
	logger.Debugln("peek =", p.r.PeekToken())

	expr := p.parseExpr_unary()
	for {
		switch p.r.PeekToken().Type {
		case token.SEMICOLON:
			return expr
		case token.MUL, token.QUO:
			tok := p.r.ReadToken()
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
	logger.Debugln("peek =", p.r.PeekToken())

	if _, ok := p.r.AcceptToken(token.ADD); ok {
		return p.parseExpr_primary()
	}
	if _, ok := p.r.AcceptToken(token.SUB); ok {
		return &ast.UnaryExpr{
			X: p.parseExpr_primary(),
		}
	}
	return p.parseExpr_primary()
}

func (p *parser) parseExpr_primary() ast.Expr {
	logger.Debugln("peek =", p.r.PeekToken())

	peek := p.r.PeekToken()

	switch peek.Type {
	case token.IDENT:
		ident := p.r.ReadToken()
		if lparen, ok := p.r.AcceptToken(token.LPAREN); ok {
			var args []ast.Expr
			for {
				if rparen, ok := p.r.AcceptToken(token.RPAREN); ok {
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
				p.r.AcceptToken(token.COMMA)
			}
		}
		return &ast.Ident{
			NamePos: ident.Pos,
			Name:    ident.IdentName(),
		}
	case token.INT:
		tok := p.r.ReadToken()
		return &ast.Number{
			ValuePos: tok.Pos,
			Value:    tok.IntValue(),
			ValueEnd: tok.EndPos(),
		}
	case token.FLOAT:
		tok := p.r.ReadToken()
		return &ast.Number{
			ValuePos: tok.Pos,
			Value:    tok.FloatValue(),
			ValueEnd: tok.EndPos(),
		}

	case token.LPAREN:
		p.r.ReadToken()
		expr := p.parseExpr()
		if _, ok := p.r.AcceptToken(token.RPAREN); !ok {
			p.err = fmt.Errorf("todo")
			panic(p.err)
		}
		p.r.ReadToken()
		return expr
	default:
		panic(fmt.Errorf("expr: %v", p.r.PeekToken()))
	}
}

func (p *parser) parseExpr_call() *ast.CallExpr {
	tokIdent := p.r.MustAcceptToken(token.IDENT)
	tokLparen := p.r.MustAcceptToken(token.LPAREN)
	args := p.parseExprList()
	tokRparen := p.r.MustAcceptToken(token.RPAREN)

	return &ast.CallExpr{
		Fun: &ast.Ident{
			NamePos: tokIdent.Pos,
			Name:    tokIdent.IdentName(),
		},
		Lparen: tokLparen.Pos,
		Args:   args,
		Rparen: tokRparen.Pos,
	}
}
