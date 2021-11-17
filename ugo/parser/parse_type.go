package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

// type MyInt int
// type MyInt = int
// type Point struct { ... }
// type Reader interface { ... }

func (p *parser) parseType() *ast.TypeSpec {
	tokType := p.mustAcceptToken(token.TYPE)
	tokIdent := p.mustAcceptToken(token.IDENT)

	var typeSpec = &ast.TypeSpec{
		TypePos: tokType.Pos,
	}

	typeSpec.Name = &ast.Ident{
		NamePos: tokIdent.Pos,
		Name:    tokIdent.IdentName(),
	}

	if tok, ok := p.acceptToken(token.ASSIGN); ok {
		typeSpec.Assign = tok.Pos
	}

	switch p.peekTokenType() {
	case token.IDENT:
		ident := p.nextToken()
		typeSpec.Type = &ast.Ident{
			NamePos: ident.Pos,
			Name:    ident.IdentName(),
		}

	case token.STRUCT:
		p.errorf("unsupport struct")
	case token.INTERFACE:
		p.errorf("unsupport interface")
	default:
		p.errorf("invalid token = %v", token.IDENT, p.peekToken())
	}

	p.acceptTokenRun(token.SEMICOLON)
	return typeSpec
}
