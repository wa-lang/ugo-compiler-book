package parser

import (
	"github.com/wa-lang/ugo/ast"
	"github.com/wa-lang/ugo/token"
)

func (p *parser) parseStmt_return() *ast.ReturnStmt {
	tokReturn := p.r.MustAcceptToken(token.RETURN)
	exprs := p.parseExprList()

	return &ast.ReturnStmt{
		Result:  tokReturn,
		Results: exprs,
	}
}
