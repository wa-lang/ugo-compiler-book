package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

// import "path/to/pkg"
// import name "path/to/pkg"
func (p *Parser) parseImport() *ast.ImportSpec {
	tokImport := p.MustAcceptToken(token.IMPORT)

	var importSpec = &ast.ImportSpec{
		ImportPos: tokImport.Pos,
	}

	asName, ok := p.AcceptToken(token.IDENT)
	if ok {
		importSpec.Name = &ast.Ident{
			NamePos: asName.Pos,
			Name:    asName.Literal,
		}
	}

	if pkgPath, ok := p.AcceptToken(token.STRING); ok {
		importSpec.Path = &ast.Ident{
			NamePos: pkgPath.Pos,
			Name:    pkgPath.Literal,
		}
	}

	return importSpec
}
