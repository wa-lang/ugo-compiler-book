package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/logger"
	"github.com/chai2010/ugo/token"
)

// const x = 1+2
// const x int = 1+2

func (p *parser) parseConst() {
	logger.Debugln("peek =", p.peekToken())

	tok, ok := p.acceptToken(token.CONST)
	if !ok {
		return
	}

	var constSpec = ast.ConstSpec{
		ConstPos: tok.Pos,
	}

	name, ok := p.acceptToken(token.IDENT)
	if !ok {
		p.errorf("export %v, got = %v", token.IDENT, name)
	}
	constSpec.Name = &ast.Ident{
		NamePos: name.Pos,
		Name:    name.IdentName(),
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

	p.acceptToken(token.SEMICOLON)

	p.file.Consts = append(p.file.Consts, &constSpec)
}
