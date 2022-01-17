package parser

import (
	"github.com/wa-lang/ugo/ast"
	"github.com/wa-lang/ugo/logger"
	"github.com/wa-lang/ugo/token"
)

// const x = 1+2
// const x int = 1+2

func (p *Parser) parseStmt_const() *ast.ConstSpec {
	logger.Debugln(p.PeekToken())

	tokConst := p.MustAcceptToken(token.CONST)
	tokIdent := p.MustAcceptToken(token.IDENT)

	var constSpec = &ast.ConstSpec{
		ConstPos: tokConst.Pos,
	}

	constSpec.Name = &ast.Ident{
		NamePos: tokIdent.Pos,
		Name:    tokIdent.Literal,
	}

	if typ, ok := p.AcceptToken(token.IDENT); ok {
		constSpec.Type = &ast.Ident{
			NamePos: typ.Pos,
			Name:    typ.Literal,
		}
	}

	if _, ok := p.AcceptToken(token.ASSIGN); ok {
		constSpec.Value = p.parseExpr()
	}

	p.AcceptTokenList(token.SEMICOLON)
	return constSpec
}
