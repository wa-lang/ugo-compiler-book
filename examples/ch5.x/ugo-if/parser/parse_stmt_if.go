package parser

import (
	"github.com/wa-lang/ugo/ast"
	"github.com/wa-lang/ugo/token"
)

func (p *Parser) parseStmt_if() *ast.IfStmt {
	tokIf := p.MustAcceptToken(token.IF)

	ifStmt := &ast.IfStmt{
		If: tokIf.Pos,
	}

	stmt := p.parseStmt()
	if _, ok := p.AcceptToken(token.SEMICOLON); ok {
		ifStmt.Init = stmt
		ifStmt.Cond = p.parseExpr()
		ifStmt.Body = p.parseStmt_block()
	} else {
		ifStmt.Init = nil
		if cond, ok := stmt.(*ast.ExprStmt); ok {
			ifStmt.Cond = cond.X
		} else {
			p.errorf(tokIf.Pos, "if cond expect expr: %#v", stmt)
		}
		ifStmt.Body = p.parseStmt_block()
	}

	if _, ok := p.AcceptToken(token.ELSE); ok {
		switch p.PeekToken().Type {
		case token.IF: // else if
			ifStmt.Else = p.parseStmt_if()
		default:
			ifStmt.Else = p.parseStmt_block()
		}
	}

	return ifStmt
}
