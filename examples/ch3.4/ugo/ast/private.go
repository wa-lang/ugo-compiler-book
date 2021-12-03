package ast

var (
	_ Stmt = (*BlockStmt)(nil)
	_ Stmt = (*ExprStmt)(nil)

	_ Expr = (*Number)(nil)
	_ Expr = (*BinaryExpr)(nil)
	_ Expr = (*UnaryExpr)(nil)
	_ Expr = (*ParenExpr)(nil)
	_ Expr = (*CallExpr)(nil)
)

func (p *BlockStmt) stmt_type() {}
func (p *ExprStmt) stmt_type()  {}

func (p *Number) expr_type()     {}
func (p *BinaryExpr) expr_type() {}
func (p *UnaryExpr) expr_type()  {}
func (p *ParenExpr) expr_type()  {}
func (p *CallExpr) expr_type()   {}

func (p *BlockStmt) Pos() int { return 0 }
func (p *ExprStmt) Pos() int  { return 0 }

func (p *Number) Pos() int     { return 0 }
func (p *BinaryExpr) Pos() int { return 0 }
func (p *UnaryExpr) Pos() int  { return 0 }
func (p *ParenExpr) Pos() int  { return 0 }
func (p *CallExpr) Pos() int   { return 0 }

func (p *BlockStmt) End() int { return 0 }
func (p *ExprStmt) End() int  { return 0 }

func (p *Number) End() int     { return 0 }
func (p *BinaryExpr) End() int { return 0 }
func (p *UnaryExpr) End() int  { return 0 }
func (p *ParenExpr) End() int  { return 0 }
func (p *CallExpr) End() int   { return 0 }

/*
BlockStmt

ExprStmt

Number
BinaryExpr
UnaryExpr
ParenExpr
CallExpr


	Pos() int
	End() int
	stmt_type()
*/
