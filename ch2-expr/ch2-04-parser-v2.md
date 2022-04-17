# 2.4. 重构解析器

在前面的例子中，我们直接用string表示了词法记号。但是这样有一些问题：比如关键字“if”和“abc”变量名的类型就不好区分，同时词法记号的位置信息也丢失了。常见做法是将词法记号定义为一个int区分的类别，同时词法记号还携带一个原始字符串表示的面值（比如一个相同值到整数可能有不同的写法），同时再辅助一些位置信息。

## 2.4.1 定义词法结构

重新定义词法结构如下：

```go
// 词法记号类型
type TokenType int

// 记号值
type Token struct {
	Type TokenType // 记号类型
	Val  string  // 记号原始字面值
	Pos  int     // 开始位置
}

// 记号类型
const (
	EOF TokenType = iota
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	LPAREN // (
	RPAREN // )
)
```

每个记号类型配有一个二元表达式的优先级(其中0表示不是二元运算符)：

```go
func (op TokenType) Precedence() int {
	switch op {
	case ADD, SUB:
		return 1
	case MUL, DIV:
		return 2
	}
	return 0
}
```

同时提供一些辅助函数：

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

## 2.4.2 定义语法树结构

词法定义从string变化为Token，AST结构调整如下：

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

现在的Token包含了完整的信息，对应每个终结字符。


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

## 2.4.4 Token流读取器

有了`Lex`函数之后我们就可以将代码转换为扁平的`Token`序列，然后在此基础之上通过语法解析器构造结构化的语法树。为了简化解析器的工作，我们再包装一个Token流读取器。

```go
type TokenReader struct {
	tokens []Token
	pos    int
	width  int
}

func NewTokenReader(input []Token) *TokenReader {
	return &TokenReader{tokens: input}
}
```

NewTokenReader 函数构造Token流读取器。

然后提供最常用的Peek、Read、Unread等类似的方法：

```go
func (p *TokenReader) PeekToken() Token {
	tok := p.ReadToken()
	p.UnreadToken()
	return tok
}

func (p *TokenReader) ReadToken() Token {
	if p.pos >= len(p.tokens) {
		p.width = 0
		return Token{Type: EOF}
	}
	tok := p.tokens[p.pos]
	p.width = 1
	p.pos += p.width
	return tok
}

func (p *TokenReader) UnreadToken() {
	p.pos -= p.width
	return
}
```

为了方便解析器工作，再定义 AcceptToken、MustAcceptToken两个方法：

```go
func (p *TokenReader) AcceptToken(expectTypes ...TokenType) (tok Token, ok bool) {
	tok = p.ReadToken()
	for _, x := range expectTypes {
		if tok.Type == x {
			return tok, true
		}
	}
	p.UnreadToken()
	return tok, false
}

func (p *TokenReader) MustAcceptToken(expectTypes ...TokenType) (tok Token) {
	tok, ok := p.AcceptToken(expectTypes...)
	if !ok {
		panic(fmt.Errorf("token.Reader.MustAcceptToken(%v) failed", expectTypes))
	}
	return tok
}
```

AcceptToken方法可以用于适配可选的Token符号，比如if之后可选的else可以用`if _, ok := r.AcceptToken(token.ELSE); ok { ... }`方式处理。而MustAcceptToken则必须匹配相应的Token，比如`t.MustAcceptToken(token.RPAREN)`强制匹配右小括弧。

## 2.4.5 二元表达式解析简化

BNF语法可以实现表达式的多优先级支持，比如前面支持加减乘除四则运算的EBNF如下：

```bnf
expr    = mul ("+" mul | "-" mul)*
mul     = primary ("*" primary | "/" primary)*
primary = num | "(" expr ")"
```

其中通过引入mul和primary规则来表示2个不同的优先级。而Go语言的二元表达式有`||`、`&&`、`==`、`+`和`*`等5中不同的优先级。如果完全通过EBNF来表示优先级则需要构造更为复杂的规则：

```bnf
expr       = logic_or
logic_or   = logic_and ("||" logic_and)*
logic_and  = equality ("&&" relational)*
equality   = relational ("==" relational | "!=" relational)*
add        = mul ("+" mul | "-" mul)*
mul        = unary ("*" unary | "/" unary)*
unary      = ("+" | "-")? primary
primary    = num | "(" expr ")"
```

这种复杂性和Go语言推崇的少即是多的哲学是相悖的！其实Go语言在设计表达式时有意无意地忽略了对右结合二元表达式的支持，如果配合运算符的优先级可以实现更简单的二元表达式解析。

下面我们看看如何简化二元表达式解析。四则运算表达式简化的ENBNF语法如下：

```bnf
expr  = unary ("+" | "-" | "*" | "/") unary)*
unary = ("+" | "-")? primary
```

不在区分优先级，只有二元和一元表达式之分。因为二元表达式只有左结合一种，配合运算符优先级可以控制剩余表达式左结合的时机。

下面重新实现ParseExpr函数：

```go
func ParseExpr(input []Token) *ExprNode {
	r := NewTokenReader(input)
	return parseExpr(r)
}

func parseExpr(r *TokenReader) *ExprNode {
	return parseExpr_binary(r, 1)
}
```

内部将Token列表转换为TokenReader，然后调用内部parseExpr函数。parseExpr函数以优先级1为参数调用parseExpr_binary解析二元表达式。

parseExpr_binary实现如下：

```go
func parseExpr_binary(r *TokenReader, prec int) *ExprNode {
	x := parseExpr_unary(r)
	for {
		op := r.PeekToken()
		if op.Type.Precedence() < prec {
			return x
		}

		r.MustAcceptToken(op.Type)
		y := parseExpr_binary(r, op.Type.Precedence()+1)
		x = &ExprNode{Token: op, Left: x, Right: y}
	}
	return nil
}
```

首先是parseExpr_unary产生一个一元的表达式，然后根据运算符op优先级和当前处理的优先级大小控制剩余表达式的递归时机。如果op比当前函数处理的优先级更高，则继续将下一个表达式递归左结合到x中，否则结束当前左结合（如果保持原有的`op.Type.Precedence()`优先级递归调用，则当前的运算符会被处理为右结合）。

一元表达式的解析如下：

```go
func parseExpr_unary(r *TokenReader) *ExprNode {
	if _, ok := r.AcceptToken(ADD); ok {
		return parseExpr_primary(r)
	}
	if _, ok := r.AcceptToken(SUB); ok {
		return &ExprNode{
			Token: Token{Type: SUB},
			Left:  &ExprNode{Token: Token{Type: NUMBER, Val: "0"}},
			Right: parseExpr_primary(r),
		}
	}
	return parseExpr_primary(r)
}
```

如果是`+x`则用`x`表示，如果是`-x`则用`0-x`表示。parseExpr_primary 则表示一个数值或小括弧包含的表达式，实现如下：

```go
func parseExpr_primary(r *TokenReader) *ExprNode {
	if _, ok := r.AcceptToken(LPAREN); ok {
		expr := parseExpr(r)
		r.MustAcceptToken(RPAREN)
		return expr
	}
	return &ExprNode{
		Token: r.MustAcceptToken(NUMBER),
	}
}
```

现在我们得到了一个更加简洁的支持多优先级只有左结合二元表达式的解析器。以后如果要支持相等更多的优先级运算符，则需要更新Token类型的优先级即可，解析器部分的代码不用变化。

