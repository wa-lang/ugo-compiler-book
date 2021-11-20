package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

func (p *parser) parseStmt_defer() *ast.DeferStmt {
	tokDefer := p.r.MustAcceptToken(token.DEFER)
	callExpr := p.parseExpr_call()

	return &ast.DeferStmt{
		Defer: tokDefer,
		Call:  callExpr,
	}
}
