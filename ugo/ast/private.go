package ast

var (
	_ Node = (*File)(nil)
	_ Node = (*PackageSpec)(nil)
	_ Node = (Stmt)(nil)
	_ Node = (Expr)(nil)
)

var (
	_ Stmt = (*ImportSpec)(nil)
	_ Stmt = (*TypeSpec)(nil)
	_ Stmt = (*ConstSpec)(nil)
	_ Stmt = (*VarSpec)(nil)
	_ Stmt = (*Func)(nil)

	_ Stmt = (*BlockStmt)(nil)
	_ Stmt = (*IfStmt)(nil)
	_ Stmt = (*ForStmt)(nil)
	_ Stmt = (*AssignStmt)(nil)
)

var (
	_ Expr = (*Ident)(nil)
	_ Expr = (*Number)(nil)
	_ Expr = (*BinaryExpr)(nil)
	_ Expr = (*UnaryExpr)(nil)
	_ Expr = (*ParenExpr)(nil)
	_ Expr = (*CallExpr)(nil)
)

func (p *File) node_private()        {}
func (p *PackageSpec) node_private() {}

func (p *ImportSpec) node_private() {}
func (p *TypeSpec) node_private()   {}
func (p *ConstSpec) node_private()  {}
func (p *VarSpec) node_private()    {}
func (p *Func) node_private()       {}

func (p *BlockStmt) node_private()  {}
func (p *IfStmt) node_private()     {}
func (p *ForStmt) node_private()    {}
func (p *AssignStmt) node_private() {}

func (p *Ident) node_private()      {}
func (p *Number) node_private()     {}
func (p *BinaryExpr) node_private() {}
func (p *UnaryExpr) node_private()  {}
func (p *ParenExpr) node_private()  {}
func (p *CallExpr) node_private()   {}

func (p *Comment) node_private()      {}
func (p *CommentGroup) node_private() {}

func (p *ImportSpec) stmt_private() {}
func (p *TypeSpec) stmt_private()   {}
func (p *ConstSpec) stmt_private()  {}
func (p *VarSpec) stmt_private()    {}
func (p *Func) stmt_private()       {}

func (p *BlockStmt) stmt_private()  {}
func (p *IfStmt) stmt_private()     {}
func (p *ForStmt) stmt_private()    {}
func (p *AssignStmt) stmt_private() {}

func (p *Ident) expr_private()      {}
func (p *Number) expr_private()     {}
func (p *BinaryExpr) expr_private() {}
func (p *UnaryExpr) expr_private()  {}
func (p *ParenExpr) expr_private()  {}
func (p *CallExpr) expr_private()   {}
