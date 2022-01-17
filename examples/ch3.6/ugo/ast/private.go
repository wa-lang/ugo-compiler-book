package ast

import "github.com/wa-lang/ugo/token"

var (
	_ Node = Expr(nil)
	_ Node = Stmt(nil)

	_ Node = (*File)(nil)

	_ Node = (*Package)(nil)
	_ Node = (*Func)(nil)

	_ Stmt = (*BlockStmt)(nil)
	_ Stmt = (*ExprStmt)(nil)

	_ Expr = (*Ident)(nil)
	_ Expr = (*Number)(nil)
	_ Expr = (*BinaryExpr)(nil)
	_ Expr = (*UnaryExpr)(nil)
	_ Expr = (*ParenExpr)(nil)
	_ Expr = (*CallExpr)(nil)
)

func (p *File) Pos() token.Pos { return token.NoPos }
func (p *File) End() token.Pos { return token.NoPos }
func (p *File) node_type()     {}

func (p *Package) Pos() token.Pos { return token.NoPos }
func (p *Package) End() token.Pos { return token.NoPos }
func (p *Package) node_type()     {}

func (p *Func) Pos() token.Pos { return token.NoPos }
func (p *Func) End() token.Pos { return token.NoPos }
func (p *Func) node_type()     {}

func (p *BlockStmt) node_type() {}
func (p *ExprStmt) node_type()  {}

func (p *Ident) node_type()      {}
func (p *Number) node_type()     {}
func (p *BinaryExpr) node_type() {}
func (p *UnaryExpr) node_type()  {}
func (p *ParenExpr) node_type()  {}
func (p *CallExpr) node_type()   {}

func (p *BlockStmt) stmt_type() {}
func (p *ExprStmt) stmt_type()  {}

func (p *Ident) expr_type()      {}
func (p *Number) expr_type()     {}
func (p *BinaryExpr) expr_type() {}
func (p *UnaryExpr) expr_type()  {}
func (p *ParenExpr) expr_type()  {}
func (p *CallExpr) expr_type()   {}

func (p *BlockStmt) Pos() token.Pos { return token.NoPos }
func (p *ExprStmt) Pos() token.Pos  { return token.NoPos }

func (p *Ident) Pos() token.Pos      { return token.NoPos }
func (p *Number) Pos() token.Pos     { return token.NoPos }
func (p *BinaryExpr) Pos() token.Pos { return token.NoPos }
func (p *UnaryExpr) Pos() token.Pos  { return token.NoPos }
func (p *ParenExpr) Pos() token.Pos  { return token.NoPos }
func (p *CallExpr) Pos() token.Pos   { return token.NoPos }

func (p *BlockStmt) End() token.Pos { return token.NoPos }
func (p *ExprStmt) End() token.Pos  { return token.NoPos }

func (p *Ident) End() token.Pos      { return token.NoPos }
func (p *Number) End() token.Pos     { return token.NoPos }
func (p *BinaryExpr) End() token.Pos { return token.NoPos }
func (p *UnaryExpr) End() token.Pos  { return token.NoPos }
func (p *ParenExpr) End() token.Pos  { return token.NoPos }
func (p *CallExpr) End() token.Pos   { return token.NoPos }
