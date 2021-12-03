package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

func (p *Parser) parseFunc() *ast.Func {
	// func main()
	tokFunc := p.MustAcceptToken(token.FUNC)
	tokFuncIdent := p.MustAcceptToken(token.IDENT)
	p.MustAcceptToken(token.LPAREN) // (
	p.MustAcceptToken(token.RPAREN) // )

	return &ast.Func{
		FuncPos: tokFunc.Pos,
		NamePos: tokFuncIdent.Pos,
		Name:    tokFuncIdent.Literal,
		Body:    p.parseStmt_block(), // {}
	}
}
