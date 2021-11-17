package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/logger"
	"github.com/chai2010/ugo/token"
)

func (p *parser) parseFile() {
	logger.Debugln("peek =", p.peekToken())

	p.file = &ast.File{
		Name: p.filename,
		Data: []byte(p.src),
	}

	// package xxx
	p.parsePackage()

LoopImport:
	for {
		switch tok := p.peekToken(); tok.Type {
		case token.EOF:
			return
		case token.ILLEGAL:
			panic(tok)
		case token.SEMICOLON:
			p.acceptTokenRun(token.SEMICOLON)

		case token.IMPORT:
			p.file.Imports = append(p.file.Imports, p.parseImport())

		default:
			break LoopImport
		}
	}

	for {
		switch tok := p.peekToken(); tok.Type {
		case token.EOF:
			return
		case token.ILLEGAL:
			panic(tok)
		case token.SEMICOLON:
			p.acceptTokenRun(token.SEMICOLON)

		case token.CONST:
			p.file.Consts = append(p.file.Consts, p.parseConst())
		case token.TYPE:
			p.file.Types = append(p.file.Types, p.parseType())
		case token.VAR:
			p.file.Globals = append(p.file.Globals, p.parseVar())
		case token.FUNC:
			p.file.Funcs = append(p.file.Funcs, p.parseFunc())

		default:
			p.errorf("unknown token: %v", tok)
		}
	}
}

func (p *parser) parsePackage() {
	pkg := p.mustAcceptToken(token.PACKAGE)
	ident := p.mustAcceptToken(token.IDENT)

	p.file.Pkg = &ast.PackageSpec{}

	p.file.Pkg.Pkg = pkg
	p.file.Pkg.PkgName = &ast.Ident{
		Name:    ident.IdentName(),
		NamePos: ident.Pos,
	}
}
