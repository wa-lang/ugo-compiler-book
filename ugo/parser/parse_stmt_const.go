package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

// const x = 1+2
// const x int = 1+2

func (p *parser) parseStmt_const() *ast.ConstSpec {
	tokConst := p.r.MustAcceptToken(token.CONST)
	tokIdent := p.r.MustAcceptToken(token.IDENT)

	var constSpec = &ast.ConstSpec{
		ConstPos: tokConst.Pos,
	}

	constSpec.Name = &ast.Ident{
		NamePos: tokIdent.Pos,
		Name:    tokIdent.IdentName(),
	}

	if typ, ok := p.r.AcceptToken(token.IDENT); ok {
		constSpec.Type = &ast.Ident{
			NamePos: typ.Pos,
			Name:    typ.IdentName(),
		}
	}

	if _, ok := p.r.AcceptToken(token.ASSIGN); ok {
		constSpec.Value = p.parseConstExpr()
	}

	p.r.AcceptTokenList(token.SEMICOLON)
	return constSpec
}
