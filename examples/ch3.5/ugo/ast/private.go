package ast

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

func (p *File) Pos() int   { return 0 }
func (p *File) End() int   { return 0 }
func (p *File) node_type() {}

func (p *Package) Pos() int   { return 0 }
func (p *Package) End() int   { return 0 }
func (p *Package) node_type() {}

func (p *Func) Pos() int   { return 0 }
func (p *Func) End() int   { return 0 }
func (p *Func) node_type() {}

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

func (p *BlockStmt) Pos() int { return 0 }
func (p *ExprStmt) Pos() int  { return 0 }

func (p *Ident) Pos() int      { return 0 }
func (p *Number) Pos() int     { return 0 }
func (p *BinaryExpr) Pos() int { return 0 }
func (p *UnaryExpr) Pos() int  { return 0 }
func (p *ParenExpr) Pos() int  { return 0 }
func (p *CallExpr) Pos() int   { return 0 }

func (p *BlockStmt) End() int { return 0 }
func (p *ExprStmt) End() int  { return 0 }

func (p *Ident) End() int      { return 0 }
func (p *Number) End() int     { return 0 }
func (p *BinaryExpr) End() int { return 0 }
func (p *UnaryExpr) End() int  { return 0 }
func (p *ParenExpr) End() int  { return 0 }
func (p *CallExpr) End() int   { return 0 }
