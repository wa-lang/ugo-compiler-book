package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

// const x = 1+2
// const x int = 1+2

func (p *parser) parseConst() *ast.ConstSpec {
	tokConst := p.mustAcceptToken(token.CONST)
	tokIdent := p.mustAcceptToken(token.IDENT)

	var constSpec = &ast.ConstSpec{
		ConstPos: tokConst.Pos,
	}

	constSpec.Name = &ast.Ident{
		NamePos: tokIdent.Pos,
		Name:    tokIdent.IdentName(),
	}

	if typ, ok := p.acceptToken(token.IDENT); ok {
		constSpec.Type = &ast.Ident{
			NamePos: typ.Pos,
			Name:    typ.IdentName(),
		}
	}

	if _, ok := p.acceptToken(token.ASSIGN); ok {
		constSpec.Value = p.parseExpr()
	}

	p.acceptTokenRun(token.SEMICOLON)
	return constSpec
}
