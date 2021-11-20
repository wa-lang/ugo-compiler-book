package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

func (p *parser) parseStmt() ast.Stmt {
	switch tok := p.r.PeekToken(); tok.Type {
	case token.EOF:
		return nil
	case token.ILLEGAL:
		panic(tok)

	case token.SEMICOLON:
		return nil
	case token.LBRACE: // {}
		return p.parseStmt_block()

	case token.TYPE: // const x = ...
		return p.parseStmt_type()
	case token.CONST: // const x = ...
		return p.parseStmt_const()
	case token.VAR: // var x = ...
		return p.parseStmt_var()

	case token.IF: // if ...
		return p.parseStmt_if()
	case token.FOR: // for ...
		return p.parseStmt_for()
	case token.DEFER: // defer ...
		return p.parseStmt_defer()
	case token.RETURN: // return ...
		return p.parseStmt_return()

	default:
		// exprList ;
		// exprList := exprList;
		// exprList = exprList;
		exprList := p.parseExprList()
		switch tok := p.r.PeekToken(); tok.Type {
		case token.SEMICOLON:
			return &ast.ExprStmt{
				X: exprList[0],
			}
		case token.DEFINE:
			valueList := p.parseExprList()
			return &ast.AssignStmt{
				Target: exprList[0],
				Value:  valueList[0],
			}
		case token.ASSIGN:
			valueList := p.parseExprList()
			return &ast.AssignStmt{
				Target: exprList[0],
				Value:  valueList[0],
			}

		default:
			panic("aa")
		}
	}
}
