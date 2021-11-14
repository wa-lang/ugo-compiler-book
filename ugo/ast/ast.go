package ast

import (
	"github.com/chai2010/ugo/token"
)

// File 表示 µGo 文件对应的语法树.
type File struct {
	Name     string          // 文件名
	Data     []byte          // 源文件
	Doc      *CommentGroup   // 文档注释
	List     []Stmt          // 语句列表
	Comments []*CommentGroup // 文件中的全部注释
}

func (p *File) Pos() token.Pos {
	if len(p.List) > 0 {
		return p.List[0].Pos()
	}
	return token.NoPos
}
func (p *File) End() token.Pos {
	if n := len(p.List); n > 0 {
		return p.List[n-1].End()
	}
	return token.NoPos
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
	Target *Ident    // 要赋值的目标
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
	ValuePos token.Pos // 数值的开始位置
	ValueEnd token.Pos // 数值的结束位置
	Value    int       // 数值
}

func (p *Number) Pos() token.Pos { return p.ValuePos }
func (p *Number) End() token.Pos { return p.ValueEnd }

// BinaryExpr 表示一个二元表达式.
type BinaryExpr struct {
	X     Expr        // 左边的运算对象
	OpPos token.Pos   // 运算符的位置
	Op    token.Token // 运算符
	Y     Expr        // 右边的运算对象
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
