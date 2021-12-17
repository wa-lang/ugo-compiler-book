package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

func (p *Parser) parseStmt_var() *ast.VarSpec {
	tokVar := p.MustAcceptToken(token.VAR)
	tokIdent := p.MustAcceptToken(token.IDENT)

	var varSpec = &ast.VarSpec{
		VarPos: tokVar.Pos,
	}

	varSpec.Name = &ast.Ident{
		NamePos: tokIdent.Pos,
		Name:    tokIdent.Literal,
	}

	if typ, ok := p.AcceptToken(token.IDENT); ok {
		varSpec.Type = &ast.Ident{
			NamePos: typ.Pos,
			Name:    typ.Literal,
		}
	}

	if _, ok := p.AcceptToken(token.ASSIGN); ok {
		varSpec.Value = p.parseExpr()
	}

	p.AcceptTokenList(token.SEMICOLON)
	return varSpec
}
