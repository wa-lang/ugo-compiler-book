package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

func (p *parser) parseStmt_return() *ast.ReturnStmt {
	tokReturn := p.r.MustAcceptToken(token.RETURN)
	exprs := p.parseExprList()

	return &ast.ReturnStmt{
		Result:  tokReturn,
		Results: exprs,
	}
}
