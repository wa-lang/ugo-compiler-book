package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

func (p *parser) parseStmt_if() *ast.IfStmt {
	tokIf := p.r.MustAcceptToken(token.IF)
	cond := p.parseExpr()
	body := p.parseStmt_block()

	ifStmt := &ast.IfStmt{
		If:   tokIf.Pos,
		Cond: cond,
		Body: body,
	}

	if _, ok := p.r.AcceptToken(token.ELSE); ok {
		switch p.r.PeekToken().Type {
		case token.IF: // else if
			ifStmt.Else = p.parseStmt_if()
		default:
			ifStmt.Else = p.parseStmt_block()
		}
	}

	return ifStmt
}
