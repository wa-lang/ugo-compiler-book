package parser

import (
	"strconv"

	"github.com/wa-lang/ugo/ast"
	"github.com/wa-lang/ugo/token"
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
		path, _ := strconv.Unquote(pkgPath.Literal)
		importSpec.Path = &ast.BasicLit{
			ValuePos:  pkgPath.Pos,
			ValueType: token.STRING,
			ValueLit:  pkgPath.Literal,
			Value:     path,
		}
	}

	return importSpec
}
