package ast

import (
	"strconv"

	"github.com/chai2010/ugo/token"
)

// File 表示 µGo 文件对应的语法树.
type File struct {
	Pkg   *PackageSpec // 包信息
	Funcs []*Func      // 函数列表
}

// Node 表示一个语法树节点.
type Node interface {
	node_private()
}

// Stmt 表示一个语句节点.
type Stmt interface {
	Node
	stmt_private()
}

// Expr 表示一个表达式节点。
type Expr interface {
	Node
	expr_private()
}

// 包信息
type PackageSpec struct {
	Pkg     token.Token
	PkgName *Ident
}

// 函数对象
type Func struct {
	Name *Ident     // 变量名字
	Body *BlockStmt // 函数体
}

// BlockStmt 表示一个语句块节点.
type BlockStmt struct {
	List []Stmt // 语句块中的语句列表
}

// ExprStmt 表示单个表达式语句
type ExprStmt struct {
	X Expr
}

// Ident 表示一个标识符节点.
type Ident struct {
	Name string // 标识符的名字
}

// Number 表示一个数值.
type Number struct {
	Value token.Token
}

func (p *Number) IntValue() int {
	v, _ := strconv.Atoi(p.Value.Literal)
	return v
}

// BinaryExpr 表示一个二元表达式.
type BinaryExpr struct {
	Op token.Token // 运算符
	X  Expr        // 左边的运算对象
	Y  Expr        // 右边的运算对象
}

// UnaryExpr 表示一个一元表达式.
type UnaryExpr struct {
	Op token.Token // 运算符
	X  Expr        // 运算对象
}

// ParenExpr 表示一个圆括弧表达式.
type ParenExpr struct {
	X Expr // 圆括弧内的表达式对象
}

// CallExpr 表示一个函数调用
type CallExpr struct {
	Fun  *Ident // 函数名字
	Args []Expr // 调用参数列表
}
