package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

func (p *parser) parseStmt_var() *ast.VarSpec {
	tokVar := p.r.MustAcceptToken(token.VAR)
	tokIdent := p.r.MustAcceptToken(token.IDENT)

	var varSpec = &ast.VarSpec{
		VarPos: tokVar.Pos,
	}

	varSpec.Name = &ast.Ident{
		NamePos: tokIdent.Pos,
		Name:    tokIdent.IdentName(),
	}

	switch p.r.PeekToken().Type {
	case token.IDENT:
	case token.LBRACK: // []T
	case token.STRUCT:
	case token.MAP:
	case token.INTERFACE:
	default:
	}

	if typ, ok := p.r.AcceptToken(token.IDENT); ok {
		varSpec.Type = &ast.Ident{
			NamePos: typ.Pos,
			Name:    typ.IdentName(),
		}
	}

	if _, ok := p.r.AcceptToken(token.ASSIGN); ok {
		varSpec.Value = p.parseExpr()
	}

	p.r.AcceptTokenList(token.SEMICOLON)
	return varSpec
}
