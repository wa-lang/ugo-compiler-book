package types

import "github.com/wa-lang/ugo/ast"

type PackageInfo struct {
	Pkg    *ast.Package
	Scopes map[ast.Node]*Scope

	Types map[ast.Expr]*Object
	Defs  map[*ast.Ident]*Object
	Uses  map[*ast.Ident]*Object
}

func BuildPackageInfo(pkgs map[string]*PackageInfo, pkgpath string) (*PackageInfo, error) {
	return &PackageInfo{}, nil // TODO
}
