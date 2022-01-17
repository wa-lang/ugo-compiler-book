package parser

import (
	"github.com/wa-lang/ugo/ast"
	"github.com/wa-lang/ugo/token"
)

func (p *Parser) parseFile() {
	p.file = &ast.File{}

	// package xxx
	p.file.Pkg = p.parsePackage()

	for {
		switch tok := p.PeekToken(); tok.Type {
		case token.EOF:
			return
		case token.ERROR:
			panic(tok)
		case token.SEMICOLON:
			p.AcceptTokenList(token.SEMICOLON)

		case token.FUNC:
			p.file.Funcs = append(p.file.Funcs, p.parseFunc())

		default:
			p.errorf(tok.Pos, "unknown token: %v", tok)
		}
	}
}

func (p *Parser) parsePackage() *ast.Package {
	tokPkg := p.MustAcceptToken(token.PACKAGE)
	tokPkgIdent := p.MustAcceptToken(token.IDENT)

	return &ast.Package{
		PkgPos:  tokPkg.Pos,
		NamePos: tokPkgIdent.Pos,
		Name:    tokPkgIdent.Literal,
	}
}
