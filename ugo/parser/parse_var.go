package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

// var x int
// var x int = 2

func (p *parser) parseVar() *ast.VarSpec {
	tokVar := p.mustAcceptToken(token.VAR)
	tokIdent := p.mustAcceptToken(token.IDENT)

	var varSpec = &ast.VarSpec{
		VarPos: tokVar.Pos,
	}

	varSpec.Name = &ast.Ident{
		NamePos: tokIdent.Pos,
		Name:    tokIdent.IdentName(),
	}

	switch p.peekTokenType() {
	case token.IDENT:
	case token.LBRACK: // []T
	case token.STRUCT:
	case token.MAP:
	case token.INTERFACE:
	default:
	}

	if typ, ok := p.acceptToken(token.IDENT); ok {
		varSpec.Type = &ast.Ident{
			NamePos: typ.Pos,
			Name:    typ.IdentName(),
		}
	}

	if _, ok := p.acceptToken(token.ASSIGN); ok {
		varSpec.Value = p.parseExpr()
	}

	p.acceptTokenRun(token.SEMICOLON)
	return varSpec
}
