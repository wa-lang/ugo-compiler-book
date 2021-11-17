package ast

import (
	"github.com/chai2010/ugo/token"
)

// File 表示 µGo 文件对应的语法树.
type File struct {
	Name string // 文件名
	Data []byte // 源文件

	Pkg     *PackageSpec  // 包信息
	Imports []*ImportSpec // 导入包信息
	Types   []*TypeSpec   // 类型信息
	Consts  []*ConstSpec  // 常量信息
	Globals []*VarSpec    // 全局变量
	Funcs   []*Func       // 函数列表
}

func (p *File) Pos() token.Pos {
	return token.Pos(1)
}
func (p *File) End() token.Pos {
	return token.Pos(1 + len(p.Data))
}

// Node 表示一个语法树节点.
type Node interface {
	Pos() token.Pos // 开始位置
	End() token.Pos // 结束位置

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
	Doc     *CommentGroup
	Pkg     token.Token // packaage 位置
	PkgName *Ident      // 包名
}

func (p *PackageSpec) Pos() token.Pos {
	return token.NoPos
}
func (p *PackageSpec) End() token.Pos {
	return token.NoPos
}

// TypeSpec 表示一个类型信息
type TypeSpec struct {
	TypePos token.Pos
	Assign  token.Pos // =
	Name    *Ident
	Type    Expr
}

func (p *TypeSpec) Pos() token.Pos {
	return token.NoPos
}
func (p *TypeSpec) End() token.Pos {
	return token.NoPos
}

// ImportSpec 表示一个导入包
type ImportSpec struct {
	ImportPos token.Pos
	Name      *Ident
	Path      *Ident
}

func (p *ImportSpec) Pos() token.Pos {
	return token.NoPos
}
func (p *ImportSpec) End() token.Pos {
	return token.NoPos
}

// 常量信息
type ConstSpec struct {
	ConstPos token.Pos // const 关键字位置
	Name     *Ident    // 常量名字
	Type     *Ident    // 常量类型, 可省略
	Value    Expr      // 常量表达式
}

func (p *ConstSpec) Pos() token.Pos {
	return token.NoPos
}
func (p *ConstSpec) End() token.Pos {
	return token.NoPos
}

// 变量信息
type VarSpec struct {
	VarPos token.Pos // var 关键字位置
	Name   *Ident    // 变量名字
	Type   *Ident    // 变量类型, 可省略
	Value  Expr      // 变量表达式
}

func (p *VarSpec) Pos() token.Pos {
	return token.NoPos
}
func (p *VarSpec) End() token.Pos {
	return token.NoPos
}

// 函数对象
type Func struct {
	FuncPos token.Pos  // var 关键字位置
	Self    *Field     // 方法所属类型
	Name    *Ident     // 变量名字
	Args    []*Field   // 函数参数
	Returns []*Field   // 返回值列表
	Body    *BlockStmt // 函数体
}

func (p *Func) Pos() token.Pos {
	return token.NoPos
}
func (p *Func) End() token.Pos {
	return token.NoPos
}

// 字段信息, 用于函数参数
type Field struct {
	Name *Ident
	Type Expr
}

func (p *Field) Pos() token.Pos {
	return token.NoPos
}
func (p *Field) End() token.Pos {
	return token.NoPos
}

// BlockStmt 表示一个语句块节点.
type BlockStmt struct {
	List []Stmt // 语句块中的语句列表
}

func (p *BlockStmt) Pos() token.Pos {
	if len(p.List) > 0 {
		return p.List[0].Pos()
	}
	return token.NoPos
}
func (p *BlockStmt) End() token.Pos {
	if n := len(p.List); n > 0 {
		return p.List[n-1].End()
	}
	return token.NoPos
}

// IfStmt 表示一个 if 语句节点.
type IfStmt struct {
	If   token.Pos  // if 关键字的位置
	Cond Expr       // if 条件, *BinaryExpr
	Body *BlockStmt // if 为真时对应的语句列表
	Else Stmt       // else 对应的语句
}

func (p *IfStmt) Pos() token.Pos {
	return p.If
}
func (p *IfStmt) End() token.Pos {
	if p.Else != nil {
		return p.Else.End()
	}
	return p.Body.End()
}

// ForStmt 表示一个 for 语句节点.
type ForStmt struct {
	For  token.Pos  // for 关键字的位置
	Body *BlockStmt // 循环对应的语句列表
}

func (p *ForStmt) Pos() token.Pos { return p.For }
func (p *ForStmt) End() token.Pos { return p.Body.End() }

// AssignStmt 表示一个赋值语句节点.
type AssignStmt struct {
	Target Expr      // 要赋值的目标
	TokPos token.Pos // ':=' 的位置
	Value  Expr      // 值
}

func (p *AssignStmt) Pos() token.Pos { return p.Target.Pos() }
func (p *AssignStmt) End() token.Pos { return p.Value.End() }

// Ident 表示一个标识符节点.
type Ident struct {
	NamePos token.Pos // 标识符的位置
	Name    string    // 标识符的名字
}

func (p *Ident) Pos() token.Pos { return p.NamePos }
func (p *Ident) End() token.Pos { return p.NamePos + token.Pos(len(p.Name)) }

// Number 表示一个数值.
type Number struct {
	ValuePos token.Pos   // 数值的开始位置
	ValueEnd token.Pos   // 数值的结束位置
	Value    interface{} // 数值: int/float64/str
}

func (p *Number) Pos() token.Pos { return p.ValuePos }
func (p *Number) End() token.Pos { return p.ValueEnd }

func (p *Number) IntValue() int {
	v, _ := p.Value.(int)
	return v
}

func (p *Number) FloatValue() float64 {
	v, _ := p.Value.(float64)
	return v
}

// BinaryExpr 表示一个二元表达式.
type BinaryExpr struct {
	Op token.Token // 运算符
	X  Expr        // 左边的运算对象
	Y  Expr        // 右边的运算对象
}

func (p *BinaryExpr) Pos() token.Pos { return p.X.Pos() }
func (p *BinaryExpr) End() token.Pos { return p.Y.End() }

// UnaryExpr 表示一个一元表达式.
type UnaryExpr struct {
	OpPos token.Pos   // 运算符的位置
	Op    token.Token // 运算符
	X     Expr        // 运算对象
}

func (p *UnaryExpr) Pos() token.Pos { return p.OpPos }
func (p *UnaryExpr) End() token.Pos { return p.X.End() }

// ParenExpr 表示一个圆括弧表达式.
type ParenExpr struct {
	Lparen token.Pos // "(" 的位置
	X      Expr      // 圆括弧内的表达式对象
	Rparen token.Pos // ")" 的位置
}

func (p *ParenExpr) Pos() token.Pos { return p.Lparen }
func (p *ParenExpr) End() token.Pos { return p.Rparen }

// CallExpr 表示一个函数调用
type CallExpr struct {
	Fun    *Ident    // 函数名字
	Lparen token.Pos //  "(" 的位置
	Args   []Expr    // 调用参数列表
	Rparen token.Pos // ")" 的位置
}

func (p *CallExpr) Pos() token.Pos { return p.Fun.Pos() }
func (p *CallExpr) End() token.Pos { return p.Rparen }

// Comment 表示一个注释
type Comment struct {
	Slash token.Pos // position of "/" starting the comment
	Text  string    // comment text (excluding '\n' for //-style comments)
}

func (p *Comment) Pos() token.Pos { return p.Slash }
func (p *Comment) End() token.Pos { return p.Slash + token.Pos(len(p.Text)) }

// CommentGroup 表示注释组
type CommentGroup struct {
	List []*Comment // len(List) > 0
}

func (p *CommentGroup) Pos() token.Pos {
	if n := len(p.List); n > 0 {
		return p.List[n-1].Pos()
	}
	return token.NoPos
}
func (p *CommentGroup) End() token.Pos {
	if n := len(p.List); n > 0 {
		return p.List[n-1].End()
	}
	return token.NoPos
}

func (p *CommentGroup) Text() string {
	var txt string
	for _, s := range p.List {
		txt += s.Text
	}
	return txt
}
