package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/logger"
	"github.com/chai2010/ugo/token"
)

// var x int
// var x int = 2

func (p *parser) parseVar() {
	logger.Debugln("peek =", p.peekToken())

	tok, ok := p.acceptToken(token.VAR)
	if !ok {
		return
	}

	var varSpec = ast.VarSpec{
		VarPos: tok.Pos,
	}

	name, ok := p.acceptToken(token.IDENT)
	if ok {
		varSpec.Name = &ast.Ident{
			NamePos: name.Pos,
			Name:    name.IdentName(),
		}
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

	p.acceptToken(token.SEMICOLON)

	p.file.Globals = append(p.file.Globals, &varSpec)
}
