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
		case token.SEMICOLON:
			p.nextToken()
			continue
		case token.IMPORT:
			p.parseImport()
		default:
			break LoopImport
		}
	}

	for {
		switch tok := p.peekToken(); tok.Type {
		case token.EOF:
			return
		case token.SEMICOLON:
			p.nextToken()
			continue
		case token.CONST:
			p.parseConst()
		case token.TYPE:
			p.parseType()
		case token.VAR:
			p.parseVar()
		case token.FUNC:
			p.parseFunc()
		default:
			p.errorf("unknown token: %v", tok)
		}
	}
}

func (p *parser) parsePackage() {
	logger.Debugln("peek =", p.peekToken())

	pkg, ok := p.acceptToken(token.PACKAGE)
	if !ok {
		p.errorf("missing package")
	}

	ident, ok := p.acceptToken(token.IDENT)
	if !ok {
		p.errorf("missing package name")
	}

	p.file.Pkg = &ast.PackageSpec{}

	p.file.Pkg.Pkg = pkg
	p.file.Pkg.PkgName = &ast.Ident{
		Name:    ident.IdentName(),
		NamePos: ident.Pos,
	}
}
