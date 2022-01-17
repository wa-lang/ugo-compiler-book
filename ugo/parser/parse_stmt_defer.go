package parser

import (
	"github.com/wa-lang/ugo/ast"
	"github.com/wa-lang/ugo/token"
)

func (p *parser) parseStmt_defer() *ast.DeferStmt {
	tokDefer := p.r.MustAcceptToken(token.DEFER)
	callExpr := p.parseExpr_call()

	return &ast.DeferStmt{
		Defer: tokDefer,
		Call:  callExpr,
	}
}
