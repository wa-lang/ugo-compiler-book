package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

func (p *parser) parseStmt_assign(block *ast.BlockStmt) {
	exprs := p.parseExprList()
	if len(exprs) > 1 {
		p.r.MustAcceptToken(token.ASSIGN, token.DEFINE)
		exprsRight := p.parseExprList()

		block.List = append(block.List, &ast.AssignStmt{
			Target: exprs[0],
			Value:  exprsRight[0],
		})
	} else {
		if _, ok := p.r.AcceptToken(token.ASSIGN, token.DEFINE); ok {
			exprsRight := p.parseExprList()
			block.List = append(block.List, &ast.AssignStmt{
				Target: exprs[0],
				Value:  exprsRight[0],
			})
		} else {
			block.List = append(block.List, &ast.ExprStmt{
				X: exprs[0],
			})
		}
	}
}
