package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

func (p *Parser) parseFile() {
	p.file = &ast.File{
		Filename: p.Filename(),
		Source:   p.Source(),
	}

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

		case token.VAR:
			p.file.Globals = append(p.file.Globals, p.parseStmt_var())
		case token.FUNC:
			p.file.Funcs = append(p.file.Funcs, p.parseFunc())

		default:
			p.errorf(tok.Pos, "unknown token: %v", tok)
		}
	}
}

func (p *Parser) parsePackage() *ast.PackageSpec {
	tokPkg := p.MustAcceptToken(token.PACKAGE)
	tokPkgIdent := p.MustAcceptToken(token.IDENT)

	return &ast.PackageSpec{
		PkgPos:  tokPkg.Pos,
		NamePos: tokPkgIdent.Pos,
		Name:    tokPkgIdent.Literal,
	}
}
