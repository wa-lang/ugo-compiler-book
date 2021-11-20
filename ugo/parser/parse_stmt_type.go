package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

func (p *parser) parseStmt_type() *ast.TypeSpec {
	tokType := p.r.MustAcceptToken(token.TYPE)
	tokIdent := p.r.MustAcceptToken(token.IDENT)

	var typeSpec = &ast.TypeSpec{
		TypePos: tokType.Pos,
	}

	typeSpec.Name = &ast.Ident{
		NamePos: tokIdent.Pos,
		Name:    tokIdent.IdentName(),
	}

	if tok, ok := p.r.AcceptToken(token.ASSIGN); ok {
		typeSpec.Assign = tok.Pos
	}

	switch tok := p.r.PeekToken(); tok.Type {
	case token.IDENT:
		ident := p.r.ReadToken()
		typeSpec.Type = &ast.Ident{
			NamePos: ident.Pos,
			Name:    ident.IdentName(),
		}

	case token.STRUCT:
		p.errorf(tok.Pos, "unsupport struct")
	case token.INTERFACE:
		p.errorf(tok.Pos, "unsupport interface")
	default:
		p.errorf(tok.Pos, "invalid token = %v", tok)
	}

	p.r.AcceptTokenList(token.SEMICOLON)
	return typeSpec
}
