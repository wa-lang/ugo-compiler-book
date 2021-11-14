package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/logger"
	"github.com/chai2010/ugo/token"
)

// type MyInt int
// type MyInt = int
// type Point struct { ... }
// type Reader interface { ... }

func (p *parser) parseType() {
	logger.Debugln("peek =", p.peekToken())

	tok, ok := p.acceptToken(token.TYPE)
	if !ok {
		return
	}

	var typeSpec = ast.TypeSpec{
		TypePos: tok.Pos,
	}

	name, ok := p.acceptToken(token.IDENT)
	if !ok {
		p.errorf("export %v, got = %v", token.IDENT, name)
	}
	typeSpec.Name = &ast.Ident{
		NamePos: name.Pos,
		Name:    name.IdentName(),
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

	p.acceptToken(token.SEMICOLON)

	p.file.Types = append(p.file.Types, &typeSpec)
}
