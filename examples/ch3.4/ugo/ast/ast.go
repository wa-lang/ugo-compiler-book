package ast

import (
	"github.com/wa-lang/ugo/token"
)

// File 表示 µGo 文件对应的语法树.
type File struct {
	Pkg   *Package // 包信息
	Funcs []*Func  // 函数列表
}

// 包信息
type Package struct {
	PkgPos  int    // package 关键字位置
	NamePos int    // 包名位置
	Name    string // 包名
}

// 函数信息
type Func struct {
	FuncPos int
	NamePos int
	Name    string
	Body    *BlockStmt
}

// 块语句
type BlockStmt struct {
	Lbrace int // '{'
	List   []Stmt
	Rbrace int // '}'
}

// Stmt 表示一个语句节点.
type Stmt interface {
	Pos() int
	End() int
	stmt_type()
}

// 表达式语句
type ExprStmt struct {
	X Expr
}

// Expr 表示一个表达式节点。
type Expr interface {
	Pos() int
	End() int
	expr_type()
}

// Number 表示一个数值.
type Number struct {
	ValuePos int
	ValueEnd int
	Value    int
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
	FuncPos  int    // Func 关键字位置
	FuncName string // 函数名字
	Lparen   int    // '(' 位置
	Args     []Expr // 调用参数列表
	Rparen   int    // ')' 位置
}
