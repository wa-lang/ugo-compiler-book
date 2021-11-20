package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/token"
)

// import "path/to/pkg"
// import name "path/to/pkg"
func (p *parser) parseImport() *ast.ImportSpec {
	tokImport := p.r.MustAcceptToken(token.IMPORT)

	var importSpec = &ast.ImportSpec{
		ImportPos: tokImport.Pos,
	}

	asName, ok := p.r.AcceptToken(token.IDENT)
	if ok {
		importSpec.Name = &ast.Ident{
			NamePos: asName.Pos,
			Name:    asName.IdentName(),
		}
	}

	if pkgPath, ok := p.r.AcceptToken(token.STRING); ok {
		importSpec.Path = &ast.Ident{
			NamePos: pkgPath.Pos,
			Name:    pkgPath.StringValue(),
		}
	}

	return importSpec
}
