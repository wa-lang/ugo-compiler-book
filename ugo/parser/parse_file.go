package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/logger"
	"github.com/chai2010/ugo/token"
)

func (p *parser) parseFile() {
	logger.Debugln("peek =", p.r.PeekToken())

	p.file = &ast.File{
		Name: p.filename,
		Data: []byte(p.src),
	}

	// package xxx
	p.parsePackage()

LoopImport:
	for {
		switch tok := p.r.PeekToken(); tok.Type {
		case token.EOF:
			return
		case token.ILLEGAL:
			panic(tok)
		case token.SEMICOLON:
			p.r.AcceptTokenList(token.SEMICOLON)

		case token.IMPORT:
			p.file.Imports = append(p.file.Imports, p.parseImport())

		default:
			break LoopImport
		}
	}

	for {
		switch tok := p.r.PeekToken(); tok.Type {
		case token.EOF:
			return
		case token.ILLEGAL:
			panic(tok)
		case token.SEMICOLON:
			p.r.AcceptTokenList(token.SEMICOLON)

		case token.CONST:
			p.file.Consts = append(p.file.Consts, p.parseStmt_const())
		case token.TYPE:
			p.file.Types = append(p.file.Types, p.parseStmt_type())
		case token.VAR:
			p.file.Globals = append(p.file.Globals, p.parseStmt_var())
		case token.FUNC:
			p.file.Funcs = append(p.file.Funcs, p.parseFunc())

		default:
			p.errorf(tok.Pos, "unknown token: %v", tok)
		}
	}
}

func (p *parser) parsePackage() {
	pkg := p.r.MustAcceptToken(token.PACKAGE)
	ident := p.r.MustAcceptToken(token.IDENT)

	p.file.Pkg = &ast.PackageSpec{}

	p.file.Pkg.Pkg = pkg
	p.file.Pkg.PkgName = &ast.Ident{
		Name:    ident.IdentName(),
		NamePos: ident.Pos,
	}
}
