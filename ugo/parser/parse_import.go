package parser

import (
	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/logger"
	"github.com/chai2010/ugo/token"
)

// import "path/to/pkg"
// import name "path/to/pkg"
func (p *parser) parseImport() {
	logger.Debugln("peek =", p.peekToken())

	tok, ok := p.acceptToken(token.IMPORT)
	if !ok {
		return
	}

	var importSpec = ast.ImportSpec{
		ImportPos: tok.Pos,
	}

	asName, ok := p.acceptToken(token.IDENT)
	if ok {
		importSpec.Name = &ast.Ident{
			NamePos: asName.Pos,
			Name:    asName.IdentName(),
		}
	}

	pkgPath, ok := p.acceptToken(token.STRING)
	if !ok {
		return
	}
	importSpec.Path = &ast.Ident{
		NamePos: pkgPath.Pos,
		Name:    pkgPath.StringValue(),
	}

	p.file.Imports = append(p.file.Imports, &importSpec)
}
