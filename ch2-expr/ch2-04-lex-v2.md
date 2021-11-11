# 2.4. 重构词法解析器

在前面的例子中，我们直接用string表示了词法记号。但是这样有一些问题：比如关键字“if”和“abc”变量名的类型就不好区分，同时词法记号的位置信息也丢失了。常见做法是将词法记号定义为一个int区分的类别，同时词法记号还携带一个原始字符串表示的面值（比如一个相同值到整数可能有不同的写法），同时再辅助一些位置信息。

## 2.4.1 定义词法结构

重新定义词法结构如下：

```go
// 词法记号类型
type lexType int

// 记号值
type Token struct {
	Type lexType // 记号类型
	Val  string  // 记号原始字面值
	Pos  int     // 开始位置
}
```

我们基于int定义了一个内部的 lexType 类型，这样方便提供一些辅助函数：

```go
func (t lexType) String() string {
	switch t {
	case EOF:
		return "EOF"
	// ...
	default:
		return "UNKNWON"
	}
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%v:%v)", t.Type, t.Val)
}
```

## 2.4.2 定义词法记号值

记号类别的具体值可以通过 goyacc 自动生成（生成的是int类型常量），也可以手工定义。比如在goyacc文件定义以下记号：

```yacc
%token ILLEGAL

%token <tok> NUMBER

%token <tok> ADD // +
%token <tok> SUB // -
%token <tok> MUL // *
%token <tok> DIV // /

%token LPAREN // (
%token RPAREN // )
```

最终会生成对应无类型的常量（值避开了ASCII值部分，因此单个ASCII字符的值也可以用作记号值）：

```go
const ILLEGAL = 57346
const NUMBER = 57347
const ADD = 57348
const SUB = 57349
const MUL = 57350
const DIV = 57351
const LPAREN = 57352
const RPAREN = 57353
```

生成的记号常量和lexType类型的关联需要在产生Token对象时手工转换。

需要注意到是程序中的注释是作为一种Token，但是在进行语法解析时或提前过滤掉注释，然后在最终的AST中再重新为每个节点组装对应的注释（如果在BNF添加注释，将会带来大量的噪声，得不偿失）。

## 2.4.3 词法解析器重构

有了Token结构之后，我们只需要将输入到代码字符串解析为Token序列即可。之前的词法是通过`strings.IndexAny`做分词，现在我们改用`text/scanner`包来辅助解析：

```go
func Lex(code string) (tokens []Token) {
	var s scanner.Scanner
	s.Init(strings.NewReader(code))
	for x := s.Scan(); x != scanner.EOF; x = s.Scan() {
		var tok = Token{
			Val: s.TokenText(),
			Pos: s.Pos().Offset,
		}
		switch x {
		case scanner.Int:
			tok.Type = NUMBER
		default:
			switch s.TokenText() {
			case "+":
				tok.Type = ADD
			case "-":
				tok.Type = SUB
			case "*":
				tok.Type = MUL
			case "/":
				tok.Type = DIV
			case "(":
				tok.Type = LPAREN
			case ")":
				tok.Type = RPAREN
			default:
				tok.Type = ILLEGAL
				tokens = append(tokens, tok)
				return
			}
		}

		tokens = append(tokens, tok)
	}

	tokens = append(tokens, Token{Type: EOF})
	return
}
```

`scanner.Scanner`可以解析常见的词法，我们只提取表达式需要的整数、四则运算和小括弧，其他类型记号用`ILLEGAL`表示无效值。

`scanner.Scanner`对于解析真正的Go代码是稍显不足的。不过词法解析是相对容易实现的工作，大家可以选择自己喜欢的方式实现，甚至直接参考Go的`go/token`和`go/scanner`包也可以。

## 2.4.4 同步调整其他部分

词法定义从string变化为Token，AST、goyacc语法文件和编译器都需要做相应的调整（程序结构可以不用变化）。

首先是语法树的值变成了Token：

```go
type ExprNode struct {
	Token           // +, -, *, /, 123
	Left  *ExprNode `json:",omitempty"`
	Right *ExprNode `json:",omitempty"`
}

func NewExprNode(token Token, left, right *ExprNode) *ExprNode {
	return &ExprNode{
		Token: token,
		Left:  left,
		Right: right,
	}
}
```

现在的Token包含了完整的信息，对应每个终结字符。我们在表示节点状态的结构体中增加tok字段：

```yacc
// expr.y
%union {
	node *ExprNode
	tok  Token
}
```

而用于和goyacc适配的词法解析函数针对每个终结字符填充tok字段：

```go
const EOF = 0 // EOF 必须为 0 (类似C语言字符串的'\0'结尾), 表示结束

func (p *exprLex) Lex(yylval *yySymType) int {
	if p.pos >= len(p.tokens) {
		return EOF
	}
	
	yylval.tok = p.tokens[p.pos]
	p.pos++

	return int(yylval.tok.Type)
}
```

词法分析器产生的都是终结记号、或者类似解析树的叶子结点，非叶子结点则对应BNF中定义的规则。终结记号的信息通过`exprLex.Lex`方法填充到`yylval.tok`（就是`%union`新添加到字段）。而BNF定义的规则对应中间的节点依然是node字段，通过以下代码生成：

```yacc
top: expr { yyrcvr.lval.node = $1 }

expr: mul { $$ = $1 }
	| expr ADD mul { $$ = NewExprNode($2, $1, $3) }
	| expr SUB mul { $$ = NewExprNode($2, $1, $3) }

mul: primary { $$ = $1 }
	| mul MUL primary { $$ = NewExprNode($2, $1, $3) }
	| mul DIV primary { $$ = NewExprNode($2, $1, $3) }

primary: NUMBER { $$ = NewExprNode($1, nil, nil) }
	| '(' expr ')' { $$ = $2 }
```

比如 primary 对应 NUMBER 时，`$$ = NewExprNode($1, nil, nil)` 代码中 `$$` 表示 primary 对应的节点内字段（就是 node），而 `$1` 则对应 NUMBER 终结记号对应的字段（对应新加的 tok 字段）。然后通过 NewExprNode 递归从底往上构造语法树。在最顶层的 top 规则，直接通过 `yyrcvr.lval.node = $1` 将构造的语法树记录到 `parser.lval.node` 中。

最后将 compiler 相关部分也响应的调整即可，这样我们就实现了相对完整的四则运算的编译器。
